---
inclusion: auto
name: macos
description: Instruction how to run some commands in the macOS terminal. Apply this skill when running commands in the terminal.
---

## macOS terminal commands

When running shell commands on macOS follow theses instructions:
- To kill a process running in a specific port, use `kill $(lsof -ti:PORT)`. Example: `kill $(lsof -ti:3000)`.
- To run a process in background, use `nohup <cmd> > <out>.log 2>&1 &`. Example: `nohup pnpm dev > out.log 2>&1 &`.
- When running `git diff` avoid pagination by adding the argument `--no-pager`. Example: `git --no-pager diff --stat`.