package irc

import (
	"bufio"
	"log"
	"net"
)

const delim byte = '\n'

var endline = []byte("\r\n")

type Connection struct {
	Network  string
	Input    chan Message
	Output   chan Message
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	QuitSend chan struct{}
	QuitRecv chan struct{}
}

func (c *Connection) Sender() {
	for {
		select {
		case msg := <-c.Input:
			c.Writer.WriteString(msg.String())
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
		raw, err := c.Reader.ReadString(delim)
		if err != nil {
			log.Println(c.Network, "error reading message", err.Error())
		}
		msg, err := ParseMessage(raw)
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

func (c *Connection) Dial(server string) error {

	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Println("Cannot connect to", server, "error:", err.Error())
		return err
	}
	c.Writer = bufio.NewWriter(conn)
	c.Reader = bufio.NewReader(conn)

	go c.Sender()
	go c.Receiver()

	return nil
}
