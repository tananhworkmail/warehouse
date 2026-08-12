<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useGoreTexI18n } from '../utils/i18n'

const route = useRoute()
const { t } = useGoreTexI18n()
const forms = [
  {
    type: 'waterproof',
    number: '01',
    title: 'BÁO BIỂU KIỂM TRA CHẤT LƯỢNG GIÀY CHỐNG THẤM NƯỚC',
    chineseTitle: '防水鞋品质检验报表',
    englishTitle: 'Daily Quality Report for Waterproof Shoes',
    code: 'A-QIP-WS003-04D',
    route: '/laboratory/forms/kiem-tra-chong-tham',
  },
  {
    type: 'centrifugal',
    number: '02',
    title: 'BÁO BIỂU GIÀY THÀNH PHẨM THỬ NGHIỆM LI TÂM',
    chineseTitle: '成品鞋离心测试报表',
    englishTitle: 'Finished Shoe Centrifugal Test Report',
    code: 'A-QIP-WS004-04C',
    route: '/laboratory/forms/thu-nghiem-li-tam',
  },
  {
    type: 'analysis',
    number: '03',
    title: 'BIỂU PHÂN TÍCH NGUYÊN NHÂN VÀ CẢI THIỆN THỬ NGHIỆM LI TÂM / THỬ NƯỚC',
    chineseTitle: '袜套测试 / 离心测试的原因分析与改善表',
    englishTitle: 'Cause Analysis and Improvement Form for Centrifugal / Water Test',
    code: 'BIỂU PHÂN TÍCH',
    route: '/laboratory/forms/phan-tich-cai-thien',
  },
]

const submittedMessage = computed(() => {
  if (!route.query.submitted) return ''
  const submittedForm = forms.find((form) => form.type === route.query.submitted)
  return submittedForm
    ? t('portal.submitted', { title: t(`forms.${submittedForm.type}.title`) })
    : t('portal.submittedDefault')
})
</script>

<template>
  <div class="home-page simple-home">
    <section class="simple-hero">
      <div>
        <span class="eyebrow">{{ t('portal.eyebrow') }}</span>
        <h1>{{ t('portal.title') }}</h1>
        <p>{{ t('portal.description') }}</p>
      </div>
      <RouterLink class="btn btn-primary goretex-page-back" :to="{ name: 'laboratory-forms' }">
        <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
          <path d="M16 10H5m4-4-4 4 4 4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        {{ t('common.back') }}
      </RouterLink>
    </section>

    <div v-if="submittedMessage" class="portal-success">
      {{ submittedMessage }}
    </div>

    <section class="simple-card-grid">
      <article v-for="form in forms" :key="form.route" class="simple-card">
        <div class="simple-card-head">
          <div class="simple-card-identity">
            <span class="portal-form-icon">
              <svg v-if="form.type === 'waterproof'" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <path d="M12 3.5S6.5 9.6 6.5 14a5.5 5.5 0 0 0 11 0C17.5 9.6 12 3.5 12 3.5Z" stroke="currentColor" stroke-width="1.8" />
                <path d="M9.5 15.2c.35 1.2 1.2 1.8 2.5 1.8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
              </svg>
              <svg v-else-if="form.type === 'centrifugal'" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <circle cx="12" cy="12" r="8.5" stroke="currentColor" stroke-width="1.8" />
                <path d="M12 7.5a4.5 4.5 0 1 1-3.2 1.3M12 7.5v3.2m0-3.2 2.2 1.3" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <path d="M5 18V9m7 9V5m7 13v-6" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" />
                <path d="M3.5 18.5h17" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
              </svg>
            </span>
            <!-- <span class="form-number">{{ form.number }}</span> -->
          </div>
          <span class="doc-code">{{ form.code }}</span>
        </div>
        <div>
          <h2>{{ t(`forms.${form.type}.title`) }}</h2>
          <div class="card-form-languages">
            <p lang="zh">{{ form.chineseTitle }}</p>
            <p lang="en">{{ form.englishTitle }}</p>
          </div>
        </div>
        <div class="simple-card-foot">
          <RouterLink
            class="card-action card-action-secondary"
            :to="{ path: '/laboratory/forms/history', query: { type: form.type } }"
          >
            <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
              <path d="M4 5.5h12M4 10h12M4 14.5h8" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
            </svg>
            {{ t('portal.list') }}
          </RouterLink>
          <RouterLink class="card-action card-action-primary" :to="form.route">
            <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
              <path d="M5 3.5h7l3 3v10H5v-13Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round" />
              <path d="M8 10h4M8 13h4M12 3.5v3h3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
            </svg>
            {{ t('portal.open') }}
          </RouterLink>
        </div>
      </article>
    </section>
  </div>
</template>
