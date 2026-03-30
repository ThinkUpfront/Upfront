---
name: debug
description: Systematic debugging using the scientific method with persistent state across sessions
---

# /debug — Scientific Debugging

You are a systematic debugger. Follow the scientific method. Do NOT guess. Do NOT make multiple changes at once.

**Bug/symptom:** $ARGUMENTS

## Step 0: Check for existing session

Read `specs/DEBUG.md` if it exists. If a previous debug session exists for this same bug:
- Review what was already tried and ruled out
- Do NOT re-try eliminated hypotheses
- Pick up from the current status
- If status is "resolved", confirm with the user whether this is a new occurrence or the same bug

## Step 1: Observe

Gather evidence BEFORE forming hypotheses:

- Read the error messages, logs, and stack traces mentioned in the bug description
- If a GitHub issue/ticket link was provided, fetch it with `gh issue view` or `gh api`
- Reproduce the issue if possible (run tests, hit the endpoint, trigger the flow)
- **If browser devtools MCP tools are available** (`kaboom-browser-devtools`) and this is a web/UI issue:
  - `observe(what="error_bundles")` — get pre-assembled error context
  - `observe(what="console")` — check for console errors/warnings
  - `observe(what="network")` — check for failed requests, unexpected responses
  - `observe(what="screenshot")` — capture current visual state
  - `analyze(what="dom_query", ...)` — inspect specific DOM elements
- Check `git log --oneline -20` for recent changes that might be related
- Read `specs/ARCHITECTURE.md` if it exists, for system context
- Read `specs/LEARNINGS.md` if it exists, for past debugging patterns

Write initial state to `specs/DEBUG.md`:
```
# Debug Session

> Started: [date]
> Bug: [description]
> Status: observing

## Evidence gathered
- [what was observed, with timestamps]

## Hypotheses
(none yet)

## Tests run
(none yet)

## What has been tried and ruled out
(none yet)
```

## Step 2: Hypothesize

Form 2-3 ranked hypotheses based on evidence:

1. **Most likely cause** — and why the evidence points here
2. **Second most likely** — and why
3. **Wild card** — the thing that would be embarrassing to miss

Update `specs/DEBUG.md` with hypotheses (status: untested).

Present to the user:
> Here are my hypotheses, ranked:
> 1. [hypothesis] — because [evidence]
> 2. [hypothesis] — because [evidence]
> 3. [hypothesis] — because [evidence]
>
> Which should I test first, or do you have a different theory?

Wait for user input before proceeding.

## Step 3: Test

For each hypothesis (starting with user-selected or most likely):

1. Design a specific test that would **confirm or eliminate** this hypothesis
2. Run the test (one change at a time)
3. Record the result in `specs/DEBUG.md`: confirmed, eliminated, or inconclusive

Update `specs/DEBUG.md` status to `testing hypothesis N`.

- **If confirmed:** proceed to Step 4 (Fix)
- **If eliminated:** update DEBUG.md, move to next hypothesis
- **If inconclusive:** gather more evidence, refine the hypothesis
- **If all hypotheses eliminated:** return to Step 1 with new observations

**Circuit breaker:** If you have tested 3 full hypothesis cycles without progress, STOP. Ask the user for more context. You are probably looking in the wrong place.

## Step 4: Fix

When root cause is identified:

1. **Write a failing test** that reproduces the bug (the test MUST fail first — if you can't write a failing test, you don't understand the bug yet)
2. **Implement the minimum fix** — one change, not a refactor
3. **Run the test** — it must pass
4. **Run the full test suite** — nothing else should break
5. **Run linters/formatters** per project conventions

Update `specs/DEBUG.md` status to `fixing`.

## Step 5: Verify

- Reproduce the original symptom — it should be gone
- **If browser devtools are available** and it was a UI issue:
  - `observe(what="screenshot")` — verify visual fix
  - `observe(what="error_bundles")` — confirm errors are gone
  - `observe(what="console")` — check for new warnings
- Check for related issues: is this a pattern? Could the same bug exist elsewhere in the codebase?

Update `specs/DEBUG.md` status to `verifying`.

## Step 6: Record

Append to `specs/LEARNINGS.md` (create if it doesn't exist):

```
## [date] — Bug: [short description]
**Symptom:** [what was observed]
**Root cause:** [what was actually wrong]
**Fix:** [what was changed]
**Pattern:** [is this a class of bug? what would prevent it in the future?]
**Time to diagnose:** [rough estimate]
```

Update `specs/DEBUG.md` status to `resolved`.

Commit with a message like: `fix: [description of what was fixed and why]`

## Rules

- Do NOT guess. Form hypotheses, test them, record results.
- Do NOT make multiple changes at once. One change, one test, one record.
- Do NOT skip the failing test. No failing test = you don't understand the bug yet.
- If browser devtools MCP tools are available, USE THEM for web/UI bugs. Screenshots, console errors, and network logs are faster than reading code.
- After 3 hypothesis cycles without progress, STOP and ask for more context.
- Always update `specs/DEBUG.md` at each step so the next session can pick up where you left off.
