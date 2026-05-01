# Assumptions — Backend Requirements

> These assumptions are derived directly from the steering documents (product, tech, infra, auth, DynamoDB, REST API).
> They do NOT require user confirmation — they are stated facts from the steering docs.
> Assumptions derived from user answers to `questions.md` will be added below as they are resolved.

---

## Confirmed Assumptions (from Steering Documents)

### A1: Language and Framework

The backend is implemented in **Go 1.25** using the **Gin-Gonic** framework. No other language or framework is permitted.

_Source: tech.md §1.1_

---

### A2: Runtime Environment

The backend runs on **AWS Lambda** (ARM64, Amazon Linux 2023, 128 MB memory) behind **Amazon API Gateway**.

_Source: tech.md §1.3_

---

### A3: No Authentication

There is no traditional authentication (login, password, OAuth, sessions). Users are identified by a **UUID v7** stored in the browser's `localStorage` and sent via the `X-User-Id` HTTP header.

_Source: auth.md, product.md §11_

---

### A4: Single DynamoDB Table

All entities (rankings, attributes, items, ratings, user references) are stored in a **single DynamoDB table** using the Single Table Design pattern. Multiple tables are explicitly forbidden.

_Source: tech.md §2.2, dynamodb.md_

---

### A5: S3 for Images

All item images are stored in **Amazon S3**. The backend is responsible for downloading (from URL), receiving (file upload), optimizing, generating thumbnails, and serving images via **pre-signed URLs**. Public S3 access is forbidden.

_Source: tech.md §2.3_

---

### A6: Stateless API

The backend is fully stateless. All persistence is handled via DynamoDB and S3. No in-memory state, sessions, or caches are maintained between requests.

_Source: tech.md §1.1_

---

### A7: API Key Protection in Production

The REST API is protected by an **API Key** in production (via API Gateway). In local development, the API Key is disabled.

_Source: tech.md §1.4_

---

### A8: Swagger/OpenAPI Documentation

All endpoints must be documented using **swaggo/gin-swagger** annotations. An undocumented endpoint is considered incomplete. Swagger UI is exposed in all environments (protected by API Key in production).

_Source: tech.md §1.2_

---

### A9: Rating Scale

Ratings are on a scale of **0 to 100** per attribute per item per user.

_Source: product.md §6_

---

### A10: Rating Save Rule

Ratings for an item are saved **only when the user fills all attributes** for that item. Partial ratings are not persisted.

_Source: product.md §6_

---

### A11: Overall Score Calculation

The **Overall** score is the arithmetic mean of all attribute scores. It is **never stored** — always calculated at read time.

_Source: product.md §5_

---

### A12: Maximum 6 Attributes per Ranking

A ranking can have at most **6 attributes**. This is a hard business rule.

_Source: product.md §3_

---

### A13: Visibility Model

Rankings are either **public** (discoverable, searchable, listed on home) or **private** (accessible only via direct URL with the ranking UUID).

_Source: product.md §8_

---

### A14: UUID v7 for IDs

All entity IDs (rankings, items, users) use **UUID v7** format, which provides time-ordered uniqueness.

_Source: auth.md §1, product.md §8_

---

### A15: Sorting Convention

Sorting follows the **JSON:API specification**: `sort` query parameter, `-` prefix for descending, comma-separated for multiple fields.

_Source: rest-api.md (Sorting Convention section)_

---

### A16: DynamoDB Access Pattern First

All DynamoDB schema design must be driven by access patterns. Every feature must answer: "What exact Query will this feature execute?" No Scans in production code.

_Source: dynamodb.md_

---

### A17: Infrastructure as Code

All infrastructure is provisioned via **AWS CDK (TypeScript)**. No manual AWS Console configuration.

_Source: tech.md §3, infra.md_

---

### A18: Local Development with LocalStack

Local development uses **LocalStack** to simulate DynamoDB and S3. The backend must run fully offline.

_Source: tech.md §2.4, infra.md §5_

---

### A19: Two View Modes

The ranking table supports two modes: **"Average of all users"** (aggregated scores) and **"My scores"** (individual user scores). The backend must support both.

_Source: product.md §6_

---

### A20: Home Page Discovery

The home page displays the **10 most recently created public rankings**. Users can search rankings by name (partial) or by tags.

_Source: product.md §9_

---

### A21: Item Default Image

If an item has no photo, the frontend displays a default image with the item's initials. This is a frontend concern — the backend simply returns `null` or empty for the image field.

_Source: product.md §4_

---

### A22: Clean Architecture

The backend follows **Clean Architecture** with DDD principles: domain layer (entities, value objects, repository interfaces), application layer (use cases), infrastructure layer (DynamoDB, S3 implementations), and interfaces layer (HTTP handlers).

_Source: ddd.md, clean-arch.md_

---

### A23: Structured JSON Logging

All logging uses **slog** with JSON-formatted structured logs including requestId, env, service, and operation fields.

_Source: golang.md, infra.md §8_

---

### A24: On-Demand DynamoDB Billing

