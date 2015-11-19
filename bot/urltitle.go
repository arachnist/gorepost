// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xpath"

	cfg "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func getUrlTitle(l string) string {
	var buf []byte
	tr := &http.Transport{
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(l)
	if err != nil {
		return fmt.Sprintf("error:", err)
	}

	// 5MiB
	if resp.ContentLength > 5*1024*1024 || resp.ContentLength < 0 {
		buf = make([]byte, 5*1024*1024)
	} else if resp.ContentLength == 0 {
		return "empty"
	} else {
		buf = make([]byte, resp.ContentLength)
	}

	i, err := io.ReadFull(resp.Body, buf)
	if err == io.ErrUnexpectedEOF {
		buf = buf[:i]
	} else if err != nil {
		return fmt.Sprintf("error:", err)
	}

	doc, err := gokogiri.ParseHtml(buf)
	defer doc.Free()
	if err != nil {
		return fmt.Sprintf("error:", err)
	}

	xpath := xpath.Compile("//head/title")
	sr, err := doc.Root().Search(xpath)

	if len(sr) > 0 {
		return sr[0].InnerHtml()
	} else {
		return "no title"
	}
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
			r = append(r, getUrlTitle(s))
		}
	}

	if len(r) > 0 {
		t := strings.Join(
			append([]string{cfg.LookupString(msg.Context, "LinkTitlePrefix")}, r...),
			cfg.LookupString(msg.Context, "LinkTitleDelimiter"),
		)

		if msg.Params[0] == cfg.LookupString(msg.Context, "Nick") {
			output <- irc.Message{
				Command:  "PRIVMSG",
				Params:   []string{msg.Prefix.Name},
				Trailing: t,
			}
		} else {
			output <- irc.Message{
				Command:  "PRIVMSG",
				Params:   msg.Params,
				Trailing: t,
			}
		}
	}
}

func init() {
	addCallback("PRIVMSG", "LINKTITLE", linktitle)
}
