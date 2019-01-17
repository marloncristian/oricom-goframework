package authentication

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CreateToken : creates a jwt token for a especific user
func CreateToken(sub string, name string, role []string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  sub,
		"name": name,
		"role": role,
		"exp":  time.Now().Add(time.Hour * time.Duration(24)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken validates a token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Invalid token")
}

// ParseTokenFromHeader returns the claims from a authorization header
func ParseTokenFromHeader(r *http.Request) (jwt.MapClaims, error) {
	authorizationHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authorizationHeader) != 2 || strings.ToLower(authorizationHeader[0]) != "bearer" {
		return nil, errors.New("No authorization header")
	}
	return ParseToken(authorizationHeader[1])
}

// GetTokenRole Gets the role claim of the token
func GetTokenRole(tkn jwt.MapClaims) []string {
	val := []string{}
	for _, r := range tkn["role"].([]interface{}) {
		val = append(val, r.(string))
	}
	return val
}

// GetTokenRoleFromHeader Gets the role claim of the token from header
func GetTokenRoleFromHeader(r *http.Request) ([]string, error) {
	token, err := ParseTokenFromHeader(r)
	if err != nil {
		return nil, err
	}
	return GetTokenRole(token), nil
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
	if err != nil || len(roles) == 0 {
		return false
	}
	for _, rle := range roles {
		if strings.ToLower(rle) == strings.ToLower(role) {
			return true
		}
	}
	return false
}

// GetTokenSub Gets the sub claim of the token
func GetTokenSub(tkn jwt.MapClaims) string {
	return tkn["sub"].(string)
}

// GetTokenSubFromHeader Gets the sub claim of the token from header
func GetTokenSubFromHeader(r *http.Request) (string, error) {
	token, err := ParseTokenFromHeader(r)
	if err != nil {
		return "", err
	}
	return GetTokenSub(token), nil
}
