package bot

import (
	"os"
	"reflect"
	"sync"
	"testing"

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
			Prefix: &irc.Prefix{
				Name: "idontexist",
			},
		},
		expectedOut: []irc.Message{
			{
				Command:  "PRIVMSG",
				Params:   []string{"idontexist"},
				Trailing: "pingity pong",
			},
		},
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
	output := make(chan irc.Message, 1)
	quitCollector := make(chan struct{}, 1)
	var r []irc.Message
	var wg sync.WaitGroup

	for _, e := range eventTests {
		r = r[:0]

		wg.Add(len(e.expectedOut))
		go func(quit chan struct{}, input chan irc.Message) {
			for {
				select {
				case msg := <-input:
					wg.Done()
					r = append(r, msg)
				case <-quit:
				}
			}
		}(quitCollector, output)

		Dispatcher(output, e.in)

		wg.Wait()
		quitCollector <- struct{}{}

		if !reflect.DeepEqual(r, e.expectedOut) {
			t.Logf("expected: %+v\n", e.expectedOut)
			t.Logf("result: %+v\n", r)
			t.Fail()
		}
	}
}

func TestMain(m *testing.M) {
	cfg.SetFileListBuilder(configLookupHelper)
	os.Exit(m.Run())
}

func configLookupHelper(map[string]string) []string {
	return []string{".testconfig.json"}
}
