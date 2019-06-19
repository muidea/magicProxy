package server

import (
	"fmt"

	"github.com/muidea/magicProxy/backend"
	"github.com/muidea/magicProxy/mysql"
)

func (c *ClientConn) handleUseDB(dbName string) error {
	var co *backend.BackendConn
	var err error

	if len(dbName) == 0 {
		return fmt.Errorf("must have database, the length of dbName is zero")
	}
	if c.schema == nil {
		return mysql.NewDefaultError(mysql.ER_NO_DB_ERROR)
	}

	nodeName := c.schema.rule.DefaultRule.Nodes[0]

	n := c.proxy.GetNode(nodeName)
	//get the connection from slave preferentially
	co, err = c.getBackendConn(n, true)
	defer c.closeConn(co, false)
	if err != nil {
		return err
	}

	if err = co.UseDB(dbName); err != nil {
		//reset the client database to null
		c.db = ""
		return err
	}
	c.db = dbName
	return c.writeOK(nil)
}
