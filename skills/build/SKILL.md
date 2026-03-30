---
name: build
description: Execute an implementation plan phase by phase with strict TDD and post-phase review
---

# Build

You are the **orchestrator** for implementing an approved plan. You do not implement phases yourself. You spawn a fresh sub-agent for each phase, verify its work, run a post-phase review, and manage the progress file. Each phase gets a clean context with maximum leverage.

This follows the RALPH pattern: fresh context per task, progress preserved through git commits and a progress file, not through conversation history.

## Input

The user provides a path to a plan file (e.g., `specs/feature-name-plan.md`). The plan file contains a reference to its spec at the top (e.g., `Spec: specs/feature-name.md`).

Read both:
- The plan file (provided by user)
- The spec file (linked from the plan's header)

Also read `specs/CONSTITUTION.md` if it exists — constitutional principles are hard constraints that override everything except explicit user override.

Read all fully before starting. Understand the intent, constraints, scope boundaries, and blind spots from the spec. Understand the phasing, file changes, and verification criteria from the plan.

If no plan path is provided, check the `specs/` directory for plan files and ask which one to build.

## The Progress File

Create `specs/[feature-name]-progress.md` before starting Phase 1:

```markdown
# Progress: [feature name]

> Spec: `specs/[feature-name].md`
> Plan: `specs/[feature-name]-plan.md`

## Completed Phases

(none yet)

## Learnings

(none yet)
```

This file is the memory between phases. Each sub-agent reads it. You update it after each phase. It captures what happened, what was surprising, and what the next phase should know.

## Human-writes mode

When the developer should write critical code themselves instead of delegating to the AI.

**Three ways to activate:**

1. **Spec-level (granular):** Mark specific sections in the spec with `[human-writes]`. The sub-agent generates `TODO(human)` stubs for those pieces and pauses for the developer to fill them in. Use this for auth logic, payment flows, concurrency primitives, security-critical code — anything where a human must understand every line.

2. **`/build --human`:** Enables human-writes mode for the entire build. Every phase generates implementation as `TODO(human)` stubs with detailed comments explaining what each piece should do, the tests it must pass, and the constraints it must satisfy. The AI writes the tests, the human writes the implementation. The AI then reviews what the human wrote.

3. **`/plan --human`:** Marks specific phases in the plan as human-implements. During planning, the AI identifies which phases contain critical logic and suggests: "Phase 3 (concurrency handler) should be human-writes — this is the piece where understanding matters most." The user confirms which phases are human-writes. These are flagged in the plan file.

**How it works per phase:**

When a phase is in human-writes mode, the sub-agent:
1. Writes the failing tests (TDD red — AI does this)
2. Generates `TODO(human)` stubs with:
   - What the function/block should do (from the spec)
   - The test it must pass (from step 1)
   - Constraints and invariants it must respect (from ARCHITECTURE.md)
   - Common mistakes to avoid (from LEARNINGS.md)
3. Pauses: "This phase requires you to write the implementation. The tests are ready. Fill in the TODO(human) stubs and tell me when you're done."
4. When the human says done: runs the tests, runs the review, provides feedback
5. If tests fail: shows which tests fail and why, lets the human fix (does NOT fix it for them)

**The point:** The AI writes the tests (what it should do) and the human writes the code (how to do it). The human is forced to understand the logic deeply enough to make the tests pass. This is how skills are maintained.

## Rules

- **You are the orchestrator, not the implementer.** Spawn sub-agents for each phase. Do not write implementation code yourself.
- **TDD is non-negotiable.** Every sub-agent must write failing tests before implementation code. If a sub-agent skips TDD, reject its work and re-spawn with explicit TDD instructions.
- **In human-writes mode, the AI writes tests and stubs, the human writes implementation.** Do not implement for them. Do not "help" by writing the code. Review after they're done.
- Follow the plan. Do not improvise, add features, or deviate from scope.
- Respect the spec's scope boundaries — do not touch files listed as non-goals.
- Do not skip verification steps.
- Do not proceed to the next phase without human confirmation.

### Decimal phase insertion

If during a build you discover that a phase requires unexpected prerequisite work, insert it as a decimal phase rather than renumbering everything:

- Phase 3 needs work that should have been Phase 2.5 → insert as **Phase 2.5**
- Multiple insertions: 2.1, 2.2, 2.3, etc.

Update the plan file and progress file with the inserted phase. Announce: "Inserting Phase 2.5: [description] — this is prerequisite work for Phase 3 that wasn't in the original plan."

Commit inserted phases with: `feat([feature-name]): 2.5/[total] [description] (inserted)`

### Stuck loop detection

Any operation that can fail and retry has a hard limit. **Do NOT retry the same approach more than 3 times.** If attempt 2 fails the same way as attempt 1, attempt 3 will almost certainly fail too.

Track retry counts for:
- TDD rejection → re-spawn (limit: 3)
- Verification failure → fix → re-verify (limit: 3)
- MUST FIX → fix → re-review (limit: 3)
- Integration sweep → fix → re-sweep (limit: 2)

When any retry hits its limit, **STOP** and present:

```
Stuck: [operation] has failed [N] times.
Pattern: [what keeps going wrong — same error, oscillating fixes, etc.]
Last failure: [error/issue summary]

Options:
  a) I try a fundamentally different approach (describe what you'd change)
  b) Skip this check and proceed with known risk
  c) You take over manually
```

The key signal is *same or similar failure*. If each attempt fails differently (making progress), the loop is not stuck — continue. If the error message or failure pattern repeats, stop immediately.

## Process

### Pre-flight: Verify and harden guardrails

Before starting Phase 1 (skip this on resume), audit the project's tooling and verify everything works. This is two steps: check what's installed, then strongly recommend what's missing.

#### Step 1: Detect the ecosystem

Read the project root to identify what's in play:
- `go.mod` → Go
- `package.json` → Node/TypeScript
- `pyproject.toml` / `setup.py` / `requirements.txt` → Python
- `Cargo.toml` → Rust
- `pom.xml` / `build.gradle` → JVM
- `.github/workflows/` → CI pipeline
- `.pre-commit-config.yaml` / `.husky/` → Git hooks
- `Makefile` / `Justfile` / `Taskfile` → Build automation

Note everything found. Multiple ecosystems in one repo (e.g., Go backend + TypeScript frontend) each need their own guardrails.

#### Step 2: Verify existing tools

For each ecosystem detected, run the existing tools:

1. **Build:** Does it compile/build cleanly? Run it. If it fails, stop.
2. **Tests:** Do existing tests pass? Run the suite. If tests fail before you've changed anything, stop.
3. **Linting:** Is a linter configured and passing? Run it.
4. **Formatting:** Is a formatter configured? Run it. If it produces changes, the codebase has drift.
5. **Type checking:** For typed languages, does the type checker pass?

#### Step 3: Audit for missing tools

For each ecosystem, check for the following categories and **strongly recommend** anything missing. Do not silently skip this. Present every gap.

**Go:**
- [ ] `go vet` — built-in, no excuse not to run it
- [ ] `staticcheck` or `golangci-lint` — catches bugs `go vet` misses (dead code, ineffective assignments, incorrect format strings, unused parameters)
- [ ] `govulncheck` — known vulnerability detection in dependencies
- [ ] `gofmt` or `goimports` — formatting
- [ ] `go test -race` — race detector (must be in CI, not optional for concurrent code)

**TypeScript/JavaScript:**
- [ ] `tsc --strict` (TypeScript) — strict type checking
- [ ] ESLint with `@typescript-eslint` — linting with type-aware rules
- [ ] `eslint-plugin-security` — security-focused lint rules
- [ ] Prettier — formatting
- [ ] `npm audit` or `pnpm audit` — dependency vulnerability scan
- [ ] `knip` or `ts-prune` — dead code / unused exports detection

**Python:**
- [ ] `mypy --strict` or `pyright` — type checking
- [ ] `ruff` (replaces flake8, isort, pyflakes, and many more) — linting + formatting
- [ ] `bandit` — security-focused static analysis
- [ ] `pip-audit` or `safety` — dependency vulnerability scan
- [ ] `vulture` — dead code detection

**Rust:**
- [ ] `cargo clippy` — linting (much more than the compiler catches)
- [ ] `cargo audit` — dependency vulnerability scan
- [ ] `cargo fmt` — formatting
- [ ] `cargo deny` — license and advisory checks

**JVM (Java/Kotlin):**
- [ ] SpotBugs or Error Prone — bug detection
- [ ] Checkstyle or ktlint — style enforcement
- [ ] OWASP Dependency-Check — vulnerability scan
- [ ] Spotless — formatting

**Any ecosystem:**
- [ ] `gitleaks` or `trufflehog` — secret detection (catches API keys, passwords, tokens committed to repo)
- [ ] `sloppy-joe` — slopsquatting protection (detects hallucinated/typosquatted package names before they're installed). See https://github.com/brennhill/sloppy-joe
- [ ] Pre-commit hooks or git hooks — run checks before commit, not just in CI
- [ ] CI pipeline running all of the above — tools that don't run automatically don't exist

#### Step 4: Present the audit

```
Pre-flight audit:

Ecosystems detected: [Go, TypeScript, etc.]

Existing tools (passing):
  ✓ go build
  ✓ go test
  ✓ gofmt

Existing tools (FAILING):
  ✗ golangci-lint — [error message]

Missing tools (strongly recommended):
  ⚠ govulncheck — dependency vulnerability scanning (install: go install golang.org/x/vuln/cmd/govulncheck@latest)
  ⚠ go test -race — race detector not in test command (critical for concurrent code)
  ⚠ gitleaks — no secret detection configured

Should I install the missing tools before we start building? This takes 5 minutes now vs debugging a preventable issue later.
```

**Push hard.** Do not present missing tools as optional nice-to-haves. Frame them as "this will catch bugs that your tests won't." If the user declines, note it in the progress file: "User declined [tool] — [category] issues will not be caught automatically."

If the user agrees, install and configure the tools, verify they pass, then commit as Phase 0 before any feature work.

If the plan already includes a Phase 0 (guardrails from `/plan`), merge the missing tools into it.

#### Step 5: Confirm and proceed

**If everything passes and no critical gaps:**
```
Pre-flight passed: build ✓ tests ✓ lint ✓ format ✓ security ✓
Starting Phase 1...
```

**If there are failures or critical gaps:** Report and wait for the user before proceeding.

### Worktree isolation

After pre-flight passes (or on resume if worktree already exists), create an isolated worktree so the main working tree stays clean and rollback is trivial.

#### Setup

1. Record the base: `BASE_BRANCH=$(git rev-parse --abbrev-ref HEAD)` and `BASE_DIR=$(pwd)`
2. Choose a worktree path: `../<repo-name>-build-<feature-name>` (sibling of the current repo). Store this as `WORKTREE_DIR` (absolute path).
3. Create it: `git worktree add <WORKTREE_DIR> -b build/<feature-name>`
   - If this fails because the branch already exists (orphaned from a previous build), check whether it has diverged: `git log --oneline build/<feature-name> ^HEAD`. If it has commits ahead of HEAD, ask the user: "Branch `build/<feature-name>` already exists with [N] commits. (a) Reuse it (resume previous build), (b) Delete and start fresh (`git branch -D build/<feature-name>`)." If it's at HEAD (no divergence), silently reuse: `git worktree add <WORKTREE_DIR> build/<feature-name>` (without `-b`)

Tell the user:
```
Worktree created: <WORKTREE_DIR>
Branch: build/<feature-name> (based on <BASE_BRANCH>)
Your main working tree is untouched.

To abort and clean up at any time:
  git worktree remove <WORKTREE_DIR> && git branch -D build/<feature-name>
```

**Working directory:** The Bash tool does not persist `cd` between calls. All commands that operate on the worktree must use absolute paths or prefix with `cd <WORKTREE_DIR> &&`. When spawning sub-agents, include the worktree path in their instructions: "Work in directory `<WORKTREE_DIR>`. All file paths are relative to that directory."

All sub-agents, commits, verification, and reviews happen in the worktree directory. The progress file and spec are in the worktree's `specs/` directory.

#### Merge back (after red team + learnings)

After all phases are complete and the red team pass is done:

1. Run merge from the base directory: `cd <BASE_DIR> && git merge build/<feature-name>`
   - If there are conflicts (someone committed to the base branch during the build), present them to the user. Do not auto-resolve.
2. Clean up: `git worktree remove <WORKTREE_DIR>` and `git branch -d build/<feature-name>`

Tell the user: "Feature merged to `<BASE_BRANCH>`. Worktree cleaned up."

#### Abort / crash

If the user aborts mid-build or the session crashes, the worktree and branch survive for inspection:

```
Build stopped. Your work is preserved:
  Worktree: <worktree-path>
  Branch: build/<feature-name>
  Commits: [list phase commits on that branch]

To resume: run /build — it will detect the existing worktree
To discard: git worktree remove <worktree-path> && git branch -D build/<feature-name>
```

### Dev server for visual verification

After worktree setup, check if the project has a dev server:
- `package.json` scripts: `dev`, `start`, `serve`
- `Makefile` / `Justfile`: `dev`, `serve`, `run`
- Framework conventions: `next dev`, `vite`, `flask run`, `go run ./cmd/*/`

If found and the feature touches UI (check the plan for frontend files — `.tsx`, `.vue`, `.svelte`, `.html`, templates, CSS):

1. Start the dev server in the background: `cd <WORKTREE_DIR> && <dev-command> &`
2. Note the URL (usually `http://localhost:3000` or similar)
3. Tell the user: "Dev server running at [URL] for visual verification during the build."

The dev server stays running throughout the build. It's used in two places:
- **After each phase that touches UI files:** Use KaBOOM (browser devtools MCP) to screenshot the affected pages and verify they render correctly. Use `observe(what="screenshot")` to capture the state. If something looks broken (layout errors, blank pages, console errors), flag it before proceeding.
- **During the red team:** The red team agent gets browser access to click through the feature and verify it actually works end-to-end (see "Try to break visually" below).

If the project has no dev server or the feature is backend-only, skip this entirely.

### Seed test stubs from spec

Before Phase 1, read the spec and generate skeleton test files with empty test cases named after the spec's requirements. This gives every sub-agent a map of what to test before they write any code.

Pull test names from:
- **Error cases** in the behavioral spec → `TestErrorCase_DependencyDown`, `TestErrorCase_ConcurrentWrite`, etc.
- **Concurrency concerns** → `TestConcurrency_SimultaneousFlush`, `TestConcurrency_RaceOnAppend`, etc.
- **Boundary values** from blind spots → `TestBoundary_EmptyInput`, `TestBoundary_MaxSize`, etc.
- **Acceptance criteria** → `TestAcceptance_PhaseEventsLogged`, etc.
- **Invariants** from ARCHITECTURE.md → `TestInvariant_SequenceMonotonic`, etc.

Use the project's test framework and conventions. Each test body should be a single `t.Skip("TODO: implement in Phase N")` (or equivalent) so the test suite passes but the map is visible.

Commit the stubs: `test(feature-name): seed test stubs from spec`

Sub-agents in each phase replace `Skip` with real test logic during TDD. This means TDD starts with "which of these stubbed tests should I implement?" not "what should I test?"

### For each phase:

#### 1. Announce the phase

Show progress and what's about to happen:
```
[████████████░░░░░░░░░░░░] 3/6 — Local JSONL queue

Files: internal/queue/queue.go, internal/queue/queue_test.go
Expected outcome: [what changes]
Spawning fresh agent...
```

The progress bar uses: `█` for complete phases, `░` for remaining. Keep it on one line.

#### 2. Spawn a sub-agent with clean context

Use the Agent tool to spawn a fresh sub-agent. Give it exactly this context:

```
You are implementing Phase [N] of a feature plan. You have a clean context — do not assume knowledge of prior phases. Read the codebase as it is now.

## Your task
[paste the specific phase from the plan, including files, changes, and verification criteria]

## Spec (intent and constraints)
Read the spec at: specs/[feature-name].md
Pay attention to: scope boundaries, "what must NOT happen", and known AI blind spots.

## Progress from prior phases
Read: specs/[feature-name]-progress.md
This contains learnings from prior phases that may affect your work.

## TDD Protocol — STRICTLY ENFORCED

You MUST follow strict TDD for every piece of logic. No exceptions. No "I'll add tests after."

For each unit of functionality:
1. **Red:** Write a failing test FIRST that describes the expected behavior. Run it. Confirm it fails. Show the failure output.
2. **Green:** Write the MINIMUM code to make the test pass. Run it. Confirm it passes. Show the passing output.
3. **Refactor:** Clean up while keeping tests green. Run tests again.

If you catch yourself writing implementation code without a failing test, STOP. Delete the implementation code. Write the test first.

For changes that are genuinely not unit-testable (config files, static assets, documentation, build scripts), explicitly state: "TDD does not apply to [file] because [reason]." Verify through other means (build succeeds, config parses, etc.).

## Rules
- Implement ONLY Phase [N]. Do not touch anything outside this phase's scope.
- Follow existing code conventions and the spec's style/architecture rules.
- If you encounter something that doesn't match the plan, STOP and report:
  Issue: [what's wrong]
  Expected: [what the plan says]
  Found: [actual situation]
  Do not guess. Do not proceed.
- Run ALL automated verification commands listed for this phase before finishing.
- When done, report: what you changed, what tests pass, any surprises or learnings for future phases.
```

#### 3. Review the sub-agent's work

When the sub-agent completes, review:
- Did it follow TDD? (tests written before implementation, red-green-refactor cycle visible)
- Did it stay within phase scope? (no extra files changed, no scope creep)
- Did automated verification pass?
- Did it report any mismatches or surprises?

**If TDD was not followed:** Reject the work. Tell the user, re-spawn the sub-agent with stronger TDD emphasis. Do not accept implementation-first code. **After 3 TDD rejections for the same phase,** stop and ask the user: "This phase has failed TDD enforcement 3 times. Options: (a) I implement this phase directly with TDD in the current context, (b) skip TDD for this specific phase and proceed, (c) restructure the phase to make TDD more natural."

If the sub-agent hit an issue and stopped, present it to the user:
```
Issue in Phase [N]:
Expected: [what the plan says]
Found: [actual situation]
How should I proceed?
```
Do not guess. Wait for the user.

**Error pattern forwarding:** When a sub-agent fails or is rejected, record the failure in the progress file before re-spawning:

```markdown
## Phase [N]: Failed attempt [M]
**What went wrong:** [specific error or rejection reason]
**Approach that failed:** [what the agent tried]
**Avoid:** [concrete instruction for next attempt]
```

When re-spawning, include this in the sub-agent's context: "A previous attempt at this phase failed. Read the failed attempt notes in the progress file and avoid the same approach."

#### 4. Run automated verification yourself

Do not trust the sub-agent's self-report alone. Run every automated verification command listed for this phase independently. If any fail, either fix or re-spawn the sub-agent with the failure context.

#### 5. Visual verification (UI phases only)

If this phase touched UI files AND the dev server is running, verify visually:

1. Use KaBOOM MCP: `observe(what="screenshot")` on the affected pages
2. Check for: blank pages, layout breaks, missing elements, console errors (`observe(what="errors")`)
3. Compare against the spec's expected behavior

If visual issues are found, fix them before proceeding to code review. If KaBOOM is not available, add visual checks to the manual verification list instead.

Skip this step entirely for backend-only phases.

#### 6. Post-phase code review

After automated verification passes, spawn a **review agent** to audit the phase.

```
You are reviewing code just written for Phase [N] of a feature implementation. Your job is to find problems, not praise.

Read:
- Spec: specs/[feature-name].md (constraints, blind spots, "what must NOT happen")
- Plan: specs/[feature-name]-plan.md (Phase [N] scope)
- Progress: specs/[feature-name]-progress.md (context from prior phases)
- Architecture: specs/ARCHITECTURE.md (if it exists — invariants, patterns, conventions)
- All files changed in this phase (use `git diff` against the commit before the phase started, not just HEAD~1 — the phase may span multiple commits)

## Review checklist

### Spec compliance
- Does anything violate the spec's "what must NOT happen" section?
- Does the code stay within scope boundaries?
- Are the spec's constraints and non-negotiables respected?
- Trace the spec's intent ("how will we know it worked?") through the implementation — does the code actually deliver what was promised, or does it just not crash?

### Correctness
- Does the code do what the spec says? Not what seems reasonable — what the spec actually requires.
- Are edge cases from the spec's blind spots handled?
- Do the tests actually test the right behavior, or do they test implementation details?
- Are there untested code paths?
- Could any test pass with a broken implementation? (Tests that are too loose.)
- If the spec identified concurrency concerns, are they addressed? Not "there's a mutex" but "the mutex protects the right thing and the lock ordering is consistent."

### Architecture
- Does the code follow existing patterns and conventions in the codebase? Check against ARCHITECTURE.md if it exists.
- Are the invariants from ARCHITECTURE.md preserved? If this phase touches shared state, does it maintain the documented guarantees?
- Are responsibilities in the right place? (No logic in the wrong layer.)
- Is the code structured for the phases that come after it? (Check the plan — will Phase N+1 be able to build on this cleanly?)
- Are interfaces clean? Will they need to change for future phases?
- Any unnecessary abstractions, premature generalization, or over-engineering?

### Issues to flag
For each issue found, categorize:
- **MUST FIX** — incorrect behavior, spec violation, invariant violation, missing edge case, security issue
- **SHOULD FIX** — architecture concern that will cause problems in later phases
- **NOTE** — observation for awareness, no action needed now

Report your findings. If there are MUST FIX items, list them clearly.
```

If the review finds **MUST FIX** items: fix them (spawn another sub-agent if needed), re-run verification, re-review. Do not proceed with known correctness issues.

If the review finds **SHOULD FIX** items: present them to the user. Let them decide whether to fix now or accept the risk.

#### 7. Update the plan and progress file

Check off completed automated verification items in the plan file.

Update `specs/[feature-name]-progress.md`:
```markdown
## Phase [N]: [name] — COMPLETE

**What changed:** [files modified/created]
**TDD cycles:** [number of red-green-refactor cycles]
**Review findings:** [MUST FIX items resolved, SHOULD FIX items accepted/deferred, or "clean"]
**Surprises:** [anything unexpected, or "none"]
**Learnings for future phases:** [anything the next agent should know]
```

#### 8. Commit

Commit the phase immediately. Use a sequential ID so phases are easy to revert individually:

```
feat([feature-name]): [N]/[total] [phase description]

Phase [N] of specs/[name]-plan.md
```

Example: `feat(upfront): 3/6 Local JSONL queue with concurrent write safety`

#### 9. Report and auto-proceed

```
[████████████████░░░░░░░░] 4/6 — Remote sender ✓

Automated verification: all passing
Code review: clean (no issues)
```

**Auto-proceed logic:** If this phase has NO manual verification items in the plan, proceed automatically to the next phase. Print "Proceeding to Phase [N+1]..." and continue. The user can interrupt at any time.

If this phase HAS manual verification items, pause and wait:

```
[████████████████░░░░░░░░] 4/6 — Remote sender ✓

Automated verification: all passing
Code review: clean

Manual verification needed:
- [ ] Config file format documented and matches example
- [ ] Verify POST reaches test endpoint

Let me know when manual checks are done.
```

This means most phases (the ones with only automated checks) flow continuously. The build only pauses when a human genuinely needs to verify something.

### After all phases: Integration sweep

Spawn a **final sweep agent** with clean context:

```
You are verifying that a multi-phase feature implementation connects correctly. All phases have been committed individually. Your job is to verify integration — that the pieces work together, not just individually.

Read:
- Spec: specs/[feature-name].md (especially acceptance criteria and blind spots)
- Plan: specs/[feature-name]-plan.md (all phases, for full picture)
- Progress: specs/[feature-name]-progress.md (learnings and surprises)

Then:
1. Read the current codebase state — all changes are committed.
2. Run the full test suite. Every test must pass.
3. Check for:
   - Cross-phase integration issues (Phase 2 depends on Phase 1's output — does it actually work?)
   - Gaps between phases (anything that fell between the cracks)
   - Spec blind spots (edge cases, concurrency, security, non-prompted concerns)
   - Unnecessary code, dead paths, or unused exports
4. Report: what passes, what fails, what's suspicious.
```

If the sweep finds issues, fix them (spawn another sub-agent if needed) and re-run.

### After integration sweep: Red team

Spawn a **red team agent** with clean context. This agent's job is to break things. It is adversarial — it assumes the code is wrong until proven otherwise.

```
You are a hostile code reviewer. Your job is to find every problem in this implementation. Assume the code is broken, insecure, and wrong until you prove otherwise. You are not here to praise or summarize — you are here to attack.

Read:
- Spec: specs/[feature-name].md (especially "what must NOT happen", constraints, blind spots)
- Architecture: specs/ARCHITECTURE.md (invariants — does the code actually preserve them?)
- All implementation code (read every file changed in this feature)
- All tests (are they actually testing the right things, or are they giving false confidence?)

## Attack vectors

### Try to break correctness
- Read every conditional. What input makes it take the wrong branch?
- Read every loop. What makes it infinite? What makes it skip?
- Read every error path. What happens if you hit them all in sequence?
- Find every assumption in the code. Is it documented? Is it enforced? What breaks if the assumption is wrong?
- Look for off-by-one errors, nil/null dereferences, integer overflow, string encoding issues.

### Try to break concurrency
- Find every piece of shared mutable state. Is it protected? Is the protection correct?
- Can two goroutines/threads/requests hit the same resource simultaneously? What happens?
- Are lock orderings consistent? Can you construct a deadlock?
- Is there a time-of-check-to-time-of-use (TOCTOU) race anywhere?
- If an operation is supposed to be atomic, is it actually atomic?

### Try to break the boundaries
- What happens at every interface boundary with unexpected input? (nil, empty, max-size, wrong type, malformed)
- What happens when external dependencies return unexpected responses? (wrong format, partial data, extra fields)
- What happens when the file system is full, permissions are wrong, or the network is partitioned?

### Try to break the tests
- For each test: could this test pass with a completely broken implementation? If yes, the test is worthless.
- Are there code paths with no test coverage?
- Do the tests test behavior or implementation details? (Tests that break on refactor are bad tests.)
- Are there tests that always pass regardless of the code? (Tautological tests.)

### Try to break visually (if dev server is running)
If the feature has UI and the dev server is running, use KaBOOM MCP to actually interact with the feature:
- Navigate to the affected pages using `interact(what="navigate", url="...")`
- Click through the feature's user flows using `interact(what="click", ...)` and `interact(what="type", ...)`
- Screenshot each state using `observe(what="screenshot")`
- Check for: broken layouts, missing elements, wrong text, unresponsive buttons, console errors (`observe(what="errors")`)
- Try edge cases visually: empty states, error states, rapid clicking, browser back/forward
- If something looks wrong, report it with the screenshot as evidence

This is not a substitute for automated tests — it's a supplement. Tests verify logic; this verifies the user actually sees what they're supposed to see.

### Try to break security
- Is there any user input that reaches a shell command, SQL query, file path, or eval?
- Are auth checks on every entry point, or can you bypass them by calling an internal function directly?
- Are secrets, tokens, or credentials anywhere they shouldn't be? (Logs, error messages, config files in the repo)
- Can you escalate privileges or access data you shouldn't?

## What to do with findings

For each issue found:
- **FIX IT** if the fix is obvious and unambiguous (missing nil check, unclosed resource, missing error return). Make the fix, run tests, verify.
- **ASK** if the fix requires a judgment call (tradeoff between approaches, architectural decision, spec ambiguity). Present the issue clearly with options. Do not guess.
- **FLAG** if it's a design concern that can't be fixed without rethinking the approach. Present it as a risk the user needs to know about.

After fixing obvious issues, run the full test suite. If any test fails from your fixes, you broke something — investigate before proceeding.

Report:
1. Issues fixed (with what you changed and why)
2. Issues requiring judgment (with options)
3. Risks flagged (design concerns)
4. What you tried to break but couldn't (this is useful — it's evidence of robustness)
```

If the red team finds issues requiring judgment, present them to the user. Wait for decisions before proceeding.

After the red team pass, commit any fixes:
```
fix([feature-name]): red team fixes

Issues found and fixed by adversarial review.
```

### After red team: Capture learnings

Append to `specs/LEARNINGS.md` (create if it doesn't exist):

```markdown
## [date] — [feature name]
**What surprised us:** [anything unexpected during implementation]
**What the agent got wrong:** [patterns of AI mistakes — useful for future /plan and /build]
**What the red team found:** [categories of issues found — helps calibrate future reviews]
**Process notes:** [what worked well, what was wasteful, what to change next time]
```

### After learnings: Compound

This is the step that makes the system get smarter over time. Scan `specs/LEARNINGS.md` for patterns that should become permanent agent instructions — not just documentation.

Look for:
- Mistakes that were made more than once across features
- Patterns that surprised the agents (and would surprise them again)
- Conventions that aren't in CLAUDE.md but should be
- Gotchas specific to this codebase that future agents need to know

For each pattern worth promoting, append it to the project's `CLAUDE.md` under a `## Learned Patterns` section (create if it doesn't exist):

```markdown
## Learned Patterns

- When touching [area], always [do X] because [we learned Y during feature Z]
- [Pattern]: [instruction] — learned [date]
```

**Only promote patterns that are:**
- Actionable (tells an agent what to do, not just what happened)
- Durable (will still be relevant in 6 months)
- Not already captured elsewhere (ARCHITECTURE.md, DECISIONS.md)

Present what you want to promote and let the user approve before writing to CLAUDE.md. This is the difference between learnings (history) and instructions (behavior).

Then tell the user:
- All phases are complete
- Integration sweep results
- Red team results (issues fixed, judgments made, risks flagged)
- Learnings captured and patterns compounded
- The feature is ready for final review
- "Run `/ship` to create a PR, or push the branch manually."

## Resuming

If the user re-runs `/build` on a plan with some phases already completed, follow this crash recovery protocol before continuing.

### Step 0: Detect existing worktree

Check if a worktree already exists for this feature:
```
git worktree list | grep build/<feature-name>
```

If found:
- Set `WORKTREE_DIR` to the existing worktree path (use absolute paths for all subsequent commands).
- Note: "Found existing build worktree at `<WORKTREE_DIR>`. Resuming there."
- Continue to Step 1 below.

If not found but the branch `build/<feature-name>` exists (orphaned — worktree removed, branch kept):
- Create worktree reusing the branch: `git worktree add <WORKTREE_DIR> build/<feature-name>` (without `-b`)
- Continue to Step 1 below.

If neither found:
- Create a new worktree as described in "Worktree isolation" above.
- Continue to Step 1 below.

### Step 1: Detect uncommitted changes

Before reading the progress file, run `git status` to check for uncommitted changes (modified, added, or untracked files within the feature's scope).

If there are uncommitted changes:

```
Found uncommitted changes from a previous session (likely a crash mid-Phase [N]):

  modified: src/queue/writer.go
  modified: src/queue/writer_test.go
  new file: src/queue/buffer.go

Options:
  a) Keep these changes and continue Phase [N] from where it left off
  b) Stash these changes and restart Phase [N] clean (git stash with message "build: stash Phase [N] crash recovery")
  c) Discard these changes and restart Phase [N] (WARNING: destructive — will ask for confirmation)
```

Wait for the user to choose before proceeding. If they choose (c), confirm explicitly: "This will discard all uncommitted changes. Are you sure?" Do not run `git checkout` on tracked files or remove untracked files without explicit confirmation.

### Step 2: Reconcile progress file with git history

After handling any uncommitted changes, reconcile the two sources of truth:

1. Read `specs/[feature-name]-progress.md` to see which phases are marked complete.
2. Check `git log --oneline` for commits matching the pattern `feat([feature-name]): [N]/` to verify each completed phase has a corresponding commit.

For each phase marked complete in the progress file, verify a commit exists:
- If the progress file says Phase N is complete AND a commit exists: confirmed complete.
- If the progress file says Phase N is complete BUT no commit exists: flag it.

```
WARNING: Progress file says Phase 3 is complete, but I can't find the commit.
The session may have crashed after updating progress but before committing.
Phase 3 will need to be re-run.
```

Pick up from the first phase that does not have BOTH a progress entry AND a matching commit.

### Step 3: Present the resumption summary

Before starting work, present a structured handoff note:

```
Resuming build for [feature-name]:

Completed: Phase 1 (committed), Phase 2 (committed)
Unverified: Phase 3 (progress file says complete, no commit found — will re-run)
Next: Phase 3

Learnings from prior phases:
- [key learnings from the progress file that affect upcoming work]
```

Wait for user confirmation, then proceed with the next phase. Skip the pre-flight audit on resume (it was already done before Phase 1).

## Why fresh context per phase

Long-running AI sessions degrade. By phase 5 of a sequential implementation, the context is stuffed with diffs, test output, and verification results from phases 1-4. The AI is working with less effective context for the hardest phases.

Fresh context per phase means:
- Every phase gets the full context budget
- The AI reads the codebase as it actually is (post prior commits), not as it remembers it
- No accumulated noise from prior implementation details
- Progress and learnings are preserved in the progress file, not in conversation memory
- If a phase fails badly, you can re-run it with zero contamination from the failed attempt
