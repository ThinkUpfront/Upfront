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

Restart Claude Code. All 22 `/upfront:*` skills will be available in every project.

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
| `/upfront:assess` | Interactive problem-solving for a specific concern or tradeoff |
| `/upfront:ideate` | Divergent brainstorming before `/upfront:feature` |
| `/upfront:quick` | Small changes without full ceremony |
| `/upfront:patch` | Bug fixes from issues or problem statements |
| `/upfront:explore` | Document codebase and operational context |
| `/upfront:enlighten` | Audit and improve CLAUDE.md/AGENTS.md with stack-specific examples |
| `/upfront:refine` | Iterate on a spec with challenge |
| `/upfront:debug` | Scientific method debugging with persistent state |
| `/upfront:teach` | Walk through codebase with comprehension checks |
| `/upfront:note` | Zero-friction idea/todo capture |
| `/upfront:pause` | Structured handoff for next session |
| `/upfront:resume` | Restore context from a pause |
| `/upfront:architect` | Architecture review with evolution plan |
| `/upfront:re-architect` | Execute an architecture evolution phase by phase |
| `/upfront:up` | Smart router — figures out what you need |

## Telemetry

Upfront sends anonymous usage events to help prioritize development: plugin version, skill name, and a hashed project identifier (derived from your git remote URL). No personally identifiable information is collected — no IP addresses, repo names, file paths, or code. The telemetry implementation is in [`plugin/hooks/hooks.json`](https://github.com/ThinkUpfront/Upfront/blob/main/plugin/hooks/hooks.json) — it's open source, you can verify exactly what's sent.

To disable, set the environment variable:

```bash
export DO_NOT_TRACK=1
```

## Uninstall

Edit `~/.claude/settings.json` and remove `"upfront@thinkupfront": true` from `enabledPlugins`. Optionally remove the `"thinkupfront"` entry from `extraKnownMarketplaces`.

Or use the clean script from the repo:

```bash
git clone https://github.com/ThinkUpfront/Upfront.git
bash Upfront/scripts/clean-plugin.sh
```

