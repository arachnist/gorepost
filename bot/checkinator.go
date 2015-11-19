// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
	var now []string
	var recently []string

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

	rmsg = "at:"

	for _, u := range values.Users {
		t := time.Unix(int64(u.Timestamp), 0)
		if t.Add(time.Minute * 10).After(time.Now()) {
			now = append(now, u.Login)
		} else {
			recently = append(recently, u.Login)
		}
	}

	if len(now) > 0 {
		rmsg += " now: "
		rmsg += strings.Join(now, ", ")
	}
	if len(recently) > 0 {
		rmsg += " recently: "
		rmsg += strings.Join(recently, ", ")
	}
	if len(now) == 0 && len(recently) == 0 {
		rmsg += " Wieje sandałem, z masłem"
	}
	if values.Kektops > 0 {
		rmsg += fmt.Sprintf("; kektops: %d", values.Kektops)
	}
	if values.Unknown > 0 {
		rmsg += fmt.Sprintf("; unknown: %d", values.Unknown)
	}

	output <- reply(msg, rmsg)
}

func init() {
	addCallback("PRIVMSG", "at", at)
}
