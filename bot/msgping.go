// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"strings"

	"github.com/arachnist/gorepost/irc"
)

func ping(output func(irc.Message), msg irc.Message) {
	if strings.Split(msg.Trailing, " ")[0] != ":ping" {
		return
	}

	output(reply(msg, "pingity pong"))
}

func init() {
	addCallback("PRIVMSG", "msgping", ping)
}
