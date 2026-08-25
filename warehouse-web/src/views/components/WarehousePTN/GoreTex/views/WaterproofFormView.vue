<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import FormActions from '../components/FormActions.vue'
import NoticeToast from '../components/NoticeToast.vue'
import ReviewModeBar from '../components/ReviewModeBar.vue'
import EditHistoryDialog from '../components/EditHistoryDialog.vue'
import { getApiErrorMessage, getWaterproofForm, submitWaterproofForm } from '../utils/api'
import { cloneFormData, createLocalRevision, getCurrentFormData, getEditHistory } from '../utils/revisions'
import { clearDraft, loadDraft, saveDraft, todayString } from '../utils/storage'
import { useGoreTexI18n } from '../utils/i18n'
import { clearMissingFields, isNonNegativeInteger, revealMissingFields } from '../utils/validation'
import { revisionChangeDirective } from '../utils/revisionHighlight'
import { cellGuideDirective } from '../utils/cellGuide'

const vRevisionChange = revisionChangeDirective
const vCellGuide = cellGuideDirective
const { t, dateLocale } = useGoreTexI18n()

function dispatchZeroInput(element, value) {
  element.value = value
  element.dispatchEvent(new Event('input', { bubbles: true }))
}

const vZeroInput = {
  mounted(element) {
    element.type = 'number'
    element.min = '0'
    element.step = '1'
    element.dataset.nonnegative = 'true'
    element.dataset.integer = 'true'
    element.inputMode = 'numeric'
    const keydown = (event) => {
      if (['-', '+', '.', ',', 'e', 'E'].includes(event.key)) event.preventDefault()
    }
    const input = () => {
      const value = String(element.value).trim()
      if (value !== '' && Number(value) < 0) dispatchZeroInput(element, '0')
    }
    const focus = () => {
      if (String(element.value).trim() === '0') dispatchZeroInput(element, '')
    }
    const blur = () => {
      if (String(element.value).trim() === '' && element.dataset.optional !== 'true') dispatchZeroInput(element, '0')
    }
    element._goreTexZeroInput = { keydown, input, focus, blur }
    element.addEventListener('keydown', keydown)
    element.addEventListener('input', input)
    element.addEventListener('focus', focus)
    element.addEventListener('blur', blur)
  },
  beforeUnmount(element) {
    const handlers = element._goreTexZeroInput
    if (!handlers) return
    element.removeEventListener('keydown', handlers.keydown)
    element.removeEventListener('input', handlers.input)
    element.removeEventListener('focus', handlers.focus)
    element.removeEventListener('blur', handlers.blur)
    delete element._goreTexZeroInput
  },
}

const timeSlots = [
  { id: '0830', label: '08:30' }, { id: '0930', label: '09:30' },
  { id: '1030', label: '10:30' }, { id: '1130', label: '11:30' },
  { id: '1330', label: '13:30' }, { id: '1430', label: '14:30' },
  { id: '1530', label: '15:30' }, { id: '1630', label: '16:30' },
  { id: 'overtime', label: 'Tăng ca / 加班 / Over time' },
]

