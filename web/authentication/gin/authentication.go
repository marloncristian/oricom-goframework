package gin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marloncristian/oricom-goframework/web/authentication"
)

// Authenticate middleware for authentication and authorization check
func Authenticate(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		if len(authorizationHeader) != 2 || strings.ToLower(authorizationHeader[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token header"})
			c.Abort()
			return
		}
		token := &authentication.Token{}
		if tknErr := token.Decode(authorizationHeader[1]); tknErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": tknErr.Error()})
			c.Abort()
			return
		}

		if len(role) > 0 {
			if !token.Role.Check(role) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
		}

		// add session verification here, like checking if the user and authType
		// combination actually exists if necessary. Try adding caching this (redis)
		// since this middleware might be called a lot
		c.Next()
	}
}
