// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"

	cfg "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

var trimTitle *regexp.Regexp

func getURLTitle(l string) string {
	title, err := httpGetXpath(l, "//head/title")
	if err == errElementNotFound {
		return "no title"
	} else if err != nil {
		return fmt.Sprint("error:", err)
	}

	return string(trimTitle.ReplaceAll([]byte(title), []byte{' '})[:])
}

func linktitle(output chan irc.Message, msg irc.Message) {
	var r []string

	for _, s := range strings.Split(msg.Trailing, " ") {
		buffer := new(bytes.Buffer)
		buffer.WriteString(s)

		b, err := regexp.Match("https?://", buffer.Bytes())
		if err != nil {
			log.Println("Context:", msg.Context, "linktitle regex error:", err)
			return
		}

		if b {
			t := getURLTitle(s)
			if t != "no title" {
				r = append(r, t)
			}
		}
	}

	if len(r) > 0 {
		t := cfg.LookupString(msg.Context, "LinkTitlePrefix") + strings.Join(r, cfg.LookupString(msg.Context, "LinkTitleDelimiter"))

		output <- reply(msg, t)
	}
}

func init() {
	trimTitle, _ = regexp.Compile("[\\s]+")
	addCallback("PRIVMSG", "LINKTITLE", linktitle)
}
