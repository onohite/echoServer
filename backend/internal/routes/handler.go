package routes

import (
	"backend/docs/docs"
	"backend/internal/config"
	"backend/internal/routes/oauth"
	v1 "backend/internal/routes/v1"
	"backend/internal/service"
	"fmt"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

type Handler struct {
	Services *service.Service
	AuthType *config.AuthType
}

func NewHandler(service *service.Service, cfg *config.AuthType) *Handler {
	return &Handler{service, cfg}
}

func (h *Handler) Init(cfg *config.Config) *echo.Echo {
	// Init echo handler
	router := echo.New()
	router.Renderer = echoview.Default()

	// Init middleware
	router.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}, bytes_in=${bytes_in}, bytes_out=${bytes_out}\n",
			Output: router.Logger.Output()}),
		middleware.Recover())

	// Init log level
	router.Debug = cfg.ServerMode != config.Dev

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	if cfg.ServerMode != config.Dev {
		docs.SwaggerInfo.Host = cfg.Host
	}

	if cfg.ServerMode != config.Prod {
		router.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// Init router
	router.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)
	h.oauthAPI(router)

	return router
}

func (h *Handler) initAPI(router *echo.Echo) {
	handlerV1 := v1.NewHandler(h.Services)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}

func (h *Handler) oauthAPI(router *echo.Echo) {
	handlerOauth := oauth.NewHandler(h.Services, h.AuthType)
	api := router.Group("/oauth")
	{
		handlerOauth.Init(api)
	}
}
