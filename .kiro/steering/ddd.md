---
inclusion: auto
name: ddd
description: Guide how to apply Domain Driven Design (DDD) in Go/Golang. Apply this skill when implementing the backend.
---

## Purpose

This document defines the rules and best practices for applying **Domain-Driven Design (DDD)** in Golang projects within this monorepo. It ensures consistent architectural decisions, code organization, and development practices across the team.

---

## Project Structure

```
backend/
  ├── .claude/
  │   └── CLAUDE.md
  ├── cmd/                # Application entry points (e.g., main.go)
  ├── internal/            # Internal packages (not to be imported by external services)
  │   ├── domain/          # Core domain logic and aggregates
  │   │   ├── <bounded-context>/
  │   │   │   ├── entity.go
  │   │   │   ├── valueobject.go
  │   │   │   ├── repository.go
  │   │   │   └── service.go
  │   ├── application/     # Use cases and orchestrators
  │   ├── infrastructure/  # External integrations (DB, APIs, etc.)
  │   └── interfaces/      # HTTP handlers, gRPC endpoints, CLI, etc.
  └── pkg/                 # Shared, generic, reusable code
```

> **Rule:** All domain code must be free of external dependencies (e.g., database, HTTP, frameworks). Only `infrastructure` may depend on external libraries.

---

## Core DDD Principles

1. **Ubiquitous Language**

   - Use consistent naming throughout code and documentation.
   - Names in the domain model must reflect business concepts.
   - Example: If the domain uses "Customer" and "Order", never alias them as `User` or `Purchase`.

2. **Bounded Contexts**

   - Each `bounded-context` must be a separate package under `internal/domain/`.
   - Dependencies **must not** cross between bounded contexts directly.
   - If communication is needed, use interfaces and application services.

3. **Entities and Value Objects**

   - **Entity:** Has an identity that persists through state changes.
   - **Value Object:** Defined by its properties, immutable, and without identity.

   ```go
   // Entity example
   type Customer struct {
       ID   uuid.UUID
       Name string
   }

   // Value Object example
   type Email struct {
       Address string
   }
   ```

4. **Domain Services**

   - Use domain services for domain logic that does **not naturally belong to a single entity or value object**.
   - Keep domain services free of infrastructure concerns.

5. **Repositories**

   - Define repository interfaces inside the domain layer.
   - Implement repository logic in the `infrastructure` layer.
   - Example:

     ```go
     // domain/repository.go
     type CustomerRepository interface {
         Save(ctx context.Context, customer Customer) error
         FindByID(ctx context.Context, id uuid.UUID) (Customer, error)
     }
     ```

6. **Application Services (Use Cases)**

   - Located in `internal/application/`.
   - Coordinate domain objects, repositories, and external systems.
   - Contain **no business rules**, only orchestration logic.

7. **Infrastructure Layer**

   - Located in `internal/infrastructure/`.
   - Contains:

     - Repository implementations
     - Database connections
     - HTTP clients
     - External integrations

   - **Rule:** The domain must not depend on infrastructure.

---

## Golang Specific Best Practices

- **Error Handling:**

  - Always return errors as `error` types, never panic.
  - Use `errors.Wrap` or `fmt.Errorf` with context messages.

- **Interfaces:**

  - Define interfaces **only where they are consumed**, not where they are implemented.
  - Keep interfaces small, usually 1-3 methods.

- **Dependency Injection:**

  - Use constructor functions for struct initialization.
  - Avoid global variables.

- **Package Visibility:**

  - Export only what is necessary.
  - Start internal functions with lowercase names.

---

## File References

- Application entry point: cmd/main.go
- Example domain entity: internal/domain/customer/entity.go
- Example repository interface: internal/domain/customer/repository.go
- Repository implementation: internal/infrastructure/customer/postgres_repository.go
- Example use case: internal/application/customer/create_customer.go

---

## Enforcement Rules for Claude Code

- Claude must **not** generate domain code that directly depends on infrastructure.
- Claude must respect the folder structure above when generating code.
- When adding a new entity or value object, Claude must:

  1. Place it inside the correct bounded context under `internal/domain/`.
  2. Generate proper unit tests.
  3. Avoid using external packages unless explicitly allowed.

- When generating repository implementations, Claude must:

  - Place them inside `internal/infrastructure/<context>/`.
  - Keep them consistent with the domain-defined repository interfaces.

- When generating API handlers, Claude must:

  - Place them inside `internal/interfaces/`.
  - Call application services instead of domain objects directly.

---

## Naming Conventions

- **Packages:** lowercase, no underscores (e.g., `customer`, `order`, `billing`)
- **Files:** lowercase, underscore-separated if needed (e.g., `repository_test.go`)
- **Structs and Types:** `PascalCase`
- **Functions:** `PascalCase` for exported, `camelCase` for internal
- **Constants:** `PascalCase`
- **Variables:** `camelCase`

---

## Summary

This `CLAUDE.md` enforces a clear separation between domain, application, infrastructure, and interfaces. Following these rules ensures that:

- Business logic is isolated and testable.
- Infrastructure can be replaced with minimal impact.
- The codebase remains maintainable as the system grows.