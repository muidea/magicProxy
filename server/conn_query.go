package server

import (
	"strings"

	"github.com/muidea/magicProxy/core/errors"
	"github.com/muidea/magicProxy/core/hack"
	"github.com/muidea/magicProxy/mysql"
)

func (c *ClientConn) handleQuery(sql string) (ret bool, err error) {
	tokens := strings.FieldsFunc(sql, hack.IsSqlSep)

	if len(tokens) == 0 {
		return false, errors.ErrCmdUnsupport
	}

	var dbName string
	tokenID, ok := mysql.PARSE_TOKEN_MAP[strings.ToLower(tokens[0])]
	if ok == true {
		switch tokenID {
		case mysql.TK_ID_USE:
			dbName = strings.Trim(tokens[1], "`")
			ret, err = c.handleUseDB(dbName)
		case mysql.TK_ID_BEGIN:
			ret, err = c.handleBegin()
		case mysql.TK_ID_START:
			secondID, secondOK := mysql.PARSE_TOKEN_MAP[strings.ToLower(tokens[1])]
			if secondOK && secondID == mysql.TK_ID_TRANSACTION {
				ret, err = c.handleBegin()
			}
		case mysql.TK_ID_COMMIT:
			ret, err = c.handleCommit()
		case mysql.TK_ID_ROLLBACK:
			ret, err = c.handleRollback()
		default:
			return false, nil
		}
	}

	return
}
