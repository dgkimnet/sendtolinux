# Switching to a Single-App UI (AppIndicator)

This note summarizes the earlier discussion about replacing the GNOME Shell extension with an AppIndicator-based UI.

## Key differences vs GNOME Shell extension
- GNOME Shell panel button/popover can **only** be implemented in a GNOME Shell extension (GJS).
- AppIndicator uses the tray system (StatusNotifier/AppIndicator) and **cannot** add a native top-panel button/popover.
- On GNOME Shell, AppIndicator icons are **not visible by default** unless the user installs an AppIndicator/KStatusNotifierItem extension.
- AppIndicator menus are more limited (simple menu items) and can’t replicate a rich popover (e.g., QR image layout) as cleanly.
- AppIndicator requires a separate GUI process (GTK/Qt/etc.) to host the tray icon.

## What remains intact
- Backend features (HTTP server, file I/O, D-Bus service, QR generation) are unchanged.
- Notifications still work via `org.freedesktop.Notifications` from either the GUI or backend process.

## “Copy to clipboard” in notifications
- The copy action is **not** tied to the GNOME Shell extension; it depends on the notification server’s support for actions.
- GNOME Shell supports actions, so a GUI app can still provide a “Copy” button.
- On some desktops/notification daemons, action buttons may be hidden or ignored.

## Language choice for the backend
- Python or Rust can implement **all current backend functionality** (D-Bus, HTTP uploads, file handling, QR generation).
- The GNOME Shell panel UI **cannot** be rewritten in Python/Rust; it must remain a GNOME Shell extension if you want a top-panel button.

## Desktop targeting summary
- GNOME-only: keep the GNOME Shell extension for the best UX.
- Cross-desktop: AppIndicator is more portable but less GNOME-native and not visible by default on GNOME.
