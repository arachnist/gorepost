package main

import (
	"log"
	"os"
	"path"

	"github.com/arachnist/gorepost/bot"
	. "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func FileListFuncBuilder(basedir, common string) func(map[string]string) []string {
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

	C.BuildFileList = FileListFuncBuilder(os.Args[1], "common.json")

	logfile, err := os.OpenFile(C.Lookup(context, "Logpath").(string), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("Error opening", C.Lookup(context, "Logpath").(string), "for writing, error:", err.Error())
	}
	log.SetOutput(logfile)

	networks := C.Lookup(context, "Networks").([]interface{})

	log.Println("Configured networks:", len(networks), networks)

	for i, conn := range make([]irc.Connection, len(networks)) {
		conn := conn
		log.Println("Setting up", networks[i].(string), "connection")
		conn.Setup(bot.Dispatcher, networks[i].(string))
	}
	<-exit
}
