---
description: Use when the user says "check for updates", "is everything up to date", "upgrade", or at the start of a project to verify tooling is current. Checks plugin version, project guardrails, and instruction file health.
user-invocable: true
---

# Upgrade

You are checking that Upfront and the project's tooling are current and healthy. Run through each check, report status, and offer to fix anything that's out of date or missing.

## Step 1: Plugin version

Check the latest Upfront release:

```bash
curl -sf --max-time 5 https://api.github.com/repos/ThinkUpfront/Upfront/releases/latest | grep tag_name
```

Compare with the installed version (currently `0.3.7`). Report:
- Current: `0.3.7`
- Latest: `v[X.Y.Z]`
- Status: UP TO DATE or UPDATE AVAILABLE

If an update is available, tell the user:
```
Run this in your terminal to update:
  claude plugin marketplace update thinkupfront
Then restart Claude Code.
```

## Step 2: Project guardrails

Check what's installed on this machine and in this project:

| Tool | Check | Purpose |
|------|-------|---------|
| sloppy-joe | `which sloppy-joe` | Supply chain protection — blocks hallucinated package names |
| gitleaks | `which gitleaks` | Secret detection — catches API keys in commits |
| Test runner | Detect from project (go test, npm test, pytest, etc.) | Verify it runs cleanly |
| Linter | Detect from project (golangci-lint, eslint, ruff, etc.) | Verify it's configured |
| Formatter | Detect from project (gofmt, prettier, black, etc.) | Verify it's configured |
| Pre-commit hooks | Check `.pre-commit-config.yaml` or `.husky/` or `.git/hooks/` | Verify hooks are installed |

For each tool, report:
- INSTALLED or MISSING
- If missing, offer to install: "Want me to install sloppy-joe? `brew install brennhill/tap/sloppy-joe`"

Run the test runner and linter to verify they pass. If they fail, report the failure — don't fix it, just flag it.

## Step 3: Instruction file health

Check for `CLAUDE.md` or `AGENTS.md` in the project root.

**If missing:** "No AI instruction file found. Want me to create one? It makes every AI interaction in this repo better." If they confirm, immediately launch `/upfront:enlighten`.

**If exists:** Quick audit — does it have:
- [ ] Build/test/lint commands near the top?
- [ ] Three-tier boundaries (always / ask first / never)?
- [ ] Stack-specific code examples?
- [ ] References to specs/ files (`@specs/ARCHITECTURE.md`)?

Report how many of the 4 checks pass. If less than 3: "Your instruction file is thin. Want me to improve it?" If they confirm, immediately launch `/upfront:enlighten`.

## Step 4: Architecture health

Check `specs/ARCHITECTURE.md`:
- **Missing:** "No architecture doc. Want me to create one?" If they confirm, immediately launch `/upfront:explore`.
- **Exists:** Check the "Last reviewed" date. If stale (>30 days with commits since), flag it: "Architecture doc is [N] days old with [N] commits since. Consider running `/upfront:architect` to review."

## Step 5: Debt balance

Check if `specs/DEBT.md` exists. If it does, count open items by severity:
- Severity 3 (structural risk): **flag prominently**
- Severity 2 (needs fixing before production): count
- Severity 1 (cosmetic): count

Include in the summary below.

## Step 6: Summary

Present a clean report:

```
Upfront Health Check
====================

Plugin:        0.3.7 ✓ (latest)
sloppy-joe:    installed ✓
gitleaks:      MISSING — brew install gitleaks
Test runner:   npm test ✓ (passing)
Linter:        eslint ✓ (configured)
Formatter:     prettier ✓ (configured)
Pre-commit:    MISSING — recommend husky + lint-staged
CLAUDE.md:     exists, 2/4 checks pass — run /upfront:enlighten
ARCHITECTURE:  exists, reviewed 12 days ago ✓
Tech debt:     3 open items (1 severity 3, 2 severity 2)

Action needed:
  1. Install gitleaks: brew install gitleaks
  2. Set up pre-commit hooks
  3. Improve CLAUDE.md: /upfront:enlighten
```

Ask: "Want me to fix any of these?"

If they say yes, handle what you can directly (install tools, set up hooks). For things that need their own workflow (enlighten, explore, architect), tell them which command to run.

## Rules

- **Check, don't assume.** Run the actual commands to verify. Don't guess from file existence.
- **Don't install without asking.** Present what's missing, offer to install, wait for confirmation.
- **Be fast.** This is a health check, not a deep audit. Keep it under 2 minutes.
- **Project-specific.** Only recommend tools relevant to the detected tech stack. Don't suggest eslint for a Go project.
