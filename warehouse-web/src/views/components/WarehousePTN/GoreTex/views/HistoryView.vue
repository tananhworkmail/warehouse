<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getApiErrorMessage, getGoreTexForms } from "../utils/api";
import { useGoreTexI18n } from "../utils/i18n";

const route = useRoute();
const router = useRouter();
const items = ref([]);
const loading = ref(false);
const errorMessage = ref("");
const { t, dateLocale } = useGoreTexI18n();

const filters = computed(() => [
  {
    value: "waterproof",
    label: t("forms.waterproof.title"),
    shortLabel: t("forms.waterproof.short"),
    code: "A-QIP-WS003-04D",
  },
  {
    value: "centrifugal",
    label: t("forms.centrifugal.title"),
    shortLabel: t("forms.centrifugal.short"),
    code: "A-QIP-WS004-04C",
  },
  {
    value: "analysis",
    label: t("forms.analysis.title"),
    shortLabel: t("forms.analysis.short"),
    code: "BIỂU PHÂN TÍCH",
  },
]);
const formTypes = ["waterproof", "centrifugal", "analysis"];
const selectedType = ref(
  formTypes.includes(route.query.type) ? route.query.type : "waterproof",
);

const filteredItems = computed(() =>
  items.value.filter((item) => item.formType === selectedType.value),
);
const showWaterproofColumns = computed(
  () => selectedType.value === "waterproof",
);
const showAnalysisColumns = computed(
  () => selectedType.value === "analysis",
);
const selectedForm = computed(
  () =>
    filters.value.find((filter) => filter.value === selectedType.value) || filters.value[0],
);

function selectType(type) {
  selectedType.value = type;
  router.replace({ query: { ...route.query, type } });
}

watch(
  () => route.query.type,
  (type) => {
    if (formTypes.includes(type)) selectedType.value = type;
  },
);

function reviewRoute(item) {
  if (item.formType === "waterproof") {
    return {
      path: "/laboratory/forms/kiem-tra-chong-tham",
      query: {
        review: "1",
        line: item.line,
        styleName: item.styleName,
      },
    };
  }
  if (item.formType === "centrifugal") {
    return {
      path: "/laboratory/forms/thu-nghiem-li-tam",
      query: { review: "1", inspectionDate: item.inspectionDate },
    };
  }
  return {
    path: "/laboratory/forms/phan-tich-cai-thien",
    query: { review: "1", id: item.analysisId },
  };
}

function formatDateTime(value) {
  if (!value) return "—";
  return new Intl.DateTimeFormat(dateLocale(), {
    hour: "2-digit",
    minute: "2-digit",
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  }).format(new Date(value));
}

function formatDate(value) {
  if (!value) return "—";
  const [year, month, day] = String(value).slice(0, 10).split("-");
  if (!year || !month || !day) return value;
  return new Intl.DateTimeFormat(dateLocale(), {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  }).format(new Date(Number(year), Number(month) - 1, Number(day)));
}

async function loadItems() {
  loading.value = true;
  errorMessage.value = "";
  try {
    items.value = await getGoreTexForms();
  } catch (error) {
    errorMessage.value = getApiErrorMessage(error);
  } finally {
    loading.value = false;
  }
}

function goBack() {
  if (window.history.state?.back) router.back();
  else router.push({ name: "laboratory-forms-portal" });
}

onMounted(loadItems);
</script>

