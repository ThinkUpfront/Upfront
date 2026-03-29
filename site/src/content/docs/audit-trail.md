---
title: Audit Trail
description: Durable JSONL log of thinking records for team visibility.
---

The `upfront` binary captures every thinking record produced during `/feature` runs and writes structured events to a local JSONL queue.

## How it works

A PostToolUse hook fires after every `/feature` skill invocation. The hook:

1. Parses the thinking record from the tool response
2. Writes a structured event to `.upfront/audit.jsonl`
3. Optionally flushes to a remote endpoint

Events are stored locally first (durable — survives network failures), then flushed to your observability stack.

## CLI commands

```bash
upfront status                            # Queue depth, last event, config
upfront log                               # View audit events (last 50)
upfront log --feature checkout --phase 1  # Filter by feature/phase
upfront flush                             # Push queued events to remote
upfront purge                             # Delete events older than TTL
```

## Remote integration

Configure `.upfront/config.json` or `~/.upfront/config.json`:

```json
{
  "endpoint": "https://your-langfuse.example.com/api/public/ingestion",
  "auth_header": "Bearer pk-lf-...",
  "ttl_days": 90,
  "project_name": "my-project"
}
```

Compatible with Langfuse, Arize Phoenix, Helicone, Portkey, or any endpoint that accepts JSON POST.

## What managers see

Three metrics without reading a single spec:

1. **Adoption** — percentage of features that went through `/feature`
2. **Depth** — how many phases were completed (all 4 = thorough, 1-2 = bailed early)
3. **Effectiveness** — rework rate difference between spec'd and unspec'd features

The audit trail is 30-second triage: decide whether to read the spec or send it back.
