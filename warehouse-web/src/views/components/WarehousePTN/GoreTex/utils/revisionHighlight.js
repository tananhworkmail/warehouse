const directiveState = new WeakMap()
let activeHost = null

function clearHighlight(element) {
  const state = directiveState.get(element)
  if (!state) return

  state.host.classList.remove('revision-changed-cell', 'revision-tooltip-pinned')
  if (activeHost === state.host) activeHost = null
  state.host.removeAttribute('tabindex')
  state.tooltip.remove()
  state.host.removeEventListener('click', state.toggle)
  directiveState.delete(element)
}

function applyHighlight(element, change) {
  clearHighlight(element)
  if (!change) return

  const host = element.closest('td, label') || element.parentElement
  if (!host) return

  const tooltip = document.createElement('span')
  tooltip.className = 'revision-cell-tooltip'
  tooltip.setAttribute('role', 'tooltip')

  const content = document.createElement('strong')
  content.textContent = change.before
  tooltip.append(content)

  const toggle = (event) => {
    if (event.target.closest('.revision-cell-tooltip')) return
    const shouldOpen = !host.classList.contains('revision-tooltip-pinned')
    if (activeHost && activeHost !== host) {
      activeHost.classList.remove('revision-tooltip-pinned')
    }
    host.classList.toggle('revision-tooltip-pinned', shouldOpen)
    activeHost = shouldOpen ? host : null
  }

  host.classList.add('revision-changed-cell')
  host.setAttribute('tabindex', '0')
  host.append(tooltip)
  host.addEventListener('click', toggle)
  directiveState.set(element, { host, tooltip, toggle })
}

export const revisionChangeDirective = {
  mounted(element, binding) {
    applyHighlight(element, binding.value)
  },
  updated(element, binding) {
    if (binding.value !== binding.oldValue) applyHighlight(element, binding.value)
  },
  beforeUnmount(element) {
    clearHighlight(element)
  },
}
