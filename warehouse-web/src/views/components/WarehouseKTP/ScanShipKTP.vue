<template>
  <main class="scan-ship-page">
    <header class="topbar">
      <h1>Scan Ship KTP</h1>
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

    <section class="setup-toolbar">
      <strong class="live-clock">{{ clockDisplay }}</strong>
      <el-input v-model="setup.containerNo" class="setup-input" :placeholder="text.container" :disabled="setupLocked || busy" />
      <el-input v-model="setup.plateId" class="setup-input" :placeholder="text.plate" :disabled="setupLocked || busy" />
      <el-button v-if="!scaleConnected" :loading="scaleConnecting" :disabled="!isSerialSupported" @click="connectToScale">
        {{ text.connectScale }}
      </el-button>
      <div class="counter counter-carton"><span>{{ text.cartons }}</span><b>{{ sessionScans.length }}</b></div>
      <div class="counter counter-pairs"><span>{{ text.pairs }}</span><b>{{ formatNumber(totalQty) }}</b></div>
    </section>

    <section class="scan-toolbar">
      <el-input
        ref="scannerRef"
        v-model="scanInput"
        class="scan-input"
        size="large"
        clearable
        autocomplete="off"
        :disabled="busy || !setupReady"
        :placeholder="scanPlaceholder"
        @input="queueScan"
        @keyup.enter="scanCurrentInput"
      >
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
      <strong
        class="weight-reading"
        :class="{ 'is-disconnected': !scaleConnected }"
        :title="scaleConnected ? text.scaleConnected : text.scaleDisconnected"
      >
        <i class="scale-status-dot" aria-hidden="true"></i>
        <span>{{ scaleDisplay }} kg</span>
      </strong>
      <el-button type="primary" size="large" :loading="loadingCarton" :disabled="busy || !setupReady || !scanInput.trim()" @click="scanCurrentInput">
        {{ text.scan }}
      </el-button>
      <el-button size="large" :icon="RefreshRight" :disabled="busy || (!sessionScans.length && !setup.containerNo && !setup.plateId)" @click="refreshSession">
        {{ text.refresh }}
      </el-button>
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
          <el-table-column prop="cartonBar" label="CARTONBAR" min-width="200" show-overflow-tooltip />
          <el-table-column prop="qty" label="Qty" width="68" align="right">
            <template #default="scope">{{ formatNumber(scope.row.qty) }}</template>
          </el-table-column>
          <el-table-column prop="lpnBarCode" label="LPNBarCode" min-width="155" show-overflow-tooltip />
          <el-table-column prop="ddbh" label="DDBH" width="130" show-overflow-tooltip />
          <el-table-column prop="receivedWeight" label="rgw" width="72" align="right" />
          <el-table-column prop="outgoingWeight" label="ogw" width="72" align="right" />
        </el-table>
      </div>

      <div class="table-panel detail-panel">
        <div class="panel-title">
          {{ text.detail }}
          <b v-if="selectedCarton">{{ selectedCarton.cartonBar }}</b>
        </div>
        <el-table :data="selectedDetails" border height="100%" :empty-text="text.emptyDetail">
          <el-table-column prop="orderNo" label="OrderNo" width="130" show-overflow-tooltip />
          <el-table-column prop="sku" label="SKU#" width="120" show-overflow-tooltip />
          <el-table-column prop="modelName" label="Model Name" min-width="155" show-overflow-tooltip />
          <el-table-column prop="xh" label="XH" width="58" align="center" />
          <el-table-column prop="orderSize" label="OrdSize" width="78" align="center" />
          <el-table-column prop="modelSize" label="ModSize" width="82" align="center" />
          <el-table-column prop="qty" label="Qty" width="65" align="right" />
          <el-table-column prop="from" label="From" width="65" align="center" />
          <el-table-column prop="to" label="To" width="65" align="center" />
        </el-table>
      </div>
    </section>
  </main>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { ElMessage, ElMessageBox } from "element-plus";
