// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"log"

	"github.com/arachnist/gorepost/irc"
)

// channeljoin joins configured IRC channels when IRC server confirms we're good
// to go.
func channeljoin(output func(irc.Message), msg irc.Message) {
	for _, channel := range cfg.LookupStringSlice(msg.Context, "Channels") {
		log.Println(msg.Context["Network"], "joining channel", channel)
		output(irc.Message{
			Command: "JOIN",
			Params:  []string{channel},
		})
	}
}

func init() {
	addCallback("001", "channel join", channeljoin)
}
