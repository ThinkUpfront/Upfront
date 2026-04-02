package format

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// DefaultPhasesTotal is the number of phases in a /feature run.
const DefaultPhasesTotal = 4

// Event represents an upfront audit event, extending the agent-monitoring
// trace format from the Delivery-Gap-Toolkit.
type Event struct {
	SessionID        string   `json:"session_id"`
	ProjectID        string   `json:"project_id"`
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

// ProjectID returns a stable hash derived from the git remote origin URL.
// Falls back to a hash of the working directory if not in a git repo.
func ProjectID(cwd string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "-C", cwd, "remote", "get-url", "origin") //nolint:gosec // args are fixed strings, cwd is from hook stdin
	out, err := cmd.Output()
	source := cwd
	if err == nil {
		source = strings.TrimSpace(string(out))
	}
	h := sha256.Sum256([]byte(source))
	return hex.EncodeToString(h[:8])
}

// NewEvent creates an Event with default values for agent_id, action_type,
// phases_total, result, and an empty skipped_questions slice.
func NewEvent(sessionID, projectID string, phase int, phaseName, thinkingSummary, target, featureName string) Event {
	return Event{
		SessionID:        sessionID,
		ProjectID:        projectID,
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
