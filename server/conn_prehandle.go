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

	err := c.checkSQLTokens(tokens, sql)
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

	return false, nil
}

//if sql need shard return nil, else return the unshard db
func (c *ClientConn) checkSQLTokens(tokens []string, sql string) error {
	tokensLen := len(tokens)
	if 0 < tokensLen {
		tokenID, ok := mysql.PARSE_TOKEN_MAP[strings.ToLower(tokens[0])]
		if ok == true {
			switch tokenID {
			case mysql.TK_ID_SELECT:
				return c.checkSelect(sql, tokens, tokensLen)
			case mysql.TK_ID_DELETE:
				return c.checkDelete(sql, tokens, tokensLen)
			case mysql.TK_ID_INSERT, mysql.TK_ID_REPLACE:
				return c.checkInsertOrReplace(sql, tokens, tokensLen)
			case mysql.TK_ID_UPDATE:
				return c.checkUpdate(sql, tokens, tokensLen)
			case mysql.TK_ID_SET:
				return c.checkSet(sql, tokens, tokensLen)
			case mysql.TK_ID_SHOW:
				return c.checkShow(sql, tokens, tokensLen)
			case mysql.TK_ID_TRUNCATE:
				return c.checkTruncate(sql, tokens, tokensLen)
			default:
				return nil
			}
		}
	}

	return nil
}

//get the execute database for select sql
func (c *ClientConn) checkSelect(sql string, tokens []string, tokensLen int) error {
	return nil
}

//get the execute database for delete sql
func (c *ClientConn) checkDelete(sql string, tokens []string, tokensLen int) error {
	return nil
}

//get the execute database for insert or replace sql
func (c *ClientConn) checkInsertOrReplace(sql string, tokens []string, tokensLen int) error {
	return nil
}

//get the execute database for update sql
func (c *ClientConn) checkUpdate(sql string, tokens []string, tokensLen int) error {
	return nil
}

//get the execute database for set sql
func (c *ClientConn) checkSet(sql string, tokens []string, tokensLen int) error {
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
			return nil
		}

		//SET [gobal/session] TRANSACTION ISOLATION LEVEL SERIALIZABLE
		//ignore this sql
		if 3 <= len(uncleanWord) {
			if strings.ToLower(uncleanWord[1]) == mysql.TK_STR_TRANSACTION ||
				strings.ToLower(uncleanWord[2]) == mysql.TK_STR_TRANSACTION {
				return errors.ErrIgnoreSQL
			}
		}
	}

	return nil
}

//get the execute database for show sql
//choose slave preferentially
//tokens[0] is show
func (c *ClientConn) checkShow(sql string, tokens []string, tokensLen int) error {
	return nil
}

//get the execute database for truncate sql
//sql: TRUNCATE [TABLE] tbl_name
func (c *ClientConn) checkTruncate(sql string, tokens []string, tokensLen int) error {
	return nil
}
