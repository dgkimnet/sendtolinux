#!/usr/bin/env bash
set -euo pipefail

APP_ID="net.dgkim.SendToLinux.Backend"

flatpak uninstall --user -y "${APP_ID}"
