package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// ConfigSpec contains mapping of environment variable names to config values
type ConfigSpec struct {
	Logger  *logrus.Logger
	NoColor bool
	ApiUrl  string `envconfig:"INFRACOST_API_URL"  required:"true"  default:"https://pricing.infracost.io"`
}

func (c *ConfigSpec) SetLogger(logger *logrus.Logger) {
	c.Logger = logger
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// loadConfig loads the config struct from environment variables
func loadConfig() *ConfigSpec {
	var config ConfigSpec
	var err error

	config.NoColor = false

	envLocalPath := filepath.Join(rootDir(), ".env.local")
	if fileExists(envLocalPath) {
		err = godotenv.Load(envLocalPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	if fileExists(".env") {
		err = godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
	}

	err = envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}

var Config = loadConfig()
