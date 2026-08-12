<template>
  <WarehouseLayoutPTN>
    <template #default>
      <el-card class="box-card">
        <div class="lab-page">
          <div class="main-layout">
            <div class="map-container">
              <div class="map-wrapper" ref="wrapperRef">
                <img
                  :src="labImage"
                  :alt="t('laboratory.alt')"
                  class="lab-image"
                  ref="imgRef"
                  @load="onImageLoad"
                />

                <div
                  class="room-overlay"
                  :style="getOverlayStyle(261, 657, 656, 1100)"
                  :title="t('laboratory.conditionRoom')"
                >
                  <div class="overlay-card">
                    <div class="ov-room-name">
                      {{ t("laboratory.conditionRoom") }}
                    </div>
                    <div class="ov-env-row">
                      <div class="ov-env-item">
                        <span class="ov-ico">🌡️</span>
                        <span
                          class="ov-val temp-col"
                          :class="{
                            'text-danger': isTempAbnormal(
                              envMap.ThiNghiem9527_L4.temperature,
                            ),
                          }"
                        >
                          {{
                            envMap.ThiNghiem9527_L4.temperature === "--"
                              ? "--"
                              : `${envMap.ThiNghiem9527_L4.temperature}${t('laboratory.temperatureUnit')}`
                          }}
                        </span>
                      </div>
                      <div class="ov-sep">|</div>
                      <div class="ov-env-item">
                        <span class="ov-ico">💧</span>
                        <span
                          class="ov-val hum-col"
                          :class="{
                            'text-danger': isHumAbnormal(
                              envMap.ThiNghiem9527_L4.humidity,
                            ),
                          }"
                        >
                          {{
                            envMap.ThiNghiem9527_L4.humidity === "--"
                              ? "--"
                              : `${envMap.ThiNghiem9527_L4.humidity}%`
                          }}
                        </span>
                      </div>
                    </div>
                    <div
                      class="ov-time"
                      v-if="envMap.ThiNghiem9527_L4.recordTime"
                    >
                      {{ envMap.ThiNghiem9527_L4.recordTime }}
                    </div>
                  </div>
                </div>

                <div
                  class="room-overlay"
                  :style="getOverlayStyle(265, 1097, 1178, 2253)"
                  :title="t('laboratory.bondingRoom')"
                >
                  <div class="overlay-card overlay-card--lg">
                    <div class="ov-room-name">
                      {{ t("laboratory.bondingRoom") }}
                    </div>
                    <div class="ov-env-row ov-env-row--lg">
                      <div class="ov-env-item ov-env-item--lg">
                        <span class="ov-ico ov-ico--lg">🌡️</span>
                        <span
                          class="ov-val ov-val--lg temp-col"
                          :class="{
                            'text-danger': isTempAbnormal(
                              envMap.ThiNghiem9517_L4.temperature,
                            ),
                          }"
                        >
                          {{
                            envMap.ThiNghiem9517_L4.temperature === "--"
                              ? "--"
                              : `${envMap.ThiNghiem9517_L4.temperature}${t('laboratory.temperatureUnit')}`
                          }}
                        </span>
                      </div>
                      <div class="ov-sep ov-sep--lg">|</div>
                      <div class="ov-env-item ov-env-item--lg">
                        <span class="ov-ico ov-ico--lg">💧</span>
                        <span
                          class="ov-val ov-val--lg hum-col"
                          :class="{
                            'text-danger': isHumAbnormal(
                              envMap.ThiNghiem9517_L4.humidity,
                            ),
                          }"
                        >
                          {{
                            envMap.ThiNghiem9517_L4.humidity === "--"
                              ? "--"
                              : `${envMap.ThiNghiem9517_L4.humidity}%`
                          }}
                        </span>
                      </div>
                    </div>
                    <div
                      class="ov-time"
                      v-if="envMap.ThiNghiem9517_L4.recordTime"
                    >
                      {{ envMap.ThiNghiem9517_L4.recordTime }}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="left-column">
              <div class="toolbar-card">
                <div class="toolbar-left">
                  <div class="label-with-icon">
                    <svg
                      viewBox="0 0 24 24"
                      width="16"
                      height="16"
                      stroke="currentColor"
                      stroke-width="2"
                      fill="none"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      class="icon-svg"
                    >
                      <rect
                        x="3"
                        y="4"
                        width="18"
                        height="18"
                        rx="2"
                        ry="2"
                      ></rect>
                      <line x1="16" y1="2" x2="16" y2="6"></line>
                      <line x1="8" y1="2" x2="8" y2="6"></line>
                      <line x1="3" y1="10" x2="21" y2="10"></line>
                    </svg>
                    <span class="toolbar-label">{{ t("laboratory.date") }}</span>
                  </div>

                  <el-date-picker
                    v-model="selectedDate"
                    type="date"
                    :placeholder="t('laboratory.datePlaceholder')"
                    :disabled-date="disableFuture"
                    format="DD/MM/YYYY"
                    value-format="YYYY-MM-DD"
                    size="default"
                    clearable
                    class="custom-date-picker"
                    style="width: 140px"
                    @change="onDateChange"
                    @clear="onClear"
                  />

                  <div
                    class="toolbar-right"
                    @click="filterAbnormal = !filterAbnormal"
                  >
                    <span
                      class="switch-text"
                      :class="{ 'text-danger': filterAbnormal }"
                    >
                      {{ t("laboratory.abnormalOnly") }}
                    </span>
                    <el-switch
                      v-model="filterAbnormal"
                      class="custom-switch"
                      style="
                        --el-switch-on-color: #ef4444;
                        --el-switch-off-color: #cbd5e1;
                      "
                    />
                  </div>
                  <button
                    class="icon-btn"
                    @click="resetFilters"
                    :title="t('common.reset')"
                  >
                    <svg
                      viewBox="0 0 24 24"
                      width="16"
                      height="16"
                      stroke="currentColor"
                      stroke-width="2"
                      fill="none"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <path
                        d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"
                      ></path>
                      <path d="M3 3v5h5"></path>
                    </svg>
                  </button>
                </div>
              </div>

              <div class="side-panel">
                <div class="panel-header temp-header">
                  <span class="panel-icon">🏠</span>
                  <div>
                    <div class="panel-title">
                      {{ t("laboratory.conditionRoom") }}
                    </div>
                    <div class="panel-subtitle">
                      {{ t("laboratory.history") }}
                    </div>
                  </div>
                </div>

                <div class="room-table-wrapper">
                  <div v-if="loading" class="panel-empty">
                    <div class="empty-icon">⏳</div>
                    <div>{{ t("laboratory.loading") }}</div>
                  </div>

                  <table
                    v-else-if="filteredHistoryByRoom.ThiNghiem9527_L4?.length"
                    class="history-table"
                  >
                    <thead>
                      <tr>
                        <th>🌡️ {{ t("laboratory.temperature") }}</th>
                        <th>💧 {{ t("laboratory.humidity") }}</th>
                        <th>🕒 {{ t("laboratory.time") }}</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr
                        v-for="(
                          row, j
                        ) in filteredHistoryByRoom.ThiNghiem9527_L4"
                        :key="j"
                        :class="{ 'row-latest': j === 0 && isToday }"
                      >
                        <td
                          class="val-temp"
                          :class="{
                            'text-danger': isTempAbnormal(row.temperature),
                          }"
                        >
                          {{ row.temperature }}{{ t("laboratory.temperatureUnit") }}
                        </td>
                        <td
                          class="val-hum"
                          :class="{
                            'text-danger': isHumAbnormal(row.humidity),
                          }"
                        >
                          {{ row.humidity }}%
                        </td>
                        <td class="val-time">{{ row.recordTime }}</td>
                      </tr>
                    </tbody>
                  </table>

                  <div v-else class="panel-empty">
                    <div class="empty-icon">📡</div>
                    <div>
                      {{ t("laboratory.noData") }}
                    </div>
                  </div>
                </div>
              </div>

              <div class="side-panel">
                <div class="panel-header hum-header">
                  <span class="panel-icon">🏠</span>
                  <div>
                    <div class="panel-title">
                      {{ t("laboratory.bondingRoom") }}
                    </div>
                    <div class="panel-subtitle">
                      {{ t("laboratory.history") }}
                    </div>
                  </div>
                </div>

                <div class="room-table-wrapper">
                  <div v-if="loading" class="panel-empty">
                    <div class="empty-icon">⏳</div>
                    <div>{{ t("laboratory.loading") }}</div>
                  </div>

                  <table
                    v-else-if="filteredHistoryByRoom.ThiNghiem9517_L4?.length"
                    class="history-table"
                  >
                    <thead>
                      <tr>
                        <th>🌡️ {{ t("laboratory.temperature") }}</th>
                        <th>💧 {{ t("laboratory.humidity") }}</th>
                        <th>🕒 {{ t("laboratory.time") }}</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr
                        v-for="(
                          row, j
                        ) in filteredHistoryByRoom.ThiNghiem9517_L4"
                        :key="j"
                        :class="{ 'row-latest': j === 0 && isToday }"
                      >
                        <td
                          class="val-temp"
                          :class="{
                            'text-danger': isTempAbnormal(row.temperature),
                          }"
                        >
                          {{ row.temperature }}{{ t("laboratory.temperatureUnit") }}
                        </td>
                        <td
                          class="val-hum"
                          :class="{
                            'text-danger': isHumAbnormal(row.humidity),
                          }"
                        >
                          {{ row.humidity }}%
                        </td>
                        <td class="val-time">{{ row.recordTime }}</td>
                      </tr>
                    </tbody>
                  </table>

                  <div v-else class="panel-empty">
                    <div class="empty-icon">📡</div>
                    <div>
                      {{ t("laboratory.noData") }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </el-card>
    </template>
  </WarehouseLayoutPTN>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref, reactive, computed } from "vue";
