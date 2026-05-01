---
inclusion: auto
name: next
description: Best practices for writing applications in Next.js. Apply this skill when writing code in the frontend.
---

## Best Practices for Next.js Development (Updated for Next.js 15)

### 1. Environment / Configuration Handling (Updated)

Next.js 15 deprecates the use of **Runtime Configuration** (`serverRuntimeConfig` and `publicRuntimeConfig`).
The official and recommended approach is now **environment variables**.

#### **1.1. Core Rules**

- **Do not** use `process.env` directly in client-side code.
  Only **server-side code** may access `process.env.*` variables.

- For environment variables that must be accessible on the client, prefix them with:
  **`NEXT_PUBLIC_`**
  These become part of the client bundle and must **never contain secrets**.

- Environment variables must be defined in:

  - `.env`, `.env.local`, `.env.production`, etc.
  - Deployment environment configuration (e.g., Vercel project settings, AWS, etc.)

#### **1.2. Server-Only Environment Variables**

Available **only on the server**:

```env
DATABASE_URL=...
SECRET_KEY=...
API_PRIVATE_TOKEN=...
```

Usage (server components, server actions, API routes, or any server-only code):

```js
const db = connect(process.env.DATABASE_URL);
```

Never expose these to the client.

#### **1.3. Client-Exposed Environment Variables**

Must start with `NEXT_PUBLIC_*`:

```env
NEXT_PUBLIC_API_BASE_URL=/api
NEXT_PUBLIC_FEATURE_FLAG=true
```

Usage in client components or browser code:

```js
const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
```

These values are embedded in the client bundle and must be safe.

#### **1.4. Example: `.env.local`**

```env
# Server-side secrets
DATABASE_URL=postgres://user:pass@host/db
EXTERNAL_API_SECRET=your-secret-token

# Client-exposed variables (safe)
NEXT_PUBLIC_API_BASE_PATH=/api
NEXT_PUBLIC_ENABLE_NEW_UI=true
```

#### **1.5. Summary**

| Type                    | Where used                                    | How defined                     | Safety                    |
| ----------------------- | --------------------------------------------- | ------------------------------- | ------------------------- |
| Server-only env vars    | Server components, server actions, API routes | `.env*`, deployment env         | **Must contain secrets**  |
| Client-exposed env vars | Client components / browser                   | Must start with `NEXT_PUBLIC_*` | **Must be non-sensitive** |

---

### 2. API / Networking Architecture

- **Never allow** frontend (client-side) code to call external APIs directly.
  It must always call **your internal endpoints** (API Routes, server actions, or backend services).

- Benefits:

  - Secrets stay server-side.
  - Centralized auth, validation, error handling.
  - Reduced client exposure to third-party errors.
  - Easier caching and performance optimization.

**Flow:**

```
Client Component → Your API Route / Server Action
    → External API (with secrets on server)
    → Process / sanitize response
    → Return to client
```

Always proxy external calls through your backend.

---

### 3. Client-side Data Fetching / State / Caching

- Use **TanStack Query (React Query)** for client-side fetching and caching.
- The fetcher function must call your **API Route** or **server action**, not an external API directly.

```js
import { useQuery } from "@tanstack/react-query";

function MyComponent() {
  const { data, error, isLoading } = useQuery(["todos"], fetchTodos);
}
```

Handle loading and errors properly, and prefer stable query keys.

---

### 4. Server-side Logic

Place sensitive logic only in server contexts:

- API Routes (`app/api/...` or `pages/api/...`)
- Server Components (App Router)
- Server Actions
- Backend service modules

Operations that **must run server-side**:

- Authentication / authorization
- Secret handling
- Database access
- External API calls with secret tokens

No secrets in client-side code.

---

### 5. Security / Data Leakage

- Secrets must **never** appear in the client bundle — keep them in `process.env.*` (without `NEXT_PUBLIC_`).
- Validate and sanitize all incoming API requests.
- Implement proper authentication and authorization in backend endpoints.

---

### 6. Performance, Build, and Rendering

- Understand which parts run on the server vs. client (App Router makes server default).
- Prefer:

  - ISR (Incremental Static Regeneration)
  - Static generation when possible
  - Server-side fetching instead of large client requests

- Avoid unnecessary client fetches.

---

### 7. Code Structure / Organization

- Organize backend logic in:

  - `app/api/` or `pages/api/`
  - Dedicated server-side services (database, external APIs, etc.)

- Keep UI separate from data-fetching logic.

**Separation of concerns:**

- UI components → render data
- Data fetching → TanStack Query (client) or server actions / API routes (server)
- Backend logic → server-side modules for external APIs, DB, secrets

---

### 8. Testing / Environment Parity

- Test API routes and server logic with mocks/stubs for external APIs.
- Ensure that all required environment variables are available in test/staging environments.
- Mock backend fetchers during client component tests.