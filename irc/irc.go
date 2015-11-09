package irc

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	. "github.com/arachnist/gorepost/config"
)

const delim byte = '\n'
const endline string = "\r\n"

type Connection struct {
	Network        string
	Nick           string
	User           string
	RealName       string
	Input          chan Message
	Output         chan Message
	reader         *bufio.Reader
	writer         *bufio.Writer
	dispatcher     func(chan struct{}, Context, chan Message, chan Message)
	conn           net.Conn
	reconnect      chan struct{}
	Quit           chan struct{}
	quitsend       chan struct{}
	quitrecv       chan struct{}
	quitdispatcher chan struct{}
	l              sync.Mutex
}

func (c *Connection) Sender() {
	log.Println(c.Network, "spawned Sender")
	for {
		select {
		case msg := <-c.Input:
			c.writer.WriteString(msg.String() + endline)
			log.Println(c.Network, "-->", msg.String())
			c.writer.Flush()
		case <-c.quitsend:
			log.Println(c.Network, "closing Sender")
			return
		}
	}
}

func (c *Connection) Receiver() {
	log.Println(c.Network, "spawned Receiver")
	for {
		c.conn.SetReadDeadline(time.Now().Add(time.Second * 600))
		raw, err := c.reader.ReadString(delim)
		if err != nil {
			log.Println(c.Network, "error reading message", err.Error())
			log.Println(c.Network, "closing Receiver")
			c.Quit <- struct{}{}
			log.Println(c.Network, "sent quit message from Receiver")
			return
		}
		msg, err := ParseMessage(raw)
		if err != nil {
			log.Println(c.Network, "error decoding message", err.Error())
			log.Println(c.Network, "closing Receiver")
			c.Quit <- struct{}{}
			log.Println(c.Network, "sent quit message from Receiver")
			return
		} else {
			log.Println(c.Network, "<--", msg.String())
		}
		select {
		case c.Output <- *msg:
		case <-c.quitrecv:
			log.Println(c.Network, "closing Receiver")
			return
		}
	}
}

func (c *Connection) Cleaner() {
	log.Println(c.Network, "spawned Cleaner")
	for {
		<-c.Quit
		log.Println(c.Network, "received quit message")
		c.l.Lock()
		log.Println(c.Network, "cleaning up!")
		c.quitsend <- struct{}{}
		c.quitrecv <- struct{}{}
		c.quitdispatcher <- struct{}{}
		c.reconnect <- struct{}{}
		c.conn.Close()
		log.Println(c.Network, "closing Cleaner")
		c.l.Unlock()
	}
}

func (c *Connection) Keeper(servers []string) {
	log.Println(c.Network, "spawned Keeper")
	for {
		<-c.reconnect
		c.l.Lock()
		if c.Input != nil {
			close(c.Input)
			close(c.Output)
			close(c.quitsend)
			close(c.quitrecv)
			close(c.quitdispatcher)
		}
		c.Input = make(chan Message, 1)
		c.Output = make(chan Message, 1)
		c.quitsend = make(chan struct{}, 1)
		c.quitrecv = make(chan struct{}, 1)
		c.quitdispatcher = make(chan struct{}, 1)
		server := servers[rand.Intn(len(servers))]
		log.Println(c.Network, "connecting to", server)
		err := c.Dial(server)
		c.l.Unlock()
		if err == nil {
			go c.Sender()
			go c.Receiver()
			go c.dispatcher(c.quitdispatcher, c.Input, c.Output)

			log.Println(c.Network, "Initializing IRC connection")
			c.Input <- Message{
				Command:  "NICK",
				Trailing: c.Nick,
			}
			c.Input <- Message{
				Command:  "USER",
				Params:   []string{c.User, "0", "*"},
				Trailing: c.RealName,
			}
		} else {
			log.Println(c.Network, "connection error", err.Error())
			c.reconnect <- struct{}{}
		}
	}
}

func (c *Connection) Setup(dispatcher func(chan struct{}, Context, chan Message, chan Message), network string, servers []string, nick string, user string, realname string) {
	rand.Seed(time.Now().UnixNano())

	c.reconnect = make(chan struct{}, 1)
	c.Quit = make(chan struct{}, 1)
	c.Nick = nick
	c.User = user
	c.RealName = realname
	c.Network = network
	c.dispatcher = dispatcher

	c.reconnect <- struct{}{}
	go c.Keeper(servers)
	go c.Cleaner()
	return
}

func (c *Connection) Dial(server string) error {
	conn, err := net.DialTimeout("tcp", server, time.Second*30)
	if err != nil {
		log.Println(c.Network, "Cannot connect to", server, "error:", err.Error())
		return err
	}
	log.Println(c.Network, "Connected to", server)
	c.writer = bufio.NewWriter(conn)
	c.reader = bufio.NewReader(conn)
	c.conn = conn

	return nil
}
