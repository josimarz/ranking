package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondOK sends a 200 JSON response.
func RespondOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

// RespondCreated sends a 201 JSON response.
func RespondCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}

// RespondNoContent sends a 204 response with no body.
func RespondNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
