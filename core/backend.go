package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/muidea/magicProxy/common/sql-parser/mysql"
)

type Executor interface {
	Execute(command string, args ...interface{}) (*mysql.Result, error)
}

type Backend interface {
	Close() error
	Executor

	UseDB(string) error
	Prepare(string) (*backStmt, error)
	ClosePrepare(uint32) error
	FieldList(table string, wildcard string) ([]*mysql.Field, error)
	SetAutoCommit(n uint8) error
	Begin() error
	Commit() error
	Rollback() error
	Connect(addr string, user string, password string) error
	ReConnect() error
	Ping() error
}

var (
	_ Backend = &backConn{}

	pingPeriod = int64(time.Second * 16)
)

//proxy <-> mysql server
type backConn struct {
	conn net.Conn

	pkg *mysql.PacketIO

	addr     string
	user     string
	password string
	db       string

	capability uint32

	status uint16

	collation mysql.CollationId
	charset   string
	salt      []byte

	pushTimestamp int64
	pkgErr        error
}

func (c *backConn) Connect(addr string, user string, password string) error {
	c.addr = addr
	c.user = user
	c.password = password

	//use utf8
	c.collation = mysql.DEFAULT_COLLATION_ID
	c.charset = mysql.DEFAULT_CHARSET

	return c.ReConnect()
}

func (c *backConn) ReConnect() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	n := "tcp"
	if strings.Contains(c.addr, "/") {
		n = "unix"
	}

	netConn, err := net.Dial(n, c.addr)
	if err != nil {
		return err
	}

	tcpConn := netConn.(*net.TCPConn)

	//SetNoDelay controls whether the operating system should delay packet transmission
	// in hopes of sending fewer packets (Nagle's algorithm).
	// The default is true (no delay),
	// meaning that data is sent as soon as possible after a Write.
	//I set this option false.
	tcpConn.SetNoDelay(false)
	tcpConn.SetKeepAlive(true)
	c.conn = tcpConn
	c.pkg = mysql.NewPacketIO(tcpConn)

	if err := c.readInitialHandshake(); err != nil {
		c.conn.Close()
		c.conn = nil
		return err
	}

	if err := c.writeAuthHandshake(); err != nil {
		c.conn.Close()
		c.conn = nil

		return err
	}

	if _, err := c.readOK(); err != nil {
		c.conn.Close()
		c.conn = nil

		return err
	}

	//we must always use autocommit
	if !c.IsAutoCommit() {
		if _, err := c.exec("set autocommit = 1;"); err != nil {
			c.conn.Close()
			c.conn = nil

			return err
		}
	}

	if _, err := c.exec("set global wait_timeout=288000;"); err != nil {
		c.conn.Close()
		c.conn = nil

		return err
	}
	if _, err := c.exec("set global interactive_timeout=288000;"); err != nil {
		c.conn.Close()
		c.conn = nil

		return err
	}
	return nil
}

func (c *backConn) Close() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
		c.salt = nil
		c.pkgErr = nil
	}

	return nil
}

func (c *backConn) readPacket() ([]byte, error) {
	d, err := c.pkg.ReadPacket()
	c.pkgErr = err
	return d, err
}

func (c *backConn) writePacket(data []byte) error {
	err := c.pkg.WritePacket(data)
	c.pkgErr = err
	return err
}

