---
title: /upfront:architect
description: Architecture review that produces an evolution plan.
---

## What it does

`/upfront:architect` performs a structural code review focused on finding architectural problems — not bugs, not features. It maps the codebase, identifies god modules, layer violations, and pattern inconsistencies, then produces an evolution plan.

## When to use it

- Before a major feature when you suspect the codebase has structural debt
- After a period of rapid AI-assisted development (when structure tends to drift)
- When onboarding to an unfamiliar codebase and wanting to understand its health
- Periodically as a health check

## How it works

1. **Load context** — reads `specs/ARCHITECTURE.md`, `specs/DECISIONS.md`, `specs/LEARNINGS.md`
2. **Structural survey** — parallel sub-agents map module inventory, dependency graph, and layer analysis
3. **God module detection** — flags modules with multiple responsibilities or excessive size
4. **Layer violation scan** — identifies improper cross-layer dependencies
5. **Pattern audit** — checks for inconsistent patterns and unnecessary abstractions
6. **Evolution plan** — writes `specs/ARCHITECTURE-EVOLUTION.md` with phased restructuring steps
7. **Decision record** — appends the review findings to `specs/DECISIONS.md`

## Key principles

- **Restructuring only** — no behavior changes, no new features
- **God modules first** — fixing patterns inside a god module is wasted effort if it's getting split
- **Don't add complexity** — patterns are for managing complexity that exists, not complexity you're introducing
- **Clean is clean** — if the architecture is healthy, it says so and moves on

## Output

- `specs/ARCHITECTURE-EVOLUTION.md` — phased plan for structural improvements
- Updated `specs/DECISIONS.md` — rationale for proposed changes
- Summary of findings

## Next step

Run `/upfront:re-architect` to execute the evolution plan phase by phase.
