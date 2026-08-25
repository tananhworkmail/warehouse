// src/hooks/i18n.js
// Shared translation instance used by the app and by components that expect a
// named locale export from this module.
import { createI18n, useI18n as vueUseI18n } from "vue-i18n";
import { ref, watch } from "vue";
import vi from "@/locales/vi.js";
import en from "@/locales/en.js";
import zh from "@/locales/zh.js";

const SUPPORTED_LOCALES = ["vi", "en", "zh"];
const savedLocale =
  typeof window !== "undefined" ? localStorage.getItem("APP_LOCALE") : null;
const defaultLocale = SUPPORTED_LOCALES.includes(savedLocale)
  ? savedLocale
  : "vi";

// The N8d pages were written before their messages were nested under
// `warehouse.n8d`. Expose compatible aliases while keeping the canonical
// namespace used by newer pages.
const withN8dAliases = (messages) => ({
  ...messages,
  common: {
    ...messages.common,
    ...messages.warehouse?.n8d?.common,
  },
  trackingPage: messages.warehouse?.n8d?.trackingPage || {},
  inventoryPage: messages.warehouse?.n8d?.inventoryPage || {},
});

const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: defaultLocale,
  fallbackLocale: "vi",
  messages: {
    vi: withN8dAliases(vi),
    en: withN8dAliases(en),
    zh: withN8dAliases(zh),
  },
});

export const locale = ref(defaultLocale);

watch(
  locale,
  (nextLocale) => {
    if (!SUPPORTED_LOCALES.includes(nextLocale)) return;
    i18n.global.locale.value = nextLocale;
    if (typeof window !== "undefined")
      localStorage.setItem("APP_LOCALE", nextLocale);
  },
  { immediate: true },
);

watch(
  () => i18n.global.locale.value,
  (nextLocale) => {
    if (nextLocale && nextLocale !== locale.value) locale.value = nextLocale;
  },
  { immediate: true },
);

export function useI18n(options) {
  return vueUseI18n(options);
}

export { vueUseI18n };
export default i18n;
