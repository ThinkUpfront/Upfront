---
title: Why Upfront
description: The problem Upfront solves and how it solves it.
---

AI makes writing code effortless. That's the problem.

When code is cheap to produce, it's easy to stop thinking. Skip the spec. Skip the design. Let the AI figure it out. Ship it, fix it later. Except "later" is now a codebase full of confident mistakes that passed every automated check and got rubber-stamped in review.

## How teams lose their skills

Not all at once — gradually.

The senior engineer who used to catch architectural issues in review now approves AI-generated PRs because the code "looks right." The mid-level engineer who used to design systems now prompts for them. The junior engineer who should be learning to think never has to, because the AI thinks for them.

Nobody notices until something breaks and no one on the team understands the system well enough to fix it.

## What the research says

- **Anthropic (2026):** Developers who delegated code generation scored 40% lower on comprehension than those who used AI for conceptual questions. The interaction pattern matters more than whether AI was used.

- **METR (2025):** AI tooling increased completion time by 19% for experienced developers on familiar codebases. Speed is not guaranteed.

- **DORA (2025):** AI adoption increases throughput but also increases delivery instability. Faster is not more reliable.

- **Microsoft Research (2024):** Developers using AI assistants believed their code was more secure than it actually was. They reviewed AI code less carefully.

- **GitClear (2024):** Across 211 million changed lines, code churn doubled while refactoring collapsed. More code, less of it survives production.

[Full research citations](/research/)

## What Upfront does

Upfront makes the thinking process explicit, challenging, and auditable. Every command is designed to force humans to engage — not fill out templates, not click "yes" to AI suggestions, but actually think through what they're building and why.

**Challenge first, decorate second.** The AI never leads with suggestions. It asks an open question, waits for your answer, and then adds what you missed. The difference between "do you approve this?" (rubber stamp) and "walk me through what happens when two of these run at the same time" (forcing function).

**Thinking records.** Every phase of `/feature` produces a record of what was decided, why, what was rejected, and what was skipped. The spec is the audit trail of the thinking, not just the conclusions. A reviewer can tell in 30 seconds whether real thinking happened.

**Human-writes mode.** For critical code — concurrency, security, core business logic — the AI writes the tests and the human writes the implementation. You cannot ship code you don't understand. [Read more](/human-first/)

**Constitutional governance.** Project-level invariants that gate every command. If a change would violate a constitutional principle, the system flags it. You can override, but you can't accidentally violate.

**Closed feedback loop.** `/retro` goes back to the spec's predictions and checks them against production reality. Did it actually work? Where were we wrong? Lessons feed forward into future specs.
