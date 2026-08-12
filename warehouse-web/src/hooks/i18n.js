// src/hooks/i18n.js
import { ref } from 'vue'
import vi from '@/locales/vi'
import en from '@/locales/en'
import zh from '@/locales/zh'

// ─── Singleton locale — EXPORT để useWarehouseMapI18n dùng chung ────────────
// Quan trọng: chỉ có 1 ref duy nhất cho toàn app.
// Nếu mỗi hook tự tạo ref riêng thì đổi ngôn ngữ ở header sẽ không
// cập nhật các component khác (phải reload trang).
export const locale = ref(localStorage.getItem('lang') || 'vi')

const messages = { vi, en, zh }

function resolve(obj, path, params) {
  const val = path.split('.').reduce((o, k) => (o != null ? o[k] : undefined), obj)
  if (val == null) return path
  if (params && typeof val === 'string') {
    return val.replace(/\{(\w+)\}/g, (_, k) => (params[k] != null ? params[k] : `{${k}}`))
  }
  return val
}

export function useI18n() {
  const t = (key, params) => {
    const lang = messages[locale.value] || messages.vi
    return resolve(lang, key, params)
  }

  return { t, locale }
}

export default { useI18n }