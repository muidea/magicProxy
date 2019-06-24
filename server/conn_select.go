package server

import (
	"fmt"

	"github.com/muidea/magicProxy/core/golog"
	"github.com/muidea/magicProxy/sqlparser"
)

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
