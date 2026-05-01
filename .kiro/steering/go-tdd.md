---
inclusion: auto
name: go-tdd
description: Guide for Test Driven Development in Go/Golang. Apply this skill when writing code in the backend.
---

# Best Practices for Writing Code with Test-Driven Development (TDD)

Test-Driven Development (TDD) is a software development practice where tests are written before the implementation code. It follows a simple cycle often referred to as **Red → Green → Refactor**:

1. **Red** – Write a failing test that defines a desired behavior.  
2. **Green** – Write the minimal code required to make the test pass.  
3. **Refactor** – Clean up the code while ensuring tests remain green.  

This approach leads to more reliable, maintainable, and well-structured code. Below are the best practices for applying TDD effectively.

---

## 1. Write Small, Incremental Tests
- Keep each test focused on **a single behavior**.  
- Avoid large, complex test cases that try to cover too much at once.  
- Smaller tests make it easier to identify the cause of a failure.  

✅ Example:  
- Write one test for a method returning the correct sum.  
- Write another test for how it handles negative numbers.  

---

## 2. Start with the Simplest Failing Test
- Always begin with the **easiest test case** to implement.  
- Don’t over-engineer the solution for future requirements.  
- Let tests guide the evolution of the code.  

---

## 3. Follow the Red → Green → Refactor Cycle
- **Red:** Ensure the test fails first (proves the test is valid).  
- **Green:** Implement just enough code to make the test pass.  
- **Refactor:** Eliminate duplication, improve readability, and optimize without altering behavior.  

---

## 4. Write Tests That Express Intent
- Use meaningful test names that describe **behavior**, not implementation.  
- Prefer descriptive test methods like `test_user_cannot_login_with_wrong_password()` over vague names like `test_login()`.  

---

## 5. Keep Tests Independent and Deterministic
- Tests should not depend on each other or external states.  
- Each test should run in isolation and always produce the same result.  
- Avoid hidden dependencies such as databases, file systems, or network calls. Use mocks or stubs when necessary.  

---

## 6. Maintain a Fast Feedback Loop
- Tests should run quickly to encourage frequent execution.  
- A slow test suite discourages developers from running it often.  
- Keep integration and end-to-end tests separate from fast unit tests.  

---

## 7. Embrace Refactoring as Part of TDD
- Refactoring is not optional—it’s a **core step** in TDD.  
- Use refactoring to:  
  - Remove duplication.  
  - Improve code readability.  
  - Simplify design while keeping tests green.  

---

## 8. Write Tests Before Fixing Bugs
- When a bug is reported, first write a test that reproduces the bug.  
- Run the test to see it fail (Red).  
- Fix the bug and watch the test pass (Green).  
- Refactor if needed.  

This ensures the bug never reappears unnoticed.  

---

## 9. Keep the Test Suite Clean
- Remove redundant or outdated tests.  
- Avoid overly complex test setups.  
- Refactor tests just like production code—clean tests are easier to maintain.  

---

## 10. Balance Unit, Integration, and End-to-End Tests
- Focus mainly on **unit tests** for fast feedback.  
- Add **integration tests** for critical interactions between components.  
- Use **end-to-end tests** sparingly to validate the system as a whole.  

---

## Conclusion
TDD is more than writing tests first—it’s a discipline that shapes how you design, implement, and maintain software. By following these best practices, you can:  

- Build confidence in your code.  
- Encourage clean, modular design.  
- Catch bugs early and prevent regressions.  
- Develop software that is easier to maintain and extend.  

Adopting TDD requires practice and consistency, but the long-term benefits for software quality and developer productivity are significant.