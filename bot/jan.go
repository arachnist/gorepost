// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	cfg "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

var objects []string
var predicates []string
var janLock sync.Mutex

func jan(output func(irc.Message), msg irc.Message) {
	args := strings.Split(msg.Trailing, " ")
	if args[0] != ":jan" {
		return
	}
	var predicate string
	var object string

	if len(args) > 1 {
		if strings.HasSuffix(args[1], "ł") {
			predicate = args[1]
			object = objects[rand.Intn(len(objects))]
		} else {
			object = args[1]
			predicate = predicates[rand.Intn(len(predicates))]
		}
	} else {
		predicate = predicates[rand.Intn(len(predicates))]
		object = objects[rand.Intn(len(objects))]
	}

	str := "Jan Paweł II " + predicate + " małe " + object

	output(reply(msg, str))
}

func lazyJanInit() {
	defer janLock.Unlock()
	var err error
	rand.Seed(time.Now().UnixNano())
	objects, err = readLines(cfg.LookupString(nil, "DictionaryObjects"))
	if err != nil {
		log.Println("failed to read objects", err)
		return
	}
	predicates, err = readLines(cfg.LookupString(nil, "DictionaryVerbs"))
	if err != nil {
		log.Println("failed to read predicates", err)
		return
	}
	addCallback("PRIVMSG", "jan", jan)
}

func init() {
	janLock.Lock()
	log.Println("Defering \"jan\" initialization")
	go lazyJanInit()
}
