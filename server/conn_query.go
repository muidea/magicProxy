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
		default:
			return false, nil
		}
	}

	return
}
