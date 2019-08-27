package server

import (
	"github.com/muidea/magicProxy/backend"
	"github.com/muidea/magicProxy/core/errors"
)

func (c *ClientConn) handleUseDB(dbName string) (ret bool, err error) {
	if dbName == "" {
		return true, errors.ErrNoDatabase
	}

	var co *backend.BackendConn

	//get the connection from slave preferentially
	co, err = c.getBackendConn()
	defer c.closeConn(co, false)
	if err != nil {
		return true, err
	}

	err = co.UseDB(dbName)
	c.currentDB = dbName

	//c.status = co.Status()

	c.checkStatus(co)

	return true, err
}
