package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/muidea/magicProxy/core"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	//进程退出信号传入各模块，以保证进程安全退出
	//从此下面的程序退出，都要开始做退出处理
	closeChan := make(chan bool)
	var wg sync.WaitGroup

	//start service!
	wg.Add(1)
	go func() {
		core.Run(closeChan, &wg)
		log.Printf("shutdown proxy")
		defer wg.Done()
	}()

	wg.Wait()
	<-closeChan
	log.Printf("shutdown now")
	//预留一点时间给log模块
	time.Sleep(time.Second)
}
