package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func ResellerAuth() gin.HandlerFunc {
	expectedToken := os.Getenv("RESELLER_TOKEN")

	return func(c *gin.Context) {
		if strings.TrimSpace(expectedToken) == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_code": "SERVER_MISCONFIGURED",
				"message":    "Reseller token is not configured",
			})
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code": "UNAUTHORIZED",
				"message":    "Missing authorization header",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid authorization format",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid bearer token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
