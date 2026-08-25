<template>
  <div class="n8d-home">
    <!-- Hiệu ứng nền gradient nhẹ nhàng cho giao diện sáng -->
    <div class="bg-mesh" aria-hidden="true">
      <div class="glow glow--tracking"></div>
      <div class="glow glow--month"></div>
      <div class="glow glow--day"></div>
    </div>

    <div
      v-if="syncingMonth"
      class="sync-blocking-overlay"
      role="dialog"
      aria-modal="true"
      :aria-label="copy.syncDialogTitle"
    >
      <div class="sync-blocking-dialog">
        <span class="sync-dialog-spinner" aria-hidden="true" />
        <strong>{{ copy.syncDialogTitle }}</strong>
        <span>{{ copy.syncDialogMessage }}</span>
        <span class="sync-dialog-period"
          >{{ copy.syncDialogPeriod }} {{ getLocalMonth() }}</span
        >
      </div>
    </div>

    <N8dHeader
      show-login
      :is-logged-in="isLoggedIn"
      :user-label="userLabel"
      @login="goToLogin"
      @logout="confirmLogout"
    >
      <main class="home-content">
        <div class="hero">
          <span class="hero-eyebrow">{{ copy.brandSubtitle }}</span>
          <h1>{{ copy.heroTitle }}</h1>
        </div>

        <section class="module-bento">
          <!-- Feature Card: Kiểm Kê Tháng (nhiều chức năng con, cần không gian lớn) -->
          <div class="module-card module-card--feature module-card--month">
            <div class="feature-top">
              <span class="card-icon card-icon--lg" aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                  <path
                    d="M16 2v4M8 2v4M3 10h18M8 14h.01M12 14h.01M16 14h.01M8 18h.01M12 18h.01M16 18h.01"
                  />
                </svg>
              </span>
            </div>

            <span class="ticket-divider" aria-hidden="true"></span>

            <div class="feature-body">
              <h3>{{ copy.inventoryMonth }}</h3>

              <div class="month-sync-row">
                <button
                  type="button"
                  class="month-action-btn month-action-btn--sync"
                  :disabled="syncingMonth"
                  @click="syncInventoryMonth"
                >
                  <svg v-if="!syncingMonth" viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M21 12a9 9 0 1 1-2.64-6.36" />
                    <path d="M21 3v6h-6" />
                  </svg>
                  <span v-else class="spinner" />
                  <span>{{ syncingMonth ? copy.syncing : copy.syncAction }}</span>
                </button>
                <span class="last-sync-badge">{{ lastSyncText }}</span>
              </div>

              <div class="month-actions">
                <button
                  type="button"
                  class="month-action-btn month-action-btn--unit"
                  @click="selectMonthRole('unit')"
                >
                  <span>{{ copy.monthRoleUnit }}</span>
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M5 12h14m-7-7 7 7-7 7" />
                  </svg>
                </button>
                <button
                  type="button"
                  class="month-action-btn month-action-btn--accounting"
                  @click="selectMonthRole('accounting')"
                >
                  <span>{{ copy.monthRoleAccounting }}</span>
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M5 12h14m-7-7 7 7-7 7" />
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <!-- Stacked Cards -->
          <div class="bento-stack">
            <!-- Tracking Card -->
            <button
              class="module-card module-card--row module-card--tracking"
              type="button"
              @click="openModule(modules[0])"
            >
              <span class="card-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <path
                    d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
                  />
                  <path d="M14 2v6h6M16 13H8M16 17H8M10 9H8" />
                </svg>
              </span>

              <div class="row-body">
                <h3>{{ modules[0].title }}</h3>
                <p class="card-desc card-desc--row">
                  {{ modules[0].description }}
                </p>
              </div>

              <span class="card-arrow" aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <path d="M5 12h14m-7-7 7 7-7 7" />
                </svg>
              </span>
            </button>

            <!-- Daily Inventory Card -->
            <button
              class="module-card module-card--row module-card--day"
              type="button"
              @click="openModule(dayModule)"
            >
              <span class="card-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <path
                    d="M12 2v4m0 12v4M2 12h4m12 0h4M4.93 4.93l2.83 2.83m8.48 8.48l2.83 2.83M19.07 4.93l-2.83 2.83m-8.48 8.48l-2.83 2.83"
                  />
                </svg>
              </span>

              <div class="row-body">
                <h3>{{ copy.inventoryDay }}</h3>
              </div>

              <span class="card-arrow" aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <path d="M5 12h14m-7-7 7 7-7 7" />
                </svg>
              </span>
            </button>
          </div>
        </section>
      </main>

      <Transition name="toast">
        <div v-if="snackbar" class="auth-toast" role="status">
          <span>{{ copy.notLoggedIn }}</span>
          <button type="button" @click="goToLogin">{{ copy.login }}</button>
        </div>
      </Transition>

      <Transition name="toast">
        <div
          v-if="syncToast.show"
          class="sync-toast"
          :class="syncToast.color"
          role="status"
        >
          <svg
            v-if="syncToast.color === 'success'"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <circle cx="12" cy="12" r="9" />
            <path d="m8 12 2.5 2.5L16 9" />
          </svg>
          <svg v-else viewBox="0 0 24 24" aria-hidden="true">
            <path d="M12 3 2.8 20h18.4L12 3Z" />
            <path d="M12 9v5m0 3h.01" />
          </svg>
          <span>{{ syncToast.message }}</span>
          <button type="button" @click="syncToast.show = false">×</button>
        </div>
      </Transition>
    </N8dHeader>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import axios from "axios";
