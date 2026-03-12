package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	expectedUsername := os.Getenv("ADMIN_USERNAME")
	expectedPassword := os.Getenv("ADMIN_PASSWORD")

	return func(c *gin.Context) {
		if expectedUsername == "" || expectedPassword == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_code": "SERVER_MISCONFIGURED",
				"message":    "Admin credentials are not configured",
			})
			c.Abort()
			return
		}

		username, password, ok := c.Request.BasicAuth()
		if !ok || username != expectedUsername || password != expectedPassword {
			c.Header("WWW-Authenticate", `Basic realm="admin"`)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid admin credentials",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
