---
title: Install
description: Get Upfront running in your project.
---

## Claude Code plugin (recommended)

```
/plugin marketplace add brennhill/Upfront
/plugin install upfront@upfront
```

This gives you all slash commands and the PostToolUse audit hook. For the audit binary:

```bash
brew install brennhill/tap/upfront
```

## Audit binary only

**Homebrew (macOS / Linux):**
```bash
brew install brennhill/tap/upfront
```

**Quick install (macOS / Linux):**
```bash
curl -fsSL https://raw.githubusercontent.com/brennhill/upfront/main/install-upfront.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/brennhill/upfront/main/install-upfront.ps1 | iex
```

**From source:**
```bash
git clone https://github.com/brennhill/upfront.git
cd upfront
go build -o upfront ./cmd/upfront/
./install.sh
```

## Manual install (no plugin system)

Copy commands into any project:

```bash
git clone https://github.com/brennhill/upfront.git
cp -r upfront/.claude/commands/ your-project/.claude/commands/
```

Or install globally:

```bash
cp -r upfront/.claude/commands/ ~/.claude/commands/
```

The commands are markdown files — no dependencies, no build step, no API keys.

## What you get

| File | Purpose |
|------|---------|
| `.claude/commands/*.md` | 16 slash commands |
| `specs/ARCHITECTURE.md` | Created by `/explore`, read by everything |
| `specs/DECISIONS.md` | Created by `/feature`, append-only |
| `specs/LEARNINGS.md` | Created by `/build`, grows over time |
| `specs/CONSTITUTION.md` | Optional — project-level invariants |
| `.upfront/audit.jsonl` | Local audit queue (if binary installed) |
