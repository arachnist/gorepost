// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/arachnist/gorepost/irc"
)

type tip struct {
	Number int    `json:"number"`
	Tip    string `json:"tip"`
}

type tips struct {
	Tips []tip `json:"tips"`
}

func frog(output func(irc.Message), msg irc.Message) {
	var values tips

	if strings.Split(msg.Trailing, " ")[0] != ":frog" {
		return
	}

	data, err := httpGet("http://frog.tips/api/1/tips/")
	if err != nil {
		output(reply(msg, fmt.Sprint("error:", err)))
		return
	}

	err = json.Unmarshal(data, &values)
	if err != nil {
		output(reply(msg, fmt.Sprint("error:", err)))
		return
	}

	tip := values.Tips[rand.Intn(len(values.Tips))]

	output(reply(msg, fmt.Sprintf("frog tip #%d: %s", tip.Number, tip.Tip)))
}

func init() {
	addCallback("PRIVMSG", "frog", frog)
}
