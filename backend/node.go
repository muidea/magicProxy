package backend

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/muidea/magicProxy/config"
	"github.com/muidea/magicProxy/core/errors"
	"github.com/muidea/magicProxy/core/golog"
)

const (
	Master      = "master"
	Slave       = "slave"
	SlaveSplit  = ","
	WeightSplit = "@"
)

type Node struct {
	Cfg config.NodeConfig

	sync.RWMutex
	Master *DB

	LastSlaveIndex int
	RoundRobinQ    []int
	SlaveWeights   []int

	DownAfterNoAlive time.Duration

	Online bool
}

func (n *Node) CheckNode() {
	//to do
	//1 check connection alive
	for n.Online {
		n.checkMaster()
		time.Sleep(16 * time.Second)
	}
}

func (n *Node) String() string {
	return n.Cfg.Name
}

func (n *Node) GetMasterConn() (*BackendConn, error) {
	db := n.Master
	if db == nil {
		return nil, errors.ErrNoMasterConn
	}
	if atomic.LoadInt32(&(db.state)) == Down {
		return nil, errors.ErrMasterDown
	}

	return db.GetConn()
}

func (n *Node) checkMaster() {
	db := n.Master
	if db == nil {
		golog.Error("Node", "checkMaster", "Master is no alive", 0)
		return
	}

	if err := db.Ping(); err != nil {
		golog.Error("Node", "checkMaster", "Ping", 0, "db.Addr", db.Addr(), "error", err.Error())
	} else {
		if atomic.LoadInt32(&(db.state)) == Down {
			golog.Info("Node", "checkMaster", "Master up", 0, "db.Addr", db.Addr())
			err := n.UpMaster(db.addr)
			if err != nil {
				golog.Error("Node", "checkMaster", "UpMaster", 0, "db.Addr", db.Addr(), "error", err.Error())
				return
			}
		}
		db.SetLastPing()
		if atomic.LoadInt32(&(db.state)) != ManualDown {
			atomic.StoreInt32(&(db.state), Up)
		}
		return
	}

	if int64(n.DownAfterNoAlive) > 0 && time.Now().Unix()-db.GetLastPing() > int64(n.DownAfterNoAlive/time.Second) {
		golog.Info("Node", "checkMaster", "Master down", 0,
			"db.Addr", db.Addr(),
			"Master_down_time", int64(n.DownAfterNoAlive/time.Second))
		n.DownMaster(db.addr, Down)
	}
}

func (n *Node) OpenDB(addr string) (*DB, error) {
	db, err := Open(addr, n.Cfg.User, n.Cfg.Password, "", n.Cfg.MaxConnNum)
	return db, err
}

func (n *Node) UpDB(addr string) (*DB, error) {
	db, err := n.OpenDB(addr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		atomic.StoreInt32(&(db.state), Down)
		return nil, err
	}
	atomic.StoreInt32(&(db.state), Up)
	return db, nil
}

func (n *Node) UpMaster(addr string) error {
	db, err := n.UpDB(addr)
	if err != nil {
		golog.Error("Node", "UpMaster", err.Error(), 0)
		return err
	}
	n.Master = db
	return err
}

func (n *Node) DownMaster(addr string, state int32) error {
	db := n.Master
	if db == nil || db.addr != addr {
		return errors.ErrNoMasterDB
	}

	db.Close()
	atomic.StoreInt32(&(db.state), state)
	return nil
}

func (n *Node) ParseMaster(masterStr string) error {
	var err error
	if len(masterStr) == 0 {
		return errors.ErrNoMasterDB
	}

	n.Master, err = n.OpenDB(masterStr)
	return err
}
