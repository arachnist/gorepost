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
	network        string
	input          chan Message
	output         chan Message
	reader         *bufio.Reader
	writer         *bufio.Writer
	dispatcher     func(chan struct{}, chan Message, chan Message)
	conn           net.Conn
	reconnect      chan struct{}
	Quit           chan struct{}
	quitsend       chan struct{}
	quitrecv       chan struct{}
	quitdispatcher chan struct{}
	l              sync.Mutex
}

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

func (c *Connection) Receiver() {
	log.Println(c.network, "spawned Receiver")
	for {
		c.conn.SetReadDeadline(time.Now().Add(time.Second * 600))
		raw, err := c.reader.ReadString(delim)
		var src, tgt string
		if err != nil {
			log.Println(c.network, "error reading message", err.Error())
			log.Println(c.network, "closing Receiver")
			c.Quit <- struct{}{}
			log.Println(c.network, "sent quit message from Receiver")
			return
		}
		msg, err := ParseMessage(raw)
		if err != nil {
			log.Println(c.network, "error decoding message", err.Error())
			log.Println(c.network, "closing Receiver")
			c.Quit <- struct{}{}
			log.Println(c.network, "sent quit message from Receiver")
			return
		} else {
			log.Println(c.network, "<--", msg.String())
		}
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
		select {
		case c.output <- *msg:
		case <-c.quitrecv:
			log.Println(c.network, "closing Receiver")
			return
		}
	}
}

func (c *Connection) Cleaner() {
	log.Println(c.network, "spawned Cleaner")
	for {
		<-c.Quit
		log.Println(c.network, "received quit message")
		c.l.Lock()
		log.Println(c.network, "cleaning up!")
		c.quitsend <- struct{}{}
		c.quitrecv <- struct{}{}
		c.quitdispatcher <- struct{}{}
		c.reconnect <- struct{}{}
		c.conn.Close()
		log.Println(c.network, "closing Cleaner")
		c.l.Unlock()
	}
}

func (c *Connection) Keeper() {
	log.Println(c.network, "spawned Keeper")
	context := make(map[string]string)
	context["Network"] = c.network
	for {
		<-c.reconnect
		c.l.Lock()
		if c.input != nil {
			close(c.input)
			close(c.output)
			close(c.quitsend)
			close(c.quitrecv)
			close(c.quitdispatcher)
		}
		c.input = make(chan Message, 1)
		c.output = make(chan Message, 1)
		c.quitsend = make(chan struct{}, 1)
		c.quitrecv = make(chan struct{}, 1)
		c.quitdispatcher = make(chan struct{}, 1)
		servers := C.Lookup(context, "Servers").([]interface{})

		server := servers[rand.Intn(len(servers))].(string)
		log.Println(c.network, "connecting to", server)
		err := c.Dial(server)
		c.l.Unlock()
		if err == nil {
			go c.Sender()
			go c.Receiver()
			go c.dispatcher(c.quitdispatcher, c.input, c.output)

			log.Println(c.network, "Initializing IRC connection")
			c.input <- Message{
				Command:  "NICK",
				Trailing: C.Lookup(context, "Nick").(string),
			}
			c.input <- Message{
				Command:  "USER",
				Params:   []string{C.Lookup(context, "User").(string), "0", "*"},
				Trailing: C.Lookup(context, "RealName").(string),
			}
		} else {
			log.Println(c.network, "connection error", err.Error())
			c.reconnect <- struct{}{}
		}
	}
}

func (c *Connection) Setup(dispatcher func(chan struct{}, chan Message, chan Message), network string) {
	rand.Seed(time.Now().UnixNano())

	c.reconnect = make(chan struct{}, 1)
	c.Quit = make(chan struct{}, 1)
	c.network = network
	c.dispatcher = dispatcher

	c.reconnect <- struct{}{}
	go c.Keeper()
	go c.Cleaner()
	return
}

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