DynamoDB uses **PAY_PER_REQUEST** (on-demand) billing mode.

_Source: infra.md §4_

---

---

## Assumptions from User Answers

### A25: Maximum 10 Tags per Ranking

The backend enforces a maximum of **10 tags** per ranking.

_Source: Question [1] — Option A_

---

### A26: Full Ranking Metadata Editing

The ranking creator can edit **all metadata** (name, description, tags, visibility) at any time after creation.

_Source: Question [2] — Option A_

---

### A27: Ranking Deletion (Cascade)

The ranking creator can **delete a ranking entirely**. Deletion cascades to all associated items and ratings.

_Source: Question [3] — Option A_

---

### A28: Ranking Field Lengths

Ranking name: max **100 characters**. Ranking description: max **500 characters**.

_Source: Question [4] — Option A_

---

### A29: Attribute Modification After Ratings

Attributes can be **added, renamed, or removed freely** even after ratings exist. Existing ratings for removed attributes are **soft-deleted** (kept in DB but excluded from calculations).

_Source: Question [5] — Option A_

---

### A30: Attribute Field Lengths

Attribute name: max **50 characters**. Attribute description: max **200 characters**.

_Source: Question [6] — Option A_

---

### A31: Item Edit/Delete Permissions

The **ranking creator** can edit or delete any item. The **item creator** can edit or delete their own items.

_Source: Question [7] — Option A_

---

### A32: Maximum 100 Items per Ranking

A ranking can have at most **100 items**.

_Source: Question [8] — Option A_

---

### A33: Item Name Length

Item name: max **100 characters**.

_Source: Question [9] — Option A_

---

### A34: Accepted Image Formats

The backend accepts **JPEG, PNG, and WebP** image formats only.

_Source: Question [10] — Option A_

---

### A35: Maximum Image File Size

Maximum upload file size is **10 MB**.

_Source: Question [11] — Option B_

---

### A36: Image Optimization Dimensions

Optimized original: max **800×800px**, WebP format, quality 80. Thumbnail: **150×150px**, WebP format, quality 75.

_Source: Question [12] — Option A_

---

### A37: Pre-Signed URL Validity

Pre-signed URLs are valid for **1 hour**.

_Source: Question [13] — Option A_

---

### A38: Rating Updates Allowed

Users can **update their ratings** at any time by re-submitting all attribute scores for an item.

_Source: Question [14] — Option A_

---

### A39: Integer-Only Ratings

Rating values must be **integers** in the range 0–100. Decimals are not accepted.

_Source: Question [15] — Option A_

---

### A40: Case and Accent Insensitive Search

Search is **case-insensitive and accent-insensitive** (e.g., "video" matches "Vídeo").

_Source: Question [16] — Option A_

---

### A41: Dedicated Recent Rankings Endpoint

The home page uses a dedicated endpoint `GET /api/v1/rankings/recent` returning the **10 most recently created public rankings**.

_Source: Question [17] — Option A_

---

### A42: Exact Tag Match

Tag search uses **exact match only**. Tags are discrete keywords.

_Source: Question [18] — Option A_

---

### A43: Cursor-Based Pagination

The backend uses **cursor-based pagination** leveraging DynamoDB's `LastEvaluatedKey`.

_Source: Question [19] — Option A_

---

### A44: Page Size Defaults

Default page size: **20**. Maximum page size: **50**.

_Source: Question [20] — Option A_

---

### A45: API Version Prefix

All routes are prefixed with **`/api/v1/`**.

_Source: Question [21] — Option A_

---

### A46: X-User-Id Header Policy

`X-User-Id` is **required** for write operations (create, update, delete, rate) and "My Rankings". It is **optional** for read-only public endpoints (list, get, search). Missing or invalid header on required endpoints returns `400 Bad Request`.

_Source: Question [22] — Option A_

---

### A47: Backend Calculates Overall

The backend **calculates and returns the `overall` score** in API responses for both "all users average" and "my scores" modes.

_Source: Question [23] — Option A_

---

### A48: Structured Error Response

API errors use a structured JSON format: `{ "error": { "code": "<ERROR_CODE>", "message": "<human-readable message>", "status": <HTTP status> } }`.

_Source: Question [24] — Option A_

---

### A49: Sortable Fields — Rankings Listing

Rankings listing supports sorting by: **`name`**, **`createdAt`**, **`updatedAt`**.

_Source: Question [25] — Option A_

---

### A50: Sortable Fields — Items in Ranking

Items within a ranking support sorting by: **`name`**, **`overall`**, and **each individual attribute score**.

_Source: Question [26] — Option A_

---

### A51: Item Creator Tracking

Each item tracks its **creator** via a `createdBy` field (userId).

_Source: Question [27] — Option A_

---

### A52: Open Item Addition

Any user can add items to any **public** ranking. For **private** rankings, any user with the ranking URL (ID) can add items.

_Source: Question [28] — Option A_

---

### A53: Health Check Endpoint

The backend exposes `GET /api/v1/health` returning `200 OK` with version, environment, and DynamoDB connectivity status.

_Source: Question [29] — Option A_
