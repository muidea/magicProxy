package server

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/muidea/magicProxy/core/errors"
	"github.com/muidea/magicProxy/core/golog"
	"github.com/muidea/magicProxy/core/hack"
	"github.com/muidea/magicProxy/mysql"
	"github.com/muidea/magicProxy/sqlparser"
)

/*处理query语句*/
func (c *ClientConn) handleQuery(sql string) (ret bool, err error) {
	defer func() {
		if e := recover(); e != nil {
			golog.OutputSql("Error", "err:%v,sql:%s", e, sql)

			if err, ok := e.(error); ok {
				const size = 4096
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]

				golog.Error("ClientConn", "handleQuery",
					err.Error(), 0,
					"stack", string(buf), "sql", sql)
			}

			err = errors.ErrInternalServer
			return
		}
	}()

	sql = strings.TrimRight(sql, ";") //删除sql语句最后的分号
	tokens := strings.FieldsFunc(sql, hack.IsSqlSep)
	tokensLen := len(tokens)
	if 0 < tokensLen {
		tokenID, ok := mysql.PARSE_TOKEN_MAP[strings.ToLower(tokens[0])]
		if ok == true {
			switch tokenID {
			case mysql.TK_ID_SELECT:
				return false, nil
			case mysql.TK_ID_DELETE:
				return false, nil
			case mysql.TK_ID_INSERT, mysql.TK_ID_REPLACE:
				return false, nil
			case mysql.TK_ID_UPDATE:
				return false, nil
			case mysql.TK_ID_SET:
				//return c.getSetExecDB(sql, tokens, tokensLen)
				return c.handleSet(sql)
			case mysql.TK_ID_SHOW:
				return c.handleShow(sql)
			case mysql.TK_ID_TRUNCATE:
				return false, nil
			}
		}
	}

	var stmt sqlparser.Statement
	stmt, err = sqlparser.Parse(sql) //解析sql语句,得到的stmt是一个interface
	if err != nil {
		golog.Error("ClientConn", "Parse", err.Error(), 0, "sql", sql)
		return false, err
	}

	switch v := stmt.(type) {
	case *sqlparser.DDL:
		return false, nil
	case *sqlparser.Select:
		return false, nil
	case *sqlparser.Insert:
		return false, nil
	case *sqlparser.Update:
		return false, nil
	case *sqlparser.Delete:
		return false, nil
	case *sqlparser.Replace:
		return false, nil
	case *sqlparser.Set:
		return false, nil
	case *sqlparser.Begin:
		return c.handleBegin()
	case *sqlparser.Commit:
		return c.handleCommit()
	case *sqlparser.Rollback:
		return c.handleRollback()
	case *sqlparser.UseDB:
		return c.handleUseDB(v.DB)
	case *sqlparser.Truncate:
		return false, nil
	default:
		return false, nil
	}
}

func (c *ClientConn) newEmptyResultset(stmt *sqlparser.Select) *mysql.Resultset {
	r := new(mysql.Resultset)
	r.Fields = make([]*mysql.Field, len(stmt.SelectExprs))

	for i, expr := range stmt.SelectExprs {
		r.Fields[i] = &mysql.Field{}
		switch e := expr.(type) {
		case *sqlparser.StarExpr:
			r.Fields[i].Name = []byte("*")
		case *sqlparser.NonStarExpr:
			if e.As != nil {
				r.Fields[i].Name = e.As
				r.Fields[i].OrgName = hack.Slice(nstring(e.Expr))
			} else {
				r.Fields[i].Name = hack.Slice(nstring(e.Expr))
			}
		default:
			r.Fields[i].Name = hack.Slice(nstring(e))
		}
	}

	r.Values = make([][]interface{}, 0)
	r.RowDatas = make([]mysql.RowData, 0)

	return r
}

func (c *ClientConn) handleExec(stmt sqlparser.Statement, args []interface{}) error {
	conns, err := c.getBackendConn()
	defer c.closeConn(conns, false)
	if err != nil {
		golog.Error("ClientConn", "handleExec", err.Error(), c.connectionID)
		return err
	}
	if conns == nil {
		err = fmt.Errorf("can't get backend connection")
		golog.Error("ClientConn", "handleExec", err.Error(), c.connectionID)
		return err
	}

	rs, err := c.executeInConn(conns, "sql", args)
	if err == nil {
		return c.writeOK(rs)
	}

	return err
}
