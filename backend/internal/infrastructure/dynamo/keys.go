package dynamo

import (
	"fmt"
	"time"
)

// RankingPK returns the partition key for a ranking.
func RankingPK(id string) string {
	return fmt.Sprintf("RANKING#%s", id)
}

// MetadataSK returns the sort key for ranking metadata.
func MetadataSK() string {
	return "METADATA"
}

// ItemSK returns the sort key for an item.
func ItemSK(itemID string) string {
	return fmt.Sprintf("ITEM#%s", itemID)
}

// TagSK returns the sort key for a tag.
func TagSK(tag string) string {
	return fmt.Sprintf("TAG#%s", tag)
}

// RatingPK returns the partition key for a rating.
func RatingPK(rankingID, itemID string) string {
	return fmt.Sprintf("RATING#%s#%s", rankingID, itemID)
}

// UserSK returns the sort key for a user.
func UserSK(userID string) string {
	return fmt.Sprintf("USER#%s", userID)
}

// OwnerGSI1PK returns the GSI1 partition key for an owner.
func OwnerGSI1PK(userID string) string {
	return fmt.Sprintf("OWNER#%s", userID)
}

// RankingGSI1SK returns the GSI1 sort key for a ranking using RFC3339Nano.
func RankingGSI1SK(createdAt time.Time) string {
	return fmt.Sprintf("RANKING#%s", createdAt.Format(time.RFC3339Nano))
}

// PublicGSI2PK returns the GSI2 partition key for public rankings.
func PublicGSI2PK() string {
	return "PUBLIC"
}

// TagGSI3PK returns the GSI3 partition key for a tag.
func TagGSI3PK(tag string) string {
	return fmt.Sprintf("TAG#%s", tag)
}

// TagGSI3SK returns the GSI3 sort key combining timestamp and ranking ID.
func TagGSI3SK(createdAt time.Time, rankingID string) string {
	return fmt.Sprintf("%s#%s", createdAt.Format(time.RFC3339Nano), rankingID)
}
