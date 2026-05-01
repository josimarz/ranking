package middleware

import (
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	googleuuid "github.com/google/uuid"
)

// Logging logs request/response details with a generated requestId.
func Logging(w io.Writer, env, service string) gin.HandlerFunc {
	var handler slog.Handler
	if env == "local" {
		handler = slog.NewTextHandler(w, nil)
	} else {
		handler = slog.NewJSONHandler(w, nil)
	}
	logger := slog.New(handler)

	return func(c *gin.Context) {
		start := time.Now()
		requestID := googleuuid.New().String()

		c.Next()

		logger.Info("request",
			slog.String("requestId", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("statusCode", c.Writer.Status()),
			slog.String("latency", time.Since(start).String()),
			slog.String("userId", GetUserID(c)),
			slog.String("env", env),
			slog.String("service", service),
		)
	}
}
