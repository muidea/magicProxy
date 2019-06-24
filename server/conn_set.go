package server

import (
	"github.com/muidea/magicProxy/sqlparser"
)

var nstring = sqlparser.String

func (c *ClientConn) handleSet(sql string) (bool, error) {
	return true, c.writeOK(nil)
}

func (c *ClientConn) handleSetOption(sql string) (bool, error) {
	return true, c.writeEOF(0)
}
