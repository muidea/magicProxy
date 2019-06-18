package core

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"

	"github.com/muidea/magicProxy/common/errors"
	"github.com/muidea/magicProxy/common/sql-parser/mysql"
	"github.com/muidea/magicProxy/common/sql-parser/sqlparser"
)

func (c *Conn) handleStmtExecute(data []byte) error {
	if len(data) < 9 {
		return mysql.ErrMalformPacket
	}

	pos := 0
	id := binary.LittleEndian.Uint32(data[0:4])
	pos += 4

	s, ok := c.stmts[id]
	if !ok {
		return mysql.NewDefaultError(mysql.ER_UNKNOWN_STMT_HANDLER,
			strconv.FormatUint(uint64(id), 10), "stmt_execute")
	}

	flag := data[pos]
	pos++
	//now we only support CURSOR_TYPE_NO_CURSOR flag
	if flag != 0 {
		return mysql.NewError(mysql.ER_UNKNOWN_ERROR, fmt.Sprintf("unsupported flag %d", flag))
	}

	//skip iteration-count, always 1
	pos += 4

	var nullBitmaps []byte
	var paramTypes []byte
	var paramValues []byte

	paramNum := s.params

	if paramNum > 0 {
		nullBitmapLen := (s.params + 7) >> 3
		if len(data) < (pos + nullBitmapLen + 1) {
			return mysql.ErrMalformPacket
		}
		nullBitmaps = data[pos : pos+nullBitmapLen]
		pos += nullBitmapLen

		//new param bound flag
		if data[pos] == 1 {
			pos++
			if len(data) < (pos + (paramNum << 1)) {
				return mysql.ErrMalformPacket
			}

			paramTypes = data[pos : pos+(paramNum<<1)]
			pos += (paramNum << 1)

			paramValues = data[pos:]
		}

		if err := c.bindStmtArgs(s, nullBitmaps, paramTypes, paramValues); err != nil {
			return err
		}
	}

	var err error

	switch stmt := s.s.(type) {
	case *sqlparser.Select:
		err = c.handlePrepareSelect(stmt, s.sql, s.args)
	case *sqlparser.Insert:
		err = c.handlePrepareExec(s.s, s.sql, s.args)
	case *sqlparser.Update:
		err = c.handlePrepareExec(s.s, s.sql, s.args)
	case *sqlparser.Delete:
		err = c.handlePrepareExec(s.s, s.sql, s.args)
	default:
		err = fmt.Errorf("command %T not supported now", stmt)
	}

	s.ResetParams()

	return err
}

