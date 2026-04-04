---
description: The entry point. Use when starting a session — "what should I work on", "let's go", "start", or when describing a new idea. Checks project state and routes to the right skill.
user-invocable: true
---

# Start

You are the front door. Figure out where the user is and get them moving fast.

**Input:** $ARGUMENTS

## Step 0: Version check

Silently check if a newer version is available:

```bash
curl -sf --max-time 3 https://api.github.com/repos/ThinkUpfront/Upfront/releases/latest | grep -o '"tag_name":"[^"]*"'
```

Compare with the installed version from:
```bash
cat ~/.claude/plugins/cache/thinkupfront/upfront/*/.claude-plugin/plugin.json 2>/dev/null
```

If a newer version exists, mention it briefly before routing — don't block:

"Update available: Upfront [current] → [latest]. Run `upfront update` in your terminal when you get a chance."

If the check fails or times out, skip silently. Never block the user for a version check.

## Step 1: Read the room

Silently check what exists:
- `specs/HANDOFF.md` — was a session paused?
- `specs/DEBT.md` — any validated spikes waiting for solidification?
- `specs/*-vision.md` — any active visions with experiments?
- `specs/*.md` (not ARCHITECTURE/DECISIONS/LEARNINGS/DEBT/HANDOFF/TODO) — any feature specs without plans?
- `specs/*-plan.md` — any plans not yet built?
- `specs/*-progress.md` — any builds in progress?
- `specs/TODO.md` — any open items?

Run `git log --oneline -5` to see recent activity.

## Step 2: Route

### If there's a handoff

Something was paused. Say:

"You left off mid-session. Want to pick up where you stopped?"

If yes, immediately launch `/upfront:resume`. If no, continue to routing below.

### If the user provided $ARGUMENTS (they know what they want)

Evaluate what they described:

- **Vague / exploratory** ("I'm thinking about...", "maybe something like...", no clear problem): → "Sounds like you're still exploring. Let's sharpen the idea." → Launch `/upfront:ideate`

- **Big / multi-feature** ("build me an app", "I want a platform", multiple concerns): → "This is bigger than one feature. Let's figure out what to test first." → Launch `/upfront:vision`

- **Clear idea, untested** (specific problem, has a solution in mind, but hasn't tried it): → "Let's test this. Build the minimum, see if it works." → Launch `/upfront:spike`

- **Validated / confident** (references a spike, has prior art, knows it works, ready to build properly): → "You know what you want. Let's build it right." → Launch `/upfront:feature`

### If no $ARGUMENTS (they just said "start" or "what's next")

Check project state and suggest the most logical next step:

1. **Validated spike exists** (DEBT.md has status "validated", no corresponding feature spec): → "You have a validated spike: [description]. Ready to solidify it?" → Launch `/upfront:feature`

2. **Feature spec exists without a plan**: → "You have a spec ready: `specs/[name].md`. Ready to plan the build?" → Launch `/upfront:plan`

3. **Plan exists, not yet built**: → "You have a plan ready: `specs/[name]-plan.md`. Ready to build?" → Launch `/upfront:build`

4. **Vision exists with a current experiment**: → "Your vision has experiment [N] queued: [description]. Ready to spike it?" → Launch `/upfront:spike`

5. **Build in progress** (progress file exists): → "You have a build in progress." → Launch `/upfront:build` with the plan path

6. **TODO has open items**: → "You have [N] open items in TODO. Want to tackle one, or start something new?"

7. **Nothing exists**: → "Fresh start. What do you want to build?"

Then route based on their answer (same logic as the $ARGUMENTS path above).

## Rules

- **One question max.** Don't interrogate. Read the state, make a suggestion, let them confirm or redirect.
- **Always launch the skill.** Don't describe what the skill does — just launch it. The user doesn't need a menu, they need momentum.
- **Respect the loop.** The natural flow is: spike → learn → solidify → build. If you see evidence of where they are in this loop, suggest the next step, not a restart.
- **Don't block.** If the user wants to skip your suggestion and do something else, let them. Say "Sure" and route where they want to go.
