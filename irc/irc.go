package irc

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
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
	Reader         *bufio.Reader
	Writer         *bufio.Writer
	conn           net.Conn
	Reconnect      chan struct{}
	Quit           chan struct{}
	QuitSend       chan struct{}
	QuitRecv       chan struct{}
	QuitDispatcher chan struct{}
	m              sync.Mutex
}

func (c *Connection) Sender() {
	log.Println(c.Network, "Spawned sender loop")
	for {
		select {
		case msg := <-c.Input:
			c.Writer.WriteString(msg.String() + endline)
			log.Println(c.Network, "-->", msg.String())
			c.Writer.Flush()
		case <-c.QuitSend:
			log.Println(c.Network, "closing Sender")
			return
		}
	}
}

func (c *Connection) Receiver() {
	log.Println(c.Network, "Spawned receiver loop")
	for {
		raw, err := c.Reader.ReadString(delim)
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
		case <-c.QuitRecv:
			log.Println(c.Network, "closing Receiver")
			return
		}
	}
}

func (c *Connection) Dispatcher() {
	log.Println(c.Network, "spawned Dispatcher")
	for {
		// just sink everything for now
		select {
		case <-c.Output:
		case <-c.QuitDispatcher:
			log.Println(c.Network, "closing Dispatcher")
			return
		}
	}
}

func (c *Connection) Cleaner() {
	log.Println(c.Network, "spawned Cleaner")
	for {
		<-c.Quit
		log.Println(c.Network, "Received quit message")
		c.m.Lock()
		log.Println(c.Network, "cleaning up!")
		c.QuitSend <- struct{}{}
		c.QuitRecv <- struct{}{}
		c.QuitDispatcher <- struct{}{}
		c.Reconnect <- struct{}{}
		c.conn.Close()
		log.Println(c.Network, "closing Cleaner")
		c.m.Unlock()
	}
}

func (c *Connection) Keeper(servers []string) {
	log.Println(c.Network, "spawned Keeper")
	for {
		<-c.Reconnect
		c.m.Lock()
		c.Input = make(chan Message, 1)
		c.Output = make(chan Message, 1)
		c.QuitSend = make(chan struct{}, 1)
		c.QuitRecv = make(chan struct{}, 1)
		c.QuitDispatcher = make(chan struct{}, 1)
		server := servers[rand.Intn(len(servers))]
		log.Println(c.Network, "connecting to", server)
		c.Dial(server)
		c.m.Unlock()

		go c.Sender()
		go c.Receiver()
		go c.Dispatcher()

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

	}
}

func (c *Connection) Setup(network string, servers []string, nick string, user string, realname string) {
	rand.Seed(time.Now().UnixNano())

	c.Reconnect = make(chan struct{}, 1)
	c.Quit = make(chan struct{}, 1)
	c.Nick = nick
	c.User = user
	c.RealName = realname

	c.Reconnect <- struct{}{}
	c.Network = network
	go c.Keeper(servers)
	go c.Cleaner()
	return
}

func (c *Connection) Dial(server string) error {

	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Println(c.Network, "Cannot connect to", server, "error:", err.Error())
		return err
	}
	log.Println(c.Network, "Connected to", server)
	c.Writer = bufio.NewWriter(conn)
	c.Reader = bufio.NewReader(conn)
	c.conn = conn

	return nil
}
