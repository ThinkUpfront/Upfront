# Upfront

Audit trail for AI-assisted feature definition. Upfront hooks into Claude Code's PostToolUse events, captures thinking record summaries from `/feature` skill invocations, and maintains a durable local JSONL queue with optional remote flush to observability tools like Langfuse, Arize Phoenix, Helicone, or Portkey.

A complete `/feature` run produces 4 phase events (Intent, Behavioral Spec, Design Approach, Implementation Design). Partial runs produce fewer — the count tells you where they stopped. This gives managers three metrics without reading every spec: adoption (did they run `/feature`?), depth (did they finish all 4 phases?), and effectiveness (do spec'd features have less rework?).

## Install

Requirements: Go 1.21+, `jq`

```bash
git clone https://github.com/brennhill/upfront.git
cd upfront
./install.sh
```

The install script:
1. Builds the `upfront` binary
2. Copies it to `~/.local/bin/` (override with `INSTALL_DIR=/your/path ./install.sh`)
3. Registers the PostToolUse hook in `~/.claude/settings.json` (backs up existing settings first)

The script is idempotent — running it again will not duplicate the hook.

### Manual install

```bash
go build -o upfront ./cmd/upfront/
cp upfront ~/.local/bin/
```

Then add the hook to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Skill",
        "hooks": [{ "type": "command", "command": "upfront hook" }]
      }
    ]
  }
}
```

## Configuration

Create `.upfront/config.json` in your project root, or `~/.upfront/config.json` for user-level defaults. Project-level config takes precedence.

```json
{
  "endpoint": "https://your-langfuse-instance.com/api/public/ingestion",
  "auth_header": "Bearer pk-lf-...",
  "ttl_days": 90,
  "project_name": "my-project"
}
```

| Field | Description | Default |
|-------|-------------|---------|
| `endpoint` | URL to POST event batches to. Leave empty for local-only mode. | (none) |
| `auth_header` | Value for the `Authorization` header. | (none) |
| `ttl_days` | Days to keep local events before `purge` deletes them. | 90 |
| `project_name` | Project identifier included for your own reference. | (none) |

If no config file exists, upfront runs in local-only mode — events are written to `.upfront/audit.jsonl` and never sent anywhere.

**Important:** If your config contains API keys or tokens, add it to `.gitignore`:

```
.upfront/config.json
```

The `.upfront/` directory is already in this repo's `.gitignore`.

## CLI Reference

### `upfront hook`

Process Claude Code PostToolUse hook input from stdin. This is called automatically by the hook — you don't run it manually.

Reads JSON from stdin, extracts thinking records from `/feature` skill responses, appends events to the local queue, and attempts a best-effort remote flush if configured.

### `upfront flush`

Manually flush queued events to the remote endpoint.

```bash
upfront flush
```

Uses rename-and-swap for atomic draining. If the remote send fails, events are re-queued. If no remote endpoint is configured, prints a message and exits.

### `upfront log`

Print audit events from the local queue.

```bash
upfront log                        # last 50 events
upfront log --limit 10             # last 10 events
upfront log --feature my-feature   # filter by feature name
upfront log --phase 1              # filter by phase number
upfront log --feature auth --phase 3  # combine filters
```

| Flag | Description | Default |
|------|-------------|---------|
| `--feature` | Filter events by feature name (exact match) | (all) |
| `--phase` | Filter events by phase number (1-4) | (all) |
| `--limit` | Maximum number of events to display | 50 |

### `upfront purge`

Delete events older than the configured TTL.

```bash
upfront purge
```

Uses `ttl_days` from config (default: 90 days). Rewrites the queue file atomically.

### `upfront status`

Show queue status and configuration.

```bash
upfront status
```

Displays: queue file path, event count, last event details, and remote endpoint (if configured).

## Compatible Remote Tools

Upfront's event format extends the agent-monitoring trace schema from the Delivery-Gap-Toolkit, which is compatible with:

- **[Langfuse](https://langfuse.com/)** — open-source LLM observability. POST events to the ingestion API.
- **[Arize Phoenix](https://phoenix.arize.com/)** — LLM tracing and evaluation. Accepts JSON trace spans.
- **[Helicone](https://helicone.ai/)** — LLM monitoring and analytics.
- **[Portkey](https://portkey.ai/)** — AI gateway with built-in observability.

Any tool that accepts JSON event arrays via HTTP POST can be used as a remote endpoint.

## Event Format

Each event is a JSON object written as one line in the JSONL queue:

```json
{
  "session_id": "claude-session-uuid",
  "timestamp": "2026-03-28T10:15:00Z",
  "agent_id": "upfront",
  "action_type": "upfront_phase_complete",
  "action_detail": "Phase 1: Intent",
  "target": "/path/to/project",
  "feature_name": "checkout-timeout-fix",
  "phase": 1,
  "phase_name": "Intent",
  "phases_total": 4,
  "thinking_summary": "Decided: Fix checkout timeout errors...",
  "skipped_questions": [],
  "duration_ms": null,
  "result": "success"
}
```

See `examples/sample-audit.jsonl` for complete examples.

## Uninstall

1. Remove the hook entry from `~/.claude/settings.json`
2. Delete the binary: `rm ~/.local/bin/upfront`
3. Optionally remove local data: `rm -rf .upfront/`
