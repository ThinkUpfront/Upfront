---
title: Install
description: Get Upfront running in your project.
---

## Quick install (recommended)

Clone the repo and ask Claude to install the skills:

```bash
git clone https://github.com/brennhill/upfront.git
```

Then in Claude Code, say:

> "Install the Upfront skills from the `upfront/` directory into my global commands (`~/.claude/commands/`)"

Claude will copy the markdown files and you'll have all slash commands available in every project.

## Manual install

Copy commands into a specific project:

```bash
git clone https://github.com/brennhill/upfront.git
cp -r upfront/.claude/commands/ your-project/.claude/commands/
```

Or install globally (available in every project):

```bash
cp -r upfront/.claude/commands/ ~/.claude/commands/
```

The commands are markdown files — no dependencies, no build step, no API keys.

## Audit binary (optional)

The audit binary captures thinking records from `/feature` runs into a JSONL queue. It's optional — the slash commands work without it.

**Homebrew (macOS / Linux):**
```bash
brew install brennhill/tap/upfront
```

**From source:**
```bash
cd upfront
go build -o upfront ./cmd/upfront/
./install.sh
```

## What you get

| File | Purpose |
|------|---------|
| `.claude/commands/*.md` | 16 slash commands |
| `specs/ARCHITECTURE.md` | Created by `/explore`, read by everything |
| `specs/DECISIONS.md` | Created by `/feature`, append-only |
| `specs/LEARNINGS.md` | Created by `/build`, grows over time |
| `specs/CONSTITUTION.md` | Optional — project-level invariants |
| `.upfront/audit.jsonl` | Local audit queue (if binary installed) |
