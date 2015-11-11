package bot

import (
	"github.com/arachnist/gorepost/irc"
)

func pingpong(output chan irc.Message, msg irc.Message) {
	output <- irc.Message{
		Command:  "PONG",
		Trailing: msg.Trailing,
	}
}

func init() {
	addCallback("PING", "ping", pingpong)
}
