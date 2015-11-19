// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/arachnist/gorepost/irc"
)

type user struct {
	Timestamp   float64
	Login       string
	Pretty_time string
}

type checkinator struct {
	Kektops int
	Unknown int
	Users   []user
}

func at(output chan irc.Message, msg irc.Message) {
	var rmsg string
	var values checkinator

	if strings.Split(msg.Trailing, " ")[0] != ":at" {
		return
	}

	data, err := httpGet("https://at.hackerspace.pl/api")
	if err != nil {
		output <- reply(msg, fmt.Sprint("error:", err))
		return
	}

	err = json.Unmarshal(data, &values)
	if err != nil {
		output <- reply(msg, fmt.Sprint("error:", err))
		return
	}

	rmsg = fmt.Sprintf("%+v", values)

	output <- reply(msg, rmsg)
}

func init() {
	addCallback("PRIVMSG", "at", at)
}
