// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudflare/gokabinet/kt"

	"github.com/arachnist/gorepost/irc"
)

var k *kt.Conn

type seenRecord struct {
	Network string
	Target  string
	Action  string
	Time    time.Time
	Text    string
}

func seenrecord(output func(irc.Message), msg irc.Message) {
	var target string

	if msg.Params == nil || msg.Command == "QUIT" {
		target = ""
	} else {
		target = msg.Params[0]
	}

	v := seenRecord{
		Network: msg.Context["Network"],
		Target:  target,
		Action:  msg.Command,
		Time:    time.Now(),
		Text:    msg.Trailing,
	}

	b, _ := json.Marshal(v)

	err := k.Set("seen/"+msg.Prefix.Name, b)
	if err != nil {
		log.Println("Context:", msg.Context, "error recording seen record:", err)
	}
}

func seen(output func(irc.Message), msg irc.Message) {
	var v seenRecord
	var r string

	args := strings.Split(msg.Trailing, " ")

	if args[0] != ":seen" {
		return
	}
	if len(args) < 2 {
		return
	}

	b, err := k.GetBytes("seen/" + args[1])
	if err == kt.ErrNotFound {
		output(reply(msg, cfg.LookupString(msg.Context, "NotSeenMessage")))
		return
	} else if err != nil {
		output(reply(msg, fmt.Sprint("error getting record for", args[1], err)))
		return
	}

	_ = json.Unmarshal(b, &v)

	r = fmt.Sprintf("Last seen %s on %s/%s at %v ", args[1], v.Network, v.Target, v.Time.Round(time.Second))

	switch v.Action {
	case "JOIN":
		r += "joining"
	case "PART":
		r += fmt.Sprint("leaving: ", v.Text)
	case "QUIT":
		r += fmt.Sprint("quitting with reasson: ", v.Text)
	case "PRIVMSG":
		r += fmt.Sprint("saying: ", v.Text)
	case "NOTICE":
		r += fmt.Sprint("noticing: ", v.Text)
	}

	output(reply(msg, r))
}

func seenInit() {
	var ktHost = cfg.LookupString(nil, "KTHost")
	var ktPort = cfg.LookupInt(nil, "KTPort")

	log.Println("seen: connecting to KT")
	k = kt.NewConn(ktHost, ktPort, 4, 2*time.Second)

	log.Println("Registering callbacks")
	addCallback("PRIVMSG", "seen", seen)
	addCallback("PRIVMSG", "seenrecord", seenrecord)
	addCallback("JOIN", "seenrecord", seenrecord)
	addCallback("PART", "seenrecord", seenrecord)
	addCallback("QUIT", "seenrecord", seenrecord)
	addCallback("NOTICE", "seenrecord", seenrecord)
}

func init() {
	log.Println("Defering \"seen\" initialization")
	addInit(seenInit)
}
