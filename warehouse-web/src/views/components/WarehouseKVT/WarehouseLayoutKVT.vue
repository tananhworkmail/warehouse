<template>
  <div class="mes-container">
    <el-container>
      <el-container class="container">
        <el-header class="header fixed" height="70px">
          <img
            src="@/assets/Logo.png"
            alt="LAIYIH"
            style="height: 60px; margin-left: 20px; margin-top: 5px"
          />
          <h1 class="currentTime">{{ currentTime }}</h1>
          <div class="action">
            <div class="lang-switch">
              <button
                v-for="item in langOptions"
                :key="item.code"
                class="lang-btn"
                :class="{ 'is-active': lang === item.code }"
                @click="setLang(item.code)"
              >
                {{ item.label }}
              </button>
            </div>

            <!-- Ẩn khi ở trang /laboratory -->
            <div class="search" v-if="!isLaboratoryPage" @click="handleSearch">
              <span>{{ t("common.search") }}</span>
              <el-icon :size="24">
                <Search />
              </el-icon>
            </div>

            <div
              class="alert-btn"
              v-if="!isLaboratoryPage"
              @click="handleAlert"
            >
              <span>{{ t("common.alert") }}</span>
              <el-badge is-dot :hidden="!isBellActive">
                <el-icon :size="24" :class="{ shaking: isBellActive }">
                  <Bell />
                </el-icon>
              </el-badge>
            </div>

            <el-dropdown @command="handleCommand">
              <span class="dropdown-trigger">
                {{ t("common.menu") }}
                <el-icon :size="24">
                  <Menu />
                </el-icon>
              </span>

              <template #dropdown>
                <el-dropdown-menu>
                  <template v-if="isWarehousePage">
                    <el-dropdown-item command="organizational">
                      {{ t("kvt.organization") }}
                    </el-dropdown-item>
                    <el-dropdown-item command="schedule">
                      {{ t("kvt.supplier") }}
                    </el-dropdown-item>
                    <el-dropdown-item command="inventory">
                      {{ t("kvt.inventory") }}
                    </el-dropdown-item>
                    <el-dropdown-item command="qrmake">
                      {{ t("kvt.qrMake") }}
                    </el-dropdown-item>
                  </template>

                  <el-dropdown-item command="back">
                    <div class="dropdown-back">
                      <img
                        src="@/assets/undo.png"
                        :alt="t('common.back')"
                        class="undo-icon"
                      />
                      <span>{{ t("common.back") }}</span>
                    </div>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>

            <!-- Nút phóng to / thu nhỏ toàn màn hình -->
            <div
              class="fullscreen-btn"
              @click="toggleFullscreen"
              :title="
                isFullscreen
                  ? t('common.exitFullscreen')
                  : t('common.fullscreen')
              "
            >
              <el-icon :size="24" v-if="!isFullscreen"><FullScreen /></el-icon>
              <el-icon :size="24" v-else><CloseBold /></el-icon>
            </div>
          </div>
        </el-header>

        <el-main class="main">
          <div class="app-main-container">
            <slot></slot>
          </div>
        </el-main>
      </el-container>
    </el-container>

    <!-- Search Dialog -->
    <el-dialog
      :loading="isLoading"
      v-model="state.searchDialogVisible"
      :title="activeSearchTab === 'n31' ? 'Search N31' : t('kvt.searchTitle')"
      width="95%"
      @close="state.searchDialogVisible = false"
    >
      <el-tabs v-model="activeSearchTab" class="search-tabs">
        <el-tab-pane :label="t('kvt.materialSearchTab')" name="material">
          <div class="search-toolbar">
            <el-input
              v-model="searchCLBH"
              :placeholder="t('kvt.searchPlaceholder')"
              style="width: 360px"
              clearable
              @keyup.enter="fetchSearchResults"
            >
              <template #prepend>
                <el-button type="danger" @click="clearSearchResults">
                  {{ t("common.clear") }}
                </el-button>
              </template>
              <template #append>
                <el-button @click="fetchSearchResults">
                  {{ t("common.query") }}
                </el-button>
              </template>
            </el-input>
          </div>

          <div class="search-result-table" v-loading="isLoading">
            <el-table
              :data="filteredResults"
              border
              stripe
              style="width: 100%"
              max-height="64vh"
            >
              <el-table-column
                prop="clbh"
                :label="t('rackDialog.table.columns.clbh')"
              />
              <el-table-column
                prop="zsbh"
                :label="t('rackDialog.table.columns.zsbh')"
              />
              <el-table-column prop="dwbh" :label="t('kvt.unit')" />
              <el-table-column prop="qty" :label="t('kvt.quantity')" />
              <el-table-column prop="remqty" :label="t('kvt.stock')" />
              <el-table-column prop="dqty" :label="t('kvt.outbound')" />
              <el-table-column prop="pack" label="Pack" sortable />
              <el-table-column prop="barcode" label="Barcode" />
              <el-table-column prop="kcbh" :label="t('kvt.rackCode')" />
            </el-table>
          </div>
        </el-tab-pane>

        <el-tab-pane label="Search N31" name="n31">
          <div class="search-toolbar n31-toolbar">
            <el-input
              v-model="n31SearchKHPO"
              :placeholder="t('kvt.n31POPlaceholder')"
              class="n31-search-input"
              clearable
              @keyup.enter="fetchN31SearchResults"
            >
              <template #prepend>{{ t("kvt.n31PO") }}</template>
            </el-input>

            <el-input
              v-model="n31SearchZLBH"
              :placeholder="t('kvt.n31OrderPlaceholder')"
              class="n31-search-input"
              clearable
              @keyup.enter="fetchN31SearchResults"
            >
              <template #prepend>{{ t("kvt.n31OrderNo") }}</template>
            </el-input>

            <el-button type="primary" @click="fetchN31SearchResults">
              {{ t("common.query") }}
            </el-button>
            <el-button type="danger" @click="clearN31SearchResults">
              {{ t("common.clear") }}
            </el-button>
            <el-checkbox v-model="n31RackOnly">
              {{ t("kvt.n31RackOnly") }}
            </el-checkbox>
          </div>

          <div class="search-result-table" v-loading="isLoading">
            <el-table
              :data="filteredN31SearchResults"
              border
              stripe
              style="width: 100%"
              max-height="64vh"
            >
              <el-table-column
                v-for="column in n31Columns"
                :key="column.prop"
                :prop="column.prop"
                :label="column.label"
                :min-width="column.minWidth"
                show-overflow-tooltip
              />
            </el-table>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>

    <OrganizationalDialog
      v-if="state.organizationalDialogVisible"
      :organizationalDialogVisible="state.organizationalDialogVisible"
      @close="(value) => (state.organizationalDialogVisible = value)"
    />

    <TraceDialog
      v-if="state.traceDialogVisible"
      :traceDialogVisible="state.traceDialogVisible"
      @close="(value) => (state.traceDialogVisible = value)"
    />

    <ScheduleDialog
      v-if="state.scheduleDialogVisible"
      :scheduleDialogVisible="state.scheduleDialogVisible"
      @close="(value) => (state.scheduleDialogVisible = value)"
    />

    <LossTemSize
      v-if="state.lossTemSizeDialogVisible"
      :lossTemSizeDialogVisible="state.lossTemSizeDialogVisible"
      @close="(value) => (state.lossTemSizeDialogVisible = value)"
    />

    <AlertKVT v-model="state.alertDialogVisible" />
  </div>
