#!/usr/bin/env bash
set -euo pipefail

EXT_DIR="${HOME}/.local/share/gnome-shell/extensions/send-to-linux@dgkim"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

mkdir -p "${EXT_DIR}"
rsync -a --delete \
  --exclude "install.sh" \
  --exclude "README.md" \
  "${SCRIPT_DIR}/" "${EXT_DIR}/"

echo "Installed extension to ${EXT_DIR}"
