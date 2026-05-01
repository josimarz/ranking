package dynamo

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/josimar/ranking/backend/pkg/pagination"
)

// Expression attribute value key constants.
const (
	exprPK       = ":pk"
	exprSKPrefix = ":sk_prefix"
	exprTerm     = ":term"
	gsi1         = "GSI1"
	gsi2         = "GSI2"
	gsi3         = "GSI3"
)

// RankingRepository implements ranking.Repository using DynamoDB.
type RankingRepository struct {
	client    API
	tableName string
}

// NewRankingRepository creates a new RankingRepository.
func NewRankingRepository(client API, tableName string) *RankingRepository {
	return &RankingRepository{client: client, tableName: tableName}
}

// Save persists a new ranking with its tags using a transactional write.
func (r *RankingRepository) Save(ctx context.Context, rk *ranking.Ranking) error {
	items := make([]types.TransactWriteItem, 0, 1+len(rk.Tags))

	items = append(items, types.TransactWriteItem{
		Put: &types.Put{
			TableName: aws.String(r.tableName),
			Item:      MarshalRanking(rk),
		},
	})

	for _, tag := range rk.Tags {
		items = append(items, types.TransactWriteItem{
			Put: &types.Put{
				TableName: aws.String(r.tableName),
				Item:      MarshalRankingTag(rk.ID, tag, rk.CreatedAt, rk.Name),
			},
		})
	}

	_, err := r.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: items,
	})
	if err != nil {
		return fmt.Errorf("save ranking: %w", err)
	}

	return nil
}

// FindByID retrieves a ranking by its ID.
func (r *RankingRepository) FindByID(ctx context.Context, id string) (*ranking.Ranking, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: RankingPK(id)},
			"SK": &types.AttributeValueMemberS{Value: MetadataSK()},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("find ranking by id: %w", err)
	}

	if out.Item == nil {
		return nil, ranking.ErrRankingNotFound
	}

	return UnmarshalRanking(out.Item)
}

// Update overwrites ranking metadata and reconciles tags.
func (r *RankingRepository) Update(ctx context.Context, rk *ranking.Ranking) error {
	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      MarshalRanking(rk),
	})
	if err != nil {
		return fmt.Errorf("update ranking metadata: %w", err)
	}

	if err := r.reconcileTags(ctx, rk); err != nil {
		return fmt.Errorf("update ranking tags: %w", err)
	}

	return nil
}

// Delete removes all items under a ranking partition key.
func (r *RankingRepository) Delete(ctx context.Context, id string) error {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			exprPK: &types.AttributeValueMemberS{Value: RankingPK(id)},
		},
		ProjectionExpression: aws.String("PK, SK"),
	})
	if err != nil {
		return fmt.Errorf("delete ranking query: %w", err)
	}

	if len(out.Items) == 0 {
		return nil
	}

	return r.batchDelete(ctx, out.Items)
}

// ListPublic queries public rankings from GSI2.
func (r *RankingRepository) ListPublic(ctx context.Context, limit int, cursor string, _ string, sortDir string) ([]*ranking.Ranking, string, error) {
	forward := sortDir != "desc"

	startKey, err := pagination.DecodeCursor(cursor)
	if err != nil {
		return nil, "", err
	}

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String(gsi2),
		KeyConditionExpression: aws.String("GSI2PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			exprPK: &types.AttributeValueMemberS{Value: PublicGSI2PK()},
		},
		ScanIndexForward:  aws.Bool(forward),
		Limit:             aws.Int32(int32(limit)),
		ExclusiveStartKey: startKey,
	})
	if err != nil {
		return nil, "", fmt.Errorf("list public rankings: %w", err)
	}

	return r.unmarshalResults(out)
}

// ListByOwner queries rankings owned by a user from GSI1.
func (r *RankingRepository) ListByOwner(ctx context.Context, userID string, limit int, cursor string) ([]*ranking.Ranking, string, error) {
	startKey, err := pagination.DecodeCursor(cursor)
	if err != nil {
		return nil, "", err
	}

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String(gsi1),
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			exprPK: &types.AttributeValueMemberS{Value: OwnerGSI1PK(userID)},
		},
		ScanIndexForward:  aws.Bool(false),
		Limit:             aws.Int32(int32(limit)),
		ExclusiveStartKey: startKey,
	})
	if err != nil {
		return nil, "", fmt.Errorf("list by owner: %w", err)
	}

	return r.unmarshalResults(out)
}

// SearchByName queries GSI2 with a filter on NameLower.
func (r *RankingRepository) SearchByName(ctx context.Context, term string, limit int, cursor string) ([]*ranking.Ranking, string, error) {
	startKey, err := pagination.DecodeCursor(cursor)
	if err != nil {
		return nil, "", err
	}

	lowerTerm := strings.ToLower(term)

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String(gsi2),
		KeyConditionExpression: aws.String("GSI2PK = :pk"),
		FilterExpression:       aws.String("contains(NameLower, :term)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			exprPK:   &types.AttributeValueMemberS{Value: PublicGSI2PK()},
			exprTerm: &types.AttributeValueMemberS{Value: lowerTerm},
		},
		ScanIndexForward:  aws.Bool(false),
		Limit:             aws.Int32(int32(limit)),
		ExclusiveStartKey: startKey,
	})
	if err != nil {
		return nil, "", fmt.Errorf("search by name: %w", err)
	}

	return r.unmarshalResults(out)
}

