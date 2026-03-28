package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain builds the binary once for all integration tests.
var binaryPath string

func TestMain(m *testing.M) {
	tmp, err := os.CreateTemp("", "upfront-test-*")
	if err != nil {
		panic(err)
	}
	_ = tmp.Close()
	binaryPath = tmp.Name()

	cmd := exec.CommandContext(context.Background(), "go", "build", "-o", binaryPath, ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic("build failed: " + string(out))
	}

	code := m.Run()
	_ = os.Remove(binaryPath)
	os.Exit(code)
}

func runCmd(t *testing.T, stdin string, args ...string) (string, string, int) { //nolint:gocritic
	t.Helper()
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, binaryPath, args...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	dir := t.TempDir()
	cmd.Dir = dir
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	code := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			code = exitErr.ExitCode()
		} else {
			t.Fatalf("unexpected error running command: %v", err)
		}
	}
	return outBuf.String(), errBuf.String(), code
}

func TestStatusExitsZero(t *testing.T) {
	_, _, code := runCmd(t, "", "status")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
}

func TestHookEmptyJSON(t *testing.T) {
	_, _, code := runCmd(t, "{}", "hook")
	if code != 0 {
		t.Fatalf("expected exit 0 for empty JSON hook input, got %d", code)
	}
}

func TestLogEmptyQueue(t *testing.T) {
	_, _, code := runCmd(t, "", "log")
	if code != 0 {
		t.Fatalf("expected exit 0 for log with empty queue, got %d", code)
	}
}

func TestHookWithFeatureInput(t *testing.T) {
	input := `{
		"session_id": "test-session-123",
		"cwd": "/tmp/project",
		"tool_name": "Skill",
		"tool_input": {"skill_name": "feature", "args": "checkout-flow"},
		"tool_response": "### Thinking Record: Intent\n**Decided:** Build checkout flow.\n**Reasoning:** Users need to buy things.\n---\n### Thinking Record: Behavioral Spec\n**Decided:** Define behaviors.\n**Reasoning:** Need clear states.\n---"
	}`
	stdout, _, code := runCmd(t, input, "hook")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "queued 2 events") {
		t.Errorf("expected 'queued 2 events' in stdout, got: %s", stdout)
	}
}

func TestHookQueuesEventsToFile(t *testing.T) {
	input := `{
		"session_id": "test-session-456",
		"cwd": "/tmp/project",
		"tool_name": "Skill",
		"tool_input": {"skill_name": "feature", "args": "search-feature"},
		"tool_response": "### Thinking Record: Intent\n**Decided:** Build search.\n**Reasoning:** Users need search.\n---"
	}`
	dir := t.TempDir()
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, binaryPath, "hook")
	cmd.Stdin = strings.NewReader(input)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hook failed: %s (output: %s)", err, out)
	}

	queueFile := filepath.Join(dir, ".upfront", "audit.jsonl")
	data, err := os.ReadFile(queueFile)
	if err != nil {
		t.Fatalf("queue file not created: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 event in queue, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "search-feature") {
		t.Errorf("expected feature_name 'search-feature' in event, got: %s", lines[0])
	}
}

func TestStatusShowsQueueInfo(t *testing.T) {
	stdout, _, code := runCmd(t, "", "status")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "Queue") {
		t.Errorf("expected 'Queue' in status output, got: %s", stdout)
	}
}

func TestUnknownSubcommand(t *testing.T) {
	_, stderr, code := runCmd(t, "", "bogus")
	if code == 0 {
		t.Fatal("expected non-zero exit for unknown subcommand")
	}
	if !strings.Contains(stderr, "unknown subcommand") {
		t.Errorf("expected 'unknown subcommand' in stderr, got: %s", stderr)
	}
}
