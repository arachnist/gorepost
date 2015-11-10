package main

import (
	"log"
	"os"
	"path"

	"github.com/arachnist/gorepost/bot"
	. "github.com/arachnist/gorepost/config"
	"github.com/arachnist/gorepost/irc"
)

func main() {
	var exit chan struct{}
	context := make(map[string]string)
	var networks []string

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

	C.BuildFileList = func(c map[string]string) []string {
		var r []string

		if c["Network"] != "" {
			if c["Source"] != "" {
				if c["Target"] != "" {
					r = append(r, path.Join(os.Args[1], c["Network"], c["Source"], c["Target"]+".json"))
				}
				r = append(r, path.Join(os.Args[1], c["Network"], c["Source"]+".json"))
			}
			r = append(r, path.Join(os.Args[1], c["Network"]+".json"))
		}

		return append(r, path.Join(os.Args[1], "common.json"))
	}

	logfile, err := os.OpenFile(C.Lookup(context, "Logpath").(string), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Error opening", C.Lookup(context, "Logpath").(string), "for writing, error:", err.Error())
	}
	log.SetOutput(logfile)

	rawNetworks := C.Lookup(context, "Networks").([]interface{})
	for _, n := range rawNetworks {
		networks = append(networks, n.(string))
	}

	log.Println("Configured networks:", len(networks), networks)

	connections := make([]irc.Connection, len(networks))
	for i, _ := range connections {
		var servers []string
		context["Network"] = networks[i]

		rawServers := C.Lookup(context, "Servers").([]interface{})
		log.Println("Rawservers:", rawServers)
		for _, n := range rawServers {
			servers = append(servers, n.(string))
		}
		log.Println(context["Network"], "Configured servers", len(servers), servers)
		connections[i].Setup(bot.Dispatcher, networks[i],
			servers,
			C.Lookup(context, "Nick").(string),
			C.Lookup(context, "User").(string),
			C.Lookup(context, "RealName").(string))
	}
	<-exit
}
