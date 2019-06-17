.Phony:

code: 
	@goyacc -o ./sqlparser/sql.go ./sqlparser/sql.y
	@gofmt -w ./sqlparser/sql.go
	@echo 'code done'

clean:
	@rm -rf sqlparser/sql.go
	@rm -rf y.output coverage.html cover.out
	@echo 'clean done'

test: .Phony
	@go test -v -coverprofile=cover.out test/parser_test.go
	@go tool cover  -html=cover.out -o coverage.html 

ptest: clean code
	@go test -timeout 30s supos.ai/data-lake/common/sql-parser/sqlparser
	@go test -timeout 30s supos.ai/data-lake/common/sql-parser/test -run ^\(TestParser\)$$
