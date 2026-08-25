<template>
  <main class="login-page">
    <!-- ═══ BÊN TRÁI: ảnh + nội dung ═══ -->
    <aside class="split split-left">
      <div class="left-bg" aria-hidden="true" />
      <div class="left-grid" aria-hidden="true" />
      <div class="frame-corners" aria-hidden="true" />
      <div class="left-watermark" aria-hidden="true">
        <span class="crate crate-1" />
        <span class="crate crate-2" />
        <span class="crate crate-3" />
      </div>

      <div class="left-content">
        <div class="brand anim" style="--i: 1">
          <span class="brand-mark">{{ isLossLogin ? "LT" : "N8" }}</span>
          <span class="brand-text">{{ copy.brandTitle }}</span>
        </div>

        <div class="left-body">
          <p class="eyebrow anim" style="--i: 2">
            <i class="pulse-dot" aria-hidden="true" />{{ copy.eyebrow }}
          </p>

          <h1 class="anim" style="--i: 3">
            {{ copy.headlinePre }}<br />
            <span :key="rotatingIndex" class="rotator">{{ currentWord }}</span>
          </h1>

          <span class="left-divider anim" style="--i: 4" aria-hidden="true" />

          <ul class="feature-list">
            <li class="anim" style="--i: 5">
              <span class="f-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <path d="M12 3 5 6v5c0 4.5 2.8 8 7 10 4.2-2 7-5.5 7-10V6l-7-3Z" />
                </svg>
              </span>
              <span class="f-text">{{ copy.featureScan }}</span>
            </li>
            <li class="anim" style="--i: 6">
              <span class="f-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <path d="M4 20V10m6 10V4m6 16v-7" />
                </svg>
              </span>
              <span class="f-text">{{ copy.featureTrack }}</span>
            </li>
            <li class="anim" style="--i: 7">
              <span class="f-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <rect x="5" y="10" width="14" height="10" rx="2" />
                  <path d="M8 10V7a4 4 0 0 1 8 0v3" />
                </svg>
              </span>
              <span class="f-text">{{ copy.featureSession }}</span>
            </li>
          </ul>
        </div>
      </div>
    </aside>

    <!-- ═══ BÊN PHẢI: form login ═══ -->
    <section class="split split-right">
      <div class="right-inner">
        <header class="panel-head">
          <div>
            <p class="panel-eyebrow">{{ copy.panelEyebrow }}</p>
            <h2>{{ copy.login }}</h2>
          </div>

          <el-select
            :model-value="locale"
            class="language-control"
            :aria-label="copy.language"
            :teleported="false"
            popper-class="language-control-popper"
            @change="changeLocale"
          >
            <template #prefix>
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <circle cx="12" cy="12" r="9" />
                <path d="M3 12h18M12 3a15 15 0 0 1 0 18M12 3a15 15 0 0 0 0 18" />
              </svg>
            </template>
            <el-option
              v-for="item in locales"
              :key="item.code"
              :label="item.label"
              :value="item.code"
            />
          </el-select>
        </header>

        <div v-if="errorMessage" class="error-message" role="alert">
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <circle cx="12" cy="12" r="9" />
            <path d="M12 8v5m0 3h.01" />
          </svg>
          <span>{{ errorMessage }}</span>
        </div>

        <form
          class="login-form"
          :class="{ shake: errorShake }"
          @submit.prevent="submit"
        >
          <label>
            <span>{{ copy.branch }}</span>
            <span class="input-shell">
              <span class="input-icon">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <path d="M4 20V8l8-4 8 4v12M8 20v-4h8v4M8 10h.01M12 10h.01M16 10h.01" />
                </svg>
              </span>
              <input
                v-model.trim="branch"
                autocomplete="organization"
                :placeholder="copy.branchPlaceholder"
                :disabled="loading || phase !== 'idle'"
                autocapitalize="characters"
              />
              <svg class="chevron" viewBox="0 0 24 24" aria-hidden="true">
                <path d="m6 9 6 6 6-6" />
              </svg>
            </span>
          </label>

          <label>
            <span>{{ copy.userId }}</span>
            <span class="input-shell">
              <span class="input-icon">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <circle cx="12" cy="8" r="4" />
                  <path d="M4.5 20c.8-4 3.3-6 7.5-6s6.7 2 7.5 6" />
                </svg>
              </span>
              <input
                ref="userInput"
                v-model.trim="userId"
                autocomplete="username"
                :placeholder="copy.userPlaceholder"
                :disabled="loading || phase !== 'idle'"
              />
            </span>
          </label>

          <label>
            <span>{{ copy.password }}</span>
            <span class="input-shell">
              <span class="input-icon">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <rect x="5" y="10" width="14" height="10" rx="2" />
                  <path d="M8 10V7a4 4 0 0 1 8 0v3M12 14v2" />
                </svg>
              </span>
              <input
                v-model="password"
                autocomplete="current-password"
                :type="showPassword ? 'text' : 'password'"
                :placeholder="copy.passwordPlaceholder"
                :disabled="loading || phase !== 'idle'"
              />
              <button
                class="password-toggle"
                type="button"
                :aria-label="showPassword ? copy.hidePassword : copy.showPassword"
                :disabled="loading || phase !== 'idle'"
                @click="showPassword = !showPassword"
              >
                <svg v-if="!showPassword" viewBox="0 0 24 24" aria-hidden="true">
                  <path d="M2 12s3.6-7 10-7 10 7 10 7-3.6 7-10 7-10-7-10-7Z" />
                  <circle cx="12" cy="12" r="3" />
                </svg>
                <svg v-else viewBox="0 0 24 24" aria-hidden="true">
                  <path d="M3 3l18 18M10.6 10.6a3 3 0 0 0 4.2 4.2M6.6 6.6C4.3 8 2.7 10 2 12c0 0 3.6 7 10 7 1.9 0 3.5-.5 4.9-1.3M9.9 4.2A10.4 10.4 0 0 1 12 4c6.4 0 10 7 10 7-.5 1-1.2 2.1-2.1 3.1" />
                </svg>
              </button>
            </span>
          </label>

          <button
            class="submit-button"
            :class="{ ok: phase !== 'idle' }"
            type="submit"
            :disabled="!canSubmit || loading || phase !== 'idle'"
          >
            <template v-if="phase === 'success'">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M4 12.5l5 5L20 6.5" />
              </svg>
              {{ copy.loginOk }}
            </template>
            <template v-else-if="loading">
              <span class="spinner" aria-hidden="true" />
              {{ copy.loggingIn }}
            </template>
            <template v-else>
              {{ copy.login }}
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M5 12h14M13 6l6 6-6 6" />
              </svg>
            </template>
          </button>
        </form>
      </div>
    </section>
  </main>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import { useLogin } from "@/hooks/useLogin";

