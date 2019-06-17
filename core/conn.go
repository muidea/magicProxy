package core

import (
	"log"
	"net"

	"github.com/muidea/magicProxy/common/sql-parser/mysql"
)

// Conn proxy connection
type Conn struct {
	packetIO *mysql.PacketIO

	backend     *Backend
	config      *Config
	rawConn     net.Conn
	runningFlag bool
}

// NewConnect new proxy connection
func NewConnect(conn net.Conn, cfg *Config) (ret *Conn) {
	backend := NewBackend(cfg.BackendAddr)
	err := backend.Connect()
	if err != nil {
		return
	}

	packetIO := mysql.NewPacketIO(conn)

	ret = &Conn{rawConn: conn, config: cfg, packetIO: packetIO, backend: backend}

	return
}

func (s *Conn) initShakeHands() {
	for {
		data, dataErr := s.backend.ReadData()
		if dataErr != nil {
			log.Printf("recv data from backend failed, err:%s", dataErr.Error())
			break
		}
		log.Printf("recv data from backend, data size:%d", len(data))

		sendErr := s.packetIO.WritePacket(data)
		if sendErr != nil {
			log.Printf("send data to client failed, err:%s", sendErr.Error())
			break
		}
		break
	}
}

// Run run connection
func (s *Conn) Run() {
	s.runningFlag = true

	s.initShakeHands()

	for s.runningFlag {
		data, dataErr := s.packetIO.ReadPacket()
		if dataErr != nil {
			log.Printf("recv data from client failed, err:%s", dataErr.Error())
			break
		}
		log.Printf("recv data from client, data size:%d", len(data))

		sendErr := s.backend.SendData(data)
		if sendErr != nil {
			log.Printf("send data to backend failed, err:%s", sendErr.Error())
			s.backend.Connect()
			continue
		}

		data, dataErr = s.backend.ReadData()
		if dataErr != nil {
			log.Printf("recv data from backend failed, err:%s", dataErr.Error())
			s.backend.Connect()
			continue
		}
		log.Printf("recv data from backend, data size:%d", len(data))

		sendErr = s.packetIO.WritePacket(data)
		if sendErr != nil {
			log.Printf("send data to client failed, err:%s", sendErr.Error())
			break
		}
	}
}
