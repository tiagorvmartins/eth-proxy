package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Identifier middleware
func Identifier() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate UUID
		id := uuid.New().String()
		// Set context variable
		c.Set("x-identifier-id", id)
		// Set header
		c.Header("x-identifier-id", id)
		c.Next()
	}
}
