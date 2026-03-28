# Progress: Upfront Audit Hooks

> Spec: `specs/upfront-audit-hooks.md`
> Plan: `specs/upfront-audit-hooks-plan.md`

## Completed Phases

## Phase 1: Go module + data model + event formatting — COMPLETE

**What changed:** `internal/format/event.go`, `internal/format/event_test.go`
**TDD cycles:** 3 (defaults, round-trip, field names), then 4 more after review fixes (error_message omit/include, duration_ms null presence, timestamp no-nanoseconds, constant)
**Review findings:** 3 MUST FIX resolved (timestamp changed from time.Time to string for clean ISO-8601, PhaseName capitalization fixed to match spec, ErrorMessage field added with omitempty). 1 SHOULD FIX resolved (PhasesTotal extracted to DefaultPhasesTotal constant). 2 SHOULD FIX accepted (input validation deferred to Phase 2 parse boundary, SkippedQuestions population is Phase 2 concern).
**Surprises:** None
**Learnings for future phases:**
- PhaseName should be passed capitalized (e.g., "Intent", "Behavioral Spec") — that's how the `### Thinking Record:` headings appear in `/feature` output
- ErrorMessage field exists with omitempty — Phase 2 should set it for parse errors
- Timestamp is a clean RFC3339 string (no nanoseconds), not a time.Time

## Phase 2: Hook stdin parser + thinking record extractor — COMPLETE

**What changed:** `internal/hook/handler.go`, `internal/hook/handler_test.go`
**TDD cycles:** 9 (malformed JSON, valid parse, non-Skill, empty response, no records, 1 phase, 4 phases, unknown phase, heading termination)
**Review findings:** 3 MUST FIX resolved (regex truncated multi-line summaries to first line — removed (?m) flag; tests didn't verify multi-line capture — added Reasoning/Risks assertions; FeatureName used SkillName instead of Args). 1 SHOULD FIX accepted (Target uses cwd as placeholder — spec file path not available in hook payload). parse_error path deferred to Phase 5 CLI layer.
**Surprises:** The (?m) regex flag caused $ to match end-of-line, making lazy [\s\S]*? stop at first line boundary. Critical bug caught by review.
**Learnings for future phases:**
- FeatureName comes from ToolInput.Args (the feature being defined), not SkillName
- Target is currently cwd — Phase 5 may want to derive spec file path from args
- parse_error events are not yet generated — Phase 5 CLI should handle this at the pipeline level
- Renamed HookInput to Input to avoid stutter (hook.Input vs hook.HookInput)

## Phase 3: Local JSONL queue with concurrent write safety — COMPLETE

**What changed:** `internal/queue/queue.go`, `internal/queue/queue_test.go`
**TDD cycles:** 11 (EnsureDir, Append JSONL, concurrent writes, ReadAll, ReadAll empty/missing, Flush, Flush missing/empty, Purge TTL, Purge missing, Flush crash recovery, corrupt line skip)
**Review findings:** 2 MUST FIX resolved (pre-existing .flushing file crash recovery — now reads prior events before rename; corrupt JSONL lines skip instead of bricking queue). 1 SHOULD FIX resolved (Purge temp file permissions 0o600). Filter API deferred to Phase 5 caller.
**Surprises:** 15 golangci-lint issues on first pass — Event struct is 224 bytes so gocritic flags pass-by-value. Changed Append to take *Event.
**Learnings for future phases:**
- Append takes *format.Event (pointer) due to struct size
- Corrupt JSONL lines are silently skipped — Phase 5 could surface warnings
- Flush recovers events from prior crashed flushes automatically
- Scanner buffer set to 1MB max line to handle large thinking summaries
- No filter API on queue — Phase 5 CLI should filter after ReadAll()

## Learnings

(see per-phase learnings above)
