// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/arachnist/gorepost/irc"
)

var trimTitle *regexp.Regexp
var trimLink *regexp.Regexp
var enc = charmap.ISO8859_2

func getURLTitle(l string) string {
	title, err := httpGetXpath(l, "//head/title")
	if err == errElementNotFound {
		return "no title"
	} else if err != nil {
		return fmt.Sprint("error:", err)
	}

	title = string(trimTitle.ReplaceAll([]byte(title), []byte{' '})[:])
	if !utf8.ValidString(title) {
		title, _, err = transform.String(enc.NewDecoder(), title)
		if err != nil {
			return fmt.Sprint("error:", err)
		}
	}

	return title
}

func linktitle(output func(irc.Message), msg irc.Message) {
	var r []string

	for _, s := range strings.Split(strings.Trim(msg.Trailing, "\001"), " ") {
		if s == "notitle" {
			return
		}

		b, _ := regexp.Match("https?://", []byte(s))

		s = string(trimLink.ReplaceAll([]byte(s), []byte("http"))[:])

		if b {
			t := getURLTitle(s)
			if t != "no title" {
				r = append(r, t)
			}
		}
	}

	if len(r) > 0 {
		t := cfg.LookupString(msg.Context, "LinkTitlePrefix") + strings.Join(r, cfg.LookupString(msg.Context, "LinkTitleDelimiter"))

		output(reply(msg, t))
	}
}

func init() {
	trimTitle, _ = regexp.Compile("[\\s]+")
	trimLink, _ = regexp.Compile("^.*?http")
	addCallback("PRIVMSG", "LINKTITLE", linktitle)
}
