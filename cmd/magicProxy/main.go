package main

import (
	"flag"
	"log"

	"github.com/muidea/magicProxy/core"
)

var bindPort = "8081"

func main() {
	flag.StringVar(&bindPort, "ListenPort", bindPort, "magicProxy listen port.")
	flag.Parse()

	log.Println("MagicProxy V1.0")
	cfg := &core.Config{BindAddr: ":3308", BackendAddr: "127.0.0.1:3306"}
	server := core.NewServer(cfg)
	server.Run()
}
