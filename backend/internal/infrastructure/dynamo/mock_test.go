package dynamo_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const testTable = "test-table"

// mockDynamoAPI implements dynamo.API for testing.
type mockDynamoAPI struct {
	putItemFn            func(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	getItemFn            func(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	deleteItemFn         func(ctx context.Context, input *dynamodb.DeleteItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	queryFn              func(ctx context.Context, input *dynamodb.QueryInput, opts ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	batchWriteItemFn     func(ctx context.Context, input *dynamodb.BatchWriteItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
	transactWriteItemsFn func(ctx context.Context, input *dynamodb.TransactWriteItemsInput, opts ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

func (m *mockDynamoAPI) PutItem(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.putItemFn(ctx, input, opts...)
}

func (m *mockDynamoAPI) GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m.getItemFn(ctx, input, opts...)
}

func (m *mockDynamoAPI) DeleteItem(ctx context.Context, input *dynamodb.DeleteItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return m.deleteItemFn(ctx, input, opts...)
}

func (m *mockDynamoAPI) Query(ctx context.Context, input *dynamodb.QueryInput, opts ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return m.queryFn(ctx, input, opts...)
}

func (m *mockDynamoAPI) BatchWriteItem(ctx context.Context, input *dynamodb.BatchWriteItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
	return m.batchWriteItemFn(ctx, input, opts...)
}

func (m *mockDynamoAPI) TransactWriteItems(ctx context.Context, input *dynamodb.TransactWriteItemsInput, opts ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
	return m.transactWriteItemsFn(ctx, input, opts...)
}
