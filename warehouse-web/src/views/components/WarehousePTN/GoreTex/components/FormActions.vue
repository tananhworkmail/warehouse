<script setup>
import { useGoreTexI18n } from '../utils/i18n'

const { t } = useGoreTexI18n()

defineProps({
  savedLabel: { type: String, default: '' },
  valid: { type: Boolean, default: true },
  submitting: { type: Boolean, default: false },
})

defineEmits(['cancel', 'save', 'submit', 'reset'])
</script>

<template>
  <div class="form-actions">
    <div class="action-status">
      <div class="save-indicator">
        <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
          <path d="M10 17a7 7 0 1 0 0-14 7 7 0 0 0 0 14Z" stroke="currentColor" stroke-width="1.5" />
          <path d="M10 6v4l2.5 1.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
        </svg>
        {{ savedLabel || t('actions.noDraft') }}
      </div>
      <span v-if="!valid" class="required-hint">
        {{ t('actions.required') }}
      </span>
    </div>
    <div class="action-buttons">
      <button type="button" class="btn btn-secondary" :disabled="submitting" @click="$emit('cancel')">
       ❌ {{ t('common.cancel') }}
      </button>
      <button type="button" class="btn btn-ghost" :disabled="submitting" @click="$emit('reset')">{{ t('common.reset') }}</button>
      <button type="button" class="btn btn-secondary" :disabled="submitting" @click="$emit('save')">{{ t('common.draft') }}</button>
      <button type="button" class="btn btn-primary" :disabled="submitting" @click="$emit('submit')">
        {{ submitting ? t('common.submitting') : t('common.submit') }}
        <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
          <path d="m7 4 6 6-6 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </button>
    </div>
  </div>
</template>
