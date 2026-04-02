---
title: /upfront:upgrade
description: Check plugin version, project guardrails, and instruction file health.
---

## What it does

`/upfront:upgrade` runs a health check on everything: plugin version, installed guardrails, instruction file quality, and architecture doc freshness. Reports what's current, what's missing, and offers to fix it.

## When to use it

- Start of a new project — verify tooling is set up
- Periodically — make sure nothing has drifted
- After an update notification — check what changed

## What it checks

| Check | What it does |
|-------|-------------|
| Plugin version | Compares installed version against latest GitHub release |
| sloppy-joe | Supply chain protection — blocks hallucinated package names |
| gitleaks | Secret detection — catches API keys in commits |
| Test runner | Detects and verifies the project's test runner works |
| Linter / formatter | Checks they're configured and running |
| Pre-commit hooks | Verifies hooks are installed |
| CLAUDE.md / AGENTS.md | Quick audit — commands, boundaries, examples, references |
| ARCHITECTURE.md | Checks freshness — flags if stale |

## Output

A clean summary with status indicators and specific actions for anything that needs fixing. Offers to install missing tools directly.
