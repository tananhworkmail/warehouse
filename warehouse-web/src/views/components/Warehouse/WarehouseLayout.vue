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

            <div class="search" @click="handleSearch">
              <span>{{ t("common.search") }}</span>
              <el-icon :size="24">
                <Search />
              </el-icon>
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
                  <el-dropdown-item command="schedule">{{
                    t("warehouse.schedule")
                  }}</el-dropdown-item>
                  <el-dropdown-item command="trace">{{
                    t("warehouse.trace")
                  }}</el-dropdown-item>
                  <el-dropdown-item command="ncu">{{
                    t("warehouse.supplier")
                  }}</el-dropdown-item>

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
            <footer class="footer-copyright">{{ copyrightStr }}</footer>
          </div>
        </el-main>
      </el-container>
    </el-container>

    <el-dialog
      :loading="isLoading"
      v-model="state.searchDialogVisible"
      :title="t('warehouse.searchTitle')"
      width="80%"
      @close="state.searchDialogVisible = false"
    >
      <div class="search-toolbar">
        <el-input
          v-model="searchDDBH"
          :placeholder="t('warehouse.searchPlaceholder')"
          style="width: 300px"
          clearable
          @keyup.enter="fetchSearchResults"
        >
          <template #prepend>
            <el-button type="danger" @click="clearSearchResults">{{
              t("common.clear")
            }}</el-button>
          </template>
          <template #append>
            <el-button @click="fetchSearchResults">{{
              t("common.query")
            }}</el-button>
          </template>
        </el-input>

        <div class="summary-box">
          <div class="summary-item">
            <span class="label">{{ t("warehouse.totalOrder") }}</span>
            <span class="value">{{ totalorderqty.toLocaleString() }}</span>
          </div>
          <div class="summary-item">
            <span class="label">{{ t("warehouse.imported") }}</span>
            <span class="value">{{ totalqtyin.toLocaleString() }}</span>
          </div>
          <div class="summary-item">
            <span class="label">{{ t("warehouse.notImported") }}</span>
            <span class="value">{{ chuanhap.toLocaleString() }}</span>
          </div>
          <div class="summary-item">
            <span class="label">{{ t("warehouse.exported") }}</span>
            <span class="value">{{ totalqtyout.toLocaleString() }}</span>
          </div>
          <div class="summary-item">
            <span class="label">{{ t("warehouse.stock") }}</span>
            <span class="value">{{ tonkho.toLocaleString() }}</span>
          </div>
        </div>
      </div>

      <div class="search-result-table" v-loading="isLoading">
        <el-table
          :data="filteredResults"
          border
          stripe
          style="width: 100%"
          max-height="70vh"
        >
          <el-table-column prop="ddbh" label="DDBH" />
          <el-table-column prop="XieMing" label="XieMing" />
          <el-table-column prop="xxcc" label="XXCC" sortable />
          <el-table-column prop="qtyin" label="Qty In" />
          <el-table-column prop="qtyout" label="Qty Out" />
          <el-table-column prop="tonkho" :label="t('warehouse.stockColumn')" sortable />
          <el-table-column prop="make" :label="t('warehouse.rackCode')">
            <template #header>
              <el-input
                v-model="makeFilter"
                :placeholder="t('warehouse.rackCode')"
                size="medium"
                style="font-weight: bold"
                clearable
              />
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>

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
    <NCUDialog
      v-if="state.ncuDialogVisible"
      :ncuDialogVisible="state.ncuDialogVisible"
      @close="(value) => (state.ncuDialogVisible = value)"
    />
  </div>
</template>

<script>
export default {
  name: "WarehouseLayout",
};
</script>

<script setup>
import {
  ref,
  reactive,
  onBeforeMount,
  onMounted,
  onBeforeUnmount,
  computed,
} from "vue";
import ScheduleDialog from "@/views/components/Warehouse/ScheduleDialog/index.vue";
import NCUDialog from "@/views/components/Warehouse/NCUDialog/index.vue";
import TraceDialog from "@/views/components/Warehouse/Trace/Scan/index.vue";
import { Search, Menu, FullScreen, CloseBold } from "@element-plus/icons-vue";
import axios from "axios";
import dayjs from "dayjs";
import { useRouter } from "vue-router";
import { useWarehouseMapI18n } from "@/hooks/useWarehouseMapI18n";

const router = useRouter();
const API_URL = import.meta.env.VITE_API_URL;
const { lang, langOptions, setLang, t } = useWarehouseMapI18n();

// Các state cơ bản
const isFullscreen = ref(false);
const isLoading = ref(false);
const makeFilter = ref("");
const currentTime = ref(dayjs().format("YYYY/MM/DD HH:mm:ss"));
const copyrightStr = ref("© 2025 YIHSHUO. All rights reserved.");

const state = reactive({
  searchDialogVisible: false,
  traceDialogVisible: false,
  scheduleDialogVisible: false,
  ncuDialogVisible: false,
});

// Các hàm xử lý Dialog
const handleSearch = () => {
  state.searchDialogVisible = true;
};

const openTrace = () => {
  state.traceDialogVisible = true;
};

const openSchedule = () => {
  state.scheduleDialogVisible = true;
};

