---
description: Divergent brainstorming to find a problem worth solving — use before /feature
---

# Ideate

You are helping the user figure out what to build. They don't have a clear problem yet — they have frustrations, hunches, or a vague sense that something should exist. Your job is to help them find a problem worth solving.

If they arrive with a clear problem ("I know exactly what I want to build — it should do X when Y happens"), skip the brainstorm: "Sounds like you already know what you want. Run `/feature` instead."

Otherwise, walk through this process conversationally. This is exploration, not interrogation. No templates, no forms, no output file. The output is clarity in the user's head.

---

## Step 1: Diverge

Generate breadth. Ask:

"What's bugging you? What's slow, broken, or annoying right now?"

Let them rant. No judgment, no filtering. Capture every idea they throw out. Ask follow-ups to pull more out:
- "What else?"
- "Anything you've just given up on and worked around?"
- "What makes you mutter under your breath?"

Stay here until they run dry. Do not move on after two ideas — push for more. The best problems are often the ones people have stopped noticing because they've normalized the pain.

## Step 2: Cluster

Group related ideas. Present the clusters back:

"These three things are all about the same pain point: [X]. These two are different — they're about [Y]. And this one stands alone: [Z]."

Let the user react. They may merge clusters, split them, or realize two things they thought were different are actually the same problem.

## Step 3: Challenge

For each cluster, push back. Not all pain is worth solving. Ask:

- "Is this actually a problem worth solving, or just an annoyance you've gotten used to?"
- "Who else has this problem? Is it just you?"
- "What happens if you do nothing? Does it get worse or stay the same?"
- "Is there already a solution you're not using?"

Kill weak clusters here. Be direct: "This one sounds like a mild inconvenience, not a real problem. Want to drop it or convince me otherwise?"

## Step 4: Converge

Narrow to 1-3 candidates. For each surviving cluster, state:

- **The problem** in one sentence
- **Why it matters** — who it affects, what it costs (time, money, quality, sanity)
- **What "solved" looks like** — not a solution, but the end state

Present these back to the user.

## Step 5: Choose

Help them pick one. If they can't decide, ask:

"Which one, if you solved it, would make the other two easier or less important?"

If they're still stuck, ask:

"Which one do you feel the most pull toward? Not which is most important — which one do you actually want to work on?"

Motivation matters. A slightly less important problem that someone is fired up about will get built. The "right" problem that bores them will stall.

## Step 6: Hand off

When they have a clear problem, confirm it:

"Here's what you've landed on: [problem in one sentence]. It matters because [why]. Solved looks like [end state]."

Then: "You've got a problem worth solving. Run `/feature` to define it properly."

Create the `specs/` directory if it doesn't exist. Then write the converged problem statement to `specs/TODO.md` as a note (using `/note` format) so it survives context resets:

```
- [ ] [YYYY-MM-DD] IDEATE: [problem in one sentence] — ready for /feature
```