func (c *backConn) readInitialHandshake() error {
	data, err := c.readPacket()
	if err != nil {
		return err
	}

	if data[0] == mysql.ERR_HEADER {
		return errors.New("read initial handshake error")
	}

	if data[0] < mysql.MinProtocolVersion {
		return fmt.Errorf("invalid protocol version %d, must >= 10", data[0])
	}

	//skip mysql version and connection id
	//mysql version end with 0x00
	//connection id length is 4
	pos := 1 + bytes.IndexByte(data[1:], 0x00) + 1 + 4

	c.salt = append(c.salt, data[pos:pos+8]...)

	//skip filter
	pos += 8 + 1

	//capability lower 2 bytes
	c.capability = uint32(binary.LittleEndian.Uint16(data[pos : pos+2]))

	pos += 2

	if len(data) > pos {
		//skip server charset
		//c.charset = data[pos]
		pos += 1

		c.status = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2

		c.capability = uint32(binary.LittleEndian.Uint16(data[pos:pos+2]))<<16 | c.capability

		pos += 2

		//skip auth data len or [00]
		//skip reserved (all [00])
		pos += 10 + 1

		// The documentation is ambiguous about the length.
		// The official Python library uses the fixed length 12
		// mysql-proxy also use 12
		// which is not documented but seems to work.
		c.salt = append(c.salt, data[pos:pos+12]...)
	}

	return nil
}

func (c *backConn) writeAuthHandshake() error {
	// Adjust client capability flags based on server support
	capability := mysql.CLIENT_PROTOCOL_41 | mysql.CLIENT_SECURE_CONNECTION |
		mysql.CLIENT_LONG_PASSWORD | mysql.CLIENT_TRANSACTIONS | mysql.CLIENT_LONG_FLAG

	capability &= c.capability

	//packet length
	//capbility 4
	//max-packet size 4
	//charset 1
	//reserved all[0] 23
	length := 4 + 4 + 1 + 23

	//username
	length += len(c.user) + 1

	//we only support secure connection
	auth := mysql.CalcPassword(c.salt, []byte(c.password))

	length += 1 + len(auth)

	c.capability = capability

	data := make([]byte, length+4)

	//capability [32 bit]
	data[4] = byte(capability)
	data[5] = byte(capability >> 8)
	data[6] = byte(capability >> 16)
	data[7] = byte(capability >> 24)

	//MaxPacketSize [32 bit] (none)
	//data[8] = 0x00
	//data[9] = 0x00
	//data[10] = 0x00
	//data[11] = 0x00

	//Charset [1 byte]
	data[12] = byte(c.collation)

	//Filler [23 bytes] (all 0x00)
	pos := 13 + 23

	//User [null terminated string]
	if len(c.user) > 0 {
		pos += copy(data[pos:], c.user)
	}
	//data[pos] = 0x00
	pos++

	// auth [length encoded integer]
	data[pos] = byte(len(auth))
	pos += 1 + copy(data[pos+1:], auth)

	// db [null terminated string]
	if len(c.db) > 0 {
		pos += copy(data[pos:], c.db)
		//data[pos] = 0x00
	}

	return c.writePacket(data)
}

func (c *backConn) writeCommand(command byte) error {
	c.pkg.Sequence = 0

	return c.writePacket([]byte{
		0x01, //1 bytes long
		0x00,
		0x00,
		0x00, //sequence
		command,
	})
}

func (c *backConn) writeCommandBuf(command byte, arg []byte) error {
	c.pkg.Sequence = 0

	length := len(arg) + 1

	data := make([]byte, length+4)

	data[4] = command

	copy(data[5:], arg)

	return c.writePacket(data)
}

func (c *backConn) writeCommandStr(command byte, arg string) error {
	c.pkg.Sequence = 0

	length := len(arg) + 1

	data := make([]byte, length+4)

	data[4] = command

	copy(data[5:], arg)

	return c.writePacket(data)
}

func (c *backConn) writeCommandUint32(command byte, arg uint32) error {
	c.pkg.Sequence = 0

	return c.writePacket([]byte{
		0x05, //5 bytes long
		0x00,
		0x00,
		0x00, //sequence

		command,

		byte(arg),
		byte(arg >> 8),
		byte(arg >> 16),
		byte(arg >> 24),
	})
}

func (c *backConn) writeCommandStrStr(command byte, arg1 string, arg2 string) error {
	c.pkg.Sequence = 0

	data := make([]byte, 4, 6+len(arg1)+len(arg2))

	data = append(data, command)
	data = append(data, arg1...)
	data = append(data, 0)
	data = append(data, arg2...)

	return c.writePacket(data)
}

