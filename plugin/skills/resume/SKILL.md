---
description: Resume from a previous pause with full context restoration
user-invocable: true
---

# Resume

You are restoring context from a previous session's handoff. Read everything, present a structured briefing, and wait for confirmation before doing anything.

## Step 1: Find the handoff

Read `specs/HANDOFF.md`.

If it does not exist, say:

> No handoff file found. What were you working on?

Stop here. Do not guess or improvise.

## Step 2: Load supporting context

Silently read whatever exists:
- `specs/ARCHITECTURE.md`
- `specs/DECISIONS.md`
- `specs/LEARNINGS.md`

If the handoff references `/build`, also read:
- The plan file mentioned in the handoff
- The spec file linked from the plan
- The `specs/*-progress.md` file for that feature

## Step 3: Check what changed since the pause

Run `git status` and `git log --oneline -10`.

Compare against the git state recorded in the handoff:
- Same branch? If not, note it.
- New commits since the pause? If so, list them — someone (or something) made changes.
- Uncommitted changes that weren't in the handoff? Flag them.
- Uncommitted changes from the handoff that are now gone? Flag that too (someone committed or discarded them).

## Step 4: Present the briefing

```
Resuming from pause on [date from handoff]:

You were running: [command from handoff]
Phase/step: [phase from handoff]

Completed:
- [completed items from handoff]

In progress: [what was being worked on]
Next action: [next action from handoff]

Key context from last session:
- [decisions, in brief]
- [gotchas, in brief]

Git state:
- Branch: [current branch]
- Uncommitted changes: [current state, noting any differences from handoff]
- Changes since pause: [new commits or "none"]

Active files:
- [file list from handoff]
```

If the user left notes in the handoff, include them:

```
Notes from last session:
- [user notes]
```

## Step 5: Check for stashed work

If the handoff's "Git state" section mentions a stash (from `/pause`), remind the user:

> "There's a stash from the previous session. Run `git stash pop` to restore your changes before continuing."

This applies regardless of which command was running — `/build`, `/quick`, `/feature`, or freeform.

## Step 6: Check for branch mismatch

If the current branch differs from the handoff's recorded branch:

- **If the handoff records a worktree** (e.g., "Worktree: ../repo-build-feature"): The mismatch is expected — the user's terminal is on the base branch while the build work lives in the worktree. Do NOT warn about this. Verify the worktree still exists by running `git worktree list` and checking for the recorded path. If it exists, note: "Build worktree at `[worktree path]` on branch `[worktree branch]`. `/build` will resume there." If it's gone, note: "The worktree at `[path]` no longer exists. `/build` will recreate it when you resume."
- **Otherwise:** Warn explicitly: "You're on branch `[current]` but the handoff was created on `[handoff branch]`. Switch back with `git checkout [handoff branch]` before continuing, or confirm you want to work on this branch instead." Wait for the user to resolve before proceeding.

## Step 7: Handle /build resumption

If the handoff says `/build` was running:

1. Present the briefing above first.
2. After user confirms, tell them: "Run `/build [plan-file-path]` to resume. It will detect completed phases and pick up where you left off." **Do not attempt to continue the build yourself** — `/build` has its own crash recovery protocol that handles uncommitted changes, progress reconciliation, and phase resumption. Your job is the context bridge, not the build executor.

## Step 8: Wait

After presenting the briefing, ask:

> Ready to continue?

Do not start working until the user confirms. They may want to adjust priorities, add context, or change direction based on the briefing.