</template>

<script>
export default {
  name: "WarehouseLayoutKVT",
};
</script>

<script setup>
import {
  ref,
  reactive,
  computed,
  onBeforeMount,
  onUnmounted,
  onMounted,
} from "vue";
import { useRouter, useRoute } from "vue-router";

import OrganizationalDialog from "@/views/components/WarehouseKVT/OrganizationalDialog/index.vue";
import ScheduleDialog from "@/views/components/WarehouseKVT/ScheduleDialogKVT/index.vue";
import TraceDialog from "@/views/components/WarehouseKVT/TraceKVT/Scan/index.vue";
import AlertKVT from "@/views/components/WarehouseKVT/AlertDialogKVT/AlertKVT.vue";
import LossTemSize from "@/views/components/WarehouseKVT/LossTemSize/index.vue";

import {
  Search,
  Menu,
  Bell,
  FullScreen,
  CloseBold,
} from "@element-plus/icons-vue";
import axios from "axios";
import dayjs from "dayjs";
// ✅ SỬA: chỉ import 1 lần từ đúng hook, bỏ import trùng từ i18n
import { useWarehouseMapI18n } from "@/hooks/useWarehouseMapI18n";

const router = useRouter();
const route = useRoute();
const API_URL = import.meta.env.VITE_API_URL;
const { lang, langOptions, setLang, t } = useWarehouseMapI18n();

