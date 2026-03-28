package format

import (
	"fmt"
	"time"
)

// DefaultPhasesTotal is the number of phases in a /feature run.
const DefaultPhasesTotal = 4

// Event represents an upfront audit event, extending the agent-monitoring
// trace format from the Delivery-Gap-Toolkit.
type Event struct {
	SessionID        string   `json:"session_id"`
	Timestamp        string   `json:"timestamp"`
	AgentID          string   `json:"agent_id"`
	ActionType       string   `json:"action_type"`
	ActionDetail     string   `json:"action_detail"`
	Target           string   `json:"target"`
	FeatureName      string   `json:"feature_name"`
	Phase            int      `json:"phase"`
	PhaseName        string   `json:"phase_name"`
	PhasesTotal      int      `json:"phases_total"`
	ThinkingSummary  string   `json:"thinking_summary"`
	SkippedQuestions []string `json:"skipped_questions"`
	DurationMS       *int64   `json:"duration_ms"`
	Result           string   `json:"result"`
	ErrorMessage     string   `json:"error_message,omitempty"`
}

// NewEvent creates an Event with default values for agent_id, action_type,
// phases_total, result, and an empty skipped_questions slice.
func NewEvent(sessionID string, phase int, phaseName, thinkingSummary, target, featureName string) Event {
	return Event{
		SessionID:        sessionID,
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
		AgentID:          "upfront",
		ActionType:       "upfront_phase_complete",
		ActionDetail:     fmt.Sprintf("Phase %d: %s", phase, phaseName),
		Target:           target,
		FeatureName:      featureName,
		Phase:            phase,
		PhaseName:        phaseName,
		PhasesTotal:      DefaultPhasesTotal,
		ThinkingSummary:  thinkingSummary,
		SkippedQuestions: []string{},
		DurationMS:       nil,
		Result:           "success",
	}
}
