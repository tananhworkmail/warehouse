<template>
  <el-dialog
    :model-value="dialogVisible"
    :title="`Rack: ${rackNo}`"
    width="80%"
    @close="emit('close')"
  >
    <div class="rack-grid-wrapper">
      <el-tabs
        v-model="state.activeGroupTab"
        class="rack-group-tabs"
        v-if="rackGroupTabs.length"
      >
        <el-tab-pane
          v-for="group in rackGroupTabs"
          :key="group.group"
          :name="group.group"
          :label="group.group"
        >
          <div class="grid-and-gauge-row">
            <div class="rack-grid-area">
              <div class="rack-grid">
                <el-button
                  v-for="rack in group.items"
                  :key="rack"
                  @click="handleClick(rack)"
                  class="rack-button custom-rack-btn"
                  :class="{ 'rack-selected': selectedRack === rack }"
                >
                  <div class="rack-button-content">
                    <div class="rack-code">{{ rack }}</div>
                    <div
                      class="rack-ton"
                      :class="{
                        'zero-ton': tonKhoMap[rack.replace('-', '')] == 0,
                      }"
                      v-show="tonKhoMap[rack.replace('-', '')] > 0"
                    >
                      {{ t("rackDialog.toolbar.tonLabel") }}:
                      {{ tonKhoMap[rack.replace("-", "")] ?? 0 }}
                    </div>
                  </div>
                </el-button>
              </div>
            </div>

            <div class="gauge-wrapper">
              <el-progress
                type="dashboard"
                :percentage="
                  Math.round(
                    tempHumidityByGroup[group.group]?.Tem ??
                      temperatureAndHumidity?.Tem ??
                      0,
                  )
                "
                :color="warehouseTemperatureColor"
                :width="120"
              >
                <template #default="{ percentage }">
                  <span class="percentage-value">{{ percentage }}°C</span>
                  <span class="percentage-label"
                    ><br />{{ t("rackDialog.grid.temperatureLabel") }}</span
                  >
                </template>
              </el-progress>

              <el-progress
                type="dashboard"
                :percentage="
                  Math.round(
                    tempHumidityByGroup[group.group]?.Hum ??
                      temperatureAndHumidity?.Hum ??
                      0,
                  )
                "
                :color="warehouseHumidityColor"
                :width="120"
              >
                <template #default="{ percentage }">
                  <span class="percentage-value">{{ percentage }}%</span>
                  <span class="percentage-label"
                    ><br />{{ t("rackDialog.grid.humidityLabel") }}</span
                  >
                </template>
              </el-progress>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <div class="dialog-scrollable-content">
      <div
        class="rack-toolbar"
        style="display: flex; align-items: center; margin: 10px 0 0 0px"
      >
        <el-checkbox
          v-model="state.hideZero"
          :disabled="true"
          style="margin: 10px 0 0 10px"
        >
          {{ t("rackDialog.toolbar.hideZeroLabel") }}
        </el-checkbox>

        <div
          style="
            margin-left: 20px;
            display: flex;
            align-items: center;
            gap: 5px;
            margin-right: 10px;
          "
        >
          <el-input
            v-model="state.search"
            :placeholder="t('rackDialog.toolbar.searchPlaceholder')"
            size="medium"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>

          <div
            style="
              display: flex;
              align-items: center;
              gap: 6px;
              font-size: 13px;
            "
          >
            <span
              style="
                margin-left: 20px;
                width: 25px;
                height: 25px;
                background: #fff3cd;
                border: 1px solid #e6d8a2;
                display: inline-block;
              "
            ></span>
            <span style="white-space: nowrap">{{
              t("rackDialog.toolbar.legend3Day")
            }}</span>
          </div>

          <div
            style="
              margin-left: 20px;
              display: flex;
              align-items: center;
              gap: 6px;
              font-size: 13px;
            "
          >
            <span
              style="
                width: 25px;
                height: 25px;
                background: #fecaca;
                border: 1px solid #fca5a5;
                display: inline-block;
              "
            ></span>
            <span style="white-space: nowrap">{{
              t("rackDialog.toolbar.legendOver6M")
            }}</span>
          </div>
        </div>
      </div>
    </div>

    <div>
      <el-tabs v-model="state.tabsActiveName" v-if="state.rackTableData.length">
        <el-tab-pane :label="t('rackDialog.tabs.rackInfo')" name="rack">
          <el-table
            :key="rackNo"
            ref="rackTableRef"
            v-loading="loading"
            :data="filteredTableData"
            :row-class-name="tableRowClassName"
            border
            highlight-current-row
            max-height="50vh"
          >
            <el-table-column
              v-for="item in materialTableFieldList"
              :key="item.prop"
              :prop="item.prop"
              :label="item.label"
              :min-width="item.minWidth || '100px'"
              align="center"
              header-align="center"
            />
            <el-table-column
              :label="t('rackDialog.table.columns.ddbh')"
              width="90"
              align="center"
            >
              <template #default="{ row }">
                <el-button
                  type="primary"
                  size="small"
                  :icon="Search"
                  @click="viewDDBH(row)"
                  :disabled="!row.memo_ry"
                >
                  {{ t("rackDialog.table.viewButton") }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </div>
  </el-dialog>

  <el-dialog
    v-model="ddbhDialog"
    :title="t('rackDialog.ddbhDialog.title')"
    width="600px"
  >
    <div style="margin-bottom: 10px">
      <b>{{ t("rackDialog.ddbhDialog.cgno") }}:</b> {{ ddbhData.cgno }}
    </div>
    <div style="margin-bottom: 10px">
      <b>{{ t("rackDialog.ddbhDialog.clbh") }}:</b> {{ ddbhData.clbh }}
    </div>
    <div class="ddbh-tag-wrapper">
      <el-tag
        v-for="item in ddbhData.list"
        :key="item"
        size="large"
        class="ddbh-tag"
        effect="light"
      >
        {{ item }}
      </el-tag>
    </div>
  </el-dialog>
</template>

<script setup>
import {
  computed,
  defineProps,
  defineEmits,
  reactive,
  ref,
  watch,
  nextTick,
  onUnmounted,
} from "vue";
import axios from "axios";
import dayjs from "dayjs";
import { useLoading } from "@/hooks/useLoading";
import { useI18n } from "@/hooks/i18n";
import { Search } from "@element-plus/icons-vue";

const API_URL = import.meta.env.VITE_API_URL;
const { loading, showLoading, hideLoading } = useLoading();
const { t } = useI18n();
const rackTableRef = ref(null);

const props = defineProps({
  rackNo: String,
  dialogVisible: Boolean,
  rackDataList: {
    type: Array,
    default: () => [],
  },
  autoSelectRack: { type: String, default: "" },
});

const emit = defineEmits(["close"]);

const normalizeRackCode = (code) =>
  String(code ?? "")
    .replace(/-/g, "")
    .trim();

const rackGroupTabs = computed(() => {
  const groupMap = new Map();

  (props.rackDataList || []).forEach((code) => {
    const normalized = normalizeRackCode(code);
    if (!normalized) return;

    const group = normalized.slice(0, 5);
    const sub = normalized.slice(5, 7);
    const displayCode = sub ? `${group}-${sub}` : group;

    if (!groupMap.has(group)) {
      groupMap.set(group, new Set());
    }

    groupMap.get(group).add(displayCode);
  });

  return Array.from(groupMap.entries())
    .map(([group, set]) => ({
      group,
      items: Array.from(set).sort((a, b) =>
        normalizeRackCode(a).localeCompare(normalizeRackCode(b), "en", {
          numeric: true,
          sensitivity: "base",
        }),
      ),
    }))
    .sort((a, b) =>
      normalizeRackCode(a.group).localeCompare(
        normalizeRackCode(b.group),
        "en",
        {
          numeric: true,
          sensitivity: "base",
        },
      ),
    );
});

const tonKhoMap = ref({});

const materialTableFieldList = computed(() => [
  { label: t("rackDialog.table.columns.cgno"), prop: "cgno" },
  { label: t("rackDialog.table.columns.zsbh"), prop: "zsbh" },
  { label: t("rackDialog.table.columns.clbh"), prop: "clbh" },
  { label: t("rackDialog.table.columns.dwbh"), prop: "dwbh", minWidth: "70px" },
  { label: t("rackDialog.table.columns.qty"), prop: "qty" },
  { label: t("rackDialog.table.columns.remqty"), prop: "remqty" },
  { label: t("rackDialog.table.columns.dqty"), prop: "dqty" },
  { label: t("rackDialog.table.columns.pack"), prop: "pack", minWidth: "60px" },
  { label: t("rackDialog.table.columns.kcbh"), prop: "kcbh", minWidth: "80px" },
  { label: t("rackDialog.table.columns.barcode"), prop: "barcode" },
  { label: t("rackDialog.table.columns.cfmdate"), prop: "cfmdate" },
]);

const state = reactive({
  tabsActiveName: "rack",
  activeGroupTab: "",
  rackTableData: [],
  hideZero: true,
  search: "",
});

const selectedRack = ref(null);

const ddbhDialog = ref(false);
const ddbhData = reactive({
  scno: "",
  cgno: "",
  clbh: "",
  list: [],
});

const viewDDBH = (row) => {
  ddbhData.scno = row.scno;
  ddbhData.cgno = row.cgno;
  ddbhData.clbh = row.clbh;
  ddbhData.list = row.memo_ry
    ? row.memo_ry
        .split(",")
        .map((i) => i.trim())
        .filter(Boolean)
    : [];
  ddbhDialog.value = true;
};

const handleClick = async (displayRackCode) => {
  selectedRack.value = displayRackCode;
  const rackCode = displayRackCode.replace("-", "");
  rack3DaySet.value = new Set();
  await getRackInformation(rackCode);
};

const filteredTableData = computed(() => {
  let data = [...state.rackTableData];

  if (state.hideZero) {
    data = data.filter((item) => item.remqty !== 0);
  }

  if (state.search?.trim()) {
    const keyword = state.search.trim().toLowerCase();
    data = data.filter((item) =>
      Object.values(item).some((v) =>
        String(v).toLowerCase().includes(keyword),
      ),
    );
  }

  data.sort((a, b) => {
    const a3 = is3DayItem(a);
    const b3 = is3DayItem(b);
    const a180 = a.is_over_180;
    const b180 = b.is_over_180;

    if (a3 !== b3) return a3 ? -1 : 1;
    if (a180 !== b180) return a180 ? -1 : 1;
    return 0;
  });

  return data;
});

const fetchTonKhoInRack = async () => {
  try {
    const res = await axios.get(`${API_URL}/warehouse/tonkhorackkvt`, {
      params: { rackNo: props.rackNo },
    });
    const data = res.data.data || [];
    tonKhoMap.value = data.reduce((acc, row) => {
      acc[row.make] = row.tonkhorack;
      return acc;
    }, {});
  } catch (err) {
    console.error("❌ Lỗi lấy tổng tồn kho:", err);
  }
};

const getRackInformation = async (rackNo) => {
  showLoading();
  rack3DaySet.value = new Set();

  try {
    const res = await axios.get(`${API_URL}/warehouse/rackinforkvt`, {
      params: { rackNo: rackNo.trim() },
    });

    const rawData = Array.isArray(res.data.data) ? res.data.data : [];

    state.rackTableData = rawData.map((item) => {
      const isOver180 = item.cfmdate
        ? dayjs()
            .startOf("day")
            .diff(dayjs(item.cfmdate).startOf("day"), "day") > 180
        : false;

      return {
        ...item,
        remqty: Number(item.remqty),
        cfmdate: item.cfmdate ? dayjs(item.cfmdate).format("YYYY-MM-DD") : "",
        is_over_180: isOver180,
      };
    });

    await get3DayRackInformation(rackNo);
  } catch (err) {
    console.error("❌ Failed to load rack info:", err);
  } finally {
    hideLoading();
  }
};

const rack3DaySet = ref(new Set());

const get3DayRackInformation = async (rackNo) => {
  try {
    const res = await axios.get(`${API_URL}/warehouse/3dayrackinforkvt`, {
      params: { rackNo: rackNo.trim() },
    });
    const rawData = Array.isArray(res.data.data) ? res.data.data : [];
    rack3DaySet.value = new Set(
      rawData.map((item) => `${item.scno}_${item.clbh}`),
    );
  } catch (err) {
    console.error("❌ Failed to load 3day rack info:", err);
  }
};

const is3DayItem = (row) => {
  const key = `${row.scno}_${row.clbh}`;
  return rack3DaySet.value.has(key);
};

const tableRowClassName = ({ row }) => {
  if (row.is_over_180) return "row-over-180";
  if (is3DayItem(row)) return "row-3day";
  return "";
};

const temperatureAndHumidity = ref(null);
const tempHumidityByGroup = reactive({});
let tempInterval = null;

const warehouseTemperatureColor = (val) => {
  if (val < 25) return "#7dd3fc";
  if (val <= 30) return "#22c55e";
  return "#ef4444";
};

const warehouseHumidityColor = (val) => {
  if (val < 45) return "#f97316";
  if (val <= 60) return "#10b981";
  return "#0000F7";
};

const fetchTempHumidityByRack = async (rackCode) => {
  if (!rackCode) return;
  try {
    const res = await axios.get(`${API_URL}/warehouse/temphumidity`, {
      params: { rackNo: rackCode, date: dayjs().format("YYYY-MM-DD") },
    });
    tempHumidityByGroup[rackCode] = res.data.data || null;
  } catch (err) {
    console.error("❌ Lỗi lấy nhiệt độ / độ ẩm:", err);
    tempHumidityByGroup[rackCode] = null;
  }
};

// ── Watch nhiệt độ theo tab ──────────────────────────────────────────────────
watch(
  () => state.activeGroupTab,
  async (newGroup) => {
    if (!newGroup) return;
    await fetchTempHumidityByRack(newGroup);
    if (tempInterval) clearInterval(tempInterval);
    tempInterval = setInterval(
      () => fetchTempHumidityByRack(newGroup),
      5 * 60 * 1000,
    );
  },
  { immediate: true },
);

// ── Watch tồn kho khi dialog mở ─────────────────────────────────────────────
watch(
  () => props.dialogVisible,
  (open) => {
    if (open && props.rackNo) fetchTonKhoInRack();
  },
  { immediate: true },
);

// ── Watch chính: khi rackGroupTabs có data thì auto-select ──────────────────
watch(
  rackGroupTabs,
  async (tabs) => {
    if (!tabs.length) return;
    await nextTick();

    // Có autoSelectRack → tìm đúng tab + rack
    if (props.autoSelectRack) {
      const code = props.autoSelectRack.toUpperCase().replace(/-/g, "");
      const groupKey = code.slice(0, 5); // ví dụ: A3102

      const foundGroup = tabs.find((g) => g.group === groupKey);
      if (foundGroup) {
        state.activeGroupTab = foundGroup.group;

        const cellCode =
          code.length === 7
            ? `${code.slice(0, 5)}-${code.slice(5, 7)}` // A3102-01
            : null;

        const item = cellCode
          ? foundGroup.items.find((i) => i === cellCode) ?? foundGroup.items[0]
          : foundGroup.items[0];

        selectedRack.value = item;
        rack3DaySet.value = new Set();
        await getRackInformation(item.replace("-", ""));
        return;
      }
    }

    // Fallback: chọn group + item đầu tiên
    const firstGroup = tabs[0];
    const firstItem = firstGroup?.items?.[0];
    if (!firstItem) return;

    state.activeGroupTab = firstGroup.group;
    selectedRack.value = firstItem;
    rack3DaySet.value = new Set();
    await getRackInformation(firstItem.replace("-", ""));
  },
  { immediate: true },
);

// ── Watch clear table row khi filter thay đổi ───────────────────────────────
watch(
  () => filteredTableData.value,
  () => {
    rackTableRef.value?.setCurrentRow(null);
  },
);

onUnmounted(() => {
  if (tempInterval) clearInterval(tempInterval);
});
</script>

<style scoped lang="scss">
.el-table {
  user-select: text;
}

:deep(.el-table thead th) {
  background-color: #3b82f6 !important;
  color: #fff !important;
  font-weight: bold;
  text-align: center;
}

.rack-grid-wrapper {
  position: sticky;
  top: 0;
  z-index: 10;
  background: #ffffff;
  padding: 8px 0 4px 0;
  border-bottom: 1px solid #ebeef5;
}

.rack-toolbar {
  display: flex;
  align-items: center;
  margin-top: 6px;
  padding-left: 10px;
  gap: 16px;
}

.grid-and-gauge-row {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding-top: 8px;
}

.rack-grid-area {
  flex: 1;
  min-width: 0;
}

.rack-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 10px;
  padding: 10px;
}

