package bot

import (
	"github.com/arachnist/gorepost/irc"
)

func papiez(output chan irc.Message, msg irc.Message) {
	output <- irc.Message{
		Command:  "PRIVMSG",
		Params:   []string{msg.Prefix.Name},
		Trailing: "pingity pong",
	}
}

func init() {
	AddMSGCallback(":ping", papiez)
}
