---
inclusion: auto
name: ts
description: Guide and best practices for writing TypeScript code. Apply this skill when writing code in the frontend or infra.
---

# TypeScript Code Quality Rules

This document outlines **best practices** for writing TypeScript with **strong, explicit typing**. It is designed to promote **type safety**, **readability**, and **maintainability**.

---

## 1. **Compiler Settings**

Always configure TypeScript for strict type checking:

```json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true,
    "strictBindCallApply": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "exactOptionalPropertyTypes": true,
    "alwaysStrict": true
  }
}
```

> **Why:** This ensures that all code is explicitly typed and prevents accidental `any` leaks.

---

## 2. **Types vs Interfaces**

* Use `type` for **unions, intersections, primitives**, and **complex type transformations**.
* Use `interface` for **object shapes** that may be **extended**.

```typescript
// ✅ Good
interface User {
  id: string;
  name: string;
}

type Status = "active" | "inactive";
```

> **Why:** Interfaces provide better merging and readability for object contracts.

---

## 3. **Avoid `any`**

* Never use `any`.
* If you must represent an unknown value, use `unknown` and **narrow it with type guards**.

```typescript
// ❌ Bad
function process(data: any) {
  console.log(data.name);
}

// ✅ Good
function process(data: unknown) {
  if (typeof data === "object" && data !== null && "name" in data) {
    console.log((data as { name: string }).name);
  }
}
```

---

## 4. **Prefer Explicit Types**

Always define **function parameter and return types**, even if TypeScript can infer them.

```typescript
// ❌ Bad
function sum(a, b) {
  return a + b;
}

// ✅ Good
function sum(a: number, b: number): number {
  return a + b;
}
```

> **Why:** Explicit types improve clarity and prevent accidental type mismatches.

---

## 5. **Readonly and Immutability**

Use `readonly` for values that should **never change**.

```typescript
// ✅ Good
interface Config {
  readonly apiKey: string;
  readonly timeout: number;
}

const config: Config = {
  apiKey: "12345",
  timeout: 5000,
};
```

> **Why:** This helps prevent accidental mutations.

---

## 6. **Avoid `null`, Prefer `undefined`**

* Prefer `undefined` over `null` unless interfacing with external APIs that require `null`.
* Use `strictNullChecks` to enforce proper handling.

```typescript
// ✅ Good
interface UserProfile {
  email?: string; // optional
}

function getEmail(user: UserProfile): string | undefined {
  return user.email;
}
```

---

## 7. **Narrow Types Early**

Use **type guards** and **narrowing** to avoid unsafe casting.

```typescript
function printId(id: string | number) {
  if (typeof id === "string") {
    console.log(id.toUpperCase());
  } else {
    console.log(id);
  }
}
```

> **Why:** This prevents unsafe `as` casts and runtime errors.

---

## 8. **Generics for Reusability**

Use generics to write flexible and reusable code, with constraints where necessary.

```typescript
// ✅ Good
function getFirstElement<T>(arr: T[]): T | undefined {
  return arr[0];
}

// With constraints
function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
  return obj[key];
}
```

> **Why:** Generics ensure type safety across different inputs.

---

## 9. **Avoid Overuse of `as`**

* Casting with `as` should be **last resort**.
* Prefer **type guards** or **refined type definitions**.

```typescript
// ❌ Bad
const name = (user as any).name;

// ✅ Good
if ("name" in user) {
  console.log(user.name);
}
```

---

## 10. **Use `never` for Exhaustiveness**

When using discriminated unions, leverage `never` to catch unhandled cases.

```typescript
type Shape = 
  | { kind: "circle"; radius: number }
  | { kind: "square"; size: number };

function area(shape: Shape): number {
  switch (shape.kind) {
    case "circle": return Math.PI * shape.radius ** 2;
    case "square": return shape.size * shape.size;
    default: {
      const _exhaustiveCheck: never = shape;
      return _exhaustiveCheck;
    }
  }
}
```

---

## 11. **Consistent Naming Conventions**

* **Interfaces & Types:** `PascalCase`
* **Variables & Functions:** `camelCase`
* **Constants:** `UPPER_CASE`

Example:

```typescript
interface UserProfile {
  id: string;
  fullName: string;
}

const DEFAULT_TIMEOUT = 3000;

function getUserName(user: UserProfile): string {
  return user.fullName;
}
```

---

## 12. **Avoid `object` and `{}`**

* Do not use the vague `object` or `{}` types.
* Use `Record<string, unknown>` or define exact shapes instead.

```typescript
// ❌ Bad
function process(data: object) {}

// ✅ Good
function process(data: Record<string, unknown>) {}
```

---

## 13. **Leverage Utility Types**

Use built-in utility types to avoid reinventing the wheel.

```typescript
type User = {
  id: string;
  name: string;
  email: string;
};

type PartialUser = Partial<User>;
type ReadonlyUser = Readonly<User>;
type UserWithoutEmail = Omit<User, "email">;
```

---

## 14. **Prefer Enums over Literal Strings (if needed)**

For large sets of constants, use `enum` or `as const`.

```typescript
// ✅ Good
enum Status {
  Active = "active",
  Inactive = "inactive",
}

// Alternative with const
const Statuses = {
  Active: "active",
  Inactive: "inactive",
} as const;

type StatusType = (typeof Statuses)[keyof typeof Statuses];
```

---

## 15. **Error Handling with Typed Errors**

Always type errors explicitly, avoid relying on `unknown` alone.

```typescript
interface AppError extends Error {
  code: string;
}

function throwError(): never {
  const error: AppError = { name: "AppError", message: "Something went wrong", code: "500" };
  throw error;
}
```

---

## 16. **Strict Linting**

Use `eslint` and `typescript-eslint` with strict rules:

Example `.eslintrc.json`:

```json
{
  "extends": [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:@typescript-eslint/strict"
  ],
  "rules": {
    "@typescript-eslint/explicit-function-return-type": "error",
    "@typescript-eslint/no-explicit-any": "error",
    "@typescript-eslint/no-unused-vars": ["error", { "argsIgnorePattern": "^_" }],
    "@typescript-eslint/consistent-type-imports": "error"
  }
}
```

---

## 17. **Imports and Exports**

* Use `type` imports when importing only types.
* Keep imports **clean and grouped**.

```typescript
// ✅ Good
import type { User } from "./types";

import { getUser } from "./api";
```

---

## 18. **Prefer DTOs for API Boundaries**

Define **Data Transfer Objects (DTOs)** to clearly separate **API input/output** from internal logic.

```typescript
interface CreateUserDTO {
  name: string;
  email: string;
}

interface UserEntity extends CreateUserDTO {
  id: string;
  createdAt: Date;
}
```

---

## 19. **Testing Types**

Use `tsd` or `vitest` type tests to ensure type integrity over time.

Example with `tsd`:

```typescript
import { expectType } from 'tsd';
import { getUser } from './user';

expectType<string>(getUser().name);
```

---

## 20. **Documentation and Comments**

* Add JSDoc comments for public functions and complex types.
* Keep inline comments short and meaningful.

```typescript
/**
 * Calculates the total price with tax.
 * @param price - Base price
 * @param taxRate - Tax rate as a decimal (e.g., 0.2 for 20%)
 */
function calculateTotal(price: number, taxRate: number): number {
  return price + price * taxRate;
}
```