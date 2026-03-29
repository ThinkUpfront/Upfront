#!/usr/bin/env bash
set -euo pipefail

# install-upfront.sh — Download and install the upfront binary.
# Usage: curl -fsSL https://raw.githubusercontent.com/brennhill/upfront/main/install-upfront.sh | bash

REPO="brennhill/upfront"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

info()  { printf '\033[1;34m==>\033[0m %s\n' "$1"; }
warn()  { printf '\033[1;33mwarning:\033[0m %s\n' "$1"; }
die()   { printf '\033[1;31merror:\033[0m %s\n' "$1" >&2; exit 1; }

# Detect OS and architecture.
detect_platform() {
  OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
  ARCH="$(uname -m)"

  case "$OS" in
    darwin) OS="darwin" ;;
    linux)  OS="linux" ;;
    *)      die "Unsupported OS: $OS" ;;
  esac

  case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)             die "Unsupported architecture: $ARCH" ;;
  esac
}

# Get the latest release tag from GitHub.
get_latest_version() {
  if command -v curl >/dev/null 2>&1; then
    VERSION="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed 's/.*"v\(.*\)".*/\1/')"
  elif command -v wget >/dev/null 2>&1; then
    VERSION="$(wget -qO- "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed 's/.*"v\(.*\)".*/\1/')"
  else
    die "curl or wget is required"
  fi

  if [ -z "$VERSION" ]; then
    die "Could not determine latest version. Check https://github.com/$REPO/releases"
  fi
}

download_and_install() {
  ARCHIVE="upfront_${VERSION}_${OS}_${ARCH}.tar.gz"
  URL="https://github.com/$REPO/releases/download/v${VERSION}/${ARCHIVE}"

  info "Downloading upfront v${VERSION} for ${OS}/${ARCH}..."

  TMPDIR="$(mktemp -d)"
  trap 'rm -rf "$TMPDIR"' EXIT

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$URL" -o "$TMPDIR/$ARCHIVE"
  else
    wget -q "$URL" -O "$TMPDIR/$ARCHIVE"
  fi

  tar -xzf "$TMPDIR/$ARCHIVE" -C "$TMPDIR"

  mkdir -p "$INSTALL_DIR"
  mv "$TMPDIR/upfront" "$INSTALL_DIR/upfront"
  chmod +x "$INSTALL_DIR/upfront"

  info "Installed upfront to $INSTALL_DIR/upfront"
}

verify_installation() {
  if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
    warn "$INSTALL_DIR is not on your PATH"
    warn "Add this to your shell profile:"
    warn "  export PATH=\"$INSTALL_DIR:\$PATH\""
  fi

  if command -v upfront >/dev/null 2>&1; then
    info "Installed: $(upfront version)"
  fi
}

register_hook() {
  SETTINGS_FILE="$HOME/.claude/settings.json"
  HOOK_ENTRY='{"matcher":"Skill","hooks":[{"type":"command","command":"upfront hook"}]}'

  if ! command -v jq >/dev/null 2>&1; then
    warn "jq not found — skipping Claude Code hook registration"
    warn "See https://github.com/$REPO#manual-install for manual setup"
    return
  fi

  mkdir -p "$(dirname "$SETTINGS_FILE")"

  if [ ! -f "$SETTINGS_FILE" ]; then
    echo '{}' > "$SETTINGS_FILE"
  fi

  if jq -e '.hooks.PostToolUse[]? | select(.matcher == "Skill") | .hooks[]? | select(.command == "upfront hook")' "$SETTINGS_FILE" >/dev/null 2>&1; then
    info "PostToolUse hook already registered"
  else
    BACKUP="$SETTINGS_FILE.backup.$(date +%Y%m%d%H%M%S)"
    cp "$SETTINGS_FILE" "$BACKUP"
    info "Backed up settings to $BACKUP"

    jq --argjson hook "$HOOK_ENTRY" '
      .hooks //= {} |
      .hooks.PostToolUse //= [] |
      .hooks.PostToolUse += [$hook]
    ' "$SETTINGS_FILE" > "$SETTINGS_FILE.tmp" && mv "$SETTINGS_FILE.tmp" "$SETTINGS_FILE"

    info "Registered PostToolUse hook in Claude Code"
  fi
}

main() {
  info "Installing upfront — audit trail for AI-assisted feature definition"
  echo ""

  detect_platform
  get_latest_version
  download_and_install
  verify_installation
  register_hook

  echo ""
  info "Done! Run 'upfront status' to verify."
  echo ""
  echo "  Docs:  https://github.com/$REPO"
  echo "  Book:  https://upfront.dev/book"
  echo ""
}

main
