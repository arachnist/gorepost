// +build go1.4

// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	cfg "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

var eventTests = []struct {
	in          irc.Message
	expectedOut []irc.Message
}{
	{ // "ping"
		in: irc.Message{
			Command:  "PING",
			Trailing: "foobar",
		},
		expectedOut: []irc.Message{
			{
				Command:  "PONG",
				Trailing: "foobar",
			},
		},
	},
	{ // "invitki"
		in: irc.Message{
			Command:  "INVITE",
			Trailing: "#test-channel",
		},
		expectedOut: []irc.Message{
			{
				Command: "JOIN",
				Params:  []string{"#test-channel"},
			},
		},
	},
	{ // "channel join"
		in: irc.Message{
			Command: "001",
			Context: map[string]string{
				"Network": "TestNetwork",
			},
		},
		expectedOut: []irc.Message{
			{
				Command: "JOIN",
				Params:  []string{"#testchan-1"},
			},
			{
				Command: "JOIN",
				Params:  []string{"#testchan-2"},
			},
		},
	},
	{ // "msgping",
		in: irc.Message{
			Command:  "PRIVMSG",
			Trailing: ":ping",
			Params:   []string{"#testchan-1"},
			Prefix: &irc.Prefix{
				Name: "idontexist",
			},
		},
		expectedOut: []irc.Message{
			{
				Command:  "PRIVMSG",
				Params:   []string{"#testchan-1"},
				Trailing: "pingity pong",
			},
		},
	},
	{ // "nickserv"
		in: irc.Message{
			Command:  "NOTICE",
			Params:   []string{"gorepost"},
			Trailing: "This nickname is registered. Please choose a different nickname, or identify via …",
			Prefix: &irc.Prefix{
				Name: "NickServ",
				User: "NickServ",
				Host: "services.",
			},
		},
		expectedOut: []irc.Message{
			{
				Command:  "PRIVMSG",
				Params:   []string{"NickServ"},
				Trailing: "IDENTIFY test_password",
			},
		},
	},
	{ // "pick",
		in: irc.Message{
			Command:  "PRIVMSG",
			Trailing: ":pick test",
			Params:   []string{"#testchan-1"},
			Prefix: &irc.Prefix{
				Name: "idontexist",
			},
		},
		expectedOut: []irc.Message{
			{
				Command:  "PRIVMSG",
				Params:   []string{"#testchan-1"},
				Trailing: "test",
			},
		},
	},
	{ // "google",
		in: irc.Message{
			Command:  "PRIVMSG",
			Trailing: ":g google.com",
			Params:   []string{"#testchan-1"},
			Prefix: &irc.Prefix{
				Name: "idontexist",
			},
		},
		expectedOut: []irc.Message{
			{
				Command:  "PRIVMSG",
				Params:   []string{"#testchan-1"},
				Trailing: "Google https://www.google.com/",
			},
		},
	},
	{ // "linktitle",
		in: irc.Message{
			Command:  "PRIVMSG",
			Trailing: "https://www.google.com/",
			Params:   []string{"#testchan-1"},
			Prefix: &irc.Prefix{
				Name: "idontexist",
			},
		},
		expectedOut: []irc.Message{
			{
				Command:  "PRIVMSG",
				Params:   []string{"#testchan-1"},
				Trailing: "↳ title: Google",
			},
		},
	},
	{ // "nickserv channeljoin"
		in: irc.Message{
			Command:  "NOTICE",
			Params:   []string{"gorepost"},
			Trailing: "You are now identified",
			Prefix: &irc.Prefix{
				Name: "NickServ",
				User: "NickServ",
				Host: "services.",
			},
		},
		expectedOut: []irc.Message{
			{
				Command: "JOIN",
				Params:  []string{"#securechan-1"},
			},
			{
				Command: "JOIN",
				Params:  []string{"#securechan-2"},
			},
		},
	},
	{ // "linktitle notitle",
		in: irc.Message{
			Command:  "PRIVMSG",
			Trailing: "https://www.google.com/ notitle",
			Params:   []string{"#testchan-1"},
			Prefix: &irc.Prefix{
				Name: "idontexist",
			},
		},
		expectedOut: []irc.Message{},
	},
	{ // "nickserv spoof"
		in: irc.Message{
			Command:  "NOTICE",
			Params:   []string{"gorepost"},
			Trailing: "This nickname is registered. Please choose a different nickname, or identify via …",
			Prefix: &irc.Prefix{
				Name: "NickServ",
				User: "NickServ",
				Host: "fake.",
			},
		},
		expectedOut: []irc.Message{},
	},
	{ // "nickserv other message"
		in: irc.Message{
			Command:  "NOTICE",
			Params:   []string{"gorepost"},
			Trailing: "Some other random message…",
			Prefix: &irc.Prefix{
				Name: "NickServ",
				User: "NickServ",
				Host: "services.",
			},
		},
		expectedOut: []irc.Message{},
	},
	{ // non-matching
		in: irc.Message{
			Command:  "PRIVMSG",
			Trailing: "foo bar baz",
		},
		expectedOut: []irc.Message{},
	},
}

func TestPlugins(t *testing.T) {
	var r []irc.Message
	var wg sync.WaitGroup

	// fake irc.Conn Sender replacement
	output := func(msg irc.Message) {
		wg.Done()
		r = append(r, msg)
	}

	for _, e := range eventTests {
		r = r[:0]

		wg.Add(len(e.expectedOut))

		Dispatcher(output, e.in)

		time.Sleep(3000000 * time.Nanosecond)

		wg.Wait()

		if fmt.Sprintf("%+v", r) != fmt.Sprintf("%+v", e.expectedOut) {
			t.Logf("expected: %+v\n", e.expectedOut)
			t.Logf("result: %+v\n", r)
			t.Fail()
		}
	}
}

func configLookupHelper(map[string]string) []string {
	return []string{".testconfig.json"}
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	cfg.SetFileListBuilder(configLookupHelper)
	os.Exit(m.Run())
}
