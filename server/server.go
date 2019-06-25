package server

import (
	"log"
	"net"
	"runtime"
	"sync/atomic"

	"github.com/muidea/magicProxy/backend"
	"github.com/muidea/magicProxy/mysql"

	"github.com/muidea/magicProxy/config"
	"github.com/muidea/magicProxy/core/errors"
	"github.com/muidea/magicProxy/core/golog"
)

func parseNode(cfg config.NodeConfig) (*backend.Node, error) {
	n := new(backend.Node)
	n.Cfg = cfg

	n.Online = true
	go n.CheckNode()

	return n, nil
}

// Server server define
type Server struct {
	cfg  *config.Config
	addr string

	databaseNode *backend.Node

	listener net.Listener
	running  bool
}

// NewServer new server
func NewServer(cfg *config.Config) (*Server, error) {
	s := new(Server)

	s.cfg = cfg
	s.addr = cfg.Addr

	if len(cfg.Charset) == 0 {
		cfg.Charset = mysql.DEFAULT_CHARSET //utf8
	}
	cid, ok := mysql.CharsetIds[cfg.Charset]
	if !ok {
		return nil, errors.ErrInvalidCharset
	}
	//change the default charset
	mysql.DEFAULT_CHARSET = cfg.Charset
	mysql.DEFAULT_COLLATION_ID = cid
	mysql.DEFAULT_COLLATION_NAME = mysql.Collations[cid]

	var err error
	netProto := "tcp"

	s.listener, err = net.Listen(netProto, s.addr)

	if err != nil {
		return nil, err
	}

	golog.Info("server", "NewServer", "Server running", 0,
		"netProto",
		netProto,
		"address",
		s.addr)
	return s, nil
}

func (s *Server) newClientConn(co net.Conn) *ClientConn {
	c := new(ClientConn)
	tcpConn := co.(*net.TCPConn)

	//SetNoDelay controls whether the operating system should delay packet transmission
	// in hopes of sending fewer packets (Nagle's algorithm).
	// The default is true (no delay),
	// meaning that data is sent as soon as possible after a Write.
	//I set this option false.
	tcpConn.SetNoDelay(false)
	c.rawConn = tcpConn

	c.pkg = mysql.NewPacketIO(tcpConn)
	c.proxy = s

	c.pkg.Sequence = 0

	c.connectionID = atomic.AddUint32(&baseConnID, 1)

	c.status = mysql.SERVER_STATUS_AUTOCOMMIT

	c.salt, _ = mysql.RandomBuf(20)

	c.closed = false

	c.charset = mysql.DEFAULT_CHARSET
	c.collation = mysql.DEFAULT_COLLATION_ID

	return c
}

func (s *Server) onConn(c net.Conn) {
	conn := s.newClientConn(c) //新建一个conn

	defer func() {
		err := recover()
		if err != nil {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)] //获得当前goroutine的stacktrace
			golog.Error("server", "onConn", "error", 0,
				"remoteAddr", c.RemoteAddr().String(),
				"stack", string(buf),
			)
		}

		conn.Close()
	}()

	if err := conn.Handshake(); err != nil {
		golog.Error("server", "onConn", err.Error(), 0)
		conn.writeError(err)
		conn.Close()
		return
	}

	golog.Info("server", "onConn", "new connection", 0, "remote address", c.RemoteAddr().String(), "database", conn.currentDB)

	conn.Run()
}

// Run run server
func (s *Server) Run() error {
	s.running = true

	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			golog.Error("server", "Run", err.Error(), 0)
			continue
		}

		go s.onConn(conn)
	}

	return nil
}

// Close close server
func (s *Server) Close() {
	s.running = false
	if s.listener != nil {
		s.listener.Close()
	}
}

// GetBackendNode get backend node
func (s *Server) GetBackendNode() (ret *backend.Node) {
	if s.databaseNode == nil {
		node, err := parseNode(s.cfg.Node)
		if err != nil {
			log.Printf("parse database Node failed,")
			return
		}

		s.databaseNode = node
	}

	ret = s.databaseNode
	return
}