import N8dHeader from "@/views/components/N8d/N8dHeader.vue";
import { clearSession, getToken } from "@/utils/appStoreAndSecurity";

const router = useRouter();
const { locale } = useI18n();

const isLoggedIn = ref(false);
const userId = ref("");
const branch = ref("");
const sn8 = ref("");
const snackbar = ref(false);
const syncingMonth = ref(false);
const lastSync = ref({ SyncDate: "", SyncUserID: "" });
const syncToast = ref({ show: false, message: "", color: "success" });
let snackbarTimer;
let syncToastTimer;
let authSyncTimer;

const messages = {
  vi: {
    brandSubtitle: "Material Management",
    heroTitle: "Chọn module để bắt đầu",
    tracking: "TRACKING",
    trackingDesc:
      "Luồng vật tư theo thời gian thực, cập nhật liên tục theo từng lô hàng.",
    inventoryMonth: "Kiểm Kê Tháng",
    monthRoleUnit: "Kiểm kê đơn vị",
    monthRoleAccounting: "Kiểm kê kế toán",
    inventoryDay: "Kiểm Kê Ngày",
    inventoryDayDesc: "Chốt số liệu mỗi ngày",
    login: "Đăng Nhập",
    logout: "Đăng Xuất",
    notLoggedIn: "Vui lòng đăng nhập để tiếp tục.",
    logoutConfirm: "Xác nhận đăng xuất khỏi N8D?",
    ready: "Sẵn sàng",
    syncAction: "Đồng bộ",
    syncing: "Đang đồng bộ...",
    noPermission: "Bạn chưa có quyền SN8 để đồng bộ.",
    syncDialogTitle: "Đang đồng bộ",
    syncDialogMessage: "Vui lòng chờ đồng bộ xong.",
    syncDialogPeriod: "Tháng đồng bộ:",
    syncSuccess: "Đã đồng bộ {count} mục.",
    syncWarning: "Đồng bộ hoàn tất nhưng có cảnh báo.",
    syncFail: "Đồng bộ thất bại, vui lòng thử lại.",
    sessionExpired: "Phiên đăng nhập đã hết hạn, vui lòng đăng nhập lại.",
    requestFailed: "Có lỗi xảy ra, vui lòng thử lại.",
    lastSyncLabel: "Đồng bộ lần cuối: {date}",
    lastSyncNone: "Chưa đồng bộ",
  },
  en: {
    brandSubtitle: "Material Management",
    heroTitle: "Pick a module to get started",
    tracking: "TRACKING",
    trackingDesc: "Real-time material flow, updated continuously per shipment.",
    inventoryMonth: "Monthly Inventory",
    monthRoleUnit: "Unit check",
    monthRoleAccounting: "Accounting check",
    inventoryDay: "Daily Inventory",
    inventoryDayDesc: "Close daily figures",
    login: "Sign In",
    logout: "Sign Out",
    notLoggedIn: "Please sign in to continue.",
    logoutConfirm: "Confirm sign out of N8D?",
    ready: "Ready",
    syncAction: "Sync",
    syncing: "Syncing...",
    noPermission: "You do not have SN8 permission to sync.",
    syncDialogTitle: "Syncing",
    syncDialogMessage: "Please wait until synchronization is complete.",
    syncDialogPeriod: "Sync month:",
    syncSuccess: "Synced {count} items.",
    syncWarning: "Sync completed with a warning.",
    syncFail: "Sync failed, please try again.",
    sessionExpired: "Your session has expired, please sign in again.",
    requestFailed: "Something went wrong, please try again.",
    lastSyncLabel: "Last synced: {date}",
    lastSyncNone: "Not synced yet",
  },
  zh: {
    brandSubtitle: "Material Management",
    heroTitle: "选择模块以开始",
    tracking: "TRACKING",
    trackingDesc: "实时物料流程，随每批货物持续更新。",
    inventoryMonth: "月度盘点",
    monthRoleUnit: "单位盘点",
    monthRoleAccounting: "会计盘点",
    inventoryDay: "每日盘点",
    inventoryDayDesc: "每日数据结算",
    login: "登录",
    logout: "登出",
    notLoggedIn: "请登录以继续。",
    logoutConfirm: "确认登出 N8D？",
    ready: "准备就绪",
    syncAction: "同步",
    syncing: "同步中...",
    noPermission: "您没有 SN8 同步权限。",
    syncDialogTitle: "正在同步",
    syncDialogMessage: "请等待同步完成。",
    syncDialogPeriod: "同步月份：",
    syncSuccess: "已同步 {count} 项。",
    syncWarning: "同步完成，但有警告。",
    syncFail: "同步失败，请重试。",
    sessionExpired: "登录已过期，请重新登录。",
    requestFailed: "发生错误，请重试。",
    lastSyncLabel: "上次同步：{date}",
    lastSyncNone: "尚未同步",
  },
};

