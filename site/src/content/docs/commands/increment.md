---
title: /upfront:increment
description: Retro on a shipped increment and steer the next one.
---

## What it does

`/upfront:increment` runs a structured retro after an increment ships. It forces reflection on what happened — what worked and why, what surprised you, whether the architecture still holds — before deciding what comes next. It's the steering wheel between increments.

## When to use it

- After shipping an increment defined by `/upfront:vision`
- Between major phases of work when you need to reflect before continuing
- When you want to check key assumptions and kill criteria against reality
- Works with or without a vision file — brownfield projects and tactical work benefit too

## How it works

1. **What shipped** — summarizes delivered work from git history and completed specs
2. **What worked and WHY** — rejects surface answers like "it went well." Forces the causal chain: why did it work?
3. **Lessons learned** — probes for estimation misses, assumption failures, process friction, and technical surprises
4. **Architecture check** — verifies `specs/ARCHITECTURE.md` still matches reality and flags accumulated debt
5. **What to adjust** — checks key assumptions (confirmed or busted), evaluates kill criteria honestly, and identifies scope or priority shifts
6. **Next increment** — proposes the next increment based on what was learned, not just what was originally planned

## Key principles

- **"It went well" is not a retro** — push for why something worked, not just that it worked
- **Kill criteria are checked honestly** — if the evidence says stop, it says stop
- **The vision serves the user** — if reality doesn't match the vision, update the vision
- **Retros are not optional** — 2 minutes of reflection prevents the next increment from repeating mistakes

## Output

- Updated vision file — completed increment marked SHIPPED with retro summary, assumptions updated, future increments adjusted
- Updated `specs/LEARNINGS.md` — transferable insights, surprises, architecture changes, and process notes
- Proposed next increment with value delivery, assumption testing, and learning goals

## Next step

When ready to start the next increment, run `/upfront:feature` for the first feature.
