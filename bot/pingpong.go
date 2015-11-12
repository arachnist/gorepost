package bot

import (
	"github.com/arachnist/gorepost/irc"
)

// pingpong responds to server pings. IRC servers disconnect idle clients that
// don't respond to PINGs.
func pingpong(output chan irc.Message, msg irc.Message) {
	output <- irc.Message{
		Command:  "PONG",
		Trailing: msg.Trailing,
	}
}

func init() {
	addCallback("PING", "ping", pingpong)
}
