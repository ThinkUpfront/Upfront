---
title: /upfront:debug
description: Scientific method debugging with persistent state.
---

Hypothesis-driven debugging that survives session crashes.

## When to use

- Something is broken and you don't know why
- Previous fix attempts haven't worked
- Need structured investigation, not guessing

## What it does

1. **Gather evidence** — read the error, check logs, check browser devtools (KaBOOM MCP for web issues)
2. **Form hypotheses** — ranked by likelihood
3. **Test** — design the smallest experiment that eliminates a hypothesis
4. **Narrow** — update hypothesis rankings based on results
5. **Fix** — once root cause is confirmed, fix with TDD

State persists in `specs/DEBUG.md`. If the session dies mid-debug, the next session reads the file and continues — never re-tries eliminated hypotheses.

**Circuit breaker:** After 3 failed hypothesis cycles, stops and asks for more context.

## Output

A committed fix, and `specs/DEBUG.md` with the full investigation trail. Appends to `specs/LEARNINGS.md`.