<template>
  <div class="history-page">
    <section class="history-heading">
      <h1>{{ t('history.title') }}</h1>
      <div class="history-heading-actions">
        <button class="btn btn-secondary" type="button" @click="loadItems">
          <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
            <path d="M15.5 7A6 6 0 1 0 16 12" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
            <path d="M12.5 7H16V3.5" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          {{ t('common.refresh') }}
        </button>
        <button class="btn btn-secondary" type="button" @click="goBack">
          <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
            <path d="M16 10H5M9 6l-4 4 4 4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          {{ t('common.back') }}
        </button>
      </div>
    </section>

    <div class="history-workspace">
      <nav class="history-filter-panel" :aria-label="t('history.chooseForm')">
        <div class="history-filters" role="group" :aria-label="t('history.filter')">
          <button
            v-for="filter in filters"
            :key="filter.value"
            type="button"
            :class="{ active: selectedType === filter.value }"
            @click="selectType(filter.value)"
          >
            <span class="history-filter-icon">
              <svg v-if="filter.value === 'waterproof'" viewBox="0 0 20 20" fill="none" aria-hidden="true">
                <path d="M10 2.8S5.2 8 5.2 11.7a4.8 4.8 0 0 0 9.6 0C14.8 8 10 2.8 10 2.8Z" stroke="currentColor" stroke-width="1.5" />
              </svg>
              <svg v-else-if="filter.value === 'centrifugal'" viewBox="0 0 20 20" fill="none" aria-hidden="true">
                <path d="M15.5 7A6 6 0 1 0 16 12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                <path d="M12.5 7H16V3.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
                <circle cx="10" cy="10" r="1.7" stroke="currentColor" stroke-width="1.5" />
              </svg>
              <svg v-else viewBox="0 0 20 20" fill="none" aria-hidden="true">
                <path d="M4 15.5V11m4 4.5V7m4 8.5V9m4 6.5V4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
              </svg>
            </span>
            <span class="history-filter-copy">
              <strong>{{ filter.shortLabel }}</strong>
            </span>
          </button>
        </div>
      </nav>

      <section class="history-results">
        <header class="history-result-heading">
          <h2>{{ selectedForm.label }}</h2>
          <strong class="history-count">
            {{ filteredItems.length }}
            <small>{{ t('history.records') }}</small>
          </strong>
        </header>

        <div v-if="loading" class="history-state">
          <span class="history-loader"></span>
          {{ t('history.loading') }}
        </div>
        <div v-else-if="errorMessage" class="history-state history-error">
          {{ errorMessage }}
        </div>
        <div v-else-if="!filteredItems.length" class="history-state">
          <strong>{{ t('history.emptyTitle') }}</strong>
          <span>{{ t('history.emptyText') }}</span>
        </div>

        <div v-else class="history-table-wrap">
          <table class="history-list">
            <thead>
              <tr>
                <th class="history-index-column">STT</th>
                <th v-if="showWaterproofColumns">{{ t('history.line') }}</th>
                <th v-if="showWaterproofColumns">{{ t('history.style') }}</th>
                <th v-if="showAnalysisColumns">{{ t('history.testDate') }}</th>
                <th v-if="showAnalysisColumns">{{ t('history.improvementDate') }}</th>
                <th v-else>{{ t('history.inspectionDate') }}</th>
                <th>{{ t('history.updated') }}</th>
                <th class="history-action-column"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in filteredItems" :key="`${item.formType}:${item.recordKey}`">
                <td class="history-index">{{ String(index + 1).padStart(2, "0") }}</td>
                <td v-if="showWaterproofColumns">
                  <strong>{{ item.line || "—" }}</strong>
                </td>
                <td v-if="showWaterproofColumns">{{ item.styleName || "—" }}</td>
                <td v-if="showAnalysisColumns">
                  <div class="history-date-list">
                    <span v-for="date in item.testDates || []" :key="date">
                      {{ formatDate(date) }}
                    </span>
                    <span v-if="!item.testDates?.length">—</span>
                  </div>
                </td>
                <td v-if="showAnalysisColumns">
                  <div class="history-date-list">
                    <span v-for="date in item.improvementDates || []" :key="date">
                      {{ formatDate(date) }}
                    </span>
                    <span v-if="!item.improvementDates?.length">—</span>
                  </div>
                </td>
                <td v-else>
                  <span class="history-date">{{ formatDate(item.inspectionDate) }}</span>
                </td>
                <td>
                  <span class="history-updated">{{ formatDateTime(item.updatedAt) }}</span>
                </td>
                <td>
                  <RouterLink class="review-link" :to="reviewRoute(item)">
                    {{ t('history.details') }}
                    <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
                      <path d="M4 10h11M11 6l4 4-4 4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
                    </svg>
                  </RouterLink>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>
  </div>
</template>
