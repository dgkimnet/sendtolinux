# Changelog

## Backend 1.1.0 - 2026-01-23
- Replace GSettings reads with CLI flags (env defaults still supported).
- Backend now accepts --bind/--port/--dir/--max-upload-mb on startup.

## Extension 1.1.0 - 2026-01-23
- Add preferences UI for backend settings.
- Start backend with CLI args derived from preferences.
- Add Preferences item to the panel menu.

## Backend 1.0.1 - 2026-01-21
- Emit one ItemReceived signal for multi-file uploads.
- Redirect GET /text and /file back to /.

## Backend 1.0.0 - 2026-01-20
- Initial release with local HTTP upload server (text + file), QR generation, and D-Bus events.

## Extension 1.0.0 - 2026-01-20
- Initial release with GNOME menu, QR display, and ItemReceived notifications/actions.
