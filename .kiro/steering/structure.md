---
inclusion: always
---

# Monorepo Project Structure

The **Ranking** product is organized as a **monorepo**.  
All code for the system lives in a single repository, divided into three top-level domains:

```

/backend   → Backend application (Go)
/frontend  → Frontend application (Next.js)
/infra     → Infrastructure as Code (AWS CDK)

```

This structure provides:

- Clear separation of responsibilities
- Single versioned source of truth for the entire system
- Atomic changes across frontend, backend, and infrastructure
- Simpler coordination between layers

---

## Root Layout

```

.
├── backend/
├── frontend/
├── infra/
├── README.md
└── .gitignore

```

The repository does **not** use workspace managers (pnpm, yarn workspaces, etc.).  
Each domain is self-contained and manages its own dependencies and tooling.

The root is responsible only for:

- High-level documentation
- Global git configuration
- CI/CD entry points

---

## `/backend` — Go Backend

Contains all server-side code written in **Go**.

Responsibilities:

- Expose HTTP APIs for:
  - Rankings
  - Items
  - Attributes
  - Ratings
  - Search and discovery
- Enforce business rules:
  - Attribute limits
  - Ownership
  - Visibility (public vs private)
- Persist and query data
- Aggregate scores for “All Users” mode
- Treat the `X-User-Id` header as the identity boundary
- Remain stateless regarding authentication

The backend is the single source of truth for all domain logic.

---

## `/frontend` — Next.js Application

Contains the full Next.js frontend.

Responsibilities:

- Render all pages:
  - Home
  - Ranking view
  - Search and tags
  - My Rankings
- Generate and persist the browser User ID in `localStorage`
- Include the User ID in every API call
- Implement:
  - SEO and metadata
  - Static generation and ISR
  - Ranking table UI
  - “All Users” and “My Scores” modes
- Handle social sharing and previews

The frontend is responsible for user experience, SEO, and zero-friction onboarding.

---

## `/infra` — AWS CDK

Contains all infrastructure code using **AWS CDK**.

Responsibilities:

- Define and provision:
  - Compute for backend
  - Hosting and CDN for frontend
  - Databases
  - Networking and security
  - Observability (logs, metrics)
- Manage environments:
  - `dev`
  - `staging`
  - `prod`
- Deploy backend and frontend artifacts
- Configure environment variables and secrets

Principles:

- No manual infrastructure
- Fully reproducible environments
- Infrastructure evolves together with application code

---

## Why a Monorepo

This structure ensures that **Ranking** evolves as a single cohesive product:

- A feature can change:
  - API contracts
  - UI
  - Infrastructure  
  in one atomic pull request.
- No version drift between layers.
- Easier onboarding for contributors.
- Clear ownership boundaries without fragmentation.

The monorepo reflects the product itself: simple, unified, and purpose-driven.