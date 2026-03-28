---
description: Resume from a previous pause with full context restoration
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

## Step 5: Handle /build resumption

If the handoff says `/build` was running:

1. Present the briefing above first.
2. After user confirms, follow `/build`'s crash recovery protocol:
   - Check for uncommitted changes and present options (keep/stash/discard)
   - Reconcile the progress file with git history
   - Pick up from the first unverified phase
3. Do not re-run the pre-flight audit.

This integrates with `/build`'s existing resume logic — you are just providing the context bridge between sessions.

## Step 6: Wait

After presenting the briefing, ask:

> Ready to continue?

Do not start working until the user confirms. They may want to adjust priorities, add context, or change direction based on the briefing.
