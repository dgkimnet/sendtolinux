# Codex Talks Log

## 2025-02-14
- User asked to start the project step by step and first create a minimal backend implementation.
- Agreed to implement a minimal Go D-Bus service (`net.dgkim.SendToLinux`) with `GetStatus`, `GetRecentItems`, and an optional test signal controlled by `STL_EMIT_TEST=1`.

## 2025-02-14 (cont.)
- User requested `.gitignore` and `backend/README.md`.
- Added Go-appropriate ignore rules and a backend README with run instructions and next steps.

## 2025-02-14 (cont. 2)
- User requested moving `go.mod` into `backend/` since the project uses both Go and GJS.
- Moved `go.mod` from repo root to `backend/go.mod`.

## 2025-02-14 (cont. 3)
- User requested updating `backend/README.md` to note the Go module location.
- Added a module layout section with guidance on running commands from `backend/`.

## 2025-02-14 (cont. 4)
- User requested implementing the backend HTTP server.
- Added GET `/` and POST `/text`, saving text to `~/Downloads/SendToLinux`, emitting `ItemReceived`, and wiring status/recent items.

## 2025-02-14 (cont. 5)
- User requested modularizing the backend: HTTP handlers and D-Bus service moved into `backend/internal/`.
- Added `internal/dbussvc` for D-Bus logic and `internal/httpserver` for HTTP server and handlers.

## 2025-02-14 (cont. 6)
- User requested documentation updates to reflect current backend modularization and HTTP support.
- Updated `backend/README.md` and `docs/DEV.md` to reflect module layout and `/text` endpoint.
