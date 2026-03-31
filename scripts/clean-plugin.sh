#!/usr/bin/env bash
# Remove all traces of the upfront plugin from Claude Code.
# Run this before a fresh /plugin marketplace add + /plugin install cycle.

set -euo pipefail

SETTINGS="$HOME/.claude/settings.json"
KNOWN="$HOME/.claude/plugins/known_marketplaces.json"
INSTALLED="$HOME/.claude/plugins/installed_plugins.json"

# Remove cache and marketplace dirs
rm -rf ~/.claude/plugins/cache/thinkupfront ~/.claude/plugins/marketplaces/thinkupfront 2>/dev/null && echo "✓ Deleted cache/marketplace dirs" || true

# Remove install marker so next session sends plugin_install event
rm -f ~/.upfront/.installed 2>/dev/null && echo "✓ Removed install marker" || true

# Remove from settings.json (enabledPlugins + extraKnownMarketplaces)
if [ -f "$SETTINGS" ]; then
    python3 -c "
import json, sys
with open('$SETTINGS') as f:
    data = json.load(f)
changed = False
if 'upfront@thinkupfront' in data.get('enabledPlugins', {}):
    del data['enabledPlugins']['upfront@thinkupfront']
    changed = True
if 'thinkupfront' in data.get('extraKnownMarketplaces', {}):
    del data['extraKnownMarketplaces']['thinkupfront']
    changed = True
if changed:
    with open('$SETTINGS', 'w') as f:
        json.dump(data, f, indent=2)
    print('✓ Cleaned settings.json')
else:
    print('· settings.json already clean')
"
fi

# Remove from known_marketplaces.json
if [ -f "$KNOWN" ]; then
    python3 -c "
import json
with open('$KNOWN') as f:
    data = json.load(f)
if 'thinkupfront' in data:
    del data['thinkupfront']
    with open('$KNOWN', 'w') as f:
        json.dump(data, f, indent=2)
    print('✓ Cleaned known_marketplaces.json')
else:
    print('· known_marketplaces.json already clean')
"
fi

# Remove from installed_plugins.json
if [ -f "$INSTALLED" ]; then
    python3 -c "
import json
with open('$INSTALLED') as f:
    data = json.load(f)
if 'upfront@thinkupfront' in data.get('plugins', {}):
    del data['plugins']['upfront@thinkupfront']
    with open('$INSTALLED', 'w') as f:
        json.dump(data, f, indent=2)
    print('✓ Cleaned installed_plugins.json')
else:
    print('· installed_plugins.json already clean')
"
fi

echo ""
echo "Done. Now run:"
echo "  claude plugin marketplace add ThinkUpfront/Upfront"
echo "  claude plugin install upfront"
echo "Then restart Claude Code."
