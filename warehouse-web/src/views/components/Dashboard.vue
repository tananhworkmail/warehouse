<template>
  <div class="dash-root" ref="rootEl">
    <!-- Three.js Background -->
    <canvas ref="bgCanvas" class="bg-canvas"></canvas>

    <!-- Scanlines + Vignette Overlay -->
    <div class="fx-overlay"></div>

    <div class="dash-content">
      <!-- ========== HEADER ========== -->
      <header class="dash-header glass-panel">
        <div class="brand">
          <img ref="logoEl" src="@/assets/Logo.png" alt="Logo" class="brand-logo" />
          <div class="brand-text">
            <span class="brand-eyebrow">ENVIRONMENT MONITORING SYSTEM</span>
            <h1 class="brand-title">
              <span class="title-neon">NHIỆT ẨM</span> DASHBOARD
            </h1>
          </div>
        </div>

        <!-- PHẦN BÊN PHẢI ĐÃ ĐƯỢC THU GỌN & ĐỒNG BỘ -->
        <div class="header-right">
          <div class="header-controls">
            <!-- Date picker -->
            <div class="date-picker glass-panel-inner">
              <input type="date" v-model="selectedDate" @change="onDateChange" />
            </div>
            <!-- Nút refresh -->
            <button class="btn-glass btn-icon" @click="refreshData" title="Làm mới dữ liệu">↻</button>
            <!-- Đồng hồ -->
            <div class="clock-block glass-panel-inner">
              <div class="clock-main">{{ clockLabel }}</div>
              <div class="clock-date">{{ dateLabel }}</div>
            </div>
            <!-- Đèn trạng thái -->
            <div class="power-lamp" :class="{ 'power-lamp--error': fetchError }">
              <span class="lamp-dot"></span>
              {{ fetchError ? 'LỖI' : 'LIVE' }}
            </div>
            <!-- Các nút điều hướng -->
            <button class="btn-glass btn-icon" @click="goBack" title="Quay lại">←</button>
            <button class="btn-glass btn-icon" @click="toggleFullscreen" title="Toàn màn hình">⛶</button>
          </div>
        </div>
      </header>

      <!-- ========== ALERT TICKER ========== -->
      <div class="alert-ticker glass-panel-inner" v-if="alerts.length">
        <span class="ticker-label">⚠ ALERTS</span>
        <div class="ticker-wrap">
          <div class="ticker-track" ref="tickerTrack">
            <span v-for="a in alerts" :key="a.id" class="ticker-item" :class="a.level">{{ a.message }}</span>
          </div>
        </div>
      </div>

      <!-- ========== STATS ROW ========== -->
      <section class="stats-row">
        <div class="stat-tile glass-panel">
          <span class="stat-icon">🏭</span>
          <div class="stat-value">{{ rooms.length }}</div>
          <span class="stat-label">Khu vực</span>
        </div>
        <div class="stat-tile glass-panel">
          <span class="stat-icon">📡</span>
          <div class="stat-value">{{ totalSensors }}</div>
          <span class="stat-label">Cảm biến</span>
        </div>
        <div class="stat-tile glass-panel">
          <span class="stat-icon">📊</span>
          <div class="stat-value" :class="toneClass(overallStatus)">{{ overallLabel }}</div>
          <span class="stat-label">Tổng thể</span>
        </div>
        <div class="stat-tile glass-panel">
          <span class="stat-icon">⏳</span>
          <div class="stat-value mono-sm">{{ uptime }}</div>
          <span class="stat-label">Uptime</span>
        </div>
      </section>
      <!-- ========== ROOM CARDS THEO NHÓM ========== -->
      <template v-for="group in roomGroups" :key="group.key">
        <section class="group-header-row" v-if="group.items.length">
          <h2 class="group-title">
            <span class="group-title-icon">{{ group.icon }}</span>
            {{ group.title }}
          </h2>
        </section>
        <section class="rooms-grid">
          <article
            v-for="room in group.items"
            :key="room.id"
            class="room-card glass-panel"
            :class="'status-' + room.status"
            :ref="(el) => setCardRef(room.id, el)"
            @mousemove="onCardMove($event, room.id)"
            @mouseleave="onCardLeave(room.id)"
          >
            <!-- Card header -->
            <div class="card-header">
              <div class="nameplate">
                <span class="room-icon" v-html="room.icon"></span>
                <div>
                  <h3>{{ room.name }}</h3>
                  <span class="room-code">{{ room.code }}</span>
                </div>
              </div>
              <div class="status-chip" :class="'chip-' + room.status">
                {{ statusLabel(room.status) }}
              </div>
            </div>

            <!-- Dual Gauges Row -->
            <div class="gauges-row">
              <div class="gauge-wrapper">
                <v-chart class="gauge-chart" :option="gaugeOption(room.temp, 'NHIỆT ĐỘ · °C', '°')" autoresize />
              </div>
              <div class="gauge-wrapper">
                <v-chart class="gauge-chart" :option="gaugeOption(room.humidity, 'ĐỘ ẨM · %RH', '%')" autoresize />
              </div>
            </div>

            <!-- Sparkline history -->
            <div class="recorder">
              <div class="recorder-head">
                <span>📋 Lịch sử {{ HISTORY_LEN }} mẫu</span>
                <span class="recorder-time">{{ relativeTime(room.updatedAt) }}</span>
              </div>
              <div class="recorder-grid">
                <!-- Temperature mini-chart -->
                <div class="recorder-panel">
                  <div class="recorder-panel__label"><i class="dot temp"></i> Nhiệt độ</div>
                  <svg viewBox="0 0 300 100" class="recorder-svg" preserveAspectRatio="none">
                    <defs>
                      <linearGradient :id="'tempGrad-' + room.id" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="var(--pen-temp)" stop-opacity="0.28" />
                        <stop offset="100%" stop-color="var(--pen-temp)" stop-opacity="0.02" />
                      </linearGradient>
                    </defs>
                    <g transform="translate(0,22)">
                      <g class="grid-lines">
                        <line v-for="gx in GRID_XS" :key="'tgx' + gx" :x1="gx" y1="0" :x2="gx" y2="60" />
                        <line v-for="gy in GRID_YS_SPARK" :key="'tgy' + gy" x1="0" :y1="gy" x2="300" :y2="gy" />
                      </g>
                      <path class="pen-area pen-area--temp" :d="smoothAreaPath(room.temp.history, room.temp.min, room.temp.max, 300, 60)" :fill="'url(#tempGrad-' + room.id + ')'" />
                      <path class="pen-trace pen-temp" :d="smoothLinePath(room.temp.history, room.temp.min, room.temp.max, 300, 60)" :ref="(el) => setTraceRef(room.id, 'temp', el)" />
                      <template v-for="(d, i) in sparkDots(room.temp.history, room.temp.timeHistory, room.temp.min, room.temp.max, 300, 60)" :key="'tdot' + room.id + i">
                        <circle :cx="d.x" :cy="d.y" r="2.8" class="spark-dot spark-dot--temp"><title>{{ d.label }}°C</title></circle>
                        <text v-if="i % 2 === 0" :x="d.x" :y="d.y - 8" class="dot-label dot-label--temp" text-anchor="middle">
                          <tspan :x="d.x">{{ d.value }}°</tspan>
                          <tspan :x="d.x" dy="10">{{ d.shortTime }}</tspan>
                        </text>
                      </template>
                    </g>
                  </svg>
                  <div class="recorder-axis">
                    <span v-for="(lbl, i) in axisTimeLabels(room.temp.timeHistory)" :key="'taxis' + i">{{ lbl }}</span>
                  </div>
                </div>

                <!-- Humidity mini-chart -->
                <div class="recorder-panel">
                  <div class="recorder-panel__label"><i class="dot hum"></i> Độ ẩm</div>
                  <svg viewBox="0 0 300 100" class="recorder-svg" preserveAspectRatio="none">
                    <defs>
                      <linearGradient :id="'humGrad-' + room.id" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="var(--pen-hum)" stop-opacity="0.28" />
                        <stop offset="100%" stop-color="var(--pen-hum)" stop-opacity="0.02" />
                      </linearGradient>
                    </defs>
                    <g transform="translate(0,22)">
                      <g class="grid-lines">
                        <line v-for="gx in GRID_XS" :key="'hgx' + gx" :x1="gx" y1="0" :x2="gx" y2="60" />
                        <line v-for="gy in GRID_YS_SPARK" :key="'hgy' + gy" x1="0" :y1="gy" x2="300" :y2="gy" />
                      </g>
                      <path class="pen-area pen-area--hum" :d="smoothAreaPath(room.humidity.history, room.humidity.min, room.humidity.max, 300, 60)" :fill="'url(#humGrad-' + room.id + ')'" />
                      <path class="pen-trace pen-hum" :d="smoothLinePath(room.humidity.history, room.humidity.min, room.humidity.max, 300, 60)" :ref="(el) => setTraceRef(room.id, 'humidity', el)" />
                      <template v-for="(d, i) in sparkDots(room.humidity.history, room.humidity.timeHistory, room.humidity.min, room.humidity.max, 300, 60)" :key="'hdot' + room.id + i">
                        <circle :cx="d.x" :cy="d.y" r="2.8" class="spark-dot spark-dot--hum"><title>{{ d.label }}%</title></circle>
                        <text v-if="i % 2 === 0" :x="d.x" :y="d.y - 8" class="dot-label dot-label--hum" text-anchor="middle">
                          <tspan :x="d.x">{{ d.value }}%</tspan>
                          <tspan :x="d.x" dy="10">{{ d.shortTime }}</tspan>
                        </text>
                      </template>
                    </g>
                  </svg>
                  <div class="recorder-axis">
                    <span v-for="(lbl, i) in axisTimeLabels(room.humidity.timeHistory)" :key="'haxis' + i">{{ lbl }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Threshold footer -->
            <div class="card-footer">
              🎯 Lý tưởng:
              <b>{{ room.temp.idealMin }}–{{ room.temp.idealMax }}°C</b> ·
              <b>{{ room.humidity.idealMin }}–{{ room.humidity.idealMax }}%RH</b>
            </div>
          </article>
        </section>
      </template>

      <!-- ========== FOOTER ========== -->
      <footer class="dash-footer glass-panel">
        <span>© 2026 ENV-MON v3.1 · Kết nối bảo mật · Chu kỳ {{ POLL_INTERVAL_SEC / 60 }} phút</span>
        <span class="footer-right">Cập nhật tiếp: {{ countdownLabel }}</span>
      </footer>
    </div>
  </div>
</template>

<script setup>
// ==================== IMPORTS ====================
import { reactive, ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import gsap from 'gsap'
import * as THREE from 'three'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { GaugeChart } from 'echarts/charts'

use([CanvasRenderer, GaugeChart])

// ==================== CONSTANTS ====================
const API_URL = (import.meta.env.VITE_API_URL || '/api/v1').replace(/\/$/, '')
const HISTORY_LEN = 16
const OVERVIEW_MAX_POINTS = 48
const OVERVIEW_INTERVAL_MINUTES = 30
const GRID_XS = [0, 50, 100, 150, 200, 250, 300]
const GRID_YS_SPARK = [0, 15, 30, 45, 60]
const POLL_INTERVAL_SEC = 300

// ==================== UTILITY FUNCTIONS ====================
function round1(n) { return Math.round(n * 10) / 10 }
function clamp(n, min, max) { return Math.min(max, Math.max(min, n)) }

function metricStatus(metric) {
  const spread = (metric.idealMax - metric.idealMin) * 0.35
  if (metric.value < metric.idealMin - spread || metric.value > metric.idealMax + spread) return 'critical'
  if (metric.value < metric.idealMin || metric.value > metric.idealMax) return 'warning'
  return 'normal'
}
function worseStatus(a, b) {
  const rank = { normal: 0, warning: 1, critical: 2 }
  return rank[a] >= rank[b] ? a : b
}
function toneClass(status) {
  return 'tone-' + (status === 'normal' ? 'ok' : status === 'warning' ? 'warn' : 'crit')
}
function statusLabel(s) {
  const map = { normal: 'ỔN ĐỊNH', warning: '⚠ CẢNH BÁO', critical: '⛔ NGUY CẤP' }
  return map[s]
}

// SVG path helpers
function toXY(history, min, max, w, h) {
  const pad = 6
  if (!history.length) return []
  const stepX = (w - 2 * pad) / Math.max(history.length - 1, 1)
  return history.map((v, i) => {
    const x = pad + i * stepX
    const pct = clamp((v - min) / (max - min), 0, 1)
    const y = h - pad - pct * (h - 2 * pad)
    return { x, y }
  })
}
function catmullRomPath(points) {
  if (points.length < 2) return ''
  if (points.length === 2) return `M ${points[0].x},${points[0].y} L ${points[1].x},${points[1].y}`
  let d = `M ${points[0].x},${points[0].y}`
  for (let i = 0; i < points.length - 1; i++) {
    const p0 = points[i - 1] || points[i]
    const p1 = points[i]
    const p2 = points[i + 1]
    const p3 = points[i + 2] || p2
    const c1x = p1.x + (p2.x - p0.x) / 6
    const c1y = p1.y + (p2.y - p0.y) / 6
    const c2x = p2.x - (p3.x - p1.x) / 6
    const c2y = p2.y - (p3.y - p1.y) / 6
    d += ` C ${c1x.toFixed(2)},${c1y.toFixed(2)} ${c2x.toFixed(2)},${c2y.toFixed(2)} ${p2.x.toFixed(2)},${p2.y.toFixed(2)}`
  }
  return d
}
function smoothLinePath(history, min, max, w, h) { return catmullRomPath(toXY(history, min, max, w, h)) }
function smoothAreaPath(history, min, max, w, h) {
  const pts = toXY(history, min, max, w, h)
  if (!pts.length) return ''
  const line = catmullRomPath(pts)
  const last = pts[pts.length - 1]
  const first = pts[0]
  return `${line} L ${last.x.toFixed(2)},${h} L ${first.x.toFixed(2)},${h} Z`
}
function axisTimeLabels(timeHistory) {
  if (!timeHistory?.length) return ['', '', '']
  const fmt = (t) => t.toLocaleTimeString('vi-VN', { hour12: false })
  return [fmt(timeHistory[0]), fmt(timeHistory[Math.floor(timeHistory.length / 2)]), fmt(timeHistory[timeHistory.length - 1])]
}
function sparkDots(history, timeHistory, min, max, w, h) {
  const pad = 6
  if (!history.length) return []
  const stepX = (w - 2 * pad) / Math.max(history.length - 1, 1)
  return history.map((v, i) => {
    const x = pad + i * stepX
    const pct = clamp((v - min) / (max - min), 0, 1)
    const y = h - pad - pct * (h - 2 * pad)
    const t = timeHistory?.[i]
    const timeStr = t ? t.toLocaleTimeString('vi-VN', { hour12: false }) : ''
    const shortTime = t ? t.toLocaleTimeString('vi-VN', { hour12: false, hour: '2-digit', minute: '2-digit' }) : ''
    return { x, y, value: v, shortTime, label: `${timeStr} — ${v}` }
  })
}

// Gauge helpers
function gaugeColor(metric) {
  const st = metricStatus(metric)
  if (st === 'normal') return '#4ee4a0'
  if (st === 'warning') return '#e0a940'
  return '#f04a5a'
}
function gaugeOption(metric, title, unit) {
  const color = gaugeColor(metric)
  return {
    series: [{
      type: 'gauge',
      startAngle: 225,
      endAngle: -45,
      min: metric.min,
      max: metric.max,
      radius: '92%',
      pointer: { width: 4, length: '55%', itemStyle: { color: '#e8edf4' } },
      progress: { show: true, width: 10, itemStyle: { color, shadowBlur: 8, shadowColor: color } },
      axisLine: { lineStyle: { width: 10, color: [[1, 'rgba(255,255,255,0.06)']] } },
      axisTick: { distance: -10, length: 5, lineStyle: { color: 'rgba(255,255,255,0.18)', width: 1 } },
      splitLine: { distance: -10, length: 10, lineStyle: { color: 'rgba(255,255,255,0.25)', width: 1.5 } },
      axisLabel: { distance: 16, color: '#7f8ea3', fontSize: 8, fontFamily: 'JetBrains Mono, monospace' },
      anchor: { show: true, size: 12, itemStyle: { color: '#e8edf4', borderWidth: 0 }, showAbove: true },
      title: { show: true, offsetCenter: [0, '68%'], color: '#7f8ea3', fontSize: 7, fontFamily: 'JetBrains Mono, monospace', fontWeight: 600 },
      detail: {
        valueAnimation: true,
        offsetCenter: [0, '32%'],
        formatter: (v) => v.toFixed(1) + unit,
        color: '#e8edf4',
        fontSize: 20,
        fontWeight: 800,
        fontFamily: 'JetBrains Mono, monospace'
      },
      data: [{ value: metric.value, name: title }],
      animationDuration: 700,
      animationEasing: 'cubicOut'
    }]
  }
}

// ==================== SVG ICONS ====================
const ICON_BEAKER = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H17.5a2.5 2.5 0 0 1 0 5H6.5A2.5 2.5 0 0 1 4 19.5Z"/><path d="M8 3v4.5l-3 9"/><path d="M16 3v4.5l3 9"/></svg>'
const ICON_DROP = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M12 3c4 5 6.5 8.6 6.5 11.6A6.5 6.5 0 0 1 5.5 14.6C5.5 11.6 8 8 12 3Z"/></svg>'
const ICON_BOX = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M3 8.5 12 4l9 4.5-9 4.5-9-4.5Z"/><path d="M3 8.5V16l9 4.5V13"/><path d="M21 8.5V16l-9 4.5"/></svg>'

// ==================== REACTIVE STATE ====================
const rootEl = ref(null)
const bgCanvas = ref(null)
const tickerTrack = ref(null)
const logoEl = ref(null)

const prefersReducedMotion = window.matchMedia?.('(prefers-reduced-motion: reduce)').matches ?? false

const now = ref(new Date())
const clockLabel = computed(() => now.value.toLocaleTimeString('vi-VN', { hour12: false }))
const dateLabel = computed(() => now.value.toLocaleDateString('vi-VN', { weekday: 'short', day: '2-digit', month: '2-digit', year: 'numeric' }))

const countdown = ref(POLL_INTERVAL_SEC)
const countdownLabel = computed(() => {
  const m = String(Math.floor(countdown.value / 60)).padStart(2, '0')
  const s = String(countdown.value % 60).padStart(2, '0')
  return `${m}:${s}`
})

const uptime = ref('00:00:00')
let startTime = Date.now()

const todayStr = () => new Date().toISOString().split('T')[0]
const selectedDate = ref(todayStr())
const selectedDateLabel = computed(() => selectedDate.value === todayStr() ? 'hôm nay' : selectedDate.value)

const isLoadingData = ref(false)
const fetchError = ref('')

const alerts = reactive([])
let alertIdCounter = 0
function pushAlert(room, type) {
  const level = metricStatus(room[type])
  if (level === 'normal') return
  const value = room[type].value
  const msg = `${room.name}: ${type === 'temp' ? 'Nhiệt độ' : 'Độ ẩm'} ${value.toFixed(1)}${type === 'temp' ? '°C' : '%'} vượt ngưỡng`
  if (alerts.length && alerts[0].message === msg) return
  alerts.unshift({ id: ++alertIdCounter, level, message: msg })
  if (alerts.length > 8) alerts.pop()
}

// ==================== ROOMS DATA ====================
const rooms = reactive([
  {
    id: 'kvt02', group: 'warehouse', deviceName: 'KVT02',
    name: 'Kho Vật Tư 02', code: 'KVT02', icon: ICON_BOX, status: 'normal', updatedAt: null,
    temp: { value: 0, min: 10, max: 45, idealMin: 25, idealMax: 30, history: [], timeHistory: [] },
    humidity: { value: 0, min: 15, max: 95, idealMin: 45, idealMax: 60, history: [], timeHistory: [] }
  },
  {
    id: 'kvt03', group: 'warehouse', deviceName: 'KVT03',
    name: 'Kho Vật Tư 03', code: 'KVT03', icon: ICON_BOX, status: 'normal', updatedAt: null,
    temp: { value: 0, min: 10, max: 45, idealMin: 25, idealMax: 30, history: [], timeHistory: [] },
    humidity: { value: 0, min: 15, max: 95, idealMin: 45, idealMax: 60, history: [], timeHistory: [] }
  },
  {
    id: 'condition-room', group: 'lab', deviceName: 'ThiNghiem9527_L4',
    name: 'Condition Room', code: 'THINGHIEM9527_L4', icon: ICON_BEAKER, status: 'normal', updatedAt: null,
    temp: { value: 0, min: 10, max: 40, idealMin: 21, idealMax: 25, history: [], timeHistory: [] },
    humidity: { value: 0, min: 15, max: 95, idealMin: 45, idealMax: 65, history: [], timeHistory: [] }
  },
  {
    id: 'bonding-room', group: 'lab', deviceName: 'ThiNghiem9517_L4',
    name: 'Bonding Room', code: 'THINGHIEM9517_L4', icon: ICON_DROP, status: 'normal', updatedAt: null,
    temp: { value: 0, min: 10, max: 40, idealMin: 21, idealMax: 25, history: [], timeHistory: [] },
    humidity: { value: 0, min: 15, max: 95, idealMin: 45, idealMax: 65, history: [], timeHistory: [] }
  }
])

const roomGroups = computed(() => [
  { key: 'warehouse', title: 'Kho Vật Tư', icon: '📦', items: rooms.filter(r => r.group === 'warehouse') },
  { key: 'lab', title: 'Phòng Thí Nghiệm', icon: '🧪', items: rooms.filter(r => r.group === 'lab') }
])
const totalSensors = rooms.length * 2
const overallStatus = computed(() => rooms.reduce((acc, r) => worseStatus(acc, r.status), 'normal'))
const overallLabel = computed(() => statusLabel(overallStatus.value))

// ==================== OVERVIEW CHART DATA ====================
const overviewHistory = ref([])
function rebuildOverviewFromHistory() {
  const intervalMs = OVERVIEW_INTERVAL_MINUTES * 60 * 1000
  const points = []
  rooms.forEach(room => {
    room.temp.history.forEach((val, i) => {
      const time = room.temp.timeHistory[i]
      if (time && val != null) points.push({ temp: val, time: new Date(time).getTime() })
    })
  })
  if (!points.length) { overviewHistory.value = []; return }
  points.sort((a, b) => a.time - b.time)
  const buckets = new Map()
  points.forEach(p => {
    const bucketTime = Math.floor(p.time / intervalMs) * intervalMs
    if (!buckets.has(bucketTime)) buckets.set(bucketTime, { sum: 0, count: 0 })
    const b = buckets.get(bucketTime); b.sum += p.temp; b.count++
  })
  const sorted = [...buckets.entries()].sort(([a], [b]) => a - b).slice(-OVERVIEW_MAX_POINTS)
  overviewHistory.value = sorted.map(([timeMs, { sum, count }]) => ({ avg: round1(sum / count), time: new Date(timeMs) }))
}

const overviewLinePoints = computed(() => {
  const w = 800, h = 160, pad = 20
  const arr = overviewHistory.value
  if (!arr.length) return ''
  const stepX = (w - 2 * pad) / Math.max(arr.length - 1, 1)
  const allVals = arr.map(p => p.avg)
  const gMin = Math.min(...allVals) - 1
  const gMax = Math.max(...allVals) + 1
  return arr.map((p, i) => {
    const x = pad + i * stepX
    const pct = clamp((p.avg - gMin) / (gMax - gMin), 0, 1)
    const y = h - pad - pct * (h - 2 * pad)
    return `${x.toFixed(1)},${y.toFixed(1)}`
  }).join(' ')
})
const overviewAreaPoints = computed(() => {
  const line = overviewLinePoints.value
  return line ? `20,140 ${line} 780,140` : ''
})
const overviewDots = computed(() => {
  const w = 800, h = 160, pad = 20
  const arr = overviewHistory.value
  if (!arr.length) return []
  const stepX = (w - 2 * pad) / Math.max(arr.length - 1, 1)
  const allVals = arr.map(p => p.avg)
  const gMin = Math.min(...allVals) - 1
  const gMax = Math.max(...allVals) + 1
  return arr.map((p, i) => {
    const x = pad + i * stepX
    const pct = clamp((p.avg - gMin) / (gMax - gMin), 0, 1)
    const y = h - pad - pct * (h - 2 * pad)
    const timeStr = p.time ? p.time.toLocaleTimeString('vi-VN', { hour12: false }) : ''
    return { x, y, label: `${timeStr} — ${p.avg.toFixed(1)}°C` }
  })
})
const overviewGridX = [0, 100, 200, 300, 400, 500, 600, 700, 800]
const overviewGridY = [0, 40, 80, 120, 160]

// ==================== DOM REFS & ANIMATION ====================
const cardRefs = {}
const traceRefs = {}
const setCardRef = (id, el) => { if (el) cardRefs[id] = el }
const setTraceRef = (id, type, el) => { if (el) traceRefs[id + '-' + type] = el }

function drawTrace(room, type) {
  const el = traceRefs[room.id + '-' + type]
  if (!el) return
  const len = el.getTotalLength()
  el.style.strokeDasharray = len
  el.style.strokeDashoffset = len
  if (prefersReducedMotion) { el.style.strokeDashoffset = 0; return }
  gsap.to(el, { strokeDashoffset: 0, duration: 1.4, delay: 0.25, ease: 'power2.out' })
}
const pulseTweens = {}
function syncCardPulse(room) {
  const card = cardRefs[room.id]
  if (!card || prefersReducedMotion) return
  if (room.status !== 'normal') {
    if (!pulseTweens[room.id]) {
      const color = room.status === 'critical' ? 'rgba(240,74,90,0.5)' : 'rgba(224,169,64,0.5)'
      pulseTweens[room.id] = gsap.to(card, {
        boxShadow: `0 0 35px 5px ${color}, inset 0 0 20px ${color}`,
        repeat: -1, yoyo: true, duration: 1.2, ease: 'sine.inOut'
      })
    }
  } else if (pulseTweens[room.id]) {
    pulseTweens[room.id].kill()
    gsap.to(card, { boxShadow: '', duration: 0.6 })
    delete pulseTweens[room.id]
  }
}
function onCardMove(e, id) {
  if (prefersReducedMotion) return
  const card = cardRefs[id]
  if (!card) return
  const rect = card.getBoundingClientRect()
  const relX = (e.clientX - rect.left) / rect.width - 0.5
  const relY = (e.clientY - rect.top) / rect.height - 0.5
  gsap.to(card, { rotateY: relX * 3, rotateX: -relY * 3, y: -4, duration: 0.5, ease: 'power2.out' })
}
function onCardLeave(id) {
  const card = cardRefs[id]
  if (!card) return
  gsap.to(card, { rotateY: 0, rotateX: 0, y: 0, duration: 0.6, ease: 'power3.out' })
}
function relativeTime(date) {
  if (!date) return 'chưa có dữ liệu'
  const diff = Math.max(0, Math.floor((now.value - date) / 1000))
  if (diff < 5) return 'vừa xong'
  if (diff < 60) return diff + ' giây trước'
  return Math.floor(diff / 60) + ' phút trước'
}

// ==================== DATA FETCHING ====================
function normalizeList(raw) {
  if (Array.isArray(raw)) return raw
  if (raw && Array.isArray(raw.data)) return raw.data
  if (raw && typeof raw === 'object') return Object.values(raw)
  return []
}
async function fetchWarehouseData(useRealtime) {
  const warehouseRooms = rooms.filter(r => r.group === 'warehouse')
  if (!warehouseRooms.length) return
  const dateParam = new URLSearchParams({ date: selectedDate.value }).toString()
  const logUrl = `${API_URL}/warehouse/temphumiditylatestlogbydevices?${dateParam}`
  const curUrl = `${API_URL}/warehouse/temphumiditybydevices`
  const requests = [fetch(logUrl)]
  if (useRealtime) requests.push(fetch(curUrl))
  const [logRes, curRes] = await Promise.all(requests)
  const logJson = await logRes.json()
  const curJson = useRealtime ? await curRes.json() : null
  const curMap = {}
  if (curJson) normalizeList(curJson?.data).forEach(item => { if (item?.DeviceName) curMap[item.DeviceName] = item })
  const grouped = {}
  normalizeList(logJson?.data).forEach(item => {
    const device = item?.DeviceName
    if (!device) return
    const t = item.Temp ?? item.Tem
    if (!grouped[device]) grouped[device] = []
    grouped[device].push({
      temp: t != null ? round1(parseFloat(t)) : null,
      hum: item.Hum != null ? round1(parseFloat(item.Hum)) : null,
      time: item.RecordTime ? new Date(item.RecordTime) : new Date()
    })
  })
  Object.values(grouped).forEach(rows => rows.sort((a, b) => a.time - b.time))
  warehouseRooms.forEach(room => {
    const cur = curMap[room.deviceName]
    if (cur) {
      if (cur.Tem != null) room.temp.value = round1(parseFloat(cur.Tem))
      if (cur.Hum != null) room.humidity.value = round1(parseFloat(cur.Hum))
      if (cur.RecordTime) room.updatedAt = new Date(cur.RecordTime)
    }
    const rows = grouped[room.deviceName]
    if (rows?.length) {
      const tempRows = rows.filter(r => r.temp != null)
      const humRows = rows.filter(r => r.hum != null)
      room.temp.history = tempRows.slice(-HISTORY_LEN).map(r => r.temp)
      room.temp.timeHistory = tempRows.slice(-HISTORY_LEN).map(r => r.time)
      room.humidity.history = humRows.slice(-HISTORY_LEN).map(r => r.hum)
      room.humidity.timeHistory = humRows.slice(-HISTORY_LEN).map(r => r.time)
      if (!cur) {
        const latest = rows[rows.length - 1]
        if (latest.temp != null) room.temp.value = latest.temp
        if (latest.hum != null) room.humidity.value = latest.hum
        room.updatedAt = latest.time
      }
    }
    room.status = worseStatus(metricStatus(room.temp), metricStatus(room.humidity))
    pushAlert(room, 'temp')
    pushAlert(room, 'humidity')
  })
}
async function fetchLabData() {
  const labRooms = rooms.filter(r => r.group === 'lab')
  if (!labRooms.length) return
  const params = new URLSearchParams({ date: selectedDate.value }).toString()
  const res = await fetch(`${API_URL}/warehouse/laboratory/environment?${params}`)
  const json = await res.json()
  const arr = json.data ?? []
  const grouped = {}
  arr.forEach(x => {
    const device = x.DeviceName
    if (!device) return
    if (!grouped[device]) grouped[device] = []
    grouped[device].push({
      temp: x.Tem != null ? round1(parseFloat(x.Tem)) : null,
      hum: x.Hum != null ? round1(parseFloat(x.Hum)) : null,
      time: x.RecordTime ? new Date(x.RecordTime) : new Date()
    })
  })
  Object.values(grouped).forEach(rows => rows.sort((a, b) => a.time - b.time))
  labRooms.forEach(room => {
    const rows = grouped[room.deviceName]
    if (!rows?.length) return
    const tempRows = rows.filter(r => r.temp != null)
    const humRows = rows.filter(r => r.hum != null)
    room.temp.history = tempRows.slice(-HISTORY_LEN).map(r => r.temp)
    room.temp.timeHistory = tempRows.slice(-HISTORY_LEN).map(r => r.time)
    room.humidity.history = humRows.slice(-HISTORY_LEN).map(r => r.hum)
    room.humidity.timeHistory = humRows.slice(-HISTORY_LEN).map(r => r.time)
    const latest = rows[rows.length - 1]
    if (latest) {
      if (latest.temp != null) room.temp.value = latest.temp
      if (latest.hum != null) room.humidity.value = latest.hum
      room.updatedAt = latest.time
    }
    room.status = worseStatus(metricStatus(room.temp), metricStatus(room.humidity))
    pushAlert(room, 'temp')
    pushAlert(room, 'humidity')
  })
}
async function fetchDashboardData() {
  if (isLoadingData.value) return
  isLoadingData.value = true
  fetchError.value = ''
  try {
    const isToday = selectedDate.value === todayStr()
    await Promise.all([fetchWarehouseData(isToday), fetchLabData()])
    rebuildOverviewFromHistory()
    await nextTick()
    rooms.forEach(room => {
      drawTrace(room, 'temp')
      drawTrace(room, 'humidity')
      syncCardPulse(room)
    })
  } catch (err) {
    console.error('Lỗi khi lấy dữ liệu:', err)
    fetchError.value = 'Không lấy được dữ liệu'
  } finally {
    isLoadingData.value = false
  }
}
function onDateChange() { fetchDashboardData(); countdown.value = POLL_INTERVAL_SEC }
async function refreshData() { await fetchDashboardData(); countdown.value = POLL_INTERVAL_SEC }

// ==================== THREE.JS BACKGROUND ====================
function initThreeBackground(canvas) {
  const renderer = new THREE.WebGLRenderer({ canvas, alpha: true, antialias: true })
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
  renderer.setSize(window.innerWidth, window.innerHeight)
  const scene = new THREE.Scene()
  const camera = new THREE.PerspectiveCamera(60, window.innerWidth / window.innerHeight, 0.1, 100)
  camera.position.set(0, 1, 18); camera.lookAt(0, 0, 0)
  const COUNT = Math.min(500, Math.floor(window.innerWidth * 0.4))
  const positions = new Float32Array(COUNT * 3)
  const connections = []
  for (let i = 0; i < COUNT; i++) {
    positions[i * 3] = (Math.random() - 0.5) * 40
    positions[i * 3 + 1] = (Math.random() - 0.5) * 25
    positions[i * 3 + 2] = (Math.random() - 0.5) * 20
  }
  const geo = new THREE.BufferGeometry()
  geo.setAttribute('position', new THREE.BufferAttribute(positions, 3))
  const mat = new THREE.PointsMaterial({
    color: 0x4ee4d4, size: 0.06, transparent: true, opacity: 0.5,
    sizeAttenuation: true, blending: THREE.AdditiveBlending, depthWrite: false
  })
  const points = new THREE.Points(geo, mat); scene.add(points)
  const lineGeo = new THREE.BufferGeometry()
  const linePositions = []
  const MAX_LINES = Math.min(60, COUNT)
  for (let i = 0; i < MAX_LINES; i++) {
    const a = Math.floor(Math.random() * COUNT)
    let b = Math.floor(Math.random() * COUNT)
    while (b === a) b = Math.floor(Math.random() * COUNT)
    linePositions.push(
      positions[a * 3], positions[a * 3 + 1], positions[a * 3 + 2],
      positions[b * 3], positions[b * 3 + 1], positions[b * 3 + 2]
    )
    connections.push({ a, b })
  }
  lineGeo.setAttribute('position', new THREE.BufferAttribute(new Float32Array(linePositions), 3))
  const lineMat = new THREE.LineBasicMaterial({
    color: 0x4ee4d4, transparent: true, opacity: 0.08, blending: THREE.AdditiveBlending, depthWrite: false
  })
  const lines = new THREE.LineSegments(lineGeo, lineMat); scene.add(lines)
  let raf; const clock = new THREE.Clock()
  function animate() {
    raf = requestAnimationFrame(animate)
    const dt = Math.min(clock.getDelta(), 0.1)
    const pos = geo.attributes.position.array
    for (let i = 0; i < COUNT; i++) {
      pos[i * 3 + 1] += Math.sin(i * 0.05 + Date.now() * 0.001) * 0.015
      pos[i * 3] += Math.cos(i * 0.03 + Date.now() * 0.0008) * 0.01
      if (pos[i * 3 + 1] > 12) pos[i * 3 + 1] = -12
      if (pos[i * 3 + 1] < -12) pos[i * 3 + 1] = 12
    }
    geo.attributes.position.needsUpdate = true
    const lp = lines.geometry.attributes.position.array
    for (let i = 0; i < MAX_LINES; i++) {
      const { a, b } = connections[i]
      lp[i * 6] = pos[a * 3]; lp[i * 6 + 1] = pos[a * 3 + 1]; lp[i * 6 + 2] = pos[a * 3 + 2]
      lp[i * 6 + 3] = pos[b * 3]; lp[i * 6 + 4] = pos[b * 3 + 1]; lp[i * 6 + 5] = pos[b * 3 + 2]
    }
    lines.geometry.attributes.position.needsUpdate = true
    points.rotation.y += dt * 0.03; lines.rotation.y = points.rotation.y
    renderer.render(scene, camera)
  }
  if (!prefersReducedMotion) animate(); else renderer.render(scene, camera)
  function onResize() {
    camera.aspect = window.innerWidth / window.innerHeight
    camera.updateProjectionMatrix()
    renderer.setSize(window.innerWidth, window.innerHeight)
  }
  window.addEventListener('resize', onResize)
  return () => {
    cancelAnimationFrame(raf); window.removeEventListener('resize', onResize)
    geo.dispose(); mat.dispose(); lineGeo.dispose(); lineMat.dispose(); renderer.dispose()
  }
}

// ==================== FULLSCREEN & NAVIGATION ====================
function toggleFullscreen() {
  if (!document.fullscreenElement) document.documentElement.requestFullscreen().catch(() => {})
  else document.exitFullscreen()
}
function goBack() { window.location.href = 'http://192.168.71.79:8084/' }

// ==================== LIFECYCLE ====================
let cleanupThree = () => {}
let pollInterval = null, clockInterval = null, countdownInterval = null, uptimeInterval = null

onBeforeUnmount(() => {
  if (pollInterval) clearInterval(pollInterval)
  if (clockInterval) clearInterval(clockInterval)
  if (countdownInterval) clearInterval(countdownInterval)
  if (uptimeInterval) clearInterval(uptimeInterval)
  cleanupThree()
  Object.values(pulseTweens).forEach(t => t.kill())
})

onMounted(async () => {
  clockInterval = setInterval(() => { now.value = new Date() }, 1000)
  uptimeInterval = setInterval(() => {
    const diff = Math.floor((Date.now() - startTime) / 1000)
    const h = String(Math.floor(diff / 3600)).padStart(2, '0')
    const m = String(Math.floor((diff % 3600) / 60)).padStart(2, '0')
    const s = String(diff % 60).padStart(2, '0')
    uptime.value = `${h}:${m}:${s}`
  }, 1000)
  countdownInterval = setInterval(() => {
    countdown.value = countdown.value > 1 ? countdown.value - 1 : POLL_INTERVAL_SEC
  }, 1000)

  await fetchDashboardData()

  const tl = gsap.timeline({ defaults: { ease: 'power3.out' } })
  tl.from('.dash-header', { y: -30, opacity: 0, duration: 0.7 })
    .from(logoEl.value, { scale: 0.3, opacity: 0, rotate: -15, duration: 0.6, ease: 'back.out(1.7)' }, '-=0.5')
    .from('.stat-tile', { y: 20, opacity: 0, duration: 0.5, stagger: 0.08 }, '-=0.3')
    .from('.overview-chart', { y: 20, opacity: 0, duration: 0.6 }, '-=0.2')
    .from('.room-card', { scale: 0.95, y: 30, opacity: 0, duration: 0.7, stagger: 0.12 }, '-=0.3')

  if (!prefersReducedMotion) {
    const statsEls = [
      { el: document.querySelector('.stat-tile:nth-child(1) .stat-value'), to: rooms.length },
      { el: document.querySelector('.stat-tile:nth-child(2) .stat-value'), to: totalSensors }
    ]
    statsEls.forEach(({ el, to }) => {
      if (!el) return
      const state = { v: 0 }
      gsap.to(state, { v: to, duration: 1.0, ease: 'power2.out', onUpdate() { el.textContent = Math.round(state.v) } })
    })
    gsap.to('.power-lamp .lamp-dot', { scale: 1.5, opacity: 0.3, duration: 1.2, repeat: -1, yoyo: true, ease: 'sine.inOut' })
    if (logoEl.value) gsap.to(logoEl.value, { scale: 1.05, duration: 2, repeat: -1, yoyo: true, ease: 'sine.inOut' })
    if (tickerTrack.value) {
      gsap.to(tickerTrack.value, {
        x: '-=500', duration: 18, repeat: -1, ease: 'linear',
        modifiers: { x: gsap.utils.unitize(x => parseFloat(x) % (tickerTrack.value.scrollWidth / 2)) }
      })
    }
  }

  if (bgCanvas.value) cleanupThree = initThreeBackground(bgCanvas.value)

  pollInterval = setInterval(async () => {
    await fetchDashboardData()
    countdown.value = POLL_INTERVAL_SEC
  }, POLL_INTERVAL_SEC * 1000)
})
</script>

<style scoped>
/* ==================== CSS CUSTOM PROPERTIES ==================== */
.dash-root {
  --bg: #060a10;
  --bg-glass: rgba(15, 22, 33, 0.7);
  --bg-glass-hover: rgba(20, 30, 45, 0.8);
  --border-glass: rgba(255, 255, 255, 0.08);
  --accent: #4ee4d4;
  --accent-dim: #1e8a82;
  --ok: #4ee4a0;
  --warn: #e0a940;
  --crit: #f04a5a;
  --text: #e8edf4;
  --text-dim: #7f8ea3;
  --pen-temp: #f0a860;
  --pen-hum: #5b9ed4;
  --mono: "JetBrains Mono", "Cascadia Code", "SFMono-Regular", Consolas, monospace;
  --sans: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;

  position: relative;
  min-height: 100vh;
  background: var(--bg);
  color: var(--text);
  font-family: var(--sans);
  -webkit-font-smoothing: antialiased;
  overflow-x: hidden;
}

/* ==================== OVERLAYS ==================== */
.fx-overlay {
  position: fixed; inset: 0; z-index: 2; pointer-events: none;
  background: repeating-linear-gradient(0deg, transparent, transparent 3px, rgba(0,0,0,0.03) 3px, rgba(0,0,0,0.03) 6px),
              radial-gradient(ellipse at center, transparent 60%, rgba(0,0,0,0.5) 100%);
}
.bg-canvas { position: fixed; inset: 0; z-index: 0; }

/* ==================== GLASS PANELS ==================== */
.glass-panel {
  background: var(--bg-glass); backdrop-filter: blur(16px); -webkit-backdrop-filter: blur(16px);
  border: 1px solid var(--border-glass); border-radius: 12px;
}
.glass-panel-inner {
  background: rgba(20,30,48,0.6); backdrop-filter: blur(8px); -webkit-backdrop-filter: blur(8px);
  border: 1px solid var(--border-glass); border-radius: 8px;
}

/* ==================== MAIN CONTENT ==================== */
.dash-content {
  position: relative; z-index: 3;
  width: 100%; max-width: 100%; margin: 0; padding: 24px clamp(20px,2.5vw,48px) 40px;
  box-sizing: border-box;
}

/* ==================== HEADER ==================== */
.dash-header {
  display: flex; align-items: center; justify-content: space-between; gap: 20px;
  padding: 16px 20px; margin-bottom: 20px; flex-wrap: wrap;
}
.brand { display: flex; align-items: center; gap: 16px; }
.brand-logo { height: 32px; object-fit: contain; flex-shrink: 0; }
.brand-text { display: flex; flex-direction: column; }
.brand-eyebrow { font-family: var(--mono); font-size: 10px; letter-spacing: 3px; color: var(--accent); display: block; margin-bottom: 6px; text-shadow: 0 0 8px var(--accent-dim); }
.brand-title { margin: 0; font-size: 28px; font-weight: 800; letter-spacing: 1px; }
.title-neon { color: var(--accent); text-shadow: 0 0 15px var(--accent), 0 0 30px var(--accent-dim); }
.brand-sub { margin: 4px 0 0; color: var(--text-dim); font-size: 12px; }

/* Header right - các control được đồng bộ kích thước */
.header-right { display: flex; align-items: center; flex-shrink: 0; }
.header-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

/* Date picker */
.date-picker {
  display: flex;
  align-items: center;
  padding: 0 10px;
  height: 36px;
}
.date-picker input[type="date"] {
  background: transparent;
  border: none;
  color: var(--text);
  font-family: var(--mono);
  font-size: 12px;
  outline: none;
  height: 100%;
}

/* Nút icon đồng bộ */
.btn-icon {
  width: 36px;
  height: 36px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  line-height: 1;
  flex-shrink: 0;
}

/* Đồng hồ */
.clock-block {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0 10px;
  height: 36px;
  line-height: 1.2;
}
.clock-main {
  font-family: var(--mono);
  font-size: 18px;
  font-weight: 700;
  color: var(--accent);
}
.clock-date {
  font-family: var(--mono);
  font-size: 9px;
  color: var(--text-dim);
  text-transform: uppercase;
}

/* Đèn trạng thái */
.power-lamp {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: var(--mono);
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 1px;
  color: var(--ok);
  text-transform: uppercase;
  white-space: nowrap;
}
.lamp-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--ok);
  box-shadow: 0 0 10px rgba(78,228,160,0.7);
}
.power-lamp--error { color: var(--crit); }
.power-lamp--error .lamp-dot {
  background: var(--crit);
  box-shadow: 0 0 10px rgba(240,74,90,0.8);
}

