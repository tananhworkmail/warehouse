<script setup>
import { computed, onMounted, ref } from 'vue'
import ComparisonChart from '../components/dashboard/ComparisonChart.vue'
import HorizontalRateChart from '../components/dashboard/HorizontalRateChart.vue'
import ParetoChart from '../components/dashboard/ParetoChart.vue'
import RateTrendChart from '../components/dashboard/RateTrendChart.vue'
import { getApiErrorMessage, getGoreTexWeeklyDashboard } from '../utils/api'
import { useGoreTexI18n } from '../utils/i18n'

const { t, dateLocale } = useGoreTexI18n()

const loading = ref(false)
const error = ref('')
const dashboard = ref(null)

function isoWeekValue(date = new Date()) {
  const utc = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()))
  const day = utc.getUTCDay() || 7
  utc.setUTCDate(utc.getUTCDate() + 4 - day)
  const yearStart = new Date(Date.UTC(utc.getUTCFullYear(), 0, 1))
  const week = Math.ceil((((utc - yearStart) / 86400000) + 1) / 7)
  return `${utc.getUTCFullYear()}-W${String(week).padStart(2, '0')}`
}

const selectedWeek = ref(isoWeekValue())
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
const paretoItems = computed(() => (dashboard.value?.visualizationResults || []).map((item) => ({
  ...item,
  label: defectKeys[item.label] ? t(`defects.${defectKeys[item.label]}`) : item.label,
})))

async function loadDashboard() {
  if (!weekParts.value || loading.value) return
  loading.value = true
  error.value = ''
  try {
    dashboard.value = await getGoreTexWeeklyDashboard(weekParts.value.year, weekParts.value.week)
  } catch (requestError) {
    error.value = getApiErrorMessage(requestError)
  } finally {
    loading.value = false
  }
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
          <input v-model="selectedWeek" type="week" @change="loadDashboard" />
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
        <header><span class="chart-number">02</span><div><h2>{{ t('dashboard.chart2') }}</h2><p>SUTER PASS RATE TREND</p></div></header>
        <RateTrendChart :items="dashboard.suterTrend" :empty-text="t('dashboard.noTrend')" :target-text="t('dashboard.target')" :passed-text="t('dashboard.passed')" />
      </article>
      <article class="dashboard-chart-card">
        <header><span class="chart-number">03</span><div><h2>{{ t('dashboard.chart3') }}</h2><p>SUTER COMPARISON</p></div></header>
        <ComparisonChart :series="comparisonSeries" :empty-text="t('dashboard.noCompare')" :passed-text="t('dashboard.passed')" />
      </article>
      <article class="dashboard-chart-card">
        <header><span class="chart-number">04</span><div><h2>{{ t('dashboard.chart4') }}</h2><p>R.RDY PASS RATE TREND</p></div></header>
        <RateTrendChart :items="dashboard.rRdyTrend" :empty-text="t('dashboard.noTrend')" :target-text="t('dashboard.target')" :passed-text="t('dashboard.passed')" />
      </article>
      <article class="dashboard-chart-card dashboard-pareto-card">
        <header><span class="chart-number">05</span><div><h2>{{ t('dashboard.chart5') }}</h2><p>VISUALIZATION RESULTS · PARETO</p></div></header>
        <ParetoChart :items="paretoItems" :empty-text="t('dashboard.noDefect')" :defects-text="t('dashboard.defects')" :cumulative-text="t('dashboard.cumulative')" />
      </article>
    </section>
  </div>
</template>
