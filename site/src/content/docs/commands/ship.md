---
title: /ship
description: Create a PR with spec-derived context.
---

Auto-populates PR descriptions from the spec so reviewers understand intent without reading the full spec.

## When to use

- Feature is built and verified
- After `/build` completes

## What it does

Creates a PR with:
- **Why** — intent from the spec
- **What** — behavioral summary
- **Key decisions** — from DECISIONS.md
- **Constraints** — what must NOT happen
- **Verification checklist** — from the plan

Links to the spec file for deep dives.

## Output

A GitHub PR with full context.
