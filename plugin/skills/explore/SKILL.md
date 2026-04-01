---
description: Use when the user says "document this codebase", "what does this project do", "set up for AI", or needs to understand an unfamiliar codebase. Documents the codebase and its operational context.
user-invocable: true
---

# Explore

You are documenting a codebase so that future `/upfront:feature`, `/upfront:plan`, and `/upfront:build` sessions have full context. This is a one-time deep investigation that produces `specs/ARCHITECTURE.md` — the shared reference document for all future work.

If `specs/ARCHITECTURE.md` already exists, read it first. You are updating it, not starting from scratch. Present what's changed since the last review.

## Philosophy

The AI can read code. What it cannot do is understand the system in context — why the architecture is the way it is, what external systems it depends on, what the performance constraints are, who owns what, and what's about to change. This exploration fills those gaps through a structured conversation.

**Challenge first, decorate second.** Read the codebase, form your understanding, present it, and ask the user to correct you. Do not ask the user to explain things you can figure out from the code.

---

## Greenfield Check

Before diving in, determine whether there's anything to explore.

If the repo is empty or near-empty (only scaffolding, go.mod, package.json, README, and no substantial source code):

Say: "This looks like a greenfield project — there's no existing architecture to document yet. I'll create a minimal `specs/ARCHITECTURE.md` with what we know (language, toolchain, directory structure) and set up `specs/DECISIONS.md` and `specs/LEARNINGS.md` for future use."

Create the minimal files, then say: "Ready to define your first feature? The architecture will build up naturally as you make decisions." If they confirm, immediately launch `/upfront:feature` — don't tell them to type it.

Do not force the user through 5 phases of exploration on an empty repo. Exit early.

---

## Phase 1: What's in the repo

**Your role: Investigate, then present for correction.**

Read the codebase thoroughly by spawning **parallel Haiku sub-agents** for each area of investigation. Use `model: "haiku"` when spawning these agents. The sub-agents do the heavy searching (Grep, Glob, Read) on the cheap model and return compressed summaries. You (on the main model) synthesize their findings into the architecture doc.

Spawn explorers in parallel for:
1. **Stack & tooling** — language, framework, go.mod/package.json/pyproject.toml, build/test/lint/format commands, CI config
2. **Structure & boundaries** — directory layout, module boundaries, what each package owns
3. **Dependencies & patterns** — external deps and versions, error handling, logging, config, naming conventions
4. **Test infrastructure** — framework, coverage, patterns, test/source file conventions

Each explorer should return a structured summary under 200 lines — not raw search results. Instruct them: "Return findings as a concise summary for an AI reader. No raw file listings. Focus on patterns, conventions, and surprises."

Present what you found concisely:

"Here's what I see in this codebase: [structured summary]. Does this match reality, or am I missing something?"

Specifically ask:
- "Are any of these directories dead/abandoned code I should ignore?"
- "Are any of these dependencies scheduled to be replaced?"
- "Is anything here about to change significantly that I should note?"

Wait for confirmation before proceeding.

---

## Phase 2: Internal architecture

**Your role: Present your understanding of how the pieces fit together. Challenge your own assumptions.**

Based on what you read, present:

### Module boundaries and responsibilities
- What does each major module/package own?
- Where are the boundaries between them?
- Are the boundaries clean or leaky?

### Data models and persistence
- What are the core entities?
- How are they stored? (database, file system, in-memory, external service)
- What are the relationships and lifecycle states?
- What constraints exist at the data level? (unique keys, foreign keys, soft deletes, etc.)

### Internal communication patterns
- How do modules talk to each other? (function calls, interfaces, channels, events)
- Is there a consistent pattern or does every module do it differently?
- What are the invariants? What must ALWAYS be true regardless of what changes?

### Error handling and propagation
- How do errors flow through the system?
- Is there a consistent error type/wrapping pattern?
- What happens when things fail? Retry? Propagate? Swallow?

Present all of this, then ask:

"Here's how I think the internals work. What am I getting wrong? What invariants am I missing?"

Wait for confirmation before proceeding.

---

## Phase 3: External connections

**Your role: Ask about everything the system touches outside this repo. The AI cannot discover this from code alone.**

Start with: "What does this system connect to? Walk me through every external dependency — services it calls, services that call it, databases, queues, file systems, third-party APIs, anything outside this repo."

Wait for the user's answer. Then probe what they missed:

### For each external connection, ask:

**Identity and ownership:**
- What is it? (service name, provider, URL/endpoint)
- Who owns it? (your team, another team, third party, managed service)
- Is there a contract? (API spec, schema, SLA)

**Data flow:**
- What data goes in? What comes back?
- What format? (JSON, protobuf, CSV, binary)
- What volume? (requests per second, messages per day, bytes per batch)
- Is it synchronous or asynchronous?

**Performance characteristics:**
- What's the expected latency? (p50, p95, p99 if known)
- What's the timeout configured? What happens when it's exceeded?
- Is there a rate limit? What happens when you hit it?
- Is there backpressure? How does the system handle it?

**Failure modes:**
- What happens when this dependency is down?
- Is there a circuit breaker, fallback, or retry policy?
- What's the blast radius? (feature degrades, page errors, data loss)
- How long can you tolerate it being down?

**Service tier:**
- What's the SLA? (99.9%, 99.99%, best-effort)
- Is there an on-call? Who gets paged?
- Is this in the critical path or is it fire-and-forget?

Do not accept "it's just a REST API" without the details above. Every external connection is a failure mode. Push for specifics.

### Also ask about:

