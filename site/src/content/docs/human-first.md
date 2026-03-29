---
title: Human-First Development
description: How Upfront keeps developers strong instead of replacing them.
---

AI writes code. Humans understand systems. When the AI writes all the code, humans stop understanding the systems they're responsible for.

## The problem with full delegation

When an AI generates code and a human approves it, two things happen:

1. The code ships faster.
2. The human's understanding of the system gets shallower.

This works until it doesn't. The system breaks in production, and the person on-call doesn't understand the code they approved.

The research is clear: developers who delegate code generation score 40% lower on comprehension than those who use AI for conceptual questions (Anthropic, 2026). AI-generated code gets reviewed less carefully despite having comparable vulnerability rates (Microsoft Research, 2024).

Speed is not the bottleneck. Understanding is.

## Six mechanisms

### 1. Challenge before suggestion

Every command follows the same pattern: **ask, wait, then decorate.**

The AI does not present options and wait for approval. It asks open questions, waits for the human's answer, and then fills the gaps the human missed. The human thinks first. The AI adds what they didn't think of.

"Do you approve this design?" produces a yes. "What happens when two of these run at the same time? Walk me through it." produces understanding.

### 2. Human-writes mode

The most direct intervention. When a phase contains critical logic — concurrency, security, core business rules — the AI writes the tests and the human writes the implementation.

1. The AI writes failing tests that describe what the code should do.
2. The AI generates `TODO(human)` stubs with constraints and common mistakes to avoid.
3. The human fills in the stubs.
4. The AI runs the tests and reviews what the human wrote.
5. If tests fail, the AI shows which tests fail and why — but does not fix it for them.

The commit is tagged `[human-writes]` so git history distinguishes human-authored critical code from AI-authored scaffolding.

**What gets tracked:**
- Number of attempts to make all tests pass
- Which tests caught real implementation bugs
- What the human did differently from what the AI would have done
- Time from stubs to green

**Three activation levels:**
- `[human-writes]` markers in the spec — granular, per-section
- `/plan --human` — the planner identifies which phases should be human-writes
- `/build --human` — every phase is human-writes

### 3. Thinking records

Every phase of `/feature` produces a thinking record:

- What was decided
- Why (the argument that won)
- What was considered and rejected
- What data was consulted
- What assumptions were made
- What questions were skipped (and whether the skip was justified)

These are not documentation. They are evidence that thinking happened. A spec with shallow thinking records is visibly different from one with deep records — and a reviewer can tell in 30 seconds which is which.

### 4. Constitutional principles

`specs/CONSTITUTION.md` defines project-level invariants — things that must always be true regardless of what feature is being built.

Every command reads the constitution. If a change would violate a principle, the system flags it and requires explicit approval to override. You can't accidentally violate a principle.

### 5. Teach mode

Understanding decays. `/teach` walks developers through a codebase in layers: context, happy path, failure modes, invariants, connections. Optionally quizzes to verify understanding — not trivia, but questions that test whether you could debug a production issue.

### 6. Retro closes the loop

After a feature ships, `/retro` goes back to the spec's predictions and checks them against reality. Each prediction is scored: hit, partial, miss, or unknown. Misses are analyzed. Lessons feed forward.

## What this means for teams

**For individual contributors:** You maintain your ability to understand and debug the systems you build. The AI handles scaffolding. You handle the logic.

**For tech leads:** Visibility into thinking quality without reading every spec. Human-writes reporting shows which engineers are getting faster at critical code.

**For engineering managers:** Three metrics — adoption rate, thinking depth, human-writes ratio. Combined with retro data, this shows whether the team is getting smarter or just faster.

**For the organization:** The knowledge is in the people, not the AI. If the AI tool goes down, your team can still build and debug software.
