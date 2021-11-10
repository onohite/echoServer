package oauth

import (
	"backend/internal/config"
	"backend/internal/routes/oauth/discord"
	"backend/internal/routes/oauth/google"
	"backend/internal/routes/oauth/vk"
	"backend/internal/service"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Handler struct {
	Services *service.Service
	cfg      *config.Config
}

type Token struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUuid string
	UserId     string
}

func NewHandler(services *service.Service, cfg *config.Config) *Handler {
	return &Handler{services, cfg}
}

func (h *Handler) Init(api *echo.Group) {
	api.GET("", h.index)
	api.POST("/login", h.login)
	api.POST("/logout", h.logout)
	api.POST("/refresh", h.refresh)
	h.InitVK(api)
	h.InitGoogle(api)
	h.InitDiscord(api)
}

func (h Handler) index(c echo.Context) error {
	title := "Выберите способ авторизации"
	err := c.Render(200, "default_auth.html", title)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (h *Handler) login(c echo.Context) error {
	var uid string
	body, err := ioutil.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()

	if err != nil {
		log.Print(err)
		return err
	}
	uid = string(body)
	log.Printf("get uuid: %s", uid)
	token, err := createToken(uid)
	log.Print(token)
	if err != nil {
		log.Print(err)
		return err
	}
	err = createAuth(h.Services, uid, token)
	if err != nil {
		log.Print(err)
		return err
	}
	out := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}
	return c.JSON(http.StatusOK, out)
}

func createToken(userid string) (*Token, error) {
	var err error
	td := new(Token)
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	accessUUID, _ := uuid.NewV4()
	td.AccessUuid = accessUUID.String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	refreshUUID, _ := uuid.NewV4()
	td.RefreshUuid = refreshUUID.String()

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func createAuth(service *service.Service, userid string, t *Token) error {
	at := time.Unix(t.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(t.RtExpires, 0)
	now := time.Now()

	err := service.Cache.Set(t.AccessUuid, userid, at.Sub(now))
	if err != nil {
		return err
	}
	err = service.Cache.Set(t.RefreshUuid, userid, rt.Sub(now))
	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (h *Handler) refresh(c echo.Context) error {
	mapToken := map[string]string{}
	if err := c.Bind(&mapToken); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	refreshToken := mapToken["refresh_token"]

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Refresh token expired")
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return c.JSON(http.StatusUnauthorized, errors.New("not valid token claims"))
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return c.JSON(http.StatusUnprocessableEntity, errors.New("cant get refresh uuid").Error())
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return c.JSON(http.StatusUnprocessableEntity, errors.New("cant get refresh uuid").Error())
		}
		//Delete the previous Refresh Token
		deleted, delErr := DeleteAuth(h, refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := createToken(userId)
		if createErr != nil {
			return c.JSON(http.StatusForbidden, createErr.Error())
		}
		//save the tokens metadata to redis
		saveErr := createAuth(h.Services, userId, ts)
		if saveErr != nil {
			return c.JSON(http.StatusForbidden, saveErr.Error())
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return c.JSON(http.StatusCreated, tokens)
	} else {
		return c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("bad access token au")
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("bad access token ui")
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func FetchAuth(service *service.Service, authD *AccessDetails) (string, error) {
	userid, err := service.Cache.Get(authD.AccessUuid)
	if err != nil {
		return "", err
	}
	return userid, nil
}

func DeleteAuth(h *Handler, givenUuid string) (int64, error) {
	deleted, err := h.Services.Cache.Del(givenUuid)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (h *Handler) logout(c echo.Context) error {
	au, err := ExtractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	deleted, delErr := DeleteAuth(h, au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	return c.JSON(http.StatusOK, "Successfully logged out")
}

func (h *Handler) InitVK(api *echo.Group) {
	vkGroup := vk.NewHandler(h.Services, h.cfg)
	vkGroup.Init(api)
}

func (h *Handler) InitGoogle(api *echo.Group) {
	gglGroup := google.NewHandler(h.Services, h.cfg)
	gglGroup.Init(api)
}

func (h *Handler) InitDiscord(api *echo.Group) {
	dscGroup := discord.NewHandler(h.Services, h.cfg)
	dscGroup.Init(api)
}
