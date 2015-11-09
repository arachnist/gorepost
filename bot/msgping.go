package bot

import (
	"strings"

	. "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func ping(c Context, output chan irc.Message, msg irc.Message) {
	if strings.Split(msg.Trailing, " ")[0] != ":ping" {
		return
	}

	output <- irc.Message{
		Command:  "PRIVMSG",
		Params:   []string{msg.Prefix.Name},
		Trailing: "pingity pong",
	}
}

func init() {
	AddCallback("PRIVMSG", "ping", ping)
}
