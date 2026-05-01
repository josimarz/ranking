---
inclusion: auto
name: next-testing
description: Best practices for testing code in Next.js. Apply this skill when writing tests in the frontend.
---

# Frontend Feature Testing Best Practices with MCP Playwright

This guide outlines best practices for testing new frontend functionalities using **MCP Playwright**, where an AI automatically executes tests for each new feature in real-time, without the need for manual browser interaction or writing test code.

---

## 1. Purpose of MCP Playwright

MCP Playwright provides:

- **Automated real-time frontend testing** executed by AI.
- **Simulation of user interactions** and functional validation.
- **No manual test code or browser interaction required**.
- **Integration with modern frameworks like Next**.

---

## 2. Workflow for Testing New Features

Whenever a new frontend feature is completed:

1. The AI detects the feature as completed.
2. MCP Playwright automatically:

   - Runs real-time functional checks against the development or staging environment.
   - Validates user interactions, UI behaviors, and expected outcomes.
   - Reports any issues, errors, or regressions.

3. Results are logged, allowing developers to review issues **without manually testing in a browser**.

This provides immediate feedback on the feature's functionality in a live environment.

---

## 3. Best Practices

- **Feature Readiness:** Ensure that the feature is fully implemented and accessible before triggering MCP Playwright tests.
- **Stable Environment:** Test in a consistent development or staging environment to avoid false negatives.
- **Use Consistent Selectors:** Employ stable attributes such as `data-testid` to allow reliable automated checks.
- **Monitor AI Test Reports:** Regularly review AI-generated test results to catch functional issues early.
- **Iterative Validation:** Run MCP Playwright tests for every new feature or significant update.
- **CI/CD Integration:** Incorporate MCP Playwright validation into your CI/CD pipeline for real-time functional verification.

---

## 4. Summary

- MCP Playwright provides **automated, real-time testing** of frontend functionalities.
- **No manual test scripts or browser interaction** are required.
- This testing **supplements, but does not replace**, unit, integration, and E2E tests.
- Provides rapid feedback to developers, ensuring that new features behave as expected before manual review.

---

**Author:** Internal Front-End Standards Team
**Last Updated:** September 2025