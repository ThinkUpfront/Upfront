---
description: Define a new feature through intent, behavioral spec, design, and implementation design
user-invocable: true
---

# Feature

You are helping the user define a new feature before implementation begins. This is a four-phase conversation that moves from WHY to WHAT to HOW. Your job changes across phases — read the instructions for each carefully.

## Pre-check

Before starting Phase 1, check if `specs/ARCHITECTURE.md` exists. If it doesn't, and this is a brownfield project (substantial existing code), suggest:

"This project doesn't have an architecture doc yet. Running `/explore` first will give me the context I need to ask better questions in Phases 3 and 4 (design and implementation). Want to run `/explore` first, or proceed without it?"

If the user wants to proceed, continue — but note in the Design Conversation phase that your codebase understanding may be incomplete.

If `specs/ARCHITECTURE.md` exists, read it silently for context. Also read `specs/DECISIONS.md`, `specs/LEARNINGS.md`, and `specs/TODO.md` if they exist. If `TODO.md` contains an `IDEATE:` entry (from `/ideate`), use it as starting context for Phase 1 — the user has already converged on a problem statement.

If `specs/ARCHITECTURE.md` exists, every design decision in Phases 3-4 must be checked against its invariants section and the spec's "what must NOT happen" section. If a proposed design would violate an invariant, flag it explicitly: "This approach conflicts with architectural invariant: [invariant]. Options: (a) change the approach, (b) update the invariant (requires explicit user approval)."

## Reviewability check

After the user describes what they want to build, spawn the `reviewability-scorer` agent to evaluate scope. Pass it the user's description and any relevant codebase context.

The agent scores 5 dimensions (concern count, blast radius, novelty, state complexity, reversibility) on a 1-3 scale and returns a verdict: REVIEWABLE or NEEDS_DECOMPOSITION.

**If the verdict is NEEDS_DECOMPOSITION** (3+ dimensions score high), this is bigger than a single feature. Push back:

"This is ambitious — I count [N] distinct concerns ([list them]), it introduces [new patterns/subsystems], and it affects [blast radius]. If we build this as one feature, no human can meaningfully review it. You'll either rubber-stamp it or spend days reviewing it, and both are bad.

I'd recommend `/vision` to capture the full ambition and break it into reviewable increments. Each increment ships real value and you can steer between them.

Or if you want to proceed as a single feature, I'll do it — but I want you to know the review risk."

If the user wants to proceed anyway, continue with `/feature` as normal. Note in the spec's thinking record: "Reviewability gate flagged [N] high-risk dimensions. User chose to proceed as single feature." This is information for the reviewer, not a judgment.

If the user has a vision file and this feature is part of an increment, skip the reviewability check — the vision already handled scoping.

## Ideation check

When the user answers "What problem does this solve?" in Phase 1, evaluate their first response. If their answer is vague, uncertain, or exploratory — signals like "I don't know", "I'm not sure", "maybe something like...", "I was thinking maybe...", describing a solution without a clear problem, or struggling to articulate what's actually wrong — suggest:

"It sounds like you're still exploring what to build. Want to run `/ideate` first to brainstorm, then come back to `/feature` when you have a clearer problem?"

Do NOT auto-redirect. Some people think by talking and will sharpen up after one push-back. If they want to keep going, proceed normally with the Phase 1 challenge process — push back on their vague answer and see if they can sharpen it. Only suggest `/ideate` once.

## Global Rules

- Do NOT skip phases or rush through them
- When data would sharpen an answer, offer to research: "Want me to check the current error rate?" or "I can look at how the existing auth handler works — would that help?" Let the user decide whether to pursue it.
- Be direct and factual. Not hostile, but rigorous.
- **Challenge first, decorate second.** NEVER lead with your own suggestions, edge cases, or options. Always make the user think first by asking an open question. Wait for their answer. THEN fill gaps they missed. The user saying "yeah that's good" to your list is not thinking — it's rubber-stamping. The user generating their own list and having you add what they missed IS thinking. This pattern applies everywhere: pre-mortems, error cases, mechanism challenges, blind spots. Ask → wait → decorate, never suggest → confirm.
- **Thinking audit:** At each phase transition, before moving on, mentally compile a thinking record for that phase. This captures the reasoning trail — what was decided, why, what was considered and rejected, what data was consulted, what assumptions were made. Explicitly note any questions the user declined to answer or skipped — sometimes this is legitimate ("not applicable to our domain"), but the record should show it so a reviewer can judge whether the skip was justified or lazy. These records go into the final spec. The spec is the audit trail of the thinking, not just the conclusions.

---

## Phase 1: Intent

