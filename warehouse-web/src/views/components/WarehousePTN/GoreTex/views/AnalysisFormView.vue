<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import FormActions from '../components/FormActions.vue'
import NoticeToast from '../components/NoticeToast.vue'
import ReviewModeBar from '../components/ReviewModeBar.vue'
import EditHistoryDialog from '../components/EditHistoryDialog.vue'
import { getAnalysisForm, getApiErrorMessage, submitAnalysisForm } from '../utils/api'
import { cloneFormData, createLocalRevision, getCurrentFormData, getEditHistory } from '../utils/revisions'
import { clearDraft, loadDraft, saveDraft } from '../utils/storage'
import { useGoreTexI18n } from '../utils/i18n'
import { clearMissingFields, revealMissingFields } from '../utils/validation'
import { revisionChangeDirective } from '../utils/revisionHighlight'

const vRevisionChange = revisionChangeDirective
const { t, dateLocale } = useGoreTexI18n()

let recordSequence = 1

function newRecord() {
  return {
    id: `record-${Date.now()}-${recordSequence++}`,
    testDate: '',
    styleName: '',
    testedCount: '',
    failedCount: '',
    rate: '',
    leakagePosition: '',
    cause: '',
    action: '',
    unit: '',
    improvementDate: '',
    qc: '',
  }
}

function emptyData() {
  return { records: [newRecord()] }
}

function normalizedData(data = {}) {
  const records = data.records?.length ? data.records : [newRecord()]
  return {
    records: records.map((record) => ({
      id: record.id || `record-${Date.now()}-${recordSequence++}`,
      testDate: record.testDate || '',
      styleName: record.styleName || '',
      testedCount: record.testedCount ?? '',
      failedCount: record.failedCount ?? '',
      rate: record.rate ?? '',
      leakagePosition: record.leakagePosition || '',
      cause: record.cause || '',
      action: record.action || '',
      unit: record.unit || '',
      improvementDate: record.improvementDate || '',
      qc: record.qc || '',
    })),
  }
}

const loaded = loadDraft('analysis', emptyData())
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
  form.records.length > 0 &&
  form.records.every((record) =>
    hasValue(record.testDate) &&
    hasValue(record.styleName) &&
    hasValue(record.testedCount) &&
    hasValue(record.failedCount) &&
    hasValue(record.rate) &&
    hasValue(record.leakagePosition) &&
    hasValue(record.cause) &&
    hasValue(record.action) &&
    hasValue(record.unit) &&
    hasValue(record.improvementDate) &&
    hasValue(record.qc),
  ),
)
const savedLabel = computed(() => savedAt.value
  ? t('actions.savedAt', { time: new Intl.DateTimeFormat(dateLocale(), { hour: '2-digit', minute: '2-digit', day: '2-digit', month: '2-digit' }).format(new Date(savedAt.value)) })
  : t('actions.noDraft'))

function addRecord() {
  form.records.push(newRecord())
}

