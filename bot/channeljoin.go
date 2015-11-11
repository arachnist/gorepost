package bot

import (
	"log"

	cfg "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func channeljoin(output chan irc.Message, msg irc.Message) {
	for _, channel := range cfg.Lookup(msg.Context, "Channels").([]interface{}) {
		log.Println(msg.Context["Network"], "joining channel", channel)
		output <- irc.Message{
			Command: "JOIN",
			Params:  []string{channel.(string)},
		}
	}
}

func init() {
	addCallback("001", "channel join", channeljoin)
}
