package http

import (
	"context"
	"net"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	middle "github.com/lzaxel/zero-manga-backend/internal/handler/http/middleware"
	"github.com/lzaxel/zero-manga-backend/internal/logger"
	"github.com/lzaxel/zero-manga-backend/internal/models"
	"github.com/lzaxel/zero-manga-backend/internal/service"
)

type Config struct {
	Host string `yaml:"host" env:"HOST"`
	Port uint   `yaml:"port" env:"PORT"`
}
type Handler struct {
	jwtValidator JWTValidator
	services     *service.Services
	server       *echo.Echo
	config       Config
	logger       logger.Logger
}

func New(config Config, services *service.Services, logger logger.Logger, jwtValidator JWTValidator) *Handler {
	echo := echo.New()
	echo.HideBanner = true
	echo.HidePort = true
	handler := Handler{
		server:       echo,
		config:       config,
		services:     services,
		logger:       logger,
		jwtValidator: jwtValidator,
	}
	handler.initMiddlewares()
	handler.initRoutes()

	return &handler
}

func (h *Handler) initMiddlewares() {
	h.server.Use(
		middleware.RequestID(),
		middleware.Recover(),
		middle.Logger(h.logger),
	)
}

func (h *Handler) initRoutes() {
	api := h.server.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	auth := v1.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refreshTokens)
	}

	user := v1.Group("/user", h.Authorized())
	{
		user.GET("", h.getAllUsers, h.WithPagination(), h.RequireUserType(models.UserTypeAdmin))
		user.GET("/id/:id", h.getUserByID)
		user.GET("/self", h.getSelfUser)
	}

	manga := v1.Group("/manga")
	{
		manga.POST("", h.createManga)
		manga.GET("", h.getAllManga, h.WithPagination())
		manga.GET("/one", h.getManga)
		manga.PATCH("/:id", h.updateManga)
		manga.DELETE("/:id", h.deleteManga)
	}

	chapter := v1.Group("/chapter")
	{
		chapter.POST("", h.createChapter, h.Authorized())
		chapter.GET("/all/:manga_id", h.getChapterByManga, h.WithPagination())
		chapter.GET("/:id", h.getChapter)
	}

	tag := v1.Group("/tag")
	{
		tag.POST("", h.createTag, h.Authorized(), h.RequireUserType(models.UserTypeAdmin))
		tag.GET("", h.getAllTags)
		tag.PATCH("/:id", h.updateTag, h.Authorized(), h.RequireUserType(models.UserTypeAdmin))
		tag.DELETE("/:id", h.deleteTag, h.Authorized(), h.RequireUserType(models.UserTypeAdmin))
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
