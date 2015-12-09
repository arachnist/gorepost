// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"fmt"
	"log"
	"regexp"

	cfg "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func nickserv(output func(irc.Message), msg irc.Message) {
	if msg.Prefix.String() != cfg.LookupString(msg.Context, "NickServPrefix") {
		log.Println("Context:", msg.Context, "Someone is spoofing nickserv!")
		return
	}

	regexStr := cfg.LookupString(msg.Context, "NickServRegex")
	if regexStr == "" {
		regexStr = "^This nickname is registered"
	}

	b, _ := regexp.Match(regexStr, []byte(msg.Trailing))

	if !b {
		return
	}

	log.Println("Context:", msg.Context, "Identifying to nickserv!")
	output(reply(msg, fmt.Sprintf("IDENTIFY %s", cfg.LookupString(msg.Context, "NickServPassword"))))
}

func joinsecuredchannels(output func(irc.Message), msg irc.Message) {
	if msg.Prefix.String() != cfg.LookupString(msg.Context, "NickServPrefix") {
		log.Println("Context:", msg.Context, "Someone is spoofing nickserv!")
		return
	}

	regexStr := cfg.LookupString(msg.Context, "NickServRegexOK")
	if regexStr == "" {
		regexStr = "^You are now identified"
	}

	b, _ := regexp.Match(regexStr, []byte(msg.Trailing))

	if !b {
		return
	}

	channels := cfg.LookupStringSlice(msg.Context, "SecuredChannels")
	if len(channels) < 1 || channels[0] == "" {
		return
	}

	for _, channel := range channels {
		log.Println(msg.Context["Network"], "joining channel", channel)
		output(irc.Message{
			Command: "JOIN",
			Params:  []string{channel},
		})
	}
}

func init() {
	addCallback("NOTICE", "nickserv", nickserv)
	addCallback("NOTICE", "join +i-only channels", joinsecuredchannels)
}
