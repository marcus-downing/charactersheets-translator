package main

import (
	"./config"
	"./server"
)

func main() {
	server.RunTranslator(config.Config.Server.Host(), config.Config.Debug)
}
