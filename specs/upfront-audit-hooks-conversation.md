# /feature Conversation: Upfront Audit Hooks

> Full conversation transcript from `/feature` session on 2026-03-28.

---

## Phase 1: Intent

### 1. What problem does this solve?

**Q:** What problem goes away if this exists?

**User:** The fundamental problem is AI-written specs are unthinking. We need to FORCE humans to think. Many will just be lazy. And we need an audit trail. Use cases:
- Track which features went through `/feature` vs were built ad-hoc
- Measure how often the full workflow (`/feature` → `/plan` → `/build`) is followed vs skipped
- Pair with PostToolUse to capture duration per phase
- Feed into the book's argument about verification culture — teams that use the process vs teams that skip it

**Challenge:** An audit trail that logs "skill X was invoked at time Y" proves the command ran, not that thinking happened. Someone can `/feature` their way through all four phases saying "yeah sure" to everything and the hook logs it as a success. What specifically does the audit trail need to capture to distinguish genuine thinking from rubber-stamping?

**User:** We need to know if they answer. The skill already writes a summary at the end. Ideally at each phase we post up the summary to some audit log that can be configured.

**Refined problem:** Teams adopt the `/feature` workflow but there's no way to know if they're actually using it, skipping it, or rubber-stamping through it. Without evidence, you can't distinguish between a team that thinks through every feature and one that types `/build` directly.

### 2. How will we know it worked?

**Q:** What existing metric moves?

**User initial answer:** Managers can see how many features are handled upfront with thinking and which aren't, what phases are skipped, etc., and can go reconstruct who is slacking off.

**Challenge:** That's the use case, not a metric. "Managers can see" is a capability. What number moves?

**User:** Three metrics:
1. Percentage of merged PRs that have a corresponding `/feature` spec
2. Rework rate difference between spec'd and unspec'd features
3. Number of features that completed all 4 phases vs bailed early

**Challenge:** Metric #1 is gameable — run `/feature`, rubber-stamp every phase, get the checkmark. Do you care about that?

**User:** Rework rate also matters — that's the real signal. Spec'd features should have lower rework rates than unspec'd ones. If they don't, either the specs are rubber-stamped or the process isn't working.

**Resolution:** All three metrics work together — adoption (are people using it), effectiveness (does it reduce rework), depth (are they doing the full thinking or bailing early). If adoption is high but effectiveness is flat, specs are rubber-stamped. If effectiveness is high but adoption is low, only disciplined engineers use it.

### 3. What is out of scope?

**User:**
- Doesn't need to record the whole spec, just the title and the summary of thinking receipts. Manager goes to the spec file for details.
- Doesn't need a big dashboard, at least for now.

**Follow-up Q:** Does this enforce anything or is it purely observational? Per-project or centralized?

**User:** Pure observation. Per-project. Centralization should happen at the collector level.

### 4. What must NOT happen?

**User:** Actively harmful would be leaking sensitive info under AI Act or similar. We want to provide the audit trail but leave it to central/corporate other controls to do all that stuff.

**Challenge:** The thinking records from `/feature` could contain business-sensitive details — "we're solving the problem of losing enterprise customers to Competitor X" or specific revenue numbers. If the hook captures the full thinking record summary, that could leak into a log file that's less protected than the spec itself. Where do you draw the line?

**User:** Capture as-is. In a corporate env that's good enough. It keeps it limited enough. The full details would be hidden but trivial to find for those with appropriate access. Most companies know if they are winning or losing in abstract against competitors — the details and strategy in depth is what matters.

### 5. Pre-mortem findings

**Q:** It's six months from now and this feature failed. Why?

**User:**
1. People angry about being spied on
2. Bad data that's wrong so people have performance conversations based on it and it's bullshit — lack of reliability in numbers
3. People use it and see zero impact on rework or escape rates

**Challenge on #1:** "Reconstruct who is slacking off" IS surveillance. The surveillance framing and team health framing are in tension. How do you prevent that framing?

**User:** You could re-frame it as some sort of team health amplifier and audit. Makes sure people aren't just zoned out from exhaustion delegating to AI. This is marketing though.

**Challenge on #2:** What makes the data wrong?

**User:** Bad data is any sort of fuckup on the hooks where they do not report in a trustworthy way. Any substantial issues mean people will ignore all data, fair or unfair.

