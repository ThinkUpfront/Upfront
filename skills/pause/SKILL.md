---
description: Pause mid-work with a structured handoff for the next session
---

# Pause

You are creating a handoff file so the next session can pick up exactly where this one left off. Be fast — read the conversation context and progress files, capture everything, get out.

**User notes:** $ARGUMENTS

## Step 1: Gather context

Silently read whatever exists:
- `specs/ARCHITECTURE.md`
- `specs/DECISIONS.md`
- `specs/LEARNINGS.md`
- Any `specs/*-progress.md` file (if `/build` was running)

Run `git status`, `git log --oneline -5`, and `git worktree list` to capture current git state (including any active build worktrees).

Do not ask the user questions. Extract everything from the conversation history and the files above.

## Step 2: Determine what was running

Figure out which workflow was active:
- `/feature` — defining a feature spec
- `/plan` — breaking a spec into phases
- `/build` — executing phases from a plan
- `/quick` — small scoped change
- Freeform — no structured command, just working

If `/build` was running, identify the plan file path and which phase was in progress.

## Step 3: Write the handoff

Create `specs/HANDOFF.md` (overwrite if it exists):

```markdown
# Handoff

> Paused: [date and time]

## What was running

**Command:** [/feature | /plan | /build path/to/plan.md | /quick | freeform]
**Phase/step:** [e.g., "Phase 3 — hook parser" or "Step 2: scope check" or "implementing the retry logic"]

## Completed

- [List what got done this session, be specific]
- [Include commits made: hash + message]

## In progress

[What was actively being worked on when we stopped. Be specific enough that a fresh session can pick up mid-task.]

## Next action

[The single most immediate thing the next session should do. Not a list of everything remaining — just the next step.]

## Decisions made

- [Key decisions from this session that affect upcoming work]
- [Include the WHY, not just the WHAT]

## Gotchas

- [Things that took time to discover and would be painful to re-learn]
- [Unexpected API behavior, edge cases found, wrong assumptions corrected]
- [Include "X does NOT work because Y" style warnings]

## Git state

- **Branch:** [branch name]
- **Worktree:** [path if /build created a worktree, or "none"]
- **Uncommitted changes:** [list of modified/new files, or "clean"]
- **Recent commits:**
  - [hash] [message]
  - [hash] [message]

## Active files

[Files that were being actively edited or are central to the current task]

- `path/to/file.go`
- `path/to/file_test.go`

## User notes

[Anything the user said to remember, or "none"]
```

## Step 4: Protect uncommitted work

If `git status` shows uncommitted changes (modified or new files), offer to stash them:

```
You have uncommitted changes. Want me to stash them for safety?
  git stash push -m "pause: [brief description of current work]"

This protects the work if the terminal closes or the branch changes. /resume will remind you to pop the stash.
```

If the user agrees, run the stash and note the stash ref in the handoff under "Git state." If they decline, note "User declined stash — uncommitted changes at risk" in the handoff.

## Step 5: Append user notes

If the user provided $ARGUMENTS (e.g., "Also remember that the webhook endpoint changes next sprint"), append those under "User notes" in the handoff.

## Step 6: Confirm

Tell the user:

```
Handoff saved to specs/HANDOFF.md

Next session: run /resume to pick up where you left off.
```

Keep it short. The user is stopping — don't make them read a wall of text.
