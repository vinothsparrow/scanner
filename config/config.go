package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Error("error on parsing configuration file")
	}
	config.SetDefault("http.server", getenv("SCANNER_SERVER", "0.0.0.0"))
	config.SetDefault("http.port", getenv("SCANNER_PORT", "8000"))
	config.SetDefault("auth.key", getenv("SCANNER_AUTH_KEY", ""))
	key := config.GetString("auth.key")
	if len(key) == 0 {
		uuidObj, _ := uuid.NewRandom()
		id := uuidObj.String()
		config.SetDefault("auth.key", id)
		log.Println("API Key is", id)
	}
	config.SetDefault("http.env", getenv("SCANNER_ENV", "release"))
	workers, _ := strconv.Atoi(getenv("SCANNER_WORKER_COUNT", "2"))
	config.SetDefault("worker.count", workers)
}

func relativePath(basedir string, path *string) {
	p := *path
	if len(p) > 0 && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

func GetConfig() *viper.Viper {
	return config
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
