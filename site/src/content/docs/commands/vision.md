---
title: /upfront:vision
description: Capture a big ambition, force strategic clarity, and break it into reviewable increments.
---

## What it does

`/upfront:vision` helps you define something bigger than a single feature — an app, a system, a product, a major initiative. It acts as a strategic sparring partner, forcing clarity on who you're building for, why it matters, and what success looks like before any code is written. The output is a living strategy document that guides incremental delivery.

## When to use it

- Starting a new product or major initiative from scratch
- When you have an ambitious idea but haven't nailed down who it's for or why now
- Before breaking work into increments — vision defines what the increments are
- When you need to pressure-test strategy before committing resources

## How it works

1. **Who and Why** — forces specificity on the target audience and their pain. Rejects vague answers like "engineers who want to be more productive"
2. **Diagnosis** — identifies the structural root cause of the problem, not just symptoms. Names the core difficulty — why this hasn't been solved already
3. **Strategy** — extracts guiding policies, key assumptions, anti-vision (what you're NOT building), and kill criteria (when to stop)
4. **Coherent Actions** — sequences the vision into 3-5 increments, each delivering value to a real person and testing a key assumption

## Key principles

- **No vague answers** — if you can't be specific about who, why, or what success looks like, you're not ready to build
- **Value-first sequencing** — increment 1 delivers value to a real person, not infrastructure
- **Every increment is an experiment** — each one tests an assumption and has a learning goal
- **Kill criteria are mandatory** — a project without kill criteria is a zombie
- **Anti-vision prevents scope creep** — every new idea gets checked against what you explicitly said you wouldn't build

## Output

- `specs/[name]-vision.md` — living strategy document with audience, diagnosis, strategy, and sequenced increments
- Updated `specs/DECISIONS.md` — strategic decisions and rationale
- Thinking record capturing what was challenged and sharpened during the conversation

## Next step

Start increment 1: run `/upfront:feature` for the first feature in the increment. After the increment ships, run `/upfront:increment` for a retro before starting the next one.
