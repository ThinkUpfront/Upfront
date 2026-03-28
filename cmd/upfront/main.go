package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/brennhill/upfront/internal/format"
	"github.com/brennhill/upfront/internal/hook"
	"github.com/brennhill/upfront/internal/queue"
	"github.com/brennhill/upfront/internal/remote"
)

const usage = `Usage: upfront <command> [options]

Commands:
  hook    Process Claude Code PostToolUse hook input from stdin
  flush   Flush queued events to remote endpoint
  log     Print audit events (--feature, --phase, --limit)
  purge   Delete events older than TTL
  status  Show queue status and configuration
`

const (
	defaultTTLDays  = 90
	defaultLogLimit = 50
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 1
	}

	switch args[0] {
	case "help", "--help", "-h":
		fmt.Fprint(stdout, usage)
		return 0
	case "hook":
		return cmdHook(stdin, stdout, stderr)
	case "flush":
		return cmdFlush(stdout, stderr)
	case "log":
		return cmdLog(args[1:], stdout, stderr)
	case "purge":
		return cmdPurge(stdout, stderr)
	case "status":
		return cmdStatus(stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown subcommand: %s\n", args[0])
		fmt.Fprint(stderr, usage)
		return 1
	}
}

func queuePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return filepath.Join(cwd, ".upfront", "audit.jsonl")
}

func loadCfg(stderr io.Writer) *remote.Config {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	cfg, err := remote.LoadConfig(cwd)
	if err != nil {
		fmt.Fprintf(stderr, "warning: config error: %v\n", err)
		return nil
	}
	return cfg
}

func cmdHook(stdin io.Reader, stdout, stderr io.Writer) int {
	const maxStdinBytes = 10 << 20 // 10MB
	data, err := io.ReadAll(io.LimitReader(stdin, maxStdinBytes+1))
	if err != nil {
		fmt.Fprintf(stderr, "error reading stdin: %v\n", err)
		return 1
	}

	// Gracefully handle empty input.
	if len(data) == 0 {
		return 0
	}

	if len(data) > maxStdinBytes {
		fmt.Fprintf(stderr, "warning: stdin exceeds %dMB, skipping\n", maxStdinBytes>>20)
		return 0
	}

	input, err := hook.ParseInput(data)
	if err != nil {
		// Don't crash on bad input -- this is the hot path.
		fmt.Fprintf(stderr, "warning: %v\n", err)
		return 0
	}

	events := hook.ExtractEvents(&input)
	if len(events) == 0 {
		return 0
	}

	q := queue.New(queuePath())
	if err := q.EnsureDir(); err != nil {
		fmt.Fprintf(stderr, "error creating queue dir: %v\n", err)
		return 1
	}

	for i := range events {
		if err := q.Append(&events[i]); err != nil {
			fmt.Fprintf(stderr, "error appending event: %v\n", err)
			return 1
		}
	}

	fmt.Fprintf(stdout, "queued %d events\n", len(events))

	// Attempt remote flush (best-effort, only if remote configured).
	tryRemoteFlush(q, stderr)

	return 0
}

func tryRemoteFlush(q *queue.Queue, stderr io.Writer) {
	cfg := loadCfg(stderr)
	if cfg == nil || cfg.Endpoint == "" {
		return
	}
	sender := remote.NewSender(cfg)
	flushed, err := q.Flush()
	if err != nil {
		fmt.Fprintf(stderr, "warning: flush failed: %v\n", err)
		return
	}
	if len(flushed) == 0 {
		return
	}
	if err := sender.Send(flushed); err != nil {
		fmt.Fprintf(stderr, "warning: remote send failed, events staged in .flushing for retry: %v\n", err)
		q.NackFlush()
		return
	}
	q.AckFlush()
}