**Challenge on #3:** If teams adopt the full 4-phase process and rework doesn't drop, what do you conclude?

**User:** Hard to know. Thesis could be wrong. This is an acceptable risk. We're dogfooding right now and at least I notice a substantial difference.

**AI-added failure modes:**
- Hook fragility: Claude Code updates could change how tools are dispatched, breaking the hook silently
- Phase attribution: Running `/feature` twice for same feature produces duplicate/partial records
- Spec-to-PR linkage: Connecting "this PR was spec'd" requires a linking mechanism not yet built
- Log destination portability: Wrong format or non-configurable destination makes every corporate integration custom

**User:** Acceptable risks, we're at v1 beta.

---

## Phase 2: Behavioral Spec

### Level 1: Stories

**Q:** Walk me through what happens today when a manager wants to understand whether their team is actually thinking through features before building them.

**User:** You get a PRD and review it and try to pick it apart yourself or delegate to a staff engineer. Many times the PRD/RFC is an insane mess and really you just need it re-done entirely.

**Q:** Now walk me through what you want the experience to be instead.

**User:** They still probably look at the spec but BEFORE the spec they could go look and see if there is a logged decision audit trail with the summary. If there is, they at least know there was AI pushback. The odds of a complete mess are low OR they can see they skipped everything and just immediately tell them to follow process.

**Resolution:** The flow is: Manager gets notified → checks audit log (30 seconds) → sees one of three outcomes (all phases complete, partial, missing) → this is the triage layer before reading the full spec.

### Level 2: Mechanism

**Q:** Why do you think logging phase completion events with thinking summaries will actually solve this? Walk me through the causal chain.

**User:** The audit log and "teams think harder" are loosely connected. Teams think harder because the tool challenges them. They use the tool because they know if they don't, it's obvious they are slacking off really fast.

**Resolution:** Causal chain is: `/feature` forces thinking (tool's job) → audit log makes skipping visible (hook's job) → visibility creates social pressure → pressure drives tool use → tool forces thinking. The audit log doesn't make people think harder — `/feature` does that. The audit log closes the accountability loop.

**Challenge:** This only works if managers actually check the log. If nobody looks, social pressure disappears.

**User:** Out of scope, that's another layer of infra. Someone has to at least try to do their job. We can just make it easier.

### Level 3: States and transitions

**Q:** What are the states this system can be in?

**User:**
- Not installed
- Installed
- Not triggered and ignored
- Triggered
- Triggered and abandoned halfway through
- Triggered and completed
- In addition: able to transmit data / unable to transmit data (network down, outside VPN)

**Follow-up Q:** How does the hook know "abandoned halfway" vs "still thinking, just slow"?

**User:** We should post at the end of each phase once there is a thinking log. There should be a few events for each `/feature` run.

**Resolution:** Hook fires at each phase transition when thinking record is written. Complete run = 4 events. Count = how far they got. No timeout needed.

**Follow-up Q:** What about `/plan` and `/build`?

**User:** Out of scope but yes we should audit them as well — those are doing a different fundamental task.

### Level 4: Error cases and edges

**Q:** What breaks?

**User:**
- Could fail to hook entirely — can we hook multiple times in a skill?
- Could fail to post — network, etc. Should log and have it as "pending" and send next time the hook fires (send all in queue)
- Could just be broken — nothing works, nothing written to disk, nothing sent. All different failure modes
- For consistency we should probably add to disk queue, then empty queue in a post. Creates durable log

**Confirmed:** Yes, hooks fire on every PostToolUse event, so multiple times within a single `/feature` run.

**AI-added failure modes:**
- Context resets: User hits context limit mid-`/feature`, starts new session. Two partial trails for one feature.
- Multiple `/feature` runs: User reruns, two complete trails for one feature.
- Disk full / permissions: Queue can't be written, silent failure.
- Log rotation / size: Queue grows unbounded if remote never configured.

**User responses:**
- Context resets / multiple runs: Log both, reader interprets. Multiple runs are useful signal.
- Disk full: Development stops anyway, throw an error. Plus manager sees zero activity so it still works.
- Log rotation: Throw away after 90 days. If nobody notices after 3 months, no adults are in the room anyway.

---

## Phase 3: Design Conversation

### Codebase research

