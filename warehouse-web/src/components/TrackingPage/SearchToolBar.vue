<template>
  <section class="tracking-toolbar" :aria-label="copy.flow">
    <form class="tracking-form" :aria-busy="isloading" @submit.prevent="runTracking">
      <div class="tracking-form__grid">
        <!-- Thêm class dateIn và dateOut để CSS định vị -->
        <fieldset class="tracking-field tracking-field--date tracking-field--dateIn">
          <legend>{{ dateInLabel }}</legend>
          <DatePicker
            :index="0"
            :from-label="copy.from"
            :to-label="copy.to"
            @Set_Date_Range="setDateRange"
          />
        </fieldset>

        <fieldset class="tracking-field tracking-field--date tracking-field--dateOut">
          <legend>{{ dateOutLabel }}</legend>
          <DatePicker
            :index="1"
            :from-label="copy.from"
            :to-label="copy.to"
            @Set_Date_Range="setDateRange"
          />
        </fieldset>

        <label
          v-for="field in textFields"
          :key="field.index"
          class="tracking-field"
          :class="[
            `tracking-field--${field.key}`,
            {
              'tracking-field--conflict': hasConflictingSearchKeys && (field.index === 8 || field.index === 9),
            }
          ]"
        >
          <span>{{ field.label }}</span>
          <input
            :id="`tracking-field-${field.index}`"
            v-model="inputValues[field.index]"
            type="text"
            autocomplete="off"
            :placeholder="field.placeholder"
            @input="clearValidation"
          />
        </label>
      </div>

      <footer class="tracking-form__footer">
        <div class="form-message">
          <p v-if="validationMessage" class="tracking-validation" role="alert">
            <span aria-hidden="true">!</span>
            {{ validationMessage }}
          </p>
          <p v-else class="tracking-helper">{{ copy.helper }}</p>
        </div>

        <div class="tracking-form__buttons">
          <button class="action-button action-button--ghost" type="button" :disabled="isloading" @click="clearFilters">
            {{ copy.clear }}
          </button>
          <button class="action-button action-button--print" type="button" :disabled="isloading" @click="runPrint">
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path d="M7 9V3h10v6M7 18H5a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2M7 14h10v7H7z" />
            </svg>
            {{ printLabel }}
          </button>
          <button class="action-button action-button--primary" type="submit" :disabled="isloading">
            <span v-if="isloading" class="button-spinner" aria-hidden="true" />
            <svg v-else viewBox="0 0 24 24" aria-hidden="true">
              <circle cx="10.5" cy="10.5" r="6.5" />
              <path d="m15.5 15.5 5 5" />
            </svg>
            {{ isloading ? copy.searching : trackingLabel }}
          </button>
        </div>
      </footer>
    </form>
  </section>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import DatePicker from '@/components/DatePicker.vue'
import useTracking from '@/hooks/useTracking'

const { t, locale } = useI18n()
const { getTrackingData, getTrackDataToPrint, isloading } = useTracking()

const validationMessage = ref('')

const pad2 = (value) => String(value).padStart(2, '0')
const getLocalDate = () => {
  const date = new Date()
  return `${date.getFullYear()}-${pad2(date.getMonth() + 1)}-${pad2(date.getDate())}`
}

const today = getLocalDate()
const defaultDateRange = `${today}~${today}`
const inputValues = reactive([
  defaultDateRange,
  defaultDateRange,
  '',
  '',
  '',
  '',
  '',
  '',
  '',
  '',
])

const localeCopies = {
  vi: {
    from: 'Từ ngày',
    to: 'Đến ngày',
    helper: 'Chọn luồng nhập hoặc xuất, nhập điều kiện rồi bấm Tracking. Chỉ dùng một trong hai: MAKE hoặc Barcode.',
    flow: 'Luồng tracking',
    clear: 'Xóa lọc',
    searching: 'Đang tìm...',
  },
  en: {
    from: 'From date',
    to: 'To date',
    helper: 'Choose inbound or outbound, enter filters, then select Tracking. Use either MAKE or Barcode.',
    flow: 'Tracking flow',
    clear: 'Clear',
    searching: 'Searching...',
  },
  zh: {
    from: '开始日期',
    to: '结束日期',
    helper: '请选择入库或出库，输入条件后点击追踪。MAKE 与 Barcode 只能填写一个。',
    flow: '追踪流程',
    clear: '清除',
    searching: '查询中...',
  },
}

