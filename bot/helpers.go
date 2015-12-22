// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xpath"

	"github.com/arachnist/gorepost/irc"
)

var errElementNotFound = errors.New("element not found in document")

func httpGet(link string) ([]byte, error) {
	var buf []byte
	tr := &http.Transport{
		TLSHandshakeTimeout:   20 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	// 5MiB
	if resp.ContentLength > 5*1024*1024 || resp.ContentLength < 0 {
		buf = make([]byte, 5*1024*1024)
	} else if resp.ContentLength == 0 {
		return []byte{}, nil
	} else {
		buf = make([]byte, resp.ContentLength)
	}

	i, err := io.ReadFull(resp.Body, buf)
	if err == io.ErrUnexpectedEOF {
		buf = buf[:i]
	} else if err != nil {
		return []byte{}, err
	}

	return buf, nil
}

func httpGetXpath(link, xpathStr string) (string, error) {
	buf, err := httpGet(link)
	if err != nil {
		return "", err
	}

	doc, err := gokogiri.ParseHtml(buf)
	defer doc.Free()
	if err != nil {
		return "", err
	}
	if doc.Root() == nil {
		return "", errElementNotFound
	}

	xpath := xpath.Compile(xpathStr)
	defer xpath.Free()
	sr, err := doc.Root().Search(xpath)
	if err != nil {
		return "", err
	}

	if len(sr) > 0 {
		return sr[0].InnerHtml(), nil
	}

	return "", errElementNotFound
}

func reply(msg irc.Message, text string) irc.Message {
	if msg.Params[0] == cfg.LookupString(msg.Context, "Nick") {
		return irc.Message{
			Command:  "PRIVMSG",
			Params:   []string{msg.Prefix.Name},
			Trailing: text,
		}
	}

	return irc.Message{
		Command:  "PRIVMSG",
		Params:   msg.Params,
		Trailing: text,
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
