<template>
  <WarehouseLayoutKTP>
    <template #default>
      <section class="ktp-page">
        <div ref="scrollRef" v-loading="loading" class="warehouse-scroll">
          <div
            class="warehouse-board"
            :style="{ width: `${ktpLayout.width}px`, height: `${ktpLayout.height}px` }"
          >
            <svg
              class="layout-layer"
              :viewBox="`0 0 ${ktpLayout.width} ${ktpLayout.height}`"
              aria-hidden="true"
            >
              <template v-for="shape in ktpLayout.shapes" :key="shape.id">
                <line
                  v-if="shape.type === 'line'"
                  class="wall-line"
                  :x1="shape.x1"
                  :y1="shape.y1"
                  :x2="shape.x2"
                  :y2="shape.y2"
                />
                <g v-else-if="shape.type === 'rect'">
                  <rect
                    :class="['shape-rect', `is-${shape.variant}`]"
                    :x="shape.x"
                    :y="shape.y"
                    :width="shape.w"
                    :height="shape.h"
                  />
                  <text
                    v-if="shape.text"
                    class="shape-text"
                    :x="shape.x + shape.w / 2"
                    :y="shape.y + shape.h / 2 - ((shape.text.split('\n').length - 1) * 7)"
                    text-anchor="middle"
                  >
                    <tspan
                      v-for="(line, index) in shape.text.split('\n')"
                      :key="`${shape.id}-${index}`"
                      :x="shape.x + shape.w / 2"
                      :dy="index === 0 ? 0 : 14"
                    >
                      {{ line }}
                    </tspan>
                  </text>
                </g>
                <circle
                  v-else-if="shape.type === 'circle'"
                  :class="['shape-circle', `is-${shape.variant}`]"
                  :cx="shape.cx"
                  :cy="shape.cy"
                  :r="shape.r"
                />
                <g v-else-if="shape.type === 'door'" class="door-shape">
                  <rect
                    :x="shape.x"
                    :y="doorPanelY(shape)"
                    :width="shape.size"
                    :height="doorPanelThickness"
                  />
                  <line
                    :x1="shape.x + 4"
                    :y1="doorPanelY(shape) + doorPanelThickness / 2"
                    :x2="shape.x + shape.size - 4"
                    :y2="doorPanelY(shape) + doorPanelThickness / 2"
                  />
                </g>
                <g v-else-if="shape.type === 'stairs'" class="stairs-shape">
                  <rect
                    :x="shape.x"
                    :y="shape.y"
                    :width="shape.w"
                    :height="shape.h"
                    rx="2"
                  />
                  <line
                    v-for="step in shape.steps"
                    :key="`${shape.id}-${step}`"
                    :x1="shape.x"
                    :y1="shape.y + (shape.h / shape.steps) * step"
                    :x2="shape.x + shape.w"
                    :y2="shape.y + (shape.h / shape.steps) * step"
                  />
                  <polyline :points="stairsStepPoints(shape)" />
                </g>
                <text
                  v-else-if="shape.type === 'text'"
                  :class="['layout-text', `is-${shape.variant}`]"
                  :x="shape.x"
                  :y="shape.y"
                  :text-anchor="shape.anchor || 'middle'"
                >
                  <tspan
                    v-for="(line, index) in shape.value.split('\n')"
                    :key="`${shape.id}-${index}`"
                    :x="shape.x"
                    :dy="index === 0 ? 0 : 16"
                  >
                    {{ line }}
                  </tspan>
                </text>
              </template>

              

              <polyline
                v-for="route in ktpLayout.escapeRoutes"
                :key="route.id"
                class="escape-route"
                :points="routeBodyPoints(route.points)"
              />
              <polygon
                v-for="route in ktpLayout.escapeRoutes"
                :key="`${route.id}-head`"
                class="escape-arrow-head"
                :points="routeArrowHeadPoints(route.points)"
              />
            </svg>

            <div
              class="capacity-summary"
              role="group"
              :aria-label="t('ktp.capacity.summary')"
            >
              <div
                class="capacity-chart"
                :style="capacityChartStyle"
                :title="`${warehouseUsagePercent}%`"
              >
                <div class="capacity-chart-center">
                  <strong>{{ warehouseUsagePercent }}%</strong>
                  <span>{{ t("ktp.capacity.usedShort") }}</span>
                </div>
              </div>

              <div class="capacity-metrics">
                <div class="capacity-metric is-total">
                  <span>{{ t("ktp.capacity.total") }}</span>
                  <strong>{{ formatNumber(warehouseCapacity) }} {{ t("ktp.capacity.pairs") }}</strong>
                </div>
                <div class="capacity-metric is-used">
                  <span>{{ t("ktp.capacity.used") }}</span>
                  <strong>{{ formatNumber(warehouseUsed) }} {{ t("ktp.capacity.pairs") }}</strong>
                </div>
                <div class="capacity-metric is-remaining">
                  <span>{{ t("ktp.capacity.remaining") }}</span>
                  <strong>{{ formatNumber(warehouseRemaining) }} {{ t("ktp.capacity.pairs") }}</strong>
                </div>
              </div>
            </div>

            <div
              v-for="sensor in tempHumiditySensors"
              :key="sensor.deviceName"
              class="sensor-hotspot"
              :style="sensorHotspotStyle(sensor)"
            >
              <span class="sensor-dot"></span>
              <div class="sensor-card">
              <span class="sensor-name">{{ sensor.deviceName }}</span>
              <span class="sensor-value">{{ formatSensorValue(sensor.deviceName, "Tem", "°C") }}</span>
              <span class="sensor-value">{{ formatSensorValue(sensor.deviceName, "Hum", "%") }}</span>
            </div>
            </div>

            <button
              v-for="item in rackItems"
              :key="item.id"
              :ref="(el) => setRackRef(item.id, el)"
              type="button"
              class="rack-cell"
              :class="[
                rackLabelClass(item.code),
                {
                  'is-horizontal': item.direction === 'horizontal',
                  'is-empty': getRack(item.code).status === 'empty',
                  'is-occupied': getRack(item.code).status === 'occupied',
                  'is-selected': selectedRack?.layoutId === item.id,
                  'is-over-capacity': getRackUsagePercentage(item.code) > 100,
                },
              ]"
              :style="rackStyle(item)"
              @click="openRackDialog(item)"
            >
              <span class="rack-usage-fill" aria-hidden="true"></span>
              <span class="rack-code">{{ item.code }}</span>
              <span class="qty-text">
                {{ getRackUsagePercentage(item.code) }}%
              </span>
            </button>
          </div>
        </div>

        <RackInformationDialogKTP
          v-model:dialog-visible="dialogVisible"
          :rack="selectedRack"
          :sensor-device-name="selectedRack?.sensorDeviceName || ''"
          :sensor-data="getSensorData(selectedRack?.sensorDeviceName)"
          @close="closeRackDialog"
        />

      </section>
    </template>
  </WarehouseLayoutKTP>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from "vue";