const copy = computed(() => localeCopies[locale.value] || localeCopies.vi)
const dateInLabel = computed(() => t('trackingPage.searchToolBar.scanDateIn'))
const dateOutLabel = computed(() => t('trackingPage.searchToolBar.scanDateOut'))

const effectiveModeIndex = computed(() => {
  const inChanged = inputValues[0] !== defaultDateRange
  const outChanged = inputValues[1] !== defaultDateRange
  return outChanged && !inChanged ? 1 : 0
})
const trackingLabel = computed(() => t('trackingPage.searchToolBar.tracking'))
const printLabel = computed(() => t('trackingPage.searchToolBar.print'))

const textFields = computed(() => [
  { index: 2, key: 'materialCode' },
  { index: 3, key: 'materialName' },
  { index: 4, key: 'supplierCode' },
  { index: 5, key: 'supplierName' },
  { index: 6, key: 'command' },
  { index: 7, key: 'purchaseOrder' },
  { index: 8, key: 'make' },
  { index: 9, key: 'barcode' },
].map((field) => ({
  ...field,
  label: t(`trackingPage.searchToolBar.${field.key}`),
  placeholder: t(`trackingPage.searchToolBar.${field.key}`),
})))

const hasConflictingSearchKeys = computed(() => Boolean(
  inputValues[8]?.trim() && inputValues[9]?.trim(),
))

const clearValidation = () => {
  if (!hasConflictingSearchKeys.value) validationMessage.value = ''
}

const validateLookup = () => {
  if (!hasConflictingSearchKeys.value) return true
  validationMessage.value = t('trackingPage.searchToolBar.alertMakeOrBarcode')
  return false
}

const setDateRange = (dateRange, index) => {
  inputValues[index] = dateRange
}

const runTracking = () => {
  if (!validateLookup()) return
  validationMessage.value = ''
  getTrackingData(inputValues, effectiveModeIndex.value)
}

const runPrint = () => {
  if (!validateLookup()) return
  validationMessage.value = ''
  getTrackDataToPrint(inputValues, effectiveModeIndex.value)
}

const clearFilters = () => {
  for (let index = 2; index <= 9; index += 1) inputValues[index] = ''
  validationMessage.value = ''
}
</script>

<style scoped>
.tracking-toolbar {
  min-width: 0;
  overflow: hidden;
  border: 1px solid #d7e1ee;
  border-radius: 16px;
  color: #172947;
  background: #fff;
  box-shadow: 0 10px 28px rgba(24, 49, 84, 0.08);
}

.tracking-form__footer,
.tracking-form__buttons,
.tracking-validation {
  display: flex;
  align-items: center;
}

.action-button svg {
  width: 18px;
  height: 18px;
  flex: 0 0 auto;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 1.8;
}

.tracking-form {
  padding: 11px 16px 12px;
}

/* =========================================
   GRID LAYOUT KHIẾN MÃ VT / TÊN VT THẲNG HÀNG
   ========================================= */
.tracking-form__grid {
  display: grid;
  gap: 10px;
  /* Mặc định Mobile: 1 cột, sắp xếp theo luồng nghiệp vụ */
  grid-template-columns: minmax(0, 1fr);
  grid-template-areas:
    "dateIn"
    "dateOut"
    "matCode"
    "matName"
    "supCode"
    "supName"
    "command"
    "po"
    "make"
    "barcode";
}

/* Gán vùng area cho từng ô field */
.tracking-field--dateIn { grid-area: dateIn; }
.tracking-field--dateOut { grid-area: dateOut; }
.tracking-field--materialCode { grid-area: matCode; }
.tracking-field--materialName { grid-area: matName; }
.tracking-field--supplierCode { grid-area: supCode; }
.tracking-field--supplierName { grid-area: supName; }
.tracking-field--command { grid-area: command; }
.tracking-field--purchaseOrder { grid-area: po; }
.tracking-field--make { grid-area: make; }
.tracking-field--barcode { grid-area: barcode; }

.tracking-field {
  display: grid;
  min-width: 0;
  gap: 4px;
  padding: 0;
  border: 0;
}

.tracking-field--date {
  padding: 7px 9px 8px;
  border: 1px solid #d6e0eb;
  border-radius: 10px;
  background: #f6f9fc;
}

.tracking-field > span,
.tracking-field legend {
  color: #33455c;
  font-size: 12px;
  font-weight: 800;
  line-height: 1.2;
}

.tracking-field legend {
  padding: 0 4px;
}