func (c *backConn) Ping() error {
	if err := c.writeCommand(mysql.COM_PING); err != nil {
		return err
	}

	if _, err := c.readOK(); err != nil {
		return err
	}

	return nil
}

func (c *backConn) UseDB(dbName string) error {
	if c.db == dbName || len(dbName) == 0 {
		return nil
	}

	if err := c.writeCommandStr(mysql.COM_INIT_DB, dbName); err != nil {
		return err
	}

	if _, err := c.readOK(); err != nil {
		return err
	}

	c.db = dbName
	return nil
}

func (c *backConn) GetDB() string {
	return c.db
}

func (c *backConn) GetAddr() string {
	return c.addr
}

func (c *backConn) Execute(command string, args ...interface{}) (*mysql.Result, error) {
	if len(args) == 0 {
		return c.exec(command)
	} else {
		if s, err := c.Prepare(command); err != nil {
			return nil, err
		} else {
			var r *mysql.Result
			r, err = s.Execute(args...)
			s.Close()
			return r, err
		}
	}
}

func (c *backConn) ClosePrepare(id uint32) error {
	return c.writeCommandUint32(mysql.COM_STMT_CLOSE, id)
}

func (c *backConn) Begin() error {
	_, err := c.exec("begin")
	return err
}

func (c *backConn) Commit() error {
	_, err := c.exec("commit")
	return err
}

func (c *backConn) Rollback() error {
	_, err := c.exec("rollback")
	return err
}

func (c *backConn) SetAutoCommit(n uint8) error {
	if n == 0 {
		if _, err := c.exec("set autocommit = 0"); err != nil {
			c.conn.Close()
			c.conn = nil

			return err
		}
	} else {
		if _, err := c.exec("set autocommit = 1"); err != nil {
			c.conn.Close()
			c.conn = nil

			return err
		}
	}
	return nil
}

func (c *backConn) SetCharset(charset string, collation mysql.CollationId) error {
	charset = strings.Trim(charset, "\"'`")

	if collation == 0 {
		collation = mysql.CollationNames[mysql.Charsets[charset]]
	}

	if c.charset == charset && c.collation == collation {
		return nil
	}

	_, ok := mysql.CharsetIds[charset]
	if !ok {
		return fmt.Errorf("invalid charset %s", charset)
	}

	_, ok = mysql.Collations[collation]
	if !ok {
		return fmt.Errorf("invalid collation %s", collation)
	}

	if _, err := c.exec(fmt.Sprintf("SET NAMES %s COLLATE %s", charset, mysql.Collations[collation])); err != nil {
		return err
	} else {
		c.collation = collation
		c.charset = charset
		return nil
	}
}

func (c *backConn) FieldList(table string, wildcard string) ([]*mysql.Field, error) {
	if err := c.writeCommandStrStr(mysql.COM_FIELD_LIST, table, wildcard); err != nil {
		return nil, err
	}

	data, err := c.readPacket()
	if err != nil {
		return nil, err
	}

	fs := make([]*mysql.Field, 0, 4)
	var f *mysql.Field
	if data[0] == mysql.ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	} else {
		for {
			if data, err = c.readPacket(); err != nil {
				return nil, err
			}

			// EOF Packet
			if c.isEOFPacket(data) {
				return fs, nil
			}

			if f, err = mysql.FieldData(data).Parse(); err != nil {
				return nil, err
			}
			fs = append(fs, f)
		}
	}
	return nil, fmt.Errorf("field list error")
}

func (c *backConn) exec(query string) (*mysql.Result, error) {
	if err := c.writeCommandStr(mysql.COM_QUERY, query); err != nil {
		return nil, err
	}

	return c.readResult(false)
}

