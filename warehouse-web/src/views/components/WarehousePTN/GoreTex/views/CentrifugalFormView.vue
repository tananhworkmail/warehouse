<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import FormActions from '../components/FormActions.vue'
import NoticeToast from '../components/NoticeToast.vue'
import ReviewModeBar from '../components/ReviewModeBar.vue'
import EditHistoryDialog from '../components/EditHistoryDialog.vue'
import { getApiErrorMessage, getCentrifugalForm, submitCentrifugalForm } from '../utils/api'
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
  '07H30-08H00', '08H00-08H30', '08H30-09H00', '09H00-09H30',
  '09H30-10H00', '10H00-10H30', '10H30-11H00', '11H00-11H30',
  '11H30-12H00 / 12H30-13H00', '13H00-13H30', '13H30-14H00',
  '14H00-14H30', '14H30-15H00', '15H00-15H30', '15H30-16H00',
  '16H00-16H30', 'Tăng ca 16H30-17H / 加班 / Over time',
  'Tăng ca 17H-17H30 / 加班 / Over time', 'Tăng ca 17H30-18H / 加班 / Over time',
  'Tăng ca 18H-18H30 / 加班 / Over time', 'Tăng ca 18H30-19H / 加班 / Over time',
  'Tăng ca 19H-19H30 / 加班 / Over time', 'Tăng ca 19H30-20H / 加班 / Over time',
  'Tăng ca 20H-20H30 / 加班 / Over time',
]

const issues = [
  { id: 'toe', vi: 'Mũi thấm nước', zh: '鞋头漏水', en: 'Toe leakage' },
  { id: 'heel', vi: 'Gót thấm nước', zh: '后套漏水', en: 'Heel leakage' },
  { id: 'medial', vi: 'Hong trong thấm nước', zh: '内腰漏水', en: 'Medial quarter leakage' },
  { id: 'lateral', vi: 'Hong ngoài thấm nước', zh: '外腰漏水', en: 'Lateral quarter leakage' },
]

const resultOptions = ['PASS', 'FAIL']

function normalizeResult(value) {
  const normalized = String(value || '').trim().toUpperCase()
  return resultOptions.includes(normalized) ? normalized : ''
}

function newEntry(time, index) {
  return {
    id: `row-${index + 1}`,
    time,
    style: '',
    po: '',
    size: '',
    lean: '',
    issueValues: { toe: '0', heel: '0', medial: '0', lateral: '0' },
    otherIssue: '0',
    result: '',
    productionLeader: '',
    qcLeader: '',
  }
}

function emptyData() {
  return {
    inspectionDate: todayString(),
    entries: timeSlots.map(newEntry),
    signatures: { kcs: '', fieldOfficer: '', preparedBy: '' },
  }
}

function normalizedData(data = {}) {
  const entries = timeSlots.map((time, index) => {
    const old = data.entries?.find((entry) => entry.time === time) || data.entries?.[index] || {}
    const issueValues = { toe: '0', heel: '0', medial: '0', lateral: '0', ...(old.issueValues || {}) }
    if (!old.issueValues && Array.isArray(old.issues)) {
      old.issues.forEach((id) => { issueValues[id] = '✓' })
    }
    issues.forEach((issue) => {
      if (String(issueValues[issue.id] ?? '').trim() === '') issueValues[issue.id] = '0'
    })
    return {
      id: `row-${index + 1}`,
      time,
      style: old.style || '',
      po: old.po || '',
      size: old.size || '',
      lean: old.lean || '',
      issueValues,
      result: normalizeResult(old.result),
      otherIssue: String(old.otherIssue ?? '').trim() === '' ? '0' : String(old.otherIssue),
      productionLeader: old.productionLeader || '',
      qcLeader: old.qcLeader || '',
    }
  })
  return {
    inspectionDate: data.inspectionDate || todayString(),
    entries,
    signatures: {
      kcs: data.signatures?.kcs || '',
      fieldOfficer: data.signatures?.fieldOfficer || '',
      preparedBy: data.signatures?.preparedBy || '',
    },
  }
}

