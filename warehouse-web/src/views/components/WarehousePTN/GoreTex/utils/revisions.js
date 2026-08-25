export function cloneFormData(data) {
  return JSON.parse(JSON.stringify(data || {}))
}

export function getCurrentFormData(data) {
  const current = cloneFormData(data)
  delete current._editHistory
  return current
}

export function getEditHistory(data) {
  return Array.isArray(data?._editHistory)
    ? cloneFormData(data._editHistory)
    : []
}

export function createLocalRevision(data) {
  const now = new Date()
  const offset = now.getTimezoneOffset() * 60_000
  return {
    editedAt: new Date(now.getTime() - offset).toISOString().slice(0, 19),
    data: getCurrentFormData(data),
  }
}

const FIELD_LABELS = {
  line: 'Chuyền / Line',
  styleName: 'Dạng giày',
  inspectionDate: 'Ngày kiểm tra',
  testDate: 'Ngày thử nghiệm',
  testedCount: 'Số lượng thử nghiệm',
  failedCount: 'Số lượng không đạt',
  rate: 'Tỉ lệ',
  leakagePosition: 'Vị trí rỉ nước',
  cause: 'Phân tích nguyên nhân',
  action: 'Biện pháp xử lý',
  unit: 'Đơn vị phát sinh / cải thiện',
  improvementDate: 'Ngày cải thiện',
  qc: 'QC',
  style: 'Dạng giày',
  po: 'PO',
  size: 'Size',
  lean: 'Lean',
  otherIssue: 'Vấn đề khác',
  result: 'Kết quả',
  note: 'Ghi chú',
}

function flatten(value, prefix = '', output = {}) {
  if (Array.isArray(value)) {
    value.forEach((item, index) => flatten(item, prefix ? `${prefix}.${index}` : String(index), output))
    return output
  }
  if (value && typeof value === 'object') {
    Object.entries(value).forEach(([key, child]) => {
      if (key === '_editHistory' || key === 'id') return
      flatten(child, prefix ? `${prefix}.${key}` : key, output)
    })
    return output
  }
  if (prefix) output[prefix] = value
  return output
}

function labelForPath(path) {
  const parts = path.split('.')
  const rowIndex = parts.findIndex((part) => /^\d+$/.test(part))
  const rowLabel = rowIndex >= 0 ? `Dòng ${Number(parts[rowIndex]) + 1} – ` : ''
  const field = parts.at(-1)
  if (FIELD_LABELS[field]) return `${rowLabel}${FIELD_LABELS[field]}`

  if (path.startsWith('counts.')) {
    const countKey = parts.slice(1).join(' · ').replaceAll('-', ' / ')
    return `Số lượng – ${countKey}`
  }
  if (path.startsWith('totals.')) return `Tổng – ${parts.slice(1).join(' · ')}`
  if (path.startsWith('rates.')) return `Tỉ lệ – ${parts.slice(1).join(' · ')}`
  if (path.startsWith('issueValues.')) return `${rowLabel}Hạng mục – ${parts.slice(rowIndex + 2).join(' · ')}`

  return `${rowLabel}${field.replace(/([A-Z])/g, ' $1').trim()}`
}

function displayValue(value) {
  if (value === null || value === undefined || value === '') return '—'
  if (typeof value === 'boolean') return value ? 'Có' : 'Không'
  return String(value)
}

export function getRevisionChanges(beforeData, afterData) {
  const before = flatten(getCurrentFormData(beforeData))
  const after = flatten(getCurrentFormData(afterData))
  const paths = [...new Set([...Object.keys(before), ...Object.keys(after)])].sort()

  return paths
    .filter((path) => String(before[path] ?? '') !== String(after[path] ?? ''))
    .map((path) => ({
      path,
      label: labelForPath(path),
      before: displayValue(before[path]),
      after: displayValue(after[path]),
    }))
}
