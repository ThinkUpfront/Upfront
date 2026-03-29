package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/brennhill/upfront/internal/format"
)

func TestEnsureDir_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), ".upfront")
	q := New(filepath.Join(dir, "audit.jsonl"))

	if err := q.EnsureDir(); err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected directory, got file")
	}
}

func newTestQueue(t *testing.T) *Queue {
	t.Helper()
	dir := filepath.Join(t.TempDir(), ".upfront")
	q := New(filepath.Join(dir, "audit.jsonl"))
	if err := q.EnsureDir(); err != nil {
		t.Fatalf("EnsureDir: %v", err)
	}
	return q
}

func testEvent(session string, phase int, phaseName string) *format.Event {
	e := format.NewEvent(session, phase, phaseName, "summary for "+phaseName, "/tmp/target", "test-feature")
	return &e
}

func appendOrFail(t *testing.T, q *Queue, e *format.Event) {
	t.Helper()
	if err := q.Append(e); err != nil {
		t.Fatalf("Append: %v", err)
	}
}

func TestAppend_WritesValidJSONL(t *testing.T) {
	q := newTestQueue(t)

	e1 := testEvent("s1", 1, "Intent")
	e2 := testEvent("s1", 2, "Behavioral Spec")

	if err := q.Append(e1); err != nil {
		t.Fatalf("Append e1: %v", err)
	}
	if err := q.Append(e2); err != nil {
		t.Fatalf("Append e2: %v", err)
	}

	data, err := os.ReadFile(q.path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), string(data))
	}
	if !strings.Contains(lines[0], `"session_id":"s1"`) {
		t.Errorf("line 0 missing session_id: %s", lines[0])
	}
	if !strings.Contains(lines[1], `"phase_name":"Behavioral Spec"`) {
		t.Errorf("line 1 missing phase_name: %s", lines[1])
	}
}

func TestAppend_ConcurrentWritesNoCurruption(t *testing.T) {
	q := newTestQueue(t)
	const n = 50

	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e := testEvent(fmt.Sprintf("s%d", i), 1, "Intent")
			if err := q.Append(e); err != nil {
				t.Errorf("Append %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	data, err := os.ReadFile(q.path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != n {
		t.Fatalf("expected %d lines, got %d", n, len(lines))
	}

	// Every line must be valid JSON
	for i, line := range lines {
		var e format.Event
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Errorf("line %d not valid JSON: %v\n%s", i, err, line)
		}
	}
}

func TestReadAll_ReturnsAllEvents(t *testing.T) {
	q := newTestQueue(t)

	appendOrFail(t, q, testEvent("s1", 1, "Intent"))
	appendOrFail(t, q, testEvent("s1", 2, "Behavioral Spec"))

	events, err := q.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].PhaseName != "Intent" {
		t.Errorf("events[0].PhaseName = %q, want Intent", events[0].PhaseName)
	}
	if events[1].PhaseName != "Behavioral Spec" {
		t.Errorf("events[1].PhaseName = %q, want Behavioral Spec", events[1].PhaseName)
	}
}

