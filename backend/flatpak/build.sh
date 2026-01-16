#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FLATPAK_DIR="${ROOT_DIR}/flatpak"
REPO_DIR="${FLATPAK_DIR}/repo"
BUILD_DIR="${FLATPAK_DIR}/build-dir"
BUNDLE_PATH="${FLATPAK_DIR}/send-to-linux-backend.flatpak"
MANIFEST="${FLATPAK_DIR}/net.dgkim.SendToLinux.Backend.json"
APP_ID="net.dgkim.SendToLinux.Backend"

flatpak-builder --force-clean --repo="${REPO_DIR}" "${BUILD_DIR}" "${MANIFEST}"
flatpak build-bundle "${REPO_DIR}" "${BUNDLE_PATH}" "${APP_ID}"

echo "Built bundle: ${BUNDLE_PATH}"
