// src/hooks/useWarehouseMapI18n.js
import vi from '@/locales/vi'
import en from '@/locales/en'
import zh from '@/locales/zh'

// Use the shared locale singleton so header language changes stay reactive.
import { locale } from '@/hooks/i18n'

const supportedLangs = ['vi', 'en', 'zh']

export const langOptions = [
  { code: 'vi', label: 'VI' },
  { code: 'en', label: 'EN' },
  { code: 'zh', label: '中文' },
]

const messages = { vi, en, zh }

function resolve(obj, path) {
  return path.split('.').reduce((o, k) => o?.[k], obj) ?? path
}

export function useWarehouseMapI18n() {
  const t = (path) => {
    const lang = messages[locale.value] || messages.vi
    return resolve(lang, path)
  }

  const setLang = (code) => {
    if (!supportedLangs.includes(code)) return
    locale.value = code
    localStorage.setItem('lang', code)
  }

  return {
    t,
    lang: locale,
    locale,
    langOptions,
    setLang,
  }
}

export default { useWarehouseMapI18n }
