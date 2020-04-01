package gin

import (
	"github.com/labstack/echo/v4"
)

// Authenticate middleware for authentication and authorization check
func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := GetTokenFromHeader(c)
		if err != nil {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

// AuthenticateAdmin middeware for authentication for admin role
func AuthenticateAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := GetTokenFromHeader(c)
		if err != nil {
			return echo.ErrUnauthorized
		}
		if !token.Role.Check("admin") {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