const copy = computed(() => messages[locale.value] || messages.vi);
const userLabel = computed(
  () =>
    [branch.value, userId.value].filter(Boolean).join(" · ") ||
    copy.value.ready,
);
const hasSN8Permission = computed(() =>
  ["1", "Y", "YES", "TRUE"].includes(sn8.value.trim().toUpperCase()),
);

const modules = computed(() => [
  {
    id: "tracking",
    title: copy.value.tracking,
    description: copy.value.trackingDesc,
    route: "/n8d/tracking",
    requiresAuth: false,
  },
]);

const dayModule = computed(() => ({
  id: "day",
  title: copy.value.inventoryDay,
  description: copy.value.inventoryDayDesc,
  inventoryType: "day",
  requiresAuth: true,
}));

const decodeJwtPayload = (token) => {
  try {
    const encoded = token.split(".")[1];
    if (!encoded) return null;
    const normalized = encoded.replace(/-/g, "+").replace(/_/g, "/");
    const padded = normalized.padEnd(Math.ceil(normalized.length / 4) * 4, "=");
    return JSON.parse(atob(padded));
  } catch {
    return null;
  }
};

const syncSession = () => {
  const token = getToken();
  const payload = token ? decodeJwtPayload(token) : null;
  userId.value =
    payload?.USERID || payload?.userid || localStorage.getItem("USERID") || "";
  branch.value =
    payload?.GSBH || payload?.gsbh || localStorage.getItem("GSBH") || "";
  sn8.value = String(payload?.SN8 || localStorage.getItem("SN8") || "");
  isLoggedIn.value = Boolean(token && userId.value);
  if (userId.value) localStorage.setItem("USERID", userId.value);
  if (branch.value) localStorage.setItem("GSBH", branch.value);
};

const showLoginRequired = () => {
  snackbar.value = true;
  window.clearTimeout(snackbarTimer);
  snackbarTimer = window.setTimeout(() => {
    snackbar.value = false;
  }, 3500);
};

const openModule = (module) => {
  if (module.requiresAuth && !isLoggedIn.value) {
    showLoginRequired();
    return;
  }
  if (module.inventoryType) {
    router.push({
      path: "/n8d/inventory",
      query: { type: module.inventoryType },
    });
    return;
  }
  router.push(module.route);
};

const selectMonthRole = (role) => {
  if (!isLoggedIn.value) {
    showLoginRequired();
    return;
  }
  router.push({ path: "/n8d/inventory", query: { type: "month", role } });
};

const pad2 = (v) => String(v).padStart(2, "0");
const getLocalMonth = () => {
  const d = new Date();
  return `${d.getFullYear()}-${pad2(d.getMonth() + 1)}`;
};

const getRequestErrorMessage = (e) => {
  if (e?.response?.status === 401) return copy.value.sessionExpired;
  if (e?.response?.status === 403) return copy.value.noPermission;
  return copy.value.requestFailed;
};

const showSyncToast = (msg, color = "success") => {
  window.clearTimeout(syncToastTimer);
  syncToast.value = { show: true, message: msg, color };
  syncToastTimer = window.setTimeout(() => {
    syncToast.value.show = false;
  }, 3500);
};

