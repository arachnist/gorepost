package main

import (
	"log"
	"os"

	"github.com/arachnist/gorepost/bot"
	"github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func main() {
	var exit chan struct{}

	if len(os.Args) < 2 {
		log.Fatalln("Usage:", os.Args[0], "<config-file.json>")
	}
	config, err := config.ReadConfig(os.Args[1])
	if err != nil {
		log.Fatalln("Error reading configuration from", os.Args[1], "error:", err.Error())
	}

	logfile, err := os.OpenFile(config.Logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Error opening", config.Logpath, "for writing, error:", err.Error())
	}
	log.SetOutput(logfile)

	connections := make([]irc.Connection, len(config.Networks))
	for i, _ := range connections {
		network := config.Networks[i]
		connections[i].Setup(bot.Dispatcher, network, config.Servers[network], config.Nick, config.User, config.RealName)

		bot.AddCallback("001", "channel join", func(output chan irc.Message, msg irc.Message) {
			for _, channel := range config.Channels[network] {
				log.Println(network, "joining channel", channel)
				output <- irc.Message{
					Command: "JOIN",
					Params:  []string{channel},
				}
			}
		})
	}
	<-exit
}
