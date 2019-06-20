package backend

import (
	"testing"
)

func TestDB_Ping(t *testing.T) {
	db, err := Open("127.0.0.1:3308", "magicbatis", "magicbatis", "testDB", 1)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.Ping()
}
