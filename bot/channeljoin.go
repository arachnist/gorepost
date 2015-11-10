package bot

import (
	"log"

	. "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func channeljoin(output chan irc.Message, msg irc.Message) {
	for _, channel := range C.Lookup(*msg.Context, "Channels").([]interface{}) {
		log.Println(msg.Context.Network, "joining channel", channel)
		output <- irc.Message{
			Command: "JOIN",
			Params:  []string{channel.(string)},
		}
	}
}

func init() {
	AddCallback("001", "channel join", channeljoin)
}
