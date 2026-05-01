package dynamo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/josimar/ranking/backend/internal/domain/item"
	"github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/josimar/ranking/backend/internal/domain/rating"
)

// DynamoDB attribute key constants.
const (
	attrCreatedAt = "CreatedAt"
	attrUpdatedAt = "UpdatedAt"
	attrType      = "Type"
	attrRankingID = "RankingID"
)

// MarshalRanking converts a Ranking to a DynamoDB item.
func MarshalRanking(r *ranking.Ranking) map[string]types.AttributeValue {
	attrs := marshalAttributes(r.Attributes)
	tags := marshalStringList(r.Tags)

	m := map[string]types.AttributeValue{
		"PK":          &types.AttributeValueMemberS{Value: RankingPK(r.ID)},
		"SK":          &types.AttributeValueMemberS{Value: MetadataSK()},
		"GSI1PK":      &types.AttributeValueMemberS{Value: OwnerGSI1PK(r.OwnerUserID)},
		"GSI1SK":      &types.AttributeValueMemberS{Value: RankingGSI1SK(r.CreatedAt)},
		"ID":          &types.AttributeValueMemberS{Value: r.ID},
		"Name":        &types.AttributeValueMemberS{Value: r.Name},
		"NameLower":   &types.AttributeValueMemberS{Value: r.NameLower},
		"Description": &types.AttributeValueMemberS{Value: r.Description},
		"Visibility":  &types.AttributeValueMemberS{Value: string(r.Visibility)},
		"Tags":        tags,
		"Attributes":  &types.AttributeValueMemberS{Value: attrs},
		"OwnerUserID": &types.AttributeValueMemberS{Value: r.OwnerUserID},
		attrCreatedAt: &types.AttributeValueMemberS{Value: r.CreatedAt.Format(time.RFC3339Nano)},
		attrUpdatedAt: &types.AttributeValueMemberS{Value: r.UpdatedAt.Format(time.RFC3339Nano)},
		attrType:      &types.AttributeValueMemberS{Value: "RANKING"},
	}

	if r.Visibility == ranking.Public {
		m["GSI2PK"] = &types.AttributeValueMemberS{Value: PublicGSI2PK()}
		m["GSI2SK"] = &types.AttributeValueMemberS{Value: RankingGSI1SK(r.CreatedAt)}
	}

	return m
}

