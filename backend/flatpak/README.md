# Backend Flatpak Packaging

## Build
Build a local repo and bundle:
```bash
./backend/flatpak/build.sh
```

## Install
Install the bundle for the current user:
```bash
./backend/flatpak/install.sh
```

## Uninstall
Remove the app:
```bash
./backend/flatpak/uninstall.sh
```

## Notes
- App ID: `net.dgkim.SendToLinux.Backend`
- Manifest: `backend/flatpak/net.dgkim.SendToLinux.Backend.json`
- Requires the Go SDK extension: `org.freedesktop.Sdk.Extension.golang//23.08`
