<template>
  <section class="tracking-results" aria-labelledby="tracking-results-title">
    <header class="results-header">
      <div>
        <p class="results-eyebrow">N8D · TRACKING DATA</p>
        <h2 id="tracking-results-title" class="results-title">
          {{ t('trackingPage.leftTable.resultsTitle') }}
        </h2>
      </div>

      <span class="results-count">
        {{ t('trackingPage.leftTable.recordCount', { count: Track_data_list.length }) }}
      </span>
    </header>

    <div class="results-summary" aria-label="Tracking totals">
      <div class="summary-item summary-item--in">
        <span>{{ t('trackingPage.leftTable.qtyIn') }}</span>
        <strong>{{ formatQuantity(totalQty) }}</strong>
      </div>
      <div class="summary-item summary-item--out">
        <span>{{ t('trackingPage.leftTable.qtyOut') }}</span>
        <strong>{{ formatQuantity(totalOutQty) }}</strong>
      </div>
      <div class="summary-item summary-item--stock">
        <span>{{ t('trackingPage.leftTable.qtyRemaining') }}</span>
        <strong>{{ formatQuantity(totalRemainingQty) }}</strong>
      </div>
    </div>

    <div class="results-body" aria-live="polite">
      <!-- Container xử lý cả scroll ngang lẫn scroll dọc của bảng (không còn chứa pagination) -->
      <div class="table-shell desktop-results" tabindex="0" :aria-label="scrollHint">
        <table class="tracking-table">
          <thead>
            <tr>
              <th v-for="header in headers" :key="header.key" scope="col">
                {{ header.title }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="isloading" class="state-row">
              <td :colspan="headers.length">
                <div class="result-state">
                  <span class="loading-spinner" aria-hidden="true" />
                  <strong>{{ t('trackingPage.leftTable.loading') }}</strong>
                </div>
              </td>
            </tr>

            <tr v-else-if="trackingError" class="state-row">
              <td :colspan="headers.length">
                <div class="result-state result-state--error" role="alert">
                  <span class="state-symbol" aria-hidden="true">!</span>
                  <strong>{{ t('trackingPage.leftTable.errorTitle') }}</strong>
                  <p>{{ trackingError }}</p>
                  <button class="retry-button" type="button" @click="refreshTrackingData">
                    {{ retryLabel }}
                  </button>
                </div>
              </td>
            </tr>

            <tr v-else-if="!Track_data_list.length" class="state-row">
              <td :colspan="headers.length">
                <div class="result-state">
                  <svg class="empty-icon" viewBox="0 0 40 40" aria-hidden="true">
                    <circle cx="18" cy="18" r="11" />
                    <path d="m26 26 8 8M13 18h10M18 13v10" />
                  </svg>
                  <strong>{{ t('trackingPage.leftTable.emptyTitle') }}</strong>
                  <p>{{ t('trackingPage.leftTable.emptyHint') }}</p>
                </div>
              </td>
            </tr>

            <tr v-for="(row, rowIndex) in pagedList" v-else :key="rowKey(row, rowIndex)">
              <td v-for="header in headers" :key="header.key" :class="`cell-${header.key}`">
                <button
                  v-if="header.key === 'File_Name' && row[header.key]"
                  class="report-link"
                  type="button"
                  :title="t('trackingPage.tableCell.viewImqcReport')"
                  @click="openIMQCReport(row[header.key])"
                >
                  {{ row[header.key] }}
                </button>
                <span v-else>{{ displayValue(row[header.key]) }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Mobile Results giữ nguyên như cũ -->
      <div class="mobile-results">
        <div v-if="isloading" class="result-state">
          <span class="loading-spinner" aria-hidden="true" />
          <strong>{{ t('trackingPage.leftTable.loading') }}</strong>
        </div>

        <div v-else-if="trackingError" class="result-state result-state--error" role="alert">
          <span class="state-symbol" aria-hidden="true">!</span>
          <strong>{{ t('trackingPage.leftTable.errorTitle') }}</strong>
          <p>{{ trackingError }}</p>
          <button class="retry-button" type="button" @click="refreshTrackingData">
            {{ retryLabel }}
          </button>
        </div>

        <div v-else-if="!Track_data_list.length" class="result-state">
          <svg class="empty-icon" viewBox="0 0 40 40" aria-hidden="true">
            <circle cx="18" cy="18" r="11" />
            <path d="m26 26 8 8M13 18h10M18 13v10" />
          </svg>
          <strong>{{ t('trackingPage.leftTable.emptyTitle') }}</strong>
          <p>{{ t('trackingPage.leftTable.emptyHint') }}</p>
        </div>

        <article v-for="(row, rowIndex) in pagedList" v-else :key="`mobile-${rowKey(row, rowIndex)}`" class="result-card">
          <header>
            <div>
              <small>Barcode</small>
              <strong class="barcode-value">{{ displayValue(row.Barcode) }}</strong>
            </div>
            <span>{{ displayValue(row.GSBH) }}</span>
          </header>
          <div class="result-card__material">
            <strong>{{ displayValue(row.CLBH) }}</strong>
            <p>{{ displayValue(row.YWPM) }}</p>
          </div>
          <dl>
            <div><dt>{{ t('trackingPage.leftTable.qtyIn') }}</dt><dd>{{ displayValue(row.Qty) }}</dd></div>
            <div><dt>{{ t('trackingPage.leftTable.qtyOut') }}</dt><dd>{{ displayValue(row.DQty) }}</dd></div>
            <div><dt>{{ t('trackingPage.leftTable.qtyRemaining') }}</dt><dd>{{ displayValue(row.RemQty) }}</dd></div>
            <div><dt>Pack</dt><dd>{{ displayValue(row.Pack) }}</dd></div>
            <div><dt>CGNO</dt><dd>{{ displayValue(row.CGNO) }}</dd></div>
            <div><dt>CFMDATE</dt><dd>{{ displayValue(row.CFMDATE || row.CFMDATE1) }}</dd></div>
          </dl>
          <button
            v-if="row.File_Name"
            class="report-link report-link--mobile"
            type="button"
            @click="openIMQCReport(row.File_Name)"
          >
            IMQC Report: {{ row.File_Name }}
          </button>
        </article>

        <div v-if="!isloading && !trackingError && Track_data_list.length" class="table-pagination table-pagination--mobile">
          <div class="page-controls">
            <button class="page-btn" type="button" :disabled="currentPage === 1" @click="goToPage(currentPage - 1)">‹</button>
            <button
              v-for="p in pageNumbers"
              :key="p === '...' ? `mobile-ellipsis-${Math.random()}` : `mobile-page-${p}`"
              class="page-btn"
              :class="{ active: p === currentPage, ellipsis: p === '...' }"
              type="button"
              :disabled="p === '...'"
              @click="p !== '...' && goToPage(p)"
            >{{ p }}</button>
            <button class="page-btn" type="button" :disabled="currentPage === totalPages" @click="goToPage(currentPage + 1)">›</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Thanh phân trang luôn hiển thị, KHÔNG nằm trong vùng cuộn ngang/dọc của bảng -->
    <div v-if="!isloading && !trackingError && Track_data_list.length" class="table-pagination desktop-results">
      <span class="page-info">{{ pageRangeLabel }}</span>
      <div class="page-controls">
        <button class="page-btn" type="button" :disabled="currentPage === 1" @click="goToPage(currentPage - 1)">‹</button>
        <button
          v-for="p in pageNumbers"
          :key="p === '...' ? `ellipsis-${Math.random()}` : `page-${p}`"
          class="page-btn"
          :class="{ active: p === currentPage, ellipsis: p === '...' }"
          type="button"
          :disabled="p === '...'"
          @click="p !== '...' && goToPage(p)"
        >{{ p }}</button>
        <button class="page-btn" type="button" :disabled="currentPage === totalPages" @click="goToPage(currentPage + 1)">›</button>
      </div>
    </div>

    <p class="scroll-hint">{{ scrollHint }}</p>
  </section>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import useTracking from '@/hooks/useTracking'

const API_PREFIX = (import.meta.env.VITE_API_URL || '/api/v1').replace(/\/+$/, '')

const { t, locale } = useI18n()
const { Track_data_list, isloading, trackingError, refreshTrackingData } = useTracking()

const columnLabels = [
  { key: 'Memo_RY', zh: '訂單', en: 'Order', vi: 'Lệnh' },
  { key: 'GSBH', zh: '廠別', en: 'Factory', vi: 'Xưởng' },
  { key: 'CKBH', zh: '倉庫', en: 'Storehouse', vi: 'Kho' },
  { key: 'MAKE', zh: '貨架號', en: 'Shelf No', vi: 'Kệ' },
  { key: 'CLBH', zh: '材料編號', en: 'Material Code', vi: 'Mã vật tư' },
  { key: 'YWPM', zh: '材料名稱', en: 'Material Name', vi: 'Tên vật tư' },
  { key: 'DWBH', zh: '單位', en: 'Unit', vi: 'Đơn vị' },
  { key: 'Pack', zh: '包裝', en: 'Package', vi: 'Cuộn' },
  { key: 'Qty', zh: '數量', en: 'Quantity', vi: 'SL nhập' },
  { key: 'DQty', zh: '發料數量', en: 'Dispatch Quantity', vi: 'SL xuất' },
  { key: 'RemQty', zh: '入庫數量', en: 'Inbound Quantity', vi: 'SL tồn' },
  { key: 'CGNO', zh: '採購單號', en: 'PO Number', vi: 'Phiếu mua' },
  { key: 'DOCNO', zh: '貨櫃INV', en: 'Container', vi: 'DOCNO' },
  { key: 'MEMO', zh: '備註', en: 'Remark', vi: 'MEMO' },
  { key: 'ZSBH', zh: '廠商', en: 'Supplier Code', vi: 'Mã N.C.U' },
  { key: 'ZSYWJC', zh: '廠商名稱', en: 'Supplier Name', vi: 'Tên N.C.U' },
  { key: 'CFMID', zh: '確認ID', en: 'Scan ID', vi: 'ID Scan' },
  { key: 'CFMDATE', zh: '確認日期', en: 'Scan Date', vi: 'Ngày Scan' },
  { key: 'File_Name', zh: 'IMQC', en: 'Report', vi: 'Report' },
]

const headers = computed(() => columnLabels.map((column) => ({
  key: column.key,
  title: column[locale.value] || column.vi,
})))

const localeText = computed(() => ({
  vi: { retry: 'Thử lại', scroll: 'Kéo ngang bảng để xem đầy đủ các cột' },
  en: { retry: 'Retry', scroll: 'Scroll horizontally to view all columns' },
  zh: { retry: '重试', scroll: '横向滚动以查看全部栏位' },
}[locale.value] || { retry: 'Thử lại', scroll: 'Kéo ngang bảng để xem đầy đủ các cột' }))

const retryLabel = computed(() => localeText.value.retry)
const scrollHint = computed(() => localeText.value.scroll)

const sumField = (field) => Track_data_list.value.reduce((total, row) => {
  const value = Number.parseFloat(row?.[field])
  return Number.isFinite(value) ? total + value : total
}, 0)

const totalQty = computed(() => sumField('Qty'))
const totalOutQty = computed(() => sumField('DQty'))
const totalRemainingQty = computed(() => sumField('RemQty'))

// --- Phân trang phía client ---
const PAGE_SIZE = 50
const currentPage = ref(1)

const totalPages = computed(() =>
  Math.max(1, Math.ceil((Track_data_list.value.length || 0) / PAGE_SIZE)),
)

const pagedList = computed(() => {
  const start = (currentPage.value - 1) * PAGE_SIZE
  return Track_data_list.value.slice(start, start + PAGE_SIZE)
})

const pageRangeLabel = computed(() => {
  const total = Track_data_list.value.length
  if (!total) return `0 / 0`
  const start = (currentPage.value - 1) * PAGE_SIZE + 1
  const end = Math.min(currentPage.value * PAGE_SIZE, total)
  return `${start} - ${end} / ${total}`
})

// Danh sách nút số trang hiển thị, có dấu "..." khi nhiều trang
const pageNumbers = computed(() => {
  const total = totalPages.value
  const current = currentPage.value
  const delta = 1
  const range = []
  const withDots = []
  let last

  for (let i = 1; i <= total; i++) {
    if (i === 1 || i === total || (i >= current - delta && i <= current + delta)) {
      range.push(i)
    }
  }

  for (const i of range) {
    if (last) {
      if (i - last === 2) withDots.push(last + 1)
      else if (i - last > 2) withDots.push('...')
    }
    withDots.push(i)
    last = i
  }

  return withDots
})

const goToPage = (page) => {
  const target = Math.min(Math.max(1, page), totalPages.value)
  currentPage.value = target
}

// Reset về trang 1 mỗi khi dữ liệu tracking thay đổi (search mới, refresh...)
watch(
  () => Track_data_list.value,
  () => {
    currentPage.value = 1
  },
)

const formatQuantity = (value) => new Intl.NumberFormat(locale.value, {
  maximumFractionDigits: 2,
}).format(value)

const displayValue = (value) => {
  if (value === null || value === undefined || value === '') return '—'
  return value
}

const rowKey = (row, index) => `${row?.Barcode || 'row'}-${row?.CLBH || ''}-${index}`

const openIMQCReport = (fileName) => {
  const url = `${API_PREFIX}/scan/entry/imqcreport?File_Name=${encodeURIComponent(fileName)}`
  window.open(url, '_blank', 'noopener,noreferrer')
}
</script>

<style scoped>
.tracking-results {
  display: flex;
  min-width: 0;
  min-height: 0;
  height: auto;
  /* Chiều cao tối đa của cả khối bảng: co theo số dòng dữ liệu, vượt mốc này mới cuộn dọc bên trong */
  max-height: 80vh;
  flex-direction: column;
  overflow: hidden;
  color: #172033;
  background: #fff;
}

.results-header {
  display: flex;
  min-height: 58px;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 9px 16px;
  border-bottom: 1px solid #dce4ef;
  background: #fff;
  flex-shrink: 0;
}

.results-eyebrow {
  margin: 0 0 1px;
  color: #58708e;
  font-size: 10px;
  font-weight: 800;
  letter-spacing: 0.14em;
  line-height: 1.2;
}

.results-title {
  margin: 0;
  color: #172947;
  font-size: 20px;
  font-weight: 800;
  line-height: 1.2;
}

.results-count {
  flex: 0 0 auto;
  padding: 6px 12px;
  border: 1px solid #c9d7e8;
  border-radius: 999px;
  color: #34516f;
  background: #f5f8fc;
  font-size: 13px;
  font-weight: 800;
}

.results-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 7px;
  padding: 7px 16px;
  border-bottom: 1px solid #dde5ef;
  background: #f7f9fc;
  flex-shrink: 0;
}

.summary-item {
  display: flex;
  min-width: 0;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 6px 10px;
  border: 1px solid #e0e7ef;
  border-left: 5px solid #7592b2;
  border-radius: 7px;
  color: #45566c;
  background: #fff;
  font-size: 12px;
  font-weight: 800;
}

.summary-item strong {
  color: #183d69;
  font-size: 19px;
  font-weight: 900;
}

.summary-item--in { border-color: #7fd9bb; border-left-color: #0d8f66; background: #d6f5e7; }
.summary-item--out { border-color: #f3c48f; border-left-color: #c9701f; background: #fbe3c4; }
.summary-item--stock { border-color: #a8c8ef; border-left-color: #1c5fb0; background: #d7e7fb; }
.summary-item--in span { color: #08653f; }
.summary-item--out span { color: #8a4a0d; }
.summary-item--stock span { color: #123f75; }
.summary-item--in strong { color: #055c39; }
.summary-item--out strong { color: #833f09; }
.summary-item--stock strong { color: #0e356b; }

.results-body {
  display: flex;
  min-width: 0;
  min-height: 0;
  flex: 1;
  /* Scroll thực sự diễn ra bên trong .table-shell (cả ngang lẫn dọc), ở đây chỉ giới hạn chiều cao */
  overflow: hidden;
}

.table-shell {
  min-width: 100%;
  width: 100%;
  /* table-shell tự quản lý CẢ 2 chiều scroll: nhờ vậy thanh cuộn ngang luôn nằm ở mép dưới
     vùng đang nhìn thấy (không phải mép dưới của toàn bộ dữ liệu) -> không cần kéo xuống mới thấy */
  overflow: auto;
  outline: none;
  scrollbar-color: #8ca3bd #e8eef5;
  scrollbar-width: thin;
  position: relative;
}

.table-shell:focus-visible {
  box-shadow: inset 0 0 0 3px rgba(36, 119, 207, 0.23);
}

.tracking-table {
  width: max-content;
  min-width: 2250px;
  border-spacing: 0;
  border-collapse: separate;
  color: #1b2740;
  background: #fff;
  font-size: 14px;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-rendering: optimizeLegibility;
}

.tracking-table th {
  position: sticky;
  z-index: 2;
  top: 0;
  min-width: 104px;
  height: 44px;
  padding: 0 12px;
  border-right: 1px solid #315d89;
  border-bottom: 1px solid #244f7b;
  color: #fff;
  background: #173f70;
  font-size: 13px;
  font-weight: 800;
  text-align: left;
  white-space: nowrap;
}

.tracking-table th:first-child,
.tracking-table td:first-child {
  position: sticky;
  z-index: 1;
  left: 0;
  min-width: 190px;
}

.tracking-table th:first-child {
  z-index: 3;
  background: #0f345f;
}

.tracking-table td {
  height: 40px;
  max-width: 320px;
  overflow: hidden;
  padding: 0 12px;
  border-right: 1px solid #e2e8f0;
  border-bottom: 1px solid #e4eaf1;
  background: #fff;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tracking-table tbody tr:nth-child(even) td {
  background: #f7f9fc;
}

.tracking-table tbody tr:hover td {
  background: #eaf3ff;
}

.tracking-table .state-row td {
  position: static;
  height: 210px;
  background: #f8fafc;
}

/* PHÂN TRANG: luôn hiển thị cố định bên dưới bảng, không nằm trong vùng cuộn nên không cần kéo xuống mới thấy */
.table-pagination {
  flex-shrink: 0;
  z-index: 4;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 16px;
  background: #f8fafc;
  border-top: 2px solid #cdd8e8;
  box-shadow: 0 -2px 6px rgba(24, 49, 84, 0.04);
}

.page-btn.ellipsis {
  border-color: transparent;
  background: transparent;
  cursor: default;
}

.page-info {
  font-size: 13px;
  color: #5b708b;
  font-weight: 600;
}

.page-controls {
  display: flex;
  gap: 4px;
}

.page-btn {
  min-width: 32px;
  height: 32px;
  padding: 0 8px;
  border: 1px solid #c9d7e8;
  border-radius: 6px;
  background: #fff;
  color: #34516f;
  font: inherit;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.page-btn:hover:not(:disabled) {
  background: #eaf3ff;
  border-color: #2471c2;
  color: #2471c2;
}

.page-btn.active {
  background: #173f70;
  border-color: #173f70;
  color: #fff;
  cursor: default;
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.barcode-value {
  color: #075eaa;
  font-family: Consolas, "Courier New", monospace;
  font-weight: 800;
}

.result-state {
  display: flex;
  width: 100%;
  min-height: 190px;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 7px;
  padding: 22px;
  color: #5b708b;
  text-align: center;
}

.result-state strong {
  color: #243d5d;
  font-size: 16px;
}

.result-state p {
  max-width: 520px;
  margin: 0;
  font-size: 13px;
}

.result-state--error strong,
.result-state--error .state-symbol {
  color: #b42318;
}

.state-symbol {
  display: grid;
  width: 34px;
  height: 34px;
  place-items: center;
  border: 2px solid currentColor;
  border-radius: 50%;
  color: #6a7f9a;
  font-size: 18px;
  font-weight: 800;
}

.empty-icon {
  width: 39px;
  height: 39px;
  fill: none;
  stroke: #6c83a0;
  stroke-linecap: round;
  stroke-width: 1.7;
}

.loading-spinner {
  width: 34px;
  height: 34px;
  border: 3px solid #d4e1ef;
  border-top-color: #1971c2;
  border-radius: 50%;
  animation: table-spin 0.75s linear infinite;
}

@keyframes table-spin {
  to { transform: rotate(360deg); }
}

.retry-button {
  margin-top: 3px;
  padding: 7px 14px;
  border: 1px solid #b8cce1;
  border-radius: 8px;
  color: #15578f;
  background: #edf5fc;
  font: inherit;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
}

.report-link {
  max-width: 230px;
  overflow: hidden;
  padding: 0;
  border: 0;
  color: #075ca8;
  background: transparent;
  font: inherit;
  font-weight: 750;
  text-decoration: underline;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
}

.report-link:hover {
  color: #003d78;
}

.mobile-results {
  display: none;
}

.scroll-hint {
  margin: 0;
  padding: 4px 16px;
  border-top: 1px solid #e2e8f0;
  color: #71839a;
  background: #f8fafc;
  font-size: 10px;
  text-align: right;
  flex-shrink: 0;
}

@media (max-width: 680px) {
  .results-header {
    padding: 9px 11px;
  }

  .results-summary {
    gap: 5px;
    padding: 6px 10px;
  }

  .summary-item {
    align-items: flex-start;
    flex-direction: column;
    gap: 0;
  }

  .summary-item strong {
    font-size: 14px;
  }

  .desktop-results,
  .scroll-hint {
    display: none;
  }

  .mobile-results {
    display: flex;
    flex-direction: column;
    width: 100%;
    overflow-y: auto;
    padding: 8px;
    background: #f1f5f9;
  }

  .result-card {
    overflow: hidden;
    margin-bottom: 8px;
    border: 1px solid #d8e2ee;
    border-radius: 11px;
    background: #fff;
    box-shadow: 0 4px 10px rgba(28, 49, 77, 0.05);
  }

  .result-card > header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    padding: 8px 10px;
    border-bottom: 1px solid #e3e9f1;
    background: #f7faff;
  }

  .result-card > header div {
    display: grid;
  }

  .result-card > header small {
    color: #6b7d93;
    font-size: 10px;
    font-weight: 700;
  }

  .result-card > header > span {
    padding: 3px 7px;
    border-radius: 999px;
    color: #155a75;
    background: #dff3f4;
    font-size: 12px;
    font-weight: 800;
  }

  .result-card__material {
    padding: 8px 10px 4px;
  }

  .result-card__material strong {
    color: #183f6a;
    font-size: 14px;
  }

  .result-card__material p {
    margin: 2px 0 0;
    color: #5b6e85;
    font-size: 12px;
  }

  .result-card dl {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 5px;
    margin: 0;
    padding: 6px 10px 9px;
  }

  .result-card dl div {
    min-width: 0;
    padding: 5px 6px;
    border-radius: 7px;
    background: #f3f6fa;
  }

  .result-card dt {
    overflow: hidden;
    color: #72849a;
    font-size: 9px;
    font-weight: 700;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .result-card dd {
    overflow: hidden;
    margin: 1px 0 0;
    color: #273d59;
    font-size: 13px;
    font-weight: 800;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .report-link--mobile {
    display: block;
    max-width: calc(100% - 20px);
    margin: 0 10px 9px;
    font-size: 12px;
  }

  .table-pagination--mobile {
    position: static;
    box-shadow: none;
    border-top: none;
    margin-top: auto;
    padding: 12px 10px;
    justify-content: center;
  }
}
</style>