import WarehouseLayoutPTN from "@/views/components/WarehousePTN/WarehouseLayoutPTN.vue";
import labImage from "@/assets/laboratory.jpg";
import { useWarehouseMapI18n } from "@/hooks/useWarehouseMapI18n";

const API_URL = (import.meta.env.VITE_API_URL || "/api/v1").replace(/\/$/, "");
const { t } = useWarehouseMapI18n();

// ─── helpers ───────────────────────────────────────────────────────────────

const todayStr = () => {
  const d = new Date();
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
};

// ─── state ─────────────────────────────────────────────────────────────────

const selectedDate = ref(todayStr());
const loading = ref(false);
const filterAbnormal = ref(false);

const isToday = computed(
  () => !selectedDate.value || selectedDate.value === todayStr(),
);

// Hàm Reset toàn bộ trạng thái về mặc định
const resetFilters = () => {
  selectedDate.value = todayStr();
  filterAbnormal.value = false;
  fetchData();
};

const onClear = () => {
  selectedDate.value = todayStr();
  fetchData();
};

const emptyEnv = () => ({ temperature: "--", humidity: "--", recordTime: "" });

const envMap = reactive({
  ThiNghiem9517_L4: emptyEnv(),
  ThiNghiem9527_L4: emptyEnv(),
});

