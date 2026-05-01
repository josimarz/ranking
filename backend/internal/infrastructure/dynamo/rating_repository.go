package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/josimar/ranking/backend/internal/domain/rating"
)

const batchWriteMaxItems = 25

// RatingRepository implements rating.Repository using DynamoDB.
type RatingRepository struct {
	client    API
	tableName string
}

// NewRatingRepository creates a new RatingRepository.
func NewRatingRepository(client API, tableName string) *RatingRepository {
	return &RatingRepository{client: client, tableName: tableName}
}

// Save persists a rating.
func (r *RatingRepository) Save(ctx context.Context, rt *rating.Rating) error {
	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      MarshalRating(rt),
	})
	if err != nil {
		return fmt.Errorf("save rating: %w", err)
	}

	return nil
}

// FindByUserAndItem retrieves a rating by ranking, item, and user.
func (r *RatingRepository) FindByUserAndItem(ctx context.Context, rankingID, itemID, userID string) (*rating.Rating, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: RatingPK(rankingID, itemID)},
			"SK": &types.AttributeValueMemberS{Value: UserSK(userID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("find rating: %w", err)
	}

	if out.Item == nil {
		return nil, nil
	}

	return UnmarshalRating(out.Item)
}

// FindAllByItem returns all ratings for a specific item, handling pagination.
func (r *RatingRepository) FindAllByItem(ctx context.Context, rankingID, itemID string) ([]*rating.Rating, error) {
	var ratings []*rating.Rating
	var lastKey map[string]types.AttributeValue

	for {
		out, err := r.client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(r.tableName),
			KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				exprPK:       &types.AttributeValueMemberS{Value: RatingPK(rankingID, itemID)},
				exprSKPrefix: &types.AttributeValueMemberS{Value: "USER#"},
			},
			ExclusiveStartKey: lastKey,
		})
		if err != nil {
			return nil, fmt.Errorf("find all ratings by item: %w", err)
		}

		for _, m := range out.Items {
			rt, err := UnmarshalRating(m)
			if err != nil {
				return nil, fmt.Errorf("unmarshal rating: %w", err)
			}
			ratings = append(ratings, rt)
		}

		if out.LastEvaluatedKey == nil {
			break
		}
		lastKey = out.LastEvaluatedKey
	}

	return ratings, nil
}

// DeleteByItem removes all ratings for a specific item.
func (r *RatingRepository) DeleteByItem(ctx context.Context, rankingID, itemID string) error {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			exprPK:       &types.AttributeValueMemberS{Value: RatingPK(rankingID, itemID)},
			exprSKPrefix: &types.AttributeValueMemberS{Value: "USER#"},
		},
		ProjectionExpression: aws.String("PK, SK"),
	})
	if err != nil {
		return fmt.Errorf("delete by item query: %w", err)
	}

	if len(out.Items) == 0 {
		return nil
	}

	for i := 0; i < len(out.Items); i += batchWriteMaxItems {
		end := min(i+batchWriteMaxItems, len(out.Items))

		requests := make([]types.WriteRequest, 0, end-i)
		for _, m := range out.Items[i:end] {
			requests = append(requests, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						"PK": m["PK"],
						"SK": m["SK"],
					},
				},
			})
		}

		_, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.tableName: requests,
			},
		})
		if err != nil {
			return fmt.Errorf("delete by item batch: %w", err)
		}
	}

	return nil
}

// DeleteByRanking removes all ratings for all items in a ranking.
func (r *RatingRepository) DeleteByRanking(ctx context.Context, rankingID string, itemIDs []string) error {
	for _, itemID := range itemIDs {
		if err := r.DeleteByItem(ctx, rankingID, itemID); err != nil {
			return fmt.Errorf("delete ratings for item %s: %w", itemID, err)
		}
	}

	return nil
}
