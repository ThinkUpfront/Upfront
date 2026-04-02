---
title: /upfront:refine
description: Iterate on a spec with inline corrections and challenge.
---

Targeted revision without re-running `/upfront:feature`.

## When to use

- Spec needs updates based on new information
- Stakeholder feedback changed requirements
- You want to adjust without starting over

## What it does

Add `//` comments at the start of a line in the spec to mark corrections:

```markdown
Rework rate drops by 20% for spec'd features.
// too specific — we don't have a baseline yet
```

Then run `/upfront:refine specs/[name].md`. The AI:

1. Collects all `//` comments as agenda items
2. **Challenges each change** — is this making the spec better or backing away from a hard commitment?
3. Checks for ripple effects across sections
4. Applies agreed changes and adds refinement notes to thinking records
5. Updates `specs/DECISIONS.md` if affected

Factual corrections are accepted without challenge. Weakening changes get pushed back on.

## Output

Updated spec with refinement audit trail in the thinking records.
