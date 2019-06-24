package server

import (
	"github.com/muidea/magicProxy/backend"
)

func (c *ClientConn) handlePing() (ret bool, err error) {
	var co *backend.BackendConn

	//get the connection from slave preferentially
	co, err = c.getBackendConn()
	defer c.closeConn(co, false)
	if err != nil {
		return false, err
	}

	err = co.Ping()
	return true, err
}
