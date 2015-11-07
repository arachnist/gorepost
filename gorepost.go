package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func main() {
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
	for i, _ := range connections {
		network := config.Networks[i]
		connections[i].Setup(network, config.Servers[network], config.Nick, config.User, config.RealName)
	}
	<-exit
}
