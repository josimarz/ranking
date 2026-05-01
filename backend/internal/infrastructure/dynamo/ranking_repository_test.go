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
	"github.com/josimar/ranking/backend/internal/domain/ranking"
	ddb "github.com/josimar/ranking/backend/internal/infrastructure/dynamo"
	"github.com/stretchr/testify/require"
)

const (
	rankingTestID      = "test-id"
	rankingTestUserID  = "user-1"
	rankingTestTagGame = "games"
)

func newRepoTestRanking() *ranking.Ranking {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	return &ranking.Ranking{
		ID:          rankingTestID,
		Name:        "Test Ranking",
		NameLower:   "test ranking",
		Description: "A test ranking",
		Visibility:  ranking.Public,
		Tags:        []string{rankingTestTagGame, "fun"},
		Attributes: []ranking.Attribute{
			{ID: "attr-1", Name: "Graphics", Description: "Visual quality", Active: true},
		},
		OwnerUserID: rankingTestUserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestRankingRepo_Save(t *testing.T) {
	t.Parallel()

	t.Run("saves ranking with tags via transact write", func(t *testing.T) {
		t.Parallel()

		var capturedInput *dynamodb.TransactWriteItemsInput
		mock := &mockDynamoAPI{
			transactWriteItemsFn: func(_ context.Context, input *dynamodb.TransactWriteItemsInput, _ ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
				capturedInput = input
				return &dynamodb.TransactWriteItemsOutput{}, nil
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		err := repo.Save(context.Background(), newRepoTestRanking())

		require.NoError(t, err)
		// 1 metadata + 2 tags = 3 transact items
		require.Len(t, capturedInput.TransactItems, 3)
	})

	t.Run("returns error when transact write fails", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			transactWriteItemsFn: func(_ context.Context, _ *dynamodb.TransactWriteItemsInput, _ ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
				return nil, errors.New("dynamo error")
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		err := repo.Save(context.Background(), newRepoTestRanking())

		require.Error(t, err)
	})
}

func TestRankingRepo_FindByID(t *testing.T) {
	t.Parallel()

	t.Run("returns ranking when found", func(t *testing.T) {
		t.Parallel()

		r := newRepoTestRanking()
		item := ddb.MarshalRanking(r)

		mock := &mockDynamoAPI{
			getItemFn: func(_ context.Context, input *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				require.Equal(t, aws.String(testTable), input.TableName)
				return &dynamodb.GetItemOutput{Item: item}, nil
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		result, err := repo.FindByID(context.Background(), rankingTestID)

		require.NoError(t, err)
		require.Equal(t, r.ID, result.ID)
		require.Equal(t, r.Name, result.Name)
	})

	t.Run("returns ErrRankingNotFound when item is nil", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{Item: nil}, nil
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		_, err := repo.FindByID(context.Background(), "missing-id")

		require.ErrorIs(t, err, ranking.ErrRankingNotFound)
	})

	t.Run("returns error when get item fails", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return nil, errors.New("dynamo error")
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		_, err := repo.FindByID(context.Background(), rankingTestID)

		require.Error(t, err)
	})
}

func TestRankingRepo_Update(t *testing.T) {
	t.Parallel()

	t.Run("updates metadata and reconciles tags", func(t *testing.T) {
		t.Parallel()

		var putCalled bool
		var batchCalls int
		mock := &mockDynamoAPI{
			putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				putCalled = true
				return &dynamodb.PutItemOutput{}, nil
			},
			queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				return &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"PK":  &types.AttributeValueMemberS{Value: fmt.Sprintf("RANKING#%s", rankingTestID)},
							"SK":  &types.AttributeValueMemberS{Value: "TAG#old-tag"},
							"Tag": &types.AttributeValueMemberS{Value: "old-tag"},
						},
					},
				}, nil
			},
			batchWriteItemFn: func(_ context.Context, _ *dynamodb.BatchWriteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
				batchCalls++
				return &dynamodb.BatchWriteItemOutput{}, nil
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)
		r := newRepoTestRanking()
		r.Tags = []string{rankingTestTagGame, "fun"}

		err := repo.Update(context.Background(), r)

		require.NoError(t, err)
		require.True(t, putCalled)
		// 1 batch for deleting old-tag + 1 batch for adding new tags
		require.Equal(t, 2, batchCalls)
	})

	t.Run("returns error when put item fails", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				return nil, errors.New("dynamo error")
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		err := repo.Update(context.Background(), newRepoTestRanking())

		require.Error(t, err)
	})
}

