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

	cfg "github.com/arachnist/gorepost/config"
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

func seenrecord(output chan irc.Message, msg irc.Message) {
	if msg.Params == nil {
		return
	}

	v := seenRecord{
		Network: msg.Context["Network"],
		Target:  msg.Params[0],
		Action:  msg.Command,
		Time:    time.Now(),
		Text:    msg.Trailing,
	}

	b, err := json.Marshal(v)
	if err != nil {
		log.Println("Context:", msg.Context, "json marshal of seen record failed:", err)
		return
	}

	err = k.Set("seen/"+msg.Prefix.Name, b)
	if err != nil {
		log.Println("Context:", msg.Context, "error recording seen record:", err)
	}
}

func seen(output chan irc.Message, msg irc.Message) {
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
		output <- reply(msg, cfg.LookupString(msg.Context, "NotSeenMessage"))
		return
	} else if err != nil {
		output <- reply(msg, fmt.Sprint("error getting record for", args[1], err))
		return
	}

	err = json.Unmarshal(b, &v)
	if err != nil {
		output <- reply(msg, fmt.Sprint("error unmarshaling record for", args[1], err))
		return
	}

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
	}

	output <- reply(msg, r)
}

func init() {
	var err error
	var ktHost = "127.0.0.1"
	var ktPort = 1337

	log.Println("SEEN: connecting to KT")
	k, err = kt.NewConn(ktHost, ktPort, 4, 2*time.Second)
	if err != nil {
		log.Println("error connecting to kyoto tycoon", err)
		return
	}

	log.Println("Registering callbacks")
	addCallback("PRIVMSG", "seen", seen)
	addCallback("PRIVMSG", "seenrecord", seenrecord)
	addCallback("JOIN", "seenrecord", seenrecord)
	addCallback("PART", "seenrecord", seenrecord)
	addCallback("QUIT", "seenrecord", seenrecord)
	addCallback("NOTICE", "seenrecord", seenrecord)
}