const router = useRouter();
const route = useRoute();
const { locale } = useI18n();

const branch = ref(localStorage.getItem("GSBH") || "VDH");
const userId = ref(localStorage.getItem("USERID") || "");
const password = ref("");
const showPassword = ref(false);
const loading = ref(false);
const errorMessage = ref("");
const errorShake = ref(false);
const userInput = ref(null);

/* ── Hiệu ứng mở màn: idle → success → chuyển route ── */
const phase = ref("idle");
const SUCCESS_MS = 600; // thời gian nút hiển thị "Thành công"

/* ── Từ khóa xoay trên tiêu đề ── */
const rotatingIndex = ref(0);
let rotateTimer = null;

const locales = [
  { code: "vi", label: "Tiếng Việt" },
  { code: "en", label: "English" },
  { code: "zh", label: "中文" },
];

const messages = {
  vi: {
    brandTitle: "HỆ THỐNG QUẢN LÝ KHO VẬT TƯ",
    eyebrow: "N8D · MATERIAL WAREHOUSE",
    headlinePre: "Kiểm soát",
    rotatingWords: ["kho vận", "vật tư", "barcode"],
    lead: "Nền tảng quản lý xuất – nhập – tồn theo thời gian thực, gắn liền với luồng vận hành của nhà máy.",
    panelEyebrow: "SECURE ACCESS",
    panelSub: "Đăng nhập bằng tài khoản nhân viên để tiếp tục.",
    language: "Ngôn ngữ",
    featureScan: "Kiểm kê barcode nhanh",
    featureTrack: "Theo dõi dữ liệu tập trung",
    featureSession: "Ghi nhận theo tài khoản",
    login: "Đăng nhập",
    loginOk: "Thành công",
    branch: "Chi nhánh (CSBH)",
    branchPlaceholder: "Ví dụ: VDH",
    userId: "Mã nhân viên",
    userPlaceholder: "Nhập mã nhân viên",
    password: "Mật khẩu",
    passwordPlaceholder: "Nhập mật khẩu",
    showPassword: "Hiện mật khẩu",
    hidePassword: "Ẩn mật khẩu",
    loggingIn: "Đang xác thực…",
    failed: "Đăng nhập không thành công.",
    securityNote: "Phiên làm việc được ghi nhận theo tài khoản và chi nhánh.",
  },
  en: {
    brandTitle: "MATERIAL WAREHOUSE MANAGEMENT",
    eyebrow: "N8D · MATERIAL WAREHOUSE",
    headlinePre: "Total control of",
    rotatingWords: ["inventory", "materials", "barcode"],
    lead: "Real-time inbound – outbound – stock management, wired directly into plant operations.",
    panelEyebrow: "SECURE ACCESS",
    panelSub: "Sign in with your employee account to continue.",
    language: "Language",
    featureScan: "Fast barcode inventory",
    featureTrack: "Centralized tracking",
    featureSession: "Account-based records",
    login: "Sign in",
    loginOk: "Success",
    branch: "Branch (CSBH)",
    branchPlaceholder: "Example: VDH",
    userId: "Employee ID",
    userPlaceholder: "Enter employee ID",
    password: "Password",
    passwordPlaceholder: "Enter password",
    showPassword: "Show password",
    hidePassword: "Hide password",
    loggingIn: "Signing in…",
    failed: "Sign-in failed.",
    securityNote: "Sessions are recorded per account and branch.",
  },
  zh: {
    brandTitle: "物资仓库管理系统",
    eyebrow: "N8D · MATERIAL WAREHOUSE",
    headlinePre: "全面掌控",
    rotatingWords: ["库存", "物资", "条码"],
    lead: "实时管理入库、出库与库存，与工厂运营流程紧密衔接。",
    panelEyebrow: "安全登录",
    panelSub: "请使用员工账号登录以继续。",
    language: "语言",
    featureScan: "快速条码盘点",
    featureTrack: "集中追踪数据",
    featureSession: "按账户记录",
    login: "登录",
    loginOk: "登录成功",
    branch: "分厂 (CSBH)",
    branchPlaceholder: "例如：VDH",
    userId: "员工编号",
    userPlaceholder: "请输入员工编号",
    password: "密码",
    passwordPlaceholder: "请输入密码",
    showPassword: "Show password" /* fallback */,
    hidePassword: "隐藏密码",
    showPasswordLabel: "显示密码",
    loggingIn: "正在验证…",
    failed: "登录失败。",
    securityNote: "会话按账户与分厂记录。",
  },
};

