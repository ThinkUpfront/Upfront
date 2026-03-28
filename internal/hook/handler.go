package hook

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/brennhill/upfront/internal/format"
)

// Input represents the JSON payload from Claude Code PostToolUse hook stdin.
type Input struct {
	SessionID    string    `json:"session_id"`
	Cwd          string    `json:"cwd"`
	ToolName     string    `json:"tool_name"`
	ToolInput    ToolInput `json:"tool_input"`
	ToolResponse string    `json:"tool_response"`
}

// ToolInput represents the tool_input field from the hook payload.
type ToolInput struct {
	SkillName string `json:"skill_name"`
	Args      string `json:"args"`
}

// phaseNumbers maps thinking record phase names to their phase number.
var phaseNumbers = map[string]int{
	"Intent":                1,
	"Behavioral Spec":       2,
	"Design Approach":       3,
	"Implementation Design": 4,
}

// thinkingRecordRe matches "### Thinking Record: <PhaseName>" and captures
// the phase name and the body text up to the next "---" or "##" heading or end of input.
// Uses (?:^|\n) instead of (?m)^ to avoid (?m) making $ match end-of-line,
// which would cause the lazy [\s\S]*? to stop at the first line boundary.
var thinkingRecordRe = regexp.MustCompile(
	`(?:^|\n)### Thinking Record:\s*(.+?)\s*\n([\s\S]*?)(?:\n---|\n##|\z)`,
)

// ParseInput parses raw JSON bytes from stdin into a Input.
func ParseInput(data []byte) (Input, error) {
	var h Input
	if err := json.Unmarshal(data, &h); err != nil {
		return Input{}, fmt.Errorf("parse hook input: %w", err)
	}
	return h, nil
}

// ExtractEvents parses thinking records from a Input and returns
// a slice of Events. Returns nil if the tool is not a Skill invocation,
// if tool_response is empty, or if no thinking records are found.
func ExtractEvents(h *Input) []format.Event {
	if h.ToolName != "Skill" {
		return nil
	}
	if h.ToolInput.SkillName != "feature" {
		return nil
	}
	if h.ToolResponse == "" {
		return nil
	}

	matches := thinkingRecordRe.FindAllStringSubmatch(h.ToolResponse, -1)
	if len(matches) == 0 {
		return nil
	}

	var events []format.Event
	for _, m := range matches {
		phaseName := m[1]
		summary := strings.TrimSpace(m[2])

		phase, ok := phaseNumbers[phaseName]
		if !ok {
			continue
		}

		e := format.NewEvent(h.SessionID, phase, phaseName, summary, h.Cwd, h.ToolInput.Args)
		events = append(events, e)
	}
	return events
}
