<template>
  <section class="move-order-page">
    <div class="move-shell">
      <div class="top-bar">
        <el-button class="back-btn" @click="goBackToMap">
          <el-icon><ArrowLeft /></el-icon>
          {{ t("ktpMove.backToMap") }}
        </el-button>

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
      </div>

      <div class="scan-box">
        <div class="step-line">
          <div class="step-item is-active">
            <strong>1</strong>
            <span>{{ t("ktpMove.stepCarton") }}</span>
          </div>
          <div class="step-connector"></div>
          <div class="step-item" :class="{ 'is-active': orderInfo }">
            <strong>2</strong>
            <span>{{ t("ktpMove.stepRack") }}</span>
          </div>
        </div>

        <div class="scan-title">
          <span>{{ t("ktpMove.scanTitle") }}</span>
          <small>
            {{ orderInfo ? t("ktpMove.scanHintWithOrder") : t("ktpMove.scanHintStart") }}
          </small>
        </div>

        <el-input
          ref="scanInputRef"
          v-model="scanInput"
          class="scan-input"
          :placeholder="t('ktpMove.scanPlaceholder')"
          clearable
          inputmode="text"
          autocomplete="off"
          autocapitalize="characters"
          :disabled="isBusy"
          @keyup.enter="processScan()"
        >
          <template #append>
            <el-button
              class="send-btn"
              :loading="isBusy"
              @click="processScan()"
            >
              <el-icon><Promotion /></el-icon>
            </el-button>
          </template>
        </el-input>
      </div>

      <div class="camera-panel" :class="cameraFeedbackClass">
        <div class="scan-feedback" :class="feedbackClass">
          <el-icon><component :is="feedbackIcon" /></el-icon>
          <strong>{{ feedbackText }}</strong>
        </div>
        <div class="camera-frame">
          <video ref="videoRef" class="camera-video" muted playsinline></video>
          <div v-if="cameraScanEffect" class="camera-scan-effect">
            <el-icon><SuccessFilled /></el-icon>
            <span>{{ t("ktpMove.scanReceived") }}</span>
          </div>
          <div class="scan-guide">
            <span></span>
            <span></span>
            <span></span>
            <span></span>
          </div>
        </div>
        <p>{{ cameraStatus }}</p>
      </div>

      <div v-if="orderInfo" class="order-strip">
        <div class="order-cell order-main">
          <span>{{ t("ktpMove.orderNo") }}</span>
          <strong>{{ orderInfo.ddbh }}</strong>
        </div>
        <div class="order-cell">
          <span>{{ t("ktpMove.currentRack") }}</span>
          <strong>{{ currentRackText }}</strong>
        </div>
        <div class="order-cell">
          <span>{{ t("ktpMove.cartonCount") }}</span>
          <strong>{{ formatNumber(orderInfo.codebarCount) }}</strong>
        </div>
        <div class="order-cell">
          <span>{{ t("ktpMove.quantity") }}</span>
          <strong>{{ formatNumber(orderInfo.totalQty) }}</strong>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { ElMessage } from "element-plus";
import {
  ArrowLeft,
  CircleCloseFilled,
  Promotion,
  SuccessFilled,
  WarningFilled,
} from "@element-plus/icons-vue";
import { BrowserMultiFormatReader } from "@zxing/browser";
import { useWarehouseMapI18n } from "@/hooks/useWarehouseMapI18n";
import { ktpLayout } from "@/views/components/WarehouseKTP/ktpLayout";

const router = useRouter();
const API_URL = import.meta.env.VITE_API_URL;
const API_BASE = `${API_URL}/warehouse-ktp`;
const { lang, langOptions, setLang, t } = useWarehouseMapI18n();

const scanInput = ref("");
const activeCartonBar = ref("");
const orderInfo = ref(null);
const loadingLookup = ref(false);
const moving = ref(false);
const scanInputRef = ref(null);
const videoRef = ref(null);
const activeCamera = ref(false);
const cameraStatus = ref("");
const cameraLocked = ref(false);
const cameraControls = ref(null);
const feedbackType = ref("idle");
const feedbackText = ref(t("ktpMove.cameraGuide"));
const cameraScanEffect = ref(false);
let feedbackTimer = null;
let cameraEffectTimer = null;

const codeReader = new BrowserMultiFormatReader();
const validRackCodes = new Set(ktpLayout.racks.map((rack) => rack.code));

