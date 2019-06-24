package server

import (
	"sync"
	"testing"
	"time"

	"github.com/muidea/magicProxy/backend"
	"github.com/muidea/magicProxy/config"
)

var testServerOnce sync.Once
var testServer *Server
var testDBOnce sync.Once
var testDB *backend.DB

var testConfigData = []byte(`
addr : 127.0.0.1:9696
user : root
password : root

node :
    name : node 
    idle_conns : 16
    user: magicbatis
    password: magicbatis
    address : 127.0.0.1:3306
`)

func newTestServer(t *testing.T) *Server {
	f := func() {
		cfg, err := config.ParseConfigData(testConfigData)
		if err != nil {
			t.Fatal(err.Error())
		}

		testServer, err = NewServer(cfg)
		if err != nil {
			t.Fatal(err)
		}

		go testServer.Run()

		time.Sleep(1 * time.Second)
	}

	testServerOnce.Do(f)

	return testServer
}

func newTestDB(t *testing.T) *backend.DB {
	newTestServer(t)

	f := func() {
		testDB, _ = backend.Open("127.0.0.1:3306", "root", "rootkit", "testDB", 100)
	}

	testDBOnce.Do(f)
	return testDB
}

func newTestDBConn(t *testing.T) *backend.BackendConn {
	db := newTestDB(t)

	c, err := db.GetConn()

	if err != nil {
		t.Fatal(err)
	}

	return c
}

func TestServer(t *testing.T) {
	newTestServer(t)
}
