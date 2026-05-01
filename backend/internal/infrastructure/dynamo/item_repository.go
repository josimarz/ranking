package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/josimar/ranking/backend/internal/domain/item"
	"github.com/josimar/ranking/backend/pkg/apperror"
)

const itemSKPrefix = "ITEM#"

// ItemRepository implements item.Repository using DynamoDB.
type ItemRepository struct {
	client    API
	tableName string
}

// NewItemRepository creates a new ItemRepository.
func NewItemRepository(client API, tableName string) *ItemRepository {
	return &ItemRepository{client: client, tableName: tableName}
}

// Save persists an item.
func (r *ItemRepository) Save(ctx context.Context, i *item.Item) error {
	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      MarshalItem(i),
	})
	if err != nil {
		return fmt.Errorf("save item: %w", err)
	}

	return nil
}

// FindByID retrieves an item by ranking and item ID.
func (r *ItemRepository) FindByID(ctx context.Context, rankingID, itemID string) (*item.Item, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: RankingPK(rankingID)},
			"SK": &types.AttributeValueMemberS{Value: ItemSK(itemID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("find item by id: %w", err)
	}

	if out.Item == nil {
		return nil, apperror.NewNotFoundError(apperror.ItemNotFound, "item not found")
	}

	return UnmarshalItem(out.Item)
}

// FindByRanking retrieves all items for a ranking.
func (r *ItemRepository) FindByRanking(ctx context.Context, rankingID string) ([]*item.Item, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			exprPK:       &types.AttributeValueMemberS{Value: RankingPK(rankingID)},
			exprSKPrefix: &types.AttributeValueMemberS{Value: itemSKPrefix},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("find items by ranking: %w", err)
	}

	items := make([]*item.Item, 0, len(out.Items))
	for _, m := range out.Items {
		i, err := UnmarshalItem(m)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

// Update overwrites an item.
func (r *ItemRepository) Update(ctx context.Context, i *item.Item) error {
	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      MarshalItem(i),
	})
	if err != nil {
		return fmt.Errorf("update item: %w", err)
	}

	return nil
}

// Delete removes an item.
func (r *ItemRepository) Delete(ctx context.Context, rankingID, itemID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: RankingPK(rankingID)},
			"SK": &types.AttributeValueMemberS{Value: ItemSK(itemID)},
		},
	})
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	return nil
}

// DeleteByRanking removes all items for a ranking.
func (r *ItemRepository) DeleteByRanking(ctx context.Context, rankingID string) error {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			exprPK:       &types.AttributeValueMemberS{Value: RankingPK(rankingID)},
			exprSKPrefix: &types.AttributeValueMemberS{Value: itemSKPrefix},
		},
		ProjectionExpression: aws.String("PK, SK"),
	})
	if err != nil {
		return fmt.Errorf("delete by ranking query: %w", err)
	}

	if len(out.Items) == 0 {
		return nil
	}

	requests := make([]types.WriteRequest, 0, len(out.Items))
	for _, m := range out.Items {
		requests = append(requests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"PK": m["PK"],
					"SK": m["SK"],
				},
			},
		})
	}

	_, err = r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			r.tableName: requests,
		},
	})
	if err != nil {
		return fmt.Errorf("delete by ranking batch: %w", err)
	}

	return nil
}

// CountByRanking returns the number of items in a ranking.
func (r *ItemRepository) CountByRanking(ctx context.Context, rankingID string) (int, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			exprPK:       &types.AttributeValueMemberS{Value: RankingPK(rankingID)},
			exprSKPrefix: &types.AttributeValueMemberS{Value: itemSKPrefix},
		},
		Select: types.SelectCount,
	})
	if err != nil {
		return 0, fmt.Errorf("count by ranking: %w", err)
	}

	return int(out.Count), nil
}