func (c *Conn) bindStmtArgs(s *Stmt, nullBitmap, paramTypes, paramValues []byte) error {
	args := s.args

	pos := 0

	var v []byte
	var n int = 0
	var isNull bool
	var err error

	for i := 0; i < s.params; i++ {
		if nullBitmap[i>>3]&(1<<(uint(i)%8)) > 0 {
			args[i] = nil
			continue
		}

		tp := paramTypes[i<<1]
		isUnsigned := (paramTypes[(i<<1)+1] & 0x80) > 0

		switch tp {
		case mysql.MYSQL_TYPE_NULL:
			args[i] = nil
			continue

		case mysql.MYSQL_TYPE_TINY:
			if len(paramValues) < (pos + 1) {
				return mysql.ErrMalformPacket
			}

			if isUnsigned {
				args[i] = uint8(paramValues[pos])
			} else {
				args[i] = int8(paramValues[pos])
			}

			pos++
			continue

		case mysql.MYSQL_TYPE_SHORT, mysql.MYSQL_TYPE_YEAR:
			if len(paramValues) < (pos + 2) {
				return mysql.ErrMalformPacket
			}

			if isUnsigned {
				args[i] = uint16(binary.LittleEndian.Uint16(paramValues[pos : pos+2]))
			} else {
				args[i] = int16((binary.LittleEndian.Uint16(paramValues[pos : pos+2])))
			}
			pos += 2
			continue

		case mysql.MYSQL_TYPE_INT24, mysql.MYSQL_TYPE_LONG:
			if len(paramValues) < (pos + 4) {
				return mysql.ErrMalformPacket
			}

			if isUnsigned {
				args[i] = uint32(binary.LittleEndian.Uint32(paramValues[pos : pos+4]))
			} else {
				args[i] = int32(binary.LittleEndian.Uint32(paramValues[pos : pos+4]))
			}
			pos += 4
			continue

		case mysql.MYSQL_TYPE_LONGLONG:
			if len(paramValues) < (pos + 8) {
				return mysql.ErrMalformPacket
			}

			if isUnsigned {
				args[i] = binary.LittleEndian.Uint64(paramValues[pos : pos+8])
			} else {
				args[i] = int64(binary.LittleEndian.Uint64(paramValues[pos : pos+8]))
			}
			pos += 8
			continue

		case mysql.MYSQL_TYPE_FLOAT:
			if len(paramValues) < (pos + 4) {
				return mysql.ErrMalformPacket
			}

			args[i] = float32(math.Float32frombits(binary.LittleEndian.Uint32(paramValues[pos : pos+4])))
			pos += 4
			continue

		case mysql.MYSQL_TYPE_DOUBLE:
			if len(paramValues) < (pos + 8) {
				return mysql.ErrMalformPacket
			}

			args[i] = math.Float64frombits(binary.LittleEndian.Uint64(paramValues[pos : pos+8]))
			pos += 8
			continue

		case mysql.MYSQL_TYPE_DECIMAL, mysql.MYSQL_TYPE_NEWDECIMAL, mysql.MYSQL_TYPE_VARCHAR,
			mysql.MYSQL_TYPE_BIT, mysql.MYSQL_TYPE_ENUM, mysql.MYSQL_TYPE_SET, mysql.MYSQL_TYPE_TINY_BLOB,
			mysql.MYSQL_TYPE_MEDIUM_BLOB, mysql.MYSQL_TYPE_LONG_BLOB, mysql.MYSQL_TYPE_BLOB,
			mysql.MYSQL_TYPE_VAR_STRING, mysql.MYSQL_TYPE_STRING, mysql.MYSQL_TYPE_GEOMETRY,
			mysql.MYSQL_TYPE_DATE, mysql.MYSQL_TYPE_NEWDATE,
			mysql.MYSQL_TYPE_TIMESTAMP, mysql.MYSQL_TYPE_DATETIME, mysql.MYSQL_TYPE_TIME:
			if len(paramValues) < (pos + 1) {
				return mysql.ErrMalformPacket
			}

			v, isNull, n, err = mysql.LengthEnodedString(paramValues[pos:])
			pos += n
			if err != nil {
				return err
			}

			if !isNull {
				args[i] = v
				continue
			} else {
				args[i] = nil
				continue
			}
		default:
			return fmt.Errorf("Stmt Unknown FieldType %d", tp)
		}
	}
	return nil
}
func (c *Conn) handlePrepareExec(stmt sqlparser.Statement, sql string, args []interface{}) error {
	co := c.down

	var rs []*mysql.Result
	rs, err := c.executeInNode(co, sql, args)

	if err != nil {
		log.Printf("Conn, handlePrepareExec, %d, %s", c.connectionId, err.Error())
		return err
	}

	status := c.status | rs[0].Status
	if rs[0].Resultset != nil {
		err = c.writeResultset(status, rs[0].Resultset)
	} else {
		err = c.writeOK(rs[0])
	}

	return err
}
func (c *Conn) executeInNode(conn Backend, sql string, args []interface{}) ([]*mysql.Result, error) {
	r, err := conn.Execute(sql, args...)
	if err != nil {
		return nil, err
	}

	return []*mysql.Result{r}, err
}

func (c *Conn) handlePrepareSelect(stmt *sqlparser.Select, sql string, args []interface{}) error {
	co := c.down
	var rs []*mysql.Result
	rs, err := c.executeInNode(co, sql, args)

	if err != nil {
		log.Printf("Conn, handlePrepareExec, %d, %s", c.connectionId, err.Error())
		return err
	}

	status := c.status | rs[0].Status
	if rs[0].Resultset != nil {
		err = c.writeResultset(status, rs[0].Resultset)
	} else {
		r := c.newEmptyResultset(stmt)
		err = c.writeResultset(status, r)
	}

	return err
}

var nstring = sqlparser.String

func (c *Conn) newEmptyResultset(stmt *sqlparser.Select) *mysql.Resultset {
	r := new(mysql.Resultset)
	r.Fields = make([]*mysql.Field, len(stmt.SelectExprs))

	for i, expr := range stmt.SelectExprs {
		r.Fields[i] = &mysql.Field{}
		switch e := expr.(type) {
		case *sqlparser.StarExpr:
			r.Fields[i].Name = []byte("*")
		// case *sqlparser.Columns:
		// 	if e.As != nil {
		// 		r.Fields[i].Name = e.As
		// 		r.Fields[i].OrgName = mysql.Slice(nstring(e.Expr))
		// 	} else {
		// 		r.Fields[i].Name = mysql.Slice(nstring(e.Expr))
		// 	}
		default:
			r.Fields[i].Name = mysql.Slice(nstring(e))
		}
	}

	r.Values = make([][]interface{}, 0)
	r.RowDatas = make([]mysql.RowData, 0)

	return r
}

