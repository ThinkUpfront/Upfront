# CLAUDE.md

## Project

Upfront — audit trail for AI-assisted feature definition workflows.

Go binary that hooks into Claude Code's PostToolUse events, captures thinking record summaries from `/feature` skill invocations, and maintains a durable local JSONL queue with optional remote flush to observability tools (Langfuse, Arize Phoenix, etc.).

## Architecture

- `cmd/upfront/main.go` — CLI entrypoint (hook, flush, log, purge, status)
- `internal/format/` — Event struct aligned to Delivery-Gap-Toolkit agent-monitoring trace format
- `internal/hook/` — Parse Claude Code hook stdin, extract thinking records from tool_response
- `internal/queue/` — Local JSONL queue with O_APPEND writes, rename-and-swap flush
- `internal/remote/` — POST events to configurable endpoint

## Build & Test

```bash
go build ./...
go test ./...
go test ./... -race    # always run with race detector
go build -o upfront ./cmd/upfront/
sloppy-joe check      # dependency supply chain check (typosquatting, hallucinated packages)
```

## Specs

All specs and plans are in `specs/`. Read these before making architectural decisions:
- `specs/upfront-audit-hooks.md` — full feature spec with thinking records
- `specs/upfront-audit-hooks-plan.md` — phased implementation plan
- `specs/upfront-audit-hooks-conversation.md` — full /feature conversation transcript

## Conventions

- Go stdlib only — zero external dependencies. Enforced by `sloppy-joe check` with canonical config at `~/.sloppy-joe/config.json`
- Event format must match the trace schema in the Delivery-Gap-Toolkit agent-monitoring guide
- Local queue file: `.upfront/audit.jsonl`
- Config file: `.upfront/config.json` (project-level), `~/.upfront/config.json` (user-level fallback)
- 90-day TTL on local logs

## Commit Style

- Descriptive commit messages with what changed and why
- No Co-Authored-By lines
