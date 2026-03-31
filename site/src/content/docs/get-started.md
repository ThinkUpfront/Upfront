---
title: Get Started
description: You only need one command. The system handles the rest.
---

## Install

```bash
claude plugin marketplace add ThinkUpfront/Upfront
claude plugin install upfront
```

Restart Claude Code. Done.

## Your first feature

Type this:

```
/upfront:feature
```

That's it. The AI takes over from here:

1. **It asks what problem you're solving** — not what you want to build, what problem goes away. It pushes back if your answer is vague.
2. **It walks you through the behavioral spec** — user stories, mechanism, states, concurrency, error cases. It challenges your thinking at each step.
3. **It researches your codebase** — finds relevant patterns, presents options, lets you make the decisions.
4. **It proposes the implementation design** — and challenges the codebase itself, flagging structural issues before you build on them.

At the end, you have a spec in `specs/` that captures every decision, why it was made, and what was rejected.

## What happens next

The AI tells you. At the end of `/upfront:feature`, it says:

> "Spec is ready. Next step: run `/upfront:plan` to break this into implementation phases."

After `/upfront:plan`, it says:

> "Plan is ready. Next step: run `/upfront:build` to start building."

You don't need to memorize the workflow. Each skill hands off to the next one.

## That might be all you need

For most work, the path is:

```
/upfront:feature → /upfront:plan → /upfront:build
```

Three commands. The AI guides each transition.

## When you need more

The system will tell you. If your feature is too big to review, `/upfront:feature` will say:

> "This is ambitious — I count 4 distinct concerns. I'd recommend `/upfront:vision` to capture the full ambition and break it into reviewable increments."

If your problem is vague, it'll suggest `/upfront:ideate`. If there's no architecture doc, it'll suggest `/upfront:explore`. You don't need to know which skill to use — describe what you want and the system routes you.

Or just type `/upfront:up` and it figures out where you are and what you need.

## Quick reference

| What you're doing | Type this |
|-------------------|-----------|
| Build a feature | `/upfront:feature` |
| Small fix or bug | `/upfront:quick` |
| Not sure where to start | `/upfront:up` |
