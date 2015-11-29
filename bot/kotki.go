// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/arachnist/gorepost/irc"
)

var errNotReally = errors.New("not an error")

func redirectError(*http.Request, []*http.Request) error {
	return errNotReally
}

func kotki(output func(irc.Message), msg irc.Message) {
	if strings.Split(msg.Trailing, " ")[0] != ":kotki" {
		return
	}

	var rmsg string
	tr := &http.Transport{
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	}
	client := &http.Client{
		Transport:     tr,
		CheckRedirect: redirectError,
	}

	resp, _ := client.Get("http://thecatapi.com/api/images/get?format=src&type=png")
	rurl, _ := resp.Location()
	rmsg = rurl.String()

	output(reply(msg, rmsg))
}

func init() {
	addCallback("PRIVMSG", "kotki", kotki)
}