const issueRows = [
  { id: 'toe-left', no: '1', noRowspan: 2, issueRowspan: 2, issueColspan: 2, vi: 'Mũi thấm nước', zh: '鞋头漏水', en: 'Toe leakage', side: 'Trái / 左', groupStart: true },
  { id: 'toe-right', side: 'Phải / 右' },
  { id: 'heel-left', no: '2', noRowspan: 2, issueRowspan: 2, issueColspan: 2, vi: 'Gót thấm nước', zh: '后套漏水', en: 'Heel leakage', side: 'Trái / 左', groupStart: true },
  { id: 'heel-right', side: 'Phải / 右' },
  { id: 'medial-left', no: '3', noRowspan: 2, issueRowspan: 2, issueColspan: 2, vi: 'Hong trong thấm nước', zh: '内腰漏水', en: 'Medial quarter leakage', side: 'Trái / 左', groupStart: true },
  { id: 'medial-right', side: 'Phải / 右' },
  { id: 'lateral-left', no: '4', noRowspan: 2, issueRowspan: 2, issueColspan: 2, vi: 'Hong ngoài thấm nước', zh: '外腰漏水', en: 'Lateral quarter leakage', side: 'Trái / 左', groupStart: true },
  { id: 'lateral-right', side: 'Phải / 右' },
  { id: 'material-left', no: '5', noRowspan: 2, issueRowspan: 2, issueColspan: 2, vi: 'Vật tư không đạt', zh: '材料不良', en: 'Poor material', side: 'Trái / 左', groupStart: true },
  { id: 'material-right', side: 'Phải / 右' },
  { id: 'attaching-left', no: '6', noRowspan: 2, issueRowspan: 2, issueColspan: 2, vi: 'Dán đế lệch', zh: '贴移位', en: 'Attaching did not follow marks', side: 'Trái / 左', groupStart: true },
  { id: 'attaching-right', side: 'Phải / 右' },
  { id: 'wrinkled-left', no: '7', noRowspan: 2, issueRowspan: 2, issueColspan: 2, vi: 'Ép đế nhăn', zh: '底皱', en: 'Wrinkled sole', side: 'Trái / 左', groupStart: true },
  { id: 'wrinkled-right', side: 'Phải / 右' },
  { id: 'zigzag-quarter-left', no: '8', noRowspan: 4, issueRowspan: 4, issueColspan: 1, vi: 'Zíc zắc hở', zh: '万能车线松开', en: 'Open zigzag stitch', subVi: 'Thân', subZh: '腰身', subEn: 'Quarter', subRowspan: 2, side: 'Trái / 左', groupStart: true },
  { id: 'zigzag-quarter-right', side: 'Phải / 右' },
  { id: 'zigzag-sole-left', subVi: 'Đế', subZh: '底', subEn: 'Sole', subRowspan: 2, side: 'Trái / 左' },
  { id: 'zigzag-sole-right', side: 'Phải / 右' },
  { id: 'bonding-toe-left', no: '9', noRowspan: 4, issueRowspan: 4, issueColspan: 1, vi: 'Hở keo đế', zh: '底开胶', en: 'Poor bonding sole', subVi: 'Mũi', subZh: '鞋头', subEn: 'Toe / vamp', subRowspan: 2, side: 'Trái / 左', groupStart: true },
  { id: 'bonding-toe-right', side: 'Phải / 右' },
  { id: 'bonding-heel-left', subVi: 'Gót', subZh: '后套', subEn: 'Heel', subRowspan: 2, side: 'Trái / 左' },
  { id: 'bonding-heel-right', side: 'Phải / 右' },
]

function emptyData() {
  const counts = {}
  const totals = {}
  const rates = {}
  issueRows.forEach((row) => {
    timeSlots.forEach((slot) => { counts[`${row.id}:${slot.id}`] = slot.id === 'overtime' ? '' : '0' })
    totals[row.id] = '0'
    rates[row.id] = '0'
  })
  const summaryCounts = { inspection: {}, defects: {}, rates: {}, passed: {}, passRates: {} }
  Object.entries(summaryCounts).forEach(([key, values]) => {
    timeSlots.forEach((slot) => { values[slot.id] = key === 'inspection' && slot.id === 'overtime' ? '' : '0' })
  })
  return {
    meta: {
      line: '', styleName: '', inspectionDate: todayString(),
    },
    counts,
    totals,
    rates,
    summaryCounts,
    summaryTotals: { inspection: '0', defects: '0', rate: '0', passed: '0', passRate: '0' },
    signatures: { productionSupervisor: '', kcs: '', fieldOfficer: '', preparedBy: '' },
  }
}

function numericMap(values = {}, defaults = {}, optionalKeys = []) {
  return Object.fromEntries(Object.keys(defaults).map((key) => [
    key,
    String(values[key] ?? '').trim() === '' ? (optionalKeys.includes(key) ? '' : '0') : String(values[key]),
  ]))
}

function formatNumber(value) {
  if (!Number.isFinite(value)) return '0'
  return Number.isInteger(value) ? String(value) : String(Number(value.toFixed(2)))
}

function calculateRate(numerator, denominator) {
  return denominator > 0 ? formatNumber((numerator / denominator) * 100) : '0'
}