const fetchLastSync = async () => {
  if (!branch.value) return;
  try {
    const res = await axios.get(
      "/api/v1/scan/number-scan/inventory-month/last-sync",
      {
        params: {
          Date: getLocalMonth(),
          GSBH: branch.value,
        },
      },
    );
    lastSync.value = {
      SyncDate: res?.data?.data?.SyncDate || "",
      SyncUserID: res?.data?.data?.SyncUserID || "",
    };
  } catch {
    // Không hiện lỗi cho phần thông tin phụ này, giữ nguyên giá trị cũ.
  }
};

const lastSyncText = computed(() => {
  if (!lastSync.value.SyncDate) return copy.value.lastSyncNone;
  return copy.value.lastSyncLabel.replace("{date}", lastSync.value.SyncDate);
});

const syncInventoryMonth = async () => {
  // Chỉ cho đồng bộ khi đã đăng nhập, y chang trang Kiểm Kê
  if (!isLoggedIn.value) {
    showLoginRequired();
    return;
  }
  if (!hasSN8Permission.value) {
    showSyncToast(copy.value.noPermission, "error");
    return;
  }
  if (syncingMonth.value) return;
  syncSession();
  if (!isLoggedIn.value) return;
  if (!hasSN8Permission.value) {
    showSyncToast(copy.value.noPermission, "error");
    return;
  }
  syncingMonth.value = true;
  try {
    const res = await axios.post(
      "/api/v1/scan/number-scan/inventory-month/sync",
      null,
      {
        params: {
          Date: getLocalMonth(),
          GSBH: branch.value,
          USERID: userId.value,
        },
      },
    );
    const count = res?.data?.data?.synced ?? 0;
    const resType = res?.data?.type;
    if (resType === "warning") showSyncToast(copy.value.syncWarning, "warning");
    else
      showSyncToast(
        copy.value.syncSuccess.replace("{count}", count),
        "success",
      );
    fetchLastSync();
  } catch (e) {
    showSyncToast(getRequestErrorMessage(e) || copy.value.syncFail, "error");
  } finally {
    syncingMonth.value = false;
  }
};

const goToLogin = () => {
  snackbar.value = false;
  router.push("/n8d/login");
};

const scheduleSessionSync = () => {
  window.clearTimeout(authSyncTimer);
  authSyncTimer = window.setTimeout(() => {
    authSyncTimer = undefined;
    syncSession();
  }, 100);
};

const confirmLogout = () => {
  if (!window.confirm(copy.value.logoutConfirm)) return;
  clearSession();
  syncSession();
};

const LOCALE_CODES = ["vi", "en", "zh"];

onMounted(() => {
  const savedLocale = localStorage.getItem("APP_LOCALE");
  if (LOCALE_CODES.includes(savedLocale)) locale.value = savedLocale;
  syncSession();
  fetchLastSync();
  window.addEventListener("auth-changed", scheduleSessionSync);
  window.addEventListener("storage", scheduleSessionSync);
});

onBeforeUnmount(() => {
  window.clearTimeout(snackbarTimer);
  window.clearTimeout(syncToastTimer);
  window.clearTimeout(authSyncTimer);
  window.removeEventListener("auth-changed", scheduleSessionSync);
  window.removeEventListener("storage", scheduleSessionSync);
});
</script>

<script>
export default {
  mounted() {
    const cards = this.$el.querySelectorAll(".module-card");
    const handleMouseMove = (e) => {
      cards.forEach((card) => {
        const rect = card.getBoundingClientRect();
        card.style.setProperty("--mouse-x", `${e.clientX - rect.left}px`);
        card.style.setProperty("--mouse-y", `${e.clientY - rect.top}px`);
      });
    };
    this.$el.addEventListener("mousemove", handleMouseMove);
    this._cleanupMouseMove = () =>
      this.$el.removeEventListener("mousemove", handleMouseMove);
  },
  beforeUnmount() {
    if (this._cleanupMouseMove) this._cleanupMouseMove();
  },
};
</script>

<style scoped>
/* Reset & Base */
* {
  box-sizing: border-box;
}
button {
  font: inherit;
  border: none;
  outline: none;
  background: none;
  color: inherit;
  cursor: pointer;
}

.n8d-home {
  --ink: #1e293b;
  --ink-soft: #475569;
  --muted: #94a3b8;
  --paper: #f1f5f9;
  --surface: #ffffff;
  --line: #e2e8f0;
  --line-strong: #cbd5e1;
  --danger: #ef4444;
  --brand: 15 23 42;

  --accent-tracking: 37 99 235;
  --accent-month: 5 150 105; /* Emerald */
  --accent-day: 217 119 6; /* Amber */

  --font-display: -apple-system, BlinkMacSystemFont, "Segoe UI", "Inter", Roboto,
    "Helvetica Neue", Arial, sans-serif;
  --font-mono: ui-monospace, "SFMono-Regular", "Roboto Mono", Menlo, Consolas,
    monospace;

  color-scheme: light;
  position: relative;
  overflow: hidden;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--paper);
  color: var(--ink);
  font-family: var(--font-display);
}