const isBusy = computed(() => loadingLookup.value || moving.value);
const feedbackClass = computed(() => `is-${feedbackType.value}`);
const cameraFeedbackClass = computed(() =>
  feedbackType.value === "idle" ? "" : `camera-${feedbackType.value}`,
);
const feedbackIcon = computed(() => {
  if (feedbackType.value === "success") return SuccessFilled;
  if (feedbackType.value === "error") return CircleCloseFilled;
  if (feedbackType.value === "warning") return WarningFilled;
  return Promotion;
});

const currentRackText = computed(() => {
  const racks = orderInfo.value?.currentRacks || [];
  if (!racks.length) return t("ktpMove.noRack");
  return racks.map((rack) => rack.rackCode).join(", ");
});

const formatMessage = (path, params = {}) => {
  let text = t(path);
  Object.entries(params).forEach(([key, value]) => {
    text = text.replaceAll(`{${key}}`, value);
  });
  return text;
};

const focusScanInput = () => {
  nextTick(() => {
    scanInputRef.value?.focus?.();
  });
};

const goBackToMap = () => {
  router.push("/warehouse-ktp");
};

const formatNumber = (value) => Number(value || 0).toLocaleString("vi-VN");

const extractRackCode = (value) => {
  const normalized = String(value || "").trim().toUpperCase();
  if (validRackCodes.has(normalized)) return normalized;

  const matches = normalized.match(/[A-Z]+\d{2}/g) || [];
  return matches.find((code) => validRackCodes.has(code)) || "";
};

const getErrorMessage = (error, fallback) => {
  return error.response?.data?.message || error.message || fallback;
};

const playBeep = (type = "success") => {
  try {
    const AudioContext = window.AudioContext || window.webkitAudioContext;
    if (!AudioContext) return;
    const context = new AudioContext();
    const masterGain = context.createGain();
    masterGain.gain.value = 0.28;
    masterGain.connect(context.destination);

    const notes =
      type === "success"
        ? [
            { frequency: 1046, start: 0, duration: 0.11 },
            { frequency: 1397, start: 0.14, duration: 0.13 },
          ]
        : [
            { frequency: 220, start: 0, duration: 0.24 },
            { frequency: 165, start: 0.26, duration: 0.24 },
          ];

    notes.forEach((note) => {
      const oscillator = context.createOscillator();
      const gain = context.createGain();
      oscillator.type = type === "success" ? "triangle" : "sawtooth";
      oscillator.frequency.value = note.frequency;
      gain.gain.setValueAtTime(0.001, context.currentTime + note.start);
      gain.gain.exponentialRampToValueAtTime(1, context.currentTime + note.start + 0.015);
      gain.gain.exponentialRampToValueAtTime(0.001, context.currentTime + note.start + note.duration);
      oscillator.connect(gain);
      gain.connect(masterGain);
      oscillator.start(context.currentTime + note.start);
      oscillator.stop(context.currentTime + note.start + note.duration);
    });

    const lastNote = notes[notes.length - 1];
    window.setTimeout(() => context.close(), (lastNote.start + lastNote.duration + 0.1) * 1000);
  } catch (error) {
    // Audio feedback is best-effort only.
  }
};

const triggerCameraDetectedEffect = () => {
  cameraScanEffect.value = false;
  if (cameraEffectTimer) {
    clearTimeout(cameraEffectTimer);
  }

  requestAnimationFrame(() => {
    cameraScanEffect.value = true;
    cameraEffectTimer = setTimeout(() => {
      cameraScanEffect.value = false;
    }, 850);
  });
};

const showSuccessMessage = (message) => {
  ElMessage({
    message,
    type: "success",
    duration: 2800,
    showClose: true,
    customClass: "move-order-success-message",
  });
};

const showFeedback = (type, text) => {
  feedbackType.value = type;
  feedbackText.value = text;
  playBeep(type);

  if (feedbackTimer) {
    clearTimeout(feedbackTimer);
  }
  feedbackTimer = setTimeout(() => {
    feedbackType.value = "idle";
    feedbackText.value = activeCamera.value ? t("ktpMove.cameraGuide") : t("ktpMove.cameraClosed");
  }, type === "success" ? 2600 : 1900);
};

