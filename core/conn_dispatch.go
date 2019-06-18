package core

import (
	"bytes"
	"fmt"
	"log"

	"github.com/muidea/magicProxy/common/sql-parser/mysql"
)

func (c *Conn) dispatch(data []byte) error {
	if data == nil {
		return nil
	}
	cmd := data[0]
	data = data[1:]
	log.Printf("--> sql (%d),db:%s, cmd:%v, sql:%s", c.connectionId, c.connectionDB, cmd, string(data))
	switch cmd {
	case mysql.COM_QUIT: // 1
		c.handleRollback()
		c.Close()
		return nil
	case mysql.COM_INIT_DB: // 2
		return c.handleUseDB(mysql.Strings(data))
	case mysql.COM_QUERY: // 3
		return c.handleQuery(mysql.Strings(data))
	case mysql.COM_FIELD_LIST: // 4
		return c.handleFieldList(data)
	case mysql.COM_PING: // 14
		return c.writeOK(nil)
	case mysql.COM_STMT_PREPARE: // 22
		return c.handleStmtPrepare(mysql.Strings(data))
	case mysql.COM_STMT_EXECUTE: // 23
		return c.handleStmtExecute(data)
	case mysql.COM_STMT_SEND_LONG_DATA: // 24
		return c.handleStmtSendLongData(data)
	case mysql.COM_STMT_CLOSE: // 25
		return c.handleStmtClose(data)
	case mysql.COM_STMT_RESET: // 26
		return c.handleStmtReset(data)
	case mysql.COM_SET_OPTION: // 27
		return c.writeEOF(0)
	default:
		msg := fmt.Sprintf("command %d not supported now", cmd)
		log.Printf("Conn, dispatch, %s", msg)
		return mysql.NewError(mysql.ER_UNKNOWN_ERROR, msg)
	}
}

func (c *Conn) handleUseDB(dbName string) error {
	var err error

	if len(dbName) == 0 {
		return fmt.Errorf("must have database, the length of dbName is zero")
	}

	if err = c.down.UseDB(dbName); err != nil {
		//reset the client database to null
		c.connectionDB = ""
		return err
	}
	c.connectionDB = dbName

	return c.writeOK(nil)
}

func (c *Conn) handleFieldList(data []byte) error {
	index := bytes.IndexByte(data, 0x00)
	table := string(data[0:index])
	wildcard := string(data[index+1:])

	co := c.down

	if err := co.UseDB(c.connectionDB); err != nil {
		//reset the database to null
		c.connectionDB = ""
		return err
	}

	if fs, err := co.FieldList(table, wildcard); err != nil {
		return err
	} else {
		return c.writeFieldList(c.status, fs)
	}
}
func (c *Conn) writeFieldList(status uint16, fs []*mysql.Field) error {
	c.affectedRows = int64(-1)
	var err error
	total := make([]byte, 0, 1024)
	data := make([]byte, 4, 512)

	for _, v := range fs {
		data = data[0:4]
		data = append(data, v.Dump()...)
		total, err = c.WritePacketBatch(total, data, false)
		if err != nil {
			return err
		}
	}

	_, err = c.writeEOFBatch(total, status, true)
	return err
}
