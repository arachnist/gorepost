package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	config, err := config.ReadConfig(os.Args[1])
	var exit chan struct{}
	if err != nil {
		fmt.Println("Error reading configuration from", os.Args[1], "error:", err.Error())
		os.Exit(1)
	}

	logfile, err := os.OpenFile(config.Logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening", config.Logpath, "for writing, error:", err.Error())
		os.Exit(1)
	}
	log.SetOutput(logfile)

	connections := make([]irc.Connection, len(config.Networks))
	for i, conn := range connections {
		network := config.Networks[i]
		server := config.Servers[network][rand.Intn(len(config.Servers[network]))]
		conn.Network = config.Networks[i]
		conn.Dial(server)
	}
	<-exit
}
