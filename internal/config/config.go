// package config loads the config struct by the function MustLoad()
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	ENV           string        `yaml:"env" env-default:"local"`
	DSN           string        `yaml:"storage_dsn" env-required:"true"`
	HTTP_Server   http_server   `yaml:"http_server"`
	AccessJwtTTL  time.Duration `yaml:"jwt_access_ttl" env-default:"1m"`
	RefreshJwtTTL time.Duration `yaml:"jwt_refresh_ttl" env-default:"1m"`
	SigningKey    string        `yaml:"signing_key"`
}

type http_server struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" || configPath == "local" {
		log.Fatal("config path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config path: %s is not exists", configPath)
	}

	var config Config

	log.Println("loading config from: ", configPath)
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config: %s, err: %s", configPath, err)
	}
	return &config
}
