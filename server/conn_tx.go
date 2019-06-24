package server

func (c *ClientConn) handleBegin() (bool, error) {
	if err := c.begin(); err != nil {
		return false, err
	}

	return true, c.writeOK(nil)
}

func (c *ClientConn) handleCommit() (bool, error) {
	if err := c.commit(); err != nil {
		return false, err
	}

	return true, c.writeOK(nil)
}

func (c *ClientConn) handleRollback() (bool, error) {
	if err := c.rollback(); err != nil {
		return true, err
	}

	return true, c.writeOK(nil)
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
