// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"math/rand"
	"strings"
	"time"

	"github.com/arachnist/gorepost/irc"
)

func pick(output func(irc.Message), msg irc.Message) {
	args := strings.Split(msg.Trailing, " ")
	if args[0] != ":pick" {
		return
	}

	// don't pick :pick
	choice := args[rand.Intn(len(args)-1)+1]

	output(reply(msg, choice))
}

func init() {
	rand.Seed(time.Now().UnixNano())
	addCallback("PRIVMSG", "pick", pick)
}
