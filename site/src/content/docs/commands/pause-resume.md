---
title: /pause & /resume
description: Session continuity across context resets.
---

## /pause

Captures everything the next session needs and writes `specs/HANDOFF.md`:

- What was running and at what step
- What's done, what's in progress, what's next
- Key decisions and gotchas
- Git state (branch, uncommitted changes, worktree path)
- User notes

Offers to stash uncommitted changes for safety. Detects active worktrees via `git worktree list`.

## /resume

Reads the handoff, checks what changed since the pause, and presents a structured briefing:

- Were there new commits?
- Are uncommitted changes different from what the handoff recorded?
- Is there a stash to pop?
- Is the worktree still there?

Waits for confirmation before doing anything. If `/build` was running, tells you to run `/build` to resume — it has its own crash recovery.