/* ==================== BUTTON GLASS ==================== */
.btn-glass {
  background: var(--bg-glass);
  border: 1px solid var(--border-glass);
  color: var(--text-dim);
  padding: 6px 10px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.2s;
}
.btn-glass:hover {
  color: var(--accent);
  border-color: var(--accent);
  box-shadow: 0 0 12px var(--accent-dim);
}

/* ==================== ALERT TICKER ==================== */
.alert-ticker { display: flex; align-items: center; margin-bottom: 16px; padding: 0 14px; height: 34px; overflow: hidden; border-left: 3px solid var(--warn); }
.ticker-label { font-family: var(--mono); font-size: 10px; letter-spacing: 2px; color: var(--warn); margin-right: 16px; font-weight: 700; white-space: nowrap; }
.ticker-wrap { flex: 1; overflow: hidden; }
.ticker-track { display: flex; gap: 40px; white-space: nowrap; }
.ticker-item { font-family: var(--mono); font-size: 10px; color: var(--text-dim); }
.ticker-item.warning { color: var(--warn); }
.ticker-item.critical { color: var(--crit); font-weight: 600; }

/* ==================== STATS ROW ==================== */
.stats-row { display: grid; grid-template-columns: repeat(4,1fr); gap: 14px; margin-bottom: 18px; }
.stat-tile { padding: 14px 16px; display: flex; flex-direction: column; gap: 4px; position: relative; overflow: hidden; }
.stat-tile::after { content: ""; position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, transparent, var(--accent-dim), transparent); opacity: 0.5; }
.stat-icon { font-size: 22px; }
.stat-value { font-family: var(--mono); font-size: 30px; font-weight: 700; line-height: 1; }
.stat-value.tone-ok { color: var(--ok); }
.stat-value.tone-warn { color: var(--warn); }
.stat-value.tone-crit { color: var(--crit); }
.stat-value.mono-sm { font-size: 22px; color: var(--accent); }
.stat-label { font-size: 10px; letter-spacing: 2px; text-transform: uppercase; color: var(--text-dim); }

