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
import { clearMissingFields, isNonNegativeNumber, revealMissingFields } from '../utils/validation'
import { revisionChangeDirective } from '../utils/revisionHighlight'

const vRevisionChange = revisionChangeDirective
const { t, dateLocale } = useGoreTexI18n()

function dispatchZeroInput(element, value) {
  element.value = value
  element.dispatchEvent(new Event('input', { bubbles: true }))
}

const vZeroInput = {
  mounted(element) {
    element.type = 'number'
    element.min = '0'
    element.step = 'any'
    element.dataset.nonnegative = 'true'
    element.inputMode = 'decimal'
    const keydown = (event) => {
      if (event.key === '-') event.preventDefault()
    }
    const input = () => {
      const value = String(element.value).trim()
      if (value !== '' && Number(value) < 0) dispatchZeroInput(element, '0')
    }
    const focus = () => {
      if (String(element.value).trim() === '0') dispatchZeroInput(element, '')
    }
    const blur = () => {
      if (String(element.value).trim() === '') dispatchZeroInput(element, '0')
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
    timeSlots.forEach((slot) => { counts[`${row.id}:${slot.id}`] = '0' })
    totals[row.id] = '0'
    rates[row.id] = '0'
  })
  const summaryCounts = { inspection: {}, defects: {}, rates: {} }
  Object.values(summaryCounts).forEach((values) => {
    timeSlots.forEach((slot) => { values[slot.id] = '0' })
  })
  return {
    meta: {
      line: '', styleName: '', inspectionDate: todayString(),
    },
    counts,
    totals,
    rates,
    summaryCounts,
    summaryTotals: { inspection: '0', defects: '0', rate: '0' },
    summaryEdges: {
      inspection: { side: '0', end: '0' },
      defects: { side: '0', end: '0' },
      rates: { side: '0', total: '0' },
    },
    signatures: { productionSupervisor: '', kcs: '', fieldOfficer: '', preparedBy: '' },
  }
}

function numericMap(values = {}, defaults = {}) {
  return Object.fromEntries(Object.keys(defaults).map((key) => [
    key,
    String(values[key] ?? '').trim() === '' ? '0' : String(values[key]),
  ]))
}

function normalizedData(data = {}) {
  const base = emptyData()
  const counts = numericMap(data.counts, base.counts)
  const totals = numericMap(data.totals, base.totals)
  issueRows.forEach((row) => {
    if (String(data.totals?.[row.id] ?? '').trim() === '') {
      totals[row.id] = sumValues(timeSlots.map((slot) => counts[`${row.id}:${slot.id}`]))
    }
  })
  const summaryCounts = {
    inspection: numericMap(data.summaryCounts?.inspection, base.summaryCounts.inspection),
    defects: numericMap(data.summaryCounts?.defects, base.summaryCounts.defects),
    rates: numericMap(data.summaryCounts?.rates, base.summaryCounts.rates),
  }
  const summaryValue = (key, fallback) => {
    const value = key === 'inspection'
      ? data.summaryTotals?.inspection ?? data.inspectionTotal
      : data.summaryTotals?.[key]
    return String(value ?? '').trim() === '' ? fallback : String(value)
  }
  return {
    ...base,
    ...data,
    meta: {
      line: data.meta?.line || '',
      styleName: data.meta?.styleName || '',
      inspectionDate: data.meta?.inspectionDate || todayString(),
    },
    counts,
    totals,
    rates: numericMap(data.rates, base.rates),
    summaryCounts,
    summaryTotals: {
      inspection: summaryValue('inspection', sumValues(Object.values(summaryCounts.inspection))),
      defects: summaryValue('defects', sumValues(Object.values(summaryCounts.defects))),
      rate: summaryValue('rate', '0'),
    },
    summaryEdges: {
      inspection: {
        side: data.summaryEdges?.inspection?.side || '0',
        end: data.summaryEdges?.inspection?.end || '0',
      },
      defects: {
        side: data.summaryEdges?.defects?.side || '0',
        end: data.summaryEdges?.defects?.end || '0',
      },
      rates: {
        side: data.summaryEdges?.rates?.side || '0',
        total: data.summaryEdges?.rates?.total || sumValues(Object.values(summaryCounts.rates)),
      },
    },
    signatures: {
      productionSupervisor: data.signatures?.productionSupervisor || '',
      kcs: data.signatures?.kcs || '',
      fieldOfficer: data.signatures?.fieldOfficer || '',
      preparedBy: data.signatures?.preparedBy || '',
    },
  }
}

function numericValue(value) {
  const parsed = Number(String(value ?? '').replace(',', '.'))
  return Number.isFinite(parsed) ? parsed : 0
}

function sumValues(values) {
  const total = values.reduce((sum, value) => sum + numericValue(value), 0)
  return Number.isInteger(total) ? String(total) : String(Number(total.toFixed(2)))
}

function updateIssueTotal(rowId) {
  form.totals[rowId] = sumValues(timeSlots.map((slot) => form.counts[`${rowId}:${slot.id}`]))
}

