---
inclusion: auto
name: gin
description: Guide and best practices to write code with Gin-Gonic framework. Apply this skill when writing code in the backend.
---

# Gin-Gonic best practices

## Purpose

This document outlines best practices for developing backend applications using the **Gin-Gonic** framework. It covers project organization, middleware design, logging strategies, and deployment considerations. These rules help maintain code quality, performance, and readability across development and production environments.

---

## Project Structure

```
backend/
  ├── cmd/
  │   └── server/
  │       └── main.go             # Entry point for the Gin application
  ├── internal/
  │   ├── interfaces/
  │   │   └── http/
  │   │       ├── handlers/       # Route handlers
  │   │       ├── middleware/     # Custom middlewares
  │   │       └── router.go       # Gin router initialization
  │   ├── application/            # Use cases / business orchestration
  │   ├── domain/                  # Core domain logic
  │   └── infrastructure/         # Database, external services, etc.
  └── pkg/
      └── logger/                  # Custom logger utilities
```

> **Rule:** The `internal/interfaces/http/` folder contains **only HTTP-related code**. Domain and application layers must **not depend** on Gin.

---

## Logging Best Practices

### 1. Colored Logs in Local Development

* When running locally, logs should be **human-friendly** and **colorized** for better readability.
* This can be achieved by enabling `gin.ForceConsoleColor()` and using Gin's default logger.

Example setup for local environment:

```go
if os.Getenv("APP_ENV") == "local" {
    gin.ForceConsoleColor()
    router.Use(gin.Logger())
} else {
    router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("[%s] %d | %s | %s\n",
            param.TimeStamp.Format(time.RFC3339),
            param.StatusCode,
            param.Method,
            param.Path,
        )
    }))
}
```

### 2. Plain Logs in Production

* In production, **avoid colors** to keep logs structured and suitable for log aggregators like Loki, ELK, or CloudWatch.
* Use `LoggerWithFormatter` to define a custom, JSON-friendly, or structured log format.

Example JSON-style production logs:

```go
router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
    return fmt.Sprintf(
        `{"timestamp":"%s","status":%d,"method":"%s","path":"%s","latency":"%s"}\n`,
        param.TimeStamp.Format(time.RFC3339),
        param.StatusCode,
        param.Method,
        param.Path,
        param.Latency,
    )
}))
```

### 3. Centralized Logger

* Create a custom logger in `pkg/logger` to standardize log behavior across environments.
* Example utility:

```go
package logger

import (
    "log"
    "os"
)

type Logger struct {
    *log.Logger
}

func NewLogger() *Logger {
    return &Logger{Logger: log.New(os.Stdout, "", log.LstdFlags)}
}

func (l *Logger) Info(msg string) {
    l.Println("INFO:", msg)
}

func (l *Logger) Error(msg string) {
    l.Println("ERROR:", msg)
}
```

> **Recommendation:** Always log in a **structured format** in production for easy parsing and filtering.

---

## Middleware Best Practices

1. **Custom Middleware**

   * Place all custom middleware in `internal/interfaces/http/middleware/`.
   * Keep middleware focused on **cross-cutting concerns**, such as:

     * Authentication & Authorization
     * Request/Response Logging
     * Rate Limiting
     * Request Validation

2. **Error Handling Middleware**

   * Centralize error handling using a global middleware to return consistent error responses.

Example:

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) > 0 {
            c.JSON(-1, gin.H{
                "error": c.Errors[0].Error(),
            })
        }
    }
}
```

---

## Router Configuration

* Define routes in `router.go`.
* Group routes by bounded context or module.
* Example:

```go
func SetupRouter() *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery())

    api := r.Group("/api/v1")
    {
        api.GET("/health", HealthHandler)
        api.POST("/users", CreateUserHandler)
    }

    return r
}
```

> **Rule:** Never register routes directly inside `main.go`. Always use `router.go`.

---

## Environment Configuration

* Control environment-specific behavior with the `APP_ENV` variable:

  * `local`: Colorized logs, detailed error messages.
  * `staging`: Structured logs, debug-level messages.
  * `production`: Structured logs, minimal output.

Example `.env`:

```
APP_ENV=local
PORT=8080
```

Load environment variables using `github.com/joho/godotenv` or similar library.

---

## Deployment Considerations

* **Graceful Shutdown:** Implement proper shutdown to close connections and clean resources.

```go
srv := &http.Server{
    Addr:    ":8080",
    Handler: router,
}

// Start server in a goroutine
if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
    log.Fatalf("listen: %s\n", err)
}

// Handle shutdown signals
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

---

## Summary

By following these practices, Gin-based backend services will be:

* Maintainable and modular
* Consistent across environments
* Ready for production logging and monitoring
* Easy to debug in local development with colorized logs

This ensures a smooth transition from local development to production deployment.