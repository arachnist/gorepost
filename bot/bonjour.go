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

func getPicDate() time.Time {
	for {
		now := time.Now()
		epoch := time.Date(2018, time.December, 10, 0, 0, 0, 0, time.UTC)
		delta := now.Sub(epoch)
		r := rand.Int63n(int64(delta.Seconds()))
		newdate := epoch.Add(time.Duration(r) * time.Second)
		if newdate.Weekday() == 0 || newdate.Weekday() == 6 {
			continue
		}
		return newdate
	}
}

func bonjour(output func(irc.Message), msg irc.Message) {
	var rmsg string

	if strings.Split(msg.Trailing, " ")[0] != ":bonjour" {
		return
	}

	d := getPicDate()
	year, month, day := d.Date()
	img, err := httpGetXpath("http://www.bonjourmadame.fr/"+fmt.Sprintf("%d/%d/%d/", year, month, day), "//div[@class='post-content']//p/img/@src")
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
