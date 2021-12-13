package v1

import (
	"backend/internal/routes/middleware"
	"backend/internal/service/db"
	"backend/internal/utils"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func (h *Handler) initRemindRoutes(api *echo.Group) {
	remind := api.Group("/remind")
	remind.Use(middleware.JWTMiddleware(h.Services))
	{
		remind.POST("", h.CreateRemind)
		remind.PATCH("", h.UpdateRemind)
		remind.GET("", h.GetReminds)
	}
}

func (h *Handler) CreateRemind(c echo.Context) error {
	sess, _ := session.Get("session", c)
	uuid, ok := sess.Values["uuid"].(string)
	if uuid == "" || !ok {
		return c.JSON(http.StatusUnauthorized, "access denied")
	}

	user, err := h.Services.DB.GetUser(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "cant found user with this uuid")
	}

	var rem db.Remind
	if err := c.Bind(&rem); err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("try with valid data, detail : %v", err))
	}

	rem.From = user.Email

	if err := h.Services.Queue.SetRemind(rem); err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("cant set request to queue : %v", err))
	}

	return c.JSON(http.StatusOK, "ok")
}

type RemindUpdateData struct {
	Id      string
	Where   string `json:"where"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

func (h *Handler) UpdateRemind(c echo.Context) error {
	sess, _ := session.Get("session", c)
	uuid, ok := sess.Values["uuid"].(string)
	if uuid == "" || !ok {
		return c.JSON(http.StatusUnauthorized, "access denied")
	}

	var data RemindUpdateData
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("try with valid data, detail : %v", err))
	}

	var errors []utils.ErrorStruct
	if data.Where != "" {
		err = h.Services.DB.UpdateRemindTo(data.Id, data.Where)
		if err != nil {
			errors = append(errors, utils.ErrorStruct{Message: "where update error", Detail: err.Error()})
		}

	}

	if data.Message != "" {
		err = h.Services.DB.UpdateRemindMessage(data.Id, data.Message)
		if err != nil {
			errors = append(errors, utils.ErrorStruct{Message: "message update error", Detail: err.Error()})
		}
	}

	if data.Date != "" {
		err = h.Services.DB.UpdateRemindDate(data.Id, data.Date)
		if err != nil {
			errors = append(errors, utils.ErrorStruct{Message: "date update error", Detail: err.Error()})
		}
	}

	if len(errors) != 0 {
		return c.JSON(http.StatusBadRequest, errors)
	}

	return c.JSON(http.StatusOK, "updated")
}

func (h *Handler) GetReminds(c echo.Context) error {
	sess, _ := session.Get("session", c)
	uuid, ok := sess.Values["uuid"].(string)
	if uuid == "" || !ok {
		return c.JSON(http.StatusUnauthorized, "access denied")
	}

	user, err := h.Services.DB.GetUser(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "cant found user with this uuid")
	}

	list, err := h.Services.DB.GetListReminds(user.Email)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("cant found reminds with this user :%v", err))
	}

	return c.JSON(http.StatusOK, list)
}