import axios from "axios";
import { ElMessage } from "element-plus";
import WarehouseLayoutKTP from "@/views/components/WarehouseKTP/WarehouseLayoutKTP.vue";
import RackInformationDialogKTP from "@/views/components/WarehouseKTP/RackInformationDialogKTP/index.vue";
import { ktpLayout } from "@/views/components/WarehouseKTP/ktpLayout";
import {
  ktpRackCapacityByCode,
  ktpWarehouseCapacity,
} from "@/views/components/WarehouseKTP/ktpCapacity";
import { useI18n } from "@/hooks/i18n";

const API_URL = import.meta.env.VITE_API_URL;
const API_BASE = `${API_URL}/warehouse-ktp`;
const { t, locale } = useI18n();

const loading = ref(false);
const dialogVisible = ref(false);
const selectedRack = ref(null);
const rackDataMap = ref(new Map());
const tempHumidityData = ref({});
const rackRefs = new Map();
let ktpEventSource = null;
let realtimeRefreshTimer = null;
let tempHumidityTimer = null;

const tempHumiditySensors = [
  { deviceName: "KTP01", zone: "Khu A", x: 663, y: 480 },
  { deviceName: "KTP02", zone: "Khu B", x: 674, y: 410 },
];

const rackItems = computed(() =>
  ktpLayout.racks.map((item, index) => ({
    ...item,
    id: `${item.code}-${index}`,
  })),
);