/* ==================== OVERVIEW CHART ==================== */
.overview-chart { padding: 14px 18px; margin-bottom: 18px; }
.chart-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.chart-header h3 { margin: 0; font-size: 14px; font-weight: 600; }
.overview-svg { width: 100%; height: 110px; background: rgba(0,0,0,0.2); border-radius: 6px; border: 1px solid var(--border-glass); }
.grid-lines line { stroke: rgba(255,255,255,0.05); stroke-width: 0.5; }
.recorder-svg .grid-lines line { stroke: rgba(255,255,255,0.1); stroke-dasharray: 2 3; }
.overview-area { fill: url(#areaGrad); }
.overview-line { fill: none; stroke: var(--accent); stroke-width: 2.5; stroke-linejoin: round; filter: drop-shadow(0 0 6px var(--accent-dim)); }
.overview-dot { fill: var(--bg); stroke: var(--accent); stroke-width: 1.5; cursor: pointer; transition: r 0.15s ease; }
.overview-dot:hover { r: 5; fill: var(--accent); }

/* ==================== NHÓM KHU VỰC ==================== */
.group-header-row { display: flex; align-items: center; margin: 6px 2px 2px; }
.group-title { display: flex; align-items: center; gap: 8px; margin: 0; font-size: 15px; font-weight: 700; letter-spacing: 0.5px; color: var(--text); text-transform: uppercase; }
.group-title-icon { font-size: 16px; }
.group-title::after { content: ""; flex: 1; height: 1px; margin-left: 10px; background: var(--border-glass); }

/* ==================== ROOM CARDS GRID ==================== */
.rooms-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(480px,1fr)); gap: 18px; margin-bottom: 18px; }
.room-card { position: relative; will-change: transform, box-shadow; border-left: 4px solid transparent; transition: border-left-color 0.5s ease; perspective: 800px; padding: 0; }
.room-card.status-normal { border-left-color: var(--ok); }
.room-card.status-warning { border-left-color: var(--warn); }
.room-card.status-critical { border-left-color: var(--crit); }
.card-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; padding: 16px 18px 12px; border-bottom: 1px solid var(--border-glass); }
.nameplate { display: flex; align-items: flex-start; gap: 10px; }
.room-icon { width: 26px; height: 26px; color: var(--accent); flex: none; margin-top: 2px; }
.room-icon :deep(svg) { width: 100%; height: 100%; }
.nameplate h3 { margin: 0; font-size: 15px; font-weight: 700; }
.room-code { display: block; margin-top: 3px; font-family: var(--mono); font-size: 10px; letter-spacing: 0.5px; color: var(--text-dim); }
.status-chip { font-family: var(--mono); font-size: 9px; letter-spacing: 1.5px; padding: 4px 10px; border-radius: 20px; text-transform: uppercase; font-weight: 700; }
.chip-normal { background: rgba(78,228,160,0.15); color: var(--ok); border: 1px solid rgba(78,228,160,0.3); }
.chip-warning { background: rgba(224,169,64,0.15); color: var(--warn); border: 1px solid rgba(224,169,64,0.3); }
.chip-critical { background: rgba(240,74,90,0.15); color: var(--crit); border: 1px solid rgba(240,74,90,0.3); animation: chipPulse 1.5s infinite; }
@keyframes chipPulse { 0%,100% { box-shadow: 0 0 8px rgba(240,74,90,0.4); } 50% { box-shadow: 0 0 16px rgba(240,74,90,0.8); } }

