package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	repoURL         = "https://github.com/ThinkUpfront/Upfront.git"
	pluginName      = "upfront"
	marketplaceName = "thinkupfront"
	pluginKey       = pluginName + "@" + marketplaceName
)

func claudeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".claude"), nil
}

type cloneResult struct {
	version   string
	commitSHA string
	tmpDir    string
}

func cmdUpdate(stdout, stderr io.Writer) int {
	fmt.Fprintln(stdout, "Upfront plugin update")
	fmt.Fprintln(stdout, "")

	claudeBase, err := claudeDir()
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	pluginsDir := filepath.Join(claudeBase, "plugins")
	cacheBase := filepath.Join(pluginsDir, "cache", marketplaceName, pluginName)

	result, err := cloneAndReadVersion(stdout, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	defer os.RemoveAll(result.tmpDir) //nolint:errcheck // best-effort cleanup

	if err := installPlugin(stdout, stderr, result, cacheBase, pluginsDir, claudeBase); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, "")
	fmt.Fprintf(stdout, "Upfront plugin %s installed successfully.\n", result.version)
	fmt.Fprintln(stdout, "Restart Claude Code to pick up the new skills.")
	return 0
}

func installPlugin(stdout, stderr io.Writer, result cloneResult, cacheBase, pluginsDir, claudeBase string) error {
	// Remove ALL old cached versions.
	fmt.Fprintln(stdout, "Removing old plugin cache...")
	if err := os.RemoveAll(cacheBase); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(stderr, "warning: could not remove old cache: %v\n", err)
	}

	// Copy plugin/ contents to new cache location.
	installPath := filepath.Join(cacheBase, result.version)
	fmt.Fprintf(stdout, "Installing to %s\n", installPath)

	srcDir := filepath.Join(result.tmpDir, "plugin")
	if err := copyDir(srcDir, installPath); err != nil {
		return fmt.Errorf("copying plugin files: %w", err)
	}

	// Update installed_plugins.json.
	fmt.Fprintln(stdout, "Updating plugin registry...")
	if err := updateInstalledPlugins(pluginsDir, installPath, result.version, result.commitSHA); err != nil {
		return fmt.Errorf("updating installed_plugins.json: %w", err)
	}

	// Update settings.json (enable plugin + marketplace).
	fmt.Fprintln(stdout, "Updating Claude Code settings...")
	if err := updateSettings(claudeBase); err != nil {
		return fmt.Errorf("updating settings.json: %w", err)
	}

	return nil
}

func cloneAndReadVersion(stdout, stderr io.Writer) (cloneResult, error) {
	fmt.Fprintln(stdout, "Cloning latest from GitHub...")
	tmpDir, err := os.MkdirTemp("", "upfront-update-*")
	if err != nil {
		return cloneResult{}, fmt.Errorf("creating temp dir: %w", err)
	}

	ctx := context.Background()

	cloneCmd := exec.CommandContext(ctx, "git", "clone", "--depth=1", repoURL, tmpDir) //nolint:gosec // repoURL is a compile-time constant, tmpDir is from os.MkdirTemp
	cloneCmd.Stderr = stderr
	if err := cloneCmd.Run(); err != nil {
		_ = os.RemoveAll(tmpDir)
		return cloneResult{}, fmt.Errorf("cloning repo: %w", err)
	}

	pluginJSONPath := filepath.Join(tmpDir, "plugin", ".claude-plugin", "plugin.json")
	pluginJSON, err := os.ReadFile(pluginJSONPath) //#nosec G304 -- path built from temp dir
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return cloneResult{}, fmt.Errorf("reading plugin.json: %w", err)
	}
	var meta struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(pluginJSON, &meta); err != nil {
		_ = os.RemoveAll(tmpDir)
		return cloneResult{}, fmt.Errorf("parsing plugin.json: %w", err)
	}

	shaCmd := exec.CommandContext(ctx, "git", "-C", tmpDir, "rev-parse", "HEAD") //nolint:gosec // tmpDir is from os.MkdirTemp
	shaOut, err := shaCmd.Output()
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return cloneResult{}, fmt.Errorf("reading commit SHA: %w", err)
	}
	commitSHA := string(shaOut[:len(shaOut)-1])

	fmt.Fprintf(stdout, "Found version: %s\n", meta.Version)
	return cloneResult{version: meta.Version, commitSHA: commitSHA, tmpDir: tmpDir}, nil
}