/* ── Button rack bình thường ── */
.custom-rack-btn {
  background-color: #45484d;
  color: #ffffff !important;
  border: none;
  height: 60px;
  width: 100%;
  padding: 6px;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: background-color 0.2s ease, box-shadow 0.2s ease;

  &:hover {
    background-color: #34495e;
  }

  .rack-button-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    font-size: 14px;
  }

  .rack-code {
    font-size: 14px;
    font-weight: 600;
  }

  .rack-ton {
    font-size: 13px;
    color: #ffd700;
    font-weight: bold;
    margin-top: 4px;
  }

  .zero-ton {
    display: none;
  }
}

/* ── Rack đang được chọn — dùng :deep để override Element Plus ── */
:deep(.rack-selected.custom-rack-btn) {
  background: linear-gradient(
    135deg,
    #1e90ff 0%,
    #0052cc 60%,
    #003a99 100%
  ) !important;
  border: 2px solid #60b0ff !important;
  box-shadow: 0 0 10px rgba(30, 144, 255, 0.65),
    0 0 22px rgba(30, 144, 255, 0.35) !important;

  .rack-code {
    font-weight: 700;
    color: #ffffff;
  }

  .rack-ton {
    color: #ffe566;
  }
}

:deep(.rack-selected.custom-rack-btn:hover) {
  background: linear-gradient(
    135deg,
    #3aa0ff 0%,
    #1a6edb 60%,
    #0a4ab8 100%
  ) !important;
}

/* ── Gauge ── */
.gauge-wrapper {
  flex-shrink: 0;
  display: flex;
  gap: 12px;
  align-items: center;
  padding-right: 8px;
  white-space: nowrap;
}

.percentage-value {
  font-size: 14px;
  font-weight: bold;
}

.percentage-label {
  font-size: 11px;
  color: #666;
  text-align: center;
}

.dialog-scrollable-content {
  max-height: 60vh;
  overflow-y: auto;
  padding-top: 10px;
}

/* ── Row highlight ── */
:deep(.el-table .row-3day > td) {
  background-color: #fff3cd;
  font-weight: 600;
}

:deep(.el-table .row-3day:hover > td) {
  background-color: #fff3cd !important;
}

:deep(.el-table .row-over-180 > td) {
  background-color: #fecaca;
  font-weight: 600;
}

:deep(.el-table .row-over-180:hover > td) {
  background-color: #fecaca !important;
}

/* ── DDBH dialog ── */
.ddbh-tag-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 10px;
  max-height: 320px;
  overflow-y: auto;
  padding: 10px;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  background: #fafafa;
}

.ddbh-tag {
  font-size: 18px;
  font-weight: 700;
  padding: 8px 16px;
  border-radius: 6px;
}
</style>
