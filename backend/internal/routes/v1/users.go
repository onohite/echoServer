package v1

import (
	"backend/internal/routes/middleware"
	"backend/internal/utils"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserUpdateData struct {
	Username   string `json:"username"`
	Email      string `json:"email" validate:"email"`
	AvatarLink string `json:"avatar_link"`
	Sex        int    `json:"sex"`
	Bdate      string `json:"bdate"`
}

func (h *Handler) initUsersRoutes(api *echo.Group) {
	users := api.Group("/user")
	users.Use(middleware.JWTMiddleware(h.Services))
	{
		users.GET("", h.GetUser)
		users.PATCH("", h.UpdateUser)
	}
}

func (h *Handler) GetUser(c echo.Context) error {
	sess, _ := session.Get("session", c)
	uuid, ok := sess.Values["uuid"].(string)
	if uuid == "" || !ok {
		return c.JSON(http.StatusUnauthorized, "access denied")
	}

	user, err := h.Services.DB.GetUser(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "cant found user with this uuid")
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	sess, _ := session.Get("session", c)
	uuid, ok := sess.Values["uuid"].(string)
	if uuid == "" || !ok {
		return c.JSON(http.StatusUnauthorized, "access denied")
	}

	var data UserUpdateData
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("try with valid data, detail : %v", err))
	}

	err = c.Validate(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("validation error %v", err))
	}

	var errors []utils.ErrorStruct
	if data.Sex != 0 {
		err = h.Services.DB.UpdateUserSex(data.Sex, uuid)
		if err != nil {
			errors = append(errors, utils.ErrorStruct{Message: "sex update error", Detail: err.Error()})
		}

	}

	if data.Bdate != "" {
		err = h.Services.DB.UpdateUserBdate(data.Bdate, uuid)
		if err != nil {
			errors = append(errors, utils.ErrorStruct{Message: "bdate update error", Detail: err.Error()})
		}
	}

	if data.AvatarLink != "" {
		err = h.Services.DB.UpdateUserAvatar(data.AvatarLink, uuid)
		if err != nil {
			errors = append(errors, utils.ErrorStruct{Message: "avatar update error", Detail: err.Error()})
		}
	}

	if data.Email != "" {
		err = h.Services.DB.UpdateUserEmail(data.Email, uuid)
		if err != nil {
			errors = append(errors, utils.ErrorStruct{Message: "email update error", Detail: err.Error()})
		}
	}

	if data.Username != "" {
		err = h.Services.DB.UpdateUserUserName(data.Username, uuid)
		if err != nil {
			errors = append(errors, utils.ErrorStruct{Message: "username update error", Detail: err.Error()})
		}
	}

	if len(errors) != 0 {
		return c.JSON(http.StatusBadRequest, errors)
	}

	return c.JSON(http.StatusOK, "updated")
}
