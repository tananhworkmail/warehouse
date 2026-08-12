<!-- src/views/LoginPage.vue -->
<template>
  <div class="login-bg">
    <div class="login-card">
      <!-- Brand -->
      <div class="login-brand">
        <div class="brand-mark">
          <img src="@/assets/Logo.png" alt="Logo" class="login-logo" />
        </div>
        <h2 class="login-title">Loss Tem Size</h2>
        <p class="login-sub">Please sign in to continue</p>
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
import { ref } from "vue"
import { useRouter, useRoute } from "vue-router"

const API_BASE_URL = (import.meta.env.VITE_API_URL || "http://192.168.71.87:3084/api/v1").replace(/\/$/, "")

const router = useRouter()
const route = useRoute()

const username = ref("")
const password = ref("")
const branch = ref("VDH")
const branches = ["VDH"]
const showPassword = ref(false)
const loading = ref(false)
const errorMsg = ref("")

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

    localStorage.setItem("token", data.token || "")
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