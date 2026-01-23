#!/usr/bin/env bash
set -euo pipefail

EXT_DIR="${HOME}/.local/share/gnome-shell/extensions/send-to-linux@dgkim"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

mkdir -p "${EXT_DIR}"
rsync -a --delete \
  --exclude "install.sh" \
  --exclude "README.md" \
  "${SCRIPT_DIR}/" "${EXT_DIR}/"

if command -v glib-compile-schemas >/dev/null 2>&1; then
  glib-compile-schemas "${EXT_DIR}/schemas"
else
  echo "Warning: glib-compile-schemas not found; prefs may not load until schemas are compiled."
fi

echo "Installed extension to ${EXT_DIR}"