// Kiểm tra route
const isLaboratoryPage = computed(() => route.path.includes("/laboratory"));
const isWarehousePage = computed(() => route.path.includes("/warehouse-kvt"));

const isBellActive = ref(false);
const isLoading = ref(false);
const makeFilter = ref("");

// Fullscreen
const isFullscreen = ref(false);

const goBack = () => {
  router.push("/");
};

let evt = null;
let timer = null;

const updateFullscreenStatus = () => {
  isFullscreen.value = !!(
    document.fullscreenElement ||
    document.webkitFullscreenElement ||
    document.msFullscreenElement
  );
};

const toggleFullscreen = () => {
  const elem = document.documentElement;
  if (!isFullscreen.value) {
    if (elem.requestFullscreen) elem.requestFullscreen();
    else if (elem.webkitRequestFullscreen) elem.webkitRequestFullscreen();
    else if (elem.msRequestFullscreen) elem.msRequestFullscreen();
  } else {
    if (document.exitFullscreen) document.exitFullscreen();
    else if (document.webkitExitFullscreen) document.webkitExitFullscreen();
    else if (document.msExitFullscreen) document.msExitFullscreen();
  }
};

onMounted(() => {
  evt = new EventSource(`${API_URL}/sse`);
  evt.onmessage = (e) => {
    const data = JSON.parse(e.data);
    if (data.type === "HIGH" || data.type === "LOW") {
      isBellActive.value = true;
    }
    if (data.type === "OFF") {
      isBellActive.value = false;
    }
  };

  document.addEventListener("fullscreenchange", updateFullscreenStatus);
  document.addEventListener("webkitfullscreenchange", updateFullscreenStatus);
  document.addEventListener("msfullscreenchange", updateFullscreenStatus);
  updateFullscreenStatus();
});

onUnmounted(() => {
  if (evt) evt.close();
  if (timer) clearInterval(timer);
  document.removeEventListener("fullscreenchange", updateFullscreenStatus);
  document.removeEventListener(
    "webkitfullscreenchange",
    updateFullscreenStatus,
  );
  document.removeEventListener("msfullscreenchange", updateFullscreenStatus);
});

const state = reactive({
  searchDialogVisible: false,
  alertDialogVisible: false,
  organizationalDialogVisible: false,
  scheduleDialogVisible: false,
  traceDialogVisible: false,
  lossTemSizeDialogVisible: false,
});

const currentTime = ref(dayjs().format("YYYY/MM/DD HH:mm:ss"));

const syncCurrentTime = () => {
  timer = setInterval(() => {
    currentTime.value = dayjs().format("YYYY/MM/DD HH:mm:ss");
  }, 1000);
};

const openSchedule = () => {
  state.scheduleDialogVisible = true;
};

const openOrganizational = () => {
  state.organizationalDialogVisible = true;
};

const handleAlert = () => {
  state.alertDialogVisible = true;
  isBellActive.value = false;
};

const handleSearch = () => {
  state.searchDialogVisible = true;
};

// Search
const activeSearchTab = ref("material");
const searchCLBH = ref("");
const searchResults = ref([]);
const n31SearchKHPO = ref("");
const n31SearchZLBH = ref("");
const n31SearchResults = ref([]);
const n31RackOnly = ref(false);
const totalqtyin = ref(0);
const totalqtyout = ref(0);
const totalorderqty = ref(0);

