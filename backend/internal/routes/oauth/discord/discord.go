package discord

import (
	"backend/internal/config"
	"backend/internal/service"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const oauthDiscordUrlAPI = "https://discord.com/api/v9/users/@me"

type Handler struct {
	Services  *service.Service
	cfg       *oauth2.Config
	cfgServer *config.Config
}

func NewHandler(services *service.Service, cfg *config.Config) *Handler {
	var discord = oauth2.Endpoint{
		AuthURL:  "https://discord.com/api/oauth2/authorize",
		TokenURL: "https://discord.com/api/oauth2/token",
	}

	authCfg := oauth2.Config{
		ClientID:     cfg.AuthType.DiscordConfig.ClientID,
		ClientSecret: cfg.AuthType.DiscordConfig.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/oauth/discord/redirect", cfg.Dns),
		Scopes:       []string{"identify", "email"},
		Endpoint:     discord,
	}
	return &Handler{services, &authCfg, cfg}
}

func (h Handler) Init(api *echo.Group) {
	vkGroup := api.Group("/discord")
	{
		vkGroup.GET("/login", h.Login)
		vkGroup.GET("/redirect", h.Redirect)
	}
}

func (h Handler) Login(c echo.Context) error {
	cookie, state := generateStateOauthCookie()

	c.SetCookie(&cookie)
	http.SetCookie(c.Response().Writer, &cookie)
	u := h.cfg.AuthCodeURL(state)
	return c.Redirect(302, u)
}

func generateStateOauthCookie() (http.Cookie, string) {
	var expiration = time.Now().Add(5 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "session", Value: state, Expires: expiration}

	return cookie, state
}

func (h Handler) Redirect(c echo.Context) error {
	checkCook, _ := c.Request().Cookie("session")
	log.Println(checkCook)
	cookie, err := c.Cookie("session")
	if err != nil {
		return err
	}

	stateTemp := c.QueryParam("state")
	if stateTemp[len(stateTemp)-1] == '}' {
		stateTemp = stateTemp[:len(stateTemp)-1]
	}
	if stateTemp == "" {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	} else if stateTemp != cookie.Value {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}

	token, err := h.cfg.Exchange(context.Background(), code)
	if err != nil {
		return err
	}
	log.Println(token.AccessToken)

	req, err := http.NewRequest("GET", oauthDiscordUrlAPI, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}

	return c.JSONBlob(200, bytes)
}