const isLossLogin = computed(() => route.path === "/login");

const copy = computed(() => {
  const baseCopy = messages[locale.value] || messages.vi;
  if (!isLossLogin.value) return baseCopy;

  const lossCopy = {
    vi: {
      brandTitle: "QUẢN LÝ HAO HỤT TEM",
      eyebrow: "LOSS TEM SIZE · QUẢN LÝ HAO HỤT TEM",
      headlinePre: "Ghi nhận",
      rotatingWords: ["hao hụt", "tem size", "theo ca"],
      lead: "Ghi nhận và đối chiếu hao hụt tem size theo chi nhánh, theo ca và theo tài khoản.",
      featureScan: "Ghi nhận nhanh theo chi nhánh",
      featureTrack: "Theo dõi dữ liệu tập trung",
      featureSession: "Bảo mật theo tài khoản",
    },
    en: {
      brandTitle: "LABEL LOSS MANAGEMENT",
      eyebrow: "LOSS TEM SIZE · LABEL LOSS MANAGEMENT",
      headlinePre: "Track every",
      rotatingWords: ["loss", "label", "shift"],
      lead: "Record and reconcile label-size loss by branch, shift and account.",
      featureScan: "Fast branch-based records",
      featureTrack: "Centralized tracking",
      featureSession: "Account-based security",
    },
    zh: {
      brandTitle: "标签损耗管理",
      eyebrow: "LOSS TEM SIZE · 标签损耗管理",
      headlinePre: "精准记录",
      rotatingWords: ["损耗", "标签", "班组"],
      lead: "按分厂、班组和账户记录并核对标签损耗。",
      featureScan: "按分厂快速记录",
      featureTrack: "集中追踪数据",
      featureSession: "账户安全管理",
    },
  };

  return { ...baseCopy, ...lossCopy[locale.value] };
});

