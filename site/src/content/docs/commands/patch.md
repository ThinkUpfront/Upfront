---
title: /upfront:patch
description: Fix a bug or small feature from a GitHub issue — more structure than /quick, less than /feature.
---

For well-understood problems that need some rigor but not the full ceremony.

## When to use

- GitHub issue with a clear problem statement
- Bug fix between 50-300 lines
- Small feature request with known scope

## What it does

1. **Understand** — reads the GitHub issue (`gh issue view`) or description, plus context files
2. **Investigate** — reads relevant code, presents diagnosis with root cause, approach, files, and risk
3. **Confirm** — waits for user approval of the approach (skipped in `--auto` mode)
4. **Implement with TDD** — failing test first, then implementation
5. **Self-review** — checks against DECISIONS.md, ARCHITECTURE.md, CONSTITUTION.md
6. **Run all checks** — full test suite, not just the new tests
7. **Commit** — with `Fixes #[issue]` reference
8. **Report** — offers to close the GitHub issue with a summary

## Auto mode

`/upfront:patch --auto #42` — runs without pausing for confirmation. Still stops for scope overflow, unclear root cause, or constitutional violations.

## Scope gates

- Under ~50 lines? Redirects to `/upfront:quick`
- Over ~300 lines? Redirects to `/upfront:feature` + `/upfront:plan` + `/upfront:build`
- Root cause unclear? Suggests `/upfront:debug` instead

## Output

A committed fix with tests, referencing the issue.
