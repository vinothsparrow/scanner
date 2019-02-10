package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/vinothsparrow/scanner/config"
	"github.com/vinothsparrow/scanner/helper"
	"github.com/vinothsparrow/scanner/server"
)

// @title .git folder scanner API
// @version 1.0
// @description This is a scanner server to scan .git folder in the website.

// @host
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in query
// @name api_key

// @contact.name API Support
// @contact.url https://github.com/vinothsparrow/scanner
// @contact.email vinothsparrow@live.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://swagger.io/terms/
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
