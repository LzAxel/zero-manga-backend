package main

import (
	"github.com/lzaxel/zero-manga-backend/internal/app"
	"github.com/lzaxel/zero-manga-backend/internal/config"
)

// @title           ZeroManga API
// @version         0.1
// @description     This is an API for ZeroManga.

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

func main() {
	cfg := config.ReadConfig()

	app := app.New(cfg)

	app.Start()
}
