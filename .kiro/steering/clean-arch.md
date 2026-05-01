---
inclusion: auto
name: clean-arch
description: Instructions how to implement Clean Archtecture in Go/Golang. Apply this skill when implement the backend.
---

# Clean Architecture Best Practices

This document defines the best practices for developing software following Clean Architecture principles. These rules aim to ensure maintainable, testable, and scalable code.

---

## 1. Project Structure

- Follow the **Clean Architecture layers**:
  - **Entities / Domain**: Core business logic and models.
  - **Use Cases / Application**: Business rules orchestration, application logic.
  - **Interface / Adapter / Controller**: Input and output handling (API, CLI, UI).
  - **Infrastructure / Framework**: Database, external services, and frameworks.
- Each layer should **depend only on the inner layers**, never on outer layers.
- Avoid placing domain logic in controllers, services, or repositories.

---

## 2. Code Organization

- Keep each package **small and focused** on a single responsibility.
- Name packages and files clearly to reflect their **domain or use case**.
- Group related **interfaces, adapters, and use cases** together for clarity.
- Avoid circular dependencies between packages.

---

## 3. Domain Layer

- Define **pure domain models** without framework dependencies.
- Business rules must reside in **entities** or **domain services**.
- Use **value objects** to encapsulate domain invariants and validations.
- Avoid leaking infrastructure or framework concerns into domain logic.

---

## 4. Use Case Layer

- Orchestrate application-specific business rules in **use case services**.
- Keep use cases **decoupled from infrastructure**.
- Use **interfaces for repositories and external services** to enable easy testing.
- Return **domain-specific outputs** instead of HTTP responses or framework objects.

---

## 5. Interface / Adapter Layer

- Controllers, presenters, and gateways should **translate external inputs to use case inputs**.
- Keep logic in adapters minimal; their main responsibility is **data transformation and communication**.
- Use dependency inversion: **inject interfaces instead of concrete implementations**.

---

## 6. Infrastructure Layer

- Implement repositories, external service clients, and database migrations here.
- Avoid business logic in infrastructure; it should only handle **technical concerns**.
- Use migration tools consistently (e.g., `golang-migrate` for database migrations in Go projects).

---

## 7. Dependency Management

- Inner layers **should never depend on outer layers**.
- Use dependency injection to provide concrete implementations at runtime.
- Keep dependencies **lightweight and focused on abstractions**.

---

## 8. Testing

- Test **domain logic and use cases** in isolation.
- Mock external dependencies in use case and controller tests.
- Aim for **high coverage on business rules** while avoiding over-testing framework code.

---

## 9. Error Handling

- Return errors in **domain-meaningful terms**.
- Map errors to framework-specific responses in the adapter layer.
- Avoid panic or swallowing errors silently.

---

## 10. Documentation and Comments

- Document **public APIs, domain models, and use cases**.
- Avoid redundant comments; code should be **self-explanatory whenever possible**.
- Keep README and architectural diagrams up to date.

---

## 11. General Best Practices

- Follow **SOLID principles** consistently.
- Keep functions small and focused on a **single responsibility**.
- Avoid business logic duplication; reuse domain services or value objects.
- Apply **consistent naming conventions** and code formatting.
- Review code for **clean architecture violations** before merging.