function recalculateData(target) {
  const inspectionTotal = timeSlots.reduce((sum, slot) => sum + numericValue(target.summaryCounts.inspection[slot.id]), 0)
  const defectsBySlot = Object.fromEntries(timeSlots.map((slot) => [slot.id, 0]))

  issueRows.forEach((row) => {
    const total = timeSlots.reduce((sum, slot) => {
      const value = numericValue(target.counts[`${row.id}:${slot.id}`])
      defectsBySlot[slot.id] += value
      return sum + value
    }, 0)
    target.totals[row.id] = formatNumber(total)
    target.rates[row.id] = calculateRate(total, inspectionTotal)
  })

  let totalDefects = 0
  let totalPassed = 0
  timeSlots.forEach((slot) => {
    const inspection = numericValue(target.summaryCounts.inspection[slot.id])
    const defects = defectsBySlot[slot.id]
    const passed = inspection - defects
    target.summaryCounts.defects[slot.id] = formatNumber(defects)
    target.summaryCounts.rates[slot.id] = calculateRate(defects, inspection)
    target.summaryCounts.passed[slot.id] = formatNumber(passed)
    target.summaryCounts.passRates[slot.id] = calculateRate(passed, inspection)
    totalDefects += defects
    totalPassed += passed
  })
  target.summaryTotals.inspection = formatNumber(inspectionTotal)
  target.summaryTotals.defects = formatNumber(totalDefects)
  target.summaryTotals.rate = calculateRate(totalDefects, inspectionTotal)
  target.summaryTotals.passed = formatNumber(totalPassed)
  target.summaryTotals.passRate = calculateRate(totalPassed, inspectionTotal)
  return target
}

function normalizedData(data = {}) {
  const base = emptyData()
  const optionalCountKeys = issueRows.map((row) => `${row.id}:overtime`)
  const counts = numericMap(data.counts, base.counts, optionalCountKeys)
  const summaryCounts = {
    inspection: numericMap(data.summaryCounts?.inspection, base.summaryCounts.inspection, ['overtime']),
    defects: numericMap(data.summaryCounts?.defects, base.summaryCounts.defects),
    rates: numericMap(data.summaryCounts?.rates, base.summaryCounts.rates),
    passed: numericMap(data.summaryCounts?.passed, base.summaryCounts.passed),
    passRates: numericMap(data.summaryCounts?.passRates, base.summaryCounts.passRates),
  }
  const result = {
    ...base,
    ...data,
    meta: {
      line: data.meta?.line || '',
      styleName: String(data.meta?.styleName || '').toLocaleUpperCase(),
      inspectionDate: data.meta?.inspectionDate || todayString(),
    },
    counts,
    totals: { ...base.totals },
    rates: { ...base.rates },
    summaryCounts,
    summaryTotals: { ...base.summaryTotals },
    signatures: {
      productionSupervisor: data.signatures?.productionSupervisor || '',
      kcs: data.signatures?.kcs || '',
      fieldOfficer: data.signatures?.fieldOfficer || '',
      preparedBy: data.signatures?.preparedBy || '',
    },
  }
  return recalculateData(result)
}

function numericValue(value) {
  const parsed = Number(String(value ?? '').replace(',', '.'))
  return Number.isFinite(parsed) ? parsed : 0
}

function recalculate() {
  recalculateData(form)
}

const loaded = loadDraft('waterproof', emptyData())
const form = reactive(normalizedData(loaded.data))
const savedAt = ref(loaded.savedAt)
const notice = reactive({ show: false, message: '', tone: 'success', timer: null })
const submitting = ref(false)
const formElement = ref(null)
const route = useRoute()
const router = useRouter()
const isReview = computed(() => route.query.review === '1')
const isEditing = ref(false)
const reviewSnapshot = ref(null)
const editHistory = ref([])
const showEditHistory = ref(false)
const viewingRevision = ref(false)
const revisionChanges = ref([])
const revisionChange = (path) => revisionChanges.value.find((change) => change.path === path)
const isReadOnly = computed(() => isReview.value && !isEditing.value)
const hasValue = (value) => String(value ?? '').trim() !== ''
const hasNegativeCalculatedValue = computed(() =>
  timeSlots.some((slot) => numericValue(form.summaryCounts.passed[slot.id]) < 0) ||
  numericValue(form.summaryTotals.passed) < 0,
)
const isValid = computed(() =>
  hasValue(form.meta.line) &&
  hasValue(form.meta.styleName) &&
  hasValue(form.meta.inspectionDate) &&
  issueRows.every((row) =>
    timeSlots.every((slot) => slot.id === 'overtime'
      ? !hasValue(form.counts[`${row.id}:${slot.id}`]) || isNonNegativeInteger(form.counts[`${row.id}:${slot.id}`])
      : isNonNegativeInteger(form.counts[`${row.id}:${slot.id}`])),
  ) &&
  timeSlots.every((slot) => slot.id === 'overtime'
    ? !hasValue(form.summaryCounts.inspection[slot.id]) || isNonNegativeInteger(form.summaryCounts.inspection[slot.id])
    : isNonNegativeInteger(form.summaryCounts.inspection[slot.id])) &&
  !hasNegativeCalculatedValue.value &&
  Object.values(form.signatures).every(hasValue),
)
const savedLabel = computed(() => savedAt.value
  ? t('actions.savedAt', { time: new Intl.DateTimeFormat(dateLocale(), { hour: '2-digit', minute: '2-digit', day: '2-digit', month: '2-digit' }).format(new Date(savedAt.value)) })
  : t('actions.noDraft'))

