import Adw from 'gi://Adw';
import Gio from 'gi://Gio';
import Gtk from 'gi://Gtk';
import { ExtensionPreferences, gettext as _ } from 'resource:///org/gnome/Shell/Extensions/js/extensions/prefs.js';

export default class SendToLinuxPreferences extends ExtensionPreferences {
    fillPreferencesWindow(window) {
        Adw.init();

        const settings = this.getSettings();

        window.set_title(_('Send to Linux'));
        window.set_default_size(520, 420);

        const page = new Adw.PreferencesPage({
            title: _('Backend'),
            icon_name: 'send-to-symbolic',
        });
        window.add(page);

        const serverGroup = new Adw.PreferencesGroup({
            title: _('Server'),
            description: _('Settings applied when the backend starts.'),
        });
        page.add(serverGroup);

        const bindRow = new Adw.EntryRow({
            title: _('Bind address'),
            text: settings.get_string('bind'),
        });
        settings.bind('bind', bindRow, 'text', Gio.SettingsBindFlags.DEFAULT);
        serverGroup.add(bindRow);

        const portAdjustment = new Gtk.Adjustment({
            lower: 0,
            upper: 65535,
            step_increment: 1,
            page_increment: 10,
            value: settings.get_int('port'),
        });
        const portRow = new Adw.SpinRow({
            title: _('Port'),
            adjustment: portAdjustment,
            climb_rate: 1,
            digits: 0,
        });
        settings.bind('port', portRow, 'value', Gio.SettingsBindFlags.DEFAULT);
        serverGroup.add(portRow);

        const dirRow = new Adw.EntryRow({
            title: _('Save folder name (saved under Downloads, e.g., SendToLinux)'),
            text: settings.get_string('dir'),
        });
        settings.bind('dir', dirRow, 'text', Gio.SettingsBindFlags.DEFAULT);
        serverGroup.add(dirRow);

        const uploadGroup = new Adw.PreferencesGroup({
            title: _('Uploads'),
        });
        page.add(uploadGroup);

        const sizeAdjustment = new Gtk.Adjustment({
            lower: 1,
            upper: 2048,
            step_increment: 1,
            page_increment: 10,
            value: settings.get_int('max-upload-mb'),
        });
        const maxUploadRow = new Adw.SpinRow({
            title: _('Max upload size (MB)'),
            adjustment: sizeAdjustment,
            climb_rate: 1,
            digits: 0,
        });
        settings.bind('max-upload-mb', maxUploadRow, 'value', Gio.SettingsBindFlags.DEFAULT);
        uploadGroup.add(maxUploadRow);
    }
}
