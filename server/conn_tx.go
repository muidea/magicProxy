package server

import (
	"github.com/muidea/magicProxy/backend"
	"github.com/muidea/magicProxy/mysql"
)

func (c *ClientConn) isInTransaction() bool {
	return c.status&mysql.SERVER_STATUS_IN_TRANS > 0
}

func (c *ClientConn) handleBegin() (bool, error) {
	var co *backend.BackendConn

	//get the connection from slave preferentially
	co, err := c.getBackendConn()
	defer c.closeConn(co, false)
	if err != nil {
		return false, err
	}

	co.SetAutoCommit(0)

	c.status |= mysql.SERVER_STATUS_IN_TRANS
	return true, co.Begin()
}

func (c *ClientConn) handleCommit() (bool, error) {
	var co *backend.BackendConn

	//get the connection from slave preferentially
	co, err := c.getBackendConn()
	defer c.closeConn(co, false)
	if err != nil {
		return false, err
	}

	defer co.SetAutoCommit(1)

	c.status &= ^mysql.SERVER_STATUS_IN_TRANS
	return true, co.Commit()
}

func (c *ClientConn) handleRollback() (bool, error) {
	var co *backend.BackendConn

	//get the connection from slave preferentially
	co, err := c.getBackendConn()
	defer c.closeConn(co, false)
	if err != nil {
		return false, err
	}

	defer co.SetAutoCommit(1)

	c.status &= ^mysql.SERVER_STATUS_IN_TRANS

	return true, co.Rollback()
}
