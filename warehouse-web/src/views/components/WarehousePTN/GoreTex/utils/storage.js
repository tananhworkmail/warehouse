const prefix = 'quality-desk:'

export function saveDraft(key, data) {
  const savedAt = new Date().toISOString()
  localStorage.setItem(`${prefix}${key}`, JSON.stringify({ data, savedAt }))
  return savedAt
}

export function loadDraft(key, fallback) {
  try {
    const value = JSON.parse(localStorage.getItem(`${prefix}${key}`))
    return value?.data ? { ...value, exists: true } : { data: fallback, savedAt: null, exists: false }
  } catch {
    return { data: fallback, savedAt: null, exists: false }
  }
}

export function clearDraft(key) {
  localStorage.removeItem(`${prefix}${key}`)
}

export function getDraftMeta(key) {
  try {
    const value = JSON.parse(localStorage.getItem(`${prefix}${key}`))
    return value?.savedAt || null
  } catch {
    return null
  }
}

export function formatSavedAt(value) {
  if (!value) return 'Chưa có bản nháp'
  return `Lưu lúc ${new Intl.DateTimeFormat('vi-VN', {
    hour: '2-digit',
    minute: '2-digit',
    day: '2-digit',
    month: '2-digit',
  }).format(new Date(value))}`
}

export function todayString() {
  const parts = new Intl.DateTimeFormat('en-CA', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }).formatToParts(new Date())
  const get = (type) => parts.find((part) => part.type === type)?.value
  return `${get('year')}-${get('month')}-${get('day')}`
}
