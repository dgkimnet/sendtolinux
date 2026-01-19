# Backend (Go)

## Purpose
The backend provides the D-Bus service and HTTP server that receives uploads. It runs as a standalone Go process and communicates with the GNOME extension over the session bus.

## Current Status
Steps 1 and 2 are partially implemented:
- D-Bus name registration (`net.dgkim.SendToLinux`)
- `GetStatus`, `GetRecentItems`, and `GetQrPath` methods (in-memory)
- Optional `ItemReceived` test signal
- HTTP server with GET `/` and POST `/text` + `/file` (saved to Downloads)
- QR PNG generation for the upload URL (requires `qrencode`)

## Run
```bash
go -C backend run .
```

Emit a one-time test signal on startup:
```bash
STL_EMIT_TEST=1 go -C backend run .
```

## Next Steps
- Persist recent items to disk (optional).
- Enforce LAN-only access if needed.

## Module Layout
The Go module lives in `backend/go.mod`. Run commands from the repo root using `go -C backend ...` or from `backend/` directly.

## Package Layout
- `internal/dbussvc` contains the D-Bus service and in-memory status/recent items.
- `internal/httpserver` contains HTTP handlers and server startup.

## Flatpak

Install golang in `flatpak` environment
```bash
flatpak install --user flathub org.freedesktop.Sdk.Extension.golang//23.08
```

The backend can be packaged as a Flatpak app with ID `net.dgkim.SendToLinux.Backend` using:
```bash
flatpak-builder --install --user build-dir flatpak/net.dgkim.SendToLinux.Backend.json
```

After installation, the extension can start/stop the backend via the panel menu.

If Flatpak builds fail fetching modules, vendor dependencies first:
```bash
go -C backend mod vendor
```
