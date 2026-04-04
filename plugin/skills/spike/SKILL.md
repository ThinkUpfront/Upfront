---
description: Use when the user needs to test an idea fast — "just build it", "let me see it", "spike this", "prototype", "I need to try something", "quick and dirty", "throwaway". The primary path for unvalidated ideas — hypothesize, build, learn.
user-invocable: true
---

# Spike

You are running an experiment. The user has a hypothesis — your job is to build the fastest possible thing that tests it. This is the scientific method applied to software: build the minimum, observe the result, learn.

This is not a shortcut or a guilty pleasure. This is how good software gets built — test the idea before you invest in it.

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

## Step 3: Name what's deferred

Be explicit about what you're skipping — not as guilt, but as a checklist for later if the experiment succeeds:

```
Deferred (expected for a spike):
- No spec — if this works, solidify with /upfront:feature
- No tests — add before production
- No architecture review — may need restructuring
- [anything else specific — no auth, no validation, hardcoded config, etc.]

Let's go find out if this idea works.
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

What did you learn? Three paths:
  a) The idea works — let's solidify it (/upfront:feature → /plan → /build)
  b) Close but wrong — tell me what's off and I'll adjust the experiment
  c) Dead end — kill it and save the lesson
```

If they choose (a): the experiment succeeded. Immediately launch `/upfront:feature` with the spike as context. The prototype answers most of the Phase 1 and Phase 2 questions already. The feature spec captures what was validated and builds it properly — real tests, real architecture, real error handling.

If they choose (b): iterate on the spike. Keep it fast. After 2-3 iterations, nudge: "We've learned enough — is the idea validated or should we kill it?"

If they choose (c): the experiment failed, which is a success — you learned something without wasting weeks. Offer to revert: "Want me to `git checkout .` to clean up, or keep it for reference?"

## Step 6: Log the experiment

After the spike is done (regardless of outcome), append to `specs/DEBT.md` (create if it doesn't exist):

```markdown
## [date] — Spike: [one-line description]

**Status:** [validated — proceeding to /feature | killed — code reverted | killed — code kept for reference]
**What we learned:** [the actual insight — this is the most important line]
**Files touched:** [list]
**Deferred items (resolve during solidification):**
- [ ] No spec — needs /upfront:feature
- [ ] No tests — needs test coverage
- [ ] No architecture review — may need restructuring
- [ ] [specific items — hardcoded config, fake auth, mock data, etc.]
**Created by:** spike
**Feature:** [name or "exploratory"]
```

If the user proceeds to `/upfront:feature`, the spec should reference the spike: "Validated by spike on [date]. The experiment confirmed [what was learned]. Deferred items tracked in `specs/DEBT.md` — resolve during solidification."

## Rules

- **Speed is the point.** If you're spending more than 2 minutes on any decision, pick the simpler option and move on.
- **Fake everything you can.** Real backend work in a spike is wasted effort 90% of the time.
- **Don't refactor.** Don't clean up. Don't add comments. Don't follow conventions. This code is disposable — the solidification phase builds it right.
- **Log what you learned.** The experiment log is not about guilt — it's about capturing the insight so the solidification phase has a head start.
- **Stop if the idea is dead.** If you realize mid-build that the idea fundamentally won't work, stop immediately and tell the user why. A killed experiment is a successful experiment — you learned something.
- **One spike at a time.** Don't let spikes accumulate. If there's already an open spike in `specs/DEBT.md` (status: validated, no corresponding feature spec), flag it: "You have a validated spike from [date] that hasn't been solidified. Want to solidify that first or start a new experiment?"
