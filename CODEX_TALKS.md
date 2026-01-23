# Codex Talks Log

## 2026-01-16
- User asked to start the project step by step and first create a minimal backend implementation.
- Agreed to implement a minimal Go D-Bus service (`net.dgkim.SendToLinux`) with `GetStatus`, `GetRecentItems`, and an optional test signal controlled by `STL_EMIT_TEST=1`.

## 2026-01-16 (cont.)
- User requested `.gitignore` and `backend/README.md`.
- Added Go-appropriate ignore rules and a backend README with run instructions and next steps.

## 2026-01-16 (cont. 2)
- User requested moving `go.mod` into `backend/` since the project uses both Go and GJS.
- Moved `go.mod` from repo root to `backend/go.mod`.

## 2026-01-16 (cont. 3)
- User requested updating `backend/README.md` to note the Go module location.
- Added a module layout section with guidance on running commands from `backend/`.

## 2026-01-16 (cont. 4)
- User requested implementing the backend HTTP server.
- Added GET `/` and POST `/text`, saving text to `~/Downloads/SendToLinux`, emitting `ItemReceived`, and wiring status/recent items.

## 2026-01-16 (cont. 5)
- User requested modularizing the backend: HTTP handlers and D-Bus service moved into `backend/internal/`.
- Added `internal/dbussvc` for D-Bus logic and `internal/httpserver` for HTTP server and handlers.

## 2026-01-16 (cont. 6)
- User requested documentation updates to reflect current backend modularization and HTTP support.
- Updated `backend/README.md` and `docs/DEV.md` to reflect module layout and `/text` endpoint.

## 2026-01-16 (cont. 7)
- User requested step 3: GNOME extension D-Bus subscribe + notification on ItemReceived.
- Added minimal extension in `extension/` with metadata and D-Bus signal subscription that triggers GNOME notifications.

## 2026-01-16 (cont. 8)
- User requested `extension/README.md` and a build/install script.
- Added extension README and `extension/install.sh` for local install.

## 2026-01-16 (cont. 9)
- User requested an uninstall script for the extension.
- Added `extension/uninstall.sh` to remove the installed extension directory.

## 2026-01-16 (cont. 10)
- User requested a task bar icon and a command to open the received folder.
- Added a GNOME panel button with a menu item to open Downloads/SendToLinux.

## 2026-01-16 (cont. 11)
- User requested Flatpak-packaged backend and extension start/stop controls.
- Added Flatpak manifest at `backend/flatpak/net.dgkim.SendToLinux.Backend.json` and extension menu actions to run/kill the Flatpak.
- Updated backend and extension docs plus `docs/DEV.md` with packaging notes.

## 2026-01-16 (cont. 12)
- User reported Flatpak build failed with `go` not found.
- Updated Flatpak manifest to use the Go SDK extension and appended PATH; noted requirement in `docs/DEV.md`.

## 2026-01-16 (cont. 13)
- User hit Flatpak build failure due to Go module download blocked.
- Vendored Go dependencies and updated Flatpak build to use `-mod=vendor`; documented in `docs/DEV.md` and `backend/README.md`.

## 2026-01-16 (cont. 14)
- User requested ignoring Flatpak build artifacts and the vendored dependencies.
- Updated `.gitignore` to exclude Flatpak build directories and `backend/vendor/`.

## 2026-01-16 (cont. 15)
- User requested Flatpak README and build/install/uninstall scripts.
- Added `backend/flatpak/README.md` plus `build.sh`, `install.sh`, and `uninstall.sh`.

## 2026-01-16 (cont. 16)
- User requested updating `docs/TODO.md` to reflect progress.
- Marked completed backend HTTP/D-Bus items and extension notification/icon/open-folder tasks; added Flatpak publish workflow to Post-MVP.

## 2026-01-16 (cont. 17)
- User requested moving hard-coded HTML into template and CSS files.
- Added HTML/CSS assets under `backend/internal/httpserver/assets` and updated the HTTP server to serve them.

## 2026-01-16 (cont. 18)
- User requested HTML response page for POST results with a success message.
- Updated template to render a notice and serve it after `/text` saves successfully.

## 2026-01-16 (cont. 19)
- User requested notification actions for copy-to-clipboard and open received folder.
- Added MessageTray notifications with Copy (text only) and Open Folder actions.

## 2026-01-16 (cont. 20)
- User reported `MessageTray.SystemNotificationSource` not a constructor.
- Switched to `MessageTray.Source` for GNOME 45 compatibility.

## 2026-01-16 (cont. 21)
- User reported MessageTray.Source constructor expects a properties object.
- Updated to pass `{ title, iconName }`.

## 2026-01-16 (cont. 22)
- User reported Flatpak runtime missing HTML assets.
- Embedded HTTP assets with Go embed, with optional STL_ASSET_DIR override for local dev.

## 2026-01-16 (cont. 23)
- User reported MessageTray notification error; constructor signature mismatch.
- Updated notification construction to use the properties object form.

## 2026-01-16 (cont. 24)
- User reported error on notification source notify call.
- Switched to `showNotification`, added NotificationPolicy and source lifecycle handling.

## 2026-01-19
- User requested backend-generated QR PNG and extension popover display.
- Added backend QR generation via `qrencode`, exposed `GetQrPath`, and wired the extension popover to show the QR and URL.
- Updated D-Bus and dev docs for the new QR path method and dependency.

## 2026-01-19 (cont.)
- User requested URL hostnames include `.local`.
- Updated backend hostname resolution to append `.local` when needed.

## 2026-01-19 (cont. 2)
- User requested backend file uploads from client.
- Added POST `/file` multipart handling, file upload form, and updated docs checklists.

## 2026-01-19 (cont. 3)
- User requested documentation review and updates.
- Aligned README, DEV, ARCHITECTURE, TODO, and component READMEs with the current QR, upload, and popover status.

## 2026-01-20
- User requested backend support for receiving multiple files and UI updates.
- Updated `/file` handler to process multiple multipart files and adjusted the upload form to allow multi-select.
- Added module version strings, created `CHANGELOG.md`, and displayed backend/extension versions in the UI.

## 2026-01-21
- User requested only one notification when multiple files are received.
- Updated backend to emit a single ItemReceived signal for multi-file uploads.
- User requested GET /text and /file redirect to the index when reopened in a browser.
- Added redirects for GET /text and GET /file to /.
- User requested backend version bump and changelog update.
- Updated backend version to 1.0.1 and added changelog entry.

## 2026-01-23
- User requested GSettings-based preferences UI for the extension.
- Added GSettings schema, prefs window, and install-time schema compilation.
- Updated extension and dev docs for preferences UI and config keys.
- User requested backend read GSettings directly and update Flatpak manifest.
- Added backend config loader with GSettings reads (env override fallback) and installed/compiled schema in Flatpak.
- Updated Flatpak permissions to access host dconf and added Preferences menu item + prefs import fixes.
- User requested backend use CLI flags instead of GSettings reads.
- Switched backend config to flags with env defaults and updated extension to pass prefs as Flatpak args.
- Bumped backend/extension versions to 1.1.0 and updated changelog.
- User clarified `--dir` should be a folder name under Downloads; updated backend resolution and prefs UI wording.
- Fixed prefs UI property compatibility and made Open Received Folder respect the configured folder name.
