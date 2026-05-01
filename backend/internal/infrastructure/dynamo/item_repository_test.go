package dynamo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/josimar/ranking/backend/internal/domain/item"
	"github.com/josimar/ranking/backend/internal/infrastructure/dynamo"
	"github.com/stretchr/testify/require"
)

const (
	testItemID    = "item-1"
	testRankingID = "ranking-1"
	testUserID    = "user-1"
)

func newTestItem() *item.Item {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	return &item.Item{
		ID:        testItemID,
		Name:      "Test Item",
		ImageKey:  "img/key.webp",
		RankingID: testRankingID,
		CreatedBy: testUserID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func marshaledTestItem() map[string]types.AttributeValue {
	return dynamo.MarshalItem(newTestItem())
}

func TestItemRepository_Save(t *testing.T) {
	t.Parallel()

	t.Run("saves item successfully", func(t *testing.T) {
		t.Parallel()

		var captured *dynamodb.PutItemInput
		mock := &mockDynamoAPI{
			putItemFn: func(_ context.Context, params *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				captured = params
				return &dynamodb.PutItemOutput{}, nil
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		err := repo.Save(context.Background(), newTestItem())

		require.NoError(t, err)
		require.Equal(t, testTable, *captured.TableName)
		require.Equal(t, "RANKING#ranking-1", captured.Item["PK"].(*types.AttributeValueMemberS).Value)
		require.Equal(t, "ITEM#item-1", captured.Item["SK"].(*types.AttributeValueMemberS).Value)
	})

	t.Run("returns error on DynamoDB failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				return nil, errors.New("dynamo error")
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		err := repo.Save(context.Background(), newTestItem())

		require.Error(t, err)
	})
}

func TestItemRepository_FindByID(t *testing.T) {
	t.Parallel()

	t.Run("returns item when found", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			getItemFn: func(_ context.Context, params *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				require.Equal(t, "RANKING#ranking-1", params.Key["PK"].(*types.AttributeValueMemberS).Value)
				require.Equal(t, "ITEM#item-1", params.Key["SK"].(*types.AttributeValueMemberS).Value)
				return &dynamodb.GetItemOutput{Item: marshaledTestItem()}, nil
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		result, err := repo.FindByID(context.Background(), testRankingID, testItemID)

		require.NoError(t, err)
		require.Equal(t, testItemID, result.ID)
		require.Equal(t, "Test Item", result.Name)
	})

	t.Run("returns not found error when item missing", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{Item: nil}, nil
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		result, err := repo.FindByID(context.Background(), testRankingID, testItemID)

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("returns error on DynamoDB failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return nil, errors.New("dynamo error")
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		result, err := repo.FindByID(context.Background(), testRankingID, testItemID)

		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestItemRepository_FindByRanking(t *testing.T) {
	t.Parallel()

	t.Run("returns items for ranking", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, params *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				require.Contains(t, *params.KeyConditionExpression, "PK = :pk")
				require.Contains(t, *params.KeyConditionExpression, "begins_with(SK, :sk_prefix)")
				return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{marshaledTestItem()}}, nil
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		items, err := repo.FindByRanking(context.Background(), testRankingID)

		require.NoError(t, err)
		require.Len(t, items, 1)
		require.Equal(t, testItemID, items[0].ID)
	})

	t.Run("returns empty slice when no items", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, nil
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		items, err := repo.FindByRanking(context.Background(), testRankingID)

		require.NoError(t, err)
		require.Empty(t, items)
	})

	t.Run("returns error on DynamoDB failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				return nil, errors.New("dynamo error")
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		items, err := repo.FindByRanking(context.Background(), testRankingID)

		require.Error(t, err)
		require.Nil(t, items)
	})
}

func TestItemRepository_Update(t *testing.T) {
	t.Parallel()

	t.Run("updates item successfully", func(t *testing.T) {
		t.Parallel()

		var captured *dynamodb.PutItemInput
		mock := &mockDynamoAPI{
			putItemFn: func(_ context.Context, params *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				captured = params
				return &dynamodb.PutItemOutput{}, nil
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		err := repo.Update(context.Background(), newTestItem())

		require.NoError(t, err)
		require.Equal(t, testTable, *captured.TableName)
	})

	t.Run("returns error on DynamoDB failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				return nil, errors.New("dynamo error")
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		err := repo.Update(context.Background(), newTestItem())

		require.Error(t, err)
	})
}

func TestItemRepository_Delete(t *testing.T) {
	t.Parallel()

	t.Run("deletes item successfully", func(t *testing.T) {
		t.Parallel()

		var captured *dynamodb.DeleteItemInput
		mock := &mockDynamoAPI{
			deleteItemFn: func(_ context.Context, params *dynamodb.DeleteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
				captured = params
				return &dynamodb.DeleteItemOutput{}, nil
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		err := repo.Delete(context.Background(), testRankingID, testItemID)

		require.NoError(t, err)
		require.Equal(t, "RANKING#ranking-1", captured.Key["PK"].(*types.AttributeValueMemberS).Value)
		require.Equal(t, "ITEM#item-1", captured.Key["SK"].(*types.AttributeValueMemberS).Value)
	})

	t.Run("returns error on DynamoDB failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			deleteItemFn: func(_ context.Context, _ *dynamodb.DeleteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
				return nil, errors.New("dynamo error")
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		err := repo.Delete(context.Background(), testRankingID, testItemID)

		require.Error(t, err)
	})
}

func TestItemRepository_DeleteByRanking(t *testing.T) {
	t.Parallel()

	t.Run("queries and batch deletes items", func(t *testing.T) {
		t.Parallel()

		queryCalled := false
		batchCalled := false

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				queryCalled = true
				return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{
					{
						"PK": &types.AttributeValueMemberS{Value: "RANKING#ranking-1"},
						"SK": &types.AttributeValueMemberS{Value: "ITEM#item-1"},
					},
				}}, nil
			},
			batchWriteItemFn: func(_ context.Context, params *dynamodb.BatchWriteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
				batchCalled = true
				reqs := params.RequestItems[testTable]
				require.Len(t, reqs, 1)
				key := reqs[0].DeleteRequest.Key
				require.Equal(t, "RANKING#ranking-1", key["PK"].(*types.AttributeValueMemberS).Value)
				require.Equal(t, "ITEM#item-1", key["SK"].(*types.AttributeValueMemberS).Value)
				return &dynamodb.BatchWriteItemOutput{}, nil
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		err := repo.DeleteByRanking(context.Background(), testRankingID)

		require.NoError(t, err)
		require.True(t, queryCalled)
		require.True(t, batchCalled)
	})

	t.Run("succeeds when no items to delete", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, nil
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		err := repo.DeleteByRanking(context.Background(), testRankingID)

		require.NoError(t, err)
	})

	t.Run("returns error on query failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				return nil, errors.New("query error")
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		err := repo.DeleteByRanking(context.Background(), testRankingID)

		require.Error(t, err)
	})

	t.Run("returns error on batch write failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{
					{
						"PK": &types.AttributeValueMemberS{Value: "RANKING#ranking-1"},
						"SK": &types.AttributeValueMemberS{Value: "ITEM#item-1"},
					},
				}}, nil
			},
			batchWriteItemFn: func(_ context.Context, _ *dynamodb.BatchWriteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
				return nil, errors.New("batch error")
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		err := repo.DeleteByRanking(context.Background(), testRankingID)

		require.Error(t, err)
	})
}

func TestItemRepository_CountByRanking(t *testing.T) {
	t.Parallel()

	t.Run("returns count of items", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, params *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				require.Equal(t, types.SelectCount, params.Select)
				return &dynamodb.QueryOutput{Count: 5}, nil
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		count, err := repo.CountByRanking(context.Background(), testRankingID)

		require.NoError(t, err)
		require.Equal(t, 5, count)
	})

	t.Run("returns zero when no items", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				return &dynamodb.QueryOutput{Count: 0}, nil
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		count, err := repo.CountByRanking(context.Background(), testRankingID)

		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("returns error on DynamoDB failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				return nil, errors.New("dynamo error")
			},
		}

		repo := dynamo.NewItemRepository(mock, testTable)
		count, err := repo.CountByRanking(context.Background(), testRankingID)

		require.Error(t, err)
		require.Equal(t, 0, count)
	})
}

func TestItemRepository_ImplementsInterface(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{}
	repo := dynamo.NewItemRepository(mock, testTable)

	var _ item.Repository = repo
}