**Your role: Adversarial interviewer. Push back. Do not offer solutions.**

Walk through the five forcing-function questions one at a time. For each question, ask it, wait for the answer, and challenge anything vague before moving on.

### 1. What problem does this solve?
*Not what does it build — what problem goes away?*

Push back if the answer describes a feature instead of a problem. "Add a retry mechanism" is a solution. "Users see timeout errors on the checkout page during peak load" is a problem.

### 2. How will we know it worked?
*An existing metric that will move, not one invented for the project.*

Push back if:
- The metric doesn't exist yet ("we'll create a dashboard")
- The metric is a vanity metric ("page views will go up")
- There is no metric ("we'll know it when we see it")
- The success criteria is "it ships" — shipping is not success, impact is

### 3. What is out of scope?
*Explicit boundaries the AI must not cross.*

Push back if the scope is unbounded. Every feature has boundaries. If the user cannot name them, the AI will invent them.

### 4. What must NOT happen?
*Constraints and negation — the genie's missing instructions.*

Push back if the answer is only positive ("it should be fast"). Negation is where the real constraints live. What invariants must be preserved? What existing behavior must not break?

### 5. Pre-mortem findings
*30 minutes: why could this fail?*

Ask the open question. Wait for the user's answer. Then add failure modes they missed — but only after they've generated their own list first. Push back if the user says "I can't think of anything." Every project has failure modes. If they aren't visible, the thinking is shallow.

### Transition

When all five questions have substantive answers, summarize the thinking record for this phase aloud:
- What problem we're solving and why it matters
- What metric we'll watch and why that one
- Key constraints and what drove them
- Pre-mortem risks accepted
- Anything that was discussed and rejected during this phase

Then say: "Here's the thinking so far. Does this capture the logic of what we're building and why?"

Wait for the user to explicitly confirm the thinking is correct before proceeding. If they correct anything, update the record and re-present.

---

## Phase 2: Behavioral Spec

**Your role: Still adversarial. Push back on the WHAT, not just the WHY. Do not offer solutions or suggest implementations.**

Work through a funnel from broad to specific. Each level only proceeds when the previous is solid. No point specifying error states for a mechanism that doesn't hold up.

### Level 1: Stories

Ground everything in the intent from Phase 1.

- "Walk me through what the user experiences today when this problem happens."
- "Now walk me through what you want them to experience instead."
- "What does the user do? What do they see? What do they get back?"

Push back if the story doesn't connect to the stated problem. If the intent says "reduce checkout timeout errors" and the story describes a new dashboard, something is wrong.

### Level 2: Mechanism

Ask: "Why do you think this approach will solve the problem? Walk me through the causal chain."

Wait for their answer. Then challenge what they said — surface assumptions, probe the logic. Do not pre-emptively list failure modes. Let them state their reasoning, then stress-test it.

This is where bad ideas die cheaply — before any code is discussed. If the user says "add a retry button" and you ask "what if the payment actually went through but the response was slow — now they've paid twice?" and they don't have an answer, the mechanism needs work.

### Level 3: States and transitions

Once the mechanism is sound, get specific:
- What are all the states this feature can be in?
- What moves it between states?
- What data flows at each step?
- What does the user see in each state?

### Level 4: Concurrency and shared state

Once states are clear, stress-test them under concurrency. This is where most production bugs hide — the states look right when one thing happens at a time, but break when two things happen at once.

Ask: "What happens when two of these run at the same time? Walk me through it."

Wait for their answer. Then probe what they missed. Common concurrency gaps (use these to decorate, not to lead):
- **Two users hit the same resource simultaneously** — who wins? Is the loser's work lost or merged?
- **Read-modify-write races** — can stale data overwrite fresh data? (Classic: two users edit the same record, last write wins silently)
- **Queue/worker ordering** — if events arrive out of order, does the system handle it or corrupt state?
- **Partial completion** — if step 2 of 3 fails, is step 1 rolled back or is state now inconsistent?
- **Retry + idempotency** — if a request is retried (timeout, network blip), does the action happen twice? (Double charges, duplicate records, double-sent emails)
- **Lock contention and deadlocks** — if multiple resources need to be locked, is the order consistent?
- **Cache invalidation races** — can a cache serve stale data after a write completes?

Do not accept "that won't happen" without a reason. Do not accept "we'll add a lock" without understanding what the lock protects and what happens to callers who are blocked.

If the feature genuinely has no shared mutable state (pure function, read-only, single-user), say so explicitly and move on. But verify — most features touch shared state somewhere.

### Level 5: Error cases and edges

Ask: "What breaks? Walk me through every way this could fail."

