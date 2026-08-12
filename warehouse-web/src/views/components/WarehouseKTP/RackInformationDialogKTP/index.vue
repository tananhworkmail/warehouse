<template>
  <el-dialog
    :model-value="dialogVisible"
    :title="dialogTitle"
    width="min(900px, 96vw)"
    class="rack-detail-dialog"
    align-center
    destroy-on-close
    @close="emitClose"
  >
    <div v-if="rack" class="rack-dialog">
      <div class="rack-name">{{ t("ktp.rackDetail.rack", { rack: rack.rackCode }) }}</div>

      <div class="sensor-panel">
        <div class="sensor-title">
          <strong>{{ sensorDeviceName || "--" }}</strong>
          <span>{{ sensorZone }}</span>
        </div>

        <div class="gauge-row">
          <el-progress
            type="dashboard"
            :percentage="sensorPercentage('Tem')"
            :color="warehouseTemperatureColor"
            :width="112"
          >
            <template #default="{ percentage }">
              <span class="percentage-value">{{ percentage }}°C</span>
              <span class="percentage-label">{{ t("ktp.rackDetail.temperature") }}</span>
            </template>
          </el-progress>

          <el-progress
            type="dashboard"
            :percentage="sensorPercentage('Hum')"
            :color="warehouseHumidityColor"
            :width="112"
          >
            <template #default="{ percentage }">
              <span class="percentage-value">{{ percentage }}%</span>
              <span class="percentage-label">{{ t("ktp.rackDetail.humidity") }}</span>
            </template>
          </el-progress>
        </div>
      </div>

      <div class="metric-row">
        <div
          class="capacity-card"
          :class="{ 'is-over-capacity': rackUsagePercentage > 100 }"
        >
          <div class="capacity-card-heading">
            <span>{{ t("ktp.rackDetail.usage") }}</span>
            <strong>{{ rackUsagePercentage }}%</strong>
          </div>
          <div class="capacity-progress" aria-hidden="true">
            <span :style="capacityProgressStyle"></span>
          </div>
          <div class="capacity-values">
            <div>
              <span>{{ t("ktp.rackDetail.totalQuantity") }}</span>
              <strong>{{ formatNumber(rack.totalQty) }} {{ t("ktp.capacity.pairs") }}</strong>
            </div>
            <div>
              <span>{{ t("ktp.rackDetail.capacity") }}</span>
              <strong>{{ formatNumber(rack.capacity) }} {{ t("ktp.capacity.pairs") }}</strong>
            </div>
          </div>
        </div>
        <div class="metric-card">
          <span>{{ t("ktp.rackDetail.cartonCount") }}</span>
          <strong>{{ formatNumber(rack.codebarCount) }}</strong>
        </div>
      </div>

      <div v-loading="detailLoading" class="quantity-panel">
        <div class="quantity-title">
          <span>{{ t("ktp.rackDetail.orderDetail") }}</span>
          <strong>
            {{ t("ktp.rackDetail.orderCount", { count: formatNumber(orderDetails.length) }) }}
          </strong>
        </div>

        <div v-if="orderDetails.length" class="order-list">
          <article
            v-for="order in orderDetails"
            :key="order.ddbh"
            class="order-item"
            :class="{ 'is-expanded': selectedOrderDDBH === order.ddbh }"
          >
            <button
              type="button"
              class="order-overview"
              :aria-expanded="selectedOrderDDBH === order.ddbh"
              @click="toggleOrderDetail(order.ddbh)"
            >
              <div class="order-identity">
                <span>{{ t("ktp.rackDetail.order") }}</span>
                <strong>{{ order.ddbh || "--" }}</strong>
              </div>
              <span class="order-toggle" aria-hidden="true">
                {{ selectedOrderDDBH === order.ddbh ? "−" : "+" }}
              </span>
            </button>

            <div
              v-if="selectedOrderDDBH === order.ddbh"
              class="order-detail-content"
            >
              <div class="order-totals">
                <span>
                  {{ t("ktp.rackDetail.totalQuantity") }}:
                  <strong>{{ formatNumber(order.totalQty) }}</strong>
                </span>
                <span>
                  {{ t("ktp.rackDetail.cartonCount") }}:
                  <strong>{{ formatNumber(order.codebarCount) }}</strong>
                </span>
              </div>

              <div class="quantity-grid">
                <button
                  v-for="status in order.statuses"
                  :key="`${order.ddbh}-${status.sb}`"
                  type="button"
                  class="quantity-card"
                  :class="[
                    statusClass(status.sb),
                    {
                      'is-active':
                        statusDialogVisible &&
                        activeOrderDDBH === order.ddbh &&
                        activeSB === String(status.sb),
                    },
                  ]"
                  :aria-pressed="
                    statusDialogVisible &&
                    activeOrderDDBH === order.ddbh &&
                    activeSB === String(status.sb)
                  "
                  @click="openStatusDetail(order.ddbh, status.sb)"
                >
                  <span>{{ statusLabels[String(status.sb)] || status.sb }}</span>
                  <strong>
                    {{ t("ktp.rackDetail.cartonUnit", { count: formatNumber(status.cartonCount) }) }}
                  </strong>
                </button>
              </div>
            </div>
          </article>
        </div>
        <div v-else-if="!detailLoading" class="carton-empty">
          {{ t("ktp.rackDetail.noOrders") }}
        </div>
      </div>

    </div>

    <template #footer>
      <el-button type="primary" @click="emitClose">{{ t("ktp.rackDetail.close") }}</el-button>
    </template>
  </el-dialog>

  <el-dialog
    v-model="statusDialogVisible"
    :title="statusDialogTitle"
    width="min(780px, 94vw)"
    class="status-detail-dialog"
    align-center
    append-to-body
    destroy-on-close
    @closed="resetActiveStatus"
  >
    <div v-loading="detailLoading" class="carton-detail">
      <div class="carton-summary">
        <div>
          <span>{{ t("ktp.rackDetail.cartonCount") }}</span>
          <strong>{{ formatNumber(activeStatusDetail.cartonCount) }}</strong>
        </div>
        <div>
          <span>{{ t("ktp.rackDetail.totalQuantity") }}</span>
          <strong>{{ formatNumber(activeStatusDetail.totalQty) }}</strong>
        </div>
      </div>

      <div v-if="activeStatusDetail.cartons.length" class="carton-table-wrap">
        <table class="carton-table">
          <thead>
            <tr>
              <th>{{ t("ktp.rackDetail.cartonNo") }}</th>
              <th>{{ t("ktp.rackDetail.quantity") }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(carton, index) in activeStatusDetail.cartons"
              :key="`${carton.cartonNo}-${index}`"
            >
              <td>{{ carton.cartonNo || "--" }}</td>
              <td>{{ formatNumber(carton.qty) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else-if="!detailLoading" class="carton-empty">
        {{ t("ktp.rackDetail.noCartons") }}
      </div>
    </div>

    <template #footer>
      <el-button type="primary" @click="statusDialogVisible = false">
        {{ t("ktp.rackDetail.close") }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import axios from "axios";
import { ElMessage } from "element-plus";
import { useI18n } from "@/hooks/i18n";

const API_URL = import.meta.env.VITE_API_URL;
const API_BASE = `${API_URL}/warehouse-ktp`;
const { t, locale } = useI18n();

const props = defineProps({
  dialogVisible: Boolean,
  rack: {
    type: Object,
    default: null,
  },
  sensorDeviceName: {
    type: String,
    default: "",
  },
  sensorData: {
    type: Object,
    default: null,
  },
});

const emit = defineEmits(["update:dialogVisible", "close"]);

const activeSB = ref("");
const activeOrderDDBH = ref("");
const selectedOrderDDBH = ref("");
const statusDialogVisible = ref(false);
const detailLoading = ref(false);
const orderDetail = ref(null);

const statusLabels = computed(() => ({
  1: t("ktp.rackDetail.inbound"),
  2: t("ktp.rackDetail.recycle"),
  4: t("ktp.rackDetail.inspection"),
}));

const dialogTitle = computed(() => t("ktp.rackDetail.title"));

const rackUsagePercentage = computed(() => {
  const capacity = Number(props.rack?.capacity || 0);
  if (capacity <= 0) return 0;
  return Math.max(Math.round((Number(props.rack?.totalQty || 0) / capacity) * 100), 0);
});

const capacityProgressStyle = computed(() => ({
  width: `${Math.min(rackUsagePercentage.value, 100)}%`,
}));

const emitClose = () => {
  statusDialogVisible.value = false;
  emit("update:dialogVisible", false);
  emit("close");
};

const formatNumber = (value) => {
  const numberLocales = { vi: "vi-VN", en: "en-US", zh: "zh-CN" };
  return Number(value || 0).toLocaleString(numberLocales[locale.value] || "vi-VN");
};

const activeStatusLabel = computed(() => statusLabels.value[activeSB.value] || "");
const orderDetails = computed(() =>
  Array.isArray(orderDetail.value?.orders) ? orderDetail.value.orders : [],
);

const activeOrder = computed(() =>
  orderDetails.value.find((order) => order.ddbh === activeOrderDDBH.value),
);

const statusDialogTitle = computed(() => {
  return t("ktp.rackDetail.statusTitle", {
    status: activeStatusLabel.value,
    order: activeOrderDDBH.value || "--",
  });
});

const activeStatusDetail = computed(() => {
  const statuses = Array.isArray(activeOrder.value?.statuses)
    ? activeOrder.value.statuses
    : [];

  return (
    statuses.find((item) => String(item.sb) === activeSB.value) || {
      cartonCount: 0,
      totalQty: 0,
      cartons: [],
    }
  );
});

const loadOrderDetail = async () => {
  const rackCode = props.rack?.rackCode;
  if (!rackCode) return;

  if (orderDetail.value?.rackCode === rackCode) {
    return;
  }

  detailLoading.value = true;
  try {
    const res = await axios.get(
      `${API_BASE}/racks/${encodeURIComponent(rackCode)}/order-detail`,
    );
    if (props.rack?.rackCode === rackCode) {
      orderDetail.value = res.data?.data || null;
    }
  } catch (err) {
    console.error("Load warehouse KTP rack order detail failed:", err);
    ElMessage.error(t("ktp.rackDetail.loadError"));
  } finally {
    detailLoading.value = false;
  }
};

const statusClass = (sb) => ({
  "1": "is-inbound",
  "2": "is-recycle",
  "4": "is-inspection",
})[String(sb)] || "";

const toggleOrderDetail = (ddbh) => {
  selectedOrderDDBH.value =
    selectedOrderDDBH.value === ddbh ? "" : ddbh;
  resetActiveStatus();
  statusDialogVisible.value = false;
};

const openStatusDetail = (ddbh, sb) => {
  activeOrderDDBH.value = ddbh;
  activeSB.value = String(sb);
  statusDialogVisible.value = true;
};

const resetActiveStatus = () => {
  activeSB.value = "";
  activeOrderDDBH.value = "";
};

watch(
  () => [props.dialogVisible, props.rack?.rackCode],
  ([isVisible]) => {
    resetActiveStatus();
    selectedOrderDDBH.value = "";
    statusDialogVisible.value = false;
    orderDetail.value = null;
    if (isVisible) {
      loadOrderDetail();
    }
  },
  { immediate: true },
);

const sensorZone = computed(() => {
  if (props.sensorDeviceName === "KTP01") return t("ktp.rackDetail.zoneA");
  if (props.sensorDeviceName === "KTP02") return t("ktp.rackDetail.zoneB");
  return "";
});

const sensorPercentage = (field) => {
  const value = Number(props.sensorData?.[field]);
  return Number.isFinite(value) ? Math.round(value) : 0;
};

const TEMP_MIN = 28;
const TEMP_MAX = 35;
const HUM_MIN = 50;
const HUM_MAX = 60;

const warehouseTemperatureColor = (val) => {
  if (val < TEMP_MIN) return "#f97316";
  if (val <= TEMP_MAX) return "#22c55e";
  return "#ef4444";
};

const warehouseHumidityColor = (val) => {
  if (val < HUM_MIN) return "#f97316";
  if (val <= HUM_MAX) return "#10b981";
  return "#0000f7";
};
</script>

<style scoped lang="scss">
:deep(.rack-detail-dialog),
:deep(.status-detail-dialog) {
  display: flex;
  flex-direction: column;
  margin: 0;
  overflow: hidden;
  border: 1px solid rgb(148 163 184 / 28%);
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 24px 70px rgb(15 23 42 / 24%);

  .el-dialog__header {
    margin: 0;
    padding: 18px 24px;
    border-bottom: 1px solid #dbeafe;
    background: linear-gradient(135deg, #eff6ff 0%, #ffffff 58%, #eef2ff 100%);
  }

  .el-dialog__title {
    color: #0f172a;
    font-size: 20px;
    font-weight: 900;
  }

  .el-dialog__headerbtn {
    top: 8px;
    right: 10px;

    .el-dialog__close {
      color: #475569;
      font-size: 21px;
    }
  }

  .el-dialog__body {
    min-height: 0;
    padding: 22px 24px 24px;
    overflow-y: auto;
    background: #f8fafc;
  }

  .el-dialog__footer {
    padding: 14px 24px 18px;
    border-top: 1px solid #e2e8f0;
    background: #ffffff;
  }

  .el-button--primary {
    min-width: 96px;
    border-radius: 9px;
    font-weight: 800;
  }
}

:deep(.rack-detail-dialog) {
  max-height: 92vh;
}

:deep(.status-detail-dialog) {
  max-height: 90vh;

  .el-dialog__body {
    min-height: 250px;
  }
}

.rack-dialog {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.rack-name {
  color: #0f172a;
  font-size: 34px;
  font-weight: 900;
  line-height: 1.1;
  text-align: center;
}

.metric-row {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(150px, 1fr);
  gap: 10px;
}

.capacity-card {
  padding: 12px 14px;
  border: 1px solid #2f855a;
  border-radius: 12px;
  background: #f0fdfa;
}

.capacity-card-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;

  span {
    color: #475569;
    font-weight: 800;
  }

  strong {
    color: #1f2937;
    font-size: 22px;
    font-weight: 900;
  }
}

.capacity-progress {
  height: 12px;
  margin: 10px 0;
  overflow: hidden;
  border-radius: 999px;
  background: #00ff73;
  box-shadow: inset 0 1px 2px rgba(15, 23, 42, 0.14);

  span {
    display: block;
    height: 100%;
    border-radius: inherit;
    background: #38bdf8;
    transition: width 0.25s ease;
  }
}

.capacity-values {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;

  div {
    display: flex;
    min-width: 0;
    flex-direction: column;
    gap: 4px;
  }

  span {
    color: #64748b;
    font-size: 12px;
    font-weight: 700;
  }

  strong {
    color: #0f172a;
    font-size: 17px;
    font-weight: 900;
    white-space: nowrap;
  }
}

.capacity-card.is-over-capacity {
  border-color: #fca5a5;
  background: #fff1f2;

  .capacity-card-heading strong,
  .capacity-values > div:first-child strong {
    color: #b91c1c;
  }

  .capacity-progress {
    box-shadow: inset 0 0 0 2px #dc2626;
  }
}

.metric-card {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 14px;
  border: 1px solid #d5dde8;
  border-radius: 12px;
  background: #f8fafc;
  color: #475569;

  span {
    font-weight: 700;
  }

  strong {
    color: #0f172a;
    font-size: 20px;
    font-weight: 900;
  }
}

.quantity-panel {
  padding: 14px;
  border: 1px solid #d5dde8;
  border-radius: 12px;
  background: #ffffff;
}

.quantity-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
  color: #334155;
  font-size: 14px;
  font-weight: 900;

  strong {
    padding: 4px 9px;
    border-radius: 999px;
    background: #e0f2fe;
    color: #0369a1;
    font-size: 12px;
    white-space: nowrap;
  }
}

.order-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.order-item {
  border: 1px solid #dbeafe;
  border-radius: 12px;
  background: #f8fbff;
  overflow: hidden;

  &.is-expanded {
    border-color: #93c5fd;
    box-shadow: 0 0 0 2px rgb(59 130 246 / 8%);
  }
}

.order-overview {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  padding: 12px;
  border: 0;
  background: transparent;
  cursor: pointer;
  font: inherit;
  text-align: left;
  transition: background-color 0.15s ease;

  &:hover {
    background: #eff6ff;
  }
}

.order-identity {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 3px;

  span {
    color: #64748b;
    font-size: 11px;
    font-weight: 800;
    text-transform: uppercase;
  }

  strong {
    color: #1d4ed8;
    font-size: 17px;
    font-weight: 900;
    overflow-wrap: anywhere;
  }
}

.order-toggle {
  display: flex;
  width: 28px;
  height: 28px;
  flex: 0 0 28px;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: #dbeafe;
  color: #1d4ed8;
  font-size: 20px;
  font-weight: 900;
  line-height: 1;
}

.order-detail-content {
  padding: 0 12px 12px;
  border-top: 1px solid #dbeafe;
  background: #ffffff;
}

.order-totals {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  gap: 6px 14px;
  padding: 10px 0;
  color: #64748b;
  font-size: 12px;
  font-weight: 700;
  text-align: left;

  strong {
    color: #0f172a;
    font-weight: 900;
  }
}

.quantity-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.quantity-card {
  display: flex;
  min-width: 0;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  padding: 11px;
  border: 1px solid transparent;
  border-radius: 12px;
  cursor: pointer;
  font: inherit;
  text-align: left;
  transition:
    border-color 0.15s ease,
    box-shadow 0.15s ease,
    transform 0.15s ease;

  &:hover {
    transform: translateY(-1px);
  }

  &.is-active {
    border-color: #2563eb;
    box-shadow: 0 0 0 2px rgb(37 99 235 / 12%);
  }

  span {
    color: #475569;
    font-size: 13px;
    font-weight: 700;
  }

  strong {
    color: #0f172a;
    font-size: 19px;
    font-weight: 900;
  }

  &.is-inbound {
    background: #ecfdf5;
  }

  &.is-recycle {
    background: #fff7ed;
  }

  &.is-inspection {
    background: #fefce8;
  }
}

.carton-detail {
  min-height: 220px;
}

.carton-summary {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
  margin-bottom: 18px;

  div {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    padding: 16px 18px;
    border: 1px solid #dbeafe;
    border-radius: 12px;
    background: linear-gradient(135deg, #ffffff, #eff6ff);
  }

  span {
    color: #64748b;
    font-size: 13px;
    font-weight: 700;
  }

  strong {
    color: #0f172a;
    font-size: 22px;
    font-weight: 900;
  }
}

.carton-table-wrap {
  max-height: 390px;
  overflow: auto;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #ffffff;
}

.carton-table {
  width: 100%;
  border-collapse: collapse;

  th,
  td {
    padding: 12px 16px;
    border-bottom: 1px solid #e2e8f0;
  }

  th {
    position: sticky;
    top: 0;
    z-index: 1;
    background: #eaf2ff;
    color: #475569;
    font-size: 13px;
    font-weight: 900;
    text-align: left;
  }

  td {
    color: #0f172a;
    font-size: 14px;
    font-weight: 700;
  }

  th:last-child,
  td:last-child {
    text-align: right;
  }

  tbody tr:last-child td {
    border-bottom: 0;
  }

  tbody tr:nth-child(even) {
    background: #f8fafc;
  }

  tbody tr:hover {
    background: #eff6ff;
  }
}

.carton-empty {
  padding: 18px;
  border-radius: 12px;
  background: #f8fafc;
  color: #64748b;
  font-size: 13px;
  font-weight: 700;
  text-align: center;
}

.sensor-panel {
  padding: 14px;
  border: 1px solid #d5dde8;
  border-radius: 14px;
  background: #ffffff;
}

.sensor-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;

  strong {
    color: #b91c1c;
    font-size: 18px;
    font-weight: 900;
  }

  span {
    color: #64748b;
    font-size: 13px;
    font-weight: 800;
  }
}

.gauge-row {
  display: flex;
  justify-content: center;
  gap: 20px;
}

@media (max-width: 520px) {
  .metric-row {
    grid-template-columns: 1fr;
  }

  .order-totals {
    justify-content: flex-start;
    text-align: left;
  }
}

.percentage-value {
  display: block;
  color: #0f172a;
  font-size: 18px;
  font-weight: 900;
}

.percentage-label {
  display: block;
  color: #64748b;
  font-size: 12px;
  font-weight: 800;
}
</style>
