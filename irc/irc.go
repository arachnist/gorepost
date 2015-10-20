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
            log.Println(c.Network, "-->", msg.String())
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
        log.Println(c.Network, "<-- RAW", raw)
		msg, err := ParseMessage(raw)
		if err != nil {
			log.Println(c.Network, "error decoding message", err.Error())
		}
		select {
		case c.Output <- *msg:
            log.Println(c.Network, "<--", msg.String())
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
		log.Println(c.Network, "Cannot connect to", server, "error:", err.Error())
		return err
	}
	log.Println(c.Network, "Connected to", server)
	c.Writer = bufio.NewWriter(conn)
	log.Println(c.Network, "Spawned bufio writer")
	c.Reader = bufio.NewReader(conn)
	log.Println(c.Network, "Spawned bufio reader")

	go c.Sender()
	go c.Receiver()

	return nil
}