func TestRankingRepo_Delete(t *testing.T) {
	t.Parallel()

	t.Run("queries and deletes all items under ranking", func(t *testing.T) {
		t.Parallel()

		var batchCalled bool
		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				return &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("RANKING#%s", rankingTestID)},
							"SK": &types.AttributeValueMemberS{Value: "METADATA"},
						},
						{
							"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("RANKING#%s", rankingTestID)},
							"SK": &types.AttributeValueMemberS{Value: "TAG#games"},
						},
					},
				}, nil
			},
			batchWriteItemFn: func(_ context.Context, input *dynamodb.BatchWriteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
				batchCalled = true
				require.Len(t, input.RequestItems[testTable], 2)
				return &dynamodb.BatchWriteItemOutput{}, nil
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		err := repo.Delete(context.Background(), rankingTestID)

		require.NoError(t, err)
		require.True(t, batchCalled)
	})

	t.Run("returns error when query fails", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				return nil, errors.New("dynamo error")
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		err := repo.Delete(context.Background(), rankingTestID)

		require.Error(t, err)
	})
}

func TestRankingRepo_ListPublic(t *testing.T) {
	t.Parallel()

	t.Run("queries GSI2 with correct parameters", func(t *testing.T) {
		t.Parallel()

		r := newRepoTestRanking()
		item := ddb.MarshalRanking(r)

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, input *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				require.Equal(t, aws.String("GSI2"), input.IndexName)
				require.False(t, *input.ScanIndexForward)
				return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		results, cursor, err := repo.ListPublic(context.Background(), 10, "", "createdAt", "desc")

		require.NoError(t, err)
		require.Len(t, results, 1)
		require.Empty(t, cursor)
	})

	t.Run("scan index forward for ascending sort", func(t *testing.T) {
		t.Parallel()

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, input *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				require.True(t, *input.ScanIndexForward)
				return &dynamodb.QueryOutput{}, nil
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		_, _, err := repo.ListPublic(context.Background(), 10, "", "createdAt", "asc")

		require.NoError(t, err)
	})
}

func TestRankingRepo_ListByOwner(t *testing.T) {
	t.Parallel()

	t.Run("queries GSI1 with owner key", func(t *testing.T) {
		t.Parallel()

		r := newRepoTestRanking()
		item := ddb.MarshalRanking(r)

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, input *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				require.Equal(t, aws.String("GSI1"), input.IndexName)
				require.False(t, *input.ScanIndexForward)
				return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		results, cursor, err := repo.ListByOwner(context.Background(), rankingTestUserID, 10, "")

		require.NoError(t, err)
		require.Len(t, results, 1)
		require.Empty(t, cursor)
	})
}

func TestRankingRepo_SearchByName(t *testing.T) {
	t.Parallel()

	t.Run("queries GSI2 with filter expression", func(t *testing.T) {
		t.Parallel()

		r := newRepoTestRanking()
		item := ddb.MarshalRanking(r)

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, input *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				require.Equal(t, aws.String("GSI2"), input.IndexName)
				require.NotNil(t, input.FilterExpression)
				require.Contains(t, *input.FilterExpression, "contains")
				return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		results, cursor, err := repo.SearchByName(context.Background(), "test", 10, "")

		require.NoError(t, err)
		require.Len(t, results, 1)
		require.Empty(t, cursor)
	})
}

func TestRankingRepo_SearchByTag(t *testing.T) {
	t.Parallel()

	t.Run("queries GSI3 then fetches full rankings", func(t *testing.T) {
		t.Parallel()

		r := newRepoTestRanking()
		fullItem := ddb.MarshalRanking(r)

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, input *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				if input.IndexName != nil && *input.IndexName == "GSI3" {
					return &dynamodb.QueryOutput{
						Items: []map[string]types.AttributeValue{
							{
								"RankingID": &types.AttributeValueMemberS{Value: rankingTestID},
							},
						},
					}, nil
				}
				return &dynamodb.QueryOutput{}, nil
			},
			getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{Item: fullItem}, nil
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		results, cursor, err := repo.SearchByTag(context.Background(), rankingTestTagGame, 10, "")

		require.NoError(t, err)
		require.Len(t, results, 1)
		require.Equal(t, rankingTestID, results[0].ID)
		require.Empty(t, cursor)
	})
}

func TestRankingRepo_FindRecent(t *testing.T) {
	t.Parallel()

	t.Run("queries GSI2 with limit 10 descending", func(t *testing.T) {
		t.Parallel()

		r := newRepoTestRanking()
		item := ddb.MarshalRanking(r)

		mock := &mockDynamoAPI{
			queryFn: func(_ context.Context, input *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
				require.Equal(t, aws.String("GSI2"), input.IndexName)
				require.Equal(t, int32(10), *input.Limit)
				require.False(t, *input.ScanIndexForward)
				return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil
			},
		}

		repo := ddb.NewRankingRepository(mock, testTable)

		results, err := repo.FindRecent(context.Background())

		require.NoError(t, err)
		require.Len(t, results, 1)
	})
}

func TestRankingRepo_ImplementsInterface(t *testing.T) {
	t.Parallel()

	mock := &mockDynamoAPI{}
	repo := ddb.NewRankingRepository(mock, testTable)

	var _ ranking.Repository = repo
}
