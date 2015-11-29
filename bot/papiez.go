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

var adverbs []string

func papiez(output func(irc.Message), msg irc.Message) {
	args := strings.Split(msg.Trailing, " ")
	if args[0] != ":papiez" && args[0] != ":papież" {
		return
	}

	choice := "Papież " + adverbs[rand.Intn(len(adverbs))]

	output(reply(msg, choice))
}

func init() {
	var err error
	rand.Seed(time.Now().UnixNano())
	adverbs, err = readLines("/home/repost/przymiotniki")
	if err != nil {
		log.Println("failed to read adverbs", err)
		return
	}
	addCallback("PRIVMSG", "papiez", papiez)
}
