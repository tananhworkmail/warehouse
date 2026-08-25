<!-- src/views/LoginPage.vue -->
<template>
  <div class="login-bg" :class="{ 'is-ktp-scan': isKTPScanInLogin }">
    <div v-if="isKTPScanInLogin" class="warehouse-vector-scene" aria-hidden="true">
      <svg viewBox="0 0 1600 900" preserveAspectRatio="xMidYMid slice">
        <defs>
          <linearGradient id="warehouse-sky" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0" stop-color="#071c19" />
            <stop offset="1" stop-color="#0b3027" />
          </linearGradient>
          <linearGradient id="warehouse-floor" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0" stop-color="#163c34" />
            <stop offset="1" stop-color="#061512" />
          </linearGradient>
          <linearGradient id="carton-color" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0" stop-color="#ffc65b" />
            <stop offset="1" stop-color="#d88a2a" />
          </linearGradient>
          <filter id="soft-glow" x="-80%" y="-80%" width="260%" height="260%">
            <feGaussianBlur stdDeviation="8" result="blur" />
            <feMerge><feMergeNode in="blur" /><feMergeNode in="SourceGraphic" /></feMerge>
          </filter>
          <pattern id="floor-grid" width="90" height="60" patternUnits="userSpaceOnUse">
            <path d="M 90 0 L 0 0 0 60" fill="none" stroke="#58b99a" stroke-opacity=".12" stroke-width="2" />
          </pattern>
        </defs>

        <rect width="1600" height="620" fill="url(#warehouse-sky)" />
        <path d="M0 610 L1600 610 L1600 900 L0 900 Z" fill="url(#warehouse-floor)" />
        <path d="M0 610 L1600 610 L1600 900 L0 900 Z" fill="url(#floor-grid)" />

        <g class="warehouse-lights" fill="#a7ffe0">
          <rect x="170" y="72" width="260" height="12" rx="6" opacity=".72" />
          <rect x="670" y="72" width="260" height="12" rx="6" opacity=".48" />
          <rect x="1170" y="72" width="260" height="12" rx="6" opacity=".72" />
        </g>

        <g class="warehouse-rack warehouse-rack--left">
          <path d="M45 145 H410 V650 H45 Z M45 315 H410 M45 480 H410" fill="none" stroke="#39a58d" stroke-width="18" />
          <g fill="url(#carton-color)" stroke="#8f591f" stroke-width="3">
            <rect x="72" y="188" width="138" height="105" rx="5" />
            <rect x="225" y="177" width="155" height="116" rx="5" />
            <rect x="72" y="348" width="170" height="110" rx="5" />
            <rect x="257" y="360" width="123" height="98" rx="5" />
            <rect x="72" y="518" width="128" height="105" rx="5" />
            <rect x="215" y="511" width="165" height="112" rx="5" />
          </g>
          <g stroke="#ffe29a" stroke-width="5" opacity=".72">
            <path d="M141 188 V293 M302 177 V293 M157 348 V458 M318 360 V458 M136 518 V623 M298 511 V623" />
          </g>
        </g>

        <g class="warehouse-rack warehouse-rack--right" opacity=".64">
          <path d="M1328 145 H1570 V585 H1328 Z M1328 292 H1570 M1328 438 H1570" fill="none" stroke="#39a58d" stroke-width="16" />
          <g fill="#bd762a" stroke="#7c4a1d" stroke-width="3">
            <rect x="1350" y="174" width="86" height="96" rx="4" />
            <rect x="1448" y="174" width="98" height="96" rx="4" />
            <rect x="1350" y="318" width="112" height="98" rx="4" />
            <rect x="1474" y="318" width="72" height="98" rx="4" />
            <rect x="1350" y="463" width="92" height="96" rx="4" />
            <rect x="1454" y="463" width="92" height="96" rx="4" />
          </g>
        </g>

        <g class="conveyor">
          <path d="M385 690 H1490 L1540 744 H435 Z" fill="#244d48" stroke="#52b7a0" stroke-width="8" />
          <path d="M435 744 H1540 V790 H435 Z" fill="#0b2421" stroke="#2b6f62" stroke-width="7" />
          <g fill="#9bcfc3" stroke="#163d37" stroke-width="5">
            <circle cx="475" cy="718" r="22" /><circle cx="570" cy="718" r="22" />
            <circle cx="665" cy="718" r="22" /><circle cx="760" cy="718" r="22" />
            <circle cx="855" cy="718" r="22" /><circle cx="950" cy="718" r="22" />
            <circle cx="1045" cy="718" r="22" /><circle cx="1140" cy="718" r="22" />
            <circle cx="1235" cy="718" r="22" /><circle cx="1330" cy="718" r="22" />
            <circle cx="1425" cy="718" r="22" />
          </g>
          <path d="M520 790 L485 884 M1435 790 L1470 884" stroke="#173d37" stroke-width="20" />
        </g>

        <g class="moving-carton">
          <animateTransform attributeName="transform" type="translate" values="150 0;-760 0" dur="8s" repeatCount="indefinite" />
          <rect x="1080" y="565" width="205" height="125" rx="9" fill="url(#carton-color)" stroke="#8f591f" stroke-width="5" />
          <path d="M1182 565 V690" stroke="#ffe29a" stroke-width="8" opacity=".78" />
          <rect x="1210" y="600" width="58" height="55" rx="4" fill="#fff8dc" />
          <g stroke="#25382f" stroke-width="4">
            <path d="M1220 610 V646 M1228 610 V646 M1238 610 V646 M1248 610 V646 M1258 610 V646" />
          </g>
        </g>

        <g class="scanner" transform="translate(1320 470) rotate(-12)">
          <path d="M115 55 C170 42 210 62 238 100 L196 131 C169 105 143 95 104 103 Z" fill="#2e6975" />
          <rect x="0" y="0" width="142" height="82" rx="23" fill="#182c32" stroke="#68d4c7" stroke-width="6" />
          <rect x="32" y="63" width="52" height="112" rx="19" fill="#132126" stroke="#477c7c" stroke-width="5" />
          <rect x="13" y="17" width="30" height="46" rx="8" fill="#ff6359" filter="url(#soft-glow)" />
        </g>

        <g class="scan-beam" filter="url(#soft-glow)">
          <path d="M1327 509 L1090 604 L1090 636 L1327 533 Z" fill="#ff594f" opacity=".18">
            <animate attributeName="opacity" values=".08;.42;.08" dur="1.25s" repeatCount="indefinite" />
          </path>
          <path d="M1327 521 L1090 620" stroke="#ff7065" stroke-width="6" stroke-linecap="round">
            <animate attributeName="stroke-width" values="3;8;3" dur="1.25s" repeatCount="indefinite" />
          </path>
        </g>

        <g class="scan-success" transform="translate(1118 520)" filter="url(#soft-glow)">
          <circle r="43" fill="#20d889" opacity=".18">
            <animate attributeName="r" values="34;53;34" dur="1.5s" repeatCount="indefinite" />
            <animate attributeName="opacity" values=".3;.05;.3" dur="1.5s" repeatCount="indefinite" />
          </circle>
          <circle r="28" fill="#0d3c30" stroke="#52f0af" stroke-width="5" />
          <path d="M-13 0 L-3 11 L15 -12" fill="none" stroke="#8bffd0" stroke-width="7" stroke-linecap="round" stroke-linejoin="round" />
        </g>
      </svg>
    </div>
    <div class="login-card">
      <!-- Brand -->
      <div class="login-brand">
        <div class="brand-mark">
          <img src="@/assets/Logo.png" alt="Logo" class="login-logo" />
        </div>
        <h2 class="login-title">{{ loginTitle }}</h2>
        <p class="login-sub">Đăng nhập để tiếp tục</p>
      </div>

      <!-- Form -->
      <div class="login-form">
        <!-- Branch -->
        <div class="login-field">
          <label class="login-label">Branch</label>
          <el-select
            v-model="branch"
            class="login-select"
            :class="{ 'is-selected': !!branch }"
            size="large"
            placeholder="Select branch"
          >
            <el-option v-for="b in branches" :key="b" :label="b" :value="b" />
          </el-select>
          <div v-if="branch" class="branch-hint">
            <span class="branch-dot"></span>
            Đã chọn: <strong>{{ branch }}</strong>
          </div>
        </div>

        <!-- User ID -->
        <div class="login-field">
          <label class="login-label">User ID</label>
          <el-input
            v-model="username"
            placeholder="Enter your User ID"
            size="large"
            class="login-input"
            @keyup.enter="handleSubmit"
          >
            <template #prefix>
              <svg
                class="input-icon"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                <circle cx="12" cy="7" r="4" />
              </svg>
            </template>
          </el-input>
        </div>

        <!-- Password -->
        <div class="login-field">
          <label class="login-label">Password</label>
          <el-input
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            placeholder="Enter your password"
            size="large"
            class="login-input"
            @keyup.enter="handleSubmit"
          >
            <template #prefix>
              <svg
                class="input-icon"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
            </template>
            <template #suffix>
              <span class="eye-btn" @click="showPassword = !showPassword">
                <svg
                  v-if="showPassword"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  class="input-icon"
                >
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                  <circle cx="12" cy="12" r="3" />
                </svg>
                <svg
                  v-else
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  class="input-icon"
                >
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
                  <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
                  <line x1="1" y1="1" x2="23" y2="23" />
                </svg>
              </span>
            </template>
          </el-input>
        </div>

        <!-- Error message -->
        <div v-if="errorMsg" class="login-error">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            class="input-icon"
          >
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="8" x2="12" y2="12" />
            <line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
          {{ errorMsg }}
        </div>

        <!-- Submit -->
        <button
          class="login-btn"
          :class="{ 'is-loading': loading }"
          :disabled="!username || !password || !branch || loading"
          @click="handleSubmit"
        >
          <span v-if="loading" class="btn-spinner" />
          <span>{{ loading ? "Signing in..." : "Sign In" }}</span>
        </button>
      </div>
    </div>

    <!-- Toast -->
    <transition name="toast-fade">
      <div v-if="showToastMsg" class="login-toast" :class="`toast--${toastColor}`">
        {{ toastMsg }}
      </div>
    </transition>
  </div>
