// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package bot file msgping.go contains a simple "ping-pong" style callback
package bot

import (
	"strings"

	"github.com/arachnist/gorepost/irc"
)

func ping(output chan irc.Message, msg irc.Message) {
	if strings.Split(msg.Trailing, " ")[0] != ":ping" {
		return
	}

	output <- irc.Message{
		Command:  "PRIVMSG",
		Params:   []string{msg.Prefix.Name},
		Trailing: "pingity pong",
	}
}

func init() {
	addCallback("PRIVMSG", "msgping", ping)
}
