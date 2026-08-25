/**
 * plugins/index.js
 *
 * Automatically included in `./src/main.js`
 */

// Plugins
import vuetify from './vuetify'
import i18n from '@/hooks/i18n'

export function registerPlugins(app) {
    app.use(vuetify)
    app.use(i18n)
}