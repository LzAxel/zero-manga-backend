package http

import (
	"context"
	"net"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lzaxel/zero-manga-backend/internal/logger"
	"github.com/lzaxel/zero-manga-backend/internal/service"
)

type Config struct {
	Host string `yaml:"host" env:"HOST"`
	Port uint   `yaml:"port" env:"PORT"`
}
type Handler struct {
	services *service.Services
	server   *echo.Echo
	config   Config
	logger   logger.Logger
}

func New(config Config, services *service.Services, logger logger.Logger) *Handler {
	echo := echo.New()
	echo.HideBanner = true
	echo.HidePort = true
	handler := Handler{
		server:   echo,
		config:   config,
		services: services,
		logger:   logger,
	}
	handler.initMiddlewares()
	handler.initRoutes()

	return &handler
}

func (h *Handler) initMiddlewares() {
	h.server.Use(
		middleware.Logger(),
		middleware.Recover(),
	)
}

func (h *Handler) initRoutes() {
	api := h.server.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	v1.Group("/auth")
	{

	}

	v1.Group("/user")
	{

	}

	v1.Group("/manga")
	{

	}

	v1.Group("/chapter")
	{

	}
}

func (h *Handler) Stop(ctx context.Context) error {
	h.logger.Infof("shutting down server")
	return h.server.Shutdown(ctx)
}

func (h *Handler) Start() error {
	h.logger.Infof("starting server on %s:%d", h.config.Host, h.config.Port)
	return h.server.Start(net.JoinHostPort(h.config.Host, strconv.Itoa(int(h.config.Port))))
}
