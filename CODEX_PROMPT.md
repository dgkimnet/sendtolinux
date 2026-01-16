# Codex Prompt

You are implementing "Send to Linux" for GNOME:
- A Go backend that runs an HTTP server for LAN uploads (text + files)
- A GNOME Shell extension (GJS) that shows a tray icon, displays a QR code for the upload URL, and shows GNOME notifications
- Real-time communication via D-Bus session bus (no polling)

Implement MVP in small steps:
1) backend: minimal D-Bus service (RequestName, GetStatus, Emit test signal)
2) backend: HTTP server with GET `/` page and POST `/text` that saves to Downloads and emits ItemReceived
3) extension: D-Bus subscribe + notification on ItemReceived
4) extension: popover UI with URL text and a placeholder QR (then implement real QR)
5) backend: file upload POST `/file` multipart, save file, emit signal with path
6) extension: recent items list + copy-to-clipboard action

Constraints:
- Keep GNOME Shell stable: extension must not do heavy I/O
- Use godbus/dbus/v5 in Go
- Use Gio.DBusProxy or Gio.DBusConnection in GJS
- Default save dir: ~/Downloads/SendToLinux/

