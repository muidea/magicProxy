package backend

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/muidea/magicProxy/config"
	"github.com/muidea/magicProxy/core/errors"
	"github.com/muidea/magicProxy/core/golog"
)

// Node backend node
type Node struct {
	Cfg config.NodeConfig

	sync.RWMutex
	Database *DB

	Online bool
}

// CheckNode check node status
func (n *Node) CheckNode() {
	//to do
	//1 check connection alive
	for n.Online {
		n.checkStatus()
		time.Sleep(16 * time.Second)
	}
}

func (n *Node) String() string {
	return n.Cfg.Name
}

// GetConn get connection
func (n *Node) GetConn() (*BackendConn, error) {
	db := n.Database
	if db == nil {
		err := n.parseDB(n.Cfg.Address)
		if err != nil {
			return nil, err
		}

		db = n.Database
	}
	if atomic.LoadInt32(&(db.state)) == Down {
		return nil, errors.ErrMasterDown
	}

	return db.GetConn()
}

func (n *Node) checkStatus() {
	db := n.Database
	if db == nil {
		golog.Error("Node", "checkStatus", "Database is no alive", 0)
		return
	}

	if err := db.Ping(); err != nil {
		golog.Error("Node", "checkStatus", "Ping", 0, "db.Addr", db.Addr(), "error", err.Error())
	}
}

func (n *Node) openDB(addr string) (*DB, error) {
	db, err := Open(addr, n.Cfg.User, n.Cfg.Password, "", n.Cfg.MaxConnNum)
	return db, err
}

func (n *Node) parseDB(dbAddress string) error {
	var err error
	if len(dbAddress) == 0 {
		return errors.ErrNoMasterDB
	}

	n.Database, err = n.openDB(dbAddress)
	return err
}
