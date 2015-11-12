// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package bot

import (
	"github.com/arachnist/gorepost/irc"
)

func invite(output chan irc.Message, msg irc.Message) {
	output <- irc.Message{
		Command: "JOIN",
		Params:  []string{msg.Trailing},
	}
}

func init() {
	addCallback("INVITE", "invitki", invite)
}
