# GNOME Extension

## Purpose
The extension listens for D-Bus signals from the backend and shows GNOME notifications when new items arrive. It also shows backend status, URL, and QR code in the panel popover.
The panel menu also shows the extension version.

## Layout
- `metadata.json` extension metadata and UUID
- `extension.js` D-Bus signal subscription and notifications
- `install.sh` helper to install the extension into the user profile
- `uninstall.sh` helper to remove the extension from the user profile

## Build
No build step is required for this simple JavaScript extension.

## Install (dev)
```bash
./extension/install.sh
```

Restart GNOME Shell and enable the extension:
- X11: Alt+F2 -> `r`
- Wayland: log out/in

Then run:
```bash
gnome-extensions enable send-to-linux@dgkim
```

## Backend Control
Use the panel menu to start or stop the backend Flatpak (`net.dgkim.SendToLinux.Backend`). The popover will show "Backend offline" until the service is running.

## Remove
```bash
rm -rf ~/.local/share/gnome-shell/extensions/send-to-linux@dgkim
```