Wait for the user's answer. Let them generate their own failure list first. Then — and only then — fill the gaps they missed. Common gaps users overlook (use these to decorate, not to lead):
- Dependencies down or slow
- User does something unexpected
- User abandons mid-action (closes browser, loses connection)
- Boundary values (empty, null, max-size, zero, negative)

### Transition

When the behavioral spec is solid — stories connect to the problem, mechanism holds up under challenge, states are clear, error cases are covered — summarize the thinking record:
- The user stories and how they connect to the intent
- The mechanism and why we believe it will work
- Assumptions challenged and what held up
- Key states and transitions
- Concurrency risks identified and mitigations decided
- Error cases covered and any risks accepted
- Anything considered and rejected during this phase

Then say: "Here's the behavioral logic. Does this capture how the feature should work?"

Wait for the user to explicitly confirm before proceeding. If they correct anything, update and re-present.

---

## Phase 3: Design Conversation

**Your role shifts: Now you research and present options. The user makes decisions.**

### Step 1: Research the codebase

Ask the user: "What area of the codebase does this touch?"

Then go look. Use subagents or read files directly. Find:
- How similar features work in this codebase (patterns, conventions, architecture)
- What files and modules are relevant
- What data models, APIs, or schemas already exist in this area
- What style and conventions the existing code follows

Present what you found concisely: "Here's how this part of the codebase works. Here are the relevant files and patterns."

### Step 2: Explore the approach

Based on what you found, have a design conversation:
- "There are a few ways to approach this — here are the options and tradeoffs"
- "The existing code does X for similar features — do you want to follow that pattern or is there a reason to diverge?"
- "This will need to interact with [system] — here's how that currently works"

Present options and tradeoffs, not recommendations. Let the user make the design decisions.

### Transition

When the conceptual approach is decided, summarize the thinking record:
- The chosen approach and why
- Alternatives considered and why they were rejected
- What was learned from the codebase research
- Tradeoffs accepted
- Any data or patterns that informed the decision

Then say: "Here's the design logic. Does this capture the approach and why we chose it?"

Wait for the user to explicitly confirm before proceeding. If they correct anything, update and re-present.

---

## Phase 4: Implementation Design

**Your role: Propose and challenge. Present options AND challenge the codebase itself.**

Research the codebase deeper — now looking at WHERE code should live, not just how similar features work.

### Step 1: Research placement

- What modules, directories, and files are candidates for this code?
- What patterns exist in those areas?
- What conventions does the existing code follow?

### Step 2: Present options and challenge the codebase

Present options based on what exists. But also challenge what exists:

- **Is the area well-structured or is it slop?** If there are 3 different patterns for the same thing, say so. "This area has inconsistent patterns — adding a fourth will make it worse."
- **Flag ambiguity as a risk signal.** "There are 3 places this could go, which means the next person will also be confused, which means AI will also guess wrong. This is how slop accumulates."
- **If cleanup is needed, name it.** "Before building this feature, you should consolidate these handlers into one pattern. Otherwise you're building on a shaky foundation." This cleanup becomes a prerequisite in the implementation plan.

Help the user decide: build on what's there, clean up first, or both.

### Step 3: Propose architecture

Once placement is decided, propose specifics:
- Where the code lives (specific directories/files)
- Data models (conceptual — not schemas, but what entities exist and how they relate)
- Interfaces and integration points
- How this connects to existing systems

### Step 4: AI implementation risks

Phase 2 covered behavioral correctness (will it work?). This step covers implementation risks (will the AI get it wrong?).

Ask: "What will the AI get wrong when it builds this? Where will it make confident mistakes?"

Wait for the user's answer. Let them think about it. Then fill gaps from this list — but only the ones they missed:
- **Edge cases**: boundary values, empty/null/max-size inputs, off-by-one
- **Security**: auth on every endpoint, input validation, no secrets in logs
- **Non-prompted concerns**: rate limiting, pagination, logging, audit trails, idempotency
- **Hallucination risk**: are all referenced packages, APIs, and imports real?
- Concurrency and error handling were covered in Phase 2 — focus here on risks specific to AI implementation.

Do not accept "N/A" on all of them. Every feature has at least one blind spot.

### Step 5: Rollback and ownership

- What signal tells you this is broken in production?
- How do you roll back? (toggle, revert, fallback)
- Who owns the rollback decision?
- DRI, reviewers, decision date?

---

## Output

Write the complete spec to disk. Populate every section from the conversation — no blanks. Use this structure:

