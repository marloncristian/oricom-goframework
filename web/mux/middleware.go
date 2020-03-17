package mux

import (
	"errors"
	"net/http"
	"strings"

	"github.com/marloncristian/oricom-goframework/web/authentication"
)

// ParseTokenFromHeader returns the claims from a authorization header
func ParseTokenFromHeader(r *http.Request) (*authentication.AuthenticationToken, error) {
	authorizationHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authorizationHeader) != 2 || strings.ToLower(authorizationHeader[0]) != "bearer" {
		return nil, errors.New("No authorization header")
	}
	token := &authentication.AuthenticationToken{}
	token.Decode(authorizationHeader[1])
	return token, nil
}

// AuthenticateRole checks the authentication for the admin role
func AuthenticateRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := ParseTokenFromHeader(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !token.Role.Check(role) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
