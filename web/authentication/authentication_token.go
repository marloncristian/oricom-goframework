package authentication

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// AuthenticationToken struct for token manipulation
type AuthenticationToken struct {
	Subscriber string
	Name       string
	Role       *AuthenticationTokenRole
	Expiration time.Time
	claims     jwt.MapClaims
}

// AuthenticationTokenRole struct for role manipulation
type AuthenticationTokenRole struct {
	Roles []string
}

// NewAuthenticationToken creates a new token object
func NewAuthenticationToken(sub string, name string, roles []string) AuthenticationToken {
	return AuthenticationToken{
		Subscriber: sub,
		Name:       name,
		Role:       NewAuthenticationTokenRole(roles),
	}
}

// NewAuthenticationTokenRole creates a new token role object
func NewAuthenticationTokenRole(roles []string) *AuthenticationTokenRole {
	return &AuthenticationTokenRole{
		Roles: roles,
	}
}

// Encode return a string encoded token
func (token *AuthenticationToken) Encode() (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  token.Subscriber,
		"name": token.Name,
		"role": token.Role.Roles,
		"exp":  time.Now().Add(time.Hour * time.Duration(24)).Unix(),
	})
	tokenString, err := tkn.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Decode decodes a encripted string to a token object
func (token *AuthenticationToken) Decode(tokenString string) error {
	tkn, err := jwt.Parse(tokenString, func(tknJwt *jwt.Token) (interface{}, error) {
		if _, ok := tknJwt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tknJwt.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return err
	}
	concatClaim := func(claim string, claims jwt.MapClaims) []string {
		val := []string{}
		for _, r := range claims[claim].([]interface{}) {
			val = append(val, r.(string))
		}
		return val
	}
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		token.Subscriber = claims["sub"].(string)
		token.Name = claims["name"].(string)
		token.Role = NewAuthenticationTokenRole(concatClaim("role", claims))
		token.Expiration = claims["exp"].(time.Time)
		token.claims = claims
		return nil
	}
	return errors.New("Invalid token")
}

// Claim returns a claim value
func (token *AuthenticationToken) Claim(claim string) (interface{}, error) {
	if token.claims == nil {
		return nil, errors.New("Claims are empty")
	}
	return token.claims[claim], nil
}

// Check verifies if a role exists in the array
func (tokenRole *AuthenticationTokenRole) Check(role string) bool {
	for _, rol := range tokenRole.Roles {
		if strings.ToLower(rol) == strings.ToLower(role) {
			return true
		}
	}
	return false
}

// Empty returns if roles are empty
func (tokenRole *AuthenticationTokenRole) Empty() bool {
	return len(tokenRole.Roles) == 0
}

/*

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

// GetTokenSub Gets the sub claim of the token
func GetTokenSub(tkn jwt.MapClaims) string {
	return tkn["sub"].(string)
}
*/