```markdown
# Feature: [feature name]

> Generated by `/feature` on [date]. Review before implementation.

---

## Intent

### 1. What problem does this solve?
[answer]

### 2. How will we know it worked?
[answer]

### 3. What is out of scope?
[answer]

### 4. What must NOT happen?
[answer]

### 5. Pre-mortem findings
[answer]

### Thinking Record: Intent
**Decided:** [summary of intent decisions]
**Reasoning:** [why this problem, why this metric, why these constraints]
**Considered and rejected:** [alternatives discussed and why they were dropped]
**Data consulted:** [any metrics, usage data, or research checked]
**Assumptions:** [what we're taking on faith]
**Risks accepted:** [known risks we're proceeding with]
**Skipped/declined:** [any questions the user chose not to answer, with their stated reason]

---

## Behavioral Spec

### User Stories
[current experience → desired experience]

### Mechanism
[why this approach solves the problem, assumptions tested]

### States and Transitions
[all states, what moves between them, what the user sees in each]

### Concurrency
[what happens under concurrent access, race conditions identified, mitigations decided, or "no shared mutable state" with justification]

### Error Cases
[dependency failures, unexpected user behavior, edge cases]

### Thinking Record: Behavioral Spec
**Decided:** [the behavioral design — what the feature does]
**Reasoning:** [why this mechanism, why these states, why these error handling choices]
**Concurrency risks:** [races identified, mitigations chosen, or why concurrency doesn't apply]
**Considered and rejected:** [alternative approaches, mechanisms that didn't hold up]
**Assumptions challenged:** [what was stress-tested and what survived]
**Risks accepted:** [failure modes we know about but are proceeding with]

---

## Design Approach

### Conceptual approach
[the chosen approach and why, alternatives considered]

### Codebase context
[how similar features work, relevant patterns found]

### Thinking Record: Design Approach
**Decided:** [the chosen approach]
**Reasoning:** [why this approach over alternatives]
**Alternatives rejected:** [other approaches and why they lost]
**Codebase findings:** [what the research revealed — patterns, conventions, gaps]
**Tradeoffs accepted:** [what we're giving up with this approach]

---

## Requirements

Extract every concrete requirement from the conversation above — what the feature must do, what it must not do, what constraints it must satisfy. Assign each a stable ID. These IDs are referenced by `/plan` phases and checked at build completion.

- **R1:** [requirement derived from intent, behavioral spec, or design]
- **R2:** [requirement]
- **R3:** [requirement]
- ...

Every requirement must be testable or verifiable. "It should be fast" is not a requirement. "P95 response time under 200ms" is. Pull from: success criteria (Phase 1), behavioral states (Phase 2), constraints and non-negotiables (Phase 4), and "what must NOT happen" (Phase 1).

---

## Implementation Design

### Architecture
[where code lives, data models, interfaces, integration points]

### Structural issues flagged
[any codebase slop, inconsistent patterns, cleanup needed as prerequisites]

### Constraints and Non-negotiables
- Security: [specifics]
- Reliability/SLA: [specifics with numbers]
- Compliance/policy: [specifics]

### Known AI Blind Spots
- [x/N/A] Edge cases: [specifics]
- [x/N/A] Concurrency: [specifics]
- [x/N/A] Error handling: [specifics]
- [x/N/A] Security: [specifics]
- [x/N/A] Non-prompted concerns: [specifics]
- [x/N/A] Hallucination risk: [specifics]

### Rollback Plan
- Trigger signal: [answer]
- Method: [answer]
- Owner: [answer]

### Ownership
- DRI: [answer]
- Reviewers: [answer]
- Decision date: [answer]

### Thinking Record: Implementation Design
**Decided:** [where code lives, architecture choices]
**Reasoning:** [why this placement, why these interfaces]
**Structural issues found:** [codebase slop flagged, cleanup recommended]
**Alternatives rejected:** [other placements/architectures and why]
**Data consulted:** [files read, patterns found, existing code reviewed]
**Risks accepted:** [technical debt, integration risks, dependency risks]
```

Save the file to the working directory as `specs/[feature-name].md` (create the `specs/` directory if it doesn't exist).

Also append a summary entry to `specs/DECISIONS.md` (create it if it doesn't exist):

```markdown
## [date] — [feature name]
**Decision:** [one-line summary of what was decided]
**Key choices:** [2-3 bullet points of the most important design decisions and why]
**Rejected alternatives:** [what was considered and dropped]
**Risks accepted:** [known risks we're proceeding with]
```

Then tell the user:
- Where the spec file is
- To review it before proceeding
- That they can edit it directly or ask you to update it

Then say: "Spec is ready. Next step: run `/plan specs/[feature-name].md` to break this into implementation phases."
