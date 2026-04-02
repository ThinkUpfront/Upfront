---
description: Use when the user needs to test an idea fast — "just build it", "let me see it", "spike this", "prototype", "I need to try something", "quick and dirty", "throwaway". Rapid prototype with compressed thinking and automatic debt tracking.
user-invocable: true
---

# Spike

You are building a rapid prototype to test an idea. Speed over quality. The goal is a clickable/runnable thing the user can react to — not production code.

## Step 1: One question

Ask exactly one thing:

"Describe what you want to build in as much detail as possible. The more you give me, the less I have to guess."

Wait for their answer. Do not ask follow-up questions yet.

## Step 2: Autocomplete and confirm

Fill in everything they didn't say. Make decisions — don't ask. Present the complete picture:

```
Here's what I'm going to build:

**What it does:** [one paragraph — what the user described plus what you inferred]

**What you'll see:**
- [screen/page/output 1]
- [screen/page/output 2]
- [screen/page/output 3]

**What's fake:**
- [hardcoded data instead of real API]
- [mock auth / no auth]
- [no persistence — refreshing loses everything]
- [no error handling]
- [no tests]

**Tech:** [simplest possible stack — if the project already has a frontend, use it. If not, plain HTML/CSS/JS or the lightest framework available]

**Estimated time:** [rough — "10 minutes", "30 minutes", etc.]

Does this look right? I'll start as soon as you say go.
```

If they correct something, update and re-confirm. Keep this to one round — don't iterate endlessly on a prototype spec.

## Step 3: Acknowledge the debt

Before building, be explicit about what's being deferred:

```
Debt we're taking on:
- No spec (will need /upfront:feature if this idea survives)
- No tests (will need to add before this goes to production)
- No architecture review (may need restructuring to fit the codebase)
- [anything else specific — no auth, no validation, hardcoded config, etc.]

This is a prototype. If the idea works, we'll do it right. Let's go.
```

## Step 4: Build the minimum

Build the fastest thing that lets the user evaluate the idea. Priorities:

1. **Fake the backend.** Use hardcoded data, JSON files, localStorage, or mock APIs. Do not build real persistence, real auth, or real integrations. The user needs to see and click — they don't need it to actually work end-to-end.

2. **Make it look real enough to react to.** The UI doesn't need to be polished, but it needs to be concrete enough that the user can say "yes this is what I want" or "no, it should work like X instead." Wireframe-quality is fine. Placeholder text is fine. But the layout and flow should be real.

3. **One happy path only.** No error states, no edge cases, no empty states. Build the golden path that demonstrates the idea.

4. **Skip everything else.** No linting, no type safety, no code review, no architecture. This code is disposable.

Work in the current branch — do not create a worktree. Spike code lives alongside the real code. It can be reverted easily with git.

While building, if you discover something that changes the shape of the idea (e.g., "this API doesn't support what you need" or "the data model makes this harder than expected"), stop and tell the user immediately. Don't silently work around it — the discovery IS the point of the spike.

## Step 5: Demo

When the prototype is runnable, tell the user how to see it:

```
Spike ready.

To see it: [exact command — e.g., "open index.html", "npm run dev → http://localhost:3000/prototype", etc.]

Try:
1. [first thing to click/do]
2. [second thing to click/do]
3. [third thing to click/do]

What do you think? Three options:
  a) This is the right idea — let's spec it properly (/upfront:feature)
  b) Close but needs changes — tell me what's wrong and I'll update the spike
  c) Wrong direction — kill it
```

If they choose (a), immediately launch `/upfront:feature` with the spike as context — the prototype answers most of the Phase 1 and Phase 2 questions already. Pass the spike description and what was learned. The feature spec should note: "Originated from spike on [date]. Prototype validated [what worked]. Spike code will be replaced by production implementation."

If they choose (b), iterate on the spike. Keep it fast — don't gold-plate. After 2-3 iterations, nudge toward (a) or (c): "We've iterated a few times. Is the idea validated enough to spec properly, or should we kill it?"

If they choose (c), offer to revert: "Want me to `git checkout .` to remove the spike code, or keep it around for reference?"

## Step 6: Log the debt

After the spike is done (regardless of outcome), append to `specs/DEBT.md` (create if it doesn't exist):

```markdown
## [date] — Spike: [one-line description]

**Status:** [kept — proceeding to /feature | killed — code reverted | killed — code kept for reference]
**Files touched:** [list]
**Debt items:**
- [ ] No spec — needs /upfront:feature if shipping
- [ ] No tests — needs test coverage before production
- [ ] No architecture review — may conflict with existing patterns
- [ ] [specific items — hardcoded config, fake auth, mock data, etc.]
**Severity:** [1 = cosmetic, 2 = needs fixing before production, 3 = structural risk]
**Created by:** spike
**Feature:** [name or "exploratory"]
```

If the user proceeds to `/upfront:feature`, the spec should reference the spike: "This feature was prototyped in a spike on [date]. The prototype validated [what was learned]. Debt items from the spike are tracked in `specs/DEBT.md`."

## Rules

- **Speed is the point.** If you're spending more than 2 minutes on any decision, pick the simpler option and move on.
- **Fake everything you can.** Real backend work in a spike is wasted effort 90% of the time.
- **Don't refactor.** Don't clean up. Don't add comments. Don't follow conventions. This code is disposable.
- **Do track the debt.** The fast part is building. The responsible part is logging what you skipped. Never skip the debt log.
- **Stop if the idea is dead.** If you realize mid-build that the idea fundamentally won't work (not "it's hard" but "it's impossible given the constraints"), stop immediately and tell the user why. Don't finish a prototype of a dead idea.
- **One spike at a time.** Don't let spikes accumulate. If there's already an open spike in `specs/DEBT.md` (status: kept, no corresponding feature spec), flag it: "You have an open spike from [date]. Want to resolve that first or add another?"
