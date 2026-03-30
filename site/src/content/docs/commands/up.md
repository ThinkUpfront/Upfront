---
title: /upfront:up
description: Smart router that figures out what you need and sends you there.
---

## What it does

`/upfront:up` is the starting point when you're not sure which skill to use. It reads your intent and routes you to the right skill, or shows you the menu if you come in blank. If you have a paused session, it offers to resume it.

## When to use it

- When you're not sure which skill fits your situation
- At the start of a session when you want to see what's available
- When you have a paused session and want to pick up where you left off

## How it works

1. **Check context** — looks for `specs/HANDOFF.md` from a previous `/upfront:pause`
2. **Get intent** — if you provide arguments, it reads them. If not, it shows the menu or offers to resume a paused session
3. **Route by intent** — matches what you mean (not just keywords) to the right skill:

| Intent | Route |
|--------|-------|
| Something big / multi-feature | `/upfront:vision` |
| Build something new | `/upfront:feature` or `/upfront:ideate` |
| Increment retro | `/upfront:increment` |
| Fix a bug | `/upfront:debug` |
| Small change | `/upfront:quick` |
| GitHub issue | `/upfront:patch` |
| Plan from a spec | `/upfront:plan` |
| Start implementing | `/upfront:build` |
| Review / ship | `/upfront:ship` |
| Understand code | `/upfront:teach` |
| Document for AI | `/upfront:explore` |
| Post-ship retro | `/upfront:retro` |
| Save progress | `/upfront:pause` |
| Resume work | `/upfront:resume` |
| Brainstorm | `/upfront:ideate` |
| Update a spec | `/upfront:refine` |
| Capture a note | `/upfront:note` |

4. **Confirm and go** — tells you where it's routing and immediately starts that skill's workflow

## Key principles

- **Route by meaning, not keywords** — understands intent even when phrased differently
- **One clarifying question max** — if ambiguous between two options, asks one short question
- **Starts the work** — doesn't tell you to type a command; it begins executing immediately

## Output

Routes you to the appropriate skill and begins its workflow.

## Next step

If you already know what you need, call the skill directly. Use `/upfront:up` when you don't.
