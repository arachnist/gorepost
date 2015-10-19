package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/arachnist/gorepost/config"
	"github.com/sorcix/irc"
)

type Connection struct {
	Network  string
	Input    chan irc.Message
	Output   chan irc.Message
	IRCConn  *irc.Conn
	QuitSend chan struct{}
	QuitRecv chan struct{}
}

func (c *Connection) Sender() {
	for {
		select {
		case msg := <-c.Input:
			c.IRCConn.Encode(&msg)
		case <-c.QuitSend:
			log.Println(c.Network, "closing Sender")
			close(c.Input)
			close(c.QuitSend)
			return
		}
	}
}

func (c *Connection) Receiver() {
	for {
		msg, err := c.IRCConn.Decode()
		if err != nil {
			log.Println(c.Network, "error decoding message", err.Error())
		}
		select {
		case c.Output <- *msg:
		case <-c.QuitRecv:
			log.Println(c.Network, "closing receiver")
			close(c.Output)
			close(c.QuitRecv)
			return
		}
	}
}

func SetupConn(network string, config config.Config, connection *Connection) error {
	rand.Seed(time.Now().UnixNano())
	server := config.Servers[network][rand.Intn(len(config.Servers[network]))]

	conn, err := irc.Dial(server)
	if err != nil {
		log.Println("Cannot connect to", network, "server:", server, "error:", err.Error())
		return err
	}
	connection.IRCConn = conn

	go connection.Sender()
	go connection.Receiver()

	// Initial commands sent to IRC server
	connection.Input <- irc.Message{
		Command:  "NICK",
		Trailing: config.Nick,
	}
	connection.Input <- irc.Message{
		Command:  "USER",
		Params:   []string{config.Nick, "3", "*"},
		Trailing: config.Nick,
	}

	return nil
}

func ConnectionKeeper(connection *Connection) {
	for {

	}
}

func main() {
	config, err := config.ReadConfig(os.Args[1])
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

}
