<script setup>
import { computed } from 'vue'

const props = defineProps({
  items: { type: Array, default: () => [] },
  emptyText: { type: String, default: 'No data' },
  targetText: { type: String, default: 'Target 99%' },
  passedText: { type: String, default: 'passed' },
})
const width = 700
const height = 300
const left = 48
const right = 26
const top = 30
const bottom = 46
const plotWidth = width - left - right
const plotHeight = height - top - bottom
const x = (index) => left + (props.items.length <= 1 ? plotWidth / 2 : index * plotWidth / (props.items.length - 1))
const y = (value) => top + (100 - Math.max(0, Math.min(100, value))) * plotHeight / 100
const points = computed(() => props.items.map((item, index) => `${x(index)},${y(item.rate)}`).join(' '))
</script>

<template>
  <div v-if="!items.length" class="dashboard-chart-empty">{{ emptyText }}</div>
  <div v-else class="dashboard-svg-scroll">
    <svg :viewBox="`0 0 ${width} ${height}`" role="img" aria-label="Xu hướng tỷ lệ đạt theo ngày">
      <g v-for="tick in [0, 25, 50, 75, 100]" :key="tick">
        <line :x1="left" :y1="y(tick)" :x2="width - right" :y2="y(tick)" class="chart-grid-line" />
        <text :x="left - 9" :y="y(tick) + 4" class="chart-axis-label" text-anchor="end">{{ tick }}%</text>
      </g>
      <line :x1="left" :y1="y(items[0]?.target ?? 99)" :x2="width - right" :y2="y(items[0]?.target ?? 99)" class="chart-target-line" />
      <text :x="width - right" :y="y(items[0]?.target ?? 99) - 7" class="chart-target-label" text-anchor="end">{{ targetText }}</text>
      <polyline :points="points" class="chart-trend-line" />
      <g v-for="(item, index) in items" :key="item.date">
        <circle :cx="x(index)" :cy="y(item.rate)" r="5" class="chart-trend-point">
          <title>{{ item.date }}: {{ item.pass }}/{{ item.total }} {{ passedText }} ({{ item.rate }}%)</title>
        </circle>
        <text :x="x(index)" :y="y(item.rate) - 12" class="chart-value-label" text-anchor="middle">{{ item.rate }}%</text>
        <text :x="x(index)" :y="height - 17" class="chart-axis-label" text-anchor="middle">{{ item.label }}</text>
      </g>
    </svg>
  </div>
</template>
