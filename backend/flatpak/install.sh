#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUNDLE_PATH="${ROOT_DIR}/flatpak/send-to-linux-backend.flatpak"

if [ ! -f "${BUNDLE_PATH}" ]; then
  echo "Bundle not found: ${BUNDLE_PATH}"
  echo "Run ./backend/flatpak/build.sh first."
  exit 1
fi

flatpak install --user --reinstall "${BUNDLE_PATH}"
