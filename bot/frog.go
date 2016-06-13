// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/arachnist/gorepost/irc"
)

type tip struct {
	Number int    `json:"number"`
	Tip    string `json:"tip"`
}

type tips struct {
	Tips []tip `json:"tips"`
	lock sync.RWMutex
}

func (tips *tips) fetchTips() error {
	tips.lock.Lock()
	defer tips.lock.Unlock()

	data, err := httpGet("http://frog.tips/api/1/tips/")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, tips)
	if err != nil {
		return err
	}

	return nil
}

func (tips *tips) popTip() string {
	if len(tips.Tips) == 0 {
		if err := tips.fetchTips(); err != nil {
			return fmt.Sprint(err)
		}
	}

	tips.lock.RLock()
	defer tips.lock.RUnlock()

	rmsg := tips.Tips[len(tips.Tips)-1].Tip
	tips.Tips = tips.Tips[:len(tips.Tips)-1]

	return rmsg
}

var t tips

func frog(output func(irc.Message), msg irc.Message) {
	if strings.Split(msg.Trailing, " ")[0] != ":frog" {
		return
	}

	output(reply(msg, t.popTip()))
}

func init() {
	addCallback("PRIVMSG", "frog", frog)
}
