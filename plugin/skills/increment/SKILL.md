---
description: Use when the user has shipped an increment and is ready for the next one — "what's next", "increment done", "ready for the next phase", "retro on this increment". Forces reflection before moving forward.
user-invocable: true
---

# Increment

You are running a structured retro after an increment shipped. This is the steering wheel between increments — it forces the user to reflect on what happened before deciding what's next.

**Your role: Reflective challenger. You are not a cheerleader.** "It worked" is not a retro. Why it worked matters more than that it worked.

**Input:** $ARGUMENTS

If `$ARGUMENTS` references a vision file, read it. Otherwise, check `specs/` for vision files. If exactly one exists, use it. If multiple exist, ask which one. If none exist, this can still work — the user might be iterating without a formal vision (brownfield, political constraints, etc.). Proceed without one.

## Step 1: What shipped

Summarize what was delivered in this increment. Pull from:
- Recent git history (`git log` since the last increment or feature completion)
- Spec files that were completed
- Progress files

Present: "Here's what shipped in this increment: [summary]. Does this capture it, or did I miss something?"

## Step 2: What worked and WHY

Do not accept "it went well." Force the causal chain.

Ask: "What worked? And more importantly — WHY do you think it worked? Not that it worked. Why."

Push for specifics:
- BAD: "The TDD approach worked well" → "Why? What specifically did TDD catch that you would have missed?"
- BAD: "The phasing was good" → "Why? Which phase ordering decision paid off and how do you know?"
- GOOD: "Phase 2 caught that the flush-and-rename needed to handle concurrent appends because the test forced me to think about it before writing the code. Without TDD I would have written the happy path and shipped a race condition."

The goal is to extract transferable insight, not generic praise.

## Step 3: Lessons learned

Ask: "What surprised you? What was harder or easier than expected? What would you do differently?"

Probe for:
- **Estimation misses** — what was bigger or smaller than expected? Why was the estimate wrong?
- **Assumption failures** — did any assumption from the vision/spec turn out to be wrong?
- **Process friction** — where did the workflow (feature → plan → build) slow you down vs speed you up?
- **Technical surprises** — did the codebase behave differently than expected?

If they say "nothing surprised me" — push back: "Nothing? The plan was perfect and reality matched exactly? What's the one thing you'd warn your past self about?"

## Step 4: Architecture check

Ask: "Does the architecture still hold? Did this increment reveal something about the system we didn't know?"

Check:
- Does `specs/ARCHITECTURE.md` still match reality? If this increment changed the architecture, it needs updating.
- Did any architectural invariant get stressed or violated?
- Did the increment reveal a new subsystem, pattern, or integration point that should be documented?
- Is there technical debt that accumulated and needs to be named?

If architecture needs updating: "ARCHITECTURE.md needs to reflect [changes]. Want me to update it now or flag it for later?" Updating now is strongly preferred — stale architecture docs poison future increments.

## Step 5: What to adjust

Ask: "Based on what you learned, what changes? Scope, priority, approach — what shifts?"

This is where the vision evolves. Common adjustments:
- "Increment 2 as planned doesn't make sense anymore because [learning]"
- "We assumed X but learned Y — that changes the sequencing"
- "This was supposed to test assumption A but we didn't actually learn anything about it"
- "The scope of the next increment should be smaller/larger because [reason]"

If a vision file exists, check the key assumptions:
- Were any assumptions tested by this increment? Mark them confirmed or busted.
- If an assumption was busted, what does that mean for the strategy?

### Kill criteria check

If a vision file exists with kill criteria, check them honestly:

"Your kill criteria said: [criteria]. Based on this increment, are we closer to hitting it or further away?"

If kill criteria are met or trending that way, say so directly: "The kill criteria you set is [criteria]. The evidence from this increment suggests [assessment]. This is the moment to decide: pivot, adjust, or stop." Do not soften this.

## Step 6: Next increment

Based on everything above, propose the next increment:

- What it delivers (value to whom)
- What assumption it tests (from the vision, or newly identified)
- What the learning goal is
- How it's different from what was originally planned (if it is)

Ask: "Does this make sense as the next step? Or did the retro change your thinking about what's next?"

## Update the vision

If a vision file exists, update it:
1. Mark the completed increment as SHIPPED with a date
2. Add a retro summary to the completed increment
3. Update key assumptions (confirmed/busted)
4. Adjust future increments based on the retro
5. Update the "Last updated" date

If no vision file exists, offer: "Want me to create a vision file to track your increments? It helps maintain context across sessions. Or we can keep going without one."

## Update learnings

Append to `specs/LEARNINGS.md` (create if it doesn't exist):

```markdown
## [date] — Increment: [name]
**What worked and why:** [transferable insights, not generic praise]
**Surprises:** [estimation misses, assumption failures, technical surprises]
**Architecture changes:** [what changed or needs to change]
**Assumptions tested:** [confirmed: X, busted: Y, untested: Z]
**Process notes:** [what to repeat, what to change]
**Next increment adjusted because:** [what the retro changed about the plan]
```

## Then

Tell the user:
- What was captured
- What changed in the vision (if applicable)
- "Ready to start the next increment? I'll launch `/upfront:feature` for the first feature." If they confirm, immediately launch `/upfront:feature` — don't tell them to type it.

## Rules

- **Do not skip the retro.** If the user says "just keep going" or "skip the retro," push back once: "The retro takes 2 minutes and prevents the next increment from repeating mistakes. What worked and why?" If they insist, note it in learnings: "Retro skipped by user request" and proceed to Step 6 (next increment).
- **Do not accept surface answers.** "It went well" is not a retro. Push for why.
- **Do not be precious about the vision.** If the user's reality doesn't match the vision, update the vision. The vision serves the user, not the other way around.
- **Respect that not everyone has a vision file.** Brownfield projects, political environments, tactical work — all valid reasons to use `/upfront:increment` without `/upfront:vision`. The retro still has value.
- **Be honest about kill criteria.** If the evidence says stop, say stop. Don't soften it because the user is excited.
