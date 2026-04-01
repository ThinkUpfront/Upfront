---
title: Skills Overview
description: All Upfront skills and when to use each one.
---

Upfront is a set of slash skills for [Claude Code](https://claude.ai/claude-code). Each skill is a markdown file — no dependencies, no build step.

## The flow

```
Vision → Feature → Plan → Build
```

## Choosing the right skill

| Situation | Skill | Ceremony |
|-----------|-------|----------|
| Big ambition, need strategic clarity | `/upfront:vision` | Rumelt's kernel: diagnosis, policies, actions |
| Between increments, need to reflect | `/upfront:increment` | Structured retro + next increment steering |
| Specific concern, risk, or tradeoff | `/upfront:assess` | Interactive problem-solving |
| Don't know what to build | `/upfront:ideate` | Conversation only |
| Need to understand a codebase | `/upfront:explore` | Produces ARCHITECTURE.md |
| Improve your AI instruction files | `/upfront:enlighten` | Audits CLAUDE.md/AGENTS.md, adds stack-specific examples |
| Need to understand a codebase you forgot | `/upfront:teach` | Walkthrough + optional quiz |
| Know the problem, need to define it | `/upfront:feature` | 4-phase forced thinking |
| Need to revise a spec | `/upfront:refine` | Challenge-based iteration |
| Spec ready, need implementation phases | `/upfront:plan` | Architecture deep-dive + phasing |
| Plan ready, need to build | `/upfront:build` | TDD + review + red team per phase |
| GitHub issue or clear bug | `/upfront:patch` | Investigate + TDD + commit |
| Tiny change (<50 lines) | `/upfront:quick` | Just do it with TDD |
| Something is broken, unclear why | `/upfront:debug` | Scientific method |
| Structural debt, need evolution plan | `/upfront:architect` | Architecture review + evolution plan |
| Evolution plan ready, need to restructure | `/upfront:re-architect` | Phase-by-phase restructuring |
| Feature built, need a PR | `/upfront:ship` | Auto-populated from spec |
| Feature shipped, need to check | `/upfront:retro` | Predictions vs reality |
| Check tooling is current | `/upfront:upgrade` | Health check + fix |
| Not sure where to start | `/upfront:up` | Smart router |
| Capture a quick idea or todo | `/upfront:note` | Append to TODO.md |
| Need to stop mid-work | `/upfront:pause` | Structured handoff |
| Resuming from a pause | `/upfront:resume` | Context restoration |

## Size guide

| Change size | Skill path |
|-------------|------------|
| <50 lines | `/upfront:quick` |
| 50-300 lines | `/upfront:patch` |
| 300+ lines | `/upfront:feature` → `/upfront:plan` → `/upfront:build` |

`/upfront:patch` and `/upfront:quick` both enforce scope gates — if a change grows beyond its boundary, the skill stops and redirects you up.
