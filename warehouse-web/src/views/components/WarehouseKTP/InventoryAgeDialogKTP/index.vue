<template>
  <el-dialog
    :model-value="dialogVisible"
    :title="t('ktp.inventoryAge.detailTitle')"
    width="min(1080px, 96vw)"
    class="inventory-age-dialog"
    align-center
    append-to-body
    destroy-on-close
    @close="emitClose"
  >
    <div class="inventory-age-table-wrap">
      <el-table
        :data="pagedRows"
        height="min(62vh, 590px)"
        stripe
        table-layout="fixed"
        :default-sort="{ prop: 'storageAgeDays', order: 'descending' }"
        @sort-change="handleSortChange"
      >
        <el-table-column
          prop="rackCode"
          :label="t('ktp.inventoryAge.storageLocation')"
          min-width="145"
        />
        <el-table-column
          prop="ddbh"
          :label="t('ktp.inventoryAge.orderNo')"
          min-width="190"
        />
        <el-table-column
          prop="qty"
          :label="t('ktp.inventoryAge.pairs')"
          min-width="120"
          align="right"
        >
          <template #default="{ row }">{{ formatNumber(row.qty) }}</template>
        </el-table-column>
        <el-table-column
          prop="lastInDate"
          :label="t('ktp.inventoryAge.inboundTime')"
          min-width="180"
        >
          <template #default="{ row }">{{ row.lastInDate || "--" }}</template>
        </el-table-column>
        <el-table-column
          prop="storageAgeDays"
          :label="t('ktp.inventoryAge.storageAge')"
          min-width="205"
          sortable="custom"
        >
          <template #default="{ row }">
            <span class="age-badge" :style="ageBadgeStyle(row.storageAgeDays)">
              {{ ageBadgeText(row.storageAgeDays) }}
            </span>
          </template>
        </el-table-column>
        <template #empty>
          <div class="empty-state">{{ t("ktp.inventoryAge.noData") }}</div>
        </template>
      </el-table>
    </div>

    <div v-if="rows.length" class="inventory-pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="rows.length"
        layout="total, sizes, prev, pager, next"
      />
    </div>

    <template #footer>
      <el-button type="primary" @click="emitClose">
        {{ t("ktp.rackDetail.close") }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { useI18n } from "@/hooks/i18n";
import { getInventoryAgeCategory } from "@/views/components/WarehouseKTP/inventoryAge";

const props = defineProps({
  dialogVisible: Boolean,
  rows: {
    type: Array,
    default: () => [],
  },
});

const emit = defineEmits(["update:dialogVisible", "close"]);
const { t, locale } = useI18n();
const currentPage = ref(1);
const pageSize = ref(10);
const sortOrder = ref("descending");

const sortedRows = computed(() => {
  const direction = sortOrder.value === "ascending" ? 1 : -1;
  return [...props.rows].sort((left, right) => {
    const leftAge =
      left.storageAgeDays !== null && Number.isFinite(Number(left.storageAgeDays))
        ? Number(left.storageAgeDays)
        : -1;
    const rightAge =
      right.storageAgeDays !== null && Number.isFinite(Number(right.storageAgeDays))
        ? Number(right.storageAgeDays)
        : -1;
    return (leftAge - rightAge) * direction;
  });
});

const pagedRows = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  return sortedRows.value.slice(start, start + pageSize.value);
});

const formatNumber = (value) => {
  const numberLocales = { vi: "vi-VN", en: "en-US", zh: "zh-CN" };
  return Number(value || 0).toLocaleString(numberLocales[locale.value] || "vi-VN");
};

const ageBadgeStyle = (days) => {
  const category = getInventoryAgeCategory(days);
  return {
    "--age-color": category.color,
    "--age-soft-color": category.softColor,
  };
};

const ageBadgeText = (days) => {
  const category = getInventoryAgeCategory(days);
  if (category.key === "unknown") return t(category.labelKey);
  return `${t("ktp.inventoryAge.dayCount", { count: days })} / ${t(category.labelKey)}`;
};

const handleSortChange = ({ order }) => {
  sortOrder.value = order || "descending";
  currentPage.value = 1;
};

const emitClose = () => {
  emit("update:dialogVisible", false);
  emit("close");
};

watch([() => props.dialogVisible, pageSize], () => {
  currentPage.value = 1;
});
</script>

<style lang="scss" scoped>
:deep(.inventory-age-dialog) {
  display: flex;
  max-height: 92vh;
  flex-direction: column;
  margin: 0;
  overflow: hidden;
  border-radius: 18px;

  .el-dialog__header {
    margin: 0;
    padding: 18px 24px;
    border-bottom: 1px solid #dbeafe;
    background: linear-gradient(135deg, #eff6ff 0%, #ffffff 64%, #ecfdf5 100%);
  }

  .el-dialog__title {
    color: #0f172a;
    font-size: 20px;
    font-weight: 900;
  }

  .el-dialog__body {
    min-height: 0;
    padding: 18px 20px;
    overflow-y: auto;
    background: #f8fafc;
  }

  .el-dialog__footer {
    padding: 13px 20px 16px;
    border-top: 1px solid #e2e8f0;
  }

  .el-table th.el-table__cell {
    background: #eaf2ff;
    color: #334155;
    font-weight: 900;
  }

  .el-table td.el-table__cell {
    color: #0f172a;
    font-weight: 650;
  }
}

.inventory-age-table-wrap {
  overflow: hidden;
  border: 1px solid #dbe2ea;
  border-radius: 12px;
  background: #ffffff;
}

.age-badge {
  display: inline-flex;
  min-height: 30px;
  align-items: center;
  padding: 5px 10px;
  border: 1px solid color-mix(in srgb, var(--age-color) 55%, transparent);
  border-radius: 7px;
  background: var(--age-soft-color);
  color: var(--age-color);
  font-size: 12px;
  font-weight: 900;
  white-space: nowrap;
}

.inventory-pagination {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
}

.empty-state {
  padding: 24px;
  color: #64748b;
  font-weight: 700;
}

@media (max-width: 720px) {
  .inventory-pagination {
    justify-content: center;
    overflow-x: auto;
  }
}
</style>