const getRackCapacity = (rackCode) =>
  Number(ktpRackCapacityByCode[rackCode] || 0);

const warehouseCapacity = ktpWarehouseCapacity;
const warehouseUsed = computed(() =>
  ktpLayout.racks.reduce(
    (total, item) => total + Number(getRack(item.code).totalQty || 0),
    0,
  ),
);
const warehouseRemaining = computed(() =>
  Math.max(warehouseCapacity - warehouseUsed.value, 0),
);

const calculateUsagePercentage = (used, capacity) => {
  const numericCapacity = Number(capacity || 0);
  if (numericCapacity <= 0) return 0;
  return Math.max(Math.round((Number(used || 0) / numericCapacity) * 100), 0);
};

const warehouseUsagePercent = computed(() =>
  calculateUsagePercentage(warehouseUsed.value, warehouseCapacity),
);
const capacityChartStyle = computed(() => {
  const chartPercent = Math.min(warehouseUsagePercent.value, 100);
  return {
    background: `conic-gradient(rgb(228, 252, 255) 0deg, rgb(0 126 255) 56%, rgb(0, 255, 115) 56%, rgb(0, 255, 115) 100%)`,
  };
});

const fetchRacks = async ({ silent = false } = {}) => {
  if (!silent) {
    loading.value = true;
  }
  try {
    const res = await axios.get(`${API_BASE}/racks`);
    applyRackData(res.data.data || []);
  } catch (err) {
    const message =
      err.response?.data?.message ||
      err.message ||
      "Khong tai duoc du lieu kho thanh pham";
    console.error("Load warehouse KTP racks failed:", err);
    ElMessage.error(message);
  } finally {
    if (!silent) {
      loading.value = false;
    }
  }
};

const applyRackData = (rows) => {
  const nextMap = new Map();
  rows.forEach((rack) => {
    nextMap.set(rack.rackCode, {
      rackCode: rack.rackCode,
      ddbh: rack.ddbh || "",
      currentCode: null,
      status: Number(rack.totalQty || 0) > 0 ? "occupied" : "empty",
      totalQty: Number(rack.totalQty || 0),
      inboundCartonCount: Number(rack.inboundCartonCount || 0),
      recycleCartonCount: Number(rack.recycleCartonCount || 0),
      inspectionCartonCount: Number(rack.inspectionCartonCount || 0),
      codebarCount: Number(rack.codebarCount || 0),
      sampleCodes: rack.sampleCodes || [],
      capacity: getRackCapacity(rack.rackCode),
    });
  });
  rackDataMap.value = nextMap;

  if (selectedRack.value?.rackCode) {
    const latestRack =
      nextMap.get(selectedRack.value.rackCode) ||
      getEmptyRack(selectedRack.value.rackCode);
    selectedRack.value = {
      ...latestRack,
      layoutId: selectedRack.value.layoutId,
      sensorDeviceName: selectedRack.value.sensorDeviceName,
    };
  }
};

const getEmptyRack = (rackCode) => ({
  rackCode,
  ddbh: "",
  currentCode: null,
  status: "empty",
  totalQty: 0,
  inboundCartonCount: 0,
  recycleCartonCount: 0,
  inspectionCartonCount: 0,
  codebarCount: 0,
  sampleCodes: [],
  capacity: getRackCapacity(rackCode),
});

const getRack = (rackCode) => {
  return rackDataMap.value.get(rackCode) || getEmptyRack(rackCode);
};

const openRackDialog = (item) => {
  const rack = getRack(item.code);
  selectedRack.value = {
    ...rack,
    rackCode: item.code,
    layoutId: item.id,
    sensorDeviceName: getRackSensorDeviceName(item.code),
    capacity: getRackCapacity(item.code),
  };
  dialogVisible.value = true;
};

const closeRackDialog = () => {
  dialogVisible.value = false;
  selectedRack.value = null;
};

