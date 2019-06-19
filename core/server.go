package core

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/muidea/magicCommon/foundation/cache"
	"github.com/muidea/magicProxy/common/sql-parser/mysql"
)

var session = uint32(0)

// Server struct
type Server struct {
	join       *sync.WaitGroup
	listenPort string
	config     *Config
	listener   net.Listener
	running    bool

	backendRegistry *Registry
	connCache       cache.KVCache
}

func (p *Server) getBackend(cfg *Config) Backend {
	backend, err := p.backendRegistry.FetchOut()
	if err != nil {
		return nil
	}

	return backend
}

// Close func
func (p *Server) Close() {
	p.running = false
	if p.listener != nil {
		p.listener.Close()
	}

	p.backendRegistry.Close()
	p.backendRegistry = nil

	keys := p.connCache.GetAll()
	for _, v := range keys {
		val, ok := p.connCache.Fetch(v)
		if ok {
			val.(*Conn).Close()
		}
	}

	p.connCache.Release()
	p.connCache = nil
}

// NewServer func
func NewServer(cfg *Config, join *sync.WaitGroup) (p *Server) {
	p = new(Server)
	p.join = join
	p.listenPort = cfg.ListenPort
	p.config = cfg
	p.connCache = cache.NewKVCache()
	p.backendRegistry = NewRegistry(cfg, join)

	var err error
	p.listener, err = net.Listen("tcp", cfg.ListenPort)
	if err != nil {
		panic(fmt.Errorf("listen %s failed %s", cfg.ListenPort, err))
	}

	return
}

// Run func
func (p *Server) Run() error {
	p.join.Add(1)
	defer p.join.Done()
	p.running = true
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

	conn.connCache = p.connCache
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

	backend := p.getBackend(p.config)
	if backend == nil {
		log.Printf("create new connect failed.")
		return
	}

	p.join.Add(1)
	defer p.join.Done()

	if err := conn.Handshake(); err != nil {
		log.Printf("Server, onConn failed, %s", err.Error())
		conn.writeErr(err)
		conn.Close()
		backend.Close()
		backend = nil

		return
	}

	conn.Run(backend)
}
