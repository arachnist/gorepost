package bot

import (
	"log"
	"strings"

	. "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

var Callbacks = make(map[string]map[string]func(chan irc.Message, irc.Message))

func AddCallback(command, name string, callback func(chan irc.Message, irc.Message)) {
	log.Println("adding callback", command, name)
	if _, ok := Callbacks[command]; !ok {
		Callbacks[command] = make(map[string]func(chan irc.Message, irc.Message))
	}
	Callbacks[strings.ToUpper(command)][strings.ToUpper(name)] = callback
}

func RemoveCallback(command, name string) {
	delete(Callbacks[command], name)
}

func elementInSlice(s []interface{}, e interface{}) bool {
	for _, se := range s {
		if se == e {
			return true
		}
	}

	return false
}

func Dispatcher(quit chan struct{}, output chan irc.Message, input chan irc.Message) {
	log.Println("spawned Dispatcher")
	for {
		select {
		case msg := <-input:
			if msg.Context["Source"] != "" {
				if elementInSlice(C.Lookup(msg.Context, "Ignore").([]interface{}), msg.Context["Source"]) {
					log.Println("Context:", msg.Context, "Ignoring", msg.Context["Source"])
					continue
				}
			}
			if Callbacks[msg.Command] != nil {
				for i, f := range Callbacks[msg.Command] {
					if elementInSlice(C.Lookup(msg.Context, "DisabledPlugins").([]interface{}), i) {
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