const fetchTempHumidityByDevices = async () => {
  try {
    const res = await axios.get(`${API_BASE}/temphumiditybydevices`);
    const raw = res.data?.data;
    const list = Array.isArray(raw) ? raw : [];
    const map = {};

    list.forEach((item) => {
      if (item?.DeviceName) {
        map[item.DeviceName] = item;
      }
    });

    tempHumidityData.value = map;
  } catch (err) {
    console.error("Load warehouse KTP temp humidity failed:", err);
  }
};

const setRackRef = (id, el) => {
  if (el) rackRefs.set(id, el);
  else rackRefs.delete(id);
};

const rackStyle = (item) => {
  const usagePercentage = getRackUsagePercentage(item.code);
  const visiblePercentage = Math.min(usagePercentage, 100);
  const labelPosition = Math.min(Math.max(visiblePercentage, 8), 92);

  return {
    left: `${item.x}px`,
    top: `${item.y}px`,
    width: `${item.w}px`,
    height: `${item.h}px`,
    "--usage-fill": `${visiblePercentage}%`,
    "--usage-label": `${labelPosition}%`,
  };
};

const sensorHotspotStyle = (sensor) => ({
  left: `${sensor.x - 10}px`,
  top: `${sensor.y - 10}px`,
});

const getSensorData = (deviceName) => {
  return tempHumidityData.value[deviceName] || null;
};

const formatSensorValue = (deviceName, field, unit) => {
  const value = tempHumidityData.value[deviceName]?.[field];
  const numberValue = Number(value);
  return Number.isFinite(numberValue) ? `${numberValue.toFixed(1)}${unit}` : `--${unit}`;
};

const getRackCodeParts = (rackCode) => {
  const match = String(rackCode).match(/^([A-Z]+)(\d+)$/);
  if (!match) {
    return { prefix: "", number: 0 };
  }

  return {
    prefix: match[1],
    number: Number(match[2]),
  };
};

const rackLabelClass = (rackCode) => {
  const { prefix, number } = getRackCodeParts(rackCode);

  if (prefix === "A" && number >= 29 && number <= 47) {
    return "label-right";
  }

  if (
    (prefix === "A" && number >= 13 && number <= 28) ||
    (prefix === "A" && number >= 48 && number <= 58) ||
    (prefix === "B" && number >= 1 && number <= 44)
  ) {
    return "label-bottom";
  }

  return "label-top";
};

const getRackSensorDeviceName = (rackCode) => {
  const { prefix } = getRackCodeParts(rackCode);
  if (prefix === "A") return "KTP01";
  if (prefix === "B") return "KTP02";
  return "";
};

const escapeArrowHeadLength = 16;
const escapeArrowHeadWidth = 18;

const formatSvgPoints = (points) => points.map(([x, y]) => `${x},${y}`).join(" ");

const getLastSegment = (points) => {
  const tip = points[points.length - 1];
  const previous = points[points.length - 2];
  const dx = tip[0] - previous[0];
  const dy = tip[1] - previous[1];
  const length = Math.hypot(dx, dy) || 1;

  return {
    tip,
    unitX: dx / length,
    unitY: dy / length,
  };
};

const routeBodyPoints = (points) => {
  if (points.length < 2) return formatSvgPoints(points);

  const { tip, unitX, unitY } = getLastSegment(points);
  const bodyEnd = [
    tip[0] - unitX * escapeArrowHeadLength,
    tip[1] - unitY * escapeArrowHeadLength,
  ];

  return formatSvgPoints([...points.slice(0, -1), bodyEnd]);
};

const routeArrowHeadPoints = (points) => {
  if (points.length < 2) return "";

  const { tip, unitX, unitY } = getLastSegment(points);
  const baseX = tip[0] - unitX * escapeArrowHeadLength;
  const baseY = tip[1] - unitY * escapeArrowHeadLength;
  const halfWidth = escapeArrowHeadWidth / 2;
  const perpendicularX = -unitY * halfWidth;
  const perpendicularY = unitX * halfWidth;

  return formatSvgPoints([
    tip,
    [baseX + perpendicularX, baseY + perpendicularY],
    [baseX - perpendicularX, baseY - perpendicularY],
  ]);
};

const doorPanelThickness = 10;

