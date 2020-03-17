package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Authenticate middleware for authentication and authorization check
func Authenticate(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := GetTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if len(role) > 0 && !token.Role.Check(role) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
