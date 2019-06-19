COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "GitNotFound")
BUILD_DATE=$(shell date '+%Y-%m-%d %H:%M:%S')

all: build

build: magicProxy
goyacc:
	go build -o ./bin/goyacc ./vendor/golang.org/x/tools/cmd/goyacc
magicProxy: goyacc
	./bin/goyacc -o ./sqlparser/sql.go ./sqlparser/sql.y
	gofmt -w ./sqlparser/sql.go
	go build -ldflags "-X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\"" -o ./bin/magicProxy ./cmd/magicProxy
clean:
	@rm -rf bin
	@rm -f ./sqlparser/y.output ./sqlparser/sql.go

test:
	go test ./go/... -race
