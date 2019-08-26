package server

func (c *ClientConn) isBindConnection() bool {
	return c.bindConnectionFlag
}

func (c *ClientConn) onBindConnection() {
	c.bindConnectionFlag = true
}

func (c *ClientConn) offBindConnection() {
	c.bindConnectionFlag = false
}

func (c *ClientConn) handleSet(sql string) (ret bool, err error) {

	c.onBindConnection()

	err = c.executeSQL(sql)

	return true, err
}
