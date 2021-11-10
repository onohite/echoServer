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
	Services *service.Service
	cfg      *oauth2.Config
}

func NewHandler(services *service.Service, cfg *config.AuthConfig) *Handler {
	var discord = oauth2.Endpoint{
		AuthURL:  "https://discord.com/api/oauth2/authorize",
		TokenURL: "https://discord.com/api/oauth2/token",
	}

	authCfg := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  "https://2b60-109-191-92-239.ngrok.io/oauth/discord/redirect",
		Scopes:       []string{"identify", "email"},
		Endpoint:     discord,
	}
	return &Handler{services, &authCfg}
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
	//url := fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%s&redirect_uri=%s&display=%s&scope=%s&response_type=code&state=%s", clientID, redirectURI, "mobile", scopeTemp, state)
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

	//fields := strings.Join([]string{"bdate", "city", "county", "sex", "games", "photo_400_orig"}, ",")
	//url := fmt.Sprintf("https://api.vk.com/method/%s?v=5.124&fields=%s&access_token=%s", "users.get", fields, token.AccessToken)
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

	//return c.Redirect(302,"https://2b60-109-191-92-239.ngrok.io/oauth/register")
	return c.JSONBlob(200, bytes)
}
