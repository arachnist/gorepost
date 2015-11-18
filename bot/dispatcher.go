// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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

func elementInSlice(s []string, e string) bool {
	for _, se := range s {
		if se == e {
			return true
		}
	}

	return false
}

// Dispatcher takes irc messages and dispatches them to registered callbacks.
//
// It will take an input message, check (based on message context), if the
// message should be dispatched, and passes it to registered callback.
func Dispatcher(output chan irc.Message, input irc.Message) {
	if input.Context["Source"] != "" {
		if elementInSlice(cfg.LookupStringSlice(input.Context, "Ignore"), input.Context["Source"]) {
			log.Println("Context:", input.Context, "Ignoring", input.Context["Source"])
			return
		}
	}
	if callbacks[input.Command] != nil {
		for i, f := range callbacks[input.Command] {
			if elementInSlice(cfg.LookupStringSlice(input.Context, "DisabledPlugins"), i) {
				log.Println("Context:", input.Context, "Plugin disabled", i)
				return
			}
			go f(output, input)
		}
	}
}
