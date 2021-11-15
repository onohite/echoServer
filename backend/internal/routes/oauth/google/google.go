package google

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
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type Handler struct {
	Services  *service.Service
	cfg       *oauth2.Config
	cfgServer *config.Config
}

func NewHandler(services *service.Service, cfg *config.Config) *Handler {

	authCfg := oauth2.Config{
		ClientID:     cfg.AuthType.GoogleConfig.ClientID,
		ClientSecret: cfg.AuthType.GoogleConfig.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/oauth/google/redirect", cfg.Dns),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	return &Handler{services, &authCfg, cfg}
}

func (h Handler) Init(api *echo.Group) {
	vkGroup := api.Group("/google")
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
		Path:     "/oauth/google/redirect",
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

	url := oauthGoogleUrlAPI + token.AccessToken
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}

	var googleResp GoogleResp
	if err := json.Unmarshal(resp.Body(), &googleResp); err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка маршалинга google")
	}

	uniqueKey := utils.GenerateKey(googleResp.Name, googleResp.ID)
	log.Println(uniqueKey)
	user := db.User{
		Name:       googleResp.Name,
		AvatarLink: googleResp.Picture,
		Email:      googleResp.Email,
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
