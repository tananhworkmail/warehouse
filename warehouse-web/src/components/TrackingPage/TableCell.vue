<template>
    <!-- Cột IMQC Report: hiển thị link click mở file -->
    <div
        v-if="row_type === 'left_table_row' && isIMQCFile"
        class="read-content imqc-link"
        @click.stop="openIMQCReport"
        :title="t('trackingPage.tableCell.viewImqcReport')"
    >
        📄 {{ content_value }}
    </div>
    <!-- Cột thường cho left_table_row -->
    <div
        v-else-if="row_type === 'left_table_row'"
        class="read-content"
        :contenteditable="false"
        @click="handleclick"
    >
        {{ content_value }}
    </div>
    <input
        v-else-if="row_type === 'right_top_table_row' && col_index === 1"
        :class="'editable-content'"
        :contenteditable="isOnEdit"
        :value="content_value"
        style="width: 100%;"
        @blur="handblur"
        ref="editableDiv"
        @input="handleInput"
    />
    <div
        v-else-if="row_type === 'right_bottom_table_row'"
        :class="'read-content'"
        @dblclick="handleOweTableRowDBclick"
    >
        {{ content_value }}
    </div> <div
        v-else
        :class="'read-content'"
    >
        {{ content_value }}
    </div>
</template>

<script setup>
import { toRefs, ref, computed, onUpdated, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import useTracking from '../../hooks/useTracking'

const API_PREFIX = (import.meta.env.VITE_API_URL || '/api/v1').replace(/\/+$/, '')

const { t } = useI18n()

const {
    getKCLLScan,
    right_table_selceted_row_colume,
    left_table_selceted_row_colume,
    KCLLScan_RFSS_row_list,
    KCLLScan_RFSSS_row_list,
    update_KCLLScan_RFSS_row_list,
    update_KCLLScan_RFSSS_row_list,
} = useTracking()
const props = defineProps([
    'value',
    'row_type',
    'row_id',
    'col_index',
    'row_index',
])
const { value, row_type, row_id, col_index, row_index } = toRefs(props)
const isOnEdit = ref(false)
const editableDiv = ref(null)
const start = ref()
const content_value = ref(value.value)

// Kiểm tra xem giá trị ô hiện tại có phải là file IMQC Report không
const isIMQCFile = computed(() => {
    if (!content_value.value || typeof content_value.value !== 'string') return false
    const val = content_value.value.toLowerCase()
    return val.endsWith('.pdf') || val.endsWith('.jpg') || val.endsWith('.png') || val.endsWith('.jpeg')
})

// Mở file IMQC Report thông qua API (giống Entry page)
const openIMQCReport = () => {
    if (!content_value.value) return
    const url = `${API_PREFIX}/scan/entry/imqcreport?File_Name=${encodeURIComponent(content_value.value)}`
    window.open(url, '_blank')
}
// 監聽 value 的變化，並更新 content_value
watch(value, (newVal) => {
    content_value.value = newVal
})

// 監聽 selceted_row_colume 的變化
watch(
    left_table_selceted_row_colume,
    (newVal) => {
        if (
            newVal.row === row_index.value &&
            newVal.colume === col_index.value
        ) {
            isOnEdit.value = true
            getKCLLScan(row_id.value)
        } else {
            isOnEdit.value = false
        }
    },
    { deep: true },
)

onUpdated(() => {
    if (isOnEdit.value && editableDiv.value) {
        editableDiv.value.focus()
    }
})

onMounted(() => {
    if (isOnEdit.value && editableDiv.value) {
        editableDiv.value.focus()
    }
})

const handleclick = () => {
    left_table_selceted_row_colume.row = row_index.value
    left_table_selceted_row_colume.colume = col_index.value
}

function handleOweTableRowDBclick() {
    right_table_selceted_row_colume.row = row_index.value
    right_table_selceted_row_colume.colume = col_index.value
    const new_KCLLScan_RFSSS_row = {
        ...KCLLScan_RFSSS_row_list.value[
            right_table_selceted_row_colume.row
        ],
        ScanQty_input: parseFloat('0'),
    }

    update_KCLLScan_RFSSS_row_list(
        new_KCLLScan_RFSSS_row,
        right_table_selceted_row_colume.row,
    )
}

function handblur(event) {
    content_value.value = event.target.value ?
        parseFloat(event.target.value) : 0
    const new_KCLLScan_RFSS_row = {
        ...KCLLScan_RFSS_row_list.value[
            right_table_selceted_row_colume.row
        ],
        ScanQty_input: parseFloat(content_value.value),
    }

    update_KCLLScan_RFSS_row_list(
        new_KCLLScan_RFSS_row,
        right_table_selceted_row_colume.row,
    )
    right_table_selceted_row_colume.row = ''
    right_table_selceted_row_colume.colume = ''
    isOnEdit.value = false
}

// 更新 content_value 當用戶輸入時
function handleInput(event) {
    right_table_selceted_row_colume.row = row_index.value
    right_table_selceted_row_colume.colume = col_index.value
    isOnEdit.value = true
    content_value.value = event.target.value ?
        parseFloat(event.target.value) : 0
}

</script>

<style lang="css" scoped>
.editable-content,
.read-content {
    white-space: nowrap;
    overflow: hidden; /* Hide content overflow */
    padding: 5px;
    height: max-content;
    /* border-style: inset; */
    font-size: 14px !important;
    overflow-x: auto; /* Enable horizontal scrolling */
    max-width: 500px;
}

/* Hide scrollbars for both editable and read-only content */
.editable-content::-webkit-scrollbar,
.read-content::-webkit-scrollbar {
    display: none;
}

.editable-content {
    background-color: lightblue;
    -ms-overflow-style: none; /* IE and Edge */
    scrollbar-width: none; /* Firefox */
}

/* Style cho cột IMQC Report - clickable */
.imqc-link {
    color: #4fc3f7 !important;
    cursor: pointer;
    text-decoration: underline;
    transition: color 0.2s ease;
}

.imqc-link:hover {
    color: #81d4fa !important;
    text-shadow: 0 0 4px rgba(79, 195, 247, 0.5);
}
</style>
