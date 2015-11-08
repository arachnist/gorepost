package bot

import (
	"log"
	"time"

	"github.com/arachnist/gorepost/irc"
)

var Callbacks = make(map[string]func(chan irc.Message, irc.Message))

func AddCallback(command string, callback func(chan irc.Message, irc.Message)) {
	Callbacks[command] = callback
}

func RemoveCallback(command string) {
	delete(Callbacks, command)
}

func Dispatcher(output *chan irc.Message, input *chan irc.Message) {
	// FIXME
	time.Sleep(time.Second * 2)
	log.Println("spawned Dispatcher")
	for {
		msg := <-*input
		if Callbacks[msg.Command] != nil {
			go Callbacks[msg.Command](*output, msg)
		}
	}
}