/*处理query语句*/
func (c *Conn) handleQuery(sql string) (err error) {
	defer func() {
		if e := recover(); e != nil {

			if err, ok := e.(error); ok {
				const size = 4096
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				log.Printf("Conn, handleQuery %s %s %s", c.raw.RemoteAddr().String(), err, string(buf))
			} else {
				log.Printf("Conn, handleQuery %s %s", c.raw.RemoteAddr().String(), err)
			}

			err = errors.ErrInternalServer
			return
		}
	}()

	if c.connectionDB != "" && c.curDB == "" {
		c.curDB = c.connectionDB

		str := fmt.Sprintf("USE %s", c.connectionDB)
		err := c.handleQuery(str)
		if err != nil {
			log.Printf("Conn, handleQuery, use database failed, err:%s", err.Error())
		}
	}

	var stmt sqlparser.Statement
	stmt, err = sqlparser.Parse(sql) //解析sql语句,得到的stmt是一个interface

	if err != nil {
		log.Printf("Conn, handleQuery, %s, %s", err.Error(), sql)
	}

	switch v := stmt.(type) {
	case *sqlparser.Select:
		return c.handleSelectExec(v, nil)
	case *sqlparser.Insert, *sqlparser.Update, *sqlparser.Delete:
		return c.handleExec(v, nil)
	case *sqlparser.Set, *sqlparser.DDL, *sqlparser.DBDDL:
		return c.handleExec2(sql, nil)
	case *sqlparser.Show:
		return c.handleExec2(sql, nil)
	case *sqlparser.Begin:
		return c.handleBegin()
	case *sqlparser.Commit:
		return c.handleCommit()
	case *sqlparser.Rollback:
		return c.handleRollback()
	case *sqlparser.Use:
		return c.handleUseDB(v.DBName.String())
	default:
		fmt.Errorf("command %T not supported now", stmt)
		return c.handleExec2(sql, nil)
	}

	return nil
}
func (c *Conn) handleExec2(sql string, args []interface{}) (err error) {
	var r *mysql.Result
	r, err = c.down.Execute(sql, args...)
	if err == nil {
		// sync if async failed
		if r.Resultset != nil {
			err = c.writeResultset(r.Status, r.Resultset)
		} else {
			err = c.writeOK(r)
		}
	}
	return
}

func (c *Conn) handleSelectExec(stmt sqlparser.Statement, args []interface{}) (err error) {
	sql := sqlparser.String(stmt)
	return c.handleExec2(sql, args)
}

func (c *Conn) handleExec(stmt sqlparser.Statement, args []interface{}) (err error) {
	sql := sqlparser.String(stmt)
	return c.handleExec2(sql, args)
}

func (c *Conn) mergeExecResult(rs []*mysql.Result) error {
	r := new(mysql.Result)
	for _, v := range rs {
		r.Status |= v.Status
		r.AffectedRows += v.AffectedRows
		if r.InsertId == 0 {
			r.InsertId = v.InsertId
		} else if r.InsertId > v.InsertId {
			//last insert id is first gen id for multi row inserted
			//see http://dev.mysql.com/doc/refman/5.6/en/information-functions.html#function_last-insert-id
			r.InsertId = v.InsertId
		}
	}

	if r.InsertId > 0 {
		c.lastInsertId = int64(r.InsertId)
	}
	c.affectedRows = int64(r.AffectedRows)

	return c.writeOK(r)
}

func (c *Conn) handleSet(stmt *sqlparser.Set, sql string) (err error) {
	if len(stmt.Exprs) != 1 && len(stmt.Exprs) != 2 {
		return fmt.Errorf("must set one item once, not %s", nstring(stmt))
	}

	k := string(stmt.Exprs[0].Name.String())
	switch strings.ToUpper(k) {
	case `AUTOCOMMIT`, `@@AUTOCOMMIT`, `@@SESSION.AUTOCOMMIT`:
		return c.handleSetAutoCommit(stmt.Exprs[0].Expr)
	// case `NAMES`,
	// 	`CHARACTER_SET_RESULTS`, `@@CHARACTER_SET_RESULTS`, `@@SESSION.CHARACTER_SET_RESULTS`,
	// 	`CHARACTER_SET_CLIENT`, `@@CHARACTER_SET_CLIENT`, `@@SESSION.CHARACTER_SET_CLIENT`,
	// 	`CHARACTER_SET_CONNECTION`, `@@CHARACTER_SET_CONNECTION`, `@@SESSION.CHARACTER_SET_CONNECTION`:
	// 	if len(stmt.Exprs) == 2 {
	// 		//SET NAMES 'charset_name' COLLATE 'collation_name'
	// 		return c.handleSetNames(stmt.Exprs[0].Expr, stmt.Exprs[1].Expr)
	// 	}
	// 	return c.handleSetNames(stmt.Exprs[0].Expr, nil)
	default:
		log.Printf("Conn, handleSet, command not supported, %d, %s", c.connectionId, sql)
		return c.writeOK(nil)
	}
}

