<script setup>
import { useGoreTexI18n } from '../utils/i18n'

const { t } = useGoreTexI18n()

defineProps({
  editing: { type: Boolean, default: false },
  valid: { type: Boolean, default: true },
  submitting: { type: Boolean, default: false },
  historyCount: { type: Number, default: 0 },
  viewingRevision: { type: Boolean, default: false },
})

defineEmits(['edit', 'cancel', 'submit', 'history', 'current', 'back'])
</script>

<template>
  <div class="review-mode-bar" :class="{ editing, revision: viewingRevision }">
    <div class="review-mode-copy">
      <span class="review-mode-dot"></span>
      <div>
        <strong>
          {{
            viewingRevision
              ? t('review.revisionMode')
              : editing
                ? t('review.editMode')
                : t('review.reviewMode')
          }}
        </strong>
        <small>
          {{
            viewingRevision
              ? t('review.revisionHint')
              : editing
              ? valid
                ? t('review.editHint')
                : t('review.requiredHint')
              : t('review.lockedHint')
          }}
        </small>
      </div>
    </div>
    <div class="review-mode-actions">
      <button
        type="button"
        class="review-mode-button secondary"
        :disabled="submitting"
        @click="$emit('history')"
      >
        <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
          <path d="M10 6v4l2.5 1.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
          <path d="M4.3 5.5A7 7 0 1 1 3 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
          <path d="M3 4v3.5h3.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        {{ t('common.history') }}
        <span class="review-history-count">{{ historyCount }}</span>
      </button>
      <button
        v-if="viewingRevision"
        type="button"
        class="review-mode-button"
        @click="$emit('current')"
      >
        <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
          <path d="m7 4 6 6-6 6" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        {{ t('review.current') }}
      </button>
      <button
        v-if="editing && !viewingRevision"
        type="button"
        class="review-mode-button secondary"
        :disabled="submitting"
        @click="$emit('cancel')"
      >
        <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
          <path d="m6 6 8 8M14 6l-8 8" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
        </svg>
        {{ t('review.cancelEdit') }}
      </button>
      <button
        v-if="editing && !viewingRevision"
        type="button"
        class="review-mode-button"
        :disabled="submitting"
        @click="$emit('submit')"
      >
        <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
          <path d="M4 10h11M11 6l4 4-4 4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        {{ submitting ? t('review.sending') : t('common.submit') }}
      </button>
      <button
        v-if="!editing && !viewingRevision"
        type="button"
        class="review-mode-button"
        @click="$emit('edit')"
      >
        <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
          <path d="m4 14.5.5-3L12.8 3.2a1.5 1.5 0 0 1 2.1 0l1.9 1.9a1.5 1.5 0 0 1 0 2.1l-8.3 8.3-3 .5-1.5-1.5Z" stroke="currentColor" stroke-width="1.4" stroke-linejoin="round" />
          <path d="m11.5 4.5 4 4" stroke="currentColor" stroke-width="1.4" />
        </svg>
        {{ t('common.edit') }}
      </button>
      <button
        v-if="!editing && !viewingRevision"
        type="button"
        class="review-mode-button secondary"
        @click="$emit('back')"
      >
        <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
          <path d="M16 10H5M9 6l-4 4 4 4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        {{ t('common.back') }}
      </button>
    </div>
  </div>
</template>
