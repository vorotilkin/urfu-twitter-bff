package http

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"twitter-bff/domain/models"
)

func JwtCookieMiddleware(secretKey string, skipPaths ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, path := range skipPaths {
				if c.Path() == path {
					return next(c)
				}
			}

			cookie, err := c.Cookie(models.JWTCookieName)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"message": "Missing token",
					})
				}

				return c.JSON(http.StatusBadRequest, map[string]string{
					"message": "Invalid cookie",
				})
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
				}
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Invalid or expired token",
				})
			}

			// Передаем данные токена дальше
			claims := token.Claims.(jwt.MapClaims)
			c.Set("user", claims)
			return next(c)
		}
	}
}