const n31Columns = computed(() => [
  { prop: "buyno", label: t("kvt.n31Columns.buyOrder"), minWidth: 130 },
  { prop: "khpo", label: t("kvt.n31Columns.customerPO"), minWidth: 140 },
  { prop: "zlbh", label: t("kvt.n31Columns.ry"), minWidth: 130 },
  { prop: "article", label: t("kvt.n31Columns.sku"), minWidth: 120 },
  { prop: "xieming", label: t("kvt.n31Columns.shoeName"), minWidth: 170 },
  { prop: "pairs", label: t("kvt.n31Columns.poQty"), minWidth: 105 },
  { prop: "clbh", label: t("kvt.n31Columns.materialCode"), minWidth: 145 },
  {
    prop: "ywpm",
    label: t("kvt.n31Columns.materialDescription"),
    minWidth: 240,
  },
  { prop: "dwbh", label: t("kvt.n31Columns.unit"), minWidth: 80 },
  { prop: "cqdh", label: t("kvt.n31Columns.location"), minWidth: 105 },
  { prop: "clsl", label: t("kvt.n31Columns.qtyUsage"), minWidth: 115 },
  { prop: "cgqty", label: t("kvt.n31Columns.qtyOrdered"), minWidth: 125 },
  { prop: "rkqty", label: t("kvt.n31Columns.qtyInput"), minWidth: 135 },
  { prop: "usestock", label: t("kvt.n31Columns.useStock"), minWidth: 110 },
  { prop: "cgno", label: t("kvt.n31Columns.materialPO"), minWidth: 125 },
  { prop: "ddgb", label: t("kvt.n31Columns.destination"), minWidth: 140 },
  { prop: "zsbh", label: t("kvt.n31Columns.vendorCode"), minWidth: 125 },
  { prop: "zsywjc", label: t("kvt.n31Columns.vendorName"), minWidth: 180 },
  { prop: "invoice", label: t("kvt.n31Columns.invoice"), minWidth: 135 },
  { prop: "make", label: t("kvt.n31Columns.rackCode"), minWidth: 105 },
]);

const filteredN31SearchResults = computed(() => {
  if (!n31RackOnly.value) return n31SearchResults.value;
  return n31SearchResults.value.filter((item) => item.make?.trim());
});

const filteredResults = computed(() => {
  if (!makeFilter.value) return searchResults.value;
  return searchResults.value.filter((item) =>
    item.make?.toLowerCase().includes(makeFilter.value.toLowerCase()),
  );
});

const fetchSearchResults = async () => {
  if (!searchCLBH.value) return;
  isLoading.value = true;
  try {
    const res = await axios.get(`${API_URL}/warehouse/searchkvt`, {
      params: { CLBH: searchCLBH.value },
    });
    searchResults.value = res.data.data || [];
  } catch (error) {
    console.error("Lỗi khi gọi API:", error);
  } finally {
    isLoading.value = false;
  }
};

const clearSearchResults = () => {
  searchCLBH.value = "";
  searchResults.value = [];
  makeFilter.value = "";
  totalqtyin.value = 0;
  totalqtyout.value = 0;
  totalorderqty.value = 0;
};

const fetchN31SearchResults = async () => {
  const po = n31SearchKHPO.value.trim();
  const orderNo = n31SearchZLBH.value.trim();
  if (!po && !orderNo) return;

  isLoading.value = true;
  try {
    const res = await axios.get(`${API_URL}/warehouse/searchn31`, {
      params: { po, orderNo },
    });
    n31SearchResults.value = res.data.data || [];
  } catch (error) {
    n31SearchResults.value = [];
    console.error("Lỗi khi gọi API Search N31:", error);
  } finally {
    isLoading.value = false;
  }
};

const clearN31SearchResults = () => {
  n31SearchKHPO.value = "";
  n31SearchZLBH.value = "";
  n31SearchResults.value = [];
  n31RackOnly.value = false;
};

const handleCommand = (command) => {
  switch (command) {
    case "back":
      goBack();
      break;
    case "inventory":
      window.open("https://192.168.71.9:8100/", "_blank");
      break;
    case "organizational":
      openOrganizational();
      break;
    case "schedule":
      openSchedule();
      break;
    case "lossTemSize":
      router.push({ path: "/loss-tem-size", query: { from: "warehouse-kvt" } });
      break;
    case "qrmake": {
      const link = document.createElement("a");
      link.href = "/QRMAKE.rar"; // đặt file vào public/QR/QRMAKE.zip
      link.download = "QRMAKE.rar";
      link.click();
      break;
    }
  }
};

onBeforeMount(() => {
  syncCurrentTime();
});
</script>

<style lang="scss" scoped>
$base-z-index-999: 999;

.fullscreen-btn {
  cursor: pointer;
  color: white;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  transition: color 0.2s;
}

.fullscreen-btn:hover {
  color: #409eff;
}

.search-toolbar {
  position: sticky;
  top: 0;
  background: white;
  z-index: 10;
  padding-bottom: 8px;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;
  border-bottom: 1px solid #e4e7ed;
}

