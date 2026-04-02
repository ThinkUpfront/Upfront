---
title: /upfront:quick
description: Small changes without ceremony — under 50 lines.
---

For well-understood changes that don't need the full workflow.

## When to use

- Config changes, renames, small fixes
- Under ~50 lines of non-test code
- No ambiguity about what to do

## What it does

1. Read context (ARCHITECTURE.md, DECISIONS.md, LEARNINGS.md, CONSTITUTION.md)
2. Scope check — bail if too big
3. TDD if testable
4. Make the change
5. Run all checks
6. Scope re-check — if it grew past 50 lines, stash and redirect to `/upfront:feature`
7. Self-review
8. Commit

## Scope gate

If the change grows beyond ~50 lines during implementation, `/upfront:quick` stops, stashes the partial work, and tells you to run `/upfront:feature` instead.

## Output

A committed change with conventional prefix (`fix:`, `chore:`, `refactor:`, `tweak:`, `docs:`).