.bg-mesh {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  opacity: 0.6;
  mask-image: radial-gradient(
    ellipse 80% 50% at 50% 0%,
    black 20%,
    transparent 70%
  );
  -webkit-mask-image: radial-gradient(
    ellipse 80% 50% at 50% 0%,
    black 20%,
    transparent 70%
  );
}

.glow {
  position: absolute;
  width: 500px;
  height: 500px;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.4;
}
.glow--tracking {
  top: -15%;
  left: 5%;
  background: rgba(var(--accent-tracking), 0.2);
}
.glow--month {
  top: -10%;
  right: 10%;
  background: rgba(var(--accent-month), 0.25);
}
.glow--day {
  top: 5%;
  left: 40%;
  background: rgba(var(--accent-day), 0.2);
}

/* --- Main Content --- */
.home-content {
  position: relative;
  z-index: 1;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  width: min(1100px, 94%);
  margin: 0 auto;
  padding: 52px 0 72px;
}
.hero {
  margin-bottom: 44px;
}
.hero-eyebrow {
  display: inline-block;
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--ink-soft);
  background: var(--surface);
  padding: 8px 14px;
  border-radius: 8px;
  margin-bottom: 20px;
  border: 1px solid var(--line);
  box-shadow: 0 2px 4px -1px rgba(0, 0, 0, 0.04);
}
.hero h1 {
  margin: 0;
  font-size: 44px;
  font-weight: 900;
  letter-spacing: -0.035em;
  color: var(--ink);
  line-height: 1.1;
  text-transform: uppercase;
}

/* --- Bento Grid --- */
.module-bento {
  display: grid;
  grid-template-columns: 1.3fr 1fr;
  gap: 24px;
  align-items: stretch;
}
.bento-stack {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* --- Chung cho cả 3 Card (Bóng đổ nền tảng) --- */
.module-card {
  position: relative;
  border-radius: 24px;
  background: var(--surface);
  border: 1px solid transparent;
  text-align: left;
  overflow: hidden;
  box-shadow: 0 0 0 1px rgba(15, 23, 42, 0.03),
    0 4px 6px -1px rgba(15, 23, 42, 0.02),
    0 10px 20px -5px rgba(15, 23, 42, 0.06),
    0 25px 50px -12px rgba(15, 23, 42, 0.15);
  transition: all 0.35s cubic-bezier(0.22, 1, 0.36, 1);
}

.module-card::after {
  content: "";
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: radial-gradient(
    600px circle at var(--mouse-x, 50%) var(--mouse-y, 50%),
    rgba(var(--accent), 0.04),
    transparent 40%
  );
  opacity: 0;
  transition: opacity 0.3s ease;
  pointer-events: none;
}

.module-card:hover::after {
  opacity: 1;
}
.module-card:focus-visible {
  outline: 2px solid rgb(var(--accent));
  outline-offset: 2px;
}

.module-card--tracking {
  --accent: var(--accent-tracking);
}
.module-card--month {
  --accent: var(--accent-month);
}
.module-card--day {
  --accent: var(--accent-day);
}

/* --- FEATURE CARD (Trái - Theo dõi) --- */
.module-card--feature {
  display: flex;
  flex-direction: column;
  padding: 36px;
}

.module-card--feature:hover {
  transform: translateY(-6px) scale(1.01);
  border-color: rgba(var(--accent), 0.3);
  box-shadow: 0 0 0 1px rgba(var(--accent), 0.1),
    0 4px 8px -2px rgba(var(--accent), 0.1),
    0 16px 32px -8px rgba(15, 23, 42, 0.1),
    0 48px 64px -16px rgba(var(--accent), 0.25);
}

.feature-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
}

.card-icon {
  display: grid;
  place-items: center;
  width: 52px;
  height: 52px;
  border-radius: 16px;
  background: rgba(var(--accent), 0.08);
  flex-shrink: 0;
  transition: all 0.35s cubic-bezier(0.22, 1, 0.36, 1);
}
.module-card:hover .card-icon {
  background: rgb(var(--accent));
  box-shadow: 0 8px 24px -4px rgba(var(--accent), 0.5);
  transform: scale(1.05);
}

