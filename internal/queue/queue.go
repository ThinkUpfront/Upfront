package queue

import (
	"bufio"
	"encoding/json"
	"errors"
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

	_, err = f.Write(data)
	return err
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

// Flush atomically drains the queue using rename-and-swap:
// 1. Recover any pre-existing .flushing file from a prior crash
// 2. Rename audit.jsonl to audit.jsonl.flushing
// 3. Read all events from the .flushing file
// 4. Delete the .flushing file
// Returns the events for the caller to send to a remote endpoint.
// If the file is missing or empty, returns an empty slice with no error.
func (q *Queue) Flush() ([]format.Event, error) {
	flushPath := q.path + ".flushing"

	// Recover events from a pre-existing .flushing file left by a prior crash.
	prior, err := readEventsFromFile(flushPath)
	if err != nil {
		return nil, err
	}

	if err := os.Rename(q.path, flushPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// No new events, but may have recovered prior ones.
			if len(prior) > 0 {
				_ = os.Remove(flushPath)
			}
			return prior, nil
		}
		return nil, err
	}

	events, err := readEventsFromFile(flushPath)
	if err != nil {
		return nil, err
	}

	// Combine any prior crash-recovered events with the current batch.
	if len(prior) > 0 {
		events = append(prior, events...)
	}

	if err := os.Remove(flushPath); err != nil {
		return events, err
	}

	return events, nil
}

// Purge rewrites the queue file keeping only events newer than the given TTL.
// If the file is missing, this is a no-op.
func (q *Queue) Purge(ttl time.Duration) error {
	events, err := readEventsFromFile(q.path)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
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

	// Write kept events to a temp file, then rename over the original.
	tmpPath := q.path + ".purging"
	f, err := os.OpenFile(filepath.Clean(tmpPath), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	for i := range kept {
		data, err := json.Marshal(kept[i])
		if err != nil {
			_ = f.Close()
			_ = os.Remove(tmpPath)
			return err
		}
		data = append(data, '\n')
		if _, err := f.Write(data); err != nil {
			_ = f.Close()
			_ = os.Remove(tmpPath)
			return err
		}
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, q.path)
}
