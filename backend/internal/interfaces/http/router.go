package http

import (
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/josimar/ranking/backend/docs" // swagger generated docs
	"github.com/josimar/ranking/backend/internal/interfaces/http/handler"
	"github.com/josimar/ranking/backend/internal/interfaces/http/middleware"
)

// NewRouter creates a Gin engine with all middleware and routes registered.
func NewRouter(
	healthHandler *handler.HealthHandler,
	rankingHandler *handler.RankingHandler,
	itemHandler *handler.ItemHandler,
	ratingHandler *handler.RatingHandler,
	corsOrigin string,
	env string,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.CORS(corsOrigin))
	r.Use(middleware.Logging(os.Stdout, env, "ranking-api"))
	r.Use(middleware.UserID())
	r.Use(middleware.ErrorHandler())

	requireUser := middleware.RequireUserID()
	aliasRankingID := paramAlias("id", "rankingId")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	v1.GET("/health", wrap(healthHandler, func(h *handler.HealthHandler) gin.HandlerFunc { return h.Check }))

	v1.POST("/rankings", requireUser, wrap(rankingHandler, func(h *handler.RankingHandler) gin.HandlerFunc { return h.Create }))
	v1.GET("/rankings", wrap(rankingHandler, func(h *handler.RankingHandler) gin.HandlerFunc { return h.List }))
	v1.GET("/rankings/recent", wrap(rankingHandler, func(h *handler.RankingHandler) gin.HandlerFunc { return h.Recent }))
	v1.GET("/rankings/search", wrap(rankingHandler, func(h *handler.RankingHandler) gin.HandlerFunc { return h.Search }))
	v1.GET("/rankings/mine", requireUser, wrap(rankingHandler, func(h *handler.RankingHandler) gin.HandlerFunc { return h.Mine }))
	v1.GET("/rankings/:id", wrap(rankingHandler, func(h *handler.RankingHandler) gin.HandlerFunc { return h.Get }))
	v1.PUT("/rankings/:id", requireUser, wrap(rankingHandler, func(h *handler.RankingHandler) gin.HandlerFunc { return h.Update }))
	v1.DELETE("/rankings/:id", requireUser, wrap(rankingHandler, func(h *handler.RankingHandler) gin.HandlerFunc { return h.Delete }))

	v1.POST("/rankings/:id/items", requireUser, aliasRankingID, wrap(itemHandler, func(h *handler.ItemHandler) gin.HandlerFunc { return h.Create }))
	v1.GET("/rankings/:id/items", aliasRankingID, wrap(itemHandler, func(h *handler.ItemHandler) gin.HandlerFunc { return h.List }))
	v1.PUT("/rankings/:id/items/:itemId", requireUser, aliasRankingID, wrap(itemHandler, func(h *handler.ItemHandler) gin.HandlerFunc { return h.Update }))
	v1.DELETE("/rankings/:id/items/:itemId", requireUser, aliasRankingID, wrap(itemHandler, func(h *handler.ItemHandler) gin.HandlerFunc { return h.Delete }))

	v1.PUT("/rankings/:id/items/:itemId/rating", requireUser, aliasRankingID, wrap(ratingHandler, func(h *handler.RatingHandler) gin.HandlerFunc { return h.Submit }))
	v1.GET("/rankings/:id/items/:itemId/rating/mine", requireUser, aliasRankingID, wrap(ratingHandler, func(h *handler.RatingHandler) gin.HandlerFunc { return h.GetMine }))

	return r
}

func wrap[T any](h *T, fn func(*T) gin.HandlerFunc) gin.HandlerFunc {
	if h == nil {
		return func(_ *gin.Context) {}
	}
	return fn(h)
}

// paramAlias creates a middleware that copies a route param under an alias name.
func paramAlias(from, to string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if v := c.Param(from); v != "" {
			c.AddParam(to, v)
		}
		c.Next()
	}
}
