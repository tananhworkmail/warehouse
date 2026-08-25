const guideClasses = ['cell-guide-row', 'cell-guide-column', 'cell-guide-active']

function clearGuide(table) {
  table.querySelectorAll(`.${guideClasses.join(', .')}`).forEach((cell) => {
    cell.classList.remove(...guideClasses)
  })
}

function mapTableCells(table) {
  const rowSpans = []
  table.querySelectorAll('tr').forEach((row) => {
    let column = 0
    Array.from(row.cells).forEach((cell) => {
      while ((rowSpans[column] || 0) > 0) column += 1
      const span = Math.max(1, Number(cell.colSpan) || 1)
      cell.dataset.guideColumnStart = String(column)
      cell.dataset.guideColumnEnd = String(column + span - 1)
      const rowSpan = Math.max(1, Number(cell.rowSpan) || 1)
      for (let index = column; index < column + span; index += 1) {
        rowSpans[index] = Math.max(rowSpans[index] || 0, rowSpan)
      }
      column += span
    })
    for (let index = 0; index < rowSpans.length; index += 1) {
      if (rowSpans[index] > 0) rowSpans[index] -= 1
    }
  })
}

const handlers = new WeakMap()

export const cellGuideDirective = {
  mounted(table) {
    const focusin = (event) => {
      if (!event.target.matches('input, textarea, select')) return
      const activeCell = event.target.closest('td, th')
      if (!activeCell || !table.contains(activeCell)) return
      clearGuide(table)
      mapTableCells(table)
      const start = Number(activeCell.dataset.guideColumnStart)
      const end = Number(activeCell.dataset.guideColumnEnd)
      Array.from(activeCell.parentElement.cells).forEach((cell) => cell.classList.add('cell-guide-row'))
      table.querySelectorAll('td, th').forEach((cell) => {
        const cellStart = Number(cell.dataset.guideColumnStart)
        const cellEnd = Number(cell.dataset.guideColumnEnd)
        if (cellStart <= end && cellEnd >= start) cell.classList.add('cell-guide-column')
      })
      activeCell.classList.add('cell-guide-active')
    }
    const focusout = (event) => {
      if (!table.contains(event.relatedTarget)) clearGuide(table)
    }
    table.addEventListener('focusin', focusin)
    table.addEventListener('focusout', focusout)
    handlers.set(table, { focusin, focusout })
  },
  unmounted(table) {
    const registered = handlers.get(table)
    if (!registered) return
    table.removeEventListener('focusin', registered.focusin)
    table.removeEventListener('focusout', registered.focusout)
    handlers.delete(table)
  },
}
