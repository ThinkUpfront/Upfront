# Upfront

**Force thinking before code.**

AI makes writing code effortless. That's the problem.

When code is cheap to produce, it's easy to stop thinking. Skip the spec. Skip the design. Let the AI figure it out. Ship it, fix it later. Except "later" is now a codebase full of confident mistakes that passed every automated check and got rubber-stamped in review.

This is how teams lose their skills. Not all at once — gradually. The senior engineer who used to catch architectural issues in review now approves AI-generated PRs because the code "looks right." The mid-level engineer who used to design systems now prompts for them. The junior engineer who should be learning to think never has to, because the AI thinks for them.

Nobody notices until something breaks and no one on the team understands the system well enough to fix it.

Upfront solves this by making the thinking process explicit, challenging, and auditable. Every command in this toolkit is designed to force humans to engage — not fill out templates, not click "yes" to AI suggestions, but actually think through what they're building and why.

If the AI can't talk you out of your approach, your approach is probably sound. If it can, you should know that before writing a single line of code.

---

## How it works

Upfront is a set of slash commands for [Claude Code](https://claude.ai/claude-code) that cover the full development lifecycle. Each command is a markdown file in `.claude/commands/` — no dependencies, no build step, no SaaS. Copy the commands into any project and they work.

The commands follow a natural flow:

```
Think → Define → Plan → Build → Ship → Learn
```

Every command that involves decisions challenges you. The AI doesn't suggest and wait for approval — it asks open questions, waits for your answer, then fills the gaps you missed. You think first. The AI decorates second.

---

## Commands

### Think

#### `/ideate` — Find a problem worth solving

For when you don't know what to build yet. Divergent brainstorming that starts wide ("what's bugging you?"), clusters related pain points, challenges whether each is worth solving, converges to 1-3 candidates, and helps you pick one. No output file — the output is clarity. When you have a problem, it hands off to `/feature`.

#### `/explore` — Document the codebase and its ecosystem

Five-phase investigation that produces `specs/ARCHITECTURE.md` — the shared reference document every other command reads. Goes beyond code: maps external connections (latency, SLAs, failure modes, data volumes), ecosystem context (upstream, downstream, blast radius), and operational reality (deployment, monitoring, who gets paged). Asks about ecosystem diagrams.

Greenfield projects get a fast exit: detects empty repos, creates minimal scaffolding, sends you to `/feature`.

#### `/teach` — Rebuild your understanding

For when you haven't touched the codebase in a while. Gauges your current level (refresher, domain expert but new to code, or completely new), walks through the system in layers (context → happy path → failure modes → invariants → connections), and optionally quizzes you with questions that test real understanding — not trivia.

Generates a study guide with your strong areas, focus areas, key files to read, invariants to memorize, and a safe first task to build familiarity.

---

### Define

#### `/feature` — Define a feature through forced thinking

The core command. Four phases, each with a different AI role:

1. **Intent** — Five forcing-function questions. What problem? How will you know it worked? What's out of scope? What must NOT happen? Pre-mortem. The AI is adversarial — it pushes back on vague answers and refuses to move on until the thinking is substantive.

2. **Behavioral Spec** — Five-level funnel: user stories → mechanism (why will this approach work?) → states and transitions → concurrency and shared state → error cases. The AI challenges your causal logic before any code is discussed.

3. **Design Conversation** — The AI researches the codebase, presents options and tradeoffs. You make the decisions.

4. **Implementation Design** — The AI proposes architecture AND challenges the codebase itself. Flags inconsistent patterns, ambiguous placement, structural rot. Cleanup becomes a prerequisite.

Every phase transition produces a thinking record — what was decided, why, what was rejected, what was skipped. The spec is the audit trail of the thinking, not just the conclusions.

If your problem statement is vague, `/feature` suggests `/ideate`. If there's no architecture doc, it suggests `/explore`.

Appends design decisions to `specs/DECISIONS.md` for future reference.

#### `/refine` — Iterate on a spec without starting over

Add `//` comments directly in the spec file to mark corrections, then run `/refine`. The AI treats each `//` comment as an agenda item — but challenges every change before applying it. Weakening a success metric? Why? Removing a constraint? What changed?

Checks for ripple effects across sections and updates thinking records with refinement notes.

---

### Plan

#### `/plan` — Break a spec into buildable phases

Starts with a three-level architectural deep-dive:

1. **System architecture** — components, communication, system-level invariants
2. **Subsystem architecture** — modules, boundaries, data models, subsystem invariants
3. **Design patterns and connections** — concrete interfaces, edge behaviors, error propagation

Each level is confirmed before proceeding. The AI challenges assumptions and pushes for well-understood behaviors: "You're building a queue with flush semantics — this is a solved problem. Are you inventing from scratch?"

If `specs/ARCHITECTURE.md` is stale (>30 days old with commits since), the AI refuses to use it at face value — actively compares it to the actual codebase and presents specific drifts.

After architecture, audits for missing guardrails (linters, type checkers, security scanners, dead code detection, secret detection, slopsquatting protection). Proposes a Phase 0 to install anything missing.

Then breaks the spec into ~400 LOC phases, each independently verifiable and committable.

Architecture persists in `specs/ARCHITECTURE.md` across features. Reviewed and updated every `/plan` run.

---

### Build

#### `/build` — Execute phases with TDD, review, and red team

The orchestrator. Spawns a fresh sub-agent for each phase (preventing context degradation), enforces strict TDD, and runs a post-phase code review.

**Pre-flight:** Detects all ecosystems in the project, verifies existing tools work, audits for missing tools (per-language checklists for Go, TypeScript, Python, Rust, JVM — linting, security, vulnerability scanning, dead code, formatting, secret detection, slopsquatting). Pushes hard for installation.

**Per phase:**
1. Fresh sub-agent with clean context (RALPH pattern)
2. Strict TDD — red/green/refactor. Work rejected if tests aren't written first.
3. Automated verification — every check from the plan runs independently
4. Post-phase code review — separate agent checks spec compliance, correctness, architecture. Optional stronger model for review.

**After all phases:**
1. Integration sweep — verify pieces connect correctly
2. Red team — adversarial agent that tries to break correctness, concurrency, boundaries, tests, and security. Fixes obvious issues, asks about judgment calls, flags design concerns.
3. Learning capture — appends to `specs/LEARNINGS.md`

**Crash recovery:** On resume, detects uncommitted changes (keep/stash/discard), reconciles progress file with git history, presents structured handoff summary.

#### `/quick` — Small changes without ceremony

For well-understood changes that don't need the full workflow. Takes a one-line description. Still enforces TDD and runs a code review, but no spec, no plan, no phases. If the change grows past ~50 lines, stops and redirects to `/feature`.

#### `/patch` — More structure than /quick, less than /feature

For bug fixes and small features from a clear problem statement or GitHub issue.

#### `/debug` — Scientific method debugging

Hypothesis → test → narrow → fix, with persistent state in `specs/DEBUG.md`. If a session dies mid-debug, the next session reads the file and picks up — never re-tries eliminated hypotheses.

Integrates with browser devtools (via [Gasoline](https://github.com/anthropics/gasoline)) for web/UI debugging: console errors, network requests, screenshots, DOM state.

Circuit breaker: after 3 failed hypothesis cycles, stops and asks for more context.

---

### Ship

#### `/ship` — Create a PR with spec context

Auto-populates the PR description from the spec: why (intent), what (behavioral summary), key decisions, constraints, and a verification checklist. Reviewers get the "why" without reading the full spec. Links to the spec file for deep dives.

---

### Learn

#### `/retro` — Check predictions against reality

After a feature ships and has production data, go back to the spec's "how will we know it worked?" and check. Scores each prediction (hit/partial/miss/unknown), analyzes why misses happened (mechanism, measurement, or environment), and extracts generalizable lessons.

Pushes for numbers, not feelings. "I think it improved" gets challenged: "Do you have the actual number?"

Feeds forward: suggests changes to `/feature`, `/plan`, or `/build` if the retro reveals a pattern.

---

### Support

#### `/note` — Zero-friction todo capture

`/note this module needs refactoring` — appends a timestamped item to `specs/TODO.md`. `/note` shows the list. `/note done 3` marks item 3 complete. `/note clear` removes completed items.

#### `/pause` — Structured session handoff

Captures everything the next session needs: what was running, what's done, what's next, key decisions, gotchas, git state, active files. Writes `specs/HANDOFF.md`. No questions asked — reads the conversation context and gets out.

#### `/resume` — Restore from handoff

Reads `specs/HANDOFF.md`, checks what changed since the pause, presents a structured briefing, waits for confirmation. Integrates with `/build`'s crash recovery.

---

## Persistent documents

These files accumulate project knowledge across features and sessions:

| File | Purpose | Updated by |
|------|---------|-----------|
| `specs/ARCHITECTURE.md` | System + subsystem + patterns + external connections | `/explore`, `/plan` |
| `specs/DECISIONS.md` | Append-only design decision register | `/feature`, `/refine` |
| `specs/LEARNINGS.md` | What surprised us, what went wrong, patterns | `/build`, `/debug`, `/retro` |
| `specs/TODO.md` | Scratchpad for ideas and tasks | `/note` |
| `specs/HANDOFF.md` | Session continuity | `/pause` |
| `specs/DEBUG.md` | Active debug session state | `/debug` |

---

## Install

Copy the commands into any project:

```bash
git clone https://github.com/brennhill/upfront.git
cp -r upfront/.claude/commands/ your-project/.claude/commands/
```

Or install globally (available in every project):

```bash
cp -r upfront/.claude/commands/ ~/.claude/commands/
```

That's it. No dependencies, no build step, no API keys. The commands are markdown files that Claude Code reads as instructions.

---

## Backed by research

Upfront's design is grounded in empirical research on AI-assisted software development:

### The skill erosion problem

- **Harvard/BCG (2023):** 758 consultants using GPT-4 — quality jumped 40% inside the AI's capability frontier, but dropped 19 percentage points *below the no-AI group* on tasks outside it. The variable was judgment, and judgment atrophies without exercise. (Dell'Acqua et al., "Navigating the Jagged Technological Frontier")

- **DORA (2025):** AI adoption increases throughput but also increases delivery instability. Teams shipping faster are simultaneously shipping less reliably. Speed without verification is not a gain. (DORA State of DevOps Report 2025)

- **Microsoft Research (2024):** Developers using AI assistants experienced a "false sense of confidence" — they believed their code was more secure than it actually was. The AI-generated code had comparable vulnerability rates but the developers reviewed it less carefully. (Perry et al.)

### The review collapse

- **SmartBear/Cisco:** Review effectiveness collapses above 400 lines of code. Reviewers stop finding defects when diffs are too large. This is why `/plan` targets ~400 LOC phases — it's the empirical limit of human attention. (SmartBear Code Review Study)

- **Faros AI (2024):** 10,000 developers across 1,255 teams — PR volume increased 98% after AI adoption, but net throughput (accepted, non-reverted changes) showed zero improvement. The volume increase was absorbed by review overhead and rework.

- **GitClear (2024):** Across 211 million changed lines, code churn doubled while refactoring collapsed. AI generates more code but less of it survives contact with production.

### Why forcing functions work

- **Nagappan & Ball (2005):** Code churn is a defect predictor with 89% accuracy. Rework rate — the metric `/feature` ties to success criteria — is one of the strongest signals of software quality. (ICSE 2005)

- **Capers Jones:** Defect removal efficiency (DRE) above 95% is the threshold for adequate quality, across 12,000+ projects. The Verification Triangle in the Delivery Gap framework uses this as the theoretical anchor.

- **Mantyla & Lassenius (2009):** 75% of defects found in code review are evolvability issues (structure, clarity, maintainability), not functional bugs. This is why `/build`'s review checks architecture and spec compliance, not just correctness.

### The specification gap

- **Montgomery et al.:** Systematic mapping of requirements quality research shows that ambiguous requirements are the single largest source of downstream defects. Clear intent — what Upfront's `/feature` forces — reduces rework more than any other intervention.

- **Jellyfish (2025):** 60% of engineering leaders cite "lack of clear metrics" as their biggest challenge with AI adoption. Only 20% measure AI's actual impact on delivery outcomes. The audit trail and `/retro` feedback loop exist to close this gap.

---

## Philosophy

Upfront is built on one belief: **the thinking is the product, not the code.**

AI can generate code. It cannot generate judgment. Every command in this toolkit exists to protect and exercise human judgment — the one thing that doesn't come back once it's gone.

The spec is not the point. The thinking the spec forces is the point.

---

## License

AGPL-3.0. See [LICENSE](LICENSE).

For commercial licensing (proprietary use, SaaS embedding, or redistribution without AGPL obligations), contact [brenn@thedeliverygap.com](mailto:brenn@thedeliverygap.com).

## Related

- [The Delivery Gap](https://thedeliverygap.com) — the book this toolkit accompanies
- [Delivery Gap Toolkit](https://github.com/brennhill/Delivery-Gap-Toolkit) — verification infrastructure (gates, policies, measurement)
- [sloppy-joe](https://github.com/brennhill/sloppy-joe) — slopsquatting protection for AI-generated package names
