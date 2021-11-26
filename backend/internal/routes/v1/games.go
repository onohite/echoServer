package v1

import (
	"backend/internal/routes/middleware"
	"backend/internal/service/db"
	"backend/internal/service/graph"
	"backend/internal/utils"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) initGamesRoutes(api *echo.Group) {
	games := api.Group("/games")
	{
		games.GET("", h.GetListGames)
		games.GET("/rank", h.GetGameRank)
		games.GET("/profile/id/", h.GetGameProfile, middleware.JWTMiddleware(h.Services))
		games.POST("/profile", h.CreateGameProfile, middleware.JWTMiddleware(h.Services))
		games.PATCH("/profile", h.PathGameProfile, middleware.JWTMiddleware(h.Services))
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

func (h *Handler) GetGameProfile(c echo.Context) error {
	profileID := c.QueryParam("id")
	if profileID == "" {
		return c.JSON(http.StatusBadRequest, "profile id was empty")
	}
	id, err := strconv.Atoi(profileID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "validation error")
	}
	//TODO добавить кеширование
	uid, err := h.Services.DB.FindGameProfile(id)
	if err != nil {
		log.Errorf("found err in FindGameProfile %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	profile, err := h.Services.Graph.GetProfile(uid)
	if err != nil {
		log.Errorf("found err in GetProfile %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, profile)
}

type CrGameProfile struct {
	GameID      int    `json:"game_id"`
	RankID      int    `json:"rank_id"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
}

func (h *Handler) CreateGameProfile(c echo.Context) error {
	sess, _ := session.Get("session", c)
	uuid, ok := sess.Values["uuid"].(string)
	if uuid == "" || !ok {
		return c.JSON(http.StatusUnauthorized, "access denied")
	}

	var profile graph.GameProfile
	var jsonReq CrGameProfile
	user, err := h.Services.DB.GetUser(uuid)
	if err != nil {
		log.Errorf("cant found user createGameProfile %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if user == nil {
		log.Error("user not found")
		return c.JSON(http.StatusUnauthorized, err)
	}

	if err = c.Bind(&jsonReq); err != nil {
		log.Errorf("cant bind jsonReq createGameProfile %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	profile = makeGameProfile(*user, jsonReq, uuid)
	//TODO добавить кеш
	uID, err := h.Services.Graph.SetProfile(profile)
	if err != nil {
		log.Errorf("cant create graph profile SetProfile %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	id, err := h.Services.DB.CreateGameProfile(uID)
	if err != nil {
		log.Errorf("cant create db profile createGameProfile %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, id)
}

type UpdateGameProfile struct {
	ID          int    `json:"id"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
}

func (h *Handler) PathGameProfile(c echo.Context) error {
	sess, _ := session.Get("session", c)
	uuid, ok := sess.Values["uuid"].(string)
	if uuid == "" || !ok {
		return c.JSON(http.StatusUnauthorized, "access denied")
	}

	var jsonReq UpdateGameProfile
	user, err := h.Services.DB.GetUser(uuid)
	if err != nil {
		log.Errorf("cant found user createGameProfile %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if user == nil {
		log.Error("user not found")
		return c.JSON(http.StatusUnauthorized, err)
	}

	if err = c.Bind(&jsonReq); err != nil {
		log.Errorf("cant bind jsonReq createGameProfile %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if jsonReq.ID != 0 {
		log.Errorf("profile id was empty")
		return c.NoContent(http.StatusBadRequest)
	}

	//TODO убрать из запросов uuid тк можно привязатся к uid graph
	var errors []utils.ErrorStruct
	if jsonReq.Contact != "" {
		err := h.Services.DB.UpdateGameProfileContact(jsonReq.Contact, uuid, jsonReq.ID)
		if err != nil {
			errors = append(errors, utils.ErrorStruct{Message: "Contact update error", Detail: err.Error()})
		}
	}

	if jsonReq.Description != "" {
		err := h.Services.DB.UpdateGameProfileDescription(jsonReq.Description, uuid, jsonReq.ID)
		if err != nil {
			errors = append(errors, utils.ErrorStruct{Message: "Description update error", Detail: err.Error()})
		}
	}

	if len(errors) != 0 {
		return c.JSON(http.StatusBadRequest, errors)
	}

	return c.JSON(http.StatusOK, "updated")
}

func makeGameProfile(user db.PublicUser, req CrGameProfile, uuid string) graph.GameProfile {
	var profile graph.GameProfile
	profile.UserID = uuid
	profile.Sex = user.Sex
	birthdate, _ := time.Parse(utils.DataLayout, user.Bdate)
	profile.Age = toAge(birthdate)
	profile.GameID = req.GameID
	profile.RankID = req.RankID
	profile.Contact = req.Contact
	profile.Description = req.Description
	return profile
}

func toAge(birthdate time.Time) int {
	today := time.Now().UTC()
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	by, bm, bd := birthdate.Date()
	birthdate = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	if today.Before(birthdate) {
		return 0
	}
	age := ty - by
	anniversary := birthdate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}
	return age
}
