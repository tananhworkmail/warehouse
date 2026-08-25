<script setup>
import { computed } from 'vue'

const props = defineProps({
  items: { type: Array, default: () => [] },
  emptyText: { type: String, default: 'No data' },
  defectsText: { type: String, default: 'defects' },
  percentageText: { type: String, default: 'percentage' },
})
const width = 800
const height = 410
const left = 54
const right = 50
const top = 28
const bottom = 150
const plotWidth = width - left - right
const plotHeight = height - top - bottom
const maxCount = computed(() => Math.max(1, ...props.items.map((item) => item.count)))
const bandWidth = computed(() => plotWidth / Math.max(1, props.items.length))
const barWidth = computed(() => Math.max(12, bandWidth.value * .62))
const x = (index) => left + index * bandWidth.value + bandWidth.value / 2
const yCount = (value) => top + (maxCount.value - value) * plotHeight / maxCount.value
const yRate = (value) => top + (100 - value) * plotHeight / 100
const points = computed(() => props.items.map((item, index) => `${x(index)},${yRate(item.percentage)}`).join(' '))
const hasDefects = computed(() => props.items.some((item) => Number(item.count) > 0))
</script>

<template>
  <div v-if="!items.length" class="dashboard-chart-empty">{{ emptyText }}</div>
  <div v-else class="dashboard-svg-scroll">
    <svg :viewBox="`0 0 ${width} ${height}`" role="img" aria-label="Biểu đồ Pareto lỗi chống thấm">
      <g v-for="tick in [0, 25, 50, 75, 100]" :key="tick">
        <line :x1="left" :y1="yRate(tick)" :x2="width - right" :y2="yRate(tick)" class="chart-grid-line" />
        <text :x="width - right + 8" :y="yRate(tick) + 4" class="chart-axis-label">{{ tick }}%</text>
      </g>
      <g v-for="(item, index) in items" :key="item.label">
        <rect :x="x(index) - barWidth / 2" :y="yCount(item.count)" :width="barWidth" :height="height - bottom - yCount(item.count)" rx="3" class="chart-pareto-bar">
          <title>{{ item.label }}: {{ item.count }} {{ defectsText }} · {{ percentageText }} {{ item.percentage }}%</title>
        </rect>
        <text :x="x(index)" :y="yCount(item.count) - 7" class="chart-value-label" text-anchor="middle">{{ item.count }}</text>
        <text
          :x="x(index) - 3"
          :y="height - bottom + 19"
          class="chart-pareto-label"
          text-anchor="end"
          :transform="`rotate(-45 ${x(index) - 3} ${height - bottom + 19})`"
        >{{ item.label }}</text>
      </g>
      <polyline :points="points" class="chart-pareto-line" :class="{ 'is-zero': !hasDefects }" />
      <g v-for="(item, index) in items" :key="`rate-${item.label}`">
        <circle v-if="Number(item.count) > 0" :cx="x(index)" :cy="yRate(item.percentage)" r="4.5" class="chart-pareto-point" />
        <text v-if="Number(item.count) > 0" :x="x(index)" :y="yRate(item.percentage) - 9" class="chart-pareto-rate" text-anchor="middle">{{ item.percentage }}%</text>
      </g>
    </svg>
  </div>
</template>
