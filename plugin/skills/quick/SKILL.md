---
description: Execute a small, well-understood change without the full /upfront:feature + /upfront:plan + /upfront:build ceremony
user-invocable: true
---

# Quick

You are executing a small, well-scoped change. No spec, no plan, no phases. This is for fixes, tweaks, renames, config changes, and small improvements — NOT new features.

**Input:** $ARGUMENTS

## Pre-check

### 1. Read context (if available)

Silently read these files if they exist:
- `specs/ARCHITECTURE.md`
- `specs/DECISIONS.md`
- `specs/LEARNINGS.md`

These give you the project's conventions, patterns, prior decisions, and known gotchas. Follow them. Architectural invariants (from ARCHITECTURE.md's invariants section) are hard constraints — if the change would violate one, stop and say so.

### 2. Scope check

Evaluate the requested change. If ANY of these are true, stop immediately:

- The change introduces a new user-facing feature (not a fix/tweak to an existing one)
- The change will touch more than ~50 lines of non-test code (estimated)
- The change requires new data models, new API endpoints, or new modules
- The change has unclear requirements that need a design conversation

If you hit any of these, escalate:

- If it's a single feature: "This is bigger than a /upfront:quick change. Run `/upfront:feature` instead to define the intent and constraints properly."
- If it spans multiple concerns, subsystems, or would take multiple features to do right: "This isn't quick — it's increment 1 of something bigger. Consider `/upfront:vision` to capture what you're really after, or at minimum `/upfront:feature` to define one piece properly."

Do not proceed. Do not try to be helpful by doing it anyway.

### 3. Understand the change area

Read the files that will be affected. Understand the existing patterns, conventions, and test coverage before touching anything. Do not guess — look.

## Execute

### Step 1: TDD (if testable)

If the change affects logic (not config, not docs, not formatting):

1. **Write a failing test first** that captures the expected behavior change. Run it. Confirm it fails.
2. **Make the change.** Run the test. Confirm it passes.
3. **Refactor if needed.** Keep tests green.

If the change is not testable (config file, documentation, rename with no behavioral change, formatting), state why TDD does not apply and proceed.

### Step 2: Make the change

- Follow existing code conventions exactly
- Do not refactor nearby code (that is a separate change)
- Do not add features beyond what was asked
- Do not "improve" things you noticed while in the area

### Step 3: Run all checks

Run every existing linter, formatter, type checker, and test suite for the affected area. Not just the test you wrote — the full suite. Fix anything you broke.

If the project has no automated checks, say so explicitly: "This project has no linters or test suite configured. The change compiles/parses but has no automated verification."

### Step 4: Scope re-check

Now that you have made the change, check: did it stay under ~50 lines of non-test code?

If it grew beyond that, stop and say:

> "This turned out to be bigger than expected (~[N] lines). The change is partially done. You should run `/upfront:feature` to define this properly, then `/upfront:plan` and `/upfront:build` to finish it. I'll stash the work-in-progress so you can pick it up later: `git stash push -m 'quick: partial [description]'`"

Stash the partial work before stopping. Do not leave uncommitted half-done code in the working tree.

### Step 5: Quick review

Do a single review pass on your own diff. Check for:
- Spec compliance (does the change do what was asked, nothing more?)
- Correctness (edge cases, nil/null, off-by-one, error handling)
- Convention compliance (does it match the existing codebase patterns?)
- Anything that violates `specs/DECISIONS.md` constraints

If you find issues, fix them and re-run checks.

### Step 6: Commit

Commit with a clear message that describes what changed and why. Use a conventional prefix:

- `fix:` for bug fixes
- `chore:` for config, dependency, or tooling changes
- `refactor:` for renames and restructuring with no behavioral change
- `tweak:` for small behavioral adjustments to existing features
- `docs:` for documentation changes

Example: `fix: increase API timeout from 5s to 10s to reduce gateway errors`

## Report

Tell the user:
- What changed (files and summary)
- What checks passed
- If TDD was applied and what the test covers

Keep it brief. This is a /upfront:quick change — the report should be quick too.
