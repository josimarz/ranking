---
inclusion: auto
name: auth
description: Define how should users be authenticated. Apply this skill when implementing authentication in the frontend.
---

# Technical Specification – User Identification (No Traditional Authentication)

This document defines how **Ranking** identifies users without using traditional authentication (login/password).

The system relies on a **browser-scoped unique identifier** stored in `localStorage`.  
This identifier represents the “user” across all interactions in the product.

The approach prioritizes:

- Zero friction
- Immediate usability
- No registration flow
- No personal data collection

---

## 1. Core Concept

Each browser instance receives a unique identifier:

- Format: **UUID v7**
- Storage: `localStorage`
- Key name: `ranking:userId`

Example:

```json
localStorage["ranking:userId"] = "019bc9c7-b857-72c5-a491-7f9686c7989a"
```

This value becomes the **User ID** for all operations:

* Ownership of rankings
* Attribution of item ratings
* “My Rankings” listing
* User-specific views (e.g., *My Scores*)

No email, password, or account is ever created.

---

## 2. Initialization Flow

On every page load, the frontend must execute:

```ts
function getOrCreateUserId(): string {
  const key = "ranking:userId";
  let userId = localStorage.getItem(key);

  if (!userId) {
    userId = generateUUIDv7();
    localStorage.setItem(key, userId);
  }

  return userId;
}
```

Flow:

1. Check `localStorage` for `ranking:userId`
2. If it exists:

   * Use it as the current user identifier
3. If it does not exist:

   * Generate a new UUIDv7
   * Persist it in `localStorage`
   * Use it for the session

This logic must run:

* On the very first visit
* Before any API interaction
* In every environment (home, ranking page, search, etc.)

---

## 3. API Contract

Every request that depends on user context must include the User ID.

Options:

* HTTP header (preferred):

  ```
  X-User-Id: 019bc9c7-b857-72c5-a491-7f9686c7989a
  ```
* Or request body field:

  ```json
  {
    "userId": "019bc9c7-b857-72c5-a491-7f9686c7989a"
  }
  ```

The backend must:

* Treat this ID as opaque
* Never infer personal data
* Never modify it
* Trust it as the identity boundary

---

## 4. Ownership Model

### Rankings

When a ranking is created:

```json
{
  "id": "ranking-uuid",
  "ownerUserId": "user-uuid",
  "name": "Video Games",
  "visibility": "public"
}
```

Rules:

* `ownerUserId` is always the current User ID
* Only the owner can:

  * Edit ranking metadata
  * Change attributes
  * Change visibility

### Ratings

Each rating is scoped by:

* `rankingId`
* `itemId`
* `userId`

Example:

```json
{
  "rankingId": "ranking-uuid",
  "itemId": "item-uuid",
  "userId": "user-uuid",
  "scores": {
    "graphics": 90,
    "sound": 85,
    "price": 70
  }
}
```

This allows:

* Aggregated averages (all users)
* Personal views (“My Scores”)

---

## 5. “My Rankings”

The “My Rankings” page is resolved entirely by `userId`.

Query:

```
GET /api/rankings?ownerUserId={userId}
```

The backend returns only rankings created with that ID.

If the user:

* Clears browser data
* Changes browser
* Switches device

They receive a **new User ID** and therefore:

* Lose access to previous “My Rankings”
* Appear as a new user

This is an **intentional business decision**.

---

## 6. Edge Cases

### Clearing Storage

If `localStorage` is cleared:

* A new UUID is generated
* The user becomes “new”
* Previous data remains in the system but is no longer associated

### Incognito Mode

* Each incognito session gets a new User ID
* Rankings created there will not persist

### Multiple Tabs

* All tabs share the same `localStorage`
* The User ID is consistent across tabs

---

## 7. Privacy & Compliance

This system:

* Stores no personal information
* Does not identify a real person
* Does not require consent for personal data processing
* Acts as a **technical session identity**, not an account

The User ID:

* Must never be shown as “account data”
* Must not be editable in the UI
* Must not be used for marketing or profiling

It exists solely to:

* Attribute content
* Enable personalization
* Maintain continuity of use

---

## 8. Security Considerations

This model is **not** designed for strong identity guarantees.

Known limitations:

* User IDs can be copied between browsers
* IDs can be manually changed in DevTools
* There is no identity verification

Mitigations:

* Treat all actions as low-risk
* Avoid sensitive operations
* Never store private personal data
* Use visibility rules (public vs private) as the primary access control

This is acceptable because:

* The product has no financial operations
* No personal data is involved
* The business explicitly favors simplicity over strict identity

---

## 9. Summary

* Every browser gets a UUIDv7 stored in `localStorage`
* This UUID is the sole user identifier
* It is created automatically and silently
* It powers:

  * Ranking ownership
  * User ratings
  * “My Rankings”
* Losing the ID means losing continuity
* This is a conscious product decision

This approach delivers **zero-friction onboarding** while preserving enough identity to support all core features of the product.