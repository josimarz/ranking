package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/josimar/ranking/backend/pkg/uuid"
)

const userIDKey = "userId"

// UserID extracts X-User-Id header, validates UUID format, and injects into context.
func UserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-User-Id")
		if id != "" && uuid.IsValid(id) {
			c.Set(userIDKey, id)
		}
		c.Next()
	}
}

// GetUserID retrieves the user ID from the Gin context.
func GetUserID(c *gin.Context) string {
	v, _ := c.Get(userIDKey)
	id, _ := v.(string)
	return id
}