const mergedHistoryByRoom = reactive({
  ThiNghiem9527_L4: [],
  ThiNghiem9517_L4: [],
});

// ─── filter helpers ────────────────────────────────────────────────────────

const isTempAbnormal = (t) => {
  if (t === "--") return false;
  const val = parseFloat(t);
  return !isNaN(val) && (val < 21 || val > 25);
};

const isHumAbnormal = (h) => {
  if (h === "--") return false;
  const val = parseFloat(h);
  return !isNaN(val) && (val < 45 || val > 65);
};

const filteredHistoryByRoom = computed(() => {
  if (!filterAbnormal.value) return mergedHistoryByRoom;

  const result = {};
  for (const room in mergedHistoryByRoom) {
    result[room] = mergedHistoryByRoom[room].filter(
      (row) => isTempAbnormal(row.temperature) || isHumAbnormal(row.humidity),
    );
  }
  return result;
});

// ─── image overlay ──────────────────────────────────────────────────────────

const imgRef = ref(null);
const naturalSize = ref({ w: 1, h: 1 });

const onImageLoad = () => {
  naturalSize.value = {
    w: imgRef.value?.naturalWidth || 1,
    h: imgRef.value?.naturalHeight || 1,
  };
};

const getOverlayStyle = (x1, y1, x2, y2) => {
  const { w, h } = naturalSize.value;
  return {
    left: `${(x1 / w) * 100}%`,
    top: `${(y1 / h) * 100}%`,
    width: `${((x2 - x1) / w) * 100}%`,
    height: `${((y2 - y1) / h) * 100}%`,
  };
};

