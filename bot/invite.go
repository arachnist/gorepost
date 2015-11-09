package bot

import (
	"github.com/arachnist/gorepost/irc"
)

func invite(output chan irc.Message, msg irc.Message) {
	output <- irc.Message{
		Command: "JOIN",
		Params:  []string{msg.Trailing},
	}
}

func init() {
	AddCallback("INVITE", "invitki", invite)
}
