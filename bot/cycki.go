// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/arachnist/gorepost/irc"
)

func cycki(output chan irc.Message, msg irc.Message) {
	var rmsg string

	if strings.Split(msg.Trailing, " ")[0] != ":cycki" {
		return
	}

	img, err := httpGetXpath("http://www.bonjourmadame.fr/page/"+string(rand.Intn(2370)+1), "//div[@class='photo post']//img/@src")
	if err != nil {
		rmsg = fmt.Sprint("error:", err)
	} else {
		rmsg = "bonjour (nsfw): " + img
	}

	output <- reply(msg, rmsg)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	addCallback("PRIVMSG", "bonjour", cycki)
}
