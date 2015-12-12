// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/arachnist/gorepost/irc"
)

var adjectives []string
var papiezLock sync.RWMutex

func papiez(output func(irc.Message), msg irc.Message) {
	args := strings.Split(msg.Trailing, " ")
	if args[0] != ":papiez" && args[0] != ":papież" {
		return
	}

	papiezLock.RLock()
	defer papiezLock.RUnlock()

	choice := "Papież " + adjectives[rand.Intn(len(adjectives))]

	output(reply(msg, choice))
}

func lazyPapiezInit() {
	defer papiezLock.Unlock()
	cfgLock.Lock()
	defer cfgLock.Unlock()
	var err error
	rand.Seed(time.Now().UnixNano())
	adjectives, err = readLines(cfg.LookupString(nil, "DictionaryAdjectives"))
	if err != nil {
		log.Println("failed to read adjectives", err)
		return
	}
	addCallback("PRIVMSG", "papiez", papiez)
}

func init() {
	papiezLock.Lock()
	log.Println("Defering \"papiez\" initialization")
	go lazyPapiezInit()
}
