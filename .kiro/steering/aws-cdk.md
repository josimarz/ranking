---
inclusion: auto
name: aws-cdk
description: Guide and best practices when writing code using AWS CDK library in TypeScript. Apply this skill when writing IaC.
---

# AWS CDK TypeScript Best Practices

This document provides guidance and best practices for writing AWS CDK code in TypeScript.

## Start a New AWS CDK Project

To start a new AWS CDK project, use the following commands:

```bash
pnpm dlx cdk init app --language typescript --generate-only
pnpm install
```

These commands initialize a TypeScript CDK project without immediately synthesizing it, and install all required dependencies.

## Best Practices

When developing with AWS CDK, follow these guidelines:

1. **Follow recommended Model Context Protocols (MCPs)**
   Reference the following MCPs to guide best practices when writing CDK code:

   * `awslabs.core-mcp-server`
   * `awslabs.cdk-mcp-server`

2. **Follow established CDK patterns**

   * Use `Stacks` to group related resources.
   * Use `Constructs` to create reusable components.
   * Avoid hardcoding configuration values; use context variables or environment variables.

3. **Keep code modular and maintainable**

   * Split large stacks into multiple files or constructs.
   * Write helper functions for repeated logic.
   * Organize resources logically to make code easy to read and maintain.

4. **Infrastructure as Code hygiene**

   * Use strong typing with TypeScript interfaces where possible.
   * Include comments explaining the purpose of resources.
   * Apply consistent naming conventions for resources.

5. **Security and compliance**

   * Follow AWS security best practices (least privilege, encryption, etc.).
   * Use CDK’s built-in methods for managing IAM policies safely.

6. **Testing and validation**

   * Write unit tests for constructs using `@aws-cdk/assert` or similar libraries.
   * Use `cdk synth` and `cdk diff` to validate changes before deploying.

7. **Deployment practices**

   * Use separate environments (dev, staging, prod) with isolated stacks.
   * Automate deployment using CI/CD pipelines.
   * Avoid manual edits to deployed resources outside of CDK.

Following these best practices and MCPs will help ensure that your AWS CDK projects are maintainable, secure, and scalable.