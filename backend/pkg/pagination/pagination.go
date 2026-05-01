package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/josimar/ranking/backend/pkg/apperror"
)

// EncodeCursor base64-encodes a DynamoDB LastEvaluatedKey into a cursor string.
func EncodeCursor(lastEvaluatedKey map[string]types.AttributeValue) string {
	if lastEvaluatedKey == nil {
		return ""
	}

	raw := make(map[string]any, len(lastEvaluatedKey))
	for k, v := range lastEvaluatedKey {
		var out any
		if err := attributevalue.Unmarshal(v, &out); err != nil {
			return ""
		}
		raw[k] = out
	}

	data, err := json.Marshal(raw)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(data)
}

// DecodeCursor decodes a cursor string back into a DynamoDB ExclusiveStartKey.
func DecodeCursor(cursor string) (map[string]types.AttributeValue, error) {
	if cursor == "" {
		return nil, nil
	}

	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, &apperror.AppError{
			Code:    apperror.InvalidCursor,
			Message: fmt.Sprintf("invalid cursor encoding: %s", err.Error()),
			Status:  http.StatusBadRequest,
		}
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, &apperror.AppError{
			Code:    apperror.InvalidCursor,
			Message: fmt.Sprintf("invalid cursor format: %s", err.Error()),
			Status:  http.StatusBadRequest,
		}
	}

	key, err := attributevalue.MarshalMap(raw)
	if err != nil {
		return nil, &apperror.AppError{
			Code:    apperror.InvalidCursor,
			Message: fmt.Sprintf("invalid cursor data: %s", err.Error()),
			Status:  http.StatusBadRequest,
		}
	}

	return key, nil
}

// ClampLimit constrains a pagination limit between default and max values.
func ClampLimit(limit, defaultLimit, maxLimit int) int {
	if limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}
