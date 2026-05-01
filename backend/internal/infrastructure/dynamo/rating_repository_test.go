package dynamo_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/josimar/ranking/backend/internal/domain/rating"
	"github.com/josimar/ranking/backend/internal/infrastructure/dynamo"
	"github.com/stretchr/testify/require"
)

const (
	testAttr1 = "attr1"
	testAttr2 = "attr2"
)

func newTestRating() *rating.Rating {
	return &rating.Rating{
		RankingID: testRankingID,
		ItemID:    testItemID,
		UserID:    testUserID,
		Scores:    map[string]int{testAttr1: 80, testAttr2: 90},
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestRatingRepo_Save_Success(t *testing.T) {
	t.Parallel()

	var captured *dynamodb.PutItemInput
	mock := &mockDynamoAPI{
		putItemFn: func(_ context.Context, params *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			captured = params
			return &dynamodb.PutItemOutput{}, nil
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	err := repo.Save(context.Background(), newTestRating())

	require.NoError(t, err)
	require.NotNil(t, captured)
	require.Equal(t, aws.String(testTable), captured.TableName)
	require.Equal(t, "RATING#ranking-1#item-1", captured.Item["PK"].(*types.AttributeValueMemberS).Value)
	require.Equal(t, "USER#user-1", captured.Item["SK"].(*types.AttributeValueMemberS).Value)
}

func TestRatingRepo_Save_Error(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("dynamo error")
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	err := repo.Save(context.Background(), newTestRating())

	require.Error(t, err)
}

func TestRatingRepo_FindByUserAndItem_Success(t *testing.T) {
	t.Parallel()

	r := newTestRating()
	marshaledItem := dynamo.MarshalRating(r)

	mock := &mockDynamoAPI{
		getItemFn: func(_ context.Context, params *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			require.Equal(t, aws.String(testTable), params.TableName)
			require.Equal(t, "RATING#ranking-1#item-1", params.Key["PK"].(*types.AttributeValueMemberS).Value)
			require.Equal(t, "USER#user-1", params.Key["SK"].(*types.AttributeValueMemberS).Value)
			return &dynamodb.GetItemOutput{Item: marshaledItem}, nil
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	result, err := repo.FindByUserAndItem(context.Background(), testRankingID, testItemID, testUserID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testRankingID, result.RankingID)
	require.Equal(t, testItemID, result.ItemID)
	require.Equal(t, testUserID, result.UserID)
}

func TestRatingRepo_FindByUserAndItem_NotFound(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{
		getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{}, nil
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	result, err := repo.FindByUserAndItem(context.Background(), testRankingID, testItemID, testUserID)

	require.NoError(t, err)
	require.Nil(t, result)
}

func TestRatingRepo_FindByUserAndItem_Error(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{
		getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return nil, errors.New("dynamo error")
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	result, err := repo.FindByUserAndItem(context.Background(), testRankingID, testItemID, testUserID)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRatingRepo_FindAllByItem_Success(t *testing.T) {
	t.Parallel()

	r1 := newTestRating()
	r2 := &rating.Rating{
		RankingID: testRankingID,
		ItemID:    testItemID,
		UserID:    "user-2",
		Scores:    map[string]int{testAttr1: 70, "attr2": 60},
		CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, params *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			require.Equal(t, aws.String(testTable), params.TableName)
			return &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					dynamo.MarshalRating(r1),
					dynamo.MarshalRating(r2),
				},
			}, nil
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	results, err := repo.FindAllByItem(context.Background(), testRankingID, testItemID)

	require.NoError(t, err)
	require.Len(t, results, 2)
	require.Equal(t, testUserID, results[0].UserID)
	require.Equal(t, "user-2", results[1].UserID)
}

func TestRatingRepo_FindAllByItem_Empty(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, nil
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	results, err := repo.FindAllByItem(context.Background(), testRankingID, testItemID)

	require.NoError(t, err)
	require.Empty(t, results)
}

func TestRatingRepo_FindAllByItem_Error(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("dynamo error")
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	results, err := repo.FindAllByItem(context.Background(), testRankingID, testItemID)

	require.Error(t, err)
	require.Nil(t, results)
}

func TestRatingRepo_FindAllByItem_QueryUsesBeginsWith(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, params *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			require.Contains(t, *params.KeyConditionExpression, "begins_with")
			require.Equal(t, "RATING#ranking-1#item-1", params.ExpressionAttributeValues[":pk"].(*types.AttributeValueMemberS).Value)
			require.Equal(t, "USER#", params.ExpressionAttributeValues[":sk_prefix"].(*types.AttributeValueMemberS).Value)
			return &dynamodb.QueryOutput{}, nil
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	_, err := repo.FindAllByItem(context.Background(), testRankingID, testItemID)

	require.NoError(t, err)
}

func TestRatingRepo_FindAllByItem_Pagination(t *testing.T) {
	t.Parallel()

	r1 := newTestRating()
	r2 := &rating.Rating{
		RankingID: testRankingID,
		ItemID:    testItemID,
		UserID:    "user-2",
		Scores:    map[string]int{testAttr1: 70},
		CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	callCount := 0
	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, params *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			callCount++
			if callCount == 1 {
				return &dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{dynamo.MarshalRating(r1)},
					LastEvaluatedKey: map[string]types.AttributeValue{"PK": &types.AttributeValueMemberS{Value: "cursor"}},
				}, nil
			}
			require.NotNil(t, params.ExclusiveStartKey)
			return &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{dynamo.MarshalRating(r2)},
			}, nil
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	results, err := repo.FindAllByItem(context.Background(), testRankingID, testItemID)

	require.NoError(t, err)
	require.Len(t, results, 2)
	require.Equal(t, 2, callCount)
}

func TestRatingRepo_DeleteByItem_Success(t *testing.T) {
	t.Parallel()

	r := newTestRating()
	queryCalled := false
	batchCalled := false

	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			queryCalled = true
			return &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{dynamo.MarshalRating(r)},
			}, nil
		},
		batchWriteItemFn: func(_ context.Context, params *dynamodb.BatchWriteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
			batchCalled = true
			require.Len(t, params.RequestItems[testTable], 1)
			delReq := params.RequestItems[testTable][0].DeleteRequest
			require.NotNil(t, delReq)
			require.Equal(t, "RATING#ranking-1#item-1", delReq.Key["PK"].(*types.AttributeValueMemberS).Value)
			require.Equal(t, "USER#user-1", delReq.Key["SK"].(*types.AttributeValueMemberS).Value)
			return &dynamodb.BatchWriteItemOutput{}, nil
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	err := repo.DeleteByItem(context.Background(), testRankingID, testItemID)

	require.NoError(t, err)
	require.True(t, queryCalled)
	require.True(t, batchCalled)
}

func TestRatingRepo_DeleteByItem_NoRatings(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, nil
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	err := repo.DeleteByItem(context.Background(), testRankingID, testItemID)

	require.NoError(t, err)
}

func TestRatingRepo_DeleteByItem_QueryError(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("query error")
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	err := repo.DeleteByItem(context.Background(), testRankingID, testItemID)

	require.Error(t, err)
}

func TestRatingRepo_DeleteByItem_BatchWriteError(t *testing.T) {
	t.Parallel()

	r := newTestRating()
	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{dynamo.MarshalRating(r)},
			}, nil
		},
		batchWriteItemFn: func(_ context.Context, _ *dynamodb.BatchWriteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
			return nil, errors.New("batch error")
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	err := repo.DeleteByItem(context.Background(), testRankingID, testItemID)

	require.Error(t, err)
}

func TestRatingRepo_DeleteByItem_BatchChunking(t *testing.T) {
	t.Parallel()

	items := make([]map[string]types.AttributeValue, 30)
	for i := range items {
		r := &rating.Rating{
			RankingID: testRankingID,
			ItemID:    testItemID,
			UserID:    fmt.Sprintf("user-%d", i),
			Scores:    map[string]int{testAttr1: 50},
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		items[i] = dynamo.MarshalRating(r)
	}

	batchCount := 0
	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: items}, nil
		},
		batchWriteItemFn: func(_ context.Context, params *dynamodb.BatchWriteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
			batchCount++
			require.LessOrEqual(t, len(params.RequestItems[testTable]), 25)
			return &dynamodb.BatchWriteItemOutput{}, nil
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	err := repo.DeleteByItem(context.Background(), testRankingID, testItemID)

	require.NoError(t, err)
	require.Equal(t, 2, batchCount)
}

func TestRatingRepo_DeleteByRanking_Success(t *testing.T) {
	t.Parallel()

	r1 := newTestRating()
	r2 := &rating.Rating{
		RankingID: testRankingID,
		ItemID:    "item-2",
		UserID:    testUserID,
		Scores:    map[string]int{testAttr1: 50},
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	queryCount := 0
	batchCount := 0

	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, params *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			queryCount++
			pk := params.ExpressionAttributeValues[":pk"].(*types.AttributeValueMemberS).Value
			switch pk {
			case "RATING#ranking-1#item-1":
				return &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{dynamo.MarshalRating(r1)},
				}, nil
			case "RATING#ranking-1#item-2":
				return &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{dynamo.MarshalRating(r2)},
				}, nil
			default:
				return &dynamodb.QueryOutput{}, nil
			}
		},
		batchWriteItemFn: func(_ context.Context, _ *dynamodb.BatchWriteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
			batchCount++
			return &dynamodb.BatchWriteItemOutput{}, nil
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	err := repo.DeleteByRanking(context.Background(), testRankingID, []string{testItemID, "item-2"})

	require.NoError(t, err)
	require.Equal(t, 2, queryCount)
	require.Equal(t, 2, batchCount)
}

func TestRatingRepo_DeleteByRanking_EmptyItemIDs(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{}
	repo := dynamo.NewRatingRepository(mock, testTable)

	err := repo.DeleteByRanking(context.Background(), testRankingID, []string{})

	require.NoError(t, err)
}

func TestRatingRepo_DeleteByRanking_QueryError(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("query error")
		},
	}

	repo := dynamo.NewRatingRepository(mock, testTable)

	err := repo.DeleteByRanking(context.Background(), testRankingID, []string{testItemID})

	require.Error(t, err)
}

func TestRatingRepo_ImplementsInterface(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{}
	repo := dynamo.NewRatingRepository(mock, testTable)

	var _ rating.Repository = repo
}