</template>

<script setup>
import { computed, ref } from "vue"
import { useRouter, useRoute } from "vue-router"

const API_BASE_URL = (import.meta.env.VITE_API_URL || "http://192.168.71.87:3084/api/v1").replace(/\/$/, "")

const router = useRouter()
const route = useRoute()
const isKTPScanInLogin = computed(() =>
  String(route.query.redirect || "").startsWith("/warehouse-ktp/scan-in")
)
const loginTitle = computed(() =>
  isKTPScanInLogin.value ? "Scan In Kho thành phẩm" : "Loss Tem Size"
)

const username = ref("")
const password = ref("")
const branch = ref("VDH")
const branches = ["VDH"]
const showPassword = ref(false)
const loading = ref(false)
const errorMsg = ref(
  route.query.reason === "no-scan-in-permission"
    ? "Tài khoản không có quyền Scan In KTP. Vui lòng đăng nhập bằng tài khoản khác."
    : ""
)
const authStorageKeys = ["token", "userid", "username", "role"]

const clearStoredLogin = () => {
  authStorageKeys.forEach((key) => localStorage.removeItem(key))
}

// Toast
const showToastMsg = ref(false)
const toastMsg = ref("")
const toastColor = ref("green")
let toastTimer = null

const showToast = (msg, color = "green") => {
  toastMsg.value = msg
  toastColor.value = color
  showToastMsg.value = true
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    showToastMsg.value = false
  }, 2500)
}

