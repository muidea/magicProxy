package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/muidea/magicProxy/config"
	"github.com/muidea/magicProxy/core/golog"
	"github.com/muidea/magicProxy/server"
)

var configFile *string = flag.String("config", "/etc/ks.yaml", "magicProxy config file")
var version *bool = flag.Bool("v", false, "the version of magicProxy")

var (
	BuildDate    string
	BuildVersion string
)

func main() {
	fmt.Print("magicProxy v1.0")
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	fmt.Printf("Git commit:%s\n", BuildVersion)
	fmt.Printf("Build time:%s\n", BuildDate)
	if *version {
		return
	}
	if len(*configFile) == 0 {
		fmt.Println("must use a config file")
		return
	}

	cfg, err := config.ParseConfigFile(*configFile)
	if err != nil {
		fmt.Printf("parse config file error:%v\n", err.Error())
		return
	}

	var svr *server.Server
	svr, err = server.NewServer(cfg)
	if err != nil {
		golog.Error("main", "main", err.Error(), 0)
		golog.GlobalSysLogger.Close()
		golog.GlobalSqlLogger.Close()
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGPIPE,
	)

	go func() {
		for {
			sig := <-sc
			if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
				golog.Info("main", "main", "Got signal", 0, "signal", sig)
				golog.GlobalSysLogger.Close()
				golog.GlobalSqlLogger.Close()
				svr.Close()
			} else if sig == syscall.SIGPIPE {
				golog.Info("main", "main", "Ignore broken pipe signal", 0)
			}
		}
	}()
	svr.Run()
}
