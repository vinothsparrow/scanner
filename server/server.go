package server

import (
	"github.com/vinothsparrow/scanner/config"
)

func Init() {
	config := config.GetConfig()

	r := NewRouter()
	r.Run(config.GetString("http.server") + ":" + config.GetString("http.port"))
}
