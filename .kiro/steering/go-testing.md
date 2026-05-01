---
inclusion: auto
name: go-testing
description: Best practices for writing test code in Go/Golang. Apply this skill when writing tests in the backend.
---

# Go Testing Best Practices

This document outlines best practices for writing tests in Go. These rules are intended to improve test readability, maintainability, and reliability.

---

## 1. Prefer Table-Driven Tests

- Use table-driven tests to cover multiple test cases in a single test function.
- This improves readability and reduces duplication.

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive numbers", 1, 2, 3},
        {"negative numbers", -1, -2, -3},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("got %d, want %d", got, tt.want)
            }
        })
    }
}
````

---

## 2. Use `t.Parallel()` for Independent Tests

* Run tests in parallel when they do not share state to speed up test execution.

```go
func TestSomething(t *testing.T) {
    t.Parallel()
    // test logic here
}
```

---

## 3. Use `require` or `assert` from `testify` for Cleaner Assertions

* Using `require` or `assert` improves readability and makes failures easier to diagnose.

```go
import (
    "testing"
    "github.com/stretchr/testify/require"
)

func TestSomething(t *testing.T) {
    result := DoSomething()
    require.Equal(t, expected, result)
}
```

---

## 4. Keep Unit Tests Small and Focused

* Each unit test should test a single behavior.
* Avoid testing multiple unrelated functionalities in the same test.

---

## 5. Use Mocking for External Dependencies

* Mock external services (HTTP, databases, etc.) for unit tests.
* Avoid making real network calls in unit tests.

---

## 6. Use Testcontainers for Integration Tests

* For integration tests requiring external dependencies (databases, message queues, etc.), use **Testcontainers**.
* This ensures tests are reproducible and isolated from the local environment.

```go
import (
    "context"
    "testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func TestDatabaseIntegration(t *testing.T) {
    ctx := context.Background()
    
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "user",
            "POSTGRES_PASSWORD": "password",
            "POSTGRES_DB":       "testdb",
        },
        WaitingFor: wait.ForListeningPort("5432/tcp"),
    }

    postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    require.NoError(t, err)
    defer postgresC.Terminate(ctx)

    host, err := postgresC.Host(ctx)
    require.NoError(t, err)

    port, err := postgresC.MappedPort(ctx, "5432")
    require.NoError(t, err)

    // Connect to DB using host:port and run integration tests
}
```

---

## 7. Name Tests Clearly

* Use descriptive names for test functions to clearly indicate what is being tested.

```go
func TestCalculateDiscount_ForVIPCustomer_ReturnsCorrectValue(t *testing.T) {
    // ...
}
```

---

## 8. Avoid Hard-Coding Values

* Use constants or variables instead of magic numbers or strings.
* Makes tests easier to maintain.

---

## 9. Clean Up Resources

* Ensure tests clean up any resources they create (temporary files, containers, database entries).
* This prevents side effects between tests.

---

## 10. Keep Tests Fast

* Unit tests should run quickly.
* Integration tests may be slower but should still aim for minimal execution time.

---