.card-icon svg {
  width: 24px;
  height: 24px;
  stroke: rgb(var(--accent));
  fill: none;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
  transition: stroke 0.3s ease;
}
.module-card:hover .card-icon svg {
  stroke: #ffffff;
}

.card-icon--lg {
  width: 68px;
  height: 68px;
  border-radius: 20px;
}
.card-icon--lg svg {
  width: 30px;
  height: 30px;
}

.ticket-divider {
  display: block;
  height: 1px;
  margin: 28px 0;
  background-image: linear-gradient(
    to right,
    var(--line-strong) 60%,
    transparent 0%
  );
  background-size: 12px 1px;
  background-repeat: repeat-x;
}

.feature-body h3 {
  margin: 0 0 12px;
  font-size: 28px;
  font-weight: 900;
  letter-spacing: -0.025em;
  color: var(--ink);
  text-transform: uppercase;
}
.card-desc {
  margin: 0;
  font-size: 15px;
  color: var(--muted);
  line-height: 1.6;
  max-width: 42ch;
  font-weight: 500;
}
.card-desc--row {
  font-size: 13px;
  line-height: 1.4;
  max-width: none;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
}

.card-foot {
  margin-top: auto;
  padding-top: 24px;
  border-top: 1px solid var(--line);
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.card-arrow {
  display: grid;
  place-items: center;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  color: var(--muted);
  background: var(--paper);
  transition: all 0.35s cubic-bezier(0.22, 1, 0.36, 1);
}
.card-arrow svg {
  width: 16px;
  height: 16px;
  stroke: currentColor;
  fill: none;
  stroke-width: 2.5;
  stroke-linecap: round;
  stroke-linejoin: round;
}
.module-card:hover .card-arrow {
  color: #fff;
  background: rgb(var(--accent));
  box-shadow: 0 4px 12px -2px rgba(var(--accent), 0.5);
  transform: translateX(4px);
}

/* ---------------------------------------------------------
   ROW CARDS (PHÍA DƯỚI - NỔI BẬT MÀU SẮC RIÊNG BIỆT)
   --------------------------------------------------------- */
.module-card--row {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 28px;
  box-shadow: none;
}

.row-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.row-body h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 800;
  letter-spacing: -0.01em;
  color: var(--ink);
  text-transform: uppercase;
}

.row-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Tăng độ đậm cho Icon của Card nhỏ */
.module-card--row .card-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: rgba(var(--accent), 0.12);
  border: 1px solid rgba(var(--accent), 0.15);
}
.module-card--row .card-icon svg {
  width: 26px;
  height: 26px;
}

/* --- CARD KIỂM KÊ THÁNG (XANH NGỌC LỤC) --- */
.module-card--month {
  background: linear-gradient(
    135deg,
    rgba(5, 150, 105, 0.08) 0%,
    rgba(255, 255, 255, 1) 50%
  );
  border: 1px solid rgba(5, 150, 105, 0.15);
  box-shadow: 0 0 0 1px rgba(5, 150, 105, 0.05),
    0 8px 20px -6px rgba(5, 150, 105, 0.2),
    0 20px 40px -10px rgba(15, 23, 42, 0.1);
  cursor: default; /* Card本月 không phải là button nữa */
}

.module-card--month:hover {
  transform: translateY(-2px); /* Hiệu ứng nhẹ hơn */
  border-color: rgba(5, 150, 105, 0.4);
  box-shadow: 0 0 0 1px rgba(5, 150, 105, 0.1),
    0 16px 30px -8px rgba(5, 150, 105, 0.3),
    0 40px 60px -15px rgba(5, 150, 105, 0.35);
}

/* Thẻ Kiểm Kê Tháng giờ là feature card (to, bên trái) — cần đủ chỗ cho 3 chức năng con + badge đồng bộ */
.module-card--month.module-card--feature {
  min-height: 340px;
  justify-content: flex-start;
}
.module-card--month .feature-body {
  display: flex;
  flex-direction: column;
  flex: 1;
}
.module-card--month .feature-body h3 {
  margin: 0 0 16px;
}

/* --- Hàng chứa nút Đồng bộ + thời gian đồng bộ lần cuối (thẳng hàng, ngang nhau) --- */
.month-sync-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 10px;
  margin-bottom: 12px;
}
.month-sync-row .month-action-btn--sync {
  width: auto;
  flex: 1 1 auto;
  aspect-ratio: auto;
  flex-direction: row;
  justify-content: center;
  gap: 8px;
  padding: 12px 20px;
  min-height: 44px;
}
.month-sync-row .month-action-btn--sync svg,
.month-sync-row .month-action-btn--sync .spinner {
  width: 17px;
  height: 17px;
}
.month-sync-row .last-sync-badge {
  flex: 0 0 auto;
  min-width: 0;
  min-height: 44px;
  justify-content: center;
}

