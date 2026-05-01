---
inclusion: auto
name: node
description: Instructions to use PNPM as package manager. Apply this skill when running commands in the terminal to interact with Node.js.
---

# Best practices when working with Node.js environment

- Always load env files using the argument `--env-file`.

## Use pnpm over npm

Default to using pnpm instead of npm

- Use `pnpm <script>` to run a script defined in `package.json`
- Use `pnpm install` to install all packages
- Use `pnpm add <dep>` to add a dependency
- Use `pnpm -D add <dep>` to add a dev dependency
- Use `pnpm dlx <dep>` to run a CLI tool for a dependency that is not installed
- Use `pnpm exec <dep>` to run a CLI tool for a dependency that is installed locally
- Never install a dependency globally, always install locally and use `pnpm exec`
- Once `pnpm deploy` is a special command for PNPM Workspaces, prefere `pnpm run deploy`