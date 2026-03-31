---
title: Upfront vs Superpowers
description: How Upfront compares to the Superpowers plugin for Claude Code.
---

Both Upfront and [Superpowers](https://github.com/obra/superpowers) are Claude Code plugins that enforce structure before code. They share the same instinct — don't let the AI jump straight to implementation. But they make different bets on what matters most.

## At a glance

| | Upfront | Superpowers |
|---|---|---|
| **Core bet** | The human's thinking is the product | The agent's discipline is the product |
| **Who thinks** | Human thinks, AI challenges | Human approves, AI plans |
| **Brainstorming** | `/upfront:ideate` — AI asks, human generates | Mandatory brainstorm gate — blocks progress until done |
| **Feature definition** | 4-phase adversarial conversation (intent → spec → design → implementation) | Spec document generated from brainstorm |
| **Planning** | `/upfront:plan` — architecture deep-dive, ~400 LOC phases | Auto-generated implementation plan from spec |
| **Building** | `/upfront:build` — TDD, fresh subagent per phase, red team | TDD with subagent-driven execution, two-stage review |
| **Human role during build** | Reviews each phase, decides on judgment calls | Approves spec, then subagents execute with automated review |
| **Reviewability** | Reviewability scoring (5 dimensions), escalates to `/upfront:vision` if too complex | No explicit reviewability check |
| **Strategic planning** | `/upfront:vision` for multi-feature ambitions, `/upfront:increment` for retros | Not addressed — operates at feature level |
| **Audit trail** | JSONL thinking records with remote flush to observability tools | None |
| **Config protection** | Blocks AI from weakening linter/formatter rules | None |
| **License** | AGPL-3.0 | MIT |

## Where Superpowers is stronger

- **Execution automation** — Superpowers' subagent-driven development mode can execute a full plan with two-stage automated review (spec compliance, then code quality) with less human involvement.
- **Structured brainstorm gate** — brainstorming is designed as a prerequisite in the workflow, though in practice it can be skipped. Both tools push back on skipping thinking but neither physically prevents it.
- **Composability** — skills activate automatically based on context. You don't need to know which command to run.
- **Maturity** — 40k+ stars, in the official Anthropic plugin marketplace, created by Jesse Vincent.

## Where Upfront is stronger

- **Adversarial thinking** — the AI doesn't suggest and wait for approval. It asks open questions, waits for your answer, then fills gaps. The goal is to make *you* think, not to generate a spec *for* you.
- **Reviewability gates** — scores changes on 5 dimensions (concern count, blast radius, novelty, state complexity, reversibility) and pushes back before any code is written if the change is too complex for meaningful review.
- **Strategic layer** — `/upfront:vision` captures multi-feature ambitions using Rumelt's kernel (diagnosis, guiding policies, coherent actions). `/upfront:increment` forces structured retros between increments. Superpowers operates at the feature level.
- **Audit trail** — every phase produces a thinking record captured as structured events. Flushable to Langfuse, Arize Phoenix, or any observability tool. Managers can see adoption, depth, and effectiveness without reading specs.
- **Security guardrails** — `/upfront:plan` audits the project for missing security tooling (linters, type checkers, vulnerability scanners, secret detection, slopsquatting protection) and pushes hard for installation as Phase 0 before any feature code. The config protection hook blocks the AI from weakening rules after they're installed.
- **Architecture as a living document** — `/upfront:plan` checks if `ARCHITECTURE.md` is stale (>30 days with commits), actively compares it to the codebase, and updates it before planning. `/upfront:architect` does full structural reviews. The architecture is revisited regularly, not written once and forgotten.

## They're complementary

Superpowers optimizes the agent's execution. Upfront optimizes the human's thinking. You could use both — Upfront for feature definition and strategic planning, Superpowers for disciplined execution.

The question is what you believe is the bottleneck: is it the AI producing bad code, or the human not thinking clearly about what to build?

## Install Upfront

```bash
claude plugin marketplace add ThinkUpfront/Upfront
claude plugin install upfront
```