/* --- 2 NÚT KIỂM KÊ ĐƠN VỊ / KẾ TOÁN (XẾP DỌC) --- */
.month-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.month-action-btn {
  width: 100%;
  aspect-ratio: auto;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 13px 16px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 700;
  text-align: left;
  border: 1px solid transparent;
  transition: all 0.2s ease;
}
.month-action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.month-action-btn svg {
  width: 16px;
  height: 16px;
  stroke: currentColor;
  fill: none;
  stroke-width: 2.5;
  stroke-linecap: round;
  stroke-linejoin: round;
  transition: transform 0.2s ease;
}
.month-action-btn span {
  line-height: 1.25;
}

/* Màu riêng cho nút Đồng bộ (Hổ phách) */
.month-action-btn--sync {
  background: rgba(217, 119, 6, 0.1);
  color: rgb(180, 83, 9);
  border-color: rgba(217, 119, 6, 0.2);
}
.month-action-btn--sync:hover:not(:disabled) {
  background: rgb(217, 119, 6);
  color: #fff;
  border-color: rgb(217, 119, 6);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px -2px rgba(217, 119, 6, 0.4);
}
.month-action-btn--sync .spinner {
  width: 22px;
  height: 22px;
  border: 3px solid rgba(180, 83, 9, 0.3);
  border-top-color: rgb(180, 83, 9);
}
.month-action-btn--sync:hover:not(:disabled) .spinner {
  border-color: rgba(255, 255, 255, 0.35);
  border-top-color: #fff;
}

/* Badge "Đồng bộ lần cuối" — nổi bật, ngang hàng với nút Đồng bộ */
.last-sync-badge {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-size: 14px;
  font-weight: 800;
  letter-spacing: -0.01em;
  color: rgb(146, 64, 14);
  background: rgba(217, 119, 6, 0.14);
  border: 1px solid rgba(217, 119, 6, 0.3);
  padding: 8px 16px;
  border-radius: 999px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  box-shadow: 0 4px 10px -3px rgba(217, 119, 6, 0.3);
}
.last-sync-badge::before {
  content: "";
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: rgb(217, 119, 6);
  flex-shrink: 0;
  animation: pulse 2s ease-in-out infinite;
}


/* Màu riêng cho nút Kiểm kê Đơn vị (Xanh ngọc lục) */
.month-action-btn--unit {
  background: rgba(5, 150, 105, 0.1);
  color: rgb(5, 150, 105);
  border-color: rgba(5, 150, 105, 0.2);
}
.month-action-btn--unit:hover {
  background: rgb(5, 150, 105);
  color: #fff;
  border-color: rgb(5, 150, 105);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px -2px rgba(5, 150, 105, 0.4);
}
.month-action-btn--unit:hover svg {
  transform: translateX(3px);
}

/* Màu riêng cho nút Kiểm kê Kế toán (Xanh dương) */
.month-action-btn--accounting {
  background: rgba(37, 99, 235, 0.1);
  color: rgb(37, 99, 235);
  border-color: rgba(37, 99, 235, 0.2);
}
.month-action-btn--accounting:hover {
  background: rgb(37, 99, 235);
  color: #fff;
  border-color: rgb(37, 99, 235);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px -2px rgba(37, 99, 235, 0.4);
}
.month-action-btn--accounting:hover svg {
  transform: translateX(3px);
}

/* --- CARD KIỂM KÊ NGÀY (CAM ĐẤT) --- */
.module-card--day {
  background: linear-gradient(
    135deg,
    rgba(217, 119, 6, 0.08) 0%,
    rgba(255, 255, 255, 1) 50%
  );
  border: 1px solid rgba(217, 119, 6, 0.15);
  box-shadow: 0 0 0 1px rgba(217, 119, 6, 0.05),
    0 8px 20px -6px rgba(217, 119, 6, 0.2),
    0 20px 40px -10px rgba(15, 23, 42, 0.1);
}

.module-card--day:hover {
  transform: translateY(-6px) scale(1.02);
  background: linear-gradient(
    135deg,
    rgba(217, 119, 6, 0.15) 0%,
    rgba(255, 255, 255, 1) 40%
  );
  border-color: rgba(217, 119, 6, 0.4);
  box-shadow: 0 0 0 1px rgba(217, 119, 6, 0.1),
    0 16px 30px -8px rgba(217, 119, 6, 0.3),
    0 40px 60px -15px rgba(217, 119, 6, 0.35);
}