func TestReadAll_EmptyFile(t *testing.T) {
	q := newTestQueue(t)

	if err := os.WriteFile(q.path, []byte{}, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	events, err := q.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on empty file: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	q := newTestQueue(t)

	events, err := q.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on missing file: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestFlush_ReturnsEventsAndStagesFlushing(t *testing.T) {
	q := newTestQueue(t)

	appendOrFail(t, q, testEvent("s1", 1, "Intent"))
	appendOrFail(t, q, testEvent("s1", 2, "Behavioral Spec"))

	events, err := q.Flush()
	if err != nil {
		t.Fatalf("Flush: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].PhaseName != "Intent" {
		t.Errorf("events[0].PhaseName = %q, want Intent", events[0].PhaseName)
	}

	// Original file should be gone.
	if _, err := os.Stat(q.path); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected queue file to be removed, got err=%v", err)
	}

	// .flushing should exist (crash recovery) until AckFlush is called.
	if _, err := os.Stat(q.path + ".flushing"); err != nil {
		t.Errorf("expected .flushing file to exist for crash recovery, got err=%v", err)
	}

	// AckFlush cleans up the .flushing file.
	q.AckFlush()
	if _, err := os.Stat(q.path + ".flushing"); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected .flushing file to be removed after AckFlush, got err=%v", err)
	}
}

func TestFlush_SecondFlushBlockedByLock(t *testing.T) {
	q := newTestQueue(t)

	appendOrFail(t, q, testEvent("s1", 1, "Intent"))
	appendOrFail(t, q, testEvent("s1", 2, "Behavioral Spec"))

	events1, err := q.Flush()
	if err != nil {
		t.Fatalf("Flush 1: %v", err)
	}
	if len(events1) != 2 {
		t.Fatalf("expected 2 events from first flush, got %d", len(events1))
	}

	// Second flush should return nil (lock held by first flush).
	events2, err := q.Flush()
	if err != nil {
		t.Fatalf("Flush 2: %v", err)
	}
	if len(events2) != 0 {
		t.Fatalf("expected 0 events from second flush (lock held), got %d", len(events2))
	}

	q.AckFlush()
}

func TestFlush_FailedSendRetainsEventsForRetry(t *testing.T) {
	q := newTestQueue(t)

	appendOrFail(t, q, testEvent("s1", 1, "Intent"))

	// First flush: simulate a failed send (don't call AckFlush).
	events1, err := q.Flush()
	if err != nil {
		t.Fatalf("Flush 1: %v", err)
	}
	if len(events1) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events1))
	}

	// Release lock without acking (simulates failed send).
	_ = os.Remove(q.path + ".lock")

	// Add another event while .flushing still exists.
	appendOrFail(t, q, testEvent("s2", 2, "Behavioral Spec"))

	// Second flush should recover events from .flushing plus new queue events.
	events2, err := q.Flush()
	if err != nil {
		t.Fatalf("Flush 2: %v", err)
	}
	if len(events2) != 2 {
		t.Fatalf("expected 2 events (1 recovered + 1 new), got %d", len(events2))
	}

	q.AckFlush()

	// After ack, .flushing should be gone.
	if _, err := os.Stat(q.path + ".flushing"); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected .flushing removed after AckFlush")
	}
}

func TestFlush_DrainFileRecovery(t *testing.T) {
	q := newTestQueue(t)

	// Simulate a crash that left a .drain file (rename succeeded, append didn't).
	drainPath := q.path + ".drain"
	e := testEvent("s-drain", 1, "Intent")
	data, _ := json.Marshal(e)
	data = append(data, '\n')
	if err := os.WriteFile(drainPath, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	events, err := q.Flush()
	if err != nil {
		t.Fatalf("Flush: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 recovered event, got %d", len(events))
	}
	if events[0].SessionID != "s-drain" {
		t.Errorf("expected s-drain, got %q", events[0].SessionID)
	}

	// .drain should be cleaned up.
	if _, err := os.Stat(drainPath); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected .drain file removed after recovery")
	}

	q.AckFlush()
}

