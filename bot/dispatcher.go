package bot

import (
	"log"

	"github.com/arachnist/gorepost/irc"
)

var Callbacks = make(map[string]func(chan irc.Message, irc.Message))

func AddCallback(command string, callback func(chan irc.Message, irc.Message)) {
	Callbacks[command] = callback
}

func RemoveCallback(command string) {
	delete(Callbacks, command)
}

func Dispatcher(quit chan struct{}, output chan irc.Message, input chan irc.Message) {
	log.Println("spawned Dispatcher")
	for {
		select {
		case msg := <-input:
			if Callbacks[msg.Command] != nil {
				go Callbacks[msg.Command](output, msg)
			}
		case <-quit:
			log.Println("closing Dispatcher")
			return
		}
	}
}
