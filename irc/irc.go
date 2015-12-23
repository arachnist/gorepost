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
	"unicode/utf8"

	"github.com/arachnist/dyncfg"
)

const delim byte = '\n'
const endline string = "\r\n"

// Connection struct. Contains basic information about this connection, and quit
// channels.
type Connection struct {
	network          string
	reader           *bufio.Reader
	writer           *bufio.Writer
	dispatcher       func(func(Message), Message)
	conn             net.Conn
	reconnect        chan struct{}
	reconnectCleanup chan struct{}
	Quit             chan struct{}
	quitrecv         chan struct{}
	quitkeeper       chan struct{}
	l                sync.Mutex
	cfg              *dyncfg.Dyncfg
}

// Sender sends IRC messages to server and logs their contents.
func (c *Connection) Sender(msg Message) {
	c.l.Lock()
	defer c.l.Unlock()
	if msg.WireLen() > maxLength {
		currLen := 0
		for i, ch := range msg.String() {
			currLen += utf8.RuneLen(ch)
			if currLen > maxLength {
				c.writer.WriteString(msg.String()[:i] + endline)
				log.Println(c.network, "-->", msg.String())
				c.writer.Flush()
				// eh, it is a bit naive to assume that we won't explode again…
				if msg.Command == "PRIVMSG" { // we don't care otherwise
					newMsg := msg
					newMsg.Trailing = "Message truncated"
					go c.Sender(newMsg)
				}
				return
			}
		}
	} else {
		c.writer.WriteString(msg.String() + endline)
		log.Println(c.network, "-->", msg.String())
		c.writer.Flush()
	}
}

// Receiver receives IRC messages from server, logs their contents, sets message
// context and initializes disconnect procedure on timeout or other errors.
func (c *Connection) Receiver() {
	log.Println(c.network, "spawned Receiver")
	for {
		c.conn.SetDeadline(time.Now().Add(time.Second * 600))
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

		c.dispatcher(c.Sender, *msg)
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
			c.quitrecv <- struct{}{}
			// there's a slight chance to hit this if quit request is received
			// before irc connection is established, possibly between reconnects
			if c.conn != nil {
				c.conn.Close()
			}
			log.Println(c.network, "closing Cleaner")
			return
		case <-c.reconnectCleanup:
			log.Println(c.network, "cleaning up before reconnect!")
			c.l.Lock()
			log.Println(c.network, "cleaning up!")
			c.quitrecv <- struct{}{}
			c.conn.Close()
			log.Println(c.network, "sending reconnect signal!")
			c.l.Unlock()
			c.reconnect <- struct{}{}
		}
	}
}

// Keeper makes sure that IRC connection is alive by reconnecting when
// requested and restarting Receiver goroutine.
func (c *Connection) Keeper() {
	log.Println(c.network, "spawned Keeper")
	context := make(map[string]string)
	context["Network"] = c.network
	for {
		select {
		case <-c.quitkeeper:
			if c.quitrecv != nil {
				close(c.quitrecv)
			}
			return
		case <-c.reconnect:
		}

		c.l.Lock()
		if c.quitrecv != nil {
			close(c.quitrecv)
		}
		c.quitrecv = make(chan struct{}, 1)
		servers := c.cfg.LookupStringSlice(context, "Servers")

		server := servers[rand.Intn(len(servers))]
		log.Println(c.network, "connecting to", server)
		err := c.Dial(server)
		c.l.Unlock()
		if err == nil {
			go c.Receiver()

			log.Println(c.network, "Initializing IRC connection")
			c.Sender(Message{
				Command:  "NICK",
				Trailing: c.cfg.LookupString(context, "Nick"),
			})
			c.Sender(Message{
				Command:  "USER",
				Params:   []string{c.cfg.LookupString(context, "User"), "0", "*"},
				Trailing: c.cfg.LookupString(context, "RealName"),
			})
		} else {
			log.Println(c.network, "connection error", err.Error())
			time.Sleep(time.Second * 3)
			c.reconnect <- struct{}{}
		}
	}
}

// Setup performs initialization tasks.
func (c *Connection) Setup(dispatcher func(func(Message), Message), network string, config *dyncfg.Dyncfg) {
	rand.Seed(time.Now().UnixNano())

	c.reconnect = make(chan struct{}, 1)
	c.reconnectCleanup = make(chan struct{}, 1)
	c.quitkeeper = make(chan struct{}, 1)
	c.Quit = make(chan struct{}, 1)
	c.network = network
	c.dispatcher = dispatcher
	c.cfg = config

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
