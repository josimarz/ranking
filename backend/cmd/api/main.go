package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	itemuc "github.com/josimar/ranking/backend/internal/application/item"
	rankinguc "github.com/josimar/ranking/backend/internal/application/ranking"
	ratinguc "github.com/josimar/ranking/backend/internal/application/rating"
	"github.com/josimar/ranking/backend/internal/infrastructure/dynamo"
	"github.com/josimar/ranking/backend/internal/infrastructure/image"
	s3store "github.com/josimar/ranking/backend/internal/infrastructure/s3"
	router "github.com/josimar/ranking/backend/internal/interfaces/http"
	"github.com/josimar/ranking/backend/internal/interfaces/http/handler"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	env := envOrDefault("ENV", "local")
	tableName := os.Getenv("DYNAMODB_TABLE_MAIN")
	corsOrigin := envOrDefault("CORS_ORIGIN", "*")
	port := envOrDefault("PORT", "8080")

	ctx := context.Background()

	// AWS clients
	dynamoClient, err := dynamo.NewDynamoClient(ctx)
	if err != nil {
		logger.Error("failed to create DynamoDB client", slog.Any("error", err))
		os.Exit(1)
	}

	s3Client, bucket, err := s3store.NewS3Client(ctx)
	if err != nil {
		logger.Error("failed to create S3 client", slog.Any("error", err))
		os.Exit(1)
	}

	// Repositories
	rankingRepo := dynamo.NewRankingRepository(dynamoClient, tableName)
	itemRepo := dynamo.NewItemRepository(dynamoClient, tableName)
	ratingRepo := dynamo.NewRatingRepository(dynamoClient, tableName)

	// Infrastructure
	presigner := s3.NewPresignClient(s3Client)
	imageStore := s3store.NewS3ImageStore(s3Client, presigner, bucket)
	imageProc := image.NewWebPProcessor()

	// Ranking use cases
	createRanking := rankinguc.NewCreateRankingUseCase(rankingRepo)
	getRanking := rankinguc.NewGetRankingUseCase(rankingRepo)
	updateRanking := rankinguc.NewUpdateRankingUseCase(rankingRepo)
	deleteRanking := rankinguc.NewDeleteRankingUseCase(rankingRepo, itemRepo, ratingRepo, imageStore)
	listPublic := rankinguc.NewListPublicRankingsUseCase(rankingRepo)
	recentRankings := rankinguc.NewGetRecentRankingsUseCase(rankingRepo)
	searchRankings := rankinguc.NewSearchRankingsUseCase(rankingRepo)
	myRankings := rankinguc.NewListMyRankingsUseCase(rankingRepo)

	// Item use cases
	addItem := itemuc.NewAddItemUseCase(rankingRepo, itemRepo, imageStore, imageProc)
	listItems := itemuc.NewListItemsUseCase(rankingRepo, itemRepo, ratingRepo, imageStore)
	updateItem := itemuc.NewUpdateItemUseCase(rankingRepo, itemRepo, imageStore, imageProc)
	deleteItem := itemuc.NewDeleteItemUseCase(rankingRepo, itemRepo, ratingRepo, imageStore)

	// Rating use cases
	submitRating := ratinguc.NewSubmitRatingUseCase(rankingRepo, itemRepo, ratingRepo)
	getMyRating := ratinguc.NewGetMyRatingUseCase(ratingRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler(tableName, dynamoClient)
	rankingHandler := handler.NewRankingHandler(
		createRanking, getRanking, updateRanking, deleteRanking,
		listPublic, recentRankings, searchRankings, myRankings,
	)
	itemHandler := handler.NewItemHandler(addItem, listItems, updateItem, deleteItem)
	ratingHandler := handler.NewRatingHandler(submitRating, getMyRating)

	// Router
	r := router.NewRouter(healthHandler, rankingHandler, itemHandler, ratingHandler, corsOrigin, env)

	if env == "local" {
		logger.Info("starting HTTP server", slog.String("port", port), slog.String("env", env))
		addr := fmt.Sprintf(":%s", port)
		if err := http.ListenAndServe(addr, r); err != nil { //nolint:gosec // local dev only
			logger.Error("server failed", slog.Any("error", err))
			os.Exit(1)
		}
	} else {
		logger.Info("starting Lambda handler", slog.String("env", env))
		ginLambda := ginadapter.New(r)
		lambda.Start(ginLambda.ProxyWithContext)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
