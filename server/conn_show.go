package server

import "log"

func (c *ClientConn) handleShow(sql string) (bool, error) {
	log.Printf("handleShow,sql:%s", sql)
	return true, c.executeSQL(sql)
}
