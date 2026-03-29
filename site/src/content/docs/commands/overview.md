---
title: Commands Overview
description: All Upfront commands and when to use each one.
---

Upfront is a set of slash commands for [Claude Code](https://claude.ai/claude-code). Each command is a markdown file — no dependencies, no build step.

## The flow

```
Think → Define → Plan → Build → Ship → Learn
```

## Choosing the right command

| Situation | Command | Ceremony |
|-----------|---------|----------|
| Don't know what to build | `/ideate` | Conversation only |
| Need to understand a codebase | `/explore` | Produces ARCHITECTURE.md |
| Need to understand a codebase you forgot | `/teach` | Walkthrough + optional quiz |
| Know the problem, need to define it | `/feature` | 4-phase forced thinking |
| Need to revise a spec | `/refine` | Challenge-based iteration |
| Spec ready, need implementation phases | `/plan` | Architecture deep-dive + phasing |
| Plan ready, need to build | `/build` | TDD + review + red team per phase |
| GitHub issue or clear bug | `/patch` | Investigate + TDD + commit |
| Tiny change (<50 lines) | `/quick` | Just do it with TDD |
| Something is broken, unclear why | `/debug` | Scientific method |
| Feature built, need a PR | `/ship` | Auto-populated from spec |
| Feature shipped, need to check | `/retro` | Predictions vs reality |

## Size guide

| Change size | Command path |
|-------------|-------------|
| <50 lines | `/quick` |
| 50-300 lines | `/patch` |
| 300+ lines | `/feature` → `/plan` → `/build` |

`/patch` and `/quick` both enforce scope gates — if a change grows beyond its boundary, the command stops and redirects you up.