const processScan = async (rawValue = scanInput.value) => {
  const scannedValue = String(rawValue || "").trim();
  scanInput.value = "";

  if (!scannedValue) {
    showFeedback("warning", t("ktpMove.noScanCode"));
    ElMessage.warning(t("ktpMove.noScanCode"));
    focusScanInput();
    return;
  }

  const rackCode = extractRackCode(scannedValue);
  if (rackCode) {
    await moveOrder(rackCode);
    return;
  }

  await lookupOrder(scannedValue);
};

const lookupOrder = async (cartonBar) => {
  loadingLookup.value = true;
  try {
    const res = await axios.get(`${API_BASE}/orders/by-carton`, {
      params: { cartonBar },
    });
    orderInfo.value = res.data.data;
    activeCartonBar.value = cartonBar;
    const message = formatMessage("ktpMove.cartonAccepted", { order: orderInfo.value.ddbh });
    showFeedback("success", message);
    showSuccessMessage(message);
    focusScanInput();
  } catch (error) {
    orderInfo.value = null;
    activeCartonBar.value = "";
    const message = getErrorMessage(error, t("ktpMove.orderNotFound"));
    showFeedback("error", message);
    ElMessage.error(message);
    focusScanInput();
  } finally {
    loadingLookup.value = false;
  }
};

const moveOrder = async (newRackCode) => {
  if (!orderInfo.value || !activeCartonBar.value) {
    showFeedback("warning", t("ktpMove.scanCartonFirst"));
    ElMessage.warning(t("ktpMove.scanCartonFirst"));
    focusScanInput();
    return;
  }

  const currentRacks = orderInfo.value.currentRacks || [];
  if (currentRacks.some((rack) => rack.rackCode === newRackCode)) {
    showFeedback("warning", t("ktpMove.sameRack"));
    ElMessage.warning(t("ktpMove.sameRack"));
    focusScanInput();
    return;
  }

  moving.value = true;
  try {
    const res = await axios.post(`${API_BASE}/orders/move`, {
      cartonBar: activeCartonBar.value,
      newRackCode,
    });
    orderInfo.value = res.data.data;
    const message = formatMessage("ktpMove.moveSuccess", {
      order: orderInfo.value.ddbh,
      rack: newRackCode,
    });
    showFeedback("success", message);
    showSuccessMessage(message);
    focusScanInput();
  } catch (error) {
    const message = getErrorMessage(error, t("ktpMove.moveFailed"));
    showFeedback("error", message);
    ElMessage.error(message);
    focusScanInput();
  } finally {
    moving.value = false;
  }
};

const startCamera = async () => {
  if (cameraControls.value) return;
  activeCamera.value = true;
  cameraLocked.value = false;
  cameraStatus.value = t("ktpMove.cameraOpening");
  feedbackText.value = t("ktpMove.cameraOpening");

  await nextTick();
  if (!videoRef.value) {
    cameraStatus.value = t("ktpMove.cameraFrameMissing");
    showFeedback("error", t("ktpMove.cameraOpenFailed"));
    return;
  }

  try {
    const controls = await codeReader.decodeFromConstraints(
      {
        audio: false,
        video: {
          facingMode: { ideal: "environment" },
          width: { ideal: 1280 },
          height: { ideal: 720 },
        },
      },
      videoRef.value,
      (result) => {
        if (!result || cameraLocked.value) return;
        cameraLocked.value = true;
        triggerCameraDetectedEffect();
        const scannedText = result.getText();
        processScan(scannedText).finally(() => {
          window.setTimeout(() => {
            cameraLocked.value = false;
          }, 1200);
        });
      },
    );

    cameraControls.value = controls;
    cameraStatus.value = t("ktpMove.cameraGuide");
    feedbackType.value = "idle";
    feedbackText.value = t("ktpMove.cameraGuide");
  } catch (error) {
    activeCamera.value = false;
    cameraStatus.value = t("ktpMove.cameraOpenFailed");
    showFeedback("error", t("ktpMove.cameraOpenFailed"));
    ElMessage.error(t("ktpMove.cameraPermissionHelp"));
    focusScanInput();
  }
};

const stopCamera = () => {
  if (cameraControls.value) {
    cameraControls.value.stop();
    cameraControls.value = null;
  }
  activeCamera.value = false;
  cameraStatus.value = "";
  cameraLocked.value = false;
};

onMounted(() => {
  focusScanInput();
  startCamera();
});

onBeforeUnmount(() => {
  if (feedbackTimer) {
    clearTimeout(feedbackTimer);
  }
  if (cameraEffectTimer) {
    clearTimeout(cameraEffectTimer);
  }
  stopCamera();
});
</script>