const doorPanelY = (door) =>
  door.direction === "up" ? door.y - doorPanelThickness : door.y;

const stairsStepPoints = (shape) => {
  const stepCount = Number(shape.steps || 5);
  const stepWidth = shape.w / stepCount;
  const stepHeight = shape.h / stepCount;
  const points = [[shape.x, shape.y + shape.h]];

  for (let index = 1; index <= stepCount; index += 1) {
    points.push([shape.x + stepWidth * index, shape.y + shape.h - stepHeight * (index - 1)]);
    points.push([shape.x + stepWidth * index, shape.y + shape.h - stepHeight * index]);
  }

  return formatSvgPoints(points);
};


const formatNumber = (value) => {
  const numberLocales = { vi: "vi-VN", en: "en-US", zh: "zh-CN" };
  return Number(value || 0).toLocaleString(numberLocales[locale.value] || "vi-VN");
};

const getRackUsagePercentage = (rackCode) =>
  calculateUsagePercentage(
    getRack(rackCode).totalQty,
    getRackCapacity(rackCode),
  );

const scheduleRealtimeRefresh = () => {
  if (realtimeRefreshTimer) {
    clearTimeout(realtimeRefreshTimer);
  }

  realtimeRefreshTimer = setTimeout(() => {
    fetchRacks({ silent: true });
  }, 250);
};

const connectWarehouseKTPRealtime = () => {
  if (!window.EventSource || !API_URL) {
    return;
  }

  ktpEventSource = new EventSource(`${API_URL}/sse`);
  ktpEventSource.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (data.type === "WAREHOUSE_KTP_RACKS_UPDATED") {
        scheduleRealtimeRefresh();
      }
    } catch (err) {
      console.warn("Invalid warehouse KTP realtime event:", err);
    }
  };
};

onMounted(() => {
  fetchRacks();
  fetchTempHumidityByDevices();
  tempHumidityTimer = setInterval(fetchTempHumidityByDevices, 5 * 60 * 1000);
  connectWarehouseKTPRealtime();
});

onUnmounted(() => {
  if (ktpEventSource) {
    ktpEventSource.close();
  }
  if (realtimeRefreshTimer) {
    clearTimeout(realtimeRefreshTimer);
  }
  if (tempHumidityTimer) {
    clearInterval(tempHumidityTimer);
  }
});
</script>

<style lang="scss" scoped>
.ktp-page {
  min-height: 100%;
  background:
    linear-gradient(180deg, #f8fafc 0%, #eef2f7 100%);
}

.warehouse-scroll {
  height: calc(100vh - 70px);
  overflow: auto;
  padding: 0;
  background: transparent;
  scrollbar-color: #94a3b8 #e2e8f0;
  scrollbar-width: thin;
}

.warehouse-scroll::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

.warehouse-scroll::-webkit-scrollbar-track {
  background: #e2e8f0;
}

.warehouse-scroll::-webkit-scrollbar-thumb {
  border: 2px solid #e2e8f0;
  border-radius: 999px;
  background: #94a3b8;
}

.warehouse-board {
  position: relative;
  margin-top: -50px;
  background:
    linear-gradient(rgba(148, 163, 184, 0.07) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148, 163, 184, 0.07) 1px, transparent 1px),
    #ffffff;
  background-size: 32px 32px;
  box-shadow: inset 0 0 0 1px rgba(148, 163, 184, 0.18);
  transform-origin: top left;
}

.layout-layer {
  position: absolute;
  inset: 0;
  z-index: 1;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.capacity-summary {
  position: absolute;
  left: 570px;
  top: 145px;
  z-index: 3;
  display: flex;
  align-items: center;
  gap: 14px;
  width: 300px;
  padding: 9px 12px;
  border: 1px solid rgba(148, 163, 184, 0.42);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.12);
  font-family: "Segoe UI", Arial, Helvetica, sans-serif;
}

.capacity-chart {
  position: relative;
  width: 78px;
  height: 78px;
  flex: 0 0 78px;
  border-radius: 50%;
  box-shadow: inset 0 0 0 1px rgba(15, 23, 42, 0.08);
}

