package bot

import (
	. "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func pingpong(c Context, output chan irc.Message, msg irc.Message) {
	output <- irc.Message{
		Command:  "PONG",
		Trailing: msg.Trailing,
	}
}

func init() {
	AddCallback("PING", "ping", pingpong)
}