func (c *Conn) handleSetAutoCommit(val sqlparser.Expr) error {
	flag := sqlparser.String(val)
	flag = strings.Trim(flag, "'`\"")
	// autocommit允许为 0, 1, ON, OFF, "ON", "OFF", 不允许"0", "1"
	if flag == `0` || flag == `1` {

		if !sqlparser.IsValue(val) {
			return fmt.Errorf("set autocommit error")
		}
	}
	switch strings.ToUpper(flag) {
	case `1`, `ON`:
		c.status |= mysql.SERVER_STATUS_AUTOCOMMIT
		if c.status&mysql.SERVER_STATUS_IN_TRANS > 0 {
			c.status &= ^mysql.SERVER_STATUS_IN_TRANS
		}
		co := c.down
		if e := co.SetAutoCommit(1); e != nil {
			co.Close()
			return fmt.Errorf("set autocommit error, %v", e)
		}
	case `0`, `OFF`:
		c.status &= ^mysql.SERVER_STATUS_AUTOCOMMIT
	default:
		return fmt.Errorf("invalid autocommit flag %s", flag)
	}

	return c.writeOK(nil)
}

// func (c *Conn) handleSetNames(ch, ci sqlparser.ValExpr) error {
// 	var cid mysql.CollationId
// 	var ok bool

// 	value := sqlparser.String(ch)
// 	value = strings.Trim(value, "'`\"")

// 	charset := strings.ToLower(value)
// 	if charset == "null" {
// 		return c.writeOK(nil)
// 	}
// 	if ci == nil {
// 		if charset == "default" {
// 			charset = mysql.DEFAULT_CHARSET
// 		}
// 		cid, ok = mysql.CharsetIds[charset]
// 		if !ok {
// 			return fmt.Errorf("invalid charset %s", charset)
// 		}
// 	} else {
// 		collate := sqlparser.String(ci)
// 		collate = strings.Trim(collate, "'`\"")
// 		collate = strings.ToLower(collate)
// 		cid, ok = mysql.CollationNames[collate]
// 		if !ok {
// 			return fmt.Errorf("invalid collation %s", collate)
// 		}
// 	}
// 	c.charset = charset
// 	c.collation = cid

// 	return c.writeOK(nil)
// }

// func (c *Conn) handleSimpleSelect(stmt *sqlparser.SelectExpr) error {
// 	nonStarExpr, _ := stmt.SelectExprs[0].(*sqlparser.NonStarExpr)
// 	var name string = mysql.Strings(nonStarExpr.As)
// 	if name == "" {
// 		name = "last_insert_id()"
// 	}
// 	var column = 1
// 	var rows [][]string
// 	var names []string = []string{
// 		name,
// 	}

// 	var t = fmt.Sprintf("%d", c.lastInsertId)
// 	rows = append(rows, []string{t})

// 	r := new(mysql.Resultset)

// 	var values [][]interface{} = make([][]interface{}, len(rows))
// 	for i := range rows {
// 		values[i] = make([]interface{}, column)
// 		for j := range rows[i] {
// 			values[i][j] = rows[i][j]
// 		}
// 	}

// 	r, _ = c.buildResultset(nil, names, values)
// 	return c.writeResultset(c.status, r)
// }

func (c *Conn) handleBegin() error {
	co := c.down
	if err := co.Begin(); err != nil {
		return err
	}
	c.status |= mysql.SERVER_STATUS_IN_TRANS
	return c.writeOK(nil)
}
func (c *Conn) handleCommit() error {
	co := c.down
	if err := co.Commit(); err != nil {
		return err
	}
	c.status |= mysql.SERVER_STATUS_IN_TRANS
	return c.writeOK(nil)
}
func (c *Conn) handleRollback() error {
	co := c.down
	if err := co.Rollback(); err != nil {
		return err
	}
	c.status |= mysql.SERVER_STATUS_IN_TRANS
	return c.writeOK(nil)
}
