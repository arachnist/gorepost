package irc

import (
	"bufio"
	"log"
	"net"
)

const delim byte = '\n'
const endline string = "\r\n"

type Connection struct {
	Network  string
	Input    chan Message
	Output   chan Message
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	QuitSend chan struct{}
	QuitRecv chan struct{}
}

func (c Connection) Sender() {
	log.Println(c.Network, "Spawned sender loop")
	for {
		select {
		case msg := <-c.Input:
			c.Writer.WriteString(msg.String() + endline)
			log.Println(c.Network, "-->", msg.String())
			c.Writer.Flush()
		case <-c.QuitSend:
			log.Println(c.Network, "closing Sender")
			close(c.Input)
			close(c.QuitSend)
			return
		}
	}
}

func (c Connection) Receiver() {
	log.Println(c.Network, "Spawned receiver loop")
	for {
		raw, err := c.Reader.ReadString(delim)
		if err != nil {
			log.Println(c.Network, "error reading message", err.Error())
		}
		msg, err := ParseMessage(raw)
		if err != nil {
			log.Println(c.Network, "error decoding message", err.Error())
		}
		log.Println(c.Network, "<--", msg.String())
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

func (c Connection) Dial(server string, nick string, user string, realname string) error {

	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Println(c.Network, "Cannot connect to", server, "error:", err.Error())
		return err
	}
	log.Println(c.Network, "Connected to", server)
	c.Writer = bufio.NewWriter(conn)
	c.Reader = bufio.NewReader(conn)

	go c.Sender()
	go c.Receiver()

    log.Println(c.Network, "Initializing IRC connection")
        c.Input <- Message{
            Command:    "NICK",
            Trailing:   nick,
        }
        c.Input <- Message{
            Command:    "USER",
            Params:     []string{user, "0", "*"},
            Trailing:   realname,
        }

	return nil
}
