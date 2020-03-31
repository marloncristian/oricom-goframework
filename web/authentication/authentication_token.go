package authentication

import (
	"errors"
	"fmt"
	"math"
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
	unixToTime := func(timeFloat float64) time.Time {
		sec, dec := math.Modf(timeFloat)
		return time.Unix(int64(sec), int64(dec*(1e9)))
	}
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		token.Subscriber = claims["sub"].(string)
		token.Name = claims["name"].(string)
		token.Role = NewAuthenticationTokenRole(concatClaim("role", claims))
		token.Expiration = unixToTime(claims["exp"].(float64))
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

// Expired verifies if the token is expired
func (token *AuthenticationToken) Expired() (bool, error) {
	clm, err := token.Claim("exp")
	if err != nil {
		return false, err
	}
	clmObj, ok := clm.(int)
	if !ok {
		return false, errors.New("unknow type")
	}
	exp := time.Unix(int64(clmObj), 0)
	return exp.After(time.Now()), nil
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
