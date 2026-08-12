<script setup>
import { computed } from 'vue'
import { getRevisionChanges } from '../utils/revisions'
import { useGoreTexI18n } from '../utils/i18n'

const { t, dateLocale } = useGoreTexI18n()

const props = defineProps({
  show: { type: Boolean, default: false },
  history: { type: Array, default: () => [] },
  currentData: { type: Object, default: () => ({}) },
})

defineEmits(['close', 'view'])

const revisions = computed(() =>
  props.history
    .map((entry, index) => ({
      ...entry,
      revisionNumber: index + 1,
      changes: getRevisionChanges(
        entry.data,
        props.history[index + 1]?.data || props.currentData,
      ),
      viewData: props.history[index + 1]?.data || props.currentData,
    }))
    .reverse(),
)

function formatDateTime(value) {
  if (!value) return t('revisions.noTime')
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat(dateLocale(), {
    hour: '2-digit',
    minute: '2-digit',
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date)
}
</script>

<template>
  <div v-if="show" class="edit-history-overlay" @click.self="$emit('close')">
    <section class="edit-history-dialog" role="dialog" aria-modal="true" aria-labelledby="edit-history-title">
      <header class="edit-history-header">
        <div>
          <span class="edit-history-heading-icon">
            <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path d="M12 7v5l3 2" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
              <path d="M4.9 6.3A9 9 0 1 1 3 12" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
              <path d="M3 5v4h4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </span>
          <div>
            <h2 id="edit-history-title">{{ t('revisions.title') }}</h2>
            <p>{{ t('revisions.versions', { count: history.length }) }}</p>
          </div>
        </div>
        <button class="edit-history-close" type="button" :aria-label="t('common.close')" @click="$emit('close')">×</button>
      </header>

      <div v-if="!revisions.length" class="edit-history-empty">
        {{ t('revisions.empty') }}
      </div>

      <div v-else class="edit-history-list">
        <article v-for="entry in revisions" :key="`${entry.editedAt}-${entry.revisionNumber}`" class="edit-history-item">
          <span class="edit-history-timeline-dot"></span>
          <div class="edit-history-item-copy">
            <strong>{{ t('revisions.version', { number: entry.revisionNumber }) }}</strong>
            <span>{{ t('revisions.savedBefore', { time: formatDateTime(entry.editedAt) }) }}</span>
            <small class="edit-history-change-count">{{ t('revisions.changes', { count: entry.changes.length }) }}</small>
          </div>
          <button type="button" class="edit-history-view" @click="$emit('view', entry)">
            <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
              <path d="M2.5 10s2.7-4.5 7.5-4.5 7.5 4.5 7.5 4.5-2.7 4.5-7.5 4.5S2.5 10 2.5 10Z" stroke="currentColor" stroke-width="1.5" />
              <circle cx="10" cy="10" r="2" stroke="currentColor" stroke-width="1.5" />
            </svg>
            {{ t('revisions.viewVersion') }}
          </button>
        </article>
      </div>
    </section>
  </div>
</template>
