# TODO

## MVP
- Backend
  - [ ] Implement HTTP server with:
    - [x] GET `/` upload UI (text + file)
    - [x] POST `/text`
    - [ ] POST `/file` (multipart)
  - [x] Save into `~/Downloads/SendToLinux/`
  - [ ] D-Bus service:
    - [x] RequestName `net.dgkim.SendToLinux`
    - [x] Export methods `GetStatus`, `GetRecentItems`
    - [x] Emit signal `ItemReceived`
  - [ ] Determine LAN IP for URL
  - [ ] Enforce LAN-only (optional but recommended)

- Extension
  - [x] Add top bar icon
  - [ ] Popover UI:
    - [ ] show QR image (and URL text)
    - [ ] recent items list
  - [ ] D-Bus client:
    - [ ] call `GetStatus`
    - [x] subscribe `ItemReceived`
  - [x] Notification on receive
  - [ ] Copy-to-clipboard for text
  - [x] Open folder / open file for file items

## Post-MVP
- [ ] systemd --user service installation (optional)
- [ ] Flatpak publish workflow (bundle/repo distribution)
- [ ] Better conflict-free naming strategy for saved files
- [ ] History persistence (`index.jsonl`)
- [ ] Upload progress / size limits in UI
- [ ] Multiple NICs handling (Wi-Fi + VPN)
- [ ] Optional PIN confirmation on desktop notification
- [ ] iOS share-sheet friendly UI (minimal, fast)
