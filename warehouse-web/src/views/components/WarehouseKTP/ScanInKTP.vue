<template>
  <main class="scan-in-page">
    <header class="topbar">
      <h1>Scan In KTP</h1>
      <div class="page-actions">
        <div class="current-user" :title="`${t('ktpScanIn.userId')}: ${currentUserId}`">
          <el-icon><UserFilled /></el-icon>
          <span>{{ t("ktpScanIn.userId") }}</span>
          <strong>{{ currentUserId }}</strong>
        </div>
        <el-tooltip :content="`${t('ktpScanIn.actions.language')}: ${currentLanguageLabel}`" placement="bottom">
          <el-dropdown trigger="click" @command="handleLanguageChange">
            <el-button class="language-button" size="small" plain :aria-label="t('ktpScanIn.actions.language')">
              🌐 {{ currentLanguageLabel }}
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item
                  v-for="option in langOptions"
                  :key="option.code"
                  :command="option.code"
                  :disabled="option.code === lang"
                >
                  {{ option.label }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </el-tooltip>
        <el-tooltip :content="t('ktpScanIn.actions.backToWarehouse')" placement="bottom">
          <el-button
            class="action-icon-button"
            type="primary"
            plain
            circle
            :aria-label="t('ktpScanIn.actions.backToWarehouse')"
            @click="handleBackToWarehouse"
          >
            <el-icon><Back /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip :content="t('ktpScanIn.actions.logout')" placement="bottom">
          <el-button
            class="action-icon-button"
            type="danger"
            plain
            circle
            :aria-label="t('ktpScanIn.actions.logout')"
            @click="handleLogout"
          >
            <el-icon><SwitchButton /></el-icon>
          </el-button>
        </el-tooltip>
      </div>
    </header>

    <section class="scan-toolbar">
      <div class="scan-mode-row">
        <span>{{ t("ktpScanIn.scanMode.label") }}</span>
        <el-radio-group v-model="scanMode" size="small" :disabled="loadingCarton" @change="handleScanModeChange">
          <el-radio-button :value="SCAN_MODE.CARTON">
            {{ t("ktpScanIn.scanMode.carton") }}
          </el-radio-button>
          <el-radio-button :value="SCAN_MODE.CARTON_INBOX">
            {{ t("ktpScanIn.scanMode.cartonInbox") }}
          </el-radio-button>
        </el-radio-group>
        <strong v-if="collectingInboxes" class="inbox-inline-progress">
          {{ t("ktpScanIn.inbox.progress") }} {{ inboxScannedQty }}/{{ inboxExpectedQty }}
        </strong>
      </div>

      <div class="operation-grid">
        <div class="setup-grid">
        <label class="field-block production-line-field">
          <span>{{ t("ktpScanIn.productionLine") }} <b>*</b></span>
          <el-select
            ref="productionLineInputRef"
            v-model="scanSetup.productionLine"
            filterable
            :loading="loadingProductionLines"
            :placeholder="t('ktpScanIn.productionLinePlaceholder')"
            @change="focusStorage"
          >
            <el-option
              v-for="line in productionLines"
              :key="line.id"
              :label="line.name"
              :value="line.id"
            >
              <span>{{ line.name }}</span>
              <small>{{ line.factory }} · {{ line.id }}</small>
            </el-option>
          </el-select>
        </label>

        <label class="field-block short-field">
          <span>{{ t("ktpScanIn.storage") }} <b>*</b></span>
          <el-input
            ref="storageInputRef"
            v-model="scanSetup.storage"
            autocomplete="off"
            placeholder="A2"
            @keyup.enter="focusLocation"
          />
        </label>

        <label class="field-block short-field">
          <span>{{ t("ktpScanIn.locationLabel") }}</span>
          <el-input
            ref="locationInputRef"
            v-model="scanSetup.location"
            autocomplete="off"
            placeholder="0"
            @keyup.enter="focusScanner"
          />
        </label>

        <div class="carton-count" title="Số carton trong phiên">
          <span>Carton</span>
          <strong>{{ sessionScans.length }}</strong>
        </div>
        </div>

        <label class="field-block barcode-field">
          <span>{{ scannerTitle }}</span>
          <el-input
            ref="cartonInputRef"
            v-model="cartonInput"
            class="barcode-input"
            inputmode="text"
            autocomplete="off"
            autocapitalize="characters"
            :disabled="!setupReady || loadingCarton"
            :placeholder="scannerPlaceholder"
            @input="queueAutomaticCartonLookup"
            @keyup.enter="handleCartonEnter"
          >
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>

        </label>

        <div class="scale-block" :class="`is-${scaleStatus}`">
          <div class="scale-reading">
            <span>{{ t("ktpScanIn.scale.currentWeight") }}</span>
            <strong>{{ scaleDisplayWeight }}</strong>
            <em>kg</em>
          </div>
          <div class="scale-detail">
            <span><i></i></span>
          </div>
          <div class="scale-actions">
            <el-button
              v-if="!scaleConnected"
              type="primary"
              size="small"
              :loading="scaleConnecting"
              :disabled="!isSerialSupported"
              @click="handleScaleConnect(hasSavedPort)"
            >
              {{ scaleConnectLabel }}
            </el-button>
          </div>
        </div>

      </div>

      <section v-if="cartonInfo" ref="cartonPreviewRef" class="carton-preview">
        <div><span>CARTONBAR</span><b>{{ cartonInfo.cartonBar || EMPTY_VALUE }}</b></div>
        <div><span>YSBH</span><b>{{ cartonInfo.orderNo || EMPTY_VALUE }}</b></div>
        <div><span>Article</span><b>{{ cartonInfo.article || EMPTY_VALUE }}</b></div>
        <div><span>XieMing</span><b>{{ cartonInfo.xieMing || EMPTY_VALUE }}</b></div>
        <div><span>Country</span><b>{{ cartonInfo.country || EMPTY_VALUE }}</b></div>
        <div><span>XH</span><b>{{ cartonInfo.xh || EMPTY_VALUE }}</b></div>
        <div><span>Qty</span><b>{{ formatNumber(cartonInfo.packageQty) }}</b></div>
      </section>

      <section v-if="isInboxMode && inboxExpectedQty > 0" class="inbox-progress-panel">
        <div class="inbox-progress-summary">
          <strong>{{ t("ktpScanIn.inbox.title") }}</strong>
          <span>{{ t("ktpScanIn.inbox.progress") }} {{ inboxScannedQty }}/{{ inboxExpectedQty }}</span>
          <b>{{ t("ktpScanIn.inbox.remaining") }}: {{ inboxRemainingQty }}</b>
        </div>
        <div class="inbox-progress-track"><i :style="{ width: `${inboxProgressPercent}%` }"></i></div>
        <div class="inbox-progress-items">
          <span
            v-for="item in inboxItems"
            :key="item.inboxBar"
            :class="{ 'is-complete': item.scannedQty >= item.expectedQty }"
          >
            {{ item.size || EMPTY_VALUE }}: {{ item.scannedQty }}/{{ item.expectedQty }}
          </span>
        </div>
      </section>

      <aside v-if="validationResult" class="condition-status">
        <el-alert
          v-for="check in failedChecks"
          :key="check.key"
          class="condition-card"
          type="error"
          show-icon
          :closable="false"
          :title="check.title"
          :description="check.message"
        />

        <section
          class="weight-validation"
          :class="`is-${currentWeightCheck.status}`"
        >
          <div class="weight-title">
            <strong>{{ t("ktpScanIn.weightCheckTitle") }}</strong>
            <span>{{ weightCheckStatusText }}</span>
          </div>
          <div class="weight-values">
            <div><span>{{ t("ktpScanIn.weight.standardCarton") }}</span><b>{{ formatKg(currentWeightCheck.standardWeight) }}</b></div>
            <div><span>{{ t("ktpScanIn.weight.actual") }}</span><b>{{ formatKg(currentWeightCheck.actualWeight) }}</b></div>
            <div><span>{{ t("ktpScanIn.weight.standardShoe") }}</span><b>{{ formatKg(currentWeightCheck.gram) }}</b></div>
            <div><span>{{ t("ktpScanIn.weight.allowedTolerance") }}</span><b>± {{ formatKg(currentWeightCheck.tolerance) }}</b></div>
            <div><span>{{ t("ktpScanIn.weight.difference") }}</span><b>{{ formatSignedKg(currentWeightCheck.difference) }}</b></div>
          </div>
          <p>{{ currentWeightCheck.message }}</p>
          <el-button
            v-if="canUnlockWeight"
            class="weight-unlock-button"
            type="warning"
            size="small"
            :icon="Unlock"
            @click="handleWeightOverride"
          >
            {{ t("ktpScanIn.weight.unlock") }}
          </el-button>
        </section>
      </aside>
    </section>

    <section ref="sessionPanelRef" class="session-panel">
      <div class="table-heading">
        <strong>{{ t("ktpScanIn.cartonList") }}</strong>
        <div class="table-heading-actions">
          <span class="session-carton-count">{{ sessionScans.length }} carton</span>
          <el-button
            class="end-session-button"
            size="small"
            type="warning"
            plain
            :icon="Finished"
            :disabled="loadingCarton"
            @click="handleEndSession"
          >
            {{ t("ktpScanIn.actions.endSession") }}
          </el-button>
        </div>
      </div>
      <div class="table-wrap">
        <el-table
          :data="sessionScans"
          height="100%"
          size="small"
          border
          table-layout="fixed"
          :empty-text="t('ktpScanIn.historyEmpty')"
        >
          <el-table-column type="index" label="#" width="44" align="center" />
          <el-table-column prop="cartonBar" label="CARTONBAR" width="205" show-overflow-tooltip />
          <el-table-column prop="orderNo" label="YSBH" width="125" show-overflow-tooltip />
          <el-table-column prop="article" label="Article" width="130" show-overflow-tooltip />
          <el-table-column prop="xieMing" label="XieMing" min-width="165" show-overflow-tooltip />
          <el-table-column prop="country" label="Country" width="130" show-overflow-tooltip />
          <el-table-column prop="xh" label="XH" width="62" align="center" />
          <el-table-column prop="quantity" label="Qty" width="68" align="right" />
          <el-table-column prop="depName" label="DepName" width="145" show-overflow-tooltip />
          <el-table-column prop="ckbh" label="CKBH" width="68" align="center" />
          <el-table-column prop="kcbh" label="KCBH" width="68" align="center" />
          <el-table-column prop="kvbh" label="KVBH" width="68" align="center" />
          <el-table-column prop="incs" label="INCS" width="62" align="center" />
          <el-table-column prop="inDate" label="INDATE" width="145" />
          <el-table-column prop="rgw" label="rgw" width="82" align="right" />
          <el-table-column prop="inUser" label="INUSER" width="90" show-overflow-tooltip />
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
import { Back, Finished, Search, SwitchButton, Unlock, UserFilled } from "@element-plus/icons-vue";
import { useWarehouseMapI18n } from "@/hooks/useWarehouseMapI18n";
import { SCALE_STATUS, useScaleSerial } from "@/hooks/useScaleSerial";
import { clearSession } from "@/utils/appStoreAndSecurity";

const router = useRouter();
const { t, lang, langOptions, setLang } = useWarehouseMapI18n();
const API_PREFIX = (import.meta.env.VITE_API_URL || "/api/v1").replace(/\/+$/, "");
const API_BASE = `${API_PREFIX}/warehouse-ktp`;
const EMPTY_VALUE = "—";
const SCAN_MODE = Object.freeze({
  CARTON: "carton",
  CARTON_INBOX: "carton_inbox",
});
const WORK_SESSION_STORAGE_VERSION = 1;

function getAuthHeaders(extraHeaders = {}) {
  const token = localStorage.getItem("token");
  return {
    ...extraHeaders,
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };
}
const {
  currentWeight,
  scaleStatus,
  scaleErrorCode,
  isSerialSupported,
  scaleConnected,
  scaleConnecting,
  hasSavedPort,
  initializeScale,
  connectScale,
  resetWeight,
  disposeScale,
} = useScaleSerial();

const scanSetup = reactive({
  warehouse: "A",
  productionLine: "",
  storage: "A2",
  location: "",
  toSx: "A",
});
const scanMode = ref(SCAN_MODE.CARTON);
const cartonInput = ref("");
const previewCartonCode = ref("");
const cartonInfo = ref(null);
const validationResult = ref(null);
const productionLines = ref([]);
const loadingProductionLines = ref(false);
const loadingCarton = ref(false);
const sessionScans = ref([]);
const inboxItems = ref([]);
const inboxScans = ref([]);
const collectingInboxes = ref(false);
const lockedCartonWeight = ref(null);
const productionLineInputRef = ref(null);
const storageInputRef = ref(null);
const locationInputRef = ref(null);
const cartonInputRef = ref(null);
const cartonPreviewRef = ref(null);
const sessionPanelRef = ref(null);
const AUTOMATIC_LOOKUP_DELAY_MS = 180;
const MINIMUM_CARTON_LENGTH = 8;
let automaticLookupTimer = null;
let scanAudioContext = null;
let suspendWorkSessionPersistence = false;
const weightOverrideConfirmed = ref(false);
const weightOverrideUserId = ref("");
const currentUserId = computed(() => loggedInUserId() || EMPTY_VALUE);

function workSessionStorageKey() {
  const user = localStorage.getItem("USERID")
    || localStorage.getItem("userid")
    || localStorage.getItem("userId")
    || "anonymous";
  return `warehouse-ktp-scan-in-session:${String(user).trim().toUpperCase()}`;
}

function persistWorkSession() {
  if (suspendWorkSessionPersistence) return;
  localStorage.setItem(workSessionStorageKey(), JSON.stringify({
    version: WORK_SESSION_STORAGE_VERSION,
    scanSetup: { ...scanSetup },
    scanMode: scanMode.value,
    sessionScans: sessionScans.value,
  }));
}

function restoreWorkSession() {
  const raw = localStorage.getItem(workSessionStorageKey());
  if (!raw) return;
  try {
    const saved = JSON.parse(raw);
    if (saved?.version !== WORK_SESSION_STORAGE_VERSION) return;
    if (saved.scanSetup && typeof saved.scanSetup === "object") {
      Object.assign(scanSetup, {
        warehouse: "A",
        productionLine: String(saved.scanSetup.productionLine || ""),
        storage: String(saved.scanSetup.storage || "A2"),
        location: String(saved.scanSetup.location || ""),
        toSx: "A",
      });
    }
    if (Object.values(SCAN_MODE).includes(saved.scanMode)) {
      scanMode.value = saved.scanMode;
    }
    if (Array.isArray(saved.sessionScans)) {
      sessionScans.value = saved.sessionScans.filter((item) => item?.cartonBar);
    }
  } catch {
    localStorage.removeItem(workSessionStorageKey());
  }
}

function ensureScanAudioContext() {
  const AudioContextClass = window.AudioContext || window.webkitAudioContext;
  if (!AudioContextClass) return null;
  if (!scanAudioContext || scanAudioContext.state === "closed") {
    scanAudioContext = new AudioContextClass();
  }
  if (scanAudioContext.state === "suspended") {
    void scanAudioContext.resume().catch(() => {});
  }
  return scanAudioContext;
}

function unlockScanAudioContext() {
  ensureScanAudioContext();
}

function scheduleScanTone(context, frequency, startsAfter, duration, volume = 0.11, type = "sine") {
  const oscillator = context.createOscillator();
  const gain = context.createGain();
  const startsAt = context.currentTime + startsAfter;
  const endsAt = startsAt + duration;

  oscillator.type = type;
  oscillator.frequency.setValueAtTime(frequency, startsAt);
  gain.gain.setValueAtTime(0.0001, startsAt);
  gain.gain.exponentialRampToValueAtTime(volume, startsAt + 0.012);
  gain.gain.exponentialRampToValueAtTime(0.0001, endsAt);
  oscillator.connect(gain);
  gain.connect(context.destination);
  oscillator.start(startsAt);
  oscillator.stop(endsAt + 0.01);
}

async function playScanSound(result) {
  const context = ensureScanAudioContext();
  if (!context) return;
  try {
    if (context.state === "suspended") await context.resume();
    if (result === "success") {
      scheduleScanTone(context, 880, 0, 0.11);
      scheduleScanTone(context, 1175, 0.13, 0.16);
      return;
    }
    if (result === "cartonReady") {
      scheduleScanTone(context, 620, 0, 0.1, 0.14, "triangle");
      scheduleScanTone(context, 820, 0.12, 0.14, 0.14, "triangle");
      return;
    }
    if (result === "inboxScanned") {
      scheduleScanTone(context, 1450, 0, 0.075, 0.16, "sine");
      return;
    }
    if (result === "inboxComplete") {
      scheduleScanTone(context, 760, 0, 0.1, 0.15, "triangle");
      scheduleScanTone(context, 1020, 0.11, 0.1, 0.15, "triangle");
      scheduleScanTone(context, 1320, 0.22, 0.18, 0.16, "triangle");
      return;
    }
    // Two short, deep and dense buzzes. Each buzz uses two layers so the low
    // tone remains loud enough on small monitor speakers.
    scheduleScanTone(context, 180, 0, 0.23, 0.26, "square");
    scheduleScanTone(context, 230, 0, 0.23, 0.2, "sawtooth");
    scheduleScanTone(context, 180, 0.34, 0.25, 0.28, "square");
    scheduleScanTone(context, 230, 0.34, 0.25, 0.22, "sawtooth");
  } catch {
    // Trình duyệt có thể chặn âm thanh cho tới tương tác đầu tiên của người dùng.
  }
}

const setupReady = computed(() => [
  scanSetup.productionLine,
  scanSetup.storage,
].every((value) => String(value || "").trim()));

const selectedProductionLine = computed(() =>
  productionLines.value.find((line) => line.id === scanSetup.productionLine),
);
const isInboxMode = computed(() => scanMode.value === SCAN_MODE.CARTON_INBOX);
const inboxExpectedQty = computed(() => inboxItems.value.reduce(
  (total, item) => total + Number(item.expectedQty || 0),
  0,
));
const inboxScannedQty = computed(() => inboxItems.value.reduce(
  (total, item) => total + Number(item.scannedQty || 0),
  0,
));
const inboxRemainingQty = computed(() => Math.max(0, inboxExpectedQty.value - inboxScannedQty.value));
const inboxProgressPercent = computed(() => (
  inboxExpectedQty.value > 0
    ? Math.min(100, Math.round((inboxScannedQty.value / inboxExpectedQty.value) * 100))
    : 0
));
const scannerTitle = computed(() => (
  collectingInboxes.value ? t("ktpScanIn.inbox.scanTitle") : t("ktpScanIn.cartonScanTitle")
));
const scannerPlaceholder = computed(() => (
  collectingInboxes.value ? t("ktpScanIn.inbox.placeholder") : t("ktpScanIn.cartonPlaceholder")
));
const currentLanguageLabel = computed(() =>
  langOptions.find((option) => option.code === lang.value)?.label || "VI",
);
const currentNumberLocale = computed(() => ({
  vi: "vi-VN",
  en: "en-US",
  zh: "zh-CN",
}[lang.value] || "vi-VN"));

const currentWeightCheck = computed(() => {
  const source = validationResult.value?.weightRule;
  if (!source) return { status: "no_data", passed: false, message: "" };
  const liveWeight = Number(currentWeight.value);
  const lockedWeight = Number(lockedCartonWeight.value);
  const actual = lockedWeight > 0 ? lockedWeight : (liveWeight > 0 ? liveWeight : source.actualWeight);
  if (!(Number(actual) > 0)) {
    return { ...source, status: "no_data", passed: false, difference: null, message: t("ktpScanIn.weight.noActual") };
  }
  if (!(source.standardWeight > 0) || !(source.tolerance > 0)) {
    return { ...source, actualWeight: actual };
  }
  const difference = source.standardWeight - actual;
  const exceeded = Math.abs(difference) > source.tolerance;
  const overridden = exceeded && weightOverrideConfirmed.value;
  return {
    ...source,
    actualWeight: actual,
    difference,
    status: overridden ? "overridden" : (exceeded ? "failed" : "passed"),
    passed: !exceeded || overridden,
    overridden,
    overrideUserId: overridden ? weightOverrideUserId.value : "",
    message: overridden
      ? t("ktpScanIn.weight.overrideConfirmed").replace("{userid}", weightOverrideUserId.value)
      : source.message,
    // message: t(exceeded ? "ktpScanIn.weight.failedMessage" : "ktpScanIn.weight.passedMessage"),
  };
});

const nonWeightChecksPassed = computed(() => {
  const result = validationResult.value;
  return Boolean(
    result
      && result.eligibility?.passed
      && result.productionCheck?.passed
      && result.materialCheck?.passed,
  );
});

const weightOverrideEligible = computed(() => Boolean(
  validationResult.value
    && currentWeightCheck.value.status === "failed"
    && nonWeightChecksPassed.value,
));
const canUnlockWeight = computed(() => weightOverrideEligible.value && !loadingCarton.value);

const failedChecks = computed(() => {
  if (!validationResult.value) return [];
  const result = validationResult.value;
  return [
    { key: "carton", title: t("ktpScanIn.validation.carton"), ...result.eligibility },
    { key: "production", title: t("ktpScanIn.production.title"), ...result.productionCheck },
    { key: "material", title: t("ktpScanIn.material.title"), ...result.materialCheck },
  ].filter((check) => check.status !== "passed");
});

const previewMainChecksPassed = computed(() => {
  const result = validationResult.value;
  return Boolean(
    result
      && currentWeightCheck.value.passed
      && result.eligibility?.passed
      && result.productionCheck?.passed
      && result.materialCheck?.passed,
  );
});
const previewCanScanIn = computed(() => Boolean(
  previewMainChecksPassed.value
    && (!isInboxMode.value || (inboxExpectedQty.value > 0 && inboxRemainingQty.value === 0)),
));

const weightCheckStatusText = computed(() => {
  if (currentWeightCheck.value.status === "overridden") {
    return t("ktpScanIn.weight.overrideStatus");
  }
  return t(`ktpScanIn.status.${currentWeightCheck.value.status || "no_data"}`);
});

const scaleDisplayWeight = computed(() => {
  if (currentWeight.value === null) return "--.--";
  return (Object.is(currentWeight.value, -0) ? 0 : currentWeight.value).toFixed(2);
});
const scaleStatusText = computed(() => t(`ktpScanIn.scale.status.${scaleStatus.value}`));
const scaleConnectLabel = computed(() =>
  t(hasSavedPort.value ? "ktpScanIn.scale.reconnect" : "ktpScanIn.scale.connect"),
);
const scaleErrorText = computed(() => {
  if (!isSerialSupported.value) return t("ktpScanIn.scale.errors.unsupported");
  return scaleErrorCode.value ? t(`ktpScanIn.scale.errors.${scaleErrorCode.value}`) : "";
});

function focusInput(inputRef) {
  nextTick(() => inputRef.value?.focus?.());
}
function focusProductionLine() { focusInput(productionLineInputRef); }
function focusStorage() { focusInput(storageInputRef); }
function focusLocation() { focusInput(locationInputRef); }
function focusScanner() {
  if (setupReady.value) focusInput(cartonInputRef);
}

function formatKg(value) {
  return Number.isFinite(Number(value)) && value !== null
    ? `${Number(value).toFixed(3)} kg`
    : EMPTY_VALUE;
}
function formatSignedKg(value) {
  if (!Number.isFinite(Number(value)) || value === null) return EMPTY_VALUE;
  const number = Number(value);
  return `${number > 0 ? "+" : ""}${number.toFixed(3)} kg`;
}
function formatNumber(value) {
  return Number(value || 0).toLocaleString(currentNumberLocale.value, { maximumFractionDigits: 2 });
}
function formatDate(value) {
  const date = value ? new Date(value) : new Date();
  return Number.isNaN(date.getTime())
    ? EMPTY_VALUE
    : date.toLocaleString(currentNumberLocale.value, { hour12: false });
}
function getErrorMessage(error, fallback) {
  return error.response?.data?.message || error.message || fallback;
}
async function handleScaleConnect(reuseSavedPort) {
  const connected = await connectScale({ reuseSavedPort });
  if (connected) {
    ElMessage.success(t("ktpScanIn.scale.connectedMessage"));
  } else if (scaleErrorCode.value !== "permissionCancelled") {
    ElMessage.error(scaleErrorText.value || t("ktpScanIn.scale.errors.connectionFailed"));
  }
}
async function loadSetupOptions() {
  loadingProductionLines.value = true;
  try {
    const response = await axios.get(`${API_BASE}/scan-in/production-lines`, { headers: getAuthHeaders() });
    productionLines.value = Array.isArray(response.data?.data)
      ? response.data.data
      : [];
  } catch (error) {
    ElMessage.error(getErrorMessage(error, t("ktpScanIn.productionLineLoadFailed")));
  } finally {
    loadingProductionLines.value = false;
  }
}

function clearPreview() {
  previewCartonCode.value = "";
  cartonInfo.value = null;
  validationResult.value = null;
  inboxItems.value = [];
  inboxScans.value = [];
  collectingInboxes.value = false;
  lockedCartonWeight.value = null;
  weightOverrideConfirmed.value = false;
  weightOverrideUserId.value = "";
}

function clearAutomaticLookupTimer() {
  if (automaticLookupTimer === null) return;
  window.clearTimeout(automaticLookupTimer);
  automaticLookupTimer = null;
}

function queueAutomaticCartonLookup(value) {
  ensureScanAudioContext();
  clearAutomaticLookupTimer();
  const cartonBar = String(value || "").trim();
  const minimumLength = collectingInboxes.value ? 1 : MINIMUM_CARTON_LENGTH;
  if (!setupReady.value || loadingCarton.value || cartonBar.length < minimumLength) return;

  automaticLookupTimer = window.setTimeout(() => {
    automaticLookupTimer = null;
    if (loadingCarton.value || String(cartonInput.value || "").trim() !== cartonBar) return;
    void handleCartonEnter();
  }, AUTOMATIC_LOOKUP_DELAY_MS);
}

async function handleCartonEnter() {
  ensureScanAudioContext();
  clearAutomaticLookupTimer();
  if (!setupReady.value) {
    void playScanSound("error");
    ElMessage.warning(t("ktpScanIn.setupIncomplete"));
    focusProductionLine();
    return;
  }
  const scannedValue = String(cartonInput.value || "").trim();
  if (!scannedValue) {
    focusScanner();
    return;
  }
  if (loadingCarton.value) return;
  if (collectingInboxes.value) {
    await handleInboxScan(scannedValue);
    return;
  }
  await loadCartonInformation(scannedValue);
}

function beginInboxCollection(inboxCheck) {
  const items = Array.isArray(inboxCheck?.items) ? inboxCheck.items : [];
  const expectedQty = Number(inboxCheck?.expectedQty || 0);
  if (!items.length || expectedQty <= 0) {
    const message = inboxCheck?.message || t("ktpScanIn.inbox.noData");
    clearPreview();
    resetWeight();
    void playScanSound("error");
    ElMessage.error(message);
    return false;
  }

  inboxItems.value = items.map((item) => ({
    inboxBar: String(item.inboxBar || "").trim(),
    size: String(item.size || "").trim(),
    expectedQty: Number(item.expectedQty || 0),
    scannedQty: 0,
  }));
  inboxScans.value = [];
  lockedCartonWeight.value = Number(currentWeightCheck.value.actualWeight || 0);
  collectingInboxes.value = true;
  void playScanSound("cartonReady");
  ElMessage.info(t("ktpScanIn.inbox.scanPrompt"));
  focusScanner();
  return true;
}

async function handleInboxScan(scannedValue) {
  const inboxBar = String(scannedValue || "").trim();
  cartonInput.value = "";
  if (!inboxBar) {
    focusScanner();
    return;
  }

  const item = inboxItems.value.find(
    (candidate) => candidate.inboxBar.toUpperCase() === inboxBar.toUpperCase(),
  );
  if (!item) {
    void playScanSound("error");
    ElMessage.error(`${t("ktpScanIn.inbox.invalid")}: ${inboxBar}`);
    focusScanner();
    return;
  }
  if (item.scannedQty >= item.expectedQty) {
    void playScanSound("error");
    ElMessage.error(`${t("ktpScanIn.inbox.exceeded")}: ${item.inboxBar}`);
    focusScanner();
    return;
  }

  item.scannedQty += 1;
  inboxScans.value.push(item.inboxBar);
  void playScanSound("inboxScanned");
  if (inboxRemainingQty.value === 0) {
    collectingInboxes.value = false;
    await confirmPreviewCarton();
    return;
  }
  focusScanner();
}

async function animateCartonPreviewToTable() {
  await nextTick();
  const preview = cartonPreviewRef.value;
  const target = sessionPanelRef.value;
  if (!preview || !target || typeof document === "undefined") return;

  const previewBox = preview.getBoundingClientRect();
  const targetBox = target.getBoundingClientRect();
  if (!previewBox.width || !previewBox.height || !targetBox.width) return;

  const flyer = document.createElement("div");
  const flyerIcon = document.createElement("span");
  const flyerText = document.createElement("div");
  const flyerTitle = document.createElement("strong");
  const flyerCode = document.createElement("b");
  flyer.className = "ktp-carton-flyer";
  flyerIcon.className = "ktp-carton-flyer__icon";
  flyerText.className = "ktp-carton-flyer__text";
  flyerIcon.textContent = "📦";
  flyerTitle.textContent = t("ktpScanIn.animationPassed");
  flyerCode.textContent = previewCartonCode.value;
  flyerText.append(flyerTitle, flyerCode);
  flyer.append(flyerIcon, flyerText);

  const flyerWidth = Math.min(310, Math.max(245, previewBox.width * 0.3));
  const startLeft = previewBox.left + Math.min(20, Math.max(0, previewBox.width - flyerWidth));
  const startTop = previewBox.top + Math.max(0, (previewBox.height - 64) / 2);
  const destinationLeft = targetBox.left + 18;
  const destinationTop = targetBox.top + 38;
  const moveX = destinationLeft - startLeft;
  const moveY = destinationTop - startTop;
  flyer.style.left = `${startLeft}px`;
  flyer.style.top = `${startTop}px`;
  flyer.style.width = `${flyerWidth}px`;
  document.body.appendChild(flyer);

  if (typeof flyer.animate !== "function") {
    flyer.remove();
    return;
  }

  const animation = flyer.animate([
    { transform: "translate3d(0, 0, 0) scale(0.72)", opacity: 0 },
    { offset: 0.12, transform: "translate3d(0, -10px, 0) scale(1)", opacity: 1 },
    {
      offset: 0.55,
      transform: `translate3d(${moveX * 0.5 + 65}px, ${moveY * 0.42 - 38}px, 0) rotate(4deg) scale(1.04)`,
      opacity: 1,
    },
    {
      offset: 0.88,
      transform: `translate3d(${moveX * 0.94}px, ${moveY * 0.9}px, 0) rotate(-2deg) scale(0.94)`,
      opacity: 1,
    },
    { transform: `translate3d(${moveX}px, ${moveY}px, 0) rotate(0) scale(0.82)`, opacity: 0.12 },
  ], {
    duration: 900,
    easing: "cubic-bezier(0.2, 0.72, 0.28, 1)",
    fill: "forwards",
  });

  try {
    await animation.finished;
  } catch {
    // The element can be removed when the page changes while the animation is running.
  } finally {
    flyer.remove();
  }
}

function loggedInUserId() {
  return String(
    localStorage.getItem("USERID")
      || localStorage.getItem("userid")
      || localStorage.getItem("userId")
      || "",
  ).trim();
}

async function requestWeightOverride({ askToContinue = false, allowWhileLoading = false } = {}) {
  if (!weightOverrideEligible.value || (!allowWhileLoading && loadingCarton.value)) return false;

  if (askToContinue) {
    void playScanSound("error");
    try {
      await ElMessageBox.confirm(
        t("ktpScanIn.weight.continueQuestion"),
        t("ktpScanIn.weight.continueTitle"),
        {
          confirmButtonText: t("ktpScanIn.weight.continueYes"),
          cancelButtonText: t("ktpScanIn.weight.continueNo"),
          type: "warning",
          closeOnClickModal: false,
        },
      );
    } catch {
      resetWeight();
      focusScanner();
      return false;
    }
  }

  const currentUser = loggedInUserId();
  if (!currentUser) {
    ElMessage.error(t("ktpScanIn.weight.unlockMissingUser"));
    return false;
  }

  try {
    await ElMessageBox.prompt(
      t("ktpScanIn.weight.unlockPrompt"),
      t("ktpScanIn.weight.unlockTitle"),
      {
        confirmButtonText: t("ktpScanIn.weight.unlockConfirm"),
        cancelButtonText: t("ktpScanIn.actions.cancel"),
        inputType: "password",
        inputPlaceholder: t("ktpScanIn.weight.unlockPlaceholder"),
        inputValidator: (value) => (
          String(value || "").trim().toUpperCase() === currentUser.toUpperCase()
            || t("ktpScanIn.weight.unlockUserMismatch")
        ),
      },
    );
  } catch {
    resetWeight();
    focusScanner();
    return false;
  }

  weightOverrideUserId.value = currentUser;
  weightOverrideConfirmed.value = true;
  ElMessage.success(t("ktpScanIn.weight.unlockSuccess"));
  await nextTick();

  if (isInboxMode.value) {
    beginInboxCollection(validationResult.value?.inboxCheck);
  } else {
    await confirmPreviewCarton();
  }
  return true;
}

async function handleWeightOverride() {
  await requestWeightOverride();
}

async function confirmPreviewCarton() {
  if (!previewCanScanIn.value) {
    void playScanSound("error");
    ElMessage.error("Thùng chưa đạt đủ điều kiện, không thể ghi xuống");
    focusScanner();
    return;
  }

  const inUser = loggedInUserId() || "WEB_KTP";
  loadingCarton.value = true;
  try {
    const response = await axios.post(`${API_BASE}/scan-in/carton`, {
      cartonBar: previewCartonCode.value,
      actualWeight: currentWeightCheck.value.actualWeight,
      warehouse: scanSetup.warehouse,
      productionLine: scanSetup.productionLine,
      storage: scanSetup.storage,
      location: scanSetup.location,
      toSx: scanSetup.toSx,
      inUser,
      scanMode: scanMode.value,
      inboxScans: inboxScans.value,
      weightOverride: weightOverrideConfirmed.value,
      weightOverrideUserId: weightOverrideUserId.value,
    }, {
      headers: getAuthHeaders({ "X-User": inUser }),
    });
    const savedCarton = response.data?.data || validationResult.value;
    void playScanSound(isInboxMode.value ? "inboxComplete" : "success");
    await animateCartonPreviewToTable();
    addSessionScan(savedCarton);
    prepareScannerAfterSuccessfulCarton();
    ElMessage.success("Đã nhập kho thùng carton vào YWCP");
  } catch (error) {
    void playScanSound("error");
    ElMessage.error(getErrorMessage(error, "Không thể cập nhật thùng carton vào YWCP"));
    clearPreview();
    resetWeight();
  } finally {
    loadingCarton.value = false;
    focusScanner();
  }
}

async function loadCartonInformation(cartonBar) {
  if (loadingCarton.value) return;
  loadingCarton.value = true;
  clearPreview();
  try {
    const response = await axios.get(`${API_BASE}/scan-in/carton/${encodeURIComponent(cartonBar)}`, {
      headers: getAuthHeaders(),
      params: {
        actualWeight: currentWeight.value > 0 ? currentWeight.value : undefined,
        scanMode: scanMode.value,
      },
    });
    const data = response.data?.data || {};
    const orderNo = String(data.carton?.ddbh || data.carton?.orderNo || "").trim();
    const hasDifferentOrder = orderNo && sessionScans.value.some((item) => {
      const sessionOrder = String(item.ddbh || item.orderNo || "").trim();
      return sessionOrder && sessionOrder.toUpperCase() !== orderNo.toUpperCase();
    });
    if (hasDifferentOrder) {
      ElMessage.warning({
        message: t("ktpScanIn.differentOrderWarning").replace("{order}", orderNo),
        duration: 4500,
        showClose: true,
      });
    }
    cartonInfo.value = data.carton || null;
    validationResult.value = data;
    previewCartonCode.value = cartonBar;
    cartonInput.value = "";
    await nextTick();
    if (previewMainChecksPassed.value) {
      if (isInboxMode.value) {
        beginInboxCollection(data.inboxCheck);
      } else {
        await confirmPreviewCarton();
      }
    } else if (weightOverrideEligible.value) {
      lockedCartonWeight.value = Number(currentWeightCheck.value.actualWeight || 0);
      await requestWeightOverride({ askToContinue: true, allowWhileLoading: true });
    } else {
      void playScanSound("error");
      ElMessage.error({
        message: "Thùng không đạt điều kiện nhập kho. Vui lòng kiểm tra các cảnh báo bên dưới và quét lại.",
        duration: 5000,
        showClose: true,
      });
      resetWeight();
    }
  } catch (error) {
    clearPreview();
    resetWeight();
    void playScanSound("error");
    ElMessage.error(getErrorMessage(error, t("ktpScanIn.lookupFailed")));
  } finally {
    cartonInput.value = "";
    loadingCarton.value = false;
    focusScanner();
  }
}

function addSessionScan(result) {
  const carton = result.carton;
  const user = localStorage.getItem("USERID")
    || localStorage.getItem("userid")
    || localStorage.getItem("userId")
    || carton.inUser
    || EMPTY_VALUE;
  sessionScans.value.unshift({
    cartonBar: carton.cartonBar,
    ddbh: carton.ddbh || "",
    orderNo: carton.orderNo || EMPTY_VALUE,
    article: carton.article || EMPTY_VALUE,
    xieMing: carton.xieMing || EMPTY_VALUE,
    country: carton.country || EMPTY_VALUE,
    xh: carton.xh || EMPTY_VALUE,
    quantity: formatNumber(carton.packageQty),
    depName: carton.depName || selectedProductionLine.value?.name || EMPTY_VALUE,
    ckbh: scanSetup.warehouse.toUpperCase(),
    kcbh: scanSetup.storage.toUpperCase(),
    kvbh: scanSetup.location.toUpperCase(),
    incs: Number(carton.incs || 0),
    inDate: formatDate(carton.inDate),
    rgw: Number(result.weightRule?.actualWeight || 0).toFixed(2),
    inUser: user,
  });
}

function prepareScannerAfterSuccessfulCarton() {
  clearAutomaticLookupTimer();
  cartonInput.value = "";
  collectingInboxes.value = false;
  resetWeight();
  focusScanner();
}

function handleLanguageChange(code) {
  setLang(code);
  focusScanner();
}

function handleScanModeChange() {
  clearAutomaticLookupTimer();
  cartonInput.value = "";
  clearPreview();
  resetWeight();
  focusScanner();
}

function resetWorkSession() {
  clearAutomaticLookupTimer();
  sessionScans.value = [];
  cartonInput.value = "";
  scanSetup.productionLine = "";
  scanSetup.location = "";
  clearPreview();
  resetWeight();
  focusProductionLine();
}

async function handleEndSession() {
  if (sessionScans.value.length) {
    try {
      await ElMessageBox.confirm(
        t("ktpScanIn.actions.endSessionConfirm"),
        t("ktpScanIn.actions.endSessionTitle"),
        {
          confirmButtonText: t("ktpScanIn.actions.endSession"),
          cancelButtonText: t("ktpScanIn.actions.cancel"),
          type: "warning",
        },
      );
    } catch {
      focusScanner();
      return;
    }
  }

  suspendWorkSessionPersistence = true;
  try {
    resetWorkSession();
    await nextTick();
    localStorage.removeItem(workSessionStorageKey());
  } finally {
    suspendWorkSessionPersistence = false;
  }
  ElMessage.success(t("ktpScanIn.actions.endSessionSuccess"));
}

function handleBackToWarehouse() {
  void router.push({ name: "warehouse-ktp" });
}

function handleLogout() {
  clearAutomaticLookupTimer();
  clearSession();
  void router.replace({
    name: "warehouse-ktp-scan-in-login",
    query: { redirect: "/warehouse-ktp/scan-in" },
  });
}

restoreWorkSession();
watch(
  [scanSetup, scanMode, sessionScans],
  persistWorkSession,
  { deep: true, flush: "sync" },
);
watch(setupReady, (ready) => {
  if (ready) focusScanner();
});
watch(scaleStatus, (status, previousStatus) => {
  if (status === SCALE_STATUS.LOST && previousStatus === SCALE_STATUS.CONNECTED) {
    ElMessage.warning(t("ktpScanIn.scale.lostMessage"));
  }
  if (status === SCALE_STATUS.ERROR && scaleErrorCode.value === "readFailed") {
    ElMessage.error(t("ktpScanIn.scale.errors.readFailed"));
  }
});

onMounted(async () => {
  window.addEventListener("pointerdown", unlockScanAudioContext, { passive: true });
  window.addEventListener("keydown", unlockScanAudioContext);
  await Promise.all([initializeScale({ autoConnect: true }), loadSetupOptions()]);
  if (setupReady.value) focusScanner();
  else focusProductionLine();
});
onBeforeUnmount(() => {
  clearAutomaticLookupTimer();
  window.removeEventListener("pointerdown", unlockScanAudioContext);
  window.removeEventListener("keydown", unlockScanAudioContext);
  if (scanAudioContext && scanAudioContext.state !== "closed") {
    void scanAudioContext.close();
  }
  scanAudioContext = null;
  void disposeScale();
});
</script>

<style scoped lang="scss">
.scan-in-page {
  display: flex;
  height: 100vh;
  min-height: 0;
  overflow: hidden;
  padding: 4px 6px 6px;
  flex-direction: column;
  gap: 5px;
  background: #edf1f4;
  color: #253142;
  font-size: 13px;
}

.scan-toolbar,
.session-panel {
  border: 1px solid #bdcbd3;
  background: #fff;
}

.scan-toolbar {
  position: relative;
  padding: 7px 9px 8px;
  background: #e0f1f1;
}

.scan-mode-row {
  display: flex;
  min-height: 30px;
  margin-bottom: 6px;
  padding-right: 300px;
  align-items: center;
  gap: 9px;
}

.scan-mode-row > span {
  color: #344054;
  font-size: 13px;
  font-weight: 800;
}

.scan-mode-row :deep(.el-radio-button__inner) {
  padding: 7px 14px;
  font-size: 13px;
  font-weight: 750;
}

.inbox-inline-progress {
  color: #067647;
  font-size: 13px;
}

.operation-grid {
  display: grid;
  grid-template-columns: max-content 380px max-content;
  align-items: end;
  justify-content: start;
  gap: 14px;
}

.page-actions {
  position: absolute;
  z-index: 3;
  top: 6px;
  right: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.current-user {
  display: inline-flex;
  height: 30px;
  max-width: 190px;
  padding: 0 9px;
  border: 1px solid #9fb5bf;
  border-radius: 4px;
  align-items: center;
  gap: 5px;
  overflow: hidden;
  background: #fff;
  color: #344054;
  font-size: 12px;
  white-space: nowrap;
}

.current-user span {
  color: #667085;
}

.current-user strong {
  overflow: hidden;
  color: #17324d;
  text-overflow: ellipsis;
}

.page-actions :deep(.action-icon-button) {
  width: 30px;
  height: 30px;
  min-height: 30px;
  padding: 0;
  margin-left: 0;
  font-size: 16px;
}

.page-actions :deep(.language-button) {
  height: 30px;
  min-height: 30px;
  margin-left: 0;
  padding: 0 9px;
  font-weight: 750;
}

.condition-status {
  display: grid;
  width: 100%;
  min-width: 0;
  margin-top: 6px;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  align-items: stretch;
  gap: 6px;
}

.setup-grid {
  display: grid;
  grid-template-columns: 280px 95px 110px 86px;
  align-items: end;
  gap: 9px;
}

.field-block {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 3px;
}

.field-block > span {
  color: #344054;
  font-size: 13px;
  font-weight: 800;
  line-height: 16px;
}

.field-block b {
  color: #d92d20;
}

.field-block :deep(.el-select) {
  width: 100%;
}

.setup-grid :deep(.el-input__wrapper),
.setup-grid :deep(.el-select__wrapper) {
  min-height: 38px;
  font-size: 15px;
}

.carton-count {
  display: flex;
  height: 38px;
  align-items: baseline;
  justify-content: center;
  gap: 6px;
  border: 1px solid #9fb8bc;
  background: #f7ffff;
}

.carton-count span {
  color: #52646a;
  font-size: 13px;
  font-weight: 700;
}

.carton-count strong {
  color: #d92d20;
  font-size: 22px;
  line-height: 28px;
}

.barcode-field {
  position: relative;
  width: 100%;
  align-self: end;
}

.barcode-input :deep(.el-input__wrapper) {
  min-height: 38px;
  box-shadow: 0 0 0 2px #409eff inset;
}

.barcode-input :deep(.el-input__inner) {
  font-family: Consolas, monospace;
  font-size: 17px;
  font-weight: 800;
}


.scale-block {
  display: grid;
  width: max-content;
  max-width: 100%;
  height: 38px;
  justify-self: start;
  min-width: 0;
  padding: 3px 5px 3px 9px;
  grid-template-columns: auto auto auto;
  align-items: center;
  gap: 9px;
  border: 1px solid #aab9c4;
  background: #fff;
}

.scale-reading {
  display: flex;
  align-items: baseline;
  gap: 4px;
  white-space: nowrap;
}

.scale-reading span {
  color: #667085;
  font-size: 12px;
  font-weight: 700;
}

.scale-reading strong {
  color: #d92d20;
  font-family: Consolas, monospace;
  font-size: 25px;
  line-height: 28px;
}

.scale-reading em {
  color: #475467;
  font-size: 13px;
  font-style: normal;
  font-weight: 800;
}

.scale-detail {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 2px;
}

.scale-detail > span {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #667085;
  font-size: 12px;
  font-weight: 700;
}

.scale-detail i {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #98a2b3;
}

.scale-block.is-connected .scale-detail i { background: #12b76a; }
.scale-block.is-error .scale-detail i,
.scale-block.is-lost .scale-detail i { background: #f04438; }

.scale-detail small {
  overflow: hidden;
  color: #475467;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.scale-actions {
  white-space: nowrap;
}

.carton-preview {
  position: relative;
  z-index: 5;
  display: grid;
  margin-top: 7px;
  padding: 7px 9px;
  grid-template-columns: 1.5fr 1fr 1fr 1.3fr 1fr 60px 60px;
  gap: 8px;
  border: 1px solid #bdcbd3;
  background: #fff;
  transform-origin: top left;
}

.carton-preview > div,
.weight-values > div {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 2px;
}

.carton-preview span,
.weight-values span {
  color: #667085;
  font-size: 11px;
  font-weight: 800;
}

.carton-preview b {
  overflow: hidden;
  color: #1d2939;
  font-size: 14px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.inbox-progress-panel {
  display: grid;
  margin-top: 5px;
  padding: 6px 9px;
  grid-template-columns: max-content minmax(180px, 1fr) auto;
  align-items: center;
  gap: 7px 12px;
  border: 1px solid #e3ae21;
  background: #fff9d7;
}

.inbox-progress-summary {
  display: flex;
  align-items: center;
  gap: 12px;
  white-space: nowrap;
}

.inbox-progress-summary strong {
  color: #7a4d00;
  font-size: 13px;
}

.inbox-progress-summary span,
.inbox-progress-summary b {
  color: #344054;
  font-size: 12px;
}

.inbox-progress-track {
  height: 8px;
  overflow: hidden;
  border-radius: 8px;
  background: #eadfb0;
}

.inbox-progress-track i {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: #12a150;
  transition: width 160ms ease;
}

.inbox-progress-items {
  display: flex;
  max-width: 420px;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 4px;
}

.inbox-progress-items span {
  padding: 2px 6px;
  border: 1px solid #d0b95c;
  border-radius: 3px;
  background: #fffef5;
  color: #684d00;
  font-family: Consolas, monospace;
  font-size: 11px;
  font-weight: 800;
}

.inbox-progress-items span.is-complete {
  border-color: #75c793;
  background: #eaf9ef;
  color: #067647;
}

:global(.ktp-carton-flyer) {
  position: fixed;
  z-index: 99999;
  display: flex;
  min-height: 64px;
  box-sizing: border-box;
  align-items: center;
  padding: 8px 13px;
  gap: 11px;
  border: 2px solid rgba(221, 161, 0, 0.78);
  border-radius: 7px;
  background: linear-gradient(145deg, rgba(255, 252, 178, 0.94) 0%, rgba(255, 221, 70, 0.9) 100%);
  box-shadow: 0 12px 25px rgba(135, 91, 0, 0.26);
  color: #553700;
  backdrop-filter: blur(2px);
  pointer-events: none;
  transform-origin: center;
}

:global(.ktp-carton-flyer__icon) {
  font-size: 36px;
  line-height: 1;
}

:global(.ktp-carton-flyer__text) {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 3px;
}

:global(.ktp-carton-flyer__text strong) {
  font-size: 13px;
  font-weight: 900;
}

:global(.ktp-carton-flyer__text b) {
  overflow: hidden;
  font-family: Consolas, monospace;
  font-size: 16px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.weight-validation {
  display: grid;
  min-width: 0;
  padding: 6px 8px;
  grid-template-columns: 1fr;
  align-content: start;
  gap: 5px;
  border: 1px solid #e3a008;
  background: #fffaeb;
}

.weight-validation.is-passed {
  border-color: #75c69a;
  background: #ecfdf3;
}

.weight-validation.is-overridden {
  border-color: #f0a331;
  background: #fff7e6;
}

.weight-validation.is-failed {
  border-color: #f19a94;
  background: #fff1f0;
}

.weight-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
}

.weight-title strong {
  font-size: 12px;
}

.weight-title span {
  color: #b54708;
  font-size: 12px;
  font-weight: 900;
}

.weight-validation.is-passed .weight-title span { color: #067647; }
.weight-validation.is-failed .weight-title span { color: #b42318; }
.weight-validation.is-overridden .weight-title span { color: #b54708; }

.weight-values {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(74px, 1fr));
  gap: 4px 6px;
}

.weight-values b {
  color: #1d2939;
  font-family: Consolas, monospace;
  font-size: 13px;
}

.weight-validation p {
  margin: 0;
  overflow: hidden;
  color: #475467;
  font-size: 11px;
  font-weight: 700;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.weight-unlock-button {
  width: max-content;
  margin: 1px 0 0;
  font-weight: 800;
}

.condition-card {
  min-width: 0;
  height: 100%;
}

.condition-status :deep(.el-alert) {
  min-width: 0;
  padding: 7px 10px;
}

.condition-status :deep(.el-alert__content) {
  min-width: 0;
}

.condition-status :deep(.el-alert__title) {
  font-size: 14px;
  font-weight: 800;
}

.condition-status :deep(.el-alert__description) {
  margin-top: 1px;
  font-size: 13px;
}

.session-panel {
  display: flex;
  min-height: 0;
  flex: 1;
  flex-direction: column;
  background: #e2fac5;
}

.table-heading {
  display: flex;
  height: 34px;
  padding: 0 9px;
  flex: 0 0 auto;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #bdcbd3;
  background: #f7f9fa;
}

.table-heading strong {
  color: #344054;
  font-size: 14px;
}

.table-heading .session-carton-count {
  color: #16803b;
  font-size: 13px;
  font-weight: 800;
  white-space: nowrap;
}

.table-heading-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.table-heading :deep(.end-session-button) {
  height: 26px;
  margin-left: 0;
  padding: 4px 9px;
  font-weight: 750;
}

.table-wrap {
  min-height: 0;
  flex: 1;
}

.table-wrap :deep(.el-table) {
  --el-table-bg-color: #e2fac5;
  --el-table-tr-bg-color: #e2fac5;
  --el-table-row-hover-bg-color: #d2f5aa;
  --el-table-header-bg-color: #eef2f2;
  --el-table-border-color: #aebcc2;
  color: #263238;
  font-size: 13px;
}

.table-wrap :deep(.el-table th.el-table__cell) {
  padding: 5px 0;
  color: #344054;
  font-weight: 800;
}

.table-wrap :deep(.el-table td.el-table__cell) {
  padding: 4px 0;
}

.table-wrap :deep(.el-table__empty-block) {
  background: #e2fac5;
}

@media (max-width: 1650px) {
  .operation-grid {
    grid-template-columns: max-content minmax(300px, 1fr) max-content;
  }

  .setup-grid {
    grid-template-columns: 250px 90px 100px 80px;
  }
}

@media (max-width: 1400px) {
  .scan-in-page { overflow: auto; }
  .operation-grid { grid-template-columns: minmax(420px, 1fr) 380px; }
  .setup-grid { grid-column: 1 / -1; }
  .session-panel { min-height: 420px; }
}

@media (max-width: 820px) {
  .operation-grid { grid-template-columns: 1fr; }
  .setup-grid { grid-template-columns: repeat(3, minmax(90px, 1fr)); }
  .production-line-field { grid-column: span 2; }
  .condition-status { grid-template-columns: 1fr; }
  .scale-block { max-width: none; }
  .carton-preview { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .weight-validation { grid-template-columns: 1fr; }
  .weight-values { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .scan-mode-row { align-items: flex-start; flex-direction: column; }
  .inbox-progress-panel { grid-template-columns: 1fr; }
  .inbox-progress-items { max-width: none; justify-content: flex-start; }
}

/* Đồng bộ bố cục với Scan Out và Scan Ship. */
.scan-in-page {
  box-sizing: border-box;
  height: 100vh;
  min-height: 0;
  overflow: hidden;
  padding: 12px;
  gap: 8px;
  background: #e8efec;
  color: #172033;
  font-size: 13px;
}

.topbar,
.page-actions {
  display: flex;
  align-items: center;
}

.topbar {
  min-height: 42px;
  flex: 0 0 auto;
  justify-content: space-between;
}

.topbar h1 {
  margin: 0;
  color: #173f32;
  font-size: 23px;
  text-shadow: 0 1px 0 #fff;
}

.page-actions {
  position: static;
  gap: 8px;
}

.current-user {
  height: 34px;
  max-width: 230px;
  padding: 0 11px;
  border-color: #d5dfda;
  border-radius: 999px;
  background: #fff;
  font-weight: 700;
}

.current-user span {
  display: none;
}

.page-actions :deep(.language-button) {
  height: 34px;
  min-height: 34px;
  border-radius: 8px;
}

.page-actions :deep(.action-icon-button) {
  width: 34px;
  height: 34px;
  min-height: 34px;
}

.scan-toolbar {
  display: flex;
  padding: 0;
  border: 0;
  flex: 0 0 auto;
  flex-direction: column;
  gap: 8px;
  background: transparent;
}

.scan-mode-row,
.operation-grid,
.carton-preview,
.inbox-progress-panel,
.condition-status {
  margin: 0;
  border: 1px solid #bdcec7;
  border-radius: 9px;
  background: #fff;
}

.scan-mode-row {
  min-height: 44px;
  padding: 7px 9px;
  background: linear-gradient(90deg, #f7fbf9, #fff);
  box-shadow: 0 4px 12px rgba(31, 60, 48, 0.06);
}

.scan-mode-row :deep(.el-radio-button__inner) {
  padding: 9px 15px;
}

.operation-grid {
  padding: 8px;
  grid-template-columns: max-content minmax(320px, 500px) max-content;
  align-items: end;
  justify-content: start;
  gap: 10px;
  box-shadow: 0 6px 18px rgba(31, 60, 48, 0.09);
}

.setup-grid {
  grid-column: auto;
  grid-template-columns: 280px 95px 110px 92px;
  gap: 8px;
}

.field-block > span {
  color: #34443d;
  font-size: 12px;
}

.setup-grid :deep(.el-input__wrapper),
.setup-grid :deep(.el-select__wrapper),
.barcode-input :deep(.el-input__wrapper) {
  min-height: 40px;
  border-radius: 7px;
  font-size: 15px;
}

.barcode-input :deep(.el-input__wrapper) {
  border: 1px solid #8faaa0;
  background: #fbfffd;
  box-shadow: 0 0 0 2px rgba(23, 107, 80, 0.08);
}

.barcode-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px #16845f inset;
}

.carton-count {
  height: 40px;
  border: 0;
  border-radius: 8px;
  align-items: center;
  background: #176b50;
}

.carton-count span,
.carton-count strong {
  color: #fff;
}

.carton-count strong {
  font-size: 25px;
}

.scale-block {
  height: 40px;
  padding: 3px 9px;
  border: 2px solid #39d98a;
  border-radius: 8px;
  background: #0b4935;
  box-shadow: 0 0 0 3px rgba(57, 217, 138, 0.14), 0 5px 14px rgba(11, 73, 53, 0.18);
}

.scale-reading span,
.scale-reading em {
  color: rgba(255, 255, 255, 0.78);
}

.scale-reading strong {
  color: #eaff62;
  font-size: 26px;
  text-shadow: 0 0 10px rgba(234, 255, 98, 0.3);
}

.scale-block.is-error,
.scale-block.is-lost,
.scale-block.is-disconnected {
  border-color: #f97066;
  background: #6b211d;
  box-shadow: 0 0 0 3px rgba(249, 112, 102, 0.14);
}

.scale-block.is-error .scale-reading strong,
.scale-block.is-lost .scale-reading strong,
.scale-block.is-disconnected .scale-reading strong {
  color: #fff2a8;
}

.scale-block.is-connecting {
  border-color: #fdb022;
  background: #654515;
  box-shadow: 0 0 0 3px rgba(253, 176, 34, 0.14);
}

.carton-preview {
  padding: 8px 10px;
  border-left: 4px solid #16845f;
  box-shadow: 0 4px 12px rgba(31, 60, 48, 0.06);
}

.inbox-progress-panel {
  border-color: #e3ae21;
  background: #fff9d7;
}

.condition-status {
  padding: 7px;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  align-items: stretch;
  gap: 7px;
  box-shadow: 0 4px 12px rgba(31, 60, 48, 0.06);
}

.condition-status :deep(.el-alert),
.weight-validation {
  min-height: 92px;
  box-sizing: border-box;
  border-radius: 7px;
}

.session-panel {
  border-color: #bdcec7;
  border-radius: 9px;
  background: #fff;
  box-shadow: 0 5px 16px rgba(31, 60, 48, 0.08);
  overflow: hidden;
}

.table-heading {
  height: 36px;
  border-bottom-color: #dce5e1;
  background: #edf5f1;
}

.table-wrap :deep(.el-table) {
  --el-table-bg-color: #fff;
  --el-table-tr-bg-color: #fff;
  --el-table-row-hover-bg-color: #edf8f3;
  --el-table-header-bg-color: #24483b;
  --el-table-border-color: #d7e0dc;
}

.table-wrap :deep(.el-table th.el-table__cell) {
  padding: 7px 0;
  background: #24483b;
  color: #fff;
}

.table-wrap :deep(.el-table__empty-block) {
  background: #fff;
}

@media (max-width: 1366px) {
  .scan-in-page { padding: 8px; gap: 6px; }
  .topbar { min-height: 37px; }
  .topbar h1 { font-size: 21px; }
  .scan-mode-row { min-height: 39px; padding: 5px 7px; }
  .scan-mode-row :deep(.el-radio-button__inner) { padding: 7px 12px; }
  .operation-grid { padding: 7px; grid-template-columns: max-content minmax(285px, 430px) max-content; gap: 8px; }
  .setup-grid { grid-template-columns: 230px 80px 90px 80px; gap: 6px; }
  .setup-grid :deep(.el-input__wrapper),
  .setup-grid :deep(.el-select__wrapper),
  .barcode-input :deep(.el-input__wrapper) { min-height: 37px; }
  .carton-count, .scale-block { height: 37px; }
  .scale-reading strong { font-size: 23px; }
  .carton-preview { padding-top: 6px; padding-bottom: 6px; }
  .condition-status { padding: 5px; gap: 5px; }
  .condition-status :deep(.el-alert), .weight-validation { min-height: 80px; }
}

@media (max-width: 1100px) {
  .scan-in-page {
    height: auto;
    min-height: 100vh;
    overflow: auto;
  }

  .operation-grid {
    grid-template-columns: minmax(360px, 500px) max-content;
  }

  .setup-grid {
    grid-column: 1 / -1;
  }

  .session-panel {
    min-height: 430px;
  }
}

@media (max-width: 820px) {
  .scan-in-page { padding: 7px; }
  .topbar { align-items: flex-start; }
  .current-user { display: none; }
  .page-actions { flex-wrap: wrap; justify-content: flex-end; }
  .scan-mode-row { padding-right: 9px; align-items: flex-start; flex-direction: column; }
  .operation-grid { grid-template-columns: 1fr; }
  .setup-grid { grid-template-columns: repeat(2, minmax(100px, 1fr)); }
  .production-line-field { grid-column: span 2; }
  .carton-preview { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .condition-status { grid-template-columns: 1fr; }
  .scale-block { width: auto; }
  .inbox-progress-panel { grid-template-columns: 1fr; }
  .inbox-progress-items { max-width: none; justify-content: flex-start; }
}
</style>
