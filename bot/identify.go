// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"strings"

	cfg "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

var identified map[string]map[string]string

func IsIdentified(msg irc.Message) bool {
	if msg.Prefix == nil {
		return false
	}
	if identified == nil {
		return false
	}
	net := msg.Context["Network"]
	if identified[net] == nil {
		return false
	}
	if identified[net][msg.Prefix.Name] == "" {
		return false
	}

	return true
}

func identify(output func(irc.Message), msg irc.Message) {
	if strings.Split(msg.Trailing, " ")[0] != ":identify" {
		return
	}

	output(irc.Message{
		Command: "WHOIS",
		Params:  []string{msg.Prefix.Name},
	})
}

func registerIdentification(output func(irc.Message), msg irc.Message) {
	net := msg.Context["Network"]

	if identified == nil {
		identified = make(map[string]map[string]string)
	}

	if _, ok := identified[net]; !ok {
		identified[net] = make(map[string]string)
	}

	identified[net][msg.Params[1]] = msg.Params[2]
}

func listIdentified(output func(irc.Message), msg irc.Message) {
	if strings.Split(msg.Trailing, " ")[0] != ":idlist" {
		return
	}

	if cfg.LookupInt(msg.Context, "AccessLevel") < 10 {
		output(reply(msg, "access denied"))
		return
	}

	var r []string

	for _, net := range identified {
		for k, v := range net {
			r = append(r, k+" identified as "+v)
		}
	}

	output(reply(msg, strings.Join(r, "; ")))
}

func init() {
	addCallback("PRIVMSG", "identify", identify)
	addCallback("PRIVMSG", "list identified", listIdentified)
	addCallback("330", "register identification", registerIdentification)
}
