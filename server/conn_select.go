package server

import (
	"bytes"
	"fmt"

	"github.com/muidea/magicProxy/core/golog"
	"github.com/muidea/magicProxy/core/hack"
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

	if fs, err := co.FieldList(table, wildcard); err != nil {
		return err
	} else {
		return c.writeFieldList(c.status, fs)
	}
}

func (c *ClientConn) writeFieldList(status uint16, fs []*mysql.Field) error {
	c.affectedRows = int64(-1)
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
		return c.writeOK(nil)
	}

	// panic("todo")
	rs, err := c.executeInNode(conns, "sql", args)
	if err == nil {
		return c.writeOK(rs)
		//err = c.mergeExecResult(rs)
	}

	return err
}

//only process last_inser_id
func (c *ClientConn) handleSimpleSelect(stmt *sqlparser.SimpleSelect) error {
	nonStarExpr, _ := stmt.SelectExprs[0].(*sqlparser.NonStarExpr)
	var name string = hack.String(nonStarExpr.As)
	if name == "" {
		name = "last_insert_id()"
	}
	var column = 1
	var rows [][]string
	var names []string = []string{
		name,
	}

	var t = fmt.Sprintf("%d", c.lastInsertID)
	rows = append(rows, []string{t})

	r := new(mysql.Resultset)

	var values [][]interface{} = make([][]interface{}, len(rows))
	for i := range rows {
		values[i] = make([]interface{}, column)
		for j := range rows[i] {
			values[i][j] = rows[i][j]
		}
	}

	r, _ = c.buildResultset(nil, names, values)
	return c.writeResultset(c.status, r)
}
