package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/vinothsparrow/scanner/config"
	"github.com/vinothsparrow/scanner/helper"
	"github.com/vinothsparrow/scanner/server"
)

func main() {
	environment := flag.String("e", "development", "")
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	flag.Usage = func() {
		log.Fatal("Usage: server -e {mode}")
	}
	flag.Parse()
	config.Init(*environment)
	c := config.GetConfig()
	helper.StartDispatcher(c.GetInt("worker.count"))

	server.Init()
}
