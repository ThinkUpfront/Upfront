---
description: Use when the user wants to improve their CLAUDE.md or AGENTS.md — "my agent instructions suck", "set up CLAUDE.md", "improve my agents file", "how do I make AI work better in this repo". Audits instruction files and generates stack-specific examples.
user-invocable: true
---

# Enlighten

You are auditing and improving the project's AI instruction files (`CLAUDE.md`, `AGENTS.md`, or equivalent). Bad instruction files are the #1 reason AI agents produce bad code — vague instructions, missing commands, no boundaries. Your job is to fix that.

**Input:** $ARGUMENTS (optional — a specific area to focus on, or blank for full audit)

## Step 1: Find and read instruction files

Check for:
- `CLAUDE.md` (Claude Code)
- `AGENTS.md` (GitHub Copilot / generic)
- `.cursorrules` (Cursor)
- `.windsurfrules` (Windsurf)
- Any other AI instruction files in the root or `.claude/` directory

Also read:
- `specs/ARCHITECTURE.md` if it exists (for tech stack context)
- `package.json`, `go.mod`, `Cargo.toml`, `pyproject.toml`, `Gemfile`, `pom.xml` — whatever exists (to detect the tech stack)
- Recent git history (`git log --oneline -20`) for commit style

If NO instruction file exists, say so and offer to create one from scratch.

## Step 2: Detect the tech stack

From the files above, identify:
- Languages and versions
- Frameworks and versions
- Package manager
- Test runner
- Linter / formatter
- Build system
- CI/CD

## Step 2b: Load stack-specific references

Based on the detected stack, read the matching reference file(s) from the plugin's `references/` directory. These contain curated do/don't examples from the best style guides (Uber Go, Google Style Guides, Bulletproof React, Rust API Guidelines, Effective Java, etc.).

| Stack detected | Reference file |
|---|---|
| React, Vite, Next.js | `references/react-vite.md` |
| TypeScript (any project) | `references/typescript.md` |
| Go | `references/golang.md` |
| Python | `references/python.md` |
| Java, Spring, Gradle, Maven | `references/java.md` |
| Kotlin, Android | `references/kotlin.md` |
| Rust | `references/rust.md` |
| C#, .NET, ASP.NET | `references/dotnet.md` |
| PHP, Laravel | `references/php.md` |

Load ALL that match (e.g., a React+TypeScript project loads both `react-vite.md` and `typescript.md`). Use the examples from these files when generating stack-specific suggestions in Steps 4-5 — they're better than generic examples.

The reference files live at: `~/.claude/plugins/marketplaces/thinkupfront/plugin/references/`

## Step 3: Audit against the six core areas

Score the existing file (or note "missing" if no file exists) against these areas. Research from 2,500+ repositories shows these are the areas that separate files that work from files that don't.

### 1. Commands (build, test, lint)

**What to check:** Are exact commands with flags specified? Are they near the top of the file?

**Bad:**
```
Run tests before committing.
```

**Good:**
```
## Build & Test
npm run build
npm test -- --coverage
npm run lint
```

**Why it matters:** Agents need to know how to verify their work. "Run tests" is vague — `npm test -- --coverage` is executable.

### 2. Project structure

**What to check:** Are key directories and files explained? Does the agent know where to put new code?

**Bad:**
```
Follow the existing project structure.
```

**Good:**
```
## Structure
- src/api/ — Express route handlers, one file per resource
- src/services/ — Business logic, no HTTP awareness
- src/models/ — Sequelize models, one per table
- tests/ — Mirrors src/ structure, suffix .test.ts
```

### 3. Code style with examples

**What to check:** Are conventions shown with code, not just described? Good AND bad examples?

**Bad:**
```
Use consistent naming conventions.
```

**Good:**
```
## Style
Functions: camelCase
Components: PascalCase
Constants: UPPER_SNAKE_CASE

// Good
function fetchUserProfile(userId: string): Promise<User> { ... }

// Bad — don't abbreviate, don't use generic names
function getUP(id: any): any { ... }
```

**Why it matters:** One real code snippet showing your style beats three paragraphs describing it.

