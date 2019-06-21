package server

import (
	"github.com/muidea/magicProxy/mysql"
)

func (c *ClientConn) handleBegin() error {
	c.status |= mysql.SERVER_STATUS_IN_TRANS
	return c.writeOK(nil)
}

func (c *ClientConn) handleCommit() (err error) {
	if err := c.commit(); err != nil {
		return err
	}

	return c.writeOK(nil)
}

func (c *ClientConn) handleRollback() (err error) {
	if err := c.rollback(); err != nil {
		return err
	}

	return c.writeOK(nil)
}

func (c *ClientConn) begin() (err error) {
	co, coErr := c.getBackendConn()
	if coErr != nil {
		err = coErr
		return
	}

	err = co.Begin()
	return
}

func (c *ClientConn) commit() (err error) {
	co, coErr := c.getBackendConn()
	if coErr != nil {
		err = coErr
		return
	}

	err = co.Commit()
	return
}

func (c *ClientConn) rollback() (err error) {
	co, coErr := c.getBackendConn()
	if coErr != nil {
		err = coErr
		return
	}

	err = co.Rollback()
	return
}
