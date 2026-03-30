package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/ThinkUpfront/Upfront/internal/format"
)

// Config holds remote sender configuration.
type Config struct {
	Endpoint    string `json:"endpoint"`
	AuthHeader  string `json:"auth_header"`
	TTLDays     int    `json:"ttl_days"`
	ProjectName string `json:"project_name"`
}

// LoadConfig loads config from projectDir/.upfront/config.json, falling back
// to ~/.upfront/config.json. Returns nil (no error) if neither exists.
func LoadConfig(projectDir string) (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return loadConfig(
		filepath.Join(projectDir, ".upfront", "config.json"),
		filepath.Join(home, ".upfront", "config.json"),
	)
}

func loadConfig(projectPath, userPath string) (*Config, error) {
	cfg, err := loadConfigFile(projectPath)
	if err != nil {
		return nil, err
	}
	if cfg != nil {
		return cfg, nil
	}
	return loadConfigFile(userPath)
}

// Sender sends events to a remote endpoint.
type Sender struct {
	config *Config
	client *http.Client
}

// NewSender creates a Sender from the given config. If config is nil, Send
// is a no-op (local-only mode).
func NewSender(cfg *Config) *Sender {
	return &Sender{
		config: cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Send POSTs a JSON array of events to the configured endpoint. Returns nil
// immediately if config is nil (local-only mode). Returns an error on network
// failure or non-2xx response status.
func (s *Sender) Send(events []format.Event) error {
	if s.config == nil {
		return nil
	}

	body, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("marshal events: %w", err)
	}

	u, err := url.Parse(s.config.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint URL: %w", err)
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return fmt.Errorf("endpoint must use http or https scheme, got %q", u.Scheme)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, s.config.Endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if s.config.AuthHeader != "" {
		req.Header.Set("Authorization", s.config.AuthHeader)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send events: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("remote returned status %d", resp.StatusCode)
	}
	return nil
}

func loadConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
