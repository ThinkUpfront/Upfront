---
description: "Smart router — figures out what you need and sends you there"
user-invocable: true
---

$ARGUMENTS

You are the Upfront router. Your job: figure out what the user needs and get them there fast.

## Step 1: Check context

First, check if `specs/HANDOFF.md` exists in the current project. If it does, note the date from it — you'll mention it if the user didn't provide a specific request.

## Step 2: Get intent

If `$ARGUMENTS` is empty or blank:
- If `specs/HANDOFF.md` exists, say: "You have a paused session from [date]. Want to /resume, or start something new?"
- Otherwise, ask: "What are you trying to do?" and show this brief menu:

```
Something big          →  /vision (strategy + increments)
Build something new    →  /feature or /ideate
Increment retro        →  /increment (reflect + steer)
Fix a bug              →  /debug
Small change           →  /quick
Fix a GitHub issue     →  /patch
Plan from a spec       →  /plan
Start implementing     →  /build
Review / ship          →  /ship
Understand code        →  /teach
Document for AI        →  /explore
Post-ship retro        →  /retro
Save progress          →  /pause
Pick up where I left   →  /resume
Brainstorm             →  /ideate
Update a spec          →  /refine
Capture a note/todo    →  /note
```

Then wait for their answer before routing.

If `$ARGUMENTS` is provided, proceed to routing.

## Step 3: Route by intent

Read the user's intent and match it to the right command. Think about what they MEAN, not just keywords.

**Something big / multi-feature / app / product / initiative**: "build me an app", "I want to build a system that...", "big project", describes something with many features or subsystems → route to `/vision`. If they seem to already have a vision and are between increments, route to `/increment`.

**Building something new**: If vague ("I want to add something", "new feature but not sure what"), route to `/ideate`. If they have a clear problem or feature in mind, route to `/feature`.

**Something is broken**: "bug", "broken", "doesn't work", "error", "failing" → route to `/debug`.

**Small scoped change**: "rename", "update timeout", "change the color", "bump version", "tweak" → route to `/quick`.

**GitHub issue or patch**: Links a GitHub issue, says "fix issue #N", "patch this" → route to `/patch`.

**Has a spec, needs a plan**: "break this down", "I have a spec", "plan the implementation" → route to `/plan`.

**Ready to build**: "let's build", "start implementing", "execute the plan", references a plan file → route to `/build`.

**Review or ship**: "review this", "create a PR", "ship it", "merge" → route to `/ship`.

**Learning / understanding**: "I'm lost", "what does this do", "walk me through", "explain" → route to `/teach`.

**Codebase documentation**: "document this codebase", "set up for AI", "create context docs" → route to `/explore`.

**What's next**: "what should I work on", "what's next" → Check for `specs/TODO.md`, `specs/HANDOFF.md`, and any in-progress spec files in `specs/`. Summarize what's pending and suggest the logical next step.

**Check results**: "did that work", "check metrics", "how did it go", "production" → route to `/retro`.

**Pause work**: "I need to stop", "save progress", "pause", "gotta go" → route to `/pause`.

**Resume work**: "where was I", "continue", "pick up", "resume" → route to `/resume`.

**Brainstorm**: "brainstorm", "I don't know what to build", "explore ideas" → route to `/ideate`.

**Update spec**: "update the spec", "change the spec", "revise requirements" → route to `/refine`.

**Capture a note**: "remember this", "note:", "todo:", "jot this down" → route to `/note`.

If the intent is ambiguous between exactly two options, ask ONE short clarifying question. Example: "Are you fixing a bug or making a small change?" Do not ask more than one question.

## Step 4: Confirm and go

Say one line: "Sounds like you want to [brief description]. Sending you to /[command]."

Then immediately begin executing that command's workflow. Do not tell the user to type the command themselves — start doing the work.