func (c *backConn) readResultset(data []byte, binary bool) (*mysql.Result, error) {
	result := &mysql.Result{
		Status:       0,
		InsertId:     0,
		AffectedRows: 0,

		Resultset: &mysql.Resultset{},
	}

	// column count
	count, _, n := mysql.LengthEncodedInt(data)

	if n-len(data) != 0 {
		return nil, mysql.ErrMalformPacket
	}

	result.Fields = make([]*mysql.Field, count)
	result.FieldNames = make(map[string]int, count)

	if err := c.readResultColumns(result); err != nil {
		return nil, err
	}

	if err := c.readResultRows(result, binary); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *backConn) readResultColumns(result *mysql.Result) (err error) {
	var i int = 0
	var data []byte

	for {
		data, err = c.readPacket()
		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
				//result.Warnings = binary.LittleEndian.Uint16(data[1:])
				//todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				c.status = result.Status
			}

			if i != len(result.Fields) {
				err = mysql.ErrMalformPacket
			}

			return
		}

		result.Fields[i], err = mysql.FieldData(data).Parse()
		if err != nil {
			return
		}

		result.FieldNames[string(result.Fields[i].Name)] = i

		i++
	}
}

func (c *backConn) readResultRows(result *mysql.Result, isBinary bool) (err error) {
	var data []byte

	for {
		data, err = c.readPacket()

		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
				//result.Warnings = binary.LittleEndian.Uint16(data[1:])
				//todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				c.status = result.Status
			}

			break
		}

		result.RowDatas = append(result.RowDatas, data)
	}

	result.Values = make([][]interface{}, len(result.RowDatas))

	for i := range result.Values {
		result.Values[i], err = result.RowDatas[i].Parse(result.Fields, isBinary)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *backConn) readUntilEOF() (err error) {
	var data []byte

	for {
		data, err = c.readPacket()

		if err != nil {
			break
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			break
		}
	}

	return
}

func (c *backConn) isEOFPacket(data []byte) bool {
	return data[0] == mysql.EOF_HEADER && len(data) <= 5
}

func (c *backConn) handleOKPacket(data []byte) (*mysql.Result, error) {
	var n int
	var pos int = 1

	r := new(mysql.Result)

	r.AffectedRows, _, n = mysql.LengthEncodedInt(data[pos:])
	pos += n
	r.InsertId, _, n = mysql.LengthEncodedInt(data[pos:])
	pos += n

	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		pos += 2

		//todo:strict_mode, check warnings as error
		//Warnings := binary.LittleEndian.Uint16(data[pos:])
		//pos += 2
	} else if c.capability&mysql.CLIENT_TRANSACTIONS > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		pos += 2
	}

	//info
	return r, nil
}

func (c *backConn) handleErrorPacket(data []byte) error {
	e := new(mysql.SqlError)

	var pos int = 1

	e.Code = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		//skip '#'
		pos++
		e.State = string(data[pos : pos+5])
		pos += 5
	}

	e.Message = string(data[pos:])

	return e
}

func (c *backConn) readOK() (*mysql.Result, error) {
	data, err := c.readPacket()
	if err != nil {
		return nil, err
	}

	if data[0] == mysql.OK_HEADER {
		return c.handleOKPacket(data)
	} else if data[0] == mysql.ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	} else {
		return nil, errors.New("invalid ok packet")
	}
}

func (c *backConn) readResult(binary bool) (*mysql.Result, error) {
	data, err := c.readPacket()
	if err != nil {
		return nil, err
	}

	if data[0] == mysql.OK_HEADER {
		return c.handleOKPacket(data)
	} else if data[0] == mysql.ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	} else if data[0] == mysql.LocalInFile_HEADER {
		return nil, mysql.ErrMalformPacket
	}

	return c.readResultset(data, binary)
}

func (c *backConn) IsAutoCommit() bool {
	return c.status&mysql.SERVER_STATUS_AUTOCOMMIT > 0
}

func (c *backConn) IsInTransaction() bool {
	return c.status&mysql.SERVER_STATUS_IN_TRANS > 0
}

func (c *backConn) GetCharset() string {
	return c.charset
}