function updateSummaryTotal(summaryKey) {
  const total = sumValues(timeSlots.map((slot) => form.summaryCounts[summaryKey][slot.id]))
  if (summaryKey === 'inspection') form.summaryTotals.inspection = total
  else if (summaryKey === 'defects') form.summaryTotals.defects = total
  else form.summaryEdges.rates.total = total
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
const isValid = computed(() =>
  hasValue(form.meta.line) &&
  hasValue(form.meta.styleName) &&
  hasValue(form.meta.inspectionDate) &&
  issueRows.every((row) =>
    timeSlots.every((slot) => isNonNegativeNumber(form.counts[`${row.id}:${slot.id}`])) &&
    isNonNegativeNumber(form.totals[row.id]) &&
    isNonNegativeNumber(form.rates[row.id]),
  ) &&
  ['inspection', 'defects', 'rates'].every((summaryKey) =>
    timeSlots.every((slot) => isNonNegativeNumber(form.summaryCounts[summaryKey][slot.id])),
  ) &&
  isNonNegativeNumber(form.summaryTotals.inspection) &&
  isNonNegativeNumber(form.summaryTotals.defects) &&
  isNonNegativeNumber(form.summaryTotals.rate) &&
  isNonNegativeNumber(form.summaryEdges.inspection.side) &&
  isNonNegativeNumber(form.summaryEdges.inspection.end) &&
  isNonNegativeNumber(form.summaryEdges.defects.side) &&
  isNonNegativeNumber(form.summaryEdges.defects.end) &&
  isNonNegativeNumber(form.summaryEdges.rates.side) &&
  isNonNegativeNumber(form.summaryEdges.rates.total) &&
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
            <label><span>Dạng giày / 型体 / Style name</span><input v-revision-change="revisionChange('meta.styleName')" v-model.trim="form.meta.styleName" type="text" required :readonly="isReview" /></label>
            <label><span>Ngày kiểm tra / 检查日 / Inspection date</span><input v-revision-change="revisionChange('meta.inspectionDate')" v-model="form.meta.inspectionDate" type="date" required /></label>
          </div>

          <table class="excel-table waterproof-grid">
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
                    step="any"
                    required
                    :aria-label="`${row.id}, ${slot.label}`"
                    @input="updateIssueTotal(row.id)"
                  />
                </td>
                <td><input v-zero-input v-revision-change="revisionChange(`totals.${row.id}`)" v-model="form.totals[row.id]" type="number" min="0" step="any" required :aria-label="`Tổng ${row.id}`" /></td>
                <td><input v-zero-input v-revision-change="revisionChange(`rates.${row.id}`)" v-model="form.rates[row.id]" type="number" min="0" step="any" required :aria-label="`Tỉ lệ ${row.id}`" /></td>
              </tr>

              <tr class="summary-row">
                <td colspan="3">Số lượng kiểm / 检验总数 / Total inspection</td>
                <td><input v-zero-input v-revision-change="revisionChange('summaryEdges.inspection.side')" v-model="form.summaryEdges.inspection.side" type="number" min="0" step="any" required aria-label="Số lượng kiểm, cột bên" /></td>
                <td v-for="slot in timeSlots" :key="slot.id"><input v-zero-input v-revision-change="revisionChange(`summaryCounts.inspection.${slot.id}`)" v-model="form.summaryCounts.inspection[slot.id]" type="number" min="0" step="any" required @input="updateSummaryTotal('inspection')" /></td>
                <td><input v-zero-input v-revision-change="revisionChange('summaryTotals.inspection')" v-model="form.summaryTotals.inspection" type="number" min="0" step="any" required /></td>
                <td><input v-zero-input v-revision-change="revisionChange('summaryEdges.inspection.end')" v-model="form.summaryEdges.inspection.end" type="number" min="0" step="any" required aria-label="Số lượng kiểm, cột cuối" /></td>
              </tr>
              <tr class="summary-row">
                <td colspan="3">Số lượng hư / 不良数 / Total defect</td>
                <td><input v-zero-input v-revision-change="revisionChange('summaryEdges.defects.side')" v-model="form.summaryEdges.defects.side" type="number" min="0" step="any" required aria-label="Số lượng hư, cột bên" /></td>
                <td v-for="slot in timeSlots" :key="slot.id"><input v-zero-input v-revision-change="revisionChange(`summaryCounts.defects.${slot.id}`)" v-model="form.summaryCounts.defects[slot.id]" type="number" min="0" step="any" required @input="updateSummaryTotal('defects')" /></td>
                <td><input v-zero-input v-revision-change="revisionChange('summaryTotals.defects')" v-model="form.summaryTotals.defects" type="number" min="0" step="any" required /></td>
                <td><input v-zero-input v-revision-change="revisionChange('summaryEdges.defects.end')" v-model="form.summaryEdges.defects.end" type="number" min="0" step="any" required aria-label="Số lượng hư, cột cuối" /></td>
              </tr>
              <tr class="summary-row">
                <td colspan="3">Tỉ lệ hàng hư / 不良百分比 / Defect rate</td>
                <td><input v-zero-input v-revision-change="revisionChange('summaryEdges.rates.side')" v-model="form.summaryEdges.rates.side" type="number" min="0" step="any" required aria-label="Tỉ lệ hàng hư, cột bên" /></td>
                <td v-for="slot in timeSlots" :key="slot.id"><input v-zero-input v-revision-change="revisionChange(`summaryCounts.rates.${slot.id}`)" v-model="form.summaryCounts.rates[slot.id]" type="number" min="0" step="any" required @input="updateSummaryTotal('rates')" /></td>
                <td><input v-zero-input v-revision-change="revisionChange('summaryEdges.rates.total')" v-model="form.summaryEdges.rates.total" type="number" min="0" step="any" required aria-label="Tỉ lệ hàng hư, tổng cộng" /></td>
                <td><input v-zero-input v-revision-change="revisionChange('summaryTotals.rate')" v-model="form.summaryTotals.rate" type="number" min="0" step="any" required /></td>
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
