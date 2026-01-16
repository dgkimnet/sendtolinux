#!/usr/bin/env bash
set -euo pipefail

EXT_DIR="${HOME}/.local/share/gnome-shell/extensions/send-to-linux@dgkim"

if [ -d "${EXT_DIR}" ]; then
  rm -rf "${EXT_DIR}"
  echo "Removed extension from ${EXT_DIR}"
else
  echo "Extension not installed at ${EXT_DIR}"
fi
