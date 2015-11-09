package bot

import (
	"log"

	. "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func channeljoin(c Context, output chan irc.Message, msg irc.Message) {
	for _, channel := range C.Lookup(c, "Channels").([]string) {
		log.Println(c.Network, "joining channel", channel)
		output <- irc.Message{
			Command: "JOIN",
			Params:  []string{channel},
		}
	}
}

func init() {
	AddCallback("001", "channel join", channeljoin)
}
