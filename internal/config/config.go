package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lzaxel/zero-manga-backend/internal/handler/http"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
)

const (
	configPath = "configs/dev.yaml"
)

type AppConfig struct {
	IsDev    bool   `yaml:"isDev" env:"IS_DEV"`
	LogLevel string `yaml:"logLevel" env:"LOG_LEVEL"`
}

type Config struct {
	Postgresql postgresql.Config `yaml:"postgres"`
	Server     http.Config       `yaml:"server"`
	App        AppConfig         `yaml:"app"`
}

var (
	config Config
	once   sync.Once
)

func ReadConfig() Config {
	once.Do(func() {
		if err := cleanenv.ReadConfig(configPath, &config); err != nil {
			panic(err)
		}
	})

	return config
}
