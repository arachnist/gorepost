package bot

import (
	"log"
	"strings"

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

func Dispatcher(quit chan struct{}, output chan irc.Message, input chan irc.Message) {
	log.Println("spawned Dispatcher")
	for {
		select {
		case msg := <-input:
			if Callbacks[msg.Command] != nil {
				for _, f := range Callbacks[msg.Command] {
					go f(output, msg)
				}
			}
		case <-quit:
			log.Println("closing Dispatcher")
			return
		}
	}
}
