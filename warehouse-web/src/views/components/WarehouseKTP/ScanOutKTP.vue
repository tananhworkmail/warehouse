<template>
  <main class="scan-out-page">
    <header class="topbar">
      <h1>Scan Out KTP</h1>
      <div class="topbar-actions">
        <span class="user-pill"><el-icon><UserFilled /></el-icon>{{ currentUser }}</span>
        <el-dropdown trigger="click" @command="setLang">
          <el-button plain>{{ currentLanguage }}</el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item v-for="item in langOptions" :key="item.code" :command="item.code">
                {{ item.label }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button circle plain type="primary" :title="text.back" @click="goBack"><el-icon><Back /></el-icon></el-button>
        <el-button circle plain type="danger" :title="text.logout" @click="logout"><el-icon><SwitchButton /></el-icon></el-button>
      </div>
    </header>

    <section class="scan-toolbar">
      <el-radio-group v-model="scanMode" size="large" :disabled="busy" @change="focusScanner">
        <el-radio-button value="recycle">{{ text.recycle }} · SB 2</el-radio-button>
        <el-radio-button value="inspection">{{ text.inspection }} · SB 4</el-radio-button>
      </el-radio-group>

      <el-input
        ref="scannerRef"
        v-model="scanInput"
        class="scan-input"
        size="large"
        clearable
        autocomplete="off"
        :disabled="busy || !scanMode"
        :placeholder="scanMode ? text.placeholder : text.chooseMode"
        @input="queueScan"
        @keyup.enter="scanCurrentInput"
      >
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>

      <el-button type="primary" size="large" :loading="processing" :disabled="busy || !scanMode || !scanInput.trim()" @click="scanCurrentInput">
        {{ text.scan }}
      </el-button>
      <el-button size="large" :icon="RefreshRight" :disabled="busy || !sessionScans.length" @click="refreshSession">
        {{ text.refresh }}
      </el-button>
      <div class="counter counter-carton"><span>{{ text.cartons }}</span><b>{{ sessionScans.length }}</b></div>
      <div class="counter counter-pairs"><span>{{ text.pairs }}</span><b>{{ formatNumber(totalQty) }}</b></div>
    </section>

    <section class="data-grid">
      <div class="table-panel session-panel">
        <div class="panel-title">{{ text.session }}</div>
        <el-table
          :data="sessionScans"
          border
          highlight-current-row
          height="100%"
          :empty-text="text.emptySession"
          @current-change="selectCarton"
        >
          <el-table-column prop="cartonBar" label="CARTONBAR" min-width="205" show-overflow-tooltip />
          <el-table-column prop="qty" label="Qty" width="75" align="right">
            <template #default="scope">{{ formatNumber(scope.row.qty) }}</template>
          </el-table-column>
          <el-table-column prop="lpnBarCode" label="LPNBarCode" min-width="165" show-overflow-tooltip />
          <el-table-column prop="ddbh" label="DDBH" width="140" show-overflow-tooltip />
        </el-table>
      </div>

      <div class="table-panel detail-panel">
        <div class="panel-title">
          {{ text.detail }}
          <b v-if="selectedCarton">{{ selectedCarton.cartonBar }}</b>
        </div>
        <el-table :data="selectedDetails" border height="100%" :empty-text="text.emptyDetail">
          <el-table-column prop="orderNo" label="OrderNo" width="135" show-overflow-tooltip />
          <el-table-column prop="sku" label="SKU#" width="125" show-overflow-tooltip />
          <el-table-column prop="modelName" label="Model Name" min-width="160" show-overflow-tooltip />
          <el-table-column prop="xh" label="XH" width="60" align="center" />
          <el-table-column prop="orderSize" label="OrdSize" width="82" align="center" />
          <el-table-column prop="modelSize" label="ModSize" width="85" align="center" />
          <el-table-column prop="qty" label="Qty" width="70" align="right">
            <template #default="scope">{{ formatNumber(scope.row.qty) }}</template>
          </el-table-column>
          <el-table-column prop="from" label="From" width="68" align="center" />
          <el-table-column prop="to" label="To" width="68" align="center" />
        </el-table>
      </div>
    </section>
  </main>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { ElMessage } from "element-plus";
import { Back, RefreshRight, Search, SwitchButton, UserFilled } from "@element-plus/icons-vue";
import { useWarehouseMapI18n } from "@/hooks/useWarehouseMapI18n";
import { clearSession } from "@/utils/appStoreAndSecurity";

const router = useRouter();
const { lang, langOptions, setLang } = useWarehouseMapI18n();
const API_PREFIX = (import.meta.env.VITE_API_URL || "/api/v1").replace(/\/+$/, "");
const API_BASE = `${API_PREFIX}/warehouse-ktp/scan-out`;
const AUTO_SCAN_DELAY = 180;
const MIN_SCAN_LENGTH = 8;
const SESSION_VERSION = 2;

const messages = {
  vi: {
    recycle: "Tái chế", inspection: "Kiểm hàng", placeholder: "Quét CARTONBAR / LPN...", chooseMode: "Chọn Tái chế hoặc Kiểm hàng",
    scan: "Scan", refresh: "Refresh phiên", cartons: "SỐ THÙNG", pairs: "SỐ ĐÔI", session: "Danh sách phiên", detail: "Chi tiết",
    emptySession: "Chưa scan carton", emptyDetail: "Chọn carton để xem chi tiết", back: "Về kho thành phẩm", logout: "Đăng xuất",
    chooseModeError: "Vui lòng chọn Tái chế hoặc Kiểm hàng", duplicate: "Carton đã có trong phiên", differentOrder: "Cảnh báo: thùng thuộc RY {order}, khác RY đang có trong danh sách phiên.", scanFailed: "Scan Out không thành công",
    scanSuccess: "Đã Scan Out {carton} - {mode}", refreshed: "Đã làm mới danh sách phiên",
    restored: "Đã khôi phục danh sách phiên Scan Out",
  },
  en: {
    recycle: "Recycle", inspection: "Inspection", placeholder: "Scan CARTONBAR / LPN...", chooseMode: "Choose Recycle or Inspection",
    scan: "Scan", refresh: "Refresh session", cartons: "CARTONS", pairs: "PAIRS", session: "Session list", detail: "Detail",
    emptySession: "No scanned cartons", emptyDetail: "Select a carton to view details", back: "Back", logout: "Log out",
    chooseModeError: "Choose Recycle or Inspection", duplicate: "Carton already exists in this session", differentOrder: "Warning: carton RY {order} differs from RYs in the session list.", scanFailed: "Scan Out failed",
    scanSuccess: "Scanned out {carton} - {mode}", refreshed: "Session list refreshed",
    restored: "Scan Out session restored",
  },
  zh: {
    recycle: "返工", inspection: "验货", placeholder: "扫描 CARTONBAR / LPN...", chooseMode: "选择返工或验货",
    scan: "扫描", refresh: "刷新本次清单", cartons: "箱数", pairs: "双数", session: "本次清单", detail: "明细",
    emptySession: "尚未扫描纸箱", emptyDetail: "选择纸箱查看明细", back: "返回", logout: "登出",
    chooseModeError: "请选择返工或验货", duplicate: "纸箱已在本次清单", differentOrder: "警告：纸箱 RY {order} 与本次清单中的 RY 不同。", scanFailed: "Scan Out 失败",
    scanSuccess: "已 Scan Out {carton} - {mode}", refreshed: "本次清单已刷新",
    restored: "已恢复 Scan Out 清单",
  },
};

const text = computed(() => messages[lang.value] || messages.vi);
const currentLanguage = computed(() => langOptions.find((item) => item.code === lang.value)?.label || "VI");
const currentUser = computed(() => localStorage.getItem("USERID") || localStorage.getItem("userid") || localStorage.getItem("userId") || "—");
const sessionKey = computed(() => `warehouse-ktp-scan-out-session:${String(currentUser.value).trim().toUpperCase()}`);
const scanMode = ref("");
const scanInput = ref("");
const sessionScans = ref([]);
const selectedCarton = ref(null);
const processing = ref(false);
const scannerRef = ref(null);
let scanTimer = null;
let audioContext = null;

const busy = computed(() => processing.value);
const selectedDetails = computed(() => selectedCarton.value?.details || []);
const totalQty = computed(() => sessionScans.value.reduce((sum, carton) => sum + Number(carton.qty || 0), 0));

function authHeaders() {
  const token = localStorage.getItem("token");
  return token ? { Authorization: `Bearer ${token}` } : {};
}
function errorMessage(error) {
  return error.response?.data?.message || error.message || text.value.scanFailed;
}
function formatNumber(value) {
  const locale = lang.value === "en" ? "en-US" : lang.value === "zh" ? "zh-CN" : "vi-VN";
  return Number(value || 0).toLocaleString(locale, { maximumFractionDigits: 2 });
}
function modeLabel(mode = scanMode.value) {
  return mode === "recycle" ? text.value.recycle : text.value.inspection;
}
function focusScanner() {
  if (scanMode.value) nextTick(() => scannerRef.value?.focus?.());
}
function clearScanTimer() {
  if (scanTimer !== null) window.clearTimeout(scanTimer);
  scanTimer = null;
}
function queueScan(value) {
  clearScanTimer();
  const code = String(value || "").trim();
  if (!scanMode.value || busy.value || code.length < MIN_SCAN_LENGTH) return;
  scanTimer = window.setTimeout(() => {
    scanTimer = null;
    if (String(scanInput.value || "").trim() === code) void scanCarton(code);
  }, AUTO_SCAN_DELAY);
}
function scanCurrentInput() {
  clearScanTimer();
  if (!scanMode.value) {
    ElMessage.warning(text.value.chooseModeError);
    return;
  }
  const code = scanInput.value.trim();
  if (code) void scanCarton(code);
}
async function scanCarton(code) {
  if (!scanMode.value) return false;
  const duplicate = sessionScans.value.some((carton) =>
    [carton.cartonBar, carton.lpnBarCode, carton.scannedCode]
      .some((value) => String(value || "").toUpperCase() === code.toUpperCase()),
  );
  if (duplicate) {
    playTone(false);
    ElMessage.error(text.value.duplicate);
    scanInput.value = "";
    focusScanner();
    return false;
  }

  processing.value = true;
  const mode = scanMode.value;
  try {
    const response = await axios.post(`${API_BASE}/carton`, { cartonBar: code, mode }, { headers: authHeaders() });
    const carton = response.data?.data;
    if (!carton?.cartonBar) throw new Error(text.value.scanFailed);
    const orderNo = String(carton.ddbh || "").trim();
    const hasDifferentOrder = orderNo && sessionScans.value.some((item) => {
      const sessionOrder = String(item.ddbh || "").trim();
      return sessionOrder && sessionOrder.toUpperCase() !== orderNo.toUpperCase();
    });
    if (hasDifferentOrder) {
      ElMessage.warning({
        message: text.value.differentOrder.replace("{order}", orderNo),
        duration: 4500,
        showClose: true,
      });
    }
    const sessionItem = { ...carton, mode, scannedAt: new Date().toISOString() };
    sessionScans.value.unshift(sessionItem);
    selectedCarton.value = sessionItem;
    scanInput.value = "";
    playTone(true);
    ElMessage.success(text.value.scanSuccess.replace("{carton}", carton.cartonBar).replace("{mode}", modeLabel(mode)));
    return true;
  } catch (error) {
    playTone(false);
    scanInput.value = "";
    ElMessage.error({ message: errorMessage(error), duration: 5000, showClose: true });
    return false;
  } finally {
    processing.value = false;
    focusScanner();
  }
}
function selectCarton(carton) {
  if (carton) selectedCarton.value = carton;
}
function refreshSession() {
  sessionScans.value = [];
  selectedCarton.value = null;
  localStorage.removeItem(sessionKey.value);
  ElMessage.success(text.value.refreshed);
  focusScanner();
}
function playTone(success) {
  try {
    const AudioContextClass = window.AudioContext || window.webkitAudioContext;
    if (!AudioContextClass) return;
    audioContext ||= new AudioContextClass();
    const frequencies = success ? [820, 1120] : [190, 190];
    frequencies.forEach((frequency, index) => {
      const oscillator = audioContext.createOscillator();
      const gain = audioContext.createGain();
      const start = audioContext.currentTime + index * 0.14;
      oscillator.type = success ? "sine" : "square";
      oscillator.frequency.value = frequency;
      gain.gain.setValueAtTime(0.12, start);
      gain.gain.exponentialRampToValueAtTime(0.0001, start + 0.12);
      oscillator.connect(gain);
      gain.connect(audioContext.destination);
      oscillator.start(start);
      oscillator.stop(start + 0.13);
    });
  } catch {
    // Audio feedback is optional.
  }
}
function restoreSession() {
  try {
    const saved = JSON.parse(localStorage.getItem(sessionKey.value) || "null");
    if (saved?.version !== SESSION_VERSION || !Array.isArray(saved.scans)) return;
    sessionScans.value = saved.scans.filter((carton) => carton?.cartonBar);
    selectedCarton.value = sessionScans.value[0] || null;
    if (sessionScans.value.length) ElMessage.info(text.value.restored);
  } catch {
    localStorage.removeItem(sessionKey.value);
  }
}
function goBack() {
  void router.push({ name: "warehouse-ktp" });
}
function logout() {
  clearSession();
  void router.replace({
    name: "warehouse-ktp-scan-out-login",
    query: { redirect: "/warehouse-ktp/scan-out" },
  });
}

watch(sessionScans, (scans) => {
  if (scans.length) localStorage.setItem(sessionKey.value, JSON.stringify({ version: SESSION_VERSION, scans }));
  else localStorage.removeItem(sessionKey.value);
}, { deep: true });

onMounted(restoreSession);
onBeforeUnmount(() => {
  clearScanTimer();
  if (audioContext && audioContext.state !== "closed") void audioContext.close();
});
</script>

<style scoped>
.scan-out-page { box-sizing: border-box; min-height: 100vh; padding: 14px; color: #172033; background: #e8efec; }
.topbar, .topbar-actions, .scan-toolbar { display: flex; align-items: center; }
.topbar { justify-content: space-between; height: 44px; margin-bottom: 10px; }
.topbar h1 { margin: 0; color: #173f32; font-size: 24px; text-shadow: 0 1px 0 #fff; }
.topbar-actions { gap: 8px; }
.user-pill { display: inline-flex; align-items: center; gap: 6px; padding: 7px 11px; border: 1px solid #d5dfda; border-radius: 999px; background: #fff; font-weight: 700; }
.scan-toolbar { gap: 8px; padding: 10px; border: 1px solid #bdcec7; border-radius: 11px; background: #fff; box-shadow: 0 6px 18px rgba(31, 60, 48, .09); }
.scan-input { width: clamp(320px, 34vw, 500px); min-width: 260px; flex: 0 1 500px; }
.scan-input :deep(.el-input__wrapper) { border: 1px solid #8faaa0; border-radius: 8px; background: #fbfffd; box-shadow: 0 0 0 2px rgba(23, 107, 80, .08); font-size: 17px; }
.scan-input :deep(.el-input__wrapper.is-focus) { border-color: #16845f; box-shadow: 0 0 0 3px rgba(22, 132, 95, .16); }
.counter { display: flex; min-width: 112px; align-items: center; justify-content: space-between; gap: 12px; padding: 6px 12px; border-radius: 8px; color: #fff; line-height: 1; }
.counter span { font-size: 11px; font-weight: 800; }
.counter b { font-size: 25px; }
.counter-carton { background: #176b50; }
.counter-pairs { background: #24569c; }
.data-grid { display: grid; grid-template-columns: minmax(390px, .82fr) minmax(610px, 1.18fr); gap: 10px; height: calc(100vh - 140px); margin-top: 10px; }
.table-panel { display: flex; min-width: 0; flex-direction: column; padding: 8px; border: 1px solid #bdcec7; border-radius: 11px; background: #fff; box-shadow: 0 5px 16px rgba(31, 60, 48, .08); overflow: hidden; }
.panel-title { display: flex; justify-content: space-between; min-height: 28px; margin: -8px -8px 8px; padding: 8px 11px; align-items: center; border-bottom: 1px solid #d4dfda; background: #edf5f1; color: #29483d; font-size: 14px; font-weight: 850; }
.panel-title b { color: #14845f; }
:deep(.el-table) { flex: 1; border-radius: 6px; }
:deep(.el-table th.el-table__cell) { background: #24483b; color: #fff; }
@media (max-width: 1366px) {
  .scan-out-page { padding: 9px; }
  .topbar { height: 38px; margin-bottom: 7px; }
  .topbar h1 { font-size: 21px; }
  .scan-toolbar { padding: 8px; }
  .scan-input { width: clamp(300px, 31vw, 430px); flex-basis: 430px; }
  .counter { min-width: 96px; padding: 6px 9px; }
  .counter b { font-size: 22px; }
  .data-grid { height: calc(100vh - 116px); margin-top: 7px; }
}
@media (max-width: 1050px) {
  .scan-toolbar { flex-wrap: wrap; }
  .scan-input { order: 2; width: calc(100% - 90px); flex-basis: calc(100% - 100px); }
  .scan-out-page { overflow: auto; }
  .data-grid { grid-template-columns: 1fr; height: auto; }
  .table-panel { height: 390px; }
}
@media (max-width: 640px) {
  .scan-out-page { padding: 8px; }
  .topbar { height: auto; align-items: flex-start; }
  .topbar-actions { flex-wrap: wrap; justify-content: flex-end; }
  .user-pill { display: none; }
  .scan-toolbar { align-items: stretch; }
  .scan-toolbar > * { flex-grow: 1; }
  .scan-input { width: 100%; flex-basis: 100%; }
  .counter { min-width: 0; }
}
@media (max-height: 720px) and (min-width: 1051px) {
  .scan-out-page { padding-top: 7px; padding-bottom: 7px; }
  .topbar { height: 34px; margin-bottom: 5px; }
  .scan-toolbar { padding-top: 7px; padding-bottom: 7px; }
  .data-grid { height: calc(100vh - 101px); margin-top: 6px; }
}
</style>
