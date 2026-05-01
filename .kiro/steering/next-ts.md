---
inclusion: auto
name: next-ts
description: Best practices for writing TypeScript code with Next.js. Apply this skill when writing code in the frontend.
---

# TypeScript Best Practices for React + Next.js

## Core Principles

When Claude generates or reviews TypeScript code for React/Next.js, it must follow these core principles:

1. **Type safety is non-negotiable**

   - Always prefer **strict typing**.
   - Avoid `any` at all costs. If absolutely unavoidable, use `unknown` instead and refine it through type guards.
   - Enable strict TypeScript options in `tsconfig.json`:

     ```json
     {
       "compilerOptions": {
         "strict": true,
         "noImplicitAny": true,
         "strictNullChecks": true,
         "noUnusedLocals": true,
         "noUnusedParameters": true,
         "noFallthroughCasesInSwitch": true
       }
     }
     ```

2. **Types are documentation**

   - Types should describe **intent and behavior**, making the code self-documenting.
   - Prefer clear, descriptive type names over terse or generic ones.

3. **Type inference over explicitness**

   - Let TypeScript infer types when it’s obvious, especially for local variables.
   - Explicitly type **function parameters**, **return values**, and **public APIs** (like exported functions, props, and hooks).

---

## File and Project Organization

| Folder            | Purpose                                              |
| ----------------- | ---------------------------------------------------- |
| `types/`          | Global and shared type definitions (`types.ts`)      |
| `src/components/` | React components with `*.tsx` files                  |
| `src/hooks/`      | Custom hooks with clear typings                      |
| `src/lib/`        | Utility functions and services                       |
| `src/context/`    | Context definitions with `ContextType` and providers |

**Rules:**

- Use `*.tsx` for all React components.
- Centralize shared interfaces and enums in `types/` or a nearby `types.ts` file.
- Avoid cyclic dependencies by structuring types hierarchically.

---

## Props and Component Typing

Always define component props with `interface` or `type` aliases.

### Example:

```tsx
interface ButtonProps {
  label: string;
  onClick: () => void;
  disabled?: boolean;
  variant?: "primary" | "secondary";
}

export function Button({
  label,
  onClick,
  disabled = false,
  variant = "primary",
}: ButtonProps) {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`btn btn-${variant}`}
    >
      {label}
    </button>
  );
}
```

---

## Avoiding `any` — Preferred Alternatives

| Situation                 | Preferred Type Instead of `any` |
| ------------------------- | ------------------------------- |
| Unknown external data     | `unknown`                       |
| Nullable value            | `string \| null` or `string?`   |
| Complex objects           | Define an explicit interface    |
| Dynamic keys              | `Record<string, T>`             |
| Values constrained to set | `union types` (`'a' \| 'b'`)    |
| Arrays of known shape     | `T[]` or `Array<T>`             |

---

## Avoid `console` for Debugging

**Bad Example:**

```ts
const data = await fetchData();
console.log("Fetched data:", data);
```

**Good Example:**

- Use proper logging or error handling instead of `console`.
- In production, remove all `console` calls.

```ts
try {
  const data = await fetchData();
  // Handle or store data instead of console.log
  updateState(data);
} catch (err) {
  handleError(err);
}
```

---

## Avoid Declaring Unused Error Variables in try…catch

**Bad Example:**

```ts
try {
  await doSomething();
} catch (err) {
  // error variable declared but not used
}
```

**Good Example:**

```ts
try {
  await doSomething();
} catch {
  // no unused variable declared
  handleError();
}
```

Or, if the error is needed:

```ts
try {
  await doSomething();
} catch (err) {
  console.error(err); // properly using the error variable
}
```

---

## API Response and Data Typing

When handling API calls, always type both **request parameters** and **response data**.

### Example: Fetching API data

```ts
// types/api.ts
export interface User {
  id: string;
  name: string;
  email: string;
}

export interface ApiError {
  message: string;
  code: number;
}
```

```ts
// src/services/userService.ts
import { User } from "@/types/api";

export async function getUsers(): Promise<User[]> {
  const response = await fetch("/api/users");
  if (!response.ok) throw new Error("Failed to fetch users");
  return response.json() as Promise<User[]>;
}
```

---

## Next.js Specific Typing

### `getServerSideProps`

```ts
import { GetServerSideProps } from "next";

interface PageProps {
  users: User[];
}

export const getServerSideProps: GetServerSideProps<PageProps> = async () => {
  const users = await getUsers();
  return { props: { users } };
};
```

---

### `getStaticProps` and `getStaticPaths`

```ts
import { GetStaticProps, GetStaticPaths } from "next";

interface Post {
  id: string;
  title: string;
}

interface Props {
  post: Post;
}

export const getStaticProps: GetStaticProps<Props> = async (context) => {
  const post = await fetchPost(context.params?.id as string);
  return { props: { post } };
};

export const getStaticPaths: GetStaticPaths = async () => {
  const posts = await fetchAllPosts();
  return {
    paths: posts.map((p) => ({ params: { id: p.id } })),
    fallback: false,
  };
};
```

---

## API Route Typing

```ts
import type { NextApiRequest, NextApiResponse } from "next";
import { User } from "@/types/api";

export default function handler(
  req: NextApiRequest,
  res: NextApiResponse<User | { error: string }>
) {
  if (req.method === "GET") {
    const user: User = { id: "1", name: "Alice", email: "alice@example.com" };
    res.status(200).json(user);
  } else {
    res.status(405).json({ error: "Method not allowed" });
  }
}
```

---

## Summary of Rules for Claude Code

1. **Never use `any`** — prefer `unknown` or explicit interfaces.
2. Use **strict compiler options** and enforce them in CI.
3. Always type public APIs: function params, return types, props, hooks.
4. Type all **Next.js APIs** (`getServerSideProps`, `API routes`, etc.).
5. Co-locate types near usage or in a shared `types/` folder.
6. Use **union types** or **enums** for fixed values instead of magic strings.
7. Ensure API responses and requests are strongly typed.
8. Leverage **generics** for reusable functions and hooks.
9. Use context safely by defining full types for state and actions.
10. Let TypeScript **infer** local variable types when obvious but **never** omit types for public interfaces.
11. **Do not use `console`** for debugging; use proper logging or error handling.
12. **Do not declare unused error variables** in `try...catch` blocks.