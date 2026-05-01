package response

import (
	"errors"

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

// RespondError sends a structured error response for AppError, or delegates to the error middleware.
func RespondError(c *gin.Context, err error) {
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

	_ = c.Error(err)
}
