package oauth

import (
	"backend/internal/config"
	"backend/internal/routes/oauth/discord"
	"backend/internal/routes/oauth/google"
	"backend/internal/routes/oauth/vk"
	"backend/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Handler struct {
	Services *service.Service
	AuthType *config.AuthType
}

func NewHandler(services *service.Service, cfg *config.AuthType) *Handler {
	return &Handler{services, cfg}
}

func (h *Handler) Init(api *echo.Group) {
	api.GET("", h.index)
	h.InitVK(api)
	h.InitGoogle(api)
	h.InitDiscord(api)
}

func (h Handler) index(c echo.Context) error {
	title := "Выберите способ авторизации"
	err := c.Render(200, "default_auth.html", title)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (h *Handler) InitVK(api *echo.Group) {
	vkGroup := vk.NewHandler(h.Services, &h.AuthType.VKconfig)
	vkGroup.Init(api)
}

func (h *Handler) InitGoogle(api *echo.Group) {
	gglGroup := google.NewHandler(h.Services, &h.AuthType.GoogleConfig)
	gglGroup.Init(api)
}

func (h *Handler) InitDiscord(api *echo.Group) {
	dscGroup := discord.NewHandler(h.Services, &h.AuthType.DiscordConfig)
	dscGroup.Init(api)
}
