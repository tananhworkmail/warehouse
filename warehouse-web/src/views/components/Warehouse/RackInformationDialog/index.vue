<template>
  <!-- Dialog: Hiển thị danh sách A01-01 → A01-26 -->
  <el-dialog
    :model-value="dialogVisible"
    :title="`Rack: ${rackNo}`"
    width="70%"
    @close="emit('close')"
  >
    <!-- Phần rack-grid luôn cố định ở trên -->
    <div class="rack-grid-wrapper">
      <div class="rack-grid">
        <el-button
          v-for="rack in subRackList"
          :key="rack"
          @click="handleClick(rack)"
          class="rack-button custom-rack-btn"
        >
          <div class="rack-button-content">
            <div class="rack-code">{{ rack }}</div>
            <div
              class="rack-ton"
              :class="{ 'zero-ton': tonKhoMap[rack.replace('-', '')] == 0 }"
              v-show="tonKhoMap[rack.replace('-', '')] > 0"
            >
              Tồn: {{ tonKhoMap[rack.replace("-", "")] ?? 0 }}
            </div>
          </div>
        </el-button>
      </div>


      <div
        class="rack-toolbar"
        style="display: flex; align-items: center; margin: 10px 0 0 0px"
      >
        <el-checkbox v-model="state.hideZero" style="margin: 10px 0 0 10px">
          Ẩn các dòng có tồn = 0
        </el-checkbox>

        <!-- Legend 3-day -->
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
              background: #fff3cd;
              border: 1px solid #e6d8a2;
              display: inline-block;
            "
          ></span>

          <span style="white-space: nowrap">
            Kế hoạch 3 ngày / 3-Day Plan
          </span>
        </div>
      </div>
    </div>

    <!-- Phần nội dung có thể scroll riêng -->
    <div class="dialog-scrollable-content">
      <el-tabs v-model="state.tabsActiveName" v-if="state.rackTableData.length">
        <el-tab-pane label="Rack Information" name="rack">
          <el-table
            v-loading="loading"
            :data="filteredTableData"
            :row-class-name="tableRowClassName"
            row-key="rowKey"
            max-height="45vh"
          >
            <el-table-column
              v-for="item in state.materialTableFieldList"
              :key="item.prop"
              :property="item.prop"
              :label="item.label"
              align="center"
              header-align="center"
            />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </div>
  </el-dialog>
</template>

<script setup>
import { computed, defineProps, defineEmits, reactive, ref, watch } from "vue";
import axios from "axios";
import { useLoading } from "@/hooks/useLoading";

const API_URL = import.meta.env.VITE_API_URL;
const { loading, showLoading, hideLoading } = useLoading();

const props = defineProps({
  rackNo: String,
  dialogVisible: Boolean,
});
const rack3DaySet = ref(new Set());

const emit = defineEmits(["close"]);

// Tạo danh sách rack con A01-01 → A01-26
const subRackList = computed(() =>
  Array.from(
    { length: 26 },
    (_, i) => `${props.rackNo}-${String(i + 1).padStart(2, "0")}`,
  ),
);

// Map chứa tồn kho: { A0101: 100, A0102: 50, ... }
const tonKhoMap = ref({});

// Dữ liệu bảng
const state = reactive({
  tabsActiveName: "rack",
  rackTableData: [],
  hideZero: true,
  materialTableFieldList: [
    { label: "DDBH", prop: "ddbh" },
    { label: "XieMing", prop: "XieMing" },
    { label: "Size", prop: "XXCC" },
    { label: "Tồn", prop: "TonKho" },
    { label: "Đơn vị", prop: "DonVi" },
    { label: "Mã kệ", prop: "Ma_Ke" },
  ],
});

