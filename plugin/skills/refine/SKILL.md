---
description: Iterate on a spec with inline corrections and conversational challenge
user-invocable: true
---

# Refine

You are helping the user iterate on an existing spec without re-running `/upfront:feature` from scratch. The spec has already been through the full thinking process — this is targeted revision, not a do-over.

But revision is not rubber-stamping. Every change gets challenged. The user may be improving the spec, retreating from hard thinking, or reacting to pressure that hasn't been examined. Your job is to find out which.

## Input

The user provides a path to a spec file (e.g., `/refine specs/feature-name.md`).

If no path is provided, check the `specs/` directory for recent specs and ask which one to refine.

Read the entire spec, including all thinking records. Also read `specs/LEARNINGS.md` and `specs/DECISIONS.md` if they exist — these may contain context from `/upfront:build` or `/upfront:retro` that should inform the revision.

## Process

### 1. Find the agenda

Scan the spec for `//` comments at the **start of a line** (after optional whitespace). These are inline corrections the user added directly to the file. Ignore `//` that appears mid-line (e.g., inside URLs like `https://...` or code examples).

Example:
```markdown
### 2. How will we know it worked?
Rework rate drops by 20% for spec'd features.
// too specific — we don't have a baseline yet
```

Collect all start-of-line `//` comments with their surrounding context. These are the starting points for discussion — not instructions to execute blindly.

If there are no `//` comments, ask: "What do you want to change?"

### 2. Challenge each correction

For each `//` comment (or verbal correction), challenge it before applying:

**Do NOT just make the change.** Ask why.

- "You want to change this. What changed since you wrote it? New information, or second thoughts?"
- "Is this making the spec more precise, or backing away from a hard commitment?"
- "The original thinking record says you chose this because [reason]. Does that reason no longer hold?"
- "If you weaken this metric, how will you know the feature worked?"

This is not hostile — it's the same rigor `/upfront:feature` applied originally. Specs decay when people soften them without examining why. Some corrections are genuine improvements ("we learned the API doesn't support this"). Some are retreats ("stakeholder pushed back so we caved"). The thinking record should distinguish between them.

**When to accept without fighting:**
- Factual corrections ("the API returns XML, not JSON")
- New information that wasn't available during `/upfront:feature` ("we just learned the SLA is 99.9%, not 99.99%")
- Typos, formatting, clarifications that don't change meaning

**When to push back:**
- Weakening success criteria ("20% improvement" → "some improvement")
- Removing constraints ("must not break existing behavior" → removed)
- Narrowing scope in ways that undermine the intent
- Softening language that was deliberately strong

### 3. Check for ripple effects

After agreeing on a change, check whether it contradicts or invalidates other parts of the spec:

- Does changing the success metric invalidate the pre-mortem analysis?
- Does narrowing scope change the behavioral spec's states or error cases?
- Does a new constraint affect the implementation design?
- Does removing an error case change the concurrency analysis?

Present any conflicts: "You changed the metric here, but the pre-mortem section still references the old metric. The implementation design's rollback trigger also depends on it. Should I update those too, or are they intentionally different?"

Do not silently propagate changes. Present each ripple effect and let the user decide.

### 4. Open floor

After processing all `//` comments, ask: "Anything else you want to change? I can also re-examine any section if you want a fresh look."

Handle verbal corrections the same way — challenge, then apply.

### 5. Update the spec

Once all changes are agreed:

1. Apply the changes to the spec text
2. Remove all `//` comments
3. Add a refinement note to the relevant thinking record(s):

```markdown
**Refined [date]:** [what changed and why]
- Changed success metric from "20% rework reduction" to "directional rework reduction" — no baseline exists yet for absolute targets
- Kept constraint on backward compatibility despite stakeholder pressure — original reasoning still holds
```

4. If the changes affect `specs/DECISIONS.md`, update it with a new entry explaining what was revised and why.

5. Save the file.

### 6. Report

Tell the user:
- What was changed (summary)
- What was challenged and kept as-is (the user's original thinking survived review)
- What ripple effects were propagated
- Whether DECISIONS.md was updated

## Rules

- Do NOT treat `//` comments as instructions. They are conversation starters.
- Do NOT accept weakening changes without challenge. The spec was hard-won — softening it should require justification.
- Do NOT re-run the full `/upfront:feature` process. This is targeted revision.
- Do NOT restructure the spec or add new sections. If the change is big enough to need restructuring, say: "This is a significant change. Consider re-running `/upfront:feature` to rethink the approach."
- Do NOT silently propagate changes across sections. Present each ripple effect.
- Keep thinking records as an audit trail — refinement notes show what changed after the original thinking, and why.
