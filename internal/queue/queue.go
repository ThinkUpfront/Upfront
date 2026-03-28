package queue

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/brennhill/upfront/internal/format"
)

// Queue manages a local JSONL audit event queue.
type Queue struct {
	path string
}

// New creates a Queue that writes to the given file path.
func New(path string) *Queue {
	return &Queue{path: path}
}

// EnsureDir creates the parent directory of the queue file if it doesn't exist.
func (q *Queue) EnsureDir() error {
	return os.MkdirAll(filepath.Dir(q.path), 0o750)
}

// Append JSON-encodes the event and appends it as a single line to the queue file.
// Uses O_APPEND for concurrent write safety on POSIX systems.
// Calls fsync to ensure durability before returning.
func (q *Queue) Append(e *format.Event) error {
	f, err := os.OpenFile(q.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	if _, err = f.Write(data); err != nil {
		return err
	}
	return f.Sync()
}

// ReadAll reads and parses all events from the queue file.
// Returns an empty slice (no error) if the file is missing or empty.
func (q *Queue) ReadAll() ([]format.Event, error) {
	return readEventsFromFile(q.path)
}

// readEventsFromFile parses JSONL events from the given path.
func readEventsFromFile(path string) ([]format.Event, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var events []format.Event
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, bufio.MaxScanTokenSize), 1<<20) // 1MB max line
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e format.Event
		if err := json.Unmarshal(line, &e); err != nil {
			// Skip corrupt lines rather than bricking the queue.
			continue
		}
		events = append(events, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

// appendEventsToFile appends JSON-encoded events to the given path with fsync.
func appendEventsToFile(path string, events []format.Event) error {
	f, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	for i := range events {
		data, err := json.Marshal(&events[i])
		if err != nil {
			return err
		}
		data = append(data, '\n')
		if _, err := f.Write(data); err != nil {
			return err
		}
	}
	return f.Sync()
}

const staleLockAge = 60 * time.Second

// tryLockFlush acquires an exclusive flush lock using O_CREATE|O_EXCL.
// Returns an unlock function and true on success, or nil and false if
// another flush is in progress. Stale locks (older than 60s) are broken.
func (q *Queue) tryLockFlush() (func(), bool) {
	lockPath := filepath.Clean(q.path + ".lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		// Check for stale lock from a crashed process.
		info, statErr := os.Stat(lockPath)
		if statErr != nil || time.Since(info.ModTime()) < staleLockAge {
			return nil, false
		}
		_ = os.Remove(lockPath)
		f, err = os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600) //nolint:gosec // lockPath derived from q.path
		if err != nil {
			return nil, false
		}
	}
	_ = f.Close()
	return func() { _ = os.Remove(lockPath) }, true
}

// Flush atomically drains the queue and stages events for remote send:
//  1. Acquire exclusive flush lock (returns nil if another flush is in progress)
//  2. Recover any .drain file from a prior crash (rename interrupted)
//  3. Rename audit.jsonl -> .drain (atomic grab of new events)
//  4. Read .drain and append to .flushing (consolidate for crash recovery)
//  5. Remove .drain
//
// The .flushing file persists until AckFlush is called after successful send.
// If the process crashes between Flush and AckFlush, the next Flush recovers
// all events from .flushing.
func (q *Queue) Flush() ([]format.Event, error) {
	unlock, ok := q.tryLockFlush()
	if !ok {
		return nil, nil // Another flush in progress.
	}

	flushPath := q.path + ".flushing"
	drainPath := q.path + ".drain"

	// Helper to unlock on error paths (AckFlush handles the success path).
	fail := func(err error) ([]format.Event, error) {
		unlock()
		return nil, err
	}

	// Recover .drain from prior crash (rename succeeded but append didn't).
	drained, drainErr := readEventsFromFile(drainPath)
	if drainErr != nil {
		return fail(fmt.Errorf("recover .drain: %w", drainErr))
	}
	if len(drained) > 0 {
		if err := appendEventsToFile(flushPath, drained); err != nil {
			return fail(err)
		}
		_ = os.Remove(drainPath)
	}

	// Atomically grab new events from the queue.
	if err := q.drainToFlushFile(drainPath, flushPath); err != nil {
		return fail(err)
	}

	// Read the consolidated .flushing file.
	events, err := readEventsFromFile(flushPath)
	if err != nil {
		return fail(err)
	}
	if len(events) == 0 {
		unlock()
	}
	return events, nil
}

// drainToFlushFile renames audit.jsonl to drainPath, reads it, appends events
// to flushPath, and removes drainPath. No-op if audit.jsonl doesn't exist.
func (q *Queue) drainToFlushFile(drainPath, flushPath string) error {
	if err := os.Rename(q.path, drainPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // No queue file — nothing to drain.
		}
		return err
	}

	newEvents, err := readEventsFromFile(drainPath)
	if err != nil {
		_ = os.Rename(drainPath, q.path)
		return err
	}
	if len(newEvents) > 0 {
		if err := appendEventsToFile(flushPath, newEvents); err != nil {
			_ = os.Rename(drainPath, q.path)
			return err
		}
	}
	_ = os.Remove(drainPath)
	return nil
}

// AckFlush removes the .flushing file and releases the flush lock after events
// have been successfully sent to the remote endpoint.
func (q *Queue) AckFlush() {
	_ = os.Remove(q.path + ".flushing")
	_ = os.Remove(q.path + ".lock")
}

// NackFlush releases the flush lock but keeps the .flushing file for
// recovery on the next flush attempt. Use when remote send fails.
func (q *Queue) NackFlush() {
	_ = os.Remove(q.path + ".lock")
}

// Purge rewrites the queue file keeping only events newer than the given TTL.
// Uses rename-first to avoid losing events from concurrent Append calls.
// If the file is missing, this is a no-op.
func (q *Queue) Purge(ttl time.Duration) error {
	// Rename-first: move the file so concurrent Append creates a new one.
	purgePath := q.path + ".purging"
	if err := os.Rename(q.path, purgePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	events, err := readEventsFromFile(purgePath)
	if err != nil {
		// Restore original file on read failure.
		_ = os.Rename(purgePath, q.path)
		return err
	}

	cutoff := time.Now().UTC().Add(-ttl)
	var kept []format.Event
	for i := range events {
		ts, err := time.Parse(time.RFC3339, events[i].Timestamp)
		if err != nil {
			// Keep events with unparseable timestamps rather than silently dropping them.
			kept = append(kept, events[i])
			continue
		}
		if ts.After(cutoff) {
			kept = append(kept, events[i])
		}
	}

	// Append kept events in a single batch BEFORE removing .purging. If the
	// process crashes during write, the .purging file still has all events.
	// If a concurrent Append created a new queue file while we held the rename,
	// both sets of events end up in the same file via O_APPEND.
	if len(kept) > 0 {
		if err := appendEventsToFile(q.path, kept); err != nil {
			return err
		}
	}
	_ = os.Remove(purgePath)
	return nil
}
