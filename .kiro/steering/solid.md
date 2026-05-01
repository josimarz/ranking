---
inclusion: auto
name: solid
description: Guide how to apply SOLID principles in Go/Golang. Apply this skill when writing code in the backend.
---

# SOLID Principles for Go Code

This guide outlines best practices for writing clean, maintainable Go code following the **SOLID principles**. Each principle includes an explanation, a **bad example**, and a **good example**.

---

## 1. Single Responsibility Principle (SRP)
**A class or function should have only one reason to change.**

### Bad Example ❌
```go
// Handles both user validation and saving to the database

type UserService struct {}

func (s *UserService) RegisterUser(name string, email string) error {
    if name == "" || email == "" {
        return fmt.Errorf("invalid data")
    }
    // Save user to database
    fmt.Println("User saved to DB")
    return nil
}
```

### Good Example ✅
```go
type UserValidator struct {}

func (v *UserValidator) Validate(name, email string) error {
    if name == "" || email == "" {
        return fmt.Errorf("invalid data")
    }
    return nil
}

type UserRepository struct {}

func (r *UserRepository) Save(name, email string) error {
    fmt.Println("User saved to DB")
    return nil
}

type UserService struct {
    validator   *UserValidator
    repository  *UserRepository
}

func (s *UserService) RegisterUser(name, email string) error {
    if err := s.validator.Validate(name, email); err != nil {
        return err
    }
    return s.repository.Save(name, email)
}
```
> **Why?** Each struct now has a single responsibility, making it easier to test and maintain.

---

## 2. Open/Closed Principle (OCP)
**Software entities should be open for extension but closed for modification.**

### Bad Example ❌
```go
func Discount(price float64, userType string) float64 {
    if userType == "VIP" {
        return price * 0.9
    } else if userType == "Regular" {
        return price * 0.95
    }
    return price
}
```

### Good Example ✅
```go
type DiscountStrategy interface {
    Apply(price float64) float64
}

type VIPDiscount struct{}
func (d VIPDiscount) Apply(price float64) float64 { return price * 0.9 }

type RegularDiscount struct{}
func (d RegularDiscount) Apply(price float64) float64 { return price * 0.95 }

type NoDiscount struct{}
func (d NoDiscount) Apply(price float64) float64 { return price }

func CalculateDiscount(price float64, strategy DiscountStrategy) float64 {
    return strategy.Apply(price)
}
```
> **Why?** Adding a new discount type doesn't require modifying the existing function.

---

## 3. Liskov Substitution Principle (LSP)
**Subtypes must be substitutable for their base types.**

### Bad Example ❌
```go
type Bird interface {
    Fly()
}

type Duck struct {}
func (d Duck) Fly() { fmt.Println("Duck flying") }

type Ostrich struct {}
func (o Ostrich) Fly() { panic("Ostrich can't fly!") }
```
> **Problem:** `Ostrich` violates expectations because it cannot fly.

### Good Example ✅
```go
type Bird interface {
    Eat()
}

type FlyingBird interface {
    Bird
    Fly()
}

type Duck struct {}
func (d Duck) Eat() { fmt.Println("Duck eating") }
func (d Duck) Fly() { fmt.Println("Duck flying") }

type Ostrich struct {}
func (o Ostrich) Eat() { fmt.Println("Ostrich eating") }
```
> **Why?** The `Ostrich` no longer breaks expectations since flying is separated.

---

## 4. Interface Segregation Principle (ISP)
**Clients should not be forced to depend on interfaces they do not use.**

### Bad Example ❌
```go
type Animal interface {
    Walk()
    Fly()
}

type Dog struct {}
func (d Dog) Walk() { fmt.Println("Dog walking") }
func (d Dog) Fly() { panic("Dog can't fly!") }
```

### Good Example ✅
```go
type Walker interface {
    Walk()
}

type Flyer interface {
    Fly()
}

type Dog struct {}
func (d Dog) Walk() { fmt.Println("Dog walking") }

// Bird implements both Walker and Flyer

type Bird struct {}
func (b Bird) Walk() { fmt.Println("Bird walking") }
func (b Bird) Fly()  { fmt.Println("Bird flying") }
```
> **Why?** Each interface is focused on a single behavior.

---

## 5. Dependency Inversion Principle (DIP)
**High-level modules should not depend on low-level modules. Both should depend on abstractions.**

### Bad Example ❌
```go
type MySQLRepository struct {}
func (r MySQLRepository) Save(data string) { fmt.Println("Saved to MySQL") }

type Service struct {
    repo MySQLRepository
}

func (s Service) Process(data string) {
    s.repo.Save(data)
}
```
> **Problem:** The `Service` is tightly coupled to `MySQLRepository`.

### Good Example ✅
```go
type Repository interface {
    Save(data string)
}

type MySQLRepository struct {}
func (r MySQLRepository) Save(data string) { fmt.Println("Saved to MySQL") }

type FileRepository struct {}
func (r FileRepository) Save(data string) { fmt.Println("Saved to file") }

type Service struct {
    repo Repository
}

func (s Service) Process(data string) {
    s.repo.Save(data)
}

// Usage
dbRepo := MySQLRepository{}
fileRepo := FileRepository{}

service := Service{repo: dbRepo}
service.Process("example data")
```
> **Why?** `Service` depends on the `Repository` abstraction, not a specific implementation.

---

## Summary of Best Practices
- **SRP:** Keep each type focused on one task.
- **OCP:** Use interfaces and polymorphism to extend behavior without modifying existing code.
- **LSP:** Ensure subtypes can replace their base types without unexpected behavior.
- **ISP:** Split large interfaces into smaller, specific ones.
- **DIP:** Depend on abstractions, not concrete implementations.

---

## Recommended Tools
- **golangci-lint** – for static code analysis
- **gofmt** – for formatting
- **go vet** – for identifying common mistakes