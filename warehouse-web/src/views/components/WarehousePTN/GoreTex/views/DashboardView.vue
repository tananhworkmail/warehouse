<script setup>
import { computed, onMounted, ref } from 'vue'
import ComparisonChart from '../components/dashboard/ComparisonChart.vue'
import HorizontalRateChart from '../components/dashboard/HorizontalRateChart.vue'
import ParetoChart from '../components/dashboard/ParetoChart.vue'
import RateTrendChart from '../components/dashboard/RateTrendChart.vue'
import { getApiErrorMessage, getGoreTexDashboardPeriod, getGoreTexWeeklyDashboard } from '../utils/api'
import { useGoreTexI18n } from '../utils/i18n'

const { t, dateLocale } = useGoreTexI18n()

const loading = ref(false)
const error = ref('')
const dashboard = ref(null)
const trendPeriod = ref('week')
const paretoPeriod = ref('week')
const periodOptions = ['day', 'week', 'month']
const trendPeriodValue = ref('')
const paretoPeriodValue = ref('')
const trendPeriodItems = ref([])
const paretoPeriodItems = ref([])
const paretoStyles = ref([])
const paretoStyle = ref('')
const trendLoading = ref(false)
const paretoLoading = ref(false)
const trendError = ref('')
const paretoError = ref('')
let trendRequestId = 0
let paretoRequestId = 0

function isoWeekValue(date = new Date()) {
  const utc = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()))
  const day = utc.getUTCDay() || 7
  utc.setUTCDate(utc.getUTCDate() + 4 - day)
  const yearStart = new Date(Date.UTC(utc.getUTCFullYear(), 0, 1))
  const week = Math.ceil((((utc - yearStart) / 86400000) + 1) / 7)
  return `${utc.getUTCFullYear()}-W${String(week).padStart(2, '0')}`
}

const selectedWeek = ref(isoWeekValue())
trendPeriodValue.value = selectedWeek.value
paretoPeriodValue.value = selectedWeek.value
const weekParts = computed(() => {
  const match = selectedWeek.value.match(/^(\d{4})-W(\d{2})$/)
  return match ? { year: Number(match[1]), week: Number(match[2]) } : null
})

const dateRange = computed(() => {
  if (!dashboard.value) return ''
  const format = (value) => {
    const [year, month, day] = String(value || '').split('-').map(Number)
    if (!year || !month || !day) return value || ''
    return new Intl.DateTimeFormat(dateLocale(), { day: '2-digit', month: '2-digit', year: 'numeric' })
      .format(new Date(year, month - 1, day))
  }
  return `${format(dashboard.value.fromDate)} – ${format(dashboard.value.toDate)}`
})

const comparisonSeries = computed(() => (dashboard.value?.suterComparison || []).map((series) => ({
  ...series,
  label: t('dashboard.week', series),
})))

const trendItems = computed(() => trendPeriodItems.value)
const trendEmptyText = computed(() => trendError.value || (trendLoading.value ? t('common.loading') : t('dashboard.noTrend')))
const paretoEmptyText = computed(() => paretoError.value || (paretoLoading.value ? t('common.loading') : t('dashboard.noDefect')))

const defectKeys = {
  'Mũi thấm nước': 'toe',
  'Gót thấm nước': 'heel',
  'Hong trong thấm nước': 'medial',
  'Hong ngoài thấm nước': 'lateral',
  'Vật tư không đạt': 'material',
  'Dán đế lệch': 'attaching',
  'Ép đế nhăn': 'wrinkled',
  'Zíc zắc hở': 'zigzag',
  'Hở keo đế': 'bonding',
}
const paretoItems = computed(() => (
  paretoPeriodItems.value
).map((item) => ({
  ...item,
  label: defectKeys[item.label] ? t(`defects.${defectKeys[item.label]}`) : item.label,
})))

function periodInputType(period) {
  return period === 'day' ? 'date' : period
}

function trendInputType(period) {
  if (period === 'day') return 'week'
  if (period === 'week') return 'month'
  return 'number'
}

const trendAnchorLabel = computed(() => {
  if (trendPeriod.value === 'day') return t('dashboard.selectTrendWeek')
  if (trendPeriod.value === 'week') return t('dashboard.selectTrendMonth')
  return t('dashboard.selectTrendYear')
})