import { Back, RefreshRight, Search, SwitchButton, UserFilled } from "@element-plus/icons-vue";
import { useWarehouseMapI18n } from "@/hooks/useWarehouseMapI18n";
import { useScaleSerial } from "@/hooks/useScaleSerial";
import { clearSession } from "@/utils/appStoreAndSecurity";

const router = useRouter();
const { lang, langOptions, setLang } = useWarehouseMapI18n();
const API_PREFIX = (import.meta.env.VITE_API_URL || "/api/v1").replace(/\/+$/, "");
const API_BASE = `${API_PREFIX}/warehouse-ktp/scan-ship`;
const AUTO_SCAN_DELAY = 180;
const MIN_SCAN_LENGTH = 8;
const SESSION_VERSION = 2;

const messages = {
  vi: {
    container: "Số container", plate: "Biển số xe", connectScale: "Kết nối cân", scaleDisconnected: "Mất kết nối cân",
    placeholder: "Quét CARTONBAR / LPN để xuất ngay...", setupRequired: "Nhập container và biển số xe", scaleRequired: "Đang chờ kết nối cân", scan: "Scan",
    refresh: "Phiên mới", cartons: "SỐ THÙNG", pairs: "SỐ ĐÔI", session: "Đã xuất trong phiên", detail: "Chi tiết", emptySession: "Chưa xuất carton",
    emptyDetail: "Chọn carton để xem chi tiết", back: "Về kho thành phẩm", logout: "Đăng xuất", setupMissing: "Thiếu container hoặc biển số xe",
    duplicate: "Carton đã có trong phiên", lookupFailed: "Không đọc được carton", invalidCarton: "Carton không đạt điều kiện Scan Ship",
    differentOrder: "Cảnh báo: thùng thuộc RY {order}, khác RY đang có trong danh sách phiên.",
    weightTitle: "Trọng lượng không đạt", weightContinue: "Trọng lượng vượt giới hạn. Vẫn thêm carton vào phiên?", weightMissing: "Chưa đọc được trọng lượng từ cân",
    confirm: "Xác nhận", cancel: "Hủy", shipped: "Đã xuất thùng {carton}", refreshed: "Đã mở phiên mới",
    scaleConnected: "Đã kết nối cân", scaleFailed: "Không kết nối được cân", restored: "Đã khôi phục danh sách phiên Scan Ship",
  },
  en: {
    container: "Container No.", plate: "Plate ID", connectScale: "Connect scale", scaleDisconnected: "Scale disconnected",
    placeholder: "Scan CARTONBAR / LPN to ship now...", setupRequired: "Enter container and plate", scaleRequired: "Waiting for scale connection", scan: "Scan",
    refresh: "New session", cartons: "CARTONS", pairs: "PAIRS", session: "Shipped this session", detail: "Detail", emptySession: "No shipped cartons",
    emptyDetail: "Select a carton to view details", back: "Back", logout: "Log out", setupMissing: "Container or plate is missing",
    duplicate: "Carton already exists in this session", lookupFailed: "Could not read carton", invalidCarton: "Carton is not eligible for Scan Ship",
    differentOrder: "Warning: carton RY {order} differs from RYs in the session list.",
    weightTitle: "Weight failed", weightContinue: "Weight exceeds the limit. Add this carton anyway?", weightMissing: "No scale reading",
    confirm: "Confirm", cancel: "Cancel", shipped: "Shipped carton {carton}", refreshed: "New session started",
    scaleConnected: "Scale connected", scaleFailed: "Could not connect scale", restored: "Scan Ship session restored",
  },
  zh: {
    container: "货柜号", plate: "车牌", connectScale: "连接电子秤", scaleDisconnected: "电子秤未连接",
    placeholder: "扫描 CARTONBAR / LPN 立即出货...", setupRequired: "输入货柜号和车牌", scaleRequired: "等待电子秤连接", scan: "扫描",
    refresh: "新建本次", cartons: "箱数", pairs: "双数", session: "本次已出货", detail: "明细", emptySession: "尚未出货",
    emptyDetail: "选择纸箱查看明细", back: "返回", logout: "登出", setupMissing: "缺少货柜号或车牌",
    duplicate: "纸箱已在本次清单", lookupFailed: "无法读取纸箱", invalidCarton: "纸箱不符合 Scan Ship 条件",
    differentOrder: "警告：纸箱 RY {order} 与本次清单中的 RY 不同。",
    weightTitle: "重量不合格", weightContinue: "重量超过范围，仍加入本次清单？", weightMissing: "尚未读取重量",
    confirm: "确认", cancel: "取消", shipped: "纸箱 {carton} 已出货", refreshed: "已新建本次",
    scaleConnected: "电子秤已连接", scaleFailed: "无法连接电子秤", restored: "已恢复本次 Scan Ship 清单",
  },
};

