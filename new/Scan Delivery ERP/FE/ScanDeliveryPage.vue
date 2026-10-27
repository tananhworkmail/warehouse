<!-- src/views/ScanDeliveryPage.vue -->
<template>
    <div class="scan-delivery-page">
        <!-- LEFT (nếu dùng thì bật lại) -->
        <!--
        <ScanDeliverLeftTable v-if="!isScanMoreTableOpen" />
        <OweTable v-else />
        -->

        <!-- RIGHT -->
        <ScanDeliverRightTable />

        <!-- PRINT LAYOUTS -->
        <PrintFileLayout v-if="print_file_layout_show" />
        <PrintExportQuantity v-if="printex_file_layout_show" />
    </div>
</template>

<script setup>
import ScanDeliverLeftTable from '@/components/ScanDeliveryPage/ScanDeliverLeftTable.vue'
import ScanDeliverRightTable from '@/components/ScanDeliveryPage/ScanDeliverRightTable.vue'
import PrintFileLayout from '@/components/ScanDeliveryPage/PrintFileLayout.vue'
import PrintExportQuantity from '@/components/ScanDeliveryPage/PrintExportQuantity.vue'
import OweTable from '@/components/ScanDeliveryPage/OweTable.vue'
import { onBeforeUnmount } from 'vue'
import useScanDelivery from '@/hooks/useScanDelivery'
const { init_status, print_file_layout_show, printex_file_layout_show, isScanMoreTableOpen } =
    useScanDelivery()
onBeforeUnmount(() => {
    init_status()
})
</script>

<style scoped>
.scan-delivery-page {
    width: 100%;
    height: 100vh;
    display: flex;
    overflow: hidden;
}
</style>
