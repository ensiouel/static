package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL"`
	GRPC     GRPC
	Static   Static
}

type GRPC struct {
	Addr string `env:"GRPC_ADDR"`
}

type Static struct {
	StorageRoot string `env:"STATIC_STORAGE_ROOT"`
	MaxFileSize int    `env:"STATIC_MAX_FILE_SIZE"`
}

func New() Config {
	var config Config
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
