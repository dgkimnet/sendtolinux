# GNOME Extension

## Purpose
The extension listens for D-Bus signals from the backend and shows GNOME notifications when new items arrive. This is the minimal implementation for step 3 of the MVP.

## Layout
- `metadata.json` extension metadata and UUID
- `extension.js` D-Bus signal subscription and notifications
- `install.sh` helper to install the extension into the user profile

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

## Remove
```bash
rm -rf ~/.local/share/gnome-shell/extensions/send-to-linux@dgkim
```
