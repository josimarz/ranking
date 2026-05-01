---
inclusion: auto
name: golang
description: Guide and best practices for write Go/Golang code. Apply this skill when writing code in the backend.
---

## Go Development Best Practices

You are an expert in Go, microservices architecture, and clean backend development practices. Your role is to ensure code is idiomatic, modular, testable, and aligned with modern best practices and design patterns.

---

### General Responsibilities:

- Guide the development of idiomatic, maintainable, and high-performance Go code.
- Enforce modular design and separation of concerns through Clean Architecture.
- Promote test-driven development (TDD), robust observability, and scalable patterns across services.

---

### Go Version:

- **Use Go version 1.25** for all projects to ensure consistency and compatibility with latest features and performance improvements.

---

### Idiomatic Go:

- Write **idiomatic Go code** following the official [Effective Go](https://go.dev/doc/effective_go) guidelines.
- Avoid unnecessary abstractions and follow common Go conventions for naming, error handling, and structuring code.

---

#### Empty Interfaces:

- **The use of `interface{}` is strictly forbidden.**

  - Always use `any` instead of `interface{}` when representing an empty interface.
  - Code containing `interface{}` **must be rejected and refactored** to use `any`.
  - Example:

    ```go
    // ✅ Correct
    func process(value any) {}

    // ❌ Incorrect - this must never be used
    func process(value interface{}) {}
    ```

---

#### String Concatenation:

- **The use of string concatenation with `+` is strictly forbidden.**

  - Always use `fmt.Sprintf` for building strings to ensure **clarity, consistency, and performance**.
  - Code containing string concatenation **must be rejected and refactored** to use `fmt.Sprintf`.
  - Example:

    ```go
    // ✅ Correct
    message := fmt.Sprintf("Hello, %s!", name)

    // ❌ Incorrect - this must never be used
    message := "Hello, " + name + "!"
    ```

---

#### Iterating Collections:

- **Prefer using `range`** for iterating over slices, arrays, and channels instead of traditional `for` loops:

  ```go
  // ✅ Good
  for i := range mySlice {
      // ...
  }

  // ❌ Bad
  for i := 0; i < len(mySlice); i++ {
      // ...
  }
  ```

---

#### Searching in Slices:

- **When checking if a slice contains a specific element, always prefer using `slices.Contains`** from Go's `slices` package instead of writing a manual `for` loop, whenever possible.
- This ensures **cleaner, more readable, and less error-prone code**.
- Example:

  ```go
  import "slices"

  func userExists(users []string, name string) bool {
      return slices.Contains(users, name)
  }
  ```

  ```go
  // ✅ Correct
  if slices.Contains(users, "Alice") {
      fmt.Println("User found!")
  }
  ```

  ```go
  // ❌ Incorrect - avoid manual loops when just checking for existence
  found := false
  for _, user := range users {
      if user == "Alice" {
          found = true
          break
      }
  }

  if found {
      fmt.Println("User found!")
  }
  ```

- Use a manual loop **only if you need additional logic beyond existence checking**, such as counting occurrences or collecting matched elements.

---

#### Linting Rules:

- **All code must strictly follow the lint rules defined in the `.golangci.yml` configuration file.**

  - Any lint violation is treated as a **blocking issue** and must be corrected before merging code.
  - The `.golangci.yml` file serves as the **single source of truth** for code style, formatting, and static analysis rules.
  - Developers must run lint checks locally before submitting a pull request using:

    ```bash
    make lint/go
    ```

    or

    ```bash
    golangci-lint run
    ```

- Example workflow:

  ```bash
  # Run linter before committing
  golangci-lint run ./...
  ```

- Code that fails linting **cannot be merged** into the main branch.
- CI/CD pipelines **must include automated lint checks** to enforce compliance.
- **No rules may be bypassed** by disabling linter checks in the code unless explicitly approved and documented in the `.golangci.yml`.

---

### Logging:

- **Must use `slog`** (Go structured logger) for all application logging.
- Logs must be **structured and JSON-formatted** for better observability and integration with systems like CloudWatch.
- Include key metadata in logs such as:

  - Request ID
  - Trace ID
  - User ID (when available)
  - Contextual information about the operation

- Use different log levels appropriately:

  - `Debug` for development and troubleshooting
  - `Info` for standard operational events
  - `Warn` for recoverable issues
  - `Error` for failures that require attention

- Example:

  ```go
  logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
  logger.Info("User login", "user_id", userID, "ip", ipAddress)
  ```

---

### REST API Development:

- When building REST APIs, use the [Gin](https://gin-gonic.com/) framework for routing and middleware.
- Organize routes clearly by grouping related endpoints.
- Use middleware for logging, authentication, and error handling.
- Return **consistent JSON responses** with appropriate HTTP status codes.
- Validate and sanitize all input data.
- Implement proper versioning for APIs (e.g., `/v1/resource`).

---

### Architecture Patterns:

- Apply **Clean Architecture** by structuring code into handlers/controllers, services/use cases, repositories/data access, and domain models.
- Use **domain-driven design** principles where applicable.
- Prioritize **interface-driven development** with explicit dependency injection.
- Prefer **composition over inheritance**; favor small, purpose-specific interfaces.
- Ensure all public functions interact with interfaces, not concrete types, to enhance flexibility and testability.

---

### Project Structure Guidelines:

- Use a consistent project layout:

  - `cmd/`: application entry points
  - `internal/`: core application logic (not exposed externally)
  - `pkg/`: shared utilities and packages
  - `api/`: gRPC/REST transport definitions and handlers
  - `configs/`: configuration schemas and loading
  - `test/`: test utilities, mocks, and integration tests

- Group code by feature when it improves clarity and cohesion.
- Keep logic decoupled from framework-specific code.

---

### Security and Resilience:

- Apply **input validation and sanitization** rigorously, especially on inputs from external sources.
- Use secure defaults for **JWT, cookies**, and configuration settings.
- Isolate sensitive operations with clear **permission boundaries**.
- Implement **retries, exponential backoff, and timeouts** on all external calls.
- Use **circuit breakers and rate limiting** for service protection.
- Consider implementing **distributed rate-limiting** to prevent abuse across services (e.g., using Redis).

**Email Validation Rule:**

- **All email validation must use the `mail.ParseAddress` function** from Go's `net/mail` package.
- Regular expressions or custom parsing functions are **strictly forbidden** for email validation.
- Example implementation:

  ```go
  import "net/mail"

  func valid(email string) bool {
      _, err := mail.ParseAddress(email)
      return err == nil
  }
  ```

---