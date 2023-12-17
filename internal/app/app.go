package app

import (
	"context"

	"github.com/lzaxel/zero-manga-backend/internal/config"
	"github.com/lzaxel/zero-manga-backend/internal/handler/http"
	"github.com/lzaxel/zero-manga-backend/internal/logger"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
	"github.com/lzaxel/zero-manga-backend/internal/service"
)

type App struct {
	services   *service.Services
	repository *repository.Repository
	hanlder    *http.Handler
	logger     logger.Logger
}

func New(config config.Config) *App {
	ctx := context.Background()
	logger := logger.NewLogrusLogger(config.App.LogLevel, config.App.IsDev)

	logger.Infof("config loaded")
	logger.Infof("connecting to postgresql on %s:%d", config.Postgresql.Host, config.Postgresql.Port)
	psql, err := postgresql.New(ctx, config.Postgresql)
	if err != nil {
		panic(err)
	}
	repository := repository.New(ctx, psql, logger)
	services := service.New(ctx, repository)
	handler := http.New(config.Server, services, logger)

	return &App{
		services:   services,
		repository: repository,
		hanlder:    handler,
		logger:     logger,
	}
}

func (app *App) Start() error {
	if err := app.hanlder.Start(); err != nil {
		app.logger.Errorf("failed to start app: %s", err)
		return err
	}

	return nil
}

func (app *App) Shutdown(context.Context) {
	if err := app.hanlder.Stop(context.Background()); err != nil {
		panic(err)
	}
}