func cmdFlush(stdout, stderr io.Writer) int {
	cfg := loadCfg(stderr)
	if cfg == nil || cfg.Endpoint == "" {
		fmt.Fprintln(stdout, "no remote endpoint configured; nothing to flush")
		return 0
	}

	q := queue.New(queuePath())
	sender := remote.NewSender(cfg)

	events, err := q.Flush()
	if err != nil {
		fmt.Fprintf(stderr, "error flushing queue: %v\n", err)
		return 1
	}
	if len(events) == 0 {
		fmt.Fprintln(stdout, "no events to flush")
		return 0
	}

	if err := sender.Send(events); err != nil {
		fmt.Fprintf(stderr, "error sending events, events staged in .flushing for retry: %v\n", err)
		q.NackFlush()
		return 1
	}
	q.AckFlush()

	fmt.Fprintf(stdout, "flushed %d events\n", len(events))
	return 0
}

func cmdLog(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("log", flag.ContinueOnError)
	fs.SetOutput(stderr)
	feature := fs.String("feature", "", "filter by feature name")
	phase := fs.Int("phase", 0, "filter by phase number")
	limit := fs.Int("limit", defaultLogLimit, "max events to display")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	q := queue.New(queuePath())
	events, err := q.ReadAll()
	if err != nil {
		fmt.Fprintf(stderr, "error reading queue: %v\n", err)
		return 1
	}

	filtered := filterEvents(events, *feature, *phase)

	// Apply limit (show most recent).
	if *limit > 0 && len(filtered) > *limit {
		filtered = filtered[len(filtered)-*limit:]
	}

	if len(filtered) == 0 {
		fmt.Fprintln(stdout, "no events found")
		return 0
	}

	for i := range filtered {
		e := &filtered[i]
		fmt.Fprintf(stdout, "%s  %-20s  phase %d/%d  %s  %s\n",
			e.Timestamp, e.FeatureName, e.Phase, e.PhasesTotal, e.PhaseName, truncate(e.ThinkingSummary, 60))
	}
	return 0
}

func filterEvents(events []format.Event, feature string, phase int) []format.Event {
	if feature == "" && phase == 0 {
		return events
	}
	var out []format.Event
	for i := range events {
		if feature != "" && events[i].FeatureName != feature {
			continue
		}
		if phase != 0 && events[i].Phase != phase {
			continue
		}
		out = append(out, events[i])
	}
	return out
}

func truncate(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}

func cmdPurge(stdout, stderr io.Writer) int {
	cfg := loadCfg(stderr)
	ttlDays := defaultTTLDays
	if cfg != nil && cfg.TTLDays > 0 {
		ttlDays = cfg.TTLDays
	}

	q := queue.New(queuePath())
	ttl := time.Duration(ttlDays) * 24 * time.Hour
	if err := q.Purge(ttl); err != nil {
		fmt.Fprintf(stderr, "error purging: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "purged events older than %d days\n", ttlDays)
	return 0
}

func cmdStatus(stdout, stderr io.Writer) int {
	qPath := queuePath()
	q := queue.New(qPath)

	events, err := q.ReadAll()
	if err != nil {
		fmt.Fprintf(stderr, "error reading queue: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Queue file:  %s\n", qPath)
	fmt.Fprintf(stdout, "Event count: %d\n", len(events))

	if len(events) > 0 {
		last := events[len(events)-1]
		fmt.Fprintf(stdout, "Last event:  %s (phase %d: %s, feature: %s)\n",
			last.Timestamp, last.Phase, last.PhaseName, last.FeatureName)
	} else {
		fmt.Fprintln(stdout, "Last event:  (none)")
	}

	cfg := loadCfg(stderr)
	if cfg != nil && cfg.Endpoint != "" {
		fmt.Fprintf(stdout, "Remote:      %s\n", cfg.Endpoint)
		if cfg.ProjectName != "" {
			fmt.Fprintf(stdout, "Project:     %s\n", cfg.ProjectName)
		}
	} else {
		fmt.Fprintln(stdout, "Remote:      (not configured)")
	}

	return 0
}
