package oauth

import (
	"backend/internal/service"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	clientID     = "7995872"
	clientSecret = "0GJJi1vWUSq2UEYXdyMV"
	redirectURI  = "https://2b60-109-191-92-239.ngrok.io/oauth/redirect"
	scope        = "account"
	state        = "12345"
)

//conf := &oauth2.Config{
//ClientID:     os.Getenv("CLIENT_ID"),
//ClientSecret: os.Getenv("CLIENT_SECRET"),
//RedirectURL:  os.Getenv("REDIRECT_URL"),
//Scopes:       []string{},
//Endpoint:     vkAuth.Endpoint,
//}

type Handler struct {
	Services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services}
}

func (h *Handler) Init(api *echo.Group) {
	api.GET("", h.index)
	api.GET("/redirect", h.redirect)
	api.GET("/login", h.login)
}

func (h Handler) index(c echo.Context) error {
	title := "Выберите способ авторизации"
	err := c.Render(200, "default_auth.html", title)
	log.Error(err)
	return err
}

func (h Handler) login(c echo.Context) error {
	scopeTemp := strings.Join([]string{"account"}, "+")
	url := fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%s&redirect_uri=%s&display=%s&scope=%s&response_type=code&state=%s", clientID, redirectURI, "mobile", scopeTemp, state)
	return c.Redirect(302, url)
}

func (h Handler) redirect(c echo.Context) error {
	stateTemp := c.QueryParam("state")
	if stateTemp[len(stateTemp)-1] == '}' {
		stateTemp = stateTemp[:len(stateTemp)-1]
	}
	if stateTemp == "" {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	} else if stateTemp != state {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}
	url := fmt.Sprintf("https://oauth.vk.com/access_token?grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		code, redirectURI, clientID, clientSecret)
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}
	defer resp.Body.Close()
	token := struct {
		AccessToken string `json:"access_token"`
	}{}
	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &token)

	fields := strings.Join([]string{"bdate", "city", "county", "sex", "games", "photo_400_orig"}, ",")
	url = fmt.Sprintf("https://api.vk.com/method/%s?v=5.124&fields=%s&access_token=%s", "users.get", fields, token.AccessToken)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}
	defer resp.Body.Close()
	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}

	return c.Render(200, "auth.html", string(bytes))
}