const canSubmit = computed(() =>
  Boolean(userId.value.trim() && password.value),
);

const rotatingWords = computed(() => copy.value.rotatingWords || []);
const currentWord = computed(
  () =>
    rotatingWords.value[rotatingIndex.value % rotatingWords.value.length] || "",
);

const changeLocale = (nextLocale) => {
  if (!locales.some((item) => item.code === nextLocale)) return;
  locale.value = nextLocale;
  localStorage.setItem("APP_LOCALE", nextLocale);
};

const wait = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

const triggerShake = () => {
  errorShake.value = false;
  requestAnimationFrame(() => {
    errorShake.value = true;
    setTimeout(() => (errorShake.value = false), 500);
  });
};

const submit = async () => {
  if (!canSubmit.value || loading.value || phase.value !== "idle") return;
  loading.value = true;
  errorMessage.value = "";

  try {
    const normalizedBranch = branch.value.trim().toUpperCase();
    const normalizedUser = userId.value.trim();
    const result = await useLogin({
      GSBH: normalizedBranch,
      USERID: normalizedUser,
      PWD: password.value,
    });

    if (!result.success) {
      errorMessage.value = result.message || copy.value.failed;
      triggerShake();
      return;
    }

    localStorage.setItem("GSBH", normalizedBranch);
    localStorage.setItem("USERID", normalizedUser);

    /* 1) Nút báo thành công */
    phase.value = "success";
    await wait(SUCCESS_MS);

    /* 2) Chuyển sang route chính (trang MainPage thật) */
    const redirect = route.query.redirect;
    const fallback = isLossLogin.value ? "/loss-tem-size" : "/n8d";
    const destination =
      typeof redirect === "string" && redirect.startsWith("/")
        ? redirect
        : fallback;
    await router.replace(destination);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  rotateTimer = setInterval(() => {
    rotatingIndex.value += 1;
  }, 2600);

  const savedLocale = localStorage.getItem("APP_LOCALE");
  if (locales.some((item) => item.code === savedLocale))
    locale.value = savedLocale;
  if (!userId.value) userInput.value?.focus();
});

onBeforeUnmount(() => {
  clearInterval(rotateTimer);
});
</script>

<style scoped>
/* ═══════════ Nền tảng ═══════════ */
.login-page {
  --ink: #0c1e33;
  --ink-2: #143352;
  --paper: #eef3f9;
  --paper-2: #ffffff;
  --line: rgb(12 30 51 / 12%);
  --muted: #52697f;
  --laser: #2f8fff;
  --sand: #bcd2e8;

  position: relative;
  display: grid;
  grid-template-columns: 1.18fr 1fr;
  min-height: 100dvh;
  overflow: hidden;
  background: var(--paper);
  color: var(--ink);
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC",
    "Microsoft YaHei", sans-serif;
}
button,
input,
select {
  font: inherit;
}

/* ═══════════ 2 panel ═══════════ */
.split {
  position: relative;
  z-index: 1;
}

