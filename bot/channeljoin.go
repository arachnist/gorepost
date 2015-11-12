// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package bot file channeljoin.go contains channeljoin callback which joins
// configured IRC channels on successful connection
package bot

import (
	"log"

	cfg "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

// channeljoin joins configured IRC channels when IRC server confirms we're good
// to go.
func channeljoin(output chan irc.Message, msg irc.Message) {
	for _, channel := range cfg.Lookup(msg.Context, "Channels").([]interface{}) {
		log.Println(msg.Context["Network"], "joining channel", channel)
		output <- irc.Message{
			Command: "JOIN",
			Params:  []string{channel.(string)},
		}
	}
}

func init() {
	addCallback("001", "channel join", channeljoin)
}
