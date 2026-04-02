---
title: /upfront:spike
description: Rapid prototyping with automatic debt tracking.
---

## What it does

Builds a throwaway prototype to test an idea. Speed over quality. You get something clickable/runnable in minutes — fake data, no backend, happy path only.

## When to use it

- You need to **see** the idea before you can spec it
- You're not sure if the approach works and need to try it
- A stakeholder needs to react to something concrete
- "Just build it and we'll figure out the details"

## How it works

1. **One question:** "Describe what you want to build in as much detail as possible."
2. **AI autocompletes:** Fills every gap, shows what it'll build, what's fake, and what tech it'll use. "Does this look right?"
3. **Acknowledges debt:** Lists everything being deferred — no spec, no tests, no architecture review.
4. **Builds the minimum:** Fake UI, mock data, hardcoded everything. One happy path. Fast.
5. **Demo:** Shows you how to see it, what to click. Then asks: keep it, iterate, or kill it?
6. **Logs debt:** Every shortcut recorded in `specs/DEBT.md`.

## What happens after

When the spike validates the idea, three options:

- **Spec it properly** — `/upfront:feature` launches with the spike as context. The approved UI/flow carries forward as a constraint. Phase 1-2 of feature definition go faster because the prototype already answered most questions.
- **Iterate** — keep tweaking the spike. After 2-3 rounds, the skill nudges you toward spec or kill.
- **Kill it** — revert the code or keep it for reference. Debt items are marked as killed.

The full path after a successful spike:

```
/upfront:spike → /upfront:feature → /upfront:plan → /upfront:build
```

The `/upfront:build` merge-back step automatically closes the spike's debt items — the feature spec replaces "no spec", TDD replaces "no tests", architecture review replaces "no arch review".

## Debt tracking

Every spike logs to `specs/DEBT.md`:
- What was built
- What files were touched
- What was skipped (spec, tests, arch review, specific items like mock auth)
- Severity rating

This debt balance surfaces in:
- `/upfront:increment` — retros show accumulated debt
- `/upfront:upgrade` — health check counts open items
