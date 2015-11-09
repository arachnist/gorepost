package main

import (
	"log"
	"os"

	"github.com/arachnist/gorepost/bot"
	. "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func main() {
	var exit chan struct{}
	var context Context

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

	C.ConfigDir = os.Args[1]

	logfile, err := os.OpenFile(C.Lookup(context, "Logpath").(string), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Error opening", C.Lookup(context, "Logpath").(string), "for writing, error:", err.Error())
	}
	log.SetOutput(logfile)

	networks := C.Lookup(context, "Networks").([]string)
	connections := make([]irc.Connection, len(networks))
	for i, conn := range connections {
		context.Network = networks[i]

		connections[i].Setup(bot.Dispatcher, networks[i],
			C.Lookup(context, "Servers").([]string),
			C.Lookup(context, "Nick").(string),
			C.Lookup(context, "User").(string),
			C.Lookup(context, "RealName").(string))
	}
	<-exit
}
