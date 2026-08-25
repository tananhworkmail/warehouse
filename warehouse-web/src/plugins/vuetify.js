/**
 * plugins/vuetify.js
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
// import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi-svg'
import {
    mdiDatabaseSearch,
    mdiPrinter,
    mdiCalendarRange,
    mdiCheck,
    mdiEye,
    mdiEyeOff,
    mdiQrcodeScan,
    mdiPower,
    mdiFileFind,
    mdiPencil,
    mdiCloseCircle,
    mdiFullscreen,
    mdiFullscreenExit,
    mdiExitToApp,
    mdiScanHelper,
    mdiMenuUp,
    mdiMenuDown,
    mdiQrcode,
    mdiClipboardText,
    mdiSquareEditOutline,
    mdiClose,
    mdiIp,
    mdiBarcode,
    mdiTruck,
} from '@mdi/js'
import { zhHans } from 'vuetify/locale'
export default createVuetify({
    locale: {
        messages: { zhHans },
    },
    icons: {
        defaultSet: 'mdi',
        aliases: {
            ...aliases,
            databasesearch: mdiDatabaseSearch,
            printer: mdiPrinter,
            calendar: mdiCalendarRange,
            checkedCircle: mdiCheck,
            eye: mdiEye,
            eyeOff: mdiEyeOff,
            qrcodescan: mdiQrcodeScan,
            power: mdiPower,
            file: mdiFileFind,
            pencil: mdiPencil,
            clear: mdiCloseCircle,
            fullscreen: mdiFullscreen,
            fullscreenExit: mdiFullscreenExit,
            exit: mdiExitToApp,
            pay: mdiScanHelper,
            qrcode: mdiQrcode,
            arrowup: mdiMenuUp,
            arrowdown: mdiMenuDown,
            order: mdiClipboardText,
            edit: mdiSquareEditOutline,
            scanhelper: mdiScanHelper,
            close: mdiClose,
            ip: mdiIp,
            barcde: mdiBarcode,
            truck: mdiTruck,
        },
        sets: {
            mdi,
        },
    },
})