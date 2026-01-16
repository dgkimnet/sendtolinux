# TODO

## MVP
- Backend
  - [ ] Implement HTTP server with:
    - [ ] GET `/` upload UI (text + file)
    - [ ] POST `/text`
    - [ ] POST `/file` (multipart)
  - [ ] Save into `~/Downloads/SendToLinux/`
  - [ ] D-Bus service:
    - [ ] RequestName `net.dgkim.SendToLinux`
    - [ ] Export methods `GetStatus`, `GetRecentItems`
    - [ ] Emit signal `ItemReceived`
  - [ ] Determine LAN IP for URL
  - [ ] Enforce LAN-only (optional but recommended)

- Extension
  - [ ] Add top bar icon
  - [ ] Popover UI:
    - [ ] show QR image (and URL text)
    - [ ] recent items list
  - [ ] D-Bus client:
    - [ ] call `GetStatus`
    - [ ] subscribe `ItemReceived`
  - [ ] Notification on receive
  - [ ] Copy-to-clipboard for text
  - [ ] Open folder / open file for file items

## Post-MVP
- [ ] systemd --user service installation (optional)
- [ ] Better conflict-free naming strategy for saved files
- [ ] History persistence (`index.jsonl`)
- [ ] Upload progress / size limits in UI
- [ ] Multiple NICs handling (Wi-Fi + VPN)
- [ ] Optional PIN confirmation on desktop notification
- [ ] iOS share-sheet friendly UI (minimal, fast)

