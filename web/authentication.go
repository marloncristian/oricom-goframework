package authentication

import (
	"net/http"
	"strings"
)

// Authenticate pipeline middleware function
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || strings.ToLower(auth[0]) != "bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err := ParseToken(auth[1])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// AuthenticateAdmin checks the authentication for the admin role
func AuthenticateAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || strings.ToLower(auth[0]) != "bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tkn, err := ParseToken(auth[1])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//verifies if the current token owner is admin role
		isAdmin := false
		for _, rl := range tkn["role"].([]interface{}) {
			if strings.ToLower(rl.(string)) == "admin" {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
