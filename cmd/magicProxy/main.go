// Copyright 2016 The magicProxy Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

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

const banner string = `
    __   _                  __                   __
   / /__(_)___  ____ ______/ /_  ____ __________/ /
  / //_/ / __ \/ __ \/ ___/ __ \ / __\/ ___/ __  /
 / ,< / / / / / /_/ (__  ) / / / /_/ / /  / /_/ /
/_/|_/_/_/ /_/\__, /____/_/ /_/\__,_/_/   \__,_/
             /____/
`

func main() {
	fmt.Print(banner)
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
