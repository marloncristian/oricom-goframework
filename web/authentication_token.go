package authentication

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	tokenSecret string
)

// Init initilizes the global variables
func Init(secret string) {
	tokenSecret = secret
}

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