function showNotice(message, tone = 'success') {
  clearTimeout(notice.timer)
  Object.assign(notice, { show: true, message, tone })
  notice.timer = setTimeout(() => { notice.show = false }, 3000)
}

function goBack() {
  if (window.history.state?.back) router.back()
  else router.push({ name: 'laboratory-forms-portal' })
}

function save() {
  savedAt.value = saveDraft('waterproof', form)
  showNotice(t('actions.savedDraft'))
}

function enableEditing() {
  if (viewingRevision.value) return
  isEditing.value = true
}

function cancelEditing() {
  if (reviewSnapshot.value) {
    Object.assign(form, normalizedData(reviewSnapshot.value))
  }
  isEditing.value = false
}

function viewRevision(entry) {
  Object.assign(form, normalizedData(entry.viewData || entry.data))
  revisionChanges.value = entry.changes || []
  viewingRevision.value = true
  isEditing.value = false
  showEditHistory.value = false
}

function showCurrentVersion() {
  if (reviewSnapshot.value) {
    Object.assign(form, normalizedData(reviewSnapshot.value))
  }
  viewingRevision.value = false
  revisionChanges.value = []
}

async function submit() {
  if (isReadOnly.value || submitting.value) return
  if (hasNegativeCalculatedValue.value) {
    formElement.value?.querySelector('.formula-negative')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    showNotice(t('actions.negativeCalculated'), 'neutral')
    return
  }
  if (!isValid.value) {
    const missingCount = revealMissingFields(formElement.value)
    showNotice(t('actions.missing', { count: missingCount }), 'neutral')
    return
  }
  clearMissingFields(formElement.value)
  submitting.value = true
  try {
    const previousRevision = isReview.value && reviewSnapshot.value
      ? createLocalRevision(reviewSnapshot.value)
      : null
    await submitWaterproofForm(form, isReview.value)
    clearDraft('waterproof')
    savedAt.value = null
    if (previousRevision) editHistory.value.push(previousRevision)
    reviewSnapshot.value = cloneFormData(form)
    isEditing.value = false
    viewingRevision.value = false
    await router.replace({
      name: 'laboratory-waterproof-form',
      query: {
        review: '1',
        line: form.meta.line,
        styleName: form.meta.styleName,
      },
    })
  } catch (error) {
    showNotice(getApiErrorMessage(error), 'neutral')
  } finally {
    submitting.value = false
  }
}

function reset() {
  if (!window.confirm(t('actions.confirmReset'))) return
  Object.assign(form, emptyData())
  clearDraft('waterproof')
  clearMissingFields(formElement.value)
  savedAt.value = null
  showNotice(t('actions.resetDone'), 'neutral')
}

onMounted(async () => {
  if (!isReview.value) return
  try {
    const detail = await getWaterproofForm(route.query.line, route.query.styleName)
    editHistory.value = getEditHistory(detail.data)
    reviewSnapshot.value = getCurrentFormData(detail.data)
    Object.assign(form, normalizedData(reviewSnapshot.value))
  } catch (error) {
    showNotice(getApiErrorMessage(error), 'neutral')
  }
})
</script>

