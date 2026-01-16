# Development Guide

## Prereqs
- GNOME Shell (GNOME 45+ recommended)
- Go 1.22+
- `gnome-extensions` CLI (optional but helpful)

## Repo layout
- `backend/` Go code
- `extension/` GNOME extension (metadata.json, extension.js, etc.)
- `docs/` design docs

## Backend dev
- `go run ./backend` should:
  - start HTTP server on a port (configurable)
  - register D-Bus name `net.dgkim.SendToLinux`
  - implement GetStatus / GetRecentItems
  - emit ItemReceived on upload

Recommended env vars:
- `STL_BIND=0.0.0.0`
- `STL_PORT=8000` (or 0 for random)
- `STL_DIR=/home/<user>/Downloads/SendToLinux`
- `STL_MAX_UPLOAD_MB=100`

## Extension dev
Install symlink (dev):
1) Copy or symlink `extension/` to:
   - `~/.local/share/gnome-shell/extensions/send-to-linux@dgkim/`
2) Restart GNOME Shell:
   - X11: Alt+F2 -> `r`
   - Wayland: log out/in (or `gnome-shell --replace` in nested session)
3) Enable:
   - `gnome-extensions enable send-to-linux@dgkim`

Logs:
- `journalctl --user -f | grep -i gnome-shell`
- Looking Glass: Alt+F2 -> `lg`

## Testing from another device
- Ensure iOS and Linux are on same Wi-Fi
- Open QR from GNOME top bar
- Upload text and file, confirm notification and saved path

## MVP checklist
- [ ] backend: D-Bus name registers
- [ ] backend: GET / serves upload page
- [ ] backend: POST /text and POST /file work
- [ ] backend: saves to Downloads/SendToLinux
- [ ] backend: emits ItemReceived on success
- [ ] extension: subscribes signals, notifies
- [ ] extension: copy-to-clipboard action works
- [ ] extension: shows QR for URL