.search-result-table {
  max-height: 64vh;
  overflow: auto;
}

.n31-toolbar {
  gap: 12px;
}

.n31-search-input {
  width: min(440px, 100%);
}

:deep(.search-tabs > .el-tabs__header) {
  margin-bottom: 16px;
}

.el-table {
  user-select: text;
}

.dropdown-back {
  display: flex;
  align-items: center;
  gap: 8px;
}

.undo-icon {
  width: 22px;
  height: 22px;
  object-fit: contain;
}

:deep(.el-dropdown-menu__item:hover .undo-icon) {
  transform: scale(1.08);
  transition: transform 0.2s ease;
}

.action {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  gap: 20px;
}

.search {
  cursor: pointer;
  color: white;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.alert-btn {
  cursor: pointer;
  color: white;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.lang-switch {
  display: flex;
  gap: 4px;
  padding: 3px;
  border: 1px solid rgba(255, 255, 255, 0.28);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.1);
}

.lang-btn {
  min-width: 38px;
  height: 28px;
  padding: 0 10px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: rgba(255, 255, 255, 0.78);
  font-size: 12px;
  font-weight: 800;
  cursor: pointer;
}

.lang-btn.is-active {
  background: #fff;
  color: #241f20;
}

.dropdown-trigger {
  cursor: pointer;
  color: white;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.dropdown-trigger:hover {
  color: #409eff;
}

.summary-box {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
  padding: 4px 8px;
  background-color: #f9f9f9;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  min-height: 30px;
  align-items: center;
}

.summary-item {
  display: flex;
  align-items: center;
  min-width: 120px;
}

.label {
  font-weight: bold;
  color: #606266;
  margin-right: 4px;
}

.value {
  color: #409eff;
  font-weight: bold;
}

.mes-container {
  user-select: none;

  .header {
    padding: 0;
    background-color: #241f20;

    .currentTime {
      position: absolute;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      color: white;
      font-weight: 900;
      text-align: center;
    }

    .menu {
      position: absolute;
      top: 5px;
      right: 20px;

      h2 {
        color: #fff;
        font-size: 18px;
        font-weight: 800;
      }
    }

    &.fixed {
      position: relative;
      width: 100%;
      z-index: $base-z-index-999;
    }
  }

  .main {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 70px);
    padding: 0;

    .app-main-container {
      flex: 1;
      overflow: auto;
      box-sizing: border-box;
      width: 100%;
      text-align: left;

      .app-main-height {
        height: 100%;
      }

      .footer-copyright {
        height: 40px;
        line-height: 40px;
        flex-shrink: 0;
        color: #999;
        text-align: center;
        border-top: 1px dashed #ccc;
      }
    }
  }
}

:deep(.el-table thead th) {
  background-color: #3b82f6 !important;
  color: white !important;
  font-weight: bold;
  text-align: center;
}

@keyframes shake {
  0% {
    transform: rotate(0deg);
  }
  25% {
    transform: rotate(15deg);
  }
  50% {
    transform: rotate(-15deg);
  }
  75% {
    transform: rotate(10deg);
  }
  100% {
    transform: rotate(0deg);
  }
}

:deep(.shaking) {
  animation: shake 0.5s infinite;
  color: #f56c6c;
}
@media (max-width: 768px) {
  .mes-container {
    .header {
      /* Thu nhỏ logo trên mobile */
      img {
        height: 40px !important;
        margin-left: 10px !important;
        margin-top: 15px !important;
      }

      /* Ẩn đồng hồ để có không gian cho các nút thao tác */
      .currentTime {
        display: none;
      }

      /* Thu gọn khu vực chứa các nút */
      .action {
        gap: 12px; /* Giảm khoảng cách giữa các nút */
        right: 10px;

        /* Ẩn chữ của các nút, chỉ hiển thị Icon để tránh rớt dòng */
        .search span,
        .alert-btn span {
          display: none;
        }

        /* Xử lý riêng cho nút Menu dropdown để ẩn chữ nhưng giữ icon */
        .dropdown-trigger {
          font-size: 0;
        }

        /* Giảm kích thước khối chuyển đổi ngôn ngữ */
        .lang-switch {
          padding: 2px;
        }

        .lang-btn {
          min-width: 30px;
          height: 24px;
          padding: 0 6px;
          font-size: 11px;
        }
      }
    }
  }
}
</style>