// ─── date picker helpers ────────────────────────────────────────────────────

const disableFuture = (date) => date > new Date();

const onDateChange = (val) => {
  if (!val) {
    selectedDate.value = todayStr();
  }
  fetchData();
};

// ─── fetch ──────────────────────────────────────────────────────────────────

let fetchController = null;
let latestRequestId = 0;

const fetchData = async ({ silent = false } = {}) => {
  const requestId = ++latestRequestId;
  const requestDate = selectedDate.value || todayStr();

  fetchController?.abort();
  fetchController = new AbortController();

  if (!silent) loading.value = true;

  try {
    const params = new URLSearchParams();
    params.set("date", requestDate);
    params.set("_", Date.now().toString());

    const res = await fetch(
      `${API_URL}/warehouse/laboratory/environment?${params.toString()}`,
      {
        cache: "no-store",
        signal: fetchController.signal,
      },
    );
    if (!res.ok) throw new Error(`HTTP ${res.status}`);

    const json = await res.json();
    if (requestId !== latestRequestId || requestDate !== selectedDate.value) {
      return;
    }

    const arr = [...(json.data ?? [])].sort(
      (a, b) => new Date(b.RecordTime || 0) - new Date(a.RecordTime || 0),
    );

    const now = new Date()
      .toLocaleString("sv-SE")
      .replace("T", " ")
      .slice(0, 19);

    Object.keys(mergedHistoryByRoom).forEach((key) => {
      mergedHistoryByRoom[key] = [];
    });
    Object.keys(envMap).forEach((key) => {
      envMap[key] = emptyEnv();
    });

    const grouped = {};
    arr.forEach((x) => {
      const device = x.DeviceName;
      if (!grouped[device]) grouped[device] = [];
      grouped[device].push({
        temperature: x.Tem != null ? parseFloat(x.Tem).toFixed(1) : "--",
        humidity: x.Hum != null ? parseFloat(x.Hum).toFixed(1) : "--",
        recordTime: x.RecordTime
          ? x.RecordTime.slice(0, 19).replace("T", " ")
          : now,
      });
    });

    Object.keys(grouped).forEach((device) => {
      const rows = grouped[device];

      if (mergedHistoryByRoom[device] !== undefined) {
        mergedHistoryByRoom[device] = rows;
      }

      const latest = rows[0];
      if (latest && envMap[device]) {
        envMap[device] = {
          temperature: latest.temperature,
          humidity: latest.humidity,
          recordTime: latest.recordTime,
        };
      }
    });
  } catch (err) {
    if (err.name === "AbortError") return;
    console.error("Fetch error:", err);
  } finally {
    if (requestId === latestRequestId) {
      loading.value = false;
      fetchController = null;
    }
  }
};

// ─── polling ────────────────────────────────────────────────────────────────

let pollTimer = null;
let pollingDate = todayStr();
const POLL_INTERVAL_MS = 5 * 60 * 1000;
const POLL_DELAY_MS = 10 * 1000;

const refreshLatestData = () => {
  const currentDate = todayStr();

  // Keep the live view on the current day when the page stays open overnight.
  if (selectedDate.value === pollingDate) {
    selectedDate.value = currentDate;
  }
  pollingDate = currentDate;

  if (selectedDate.value === currentDate) {
    fetchData({ silent: true });
  }
};

const scheduleNextPoll = () => {
  const now = Date.now();
  const nextFiveMinuteMark =
    Math.ceil(now / POLL_INTERVAL_MS) * POLL_INTERVAL_MS;
  const delay = Math.max(
    1000,
    nextFiveMinuteMark + POLL_DELAY_MS - now,
  );

  pollTimer = setTimeout(() => {
    refreshLatestData();
    scheduleNextPoll();
  }, delay);
};