/* ═══════════ Bên trái: ảnh + nội dung ═══════════ */
.split-left {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  color: #eaf2fb;
}
/* Nền màu: gradient navy nhiều lớp thay cho ảnh */
.left-bg {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(58% 50% at 88% 8%, rgb(47 143 255 / 0.34), transparent 68%),
    radial-gradient(46% 42% at 8% 96%, rgb(79 216 255 / 0.22), transparent 70%),
    linear-gradient(155deg, #0a1a30 0%, #0d2745 48%, #081525 100%);
  background-size: 140% 140%, 140% 140%, 100% 100%;
  background-position: 0% 0%, 100% 100%, center;
  animation: aurora 22s ease-in-out infinite alternate;
}
@keyframes aurora {
  from {
    background-position: 0% 0%, 100% 100%, center;
  }
  to {
    background-position: 12% 10%, 88% 88%, center;
  }
}
.left-grid {
  position: absolute;
  inset: 0;
  background-image: radial-gradient(rgb(234 242 251 / 0.14) 1px, transparent 1px);
  background-size: 30px 30px;
  mask-image: radial-gradient(78% 78% at 32% 34%, #000 46%, transparent 92%);
}

/* Watermark: các khối thùng hàng xếp lệch — gợi liên tưởng kho vận */
.left-watermark {
  position: absolute;
  right: clamp(-40px, -2vw, 10px);
  bottom: clamp(-30px, -2vw, 20px);
  z-index: 1;
  width: 260px;
  height: 220px;
  pointer-events: none;
  opacity: 0.5;
}
.crate {
  position: absolute;
  border: 1.5px solid rgb(125 187 255 / 0.4);
  border-radius: 10px;
  background: rgb(47 143 255 / 0.05);
}
.crate-1 {
  right: 10px;
  bottom: 10px;
  width: 128px;
  height: 92px;
}
.crate-2 {
  right: 92px;
  bottom: 70px;
  width: 92px;
  height: 68px;
  border-color: rgb(125 187 255 / 0.3);
}
.crate-3 {
  right: 30px;
  bottom: 118px;
  width: 66px;
  height: 50px;
  border-color: rgb(125 187 255 / 0.5);
  background: rgb(47 143 255 / 0.08);
}

/* 4 góc khung viewfinder — kiểu máy quét kho */
.frame-corners {
  position: absolute;
  inset: 22px;
  z-index: 1;
  pointer-events: none;
  opacity: 0.55;
  background-image:
    linear-gradient(#4fa8ff, #4fa8ff), linear-gradient(#4fa8ff, #4fa8ff),
    linear-gradient(#4fa8ff, #4fa8ff), linear-gradient(#4fa8ff, #4fa8ff),
    linear-gradient(#4fa8ff, #4fa8ff), linear-gradient(#4fa8ff, #4fa8ff),
    linear-gradient(#4fa8ff, #4fa8ff), linear-gradient(#4fa8ff, #4fa8ff);
  background-repeat: no-repeat;
  background-size:
    26px 2px, 2px 26px, 26px 2px, 2px 26px,
    26px 2px, 2px 26px, 26px 2px, 2px 26px;
  background-position:
    top left, top left,
    top right, top right,
    bottom left, bottom left,
    bottom right, bottom right;
}

.left-content {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 24px;
  min-height: 100%;
  padding: clamp(22px, 3.2vw, 44px);
}
.left-body {
  max-width: 560px;
  margin-block: auto;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 12px;
}
.brand-mark {
  display: grid;
  width: 42px;
  height: 42px;
  place-items: center;
  border: 1px solid rgb(230 240 250 / 35%);
  border-radius: 12px;
  background: rgb(230 240 250 / 8%);
  color: #eaf2fb;
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.04em;
  backdrop-filter: blur(6px);
  -webkit-backdrop-filter: blur(6px);
}
.brand-text {
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.14em;
}

.eyebrow {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0 0 18px;
  color: var(--sand);
  font: 600 11.5px ui-monospace, "SFMono-Regular", Menlo, Consolas, monospace;
  letter-spacing: 0.22em;
}
.pulse-dot {
  width: 7px;
  height: 7px;
  flex: 0 0 auto;
  border-radius: 50%;
  background: #2f8fff;
  box-shadow: 0 0 0 0 rgb(47 143 255 / 60%);
  animation: pulse 2.2s infinite;
}
@keyframes pulse {
  70% {
    box-shadow: 0 0 0 9px rgb(47 143 255 / 0%);
  }
  100% {
    box-shadow: 0 0 0 0 rgb(47 143 255 / 0%);
  }
}

h1 {
  margin: 0 0 18px;
  color: #eef5fc;
  font-size: clamp(38px, 4.6vw, 64px);
  font-weight: 700;
  letter-spacing: -0.02em;
  line-height: 1.04;
}
.rotator {
  display: inline-block;
  color: #4fa8ff;
  animation: wordIn 0.6s cubic-bezier(0.2, 0.8, 0.2, 1) both;
}
@keyframes wordIn {
  from {
    opacity: 0;
    transform: translateY(45%);
  }
  to {
    opacity: 1;
    transform: none;
  }
}
.left-divider {
  display: block;
  width: 46px;
  height: 3px;
  margin: 2px 0 26px;
  border-radius: 2px;
  background: linear-gradient(90deg, #4fa8ff, rgb(79 168 255 / 0));
}
.feature-list {
  display: grid;
  max-width: 460px;
  gap: 6px;
  margin: 0;
  padding: 0;
  list-style: none;
}
.feature-list li {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 13px 16px;
  border: 1px solid rgb(230 240 250 / 12%);
  border-left: 2px solid transparent;
  border-radius: 10px;
  background: rgb(8 20 36 / 32%);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition:
    border-left-color 0.25s,
    background 0.25s,
    transform 0.25s;
}
.feature-list li:hover {
  border-left-color: var(--laser);
  background: rgb(8 20 36 / 55%);
  transform: translateX(6px);
}
.f-icon {
  display: grid;
  width: 34px;
  height: 34px;
  flex: 0 0 auto;
  place-items: center;
  border: 1px solid rgb(230 240 250 / 25%);
  border-radius: 50%;
  background: rgb(230 240 250 / 7%);
  color: #eaf2fb;
}
.f-icon svg {
  width: 15px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.7;
  stroke-linecap: round;
  stroke-linejoin: round;
}
.f-text {
  color: #e7f0f9;
  font-size: 14px;
  font-weight: 600;
}

/* Vào sân khấu: chữ hiện lần lượt */
.anim {
  animation: rise 0.8s cubic-bezier(0.2, 0.75, 0.25, 1) both;
  animation-delay: calc(140ms + var(--i) * 110ms);
}
@keyframes rise {
  from {
    opacity: 0;
    transform: translateY(26px);
  }
  to {
    opacity: 1;
    transform: none;
  }
}

/* ═══════════ Bên phải: form ═══════════ */
.split-right {
  display: grid;
  place-items: center;
  padding: clamp(26px, 4vw, 56px) clamp(24px, 4.5vw, 72px);
  background-color: var(--paper);
  background-image:
    linear-gradient(rgb(12 30 51 / 4.5%) 1px, transparent 1px),
    linear-gradient(90deg, rgb(12 30 51 / 4.5%) 1px, transparent 1px);
  background-size: 42px 42px;
}
.right-inner {
  width: min(440px, 100%);
}

.panel-head {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14px 16px;
  margin-bottom: 6px;
}
.panel-eyebrow {
  margin: 0 0 10px;
  color: var(--laser);
  font: 600 11px ui-monospace, "SFMono-Regular", Menlo, Consolas, monospace;
  letter-spacing: 0.24em;
}
.panel-head h2 {
  margin: 0;
  font-size: clamp(30px, 2.6vw, 40px);
  font-weight: 700;
  letter-spacing: -0.02em;
}
.language-control {
  width: 156px;
}
.language-control :deep(svg) {
  width: 16px;
  fill: none;
  stroke: var(--ink);
  stroke-width: 1.7;
  stroke-linecap: round;
}
.language-control :deep(.el-select__wrapper),
.language-control :deep(.el-input__wrapper) {
  min-height: 42px;
  padding: 0 13px;
  border-radius: 10px;
  background: var(--paper-2);
  box-shadow: 0 0 0 1px var(--line) !important;
}
.language-control :deep(.el-select__wrapper.is-focused),
.language-control :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px var(--ink) !important;
}
.language-control :deep(.el-select__selected-item),
.language-control :deep(.el-input__inner) {
  color: var(--ink);
  font: 600 12.5px ui-monospace, "SFMono-Regular", Menlo, Consolas, monospace;
}
.language-control :deep(.el-select__caret) {
  color: var(--muted);
}

.error-message {
  display: flex;
  margin-top: 22px;
  align-items: flex-start;
  gap: 10px;
  padding: 13px 15px;
  border: 1px solid rgb(204 55 40 / 25%);
  border-radius: 12px;
  background: #fbeae6;
  color: #a3271b;
  font-size: 13px;
  line-height: 1.5;
}
.error-message svg {
  width: 19px;
  flex: 0 0 auto;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
}

.login-form {
  display: grid;
  gap: 18px;
  margin-top: 26px;
}
.login-form > label {
  display: grid;
  gap: 8px;
}
.login-form > label > span:first-child {
  color: #45607a;
  font: 600 11px ui-monospace, "SFMono-Regular", Menlo, Consolas, monospace;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}
.input-shell {
  display: flex;
  min-height: 54px;
  align-items: center;
  gap: 12px;
  padding: 0 14px;
  border-radius: 12px;
  background: var(--paper-2);
  box-shadow: 0 0 0 1px var(--line);
  transition:
    box-shadow 0.2s,
    transform 0.2s;
}
.input-shell:focus-within {
  box-shadow: 0 0 0 2px var(--ink);
  transform: translateY(-1px);
}
.input-icon {
  display: grid;
  width: 30px;
  height: 30px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 8px;
  background: rgb(12 30 51 / 6%);
  color: #4a6a85;
  transition:
    color 0.2s,
    background 0.2s;
}
.input-shell:focus-within .input-icon {
  background: var(--ink);
  color: var(--paper);
}
.input-icon svg {
  width: 16px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}
.input-shell input {
  min-width: 0;
  flex: 1;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--ink);
  font-size: 14.5px;
}
.input-shell input::placeholder {
  color: #9aa8bd;
}
.input-shell input:disabled {
  opacity: 0.55;
}
.input-shell > svg.chevron {
  width: 16px;
  flex: 0 0 auto;
  fill: none;
  stroke: #9aa8bd;
  stroke-width: 2;
  stroke-linecap: round;
}
.password-toggle {
  display: grid;
  padding: 4px;
  place-items: center;
  border: 0;
  background: transparent;
  color: #9aa8bd;
  cursor: pointer;
}
.password-toggle svg {
  width: 19px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.7;
  stroke-linecap: round;
  stroke-linejoin: round;
}
.password-toggle:hover {
  color: var(--ink);
}

.submit-button {
  position: relative;
  display: flex;
  min-height: 56px;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 8px;
  overflow: hidden;
  border: 0;
  border-radius: 12px;
  background: var(--ink);
  color: var(--paper);
  font-size: 14.5px;
  font-weight: 700;
  letter-spacing: 0.02em;
  cursor: pointer;
  transition:
    transform 0.18s ease,
    background 0.25s;
}
.submit-button::after {
  position: absolute;
  right: 0;
  bottom: 0;
  left: 0;
  height: 2px;
  background: var(--laser);
  transform: scaleX(0);
  transform-origin: left;
  transition: transform 0.3s;
  content: "";
}
.submit-button:not(:disabled):hover {
  background: var(--ink-2);
  transform: translateY(-2px);
}
.submit-button:not(:disabled):hover::after {
  transform: scaleX(1);
}
.submit-button:not(:disabled):active {
  transform: translateY(0) scale(0.99);
}
.submit-button svg {
  width: 18px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
  transition: transform 0.2s;
}
.submit-button:not(.ok):not(:disabled):hover svg {
  transform: translateX(4px);
}
.submit-button.ok {
  background: #1b6fc4;
  opacity: 1;
}
.submit-button:disabled:not(.ok) {
  opacity: 0.45;
  cursor: not-allowed;
}
.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgb(230 240 250 / 40%);
  border-top-color: var(--paper);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
.shake {
  animation: shake 0.45s;
}
@keyframes shake {
  10%,
  90% {
    transform: translateX(-2px);
  }
  20%,
  80% {
    transform: translateX(3px);
  }
  30%,
  50%,
  70% {
    transform: translateX(-5px);
  }
  40%,
  60% {
    transform: translateX(5px);
  }
}

/* ═══════════ Responsive ═══════════ */
@media (max-width: 920px) {
  .login-page {
    grid-template-columns: 1fr;
  }
  .split-left {
    min-height: 42dvh;
  }
  .feature-list {
    display: none;
  }
  .left-watermark {
    display: none;
  }
  .split-right {
    place-items: start center;
    padding-top: 40px;
  }
}
@media (max-width: 520px) {
  .split-left {
    min-height: 34dvh;
  }
  .brand-text {
    font-size: 11px;
  }
  .right-inner {
    width: 100%;
  }
  .panel-head h2 {
    font-size: 28px;
  }
  .language-control {
    width: 100%;
  }
}
@media (prefers-reduced-motion: reduce) {
  .split {
    transition: none;
  }
  .left-bg,
  .pulse-dot,
  .rotator,
  .anim {
    animation: none;
  }
}
</style>