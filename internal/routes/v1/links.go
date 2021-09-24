package v1

import (
	"echoTest/internal/service/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) initLinksRoutes(api *echo.Group) {
	links := api.Group("/links")
	{
		links.GET("/", h.GetListLink)
		links.GET("/:id", h.GetLink)
		links.POST("/", h.AddLink)
	}
}

type Links struct {
	Urls *[]db.Link `json:"Links"`
}

func (h Handler) GetListLink(c echo.Context) error {
	links, err := h.Services.DB.GetAllLinks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, Links{Urls: links})
}

func (h Handler) GetLink(c echo.Context) error {
	id := c.Param("id")

	link, err := h.Services.DB.GetLinkById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, link)
}

func (h Handler) AddLink(c echo.Context) error {
	var link db.Link
	if err := c.Bind(&link); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.Services.DB.AddLink(link); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
