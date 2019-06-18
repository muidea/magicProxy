package core

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	err error
	cfg *Config
)

func init() {
	cfg, err = ParseArgCmd()
	if err != nil {
		panic(err)
	}
}

// Run func
// sql-endpoint entry function
func Run(closeNotify chan bool, joinGroup *sync.WaitGroup) {
	proxyServer := NewServer(cfg, joinGroup)
	go proxyServer.Run()

	osSignalChan := make(chan os.Signal)
	//等待Kill信号
	signal.Notify(osSignalChan, os.Interrupt, syscall.SIGTERM)
	log.Printf("exit with signal:%v", <-osSignalChan)
	proxyServer.Close()
	// 退出完成通知
	close(closeNotify)
}
