---
title: /plan
description: Break a spec into ~400 LOC implementation phases.
---

Starts with a three-level architectural deep-dive, then breaks the spec into buildable phases.

## When to use

- You have an approved spec from `/feature`
- Before running `/build`

## What it does

### Architecture deep-dive

Three levels, each confirmed before proceeding:

1. **System** — components, communication, invariants, deployment boundaries
2. **Subsystem** — modules, data models, interfaces, subsystem invariants
3. **Patterns** — concrete interfaces, edge behaviors, concurrency model, error propagation

The AI challenges assumptions at each level. If `specs/ARCHITECTURE.md` is stale (>30 days with commits since), it refuses to use it at face value — compares it to the actual codebase and presents drifts.

### Guardrails audit

Checks for missing tools per ecosystem (linters, type checkers, security scanners, dead code detection, secret detection, slopsquatting protection). Proposes Phase 0 to install anything missing.

### Strategy before phasing

Formulates 2-3 implementation approaches varying by invasion depth, reversibility, testability, and incremental value. Presents as a comparison table. You choose.

### Phasing

Breaks into ~400 LOC phases, each independently verifiable and committable. Identifies clarity debt (code that confused the AI during research) and assigns fixes to the phases that touch that code.

### Human-writes mode (`/plan --human`)

Identifies which phases should be human-implemented: concurrency, security, core business logic, invariant enforcement.

## Output

`specs/[feature-name]-plan.md` — phases with files, changes, and verification criteria. Updates `specs/ARCHITECTURE.md`.

## Next step

`/build specs/[feature-name]-plan.md` to implement.
