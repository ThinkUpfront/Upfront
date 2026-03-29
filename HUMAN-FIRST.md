# Human-First Development

AI writes code. Humans understand systems. When the AI writes all the code, humans stop understanding the systems they're responsible for. This document explains how Upfront keeps humans strong.

---

## The problem with full delegation

When an AI generates code and a human approves it, two things happen:

1. The code ships faster.
2. The human's understanding of the system gets shallower.

This works until it doesn't. The system breaks in production, and the person on-call doesn't understand the code they approved. The senior engineer reviews a PR and rubber-stamps it because the tests pass and the code "looks right." The junior engineer never learns to design because the AI designs for them.

The research is clear: developers who delegate code generation score 40% lower on comprehension than those who use AI for conceptual questions (Anthropic, 2026). AI-generated code gets reviewed less carefully despite having comparable vulnerability rates (Microsoft Research, 2024). Teams ship faster but less reliably (DORA, 2025).

Speed is not the bottleneck. Understanding is.

---

## How Upfront strengthens humans

### 1. Challenge before suggestion

Every command that involves decisions follows the same pattern: **ask, wait, then decorate.**

The AI does not present options and wait for approval. It asks open questions, waits for the human's answer, and then fills the gaps the human missed. The human thinks first. The AI adds what they didn't think of.

This is the difference between a rubber stamp and a forcing function. "Do you approve this design?" produces a yes. "What happens when two of these run at the same time? Walk me through it." produces understanding.

**Where this happens:**
- `/feature` Phase 1: Five forcing-function questions, each challenged before moving on
- `/feature` Phase 2: Mechanism question — "Why will this approach actually work?" before discussing how
- `/plan` Level 1-3: Each architectural level is confirmed with conversation, not presented as a report
- `/refine`: Every change is challenged — "Is this making the spec more precise, or backing away from a hard commitment?"

### 2. Human-writes mode

The most direct intervention. When a phase contains critical logic — concurrency, security, core business rules, invariant enforcement — the AI writes the tests and the human writes the implementation.

**How it works:**

1. The AI writes failing tests that describe what the code should do.
2. The AI generates `TODO(human)` stubs with constraints, invariants, and common mistakes.
3. The human fills in the stubs.
4. The AI runs the tests and reviews what the human wrote.
5. If tests fail, the AI shows which tests fail and why — but does not fix it for them.

The human is forced to understand the logic deeply enough to make the tests pass. There is no shortcut. You cannot delegate the implementation back to the AI without explicitly opting out. The commit is tagged `[human-writes]` so git history distinguishes human-authored critical code from AI-authored scaffolding.

**Three activation levels:**
- `[human-writes]` markers in the spec — granular, per-section
- `/plan --human` — the planner identifies which phases should be human-writes
- `/build --human` — every phase is human-writes

**What gets tracked:**
- Number of attempts to make all tests pass
- Which tests caught real implementation bugs (these tests proved their worth)
- What the human did differently from what the AI would have done
- Time from stubs to green (for teams that want velocity data alongside quality data)

This data feeds into `/retro` and LEARNINGS.md. Over time, the team can see: which phases are humans getting faster at? Which tests consistently catch real bugs? Which patterns does the AI get wrong that humans get right?

### 3. Thinking records as audit trail

Every phase of `/feature` produces a thinking record — not just the decision, but the reasoning:

- What was decided
- Why (the argument that won)
- What was considered and rejected (the arguments that lost)
- What data was consulted
- What assumptions were made
- What questions were skipped (and whether the skip was justified)

These records are not documentation. They are evidence that thinking happened. A spec with shallow thinking records is visibly different from one with deep records — and a reviewer can tell in 30 seconds which is which.

The records also create accountability. When a feature fails in production, `/retro` goes back to the thinking record and asks: was this failure foreseeable? Was there a skipped question that would have caught it? The feedback loop is closed, not open.

### 4. Constitutional principles

`specs/CONSTITUTION.md` defines project-level invariants — things that must always be true regardless of what feature is being built. "Every write must be durable before acknowledging." "No external dependencies." "Auth checked on every entry point."

Every command reads the constitution. If a proposed change would violate a constitutional principle, the system flags it and requires explicit user approval to override. You can't accidentally violate a principle — you can only deliberately amend the constitution.

This is governance that travels with the code, not governance that lives in a wiki nobody reads.

### 5. Teach mode

`/teach` exists because understanding decays. When a developer hasn't touched a codebase in a while, or when a new team member joins, `/teach` walks them through the system in layers:

1. What does this system do? (context)
2. What happens when a request comes in? (happy path)
3. What happens when things go wrong? (failure modes)
4. What must always be true? (invariants)
5. How do the pieces connect? (integration points)

Optionally quizzes to verify understanding — not trivia questions, but questions that test whether you could debug a production issue: "If the queue stops flushing, what are the three most likely causes?" Generates a study guide with strong areas, focus areas, and a safe first task.

### 6. Retro closes the loop

`/retro` is the accountability mechanism. After a feature ships and has production data, it goes back to the spec's predictions and checks them against reality:

- "How will we know it worked?" → Did it actually work? Show the numbers.
- "What must NOT happen?" → Did any of the negative conditions occur?
- "Pre-mortem failures" → Did any of the predicted failures happen?

Each prediction is scored: hit, partial, miss, or unknown. Misses are analyzed: was the mechanism wrong, the measurement wrong, or the environment different than assumed?

The lessons feed forward. If retros consistently show that a particular type of prediction misses, that pattern gets surfaced in future `/feature` runs as a known blind spot.

---

## What this means for teams

### For individual contributors

You maintain your ability to understand and debug the systems you build. The AI handles scaffolding, test infrastructure, and repetitive code. You handle the logic that makes the system work. When production breaks at 2 AM, you understand the code well enough to fix it.

### For tech leads

You get visibility into thinking quality without reading every spec. The audit trail shows: did the team think through concurrency? Did they skip the pre-mortem? Were thinking records deep or shallow? This is 30-second triage, not multi-hour review.

Human-writes reporting shows which engineers are getting faster at critical code, and which tests consistently catch real bugs — data for coaching conversations, not performance reviews.

### For engineering managers

Three metrics without reading a single spec:

1. **Adoption** — what percentage of features went through `/feature`?
2. **Depth** — how many phases were completed? (All 4 = thorough thinking. Bailed at 2 = shortcuts.)
3. **Human-writes ratio** — what percentage of critical code was human-authored vs AI-generated?

Combined with `/retro` data (prediction accuracy over time), this shows whether the team is getting smarter or just getting faster.

### For the organization

The knowledge is in the people, not the AI. If the AI provider changes, the model degrades, or the tool goes down, your team can still build and debug software. The AI accelerated them — it didn't replace the skills they need to function without it.

---

## The principle

**The thinking is the product, not the code.**

Code is generated, reviewed, tested, refactored, and eventually deleted. What survives is the team's understanding of their systems — the judgment to know what to build, the skill to know when it's broken, and the context to know why it was built that way.

Upfront exists to protect that understanding. Not by slowing teams down, but by making sure the thinking that should happen actually happens — visibly, challengingly, and on the record.
