<script setup>
import { computed } from 'vue'

const props = defineProps({
  items: { type: Array, default: () => [] },
  emptyText: { type: String, default: 'No data' },
  passedText: { type: String, default: 'passed' },
})
const width = 720
const left = 190
const right = 56
const rowHeight = 38
const height = computed(() => Math.max(250, 70 + props.items.length * rowHeight))
const plotWidth = width - left - right
const x = (value) => left + (Math.max(0, Math.min(100, value)) / 100) * plotWidth
</script>

<template>
  <div v-if="!items.length" class="dashboard-chart-empty">{{ emptyText }}</div>
  <div v-else class="dashboard-svg-scroll">
    <svg :viewBox="`0 0 ${width} ${height}`" role="img" aria-label="Tỷ lệ đạt theo dạng giày">
      <g v-for="tick in [0, 25, 50, 75, 100]" :key="tick">
        <line :x1="x(tick)" y1="22" :x2="x(tick)" :y2="height - 34" class="chart-grid-line" />
        <text :x="x(tick)" :y="height - 13" class="chart-axis-label" text-anchor="middle">{{ tick }}%</text>
      </g>
      <g v-for="(item, index) in items" :key="item.label" :transform="`translate(0 ${36 + index * rowHeight})`">
        <text :x="left - 12" y="16" class="chart-item-label" text-anchor="end">{{ item.label }}</text>
        <rect :x="left" y="0" :width="plotWidth" height="24" rx="4" class="chart-bar-track" />
        <rect :x="left" y="0" :width="Math.max(2, x(item.rate) - left)" height="24" rx="4" class="chart-rate-bar">
          <title>{{ item.label }}: {{ item.pass }}/{{ item.total }} {{ passedText }} ({{ item.rate }}%)</title>
        </rect>
        <text :x="Math.min(width - 5, x(item.rate) + 8)" y="16" class="chart-value-label">{{ item.rate }}%</text>
      </g>
    </svg>
  </div>
</template>
