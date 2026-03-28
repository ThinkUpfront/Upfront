#!/usr/bin/env bash
set -euo pipefail

# install.sh — Build upfront and register the PostToolUse hook in Claude Code.
# Idempotent: safe to run multiple times.

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
SETTINGS_FILE="$HOME/.claude/settings.json"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# ---------- helpers ----------

info()  { printf '\033[1;34m==>\033[0m %s\n' "$1"; }
warn()  { printf '\033[1;33mwarning:\033[0m %s\n' "$1"; }
die()   { printf '\033[1;31merror:\033[0m %s\n' "$1" >&2; exit 1; }

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "'$1' is required but not found. Please install it first."
}

# ---------- preflight ----------

require_cmd go
require_cmd jq

# ---------- build ----------

info "Building upfront binary..."
(cd "$SCRIPT_DIR" && go build -o upfront ./cmd/upfront/)

# ---------- install binary ----------

mkdir -p "$INSTALL_DIR"
cp "$SCRIPT_DIR/upfront" "$INSTALL_DIR/upfront"
chmod +x "$INSTALL_DIR/upfront"
info "Installed binary to $INSTALL_DIR/upfront"

# Check if INSTALL_DIR is on PATH
if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
  warn "$INSTALL_DIR is not on your PATH."
  warn "Add this to your shell profile:  export PATH=\"$INSTALL_DIR:\$PATH\""
fi

# ---------- hook registration ----------

HOOK_ENTRY='{"matcher":"Skill","hooks":[{"type":"command","command":"upfront hook"}]}'

mkdir -p "$(dirname "$SETTINGS_FILE")"

# Create settings.json if it doesn't exist
if [ ! -f "$SETTINGS_FILE" ]; then
  info "Creating $SETTINGS_FILE"
  echo '{}' > "$SETTINGS_FILE"
fi

# Check if hook is already registered
if jq -e '.hooks.PostToolUse[]? | select(.matcher == "Skill") | .hooks[]? | select(.command == "upfront hook")' "$SETTINGS_FILE" >/dev/null 2>&1; then
  info "PostToolUse hook already registered — skipping."
else
  # Back up existing settings
  BACKUP="$SETTINGS_FILE.backup.$(date +%Y%m%d%H%M%S)"
  cp "$SETTINGS_FILE" "$BACKUP"
  info "Backed up settings to $BACKUP"

  # Add the hook using jq
  jq --argjson hook "$HOOK_ENTRY" '
    .hooks //= {} |
    .hooks.PostToolUse //= [] |
    .hooks.PostToolUse += [$hook]
  ' "$SETTINGS_FILE" > "$SETTINGS_FILE.tmp" && mv "$SETTINGS_FILE.tmp" "$SETTINGS_FILE"

  info "Registered PostToolUse hook in $SETTINGS_FILE"
fi

# ---------- done ----------

info "Installation complete."
echo ""
echo "  Verify:   upfront status"
echo "  Docs:     $SCRIPT_DIR/README.md"
echo ""
