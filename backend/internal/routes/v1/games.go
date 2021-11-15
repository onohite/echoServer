package v1

import (
	"backend/internal/service/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) initGamesRoutes(api *echo.Group) {
	users := api.Group("/games")
	{
		users.GET("", h.GetListGames)
		users.GET("/rank", h.GetGameRank)
	}
}

func (h *Handler) GetListGames(c echo.Context) error {
	var games *db.Games
	err := h.Services.Cache.GetData("/games", &games)
	if err != nil {
		games, err = h.Services.DB.GetAllGames()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		err = h.Services.Cache.SetData("/games", games, time.Hour)
		if err != nil {
			log.Error(err)
		}
	}

	return c.JSON(http.StatusOK, games)
}

//TODO добавить кеширование
func (h *Handler) GetGameRank(c echo.Context) error {
	gameID := c.QueryParam("gameID")
	if gameID == "" {
		return c.JSON(http.StatusBadRequest, "gameID was empty")
	}
	id, err := strconv.Atoi(gameID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "validation error")
	}
	ranks, err := h.Services.DB.GetGameRanks(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ranks)
}