const text = computed(() => messages[lang.value] || messages.vi);
const currentLanguage = computed(() => langOptions.find((item) => item.code === lang.value)?.label || "VI");
const currentUser = computed(() => localStorage.getItem("USERID") || localStorage.getItem("userid") || localStorage.getItem("userId") || "—");
const sessionKey = computed(() => `warehouse-ktp-scan-ship-session:${String(currentUser.value).trim().toUpperCase()}`);
const setup = reactive({ containerNo: "", plateId: "" });
const clockNow = ref(new Date());
const scanInput = ref("");
const sessionScans = ref([]);
const selectedCarton = ref(null);
const loadingCarton = ref(false);
const scannerRef = ref(null);
let scanTimer = null;
let clockTimer = null;
let audioContext = null;

const {
  currentWeight,
  scaleConnected,
  scaleConnecting,
  isSerialSupported,
  initializeScale,
  connectScale,
  resetWeight,
  disposeScale,
} = useScaleSerial();

const busy = computed(() => loadingCarton.value);
const setupReady = computed(() => [setup.containerNo, setup.plateId].every((value) => String(value || "").trim()));
const weightReady = computed(() => Number(currentWeight.value) > 0);
const scanReady = computed(() => setupReady.value && scaleConnected.value && weightReady.value);
const setupLocked = computed(() => sessionScans.value.length > 0);
const selectedDetails = computed(() => {
  const carton = selectedCarton.value;
  const details = Array.isArray(carton?.details) ? carton.details : [];
  const cartonXH = String(carton?.xh || "").trim().toUpperCase();
  const exactSize = cartonXH
    ? details.filter((detail) => String(detail.xh || "").trim().toUpperCase() === cartonXH)
    : details;
  return (exactSize.length ? exactSize : details).slice(0, 1);
});
const totalQty = computed(() => sessionScans.value.reduce((sum, carton) => sum + Number(carton.qty || 0), 0));
const scaleDisplay = computed(() => currentWeight.value === null ? "--.--" : Number(currentWeight.value).toFixed(2));
const clockDisplay = computed(() => clockNow.value.toLocaleString(
  lang.value === "en" ? "en-GB" : lang.value === "zh" ? "zh-CN" : "vi-VN",
  { hour12: false },
));
const scanPlaceholder = computed(() => {
  if (!setupReady.value) return text.value.setupRequired;
  if (!scaleConnected.value) return text.value.scaleRequired;
  if (!weightReady.value) return text.value.weightMissing;
  return text.value.placeholder;
});
function authHeaders() {
  const token = localStorage.getItem("token");
  return token ? { Authorization: `Bearer ${token}` } : {};
}
function errorMessage(error, fallback = text.value.lookupFailed) {
  return error.response?.data?.message || error.message || fallback;
}
function formatNumber(value) {
  const locale = lang.value === "en" ? "en-US" : lang.value === "zh" ? "zh-CN" : "vi-VN";
  return Number(value || 0).toLocaleString(locale, { maximumFractionDigits: 2 });
}
function focusScanner() {
  if (setupReady.value) nextTick(() => scannerRef.value?.focus?.());
}
function clearScanTimer() {
  if (scanTimer !== null) window.clearTimeout(scanTimer);
  scanTimer = null;
}
function queueScan(value) {
  clearScanTimer();
  const code = String(value || "").trim();
  if (!scanReady.value || busy.value || code.length < MIN_SCAN_LENGTH) return;
  scanTimer = window.setTimeout(() => {
    scanTimer = null;
    if (String(scanInput.value || "").trim() === code) void scanCarton(code);
  }, AUTO_SCAN_DELAY);
}
function scanCurrentInput() {
  clearScanTimer();
  if (!setupReady.value) {
    ElMessage.warning(text.value.setupMissing);
    return;
  }
  if (!scaleConnected.value) {
    ElMessage.error(text.value.scaleDisconnected);
    return;
  }
  const code = scanInput.value.trim();
  if (code) void scanCarton(code);
}
async function scanCarton(code, { silent = false } = {}) {
  const duplicate = sessionScans.value.some((carton) =>
    [carton.cartonBar, carton.lpnBarCode, carton.scannedCode]
      .some((value) => String(value || "").toUpperCase() === code.toUpperCase()),
  );
  if (duplicate) {
    playTone(false);
    if (!silent) ElMessage.error(text.value.duplicate);
    scanInput.value = "";
    focusScanner();
    return false;
  }
  if (!scaleConnected.value || !(Number(currentWeight.value) > 0)) {
    playTone(false);
    if (!silent) ElMessage.error(text.value.weightMissing);
    return false;
  }

  loadingCarton.value = true;
  try {
    const outgoingWeight = Number(currentWeight.value);
    const response = await axios.get(`${API_BASE}/carton/${encodeURIComponent(code)}`, {
      headers: authHeaders(),
      params: {
        containerNo: setup.containerNo.trim(),
        plateId: setup.plateId.trim(),
        weightEnabled: true,
        outgoingWeight,
      },
    });
    const carton = response.data?.data;
    if (!carton?.cartonBar) throw new Error(text.value.lookupFailed);
    if (!carton.eligibilityPassed || !carton.planPassed) {
      throw new Error([carton.eligibilityMessage, carton.planMessage].filter(Boolean).join(" · ") || text.value.invalidCarton);
    }

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

    let weightOverride = false;
    if (!carton.weightCheck?.passed) {
      if (!carton.weightCheck?.requiresOverride) throw new Error(carton.weightCheck?.message || text.value.weightMissing);
      try {
        await ElMessageBox.confirm(
          `${carton.weightCheck.message}. ${text.value.weightContinue}`,
          text.value.weightTitle,
          { type: "warning", confirmButtonText: text.value.confirm, cancelButtonText: text.value.cancel },
        );
        weightOverride = true;
      } catch {
        return false;
      }
    }

    const updateResponse = await axios.post(`${API_BASE}/carton`, {
      containerNo: setup.containerNo.trim(),
      plateId: setup.plateId.trim(),
      cartons: [{
        cartonBar: carton.cartonBar,
        outgoingWeight,
        weightOverride,
      }],
    }, { headers: authHeaders() });
    const savedCarton = updateResponse.data?.data?.cartons?.[0];
    if (!savedCarton?.cartonBar) throw new Error(text.value.lookupFailed);
    const sessionItem = { ...savedCarton, scannedAt: new Date().toISOString() };
    sessionScans.value.unshift(sessionItem);
    selectedCarton.value = sessionItem;
    scanInput.value = "";
    resetWeight();
    playTone(true);
    if (!silent) ElMessage.success(text.value.shipped.replace("{carton}", carton.cartonBar));
    return true;
  } catch (error) {
    scanInput.value = "";
    resetWeight();
    playTone(false);
    if (!silent) ElMessage.error({ message: errorMessage(error), duration: 5000, showClose: true });
    return false;
  } finally {
    loadingCarton.value = false;
    focusScanner();
  }
}
function selectCarton(carton) {
  if (carton) selectedCarton.value = carton;
}
function clearSessionState() {
  sessionScans.value = [];
  selectedCarton.value = null;
  scanInput.value = "";
  setup.containerNo = "";
  setup.plateId = "";
  localStorage.removeItem(sessionKey.value);
}
function refreshSession() {
  clearSessionState();
  ElMessage.success(text.value.refreshed);
}
async function connectToScale() {
  const connected = await connectScale({ reuseSavedPort: true });
  ElMessage[connected ? "success" : "error"](connected ? text.value.scaleConnected : text.value.scaleFailed);
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
function persistSession() {
  if (!sessionScans.value.length) return;
  localStorage.setItem(sessionKey.value, JSON.stringify({
    version: SESSION_VERSION,
    setup: { ...setup },
    scans: sessionScans.value,
  }));
}
function restoreSession() {
  try {
    const saved = JSON.parse(localStorage.getItem(sessionKey.value) || "null");
    if (saved?.version !== SESSION_VERSION || !Array.isArray(saved.scans)) return;
    Object.assign(setup, saved.setup || {});
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
    name: "warehouse-ktp-scan-ship-login",
    query: { redirect: "/warehouse-ktp/scan-ship" },
  });
}