.capacity-chart::after {
  position: absolute;
  inset: 15px;
  border-radius: 50%;
  background: #ffffff;
  box-shadow: 0 0 0 1px rgba(148, 163, 184, 0.25);
  content: "";
}

.capacity-chart-center {
  position: absolute;
  inset: 15px;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #1f2937;
  line-height: 1;

  strong {
    font-size: 16px;
    font-weight: 900;
  }

  span {
    margin-top: 3px;
    color: #64748b;
    font-size: 8px;
    font-weight: 800;
    text-transform: uppercase;
  }
}

.capacity-metrics {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
  gap: 7px;
}

.capacity-metric {
  display: grid;
  grid-template-columns: 78px minmax(0, 1fr);
  align-items: baseline;
  gap: 7px;
  font-size: 12px;
  line-height: 1.15;

  span,
  strong {
    font-weight: 900;
  }

  strong {
    color: #0f172a;
    white-space: nowrap;
  }

  &.is-total span {
    color: #ff0000;
  }

  &.is-used span {
    color: #0284c7;
  }

  &.is-remaining span {
    color: #16a34a;
  }
}

.wall-line {
  stroke: #0f172a;
  stroke-width: 3.5;
  stroke-linecap: square;
  vector-effect: non-scaling-stroke;
}

.escape-route {
  fill: none;
  stroke: #ef4444;
  stroke-width: 8;
  stroke-linecap: round;
  stroke-linejoin: round;
  vector-effect: non-scaling-stroke;
}

.escape-arrow-head {
  fill: #ef4444;
}

.shape-rect {
  fill: none;
  stroke: #475569;
  stroke-width: 1.8;
  vector-effect: non-scaling-stroke;
}

.shape-rect.is-machine,
.shape-rect.is-thin {
  fill: #eaf3ff;
  stroke: #1d4ed8;
  stroke-width: 1.8;
}

.shape-circle.is-temperature-sensor {
  fill: #ef4444;
  stroke: #ffffff;
  stroke-width: 2;
  vector-effect: non-scaling-stroke;
}

.sensor-hotspot {
  position: absolute;
  z-index: 3;
  width: 20px;
  height: 20px;
  pointer-events: auto;
}

.sensor-dot {
  position: absolute;
  left: 50%;
  top: 50%;
  width: 12px;
  height: 12px;
  border: 2px solid #ffffff;
  border-radius: 999px;
  background: #ef4444;
  box-shadow: 0 2px 6px rgba(15, 23, 42, 0.22);
  transform: translate(-50%, -50%);
}

.sensor-card {
  position: absolute;
  left: 18px;
  top: -10px;
  display: grid;
  grid-template-columns: 1fr;
  gap: 2px;
  min-width: 72px;
  padding: 5px 7px;
  border: 1px solid rgba(239, 68, 68, 0.45);
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 6px 14px rgba(15, 23, 42, 0.12);
  color: #0f172a;
  font-family: "Segoe UI", Arial, Helvetica, sans-serif;
  opacity: 0;
  pointer-events: none;
  transform: translateY(4px);
  transition:
    opacity 0.15s ease,
    transform 0.15s ease;
}

.sensor-hotspot:hover .sensor-card {
  opacity: 1;
  transform: translateY(0);
}

.sensor-name {
  color: #b91c1c;
  font-size: 11px;
  font-weight: 900;
  line-height: 1;
}

.sensor-value {
  color: #1e293b;
  font-size: 11px;
  font-weight: 800;
  line-height: 1.1;
}

.door-shape {
  fill: #ffffff;
  stroke: #16a34a;
  stroke-width: 2;
  stroke-linecap: square;
  vector-effect: non-scaling-stroke;
}

.door-shape line {
  stroke-width: 1.5;
}

.stairs-shape {
  fill: #f8fafc;
  stroke: #475569;
  stroke-width: 1.8;
  stroke-linecap: square;
  stroke-linejoin: round;
  vector-effect: non-scaling-stroke;
}

.stairs-shape line {
  stroke: #94a3b8;
  stroke-width: 1;
}

.stairs-shape polyline {
  fill: none;
  stroke: #0f172a;
  stroke-width: 2;
}

