---
inclusion: auto
name: next-tdd
description: Guide how to apply TDD (Test Drive Development) in Next.js. Apply this kill when writing code in the frontend.
---

# Next.js TDD Best Practices

## Core TDD Principles

Claude must ensure that all test-related recommendations align with strict **TDD methodology**:

1. **Red → Green → Refactor Loop**

   - **Red:** Write a failing test before writing any implementation code.
   - **Green:** Write the simplest code necessary to make the test pass.
   - **Refactor:** Improve the code while keeping tests green and ensuring no regression.

2. **One change at a time**

   - Only write enough code to satisfy the current failing test.
   - Avoid implementing features or edge cases not covered by failing tests.

3. **Executable specifications**

   - Tests should act as living documentation for both frontend and backend behavior.
   - Each test should be readable by developers and non-developers alike.

---

## Testing Tools for Next.js Projects

When generating or reviewing code, Claude must prefer the following testing stack:

| Layer                 | Tool                                                                                                          | Purpose                                                        |
| --------------------- | ------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------- |
| **Unit Tests**        | [Jest](https://jestjs.io/) + [Testing Library](https://testing-library.com/docs/react-testing-library/intro/) | Test isolated components and utility functions                 |
| **Integration Tests** | Jest + Testing Library                                                                                        | Test combined component behavior and API interactions (mocked) |
| **E2E Tests**         | [Playwright](https://playwright.dev/)                                                                         | Validate full app flow in a browser environment                |
| **API Tests**         | Jest (with `supertest`)                                                                                       | Test server-side routes in isolation                           |

---

## Folder Structure for Tests

Claude should always suggest a clear folder structure:

```
project-root/
  app/
    page.js
    layout.js
  src/
    components/
      Button.jsx
      Button.test.jsx
    lib/
      utils.js
      utils.test.js
    services/
      api.js
      api.test.js
  tests/
    e2e/
      home.spec.ts
  __mocks__/
    fileMock.js
    apiMock.js
  jest.setup.js
```

**Guidelines:**

- Co-locate **unit and integration tests** next to their components or modules.
- Place **E2E tests** under `tests/e2e/`.
- Use `__mocks__` folder for mocks, spies, and fixtures.

---

## Writing Tests with TDD

### 1. Component TDD Flow

Example: Building a `Button` component.

- **Step 1: Write failing test first**

```jsx
// src/components/Button.test.jsx
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Button from "./Button";

describe("Button", () => {
  it("calls onClick when clicked", async () => {
    const user = userEvent.setup();
    const handleClick = jest.fn();

    render(<Button onClick={handleClick}>Click me</Button>);
    await user.click(screen.getByRole("button", { name: /click me/i }));

    expect(handleClick).toHaveBeenCalledTimes(1);
  });
});
```

- **Step 2: Run test → it fails (Red)**
- **Step 3: Implement minimal code**

```jsx
// src/components/Button.jsx
export default function Button({ children, onClick }) {
  return <button onClick={onClick}>{children}</button>;
}
```

- **Step 4: Run test → it passes (Green)**
- **Step 5: Refactor safely if needed (Refactor)**

---

## Integration Testing API Calls

> **Rule:** The frontend must **never call external APIs directly**.
> Always go through Next.js **API routes** acting as a backend proxy.

- Example test for a service calling an internal API route:

```jsx
// src/services/api.test.js
import { fetchData } from "./api";

global.fetch = jest.fn(() =>
  Promise.resolve({
    ok: true,
    json: () => Promise.resolve({ message: "Hello" }),
  })
);

test("fetchData returns expected data", async () => {
  const data = await fetchData();
  expect(data.message).toBe("Hello");
  expect(global.fetch).toHaveBeenCalledWith("/api/data");
});
```

- The actual implementation:

```js
// src/services/api.js
export async function fetchData() {
  const res = await fetch("/api/data");
  if (!res.ok) throw new Error("Network error");
  return res.json();
}
```

---

## End-to-End (E2E) Testing with Playwright

- Store E2E tests in `tests/e2e/`.
- Configure Playwright to run headless by default, with the option to run in headed mode for debugging.
- Example Playwright test:

```ts
import { test, expect } from "@playwright/test";

test("homepage loads and displays main heading", async ({ page }) => {
  await page.goto("http://localhost:3000");
  await expect(page.getByRole("heading", { name: /welcome/i })).toBeVisible();
});
```

---

## Test Doubles and Mocks

- **Never mock what you don't own.**
  Only mock external services, APIs, and unstable dependencies.
- Prefer **realistic mock data** that reflects production scenarios.
- Place reusable mocks and fixtures in `__mocks__/`.

---

## CI/CD Integration

- Run **unit and integration tests** on every commit or pull request.
- Run **E2E tests** on a staging environment before deploying to production.
- Example GitHub Actions workflow step:

```yaml
steps:
  - run: npm run test -- --coverage
  - run: npm run test:e2e
```

---

## Code Coverage

- Always generate coverage reports:

```bash
npm run test -- --coverage
```

- Minimum threshold: **80% for statements, branches, functions, and lines**.

---

## Summary of Rules for Claude Code

1. Always follow **Red → Green → Refactor** strictly.
2. Use **Jest** for unit and integration tests.
3. Use **Testing Library** for React component behavior.
4. Use **Playwright** for E2E browser-level tests.
5. Co-locate tests with their respective modules/components.
6. Mock only external APIs or unstable dependencies.
7. Never allow the frontend to call external APIs directly.
8. Ensure all CI pipelines run tests automatically.
9. Maintain high coverage and treat failing tests as blockers.
10. Write tests as **living documentation**, readable and clear.