function removeRecord(index) {
  if (form.records.length > 1) form.records.splice(index, 1)
}

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
  savedAt.value = saveDraft('analysis', form)
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
    const result = await submitAnalysisForm(
      form,
      isReview.value ? Number(route.query.id) : 0,
      isReview.value,
    )
    clearDraft('analysis')
    savedAt.value = null
    if (previousRevision) editHistory.value.push(previousRevision)
    reviewSnapshot.value = cloneFormData(form)
    isEditing.value = false
    viewingRevision.value = false
    await router.replace({
      name: 'laboratory-analysis-form',
      query: {
        review: '1',
        id: result.recordKey,
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
  clearDraft('analysis')
  clearMissingFields(formElement.value)
  savedAt.value = null
  showNotice(t('actions.resetDone'), 'neutral')
}

onMounted(async () => {
  if (!isReview.value) return
  try {
    const detail = await getAnalysisForm(route.query.id)
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
        <div class="excel-sheet analysis-sheet">
          <h2>BIỂU PHÂN TÍCH NGUYÊN NHÂN VÀ CẢI THIỆN THỬ NGHIỆM LI TÂM / THỬ NƯỚC</h2>
          <p class="sheet-subtitle">袜套测试 / 离心测试的原因分析与改善表</p>

          <div class="analysis-table-layout">
          <table class="excel-table analysis-grid">
            <thead>
              <tr>
                <th>Ngày thử nghiệm<br /><small>测试日期</small></th>
                <th>Dạng giày<br /><small>型体 / Style name</small></th>
                <th>Số lượng thử nghiệm<br /><small>测试数量</small></th>
                <th>Số lượng không đạt<br /><small>不良率</small></th>
                <th>Tỉ lệ %<br /><small>比率</small></th>
                <th>Vị trí rỉ nước<br /><small>漏水位置</small></th>
                <th>Phân tích nguyên nhân<br /><small>原因分析</small></th>
                <th>Biện pháp xử lý<br /><small>处理措施</small></th>
                <th>Đơn vị phát sinh / cải thiện<br /><small>发生问题 / 改善单位</small></th>
                <th>Ngày cải thiện<br /><small>改善日期</small></th>
                <th>QC</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(row, index) in form.records" :key="row.id">
                <td><input v-revision-change="revisionChange(`records.${index}.testDate`)" v-model="row.testDate" type="date" required :aria-label="`Ngày thử nghiệm dòng ${index + 1}`" /></td>
                <td><input v-revision-change="revisionChange(`records.${index}.styleName`)" v-model.trim="row.styleName" type="text" required :aria-label="`Dạng giày dòng ${index + 1}`" /></td>
                <td><input v-revision-change="revisionChange(`records.${index}.testedCount`)" v-model="row.testedCount" type="text" required :aria-label="`Số lượng thử dòng ${index + 1}`" /></td>
                <td><input v-revision-change="revisionChange(`records.${index}.failedCount`)" v-model="row.failedCount" type="text" required :aria-label="`Số lượng không đạt dòng ${index + 1}`" /></td>
                <td><input v-revision-change="revisionChange(`records.${index}.rate`)" v-model="row.rate" type="text" required :aria-label="`Tỷ lệ dòng ${index + 1}`" /></td>
                <td><input v-revision-change="revisionChange(`records.${index}.leakagePosition`)" v-model.trim="row.leakagePosition" type="text" required :aria-label="`Vị trí rỉ nước dòng ${index + 1}`" /></td>
                <td><textarea v-revision-change="revisionChange(`records.${index}.cause`)" v-model.trim="row.cause" rows="4" required :aria-label="`Phân tích nguyên nhân dòng ${index + 1}`"></textarea></td>
                <td><textarea v-revision-change="revisionChange(`records.${index}.action`)" v-model.trim="row.action" rows="4" required :aria-label="`Biện pháp xử lý dòng ${index + 1}`"></textarea></td>
                <td><textarea v-revision-change="revisionChange(`records.${index}.unit`)" v-model.trim="row.unit" rows="4" required :aria-label="`Đơn vị cải thiện dòng ${index + 1}`"></textarea></td>
                <td><input v-revision-change="revisionChange(`records.${index}.improvementDate`)" v-model="row.improvementDate" type="date" required :aria-label="`Ngày cải thiện dòng ${index + 1}`" /></td>
                <td><input v-revision-change="revisionChange(`records.${index}.qc`)" v-model.trim="row.qc" type="text" required :aria-label="`QC dòng ${index + 1}`" /></td>
              </tr>
            </tbody>
          </table>
          <div class="analysis-delete-rail" aria-label="Delete rows">
            <div v-for="(row, index) in form.records" :key="`delete-${row.id}`" class="analysis-delete-slot">
              <button
                type="button"
                class="row-remove"
                :disabled="form.records.length === 1"
                :aria-label="`Xóa dòng ${index + 1}`"
                :title="`Xóa dòng ${index + 1}`"
                @click="removeRecord(index)"
              >×</button>
            </div>
          </div>
          </div>
          <button type="button" class="sheet-add" @click="addRecord">+ {{ t('actions.addRow') }}</button>
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
