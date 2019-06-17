package test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/muidea/magicProxy/common/sql-parser/sqlparser"
)

type Kase interface {
	Text() string
	Care() bool
}
type Case struct {
	raw  string
	care bool
}

func (c Case) Care() bool   { return c.care }
func (c Case) Text() string { return c.raw }

func LoadSQL(more chan bool) (out chan Kase) {
	out = make(chan Kase, 10)
	go func() {
		sample, err := os.Open("sample.sh")
		if err != nil {
			panic(err)
		}
		bio := bufio.NewScanner(sample)
		var sql string
		for bio.Scan() {
			sql = bio.Text()
			switch {
			case strings.HasPrefix(sql, "#"):
				continue
			case len(strings.TrimSpace(sql)) == 0:
				continue
			case strings.HasPrefix(sql, "!"):
				out <- Case{strings.TrimPrefix(sql, "!"), true}
			default:
				out <- Case{sql, false}
			}
			<-more
		}
		close(out)
	}()
	return
}
func TestParser(t *testing.T) {
	var p bool
	var stmt sqlparser.Statement
	var err error
	more := make(chan bool)
	for sql := range LoadSQL(more) {
		p = false
		stmt, err = sqlparser.Parse(sql.Text())
		if err != nil {
			t.Error(err)
			p = true
		}
		if sql.Care() || p {
			t.Logf("%s\n--->%s", sql.Text(), sqlparser.String(stmt))
		}
		more <- true
	}
}

func TestAccessTablename(t *testing.T) {
	var stmt sqlparser.Statement
	var err error
	stmt, err = sqlparser.Parse("Select * from item as p")
	sel := stmt.(*sqlparser.Select)
	tbl := sel.From[0]
	at := tbl.(*sqlparser.AliasedTableExpr)
	atn := at.Expr.(sqlparser.TableName)
	t.Error(err)
	t.Logf("%s\n--->%s", atn.Name.String(), sqlparser.String(stmt))

}
