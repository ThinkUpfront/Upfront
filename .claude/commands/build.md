---
description: Execute an implementation plan phase by phase with strict TDD and post-phase review
---

# Build

You are the **orchestrator** for implementing an approved plan. You do not implement phases yourself. You spawn a fresh sub-agent for each phase, verify its work, run a post-phase review, and manage the progress file. Each phase gets a clean context with maximum leverage.

This follows the RALPH pattern: fresh context per task, progress preserved through git commits and a progress file, not through conversation history.

## Input

The user provides a path to a plan file (e.g., `specs/feature-name-plan.md`). The plan file contains a reference to its spec at the top (e.g., `Spec: specs/feature-name.md`).

Read both:
- The plan file (provided by user)
- The spec file (linked from the plan's header)

Read both fully before starting. Understand the intent, constraints, scope boundaries, and blind spots from the spec. Understand the phasing, file changes, and verification criteria from the plan.

If no plan path is provided, check the `specs/` directory for plan files and ask which one to build.

## The Progress File

Create `specs/[feature-name]-progress.md` before starting Phase 1:

```markdown
# Progress: [feature name]

> Spec: `specs/[feature-name].md`
> Plan: `specs/[feature-name]-plan.md`

## Completed Phases

(none yet)

## Learnings

(none yet)
```

This file is the memory between phases. Each sub-agent reads it. You update it after each phase. It captures what happened, what was surprising, and what the next phase should know.

## Rules

- **You are the orchestrator, not the implementer.** Spawn sub-agents for each phase. Do not write implementation code yourself.
- **TDD is non-negotiable.** Every sub-agent must write failing tests before implementation code. If a sub-agent skips TDD, reject its work and re-spawn with explicit TDD instructions.
- Follow the plan. Do not improvise, add features, or deviate from scope.
- Respect the spec's scope boundaries — do not touch files listed as non-goals.
- Do not skip verification steps.
- Do not proceed to the next phase without human confirmation.

## Process

### For each phase:

#### 1. Announce the phase

Tell the user what's about to happen:
```
Starting Phase [N]: [descriptive name]
Files: [list from plan]
Expected outcome: [what changes]
Spawning fresh agent...
```

#### 2. Spawn a sub-agent with clean context

Use the Agent tool to spawn a fresh sub-agent. Give it exactly this context:

```
You are implementing Phase [N] of a feature plan. You have a clean context — do not assume knowledge of prior phases. Read the codebase as it is now.

## Your task
[paste the specific phase from the plan, including files, changes, and verification criteria]

## Spec (intent and constraints)
Read the spec at: specs/[feature-name].md
Pay attention to: scope boundaries, "what must NOT happen", and known AI blind spots.

## Progress from prior phases
Read: specs/[feature-name]-progress.md
This contains learnings from prior phases that may affect your work.

## TDD Protocol — STRICTLY ENFORCED

You MUST follow strict TDD for every piece of logic. No exceptions. No "I'll add tests after."

For each unit of functionality:
1. **Red:** Write a failing test FIRST that describes the expected behavior. Run it. Confirm it fails. Show the failure output.
2. **Green:** Write the MINIMUM code to make the test pass. Run it. Confirm it passes. Show the passing output.
3. **Refactor:** Clean up while keeping tests green. Run tests again.

If you catch yourself writing implementation code without a failing test, STOP. Delete the implementation code. Write the test first.

For changes that are genuinely not unit-testable (config files, static assets, documentation, build scripts), explicitly state: "TDD does not apply to [file] because [reason]." Verify through other means (build succeeds, config parses, etc.).

## Rules
- Implement ONLY Phase [N]. Do not touch anything outside this phase's scope.
- Follow existing code conventions and the spec's style/architecture rules.
- If you encounter something that doesn't match the plan, STOP and report:
  Issue: [what's wrong]
  Expected: [what the plan says]
  Found: [actual situation]
  Do not guess. Do not proceed.
- Run ALL automated verification commands listed for this phase before finishing.
- When done, report: what you changed, what tests pass, any surprises or learnings for future phases.
```

#### 3. Review the sub-agent's work

When the sub-agent completes, review:
- Did it follow TDD? (tests written before implementation, red-green-refactor cycle visible)
- Did it stay within phase scope? (no extra files changed, no scope creep)
- Did automated verification pass?
- Did it report any mismatches or surprises?

**If TDD was not followed:** Reject the work. Tell the user, re-spawn the sub-agent with stronger TDD emphasis. Do not accept implementation-first code.

If the sub-agent hit an issue and stopped, present it to the user:
```
Issue in Phase [N]:
Expected: [what the plan says]
Found: [actual situation]
How should I proceed?
```
Do not guess. Wait for the user.

#### 4. Run automated verification yourself

Do not trust the sub-agent's self-report alone. Run every automated verification command listed for this phase independently. If any fail, either fix or re-spawn the sub-agent with the failure context.

#### 5. Post-phase code review

After automated verification passes, spawn a **review agent** to audit the phase:

```
You are reviewing code just written for Phase [N] of a feature implementation. Your job is to find problems, not praise.

Read:
- Spec: specs/[feature-name].md (constraints, blind spots, "what must NOT happen")
- Plan: specs/[feature-name]-plan.md (Phase [N] scope)
- Progress: specs/[feature-name]-progress.md (context from prior phases)
- All files changed in this phase (git diff HEAD~1)

## Review checklist

### Correctness
- Does the code do what the spec says? Not what seems reasonable — what the spec actually requires.
- Are edge cases from the spec's blind spots handled?
- Do the tests actually test the right behavior, or do they test implementation details?
- Are there untested code paths?
- Could any test pass with a broken implementation? (Tests that are too loose.)

### Architecture
- Does the code follow existing patterns and conventions in the codebase?
- Are responsibilities in the right place? (No logic in the wrong layer.)
- Is the code structured for the phases that come after it? (Check the plan — will Phase N+1 be able to build on this cleanly?)
- Are interfaces clean? Will they need to change for future phases?
- Any unnecessary abstractions, premature generalization, or over-engineering?

### Spec compliance
- Does anything violate the spec's "what must NOT happen" section?
- Does the code stay within scope boundaries?
- Are the spec's constraints and non-negotiables respected?

### Issues to flag
For each issue found, categorize:
- **MUST FIX** — incorrect behavior, spec violation, missing edge case, security issue
- **SHOULD FIX** — architecture concern that will cause problems in later phases
- **NOTE** — observation for awareness, no action needed now

Report your findings. If there are MUST FIX items, list them clearly.
```

If the review finds **MUST FIX** items: fix them (spawn another sub-agent if needed), re-run verification, re-review. Do not proceed with known correctness issues.

If the review finds **SHOULD FIX** items: present them to the user. Let them decide whether to fix now or accept the risk.

#### 6. Update the plan and progress file

Check off completed automated verification items in the plan file.

Update `specs/[feature-name]-progress.md`:
```markdown
## Phase [N]: [name] — COMPLETE

**What changed:** [files modified/created]
**TDD cycles:** [number of red-green-refactor cycles]
**Review findings:** [MUST FIX items resolved, SHOULD FIX items accepted/deferred, or "clean"]
**Surprises:** [anything unexpected, or "none"]
**Learnings for future phases:** [anything the next agent should know]
```

#### 7. Commit

Commit the phase immediately. Use a sequential ID so phases are easy to revert individually:

```
feat([feature-name]): [N]/[total] [phase description]

Phase [N] of specs/[name]-plan.md
```

Example: `feat(upfront): 3/6 Local JSONL queue with concurrent write safety`

#### 8. Report and pause

```
Phase [N]/[total] committed.

Automated verification passed:
- [list what passed]

Code review:
- [summary of review findings — clean, or what was fixed/accepted]

Manual verification items (if any):
- [list manual items from the plan]
```

Wait for user confirmation before proceeding to the next phase. The user can interrupt, redirect, or ask for changes at any point.

### After all phases: Integration sweep

Spawn a **final sweep agent** with clean context:

```
You are verifying that a multi-phase feature implementation connects correctly. All phases have been committed individually. Your job is to verify integration — that the pieces work together, not just individually.

Read:
- Spec: specs/[feature-name].md (especially acceptance criteria and blind spots)
- Plan: specs/[feature-name]-plan.md (all phases, for full picture)
- Progress: specs/[feature-name]-progress.md (learnings and surprises)

Then:
1. Read the current codebase state — all changes are committed.
2. Run the full test suite. Every test must pass.
3. Check for:
   - Cross-phase integration issues (Phase 2 depends on Phase 1's output — does it actually work?)
   - Gaps between phases (anything that fell between the cracks)
   - Spec blind spots (edge cases, concurrency, security, non-prompted concerns)
   - Unnecessary code, dead paths, or unused exports
4. Report: what passes, what fails, what's suspicious.
```

If the sweep finds issues, fix them (spawn another sub-agent if needed) and re-run.

Then tell the user:
- All phases are complete
- Integration sweep results
- Any issues encountered during implementation
- The feature is ready for final review

## Resuming

If the user re-runs `/build` on a plan with some phases already completed:
- Read `specs/[feature-name]-progress.md` to understand what's done
- Trust that completed phases are done (they're committed)
- Pick up from the first unchecked phase
- The progress file gives the new context everything it needs — no conversation history required

## Why fresh context per phase

Long-running AI sessions degrade. By phase 5 of a sequential implementation, the context is stuffed with diffs, test output, and verification results from phases 1-4. The AI is working with less effective context for the hardest phases.

Fresh context per phase means:
- Every phase gets the full context budget
- The AI reads the codebase as it actually is (post prior commits), not as it remembers it
- No accumulated noise from prior implementation details
- Progress and learnings are preserved in the progress file, not in conversation memory
- If a phase fails badly, you can re-run it with zero contamination from the failed attempt
