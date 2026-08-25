<!-- src/components/TrackingPage/PrintFileLayoutExport.vue -->
<template>
    <div class="container">
        <!-- ================= BODY ================= -->
        <div class="content">
            <div class="file-container">
                <div v-for="(item, index) in Print_data_list" :key="item.QRCode" class="page-break">
                    <div class="table-wrapper">
                        <table class="custom-table print-only">
                            <colgroup>
                                <col class="col-qr" />
                                <col class="col-label" />
                                <col class="col-small" />
                                <col class="col-value" />
                            </colgroup>

                            <tbody>
                                <tr class="r-50">
                                    <th class="qr-cell" rowspan="3">
                                        <!-- QR -->
                                        <img class="qr-img" :src="QrcodeURL_LIST[index]" alt="QR" />
                                    </th>

                                    <td class="cell label" colspan="2">Material Code</td>
                                    <td class="cell value value-big center v-middle">
                                        <b>{{ item.CLBH }}</b>
                                    </td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell label" colspan="2">LOTNo</td>
                                    <td class="cell value">
                                        <b class="value-compact">{{ item.LotNO }}</b>
                                    </td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell value qr-text v-middle" colspan="3">
                                        <b class="value-compact">{{ item.Barcode }}</b>
                                    </td>
                                </tr>

                                <tr class="r-70">
                                    <td class="cell label">Material Name</td>
                                    <td class="cell value" colspan="3">
                                        <span class="clamp clamp-4">
                                            <b>{{ item.YWPM }}</b>
                                        </span>
                                    </td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell label">Supplier Name</td>
                                    <td class="cell value" colspan="3">
                                        <b>{{ item.ZSYWJC }}</b>
                                    </td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell label">Po.NO</td>
                                    <td class="cell value" colspan="3">
                                        <b>{{ item.CGNO }}</b>
                                    </td>
                                </tr>

                                <tr class="r-26">
                                    <td class="cell label">RY</td>
                                    <td class="cell value" colspan="3">
                                        <span class="clamp clamp-2">
                                            <b>{{ item.Memo_RY }}</b>
                                        </span>
                                    </td>
                                </tr>

                                <tr class="r-39">
                                    <td class="cell label">SKU#/BUY</td>
                                    <td class="cell value" colspan="3">
                                        <span class="clamp clamp-3">
                                            <b>{{ item.Memo_Article }}</b>
                                        </span>
                                    </td>
                                </tr>

                                <tr class="r-15">
                                    <td class="cell label">Quantity</td>
                                    <td class="cell value center" colspan="3">
                                        <b>{{ item.RemQty }} YRD</b>
                                    </td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell label">Roll No</td>
                                    <td class="cell value" colspan="3">
                                        <b>{{ item.PACK1 }}</b>
                                    </td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell label">Date Received/Inv.No</td>
                                    <td class="cell value" colspan="3">
                                        <b>{{ item.MEMO }}</b>
                                    </td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell label"></td>
                                    <td class="cell value center">Pass</td>
                                    <td class="cell value center">Fail</td>
                                    <td class="cell value center">FIFO</td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell label">Lab Test Result</td>
                                    <td class="cell value"></td>
                                    <td class="cell value"></td>
                                    <td class="cell value fifo fifo-small" rowspan="5">
                                        <!-- <b>{{ item.FIFO }}</b> -->
                                    </td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell label">Visual Result</td>
                                    <td class="cell value"></td>
                                    <td class="cell value"></td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell label">Dimensional Result</td>
                                    <td class="cell value"></td>
                                    <td class="cell value"></td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell label">Defect &amp; Grading Check</td>
                                    <td class="cell value"></td>
                                    <td class="cell value"></td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell label">Checker/Signature</td>
                                    <td class="cell value center" colspan="2">
                                        <b>77576</b>
                                    </td>
                                </tr>

                                <tr class="r-20">
                                    <td class="cell doc-no" colspan="4">
                                        <b>Document Number: QIP/V/18/03</b>
                                    </td>
                                </tr>
                            </tbody>
                        </table>

                    </div>
                </div>
            </div>
        </div>

        <!-- ================= FOOTER ================= -->
        <div class="print-footer">
            <div class="footer-cluster footer-cluster--settings">
                <v-btn class="footer-btn footer-btn--mode" @click="togglePrintMode">
                    <svg class="footer-btn__icon" viewBox="0 0 24 24" aria-hidden="true">
                        <path d="M2 8.82a15 15 0 0 1 20 0" />
                        <path d="M5 12.86a10 10 0 0 1 14 0" />
                        <path d="M8.5 16.9a5 5 0 0 1 7 0" />
                        <circle cx="12" cy="20" r="1" fill="currentColor" stroke="none" />
                    </svg>
                    <span>{{ t('trackingPage.printLayout.modeLabel') }}: {{ printMode }}</span>
                </v-btn>

                <SettingPrinter />
            </div>

            <div class="footer-cluster footer-cluster--actions">
                <v-btn class="footer-btn footer-btn--print" @click="printContent">
                    <svg class="footer-btn__icon" viewBox="0 0 24 24" aria-hidden="true">
                        <path d="M6 9V3h12v6" />
                        <rect x="4" y="9" width="16" height="8" rx="1.5" />
                        <rect x="7" y="14" width="10" height="7" rx="1" />
                    </svg>
                    <span>{{ t('trackingPage.searchToolBar.print') }}</span>
                </v-btn>

                <v-btn class="footer-btn footer-btn--print2" @click="printContent2">
                    <svg class="footer-btn__icon" viewBox="0 0 24 24" aria-hidden="true">
                        <rect x="8" y="8" width="12" height="12" rx="1.5" />
                        <path d="M4 15V5a1 1 0 0 1 1-1h10" />
                    </svg>
                    <span>{{ t('trackingPage.printLayout.print2') }}</span>
                </v-btn>

                <v-btn class="footer-btn footer-btn--cancel" @click="ExitFileLayout">
                    <svg class="footer-btn__icon" viewBox="0 0 24 24" aria-hidden="true">
                        <path d="M6 6l12 12M18 6L6 18" />
                    </svg>
                    <span>{{ t('common.cancel') }}</span>
                </v-btn>
            </div>
        </div>
    </div>
    <v-snackbar v-model="snackbar" :timeout="snackbarTimeout" :color="snackbarColor" location="top" elevation="6">
        {{ snackbarText }}
    </v-snackbar>
