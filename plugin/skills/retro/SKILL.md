---
name: retro
description: "Close the loop on a shipped feature — check spec predictions against production reality"
---

# /retro — Post-Ship Retrospective

You are running a structured retrospective that compares a feature spec's predictions against what actually happened in production. The goal is to improve future specs by learning from prediction accuracy.

## Input

The user provides a spec path as an argument: `$ARGUMENTS`

If no argument is provided or the path is empty:
1. Look for spec files in `specs/` directory (or similar spec locations in the project)
2. List the available specs with their "How will we know it worked?" summary
3. Ask which spec to review
4. Wait for the user's response before proceeding

## Process

### Step 1: Read the spec

Read the spec file and extract:
- **Success metrics** from "How will we know it worked?" (or similar success criteria section)
- **Pre-mortem risks** from "Pre-mortem findings" (or similar risk section)
- **Constraints** from "What must NOT happen" (or similar constraint section)
- **Key decisions** from thinking records, decision log, or similar sections

Present what you found:
> "Here's what the spec predicted. Let me check these against reality."

List each prediction clearly so the user can respond to them.

### Step 2: Gather evidence

For each success metric, ask:
> "The spec predicted: **[metric]**. What actually happened? Do you have the numbers?"

For pre-mortem risks, ask:
> "The pre-mortem identified these risks: **[list risks]**. Did any of these materialize?"

Then ask:
> "Were there any unexpected outcomes — good or bad — that the spec didn't predict?"

**Rules for evidence gathering:**
- If the user doesn't have data yet, stop here: "No production data available yet. Re-run `/retro` when you have **[specific metric]** data."
- Do NOT fabricate data or assume outcomes.
- If the user says "I think it improved," push back: "Do you have the actual number, or is that a feeling? Feelings are fine to note but they're not evidence."
- If the user gives partial data, work with what's available but flag what's missing.

### Step 3: Score each prediction

For every prediction the spec made, assign one score:

| Score | Meaning |
|-------|---------|
| **Hit** | The predicted outcome happened as described |
| **Partial** | Directionally correct but magnitude was off |
| **Miss** | Didn't happen, or the opposite happened |
| **Unknown** | No data available to evaluate |

Present the scorecard clearly.

### Step 4: Analyze the gaps

For each **Miss** and **Partial**, work through:
> "Why was the prediction wrong?"

Guide the user through three possible causes:
1. **Mechanism** — the approach itself didn't work as expected
2. **Measurement** — we measured the wrong thing, or the metric didn't capture the real outcome
3. **Environment** — something external changed (market, user behavior, dependencies)

Then ask:
> "Should the spec have caught this? What question in `/feature` would have surfaced it?"

For **unexpected outcomes** (good or bad):
> "Was this predictable? What signal did we miss?"

### Step 5: Extract lessons

Challenge the user to generalize. Ask these questions one at a time, not all at once:

1. "What would you do differently if you were speccing this feature again?"
2. "Does this change how you'd answer the forcing-function questions for the next feature?"
3. "Is there a pattern here that applies beyond this specific feature?"

Push for specifics. "Be more careful" is not a lesson. "Add a load test to the pre-mortem checklist" is.

### Step 6: Record

**Append to `specs/LEARNINGS.md`** (create if it doesn't exist):

```markdown
## [today's date] — Retro: [feature name]
**Predicted:** [what the spec said would happen]
**Actual:** [what actually happened]
**Score:** [N hits, N partials, N misses, N unknown]
**Why predictions missed:** [analysis from step 4]
**Lessons for future specs:** [what to do differently]
**Process notes:** [did /feature ask the right questions? did /build catch the right things?]
```

**Add a Retrospective section to the bottom of the spec file itself:**

```markdown
---
name: retro
## Retrospective ([today's date])
**Outcome:** [hit/partial/miss for each prediction]
**Key learning:** [one sentence]
```

### Step 7: Feed forward

Review the retro findings for patterns that should change how `/feature`, `/plan`, or `/build` works. If you see one, say so explicitly. Examples:

- "This is the third feature where the success metric was unmeasurable. Should `/feature` push harder on metric availability?"
- "The concurrency bug that shipped suggests `/build`'s red team should specifically test [pattern]."
- "The pre-mortem missed the biggest risk. Should `/feature` add a question about [category]?"

These are suggestions for evolving the system, not automatic changes. Present them as observations for the user to act on.

## Rules

- Do NOT run a retro before there's production data. It's worthless without evidence.
- Do NOT accept "it feels better" as evidence. Push for numbers.
- Do NOT blame the team. This is about the process and predictions, not people.
- Be honest about misses. A retro that says everything went perfectly is either lying or not looking hard enough.
- Do NOT fabricate or assume any production data. Every data point must come from the user.
- Keep the conversation focused on learning, not judgment.
