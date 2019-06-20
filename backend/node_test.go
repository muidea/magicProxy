package backend

import (
	"testing"

	"github.com/muidea/magicProxy/config"
)

func TestParse(t *testing.T) {
	node := new(Node)
	nodeConfig := config.NodeConfig{
		Name:     "node1",
		User:     "magicbatis",
		Password: "magicbatis",
		Address:  "127.0.0.1:3308",
	}
	node.Cfg = nodeConfig
	err := node.parseDB(nodeConfig.Address)
	if err != nil {
		t.Fatal(err.Error())
	}
	if node.Database.addr != "127.0.0.1:3308" {
		t.Fatal(node.Database)
	}
	t.Logf("%v\n", node.Database)
}