<style scoped lang="scss">
.move-order-page {
  min-height: 100vh;
  background: #eef2f7;
}

.move-shell {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-width: 820px;
  margin: 0 auto;
  padding: 10px;
}

.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.back-btn {
  font-family: inherit;
  flex: 0 1 auto;
  min-height: 42px;
  padding: 0 14px;
  border: 1px solid #1d4ed8;
  border-radius: 6px;
  background: #ffffff;
  color: #1d4ed8;
  font-size: 15px;
  font-weight: 900;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.12);
}

.lang-switch {
  display: flex;
  flex: 0 0 auto;
  gap: 4px;
  padding: 3px;
  border: 1px solid #cbd5e1;
  border-radius: 999px;
  background: #ffffff;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.08);
}

.lang-btn {
  min-width: 40px;
  height: 34px;
  padding: 0 10px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: #475569;
  font-size: 13px;
  font-weight: 900;
  cursor: pointer;
}

.lang-btn.is-active {
  background: #1d4ed8;
  color: #ffffff;
}

.scan-box,
.camera-panel,
.order-strip {
  border: 1px solid #d5dde8;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.08);
}

.scan-box {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 10px;
}

.step-line {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 8px;
}

.step-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #64748b;
  font-size: 13px;
  font-weight: 900;
  white-space: nowrap;
}

.step-item strong {
  display: grid;
  place-items: center;
  width: 24px;
  height: 24px;
  border-radius: 999px;
  background: #cbd5e1;
  color: #334155;
  font-size: 13px;
}

.step-item.is-active {
  color: #0f172a;
}

.step-item.is-active strong {
  background: #2563eb;
  color: #ffffff;
}

.step-connector {
  height: 2px;
  border-radius: 999px;
  background: #cbd5e1;
}

.scan-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  color: #0f172a;
  font-size: 15px;
  font-weight: 900;
}

.scan-title small {
  color: #64748b;
  font-size: 12px;
  font-weight: 900;
  text-align: right;
}

.scan-input {
  :deep(.el-input__wrapper) {
    min-height: 48px;
    padding: 0 10px;
    border-radius: 6px 0 0 6px;
    box-shadow: 0 0 0 1px #b8c3d4 inset;
  }

  :deep(.el-input__inner) {
    height: 48px;
    color: #0f172a;
    font-size: 19px;
    font-weight: 900;
  }

  :deep(.el-input-group__append) {
    padding: 0;
    border-radius: 0 6px 6px 0;
    background: #2563eb;
    box-shadow: none;
  }
}

.send-btn {
  width: 48px;
  min-height: 48px;
  margin: 0;
  border: 0;
  border-radius: 0 6px 6px 0;
  background: #2563eb;
  color: #ffffff;
  font-size: 19px;
}

.camera-panel {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  transition:
    border-color 0.15s ease,
    box-shadow 0.15s ease,
    background-color 0.15s ease;
}

.camera-success {
  border-color: #16a34a;
  box-shadow:
    0 0 0 4px rgba(22, 163, 74, 0.28),
    0 10px 26px rgba(22, 163, 74, 0.2);
}

.camera-error {
  border-color: #dc2626;
  box-shadow:
    0 0 0 3px rgba(220, 38, 38, 0.22),
    0 8px 22px rgba(220, 38, 38, 0.16);
}

.camera-warning {
  border-color: #d97706;
  box-shadow:
    0 0 0 3px rgba(217, 119, 6, 0.22),
    0 8px 22px rgba(217, 119, 6, 0.16);
}

.scan-feedback {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 46px;
  padding: 10px 12px;
  border-radius: 7px;
  background: #eef2f7;
  color: #334155;
  font-size: 16px;
  font-weight: 900;
}

.scan-feedback strong {
  min-width: 0;
  overflow-wrap: anywhere;
}

.scan-feedback .el-icon {
  flex: 0 0 auto;
  font-size: 22px;
}

.scan-feedback.is-success {
  min-height: 54px;
  border: 2px solid #16a34a;
  background: #bbf7d0;
  color: #166534;
  font-size: 18px;
  box-shadow: inset 0 0 0 1px rgba(22, 163, 74, 0.16);
}

.scan-feedback.is-success .el-icon {
  font-size: 28px;
}

.scan-feedback.is-error {
  background: #fee2e2;
  color: #991b1b;
}

.scan-feedback.is-warning {
  background: #fef3c7;
  color: #92400e;
}

