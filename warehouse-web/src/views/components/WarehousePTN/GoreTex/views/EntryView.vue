<script setup>
import { useGoreTexI18n } from '../utils/i18n'

const { t } = useGoreTexI18n()
</script>

<template>
  <div class="goretex-entry-page">
    <section class="goretex-entry-heading">
      <span class="eyebrow">{{ t('entry.eyebrow') }}</span>
      <h1>GORE-TEX</h1>
      <p>{{ t('entry.hint') }}</p>
    </section>

    <section class="goretex-entry-grid">
      <RouterLink class="goretex-entry-card" :to="{ name: 'laboratory-forms-portal' }">
        <span class="goretex-entry-icon">
          <svg viewBox="0 0 28 28" fill="none" aria-hidden="true">
            <path d="M6 4.5h11l5 5v14H6v-19Z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round" />
            <path d="M17 4.5v5h5M10 14h8M10 18h8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
          </svg>
        </span>
        <div>
          <small>{{ t('entry.portalTag') }}</small>
          <h2>{{ t('entry.portalTitle') }}</h2>
          <p>{{ t('entry.portalDescription') }}</p>
        </div>
        <span class="goretex-entry-action">
          {{ t('entry.openPortal') }}
          <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
            <path d="M4 10h11M11 6l4 4-4 4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </span>
      </RouterLink>

      <RouterLink class="goretex-entry-card dashboard" :to="{ name: 'laboratory-forms-dashboard' }">
        <span class="goretex-entry-icon">
          <svg viewBox="0 0 28 28" fill="none" aria-hidden="true">
            <path d="M5 23V13m6 10V7m6 16V11m6 12V4" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" />
            <path d="M3.5 23.5h21" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
          </svg>
        </span>
        <div>
          <small>{{ t('entry.dashboardTag') }}</small>
          <h2>Dashboard GORE-TEX</h2>
          <p>{{ t('entry.dashboardDescription') }}</p>
        </div>
        <span class="goretex-entry-action">
          {{ t('entry.openDashboard') }}
          <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
            <path d="M4 10h11M11 6l4 4-4 4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </span>
      </RouterLink>
    </section>
  </div>
</template>
