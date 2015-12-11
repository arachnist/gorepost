// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/arachnist/gorepost/irc"
)

func roll(output func(irc.Message), msg irc.Message) {
	args := strings.Split(msg.Trailing, " ")
	if args[0] != ":roll" {
		return
	}

	var err error
	rolls := 1

	if len(args) == 3 {
		rolls, err = strconv.Atoi(args[2])
		if err != nil || rolls < 1 {
			output(reply(msg, "Usage: :roll <sides int> <rolls int>, each roll is [0, n)+1, size has to be >0"))
			return
		}
	} else if len(args) != 2 {
		output(reply(msg, "Usage: :roll <sides int> <rolls int>, each roll is [0, n)+1, size has to be >0"))
		return
	}

	sides, err := strconv.Atoi(args[1])
	if err != nil {
		output(reply(msg, "Usage: :roll <sides int> <rolls int>, each roll is [0, n)+1, size has to be >0"))
		return
	}
	if sides <= 0 {
		output(reply(msg, "Usage: :roll <sides int> <rolls int>, each roll is [0, n)+1, size has to be >0"))
		return
	}
	if sides > 1000000 || rolls > 1000000 {
		output(reply(msg, "Number of rolls and dice size is limited to 1000000"))
		return
	}

	sum := rolls
	for i := 0; i < rolls; i++ {
		sum += rand.Intn(sides)
	}

	output(reply(msg, strconv.Itoa(sum)))
}

func init() {
	rand.Seed(time.Now().UnixNano())
	addCallback("PRIVMSG", "roll", roll)
}
