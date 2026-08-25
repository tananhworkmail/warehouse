<template>
  <main class="login-page" :class="`mode-${mode}`">
    <KTPScanScene />
    <button class="back-button" type="button" @click="router.push('/warehouse-ktp')">
      <span aria-hidden="true">←</span>
      Kho thành phẩm
    </button>

    <section class="login-card">
      <header class="login-brand">
        <div class="brand-mark">
          <img src="@/assets/Logo.png" alt="Lai Yih" class="login-logo" />
        </div>
        <div class="mode-badge">{{ badge }}</div>
        <h1>{{ title }}</h1>
        <p>Đăng nhập để tiếp tục thao tác quét KTP</p>
      </header>

      <form class="login-form" @submit.prevent="submit">
        <label class="login-field">
          <span>Branch</span>
          <el-select v-model="branch" class="login-control branch-control" size="large">
            <el-option v-for="item in branches" :key="item" :label="item" :value="item" />
          </el-select>
          <small><i /> Đã chọn: <strong>{{ branch }}</strong></small>
        </label>

        <label class="login-field">
          <span>User ID</span>
          <el-input
            ref="userInput"
            v-model="userId"
            class="login-control"
            size="large"
            autocomplete="username"
            placeholder="Nhập mã nhân viên"
            @keyup.enter="submit"
          >
            <template #prefix>
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                <circle cx="12" cy="7" r="4" />
              </svg>
            </template>
          </el-input>
        </label>

        <label class="login-field">
          <span>Password</span>
          <el-input
            v-model="password"
            class="login-control"
            size="large"
            :type="showPassword ? 'text' : 'password'"
            autocomplete="current-password"
            placeholder="Nhập mật khẩu"
            @keyup.enter="submit"
          >
            <template #prefix>
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
            </template>
            <template #suffix>
              <button
                class="password-toggle"
                type="button"
                :aria-label="showPassword ? 'Ẩn mật khẩu' : 'Hiện mật khẩu'"
                @click="showPassword = !showPassword"
              >
                <svg v-if="showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8Z" />
                  <circle cx="12" cy="12" r="3" />
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
                  <path d="m1 1 22 22" />
                </svg>
              </button>
            </template>
          </el-input>
        </label>

        <div v-if="errorMessage" class="login-error" role="alert">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" />
            <path d="M12 8v4m0 4h.01" />
          </svg>
          <span>{{ errorMessage }}</span>
        </div>

        <button class="submit-button" type="submit" :disabled="loading || !canSubmit">
          <i v-if="loading" class="spinner" />
          <span>{{ loading ? "Đang đăng nhập…" : submitLabel }}</span>
        </button>
      </form>
    </section>
  </main>
</template>

<script setup>
import { computed, nextTick, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useLogin } from "@/hooks/useLogin";
import KTPScanScene from "./KTPScanScene.vue";

const props = defineProps({
  mode: { type: String, required: true },
  title: { type: String, required: true },
  badge: { type: String, required: true },
  submitLabel: { type: String, required: true },
  defaultRedirect: { type: String, required: true },
  permissionMessage: { type: String, required: true },
});

const router = useRouter();
const route = useRoute();
const userId = ref("");
const password = ref("");
const branch = ref("VDH");
const branches = ["VDH"];
const showPassword = ref(false);
const loading = ref(false);
const userInput = ref(null);
const errorMessage = ref(route.query.reason === "no-permission" ? props.permissionMessage : "");
const canSubmit = computed(() => Boolean(userId.value.trim() && password.value && branch.value));

async function submit() {
  if (!canSubmit.value || loading.value) return;
  loading.value = true;
  errorMessage.value = "";
  try {
    const result = await useLogin({
      GSBH: branch.value,
      USERID: userId.value,
      PWD: password.value,
    });
    if (!result.success) {
      errorMessage.value = result.message || "Đăng nhập không thành công.";
      password.value = "";
      return;
    }
    const redirect = typeof route.query.redirect === "string" && route.query.redirect.startsWith("/")
      ? route.query.redirect
      : props.defaultRedirect;
    await router.replace(redirect);
  } catch (error) {
    errorMessage.value = error?.message || "Không thể đăng nhập. Vui lòng thử lại.";
    password.value = "";
  } finally {
    loading.value = false;
  }
}

onMounted(() => nextTick(() => userInput.value?.focus?.()));
</script>

<style scoped>
.login-page {
  --accent: #6dff9e;
  --accent-rgb: 109, 255, 158;
  position: relative;
  display: grid;
  min-height: 100vh;
  box-sizing: border-box;
  padding: 24px;
  place-items: center;
  overflow: hidden;
  background: #061713;
  color: #ecfff3;
}

.login-page.mode-out {
  --accent: #ffb45e;
  --accent-rgb: 255, 180, 94;
}

.login-page.mode-ship {
  --accent: #67dcff;
  --accent-rgb: 103, 220, 255;
}

.login-page::before {
  content: "";
  position: absolute;
  z-index: 1;
  inset: 0;
  background:
    radial-gradient(circle at center, transparent 20%, rgba(0, 8, 6, .42) 100%),
    linear-gradient(rgba(64, 255, 180, .025) 1px, transparent 1px),
    linear-gradient(90deg, rgba(64, 255, 180, .025) 1px, transparent 1px);
  background-size: auto, 56px 56px, 56px 56px;
  pointer-events: none;
}

.back-button {
  position: fixed;
  z-index: 3;
  top: 18px;
  left: 20px;
  display: inline-flex;
  height: 38px;
  padding: 0 14px;
  border: 1px solid rgba(var(--accent-rgb), .28);
  border-radius: 10px;
  align-items: center;
  gap: 7px;
  background: rgba(5, 24, 19, .72);
  backdrop-filter: blur(12px);
  color: #eafff1;
  font-weight: 800;
  cursor: pointer;
}

