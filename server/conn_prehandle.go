package server

import (
	"strings"

	"github.com/muidea/magicProxy/core/errors"
	"github.com/muidea/magicProxy/core/hack"
	"github.com/muidea/magicProxy/mysql"
)

//preprocessing sql before parse sql
func (c *ClientConn) preHandle(sql string) (bool, error) {
	if len(sql) == 0 {
		return false, errors.ErrCmdUnsupport
	}

	tokens := strings.FieldsFunc(sql, hack.IsSqlSep)
	if len(tokens) == 0 {
		return false, errors.ErrCmdUnsupport
	}

	handle, err := c.handleSQL(tokens, sql)
	if err != nil {
		//this SQL doesn't need execute in the backend.
		if err == errors.ErrIgnoreSQL {
			err = c.writeOK(nil)
			if err != nil {
				return false, err
			}
			return true, nil
		}

		return false, err
	}

	return handle, nil
}

//if sql need shard return nil, else return the unshard db
func (c *ClientConn) handleSQL(tokens []string, sql string) (bool, error) {
	tokensLen := len(tokens)
	if 0 < tokensLen {
		tokenID, ok := mysql.PARSE_TOKEN_MAP[strings.ToLower(tokens[0])]
		if ok == true {
			switch tokenID {
			case mysql.TK_ID_SELECT:
				return c.handleSelectToken(sql, tokens, tokensLen)
			case mysql.TK_ID_DELETE:
				return c.handleDeleteToken(sql, tokens, tokensLen)
			case mysql.TK_ID_INSERT, mysql.TK_ID_REPLACE:
				return c.handleInsertOrReplaceToken(sql, tokens, tokensLen)
			case mysql.TK_ID_UPDATE:
				return c.handleUpdateToken(sql, tokens, tokensLen)
			case mysql.TK_ID_SET:
				return c.handleSetToken(sql, tokens, tokensLen)
			case mysql.TK_ID_SHOW:
				return c.handleShowToken(sql, tokens, tokensLen)
			case mysql.TK_ID_TRUNCATE:
				return c.handleTruncateToken(sql, tokens, tokensLen)
			default:
				return false, nil
			}
		}
	}

	return false, nil
}

//get the execute database for select sql
func (c *ClientConn) handleSelectToken(sql string, tokens []string, tokensLen int) (bool, error) {
	selectLastInsertID := false
	selectCurrentUser := false
	selectConnectionID := false
	for i := 1; i < tokensLen; i++ {
		if strings.ToLower(tokens[i]) == mysql.TK_STR_LAST_INSERT_ID {
			selectLastInsertID = true
		}
		if strings.ToLower(tokens[i]) == mysql.TK_STR_CURRENT_USER {
			selectCurrentUser = true
		}
		if strings.ToLower(tokens[i]) == mysql.TK_STR_CONNECTION_ID {
			selectConnectionID = true
		}
	}

	if selectLastInsertID || selectCurrentUser || selectConnectionID {
		err := c.executeSQL(sql)
		return true, err
	}

	return false, nil
}

//get the execute database for delete sql
func (c *ClientConn) handleDeleteToken(sql string, tokens []string, tokensLen int) (bool, error) {
	return false, nil
}

//get the execute database for insert or replace sql
func (c *ClientConn) handleInsertOrReplaceToken(sql string, tokens []string, tokensLen int) (bool, error) {
	return false, nil
}

//get the execute database for update sql
func (c *ClientConn) handleUpdateToken(sql string, tokens []string, tokensLen int) (bool, error) {
	return false, nil
}

//get the execute database for set sql
func (c *ClientConn) handleSetToken(sql string, tokens []string, tokensLen int) (bool, error) {

	//handle three styles:
	//set autocommit= 0
	//set autocommit = 0
	//set autocommit=0
	if 2 <= tokensLen {
		before := strings.Split(sql, "=")
		//uncleanWorld is 'autocommit' or 'autocommit '
		uncleanWord := strings.Split(before[0], " ")
		secondWord := strings.ToLower(uncleanWord[1])

		if _, ok := mysql.SET_KEY_WORDS[secondWord]; ok {
			return false, nil
		}

		//SET [gobal/session] TRANSACTION ISOLATION LEVEL SERIALIZABLE
		//ignore this sql
		if 3 <= len(uncleanWord) {
			if strings.ToLower(uncleanWord[1]) == mysql.TK_STR_TRANSACTION ||
				strings.ToLower(uncleanWord[2]) == mysql.TK_STR_TRANSACTION {
				return false, errors.ErrIgnoreSQL
			}

			if strings.ToLower(uncleanWord[1]) == mysql.TK_STR_CHARACTER ||
				strings.ToLower(uncleanWord[2]) == mysql.TK_STR_CHARACTER {
				return false, errors.ErrIgnoreSQL
			}
		}
	}

	return false, nil
}

//get the execute database for show sql
//choose slave preferentially
//tokens[0] is show
func (c *ClientConn) handleShowToken(sql string, tokens []string, tokensLen int) (bool, error) {
	err := c.executeSQL(sql)
	return true, err
}

//get the execute database for truncate sql
//sql: TRUNCATE [TABLE] tbl_name
func (c *ClientConn) handleTruncateToken(sql string, tokens []string, tokensLen int) (bool, error) {
	return false, nil
}

func (c *ClientConn) executeSQL(sql string) error {
	//get connection in DB
	conn, err := c.getBackendConn()
	defer c.closeConn(conn, false)
	if err != nil {
		return err
	}

	//execute.sql may be rewritten in getShowExecDB
	rs, err := c.executeInNode(conn, sql, nil)
	if err != nil {
		return err
	}

	if rs.Resultset != nil {
		err = c.writeResultset(c.status, rs.Resultset)
	} else {
		err = c.writeOK(rs)
	}

	return err
}
