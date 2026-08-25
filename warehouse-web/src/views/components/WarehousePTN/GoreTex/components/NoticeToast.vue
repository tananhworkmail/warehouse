<script setup>
defineProps({
  show: Boolean,
  message: String,
  tone: { type: String, default: 'success' },
})
</script>

<template>
  <Transition name="toast">
    <div v-if="show" class="notice-toast" :class="`notice-${tone}`" role="status">
      <span class="notice-icon">
        <svg v-if="tone === 'success'" viewBox="0 0 20 20" fill="none">
          <path d="m5 10 3.2 3L15 6.5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        <svg v-else viewBox="0 0 20 20" fill="none">
          <path d="M10 6v5M10 14h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
        </svg>
      </span>
      {{ message }}
    </div>
  </Transition>
</template>
