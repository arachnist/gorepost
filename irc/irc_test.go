// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package irc

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	cfg "github.com/arachnist/gorepost/config"
)

var expectedOutput = []Message{
	{
		Command:  "NICK",
		Trailing: "gorepost",
	},
	{
		Command:  "USER",
		Params:   []string{"repost", "0", "*"},
		Trailing: "https://github.com/arachnist/gorepost",
	},
}

var input = []Message{
	{
		Command: "001",
		Params:  []string{"gorepost"},
	},
}

var actualOutput []Message
var actualInput []Message

func connWriter(c net.Conn) {
}

func connReader(c net.Conn) {
}

func fakeServer(t *testing.T) {
	ln, err := net.Listen("tcp", ":36667")
	if err != nil {
		t.Error("fakeServer can't start listening")
	}
	// twice, to test reconnects
	for range []int{0, 2} {
		var wg sync.WaitGroup

		conn, err := ln.Accept()
		if err != nil {
			t.Error("error accepting connection")
		}

		// writer
		go func(c net.Conn) {
			writer := bufio.NewWriter(c)
			wg.Add(len(input))
			for _, msg := range input {
				writer.WriteString(msg.String() + endline)
				writer.Flush()
				wg.Done()
			}
		}(conn)

		// reader
		go func(c net.Conn) {
			reader := bufio.NewReader(c)
			wg.Add(len(expectedOutput))
			for range expectedOutput {
				raw, err := reader.ReadString(delim)
				if err != nil {
					t.Log("Failed reading message from client:", err)
					t.Fail()
				}

				msg, err := ParseMessage(raw)
				if err != nil {
					t.Log("Failed parsing message from client:", raw)
					t.Fail()
				}

				actualOutput = append(actualOutput, *msg)
				wg.Done()
			}
		}(conn)

		wg.Wait()

		time.Sleep(1 * time.Second)

		conn.Close()
	}
}

func TestSetup(t *testing.T) {
	go fakeServer(t)

	var conn Connection
	conn.Setup(fakeDispatcher, "TestNet")

	time.Sleep(2 * time.Second)

	conn.Quit <- struct{}{}

	// since we tested a reconnect, we should expect actual results to be
	// multipled
	actualExpectedOutput := append(expectedOutput, expectedOutput...)
	actualExpectedInput := append(input, input...)

	if fmt.Sprintf("%+v", actualExpectedOutput) != fmt.Sprintf("%+v", actualOutput) {
		t.Log("Expected output does not match actual output")
		t.Logf("expected: %+v\n", actualExpectedOutput)
		t.Logf("actual  : %+v\n", actualOutput)
		t.Fail()
	}

	if fmt.Sprintf("%+v", actualExpectedInput) != fmt.Sprintf("%+v", actualInput) {
		t.Log("Expected input does not match actual input")
		t.Logf("expected: %+v\n", actualExpectedInput)
		t.Logf("actual  : %+v\n", actualInput)
		t.Fail()
	}
}

func fakeDispatcher(output chan Message, input Message) {
	// nullify Context as it isn't transmitted over the wire
	input.Context = make(map[string]string)
	actualInput = append(actualInput, input)
}

func configLookupHelper(map[string]string) []string {
	return []string{".testconfig.json"}
}

func TestMain(m *testing.M) {
	cfg.SetFileListBuilder(configLookupHelper)
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
