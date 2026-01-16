# Repository Guidelines

## Project Structure & Module Organization

This repository currently contains planning and design documentation only. The core implementation directories referenced in docs (such as `backend/` and `extension/`) are not present yet. Use the docs in `docs/` as the source of truth for architecture and development intent:
- `README.md` for scope and goals
- `docs/ARCHITECTURE.md` for system flow and storage rules
- `docs/DBUS.md` for D-Bus contract details
- `docs/DEV.md` for development expectations
- `docs/TODO.md` for the MVP checklist

If you add code, place it under the paths described in `docs/DEV.md` (e.g., `backend/` for Go, `extension/` for the GNOME Shell extension) and keep documentation in `docs/` for consistency.

## Build, Test, and Development Commands

No build or test scripts exist yet. When implementation lands, align with the commands implied by `docs/DEV.md`, for example:
- `go run ./backend` to run the Go daemon locally.
- `gnome-extensions enable send-to-linux@dgkim` to enable the extension during development.

Document any new commands in `docs/DEV.md` and keep them minimal and reproducible.

## Coding Style & Naming Conventions

There is no code style defined yet. If you introduce Go or GJS source:
- Follow standard Go formatting (`gofmt`) and idiomatic package naming.
- Use GNOME extension conventions for filenames (`extension.js`, `metadata.json`).
- Keep identifiers aligned with existing D-Bus names like `org.dgkim.SendToLinux` as described in `docs/ARCHITECTURE.md`.

## Testing Guidelines

No testing framework is configured. If tests are added, document the framework and commands in `docs/DEV.md`, and use clear naming (e.g., `*_test.go` for Go tests).

## Commit & Pull Request Guidelines

The repository has no commits yet, so there are no established message conventions. If you contribute, use clear, imperative commit subjects (e.g., "Add D-Bus status method") and include a short PR description that references relevant docs or checklist items in `docs/TODO.md`.

## Security & Configuration Tips

Respect the LAN-only design noted in `docs/ARCHITECTURE.md` (private IP ranges, configurable bind/port). Use the environment variables listed in `docs/DEV.md` (e.g., `STL_BIND`, `STL_PORT`, `STL_DIR`) for local configuration.
