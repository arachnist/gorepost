// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/arachnist/gorepost/irc"
)

func google(output func(irc.Message), msg irc.Message) {
	if strings.Split(msg.Trailing, " ")[0] != ":g" {
		return
	}

	query := strings.TrimPrefix(msg.Trailing, ":g ")

	req, _ := http.NewRequest("GET", "https://ajax.googleapis.com/ajax/services/search/web?v=1.0", nil)

	q := req.URL.Query()
	q.Set("q", query)

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		output(reply(msg, "problem connecting to google"))
		return
	}

	defer resp.Body.Close()

	var data struct {
		ResponseData struct {
			Results []struct {
				TitleNoFormatting string
				URL               string
			}
		}
	}
	if errJ := json.NewDecoder(resp.Body).Decode(&data); errJ != nil {
		output(reply(msg, "problem decoding google response"))
		return
	}
	if len(data.ResponseData.Results) > 0 {
		res := data.ResponseData.Results[0]
		link, err := url.QueryUnescape(res.URL)
		if err != nil {
			output(reply(msg, "problem decoding url"))
		}
		output(reply(msg, res.TitleNoFormatting+" "+link))
	}
}

func init() {
	addCallback("PRIVMSG", "google", google)
}