// Khi chọn A01-03, A01-15...
const handleClick = async (displayRackCode) => {
  const rackCode = displayRackCode.replace("-", ""); // A01-01 → A0101
  rack3DaySet.value = new Set();
  state.rackTableData = [];

  await Promise.all([
    getRackInformation(rackCode),
    get3DayRackInformation(rackCode),
  ]);
};
const filteredTableData = computed(() => {
  let data = state.rackTableData;

  if (state.hideZero) {
    data = data.filter((item) => item.TonKho !== 0);
  }

  data = [...data].sort((a, b) => {
    const a3 = is3DayItem(a);
    const b3 = is3DayItem(b);

    if (a3 === b3) return 0;
    return a3 ? -1 : 1;
  });

  return data;
});
const fetchTonKhoList = async () => {
  try {
    const res = await axios.get(`${API_URL}/warehouse/tonkho`, {
      params: { rackNo: props.rackNo },
    });

    const data = res.data.data || [];

    tonKhoMap.value = data.reduce((acc, row) => {
      acc[row.MA_KE] = row.TonKhoTang;
      return acc;
    }, {});
  } catch (err) {
    console.error("❌ Lỗi lấy tổng tồn kho:", err);
  }
};

const getRackInformation = async (rackNo) => {
  showLoading();
  try {
    const res = await axios.get(`${API_URL}/warehouse/rack`, {
      params: { rackNo: rackNo.trim() },
    });

    const rawData = Array.isArray(res.data.data) ? res.data.data : [];

    // Ép kiểu TonKho sang số
    state.rackTableData = rawData.map((item, index) => ({
      ...item,
      TonKho: Number(item.TonKho),
      rowKey: `${item.ysbh}_${item.Ma_Ke}_${item.XXCC}_${index}`,
    }));
  } catch (err) {
    console.error("❌ Failed to load rack info:", err);
  } finally {
    hideLoading();
  }
};
const get3DayRackInformation = async (rackNo) => {
  try {
    const res = await axios.get(`${API_URL}/warehouse/3dayrack`, {
      params: { rackNo: rackNo.trim() },
    });

    const rawData = Array.isArray(res.data.data) ? res.data.data : [];

    rack3DaySet.value = new Set(
      rawData.map((item) => `${item.ry}_${item.Ma_Ke}`),
    );
  } catch (err) {
    console.error("❌ Failed to load 3day rack info:", err);
  }
};
const is3DayItem = (row) => {
  const key = `${row.ysbh}_${row.Ma_Ke}`;
  return rack3DaySet.value.has(key);
};
const tableRowClassName = ({ row }) => {
  if (is3DayItem(row)) {
    return "row-3day";
  }
  return "";
};
// Tự động fetch tồn kho mỗi khi rackNo thay đổi
watch(
  () => props.rackNo,
  () => {
    if (props.rackNo) {
      fetchTonKhoList();
    }
  },
  { immediate: true },
);
</script>

<style scoped lang="scss">
.el-table {
  user-select: text;
}
.rack-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 10px;
  padding: 10px;
}
.custom-rack-btn {
  background-color: #45484d; // màu nền đậm dễ nhìn
  color: #ffffff !important; // chữ trắng
  border: none;
  padding: 10px;
  height: 60px;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;

  &:hover {
    background-color: #34495e; // hover nhạt hơn
  }

  .rack-button-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    font-size: 14px;
  }

  .rack-code {
    font-size: 14px;
  }

  .rack-ton {
    font-size: 13px;
    color: #ffd700;
    font-weight: bold;
    margin-top: 4px;
  }
}
.rack-grid-wrapper {
  // position: sticky;
  // top: 0;
  background: white;
  z-index: 10;
  padding: 10px 0;
  border-bottom: 1px solid #ebeef5;

  .el-checkbox {
    margin-left: 10px;
    margin-top: 8px;
  }
}

:deep(.el-table .row-3day > td) {
  background-color: #fff3cd !important;
  font-weight: 600;
}

:deep(.el-table .row-3day:hover > td) {
  background-color: #fff3cd !important;
}
:deep(.el-table thead th) {
  background-color: #3b82f6 !important; /* xanh */
  color: white !important;
  font-weight: bold;
  text-align: center;
} 
</style>
