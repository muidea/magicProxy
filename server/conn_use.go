package server

import (
	"fmt"

	"github.com/muidea/magicProxy/backend"
)

func (c *ClientConn) handleUseDB(dbName string) (ret bool, err error) {
	var co *backend.BackendConn

	if len(dbName) == 0 {
		return false, fmt.Errorf("must have database, the length of dbName is zero")
	}
	//get the connection from slave preferentially
	co, err = c.getBackendConn()
	defer c.closeConn(co, false)
	if err != nil {
		return false, err
	}

	if err = co.UseDB(dbName); err != nil {
		//reset the client database to null
		c.db = ""
		return true, err
	}
	c.db = dbName
	return true, c.writeOK(nil)
}