const loaded = loadDraft('centrifugal', emptyData())
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
  hasValue(form.inspectionDate) &&
  form.entries.every((entry) =>
    hasValue(entry.style) &&
    hasValue(entry.po) &&
    hasValue(entry.size) &&
    hasValue(entry.lean) &&
    issues.every((issue) => isNonNegativeNumber(entry.issueValues[issue.id])) &&
    isNonNegativeNumber(entry.otherIssue) &&
    hasValue(entry.result) &&
    hasValue(entry.productionLeader) &&
    hasValue(entry.qcLeader),
  ) &&
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
  savedAt.value = saveDraft('centrifugal', form)
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
    await submitCentrifugalForm(form, isReview.value)
    clearDraft('centrifugal')
    savedAt.value = null
    if (previousRevision) editHistory.value.push(previousRevision)
    reviewSnapshot.value = cloneFormData(form)
    isEditing.value = false
    viewingRevision.value = false
    await router.replace({
      name: 'laboratory-centrifugal-form',
      query: {
        review: '1',
        inspectionDate: form.inspectionDate,
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
  clearDraft('centrifugal')
  clearMissingFields(formElement.value)
  savedAt.value = null
  showNotice(t('actions.resetDone'), 'neutral')
}

onMounted(async () => {
  if (!isReview.value) return
  try {
    const detail = await getCentrifugalForm(route.query.inspectionDate)
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
        <div class="excel-sheet centrifugal-sheet">
          <h2>BÁO BIỂU GIÀY THÀNH PHẨM THỬ NGHIỆM LI TÂM</h2>
          <p class="sheet-subtitle">成品鞋离心测试报表 · DAILY QUALITY REPORT FOR WATERPROOF SHOES</p>

          <div class="sheet-meta one">
            <label><span>Ngày kiểm tra / 检查日 / Inspection date:  <input v-revision-change="revisionChange('inspectionDate')" v-model="form.inspectionDate" style ="max-width: 180px !important" type="date" required :readonly="isReview" /></span></label>
          </div>

          <table class="excel-table centrifugal-grid">
            <thead>
              <tr>
                <th rowspan="2">Thời gian<br /><small>时间 / Time</small></th>
                <th rowspan="2">Dạng giày<br /><small>型体 / Style name</small></th>
                <th rowspan="2">Lệnh<br /><small>订单 / PO #</small></th>
                <th rowspan="2">/ SIZE</th>
                <th rowspan="2">LEAN</th>
                <th colspan="5">Vấn đề / 问题 / Issue</th>
                <th rowspan="2">PASS / FAIL</th>
                <th rowspan="2">Cán bộ hiện trường kí tên<br /><small>现场干部签名 / Production leader confirmation</small></th>
                <th rowspan="2">Cán bộ QC kí tên<br /><small>品管干部签名 / QC leader confirmation</small></th>
              </tr>
              <tr>
                <th v-for="issue in issues" :key="issue.id">
                  {{ issue.vi }}<br /><small>{{ issue.zh }} / {{ issue.en }}</small>
                </th>
                <th>Vấn đề khác<br /><small>其他问题 / Other issue</small></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(entry, index) in form.entries" :key="entry.id">
                <td class="time-cell">{{ entry.time }}</td>
                <td><input v-revision-change="revisionChange(`entries.${index}.style`)" v-model.trim="entry.style" type="text" required :aria-label="`Dạng giày dòng ${index + 1}`" /></td>
                <td><input v-revision-change="revisionChange(`entries.${index}.po`)" v-model.trim="entry.po" type="text" required :aria-label="`Lệnh dòng ${index + 1}`" /></td>
                <td><input v-revision-change="revisionChange(`entries.${index}.size`)" v-model.trim="entry.size" type="text" required :aria-label="`Size dòng ${index + 1}`" /></td>
                <td><input v-revision-change="revisionChange(`entries.${index}.lean`)" v-model.trim="entry.lean" type="text" required :aria-label="`LEAN dòng ${index + 1}`" /></td>
                <td v-for="issue in issues" :key="issue.id">
                  <input v-zero-input v-revision-change="revisionChange(`entries.${index}.issueValues.${issue.id}`)" v-model="entry.issueValues[issue.id]" type="number" min="0" step="any" required :aria-label="`${issue.vi}, dòng ${index + 1}`" />
                </td>
                <td><input v-zero-input v-revision-change="revisionChange(`entries.${index}.otherIssue`)" v-model.trim="entry.otherIssue" type="number" min="0" step="any" required :aria-label="`Vấn đề khác dòng ${index + 1}`" /></td>
                <td>
                  <select
                    v-revision-change="revisionChange(`entries.${index}.result`)"
                    v-model="entry.result"
                    required
                    :aria-label="`PASS FAIL dòng ${index + 1}`"
                  >
                    <option value="" disabled>{{ t('actions.choose') }}</option>
                    <option v-for="result in resultOptions" :key="result" :value="result">
                      {{ result }}
                    </option>
                  </select>
                </td>
                <td><input v-revision-change="revisionChange(`entries.${index}.productionLeader`)" v-model.trim="entry.productionLeader" type="text" required :aria-label="`Production leader ${index + 1}`" /></td>
                <td><input v-revision-change="revisionChange(`entries.${index}.qcLeader`)" v-model.trim="entry.qcLeader" type="text" required :aria-label="`QC leader ${index + 1}`" /></td>
              </tr>
            </tbody>
          </table>

          <div class="signature-grid three">
            <label><span>CÁN BỘ KCS / 品管干部</span><input v-revision-change="revisionChange('signatures.kcs')" v-model.trim="form.signatures.kcs" type="text" required /></label>
            <label><span>CÁN BỘ HIỆN TRƯỜNG / 现场干部</span><input v-revision-change="revisionChange('signatures.fieldOfficer')" v-model.trim="form.signatures.fieldOfficer" type="text" required /></label>
            <label><span>LẬP BIỂU / 制表</span><input v-revision-change="revisionChange('signatures.preparedBy')" v-model.trim="form.signatures.preparedBy" type="text" required /></label>
          </div>
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
