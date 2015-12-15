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

func bonjour(output func(irc.Message), msg irc.Message) {
	var rmsg string

	if strings.Split(msg.Trailing, " ")[0] != ":bonjour" {
		return
	}

	t, _ := time.Parse("2006-01-02", "2015-12-01")
	max := int(time.Now().Sub(t).Hours())/24 + 1

	img, err := httpGetXpath("http://ditesbonjouralamadame.tumblr.com/page/"+fmt.Sprintf("%d", rand.Intn(max)+1), "//div[@class='photo post']//a/@href")
	if err != nil {
		rmsg = fmt.Sprint("error:", err)
	} else {
		rmsg = "bonjour (nsfw): " + img
	}

	output(reply(msg, rmsg))
}

func init() {
	rand.Seed(time.Now().UnixNano())
	addCallback("PRIVMSG", "bonjour", bonjour)
}
