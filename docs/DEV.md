# Development Guide

## Prereqs
- GNOME Shell (GNOME 45+ recommended)
- Go 1.22+
- `gnome-extensions` CLI (optional but helpful)
- `qrencode` (backend QR PNG generation)

## Repo layout
- `backend/` Go code
- `extension/` GNOME extension (metadata.json, extension.js, etc.)
- `docs/` design docs

## Backend dev
- `go -C backend run .` should:
  - start HTTP server on a port (configurable)
  - register D-Bus name `net.dgkim.SendToLinux`
  - implement GetStatus / GetRecentItems / GetQrPath (in-memory)
  - emit ItemReceived on `/text` and `/file` uploads
  - accept multiple files in a single `/file` multipart request
  - generate QR PNG for the upload URL

Flatpak packaging:
- Manifest: `backend/flatpak/net.dgkim.SendToLinux.Backend.json`
- Build/install: `flatpak-builder --install --user build-dir backend/flatpak/net.dgkim.SendToLinux.Backend.json`
- Requires: `org.freedesktop.Sdk.Extension.golang` for Go toolchain inside Flatpak SDK
- Vendored deps for offline builds: run `go -C backend mod vendor` before building

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
- [x] backend: D-Bus name registers
- [x] backend: GET / serves upload page
- [x] backend: POST /text and POST /file work
- [x] backend: saves to Downloads/SendToLinux
- [x] backend: emits ItemReceived on success
- [x] extension: subscribes signals, notifies
- [x] extension: copy-to-clipboard action works
- [x] extension: shows QR for URL