<template>
  <div class="form-page excel-page">
    <ReviewModeBar
      v-if="isReview"
      :editing="isEditing"
      :valid="isValid"
      :submitting="submitting"
      :history-count="editHistory.length"
      :viewing-revision="viewingRevision"
      @edit="enableEditing"
      @cancel="cancelEditing"
      @submit="submit"
      @history="showEditHistory = true"
      @current="showCurrentVersion"
      @back="goBack"
    />
    <form ref="formElement" class="excel-form" @submit.prevent="submit">
      <fieldset class="review-fieldset" :disabled="isReadOnly">
        <div class="sheet-scroll">
        <div class="excel-sheet waterproof-sheet">
          <h2>BÁO BIỂU KIỂM TRA CHẤT LƯỢNG GIÀY CHỐNG THẤM NƯỚC</h2>
          <p class="sheet-subtitle">防水鞋品质检验报表 · DAILY QUALITY REPORT FOR WATERPROOF SHOES</p>

          <div class="sheet-meta three">
            <label><span>Chuyền / 线别 / Line</span><input v-revision-change="revisionChange('meta.line')" v-model.trim="form.meta.line" type="text" required :readonly="isReview" /></label>
            <label><span>Dạng giày / 型体 / Style name</span><input v-revision-change="revisionChange('meta.styleName')" v-model.trim="form.meta.styleName" type="text" required :readonly="isReview" @input="form.meta.styleName = form.meta.styleName.toLocaleUpperCase()" /></label>
            <label><span>Ngày kiểm tra / 检查日 / Inspection date</span><input v-revision-change="revisionChange('meta.inspectionDate')" v-model="form.meta.inspectionDate" type="date" required /></label>
          </div>

          <table v-cell-guide class="excel-table waterproof-grid">
            <colgroup>
              <col class="col-no" />
              <col class="col-issue-main" />
              <col class="col-sub-issue" />
              <col class="col-side" />
              <col
                v-for="slot in timeSlots"
                :key="`col-${slot.id}`"
                :class="slot.id === 'overtime' ? 'col-overtime' : 'col-time'"
              />
              <col class="col-total" />
              <col class="col-rate" />
            </colgroup>
            <thead>
              <tr>
                <th class="col-no">STT<br /><small>序号 / NO</small></th>
                <th class="col-issue" colspan="3">Vấn đề<br /><small>问题 / Issue</small></th>
                <th v-for="slot in timeSlots" :key="slot.id">{{ slot.label }}</th>
                <th>Tổng cộng<br /><small>总共 / Total</small></th>
                <th>Tỉ lệ hàng hư<br /><small>不良百分比</small></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in issueRows" :key="row.id" :class="{ 'group-start': row.groupStart }">
                <td v-if="row.noRowspan" :rowspan="row.noRowspan" class="col-no">{{ row.no }}</td>
                <td
                  v-if="row.issueRowspan"
                  :rowspan="row.issueRowspan"
                  :colspan="row.issueColspan"
                  class="issue-cell"
                >
                  <strong>{{ row.vi }}</strong>
                  <small>{{ row.zh }} / {{ row.en }}</small>
                </td>
                <td v-if="row.subRowspan" :rowspan="row.subRowspan" class="sub-issue-cell">
                  <strong>{{ row.subVi }}</strong>
                  <small>{{ row.subZh }} / {{ row.subEn }}</small>
                </td>
                <td class="side-cell">{{ row.side }}</td>
                <td v-for="slot in timeSlots" :key="slot.id">
                  <input
                    v-zero-input
                    v-revision-change="revisionChange(`counts.${row.id}:${slot.id}`)"
                    v-model="form.counts[`${row.id}:${slot.id}`]"
                    type="number"
                    min="0"
                    step="1"
                    :required="slot.id !== 'overtime'"
                    :data-optional="slot.id === 'overtime' ? 'true' : null"
                    :aria-label="`${row.id}, ${slot.label}`"
                    @input="recalculate"
                  />
                </td>
                <td class="formula-cell"><span v-revision-change="revisionChange(`totals.${row.id}`)" class="calculated-value">{{ form.totals[row.id] }}</span></td>
                <td class="formula-cell"><span v-revision-change="revisionChange(`rates.${row.id}`)" class="calculated-value">{{ form.rates[row.id] }}%</span></td>
              </tr>

              <tr class="summary-row">
                <td colspan="3">Số lượng kiểm / 检验总数 / Total inspection</td>
                <td class="locked-cell"></td>
                <td v-for="slot in timeSlots" :key="slot.id"><input v-zero-input v-revision-change="revisionChange(`summaryCounts.inspection.${slot.id}`)" v-model="form.summaryCounts.inspection[slot.id]" type="number" min="0" step="1" :required="slot.id !== 'overtime'" :data-optional="slot.id === 'overtime' ? 'true' : null" @input="recalculate" /></td>
                <td class="formula-cell"><span v-revision-change="revisionChange('summaryTotals.inspection')" class="calculated-value">{{ form.summaryTotals.inspection }}</span></td>
                <td class="locked-cell"></td>
              </tr>
              <tr class="summary-row">
                <td colspan="3">Số lượng hư / 不良数 / Total defect</td>
                <td class="locked-cell"></td>
                <td v-for="slot in timeSlots" :key="slot.id" class="formula-cell"><span v-revision-change="revisionChange(`summaryCounts.defects.${slot.id}`)" class="calculated-value">{{ form.summaryCounts.defects[slot.id] }}</span></td>
                <td class="formula-cell"><span v-revision-change="revisionChange('summaryTotals.defects')" class="calculated-value">{{ form.summaryTotals.defects }}</span></td>
                <td class="locked-cell"></td>
              </tr>
              <tr class="summary-row">
                <td colspan="3">Tỉ lệ hàng hư / 不良百分比 / Defect rate</td>
                <td class="locked-cell"></td>
                <td v-for="slot in timeSlots" :key="slot.id" class="formula-cell"><span v-revision-change="revisionChange(`summaryCounts.rates.${slot.id}`)" class="calculated-value">{{ form.summaryCounts.rates[slot.id] }}%</span></td>
                <td class="formula-cell"><span v-revision-change="revisionChange('summaryTotals.rate')" class="calculated-value">{{ form.summaryTotals.rate }}%</span></td>
                <td class="locked-cell"></td>
              </tr>
              <tr class="summary-row">
                <td colspan="3">Số lần thông qua / 通过数量 / Pass quantity</td>
                <td class="locked-cell"></td>
                <td v-for="slot in timeSlots" :key="slot.id" class="formula-cell" :class="{ 'formula-negative': numericValue(form.summaryCounts.passed[slot.id]) < 0 }"><span v-revision-change="revisionChange(`summaryCounts.passed.${slot.id}`)" class="calculated-value">{{ form.summaryCounts.passed[slot.id] }}</span></td>
                <td class="formula-cell" :class="{ 'formula-negative': numericValue(form.summaryTotals.passed) < 0 }"><span v-revision-change="revisionChange('summaryTotals.passed')" class="calculated-value">{{ form.summaryTotals.passed }}</span></td>
                <td class="locked-cell"></td>
              </tr>
              <tr class="summary-row">
                <td colspan="3">Tỷ lệ thông qua / 通过率 / Pass rate</td>
                <td class="locked-cell"></td>
                <td v-for="slot in timeSlots" :key="slot.id" class="formula-cell" :class="{ 'formula-negative': numericValue(form.summaryCounts.passed[slot.id]) < 0 }"><span v-revision-change="revisionChange(`summaryCounts.passRates.${slot.id}`)" class="calculated-value">{{ form.summaryCounts.passRates[slot.id] }}%</span></td>
                <td class="formula-cell" :class="{ 'formula-negative': numericValue(form.summaryTotals.passed) < 0 }"><span v-revision-change="revisionChange('summaryTotals.passRate')" class="calculated-value">{{ form.summaryTotals.passRate }}%</span></td>
                <td class="locked-cell"></td>
              </tr>
              <tr class="confirmation-row">
                <td colspan="3">Cán bộ sản xuất xác nhận<br /><small>生产干部确认 / Production Supervisor confirmation</small></td>
                <td colspan="12"><input v-revision-change="revisionChange('signatures.productionSupervisor')" v-model.trim="form.signatures.productionSupervisor" type="text" required /></td>
              </tr>
              <tr class="signature-row">
                <td colspan="4" class="signature-name-cell">CÁN BỘ KCS<br /><small>品管干部</small><input v-revision-change="revisionChange('signatures.kcs')" v-model.trim="form.signatures.kcs" type="text" required /></td>
                <td colspan="6" class="signature-name-cell">CÁN BỘ HIỆN TRƯỜNG<br /><small>现场干部</small><input v-revision-change="revisionChange('signatures.fieldOfficer')" v-model.trim="form.signatures.fieldOfficer" type="text" required /></td>
                <td colspan="5" class="signature-name-cell">LẬP BIỂU<br /><small>制表</small><input v-revision-change="revisionChange('signatures.preparedBy')" v-model.trim="form.signatures.preparedBy" type="text" required /></td>
              </tr>
            </tbody>
          </table>
        </div>
        </div>
      </fieldset>

      <FormActions
        v-if="!isReview"
        :saved-label="savedLabel"
        :valid="isValid"
        :submitting="submitting"
        @save="save"
        @submit="submit"
        @reset="reset"
        @cancel="goBack"
      />
    </form>
    <EditHistoryDialog
      :show="showEditHistory"
      :history="editHistory"
      :current-data="reviewSnapshot || {}"
      @close="showEditHistory = false"
      @view="viewRevision"
    />
    <NoticeToast :show="notice.show" :message="notice.message" :tone="notice.tone" />
  </div>
</template>
