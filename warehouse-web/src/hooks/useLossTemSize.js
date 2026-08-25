import { ref, reactive, computed } from 'vue'

export default function useLossTemSize() {
  const reasons = ref([])

  const summaryData = ref([])
  const summaryLoading = ref(false)
  const summaryError = ref('')

  const filterMonth = ref('')
  const filterKeyword = ref('')
  const filterReason = ref('')
  const filterStatus = ref('')
  const filterSort = ref('newest')

  const summaryPage = ref(1)
  const summaryPageSize = ref(10)

  const searching = ref(false)
  const lenhError = ref('')

  const currentLenh = ref(null)
  const currentReason = ref(null)

  const mode = ref('')

  const showEntryDialog = ref(false)

  const sizeInputs = reactive({})

  const subTotal = ref(0)

  const msnvInput = ref('')

  const toNum = (v) => Number(v || 0)

  const reasonShort = (v) => {
    if (!v) return ''
    return v.length > 25 ? v.slice(0, 25) + '...' : v
  }

  const reasonTagType = (v) => {
    if (!v) return 'info'

    const text = v.toLowerCase()

    if (text.includes('in')) return 'success'
    if (text.includes('out')) return 'danger'

    return 'warning'
  }

  const sumObjectValues = (obj) => {
    if (!obj) return 0

    return Object.values(obj).reduce((a, b) => {
      return a + Number(b || 0)
    }, 0)
  }

  const summaryAllSizes = computed(() => {
    const set = new Set()

    summaryData.value.forEach((row) => {
      Object.keys(row.size_data || {}).forEach((k) => {
        set.add(k)
      })
    })

    return [...set]
  })

  const filteredSummaryData = computed(() => {
    let arr = [...summaryData.value]

    if (filterKeyword.value) {
      const kw = filterKeyword.value.toLowerCase()

      arr = arr.filter((x) => {
        return (
          x.lenh?.toLowerCase().includes(kw) ||
          x.msnv?.toLowerCase().includes(kw)
        )
      })
    }

    if (filterReason.value) {
      arr = arr.filter((x) => x.reason === filterReason.value)
    }

    if (filterStatus.value) {
      arr = arr.filter((x) => x.status === filterStatus.value)
    }

    return arr
  })

  const paginatedSummaryData = computed(() => {
    const start =
      (summaryPage.value - 1) * summaryPageSize.value

    const end = start + summaryPageSize.value

    return filteredSummaryData.value.slice(start, end)
  })

  const grandInTotal = computed(() => {
    return filteredSummaryData.value.reduce((a, b) => {
      return a + Number(b.in_total || 0)
    }, 0)
  })

  const grandOutTotal = computed(() => {
    return filteredSummaryData.value.reduce((a, b) => {
      return a + Number(b.out_total || 0)
    }, 0)
  })

  const grandNetTotal = computed(() => {
    return grandInTotal.value - grandOutTotal.value
  })

  const fetchLenhFromApi = async (lenh) => {
    return {
      lenh,
      xxccSizes: ['S', 'M', 'L', 'XL'],
    }
  }

  const loadSummaryFromAPI = async () => {
    try {
      summaryLoading.value = true

      summaryData.value = []
    } catch (err) {
      summaryError.value = err.message
    } finally {
      summaryLoading.value = false
    }
  }

  const handleSearch = async (lenh, entryMode) => {
    lenhError.value = ''

    if (!lenh.trim()) {
      lenhError.value = 'Vui lòng nhập số lệnh'
      return
    }

    searching.value = true

    try {
      const data = await fetchLenhFromApi(lenh)

      if (!data) {
        lenhError.value = `Không tìm thấy lệnh ${lenh}`
        return
      }

      mode.value = entryMode
      currentLenh.value = data
      currentReason.value = null
      showEntryDialog.value = true

      Object.keys(sizeInputs).forEach((k) => delete sizeInputs[k])

      data.xxccSizes.forEach((sz) => {
        sizeInputs[sz] = null
      })

      subTotal.value = 0
    } catch (err) {
      lenhError.value = err.message
    } finally {
      searching.value = false
    }
  }

  const updateSubTotal = () => {
    subTotal.value = Object.values(sizeInputs).reduce(
      (a, v) => a + toNum(v),
      0,
    )
  }

  return {
    reasons,
    summaryData,
    summaryLoading,
    summaryError,
    filterMonth,
    filterKeyword,
    filterReason,
    filterStatus,
    filterSort,
    summaryPage,
    summaryPageSize,
    searching,
    lenhError,
    currentLenh,
    currentReason,
    mode,
    showEntryDialog,
    sizeInputs,
    subTotal,
    msnvInput,
    filteredSummaryData,
    summaryAllSizes,
    paginatedSummaryData,
    grandInTotal,
    grandOutTotal,
    grandNetTotal,
    loadSummaryFromAPI,
    handleSearch,
    updateSubTotal,
    reasonShort,
    reasonTagType,
    sumObjectValues,
  }
}