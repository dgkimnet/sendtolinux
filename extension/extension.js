import Gio from 'gi://Gio';
import GLib from 'gi://GLib';
import St from 'gi://St';
import * as Main from 'resource:///org/gnome/shell/ui/main.js';
import * as MessageTray from 'resource:///org/gnome/shell/ui/messageTray.js';
import * as PanelMenu from 'resource:///org/gnome/shell/ui/panelMenu.js';
import * as PopupMenu from 'resource:///org/gnome/shell/ui/popupMenu.js';

const SERVICE_NAME = 'net.dgkim.SendToLinux';
const OBJECT_PATH = '/net/dgkim/SendToLinux';
const INTERFACE_NAME = 'net.dgkim.SendToLinux';
const SIGNAL_NAME = 'ItemReceived';
const EXTENSION_VERSION = '1.1.1';

export default class SendToLinuxExtension {
    constructor() {
        this._signalId = null;
        this._panelButton = null;
        this._notificationSource = null;
        this._menuOpenId = null;
        this._qrImage = null;
        this._urlLabel = null;
        this._statusLabel = null;
    }

    enable() {
        const connection = Gio.DBus.session;
        this._signalId = connection.signal_subscribe(
            SERVICE_NAME,
            INTERFACE_NAME,
            SIGNAL_NAME,
            OBJECT_PATH,
            null,
            Gio.DBusSignalFlags.NONE,
            this._onItemReceived.bind(this)
        );

        this._ensureNotificationSource();

        this._panelButton = new PanelMenu.Button(0.0, 'Send to Linux');
        const icon = new St.Icon({
            icon_name: 'send-to-symbolic',
            style_class: 'system-status-icon',
        });
        this._panelButton.add_child(icon);

        const qrItem = new PopupMenu.PopupBaseMenuItem({
            reactive: false,
            can_focus: false,
        });
        const qrBox = new St.BoxLayout({
            vertical: true,
            x_expand: true,
            y_expand: true,
        });
        this._statusLabel = new St.Label({ text: 'Loading status…' });
        this._urlLabel = new St.Label({ text: '' });
        this._qrImage = new St.Icon({
            icon_size: 160,
            visible: false,
        });
        qrBox.add_child(this._statusLabel);
        qrBox.add_child(this._urlLabel);
        qrBox.add_child(this._qrImage);
        qrItem.add_child(qrBox);
        this._panelButton.menu.addMenuItem(qrItem);

        this._panelButton.menu.addMenuItem(new PopupMenu.PopupSeparatorMenuItem());

        const versionItem = new PopupMenu.PopupMenuItem(`Extension v${EXTENSION_VERSION}`, {
            reactive: false,
            can_focus: false,
        });
        this._panelButton.menu.addMenuItem(versionItem);

        this._panelButton.menu.addMenuItem(new PopupMenu.PopupSeparatorMenuItem());

        const openItem = new PopupMenu.PopupMenuItem('Open Received Folder');
        openItem.connect('activate', () => this._openReceivedFolder());
        this._panelButton.menu.addMenuItem(openItem);

        const prefsItem = new PopupMenu.PopupMenuItem('Preferences');
        prefsItem.connect('activate', () => this._openPreferences());
        this._panelButton.menu.addMenuItem(prefsItem);

        this._panelButton.menu.addMenuItem(new PopupMenu.PopupSeparatorMenuItem());

        const startItem = new PopupMenu.PopupMenuItem('Start Backend');
        startItem.connect('activate', () => this._startBackend());
        this._panelButton.menu.addMenuItem(startItem);

        const stopItem = new PopupMenu.PopupMenuItem('Stop Backend');
        stopItem.connect('activate', () => this._stopBackend());
        this._panelButton.menu.addMenuItem(stopItem);

        Main.panel.addToStatusArea('send-to-linux', this._panelButton);

        this._menuOpenId = this._panelButton.menu.connect('open-state-changed', (_menu, isOpen) => {
            if (isOpen) {
                this._refreshStatus();
            }
        });
    }

    disable() {
        if (this._signalId !== null) {
            Gio.DBus.session.signal_unsubscribe(this._signalId);
            this._signalId = null;
        }

        if (this._notificationSource) {
            this._notificationSource.destroy();
            this._notificationSource = null;
        }

        if (this._panelButton) {
            if (this._menuOpenId !== null) {
                this._panelButton.menu.disconnect(this._menuOpenId);
                this._menuOpenId = null;
            }
            this._panelButton.destroy();
            this._panelButton = null;
        }
    }

    _onItemReceived(_connection, _sender, _path, _iface, _signal, params) {
        const [id, type, value, size] = params.deepUnpack();
        const sizeBytes = typeof size === 'bigint' ? Number(size) : size;
        const title = type === 'text' ? 'Text received' : 'File received';
        const body = type === 'text'
            ? value
            : `${value} (${sizeBytes} bytes)`;

        this._ensureNotificationSource();
        const notification = new MessageTray.Notification({
            source: this._notificationSource,
            title,
            body,
            urgency: MessageTray.Urgency.CRITICAL, // this make sure notification pops up
        });

        if (type === 'text') {
            notification.addAction('Copy', () => this._copyToClipboard(value));
        }

        notification.addAction('Open Folder', () => this._openReceivedFolder());
        this._notificationSource.addNotification(notification);
    }

