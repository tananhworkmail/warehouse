<template>
  <div class="date-range-control" role="group" :aria-label="`${fromLabel} - ${toLabel}`">
    <label class="date-input">
      <span>{{ fromLabel }}</span>
      <input
        :id="`tracking-date-from-${index}`"
        v-model="fromDate"
        type="date"
        autocomplete="off"
      />
    </label>

    <svg class="date-range-arrow" viewBox="0 0 20 20" aria-hidden="true">
      <path d="M4 10h12m-4-4 4 4-4 4" />
    </svg>

    <label class="date-input">
      <span>{{ toLabel }}</span>
      <input
        :id="`tracking-date-to-${index}`"
        v-model="toDate"
        type="date"
        autocomplete="off"
      />
    </label>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  index: {
    type: Number,
    required: true,
  },
  fromLabel: {
    type: String,
    default: 'Từ ngày',
  },
  toLabel: {
    type: String,
    default: 'Đến ngày',
  },
})

const emit = defineEmits(['Set_Date_Range'])

const pad2 = (value) => String(value).padStart(2, '0')
const getLocalDate = () => {
  const date = new Date()
  return `${date.getFullYear()}-${pad2(date.getMonth() + 1)}-${pad2(date.getDate())}`
}

const fromDate = ref(getLocalDate())
const toDate = ref(getLocalDate())

watch(
  [fromDate, toDate],
  ([from, to], previous = []) => {
    if (!from || !to) return

    if (from > to) {
      if (from !== previous[0]) toDate.value = from
      else fromDate.value = to
      return
    }

    emit('Set_Date_Range', `${from}~${to}`, props.index)
  },
  { immediate: true },
)
</script>

<style scoped>
.date-range-control {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 20px minmax(0, 1fr);
  align-items: end;
  gap: 7px;
  min-width: 0;
}

.date-input {
  display: grid;
  min-width: 0;
  gap: 4px;
}

.date-input span {
  color: #5f7087;
  font-size: 10px;
  font-weight: 700;
  line-height: 1;
}

.date-input input {
  width: 100%;
  min-width: 0;
  height: 38px;
  padding: 0 9px;
  border: 1px solid #cbd7e6;
  border-radius: 9px;
  outline: none;
  color: #172947;
  background: #fff;
  font: inherit;
  font-size: 12px;
  color-scheme: light;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.date-input input:hover {
  border-color: #9fb3cc;
}

.date-input input:focus-visible {
  border-color: #2477cf;
  box-shadow: 0 0 0 3px rgba(36, 119, 207, 0.16);
}

.date-range-arrow {
  width: 18px;
  height: 38px;
  color: #6b809a;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 1.8;
}

@media (max-width: 430px) {
  .date-range-control {
    grid-template-columns: minmax(0, 1fr);
  }

  .date-range-arrow {
    display: none;
  }
}
</style>
