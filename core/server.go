package core

import (
	"fmt"
	"log"
	"net"
)

// Server proxy server
type Server struct {
	config      *Config
	listener    net.Listener
	runningFlag bool
}

// NewServer new proxy server
func NewServer(cfg *Config) (ret *Server) {
	listener, listenErr := net.Listen("tcp", cfg.BindAddr)
	if listenErr != nil {
		panic(fmt.Errorf("listen %s failed %s", cfg.BindAddr, listenErr))
	}

	ret = &Server{config: cfg, listener: listener}
	return
}

// Run run server
func (s *Server) Run() {
	s.runningFlag = true
	for s.runningFlag {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Proxy, Run, %s", err.Error())
			continue
		}

		go s.onConn(conn)
	}
}

func (s *Server) warpConn(conn net.Conn) (ret *Conn) {

	tcpConn := conn.(*net.TCPConn)

	//SetNoDelay controls whether the operating system should delay packet transmission
	// in hopes of sending fewer packets (Nagle's algorithm).
	// The default is true (no delay),
	// meaning that data is sent as soon as possible after a Write.
	//I set this option false.
	tcpConn.SetNoDelay(false)

	ret = NewConnect(conn, s.config)

	return
}

func (s *Server) onConn(conn net.Conn) {
	log.Printf("new connect,from:%s", conn.RemoteAddr().String())

	wrapConn := s.warpConn(conn)

	wrapConn.Run()
}
