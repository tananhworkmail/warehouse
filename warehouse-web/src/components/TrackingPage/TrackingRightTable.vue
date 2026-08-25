<!-- src/components/TrackingPage/TrackingRightTable.vue -->
<template>
    <div style="
            display: flex;
            align-items: center;
            justify-content: start;
            flex-direction: column;
            background-color: gray;
        ">
        <div>
            <div class="navbar1">
                <v-dialog max-width="500" max-height="300">
                    <template v-slot:activator="{ props: activatorProps }">
                        <v-icon v-bind="activatorProps" icon="$checkedCircle" style="
                                color: white;
                                font-size: 40px;
                                transition: background-color 0.3s ease;
                                background: linear-gradient(
                                    to right,
                                    #00bfff,
                                    #1e90ff
                                );
                                border: 2px solid white;
                                margin-left: 5px;
                                border-radius: 50%;
                                margin-right: 10px;
                            "></v-icon>
                    </template>
                    <template v-slot:default="{ isActive }">
                        <v-card :title="t('trackingPage.rightTable.saveModificationTitle')">
                            <v-card-actions>
                                <v-spacer></v-spacer>
                                <v-btn :text="t('common.confirm')" style="color: green" @click="handleupdateRFSS(isActive)"></v-btn>
                                <v-btn :text="t('common.back')" style="color: gray" @click="isActive.value = false"></v-btn>
                            </v-card-actions>
                        </v-card>
                    </template>
                </v-dialog>
            </div>
            <v-data-table-virtual density="compact" :key="JSON.stringify(KCLLScan_RFSS_row_list)"
                :items="KCLLScan_RFSS_row_list" :headers="trackingpage_right_table_headars" :sort-desc="[false, true]"
                :hide-default-footer="true" style="background-color: black" class="fixed-header-table"
                :style="dynamicStyles">
                <template v-slot:headers="{ }">
                    <TableHeader :headers="trackingpage_right_table_headars"></TableHeader>
                </template>
                <template v-slot:item="{ item, index }">
                    <TableRow :rowsData="item" row_type="right_top_table_row" :row_index="index"></TableRow>
                </template>
                <template v-slot:bottom>
                    <tr>
                        <td class="summary-cell"></td>
                        <td class="data-cell">{{ t('trackingPage.rightTable.qtyOut') }}</td>
                    </tr>
                    <tr>
                        <td class="summary-cell">{{ t('trackingPage.rightTable.total') }}</td>
                        <td class="total-cell">{{ totalQtySS.toFixed(2) }}</td>
                    </tr>
                </template>
            </v-data-table-virtual>
        </div>
        <div>
            <div class="navbar1">
                <v-dialog v-model="dialog" max-width="500" max-height="300">
                    <v-card :title="t('trackingPage.rightTable.saveModificationTitle')">
                        <v-card-actions>
                            <v-spacer></v-spacer>
                            <v-btn :text="t('common.confirm')" style="color: green" @click="handleupdateRFSSS()"></v-btn>
                            <v-btn :text="t('common.back')" style="color: gray" @click="handleClose()"></v-btn>
                        </v-card-actions>
                    </v-card>
                </v-dialog>
            </div>
            <v-data-table-virtual class="fixed-header-table" density="compact"
                :key="JSON.stringify(KCLLScan_RFSSS_row_list)" :items="KCLLScan_RFSSS_row_list"
                :headers="trackingpage_right_table_headars" :sort-desc="[false, true]" :hide-default-footer="true"
                style="background-color: black" :style="dynamicStyles">
                <template #headers>
                    <TableHeader :headers="trackingpage_right_table_headars"></TableHeader>
                </template>
                <template v-slot:item="{ item, index }">
                    <TableRow :rowsData="item" row_type="right_bottom_table_row" :row_index="index"></TableRow>
                </template>

                <template v-slot:bottom>
                    <tr>
                        <td class="summary-cell"></td>
                        <td class="data-cell">{{ t('trackingPage.rightTable.qtyOutDebt') }}</td>
                    </tr>
                    <tr>
                        <td class="summary-cell">{{ t('trackingPage.rightTable.total') }}</td>
                        <td class="total-cell">{{ totalQtySSS.toFixed(2) }}</td>
                    </tr>
                </template>
            </v-data-table-virtual>
        </div>
    </div>
</template>

<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import TableRow from './TableRow.vue'
import TableHeader from '../TableHeader.vue'
import { trackingpage_right_table_headars } from '../../constants/headers'
import useTracking from '../../hooks/useTracking'
import useScreenAndResizableTable from '@/hooks/useScreen'

const { t } = useI18n()
const { mainTableHeight, mainTableWidth } = useScreenAndResizableTable(
    0.25,
    0.985,
    0.317,
    0.985,
)
const {
    KCLLScan_RFSS_row_list,
    KCLLScan_RFSSS_row_list,
    Barcode_old,
    getKCLLScan,
    UpdateKCLLScan_RFSS,
    UpdateKCLLScan_RFSSS,
} = useTracking()
onBeforeUnmount(() => {
    Clear_Table_Data()
})
function Clear_Table_Data() {
    KCLLScan_RFSS_row_list.value.length = 0
    KCLLScan_RFSSS_row_list.value.length = 0
}
const dynamicStyles = computed(() => ({
    height: mainTableHeight.value + 'px',
    width: mainTableWidth.value + 'px',
}))
function handleupdateRFSS(isActive) {
    isActive.value = false
    UpdateKCLLScan_RFSS()
}
const dialog = ref(computed(() => Barcode_old.value!=''))
function handleupdateRFSSS() {
    dialog.value = false
    UpdateKCLLScan_RFSSS()
    Barcode_old.value=''
}
function handleClose() {
    dialog.value = false,
    getKCLLScan(Barcode_old.value),
    Barcode_old.value=''
}
const totalQtySS = computed(() => {
    return KCLLScan_RFSS_row_list.value.reduce((sum, item) => {
        return sum + (item.LLNO ? parseFloat(item.SCANQty) : 0);
    }, 0.0);
});
const totalQtySSS = computed(() => {
    return KCLLScan_RFSSS_row_list.value.reduce((sum, item) => {
        return sum + (item.LLNO_OWE ? parseFloat(item.SCANQty) : 0);
    }, 0.0);
});
</script>

<style scoped>
.navbar1 {
    position: relative;
    background: linear-gradient(to bottom, #000080, #4169e1);
    display: flex;
    align-items: center;
    padding: 20px;
    border: 1px solid white;
    border-style: inset;
    height: 5vh;
}

.navbar2 {
    position: relative;
    background: linear-gradient(to bottom, #000080, #4169e1);
    display: flex;
    height: 10vh;
    align-items: center;
    width: 100%;
    padding: 20px;
    border: 1px solid white;
    border-style: inset;
}

.summary-cell {
    font-size: 16px;
    font-weight: 600;
    color: yellow;
    font-family: Arial, Helvetica, sans-serif;
    border: 1px solid white;
    min-width: 40vh;
}

.data-cell,
.total-cell {
    white-space: nowrap;
    text-align: center;
    background-color: black;
    color: cyan;
    font-weight: 600;
    font-size: 14px;
    border: 1px solid lightgray;
    width: 45%;
}

.total-cell {
    color: yellow;
}

@media (max-width: 1024px) {
    .summary-cell {
        min-width: 16.3vh;
    }

    .data-cell,
    .total-cell {
        width: 45.2%;
        padding: 5px;
    }

    .total-cell {
        font-size: 14px;
    }
}
</style>