---
inclusion: auto
name: aws-sdk
description: Guide and best practices when using AWS SDK in Go/Golang to interact with AWS resources. Apply this skill when writing code in the backend.
---

# AWS SDK for Go v2 - Best Practices

This document defines best practices for working with the **AWS SDK for Go v2** to ensure maintainability, readability, and consistency across the codebase.

---

## 1. General Coding Practices

- Always handle AWS errors explicitly using the `error` return value.
- Use **context-aware** SDK functions (e.g., `svc.SomeOperation(ctx, input)`) to support timeouts and cancellations.
- Prefer **structured logging** with request IDs and relevant AWS metadata for better observability.
- Avoid hardcoding AWS region, credentials, or endpoints.
  - Use environment variables, AWS CLI profiles, or IAM roles whenever possible.
  - Example:
    ```go
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        log.Fatalf("unable to load AWS SDK config, %v", err)
    }
    ```

---

## 2. Using Pointer Helper Functions

The AWS SDK for Go v2 uses pointers for most request and response fields.  
**Always use the built-in helper functions** defined in [`to_ptr.go`](https://github.com/aws/aws-sdk-go-v2/blob/main/aws/to_ptr.go) instead of manually taking addresses of values.

This improves code readability and consistency, and prevents common mistakes such as taking the address of temporary variables.

### Example

**✅ Good**

```go
import "github.com/aws/aws-sdk-go-v2/aws"

input := &dynamodb.PutItemInput{
    TableName: aws.String("MyTable"),
    Item: map[string]types.AttributeValue{
        "ID":   &types.AttributeValueMemberS{Value: "123"},
        "Name": &types.AttributeValueMemberS{Value: "John Doe"},
    },
}
```

**❌ Bad**

```go
input := &dynamodb.PutItemInput{
    TableName: &("MyTable"), // Avoid taking address directly
}
```

### Common Helper Functions

| Type        | Function          |
| ----------- | ----------------- |
| `string`    | `aws.String(v)`   |
| `bool`      | `aws.Bool(v)`     |
| `int`       | `aws.Int(v)`      |
| `int64`     | `aws.Int64(v)`    |
| `float64`   | `aws.Float64(v)`  |
| `time.Time` | `aws.Time(v)`     |
| `duration`  | `aws.Duration(v)` |

> **Note:**
> These functions also exist for slices and maps, e.g., `aws.StringSlice()`, `aws.BoolMap()`.
> Always prefer these over manual loops when converting data structures.

---

## 3. Handling Pagination

Many AWS SDK operations return paginated results.
**Always use the paginator utilities** provided by the SDK instead of manually managing tokens.

**✅ Good**

```go
paginator := dynamodb.NewScanPaginator(svc, &dynamodb.ScanInput{
    TableName: aws.String("MyTable"),
})

for paginator.HasMorePages() {
    page, err := paginator.NextPage(ctx)
    if err != nil {
        log.Fatalf("failed to retrieve page: %v", err)
    }
    fmt.Println("Items:", page.Items)
}
```

**❌ Bad**

```go
// Manual handling of NextToken - avoid this
output, err := svc.Scan(ctx, &dynamodb.ScanInput{
    TableName: aws.String("MyTable"),
})
```

---

## 4. Timeouts and Retries

- Use **context with timeout or deadline** to avoid hanging requests.
- Configure **retry behavior** when creating the AWS config if needed.

**Example:**

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

cfg, err := config.LoadDefaultConfig(ctx, config.WithRetryer(func() aws.Retryer {
    return retry.NewStandard(func(o *retry.StandardOptions) {
        o.MaxAttempts = 5
    })
}))
```

---

## 5. Testing with Mocks and Testcontainers

- **Unit tests** should use mocks for AWS clients.

  - The AWS SDK v2 interfaces are designed to make mocking easy.

- **Integration tests** should use [Testcontainers](https://testcontainers.com/) to spin up local AWS-compatible services like **LocalStack**.

**Example: Using LocalStack with Testcontainers**

```go
func setupLocalStack(t *testing.T) string {
    ctx := context.Background()
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image:        "localstack/localstack:latest",
            ExposedPorts: []string{"4566/tcp"},
            Env: map[string]string{
                "SERVICES": "dynamodb,s3",
            },
        },
        Started: true,
    })
    require.NoError(t, err)

    endpoint, err := container.Endpoint(ctx, "")
    require.NoError(t, err)
    return endpoint
}
```

> **Recommendation:**
> Run LocalStack in tests for realistic integration scenarios, but keep tests fast and isolated.

---

## 6. Resource Cleanup

- Always clean up resources such as S3 buckets, DynamoDB tables, or containers after tests.
- Use `defer` to ensure cleanup even in case of errors.

**Example:**

```go
bucket := "my-test-bucket"
_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
    Bucket: aws.String(bucket),
})
if err != nil {
    t.Fatalf("failed to create bucket: %v", err)
}
defer s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
    Bucket: aws.String(bucket),
})
```

---

## Summary

1. Always use pointer helpers (`aws.String()`, `aws.Int64()`, etc.).
2. Use context-aware calls with proper timeouts.
3. Handle pagination with SDK-provided paginators.
4. Mock AWS clients for unit tests, and use Testcontainers + LocalStack for integration tests.
5. Clean up resources in tests using `defer`.
6. Never hardcode credentials or regions.

Following these guidelines will ensure consistent, maintainable, and robust AWS integrations in Go.