watch([setup, sessionScans], persistSession, { deep: true });
watch(currentWeight, (weight) => {
  if (Number(weight) > 0) {
    focusScanner();
    queueScan(scanInput.value);
  }
});
onMounted(async () => {
  restoreSession();
  clockTimer = window.setInterval(() => { clockNow.value = new Date(); }, 1000);
  await initializeScale({ autoConnect: true });
  focusScanner();
});
onBeforeUnmount(() => {
  clearScanTimer();
  if (clockTimer !== null) window.clearInterval(clockTimer);
  void disposeScale();
  if (audioContext && audioContext.state !== "closed") void audioContext.close();
});
</script>

<style scoped>
.scan-ship-page { box-sizing: border-box; min-height: 100vh; padding: 12px; color: #172033; background: #e8efec; }
.topbar, .topbar-actions, .setup-toolbar, .scan-toolbar { display: flex; align-items: center; }
.topbar { justify-content: space-between; min-height: 42px; margin-bottom: 8px; }
.topbar h1 { margin: 0; color: #173f32; font-size: 23px; text-shadow: 0 1px 0 #fff; }
.topbar-actions, .setup-toolbar, .scan-toolbar { gap: 8px; }
.user-pill { display: inline-flex; align-items: center; gap: 6px; padding: 7px 11px; border: 1px solid #d5dfda; border-radius: 999px; background: #fff; font-weight: 700; }
.setup-toolbar, .scan-toolbar { padding: 8px; border: 1px solid #bdcec7; background: #fff; box-shadow: 0 5px 16px rgba(31, 60, 48, .08); }
.setup-toolbar { border-radius: 9px 9px 0 0; }
.scan-toolbar { border-top: 0; border-radius: 0 0 9px 9px; }
.setup-input { width: 210px; }
.scan-input { width: clamp(320px, 34vw, 500px); min-width: 260px; flex: 0 1 500px; }
.scan-input :deep(.el-input__wrapper) { border: 1px solid #8faaa0; border-radius: 8px; background: #fbfffd; box-shadow: 0 0 0 2px rgba(23, 107, 80, .08); font-size: 17px; }
.scan-input :deep(.el-input__wrapper.is-focus) { border-color: #16845f; box-shadow: 0 0 0 3px rgba(22, 132, 95, .16); }
.live-clock { min-width: 190px; padding: 9px 12px; border-radius: 7px; background: #e8f2ed; color: #173f32; text-align: center; font-size: 16px; }
.weight-reading {
  display: inline-flex;
  min-width: 132px;
  padding: 6px 13px;
  border: 2px solid #39d98a;
  border-radius: 8px;
  background: #0b4935;
  box-shadow: 0 0 0 3px rgba(57, 217, 138, .14), 0 5px 14px rgba(11, 73, 53, .18);
  color: #eaff62;
  align-items: center;
  justify-content: center;
  gap: 8px;
  text-align: center;
  font-family: Consolas, monospace;
  font-size: 21px;
  line-height: 1.2;
  text-shadow: 0 0 10px rgba(234, 255, 98, .3);
}
.scale-status-dot {
  width: 9px;
  height: 9px;
  flex: 0 0 9px;
  border-radius: 50%;
  background: #42f59b;
  box-shadow: 0 0 0 3px rgba(66, 245, 155, .16), 0 0 9px rgba(66, 245, 155, .72);
}
.weight-reading.is-disconnected {
  border-color: #f97066;
  background: #6b211d;
  box-shadow: 0 0 0 3px rgba(249, 112, 102, .14);
  color: #fff2a8;
}
.weight-reading.is-disconnected .scale-status-dot {
  background: #ff6b62;
  box-shadow: 0 0 0 3px rgba(255, 107, 98, .16), 0 0 9px rgba(255, 107, 98, .68);
}
.counter { display: flex; min-width: 112px; align-items: center; justify-content: space-between; gap: 12px; padding: 6px 12px; border-radius: 8px; color: #fff; line-height: 1; }
.counter span { font-size: 11px; font-weight: 800; }
.counter b { font-size: 25px; }
.counter-carton { background: #176b50; }
.counter-pairs { background: #24569c; }
.data-grid { display: grid; grid-template-columns: minmax(500px, .9fr) minmax(610px, 1.1fr); gap: 8px; height: calc(100vh - 166px); margin-top: 8px; }
.table-panel { display: flex; min-width: 0; flex-direction: column; padding: 7px; border: 1px solid #bdcec7; border-radius: 10px; background: #fff; box-shadow: 0 5px 16px rgba(31, 60, 48, .08); overflow: hidden; }
.panel-title { display: flex; justify-content: space-between; min-height: 28px; margin: -7px -7px 7px; padding: 7px 10px; align-items: center; border-bottom: 1px solid #d4dfda; background: #edf5f1; color: #29483d; font-size: 14px; font-weight: 850; }
.panel-title b { color: #14845f; }
:deep(.el-table) { flex: 1; border-radius: 5px; }
:deep(.el-table th.el-table__cell) { background: #24483b; color: #fff; }
@media (max-width: 1366px) {
  .scan-ship-page { padding: 8px; }
  .topbar { min-height: 37px; margin-bottom: 6px; }
  .topbar h1 { font-size: 21px; }
  .setup-toolbar, .scan-toolbar { padding: 7px; }
  .setup-input { width: 170px; }
  .live-clock { min-width: 164px; padding: 8px 9px; font-size: 14px; }
  .scan-input { width: clamp(300px, 31vw, 430px); flex-basis: 430px; }
  .counter { min-width: 96px; padding: 6px 9px; }
  .counter b { font-size: 22px; }
  .data-grid { height: calc(100vh - 142px); margin-top: 7px; }
}
@media (max-width: 1120px) {
  .setup-toolbar, .scan-toolbar { flex-wrap: wrap; }
  .scan-input { flex-basis: 55%; }
  .scan-ship-page { overflow: auto; }
  .data-grid { grid-template-columns: 1fr; height: auto; }
  .table-panel { height: 380px; }
}
@media (max-width: 650px) {
  .scan-ship-page { padding: 7px; }
  .topbar { align-items: flex-start; }
  .user-pill { display: none; }
  .setup-input, .scan-input { width: 100%; flex-basis: 100%; }
  .live-clock { width: 100%; box-sizing: border-box; }
  .counter { min-width: 0; }
}
@media (max-height: 720px) and (min-width: 1121px) {
  .scan-ship-page { padding-top: 6px; padding-bottom: 6px; }
  .topbar { min-height: 33px; margin-bottom: 4px; }
  .setup-toolbar, .scan-toolbar { padding-top: 6px; padding-bottom: 6px; }
  .data-grid { height: calc(100vh - 126px); margin-top: 6px; }
}
</style>
