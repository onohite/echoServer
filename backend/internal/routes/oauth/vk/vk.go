package vk

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
	"golang.org/x/oauth2/vk"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	Services  *service.Service
	cfg       *oauth2.Config
	cfgServer *config.Config
}

func NewHandler(services *service.Service, cfg *config.Config) *Handler {
	authCfg := oauth2.Config{
		ClientID:     cfg.AuthType.VKconfig.ClientID,
		ClientSecret: cfg.AuthType.VKconfig.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/oauth/vk/redirect", cfg.Dns),
		Scopes:       []string{"account"},
		Endpoint:     vk.Endpoint,
	}
	return &Handler{services, &authCfg, cfg}
}

func (h Handler) Init(api *echo.Group) {
	vkGroup := api.Group("/vk")
	{
		vkGroup.GET("/login", h.Login)
		vkGroup.GET("/redirect", h.Redirect)
	}
}

func (h Handler) Login(c echo.Context) error {
	cookie, state := generateStateOauthCookie(c)
	if err := cookie.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	param := oauth2.SetAuthURLParam("display", "mobile")
	u := h.cfg.AuthCodeURL(state, param)
	return c.Redirect(302, u)
}

func generateStateOauthCookie(c echo.Context) (*sessions.Session, string) {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
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

	log.Println(sess.Values["state"].(string))

	stateTemp := c.QueryParam("state")
	if stateTemp[len(stateTemp)-1] == '}' {
		stateTemp = stateTemp[:len(stateTemp)-1]
	}
	if stateTemp == "" {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	} else if stateTemp != sess.Values["state"].(string) {
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

	fields := strings.Join([]string{"bdate", "sex", "games", "photo_400_orig"}, ",")
	url := fmt.Sprintf("https://api.vk.com/method/%s?v=5.124&fields=%s&access_token=%s", "users.get", fields, token.AccessToken)
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ошибка авторизации")
	}
	var vkResp VkStruct
	if err := json.Unmarshal(resp.Body(), &vkResp); err != nil {
		log.Print(err)
		return err
	}

	vkR := vkResp.Response[0]
	uniqueKey := utils.GenerateKey(vkR.FirstName, vkR.LastName, strconv.Itoa(vkR.ID))
	log.Println(uniqueKey)
	user := db.User{
		Name:       vkR.FirstName + " " + vkR.LastName,
		AvatarLink: vkR.Photo400Orig,
		Sex:        vkR.Sex,
		Bdate:      vkR.Bdate,
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
