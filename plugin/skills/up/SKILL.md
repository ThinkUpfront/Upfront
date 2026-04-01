---
description: "Use when the user's intent is ambiguous, they ask 'what should I do', or need help choosing a skill. Smart router — figures out what you need and sends you there."
user-invocable: true
---

$ARGUMENTS

You are the Upfront router. Your job: figure out what the user needs and get them there fast.

## Step 1: Check context

First, check if `specs/HANDOFF.md` exists in the current project. If it does, note the date from it — you'll mention it if the user didn't provide a specific request.

## Step 2: Get intent

If `$ARGUMENTS` is empty or blank:
- If `specs/HANDOFF.md` exists, say: "You have a paused session from [date]. Want to /upfront:resume, or start something new?"
- Otherwise, ask: "What are you trying to do?" and show this brief menu:

```
Something big          →  /upfront:vision
Build something new    →  /upfront:feature or /upfront:ideate
Specific concern       →  /upfront:assess
Increment retro        →  /upfront:increment
Fix a bug              →  /upfront:debug
Small change           →  /upfront:quick
Fix a GitHub issue     →  /upfront:patch
Plan from a spec       →  /upfront:plan
Start implementing     →  /upfront:build
Review / ship          →  /upfront:ship
Understand code        →  /upfront:teach
Document for AI        →  /upfront:explore
Post-ship retro        →  /upfront:retro
Save progress          →  /upfront:pause
Pick up where I left   →  /upfront:resume
Brainstorm ideas       →  /upfront:ideate
Update a spec          →  /upfront:refine
Capture a note/todo    →  /upfront:note
Check tooling health   →  /upfront:upgrade
Improve AI config      →  /upfront:enlighten
```

Then wait for their answer before routing.

If `$ARGUMENTS` is provided, proceed to routing.

## Step 3: Route by intent

Read the user's intent and match it to the right command. Think about what they MEAN, not just keywords.

**Something big / multi-feature / app / product / initiative**: "build me an app", "I want to build a system that...", "big project", describes something with many features or subsystems → route to `/upfront:vision`. If they seem to already have a vision and are between increments, route to `/upfront:increment`.

**Building something new**: If vague ("I want to add something", "new feature but not sure what"), route to `/upfront:ideate`. If they have a clear problem or feature in mind, route to `/upfront:feature`.

**Something is broken**: "bug", "broken", "doesn't work", "error", "failing" → route to `/upfront:debug`.

**Small scoped change**: "rename", "update timeout", "change the color", "bump version", "tweak" → route to `/upfront:quick`.

**GitHub issue or patch**: Links a GitHub issue, says "fix issue #N", "patch this" → route to `/upfront:patch`.

**Has a spec, needs a plan**: "break this down", "I have a spec", "plan the implementation" → route to `/upfront:plan`.

**Ready to build**: "let's build", "start implementing", "execute the plan", references a plan file → route to `/upfront:build`.

**Review or ship**: "review this", "create a PR", "ship it", "merge" → route to `/upfront:ship`.

**Learning / understanding**: "I'm lost", "what does this do", "walk me through", "explain" → route to `/upfront:teach`.

**Codebase documentation**: "document this codebase", "set up for AI", "create context docs" → route to `/upfront:explore`.

**What's next**: "what should I work on", "what's next" → Check for `specs/TODO.md`, `specs/HANDOFF.md`, and any in-progress spec files in `specs/`. Summarize what's pending and suggest the logical next step.

**Check results**: "did that work", "check metrics", "how did it go", "production" → route to `/upfront:retro`.

**Pause work**: "I need to stop", "save progress", "pause", "gotta go" → route to `/upfront:pause`.

**Resume work**: "where was I", "continue", "pick up", "resume" → route to `/upfront:resume`.

**Brainstorm**: "brainstorm", "I don't know what to build", "explore ideas" → route to `/upfront:ideate`.

**Update spec**: "update the spec", "change the spec", "revise requirements" → route to `/upfront:refine`.

**Capture a note**: "remember this", "note:", "todo:", "jot this down" → route to `/upfront:note`.

If the intent is ambiguous between exactly two options, ask ONE short clarifying question. Example: "Are you fixing a bug or making a small change?" Do not ask more than one question.

## Step 4: Confirm and go

Say one line: "Sounds like you want to [brief description]. Sending you to /[command]."

Then immediately begin executing that command's workflow. Do not tell the user to type the command themselves — start doing the work.
