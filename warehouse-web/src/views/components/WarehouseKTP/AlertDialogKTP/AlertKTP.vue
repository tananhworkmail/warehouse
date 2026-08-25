<template>
  <el-dialog
    v-model="dialogVisible"
    :title="t('ktp.alert.title')"
    width="95%"
    top="4vh"
  >
    <div class="toolbar">
      <div class="filter-group">
        <span class="filter-label">{{ t("ktp.alert.date") }}</span>
        <el-date-picker
          v-model="selectedDate"
          type="date"
          size="default"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          style="width: 160px"
          @change="onLogFilterChange"
        />
      </div>

      <div style="display: none;" class="filter-group" hidden >
        <span class="filter-label">{{ t("ktp.alert.timeRange") }}</span>
        <el-time-picker
          v-model="timeRange"
          is-range
          size="default"
          range-separator="-"
          format="HH:mm"
          value-format="HH:mm"
          style="width: 220px"
          @change="onTimeChange"
        />
      </div>

      <div class="toolbar-badge">{{ selectedDate }}</div>
    </div>

    <div class="data-container">
      <div class="andon-row">
        <div
          v-for="kv in andonDevices"
          :key="kv.name"
          class="andon-card"
          :class="'mobile-order-' + kv.name"
          :style="cardStyle(kv.name)"
        >
          <div class="andon-left">
            <div class="andon-top">
              <span class="andon-name">
                <span :style="{ color: zoneColors[kv.name] }">
                  {{ kv.name }}
                </span>
              </span>
            </div>

            <span class="detail-btn" @click="loadZone(kv.name)">
              {{ t("ktp.alert.rackDetail") }}
            </span>
          </div>
          <div class="gauge-row">
            <el-progress
              type="dashboard"
              :percentage="Math.round(tempHumidityData[kv.name]?.Tem || 0)"
              :color="warehouseTemperatureColor"
              :width="80"
            >
              <template #default="{ percentage }">
                <span class="percentage-value">{{ percentage }}°C</span>
                <span class="percentage-label">{{ t("ktp.alert.temp") }}</span>
              </template>
            </el-progress>

            <el-progress
              type="dashboard"
              :percentage="Math.round(tempHumidityData[kv.name]?.Hum || 0)"
              :color="warehouseHumidityColor"
              :width="80"
            >
              <template #default="{ percentage }">
                <span class="percentage-value">{{ percentage }}%</span>
                <span class="percentage-label">{{ t("ktp.alert.hum") }}</span>
              </template>
            </el-progress>
          </div>
        </div>
      </div>

      <div class="tables-row">
        <div class="table-wrap mobile-order-KTP01">
          <el-table
            :data="data01"
            :max-height="tableMaxHeight"
            :row-class-name="rowClassName"
            style="width: 100%"
          >
            <el-table-column :label="t('ktp.alert.humidity')" align="center">
              <template #default="{ row }">
                <span :class="humClass(row.Hum)">{{ row.Hum }}</span>
              </template>
            </el-table-column>
            <el-table-column :label="t('ktp.alert.temp')" align="center">
              <template #default="{ row }">
                <span :class="tempClass(row.Temp)">{{ row.Temp }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="Alert_Time" :label="t('ktp.alert.time')" align="center" />
          </el-table>
        </div>

        <div class="table-wrap mobile-order-KTP02">
          <el-table
            :data="data02"
            :max-height="tableMaxHeight"
            :row-class-name="rowClassName"
            style="width: 100%"
          >
            <el-table-column :label="t('ktp.alert.humidity')" align="center">
              <template #default="{ row }">
                <span :class="humClass(row.Hum)">{{ row.Hum }}</span>
              </template>
            </el-table-column>
            <el-table-column :label="t('ktp.alert.temp')" align="center">
              <template #default="{ row }">
                <span :class="tempClass(row.Temp)">{{ row.Temp }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="Alert_Time" :label="t('ktp.alert.time')" align="center" />
          </el-table>
        </div>
      </div>
    </div>
  </el-dialog>

  <el-dialog
    v-model="detailVisible"
    :title="`${t('ktp.alert.detailTitle')}: ${selectedZone}`"
    width="560px"
  >
    <div class="rack-grid" v-if="racks.length">
      <div v-for="rack in racks" :key="rack" class="rack-item">{{ rack }}</div>
    </div>
    <div v-else class="rack-empty">{{ t("ktp.alert.noRacks") }}</div>
  </el-dialog>
</template>

<script setup>
import { ref, watch, reactive, onUnmounted } from "vue";
import axios from "axios";
import dayjs from "dayjs";
import { ktpRackCodes } from "@/views/components/WarehouseKTP/ktpLayout";
import { useWarehouseMapI18n } from "@/hooks/useWarehouseMapI18n";

const API_URL = import.meta.env.VITE_API_URL;
const API_BASE = `${API_URL}/warehouse-ktp`;
const { t } = useWarehouseMapI18n();
const props = defineProps({ modelValue: Boolean });
const emit = defineEmits(["update:modelValue"]);
const dialogVisible = ref(false);
const detailVisible = ref(false);
const selectedZone = ref("");
const racks = ref([]);
const tableMaxHeight = 420;

const selectedDate = ref(dayjs().format("YYYY-MM-DD"));
const timeRange = ref([]);
const data01 = ref([]);
const data02 = ref([]);

const andonDevices = reactive([
  { name: "KTP01" },
  { name: "KTP02" },
]);
const zoneColors = {
  KTP01: "#0ea5e9",
  KTP02: "#6366f1",
};
const TEMP_MIN = 28;
const TEMP_MAX = 35;
const HUM_MIN = 50;
const HUM_MAX = 60;

const cardStyle = (name) => {
  const color = zoneColors[name] || "#ccc";
  return {
    borderLeft: `4px solid ${color}`,
    borderRadius: "0 8px 8px 0",
  };
};
const warehouseTemperatureColor = (val) => {
  if (val < TEMP_MIN) return "#f97316";
  if (val <= TEMP_MAX) return "#22c55e";
  return "#ef4444";
};

const warehouseHumidityColor = (val) => {
  if (val < HUM_MIN) return "#f97316";
  if (val <= HUM_MAX) return "#10b981";
  return "#0000F7";
};

watch(dialogVisible, (val) => emit("update:modelValue", val));

const loadZone = (zone) => {
  const prefix = zone === "KTP01" ? "A" : "B";
  racks.value = ktpRackCodes
    .filter((rack) => rack.startsWith(prefix))
    .sort((a, b) => Number(a.slice(1)) - Number(b.slice(1)));
  selectedZone.value = zone;
  detailVisible.value = true;
};

const parseDisplayNumber = (value) =>
  parseFloat(String(value).replace("%", ""));
const rowClassName = ({ row }) =>
  parseDisplayNumber(row.Temp) > TEMP_MAX ||
  parseDisplayNumber(row.Hum) > HUM_MAX
    ? "row-high"
    : parseDisplayNumber(row.Temp) < TEMP_MIN ||
      parseDisplayNumber(row.Hum) < HUM_MIN
    ? "row-low"
    : "";
const humClass = (hum) =>
  parseDisplayNumber(hum) > HUM_MAX
    ? "hum-high"
    : parseDisplayNumber(hum) < HUM_MIN
    ? "hum-low"
    : "";
const tempClass = (temp) =>
  parseDisplayNumber(temp) > TEMP_MAX
    ? "temp-high"
    : parseDisplayNumber(temp) < TEMP_MIN
    ? "temp-low"
    : "";
const onTimeChange = () => {
  onLogFilterChange();
};

const onLogFilterChange = () => {
  if (dialogVisible.value) fetchLatestLogByDevices();
};

const tempHumidityData = ref({});
const formatRecordTime = (value) => {
  if (!value) return "";
  const parsed = dayjs(value);
  if (parsed.isValid()) return parsed.format("YYYY-MM-DD HH:mm:ss");
  return String(value).replace("T", " ").slice(0, 19);
};

const formatOneDecimal = (value) => {
  const num = Number(value);
  return Number.isFinite(num) ? num.toFixed(1) : "";
};

const toLogTableRow = (item) => {
  return {
    Hum: `${formatOneDecimal(item.Hum)}%`,
    Temp: `${formatOneDecimal(item.Temp ?? item.Tem)}°C`,
    Alert_Time: formatRecordTime(item.RecordTime),
  };
};

const fetchTempHumidityByDevices = async () => {
  try {
    const res = await axios.get(`${API_BASE}/temphumiditybydevices`);
    const raw = res.data?.data;
    let list = [];

    if (Array.isArray(raw)) {
      list = raw;
    } else if (raw && Array.isArray(raw.data)) {
      list = raw.data;
    } else if (raw && typeof raw === "object") {
      list = Object.values(raw);
    }

    const map = {};
    (list || []).forEach((item) => {
      if (item?.DeviceName) {
        map[item.DeviceName] = item;
      }
    });

    tempHumidityData.value = map;
  } catch (e) {
    console.error(e);
  }
};
let tempHumidityTimer = null;

const fetchLatestLogByDevices = async () => {
  try {
    data01.value = [];
    data02.value = [];

    const params = {
      date: selectedDate.value || dayjs().format("YYYY-MM-DD"),
      startTime: timeRange.value?.[0] || "",
      endTime: timeRange.value?.[1] || "",
    };

    const res = await axios.get(
      `${API_BASE}/temphumiditylatestlogbydevices`,
      { params },
    );

    const raw = res.data?.data;
    let list = [];

    if (Array.isArray(raw)) {
      list = raw;
    } else if (raw && Array.isArray(raw.data)) {
      list = raw.data;
    } else if (raw && typeof raw === "object") {
      list = Object.values(raw);
    }

    const rows01 = [];
    const rows02 = [];

    (list || []).forEach((item) => {
      if (item?.DeviceName === "KTP01") {
        rows01.push(toLogTableRow(item));
      }
      if (item?.DeviceName === "KTP02") {
        rows02.push(toLogTableRow(item));
      }
    });

    data01.value = rows01;
    data02.value = rows02;
  } catch (e) {
    console.error(e);
  }
};

let latestLogTimer = null;

const scheduleLatestLogReload = () => {
  clearTimeout(latestLogTimer);

  latestLogTimer = setTimeout(async () => {
    await fetchLatestLogByDevices();
    scheduleLatestLogReload();
  }, 5 * 60 * 1000);
};

const scheduleTempHumidityReload = () => {
  clearTimeout(tempHumidityTimer);
  tempHumidityTimer = setTimeout(async () => {
    await fetchTempHumidityByDevices();
    scheduleTempHumidityReload();
  }, 5 * 60 * 1000);
};

watch(
  () => props.modelValue,
  (val) => {
    dialogVisible.value = val;
    if (val) {
      fetchTempHumidityByDevices();
      fetchLatestLogByDevices();
      scheduleTempHumidityReload();
      scheduleLatestLogReload();
    } else {
      clearTimeout(tempHumidityTimer);
      clearTimeout(latestLogTimer);
    }
  },
  { immediate: true },
);

onUnmounted(() => {
  clearTimeout(tempHumidityTimer);
  clearTimeout(latestLogTimer);
});
</script>
<style scoped>
/* ?? DIALOG SHELL ?? */
:deep(.el-dialog) {
  height: 92vh;
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
}
:deep(.el-dialog__header) {
  padding: 14px 20px;
  border-bottom: 0.5px solid #e2e8f0;
  margin: 0;
}
:deep(.el-dialog__title) {
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
  letter-spacing: 0.01em;
}
:deep(.el-dialog__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  overflow-y: auto;
  padding: 14px 16px;
  background: #f8fafc;
  gap: 10px;
}

/* ?? TOOLBAR ?? */
.toolbar {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  padding: 10px 14px;
  background: #fff;
  border: 0.5px solid #e2e8f0;
  border-radius: 8px;
  flex-shrink: 0;
  flex-wrap: wrap;
}
.filter-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.filter-label {
  font-size: 11px;
  font-weight: 600;
  color: #94a3b8;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}
.toolbar-badge {
  margin-left: auto;
  font-size: 12px;
  font-weight: 600;
  color: #2563eb;
  background: #eff6ff;
  border: 0.5px solid #bfdbfe;
  padding: 4px 12px;
  border-radius: 20px;
  align-self: center;
}

/* CONTAINER CHO RESPONSIVE KHI VIEW PC */
.data-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

/* ?? ANDON ROW ?? */
.andon-row {
  display: flex;
  gap: 10px;
  flex-shrink: 0;
}
.andon-card {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 14px;
  background: #fff;
  border: 0.5px solid #e2e8f0;
  border-radius: 8px;
  min-height: 54px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.andon-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.andon-top {
  display: flex;
  align-items: center;
  gap: 6px;
}
.andon-name-wrap {
  display: inline-block;
}
.andon-name {
  font-size: 24px;
  font-weight: 600;
  display: inline-block;
  width: fit-content;
}
.andon-wrap {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.andon-badge {
  position: static;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 999px;
  white-space: nowrap;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
}
.badge-alert {
  background: #fef2f2;
  color: #b91c1c;
  border: 0.5px solid #fecaca;
}
.badge-ok {
  background: #f0fdf4;
  color: #15803d;
  border: 0.5px solid #bbf7d0;
}
.badge-stale {
  background: #fffbeb;
  color: #92400e;
  border: 0.5px solid #fde68a;
}

.btn-off {
  position: relative;
  z-index: 1;
  font-size: 16px !important;
  color: #64748b !important;
  background: #f8fafc !important;
  border: 0.5px solid #cbd5e1 !important;
  border-radius: 6px !important;
  padding: 12px 20px !important;
  height: auto !important;
  flex-shrink: 0;
  display: inline-flex !important;
  align-items: center !important;
  gap: 4px !important;
}
.btn-off:hover {
  background: #f0f5ff !important;
  color: #dc2626 !important;
  border-color: #fecaca !important;
}
.btn-off-icon {
  width: 22px;
  height: 22px;
  margin-right: 6px;
  color: #ef4444;
  flex-shrink: 0;
}

/* ?? TABLES ?? */
.tables-row {
  display: flex;
  gap: 10px;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}
.table-wrap {
  flex: 1;
  display: flex;
  flex-direction: column;
  border-radius: 8px;
  overflow: hidden;
  border: none;
}

.table-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  flex-shrink: 0;
}

.zone-title {
  font-size: 14px;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 0.03em;
}

.detail-btn {
  display: inline-flex;
  align-items: center;
  width: fit-content;
  font-size: 12px;
  font-weight: 600;
  color: #000000;
  padding: 3px 10px;
  border-radius: 6px;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  cursor: pointer;
  transition: all 0.15s ease;
}
.detail-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

:deep(.el-table) {
  border: none !important;
}
:deep(.el-table thead th.el-table__cell) {
  font-size: 15px !important;
  font-weight: 800 !important;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  padding: 10px 0 !important;
  color: #fff !important;
}

.table-wrap:nth-child(1) :deep(.el-table thead th.el-table__cell) {
  background: #0ea5e9 !important;
}
.table-wrap:nth-child(2) :deep(.el-table thead th.el-table__cell) {
  background: #6366f1 !important;
}
.table-wrap:nth-child(3) :deep(.el-table thead th.el-table__cell) {
  background: #7500fc !important;
}
:deep(.el-table td.el-table__cell) {
  font-size: 21px;
  font-weight: 800;
  color: #334155;
  border-bottom: 2px solid #e2e8f0 !important;
  padding: 14px 0 !important;
  background: #fff;
  line-height: 1.2;
}
:deep(.el-table__body tr:nth-child(even) td.el-table__cell) {
  background: #f8fafc;
}
:deep(.el-table__body tr:nth-child(odd) td.el-table__cell) {
  background: #ffffff;
}
:deep(.el-table td.el-table__cell:nth-child(3)) {
  font-size: 16px;
  font-weight: 700;
  color: #64748b;
  white-space: nowrap;
}
:deep(.el-table tr:hover td.el-table__cell) {
  background: #eef6ff !important;
}
:deep(.row-high td) {
  background: #0e9fff2d !important;
}
:deep(.row-high td:first-child) {
  border-left: 5px solid #1d4ed8 !important;
}
:deep(.row-low td) {
  background: #fffdf0 !important;
}
:deep(.row-low td:first-child) {
  border-left: 5px solid #f97316 !important;
}

.hum-high {
  color: #1d4ed8;
  font-weight: 900;
}
.hum-low {
  color: #b45309;
  font-weight: 900;
}

.temp-high {
  color: #dc2626;
  font-weight: 900;
}

.temp-low {
  color: #b45309;
  font-weight: 900;
}

/* ?? RACK GRID ?? */
.rack-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(72px, 1fr));
  gap: 8px;
  padding: 4px;
}
.rack-item {
  background: #eff6ff;
  color: #1d4ed8;
  border: 0.5px solid #bfdbfe;
  padding: 9px 6px;
  text-align: center;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: background 0.15s;
}
.rack-item:hover {
  background: #dbeafe;
}
.rack-empty {
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
  padding: 32px 0;
}
.gauge-summary {
  display: flex;
  gap: 12px;
  margin-bottom: 10px;
}
.gauge-card {
  flex: 1;
  background: #fff;
  border: 0.5px solid #e2e8f0;
  border-radius: 8px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.andon-gauge {
  display: flex;
  align-items: center;
  gap: 10px;
}
.gauge-title {
  font-size: 12px;
  font-weight: 700;
  color: #334155;
  margin-bottom: 6px;
}
.gauge-row {
  display: flex;
  gap: 12px;
  align-items: center;
  justify-content: center;
}

@keyframes blink {
  0% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.4;
    transform: scale(1.1);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

.blinking {
  animation: blink 0.8s infinite;
  border-color: #ef4444 !important;
  color: #dc2626 !important;
}

/* ?? GIAO DI廙 ?I廙 THO廕 (A -> Table A -> B -> Table B) ?? */
@media screen and (max-width: 768px) {
  :deep(.el-dialog) {
    width: 100% !important;
    height: 100vh !important;
    margin-top: 0 !important;
    border-radius: 0;
  }

  :deep(.el-dialog__body) {
    padding: 10px;
    overflow-y: auto;
  }

  .toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  .filter-group {
    width: 100%;
  }
  .toolbar .el-date-editor {
    width: 100% !important;
    max-width: 100%;
  }
  .toolbar-badge {
    margin-left: 0;
    align-self: flex-start;
  }

  /* B? QUY廕鋁 X廕銷 XEN K廕?*/
  .data-container {
    display: flex;
    flex-direction: column;
    gap: 10px;
    overflow-y: visible; /* ?廕σ b廕υ cu廙 ?廙θ tr礙n mobile */
  }

  .andon-row,
  .tables-row {
    display: contents;
  }

  .mobile-order-KTP01 {
    order: 1;
  }
  .mobile-order-KTP02 {
    order: 2;
  }

  .andon-card {
    margin-top: 15px;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 16px;
    padding: 16px;
  }
  .andon-left {
    align-items: center;
  }
  .andon-wrap {
    width: 100%;
  }
  .btn-off {
    width: 100%;
    justify-content: center;
  }

  .table-wrap {
    min-height: 250px;
  }

  :deep(.el-table td.el-table__cell) {
    font-size: 16px;
    padding: 10px 0 !important;
  }
  :deep(.el-table thead th.el-table__cell) {
    font-size: 13px !important;
  }
}
</style>
