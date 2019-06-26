package server

import (
	"github.com/muidea/magicProxy/backend"
	"github.com/muidea/magicProxy/mysql"
)

func (c *ClientConn) isInTransaction() bool {
	return c.status&mysql.SERVER_STATUS_IN_TRANS > 0
}

func (c *ClientConn) onInTransaction() {
	c.status |= mysql.SERVER_STATUS_IN_TRANS
}

func (c *ClientConn) offInTransaction() {
	c.status &= ^mysql.SERVER_STATUS_IN_TRANS
}

func (c *ClientConn) handleBegin() (ret bool, err error) {
	c.onInTransaction()

	var co *backend.BackendConn
	//get the connection from slave preferentially
	co, coErr := c.getBackendConn()
	defer c.closeConn(co, false)
	if coErr == nil {
		err = co.Begin()
		ret = true
	} else {
		err = coErr
	}

	if err != nil {
		c.offInTransaction()
	}

	return
}

func (c *ClientConn) handleCommit() (ret bool, err error) {
	var co *backend.BackendConn

	//get the connection from slave preferentially
	co, coErr := c.getBackendConn()
	defer c.closeConn(co, false)
	if coErr != nil {
		err = coErr
		return
	}

	err = co.Commit()
	ret = true

	c.offInTransaction()
	c.txConnection = nil

	return
}

func (c *ClientConn) handleRollback() (ret bool, err error) {
	var co *backend.BackendConn

	//get the connection from slave preferentially
	co, coErr := c.getBackendConn()
	defer c.closeConn(co, false)
	if coErr != nil {
		err = coErr
		return
	}

	err = co.Rollback()
	ret = true

	c.offInTransaction()
	c.txConnection = nil

	return
}
