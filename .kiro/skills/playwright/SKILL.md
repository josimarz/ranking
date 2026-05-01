---
name: playwright
description: Battle-tested Playwright patterns for E2E, API, component, visual, accessibility, and security testing. Covers locators, assertions, fixtures, network mocking, auth flows, debugging, and framework recipes for React, Next.js, Vue, and Angular. TypeScript and JavaScript.
---

# Playwright Core Testing

> Opinionated, production-tested Playwright guidance — every pattern includes when (and when *not*) to use it.

**22 reference guides** covering the full Playwright testing surface: selectors, assertions, fixtures, network mocking, auth, visual regression, accessibility, API testing, debugging, and more — with TypeScript and JavaScript examples throughout.

## Golden Rules

1. **`getByRole()` over CSS/XPath** — resilient to markup changes, mirrors how users see the page
2. **Never `page.waitForTimeout()`** — use `expect(locator).toBeVisible()` or `page.waitForURL()`
3. **Web-first assertions** — `expect(locator)` auto-retries; `expect(await locator.textContent())` does not
4. **Isolate every test** — no shared state, no execution-order dependencies
5. **`baseURL` in config** — zero hardcoded URLs in tests
6. **Retries: `2` in CI, `0` locally** — surface flakiness where it matters
7. **Traces: `'on-first-retry'`** — rich debugging artifacts without CI slowdown
8. **Fixtures over globals** — share state via `test.extend()`, not module-level variables
9. **One behavior per test** — multiple related `expect()` calls are fine
10. **Mock external services only** — never mock your own app; mock third-party APIs, payment gateways, email

## Guide Index

### Writing Tests

| What you're doing | Guide |
|---|---|
| Choosing selectors | [locators.md](references/locators.md) |
| Assertions & waiting | [assertions-and-waiting.md](references/assertions-and-waiting.md) |
| Organizing test suites | [test-organization.md](references/test-organization.md) |
| Playwright config | [configuration.md](references/configuration.md) |
| Fixtures & hooks | [fixtures-and-hooks.md](references/fixtures-and-hooks.md) |
| Test data | [test-data-management.md](references/test-data-management.md) |
| Auth & login | [authentication.md](references/authentication.md) |
| API testing (REST/GraphQL) | [api-testing.md](references/api-testing.md) |
| Visual regression | [visual-regression.md](references/visual-regression.md) |
| Accessibility | [accessibility.md](references/accessibility.md) |
| Network mocking | [network-mocking.md](references/network-mocking.md) |
| Forms & validation | [forms-and-validation.md](references/forms-and-validation.md) |
| CRUD flows | [crud-testing.md](references/crud-testing.md) |

### Debugging & Fixing

| Problem | Guide |
|---|---|
| General debugging workflow | [debugging.md](references/debugging.md) |
| Specific error message | [error-index.md](references/error-index.md) |
| Flaky / intermittent tests | [flaky-tests.md](references/flaky-tests.md) |
| Common beginner mistakes | [common-pitfalls.md](references/common-pitfalls.md) |

### Framework Recipes

| Framework | Guide |
|---|---|
| Next.js (App Router + Pages Router) | [nextjs.md](references/nextjs.md) |
| React (CRA, Vite) | [react.md](references/react.md) |
| Vue 3 / Nuxt | [vue.md](references/vue.md) |

### Additional References

| Topic | Guide |
|---|---|
| Screenshots | [screenshots.md](references/screenshots.md) |
| Page Object Model | [page-object-model.md](references/page-object-model.md) |
