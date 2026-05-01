package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josimar/ranking/backend/pkg/apperror"
)

// RequireUserID rejects requests without a user ID in context.
func RequireUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetUserID(c) == "" {
			_ = c.Error(&apperror.AppError{
				Code:    apperror.MissingUserID,
				Message: "user ID is required",
				Status:  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