func TestFlush_MissingFile(t *testing.T) {
	q := newTestQueue(t)

	events, err := q.Flush()
	if err != nil {
		t.Fatalf("Flush on missing file: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestFlush_EmptyFile(t *testing.T) {
	q := newTestQueue(t)
	if err := os.WriteFile(q.path, []byte{}, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	events, err := q.Flush()
	if err != nil {
		t.Fatalf("Flush on empty file: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestPurge_KeepsRecentRemovesOld(t *testing.T) {
	q := newTestQueue(t)

	// Write an old event (timestamp 100 days ago).
	old := testEvent("s-old", 1, "Intent")
	old.Timestamp = time.Now().UTC().Add(-100 * 24 * time.Hour).Format(time.RFC3339)
	appendOrFail(t, q, old)

	// Write a recent event (timestamp now).
	recent := testEvent("s-recent", 2, "Behavioral Spec")
	appendOrFail(t, q, recent)

	ttl := 90 * 24 * time.Hour
	if err := q.Purge(ttl); err != nil {
		t.Fatalf("Purge: %v", err)
	}

	events, err := q.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll after purge: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event after purge, got %d", len(events))
	}
	if events[0].SessionID != "s-recent" {
		t.Errorf("expected s-recent, got %q", events[0].SessionID)
	}
}

func TestPurge_MissingFile(t *testing.T) {
	q := newTestQueue(t)

	err := q.Purge(90 * 24 * time.Hour)
	if err != nil {
		t.Fatalf("Purge on missing file: %v", err)
	}
}

func TestFlush_RecoversPriorFlushingFile(t *testing.T) {
	q := newTestQueue(t)

	// Simulate a prior crash: write events to .flushing file directly.
	flushPath := q.path + ".flushing"
	prior := testEvent("s-prior", 1, "Intent")
	data, _ := json.Marshal(prior)
	data = append(data, '\n')
	if err := os.WriteFile(flushPath, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Write a new event to the main queue.
	appendOrFail(t, q, testEvent("s-new", 2, "Behavioral Spec"))

	events, err := q.Flush()
	if err != nil {
		t.Fatalf("Flush: %v", err)
	}

	// Should get both: prior recovered + new.
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].SessionID != "s-prior" {
		t.Errorf("events[0].SessionID = %q, want s-prior", events[0].SessionID)
	}
	if events[1].SessionID != "s-new" {
		t.Errorf("events[1].SessionID = %q, want s-new", events[1].SessionID)
	}
}

func TestReadAll_SkipsCorruptLines(t *testing.T) {
	q := newTestQueue(t)

	// Write a valid event, a corrupt line, and another valid event.
	appendOrFail(t, q, testEvent("s1", 1, "Intent"))

	f, err := os.OpenFile(q.path, os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	if _, err := f.WriteString("this is not json\n"); err != nil {
		t.Fatalf("Write corrupt: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	appendOrFail(t, q, testEvent("s2", 2, "Behavioral Spec"))

	events, err := q.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events (corrupt line skipped), got %d", len(events))
	}
	if events[0].SessionID != "s1" {
		t.Errorf("events[0].SessionID = %q, want s1", events[0].SessionID)
	}
	if events[1].SessionID != "s2" {
		t.Errorf("events[1].SessionID = %q, want s2", events[1].SessionID)
	}
}

func TestNackFlush_KeepsFlushingFile(t *testing.T) {
	q := newTestQueue(t)
	appendOrFail(t, q, testEvent("s1", 1, "Intent"))

	events, err := q.Flush()
	if err != nil {
		t.Fatalf("Flush: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	// NackFlush should release lock but keep .flushing.
	q.NackFlush()

	flushPath := q.path + ".flushing"
	if _, err := os.Stat(flushPath); err != nil {
		t.Errorf(".flushing file should exist after NackFlush, got: %v", err)
	}
	lockPath := q.path + ".lock"
	if _, err := os.Stat(lockPath); !errors.Is(err, os.ErrNotExist) {
		t.Errorf(".lock file should be removed after NackFlush, got: %v", err)
	}

	// Next flush should recover the .flushing events.
	recovered, err := q.Flush()
	if err != nil {
		t.Fatalf("second Flush: %v", err)
	}
	if len(recovered) != 1 {
		t.Fatalf("expected 1 recovered event, got %d", len(recovered))
	}
	if recovered[0].SessionID != "s1" {
		t.Errorf("recovered SessionID = %q, want s1", recovered[0].SessionID)
	}
	q.AckFlush()
}
