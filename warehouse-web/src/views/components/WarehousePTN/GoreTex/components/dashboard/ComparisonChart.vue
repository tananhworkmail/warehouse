<script setup>
import { computed } from 'vue'

const props = defineProps({
  series: { type: Array, default: () => [] },
  emptyText: { type: String, default: 'No data' },
  passedText: { type: String, default: 'passed' },
})
const rows = computed(() => props.series.flatMap((group, groupIndex) =>
  (group.items || []).map((item) => ({ ...item, weekLabel: group.label, groupIndex }))))
const width = 760
const left = 235
const right = 56
const rowHeight = 34
const height = computed(() => Math.max(260, 82 + rows.value.length * rowHeight))
const plotWidth = width - left - right
const x = (value) => left + Math.max(0, Math.min(100, value)) * plotWidth / 100
</script>

<template>
  <div v-if="!rows.length" class="dashboard-chart-empty">{{ emptyText }}</div>
  <div v-else class="dashboard-svg-scroll">
    <div class="dashboard-chart-legend">
      <span v-for="(group, index) in series" :key="`${group.year}-${group.week}`">
        <i :class="`legend-swatch comparison-${index % 2}`"></i>{{ group.label }}
      </span>
    </div>
    <svg :viewBox="`0 0 ${width} ${height}`" role="img" aria-label="So sánh tỷ lệ đạt giữa hai tuần">
      <g v-for="tick in [0, 25, 50, 75, 100]" :key="tick">
        <line :x1="x(tick)" y1="16" :x2="x(tick)" :y2="height - 34" class="chart-grid-line" />
        <text :x="x(tick)" :y="height - 12" class="chart-axis-label" text-anchor="middle">{{ tick }}%</text>
      </g>
      <g v-for="(item, index) in rows" :key="`${item.weekLabel}-${item.label}`" :transform="`translate(0 ${28 + index * rowHeight})`">
        <text :x="left - 12" y="15" class="chart-item-label" text-anchor="end">{{ item.label }}</text>
        <rect :x="left" y="0" :width="Math.max(2, x(item.rate) - left)" height="21" rx="3" :class="`chart-comparison-bar comparison-${item.groupIndex % 2}`">
          <title>{{ item.weekLabel }} · {{ item.label }}: {{ item.pass }}/{{ item.total }} {{ passedText }} ({{ item.rate }}%)</title>
        </rect>
        <text :x="Math.min(width - 5, x(item.rate) + 7)" y="15" class="chart-value-label">{{ item.rate }}%</text>
      </g>
    </svg>
  </div>
</template>