Explored the Delivery-Gap-Toolkit repo. Key findings:
- `quality-correctness-gates/agent-monitoring/README.md` defines a structured JSON trace format (session_id, timestamp, action_type, target, etc.) that maps to Langfuse, Arize Phoenix, Helicone, Portkey
- No existing hook infrastructure or executable logging code
- Working Python tooling exists (multi-pass review CLI) — precedent for shipping executable code
- The agent-monitoring trace format is close to OpenTelemetry span conventions

### Options presented

**Option A:** Pure hook + local JSONL file. Zero dependencies. No remote.

**Option B:** Hook + disk queue + configurable remote. What user described in Phase 2.

**Option C:** Hook + existing agent-monitoring trace format + remote queue. Reuse the toolkit's own format so audit events plug into the same tools (Langfuse, etc.) teams already use for agent monitoring.

**Q from user:** Does Option C work with existing commercially/open source available tools?

**Answer:** Yes — the trace format maps directly to what Langfuse (OSS), Arize Phoenix (OSS), Helicone (OSS), and Portkey all consume. Teams don't need a custom dashboard — Langfuse already has trace visualization, filtering, and team views.

**User chose:** Option C.

### Language decision

**Q from user:** Why not Go or Rust instead of Python/shell?

**Analysis:**
- Shell: Zero deps but JSON parsing in bash is brittle. jq is a dependency.
- Python: Works but Python version management is a mess. "Which python?" kills adoption.
- Go: Single binary, no runtime, cross-compiles. Fast. Toolkit doesn't have Go yet.
- Rust: Same single-binary advantages. User already knows it from rtk. Heavier to build for what's fundamentally parse-append-POST.

**Key question:** Who installs this? If it's distributed to teams, single binary with zero runtime deps makes adoption frictionless.

**User chose:** Go. Fast to write, single binary, broad ecosystem familiarity.

---

## Phase 4: Implementation Design

### Placement decision

Three options presented:
- `quality-correctness-gates/upfront-audit/` — alongside monitoring infra
- Top-level `upfront/` — clean extraction to own repo later
- `specs-and-process/hooks/` — alongside the commands it audits

**Structural concern flagged:** Monitoring guidance is in `quality-correctness-gates/` but the commands being monitored are in `.claude/commands/`. Putting the hook with monitoring is logically correct but physically distant. Putting it with commands is physically close but conceptually wrong.

**User chose:** Top-level `upfront/`. Clean extraction path.

### Parsing approach debate

**Issue raised:** The skill outputs free-form markdown. Parsing "Thinking Record: Intent" sections from tool_output is regex on LLM output — inherently fragile.

**Alternative proposed:** Have `/feature` write a structured JSON sidecar at each phase transition. Hook reads the sidecar.

**User pushback:** The `/feature` command is a prompt — it's instructions to an LLM. Writing a sidecar is a suggestion to the AI, not deterministic.

**Alternative proposed:** Have `/feature` emit structured JSON comment tags (`<!-- UPFRONT:{...} -->`) in the output.

**User pushback (same reason):** Still a suggestion to the LLM. Not deterministic.

**Q from user:** What about just using tool_output? It's a lot of text but so what?

**Resolution:** The thinking record has a consistent structure because the `/feature` prompt is explicit about the format. Even if formatting drifts slightly, we're capturing summaries for human consumption, not parsing into database fields. Grab everything between `### Thinking Record:` and the next `---` or `##`. 95% accuracy is fine. The fragility concern was overblown.

**Final approach:** Parse tool_output markdown with regex. Number of thinking records found = phases completed. Summary text captured as-is.

### Blind spots discussion

**Q:** What will the AI get wrong when it builds this?

**User:** Probably architecture and what data it records. If we can make the architecture clean, and vet that it really does consistently get the right data, and really does post it without crashing, memory leaks, etc., we're good. The required stuff is table stakes.

**AI-added blind spots:**
- Hook stdin schema is undocumented enough that the struct could be wrong — validate against real hook output before hardcoding
- File locking: O_APPEND atomic on POSIX for individual writes; flush needs rename-and-swap
- Parse failures: Log raw event with "parse_error" flag, don't silently drop
- Config file may contain remote endpoint URLs with tokens — document .gitignore

### Rollback and ownership

- **Signal:** Tool fail logs but no telemetry. Dogfooding is the detection mechanism.
- **Rollback:** Uninstall binary, remove hook from settings.json.
- **Owner:** Brenn Hill.