.module-card--row:active {
  transform: translateY(-2px) scale(1.005);
}

/* --- Blocking sync dialog --- */
.sync-blocking-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: grid;
  place-items: center;
  padding: 24px;
  background: rgba(15, 23, 42, 0.42);
  backdrop-filter: blur(3px);
  -webkit-backdrop-filter: blur(3px);
  cursor: wait;
}
.sync-blocking-dialog {
  display: grid;
  min-width: min(320px, 100%);
  max-width: 420px;
  justify-items: center;
  gap: 10px;
  padding: 28px 34px;
  border: 1px solid rgba(255, 255, 255, 0.7);
  border-radius: 18px;
  background: #fff;
  color: var(--ink);
  text-align: center;
  box-shadow: 0 24px 60px rgba(15, 23, 42, 0.28);
  cursor: wait;
}
.sync-blocking-dialog strong {
  font-size: 17px;
  font-weight: 800;
}
.sync-blocking-dialog > span:not(.sync-dialog-spinner) {
  color: var(--ink-soft);
  font-size: 13px;
}
.sync-dialog-period {
  margin-top: 2px;
  padding: 4px 12px;
  border-radius: 999px;
  background: rgba(79, 70, 229, 0.1);
  color: #4f46e5 !important;
  font-weight: 700;
  font-size: 12px;
  letter-spacing: 0.02em;
}
.sync-dialog-spinner {
  width: 30px;
  height: 30px;
  border: 3px solid #dbeafe;
  border-top-color: #2563eb;
  border-radius: 50%;
  animation: spin 0.75s linear infinite;
}

/* --- Toast --- */
.auth-toast {
  position: fixed;
  bottom: 32px;
  right: 32px;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 8px 8px 24px;
  background: var(--surface);
  border: 1px solid var(--line-strong);
  border-radius: 16px;
  box-shadow: 0 25px 50px -12px rgba(15, 23, 42, 0.25);
  z-index: 100;
}
.auth-toast span {
  font-size: 14px;
  color: var(--ink-soft);
  font-weight: 500;
}
.auth-toast button {
  padding: 10px 20px;
  background: var(--danger);
  color: #fff;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 700;
  box-shadow: 0 4px 12px -2px rgba(239, 68, 68, 0.4);
  transition: all 0.2s ease;
}
.auth-toast button:hover {
  filter: brightness(1.1);
  transform: translateY(-1px);
}

.sync-toast {
  position: fixed;
  bottom: 32px;
  right: 32px;
  z-index: 100;
  display: flex;
  align-items: center;
  gap: 12px;
  max-width: 420px;
  padding: 12px 16px;
  border-radius: 14px;
  background: rgba(11, 18, 32, 0.9);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  color: #f8fafc;
  font-size: 13.5px;
  font-weight: 700;
  box-shadow: 0 20px 44px -14px rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.12);
}
.sync-toast.success svg {
  width: 18px;
  fill: none;
  stroke: #a3e635;
  stroke-width: 2;
}
.sync-toast.warning svg {
  width: 18px;
  fill: none;
  stroke: #fbbf24;
  stroke-width: 2;
}
.sync-toast.error svg {
  width: 18px;
  fill: none;
  stroke: #f87171;
  stroke-width: 2;
}
.sync-toast span {
  flex: 1;
}
.sync-toast button {
  width: 24px;
  height: 24px;
  display: grid;
  place-items: center;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  color: #94a3b8;
  font-size: 16px;
  transition: 0.15s;
}
.sync-toast button:hover {
  background: rgba(255, 255, 255, 0.18);
  color: #fff;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.4s cubic-bezier(0.22, 1, 0.36, 1);
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
  }
  .sync-dialog-spinner {
    animation-duration: 1.5s !important;
  }
}

/* --- Responsive --- */
@media (max-width: 900px) {
  .module-bento {
    grid-template-columns: 1fr;
  }
  .hero h1 {
    font-size: 32px;
  }
  .home-content {
    padding: 40px 0;
  }
  .glow {
    width: 300px;
    height: 300px;
    filter: blur(80px);
  }
}

@media (max-width: 480px) {
  .module-card--feature {
    padding: 24px;
  }
  .feature-body h3 {
    font-size: 24px;
  }
  .module-card--row {
    padding: 20px;
  }
  .month-sync-row {
    flex-wrap: wrap;
  }
  .month-sync-row .month-action-btn--sync {
    flex: 1 1 100%;
  }
  .last-sync-badge {
    max-width: 100%;
    font-size: 12.5px;
    padding: 7px 12px;
  }
}
</style>