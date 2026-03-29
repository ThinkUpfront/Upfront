---
description: Break a feature spec into ~400 LOC implementation phases by researching the codebase
---

# Plan

You are breaking an approved feature spec into implementation phases. Each phase should be small enough to review (~400 lines), independently verifiable, and independently committable.

This requires understanding the actual codebase — you cannot plan in the abstract.

## Input

The user provides a path to a spec file (e.g., `specs/feature-name.md`). Read it fully.

If no path is provided, check the `specs/` directory for recent specs and ask which one to plan.

## Process

### 1. Read the spec and architecture

Read the entire spec file. Understand:
- The intent (what problem, what metric, what's out of scope, what must not happen)
- The constraint surface (scope boundaries, acceptance criteria)
- The blind spots flagged

Also read `specs/ARCHITECTURE.md` if it exists — its invariants section contains hard constraints that every proposed phase must comply with. If a phase would require violating an architectural invariant, flag it before proposing the plan.

### 2. Research the codebase

Before proposing any phases, investigate the actual code. Use subagents for parallel research where possible.

For each file, module, or component referenced in the spec (especially in the Design Approach and Implementation Design sections):
- Read the referenced files
- Understand the existing patterns, conventions, and architecture

Then answer:
- What files will need to change?
- What is the dependency order between those changes? (e.g., schema first, then store methods, then API, then UI)
- How big is each change likely to be?
- Are there existing tests to extend or will new test files be needed?
- Are there integration points where changes in one file affect another?

### 3. Architectural deep-dive

#### Staleness check

Before starting the deep-dive, check if `specs/ARCHITECTURE.md` exists.

**If it exists:**
1. Read it fully.
2. Find the "Last reviewed" (or "Last explored") date in the file.
3. Run `git log --oneline --since="[that date]" | wc -l` to count commits since that date.
4. Calculate how many days have elapsed since that date.

**If more than 30 days old AND there have been commits since:**
- Present: "ARCHITECTURE.md was last reviewed [date]. There have been [N] commits since then. This document may be stale — using it as-is risks building on wrong assumptions."
- Do NOT ask the user "is this still right?" — that invites "yes" rubber-stamping.
- Instead, actively read the codebase and COMPARE what you find to what ARCHITECTURE.md says. Walk through each section of the document and verify it against the actual code. Present specific drifts: "ARCHITECTURE.md says X, but I found Y."
- Strongly push for a full re-review: "I need to verify this is still accurate before planning. Let me walk through the architecture to check for drift."
- If major drift is found, update ARCHITECTURE.md with corrections before proceeding to phases. Present the changes to the user for confirmation.

**If less than 30 days old OR no commits since:**
- Note: "Architecture reviewed [date], [N] commits since — appears current."
- Skip the full 3-level deep-dive below. Instead, do a quick delta check: read the areas of the codebase this feature will touch and verify they match what ARCHITECTURE.md says. If they match, proceed directly to "Persisting the architecture." If there is drift in the feature-relevant areas, update those sections and proceed.

**If `specs/ARCHITECTURE.md` doesn't exist:**
- Note that it doesn't exist and proceed — the deep-dive below will create it.

---

Before proposing phases, walk through the architecture at three levels. Each level must be confirmed by the user before proceeding to the next. Wrong assumptions at any level will poison every phase.

**This is a conversation, not a presentation.** Ask questions, challenge answers, push for specifics. Do not accept vague descriptions — if the user says "it's a standard REST API," ask what "standard" means in this codebase. Do they use middleware? How is auth handled? What's the error response shape? Is there a router or is it handler-per-file?

#### Level 1: System architecture

Present your understanding of the overall system. Then ask the user to confirm or correct.

- **What are the major components?** Services, databases, queues, caches, external APIs. Draw the boundaries.
- **How do they communicate?** HTTP, gRPC, message bus, shared database, file system, stdin/stdout. What are the protocols and contracts?
- **What are the invariants?** What must ALWAYS be true about this system regardless of what changes? (e.g., "every write must be durable before acknowledging," "auth must be checked on every request," "no data crosses tenant boundaries")
- **What are the deployment boundaries?** What ships together, what ships independently? Does this feature cross a deployment boundary?

Challenge assumptions: "You said these two services communicate over HTTP — is that synchronous? What happens when the downstream service is slow? Is there a timeout? A circuit breaker? Or does the caller just hang?"

Ask: "Here's how I understand the system. Does this match reality?"

Wait for confirmation before proceeding.

#### Level 2: Subsystem architecture

Zoom into the specific subsystems this feature touches. Read the actual code — do not guess.

- **What modules, packages, or layers exist in this area?** What does each one own? Where are the boundaries between them?
- **What patterns does the existing code follow?** Repository pattern? Service layer? Handlers calling business logic calling data access? Or is it a ball of mud? Be honest about what you find.
- **What are the data models?** Not just schemas — what are the entities, what are their relationships, what are the lifecycle states? What constraints exist at the data level (unique keys, foreign keys, check constraints, soft deletes)?
- **What interfaces exist between subsystems?** Are they clean (well-defined types, clear contracts) or leaky (passing raw maps, stringly-typed, implicit assumptions)?
- **What are the subsystem invariants?** What must be true within this subsystem? (e.g., "an event's sequence number is always monotonically increasing," "a queue entry is either pending or flushed, never both," "config is loaded once at startup and never mutated")

Challenge what you find: "This module has three different error handling patterns — which one is intentional? Should we consolidate before building on top, or is there a reason they diverge?"

Ask: "Here's how the subsystems work. Does this match what you intended, or has it drifted?"

Wait for confirmation before proceeding.

#### Level 3: Design patterns and connections

Now get specific about how the pieces connect and what patterns to use.

- **What design patterns does the codebase use?** Not textbook names — actual patterns. "Handlers validate input, call a service function, and return a response. Service functions are pure logic with injected dependencies. Data access is behind interfaces." Or: "Everything is in one file and there are no patterns." Both are useful answers.
- **What patterns should this feature follow?** Match existing conventions unless there's a reason to diverge. If you recommend diverging, say why explicitly.
- **What are the connection points?** For each place this feature touches existing code, what is the exact interface? Function signatures, types, error types, return values. Not abstract — concrete.
- **What are the edge behaviors at each connection?** What happens at each boundary when inputs are empty, nil, maximum size, malformed? What happens when the connection fails? What does the caller see?
- **Concurrency model:** What is shared mutable state? What synchronization exists? What ordering guarantees are needed? If the spec identified concurrency concerns, how does the architecture address each one specifically?
- **Error propagation:** How do errors flow through the system? Are they wrapped? Typed? Do they carry context? Is there a consistent pattern or does every module do it differently?

Challenge for well-understood behaviors: "You're building a queue with flush semantics — this is a solved problem. Are you using write-ahead log semantics? Exactly-once delivery? At-least-once? What's the durability guarantee? Have you looked at how existing tools (WAL, NATS, etc.) handle this, or are you inventing from scratch?"

Push for specifics: "You said the config is loaded from a JSON file. What happens if the file doesn't exist? What if it's malformed JSON? What if it has valid JSON but missing required fields? What are the defaults? Where are they documented?"

Ask: "Here's how I think the pieces connect. Are the patterns right? Am I missing any invariants?"

Wait for confirmation before proceeding.

#### Persisting the architecture

Architecture is project knowledge, not per-feature knowledge. After all three levels are confirmed, write or update `specs/ARCHITECTURE.md`:

Use the same structure as `/explore` produces (see `specs/ARCHITECTURE.md` template in explore.md): Codebase Overview, Internal Architecture (with Module Map, Data Models, Communication Patterns, Invariants, Error Handling subsections), External Connections, Ecosystem Context, Operational Context, Known Debt and Pain Points, Risk Areas. Update the "Last reviewed" date.

**If `specs/ARCHITECTURE.md` already exists:** Read it first. Review it against the current codebase — does it still match reality? Update anything that has drifted. Add new subsystems or patterns this feature revealed. Update the "Last reviewed" date. Present changes to the user for confirmation.

**If it doesn't exist:** Create it from the conversation. This becomes the shared reference for every future `/plan` and `/feature` run.

The plan file references the architecture doc rather than duplicating it.

### 4. Audit architectural guardrails

Before proposing any implementation phases, check whether the project has tooling that enforces the architecture discussed in step 3. The architecture review is worthless if nothing prevents the code from violating it.

**Check for each category:**

- **Type safety:** Is there a type checker / compiler configured and running? (Go compiler, TypeScript strict mode, mypy, etc.) If not, the AI will produce code with type errors that won't be caught until runtime.
- **Linting:** Is there a linter configured? Does it enforce the patterns from the architecture review? (import boundaries, naming conventions, forbidden patterns) If the architecture says "no cross-module imports" but nothing enforces it, it will be violated by Phase 2.
- **Test infrastructure:** Does the test runner work? Are there existing tests that pass? If `go test ./...` or `npm test` fails before you start, you cannot do TDD.
- **Formatting:** Is there an auto-formatter? If not, every phase will introduce formatting inconsistencies that pollute diffs.
- **Build:** Does the project build cleanly right now? If it's already broken, you're building on a broken foundation.

**For each gap found, present it:**

"The architecture review says [X pattern], but there's no tooling that enforces it. Options:
1. Add [specific tool/config] as Phase 0 before the feature work
2. Accept the risk and rely on code review to catch violations
3. This doesn't matter for this feature — here's why"

Let the user decide. If they choose option 1, add a Phase 0 to the plan that installs and configures the guardrail. Phase 0 must pass before any feature code is written.

**What to look for specifically:**

- Go: `go vet`, `staticcheck` or `golangci-lint`, `go test` passing, `gofmt`/`goimports`
- TypeScript: `tsc --strict`, ESLint with import rules, `npm test` passing, Prettier
- Python: `mypy` or `pyright`, `ruff` or `flake8`, `pytest` passing, `black` or `ruff format`
- Any language: CI pipeline configured, pre-commit hooks, build scripts, [`sloppy-joe`](https://github.com/brennhill/sloppy-joe) for dependency supply chain checks (typosquatting, hallucinated packages, canonical name enforcement)

Do not silently assume guardrails exist. Check. If the project is brand new with no tooling, say so — "This project has zero guardrails. I recommend adding [X, Y, Z] as Phase 0."

### 5. Propose phases

Break the work into phases. Each phase should be:
- **~400 lines or less** — based on your estimate of actual change size, not a guess
- **Independently verifiable** — has its own automated and manual success criteria
- **Independently committable** — the codebase compiles, tests pass, and nothing is broken after this phase alone
- **Ordered by dependency** — schema before logic, logic before API, API before UI

If Phase 0 (guardrails) was agreed upon, it comes first. No feature code until the tooling is in place.

For each phase, define:
1. **What changes** — specific files and what happens to them
2. **Automated verification** — exact commands to run (e.g., `make test`, `npm run typecheck`, specific test files)
3. **Manual verification** — what a human checks before the next phase begins

#### Human-writes mode (`/plan --human`)

If the user runs `/plan --human`, or if the spec contains `[human-writes]` markers, identify which phases should be implemented by the human rather than the AI.

**Suggest human-writes for phases that contain:**
- Concurrency logic (locks, channels, atomics, transaction isolation)
- Security-critical code (auth, authorization, input validation, crypto)
- Core business logic (the mechanism from the spec — the thing that makes the feature work)
- Invariant enforcement (the code that guarantees the rules from ARCHITECTURE.md)

Present your suggestions: "I recommend these phases as human-writes: Phase 3 (queue concurrency) and Phase 5 (auth middleware). These are where understanding the code is more valuable than generating it. The rest can be AI-implemented. Agree?"

The user confirms which phases are human-writes. Mark them in the plan file:

```markdown
### Phase 3: Queue concurrency handler [human-writes]
**Files:** internal/queue/queue.go
**Changes:** Implement concurrent append with O_APPEND, flush with rename-and-swap
**Mode:** Human implements, AI writes tests and reviews
...
```

`/build` reads these markers and switches to human-writes mode for flagged phases.

### 6. Present for review

Present the plan to the user. For each phase, explain:
- What it does and why it's in this order
- Estimated size (lines changed)
- What depends on it

Ask: "Does this phasing make sense? Should I adjust the order, split any phase further, or combine small ones?"

Iterate based on feedback.

### 7. Write to disk

Once approved, write the plan as a separate file at `specs/[feature-name]-plan.md`:

```markdown
# Plan: [feature name]

> Generated by `/plan` on [date]. Review before implementation.
> Spec: `specs/[feature-name].md`

## Architecture

See `specs/ARCHITECTURE.md` (reviewed [date]).

**Feature-specific architectural notes:**
[anything specific to this feature that doesn't belong in the shared architecture doc — e.g., "this feature adds a new subsystem (upfront/internal/queue/) not yet in ARCHITECTURE.md" or "N/A — no new architectural decisions"]

---

### Phase 1: [descriptive name]
**Files:** [list of files that change]
**Changes:** [what happens in this phase]
**Estimated size:** [lines]

**Automated verification:**
- [ ] [exact command]
- [ ] [exact command]

**Manual verification:**
- [ ] [what to check]

---

### Phase 2: [descriptive name]
...
```

Then tell the user:
- Where the plan file is
- To review it before implementation
- When ready, run `/build specs/[feature-name]-plan.md` to implement phase by phase

## Rules

- Do NOT propose phases without researching the codebase first. Abstract plans lead to wrong phasing.
- Do NOT make phases too granular — a phase that changes 3 lines is not worth the overhead. Combine small related changes.
- Do NOT assume file structure — verify it by reading the codebase.
- Respect the spec's scope boundaries — do not plan changes to files listed as non-goals.
- Each phase's verification criteria should trace back to the spec's acceptance criteria and blind spots. The plan is how you deliver the spec, not a separate document.
- If the spec's Implementation Design section flags structural issues (slop, inconsistent patterns, cleanup needed), include refactoring as the first phase(s) — prerequisite work before building the feature on top. Do not build on a shaky foundation.
- If the spec is missing information you need to plan (e.g., no implementation design, unclear architecture), say what's missing and ask the user to update the spec before proceeding.
