package middleware

import (
	"backend/internal/routes/oauth"
	"backend/internal/service"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"log"
)

func JWTMiddleware(service *service.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := oauth.ExtractTokenMetadata(c.Request())
			if err != nil {
				log.Println("unauthorized MIDDLEWARE close conn/cant extract token")
				return next(c)
			}
			uuid, err := oauth.FetchAuth(service, token)
			if err != nil {
				log.Println("unauthorized MIDDLEWARE close conn/cant extract uuid")
				return next(c)
			}
			sess, _ := session.Get("session", c)
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   60 * 60 * 15,
				HttpOnly: false,
			}
			sess.Values["uuid"] = uuid
			sess.Save(c.Request(), c.Response())
			return next(c)
		}
	}
}
