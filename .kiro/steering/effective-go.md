---
inclusion: auto
name: effective-go
description: Best practices to write idiomatic Go/Golang code. Apply this skill when writing Go/Golang code.
---

# Effective Go — Writing Clear, Idiomatic Go Code

_Based on the official “Effective Go” document from go.dev_  
https://go.dev/doc/effective_go

---

## Introduction

Go is a language designed with simplicity, clarity, and efficiency in mind. Writing effective Go code means embracing Go’s conventions and idioms rather than directly translating patterns from other languages such as Java, C++, or Python.

This document consolidates best practices from the official _Effective Go_ guide and complements them with commonly adopted Go community practices. It focuses on writing readable, maintainable, testable, and idiomatic Go code.

---

## 1. Formatting

- **Use `gofmt`**: Always format your code using `gofmt` or `go fmt`.
- **Avoid manual formatting**: Do not manually align code or adjust indentation.
- **Standardized style**:
  - Tabs are used for indentation.
  - No enforced maximum line length.
  - Parentheses are used only when required by the grammar.

> Do not fight the formatter—embrace it.

---

## 2. Naming Conventions

### Package Names

- Use **short, lowercase names**.
- Avoid underscores and mixedCaps.
- Package names should describe _what the package provides_, not what it contains.

Examples:

```go
import "bytes"
import "net/http"
```

### Identifiers

- Use **MixedCaps** (camel case) for variables, functions, types, and methods.
- Avoid underscores in identifiers.

### Getters and Setters

- Do **not** prefix getters with `Get`.

```go
func (u *User) Name() string
```

- Use `SetXxx` for setters when appropriate.

```go
func (u *User) SetName(name string)
```

### Interfaces

- Interfaces with a single method typically end with `-er`:

  - `Reader`
  - `Writer`
  - `Closer`

This convention improves readability and aligns with the standard library.

---

## 3. Commentary and Documentation

- Use `//` for line comments and `/* */` for block comments.
- All **exported** identifiers should have documentation comments.
- Comments should:

  - Be written in complete sentences.
  - Start with the name of the declared item.
  - Explain **why**, not just **what**.

Example:

```go
// ParseConfig reads and validates the application configuration.
func ParseConfig(path string) (*Config, error) { ... }
```

---

## 4. Code Structure

- **Keep functions small and focused**.
- Each function should do **one thing well**.
- Avoid unnecessary nesting; return early when possible.
- Group related types and functions together in the same file.
- Prefer clarity over cleverness.

---

## 5. Semicolons

- Go inserts semicolons automatically during parsing.
- Do not write semicolons explicitly.
- The opening brace `{` must be on the same line as the statement.

Correct:

```go
if err != nil {
    return err
}
```

---

## 6. Control Structures

### If Statements

- Always use braces.
- Support optional initialization:

```go
if err := validate(); err != nil {
    return err
}
```

### For Loops

Go has a single looping construct: `for`.

```go
// Classic
for i := 0; i < n; i++ {}

// While-style
for condition {}

// Infinite
for {}
```

### Range

Use `range` to iterate over collections:

```go
for i, v := range items {
    fmt.Println(i, v)
}
```

---

## 7. Functions

### Multiple Return Values

Go functions can return multiple values, commonly used for error handling:

```go
result, err := compute()
if err != nil {
    return err
}
```

### Short Variable Declarations

Use `:=` when declaring variables inside functions for the first time:

```go
count := 10
```

Avoid overusing it when explicitness improves clarity.

---

## 8. Error Handling

- Errors are values.
- Always check errors explicitly.
- Handle errors immediately after they occur.
- Decide whether to return, wrap, log, or recover from errors based on context.

```go
if err != nil {
    return err
}
```

### Error Inspection

- Use `errors.Is` and `errors.As` for error comparison and unwrapping.
- Avoid direct equality checks for wrapped errors.

---

## 9. Data Types and Interfaces

### Interfaces

- Interfaces define behavior, not implementation.
- Prefer **small interfaces**.
- Types satisfy interfaces implicitly—no explicit declaration is required.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

### Embedding

Use embedding for composition:

```go
type ReadWriter interface {
    Reader
    Writer
}
```

Embedding promotes reuse without inheritance.

---

## 10. Concurrency

- Use **goroutines** to perform concurrent tasks.
- Use **channels** to synchronize and communicate between goroutines.
- Prefer **communication over shared memory**.
- Minimize shared state and protect it when necessary.

```go
go process(data)
```

---

## 11. Testing

- Write tests for all critical logic.
- Prefer **table-driven tests** for multiple scenarios.
- Test edge cases:

  - Empty inputs
  - Nil values
  - Boundary conditions
  - Large inputs

Example pattern:

```go
tests := []struct {
    name string
    input int
    want  int
}{ ... }
```

---

## 12. Performance

- Profile before optimizing.
- Avoid premature optimization.
- Prefer clarity and correctness first.
- Use Go’s built-in types and standard library, which are highly optimized.
- Measure improvements using benchmarks.

---

## 13. Complete Example

The official _Effective Go_ guide concludes with a complete web server example demonstrating:

- Package organization
- HTTP handlers
- Templates
- Flags
- Idiomatic error handling

It showcases how Go’s features work together in real applications.

---

## Conclusion

_Effective Go_ is not a strict rulebook but a collection of proven practices shaped by the Go community. Writing effective Go code means prioritizing simplicity, readability, explicitness, and consistency.

For deeper understanding, also consult:

- The Go Language Specification
- A Tour of Go
- Go standard library source code

---