- **Databases:** What databases does this system use? What's shared vs owned? Read replicas? Connection pooling? Migration strategy?
- **Message queues / event streams:** Kafka, SQS, RabbitMQ, etc. What topics/queues? What's the ordering guarantee? What happens with poison messages?
- **Caches:** Redis, Memcached, in-process? What's cached? What's the TTL? What happens on cache miss? What's the invalidation strategy?
- **File storage:** S3, local disk, NFS? What's stored? What's the retention policy?
- **Authentication / authorization:** How does this system authenticate to external services? How do external services authenticate to it? Where are credentials stored?

---

## Phase 4: Ecosystem context

**Your role: Understand where this system sits in the broader picture.**

Ask: "Is there an ecosystem diagram, service map, or architecture diagram that shows how this system fits into the larger infrastructure? If so, point me to it so I can reference it."

If one exists, read it and incorporate the context. If not, ask the user to describe:

- **Upstream:** What systems feed data or requests into this one?
- **Downstream:** What systems depend on this one's output?
- **Peers:** What systems operate at the same level and coordinate with this one?
- **Shared infrastructure:** What platform services does this use? (Service mesh, API gateway, config service, feature flags, observability stack)

Then ask:
- "If this system goes down, what breaks? What's the blast radius beyond this repo?"
- "If this system is slow, who notices first? Is there monitoring/alerting?"
- "Are there any upcoming changes to the ecosystem that will affect this system? (Migrations, deprecations, new dependencies)"

---

## Phase 5: Operational context

**Your role: Capture the human and process context that doesn't live in code.**

Ask about:

### Deployment
- How does this deploy? (CI/CD pipeline, manual, GitOps)
- How often? (continuous, daily, weekly, quarterly)
- How do you roll back? (revert commit, feature flag, blue/green)
- Are there deployment windows or freezes?

### Observability
- What monitoring exists? (metrics, logs, traces)
- What dashboards do people actually look at?
- What alerts fire? What's the on-call rotation?
- Where do you go first when something breaks?

### Team and ownership
- Who owns this codebase? (team, individual, shared)
- Who are the key people who understand it best?
- What's the review process? (PR review, pair programming, async)
- Are there areas of the codebase that nobody understands anymore?

### Known debt and pain points
- What's the biggest pain point in this codebase right now?
- What technical debt is actively causing problems?
- What would you fix first if you had a free week?
- What's the scariest part of the code? (The thing that makes people nervous to touch)

---

## Output

After all five phases, write or update `specs/ARCHITECTURE.md`:

```markdown
# Architecture

> Last explored: [date]
> Last reviewed during /upfront:plan: [date if applicable]

## Codebase Overview
[language, framework, toolchain, directory structure, build/test/lint commands]

## Internal Architecture

### Module Map
[modules, boundaries, responsibilities]

### Data Models
[core entities, storage, relationships, constraints]

### Communication Patterns
[how modules talk, consistency of patterns]

### Invariants
[what must always be true — the rules that cannot be violated]

### Error Handling
[how errors flow, patterns used]

## External Connections

### [Service/Dependency Name]
- **Type:** [REST API / gRPC / database / queue / cache / etc.]
- **Owner:** [your team / other team / third party]
- **Data flow:** [what goes in, what comes back, format, volume]
- **Latency:** [p50/p95/p99 if known]
- **Timeout:** [configured value, what happens on exceed]
- **Failure mode:** [what happens when it's down, blast radius]
- **SLA:** [service tier, who gets paged]

[repeat for each connection]

## Ecosystem Context
[upstream, downstream, peers, shared infrastructure, blast radius, ecosystem diagram reference]

## Operational Context

### Deployment
[how, how often, rollback strategy]

### Observability
[monitoring, dashboards, alerts, first-response playbook]

### Team
[ownership, key people, review process, knowledge gaps]

## Known Debt and Pain Points
[biggest pain points, active tech debt, scariest code, what to fix first]

## Risk Areas
[weakest integration points, fragile code, inconsistent patterns, areas where AI will get it wrong]
```

Also create or update `specs/DECISIONS.md` if it doesn't exist:

```markdown
# Decisions

> Append-only register of architectural and design decisions.
> Read by /upfront:feature and /upfront:plan for context. Do not edit past entries.

[no entries yet — future /upfront:feature sessions will append here]
```

And create or update `specs/LEARNINGS.md` if it doesn't exist:

```markdown
# Learnings

> Append-only register of what surprised us, what the agent got wrong, what patterns emerged.
> Read by /upfront:plan and /upfront:build for context. Do not edit past entries.

[no entries yet — future /upfront:build sessions will append here]
```

Then tell the user:
- Where the files are
- To review `specs/ARCHITECTURE.md` before any `/upfront:feature` or `/upfront:plan` work
- That future `/upfront:plan` runs will review and update this document
- That `/upfront:feature` and `/upfront:build` will append to `DECISIONS.md` and `LEARNINGS.md` over time

Also check if `CLAUDE.md` or `AGENTS.md` exists and is well-structured (has build/test commands, boundaries, stack-specific examples). If it's missing or thin, suggest: "Your instruction file could use some work. Want me to audit it and add stack-specific examples?" If they confirm, immediately launch `/upfront:enlighten` — don't tell them to type it.

---

## Rules

- Do NOT ask the user to explain things you can figure out from the code. Read first, present, ask for corrections.
- Do NOT accept vague answers about external connections. "It calls an API" is not an answer. Push for latency, failure modes, volume.
- Do NOT skip Phase 3 (external connections). This is where most production incidents originate and where the AI has the least visibility.
- If the user says "I don't know" about a connection's performance characteristics, note it as a risk: "Unknown latency/failure behavior for [service] — this is a blind spot."
- If the user points to an ecosystem diagram, reference its location in ARCHITECTURE.md so future sessions can find it.