.tracking-field > input {
  width: 100%;
  min-width: 0;
  height: 42px;
  padding: 0 12px;
  border: 1px solid #cbd7e6;
  border-radius: 9px;
  outline: none;
  color: #172947;
  background: #fff;
  font: inherit;
  font-size: 14px;
  font-weight: 600;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.tracking-field > input::placeholder {
  color: #98a7b9;
}

.tracking-field > input:hover {
  border-color: #9fb3cc;
}

.tracking-field > input:focus-visible {
  border-color: #2471cf;
  box-shadow: 0 0 0 3px rgba(36, 119, 207, 0.15);
}

.tracking-field--conflict > input {
  border-color: #dc4c58;
  box-shadow: 0 0 0 3px rgba(220, 76, 88, 0.12);
}

/* FORM FOOTER */
.tracking-form__footer {
  display: flex;
  align-items: center;
  min-height: 46px;
  justify-content: space-between;
  gap: 14px;
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #e0e7f0;
}

.form-message {
  min-width: 0;
  flex: 1;
}

.tracking-helper,
.tracking-validation {
  margin: 0;
  color: #5d7088;
  font-size: 13px;
  line-height: 1.35;
}

.tracking-validation {
  gap: 7px;
  color: #b4232f;
  font-weight: 750;
}

.tracking-validation span {
  display: inline-grid;
  width: 18px;
  height: 18px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 50%;
  color: #fff;
  background: #c83f4b;
  font-size: 11px;
}

.tracking-form__buttons {
  flex: 0 0 auto;
  gap: 7px;
}

.action-button {
  display: inline-flex;
  min-height: 42px;
  align-items: center;
  justify-content: center;
  gap: 7px;
  padding: 0 15px;
  border: 1px solid transparent;
  border-radius: 9px;
  font: inherit;
  font-size: 14px;
  font-weight: 800;
  cursor: pointer;
  transition: transform 0.12s ease, background 0.15s ease;
}

.action-button:hover:not(:disabled) {
  transform: translateY(-1px);
}

.action-button:focus-visible {
  outline: 3px solid rgba(36, 119, 207, 0.22);
  outline-offset: 2px;
}

.action-button:disabled {
  cursor: wait;
  opacity: 0.66;
}

.action-button--ghost {
  border-color: #d2dce8;
  color: #52667f;
  background: #fff;
}

.action-button--print {
  border-color: #a9bfd8;
  color: #184c7b;
  background: #eaf3fb;
}

.action-button--primary {
  min-width: 116px;
  color: #fff;
  background: #176fc3;
  box-shadow: 0 6px 15px rgba(23, 111, 195, 0.22);
}

.action-button--primary:hover:not(:disabled) {
  background: #0f5da8;
}

.button-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: tracking-spin 0.7s linear infinite;
}

@keyframes tracking-spin {
  to { transform: rotate(360deg); }
}

/* =========================================
   RESPONSIVE GRID AREAS
   ========================================= */

/* Tablet (>= 480px): 2 cột. Mã VT & Mã NCC ngang hàng. Tên VT & Tên NCC ngang hàng */
@media (min-width: 480px) {
  .tracking-form__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    grid-template-areas:
      "dateIn   dateOut"
      "matCode  supCode"
      "matName  supName"
      "command  po"
      "make     barcode";
  }
}

/* Desktop (>= 768px): 4 cột */
@media (min-width: 768px) {
  .tracking-form__grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
    grid-template-areas:
      "dateIn   dateIn   dateOut  dateOut"
      "matCode  matCode  supCode  supCode"
      "matName  matName  supName  supName"
      "command  po       make     barcode";
  }
}

/* Large Desktop (>= 1280px): 6 cột, giãn rộng hàng dưới cùng */
@media (min-width: 1280px) {
  .tracking-form__grid {
    grid-template-columns: repeat(6, minmax(0, 1fr));
    grid-template-areas:
      "dateIn   dateIn   dateIn   dateOut  dateOut  dateOut"
      "matCode  matCode  matCode  supCode  supCode  supCode"
      "matName  matName  matName  supName  supName  supName"
      "command  po       make     make     barcode  barcode";
  }
}

/* Xử lý nút bấm ở Mobile nhỏ */
@media (max-width: 479px) {
  .tracking-form {
    padding-right: 10px;
    padding-left: 10px;
  }
  .tracking-form__footer {
    align-items: stretch;
    flex-direction: column;
  }
  .tracking-form__buttons {
    width: 100%;
    display: grid;
    grid-template-columns: 1fr 1fr;
  }
  .action-button {
    flex: 1;
  }
  .action-button--primary {
    grid-column: 1 / -1;
  }
}
</style>