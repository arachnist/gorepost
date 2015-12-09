// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"log"
	"strings"
	"sync"

	cfg "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

var callbacks = make(map[string]map[string]func(func(irc.Message), irc.Message))
var callbackLock sync.RWMutex

// addCallback registers callbacks that can be later dispatched by Dispatcher
func addCallback(command, name string, callback func(func(irc.Message), irc.Message)) {
	callbackLock.Lock()
	defer callbackLock.Unlock()
	log.Println("adding callback", command, name)
	if _, ok := callbacks[command]; !ok {
		callbacks[command] = make(map[string]func(func(irc.Message), irc.Message))
	}
	callbacks[strings.ToUpper(command)][name] = callback
}

// Dispatcher takes irc messages and dispatches them to registered callbacks.
//
// It will take an input message, check (based on message context), if the
// message should be dispatched, and passes it to registered callback.
func Dispatcher(output func(irc.Message), input irc.Message) {
	if _, ok := cfg.LookupStringMap(input.Context, "Ignore")[input.Context["Source"]]; ok {
		log.Println("Context:", input.Context, "Ignoring", input.Context["Source"])
		return
	}

	callbackLock.RLock()
	defer callbackLock.RUnlock()
	if callbacks[input.Command] != nil {
		if len(cfg.LookupStringMap(input.Context, "WhitelistedPlugins")) > 0 {
			for i, f := range callbacks[input.Command] {
				if _, ok := cfg.LookupStringMap(input.Context, "DisabledPlugins")[i]; ok {
					log.Println("Context:", input.Context, "Plugin disabled", i)
					continue
				}
				if _, ok := cfg.LookupStringMap(input.Context, "WhitelistedPlugins")[i]; ok {
					go f(output, input)
				} else {
					log.Println("Context:", input.Context, "Plugin not whitelisted", i)
				}
			}
		} else {
			for i, f := range callbacks[input.Command] {
				if _, ok := cfg.LookupStringMap(input.Context, "DisabledPlugins")[i]; ok {
					log.Println("Context:", input.Context, "Plugin disabled", i)
					continue
				}
				go f(output, input)
			}
		}
	}
}
