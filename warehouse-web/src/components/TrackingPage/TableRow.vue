<!-- src/components/TrackingPage/TableRow.vue -->
<template>
    <tr @click="$emit('row-clicked', rowsData)" class="clickable-row" :class="isSelected ? 'selected' : 'not-select'">
        <td v-for="(item, index) in rowsDataList" :key="index" style="
                font-size: 16px;
                font-weight: 600;
                color: Yellow;
                font-family: Arial, Helvetica, sans-serif;
                border: 1px solid white;
                padding: 0%;
            ">
            <TableCell :row_id="rowsData.Barcode" :isEditableRow="isEditableRow" :value="item" :col_index="index"
                :row_index="row_index" :row_type="row_type" />
        </td>
    </tr>
</template>

<script setup>
import { computed, toRefs } from 'vue'
import TableCell from './TableCell.vue'
import useTracking from '@/hooks/useTracking'
import { useRowOperations } from '@/hooks/useRowOperations'

const { left_table_selceted_row_colume, Barcode_old } = useTracking()

const props = defineProps([
    'rowsData',
    'isEditableRow',
    'row_index',
    'row_type',
])

// Định nghĩa emit để báo cho cha biết dòng nào được click
defineEmits(['row-clicked'])

const { rowsData, isEditableRow, row_index, row_type } = toRefs(props)

const isSelected = computed(() => {
    if (row_type.value === 'left_table_row') {
        return left_table_selceted_row_colume.row === row_index.value
    } else if (row_type.value === 'right_bottom_table_row') {
        return Barcode_old.value != '';
    }
    return false
})

const DateKey = ['CFMDATE1', 'CFMDATE', 'SCANUSERDATE']
const key_need_to_be_excluded = {
    left_table_row: ['', 'Barcode', 'Pack1'],
    right_top_table_row: ['Barcode', 'CLBH', 'SCANQty', 'USERID', 'CLBH'],
    right_bottom_table_row: ['Barcode', 'SCANQty', 'USERID', 'CLBH', 'LLNO'],
}
const { createDisplayRow } = useRowOperations()

// ⚠️ Bọc trong computed() để khi rowsData thay đổi (scan mới, cập nhật dữ liệu...)
// bảng tự re-render đúng, thay vì chỉ tính 1 lần lúc component khởi tạo.
const rowsDataList = computed(() =>
    createDisplayRow(
        rowsData.value,
        key_need_to_be_excluded[props.row_type],
        DateKey,
    ),
)
</script>

<style scoped>
.clickable-row {
    cursor: pointer;
    transition: background 0.2s;
}

.clickable-row:hover {
    background-color: rgba(255, 255, 255, 0.1) !important;
}

.selected {
    background-color: #333 !important;
}
</style>