    _ensureNotificationSource() {
        if (this._notificationSource) {
            return;
        }
        this._notificationSource = new MessageTray.Source({
            title: 'Send to Linux',
            iconName: 'send-to-symbolic',
        });
        this._notificationSource.connect('destroy', () => {
            this._notificationSource = null;
        });
        Main.messageTray.add(this._notificationSource);
    }

    _refreshStatus() {
        if (!this._statusLabel || !this._urlLabel || !this._qrImage) {
            return;
        }

        this._statusLabel.text = 'Checking backend…';
        this._urlLabel.text = '';
        this._qrImage.visible = false;
        this._qrImage.gicon = null;

        Gio.DBus.session.call(
            SERVICE_NAME,
            OBJECT_PATH,
            INTERFACE_NAME,
            'GetStatus',
            null,
            null,
            Gio.DBusCallFlags.NONE,
            -1,
            null,
            (connection, result) => {
                try {
                    const reply = connection.call_finish(result);
                    const [url, _port, running] = reply.deepUnpack();
                    if (!running || !url) {
                        this._statusLabel.text = 'Backend offline';
                        return;
                    }
                    this._statusLabel.text = 'Scan to upload';
                    this._urlLabel.text = url;
                    this._loadQrPath();
                } catch (err) {
                    this._statusLabel.text = 'Backend offline';
                }
            }
        );
    }

    _loadQrPath() {
        Gio.DBus.session.call(
            SERVICE_NAME,
            OBJECT_PATH,
            INTERFACE_NAME,
            'GetQrPath',
            null,
            null,
            Gio.DBusCallFlags.NONE,
            -1,
            null,
            (connection, result) => {
                try {
                    const reply = connection.call_finish(result);
                    const [path] = reply.deepUnpack();
                    if (!path) {
                        this._statusLabel.text = 'QR unavailable';
                        this._qrImage.visible = false;
                        this._qrImage.gicon = null;
                        return;
                    }
                    const file = Gio.File.new_for_path(path);
                    this._qrImage.gicon = new Gio.FileIcon({ file });
                    this._qrImage.visible = true;
                } catch (err) {
                    this._statusLabel.text = 'QR unavailable';
                }
            }
        );
    }

    _copyToClipboard(text) {
        St.Clipboard.get_default().set_text(St.ClipboardType.CLIPBOARD, text);
    }

    _openReceivedFolder() {
        const downloads = GLib.get_user_special_dir(GLib.UserDirectory.DIRECTORY_DOWNLOAD) ||
            GLib.get_home_dir();
        let folderName = 'SendToLinux';
        try {
            const settings = new Gio.Settings({
                schema_id: 'org.gnome.shell.extensions.send-to-linux',
            });
            const dir = settings.get_string('dir');
            if (dir) {
                folderName = GLib.path_get_basename(dir);
            }
        } catch (err) {
            // Fall back to default folder name.
        }
        const folder = GLib.build_filenamev([downloads, folderName]);
        const file = Gio.File.new_for_path(folder);
        Gio.AppInfo.launch_default_for_uri(file.get_uri(), null);
    }

    _openPreferences() {
        Main.extensionManager.openExtensionPrefs('send-to-linux@dgkim', '', {});
    }

    _startBackend() {
        const argv = ['flatpak', 'run', 'net.dgkim.SendToLinux.Backend'];
        const settings = new Gio.Settings({
            schema_id: 'org.gnome.shell.extensions.send-to-linux',
        });
        const bind = settings.get_string('bind');
        const port = settings.get_int('port');
        const dir = settings.get_string('dir');
        const maxUploadMb = settings.get_int('max-upload-mb');

        if (bind) {
            argv.push('--bind', bind);
        }
        if (Number.isInteger(port)) {
            argv.push('--port', String(port));
        }
        if (dir) {
            argv.push('--dir', dir);
        }
        if (Number.isInteger(maxUploadMb)) {
            argv.push('--max-upload-mb', String(maxUploadMb));
        }

        this._runFlatpak(argv);
    }

    _stopBackend() {
        this._runFlatpak(['flatpak', 'kill', 'net.dgkim.SendToLinux.Backend']);
    }

    _runFlatpak(argv) {
        try {
            const proc = Gio.Subprocess.new(argv, Gio.SubprocessFlags.NONE);
            proc.wait_check_async(null, (subprocess, res) => {
                try {
                    subprocess.wait_check_finish(res);
                } catch (err) {
                    Main.notify('Send to Linux', `Command failed: ${err.message}`);
                }
            });
        } catch (err) {
            Main.notify('Send to Linux', `Command failed: ${err.message}`);
        }
    }
}
