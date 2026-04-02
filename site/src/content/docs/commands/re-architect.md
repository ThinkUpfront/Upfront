---
title: /upfront:re-architect
description: Execute an architecture evolution plan phase by phase.
---

## What it does

`/upfront:re-architect` executes the evolution plan produced by `/upfront:architect`. It restructures code phase by phase — moving files, updating imports, fixing layer violations — with test verification at phase boundaries.

## When to use it

- After `/upfront:architect` has produced an evolution plan
- When you have a clear structural change to make (module splits, layer fixes)

## How it works

1. **Load the plan** — reads `specs/ARCHITECTURE-EVOLUTION.md`
2. **Phase execution** — works through each phase sequentially
3. **Tests break mid-phase** — that's expected when moving files and updating imports
4. **Phase boundary verification** — all existing tests must pass at the end of each phase
5. **Commit per phase** — each completed phase gets its own commit

## Key difference from /upfront:build

`/upfront:build` runs tests after every edit and uses TDD. `/upfront:re-architect` expects tests to break during a phase as you restructure — it only verifies at phase boundaries. No new tests are written; existing tests must pass after restructuring.

## Rules

- **No behavior changes** — restructuring only. If tests fail after a phase, the restructuring introduced a bug.
- **No new features** — if you find something that needs a feature, note it and move on.
- **No new tests** — existing tests are the contract. If they pass, the restructuring is correct.
- **One phase at a time** — human reviews and approves each phase before the next begins.

## Output

- Restructured codebase with clean commits per phase
- Updated `specs/ARCHITECTURE.md` with new structure
- Updated `specs/LEARNINGS.md` with restructuring observations