const startPoll = () => {
  stopPoll();
  scheduleNextPoll();
};

const stopPoll = () => {
  if (pollTimer) {
    clearTimeout(pollTimer);
    pollTimer = null;
  }
};

const onVisibilityChange = () => {
  if (document.visibilityState !== "visible") return;
  refreshLatestData();
  startPoll();
};

onMounted(() => {
  fetchData();
  startPoll();
  document.addEventListener("visibilitychange", onVisibilityChange);
});

onBeforeUnmount(() => {
  stopPoll();
  fetchController?.abort();
  document.removeEventListener("visibilitychange", onVisibilityChange);
});
</script>

<style lang="scss" scoped>
.box-card {
  width: 100%;
  height: 100%;
}

:deep(.el-card__body) {
  padding: 0;
}

.lab-page {
  width: 100%;
  height: calc(100vh - 60px);
  display: flex;
  flex-direction: column;
  background: #f0f4fa;
  overflow: hidden;
  box-sizing: border-box;
  padding: 10px;
}

.main-layout {
  display: flex;
  gap: 10px;
  width: 100%;
  height: 100%;
  min-height: 0;
}

.left-column {
  width: 520px;
  min-width: 520px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  height: 100%;
}

// ── Control Toolbar ──────────────────────────────────────────────────────────
.toolbar-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 12px 16px;
  box-shadow: 0 2px 6px rgba(15, 23, 42, 0.04);
  flex-shrink: 0;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.label-with-icon {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #475569;
}

.icon-svg {
  color: #64748b;
  margin-bottom: 1px;
}

.toolbar-label {
  font-size: 13px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

:deep(.custom-date-picker .el-input__wrapper) {
  box-shadow: 0 0 0 1px #cbd5e1 inset;
  border-radius: 6px;
  transition: all 0.2s ease;

  &:hover {
    box-shadow: 0 0 0 1px #94a3b8 inset;
  }
  &.is-focus {
    box-shadow: 0 0 0 1px #3b82f6 inset !important;
  }
}

/* Style cho nút Reset (Làm mới) */
.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 6px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s ease;
  outline: none;

  &:hover {
    background: #e2e8f0;
    color: #0f172a;
    border-color: #cbd5e1;
  }

  &:active {
    transform: scale(0.92);
  }
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  background: #f8fafc;
  border: 1px solid #f1f5f9;
  transition: all 0.2s;

  &:hover {
    background: #f1f5f9;
  }
}

.switch-text {
  font-size: 13px;
  font-weight: 600;
  color: #64748b;
  user-select: none;
  transition: color 0.3s ease;

  &.text-danger {
    color: #ef4444;
  }
}

.custom-switch {
  pointer-events: none;
}

.text-danger {
  color: #ef4444 !important;
  font-weight: 800;
  animation: blink 2s infinite;
}

@keyframes blink {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.6;
  }
  100% {
    opacity: 1;
  }
}

