package bot

import (
	"log"
	"strings"

	"github.com/arachnist/gorepost/irc"
)

var Callbacks = make(map[string]func(chan irc.Message, irc.Message))

func AddCallback(command string, callback func(chan irc.Message, irc.Message)) {
	Callbacks[command] = callback
}

func RemoveCallback(command string) {
	delete(Callbacks, command)
}

var MSGCallbacks = make(map[string]func(chan irc.Message, irc.Message))

func AddMSGCallback(command string, callback func(chan irc.Message, irc.Message)) {
	MSGCallbacks[command] = callback
}

func RemoveMSGCallback(command string) {
	delete(MSGCallbacks, command)
}

func Dispatcher(quit chan struct{}, output chan irc.Message, input chan irc.Message) {
	log.Println("spawned Dispatcher")
	for {
		select {
		case msg := <-input:
			if msg.Command == "PRIVMSG" {
				cmd := strings.Split(msg.Trailing, " ")[0]
				if MSGCallbacks[cmd] != nil {
					go MSGCallbacks[cmd](output, msg)
				}
			}
			if Callbacks[msg.Command] != nil {
				go Callbacks[msg.Command](output, msg)
			}
		case <-quit:
			log.Println("closing Dispatcher")
			return
		}
	}
}
