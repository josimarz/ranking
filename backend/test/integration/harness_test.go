//go:build integration

package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamotypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	itemuc "github.com/josimar/ranking/backend/internal/application/item"
	rankinguc "github.com/josimar/ranking/backend/internal/application/ranking"
	ratinguc "github.com/josimar/ranking/backend/internal/application/rating"
	"github.com/josimar/ranking/backend/internal/infrastructure/dynamo"
	"github.com/josimar/ranking/backend/internal/infrastructure/image"
	s3store "github.com/josimar/ranking/backend/internal/infrastructure/s3"
	router "github.com/josimar/ranking/backend/internal/interfaces/http"
	"github.com/josimar/ranking/backend/internal/interfaces/http/handler"
)

const (
	testTableName  = "ranking-integration-test"
	testBucketName = "ranking-images-integration-test"
	testRegion     = "us-east-1"
)

type testEnv struct {
	server       *httptest.Server
	client       *http.Client
	dynamoClient *dynamodb.Client
	ctx          context.Context
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "localstack/localstack:3.8",
			ExposedPorts: []string{"4566/tcp"},
			Env: map[string]string{
				"SERVICES":       "dynamodb,s3",
				"DEFAULT_REGION": testRegion,
			},
			WaitingFor: wait.ForHTTP("/_localstack/health").WithPort("4566/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "4566")
	require.NoError(t, err)

	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(testRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "test")),
	)
	require.NoError(t, err)

	dynamoClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	createTable(ctx, t, dynamoClient)

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	createBucket(ctx, t, s3Client)

	// Wire up application (mirrors main.go)
	rankingRepo := dynamo.NewRankingRepository(dynamoClient, testTableName)
	itemRepo := dynamo.NewItemRepository(dynamoClient, testTableName)
	ratingRepo := dynamo.NewRatingRepository(dynamoClient, testTableName)

	presigner := s3.NewPresignClient(s3Client)
	imageStore := s3store.NewS3ImageStore(s3Client, presigner, testBucketName)
	imageProc := image.NewWebPProcessor()

	createRanking := rankinguc.NewCreateRankingUseCase(rankingRepo)
	getRanking := rankinguc.NewGetRankingUseCase(rankingRepo)
	updateRanking := rankinguc.NewUpdateRankingUseCase(rankingRepo)
	deleteRanking := rankinguc.NewDeleteRankingUseCase(rankingRepo, itemRepo, ratingRepo, imageStore)
	listPublic := rankinguc.NewListPublicRankingsUseCase(rankingRepo)
	recentRankings := rankinguc.NewGetRecentRankingsUseCase(rankingRepo)
	searchRankings := rankinguc.NewSearchRankingsUseCase(rankingRepo)
	myRankings := rankinguc.NewListMyRankingsUseCase(rankingRepo)

	addItem := itemuc.NewAddItemUseCase(rankingRepo, itemRepo, imageStore, imageProc)
	listItems := itemuc.NewListItemsUseCase(rankingRepo, itemRepo, ratingRepo, imageStore)
	updateItem := itemuc.NewUpdateItemUseCase(rankingRepo, itemRepo, imageStore, imageProc)
	deleteItem := itemuc.NewDeleteItemUseCase(rankingRepo, itemRepo, ratingRepo, imageStore)

	submitRating := ratinguc.NewSubmitRatingUseCase(rankingRepo, itemRepo, ratingRepo)
	getMyRating := ratinguc.NewGetMyRatingUseCase(ratingRepo)

	healthHandler := handler.NewHealthHandler(testTableName, dynamoClient)
	rankingHandler := handler.NewRankingHandler(
		createRanking, getRanking, updateRanking, deleteRanking,
		listPublic, recentRankings, searchRankings, myRankings,
	)
	itemHandler := handler.NewItemHandler(addItem, listItems, updateItem, deleteItem)
	ratingHandler := handler.NewRatingHandler(submitRating, getMyRating)

	r := router.NewRouter(healthHandler, rankingHandler, itemHandler, ratingHandler, "*", "test")

	server := httptest.NewServer(r)
	t.Cleanup(server.Close)

	return &testEnv{
		server:       server,
		client:       server.Client(),
		dynamoClient: dynamoClient,
		ctx:          ctx,
	}
}

func createTable(ctx context.Context, t *testing.T, client *dynamodb.Client) {
	t.Helper()

	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(testTableName),
		KeySchema: []dynamotypes.KeySchemaElement{
			{AttributeName: aws.String("PK"), KeyType: dynamotypes.KeyTypeHash},
			{AttributeName: aws.String("SK"), KeyType: dynamotypes.KeyTypeRange},
		},
		AttributeDefinitions: []dynamotypes.AttributeDefinition{
			{AttributeName: aws.String("PK"), AttributeType: dynamotypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("SK"), AttributeType: dynamotypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI1PK"), AttributeType: dynamotypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI1SK"), AttributeType: dynamotypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI2PK"), AttributeType: dynamotypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI2SK"), AttributeType: dynamotypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI3PK"), AttributeType: dynamotypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI3SK"), AttributeType: dynamotypes.ScalarAttributeTypeS},
		},
		GlobalSecondaryIndexes: []dynamotypes.GlobalSecondaryIndex{
			{
				IndexName: aws.String("GSI1"),
				KeySchema: []dynamotypes.KeySchemaElement{
					{AttributeName: aws.String("GSI1PK"), KeyType: dynamotypes.KeyTypeHash},
					{AttributeName: aws.String("GSI1SK"), KeyType: dynamotypes.KeyTypeRange},
				},
				Projection: &dynamotypes.Projection{ProjectionType: dynamotypes.ProjectionTypeAll},
			},
			{
				IndexName: aws.String("GSI2"),
				KeySchema: []dynamotypes.KeySchemaElement{
					{AttributeName: aws.String("GSI2PK"), KeyType: dynamotypes.KeyTypeHash},
					{AttributeName: aws.String("GSI2SK"), KeyType: dynamotypes.KeyTypeRange},
				},
				Projection: &dynamotypes.Projection{ProjectionType: dynamotypes.ProjectionTypeAll},
			},
			{
				IndexName: aws.String("GSI3"),
				KeySchema: []dynamotypes.KeySchemaElement{
					{AttributeName: aws.String("GSI3PK"), KeyType: dynamotypes.KeyTypeHash},
					{AttributeName: aws.String("GSI3SK"), KeyType: dynamotypes.KeyTypeRange},
				},
				Projection: &dynamotypes.Projection{ProjectionType: dynamotypes.ProjectionTypeAll},
			},
		},
		BillingMode: dynamotypes.BillingModePayPerRequest,
	})
	require.NoError(t, err)
}

func createBucket(ctx context.Context, t *testing.T, client *s3.Client) {
	t.Helper()

	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(testBucketName),
	})
	require.NoError(t, err)
}