/* ==================== GAUGES ==================== */
.gauges-row { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; padding: 8px 10px 0; }
.gauge-wrapper { display: flex; justify-content: center; align-items: center; }
.gauge-chart { width: 100%; max-width: 220px; height: 190px; }

/* ==================== RECORDER ==================== */
.recorder { padding: 8px 18px 10px; }
.recorder-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; font-family: var(--mono); font-size: 9px; color: var(--text-dim); letter-spacing: 0.5px; }
.recorder-time { color: var(--accent); }
.recorder-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.recorder-panel { display: flex; flex-direction: column; gap: 5px; }
.recorder-panel__label { display: flex; align-items: center; gap: 5px; font-size: 10px; color: var(--text-dim); }
.recorder-axis { display: flex; justify-content: space-between; font-family: var(--mono); font-size: 8px; color: var(--text-dim); opacity: 0.7; padding: 0 2px; }
.recorder-svg { width: 100%; height: auto; background: rgba(0,0,0,0.22); border: 1px solid var(--border-glass); border-radius: 10px; }
.pen-trace { fill: none; stroke-width: 1.8; stroke-linejoin: round; stroke-linecap: round; }
.pen-area { stroke: none; }
.pen-temp { stroke: var(--pen-temp); filter: drop-shadow(0 0 3px rgba(240,168,96,0.35)); }
.pen-hum { stroke: var(--pen-hum); filter: drop-shadow(0 0 3px rgba(91,158,212,0.35)); }
.spark-dot { cursor: pointer; opacity: 0.85; transition: r 0.15s ease, opacity 0.2s ease; }
.spark-dot--temp { fill: var(--pen-temp); }
.spark-dot--hum { fill: var(--pen-hum); }
.spark-dot:hover { r: 4; opacity: 1; }
.dot-label { font-family: var(--mono); font-size: 9px; font-weight: 700; pointer-events: none; opacity: 0.95; }
.dot-label tspan:first-child { font-weight: 800; }
.dot-label tspan:last-child { fill: var(--text-dim); font-size: 7px; font-weight: 500; }
.dot-label--temp { fill: var(--pen-temp); }
.dot-label--hum { fill: var(--pen-hum); }
.dot { width: 7px; height: 7px; border-radius: 50%; display: inline-block; }
.dot.temp { background: var(--pen-temp); }
.dot.hum { background: var(--pen-hum); }