// ── Side panels ──────────────────────────────────────────────────────────────
.side-panel {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  color: #fff;
  flex-shrink: 0;

  &.temp-header {
    background: linear-gradient(135deg, #649158 0%, #323b37 100%);
  }

  &.hum-header {
    background: linear-gradient(135deg, #1d4ed8 0%, #06b6d4 100%);
  }
}

.panel-icon {
  font-size: 20px;
  line-height: 1;
}

.panel-title {
  font-size: 20px;
  font-weight: 700;
}

.panel-subtitle {
  font-size: 10px;
  opacity: 0.85;
  line-height: 1.2;
}

.room-table-wrapper {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.panel-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 100%;
  color: #94a3b8;
  font-size: 11px;
  padding: 16px 8px;
  text-align: center;
}

.empty-icon {
  font-size: 22px;
}

// ── History table ─────────────────────────────────────────────────────────────
.history-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 15px;

  thead tr {
    position: sticky;
    top: 0;
    background: #f8fafc;
    z-index: 1;
  }

  th {
    padding: 5px 6px;
    text-align: left;
    font-weight: 600;
    color: #64748b;
    border-bottom: 1px solid #e2e8f0;
    white-space: nowrap;
  }

  td {
    padding: 4px 6px;
    border-bottom: 1px solid #f1f5f9;
    color: #374151;
  }

  tbody tr {
    transition: background 0.15s;

    &:hover {
      background: #f8fafc;
    }
  }

  .row-latest td {
    background: #fef3c7;
    font-weight: 600;
  }
}

.val-temp {
  font-weight: 700;
  color: #d97706;
  white-space: nowrap;
}

.val-hum {
  font-weight: 700;
  color: #1d4ed8;
  white-space: nowrap;
}

.val-time {
  color: #94a3b8;
  white-space: nowrap;
  font-size: 15px;
}

// ── Map ───────────────────────────────────────────────────────────────────────
.map-container {
  flex: 1;
  min-width: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: auto;
}

.map-wrapper {
  position: relative;
  display: inline-block;
  line-height: 0;
}

.lab-image {
  max-width: 100%;
  width: 1000px;
  max-height: calc(100vh - 80px);
  display: block;
  user-select: none;
}

// ── Room overlays ─────────────────────────────────────────────────────────────
.room-overlay {
  position: absolute;
  overflow: hidden;
  border: 2px solid rgba(37, 99, 235, 0.7);
  background: rgba(59, 130, 246, 0.07);
  transition: all 0.3s ease;
  animation: zonePulse 2.4s ease-in-out infinite;

  &::before {
    content: "";
    position: absolute;
    inset: -40%;
    background: linear-gradient(
      120deg,
      transparent 30%,
      rgba(255, 255, 255, 0.25),
      transparent 70%
    );
    transform: rotate(20deg) translateX(-180%);
    animation: scanFlow 3.5s linear infinite;
    pointer-events: none;
  }

  &:hover {
    border-color: #2563eb;
    background: rgba(59, 130, 246, 0.12);
    box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12),
      0 0 18px rgba(37, 99, 235, 0.22);
  }
}

@keyframes zonePulse {
  0% {
    box-shadow: 0 0 0 rgba(37, 99, 235, 0);
    opacity: 0.92;
  }
  50% {
    box-shadow: 0 0 14px rgba(37, 99, 235, 0.2), 0 0 28px rgba(37, 99, 235, 0.1);
    opacity: 1;
  }
  100% {
    box-shadow: 0 0 0 rgba(37, 99, 235, 0);
    opacity: 0.92;
  }
}

@keyframes scanFlow {
  from {
    transform: rotate(20deg) translateX(-180%);
  }
  to {
    transform: rotate(20deg) translateX(180%);
  }
}

.overlay-card {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
  padding: 10px 8px 12px;
  pointer-events: none;
}

.ov-room-name {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: clamp(9px, 1.2cqi, 13px);
  font-weight: 700;
  line-height: 1;
  color: #1e40af;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 999px;
  padding: 2px 8px;
  min-height: 20px;
  white-space: nowrap;
  box-shadow: 0 1px 5px rgba(0, 0, 0, 0.1);
}

.ov-env-row {
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 10px;
  padding: 4px 10px;
  border: 2px solid black;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.ov-val {
  font-weight: 800;
  font-size: clamp(10px, 1.5cqi, 15px);
  line-height: 1;
}

.ov-ico {
  font-size: clamp(10px, 1.2cqi, 13px);
}
.ov-sep {
  color: #cbd5e1;
  font-size: 12px;
  font-weight: 600;
}

.ov-time {
  font-size: clamp(7px, 0.9cqi, 10px);
  font-weight: 500;
  color: #64748b;
  background: rgba(255, 255, 255, 0.88);
  border-radius: 999px;
  padding: 1px 6px;
  line-height: 1.1;
  white-space: nowrap;
}

.overlay-card--lg {
  gap: 6px;

  .ov-room-name {
    font-size: clamp(10px, 1.8cqi, 16px);
    padding: 3px 12px;
  }
  .ov-time {
    font-size: clamp(8px, 1.2cqi, 12px);
  }
}

.temp-col {
  color: #d97706;
}
.hum-col {
  color: #1d4ed8;
}
</style>