// SearchByTag queries GSI3 for tag items, then fetches full rankings by ID.
func (r *RankingRepository) SearchByTag(ctx context.Context, tag string, limit int, cursor string) ([]*ranking.Ranking, string, error) {
	startKey, err := pagination.DecodeCursor(cursor)
	if err != nil {
		return nil, "", err
	}

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String(gsi3),
		KeyConditionExpression: aws.String("GSI3PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			exprPK: &types.AttributeValueMemberS{Value: TagGSI3PK(tag)},
		},
		ScanIndexForward:  aws.Bool(false),
		Limit:             aws.Int32(int32(limit)),
		ExclusiveStartKey: startKey,
	})
	if err != nil {
		return nil, "", fmt.Errorf("search by tag query: %w", err)
	}

	nextCursor := pagination.EncodeCursor(out.LastEvaluatedKey)

	rankings := make([]*ranking.Ranking, 0, len(out.Items))
	for _, item := range out.Items {
		rankingID := getS(item, "RankingID")
		if rankingID == "" {
			continue
		}

		rk, err := r.FindByID(ctx, rankingID)
		if err != nil {
			continue
		}

		rankings = append(rankings, rk)
	}

	return rankings, nextCursor, nil
}

// FindRecent returns the 10 most recent public rankings.
func (r *RankingRepository) FindRecent(ctx context.Context) ([]*ranking.Ranking, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String(gsi2),
		KeyConditionExpression: aws.String("GSI2PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			exprPK: &types.AttributeValueMemberS{Value: PublicGSI2PK()},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(10),
	})
	if err != nil {
		return nil, fmt.Errorf("find recent: %w", err)
	}

	rankings := make([]*ranking.Ranking, 0, len(out.Items))
	for _, item := range out.Items {
		rk, err := UnmarshalRanking(item)
		if err != nil {
			return nil, err
		}
		rankings = append(rankings, rk)
	}

	return rankings, nil
}

func (r *RankingRepository) reconcileTags(ctx context.Context, rk *ranking.Ranking) error {
	existing, err := r.queryExistingTags(ctx, rk.ID)
	if err != nil {
		return err
	}

	newTags := make(map[string]bool, len(rk.Tags))
	for _, tag := range rk.Tags {
		newTags[tag] = true
	}

	existingTags := make(map[string]bool, len(existing))
	for _, item := range existing {
		existingTags[getS(item, "Tag")] = true
	}

	// Delete removed tags
	var toDelete []map[string]types.AttributeValue
	for _, item := range existing {
		tag := getS(item, "Tag")
		if !newTags[tag] {
			toDelete = append(toDelete, item)
		}
	}

	if len(toDelete) > 0 {
		if err := r.batchDelete(ctx, toDelete); err != nil {
			return err
		}
	}

	// Add new tags
	var toAdd []types.WriteRequest
	for _, tag := range rk.Tags {
		if !existingTags[tag] {
			toAdd = append(toAdd, types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: MarshalRankingTag(rk.ID, tag, rk.CreatedAt, rk.Name),
				},
			})
		}
	}

	if len(toAdd) > 0 {
		_, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.tableName: toAdd,
			},
		})
		if err != nil {
			return fmt.Errorf("batch add tags: %w", err)
		}
	}

	return nil
}

func (r *RankingRepository) queryExistingTags(ctx context.Context, rankingID string) ([]map[string]types.AttributeValue, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			exprPK:       &types.AttributeValueMemberS{Value: RankingPK(rankingID)},
			exprSKPrefix: &types.AttributeValueMemberS{Value: "TAG#"},
		},
		ProjectionExpression: aws.String("PK, SK, Tag"),
	})
	if err != nil {
		return nil, fmt.Errorf("query existing tags: %w", err)
	}

	return out.Items, nil
}

func (r *RankingRepository) batchDelete(ctx context.Context, items []map[string]types.AttributeValue) error {
	requests := make([]types.WriteRequest, 0, len(items))
	for _, item := range items {
		requests = append(requests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"PK": item["PK"],
					"SK": item["SK"],
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
		return fmt.Errorf("batch delete: %w", err)
	}

	return nil
}

func (r *RankingRepository) unmarshalResults(out *dynamodb.QueryOutput) ([]*ranking.Ranking, string, error) {
	nextCursor := pagination.EncodeCursor(out.LastEvaluatedKey)

	rankings := make([]*ranking.Ranking, 0, len(out.Items))
	for _, item := range out.Items {
		rk, err := UnmarshalRanking(item)
		if err != nil {
			return nil, "", err
		}
		rankings = append(rankings, rk)
	}

	return rankings, nextCursor, nil
}
