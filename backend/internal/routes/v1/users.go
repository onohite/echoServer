package v1

import (
	"echoTest/internal/service/db"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initUsersRoutes(api *echo.Group) {
	users := api.Group("/users")
	{
		users.GET("/", h.GetUsers)
	}
}

type UsersResponse struct {
	Users *[]db.User `json:"users"`
}

func (h *Handler) GetUsers(c echo.Context) error {
	users, err := h.Services.DB.GetAllUsers()
	if err != nil {
		return c.JSON(500, err)
	}
	var resp UsersResponse
	resp.Users = users

	//body, _ := json.Marshal(resp)
	return c.JSON(200, resp)
}
