---
title: Install
description: Get Upfront running in your project.
---

## Install

```bash
brew install brennhill/tap/upfront
upfront update
```

That's it. `brew install` gets the CLI. `upfront update` installs the Claude Code plugin (skills, hooks, agents). Restart Claude Code afterward.

## Update

```bash
brew upgrade upfront
upfront update
```

`upfront update` clones the latest plugin from GitHub, nukes the old cache, and installs fresh. Run it whenever you want the latest skills.

## Verify

In Claude Code, type `/upfront:` and you should see all skills in autocomplete. Try:

```
/upfront:start
```

## What you get

| Skill | Purpose |
|-------|---------|
| `/upfront:start` | Entry point — reads project state, routes you to the right skill |
| `/upfront:vision` | Capture the hypothesis, get to the first experiment fast |
| `/upfront:spike` | Build the minimum to test an idea — the scientific method |
| `/upfront:increment` | Learn step — retro after an experiment, decide what's next |
| `/upfront:feature` | Solidify a validated idea into a production-ready spec |
| `/upfront:plan` | Break a spec into ~400 LOC implementation phases |
| `/upfront:build` | Execute phases with strict TDD and review |
| `/upfront:ship` | Create a PR with spec-derived context |
| `/upfront:retro` | Check spec predictions against reality |
| `/upfront:assess` | Interactive problem-solving for a specific concern or tradeoff |
| `/upfront:ideate` | Divergent brainstorming when you don't know what to build yet |
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
| `/upfront:upgrade` | Check plugin version, guardrails, and instruction file health |

## Telemetry

Upfront sends anonymous usage events to help prioritize development: plugin version, skill name, and a hashed project identifier (derived from your git remote URL). No personally identifiable information is collected — no IP addresses, repo names, file paths, or code. The telemetry implementation is in [`plugin/hooks/hooks.json`](https://github.com/ThinkUpfront/Upfront/blob/main/plugin/hooks/hooks.json) — it's open source, you can verify exactly what's sent.

To disable, set the environment variable:

```bash
export DO_NOT_TRACK=1
```

## Uninstall

```bash
brew uninstall upfront
```

To also remove the Claude Code plugin:

```bash
rm -rf ~/.claude/plugins/cache/thinkupfront
```

Then edit `~/.claude/settings.json` and remove `"upfront@thinkupfront": true` from `enabledPlugins`.
