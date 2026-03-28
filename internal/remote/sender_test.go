package remote

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brennhill/upfront/internal/format"
)

func TestLoadConfig_MalformedJSON(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".upfront")
	if err := os.MkdirAll(cfgDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "config.json"), []byte(`{not json`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(dir)
	if err == nil {
		t.Fatalf("LoadConfig() error = nil, want JSON parse error")
	}
	if cfg != nil {
		t.Fatalf("LoadConfig() = %+v, want nil on error", cfg)
	}
}

func TestLoadConfig_NeitherExists(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	cfg, err := loadConfig(
		filepath.Join(dir1, ".upfront", "config.json"),
		filepath.Join(dir2, ".upfront", "config.json"),
	)
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}
	if cfg != nil {
		t.Fatalf("loadConfig() = %+v, want nil", cfg)
	}
}

func TestLoadConfig_UserLevelFallback(t *testing.T) {
	projectDir := t.TempDir() // no config here
	userDir := t.TempDir()
	cfgDir := filepath.Join(userDir, ".upfront")
	if err := os.MkdirAll(cfgDir, 0o750); err != nil {
		t.Fatal(err)
	}
	cfgJSON := `{"endpoint": "https://user-level.example.com", "project_name": "fallback"}`
	if err := os.WriteFile(filepath.Join(cfgDir, "config.json"), []byte(cfgJSON), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadConfig(
		filepath.Join(projectDir, ".upfront", "config.json"),
		filepath.Join(userDir, ".upfront", "config.json"),
	)
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("loadConfig() returned nil, want non-nil")
	}
	if cfg.Endpoint != "https://user-level.example.com" {
		t.Errorf("Endpoint = %q, want %q", cfg.Endpoint, "https://user-level.example.com")
	}
	if cfg.ProjectName != "fallback" {
		t.Errorf("ProjectName = %q, want %q", cfg.ProjectName, "fallback")
	}
}

func TestLoadConfig_ProjectLevel(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".upfront")
	if err := os.MkdirAll(cfgDir, 0o750); err != nil {
		t.Fatal(err)
	}
	cfgJSON := `{
		"endpoint": "https://example.com/api/v1/traces",
		"auth_header": "Bearer pk-lf-xxxx",
		"ttl_days": 90,
		"project_name": "my-project"
	}`
	if err := os.WriteFile(filepath.Join(cfgDir, "config.json"), []byte(cfgJSON), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(dir)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig() returned nil, want non-nil")
	}
	if cfg.Endpoint != "https://example.com/api/v1/traces" {
		t.Errorf("Endpoint = %q, want %q", cfg.Endpoint, "https://example.com/api/v1/traces")
	}
	if cfg.AuthHeader != "Bearer pk-lf-xxxx" {
		t.Errorf("AuthHeader = %q, want %q", cfg.AuthHeader, "Bearer pk-lf-xxxx")
	}
	if cfg.TTLDays != 90 {
		t.Errorf("TTLDays = %d, want 90", cfg.TTLDays)
	}
	if cfg.ProjectName != "my-project" {
		t.Errorf("ProjectName = %q, want %q", cfg.ProjectName, "my-project")
	}
}

func TestSend_Success(t *testing.T) {
	var gotContentType string
	var gotAuth string
	var gotBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		gotAuth = r.Header.Get("Authorization")
		var err error
		gotBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewSender(&Config{
		Endpoint:   srv.URL,
		AuthHeader: "Bearer test-token",
	})

	events := []format.Event{
		format.NewEvent("sess-1", 1, "Intent", "summary here", "/tmp", "my-feature"),
	}
	if err := s.Send(events); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if gotContentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", gotContentType, "application/json")
	}
	if gotAuth != "Bearer test-token" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer test-token")
	}

	var decoded []format.Event
	if err := json.Unmarshal(gotBody, &decoded); err != nil {
		t.Fatalf("decoding body: %v", err)
	}
	if len(decoded) != 1 {
		t.Fatalf("len(decoded) = %d, want 1", len(decoded))
	}
	if decoded[0].SessionID != "sess-1" {
		t.Errorf("SessionID = %q, want %q", decoded[0].SessionID, "sess-1")
	}
}

func TestSend_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	s := NewSender(&Config{Endpoint: srv.URL})
	events := []format.Event{
		format.NewEvent("sess-1", 1, "Intent", "summary", "/tmp", "feat"),
	}
	err := s.Send(events)
	if err == nil {
		t.Fatal("Send() error = nil, want error on 500")
	}
}

func TestSend_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := &Sender{
		config: &Config{Endpoint: srv.URL},
		client: &http.Client{Timeout: 50 * time.Millisecond},
	}
	events := []format.Event{
		format.NewEvent("sess-1", 1, "Intent", "summary", "/tmp", "feat"),
	}
	err := s.Send(events)
	if err == nil {
		t.Fatal("Send() error = nil, want timeout error")
	}
}

func TestSend_NilConfig(t *testing.T) {
	s := NewSender(nil)
	events := []format.Event{
		format.NewEvent("sess-1", 1, "Intent", "summary", "/tmp", "feat"),
	}
	err := s.Send(events)
	if err != nil {
		t.Fatalf("Send() with nil config error = %v, want nil (no-op)", err)
	}
}

func TestSend_NoAuthHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewSender(&Config{Endpoint: srv.URL})
	events := []format.Event{
		format.NewEvent("sess-1", 1, "Intent", "summary", "/tmp", "feat"),
	}
	if err := s.Send(events); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if gotAuth != "" {
		t.Errorf("Authorization = %q, want empty (no auth configured)", gotAuth)
	}
}