### 4. Testing conventions

**What to check:** Does the file specify how to write tests? What framework? What patterns?

**Stack-specific examples to suggest:**

**React/TypeScript:**
```
## Testing
- Framework: Vitest + React Testing Library
- Test behavior, not implementation
- No snapshot tests
- Mock external APIs, never mock internal modules

// Good
test('shows error when email is invalid', async () => {
  render(<SignupForm />)
  await userEvent.type(screen.getByLabelText('Email'), 'not-an-email')
  expect(screen.getByText('Invalid email')).toBeVisible()
})

// Bad — tests implementation details
test('calls setError with correct message', () => { ... })
```

**Go:**
```
## Testing
- go test ./... -race
- Table-driven tests preferred
- No mocking frameworks — use interfaces and test doubles
- Test files next to source: foo.go → foo_test.go
```

**Python:**
```
## Testing
- pytest -v
- Fixtures over setup/teardown
- No unittest.TestCase — use plain functions
```

### 5. Git workflow

**What to check:** Commit message format? Branch naming? PR process?

Detect from recent git history what convention is already in use and codify it.

### 6. Boundaries (three-tier)

**What to check:** Does the file have explicit always/ask-first/never rules?

**Suggest this structure:**
```
## Boundaries

### Always
- Run tests before committing
- Check for secrets before staging (gitleaks)
- Follow existing patterns in the directory you're editing

### Ask first
- Adding new dependencies
- Changing database schemas
- Modifying CI/CD pipelines
- Changing auth or security code

### Never
- Commit secrets, API keys, or credentials
- Modify .env files
- Delete or weaken linter rules
- Force push to main
```

## Step 4: Present the audit

For each of the six areas, report:
- **Present / Missing / Weak**
- What's there (if anything)
- What's missing
- A specific, stack-appropriate example to add

Don't dump everything at once. Present the audit summary, then ask: "Want me to walk through each area and add what's missing?"

## Step 5: Generate improvements

For each area the user wants to improve:

1. Pull concrete examples from the loaded reference files — use their do/don't patterns directly rather than inventing generic ones
2. Adapt the examples to the project's actual types, function names, and patterns where possible
3. Explain WHY it matters (one sentence)
4. Ask if they want to add it

Write the improvements directly to the file. Don't create a separate file — improve what exists or create the file if it doesn't exist.

## Step 6: Add lazy-loading references

Instruction files shouldn't contain everything — they should REFERENCE everything. Add `@` references to files that provide context on demand, so the AI loads them when relevant instead of stuffing the instruction file with content.

Suggest adding references like:
```
@specs/ARCHITECTURE.md
@specs/DECISIONS.md
@specs/LEARNINGS.md
```

For stack-specific standards, suggest creating focused rule files in `.claude/rules/` or equivalent:
- `.claude/rules/testing.md` — testing conventions with examples
- `.claude/rules/code-style.md` — naming, formatting, patterns
- `.claude/rules/security.md` — auth, input validation, secrets handling

These are loaded by Claude Code when relevant (via path globs) without bloating the main instruction file. The main CLAUDE.md stays lean — commands, structure, boundaries — and the detail lives in referenced files.

## Step 7: Re-run offer

After improvements are applied, offer: "Want to run `/upfront:enlighten` again after you've used the AI for a few sessions? The best instruction files grow through iteration — when the agent makes a mistake, add a rule to prevent it next time."

## Rules

- **Do NOT overwrite user content.** Add to the file, don't replace existing sections unless the user explicitly asks.
- **Stack-specific examples only.** Don't suggest React testing patterns for a Go project. Detect the stack first, then generate relevant examples.
- **Commands go first.** If restructuring the file, move build/test/lint commands to the top.
- **Show, don't tell.** Every suggestion must include a concrete example, not just a description of what to add.
- **Respect existing style.** If their CLAUDE.md uses bullet points, use bullet points. If it uses headers, use headers. Match their voice.
- **Three-tier boundaries are the highest-impact addition.** If you can only add one thing, add the always/ask-first/never section.
