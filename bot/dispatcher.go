package bot

import (
	"log"
	"strings"

	cfg "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

var callbacks = make(map[string]map[string]func(chan irc.Message, irc.Message))

// addCallback registers callbacks that can be later dispatched by Dispatcher
func addCallback(command, name string, callback func(chan irc.Message, irc.Message)) {
	log.Println("adding callback", command, name)
	if _, ok := callbacks[command]; !ok {
		callbacks[command] = make(map[string]func(chan irc.Message, irc.Message))
	}
	callbacks[strings.ToUpper(command)][strings.ToUpper(name)] = callback
}

func elementInSlice(s []interface{}, e interface{}) bool {
	for _, se := range s {
		if se == e {
			return true
		}
	}

	return false
}

// Dispatcher takes irc messages and dispatches them to registered callbacks.
//
// It will take a message from input channel, check (based on message context)
// if the message should be dispatched and passes it to registered callback.
func Dispatcher(quit chan struct{}, output chan irc.Message, input chan irc.Message) {
	log.Println("spawned Dispatcher")
	for {
		select {
		case msg := <-input:
			if msg.Context["Source"] != "" {
				if elementInSlice(cfg.Lookup(msg.Context, "Ignore").([]interface{}), msg.Context["Source"]) {
					log.Println("Context:", msg.Context, "Ignoring", msg.Context["Source"])
					continue
				}
			}
			if callbacks[msg.Command] != nil {
				for i, f := range callbacks[msg.Command] {
					if elementInSlice(cfg.Lookup(msg.Context, "DisabledPlugins").([]interface{}), i) {
						log.Println("Context:", msg.Context, "Plugin disabled", i)
						continue
					}
					go f(output, msg)
				}
			}
		case <-quit:
			log.Println("closing Dispatcher")
			return
		}
	}
}
