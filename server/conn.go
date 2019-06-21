package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"

	"github.com/muidea/magicProxy/core/golog"
	"github.com/muidea/magicProxy/core/hack"
	"github.com/muidea/magicProxy/mysql"
)

// ClientConn client <-> proxy
type ClientConn struct {
	sync.Mutex

	pkg *mysql.PacketIO

	rawConn net.Conn

	proxy *Server

	capability uint32

	connectionID uint32

	status    uint16
	collation mysql.CollationId
	charset   string

	user string
	db   string

	salt []byte

	closed bool

	affectedRows int64

	stmtID uint32

	stmts map[uint32]*Stmt //prepare相关,client端到proxy的stmt
}

// DefaultCapability default capability
var DefaultCapability = uint32(mysql.CLIENT_LONG_PASSWORD | mysql.CLIENT_LONG_FLAG |
	mysql.CLIENT_CONNECT_WITH_DB | mysql.CLIENT_PROTOCOL_41 |
	mysql.CLIENT_TRANSACTIONS | mysql.CLIENT_SECURE_CONNECTION)

var baseConnID uint32 = 10000

// Handshake handshake
func (c *ClientConn) Handshake() error {
	if err := c.writeInitialHandshake(); err != nil {
		golog.Error("server", "Handshake", err.Error(),
			c.connectionID, "msg", "send initial handshake error")
		return err
	}

	if err := c.readHandshakeResponse(); err != nil {
		golog.Error("server", "readHandshakeResponse",
			err.Error(), c.connectionID,
			"msg", "read Handshake Response error")
		return err
	}

	if err := c.writeOK(nil); err != nil {
		golog.Error("server", "readHandshakeResponse",
			"write ok fail",
			c.connectionID, "error", err.Error())
		return err
	}

	c.pkg.Sequence = 0
	return nil
}

// Close close
func (c *ClientConn) Close() error {
	if c.closed {
		return nil
	}

	c.rawConn.Close()

	c.closed = true

	return nil
}

func (c *ClientConn) writeInitialHandshake() error {
	data := make([]byte, 4, 128)

	//min version 10
	data = append(data, 10)

	//server version[00]
	data = append(data, mysql.ServerVersion...)
	data = append(data, 0)

	//connection id
	data = append(data, byte(c.connectionID), byte(c.connectionID>>8), byte(c.connectionID>>16), byte(c.connectionID>>24))

	//auth-plugin-data-part-1
	data = append(data, c.salt[0:8]...)

	//filter [00]
	data = append(data, 0)

	//capability flag lower 2 bytes, using default capability here
	data = append(data, byte(DefaultCapability), byte(DefaultCapability>>8))

	//charset, utf-8 default
	data = append(data, uint8(mysql.DEFAULT_COLLATION_ID))

	//status
	data = append(data, byte(c.status), byte(c.status>>8))

	//below 13 byte may not be used
	//capability flag upper 2 bytes, using default capability here
	data = append(data, byte(DefaultCapability>>16), byte(DefaultCapability>>24))

	//filter [0x15], for wireshark dump, value is 0x15
	data = append(data, 0x15)

	//reserved 10 [00]
	data = append(data, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

	//auth-plugin-data-part-2
	data = append(data, c.salt[8:]...)

	//filter [00]
	data = append(data, 0)

	return c.writePacket(data)
}

func (c *ClientConn) readPacket() ([]byte, error) {
	return c.pkg.ReadPacket()
}

func (c *ClientConn) writePacket(data []byte) error {
	return c.pkg.WritePacket(data)
}

func (c *ClientConn) writePacketBatch(total, data []byte, direct bool) ([]byte, error) {
	return c.pkg.WritePacketBatch(total, data, direct)
}

func (c *ClientConn) readHandshakeResponse() error {
	data, err := c.readPacket()

	if err != nil {
		return err
	}

	pos := 0

	//capability
	c.capability = binary.LittleEndian.Uint32(data[:4])
	pos += 4

	//skip max packet size
	pos += 4

	//charset, skip, if you want to use another charset, use set names
	//c.collation = CollationId(data[pos])
	pos++

	//skip reserved 23[00]
	pos += 23

	//user name
	c.user = string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])

	pos += len(c.user) + 1

	//auth length and auth
	authLen := int(data[pos])
	pos++
	//auth := data[pos : pos+authLen]

	//check user
	//if _, ok := c.proxy.users[c.user]; !ok {
	//	golog.Error("ClientConn", "readHandshakeResponse", "error", 0,
	//		"auth", auth,
	//		"client_user", c.user,
	//		"config_set_user", c.user,
	//		"passworld", c.proxy.users[c.user])
	//	return mysql.NewDefaultError(mysql.ER_ACCESS_DENIED_ERROR, c.user, c.c.RemoteAddr().String(), "Yes")
	//}

	//check password
	//checkAuth := mysql.CalcPassword(c.salt, []byte(c.proxy.users[c.user]))
	//if !bytes.Equal(auth, checkAuth) {
	//	golog.Error("ClientConn", "readHandshakeResponse", "error", 0,
	//		"auth", auth,
	//		"checkAuth", checkAuth,
	//		"client_user", c.user,
	//		"config_set_user", c.user,
	//		"passworld", c.proxy.users[c.user])
	//	return mysql.NewDefaultError(mysql.ER_ACCESS_DENIED_ERROR, c.user, c.c.RemoteAddr().String(), "Yes")
	//}

	pos += authLen

	var db string
	if c.capability&mysql.CLIENT_CONNECT_WITH_DB > 0 {
		if len(data[pos:]) == 0 {
			return nil
		}

		db = string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])
		pos += len(c.db) + 1

	}
	c.db = db

	return nil
}

