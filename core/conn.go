package core

import (
	"fmt"
	"hash"
	"log"
	"net"
	"runtime"
	"strconv"

	"github.com/muidea/magicProxy/common/errors"
	"github.com/muidea/magicProxy/common/sql-parser/mysql"
)

var DEFAULT_CAPABILITY uint32 = mysql.CLIENT_LONG_PASSWORD | mysql.CLIENT_LONG_FLAG |
	mysql.CLIENT_CONNECT_WITH_DB | mysql.CLIENT_PROTOCOL_41 |
	mysql.CLIENT_TRANSACTIONS | mysql.CLIENT_SECURE_CONNECTION

var baseConnId uint32 = 10000

type Conn struct {
	hash hash.Hash
	raw  net.Conn
	*mysql.PacketIO

	down Backend

	connectionId uint32
	stmtId       uint32
	stmts        map[uint32]*Stmt //prepare相关,client端到proxy的stmt

	lastInsertId int64
	affectedRows int64

	status uint16

	collation    mysql.CollationId
	charset      string
	connectionDB string
	capability   uint32
	salt         []byte

	shutdown bool
	curDB    string
}

func (c *Conn) Close() {
	if c.shutdown {
		return
	}
	c.shutdown = true
	c.raw.Close()
	return
}

func (c *Conn) Run(down Backend) {
	c.down = down
	c.shutdown = false
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("Conn, Run, %d, %s %s", c.connectionId, err.Error(), string(buf))
		}
		c.Close()
	}()

	for {
		data, err := c.ReadPacket()

		if err != nil {
			log.Printf("Conn, Run, ReadPacket, %d, %s %s", c.connectionId, err.Error(), string(data))
			c.Close()
			break
		}

		if err = c.dispatch(data); err != nil {
			log.Printf("Conn, Run, dispatch, %d, %s %s", c.connectionId, err.Error(), string(data))
			c.writeErr(err)
			if err == mysql.ErrBadConn {
				down.ReConnect()
				c.Close()
				break
			}
		}

		if c.shutdown {
			log.Printf("Conn, Run, %d, %s", c.connectionId, "is going down")
			break
		}

		c.PacketIO.Sequence = 0
	}
}

// Handshake func
func (c *Conn) Handshake() error {
	if err := c.writeInitialHandshake(); err != nil {
		log.Printf("Conn, Handshake, %d, %s", c.connectionId, err.Error())
		return err
	}
	if err := c.readHandshakeResponse(); err != nil {
		log.Printf("Conn, Handshake, %d, %s", c.connectionId, err.Error())
		return err
	}
	if err := c.writeOK(nil); err != nil {
		log.Printf("Conn, Handshake, %d, %s", c.connectionId, err.Error())
		return err
	}

	c.PacketIO.Sequence = 0
	return nil
}

func (c *Conn) writeErr(e error) error {
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

	return c.WritePacket(data)
}

func (c *Conn) writeOK(r *mysql.Result) error {
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

	return c.WritePacket(data)
}

func (c *Conn) writeEOF(status uint16) error {
	data := make([]byte, 4, 9)

	data = append(data, mysql.EOF_HEADER)
	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		data = append(data, 0, 0)
		data = append(data, byte(status), byte(status>>8))
	}

	return c.WritePacket(data)
}

func (c *Conn) writeEOFBatch(total []byte, status uint16, direct bool) ([]byte, error) {
	data := make([]byte, 4, 9)

	data = append(data, mysql.EOF_HEADER)
	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		data = append(data, 0, 0)
		data = append(data, byte(status), byte(status>>8))
	}

	return c.WritePacketBatch(total, data, direct)
}

func (c *Conn) writeResultset(status uint16, r *mysql.Resultset) error {
	c.affectedRows = int64(-1)
	total := make([]byte, 0, 4096)
	data := make([]byte, 4, 512)
	var err error

	columnLen := mysql.PutLengthEncodedInt(uint64(len(r.Fields)))
	log.Printf("columnLen:%d", columnLen)

	data = append(data, columnLen...)
	total, err = c.WritePacketBatch(total, data, false)
	if err != nil {
		return err
	}

	for _, v := range r.Fields {
		data = data[0:4]
		data = append(data, v.Dump()...)
		total, err = c.WritePacketBatch(total, data, false)
		if err != nil {
			return err
		}
	}

	total, err = c.writeEOFBatch(total, status, false)
	if err != nil {
		return err
	}

	for _, v := range r.RowDatas {
		data = data[0:4]
		data = append(data, v...)
		total, err = c.WritePacketBatch(total, data, false)
		if err != nil {
			return err
		}
	}

	total, err = c.writeEOFBatch(total, status, true)
	total = nil
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) buildResultset(fields []*mysql.Field, names []string, values [][]interface{}) (*mysql.Resultset, error) {
	var ExistFields bool
	r := new(mysql.Resultset)

	r.Fields = make([]*mysql.Field, len(names))
	r.FieldNames = make(map[string]int, len(names))

	//use the field def that get from true database
	if len(fields) != 0 {
		if len(r.Fields) == len(fields) {
			ExistFields = true
		} else {
			return nil, errors.ErrInvalidArgument
		}
	}

	var b []byte
	var err error

	for i, vs := range values {
		if len(vs) != len(r.Fields) {
			return nil, fmt.Errorf("row %d has %d column not equal %d", i, len(vs), len(r.Fields))
		}

		var row []byte
		for j, value := range vs {
			//列的定义
			if i == 0 {
				if ExistFields {
					r.Fields[j] = fields[j]
					r.FieldNames[string(r.Fields[j].Name)] = j
				} else {
					field := &mysql.Field{}
					r.Fields[j] = field
					r.FieldNames[string(r.Fields[j].Name)] = j
					field.Name = mysql.Slice(names[j])
					if err = formatField(field, value); err != nil {
						return nil, err
					}
				}

			}
			b, err = formatValue(value)
			if err != nil {
				return nil, err
			}

			row = append(row, mysql.PutLengthEncodedString(b)...)
		}

		r.RowDatas = append(r.RowDatas, row)
	}
	//assign the values to the result
	r.Values = values

	return r, nil
}
func formatField(field *mysql.Field, value interface{}) error {
	switch value.(type) {
	case int8, int16, int32, int64, int:
		field.Charset = 63
		field.Type = mysql.MYSQL_TYPE_LONGLONG
		field.Flag = mysql.BINARY_FLAG | mysql.NOT_NULL_FLAG
	case uint8, uint16, uint32, uint64, uint:
		field.Charset = 63
		field.Type = mysql.MYSQL_TYPE_LONGLONG
		field.Flag = mysql.BINARY_FLAG | mysql.NOT_NULL_FLAG | mysql.UNSIGNED_FLAG
	case float32, float64:
		field.Charset = 63
		field.Type = mysql.MYSQL_TYPE_DOUBLE
		field.Flag = mysql.BINARY_FLAG | mysql.NOT_NULL_FLAG
	case string, []byte:
		field.Charset = 33
		field.Type = mysql.MYSQL_TYPE_VAR_STRING
	default:
		return fmt.Errorf("unsupport type %T for resultset", value)
	}
	return nil
}

func formatValue(value interface{}) ([]byte, error) {
	if value == nil {
		return mysql.Slice("NULL"), nil
	}
	switch v := value.(type) {
	case int8:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int16:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int32:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int64:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case uint8:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint16:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint32:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint64:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case float32:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case []byte:
		return v, nil
	case string:
		return mysql.Slice(v), nil
	default:
		return nil, fmt.Errorf("invalid type %T", value)
	}
}