// UnmarshalRanking converts a DynamoDB item to a Ranking.
func UnmarshalRanking(m map[string]types.AttributeValue) (*ranking.Ranking, error) {
	createdAt, err := time.Parse(time.RFC3339Nano, getS(m, attrCreatedAt))
	if err != nil {
		return nil, fmt.Errorf("parsing CreatedAt: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339Nano, getS(m, attrUpdatedAt))
	if err != nil {
		return nil, fmt.Errorf("parsing UpdatedAt: %w", err)
	}

	attrs, err := unmarshalAttributes(getS(m, "Attributes"))
	if err != nil {
		return nil, fmt.Errorf("parsing Attributes: %w", err)
	}

	return &ranking.Ranking{
		ID:          getS(m, "ID"),
		Name:        getS(m, "Name"),
		NameLower:   getS(m, "NameLower"),
		Description: getS(m, "Description"),
		Visibility:  ranking.Visibility(getS(m, "Visibility")),
		Tags:        unmarshalStringList(m["Tags"]),
		Attributes:  attrs,
		OwnerUserID: getS(m, "OwnerUserID"),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

// MarshalItem converts an Item to a DynamoDB item.
func MarshalItem(i *item.Item) map[string]types.AttributeValue {
	m := map[string]types.AttributeValue{
		"PK":          &types.AttributeValueMemberS{Value: RankingPK(i.RankingID)},
		"SK":          &types.AttributeValueMemberS{Value: ItemSK(i.ID)},
		"ID":          &types.AttributeValueMemberS{Value: i.ID},
		"Name":        &types.AttributeValueMemberS{Value: i.Name},
		attrRankingID: &types.AttributeValueMemberS{Value: i.RankingID},
		"CreatedBy":   &types.AttributeValueMemberS{Value: i.CreatedBy},
		attrCreatedAt: &types.AttributeValueMemberS{Value: i.CreatedAt.Format(time.RFC3339Nano)},
		attrUpdatedAt: &types.AttributeValueMemberS{Value: i.UpdatedAt.Format(time.RFC3339Nano)},
		attrType:      &types.AttributeValueMemberS{Value: "ITEM"},
	}

	if i.ImageKey != "" {
		m["ImageKey"] = &types.AttributeValueMemberS{Value: i.ImageKey}
	}

	return m
}

// UnmarshalItem converts a DynamoDB item to an Item.
func UnmarshalItem(m map[string]types.AttributeValue) (*item.Item, error) {
	createdAt, err := time.Parse(time.RFC3339Nano, getS(m, attrCreatedAt))
	if err != nil {
		return nil, fmt.Errorf("parsing CreatedAt: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339Nano, getS(m, attrUpdatedAt))
	if err != nil {
		return nil, fmt.Errorf("parsing UpdatedAt: %w", err)
	}

	return &item.Item{
		ID:        getS(m, "ID"),
		Name:      getS(m, "Name"),
		ImageKey:  getS(m, "ImageKey"),
		RankingID: getS(m, attrRankingID),
		CreatedBy: getS(m, "CreatedBy"),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// MarshalRating converts a Rating to a DynamoDB item.
func MarshalRating(r *rating.Rating) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"PK":          &types.AttributeValueMemberS{Value: RatingPK(r.RankingID, r.ItemID)},
		"SK":          &types.AttributeValueMemberS{Value: UserSK(r.UserID)},
		attrRankingID: &types.AttributeValueMemberS{Value: r.RankingID},
		"ItemID":      &types.AttributeValueMemberS{Value: r.ItemID},
		"UserID":      &types.AttributeValueMemberS{Value: r.UserID},
		"Scores":      marshalScores(r.Scores),
		attrCreatedAt: &types.AttributeValueMemberS{Value: r.CreatedAt.Format(time.RFC3339Nano)},
		attrUpdatedAt: &types.AttributeValueMemberS{Value: r.UpdatedAt.Format(time.RFC3339Nano)},
		attrType:      &types.AttributeValueMemberS{Value: "RATING"},
	}
}

// UnmarshalRating converts a DynamoDB item to a Rating.
func UnmarshalRating(m map[string]types.AttributeValue) (*rating.Rating, error) {
	createdAt, err := time.Parse(time.RFC3339Nano, getS(m, attrCreatedAt))
	if err != nil {
		return nil, fmt.Errorf("parsing CreatedAt: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339Nano, getS(m, attrUpdatedAt))
	if err != nil {
		return nil, fmt.Errorf("parsing UpdatedAt: %w", err)
	}

	scores, err := unmarshalScores(m["Scores"])
	if err != nil {
		return nil, fmt.Errorf("parsing Scores: %w", err)
	}

	return &rating.Rating{
		RankingID: getS(m, attrRankingID),
		ItemID:    getS(m, "ItemID"),
		UserID:    getS(m, "UserID"),
		Scores:    scores,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// MarshalRankingTag creates a DynamoDB item for a ranking tag (GSI3 support).
func MarshalRankingTag(rankingID, tag string, createdAt time.Time, rankingName string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"PK":          &types.AttributeValueMemberS{Value: RankingPK(rankingID)},
		"SK":          &types.AttributeValueMemberS{Value: TagSK(tag)},
		"GSI3PK":      &types.AttributeValueMemberS{Value: TagGSI3PK(tag)},
		"GSI3SK":      &types.AttributeValueMemberS{Value: TagGSI3SK(createdAt, rankingID)},
		attrRankingID: &types.AttributeValueMemberS{Value: rankingID},
		"Tag":         &types.AttributeValueMemberS{Value: tag},
		"RankingName": &types.AttributeValueMemberS{Value: rankingName},
		attrCreatedAt: &types.AttributeValueMemberS{Value: createdAt.Format(time.RFC3339Nano)},
		attrType:      &types.AttributeValueMemberS{Value: "TAG"},
	}
}

// getS extracts a string value from a DynamoDB attribute map.
func getS(m map[string]types.AttributeValue, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			return s.Value
		}
	}
	return ""
}

func marshalStringList(ss []string) *types.AttributeValueMemberL {
	items := make([]types.AttributeValue, len(ss))
	for i, s := range ss {
		items[i] = &types.AttributeValueMemberS{Value: s}
	}
	return &types.AttributeValueMemberL{Value: items}
}

func unmarshalStringList(av types.AttributeValue) []string {
	l, ok := av.(*types.AttributeValueMemberL)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(l.Value))
	for _, v := range l.Value {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			result = append(result, s.Value)
		}
	}
	return result
}

func marshalAttributes(attrs []ranking.Attribute) string {
	data, _ := json.Marshal(attrs)
	return string(data)
}

func unmarshalAttributes(s string) ([]ranking.Attribute, error) {
	var attrs []ranking.Attribute
	if err := json.Unmarshal([]byte(s), &attrs); err != nil {
		return nil, err
	}
	return attrs, nil
}

func marshalScores(scores map[string]int) *types.AttributeValueMemberM {
	m := make(map[string]types.AttributeValue, len(scores))
	for k, v := range scores {
		m[k] = &types.AttributeValueMemberN{Value: strconv.Itoa(v)}
	}
	return &types.AttributeValueMemberM{Value: m}
}

func unmarshalScores(av types.AttributeValue) (map[string]int, error) {
	m, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return nil, errors.New("expected M type for Scores")
	}
	result := make(map[string]int, len(m.Value))
	for k, v := range m.Value {
		n, ok := v.(*types.AttributeValueMemberN)
		if !ok {
			return nil, fmt.Errorf("expected N type for score key %s", k)
		}
		val, err := strconv.Atoi(n.Value)
		if err != nil {
			return nil, fmt.Errorf("parsing score for key %s: %w", k, err)
		}
		result[k] = val
	}
	return result, nil
}
