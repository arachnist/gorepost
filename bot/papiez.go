// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/arachnist/gorepost/irc"
)

var adjectives []string

func papiez(output func(irc.Message), msg irc.Message) {
	args := strings.Split(msg.Trailing, " ")
	if args[0] != ":papiez" && args[0] != ":papież" {
		return
	}

	choice := "Papież " + adjectives[rand.Intn(len(adjectives))]

	output(reply(msg, choice))
}

func lazyPapiezInit() {
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
	log.Println("Defering \"papiez\" initialization")
	addInit(lazyPapiezInit)
}
