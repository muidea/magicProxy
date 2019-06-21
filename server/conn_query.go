package server

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/muidea/magicProxy/backend"
	"github.com/muidea/magicProxy/core/errors"
	"github.com/muidea/magicProxy/core/golog"
	"github.com/muidea/magicProxy/core/hack"
	"github.com/muidea/magicProxy/mysql"
	"github.com/muidea/magicProxy/sqlparser"
)

/*处理query语句*/
func (c *ClientConn) handleQuery(sql string) (err error) {
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
	hasHandled, err := c.preHandle(sql)
	if err != nil {
		golog.Error("server", "preHandle", err.Error(), 0, "sql", sql, "hasHandled", hasHandled)
		return err
	}
	if hasHandled {
		return nil
	}

	var stmt sqlparser.Statement
	stmt, err = sqlparser.Parse(sql) //解析sql语句,得到的stmt是一个interface
	if err != nil {
		golog.Error("ClientConn", "Parse", err.Error(), 0, "sql", sql)
		return err
	}

	switch v := stmt.(type) {
	case *sqlparser.Select:
		return c.handleSelect(v, nil)
	case *sqlparser.Insert:
		return c.handleExec(stmt, nil)
	case *sqlparser.Update:
		return c.handleExec(stmt, nil)
	case *sqlparser.Delete:
		return c.handleExec(stmt, nil)
	case *sqlparser.Replace:
		return c.handleExec(stmt, nil)
	case *sqlparser.Set:
		return c.handleSet(v, sql)
	case *sqlparser.Begin:
		return c.handleBegin()
	case *sqlparser.Commit:
		return c.handleCommit()
	case *sqlparser.Rollback:
		return c.handleRollback()
	case *sqlparser.UseDB:
		return c.handleUseDB(v.DB)
	case *sqlparser.Truncate:
		return c.handleExec(stmt, nil)
	default:
		return fmt.Errorf("statement %T not support now", stmt)
	}
}

func (c *ClientConn) getBackendConn() (co *backend.BackendConn, err error) {
	bkNode := c.proxy.GetBackendNode()
	if bkNode == nil {
		err = fmt.Errorf("nodefine backend node")

		golog.Error("ClientConn", "GetBackendNode", err.Error(), 0)
		return
	}
	co, err = bkNode.GetConn()
	if err != nil {
		golog.Error("ClientConn", "GetConn", err.Error(), 0)
		return
	}

	if err = co.UseDB(c.db); err != nil {
		c.db = ""
		return
	}

	if err = co.SetCharset(c.charset, c.collation); err != nil {
		return
	}

	return
}

func (c *ClientConn) executeInNode(conn *backend.BackendConn, sql string, args []interface{}) (*mysql.Result, error) {
	r, err := conn.Execute(sql, args...)
	if err != nil {
		return nil, err
	}

	return r, err
}

func (c *ClientConn) closeConn(conn *backend.BackendConn, rollback bool) {
	if rollback {
		conn.Rollback()
	}

	conn.Close()
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

	rs, err := c.executeInNode(conns, "sql", args)
	if err == nil {
		return c.writeOK(rs)
	}

	return err
}
