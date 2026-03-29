---
description: Fix a bug or implement a small feature from a GitHub issue or clear problem statement — more structure than /quick, less ceremony than /feature
---

# Patch

You are implementing a fix or small feature where the problem is already understood. This is for GitHub issues, bug reports, and well-scoped tasks that don't need ideation or a full feature spec — but are too big or too risky for `/quick`.

**Input:** $ARGUMENTS

The input is one of:
- A GitHub issue URL or number (e.g., `#42`, `owner/repo#42`, `https://github.com/...`)
- A clear problem description

## Auto mode

If the input contains `--auto` (e.g., `/patch --auto #42`), run without pausing for confirmation:
- Skip the diagnosis confirmation step (Step 2 — still present the diagnosis, but proceed immediately)
- Skip the "Want me to close the issue?" prompt (close it automatically with a summary)
- Still stop if: scope exceeds 300 lines, root cause is unclear, or a constitutional principle is violated

Auto mode is for well-understood issues where you trust the pipeline. Remove `--auto` from the input before processing the rest as the issue reference.

## Scope gate

This command handles changes up to ~300 lines of non-test code. If it's clearly bigger:

> "This looks like it needs `/feature` + `/plan` + `/build`. Here's why: [reason]."

If it's clearly under ~50 lines and straightforward:

> "This is small enough for `/quick`. Want me to run that instead?"

## Process

### 1. Understand the problem

**If given a GitHub issue:** Read it with `gh issue view`. Extract:
- What's broken or missing
- Reproduction steps (if a bug)
- Acceptance criteria (if stated)
- Any constraints or context from comments

**If given a description:** Work with what you have.

Then read context files if they exist:
- `specs/ARCHITECTURE.md` (its invariants section contains hard constraints — if the fix would violate one, stop and say so)
- `specs/DECISIONS.md`
- `specs/LEARNINGS.md`

### 2. Investigate the codebase

Read the code that's relevant to the problem. Understand:
- Where the bug lives or where the feature fits
- What patterns the surrounding code follows
- What tests already exist in this area
- What could break if you change this

Present a brief diagnosis:

```
Problem: [what's wrong or missing]
Root cause: [why it happens — for bugs]
Fix approach: [what you plan to change]
Files: [list of files you'll touch]
Risk: [what could go wrong, or "low — isolated change"]
```

Wait for the user to confirm the approach. If the investigation reveals the problem is different from what the issue describes, say so.

### 3. Implement with TDD

For each unit of change:

1. **Write a failing test** that captures the expected behavior. Run it. Confirm it fails.
2. **Implement the fix.** Run the test. Confirm it passes.
3. **Run the full test suite** for the affected area. Fix anything you broke.

For changes that aren't unit-testable (config, docs, build scripts), state why and verify through other means.

If the change grows beyond ~300 lines during implementation, stop:

> "This grew to ~[N] lines. Want me to continue as-is, or should we switch to `/feature` + `/plan` + `/build` for proper phasing?"

### 4. Self-review

Review your own diff. Check:
- Does this actually fix the problem described in the issue?
- Are there edge cases the issue didn't mention but the code reveals?
- Does it follow the patterns in the surrounding code?
- Does it violate anything in `specs/DECISIONS.md` or `specs/ARCHITECTURE.md`?
- Could the fix introduce a regression? Is there a test for that?
- Did this patch introduce a design decision that future work should know about? If so, append a brief entry to `specs/DECISIONS.md`.

If you find issues, fix them and re-run tests.

If investigation reveals the root cause is unclear and requires hypothesis testing, suggest `/debug` instead.

### 5. Run all checks

Run every linter, formatter, type checker, and test suite configured for the project. Not just the tests you wrote — everything. Fix anything that fails.

### 6. Commit

Commit with a message that references the issue (if there was one):

```
fix: [description of what changed]

Fixes #[issue number]
```

Use conventional prefixes:
- `fix:` for bug fixes
- `feat:` for small features
- `refactor:` for structural changes with no behavioral change

### 7. Report

```
Patch applied:
- [what changed, in brief]
- [files modified]
- Tests: [what's covered]
- Checks: [what passed]
```

If the input was a GitHub issue, ask: "Want me to close the issue with a comment summarizing the fix?"

## Rules

- Do NOT over-engineer. This is a patch, not an architecture overhaul.
- Do NOT refactor code adjacent to the fix unless the refactoring IS the fix.
- Do NOT add features beyond what the issue asks for.
- If investigation reveals a deeper problem than the issue describes, present it to the user. Don't silently expand scope.
- If you find other bugs while investigating, note them but don't fix them. One patch, one problem.