.shape-text,
.layout-text {
  fill: #1e293b;
  font-family: "Segoe UI", Arial, Helvetica, sans-serif;
  font-size: 12px;
  font-weight: 700;
}

.layout-text.is-red-small,
.layout-text.is-delivery {
  fill: #b91c1c;
}

.layout-text.is-cfqc {
  font-family: "Times New Roman", Times, serif;
  font-size: 38px;
  font-weight: 500;
}

.layout-text.is-stairs {
  font-family: "Segoe UI", Arial, Helvetica, sans-serif;
  font-size: 14px;
  font-weight: 900;
}

.layout-text.is-map-title {
  fill: #0f172a;
  font-family: "Segoe UI", Arial, Helvetica, sans-serif;
  font-size: 27px;
  font-weight: 900;
}

.layout-text.is-map-subtitle {
  fill: #475569;
  font-family: "Microsoft YaHei", "Segoe UI", Arial, Helvetica, sans-serif;
  font-size: 24px;
  font-weight: 700;
}

.rack-cell {
  position: absolute;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1px;
  border: 1px solid #2f855a;
  border-radius: 2px;
  background: #00ff73;
  color: #1f2933;
  cursor: pointer;
  overflow: visible;
  transition:
    border-color 0.15s ease,
    box-shadow 0.15s ease,
    filter 0.15s ease,
    transform 0.15s ease;

  &:hover {
    border-color: #2f855a;
    box-shadow:
      0 0 0 2px rgba(47, 133, 90, 0.24),
      0 5px 14px rgba(15, 23, 42, 0.16);
    filter: saturate(1.05);
    transform: translateY(-1px);
    z-index: 4;
  }
}

.rack-cell.is-occupied {
  background: #00ff73;
}

.rack-cell.is-empty {
  border-color: #2f855a;
  background: #00ff73;
}

.rack-cell.is-selected {
  border-color: #1d4ed8;
  box-shadow:
    0 0 0 3px rgba(37, 99, 235, 0.26),
    0 8px 16px rgba(15, 23, 42, 0.14);
  z-index: 5;
}

.rack-cell.is-horizontal {
  justify-content: center;
}

.rack-usage-fill {
  position: absolute;
  inset: auto 0 0;
  z-index: 0;
  width: 100%;
  height: var(--usage-fill);
  max-height: 100%;
  border-radius: 1px;
  background: #489ae1;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.45);
  pointer-events: none;
  transition: height 0.25s ease;
}

.rack-cell.is-horizontal .rack-usage-fill {
  inset: 0 auto 0 0;
  width: var(--usage-fill);
  max-width: 100%;
  height: 100%;
  max-height: none;
  transition: width 0.25s ease;
}

.rack-code {
  position: absolute;
  padding: 1px 2px;
  border-radius: 2px;
  background: rgba(255, 255, 255, 0.9);
  color: #dc2626;
  font-size: 8px;
  font-weight: 900;
  line-height: 1;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.12);
  pointer-events: none;
  white-space: nowrap;
}

.rack-cell.label-top .rack-code {
  bottom: calc(100% + 2px);
  left: 50%;
  transform: translateX(-50%);
}

.rack-cell.label-right .rack-code {
  top: 50%;
  left: calc(100% + 3px);
  transform: translateY(-50%);
}

.rack-cell.label-bottom .rack-code {
  top: calc(100% + 2px);
  left: 50%;
  transform: translateX(-50%);
}

.qty-text {
  position: absolute;
  left: 50%;
  bottom: var(--usage-label);
  z-index: 1;
  padding: 1px 2px;
  border-radius: 3px;
  background: rgba(224, 242, 254, 0.96);
  color: #1f2937;
  font-size: 8px;
  font-weight: 900;
  line-height: 1;
  pointer-events: none;
  transform: translate(-50%, 50%);
  white-space: nowrap;
}

.rack-cell.is-horizontal .qty-text {
  left: var(--usage-label);
  bottom: auto;
  top: 50%;
  transform: translate(-50%, -50%);
}

.rack-cell.is-over-capacity .qty-text {
  background: rgba(254, 226, 226, 0.96);
  color: #b91c1c;
}
</style>
