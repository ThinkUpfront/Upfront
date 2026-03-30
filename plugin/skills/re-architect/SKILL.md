---
description: Execute an architecture evolution plan phase by phase — restructuring only, tests verified at phase boundaries
user-invocable: true
---

# Re-Architect

You are executing an architecture evolution plan. This is **restructuring only** — no new features, no behavior changes, no new tests. Existing tests must pass at the end of each phase.

The key difference from `/build`: tests WILL break during a phase as you move files and update imports. That's expected. You only verify at phase boundaries, not per-edit. The post-edit hook is NOT installed for re-architect runs.

## Input

The user provides a path to an evolution doc (e.g., `specs/ARCHITECTURE-EVOLUTION.md`), or it's detected automatically if only one exists.

Read it fully. Understand the priority order, the constraints, and what each phase changes.

## Rules

- **Restructuring only.** Do not change behavior. Do not add features. Do not "improve" logic while you're moving it.
- **Tests break mid-phase, pass at phase end.** Each phase leaves the codebase compiling and all existing tests passing. Between those states, things will be broken. That's fine.
- **No new tests per phase.** The behavior hasn't changed — existing tests cover it. If an existing test needs its import path updated, update it. Do not write new tests during restructuring.
- **Respect the phase order.** The evolution doc is prioritized. God module splits come before layer fixes. Layer fixes come before pattern introductions. Do not skip ahead.
- **Do not use the post-edit hook.** Restructuring phases intentionally break things mid-phase. Per-edit lint/test would block every move operation. Phase-level verification only.

## Process

### For each phase:

#### 1. Announce

```
[████████░░░░░░░░░░░░░░░░] 2/5 — Extract auth from core package

Moving: internal/core/auth.go, internal/core/auth_test.go → internal/auth/
Updating: [N] import paths across [M] files
```

#### 2. Execute the restructuring

Work through the phase's changes:
- Move files to their new locations
- Update all import paths / require statements
- Update any internal references (function calls, type references)
- Fix test imports

Do this methodically — move one module at a time, update all references, then move the next. Don't move everything at once and hope for the best.

#### 3. Verify at phase boundary

Run the full build and test suite:
```bash
[project's build command]
[project's test command]
```

**If tests pass:** Phase complete. Commit and announce:
```
refactor([area]): [phase N]/[total] [description]

Architecture evolution: specs/ARCHITECTURE-EVOLUTION.md
```

**If tests fail:** Fix the failures. These should only be import path issues or reference updates you missed. If a test failure indicates a behavioral change, you've done something wrong — revert and re-approach.

**If build fails:** Same — should only be missing imports or moved references. Fix them.

#### 4. Update ARCHITECTURE.md

After each phase, update `specs/ARCHITECTURE.md` to reflect the new structure. The architecture doc must always match reality.

#### 5. Report and proceed

```
[████████████████░░░░░░░░] 3/5 — Layer extraction ✓

Tests: [N] passing
Build: clean
ARCHITECTURE.md: updated

Proceeding to Phase 4...
```

Auto-proceed between phases (no manual verification needed — this is restructuring, not new behavior). Pause only if a phase fails verification.

### After all phases

1. Run the full test suite one final time
2. Run linters and formatters
3. Update `specs/ARCHITECTURE.md` with the final state and "Last reviewed" date
4. Commit any final ARCHITECTURE.md updates

### Capture learnings

Append to `specs/LEARNINGS.md`:

```markdown
## [date] — Architecture evolution
**What changed:** [summary of phases executed]
**Surprises:** [hidden couplings found, unexpected dependencies, things that were harder to move than expected]
**What the agent got wrong:** [if any phases needed multiple attempts, what went wrong]
```

### Record decisions

For each phase that made a structural choice, append to `specs/DECISIONS.md`:

```markdown
## [date] — Restructuring: [phase description]
**Decision:** [what was moved/split/extracted]
**Rationale:** [why this structure is better]
```

Tell the user:
- All phases complete
- Test results
- Learnings and decisions captured
- "The architecture has been restructured. Run `/architect` again in 3-5 features to check for drift."

## Stuck loop detection

If a phase's verification fails more than 3 times with the same error, stop:

```
Stuck: Phase [N] verification has failed 3 times.
Error: [what keeps failing]
This might indicate a behavioral dependency on the old structure.

Options:
  a) I investigate the dependency and adjust the approach
  b) Skip this phase and continue with the rest
  c) You look at this manually
```
