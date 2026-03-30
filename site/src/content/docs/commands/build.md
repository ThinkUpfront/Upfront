---
title: /upfront:build
description: Execute phases with TDD, review, red team, and worktree isolation.
---

The orchestrator. Spawns a fresh sub-agent for each phase, enforces strict TDD, runs post-phase review, and manages the entire build lifecycle.

## When to use

- You have an approved plan from `/upfront:plan`
- Run as: `/upfront:build specs/[feature-name]-plan.md`

## What it does

### Pre-flight

Detects all ecosystems, verifies existing tools work, audits for missing tools with per-language checklists. Pushes hard for installation.

### Worktree isolation

Creates an isolated git worktree for the build. Main working tree stays clean. Failed builds = delete the worktree. No stash dance.

### Seed test stubs

Before Phase 1, generates skeleton test files from the spec's requirements, error cases, concurrency concerns, and boundary values. Each test is a `t.Skip("TODO")` — a map of what to test before any code is written.

### Per phase (1-9 steps)

1. **Announce** with visual progress bar
2. **Spawn sub-agent** with clean context (RALPH pattern)
3. **Review work** — TDD followed? Scope respected? Error pattern forwarding for re-spawns.
4. **Automated verification** — independent of sub-agent's self-report
5. **Visual verification** — screenshots via KaBOOM MCP for UI phases
6. **Code review** — spec compliance, correctness, architecture
7. **Update progress file**
8. **Commit** with sequential IDs
9. **Auto-proceed** if no manual verification needed

### After all phases

1. **Integration sweep** — verify pieces connect correctly
2. **Red team** — adversarial agent breaks correctness, concurrency, boundaries, tests, security, and UI (with browser automation for visual attacks)
3. **Learning capture** — appends to `specs/LEARNINGS.md`
4. **Compound** — promotes durable patterns from LEARNINGS.md to CLAUDE.md as agent instructions
5. **Merge back** — merge worktree branch, clean up

### Human-writes mode (`/upfront:build --human`)

AI writes the tests. You write the implementation. Full reporting: attempts to green, tests that caught bugs, time from stubs to green. Commits tagged `[human-writes]`.

### Crash recovery

Detects existing worktrees, uncommitted changes (keep/stash/discard), reconciles progress with git history, resumes from the right phase.

## Key features

- **Stuck loop detection** — 3-retry hard limit on any repeated operation
- **Decimal phase insertion** — insert Phase 2.5 for unexpected prerequisite work
- **Constitutional enforcement** — reads `specs/CONSTITUTION.md` as a hard constraint
- **Dev server** — starts automatically for UI features, used for visual verification

## Output

Committed phases on a feature branch, merged to base on success. Progress file, learnings, compounded patterns.

## Next step

`/upfront:ship` to create a PR.
