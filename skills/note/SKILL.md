---
name: note
description: Zero-friction idea/todo capture. Append a note, view list, mark done, or clear completed.
---

# /note — Scratchpad Capture

Argument: $ARGUMENTS

## Rules

- No conversation, no questions. Just do it and confirm.
- File location: `specs/TODO.md` in the current project root.
- If `specs/TODO.md` doesn't exist, create it with a `# TODO` header first.

## Behavior

### No argument given
Read and display the current contents of `specs/TODO.md`. If the file doesn't exist, say "No TODOs yet."

### `done N` (e.g., `done 3`)
Find the Nth **unchecked** checkbox item (`- [ ]`) in `specs/TODO.md` and change it to `- [x]`. N counts only unchecked items — `done 1` marks the first open item, regardless of how many checked items precede it. Confirm which item was marked done.

### `clear`
Remove all lines matching `- [x]` from `specs/TODO.md`. Confirm how many items were cleared.

### Anything else
Treat the entire argument as a note. Append to `specs/TODO.md`:
```
- [ ] [YYYY-MM-DD] <note text>
```
Use today's date. Confirm with the note text.
