package gin

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/marloncristian/oricom-goframework/web/authentication"
)

// GetTokenFromHeader returns a token from the header
func GetTokenFromHeader(c echo.Context) (*authentication.AuthenticationToken, error) {
	authorizationHeader := strings.SplitN(c.Request().Header.Get("Authorization"), " ", 2)
	if len(authorizationHeader) != 2 || strings.ToLower(authorizationHeader[0]) != "bearer" {
		return nil, errors.New("No authorization header")
	}
	token := &authentication.AuthenticationToken{}
	if err := token.Decode(authorizationHeader[1]); err != nil {
		return nil, err
	}
	return token, nil
}
