---
description: Architecture review that produces an evolution plan — god modules first, then layers, then patterns
user-invocable: true
---

# Architect

You are performing an architectural code review. Your job is to find structural problems and propose an evolution path — not fix bugs or add features.

**Input:** $ARGUMENTS (optional — a specific area to focus on, or blank for full codebase)

## Process

### 1. Load context

Read silently:
- `specs/ARCHITECTURE.md` (if it exists — check its "Last reviewed" date. If stale (>30 days with commits since), note it: "ARCHITECTURE.md may be stale — comparing against actual codebase, not the doc.")
- `specs/DECISIONS.md` (understand why things are the way they are before proposing changes)
- `specs/LEARNINGS.md` (past mistakes and patterns)

### 2. Structural survey

Spawn **parallel Haiku sub-agents** (use `model: "haiku"` when spawning) to map the codebase. Each returns a compressed summary.

**Agent 1: Module inventory**
- List every package/module/directory with its responsibility (one line each)
- Count: files, lines, exported symbols per module
- Flag modules over 500 LOC or with more than one clear responsibility

**Agent 2: Dependency graph**
- Which modules import which other modules?
- Are there circular dependencies?
- Which modules are imported by everything? (high fan-in = potential god module)
- Which modules import everything? (high fan-out = potential orchestrator that does too much)

**Agent 3: Pattern inventory**
- What patterns exist? (repository, service layer, handler, factory, etc.)
- Are they consistent or does every module do it differently?
- Where is business logic? Is it in the right layer?
- Where is I/O? Is it separated from logic?

**Agent 4: Interface health**
- Are module boundaries clean? (well-defined types, explicit contracts)
- Are they leaky? (passing raw maps, stringly-typed, implicit assumptions)
- Are there god interfaces? (interfaces with 10+ methods = probably multiple responsibilities)

### 3. Diagnosis (top-down, three levels)

Present findings in strict order. Each level must be addressed before the next — don't propose design patterns for code that lives in the wrong module.

#### Level 1: God modules

The highest-impact structural problem. A god module is a package that:
- Has multiple unrelated responsibilities
- Is imported by most of the codebase
- Changes for unrelated reasons
- Is where new code gets dumped because "it's already there"

For each god module found:
```
GOD MODULE: [package name] ([N] LOC, [M] importers)

Responsibilities detected:
  1. [responsibility A] — [files/functions]
  2. [responsibility B] — [files/functions]
  3. [responsibility C] — [files/functions]

Proposed split:
  [package-a/] ← responsibility A
  [package-b/] ← responsibility B
  [package-c/] ← responsibility C

Impact: [which importers change, estimated effort]
```

If no god modules exist, say so and move on.

#### Level 2: Layer violations

After module boundaries are clean, check layering. Standard layers (adapt to the project's conventions):
- **Handlers / Controllers** — accept input, validate, delegate
- **Services / Business logic** — pure logic, no I/O, no framework dependencies
- **Repositories / Data access** — persistence, external APIs, I/O
- **Models / Domain types** — data structures, no behavior beyond validation

Flag:
- Business logic in handlers ("fat controllers")
- I/O in service layers (database calls mixed with logic)
- Framework dependencies leaking into business logic
- Direct database access from handlers (skipping the service layer)

For each violation:
```
LAYER VIOLATION: [description]
  Location: [file:line]
  Problem: [what's in the wrong layer]
  Fix: Extract [what] into [which layer]
```

If layers are clean, say so and move on.

#### Level 3: Design pattern opportunities

Only after levels 1-2 are addressed. Look for:
- **Builder** — complex object construction scattered across call sites
- **Strategy** — switch/if-else chains on a type that could be polymorphic
- **Factory** — object creation logic duplicated across modules
- **Observer/Event** — tight coupling where events would decouple
- **Repository** — data access not behind an interface (untestable)
- **Adapter** — external API calls not wrapped (vendor lock-in)

Only flag patterns where the benefit is clear. Don't recommend patterns for their own sake. For each:
```
PATTERN OPPORTUNITY: [pattern name]
  Location: [file:line or module]
  Current: [what it does now]
  Proposed: [what the pattern would give you]
  Benefit: [testability / flexibility / clarity — be specific]
```

### 4. Prioritize

Rank all findings by impact. Present as a table:

```
| Priority | Type | Location | Description | Effort |
|----------|------|----------|-------------|--------|
| 1 | God module | pkg/core | Split into auth, queue, config | Large |
| 2 | Layer violation | handlers/api.go | DB calls in handler | Small |
| 3 | Pattern | services/payment | Strategy for payment providers | Medium |
```

### 5. Write the evolution doc

Write `specs/ARCHITECTURE-EVOLUTION.md`:

```markdown
# Architecture Evolution: [project name]

> Generated by `/architect` on [date]
> Current state: specs/ARCHITECTURE.md (reviewed [date])

## Summary

[2-3 sentences: what's the biggest structural problem and what does the evolution achieve]

## Current pain points

[bullet list of the top issues, in priority order]

## Evolution phases

### Phase 1: [god module split / highest priority item]
**What changes:** [specific files, moves, renames]
**Why first:** [this unblocks everything else / highest impact]
**Risk:** [what could break]
**Verification:** [how to confirm it worked — tests pass, no import cycles, etc.]

### Phase 2: [next priority]
...

## Constraints

- Each phase must leave the codebase compiling and tests passing
- No behavioral changes — this is pure restructuring
- Existing tests must continue to pass (they may need import path updates)
- New tests are NOT required per phase (structure changes, not behavior changes)

## What this does NOT cover

- New features (use /feature)
- Bug fixes (use /patch or /debug)
- Performance optimization (separate concern)
```

### 6. Record the decision

Append to `specs/DECISIONS.md`:

```markdown
## [date] — Architecture review
**Decision:** [summary of what the evolution plan proposes]
**Rationale:** [top findings — god modules, layer violations, etc.]
**Alternative:** Leave as-is (accepted risk: [what continues to degrade])
```

### 7. Report

Tell the user:
- Summary of findings (god modules, layer violations, pattern opportunities)
- The evolution doc location
- "Run `/re-architect` to execute this evolution phase by phase."
- If architecture is clean: "No structural issues found. Marking as reviewed." — and update the "Last reviewed" date in `specs/ARCHITECTURE.md` so the health check nudge in `/build` knows.

## Rules

- Do NOT propose changes that alter behavior. This is restructuring only.
- Do NOT recommend patterns where the current code is simple and clear. Patterns are for managing complexity that exists, not complexity you're introducing.
- Do NOT propose changes to code you haven't read. Read first, diagnose second.
- If the architecture is already clean, say so. Not every codebase needs restructuring.
- God modules first, always. Fixing patterns inside a god module is wasted effort — it's getting split anyway.
