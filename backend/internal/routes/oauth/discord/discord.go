package discord

import (
	"backend/internal/config"
	"backend/internal/service"
	"backend/internal/service/db"
	"backend/internal/utils"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"log"
	"net/http"
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
	cookie, state := generateStateOauthCookie(&c)
	if err := cookie.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	u := h.cfg.AuthCodeURL(state)
	return c.Redirect(302, u)
}

func generateStateOauthCookie(c *echo.Context) (*sessions.Session, string) {
	sess, _ := session.Get("session", *c)
	sess.Options = &sessions.Options{
		Path:     "/oauth/discord/redirect",
		MaxAge:   60 * 60 * 5,
		HttpOnly: false,
	}

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	sess.Values["state"] = state
	return sess, state
}

func (h Handler) Redirect(c echo.Context) error {
	sess, _ := session.Get("session", c)
	state, ok := sess.Values["state"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, "state error unauthorized")
	}

	stateTemp := c.QueryParam("state")
	if stateTemp[len(stateTemp)-1] == '}' {
		stateTemp = stateTemp[:len(stateTemp)-1]
	}
	if stateTemp == "" {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	} else if stateTemp != state {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}

	sess.Options.MaxAge = -1
	err := sess.Save(c.Request(), c.Response())
	if err != nil {
		log.Print("cant delete session")
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

	resp, err := resty.New().R().SetAuthToken(token.AccessToken).Get(oauthDiscordUrlAPI)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}

	var discordResp DiscordResp
	if err := json.Unmarshal(resp.Body(), &discordResp); err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка маршалинга google")
	}

	uniqueKey := utils.GenerateKey(discordResp.Username, discordResp.ID)
	log.Println(uniqueKey)
	user := db.User{
		Name:       discordResp.Username,
		AvatarLink: fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s", discordResp.ID, discordResp.Avatar),
		Email:      discordResp.Email,
		Unique:     uniqueKey,
	}

	uuid, err := h.Services.DB.AddUser(user)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Println(uuid)

	resp, err = resty.New().R().SetBody(uuid).Post(fmt.Sprintf("%s/oauth/login", h.cfgServer.Dns))
	if err != nil {
		log.Print(err)
		return err
	}

	return c.JSONBlob(200, resp.Body())
}
