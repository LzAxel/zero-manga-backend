package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lzaxel/zero-manga-backend/internal/jwt"
)

type JWTValidator interface {
	ValidateToken(token string) (jwt.Claims, error)
}

func (h *Handler) Authorized() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			splitAuth := strings.Split(authHeader, " ")
			if len(splitAuth) != 2 {
				return echo.NewHTTPError(http.StatusUnauthorized, jwt.ErrInvalidToken)
			}

			claims, err := h.jwtValidator.ValidateToken(splitAuth[1])
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			c.Set("userID", claims.Subject)

			if err := next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}
