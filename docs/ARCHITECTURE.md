# Architecture

## Overview
The system is split to keep GNOME Shell stable:
- GNOME extension handles UI + clipboard + notifications.
- Go backend handles HTTP uploads + file I/O + QR generation.

The two communicate via D-Bus on the **session bus**.

## Data flow
1) User opens popover from GNOME top bar.
2) Extension asks backend for current upload URL + QR PNG path.
3) Extension loads the QR PNG and shows it with the URL.
4) iOS visits URL and uploads text or one/more files.
5) Backend saves data to disk and emits `ItemReceived` D-Bus signal.
6) Extension receives signal:
   - shows GNOME notification
   - updates popover recent list
   - enables "Copy" action for text / "Open folder" for file

## Storage
Default directory:
- `~/Downloads/SendToLinux/`

Naming:
- Text: `text-YYYYMMDD-HHMMSS.txt`
- File: keep original filename, add suffix if conflict

Metadata log (optional MVP+):
- `~/Downloads/SendToLinux/index.jsonl`

## Networking
- Bind: `0.0.0.0:<port>` to accept LAN access
- URL: `http://<lan-ip>:<port>/`
- LAN-only access enforcement (recommended):
  - reject non-private ranges (10/8, 172.16/12, 192.168/16)

## Process model
- Extension spawns backend if not running (optional)
- Backend also runnable standalone
- Backend registers D-Bus name:
  - `net.dgkim.SendToLinux`

## Failure modes
- Backend down: extension shows "backend offline" state
- IP change: extension refreshes URL via D-Bus `GetStatus()`
- Large uploads: backend enforces size limits
