package basic_shared_auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const authorizationHeaderName = "Authorization"

func TokenAuthorizationMiddleware(expectedToken string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if len(c.Request().Header) < 1 {
				return c.String(http.StatusUnauthorized, "Missing headers")
			}

			if len(c.Request().Header[authorizationHeaderName]) < 1 {
				return c.String(http.StatusUnauthorized, "Missing authorization header")
			}

			token := c.Request().Header[authorizationHeaderName][0]

			if token == "" {
				return c.String(http.StatusUnauthorized, "Missing token")
			}

			if token != expectedToken {
				return c.String(http.StatusUnauthorized, "Invalid token")
			}

			return next(c)
		}
	}
}