function localDateValue(date) {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function selectedWeekMonday() {
  const parts = weekParts.value
  if (!parts) return ''
  const jan4 = new Date(Date.UTC(parts.year, 0, 4))
  const monday = new Date(jan4)
  monday.setUTCDate(jan4.getUTCDate() - ((jan4.getUTCDay() + 6) % 7) + ((parts.week - 1) * 7))
  return monday.toISOString().slice(0, 10)
}

function defaultParetoPeriodValue(period) {
  if (period === 'week') return selectedWeek.value
  const mondayText = selectedWeekMonday()
  if (!mondayText) return ''
  if (period === 'day') {
    const today = new Date()
    const todayText = localDateValue(today)
    const sunday = new Date(`${mondayText}T00:00:00`)
    sunday.setDate(sunday.getDate() + 6)
    return todayText >= mondayText && todayText <= localDateValue(sunday) ? todayText : mondayText
  }
  const thursday = new Date(`${mondayText}T00:00:00`)
  thursday.setDate(thursday.getDate() + 3)
  return `${thursday.getFullYear()}-${String(thursday.getMonth() + 1).padStart(2, '0')}`
}

function defaultTrendPeriodValue(period) {
  if (period === 'day') return selectedWeek.value
  const mondayText = selectedWeekMonday()
  if (!mondayText) return ''
  const thursday = new Date(`${mondayText}T00:00:00`)
  thursday.setDate(thursday.getDate() + 3)
  if (period === 'week') {
    return `${thursday.getFullYear()}-${String(thursday.getMonth() + 1).padStart(2, '0')}`
  }
  return String(thursday.getFullYear())
}

async function loadChartPeriod(chart) {
  const isTrend = chart === 'trend'
  const period = isTrend ? trendPeriod.value : paretoPeriod.value
  const value = isTrend ? trendPeriodValue.value : paretoPeriodValue.value
  const style = isTrend ? '' : paretoStyle.value
  if (!value) return
  const loadingRef = isTrend ? trendLoading : paretoLoading
  const errorRef = isTrend ? trendError : paretoError
  const requestId = isTrend ? ++trendRequestId : ++paretoRequestId
  loadingRef.value = true
  errorRef.value = ''
  if (isTrend) trendPeriodItems.value = []
  else paretoPeriodItems.value = []
  try {
    const result = await getGoreTexDashboardPeriod(period, value, style, isTrend ? 'compare' : 'detail')
    if (requestId !== (isTrend ? trendRequestId : paretoRequestId)) return
    if (isTrend) trendPeriodItems.value = result.suterTrend || []
    else {
      paretoStyles.value = result.styles || []
      paretoStyle.value = result.selectedStyle || ''
      paretoPeriodItems.value = result.visualizationResults || []
    }
  } catch (requestError) {
    if (requestId !== (isTrend ? trendRequestId : paretoRequestId)) return
    errorRef.value = getApiErrorMessage(requestError)
  } finally {
    if (requestId === (isTrend ? trendRequestId : paretoRequestId)) loadingRef.value = false
  }
}

function changeChartPeriod(chart, period) {
  if (chart === 'trend') {
    trendPeriod.value = period
    trendPeriodValue.value = defaultTrendPeriodValue(period)
  } else {
    paretoPeriod.value = period
    paretoPeriodValue.value = defaultParetoPeriodValue(period)
    paretoStyle.value = ''
  }
  loadChartPeriod(chart)
}

function resetChartPeriods() {
  trendPeriod.value = 'week'
  paretoPeriod.value = 'week'
  trendPeriodValue.value = defaultTrendPeriodValue('week')
  paretoPeriodValue.value = selectedWeek.value
  paretoStyle.value = ''
}

async function loadDashboard() {
  if (!weekParts.value || loading.value) return
  trendRequestId += 1
  paretoRequestId += 1
  trendLoading.value = false
  paretoLoading.value = false
  loading.value = true
  error.value = ''
  try {
    const initialTrendValue = defaultTrendPeriodValue('week')
    const [weeklyDashboard, trendDashboard, paretoDashboard] = await Promise.all([
      getGoreTexWeeklyDashboard(weekParts.value.year, weekParts.value.week),
      getGoreTexDashboardPeriod('week', initialTrendValue, '', 'compare'),
      getGoreTexDashboardPeriod('week', selectedWeek.value, '', 'detail'),
    ])
    dashboard.value = weeklyDashboard
    resetChartPeriods()
    trendPeriodItems.value = trendDashboard?.suterTrend || []
    paretoPeriodItems.value = paretoDashboard?.visualizationResults || []
    paretoStyles.value = paretoDashboard?.styles || []
  } catch (requestError) {
    error.value = getApiErrorMessage(requestError)
  } finally {
    loading.value = false
  }
}

function handleGlobalWeekChange() {
  loadDashboard()
}

function changeParetoPeriodValue() {
  paretoStyle.value = ''
  loadChartPeriod('pareto')
}

onMounted(loadDashboard)
</script>

<template>
  <div class="goretex-dashboard-page">
    <header class="dashboard-heading">
      <div>
        <span class="eyebrow">{{ t('dashboard.eyebrow') }}</span>
        <h1>{{ t('dashboard.title') }}</h1>
        <p v-if="dashboard">{{ t('dashboard.week', dashboard) }} · {{ dateRange }}</p>
        <p v-else>{{ t('dashboard.description') }}</p>
      </div>
      <div class="dashboard-controls">
        <label>
          <span>{{ t('dashboard.selectWeek') }}</span>
          <input v-model="selectedWeek" type="week" @change="handleGlobalWeekChange" />
        </label>
        <button class="dashboard-refresh" type="button" :disabled="loading" @click="loadDashboard">
          <svg viewBox="0 0 20 20" fill="none" aria-hidden="true"><path d="M16.5 6.2V2.8m0 0h-3.4m3.4 0-2.3 2.3a6 6 0 1 0 1.2 7" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" /></svg>
          {{ loading ? t('common.loading') : t('common.refresh') }}
        </button>
        <RouterLink class="btn btn-primary goretex-page-back" :to="{ name: 'laboratory-forms' }">
          <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
            <path d="M16 10H5m4-4-4 4 4 4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          {{ t('common.back') }}
        </RouterLink>
      </div>
    </header>

    <div v-if="error" class="dashboard-error">
      <strong>{{ t('dashboard.loadError') }}</strong><span>{{ error }}</span>
    </div>

    <div v-else-if="loading && !dashboard" class="dashboard-loading"><i></i><span>{{ t('dashboard.aggregating') }}</span></div>

    <section v-else-if="dashboard" class="dashboard-chart-grid">
      <article class="dashboard-chart-card">
        <header><span class="chart-number">01</span><div><h2>{{ t('dashboard.chart1') }}</h2><p>SUTER PASS RATE BY ITEMS</p></div></header>
        <HorizontalRateChart :items="dashboard.suterByItems" :empty-text="t('dashboard.noPass')" :passed-text="t('dashboard.passed')" />
      </article>
      <article class="dashboard-chart-card">
        <header>
          <span class="chart-number">02</span>
          <div><h2>{{ t('dashboard.chart2') }}</h2><p>SUTER PASS RATE TREND</p></div>
          <div class="dashboard-period-switch" :aria-label="t('dashboard.viewBy')">
            <button
              v-for="period in periodOptions"
              :key="period"
              type="button"
              :class="{ active: trendPeriod === period }"
              :disabled="trendLoading"
              @click="changeChartPeriod('trend', period)"
            >{{ t(`dashboard.${period}View`) }}</button>
          </div>
          <label class="dashboard-period-picker">
            <span>{{ trendAnchorLabel }}</span>
            <input
              v-model="trendPeriodValue"
              :type="trendInputType(trendPeriod)"
              :min="trendPeriod === 'month' ? 2000 : undefined"
              :max="trendPeriod === 'month' ? 2200 : undefined"
              :disabled="trendLoading"
              @change="loadChartPeriod('trend')"
            />
          </label>
        </header>
        <RateTrendChart :items="trendItems" :empty-text="trendEmptyText" :target-text="t('dashboard.target')" :passed-text="t('dashboard.passed')" />
      </article>
      <article class="dashboard-chart-card">
        <header><span class="chart-number">03</span><div><h2>{{ t('dashboard.chart3') }}</h2><p>SUTER COMPARISON</p></div></header>
        <ComparisonChart :series="comparisonSeries" :empty-text="t('dashboard.noCompare')" :passed-text="t('dashboard.passed')" />
      </article>
      <article class="dashboard-chart-card dashboard-pareto-card">
        <header>
          <span class="chart-number">04</span>
          <div><h2>{{ t('dashboard.chart5') }}</h2><p>VISUALIZATION RESULTS · PARETO</p></div>
          <div class="dashboard-period-switch" :aria-label="t('dashboard.viewBy')">
            <button
              v-for="period in periodOptions"
              :key="period"
              type="button"
              :class="{ active: paretoPeriod === period }"
              :disabled="paretoLoading"
              @click="changeChartPeriod('pareto', period)"
            >{{ t(`dashboard.${period}View`) }}</button>
          </div>
          <label class="dashboard-period-picker">
            <span>{{ t('dashboard.selectPeriod') }}</span>
            <input
              v-model="paretoPeriodValue"
              :type="periodInputType(paretoPeriod)"
              :disabled="paretoLoading"
              @change="changeParetoPeriodValue"
            />
          </label>
          <label class="dashboard-period-picker dashboard-style-picker">
            <span>{{ t('dashboard.styleFilter') }}</span>
            <select v-model="paretoStyle" :disabled="paretoLoading" @change="loadChartPeriod('pareto')">
              <option value="">{{ t('dashboard.allStyles') }}</option>
              <option v-for="style in paretoStyles" :key="style" :value="style">{{ style }}</option>
            </select>
          </label>
        </header>
        <ParetoChart :items="paretoItems" :empty-text="paretoEmptyText" :defects-text="t('dashboard.defects')" :percentage-text="t('dashboard.percentage')" />
      </article>
    </section>
  </div>
</template>
