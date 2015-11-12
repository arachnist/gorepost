// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package bot file pingpong.go contains a callback handling IRC server PING
// command handler.
package bot

import (
	"github.com/arachnist/gorepost/irc"
)

// pingpong responds to server pings. IRC servers disconnect idle clients that
// don't respond to PINGs.
func pingpong(output chan irc.Message, msg irc.Message) {
	output <- irc.Message{
		Command:  "PONG",
		Trailing: msg.Trailing,
	}
}

func init() {
	addCallback("PING", "ping", pingpong)
}
