---
title: /feature
description: Define a feature through four phases of forced thinking.
---

The core command. Four phases, each with a different AI role.

## When to use

- You know the problem and need to define the solution
- Before any feature larger than ~300 lines
- When you need an auditable trail of design decisions

## The four phases

### Phase 1: Intent

*AI role: Adversarial interviewer*

Five forcing-function questions, asked one at a time:

1. **What problem does this solve?** Not what it builds — what problem goes away.
2. **How will we know it worked?** An existing metric that will move.
3. **What's out of scope?** Draw the boundary now, not during implementation.
4. **What must NOT happen?** Negative requirements are more durable than positive ones.
5. **Pre-mortem.** It's six months from now and this failed. What went wrong?

The AI pushes back on vague answers and refuses to move on until the thinking is substantive.

### Phase 2: Behavioral Spec

*AI role: Mechanism designer*

Five-level funnel:
1. User stories (who does what)
2. **Mechanism** — why will this approach actually work? Not how — why.
3. States and transitions
4. **Concurrency** — what happens when two of these run at the same time?
5. Error cases — what can go wrong and what happens when it does

### Phase 3: Design Conversation

*AI role: Research partner*

The AI researches the codebase, presents options and tradeoffs. You make the decisions. Reads `specs/ARCHITECTURE.md` if it exists.

### Phase 4: Implementation Design

*AI role: Codebase critic*

Proposes architecture AND challenges the codebase itself. Flags inconsistent patterns, structural rot, ambiguous placement. Cleanup becomes a prerequisite.

## Thinking records

Every phase transition produces a thinking record — what was decided, why, what was rejected, what was skipped. These go into the final spec as an audit trail.

## Output

`specs/[feature-name].md` — the spec with embedded thinking records. Also appends to `specs/DECISIONS.md`.

## Next step

`/plan` to break it into buildable phases.