// copyDir recursively copies src to dst using os.Root for safe traversal.
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o750)
		}

		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src) //#nosec G304 -- copying from controlled temp directory
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o600) //nolint:gosec // dst is built from filepath.Join of known directories
}

// installedPlugins represents ~/.claude/plugins/installed_plugins.json.
type installedPlugins struct {
	Version int                               `json:"version"`
	Plugins map[string][]installedPluginEntry `json:"plugins"`
}

type installedPluginEntry struct {
	Scope        string `json:"scope"`
	InstallPath  string `json:"installPath"`
	Version      string `json:"version"`
	InstalledAt  string `json:"installedAt"`
	LastUpdated  string `json:"lastUpdated"`
	GitCommitSHA string `json:"gitCommitSha,omitempty"`
}

func updateInstalledPlugins(pluginsDir, installPath, ver, commitSHA string) error {
	filePath := filepath.Join(pluginsDir, "installed_plugins.json")
	now := time.Now().UTC().Format(time.RFC3339)

	var data installedPlugins
	raw, err := os.ReadFile(filePath) //#nosec G304 -- well-known path
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		data = installedPlugins{Version: 2, Plugins: make(map[string][]installedPluginEntry)}
	} else {
		if err := json.Unmarshal(raw, &data); err != nil {
			return fmt.Errorf("parse installed_plugins.json: %w", err)
		}
	}

	installedAt := now
	if existing, ok := data.Plugins[pluginKey]; ok && len(existing) > 0 {
		installedAt = existing[0].InstalledAt
	}

	data.Plugins[pluginKey] = []installedPluginEntry{{
		Scope:        "user",
		InstallPath:  installPath,
		Version:      ver,
		InstalledAt:  installedAt,
		LastUpdated:  now,
		GitCommitSHA: commitSHA,
	}}

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, out, 0o600)
}

type marketplaceEntry struct {
	Source marketplaceSource `json:"source"`
}

type marketplaceSource struct {
	Source string `json:"source"`
	Repo   string `json:"repo"`
}

func updateSettings(claudeBase string) error {
	filePath := filepath.Join(claudeBase, "settings.json")

	raw, err := os.ReadFile(filePath) //#nosec G304 -- well-known path
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var full map[string]json.RawMessage
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &full); err != nil {
			return fmt.Errorf("parse settings.json: %w", err)
		}
	}
	if full == nil {
		full = make(map[string]json.RawMessage)
	}

	// Update enabledPlugins.
	var enabled map[string]bool
	if ep, ok := full["enabledPlugins"]; ok {
		if err := json.Unmarshal(ep, &enabled); err != nil {
			enabled = make(map[string]bool)
		}
	}
	if enabled == nil {
		enabled = make(map[string]bool)
	}
	enabled[pluginKey] = true
	epBytes, _ := json.Marshal(enabled)
	full["enabledPlugins"] = epBytes

	// Update extraKnownMarketplaces.
	var markets map[string]json.RawMessage
	if ekm, ok := full["extraKnownMarketplaces"]; ok {
		if err := json.Unmarshal(ekm, &markets); err != nil {
			markets = make(map[string]json.RawMessage)
		}
	}
	if markets == nil {
		markets = make(map[string]json.RawMessage)
	}
	entry := marketplaceEntry{
		Source: marketplaceSource{
			Source: "github",
			Repo:   "ThinkUpfront/Upfront",
		},
	}
	entryBytes, _ := json.Marshal(entry)
	markets[marketplaceName] = entryBytes
	ekmBytes, _ := json.Marshal(markets)
	full["extraKnownMarketplaces"] = ekmBytes

	out, err := json.MarshalIndent(full, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, out, 0o600)
}