const openNCU = () => {
  state.ncuDialogVisible = true;
};

const goBack = () => {
  router.push("/");
};

// Xử lý sự kiện Menu Dropdown
const handleCommand = (command) => {
  switch (command) {
    case "schedule":
      openSchedule();
      break;
    case "trace":
      openTrace();
      break;
    case "ncu":
      openNCU();
      break;
    case "back":
      goBack();
      break;
  }
};

// Logic Tìm kiếm DDBH
const searchDDBH = ref("");
const searchResults = ref([]);
const totalqtyin = ref(0);
const totalqtyout = ref(0);
const totalorderqty = ref(0);

const tonkho = computed(() => totalqtyin.value - totalqtyout.value);
const chuanhap = computed(() => totalorderqty.value - totalqtyin.value);

const filteredResults = computed(() => {
  if (!makeFilter.value) return searchResults.value;
  return searchResults.value.filter((item) =>
    item.make?.toLowerCase().includes(makeFilter.value.toLowerCase())
  );
});

const fetchSearchResults = async () => {
  if (!searchDDBH.value) return;

  isLoading.value = true;
  try {
    const res = await axios.get(`${API_URL}/warehouse/search`, {
      params: { DDBH: searchDDBH.value },
    });
    searchResults.value = res.data.data || [];
    await fetchSearchTotal();
  } catch (error) {
    console.error("Lỗi khi gọi API:", error);
  } finally {
    isLoading.value = false;
  }
};

const fetchSearchTotal = async () => {
  try {
    const res = await axios.get(`${API_URL}/warehouse/searchtotal`, {
      params: { DDBH: searchDDBH.value },
    });

    const totalData = res.data.data?.[0];
    if (totalData) {
      totalqtyin.value = Number(totalData.totalqtyin || 0);
      totalqtyout.value = Number(totalData.totalqtyout || 0);
      totalorderqty.value = Number(totalData.totalorderqty || 0);
    }
  } catch (error) {
    console.error("Lỗi khi gọi API:", error);
  }
};

const clearSearchResults = () => {
  searchDDBH.value = "";
  searchResults.value = [];
  makeFilter.value = "";
  totalqtyin.value = 0;
  totalqtyout.value = 0;
  totalorderqty.value = 0;
};

// Xử lý Đồng hồ
let timer = null;
const syncCurrentTime = () => {
  timer = setInterval(() => {
    currentTime.value = dayjs().format("YYYY/MM/DD HH:mm:ss");
  }, 1000);
};

// Xử lý Fullscreen
const updateFullscreenStatus = () => {
  isFullscreen.value = !!(
    document.fullscreenElement ||
    document.webkitFullscreenElement ||
    document.msFullscreenElement
  );
};

const enterFullscreen = async () => {
  const elem = document.documentElement;
  try {
    if (elem.requestFullscreen) await elem.requestFullscreen();
    else if (elem.webkitRequestFullscreen) await elem.webkitRequestFullscreen();
    else if (elem.msRequestFullscreen) await elem.msRequestFullscreen();
  } catch (error) {
    console.error("Không vào được fullscreen:", error);
  }
  updateFullscreenStatus();
};

const exitFullscreen = async () => {
  try {
    if (document.exitFullscreen) await document.exitFullscreen();
    else if (document.webkitExitFullscreen) await document.webkitExitFullscreen();
    else if (document.msExitFullscreen) await document.msExitFullscreen();
  } catch (error) {
    console.error("Không thoát được fullscreen:", error);
  }
  updateFullscreenStatus();
};

const toggleFullscreen = () => {
  isFullscreen.value ? exitFullscreen() : enterFullscreen();
};

const handleFullscreenChange = () => {
  updateFullscreenStatus();
};

// Vòng đời Component
onMounted(() => {
  document.addEventListener("fullscreenchange", handleFullscreenChange);
  document.addEventListener("webkitfullscreenchange", handleFullscreenChange);
  document.addEventListener("MSFullscreenChange", handleFullscreenChange);
  updateFullscreenStatus();
});

onBeforeUnmount(() => {
  if (timer) clearInterval(timer);
  document.removeEventListener("fullscreenchange", handleFullscreenChange);
  document.removeEventListener("webkitfullscreenchange", handleFullscreenChange);
  document.removeEventListener("MSFullscreenChange", handleFullscreenChange);
});

onBeforeMount(() => {
  syncCurrentTime();
});
</script>

<style lang="scss" scoped>
$base-z-index-999: 999;
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
  max-height: 70vh;
  overflow-y: auto;
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

.action {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  gap: 20px; // khoảng cách giữa Tìm kiếm và Menu
}

.search {
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
  $base-width: 100%;
  $app-main-min-height: 400px;
  $footer-copyright-height: 40px;
  $base-color-3: #999;
  $base-border-color: #ccc;

  .main {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 70px); // trừ phần header
    padding: 0;

    .app-main-container {
      flex: 1; // chiếm toàn bộ phần còn lại
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

.dropdown-trigger:hover {
  color: #409eff;
}
:deep(.el-table thead th) {
  background-color: #3b82f6 !important; /* xanh */
  color: white !important;
  font-weight: bold;
  text-align: center;
}
</style>
