---
title: /upfront:enlighten
description: Audit and improve your CLAUDE.md or AGENTS.md with stack-specific examples.
---

## What it does

`/upfront:enlighten` audits your project's AI instruction files (`CLAUDE.md`, `AGENTS.md`) and improves them with stack-specific examples and the three-tier boundary pattern (always / ask first / never).

## When to use it

- First time setting up AI in a repo
- After choosing a tech stack — enlighten generates conventions and examples specific to your stack
- When AI keeps making the same mistakes — add a rule to prevent it
- Periodically as the project evolves

This skill is suggested automatically when `/upfront:explore` or `/upfront:feature` detects missing or weak instruction files.

## How it works

1. **Find instruction files** — checks for CLAUDE.md, AGENTS.md, .cursorrules, etc.
2. **Detect tech stack** — from package.json, go.mod, pyproject.toml, etc.
3. **Audit six areas** — commands, project structure, code style, testing, git workflow, boundaries
4. **Present gaps** — scored as Present / Missing / Weak
5. **Generate improvements** — stack-specific examples showing good and bad patterns
6. **Write to file** — adds to existing file, never overwrites

## The six areas

Research from 2,500+ repositories shows these are what separate instruction files that work from those that don't:

1. **Commands** — exact build/test/lint commands with flags, near the top
2. **Project structure** — where code goes, one sentence per directory
3. **Code style** — with good AND bad examples in your language
4. **Testing** — framework, patterns, what to mock and what not to
5. **Git workflow** — commit format, branch naming
6. **Boundaries** — always do / ask first / never do

## Key principle

One real code snippet showing your style beats three paragraphs describing it. Every suggestion includes a concrete, stack-specific example.
