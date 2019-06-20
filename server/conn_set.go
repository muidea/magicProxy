package server

import (
	"fmt"
	"strings"

	"github.com/muidea/magicProxy/core/golog"
	"github.com/muidea/magicProxy/mysql"
	"github.com/muidea/magicProxy/sqlparser"
)

var nstring = sqlparser.String

func (c *ClientConn) handleSet(stmt *sqlparser.Set, sql string) (err error) {
	if len(stmt.Exprs) != 1 && len(stmt.Exprs) != 2 {
		return fmt.Errorf("must set one item once, not %s", nstring(stmt))
	}

	k := string(stmt.Exprs[0].Name.Name)
	switch strings.ToUpper(k) {
	case `AUTOCOMMIT`, `@@AUTOCOMMIT`, `@@SESSION.AUTOCOMMIT`:
		return c.handleSetAutoCommit(stmt.Exprs[0].Expr)
	case `NAMES`,
		`CHARACTER_SET_RESULTS`, `@@CHARACTER_SET_RESULTS`, `@@SESSION.CHARACTER_SET_RESULTS`,
		`CHARACTER_SET_CLIENT`, `@@CHARACTER_SET_CLIENT`, `@@SESSION.CHARACTER_SET_CLIENT`,
		`CHARACTER_SET_CONNECTION`, `@@CHARACTER_SET_CONNECTION`, `@@SESSION.CHARACTER_SET_CONNECTION`:
		if len(stmt.Exprs) == 2 {
			//SET NAMES 'charset_name' COLLATE 'collation_name'
			return c.handleSetNames(stmt.Exprs[0].Expr, stmt.Exprs[1].Expr)
		}
		return c.handleSetNames(stmt.Exprs[0].Expr, nil)
	default:
		golog.Error("ClientConn", "handleSet", "command not supported",
			c.connectionId, "sql", sql)
		return c.writeOK(nil)
	}
}

func (c *ClientConn) handleSetAutoCommit(val sqlparser.ValExpr) error {
	flag := sqlparser.String(val)
	flag = strings.Trim(flag, "'`\"")
	// autocommit允许为 0, 1, ON, OFF, "ON", "OFF", 不允许"0", "1"
	if flag == `0` || flag == `1` {
		_, ok := val.(sqlparser.NumVal)
		if !ok {
			return fmt.Errorf("set autocommit error")
		}
	}
	switch strings.ToUpper(flag) {
	case `1`, `ON`:
		c.status |= mysql.SERVER_STATUS_AUTOCOMMIT
		if c.status&mysql.SERVER_STATUS_IN_TRANS > 0 {
			c.status &= ^mysql.SERVER_STATUS_IN_TRANS
		}
	case `0`, `OFF`:
		c.status &= ^mysql.SERVER_STATUS_AUTOCOMMIT
	default:
		return fmt.Errorf("invalid autocommit flag %s", flag)
	}

	return c.writeOK(nil)
}

func (c *ClientConn) handleSetNames(ch, ci sqlparser.ValExpr) error {
	var cid mysql.CollationId
	var ok bool

	value := sqlparser.String(ch)
	value = strings.Trim(value, "'`\"")

	charset := strings.ToLower(value)
	if charset == "null" {
		return c.writeOK(nil)
	}
	if ci == nil {
		if charset == "default" {
			charset = mysql.DEFAULT_CHARSET
		}
		cid, ok = mysql.CharsetIds[charset]
		if !ok {
			return fmt.Errorf("invalid charset %s", charset)
		}
	} else {
		collate := sqlparser.String(ci)
		collate = strings.Trim(collate, "'`\"")
		collate = strings.ToLower(collate)
		cid, ok = mysql.CollationNames[collate]
		if !ok {
			return fmt.Errorf("invalid collation %s", collate)
		}
	}
	c.charset = charset
	c.collation = cid

	return c.writeOK(nil)
}