package app

import (
	"context"

	"github.com/lzaxel/zero-manga-backend/internal/config"
	"github.com/lzaxel/zero-manga-backend/internal/filestorage"
	"github.com/lzaxel/zero-manga-backend/internal/handler/http"
	"github.com/lzaxel/zero-manga-backend/internal/jwt"
	"github.com/lzaxel/zero-manga-backend/internal/logger"
	"github.com/lzaxel/zero-manga-backend/internal/repository"
	"github.com/lzaxel/zero-manga-backend/internal/repository/postgresql"
	"github.com/lzaxel/zero-manga-backend/internal/service"
	"github.com/lzaxel/zero-manga-backend/pkg/clock"
	"github.com/lzaxel/zero-manga-backend/pkg/uuid"
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

	if config.App.IsTesting {
		clock.InitClock(true)
		uuid.InitUUID(true)
	}

	logger.Infof("connecting to postgresql on %s:%d", config.Postgresql.Host, config.Postgresql.Port)
	psql, err := postgresql.New(ctx, config.Postgresql)
	if err != nil {
		logger.Fatalf("connect to postgresql: %s", err)
	}

	err = postgresql.Migrate(ctx, config.Postgresql)
	if err != nil {
		logger.Warnf("migrate database: %s", err)
	}
	fileStorage := filestorage.NewLocalFileStorage(config.FileStorage, logger)
	go fileStorage.Serve()

	jwt := jwt.New(config.JWT)
	repository := repository.New(ctx, psql, logger)
	services := service.New(ctx, repository, jwt, fileStorage)
	handler := http.New(config.Server, services, logger, jwt)

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
