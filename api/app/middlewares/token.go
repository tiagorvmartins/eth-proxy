package middlewares

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// Token middleware. For the sake of simplicity we consider any TOKEN, a valid and accepted token!
func Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		}
		c.Header("x-token", token)
		c.Set("x-token", token)
		c.Next()
	}
}
