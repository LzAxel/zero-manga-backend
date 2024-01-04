package main

import (
	"github.com/lzaxel/zero-manga-backend/internal/app"
	"github.com/lzaxel/zero-manga-backend/internal/config"
)

func main() {
	cfg := config.ReadConfig()

	app := app.New(cfg)

	app.Start()
}
