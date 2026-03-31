---
title: Upfront vs GSD
description: How Upfront compares to the Get Shit Done (GSD) framework for Claude Code.
---

Both Upfront and [GSD](https://github.com/gsd-build/get-shit-done) enforce structure on AI-assisted development. Both believe specs beat vibe coding. But they solve different problems and make different tradeoffs.

## At a glance

| | Upfront | GSD |
|---|---|---|
| **Core bet** | Protect human thinking and judgment | Maximize autonomous execution with persistent context |
| **Planning model** | Human-driven with AI challenge | AI-driven with human checkpoints |
| **Context strategy** | Thinking records in specs, crash recovery from git | `.planning/` directory with STATE.md, CONTEXT.md per phase |
| **Autonomous mode** | Building can be autonomous, but thinking and planning require a human | Yes — `/gsd:autonomous` runs all phases unattended |
| **Feature definition** | 4-phase adversarial conversation | `discuss-phase` with targeted questions (or `--auto`) |
| **Planning** | Architecture deep-dive, ~400 LOC phases, reviewability scoring | Research + plan with Nyquist verification (every task has automated verify) |
| **Building** | TDD, fresh subagent per phase, red team, human review | Wave-based parallel subagents, conversational UAT |
| **Verification** | Post-phase code review + red team adversarial pass | Nyquist validation (every requirement has an automated test command) |
| **Project scope** | Feature-level (with `/upfront:vision` for bigger) | Project-level (milestones, roadmaps, requirements traceability) |
| **Session persistence** | `/upfront:pause` + `/upfront:resume` with HANDOFF.md | STATE.md + CONTEXT.md survive context resets automatically |
| **Audit trail** | JSONL thinking records with remote flush | None |
| **Config protection** | Blocks AI from weakening linter rules | None |
| **Multi-runtime** | Claude Code only | Claude Code + 7 other runtimes |
| **License** | AGPL-3.0 | Check repo |

## Where GSD is stronger

- **Autonomous execution** — `/gsd:autonomous` can run an entire project from discuss through ship without human intervention. Upfront requires human engagement at every phase by design.
- **Context engineering** — GSD's `.planning/` directory is a sophisticated persistence layer. STATE.md, CONTEXT.md, REQUIREMENTS.md, and ROADMAP.md survive context resets. Upfront relies on specs/ files and git history.
- **Project-level orchestration** — milestones, roadmaps, requirements traceability, phase auditing. GSD manages an entire project lifecycle. Upfront is feature-focused (with `/upfront:vision` as an optional strategic layer).
- **Nyquist verification** — every task must have an automated verify command. This is a stronger guarantee than Upfront's post-phase review.
- **Multi-runtime** — works with Claude Code, OpenCode, Gemini, Codex, Copilot, Cursor, Windsurf. Upfront is Claude Code only.
- **Parallel execution** — wave-based subagent execution runs independent tasks concurrently.

## Where Upfront is stronger

- **Human engagement** — Upfront's core design is adversarial. The AI asks open questions, waits for your answer, then fills gaps. GSD's `--auto` mode explicitly minimizes human interaction. If you believe skill atrophy is the real risk, Upfront forces the exercise that prevents it.
- **Reviewability** — scores changes on 5 dimensions and pushes back if a change is too complex for meaningful human review. GSD doesn't gate on reviewability.
- **Thinking records** — every phase produces a structured audit trail: what was decided, why, what was rejected, what was skipped. Captures reasoning, not just conclusions. Flushable to Langfuse, Arize Phoenix, etc.
- **Strategic planning** — `/upfront:vision` uses Rumelt's kernel (diagnosis, guiding policies, coherent actions) for multi-feature ambitions. `/upfront:increment` forces structured retros between increments with kill criteria. GSD has milestones but not an opinionated strategic framework.
- **Config protection** — blocks the AI from weakening linter/formatter rules instead of fixing the code.
- **Simplicity** — 20 skills, zero dependencies, Go stdlib binary. No `.planning/` directory structure to learn, no STATE.md format to understand.

## Different philosophies

GSD asks: "How do I get the AI to ship reliably at scale?"

Upfront asks: "How do I keep the human sharp while the AI ships?"

GSD is built for throughput — autonomous execution, parallel subagents, minimal human interruption. Upfront is built for judgment — adversarial questioning, forced thinking, audit trails that prove real reasoning happened.

If your bottleneck is getting code out the door reliably, GSD is the better fit. If your concern is that your team is losing the ability to think critically about what they're building, Upfront addresses that directly.

## Install Upfront

```bash
claude plugin marketplace add ThinkUpfront/Upfront
claude plugin install upfront
```
