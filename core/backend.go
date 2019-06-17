package core

import (
	"log"
	"net"

	"github.com/muidea/magicProxy/common/sql-parser/mysql"
)

//Backend proxy backend
type Backend struct {
	packetIO *mysql.PacketIO

	conn        net.Conn
	backendAddr string
}

// NewBackend new backend
func NewBackend(backendAddr string) (ret *Backend) {
	ret = &Backend{backendAddr: backendAddr}
	err := ret.Connect()
	if err != nil {
		log.Printf("connect backend failed")
		ret = nil
	}

	return
}

// Connect connect backend
func (s *Backend) Connect() (err error) {
	netConn, netErr := net.Dial("tcp", s.backendAddr)
	if netErr != nil {
		err = netErr
		log.Printf("connect backend failed, err:%s", err.Error())
		return
	}

	tcpConn := netConn.(*net.TCPConn)

	//SetNoDelay controls whether the operating system should delay packet transmission
	// in hopes of sending fewer packets (Nagle's algorithm).
	// The default is true (no delay),
	// meaning that data is sent as soon as possible after a Write.
	//I set this option false.
	tcpConn.SetNoDelay(false)
	tcpConn.SetKeepAlive(true)
	s.conn = tcpConn

	s.packetIO = mysql.NewPacketIO(s.conn)
	return
}

// SendData send data
func (s *Backend) SendData(data []byte) (err error) {
	err = s.packetIO.WritePacket(data)
	return
}

// ReadData read data
func (s *Backend) ReadData() (ret []byte, err error) {
	ret, err = s.packetIO.ReadPacket()
	return
}
