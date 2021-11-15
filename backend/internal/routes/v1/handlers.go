package v1

import (
	"backend/internal/service"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services}
}

func (h *Handler) Init(api *echo.Group) {
	v1 := api.Group("/v1")
	{
		h.initGamesRoutes(v1)
		h.initUsersRoutes(v1)
	}
}
