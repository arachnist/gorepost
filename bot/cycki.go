// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arachnist/gorepost/irc"
)

var stripCycki *regexp.Regexp

func cycki(output func(irc.Message), msg irc.Message) {
	var rmsg string

	if strings.Split(msg.Trailing, " ")[0] != ":cycki" {
		return
	}

	img, err := httpGetXpath("http://oboobs.ru/random/", "//img/@src")
	if err != nil {
		rmsg = fmt.Sprint("error:", err)
	} else {
		rmsg = "cycki (nsfw): " + string(stripCycki.ReplaceAll([]byte(img), []byte("")))
	}

	output(reply(msg, rmsg))
}

func init() {
	stripCycki, _ = regexp.Compile("_preview")
	addCallback("PRIVMSG", "cycki", cycki)
}
