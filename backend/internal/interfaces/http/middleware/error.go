package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josimar/ranking/backend/pkg/apperror"
)

type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// ErrorHandler maps errors to structured JSON responses.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors[0].Err
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Status, errorResponse{
				Error: errorBody{
					Code:    appErr.Code,
					Message: appErr.Message,
					Status:  appErr.Status,
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse{
			Error: errorBody{
				Code:    apperror.InternalError,
				Message: "internal server error",
				Status:  http.StatusInternalServerError,
			},
		})
	}
}
