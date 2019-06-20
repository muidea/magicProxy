package backend

import (
	"testing"

	"github.com/muidea/magicProxy/config"
)

func TestParse(t *testing.T) {
	node := new(Node)
	nodeConfig := config.NodeConfig{
		Name:             "node1",
		DownAfterNoAlive: 100,
		User:             "root",
		Password:         "rootkit",
		Master:           "127.0.0.1:3306",
	}
	node.Cfg = nodeConfig
	err := node.ParseMaster(nodeConfig.Master)
	if err != nil {
		t.Fatal(err.Error())
	}
	if node.Master.addr != "127.0.0.1:3306" {
		t.Fatal(node.Master)
	}
	t.Logf("%v\n", node.RoundRobinQ)
	t.Logf("%v\n", node.SlaveWeights)
	t.Logf("%v\n", node.Master)
}
