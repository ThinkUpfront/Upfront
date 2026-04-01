---
title: /upfront:brainstorm
description: Interactive problem-solving for a specific concern, risk, or tradeoff.
---

## What it does

`/upfront:brainstorm` is a focused conversation about a specific concern. Not feature definition — problem-solving. The AI brings domain knowledge you don't have, maps the landscape, teaches you the risks, and helps you choose between tradeoffs.

## When to use it

- "X seems like a security risk, how should we handle it?"
- "What are the options for caching here?"
- "Is our auth approach solid or are we missing something?"
- "How should we handle rate limiting?"
- Any time you have a specific concern and want to think it through before acting

This skill auto-triggers on phrases like "how should we handle", "is this a problem", "what are the options for", "seems risky."

## How it works

1. **Understand** — restates your concern, asks one clarifying question if needed
2. **Map the landscape** — real risks ranked by likelihood, common approaches, counterintuitive traps
3. **Explore tradeoffs** — explicit option/benefit/cost/failure for each viable approach
4. **Land on next actions** — what to do, what to watch for, what to revisit

## Key principle

This is a teaching conversation. The AI brings knowledge — it doesn't just ask you questions. You came to it because you don't know the landscape. It maps it for you, then lets you decide.

## Next step

If the outcome is a code change, brainstorm routes you to the right skill: `/upfront:quick`, `/upfront:patch`, or `/upfront:feature` depending on size.
