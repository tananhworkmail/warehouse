function hasInputValue(field) {
  return String(field.value ?? '').trim() !== ''
}

export function isNonNegativeNumber(value) {
  const normalized = String(value ?? '').trim().replace(',', '.')
  if (normalized === '') return false
  const parsed = Number(normalized)
  return Number.isFinite(parsed) && parsed >= 0
}

function isFieldValid(field) {
  if (!hasInputValue(field)) return false
  if (field.dataset.nonnegative === 'true') return isNonNegativeNumber(field.value)
  return true
}

export function revealMissingFields(formElement) {
  if (!formElement) return 0

  const requiredFields = [...formElement.querySelectorAll('input[required], textarea[required], select[required]')]
    .filter((field) => !field.disabled)

  const missingFields = requiredFields.filter((field) => !isFieldValid(field))
  requiredFields.forEach((field) => {
    const isMissing = missingFields.includes(field)
    field.classList.toggle('field-missing', isMissing)
    field.setAttribute('aria-invalid', String(isMissing))
    if (!field.dataset.missingListener) {
      field.dataset.missingListener = 'true'
      field.addEventListener('input', () => {
        const stillMissing = !isFieldValid(field)
        field.classList.toggle('field-missing', stillMissing)
        field.setAttribute('aria-invalid', String(stillMissing))
      })
      field.addEventListener('change', () => {
        const stillMissing = !isFieldValid(field)
        field.classList.toggle('field-missing', stillMissing)
        field.setAttribute('aria-invalid', String(stillMissing))
      })
    }
  })

  if (missingFields[0]) {
    missingFields[0].scrollIntoView({ behavior: 'smooth', block: 'center', inline: 'center' })
    window.setTimeout(() => missingFields[0].focus({ preventScroll: true }), 350)
  }
  return missingFields.length
}

export function clearMissingFields(formElement) {
  formElement
    ?.querySelectorAll('.field-missing')
    .forEach((field) => {
      field.classList.remove('field-missing')
      field.removeAttribute('aria-invalid')
    })
}
