# Questions — Backend Requirements Clarification

> **Instructions:** For each question, mark your choice with an `x` inside the checkbox (e.g., `[x]`).
> If you select **Other**, please describe your choice below the option.
> These answers will drive the final `requirements.md` and `assumptions.md`.

---

## Rankings

**[1]: Should the backend enforce a maximum number of tags per ranking?**
The product doc says tags are optional keywords but does not specify a limit.

- [x] **Option A (recommended)**: Yes, limit to 10 tags per ranking — prevents abuse and keeps search performant
- [ ] **Option B**: No limit — allow any number of tags
- [ ] **Other**: _specify_

---

**[2]: Can the ranking creator edit the ranking name, description, tags, and visibility after creation?**
The product doc says the creator can change attributes and visibility, but does not explicitly mention name/description/tags editing.

- [x] **Option A (recommended)**: Yes, the creator can edit all ranking metadata (name, description, tags, visibility) at any time
- [ ] **Option B**: Only visibility and attributes can be changed after creation
- [ ] **Other**: _specify_

---

**[3]: Can the ranking creator delete a ranking entirely?**
The product doc does not mention deletion of rankings.

- [x] **Option A (recommended)**: Yes, the creator can delete a ranking — all associated items and ratings are also deleted
- [ ] **Option B**: No, rankings cannot be deleted, only made private
- [ ] **Other**: _specify_

---

**[4]: What is the maximum length for ranking name and description?**
No character limits are specified in the product doc.

- [x] **Option A (recommended)**: Name max 100 characters, description max 500 characters
- [ ] **Option B**: Name max 200 characters, description max 1000 characters
- [ ] **Other**: _specify_

---

## Attributes

**[5]: Can the creator modify attributes after ratings have already been submitted?**
The product doc says only the creator can define and alter attributes, but does not address the impact on existing ratings.

- [x] **Option A (recommended)**: Attributes can be added/renamed/removed freely; existing ratings for removed attributes are soft-deleted (kept in DB but excluded from calculations)
- [ ] **Option B**: Attributes can only be added or renamed, never removed, once any rating exists
- [ ] **Option C**: Attributes are locked once the first rating is submitted
- [ ] **Other**: _specify_

---

**[6]: What is the maximum length for attribute name and description?**

- [x] **Option A (recommended)**: Name max 50 characters, description max 200 characters
- [ ] **Option B**: Name max 100 characters, description max 500 characters
- [ ] **Other**: _specify_

---

## Items

**[7]: Can items be edited or deleted after creation?**
The product doc does not mention editing or deleting items.

- [x] **Option A (recommended)**: The ranking creator can edit or delete any item; the item creator can edit or delete their own items
- [ ] **Option B**: Only the ranking creator can edit or delete items
- [ ] **Option C**: Items cannot be deleted, only edited
- [ ] **Other**: _specify_

---

**[8]: What is the maximum number of items per ranking?**
No limit is specified in the product doc.

- [x] **Option A (recommended)**: 100 items per ranking — keeps the table view usable and queries efficient
- [ ] **Option B**: 50 items per ranking
- [ ] **Option C**: No limit
- [ ] **Other**: _specify_

---

**[9]: What is the maximum length for item name?**

- [x] **Option A (recommended)**: Max 100 characters
- [ ] **Option B**: Max 200 characters
- [ ] **Other**: _specify_

---

## Images (S3)

**[10]: What image formats should the backend accept?**

- [x] **Option A (recommended)**: JPEG, PNG, and WebP only
- [ ] **Option B**: JPEG, PNG, WebP, GIF, and AVIF
- [ ] **Other**: _specify_

---

**[11]: What should be the maximum file size for uploaded images?**

- [ ] **Option A (recommended)**: 5 MB
- [x] **Option B**: 10 MB
- [ ] **Option C**: 2 MB
- [ ] **Other**: _specify_

---

**[12]: What dimensions and format should the backend use for optimized images and thumbnails?**

- [x] **Option A (recommended)**: Optimized original resized to max 800×800px (WebP, quality 80); thumbnail 150×150px (WebP, quality 75)
- [ ] **Option B**: Optimized original max 1200×1200px (WebP); thumbnail 200×200px (WebP)
- [ ] **Other**: _specify_

---

**[13]: How long should pre-signed URLs be valid?**

- [x] **Option A (recommended)**: 1 hour — balances caching and security
- [ ] **Option B**: 15 minutes — more secure but requires more frequent regeneration
- [ ] **Option C**: 24 hours — better for caching, less secure
- [ ] **Other**: _specify_

---

## Ratings

**[14]: Can a user update their ratings after the initial submission?**
The product doc says ratings are saved when all attributes are filled, but does not mention re-rating.

- [x] **Option A (recommended)**: Yes, a user can update their ratings at any time by re-submitting all attribute scores for an item
- [ ] **Option B**: Ratings are final and cannot be changed once submitted
- [ ] **Other**: _specify_

---

**[15]: Should the backend validate that rating values are integers, or allow decimals?**
The product doc says "values from 0 to 100" but does not specify precision.

