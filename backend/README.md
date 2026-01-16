# Backend (Go)

## Purpose
The backend provides the D-Bus service and (later) the HTTP server that receives uploads. It runs as a standalone Go process and communicates with the GNOME extension over the session bus.

## Current Status
Steps 1 and 2 are partially implemented:
- D-Bus name registration (`net.dgkim.SendToLinux`)
- `GetStatus` and `GetRecentItems` methods (in-memory)
- Optional `ItemReceived` test signal
- HTTP server with GET `/` and POST `/text` (text saved to Downloads)

## Run
```bash
go run ./backend
```

Emit a one-time test signal on startup:
```bash
STL_EMIT_TEST=1 go run ./backend
```

## Next Steps
- Add POST `/file` multipart upload support.
- Persist recent items to disk (optional).
- Enforce LAN-only access if needed.

## Module Layout
The Go module lives in `backend/go.mod`. Run commands from the repo root using `go -C backend ...` or from `backend/` directly.

## Package Layout
- `internal/dbussvc` contains the D-Bus service and in-memory status/recent items.
- `internal/httpserver` contains HTTP handlers and server startup.
