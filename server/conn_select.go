package server

import (
	"bytes"
	"fmt"

	"github.com/muidea/magicProxy/core/golog"
	"github.com/muidea/magicProxy/mysql"
	"github.com/muidea/magicProxy/sqlparser"
)

func (c *ClientConn) handleFieldList(data []byte) error {
	index := bytes.IndexByte(data, 0x00)
	table := string(data[0:index])
	wildcard := string(data[index+1:])

	co, err := c.getBackendConn()
	defer c.closeConn(co, false)
	if err != nil {
		return err
	}

	if err = co.UseDB(c.db); err != nil {
		//reset the database to null
		c.db = ""
		return err
	}

	fs, err := co.FieldList(table, wildcard)
	if err != nil {
		return err
	}

	return c.writeFieldList(c.status, fs)
}

func (c *ClientConn) writeFieldList(status uint16, fs []*mysql.Field) error {
	var err error
	total := make([]byte, 0, 1024)
	data := make([]byte, 4, 512)

	for _, v := range fs {
		data = data[0:4]
		data = append(data, v.Dump()...)
		total, err = c.writePacketBatch(total, data, false)
		if err != nil {
			return err
		}
	}

	_, err = c.writeEOFBatch(total, status, true)
	return err
}

//处理select语句
func (c *ClientConn) handleSelect(stmt *sqlparser.Select, args []interface{}) error {
	conns, err := c.getBackendConn()
	defer c.closeConn(conns, false)
	if err != nil {
		golog.Error("ClientConn", "handleExec", err.Error(), c.connectionID)
		return err
	}

	if conns == nil {
		err = fmt.Errorf("can't get backend connection")
		golog.Error("ClientConn", "handleExec", err.Error(), c.connectionID)
		return err
	}

	rs, err := c.executeInConn(conns, "sql", args)
	if err == nil {
		return c.writeOK(rs)
	}

	return err
}