/* ==================== CARD FOOTER ==================== */
.card-footer { padding: 10px 18px 14px; border-top: 1px solid var(--border-glass); font-size: 10px; color: var(--text-dim); }
.card-footer b { color: var(--text); font-weight: 600; }

/* ==================== FOOTER ==================== */
.dash-footer { display: flex; justify-content: space-between; align-items: center; margin-top: 20px; padding: 10px 18px; font-family: var(--mono); font-size: 10px; color: var(--text-dim); }
.footer-right { color: var(--accent); }

/* ==================== RESPONSIVE ==================== */
@media (max-width: 1100px) {
  .rooms-grid { grid-template-columns: 1fr; }
}
@media (max-width: 700px) {
  .stats-row { grid-template-columns: repeat(2,1fr); }
  .dash-header { flex-direction: column; align-items: stretch; }
  .header-controls { justify-content: flex-start; }
  .gauges-row { grid-template-columns: 1fr; }
  .recorder-grid { grid-template-columns: 1fr; }
}
@media (max-width: 450px) {
  .stats-row { grid-template-columns: 1fr; }
  .header-controls { gap: 6px; }
  .btn-icon { width: 32px; height: 32px; }
  .clock-block { height: 32px; }
  .date-picker { height: 32px; }
}
@media (prefers-reduced-motion: reduce) {
  .room-card, .lamp-dot { transition: none; animation: none; }
}
</style>