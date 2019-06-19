package core

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/muidea/magicProxy/common/sql-parser/mysql"
)

var session = uint32(0)

// Server struct
type Server struct {
	join                *sync.WaitGroup
	listenPort, authURL string

	pool     chan Backend // mysql connections
	listener net.Listener
	running  bool
}

func (p *Server) getMySQL() (m Backend) {
	m = <-p.pool
	return
}

func (p *Server) retMySQL(m Backend) {
	if !p.running {
		m.Close()
		return
	}
	select {
	case p.pool <- m:
	default:
		m.Close()
	}
}

// KeepAlive keep alive
func (p *Server) KeepAlive() {
	var err error
	for p.running {
		time.Sleep(1 * time.Second)
		select {
		case cli, o := <-p.pool:
			if !o {
				if cli != nil {
					cli.Close()
					p.join.Done()
				}
				return
			}
			err = cli.Ping()
			p.pool <- cli
			if err != nil {
				log.Printf("Server, KeepAlive, %s", err.Error())
			}
		}
	}
}

// Close func
func (p *Server) Close() {
	p.running = false
	if p.listener != nil {
		p.listener.Close()
	}

	for {
		final := false
		select {
		case b, o := <-p.pool:
			if b == nil && o == false {
				log.Printf("backend conn close")
				final = true
				break
			}
			b.Close()
			p.join.Done()
		}

		if final {
			break
		}
	}

	close(p.pool)
}

// Auth func
func (p *Server) Auth(appid string, secret []byte) (bool, error) {
	return true, nil
}

// NewServer func
func NewServer(cfg *Config, join *sync.WaitGroup) (p *Server) {
	p = new(Server)
	p.join = join
	p.listenPort = cfg.ListenPort
	p.pool = make(chan Backend, cfg.MySQLPoolSize)
	for i := cfg.MySQLPoolSize; i > 0; i-- {
		co := new(backConn)
		if err := co.Connect(cfg.MySQLURI, cfg.MySQLUser, cfg.MySQLPswd); err != nil {
			log.Printf("NewServer, %s", err.Error())
			co.Close()
			continue
		}
		p.pool <- co
		join.Add(1)
	}

	var err error
	if p.listener, err = net.Listen("tcp", cfg.ListenPort); err != nil {
		panic(fmt.Errorf("listen %s failed %s", cfg.ListenPort, err))
	}
	return
}

// Run func
func (p *Server) Run() error {
	p.join.Add(1)
	defer p.join.Done()
	p.running = true
	go p.KeepAlive()
	for p.running {
		conn, err := p.listener.Accept()
		if err != nil {
			log.Printf("Server, Run, %s", err.Error())
			continue
		}

		go p.onConn(conn)
	}

	return nil
}

func (p *Server) warpConn(c net.Conn) (conn *Conn) {

	tcpConn := c.(*net.TCPConn)

	//SetNoDelay controls whether the operating system should delay packet transmission
	// in hopes of sending fewer packets (Nagle's algorithm).
	// The default is true (no delay),
	// meaning that data is sent as soon as possible after a Write.
	//I set this option false.
	tcpConn.SetNoDelay(false)

	conn = new(Conn)
	conn.hash = sha1.New()
	conn.connectionID = atomic.AddUint32(&session, 1)
	conn.status = mysql.SERVER_STATUS_AUTOCOMMIT
	conn.stmtID = 0
	conn.affectedRows = 0
	conn.lastInsertID = 0
	conn.raw = c
	conn.stmts = make(map[uint32]*Stmt)
	conn.salt, _ = mysql.RandomBuf(20)
	conn.PacketIO = mysql.NewPacketIO(conn.raw)
	conn.PacketIO.Sequence = 0
	return
}

func (p *Server) onConn(c net.Conn) {
	conn := p.warpConn(c) //新建一个conn
	defer func() {
		err := recover()
		if err != nil {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)] //获得当前goroutine的stacktrace
			log.Printf("Server, onConn %s %s %s", c.RemoteAddr().String(), err, string(buf))
		}
		conn.Close()
		return
	}()
	p.join.Add(1)
	defer p.join.Done()
	m := p.getMySQL()
	defer p.retMySQL(m)
	if err := conn.Handshake(); err != nil {
		log.Printf("Server, onConn, %s", err.Error())
		conn.writeErr(err)
		conn.Close()
		return
	}

	conn.Run(m)
}
