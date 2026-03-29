---
title: /explore
description: Document the codebase and its ecosystem.
---

Five-phase investigation that produces `specs/ARCHITECTURE.md`.

## When to use

- Starting on a new codebase
- Architecture doc is missing or stale
- Before `/plan` on a brownfield project

## What it does

1. **Repo contents** — What's in the codebase, what's the tech stack
2. **Internal architecture** — Components, boundaries, patterns, data models
3. **External connections** — APIs, databases, queues, third-party services (with latency, SLAs, failure modes)
4. **Ecosystem** — Upstream consumers, downstream dependencies, blast radius
5. **Operational context** — Deployment, monitoring, alerting, on-call, who owns what

Each phase asks questions — it doesn't just read code and present findings. Challenge-first, decorate-second.

Greenfield projects get a fast exit: detects empty repos, creates minimal scaffolding, sends you to `/feature`.

## Output

`specs/ARCHITECTURE.md` — the shared reference document that every other command reads.

## Next step

`/feature` to define a feature, or `/plan` if you already have a spec.
