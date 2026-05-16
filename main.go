package main

import (
	"flag"
	"log"

	"github.com/arun-builds/redis-scratch/config"
	"github.com/arun-builds/redis-scratch/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the redis server")
	flag.IntVar(&config.Port, "port", 6767, "port for redis server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Printf("starting redis")
	server.RunSyncTCPServer()
}
