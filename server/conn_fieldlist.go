package server

import (
	"bytes"

	"github.com/muidea/magicProxy/mysql"
)

func (c *ClientConn) handleFieldList(sql string) (bool, error) {
	data := []byte(sql)
	index := bytes.IndexByte(data, 0x00)
	table := string(data[0:index])
	wildcard := string(data[index+1:])

	co, err := c.getBackendConn()
	defer c.closeConn(co, false)
	if err != nil {
		return false, err
	}

	fs, err := co.FieldList(table, wildcard)
	if err != nil {
		return true, err
	}

	return true, c.writeFieldList(c.status, fs)
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