.camera-frame {
  position: relative;
  overflow: hidden;
  border-radius: 7px;
  background: #0f172a;
}

.camera-scan-effect {
  position: absolute;
  inset: 0;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background:
    radial-gradient(circle at center, rgba(34, 197, 94, 0.36), transparent 42%),
    rgba(22, 163, 74, 0.14);
  color: #ffffff;
  font-size: 22px;
  font-weight: 900;
  pointer-events: none;
  animation: cameraScanPop 0.85s ease-out forwards;
}

.camera-scan-effect::before {
  content: "";
  position: absolute;
  inset: 14%;
  border: 4px solid rgba(255, 255, 255, 0.92);
  border-radius: 10px;
  animation: cameraScanRing 0.85s ease-out forwards;
}

.camera-scan-effect .el-icon {
  font-size: 34px;
}

.camera-scan-effect span,
.camera-scan-effect .el-icon {
  position: relative;
  z-index: 1;
  text-shadow: 0 2px 10px rgba(15, 23, 42, 0.42);
}

:global(.move-order-success-message) {
  min-width: min(92vw, 460px);
  padding: 14px 18px;
  border: 2px solid #16a34a;
  background: #dcfce7;
  color: #14532d;
  font-size: 17px;
  font-weight: 900;
}

.camera-video {
  display: block;
  width: 100%;
  max-height: 40vh;
  object-fit: cover;
  aspect-ratio: 16 / 9;
}

.scan-guide {
  position: absolute;
  inset: 18%;
  pointer-events: none;
}

.scan-guide span {
  position: absolute;
  width: 32px;
  height: 32px;
  border-color: rgba(255, 255, 255, 0.92);
}

.scan-guide span:nth-child(1) {
  top: 0;
  left: 0;
  border-top: 4px solid;
  border-left: 4px solid;
}

.scan-guide span:nth-child(2) {
  top: 0;
  right: 0;
  border-top: 4px solid;
  border-right: 4px solid;
}

.scan-guide span:nth-child(3) {
  right: 0;
  bottom: 0;
  border-right: 4px solid;
  border-bottom: 4px solid;
}

.scan-guide span:nth-child(4) {
  bottom: 0;
  left: 0;
  border-bottom: 4px solid;
  border-left: 4px solid;
}

.camera-panel p {
  margin: 0;
  color: #475569;
  font-size: 12px;
  font-weight: 800;
}

@keyframes cameraScanPop {
  0% {
    opacity: 0;
    transform: scale(0.98);
  }

  18% {
    opacity: 1;
    transform: scale(1);
  }

  100% {
    opacity: 0;
    transform: scale(1.03);
  }
}

@keyframes cameraScanRing {
  0% {
    opacity: 0;
    transform: scale(0.82);
  }

  24% {
    opacity: 1;
  }

  100% {
    opacity: 0;
    transform: scale(1.08);
  }
}

.order-strip {
  display: grid;
  grid-template-columns: 1.55fr 1.1fr 0.72fr 0.78fr;
  overflow: hidden;
}

.order-cell {
  min-width: 0;
  padding: 10px;
  border-left: 1px solid #e2e8f0;
}

.order-cell:first-child {
  border-left: 0;
}

.order-cell span {
  display: block;
  color: #64748b;
  font-size: 11px;
  font-weight: 900;
  line-height: 1;
  text-transform: uppercase;
}

.order-cell strong {
  display: block;
  margin-top: 6px;
  color: #0f172a;
  font-size: 18px;
  font-weight: 900;
  line-height: 1.08;
  overflow-wrap: anywhere;
}

.order-main strong {
  color: #b91c1c;
  font-size: 19px;
}

@media (max-width: 560px) {
  .move-shell {
    padding: 7px;
  }

  .top-bar {
    align-items: stretch;
    flex-wrap: wrap;
  }

  .back-btn {
    flex: 1 1 230px;
    justify-content: center;
  }

  .lang-switch {
    flex: 1 0 auto;
    justify-content: center;
  }

  .scan-title {
    align-items: flex-start;
    flex-direction: column;
    gap: 3px;
  }

  .scan-title small {
    text-align: left;
  }

  .order-strip {
    grid-template-columns: 1fr 1fr;
  }

  .order-cell:nth-child(odd) {
    border-left: 0;
  }

  .order-cell:nth-child(n + 3) {
    border-top: 1px solid #e2e8f0;
  }
}
</style>
