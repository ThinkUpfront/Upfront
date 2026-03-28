# Handoff

> Paused: 2026-03-28

## What was running

**Command:** Freeform — building the Upfront command suite
**Phase/step:** Completed all 10 commands, discussing DS/DE/MLE variants

## Completed

- Ran `/feature` on upfront audit hooks — full 4-phase spec with thinking records
- Created implementation plan (6 phases) for Go binary
- Migrated `~/dev/upfront` from old Python project to new Go repo (old → `~/dev/upfront-archive`)
- Built and installed 10 commands globally (`~/.claude/commands/`):
  - `/ideate` — divergent brainstorming before `/feature`
  - `/explore` — 5-phase codebase + ecosystem documentation
  - `/feature` — 4-phase intent definition with concurrency analysis, ideation check, ARCHITECTURE.md pre-check
  - `/refine` — spec iteration with `//` markup + conversational challenge
  - `/plan` — 3-level architecture deep-dive, staleness check, guardrail audit, Phase 0
  - `/build` — TDD enforcement, per-phase review (optional stronger model), red team, crash recovery, learning capture, pre-flight tooling audit
  - `/quick` — small changes without ceremony (bails >50 lines)
  - `/note` — zero-friction todo capture
  - `/pause` / `/resume` — session continuity
- All commands synced across 3 locations: `~/.claude/commands/`, `~/dev/upfront/.claude/commands/`, `Delivery-Gap-Toolkit/.claude/commands/`
- Researched GSD, Superpowers, Compound Engineering, Context Engineering Kit — incorporated best patterns
- Created review.html with Mermaid workflow diagrams
- Persistent docs system: ARCHITECTURE.md, DECISIONS.md, LEARNINGS.md, TODO.md, HANDOFF.md
- Commits: multiple commits pushed to both `brennhill/upfront` and `Delivery-Gap-Toolkit`

## In progress

Discussion of DS/DE/MLE variants. User has a team copying data pipelines/tables with no reuse or docs. `/feature-ds` exists in the toolkit but was designed before the "force thinking" philosophy — needs rethink.

## Next action

Either:
1. Rethink `/feature-ds` with the same challenge-first philosophy for data work (pipelines, tables, models) — address structural rot, not just new features
2. Or start building the Go binary (`/build specs/upfront-audit-hooks-plan.md`) in `~/dev/upfront`

## Decisions made

- Commands are user-level global (`~/.claude/commands/`), not per-project — available everywhere
- Architecture docs persist across features in `specs/ARCHITECTURE.md`, reviewed every `/plan` run with staleness check (>30 days + commits = forced re-review)
- `/build` review agent can optionally use a stronger model than the implementing agent
- Red team is adversarial — fixes obvious issues silently, asks about judgment calls, flags design concerns
- Pre-flight tooling audit checks for missing linters/security scanners per ecosystem and pushes hard for installation
- `/refine` uses `//` comments as agenda items but still challenges every change — not a rubber-stamp editor
- `/feature` suggests `/ideate` for vague problem statements, `/explore` for missing ARCHITECTURE.md
- Crash recovery: stash-and-ask on uncommitted changes, reconcile progress file with git log
- DS/DE/MLE needs different commands — TDD doesn't apply, verification is experimental, "done" is ambiguous, data quality is the foundation

## Gotchas

- Agents writing to files across repos sometimes hit permission issues on the toolkit submodule — use `cp` from global to fix
- The old `brennhill/upfront` GitHub repo had different history — had to force push to replace it
- `/feature-ds` and `/challenge-ds` exist in the toolkit (`.claude/commands/`) but are NOT installed globally — they predate the current philosophy
- The Delivery-Gap-Toolkit is a git submodule under `ai-augmented-dev` — committing requires `git add Delivery-Gap-Toolkit` from the parent repo

## Git state

- **Branch:** main (both repos)
- **Uncommitted changes:** Website dist files (unrelated), some research data files
- **Recent commits:**
  - 4857eb4b Update Delivery-Gap-Toolkit submodule
  - 6ae0e865 Update Delivery-Gap-Toolkit submodule
  - 253eed18 Update Delivery-Gap-Toolkit submodule

## Active files

- `~/.claude/commands/*.md` — all 10 command files
- `~/dev/upfront/.claude/commands/*.md` — mirrored copies
- `~/dev/upfront/specs/upfront-audit-hooks.md` — feature spec for Go binary
- `~/dev/upfront/specs/upfront-audit-hooks-plan.md` — implementation plan (6 phases)
- `~/dev/upfront/specs/review.html` — Mermaid workflow diagram page

## User notes

- User has a DS/DE team with copied pipelines, no reuse, no docs — wants the same "force thinking" approach for data work
- This has evolved from "upfront audit hooks" into a full build system — name may change
- `/feature-ds` needs complete rethink, not just an update