</template>

<script setup>
import useTracking from '@/hooks/useTracking'
import qrcode from 'qrcode' // 引入 qrcode 库
import { computed, onBeforeMount, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import html2canvas from 'html2canvas'
import { usePrinter } from '@/hooks/usePrinter'

const { t } = useI18n()

const QrcodeURL_LIST = reactive([])
const onScanning_QRcode = ref(false)

const current_step = ref(0)

const Scan_QRcode_progress_rate = computed(() => {
    const rate = (100 / QrcodeURL_LIST.length) * current_step.value
    return Math.floor(rate) // 将结果向下取整
})

const { printZPL, printMode, setPrinterMode } = usePrinter()

const togglePrintMode = () => {
    setPrinterMode()
}

const printContent = async () => {
    if (printMode.value === 'wifi') {
        await printInPlace()
    } else {
        await printContentIP()
    }
}

const printContent2 = async () => {
    if (printMode.value === 'wifi') {
        await printInPlace2()
    } else {
        await printContentIP()
    }
}

onBeforeMount(() => {
    printMode.value = 'wifi'
    localStorage.setItem('printerMode', 'wifi')
})

const snackbar = ref(false)
const snackbarText = ref('')
const snackbarColor = ref('info')
const snackbarTimeout = ref(3000)

const showSnackbar = (text, color = 'info', timeout = 3000) => {
    snackbarText.value = text
    snackbarColor.value = color
    snackbarTimeout.value = timeout
    snackbar.value = true
}

const PRINT_BATCH_SIZE = 10
let currentPrintIndex = 0

const printInPlace = async () => {
    const tables = Array.from(document.querySelectorAll('.custom-table'))

    if (tables.length === 0) {
        showSnackbar(t('trackingPage.printLayout.noDataToPrint'), 'warning')
        return
    }

    const batchTables = tables.slice(
        currentPrintIndex,
        currentPrintIndex + PRINT_BATCH_SIZE
    )

    if (batchTables.length === 0) {
        showSnackbar(t('trackingPage.printLayout.allPrinted'), 'info')
        return
    }

    showSnackbar(
        t('trackingPage.printLayout.printingRange', {
            from: currentPrintIndex + 1,
            to: Math.min(currentPrintIndex + PRINT_BATCH_SIZE, tables.length),
        }),
        'info',
        2000
    )

    const clones = []

    // 1️⃣ Clone batch
    for (const table of batchTables) {
        const clone = table.cloneNode(true)
        clone.style.position = 'fixed'
        clone.style.left = '-9999px'
        clone.style.top = '0'
        clone.style.display = 'table'
        clone.style.visibility = 'visible'
        document.body.appendChild(clone)
        clones.push(clone)
    }

    await new Promise(r => setTimeout(r, 300))

    // 2️⃣ Render tuần tự
    const images = []
    for (const clone of clones) {
        const canvas = await html2canvas(clone, {
            scale: 1,
            backgroundColor: '#fff',
            useCORS: true,
        })
        images.push(canvas.toDataURL('image/png'))
    }

    // 3️⃣ Cleanup
    clones.forEach(c => document.body.removeChild(c))

    // 4️⃣ Print
    await printImages(images)

    currentPrintIndex += PRINT_BATCH_SIZE

    // 5️⃣ Thông báo trạng thái
    if (currentPrintIndex < tables.length) {
        showSnackbar(
            t('trackingPage.printLayout.printedProgress', {
                done: currentPrintIndex,
                total: tables.length,
            }),
            'success',
            4000
        )
    } else {
        showSnackbar(t('trackingPage.printLayout.allSheetsPrinted'), 'success', 4000)
        currentPrintIndex = 0
    }
}

const printImages = async (images) => {
    const iframe = document.createElement('iframe')
    iframe.style.position = 'fixed'
    iframe.style.right = '0'
    iframe.style.bottom = '0'
    iframe.style.width = '0'
    iframe.style.height = '0'
    iframe.style.border = '0'
    document.body.appendChild(iframe)

    const doc = iframe.contentWindow.document

    doc.open()
    doc.write(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<title>Print</title>
<style>
@page {
    size: 115mm 155mm;
    margin: 3mm;
}
html, body {
    margin: 0;
    padding: 0;
}
.print-page {
    width: 115mm;
    height: 155mm;
    page-break-after: always;
    display: flex;
    align-items: center;
    justify-content: center;
}
.print-page:last-child {
    page-break-after: auto;
}
.print-img {
    max-width: 100%;
    max-height: 100%;
}
</style>
</head>
<body>
${images
            .map(
                src => `
    <div class="print-page">
        <img src="${src}" class="print-img" />
    </div>
`,
            )
            .join('')}
</body>
</html>
`)
    doc.close()

    iframe.contentWindow.focus()
    iframe.contentWindow.print()

    setTimeout(() => {
        document.body.removeChild(iframe)
    }, 1000)
}


const printInPlace2 = async () => {
    // 1️⃣ Tạo iframe ẩn
    const iframe = document.createElement('iframe')
    iframe.style.position = 'fixed'
    iframe.style.right = '0'
    iframe.style.bottom = '0'
    iframe.style.width = '0'
    iframe.style.height = '0'
    iframe.style.border = '0'
    document.body.appendChild(iframe)

    const doc = iframe.contentWindow.document

    // 2️⃣ Lấy các table cần in
    const tables = document.querySelectorAll('.custom-table')

    // ❗ mỗi table = 1 page
    const pagesHTML = Array.from(tables)
        .map(t => `<div class="print-page">${t.outerHTML}</div>`)
        .join('')

    // 3️⃣ Lấy toàn bộ style hiện tại
    const styles = Array.from(
        document.querySelectorAll('style, link[rel="stylesheet"]')
    ).map(s => s.outerHTML).join('\n')

    // 4️⃣ Ghi HTML vào iframe
    doc.open()
    doc.write(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8" />
<title>Print</title>

${styles}

<style>
/* =========================
   PAGE SETUP
========================= */
@page {
    size: 115mm 155mm;
margin: 3mm;
}

/* ❗ body KHÔNG khóa height */
html, body {
    width: 115mm;
    margin: 0;
    padding: 0;
}

/* =========================
   PAGE WRAPPER
========================= */
.print-page {
    width: 115mm;
    height: 155mm;              /* ⭐ KHÓA PAGE */
    box-sizing: border-box;

    overflow: hidden;
    page-break-after: always;
    break-after: page;

    display: flex;
    align-items: flex-start;
    justify-content: center;
}

/* Trang cuối */
.print-page:last-child {
    page-break-after: auto;
    break-after: auto;
}

/* =========================
   TABLE
========================= */
.custom-table {
    width: 115mm;
    box-sizing: border-box;

    margin: 0;
    height: auto !important;
    min-height: unset !important;   /* ⭐ CHẶN CSS BLEED */
}

/* =========================
   PRINT SAFETY
========================= */
@media print {
    tr {
        page-break-inside: avoid;
        break-inside: avoid;
    }
}
</style>
</head>

<body>
${pagesHTML}
</body>
</html>
`)
    doc.close()

    // 5️⃣ Đợi ảnh load xong (QR)
    const imgs = doc.querySelectorAll('img')
    await Promise.all(
        Array.from(imgs).map(img => {
            if (img.complete) return Promise.resolve()
            return new Promise(resolve => {
                img.onload = img.onerror = resolve
            })
        })
    )

    // 6️⃣ Print
    iframe.contentWindow.focus()
    iframe.contentWindow.print()

    // 7️⃣ Cleanup
    setTimeout(() => {
        document.body.removeChild(iframe)
    }, 1000)
}

const printContentIP = async () => {
    onScanning_QRcode.value = true
    const container = document.querySelector('.file-container')

    // 等待图片加载完成
    await scrollToLoadImages(container)

    // 获取所有 .custom-table 元素
    const tables = document.querySelectorAll('.custom-table')

    // 使用 for...of 循环，确保异步操作的顺序执行
    for (const table of tables) {
        // 使用 html2canvas 渲染 DOM 到 Canvas
        const originalCanvas = await html2canvas(table)
        document.body.appendChild(originalCanvas) // 检查生成的 Canvas

        const targetWidth = 560 // 固定宽度
        const targetHeight = 700 // 固定高度

        const scaledCanvas = createScaledCanvas(
            originalCanvas,
            targetWidth,
            targetHeight,
        )

        // 将 Canvas 转换为 Base64 图像数据
        const base64Data = scaledCanvas.toDataURL('image/png')

        // 将 Base64 图像数据转换为 ZPL (Z64 格式)
        const zpl = await base64ToZPL(base64Data, targetWidth, targetHeight)
        console.log('Generated ZPL:', zpl)

        const success = await printZPL(zpl)

        if (!success) {
            alert(t('trackingPage.printLayout.printFailedZPL'))
        }

        document.body.removeChild(originalCanvas) // 清理 Canvas
    }
    onScanning_QRcode.value = false // 结束扫描
}

// 创建缩放后的 Canvas
const createScaledCanvas = (originalCanvas, targetWidth, targetHeight) => {
    const scaledCanvas = document.createElement('canvas')
    const ctx = scaledCanvas.getContext('2d')

    // 设置新的 Canvas 尺寸
    scaledCanvas.width = targetWidth
    scaledCanvas.height = targetHeight

    // 设置 Canvas 背景为透明
    ctx.clearRect(0, 0, targetWidth, targetHeight)

    // 计算缩放比例
    const scaleWidth = targetWidth / originalCanvas.width
    const scaleHeight = targetHeight / originalCanvas.height
    const scale = Math.min(scaleWidth, scaleHeight)

    // 计算居中位置
    const x = (targetWidth - originalCanvas.width * scale) / 2
    const y = (targetHeight - originalCanvas.height * scale) / 2

    // 绘制缩放后的图像，使其填满整个 Canvas
    ctx.drawImage(
        originalCanvas,
        0,
        0,
        originalCanvas.width,
        originalCanvas.height,
        x, // 新的 x 坐标
        y, // 新的 y 坐标
        originalCanvas.width * scale, // 缩放后的宽度
        originalCanvas.height * scale, // 缩放后的高度
    )

    return scaledCanvas
}

async function base64ToZPL(
    base64String,
    imageWidth,
    imageHeight,
    dpi = 300,
    padding = 20,
) {
    const labelHeight = Math.floor((dpi * 13.5) / 2.54) // 标签高度
    const labelWidth = Math.floor((dpi * 10.5) / 2.54) // 标签宽度

    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')

    // 计算缩放比例
    const scaleRatio = Math.min(
        (labelWidth - 2 * padding) / imageWidth, // 考虑左右 padding
        labelHeight / imageHeight,
    )

    // 缩放后的宽度和高度
    const scaledWidth = Math.floor(imageWidth * scaleRatio)
    const scaledHeight = Math.floor(imageHeight * scaleRatio)

    // 设置 Canvas 尺寸
    canvas.width = scaledWidth
    canvas.height = scaledHeight

    const img = new Image()
    img.src = base64String

    return new Promise((resolve, reject) => {
        img.onload = () => {
            ctx.clearRect(0, 0, canvas.width, canvas.height)
            ctx.drawImage(img, 0, 0, scaledWidth, scaledHeight)

            const imageData = ctx.getImageData(
                0,
                0,
                canvas.width,
                canvas.height,
            )
            const data = imageData.data
            const binaryData = []

            for (let y = 0; y < scaledHeight; y++) {
                let rowBytes = 0
                for (let x = 0; x < scaledWidth; x++) {
                    const idx = (y * scaledWidth + x) * 4
                    const r = data[idx]
                    const g = data[idx + 1]
                    const b = data[idx + 2]
                    const gray = 0.299 * r + 0.587 * g + 0.114 * b
                    const bit = gray > 150 ? 0 : 1
                    rowBytes = (rowBytes << 1) | bit

                    if (x % 8 === 7 || x === scaledWidth - 1) {
                        binaryData.push(rowBytes)
                        rowBytes = 0
                    }
                }
            }

            const hexArray = binaryData.map((byte) =>
                ('0' + byte.toString(16)).slice(-2).toUpperCase(),
            )
            const hexString = hexArray.join('')
            const bytesPerRow = Math.ceil(scaledWidth / 8)
            const totalBytes = bytesPerRow * scaledHeight

            // 计算左右 padding 后的偏移
            const xOffset =
                padding +
                Math.floor((labelWidth - scaledWidth - 2 * padding) / 2) // 水平偏移
            const yOffset = Math.floor((labelHeight - scaledHeight) / 2) // 垂直偏移

            const zpl = `^XA^FO${xOffset},${yOffset}^GFA,${totalBytes},${totalBytes},${bytesPerRow},${hexString}^FS^XZ`

            console.log('ZPL size:', new Blob([zpl]).size / 1024, 'KB')

            resolve(zpl)
        }

        img.onerror = (err) => reject(new Error('Image loading failed: ' + err))
    })
}

onBeforeMount(() => {
    Print_data_list.value.forEach((element) => {
        generateQRCodeBase64(element.Barcode)
    })
})

const { Print_data_list, print_file_layout_show } = useTracking()

const ExitFileLayout = () => {
    print_file_layout_show.value = false
    current_step.value = 0
}

const generateQRCodeBase64 = async (QRCode) => {
    try {
        const dataURL = await qrcode.toDataURL(QRCode)
        QrcodeURL_LIST.push(dataURL)
    } catch (error) {
        console.error('Failed to generate QR code:', error)
        return '' // 失败时返回空字符串
    }
}
const scrollToLoadImages = async (container) => {
    const step = window.innerWidth // Step size for scrolling
    const containerWidth = container.scrollWidth
    let firstStep = true // Flag to check if it's the first step
    container.scrollTo(0, 0)
    for (let i = 0; i <= containerWidth; i += step) {
        container.scrollTo(i, 0)
        if (firstStep) {
            // First step delay 1 second
            await new Promise((resolve) => setTimeout(resolve, 1000))
            firstStep = false // Set flag to false after the first step
        } else {
            // Subsequent steps delay 0.1 second
            await new Promise((resolve) => setTimeout(resolve, 1000))
        }
        current_step.value += 1
    }
    onScanning_QRcode.value = false
    current_step.value = 0
    container.scrollTo(0, 0) // Reset scroll position
}
</script>

<style lang="scss" scoped>
/* =========================================================
   GLOBAL PRINT COLOR
========================================================= */
* {
    -webkit-print-color-adjust: exact !important;
    print-color-adjust: exact !important;
}

/* =========================================================
   ROOT CONTAINER
========================================================= */
.container {
    position: fixed;
    inset: 0;
    width: 100%;
    height: 100vh;

    display: flex;
    flex-direction: column;

    background-color: rgba(255, 255, 255, 0.9);
    z-index: 4;
    overflow: hidden;
}

/* =========================================================
   CONTENT (SCROLL AREA)
========================================================= */
.content {
    flex: 1;
    overflow: hidden;
}

.file-container {
    height: 100%;
    overflow: auto;
    padding: 20px 0;

    display: flex;
    flex-direction: column;
    align-items: center;

    font-family: serif;
}

/* =========================================================
   TABLE BASE
========================================================= */
.custom-table {
    width: 115mm;
    min-height: 155mm;

    border-collapse: collapse;
    table-layout: fixed;
    box-sizing: border-box;

    background-color: #fff;
    border: 3px solid #000;
}

.custom-table th,
.custom-table td {
    border: 1px solid #000;
    padding: 1px;
    box-sizing: border-box;
}

/* Disable select */
table th,
table td {
    user-select: none;
}

/* =========================================================
   COLUMN WIDTHS
========================================================= */
.col-qr {
    width: 32mm;
}

.col-label {
    width: 10mm;
}

.col-small {
    width: 10mm;
}

.col-value {
    width: 32mm;
}

/* =========================================================
   ROW HEIGHT HELPERS
========================================================= */
.r-15 {
    height: 15px;
}

.r-20 {
    height: 20px;
}

.r-26 {
    height: 50px;
}

.r-39 {
    height: 32px;
}

.r-50 {
    height: 50px;
}

.r-70 {
    height: 70px;
}

/* =========================================================
   CELL TYPES
========================================================= */
.cell {
    padding-left: 2px;
    font-size: 11pt;
    line-height: 1.2;
    vertical-align: middle;
}

.label {
    padding-left: 4px;
    white-space: nowrap;
}

.value {
    font-weight: 400;
    padding-left: 5px;
    vertical-align: middle;
}

.center {
    text-align: center;
}

.top {
    vertical-align: top;
}

.value>b {
    display: inline-flex;
    align-items: center;
    min-height: 100%;
}


/* =========================================================
   SPECIAL TEXT
========================================================= */
.value-big {
    font-size: 23px;
    line-height: 23px;
}

.value-compact {
    display: block;
    font-size: 15px;
    line-height: 10px;
    height: 10px;
}

.qr-text {
    font-size: 15px;
    line-height: 1.1;
    padding-left: 5px;
}

.doc-no {
    text-align: right;
    padding-right: 6px;
}

/* =========================================================
   QR
========================================================= */
.qr-cell {
    padding: 0;
    text-align: center;
}

.qr-img {
    width: 38mm;
    height: 38mm;
    object-fit: contain;
    display: block;
    margin: 0 auto;
}

/* =========================================================
   MULTI LINE CLAMP (Material Name / RY / SKU)
========================================================= */
.clamp {
    width: 100%;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    word-wrap: break-word;
    white-space: normal;
    font-size: 14px;
}

/* match old logic */
.clamp-4 {
    -webkit-line-clamp: 4;
    line-height: 19px;
    height: 78px;

    display: flex;
    align-items: center;
}

.clamp-3 {
    -webkit-line-clamp: 3;
    line-height: 12px;
    height: 32px;

    display: flex;
    align-items: center;
}

.clamp-2 {
    -webkit-line-clamp: 2;
    line-height: 12px;
    height: 50px;

    display: flex;
    align-items: center;
}

/* =========================================================
   FIFO
========================================================= */
.fifo {
    font-size: 64pt;
    line-height: 1;
    text-align: center;
}

.fifo-small {
    font-size: 65px;
    line-height: 65px;
}

/* =========================================================
   FOOTER (BUTTON BAR)
========================================================= */
.print-footer {
    --pf-ink: #1E293B;
    --pf-ink-soft: #475569;
    --pf-muted: #94A3B8;
    --pf-line: #E2E8F0;
    --pf-line-strong: #CBD5E1;
    --pf-surface: #FFFFFF;
    --pf-accent: #1D4ED8;
    --pf-accent-2: #3B82F6;
    --pf-accent-glow: rgba(29, 78, 216, 0.35);
    --pf-emerald: #059669;
    --pf-emerald-dark: #047857;
    --pf-emerald-glow: rgba(5, 150, 105, 0.32);

    height: 76px;
    flex-shrink: 0;

    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 0 20px;

    background: var(--pf-surface);
    border-top: 1px solid var(--pf-line);
    box-shadow: 0 -8px 24px -16px rgba(15, 23, 42, 0.18);

    position: sticky;
    bottom: 0;
    z-index: 10;
}

.footer-cluster {
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 0;
}

/* =========================================================
   BUTTON STYLE
========================================================= */
.footer-btn {
    display: inline-flex !important;
    align-items: center;
    gap: 8px;
    height: 44px;
    padding: 0 18px !important;
    border-radius: 12px !important;
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 0.01em;
    text-transform: none !important;
    box-shadow: none !important;
    transition: filter 0.15s ease, transform 0.15s ease, background 0.15s ease,
        border-color 0.15s ease, color 0.15s ease, box-shadow 0.15s ease;
}
.footer-btn :deep(.v-btn__content) {
    display: inline-flex;
    align-items: center;
    gap: 8px;
}
.footer-btn__icon {
    width: 17px;
    height: 17px;
    flex: 0 0 auto;
    fill: none;
    stroke: currentColor;
    stroke-width: 2;
    stroke-linecap: round;
    stroke-linejoin: round;
}

/* Chip trung tính: chế độ in (WIFI/USB) */
.footer-btn--mode {
    background: #F1F5F9 !important;
    color: var(--pf-ink-soft) !important;
    border: 1px solid var(--pf-line);
}
.footer-btn--mode:hover {
    background: #E2E8F0 !important;
    border-color: var(--pf-line-strong);
    color: var(--pf-ink) !important;
}

/* Hành động chính: In */
.footer-btn--print {
    background: linear-gradient(135deg, var(--pf-accent), var(--pf-accent-2)) !important;
    color: #fff !important;
    box-shadow: 0 10px 22px -10px var(--pf-accent-glow),
        inset 0 1px 0 rgba(255, 255, 255, 0.16) !important;
}
.footer-btn--print:hover {
    filter: brightness(1.08);
    transform: translateY(-1px);
}

/* Hành động phụ: In 2 */
.footer-btn--print2 {
    background: linear-gradient(135deg, var(--pf-emerald), var(--pf-emerald-dark)) !important;
    color: #fff !important;
    box-shadow: 0 10px 22px -10px var(--pf-emerald-glow),
        inset 0 1px 0 rgba(255, 255, 255, 0.16) !important;
}
.footer-btn--print2:hover {
    filter: brightness(1.08);
    transform: translateY(-1px);
}

/* Hủy: ghost button */
.footer-btn--cancel {
    background: transparent !important;
    color: var(--pf-ink-soft) !important;
    border: 1px solid var(--pf-line);
}
.footer-btn--cancel:hover {
    background: #FEF2F2 !important;
    border-color: #FCA5A5;
    color: #DC2626 !important;
}
.footer-btn--cancel:hover .footer-btn__icon {
    color: #DC2626;
}

@media (max-width: 640px) {
    .print-footer {
        flex-wrap: wrap;
        height: auto;
        padding: 12px 14px;
        gap: 10px;
    }
    .footer-cluster {
        flex-wrap: wrap;
        width: 100%;
        justify-content: center;
    }
    .footer-btn span {
        white-space: nowrap;
    }
}

/* =========================================================
   PRINT MODE
========================================================= */
@media print {
    @page {
        size: 115mm 155mm;
        margin: 3mm;
    }

    html,
    body {
        width: 115mm;
        height: 155mm;
        margin: 0;
        padding: 0;
    }

    .custom-table {
        width: 115mm !important;
        height: 155mm !important;
        margin: 0 !important;
        page-break-after: always;
    }

    .custom-table:last-child {
        page-break-after: auto !important;
    }

    tr {
        break-inside: avoid !important;
    }

    td,
    th {
        white-space: normal !important;
        word-break: break-word !important;
        vertical-align: middle !important;
    }
}

/* =========================================================
   PROGRESS
========================================================= */
.v-progress-circular {
    margin: 1rem;
}

.page-break {
    margin-bottom: 20px;
}
</style>