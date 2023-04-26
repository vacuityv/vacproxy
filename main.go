package main

import (
	"flag"
	"log"
	"vacproxy/service"
)

func main() {
	console := flag.Bool("console", false, "run with console")

	flag.Parse()

	configPath, err := service.GetConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	config, inMatch, outMatch := service.InitConfig(configPath)
	if *console {
		server := service.NewServer(config, inMatch, outMatch)
		service.StartProxy(server, true)
	} else {
		service.ProxyService(config, inMatch, outMatch)
	}

}