.back-button:hover {
  border-color: rgba(var(--accent-rgb), .7);
  color: var(--accent);
}

.login-card {
  position: relative;
  z-index: 2;
  width: min(430px, 100%);
  box-sizing: border-box;
  padding: 30px 32px 28px;
  border: 1px solid rgba(var(--accent-rgb), .25);
  border-radius: 24px;
  background: rgba(8, 25, 20, .9);
  backdrop-filter: blur(18px);
  box-shadow: 0 0 42px rgba(var(--accent-rgb), .13), 0 26px 70px rgba(0, 0, 0, .5);
}

.login-brand {
  margin-bottom: 24px;
  text-align: center;
}

.brand-mark {
  display: grid;
  width: 72px;
  height: 72px;
  margin: 0 auto 11px;
  border: 1px solid rgba(var(--accent-rgb), .3);
  border-radius: 20px;
  place-items: center;
  background: rgba(var(--accent-rgb), .1);
  box-shadow: 0 0 24px rgba(var(--accent-rgb), .18);
}

.login-logo {
  display: block;
  width: 54px;
  height: 54px;
  object-fit: contain;
  filter: drop-shadow(0 0 9px rgba(var(--accent-rgb), .45));
}

.mode-badge {
  display: inline-flex;
  margin-bottom: 7px;
  padding: 4px 9px;
  border-radius: 999px;
  background: rgba(var(--accent-rgb), .12);
  color: var(--accent);
  font-size: 11px;
  font-weight: 900;
  letter-spacing: .08em;
}

h1 {
  margin: 0 0 5px;
  color: var(--accent);
  font-size: 27px;
  line-height: 1.15;
  text-shadow: 0 0 18px rgba(var(--accent-rgb), .34);
}

.login-brand p {
  margin: 0;
  color: rgba(229, 255, 238, .62);
  font-size: 13px;
}

.login-form,
.login-field {
  display: flex;
  flex-direction: column;
}

.login-form {
  gap: 15px;
}

.login-field {
  gap: 6px;
}

.login-field > span {
  color: rgba(var(--accent-rgb), .82);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: .07em;
  text-transform: uppercase;
}

.login-field small {
  display: flex;
  align-items: center;
  gap: 6px;
  color: rgba(225, 255, 236, .68);
  font-size: 11px;
}

.login-field small i {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 9px rgba(var(--accent-rgb), .7);
}

.login-field small strong {
  color: var(--accent);
}

.login-control :deep(.el-input__wrapper),
.login-control :deep(.el-select__wrapper) {
  border: 1.5px solid rgba(var(--accent-rgb), .24) !important;
  border-radius: 12px !important;
  background: rgba(255, 255, 255, .065) !important;
  box-shadow: none !important;
}

.login-control :deep(.el-input__wrapper:hover),
.login-control :deep(.el-select__wrapper:hover) {
  border-color: rgba(var(--accent-rgb), .55) !important;
}

.login-control :deep(.el-input__wrapper.is-focus),
.login-control :deep(.el-select__wrapper.is-focused) {
  border-color: var(--accent) !important;
  box-shadow: 0 0 0 3px rgba(var(--accent-rgb), .14) !important;
}

.login-control :deep(.el-input__inner),
.login-control :deep(.el-select__selected-item) {
  color: #effff4 !important;
}

.login-control :deep(.el-input__inner::placeholder),
.login-control :deep(.el-select__placeholder) {
  color: rgba(220, 255, 230, .38) !important;
}

.branch-control :deep(.el-select__wrapper) {
  background: rgba(var(--accent-rgb), .11) !important;
}

.input-icon,
.login-control :deep(.el-input__prefix) {
  width: 17px;
  height: 17px;
  color: rgba(var(--accent-rgb), .75);
}

.password-toggle {
  display: grid;
  width: 28px;
  height: 28px;
  padding: 0;
  border: 0;
  place-items: center;
  background: transparent;
  color: rgba(var(--accent-rgb), .72);
  cursor: pointer;
}

.password-toggle svg {
  width: 17px;
  height: 17px;
}

.login-error {
  display: flex;
  padding: 9px 11px;
  border: 1px solid rgba(248, 113, 113, .34);
  border-radius: 10px;
  align-items: flex-start;
  gap: 7px;
  background: rgba(185, 28, 28, .2);
  color: #fecaca;
  font-size: 12px;
  line-height: 1.45;
}

.login-error svg {
  width: 16px;
  min-width: 16px;
  margin-top: 1px;
}

.submit-button {
  display: flex;
  width: 100%;
  height: 48px;
  margin-top: 2px;
  border: 0;
  border-radius: 14px;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: linear-gradient(135deg, rgba(var(--accent-rgb), .72), var(--accent));
  box-shadow: 0 7px 24px rgba(var(--accent-rgb), .24);
  color: #072018;
  font-size: 15px;
  font-weight: 900;
  cursor: pointer;
}

.submit-button:hover:not(:disabled) {
  filter: brightness(1.06);
  transform: translateY(-1px);
}

.submit-button:disabled {
  opacity: .48;
  cursor: not-allowed;
}

.spinner {
  width: 17px;
  height: 17px;
  border: 2.5px solid rgba(7, 32, 24, .3);
  border-top-color: #072018;
  border-radius: 50%;
  animation: spin .7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 520px) {
  .login-page { padding: 70px 14px 18px; }
  .login-card { padding: 26px 20px 23px; border-radius: 20px; }
  .back-button { top: 14px; left: 14px; }
  h1 { font-size: 23px; }
}
</style>