const handleSubmit = async () => {
  if (!username.value || !password.value || !branch.value) return
  errorMsg.value = ""
  loading.value = true

  try {
    const res = await fetch(`${API_BASE_URL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        GSBH: branch.value,
        USERID: username.value,
        PWD: password.value,
      }),
    })

    const data = await res.json()

    if (!res.ok) {
      const msg = data?.error || data?.message || "Đăng nhập thất bại"
      errorMsg.value = msg
      showToast(msg, "red")
      return
    }

    const loginToken = data.token || ""
    if (!loginToken) {
      throw new Error("Máy chủ không trả về phiên đăng nhập hợp lệ")
    }

    if (isKTPScanInLogin.value) {
      const accessResponse = await fetch(`${API_BASE_URL}/warehouse-ktp/scan-in/access`, {
        headers: { Authorization: `Bearer ${loginToken}` },
      })
      if (!accessResponse.ok) {
        clearStoredLogin()
        const msg = accessResponse.status === 403
          ? "Tài khoản không có quyền Scan In KTP. Vui lòng đăng nhập bằng tài khoản khác."
          : "Không thể xác nhận quyền Scan In KTP. Vui lòng thử đăng nhập lại."
        errorMsg.value = msg
        showToast(msg, "red")
        username.value = ""
        password.value = ""
        return
      }
    }

    localStorage.setItem("token", loginToken)
    localStorage.setItem("userid", data.USERID || username.value)
    localStorage.setItem("username", data.USERNAME || username.value)
    localStorage.setItem("role", data.role || "")

    showToast("Đăng nhập thành công! 🎉", "green")

    const redirectTo = route.query.redirect || "/loss-tem-size"
    setTimeout(() => router.push(redirectTo), 1000)
  } catch (err) {
    errorMsg.value = err.message || "Network error"
    showToast(errorMsg.value, "red")
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-bg {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  padding: 16px;

  /* nền xanh kho/xưởng */
  background:
    radial-gradient(
      circle at top left,
      rgba(34, 197, 94, 0.22) 0%,
      transparent 28%
    ),
    radial-gradient(
      circle at top right,
      rgba(22, 163, 74, 0.12) 0%,
      transparent 32%
    ),
    radial-gradient(
      circle at bottom,
      rgba(74, 222, 128, 0.08) 0%,
      transparent 35%
    ),
    linear-gradient(
      135deg,
      #020b05 0%,
      #07160d 30%,
      #0d2215 65%,
      #061109 100%
    );
}

.login-bg.is-ktp-scan {
  background: #061713;
}

.warehouse-vector-scene {
  position: absolute;
  inset: 0;
  z-index: 0;
  overflow: hidden;
}

.warehouse-vector-scene svg {
  width: 100%;
  height: 100%;
  display: block;
}
/* lớp grid kiểu dashboard/kho */
.login-bg::before {
  content: "";
  position: absolute;
  inset: 0;

  background-image:
    linear-gradient(
      rgba(90, 255, 140, 0.05) 1px,
      transparent 1px
    ),
    linear-gradient(
      90deg,
      rgba(90, 255, 140, 0.05) 1px,
      transparent 1px
    );

  background-size: 45px 45px;
  opacity: 0.5;
}

.login-bg.is-ktp-scan::before {
  background:
    radial-gradient(circle at center, transparent 25%, rgba(0, 8, 6, 0.35) 100%),
    linear-gradient(rgba(64, 255, 180, 0.025) 1px, transparent 1px),
    linear-gradient(90deg, rgba(64, 255, 180, 0.025) 1px, transparent 1px);
  background-size: auto, 56px 56px, 56px 56px;
  opacity: 1;
}

/* lớp tối overlay */
.login-bg::after {
  content: "";
  position: absolute;
  inset: 0;
  background:
    linear-gradient(
      rgba(0, 0, 0, 0.15),
      rgba(0, 0, 0, 0.45)
    );
}

.login-bg.is-ktp-scan::after {
  background: linear-gradient(
    to bottom,
    rgba(0, 7, 6, 0.08),
    rgba(0, 7, 6, 0.28)
  );
}

/* ── Card ── */
.login-card {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 430px;
  padding: 36px 32px 30px;
  border-radius: 24px;
  background: rgba(8, 20, 12, 0.88);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  border: 1px solid rgba(74, 222, 128, 0.22);
  box-shadow:
    0 0 0 1px rgba(34, 197, 94, 0.06),
    0 0 34px rgba(34, 197, 94, 0.16),
    0 26px 60px rgba(0, 0, 0, 0.55);
  transition: box-shadow 0.3s ease, transform 0.3s ease;
}

.login-card:hover {
  box-shadow:
    0 0 0 1px rgba(74, 222, 128, 0.14),
    0 0 48px rgba(34, 197, 94, 0.22),
    0 26px 60px rgba(0, 0, 0, 0.55);
  transform: translateY(-1px);
}

/* ── Brand ── */
.login-brand {
  text-align: center;
  margin-bottom: 28px;
}

.brand-mark {
  width: 78px;
  height: 78px;
  margin: 0 auto 12px;
  border-radius: 22px;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.2), rgba(16, 185, 129, 0.08));
  border: 1px solid rgba(74, 222, 128, 0.22);
  box-shadow: 0 0 24px rgba(34, 197, 94, 0.18);
}

.login-logo {
  width: 58px;
  height: 58px;
  object-fit: contain;
  display: block;
  filter: drop-shadow(0 0 10px rgba(74, 222, 128, 0.45));
}

.login-title {
  font-size: 1.7rem;
  font-weight: 800;
  color: #6dff9e;
  text-shadow: 0 0 18px rgba(74, 222, 128, 0.45);
  margin: 0 0 6px;
  letter-spacing: 0.02em;
}

.login-sub {
  font-size: 13px;
  color: rgba(220, 255, 230, 0.62);
  margin: 0;
}

/* ── Form ── */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.login-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.login-label {
  font-size: 12px;
  font-weight: 700;
  color: rgba(134, 255, 181, 0.82);
  letter-spacing: 0.07em;
  text-transform: uppercase;
}

/* Element Plus inputs */
.login-input :deep(.el-input__wrapper),
.login-select :deep(.el-select__wrapper) {
  background: rgba(255, 255, 255, 0.06) !important;
  border: 1.5px solid rgba(74, 222, 128, 0.22) !important;
  border-radius: 12px !important;
  box-shadow: none !important;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background 0.2s ease !important;
}

.login-input :deep(.el-input__wrapper:hover),
.login-select :deep(.el-select__wrapper:hover) {
  border-color: rgba(74, 222, 128, 0.5) !important;
  background: rgba(255, 255, 255, 0.08) !important;
}

.login-input :deep(.el-input__wrapper.is-focus),
.login-select :deep(.el-select__wrapper.is-focused) {
  border-color: #6dff9e !important;
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.16) !important;
}

/* Branch selected state */
.login-select.is-selected :deep(.el-select__wrapper) {
  background: rgba(34, 197, 94, 0.13) !important;
  border-color: rgba(109, 255, 158, 0.72) !important;
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.12) !important;
}

.login-select.is-selected :deep(.el-select__selected-item),
.login-select.is-selected :deep(.el-select__placeholder) {
  color: #d7ffe3 !important;
  font-weight: 700 !important;
}

.login-input :deep(.el-input__inner),
.login-select :deep(.el-select__selected-item) {
  color: #ecfff0 !important;
  font-size: 14px !important;
}

.login-input :deep(.el-input__inner::placeholder) {
  color: rgba(220, 255, 230, 0.36) !important;
}

.login-select :deep(.el-select__placeholder) {
  color: rgba(220, 255, 230, 0.36) !important;
}

.login-input :deep(.el-input__prefix),
.login-input :deep(.el-input__suffix),
.login-select :deep(.el-select__suffix) {
  color: rgba(109, 255, 158, 0.6);
}

.login-select {
  width: 100%;
}

.branch-hint {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 2px;
  font-size: 12px;
  color: rgba(216, 255, 230, 0.72);
}

.branch-hint strong {
  color: #7dffb0;
  font-weight: 800;
}

.branch-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: #6dff9e;
  box-shadow: 0 0 10px rgba(109, 255, 158, 0.7);
  flex: 0 0 auto;
}

/* ── Input icon ── */
.input-icon {
  width: 16px;
  height: 16px;
  color: rgba(109, 255, 158, 0.65);
}

.eye-btn {
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 0 2px;
  transition: transform 0.15s ease, color 0.2s ease;
}

.eye-btn:hover {
  transform: scale(1.04);
}

.eye-btn:hover .input-icon {
  color: #6dff9e;
}

/* ── Error ── */
.login-error {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 9px 12px;
  background: rgba(220, 38, 38, 0.16);
  border: 1px solid rgba(220, 38, 38, 0.32);
  border-radius: 10px;
  color: #fecaca;
  font-size: 12px;
  font-weight: 500;
}

/* ── Button ── */
.login-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  height: 48px;
  margin-top: 4px;
  background: linear-gradient(135deg, #16a34a 0%, #22c55e 45%, #86efac 100%);
  color: #06200f;
  font-size: 15px;
  font-weight: 900;
  letter-spacing: 0.04em;
  border: none;
  border-radius: 14px;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease, filter 0.2s ease;
  box-shadow: 0 6px 22px rgba(34, 197, 94, 0.32);
}

.login-btn:hover:not(:disabled) {
  filter: brightness(1.05);
  box-shadow: 0 10px 30px rgba(34, 197, 94, 0.42);
  transform: translateY(-1px);
}

.login-btn:active:not(:disabled) {
  transform: translateY(0);
}

.login-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
  transform: none;
}

.login-btn.is-loading {
  opacity: 0.9;
}

/* Spinner */
.btn-spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2.5px solid rgba(6, 32, 15, 0.28);
  border-top-color: #06200f;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* ── Toast ── */
.login-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 9999;
  padding: 12px 20px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 700;
  box-shadow: 0 10px 28px rgba(0, 0, 0, 0.4);
  max-width: 340px;
}

.toast--green {
  background: #052e16;
  border: 1px solid #166534;
  color: #bbf7d0;
}

.toast--red {
  background: #1c0b0b;
  border: 1px solid #7f1d1d;
  color: #fecaca;
}

.toast--orange {
  background: #1c1007;
  border: 1px solid #78350f;
  color: #fed7aa;
}

/* Toast transition */
.toast-fade-enter-active,
.toast-fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.toast-fade-enter-from,
.toast-fade-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

/* ── Responsive ── */
@media (max-width: 480px) {
  .login-card {
    padding: 28px 20px 24px;
    border-radius: 20px;
  }

  .login-title {
    font-size: 1.45rem;
  }

  .brand-mark {
    width: 70px;
    height: 70px;
  }

  .login-logo {
    width: 50px;
    height: 50px;
  }
}
</style>
