package format

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestNewEvent_Defaults(t *testing.T) {
	e := NewEvent("sess-123", "p1", 1, "Intent", "Decided to build X", "specs/foo.md", "foo")

	if e.AgentID != "upfront" {
		t.Errorf("AgentID = %q, want %q", e.AgentID, "upfront")
	}
	if e.ActionType != "upfront_phase_complete" {
		t.Errorf("ActionType = %q, want %q", e.ActionType, "upfront_phase_complete")
	}
	if e.SessionID != "sess-123" {
		t.Errorf("SessionID = %q, want %q", e.SessionID, "sess-123")
	}
	if e.Phase != 1 {
		t.Errorf("Phase = %d, want %d", e.Phase, 1)
	}
	if e.PhaseName != "Intent" {
		t.Errorf("PhaseName = %q, want %q", e.PhaseName, "Intent")
	}
	if e.PhasesTotal != DefaultPhasesTotal {
		t.Errorf("PhasesTotal = %d, want %d", e.PhasesTotal, DefaultPhasesTotal)
	}
	if e.ThinkingSummary != "Decided to build X" {
		t.Errorf("ThinkingSummary = %q, want %q", e.ThinkingSummary, "Decided to build X")
	}
	if e.Target != "specs/foo.md" {
		t.Errorf("Target = %q, want %q", e.Target, "specs/foo.md")
	}
	if e.FeatureName != "foo" {
		t.Errorf("FeatureName = %q, want %q", e.FeatureName, "foo")
	}
	if e.Result != "success" {
		t.Errorf("Result = %q, want %q", e.Result, "success")
	}
	if e.ActionDetail != "Phase 1: Intent" {
		t.Errorf("ActionDetail = %q, want %q", e.ActionDetail, "Phase 1: Intent")
	}
	if e.SkippedQuestions == nil {
		t.Error("SkippedQuestions should be non-nil empty slice")
	}
	if len(e.SkippedQuestions) != 0 {
		t.Errorf("SkippedQuestions length = %d, want 0", len(e.SkippedQuestions))
	}
	if e.ErrorMessage != "" {
		t.Errorf("ErrorMessage = %q, want empty", e.ErrorMessage)
	}
	if e.Timestamp == "" {
		t.Error("Timestamp should be set")
	}
}

func TestEvent_JSONRoundTrip(t *testing.T) {
	e := NewEvent("sess-456", "p1", 2, "Behavioral Spec", "Users want faster checkout", "specs/checkout.md", "checkout")

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded.SessionID != e.SessionID {
		t.Errorf("SessionID = %q, want %q", decoded.SessionID, e.SessionID)
	}
	if decoded.AgentID != e.AgentID {
		t.Errorf("AgentID = %q, want %q", decoded.AgentID, e.AgentID)
	}
	if decoded.ActionType != e.ActionType {
		t.Errorf("ActionType = %q, want %q", decoded.ActionType, e.ActionType)
	}
	if decoded.Phase != e.Phase {
		t.Errorf("Phase = %d, want %d", decoded.Phase, e.Phase)
	}
	if decoded.PhaseName != e.PhaseName {
		t.Errorf("PhaseName = %q, want %q", decoded.PhaseName, e.PhaseName)
	}
	if decoded.PhasesTotal != e.PhasesTotal {
		t.Errorf("PhasesTotal = %d, want %d", decoded.PhasesTotal, e.PhasesTotal)
	}
	if decoded.ThinkingSummary != e.ThinkingSummary {
		t.Errorf("ThinkingSummary = %q, want %q", decoded.ThinkingSummary, e.ThinkingSummary)
	}
	if decoded.FeatureName != e.FeatureName {
		t.Errorf("FeatureName = %q, want %q", decoded.FeatureName, e.FeatureName)
	}
	if decoded.Target != e.Target {
		t.Errorf("Target = %q, want %q", decoded.Target, e.Target)
	}
	if decoded.Result != e.Result {
		t.Errorf("Result = %q, want %q", decoded.Result, e.Result)
	}
	if decoded.Timestamp != e.Timestamp {
		t.Errorf("Timestamp = %q, want %q", decoded.Timestamp, e.Timestamp)
	}
	if decoded.ActionDetail != e.ActionDetail {
		t.Errorf("ActionDetail = %q, want %q", decoded.ActionDetail, e.ActionDetail)
	}
	if len(decoded.SkippedQuestions) != len(e.SkippedQuestions) {
		t.Errorf("SkippedQuestions length = %d, want %d", len(decoded.SkippedQuestions), len(e.SkippedQuestions))
	}
}

func TestEvent_JSONFieldNames(t *testing.T) {
	e := NewEvent("sess-789", "p1", 3, "Design", "Chose option C", "specs/design.md", "design-feature")
	e.ErrorMessage = "something broke"

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map error: %v", err)
	}

	requiredFields := []string{
		"session_id", "timestamp", "agent_id", "action_type",
		"action_detail", "target", "feature_name", "phase",
		"phase_name", "phases_total", "thinking_summary",
		"skipped_questions", "duration_ms", "result", "error_message",
	}
	for _, f := range requiredFields {
		if _, ok := raw[f]; !ok {
			t.Errorf("missing JSON field %q", f)
		}
	}
}

func TestEvent_ErrorMessageOmittedWhenEmpty(t *testing.T) {
	e := NewEvent("sess-err", "p1", 1, "Intent", "summary", "specs/x.md", "x")

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if strings.Contains(string(data), "error_message") {
		t.Error("error_message should be omitted when empty")
	}

	e.ErrorMessage = "parse failed"
	data, err = json.Marshal(e)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if !strings.Contains(string(data), `"error_message":"parse failed"`) {
		t.Errorf("error_message not found in JSON: %s", data)
	}
}

func TestEvent_DurationMSNullByDefault(t *testing.T) {
	e := NewEvent("sess-abc", "p1", 1, "Intent", "summary", "specs/x.md", "x")

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if raw["duration_ms"] != nil {
		t.Errorf("duration_ms = %v, want null", raw["duration_ms"])
	}

	// Verify it IS present (not omitted)
	if !strings.Contains(string(data), `"duration_ms":null`) {
		t.Errorf("duration_ms should be explicitly null in JSON, got: %s", data)
	}
}

func TestEvent_TimestampISO8601(t *testing.T) {
	e := NewEvent("sess-time", "p1", 1, "Intent", "summary", "specs/t.md", "t")

	// Must parse as RFC3339 (ISO-8601 compatible)
	if _, err := time.Parse(time.RFC3339, e.Timestamp); err != nil {
		t.Errorf("timestamp %q is not valid RFC3339: %v", e.Timestamp, err)
	}

	// Must not contain nanoseconds
	if strings.Contains(e.Timestamp, ".") {
		t.Errorf("timestamp %q contains sub-second precision, want clean RFC3339", e.Timestamp)
	}
}

func TestEvent_DefaultPhasesTotal(t *testing.T) {
	if DefaultPhasesTotal != 4 {
		t.Errorf("DefaultPhasesTotal = %d, want 4", DefaultPhasesTotal)
	}
}
