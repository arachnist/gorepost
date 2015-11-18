// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package irc

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	cfg "github.com/arachnist/gorepost/config"
)

const delim byte = '\n'
const endline string = "\r\n"

// Connection struct. Contains basic information about this connection, as well
// as input and quit channels.
type Connection struct {
	network          string
	input            chan Message
	reader           *bufio.Reader
	writer           *bufio.Writer
	dispatcher       func(chan Message, Message)
	conn             net.Conn
	reconnect        chan struct{}
	reconnectCleanup chan struct{}
	Quit             chan struct{}
	quitsend         chan struct{}
	quitrecv         chan struct{}
	quitkeeper       chan struct{}
	l                sync.Mutex
}

// Sender sends IRC messages to server and logs their contents.
func (c *Connection) Sender() {
	log.Println(c.network, "spawned Sender")
	for {
		select {
		case msg := <-c.input:
			c.writer.WriteString(msg.String() + endline)
			log.Println(c.network, "-->", msg.String())
			c.writer.Flush()
		case <-c.quitsend:
			log.Println(c.network, "closing Sender")
			return
		}
	}
}

// Receiver receives IRC messages from server, logs their contents, sets message
// context and initializes disconnect procedure on timeout or other errors.
func (c *Connection) Receiver() {
	log.Println(c.network, "spawned Receiver")
	for {
		c.conn.SetReadDeadline(time.Now().Add(time.Second * 600))
		raw, err := c.reader.ReadString(delim)
		var src, tgt string

		if err != nil {
			log.Println(c.network, "error reading message", err.Error())
			log.Println(c.network, "closing Receiver")
			c.reconnectCleanup <- struct{}{}
			log.Println(c.network, "sent reconnect message from Receiver")
			return
		}

		msg, err := ParseMessage(raw)
		if err != nil {
			log.Println(c.network, "error decoding message", err.Error())
			log.Println(c.network, "closing Receiver")
			c.reconnectCleanup <- struct{}{}
			log.Println(c.network, "sent reconnect message from Receiver")
			return
		}

		log.Println(c.network, "<--", msg.String())

		if msg.Params == nil {
			tgt = ""
		} else {
			tgt = msg.Params[0]
		}
		if msg.Prefix == nil {
			src = ""
		} else {
			src = msg.Prefix.Name
		}
		msg.Context = map[string]string{
			"Network": c.network,
			"Source":  src,
			"Target":  tgt,
		}

		c.dispatcher(c.input, *msg)
		select {
		case <-c.quitrecv:
			log.Println(c.network, "closing Receiver")
			return
		default:
		}
	}
}

// Cleaner cleans up coroutines on IRC connection errors and initializes
// reconnection.
func (c *Connection) Cleaner() {
	log.Println(c.network, "spawned Cleaner")
	for {
		select {
		case <-c.Quit:
			log.Println(c.network, "closing connection")
			c.l.Lock()
			defer c.l.Unlock()
			log.Println(c.network, "cleaning up!")
			c.quitsend <- struct{}{}
			c.quitrecv <- struct{}{}
			c.conn.Close()
			log.Println(c.network, "closing Cleaner")
			return
		case <-c.reconnectCleanup:
			log.Println(c.network, "cleaning up before reconnect!")
			c.l.Lock()
			log.Println(c.network, "cleaning up!")
			c.quitsend <- struct{}{}
			c.quitrecv <- struct{}{}
			c.conn.Close()
			log.Println(c.network, "sending reconnect signal!")
			c.l.Unlock()
			c.reconnect <- struct{}{}
		}
	}
}

// Keeper makes sure that IRC connection is alive by reconnecting when
// requested and restarting Sender, Receiver and Dispatcher goroutines.
func (c *Connection) Keeper() {
	log.Println(c.network, "spawned Keeper")
	context := make(map[string]string)
	context["Network"] = c.network
	for {
		select {
		case <-c.quitkeeper:
			if c.input != nil {
				close(c.input)
				close(c.quitsend)
				close(c.quitrecv)
			}
			return
		case <-c.reconnect:
		}

		c.l.Lock()
		if c.input != nil {
			close(c.input)
			close(c.quitsend)
			close(c.quitrecv)
		}
		c.input = make(chan Message, 1)
		c.quitsend = make(chan struct{}, 1)
		c.quitrecv = make(chan struct{}, 1)
		servers := cfg.LookupStringSlice(context, "Servers")

		server := servers[rand.Intn(len(servers))]
		log.Println(c.network, "connecting to", server)
		err := c.Dial(server)
		c.l.Unlock()
		if err == nil {
			go c.Sender()
			go c.Receiver()

			log.Println(c.network, "Initializing IRC connection")
			c.input <- Message{
				Command:  "NICK",
				Trailing: cfg.LookupString(context, "Nick"),
			}
			c.input <- Message{
				Command:  "USER",
				Params:   []string{cfg.LookupString(context, "User"), "0", "*"},
				Trailing: cfg.LookupString(context, "RealName"),
			}
		} else {
			log.Println(c.network, "connection error", err.Error())
			c.reconnect <- struct{}{}
		}
	}
}

// Setup performs initialization tasks.
func (c *Connection) Setup(dispatcher func(chan Message, Message), network string) {
	rand.Seed(time.Now().UnixNano())

	c.reconnect = make(chan struct{}, 1)
	c.reconnectCleanup = make(chan struct{}, 1)
	c.quitkeeper = make(chan struct{}, 1)
	c.Quit = make(chan struct{}, 1)
	c.network = network
	c.dispatcher = dispatcher

	c.reconnect <- struct{}{}
	go c.Keeper()
	go c.Cleaner()
	return
}

// Dial connects to irc server and sets up bufio reader and writer.
func (c *Connection) Dial(server string) error {
	conn, err := net.DialTimeout("tcp", server, time.Second*30)
	if err != nil {
		log.Println(c.network, "Cannot connect to", server, "error:", err.Error())
		return err
	}
	log.Println(c.network, "Connected to", server)
	c.writer = bufio.NewWriter(conn)
	c.reader = bufio.NewReader(conn)
	c.conn = conn

	return nil
}