- [x] **Option A (recommended)**: Integer only (0–100) — simpler, aligns with the UI slider/input concept
- [ ] **Option B**: Decimal with one decimal place (0.0–100.0)
- [ ] **Other**: _specify_

---

## Search & Discovery

**[16]: Should search be case-insensitive and accent-insensitive (e.g., "video" matches "Vídeo")?**

- [x] **Option A (recommended)**: Yes, case-insensitive and accent-insensitive — better user experience for a Portuguese/multilingual audience
- [ ] **Option B**: Case-insensitive only, accents must match
- [ ] **Other**: _specify_

---

**[17]: How should the "10 latest public rankings" on the home page be implemented?**

- [x] **Option A (recommended)**: Dedicated endpoint `GET /api/v1/rankings/recent` returning the 10 most recently created public rankings
- [ ] **Option B**: Use the general listing endpoint `GET /api/v1/rankings?sort=-createdAt&limit=10&visibility=public`
- [ ] **Other**: _specify_

---

**[18]: Should tag search return exact matches or partial matches?**
Example: searching for tag "game" — should it match "games"?

- [x] **Option A (recommended)**: Exact tag match only — tags are discrete keywords, the user should use the correct tag
- [ ] **Option B**: Partial match (prefix-based) — "game" matches "games", "gameplay"
- [ ] **Other**: _specify_

---

## Pagination

**[19]: Which pagination strategy should the backend use?**

- [x] **Option A (recommended)**: Cursor-based pagination (using DynamoDB `LastEvaluatedKey`) — natural fit for DynamoDB, consistent with large/dynamic datasets
- [ ] **Option B**: Page-based pagination (page + size) — simpler for the frontend but requires workarounds with DynamoDB
- [ ] **Other**: _specify_

---

**[20]: What should be the default and maximum page size?**

- [x] **Option A (recommended)**: Default 20, maximum 50
- [ ] **Option B**: Default 10, maximum 100
- [ ] **Other**: _specify_

---

## API Design

**[21]: Should the API version prefix be included in all routes?**

- [x] **Option A (recommended)**: Yes, all routes under `/api/v1/` — allows future versioning
- [ ] **Option B**: No version prefix, just `/api/` — simpler, version later if needed
- [ ] **Other**: _specify_

---

**[22]: How should the backend handle the `X-User-Id` header when it is missing or invalid?**

- [x] **Option A (recommended)**: Return `400 Bad Request` for endpoints that require user context (create, rate, my-rankings); allow anonymous access for read-only public endpoints (list, get ranking, search)
- [ ] **Option B**: Always require `X-User-Id` on every request, return `400` if missing
- [ ] **Option C**: Generate a temporary ID server-side if missing
- [ ] **Other**: _specify_

---

**[23]: Should the backend return the `overall` score in API responses, or leave it entirely to the frontend?**
The product doc says overall is "always recalculated at view time" and "not stored permanently."

- [x] **Option A (recommended)**: The backend calculates and returns `overall` in the response for both "all users average" and "my scores" modes — ensures consistency and avoids duplicating logic in the frontend
- [ ] **Option B**: The backend only returns per-attribute scores; the frontend calculates overall
- [ ] **Other**: _specify_

---

## Error Handling

**[24]: What error response format should the API use?**

- [x] **Option A (recommended)**: Structured JSON error body: `{ "error": { "code": "RANKING_NOT_FOUND", "message": "Ranking with the given ID was not found", "status": 404 } }`
- [ ] **Option B**: Simple JSON: `{ "error": "Ranking not found" }`
- [ ] **Other**: _specify_

---

## Sorting

**[25]: Which fields should be sortable on the rankings listing endpoint?**

- [x] **Option A (recommended)**: `name`, `createdAt`, `updatedAt` — covers alphabetical and chronological use cases
- [ ] **Option B**: `name`, `createdAt` only
- [ ] **Other**: _specify_

---

**[26]: Which fields should be sortable on the items listing within a ranking (the ranking table view)?**

- [x] **Option A (recommended)**: `name`, `overall`, and each individual attribute score — matches the product requirement of sorting by any attribute or overall
- [ ] **Option B**: `name` and `overall` only
- [ ] **Other**: _specify_

---

## Ownership & Permissions

**[27]: Should the backend track who added each item to a ranking (item creator), or only track the ranking creator?**

- [x] **Option A (recommended)**: Track the item creator (`createdBy` field) — enables permission rules per item and audit trail
- [ ] **Option B**: Do not track item creator — only the ranking creator matters for permissions
- [ ] **Other**: _specify_

---

**[28]: Can any user add items to any public ranking, or only to rankings they created?**
The product doc says "both the creator and any user with access can add items."

- [x] **Option A (recommended)**: Any user can add items to any public ranking they can access; for private rankings, only users who have the URL (and thus the ranking ID) can add items
- [ ] **Option B**: Only the ranking creator can add items
- [ ] **Other**: _specify_

---

## Health & Operations

**[29]: Should the backend expose a health check endpoint?**

- [x] **Option A (recommended)**: Yes, `GET /api/v1/health` returning `200 OK` with basic status info (version, environment, DynamoDB connectivity)
- [ ] **Option B**: Yes, but minimal — just `200 OK` with no body
- [ ] **Other**: _specify_

---
