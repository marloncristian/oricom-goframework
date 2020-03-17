package mux

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/marloncristian/oricom-goframework/web/authentication"
)

// ParseTokenFromHeader returns the claims from a authorization header
func ParseTokenFromHeader(r *http.Request) (*authentication.Token, error) {
	authorizationHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authorizationHeader) != 2 || strings.ToLower(authorizationHeader[0]) != "bearer" {
		return nil, errors.New("No authorization header")
	}
	token := &authentication.Token{}
	token.Decode(authorizationHeader[1])
	return token, nil
}

// GetTokenRoleFromHeader Gets the role claim of the token from header
func GetTokenRoleFromHeader(r *http.Request) (*authentication.TokenRole, error) {
	token, err := ParseTokenFromHeader(r)
	if err != nil {
		return nil, err
	}
	return token.Role, nil
}

// GetTokenRole Gets the role claim of the token
func GetTokenRole(tkn jwt.MapClaims) []string {
	val := []string{}
	for _, r := range tkn["role"].([]interface{}) {
		val = append(val, r.(string))
	}
	return val
}

// CheckTokenRole returns if an role exists in the roles array
func CheckTokenRole(role string, tkn jwt.MapClaims) bool {
	roles := GetTokenRole(tkn)
	for _, rle := range roles {
		if strings.ToLower(rle) == strings.ToLower(role) {
			return true
		}
	}
	return false
}

// CheckTokenRoleFromHeader returns if an role exists in the roles array from header request
func CheckTokenRoleFromHeader(role string, r *http.Request) bool {
	roles, err := GetTokenRoleFromHeader(r)
	if err != nil || roles.Empty() {
		return false
	}
	return roles.Check(role)
}

// GetTokenSubFromHeader Gets the sub claim of the token from header
func GetTokenSubFromHeader(r *http.Request) (string, error) {
	token, err := ParseTokenFromHeader(r)
	if err != nil {
		return "", err
	}
	return token.Subscriber, nil
}

// GetTokenClaimFromHeader extracts an specific key/claim from bearer token
func GetTokenClaimFromHeader(claim string, r *http.Request) (string, error) {
	token, err := ParseTokenFromHeader(r)
	if err != nil {
		return "", err
	}
	clm, errClm := token.Claim(claim)
	if errClm != nil {
		return "", errClm
	}
	return clm.(string), nil
}

// Authenticate pipeline middleware function
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, err := ParseTokenFromHeader(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// AuthenticateRole checks the authentication for the admin role
func AuthenticateRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !CheckTokenRoleFromHeader(role, r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