// Run run
func (c *ClientConn) Run() {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]

			golog.Error("ClientConn", "Run",
				err.Error(), 0,
				"stack", string(buf))
		}

		c.Close()
	}()

	for {
		data, err := c.readPacket()
		if err != nil {
			return
		}

		err = c.dispatch(data)
		if err != nil {
			golog.Error("ClientConn", "Run",
				err.Error(), c.connectionID,
			)
			c.writeError(err)
			if err == mysql.ErrBadConn {
				c.Close()
			}
		}

		if c.closed {
			return
		}

		c.pkg.Sequence = 0
	}
}

func (c *ClientConn) dispatch(data []byte) error {
	cmd := data[0]
	data = data[1:]

	log.Printf("sql:%s", hack.String(data))

	switch cmd {
	case mysql.COM_QUIT:
		c.handleRollback()
		c.Close()
		return nil
	case mysql.COM_QUERY:
		return c.handleQuery(hack.String(data))
	case mysql.COM_PING:
		return c.writeOK(nil)
	case mysql.COM_INIT_DB:
		return c.handleUseDB(hack.String(data))
	case mysql.COM_FIELD_LIST:
		return c.handleFieldList(data)
	case mysql.COM_STMT_PREPARE:
		return c.handleStmtPrepare(hack.String(data))
	case mysql.COM_STMT_EXECUTE:
		return c.handleStmtExecute(data)
	case mysql.COM_STMT_CLOSE:
		return c.handleStmtClose(data)
	case mysql.COM_STMT_SEND_LONG_DATA:
		return c.handleStmtSendLongData(data)
	case mysql.COM_STMT_RESET:
		return c.handleStmtReset(data)
	case mysql.COM_SET_OPTION:
		return c.writeEOF(0)
	default:
		msg := fmt.Sprintf("command %d not supported now", cmd)
		golog.Error("ClientConn", "dispatch", msg, 0)
		return mysql.NewError(mysql.ER_UNKNOWN_ERROR, msg)
	}
}

func (c *ClientConn) writeOK(r *mysql.Result) error {
	if r == nil {
		r = &mysql.Result{Status: c.status}
	}
	data := make([]byte, 4, 32)

	data = append(data, mysql.OK_HEADER)

	data = append(data, mysql.PutLengthEncodedInt(r.AffectedRows)...)
	data = append(data, mysql.PutLengthEncodedInt(r.InsertId)...)

	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		data = append(data, byte(r.Status), byte(r.Status>>8))
		data = append(data, 0, 0)
	}

	return c.writePacket(data)
}

func (c *ClientConn) writeError(e error) error {
	var m *mysql.SqlError
	var ok bool
	if m, ok = e.(*mysql.SqlError); !ok {
		m = mysql.NewError(mysql.ER_UNKNOWN_ERROR, e.Error())
	}

	data := make([]byte, 4, 16+len(m.Message))

	data = append(data, mysql.ERR_HEADER)
	data = append(data, byte(m.Code), byte(m.Code>>8))

	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		data = append(data, '#')
		data = append(data, m.State...)
	}

	data = append(data, m.Message...)

	return c.writePacket(data)
}

func (c *ClientConn) writeEOF(status uint16) error {
	data := make([]byte, 4, 9)

	data = append(data, mysql.EOF_HEADER)
	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		data = append(data, 0, 0)
		data = append(data, byte(status), byte(status>>8))
	}

	return c.writePacket(data)
}

func (c *ClientConn) writeEOFBatch(total []byte, status uint16, direct bool) ([]byte, error) {
	data := make([]byte, 4, 9)

	data = append(data, mysql.EOF_HEADER)
	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		data = append(data, 0, 0)
		data = append(data, byte(status), byte(status>>8))
	}

	return c.writePacketBatch(total, data, direct)
}
