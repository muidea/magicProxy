package server

func (c *ClientConn) handleQuit() (ret bool, err error) {
	ret, err = c.handleRollback()
	c.Close()
	return true, err
}
