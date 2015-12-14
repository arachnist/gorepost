// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Gorepost is an overengineered IRC bot that i use for learning Go.
package main

import (
	"log"
	"os"
	"path"

	"github.com/arachnist/dyncfg"
	"github.com/arachnist/gorepost/bot"
	"github.com/arachnist/gorepost/irc"
)

func fileListFuncBuilder(basedir, common string) func(map[string]string) []string {
	return func(c map[string]string) []string {
		var r []string

		if c["Network"] != "" {
			if c["Target"] != "" {
				if c["Source"] != "" {
					r = append(r, path.Join(basedir, c["Network"], c["Target"], c["Source"]+".json"))
					r = append(r, path.Join(basedir, c["Network"], c["Source"]+".json"))
				}
				r = append(r, path.Join(basedir, c["Network"], c["Target"]+".json"))
			}
			r = append(r, path.Join(basedir, c["Network"]+".json"))
		}

		return append(r, path.Join(basedir, common))
	}
}

func main() {
	var exit chan struct{}
	context := make(map[string]string)

	if len(os.Args) < 2 {
		log.Fatalln("Usage:", os.Args[0], "<configuration directory>")
	}

	d, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatalln("Error reading configuration from", os.Args[1], "error:", err.Error())
	}
	if !d.IsDir() {
		log.Fatalln("Not a directory:", os.Args[1])
	}

	cfg := dyncfg.New(fileListFuncBuilder(os.Args[1], "common.json"))

	logfile, err := os.OpenFile(cfg.LookupString(context, "Logpath"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("Error opening", cfg.LookupString(context, "Logpath"), "for writing, error:", err.Error())
	}
	log.SetOutput(logfile)

	networks := cfg.LookupStringSlice(context, "Networks")

	log.Println("Configured networks:", len(networks), networks)

	bot.Initialize(cfg)
	for i, conn := range make([]irc.Connection, len(networks)) {
		conn := conn
		log.Println("Setting up", networks[i], "connection")
		conn.Setup(bot.Dispatcher, networks[i], cfg)
	}
	<-exit
}
