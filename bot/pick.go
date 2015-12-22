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
	var args []string
	if !strings.HasPrefix(msg.Trailing, ":pick ") {
		return
	}

	a := strings.TrimPrefix(msg.Trailing, ":pick ")

	if strings.Contains(a, ", ") {
		args = strings.Split(a, ", ")
	} else if strings.Contains(a, ",") {
		args = strings.Split(a, ",")
	} else {
		args = strings.Fields(a)
	}

	choice := args[rand.Intn(len(args))]

	output(reply(msg, choice))
}

func init() {
	rand.Seed(time.Now().UnixNano())
	addCallback("PRIVMSG", "pick", pick)
}
