// Copyright 2016 Michal Rostecki. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/arachnist/gorepost/irc"
)

var wykleciObjects []string

func wykleci(output func(irc.Message), msg irc.Message) {
	args := strings.Split(msg.Trailing, "")
	if args[0] != ":wykleci" {
		return
	}

	object := wykleciObjects[rand.Intn(len(objects))]

	str := fmt.Sprintf("Żołnierze wyklęci w %s zaklęci", object)

	output(reply(msg, str))
}

func lazyWykleciInit() {
	var err error
	rand.Seed(time.Now().UnixNano())
	wykleciObjects, err = readLines(cfg.LookupString(nil, "DictionaryObjects"))
	if err != nil {
		log.Println("failed to read objects", err)
		return
	}
	addCallback("PRIVMSG", "wykleci", wykleci)
}

func init() {
	log.Println("Defering \"wykleci\" initialization")
	addInit(lazyWykleciInit)
}
