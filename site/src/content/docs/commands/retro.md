---
title: /retro
description: Check spec predictions against production reality.
---

The accountability mechanism. Closes the feedback loop between what you predicted and what actually happened.

## When to use

- Feature has been in production long enough to have data
- You want to check if the spec's predictions were right

## What it does

Goes back to the spec's "how will we know it worked?" and checks:

- Each prediction scored: **hit**, **partial**, **miss**, or **unknown**
- Misses analyzed: was the mechanism wrong, the measurement wrong, or the environment different?
- Extracts generalizable lessons
- Pushes for numbers: "I think it improved" gets challenged with "do you have the actual number?"

Feeds forward: suggests changes to `/feature`, `/plan`, or `/build` if the retro reveals a pattern.

## Output

Retro report appended to `specs/LEARNINGS.md`. Pattern suggestions for future specs.
