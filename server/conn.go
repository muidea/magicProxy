package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"runtime"

	"github.com/muidea/magicProxy/backend"
	"github.com/muidea/magicProxy/core/golog"
	"github.com/muidea/magicProxy/core/hack"
	"github.com/muidea/magicProxy/mysql"
)

// ClientConn client <-> proxy
type ClientConn struct {
	//sync.Mutex

	pkg *mysql.PacketIO

	rawConn net.Conn

	proxy *Server

	capability uint32

	connectionID uint32

	status    uint16
	collation mysql.CollationId
	charset   string

	txConnection *backend.BackendConn

	bindConnectionFlag bool
	bindConnection     *backend.BackendConn

	user         string
	connectionDB string
	currentDB    string

	salt []byte

	stmtID uint32
	stmts  map[uint32]*Stmt //prepare相关,client端到proxy的stmt

	closed bool
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

	//auth_plugin_name[00]
	// magic rua !
	//data = append(data, []byte("mysql_native_password")...)
	//data = append(data, 0)

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
		pos += len(c.currentDB) + 1

	}
	c.connectionDB = db

	c.currentDB = db

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

			golog.Error("ClientConn", "Run", err.Error(), c.connectionID, "stack", string(buf))
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
			golog.Error("ClientConn", "Run", err.Error(), c.connectionID)

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
	sql := hack.String(data)

	golog.Info("ClientConn", "dispatch", "executeSQL", c.connectionID, "cmd", cmd, "sql", sql, "Sequence", c.pkg.Sequence)

	preHandle, preErr := c.preHandleSQL(cmd, sql)
	if preErr != nil || preHandle {
		golog.Info("ClientConn", "dispatch", "preHandleSQL", c.connectionID, "preHandle", preHandle)
		return preErr
	}

	return c.executeSQL(sql)
}

func (c *ClientConn) preHandleSQL(cmd byte, sql string) (ret bool, err error) {
	switch cmd {
	case mysql.COM_QUERY:
		ret, err = c.handleQuery(sql)
	case mysql.COM_INIT_DB:
		ret, err = c.handleUseDB(sql)
	case mysql.COM_PING:
		ret, err = c.handlePing()
	case mysql.COM_STMT_PREPARE:
		return c.handleStmtPrepare(sql)
	case mysql.COM_STMT_EXECUTE:
		return c.handleStmtExecute(sql)
	case mysql.COM_STMT_CLOSE:
		return c.handleStmtClose(sql)
	case mysql.COM_STMT_SEND_LONG_DATA:
		return c.handleStmtSendLongData(sql)
	case mysql.COM_STMT_RESET:
		return c.handleStmtReset(sql)
	default:
		ret = false
	}

	if ret && err == nil {
		c.writeOK(nil)
	}

	return
}

func (c *ClientConn) executeSQL(sql string) error {
	//get connection in DB
	co, err := c.getBackendConn()
	defer c.closeConn(co, false)
	if err != nil {
		return err
	}

	//execute.sql may be rewritten in getShowExecDB
	rs, err := c.executeInConn(co, sql, nil)
	if err != nil {
		return err
	}

	//c.status = co.Status()

	if rs.Resultset != nil {
		golog.Info("ClientConn", "executeSQL", "resultset", c.connectionID, "sql", sql)
		err = c.writeResultset(c.status, rs.Resultset)
	} else {
		golog.Info("ClientConn", "executeSQL", "write ok", c.connectionID, "sql", sql)
		err = c.writeOK(rs)
	}

	c.checkStatus(co)

	return err
}

func (c *ClientConn) allocConn() (co *backend.BackendConn, err error) {
	bkNode := c.proxy.GetBackendNode()
	if bkNode == nil {
		err = fmt.Errorf("nodefine backend node")

		golog.Error("ClientConn", "GetBackendNode", err.Error(), c.connectionID)
		return
	}
	co, err = bkNode.GetConn()
	if err != nil {
		golog.Error("ClientConn", "GetConn", err.Error(), c.connectionID)
		return
	}

	if err = co.UseDB(c.currentDB); err != nil {
		sqlErr, ok := err.(*mysql.SqlError)
		if ok {
			if sqlErr.Code == mysql.ER_NO_DB_ERROR || sqlErr.Code == mysql.ER_BAD_DB_ERROR {
				createSQL := fmt.Sprintf("CREATE SCHEMA `%s` DEFAULT CHARACTER SET %s COLLATE %s", c.currentDB, mysql.DEFAULT_CHARSET, mysql.DEFAULT_COLLATION_NAME)
				_, err = c.executeInConn(co, createSQL, nil)
				if err == nil {
					err = co.UseDB(c.currentDB)
				}
			}
		}

		if err != nil {
			c.currentDB = ""
			return
		}
	}

	if err = co.SetCharset(c.charset, c.collation); err != nil {
		return
	}

	c.checkStatus(co)

	return
}

func (c *ClientConn) getBackendConn() (co *backend.BackendConn, err error) {
	if c.isInTransaction() {
		if c.txConnection == nil {
			co, err = c.allocConn()
			if err != nil {
				return
			}

			co.SetAutoCommit(0)

			c.txConnection = co
		}

		co = c.txConnection

		return
	}

	if c.isBindConnection() {
		if c.bindConnection == nil {
			co, err = c.allocConn()
			if err != nil {
				return
			}

			c.bindConnection = co
		}

		co = c.bindConnection
		return
	}

	co, err = c.allocConn()
	return
}

func (c *ClientConn) executeInConn(conn *backend.BackendConn, sql string, args []interface{}) (*mysql.Result, error) {
	r, err := conn.Execute(sql, args...)
	if err != nil {
		return nil, err
	}

	return r, err
}

func (c *ClientConn) closeConn(conn *backend.BackendConn, rollback bool) {
	if c.isInTransaction() {
		return
	}

	if c.isBindConnection() {
		return
	}

	if rollback {
		conn.Rollback()
	}

	conn.Close()
}

func (c *ClientConn) checkStatus(conn *backend.BackendConn) {
	if conn == nil {
		return
	}

	clientStatus := c.status
	backendStatus := conn.Status()
	if clientStatus != backendStatus {
		golog.Error("ClientConn", "checkStatus", "mismatch status", c.connectionID, "ClientConn:", clientStatus, "BackendConn:", backendStatus)
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
