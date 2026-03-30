---
title: Install
description: Get Upfront running in your project.
---

## Plugin install (recommended)

Install via the Claude Code plugin system. Run these in your terminal:

```bash
claude plugin marketplace add ThinkUpfront/Upfront
claude plugin install upfront
```

Restart Claude Code. All 18 `/upfront:*` skills will be available in every project.

## Verify

In Claude Code, type `/upfront:` and you should see all skills in autocomplete. Try:

```
/upfront:up
```

## What you get

| Skill | Purpose |
|-------|---------|
| `/upfront:vision` | Strategic clarity for big ambitions (Rumelt's kernel) |
| `/upfront:increment` | Structured retro between increments |
| `/upfront:feature` | Define a feature through intent, spec, and design |
| `/upfront:plan` | Break a spec into ~400 LOC implementation phases |
| `/upfront:build` | Execute phases with strict TDD and review |
| `/upfront:ship` | Create a PR with spec-derived context |
| `/upfront:retro` | Check spec predictions against reality |
| `/upfront:ideate` | Divergent brainstorming before `/feature` |
| `/upfront:quick` | Small changes without full ceremony |
| `/upfront:patch` | Bug fixes from issues or problem statements |
| `/upfront:explore` | Document codebase and operational context |
| `/upfront:refine` | Iterate on a spec with challenge |
| `/upfront:debug` | Scientific method debugging with persistent state |
| `/upfront:teach` | Walk through codebase with comprehension checks |
| `/upfront:note` | Zero-friction idea/todo capture |
| `/upfront:pause` | Structured handoff for next session |
| `/upfront:resume` | Restore context from a pause |
| `/upfront:up` | Smart router — figures out what you need |

## Uninstall

Edit `~/.claude/settings.json` and remove `"upfront@thinkupfront": true` from `enabledPlugins`. Optionally remove the `"thinkupfront"` entry from `extraKnownMarketplaces`.

Or use the clean script from the repo:

```bash
git clone https://github.com/ThinkUpfront/Upfront.git
bash Upfront/scripts/clean-plugin.sh
```

## Audit binary (optional)

The audit binary captures thinking records from skill runs into a JSONL queue. It's optional — the skills work without it.

**Homebrew (macOS / Linux):**
```bash
brew install ThinkUpfront/tap/upfront
```

**From source:**
```bash
git clone https://github.com/ThinkUpfront/Upfront.git
cd Upfront
go build -o upfront ./cmd/upfront/
./install.sh
```
