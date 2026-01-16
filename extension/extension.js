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

export default class SendToLinuxExtension {
    constructor() {
        this._signalId = null;
        this._panelButton = null;
        this._notificationSource = null;
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

        const openItem = new PopupMenu.PopupMenuItem('Open Received Folder');
        openItem.connect('activate', () => this._openReceivedFolder());
        this._panelButton.menu.addMenuItem(openItem);

        this._panelButton.menu.addMenuItem(new PopupMenu.PopupSeparatorMenuItem());

        const startItem = new PopupMenu.PopupMenuItem('Start Backend');
        startItem.connect('activate', () => this._startBackend());
        this._panelButton.menu.addMenuItem(startItem);

        const stopItem = new PopupMenu.PopupMenuItem('Stop Backend');
        stopItem.connect('activate', () => this._stopBackend());
        this._panelButton.menu.addMenuItem(stopItem);

        Main.panel.addToStatusArea('send-to-linux', this._panelButton);
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
            this._panelButton.destroy();
            this._panelButton = null;
        }
    }

    _onItemReceived(_connection, _sender, _path, _iface, _signal, params) {
        const [id, type, value, size] = params.deepUnpack();
        const title = type === 'text' ? 'Text received' : 'File received';
        const body = type === 'text'
            ? value
            : `${value} (${size} bytes)`;

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

    _copyToClipboard(text) {
        St.Clipboard.get_default().set_text(St.ClipboardType.CLIPBOARD, text);
    }

    _openReceivedFolder() {
        const downloads = GLib.get_user_special_dir(GLib.UserDirectory.DIRECTORY_DOWNLOAD) ||
            GLib.get_home_dir();
        const folder = GLib.build_filenamev([downloads, 'SendToLinux']);
        const file = Gio.File.new_for_path(folder);
        Gio.AppInfo.launch_default_for_uri(file.get_uri(), null);
    }

    _startBackend() {
        this._runFlatpak(['flatpak', 'run', 'net.dgkim.SendToLinux.Backend']);
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
