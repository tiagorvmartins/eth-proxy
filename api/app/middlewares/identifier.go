package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Identifier middleware
func Identifier() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New().String()

		c.Set("x-identifier-id", id)

		c.Header("x-identifier-id", id)
		c.Next()
	}
}
