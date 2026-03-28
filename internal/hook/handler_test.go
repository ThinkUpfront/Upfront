package hook

import (
	"strings"
	"testing"
)

func TestParseInput_MalformedJSON(t *testing.T) {
	input := []byte(`this is not json`)

	_, err := ParseInput(input)
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestParseInput_ValidJSON(t *testing.T) {
	input := []byte(`{
		"session_id": "abc-123",
		"cwd": "/tmp/project",
		"tool_name": "Skill",
		"tool_input": {
			"skill_name": "feature",
			"args": "checkout flow"
		},
		"tool_response": "some output"
	}`)

	h, err := ParseInput(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.SessionID != "abc-123" {
		t.Errorf("SessionID = %q, want %q", h.SessionID, "abc-123")
	}
	if h.Cwd != "/tmp/project" {
		t.Errorf("Cwd = %q, want %q", h.Cwd, "/tmp/project")
	}
	if h.ToolName != "Skill" {
		t.Errorf("ToolName = %q, want %q", h.ToolName, "Skill")
	}
	if h.ToolInput.SkillName != "feature" {
		t.Errorf("SkillName = %q, want %q", h.ToolInput.SkillName, "feature")
	}
	if h.ToolInput.Args != "checkout flow" {
		t.Errorf("Args = %q, want %q", h.ToolInput.Args, "checkout flow")
	}
	if h.ToolResponse != "some output" {
		t.Errorf("ToolResponse = %q, want %q", h.ToolResponse, "some output")
	}
}

func TestExtractEvents_NonSkillTool(t *testing.T) {
	h := Input{
		SessionID:    "sess-1",
		Cwd:          "/tmp",
		ToolName:     "Bash",
		ToolResponse: "### Thinking Record: Intent\nsome text\n---",
	}

	events := ExtractEvents(&h)
	if len(events) != 0 {
		t.Errorf("expected 0 events for non-Skill tool, got %d", len(events))
	}
}

func TestExtractEvents_EmptyToolResponse(t *testing.T) {
	h := Input{
		SessionID:    "sess-2",
		Cwd:          "/tmp",
		ToolName:     "Skill",
		ToolInput:    ToolInput{SkillName: "feature"},
		ToolResponse: "",
	}

	events := ExtractEvents(&h)
	if len(events) != 0 {
		t.Errorf("expected 0 events for empty tool_response, got %d", len(events))
	}
}

func TestExtractEvents_NoThinkingRecords(t *testing.T) {
	h := Input{
		SessionID:    "sess-3",
		Cwd:          "/tmp",
		ToolName:     "Skill",
		ToolInput:    ToolInput{SkillName: "feature"},
		ToolResponse: "Some output without any thinking records.",
	}

	events := ExtractEvents(&h)
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestExtractEvents_OnePhase(t *testing.T) {
	h := Input{
		SessionID: "sess-4",
		Cwd:       "/tmp/project",
		ToolName:  "Skill",
		ToolInput: ToolInput{SkillName: "feature", Args: "checkout"},
		ToolResponse: `## Intent

### 1. What problem does this solve?
Some problem description.

### Thinking Record: Intent
**Decided:** Build a checkout audit trail.
**Reasoning:** The process forces thinking.
**Risks accepted:** Hook fragility.

---
`,
	}

	events := ExtractEvents(&h)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	e := events[0]
	if e.SessionID != "sess-4" {
		t.Errorf("SessionID = %q, want %q", e.SessionID, "sess-4")
	}
	if e.Phase != 1 {
		t.Errorf("Phase = %d, want 1", e.Phase)
	}
	if e.PhaseName != "Intent" {
		t.Errorf("PhaseName = %q, want %q", e.PhaseName, "Intent")
	}
	if e.Target != "/tmp/project" {
		t.Errorf("Target = %q, want %q", e.Target, "/tmp/project")
	}
	if e.FeatureName != "checkout" {
		t.Errorf("FeatureName = %q, want %q", e.FeatureName, "checkout")
	}
	if e.ThinkingSummary == "" {
		t.Error("ThinkingSummary should not be empty")
	}
	// Verify summary captures ALL lines, not just the first
	if !strings.Contains(e.ThinkingSummary, "Build a checkout audit trail") {
		t.Errorf("ThinkingSummary missing Decided line, got: %q", e.ThinkingSummary)
	}
	if !strings.Contains(e.ThinkingSummary, "The process forces thinking") {
		t.Errorf("ThinkingSummary missing Reasoning line, got: %q", e.ThinkingSummary)
	}
	if !strings.Contains(e.ThinkingSummary, "Hook fragility") {
		t.Errorf("ThinkingSummary missing Risks line, got: %q", e.ThinkingSummary)
	}
}

func TestExtractEvents_AllFourPhases(t *testing.T) {
	h := Input{
		SessionID: "sess-full",
		Cwd:       "/home/user/project",
		ToolName:  "Skill",
		ToolInput: ToolInput{SkillName: "feature", Args: "audit trail"},
		ToolResponse: `## Intent

### 1. What problem does this solve?
No audit trail exists.

### Thinking Record: Intent
**Decided:** Build audit trail for feature definitions.
**Reasoning:** Forces accountability.

---

## Behavioral Spec

### User Stories
Manager checks the audit log.

### Thinking Record: Behavioral Spec
**Decided:** Hook fires at each phase transition.
**Reasoning:** Phase-level granularity needed.

---

## Design Approach

### Conceptual approach
Hook fires on PostToolUse.

### Thinking Record: Design Approach
**Decided:** Option C — align with existing trace format.
**Reasoning:** Teams already using Langfuse get data alongside monitoring.

---

## Implementation Design

### Architecture
Go binary in top-level directory.

### Thinking Record: Implementation Design
**Decided:** Go binary with five CLI commands.
**Reasoning:** Go gives single-binary distribution.

---
`,
	}

	events := ExtractEvents(&h)
	if len(events) != 4 {
		t.Fatalf("expected 4 events, got %d", len(events))
	}

	expected := []struct {
		phase     int
		phaseName string
		decided   string
		reasoning string
	}{
		{1, "Intent", "Build audit trail", "Forces accountability"},
		{2, "Behavioral Spec", "Hook fires at each phase transition", "Phase-level granularity"},
		{3, "Design Approach", "Option C", "Teams already using Langfuse"},
		{4, "Implementation Design", "Go binary with five CLI commands", "Go gives single-binary"},
	}

	for i, want := range expected {
		e := events[i]
		if e.Phase != want.phase {
			t.Errorf("events[%d].Phase = %d, want %d", i, e.Phase, want.phase)
		}
		if e.PhaseName != want.phaseName {
			t.Errorf("events[%d].PhaseName = %q, want %q", i, e.PhaseName, want.phaseName)
		}
		if !strings.Contains(e.ThinkingSummary, want.decided) {
			t.Errorf("events[%d].ThinkingSummary missing Decided %q, got: %q", i, want.decided, e.ThinkingSummary)
		}
		if !strings.Contains(e.ThinkingSummary, want.reasoning) {
			t.Errorf("events[%d].ThinkingSummary missing Reasoning %q, got: %q", i, want.reasoning, e.ThinkingSummary)
		}
		if e.SessionID != "sess-full" {
			t.Errorf("events[%d].SessionID = %q, want %q", i, e.SessionID, "sess-full")
		}
		if e.Target != "/home/user/project" {
			t.Errorf("events[%d].Target = %q, want %q", i, e.Target, "/home/user/project")
		}
		if e.FeatureName != "audit trail" {
			t.Errorf("events[%d].FeatureName = %q, want %q", i, e.FeatureName, "audit trail")
		}
		if e.ActionType != "upfront_phase_complete" {
			t.Errorf("events[%d].ActionType = %q, want %q", i, e.ActionType, "upfront_phase_complete")
		}
		if e.AgentID != "upfront" {
			t.Errorf("events[%d].AgentID = %q, want %q", i, e.AgentID, "upfront")
		}
		if e.Result != "success" {
			t.Errorf("events[%d].Result = %q, want %q", i, e.Result, "success")
		}
	}
}

func TestExtractEvents_UnknownPhaseNameSkipped(t *testing.T) {
	h := Input{
		SessionID: "sess-unknown",
		Cwd:       "/tmp",
		ToolName:  "Skill",
		ToolInput: ToolInput{SkillName: "feature"},
		ToolResponse: `### Thinking Record: Unknown Phase
**Decided:** Something.

---

### Thinking Record: Intent
**Decided:** Real phase.

---
`,
	}

	events := ExtractEvents(&h)
	if len(events) != 1 {
		t.Fatalf("expected 1 event (unknown phase skipped), got %d", len(events))
	}
	if events[0].PhaseName != "Intent" {
		t.Errorf("PhaseName = %q, want %q", events[0].PhaseName, "Intent")
	}
}

func TestExtractEvents_TerminatedByNextHeading(t *testing.T) {
	h := Input{
		SessionID: "sess-heading",
		Cwd:       "/tmp",
		ToolName:  "Skill",
		ToolInput: ToolInput{SkillName: "feature"},
		ToolResponse: `### Thinking Record: Intent
**Decided:** Build it.
**Reasoning:** Because.

## Behavioral Spec

Some other content.
`,
	}

	events := ExtractEvents(&h)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if !strings.Contains(events[0].ThinkingSummary, "Build it") {
		t.Errorf("ThinkingSummary missing Decided line, got: %q", events[0].ThinkingSummary)
	}
	if !strings.Contains(events[0].ThinkingSummary, "Because") {
		t.Errorf("ThinkingSummary missing Reasoning line, got: %q", events[0].ThinkingSummary)
	}
	// Should NOT contain content after the ## heading
	if strings.Contains(events[0].ThinkingSummary, "Some other content") {
		t.Errorf("ThinkingSummary should not contain content after ## heading, got: %q", events[0].ThinkingSummary)
	}
}
