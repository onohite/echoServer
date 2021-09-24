package v1

import (
	"echoTest/internal/db"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB db.DBService
}

func NewHandler(db db.DBService) *Handler {
	return &Handler{db}
}

func (h *Handler) Init(api *echo.Group) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
	}
}