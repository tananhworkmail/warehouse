<template>
  <div class="landing">
    <!-- Hiệu ứng nền trang trí thêm -->
    <div class="landing__bg-glow landing__bg-glow--1"></div>
    <div class="landing__bg-glow landing__bg-glow--2"></div>

    <header class="landing__header">
      <img src="@/assets/Logo.png" alt="LAIYIH" class="landing__logo" />
      <div class="landing__header-right">
        <div class="landing__clock">
          <svg
            class="clock-icon"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <circle cx="12" cy="12" r="10"></circle>
            <polyline points="12 6 12 12 16 14"></polyline>
          </svg>
          {{ currentTime }}
        </div>
        <div class="landing__lang">
          <button
            v-for="l in langs"
            :key="l.code"
            class="lang-btn"
            :class="{ 'is-active': lang === l.code }"
            @click="setLang(l.code)"
          >
            {{ l.label }}
          </button>
        </div>
      </div>
    </header>

    <div class="landing__hero">
      <p class="landing__eyebrow">{{ i18n.eyebrow }}</p>
      <h1 class="landing__title">
        {{ i18n.titleMain }} <em>{{ i18n.titleEm }}</em>
      </h1>
    </div>

    <main class="landing__cards">
      <!-- Thẻ Vật tư -->
      <div class="wh-card wh-card--kvt" :class="{ 'is-open': openCards.kvt }">
        <div class="wh-card__head" @click="toggleCard('kvt')">
          <div class="wh-card__icon">
            <svg viewBox="0 0 64 64" fill="none">
              <rect
                x="8"
                y="28"
                width="48"
                height="28"
                rx="3"
                stroke="currentColor"
                stroke-width="3"
              />
              <path
                d="M4 28L32 8L60 28"
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
              />
              <rect
                x="24"
                y="40"
                width="16"
                height="16"
                rx="2"
                stroke="currentColor"
                stroke-width="2.5"
              />
              <line
                x1="32"
                y1="40"
                x2="32"
                y2="56"
                stroke="currentColor"
                stroke-width="2"
              />
              <rect
                x="14"
                y="34"
                width="10"
                height="10"
                rx="1.5"
                stroke="currentColor"
                stroke-width="2"
              />
              <rect
                x="40"
                y="34"
                width="10"
                height="10"
                rx="1.5"
                stroke="currentColor"
                stroke-width="2"
              />
            </svg>
          </div>
          <div class="wh-card__body">
            <p class="wh-card__label">{{ i18n.warehouseLabel }}</p>
            <h2 class="wh-card__name">{{ i18n.kvtName }}</h2>
          </div>
          <div class="wh-card__chevron">
            <svg viewBox="0 0 24 24" fill="none">
              <path
                d="M6 9l6 6 6-6"
                stroke="currentColor"
                stroke-width="2.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
        </div>
        <transition name="expand">
          <div v-if="openCards.kvt" class="wh-submenu-wrapper">
            <div class="wh-submenu">
              <button
                v-for="item in kvtMenu"
                :key="item.key"
                class="wh-submenu__item"
                @click.stop="navigate(item)"
              >
                <span class="submenu-icon">{{ item.icon }}</span>
                <span class="submenu-label">{{ item.label }}</span>
                <svg class="submenu-arrow" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M5 12h14M13 6l6 6-6 6"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </button>
            </div>
          </div>
        </transition>
      </div>

      <!-- Thẻ Đế Sinh Quản -->
      <div
        class="wh-card wh-card--khode"
        :class="{ 'is-open': openCards.khode }"
      >
        <div class="wh-card__head" @click="toggleCard('khode')">
          <div class="wh-card__icon">
            <svg viewBox="0 0 64 64" fill="none">
              <rect
                x="8"
                y="28"
                width="48"
                height="28"
                rx="3"
                stroke="currentColor"
                stroke-width="3"
              />
              <path
                d="M4 28L32 8L60 28"
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
              />
              <ellipse
                cx="32"
                cy="44"
                rx="12"
                ry="6"
                stroke="currentColor"
                stroke-width="2.5"
              />
              <path
                d="M20 44c0 3.314 5.373 6 12 6s12-2.686 12-6"
                stroke="currentColor"
                stroke-width="2"
              />
              <line
                x1="32"
                y1="38"
                x2="32"
                y2="50"
                stroke="currentColor"
                stroke-width="2"
              />
            </svg>
          </div>
          <div class="wh-card__body">
            <p class="wh-card__label">{{ i18n.warehouseLabel }}</p>
            <h2 class="wh-card__name">{{ i18n.khodeName }}</h2>
          </div>
          <div class="wh-card__chevron">
            <svg viewBox="0 0 24 24" fill="none">
              <path
                d="M6 9l6 6 6-6"
                stroke="currentColor"
                stroke-width="2.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
        </div>
        <transition name="expand">
          <div v-if="openCards.khode" class="wh-submenu-wrapper">
            <div class="wh-submenu">
              <button
                v-for="item in khodeMenu"
                :key="item.key"
                class="wh-submenu__item"
                @click.stop="navigate(item)"
              >
                <span class="submenu-icon">{{ item.icon }}</span>
                <span class="submenu-label">{{ item.label }}</span>
                <svg class="submenu-arrow" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M5 12h14M13 6l6 6-6 6"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </button>
            </div>
          </div>
        </transition>
      </div>

      <!-- [MỚI] Thẻ Vật tư Đế -->
      <div
        class="wh-card wh-card--khodevt"
        :class="{ 'is-open': openCards.khodevt }"
      >
        <div class="wh-card__head" @click="toggleCard('khodevt')">
          <div class="wh-card__icon">
            <svg viewBox="0 0 64 64" fill="none">
              <path
                d="M32 12L8 24L32 36L56 24L32 12Z"
                stroke="currentColor"
                stroke-width="3"
                stroke-linejoin="round"
              />
              <path
                d="M8 36L32 48L56 36"
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M8 48L32 60L56 48"
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
          <div class="wh-card__body">
            <p class="wh-card__label">{{ i18n.warehouseLabel }}</p>
            <h2 class="wh-card__name">{{ i18n.khodeVtName }}</h2>
          </div>
          <div class="wh-card__chevron">
            <svg viewBox="0 0 24 24" fill="none">
              <path
                d="M6 9l6 6 6-6"
                stroke="currentColor"
                stroke-width="2.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
        </div>
        <transition name="expand">
          <div v-if="openCards.khodevt" class="wh-submenu-wrapper">
            <div class="wh-submenu">
              <button
                v-for="item in khodeVtMenu"
                :key="item.key"
                class="wh-submenu__item"
                @click.stop="navigate(item)"
              >
                <span class="submenu-icon">{{ item.icon }}</span>
                <span class="submenu-label">{{ item.label }}</span>
                <svg class="submenu-arrow" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M5 12h14M13 6l6 6-6 6"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </button>
            </div>
          </div>
        </transition>
      </div>

      <!-- Thẻ Phòng thí nghiệm -->
      <div
        class="wh-card wh-card--lab"
        :class="{ 'is-open': openCards.laboratory }"
      >
        <div class="wh-card__head" @click="toggleCard('laboratory')">
          <div class="wh-card__icon">
            <svg viewBox="0 0 64 64" fill="none">
              <path
                d="M24 8v18L12 46a8 8 0 008 10h24a8 8 0 008-10L40 26V8"
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M24 18h16"
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
              />
              <path
                d="M20 42h24"
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
              />
            </svg>
          </div>
          <div class="wh-card__body">
            <p class="wh-card__label">LAB</p>
            <h2 class="wh-card__name">{{ i18n.labName }}</h2>
          </div>
          <div class="wh-card__chevron">
            <svg viewBox="0 0 24 24" fill="none">
              <path
                d="M6 9l6 6 6-6"
                stroke="currentColor"
                stroke-width="2.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
        </div>
        <transition name="expand">
          <div v-if="openCards.laboratory" class="wh-submenu-wrapper">
            <div class="wh-submenu">
              <button
                v-for="item in laboratoryMenu"
                :key="item.key"
                class="wh-submenu__item"
                @click.stop="navigate(item)"
              >
                <span class="submenu-icon">{{ item.icon }}</span>
                <span class="submenu-label">{{ item.label }}</span>
                <svg class="submenu-arrow" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M5 12h14M13 6l6 6-6 6"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </button>
            </div>
          </div>
        </transition>
      </div>

      <!-- Thẻ Thành phẩm -->
      <div class="wh-card wh-card--ktp" :class="{ 'is-open': openCards.ktp }">
        <div class="wh-card__head" @click="toggleCard('ktp')">
          <div class="wh-card__icon">
            <svg viewBox="0 0 64 64" fill="none">
              <path
                d="M12 24L32 12L52 24V40L32 52L12 40V24Z"
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M32 12V32M32 32L12 24M32 32L52 24"
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
          <div class="wh-card__body">
            <p class="wh-card__label">{{ i18n.warehouseLabel }}</p>
            <h2 class="wh-card__name">{{ i18n.ktpName }}</h2>
          </div>
          <div class="wh-card__chevron">
            <svg viewBox="0 0 24 24" fill="none">
              <path
                d="M6 9l6 6 6-6"
                stroke="currentColor"
                stroke-width="2.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
        </div>
        <transition name="expand">
          <div v-if="openCards.ktp" class="wh-submenu-wrapper">
            <div class="wh-submenu">
              <button
                v-for="item in ktpMenu"
                :key="item.key"
                class="wh-submenu__item"
                @click.stop="navigate(item)"
              >
                <span class="submenu-icon">{{ item.icon }}</span>
                <span class="submenu-label">{{ item.label }}</span>
                <svg class="submenu-arrow" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M5 12h14M13 6l6 6-6 6"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </button>
            </div>
          </div>
        </transition>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import dayjs from "dayjs";
// ✅ CHỈ import để lấy lang và setLang dùng chung singleton
import { useWarehouseMapI18n } from "@/hooks/useWarehouseMapI18n";

const router = useRouter();

// ✅ lang và setLang từ singleton — đồng bộ với các trang khác
const { lang, langOptions: langs, setLang } = useWarehouseMapI18n();

const currentTime = ref(dayjs().format("YYYY/MM/DD HH:mm:ss"));

// Đã thêm khodevt vào object để có thể đóng mở được
const openCards = reactive({
  kvt: true,
  khode: true,
  khodevt: true,
  laboratory: true,
  ktp: true,
});
let clockTimer = null;

onMounted(() => {
  clockTimer = setInterval(() => {
    currentTime.value = dayjs().format("YYYY/MM/DD HH:mm:ss");
  }, 1000);
});

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer);
});

const toggleCard = (key) => {
  openCards[key] = !openCards[key];
};

const openInventory = () => {
  window.open("https://192.168.71.9:8100/", "inventoryWindow");
};

// Hàm mở link cho Kho vật tư Đế
const openKhodeVt = () => {
  window.open("http://192.168.71.9:8082/", "khodeVtWindow");
};
const openKhodeVtIMQC = () => {
  window.open("http://192.168.71.9:8082/dashboard/imqc", "IMQCWindow");
};

// Đã cập nhật i18n
const messages = {
  vi: {
    eyebrow: "WAREHOUSE MANAGEMENT",
    titleMain: "QUẢN LÝ",
    titleEm: "KHO",
    warehouseLabel: "KHO",
    kvtName: "Vật tư Mặt",
    khodeName: "Đế Sinh Quản", // Sửa thành Đế Sinh Quản
    khodeVtName: "Vật tư Đế", // Thêm mới
    menuKvtMap: "Sơ đồ kho vật tư",
    menuKiemKe: "Kiểm kê",
    menuLossTemSize: "Loss_Tem Size",
    menuKhodeMap: "Sơ đồ kho đế sinh quản",
    menuKhodeVtMap: "Kho vật tư đế",
    menuKhodeIMQC: "IMQC Đế", // Thêm mới
    labName: "Phòng thí nghiệm",
    menuLabMap: "Sơ đồ phòng thí nghiệm",
    menuGoreTexForms: "Biểu mẫu GORE-TEX",
    ktpName: "Thành phẩm",
    menuKtpMap: "Sơ đồ kho thành phẩm",
  },
  en: {
    eyebrow: "WAREHOUSE MANAGEMENT",
    titleMain: "WAREHOUSE",
    titleEm: "MANAGEMENT",
    warehouseLabel: "WAREHOUSE",
    kvtName: "Materials",
    khodeName: "Sole (Planning)",
    khodeVtName: "Sole Materials",
    menuKvtMap: "Materials warehouse map",
    menuKiemKe: "Inventory check",
    menuLossTemSize: "Loss_Tem Size",
    menuKhodeMap: "Sole warehouse map",
    menuKhodeVtMap: "Sole materials map",
    menuKhodeIMQC: "IMQC Sole",
    labName: "Laboratory",
    menuLabMap: "Laboratory map",
    menuGoreTexForms: "GORE-TEX forms",
    ktpName: "Finished Goods",
    menuKtpMap: "Finished goods map",
  },
  zh: {
    eyebrow: "仓库管理系统",
    titleMain: "仓库",
    titleEm: "管理",
    warehouseLabel: "仓库",
    kvtName: "物料仓",
    khodeName: "生管鞋底仓",
    khodeVtName: "鞋底物料仓",
    menuKvtMap: "物料仓库平面图",
    menuKiemKe: "库存盘点",
    menuLossTemSize: "Loss_Tem Size",
    menuKhodeMap: "生管鞋底仓库平面图",
    menuKhodeVtMap: "鞋底物料仓库平面图",
    menuKhodeIMQC: "IMQC 鞋底",
    labName: "实验室",
    menuLabMap: "实验室平面图",
    menuGoreTexForms: "GORE-TEX 表单",
    ktpName: "成品仓",
    menuKtpMap: "成品仓库平面图",
  },
};

// ✅ i18n computed dùng lang từ singleton — tự động reactive khi đổi ngôn ngữ
const i18n = computed(() => messages[lang.value] || messages["vi"]);

// ── Menus ─────────────────────────────
const kvtMenu = computed(() => [
  {
    key: "main",
    icon: "🏗️",
    label: i18n.value.menuKvtMap,
    route: "/warehouse-kvt",
  },
  {
    key: "kiemke",
    icon: "📋",
    label: i18n.value.menuKiemKe,
    action: openInventory,
  },
  {
    key: "lossTemSize",
    icon: "📦",
    label: i18n.value.menuLossTemSize,
    route: "/loss-tem-size",
  },
]);

const khodeMenu = computed(() => [
  {
    key: "main",
    icon: "🦶",
    label: i18n.value.menuKhodeMap,
    route: "/warehouse-khode",
  },
]);

// Menu cho Kho Vật Tư Đế, sử dụng hàm gọi link thẳng tới URL mới
const khodeVtMenu = computed(() => [
  {
    key: "main",
    icon: "🧱",
    label: i18n.value.menuKhodeVtMap,
    action: openKhodeVt,
  },
  {
    key: "IMQC",
    icon: "📋",
    label: i18n.value.menuKhodeIMQC,
    action: openKhodeVtIMQC,
  },
]);

const laboratoryMenu = computed(() => [
  {
    key: "main",
    icon: "🧪",
    label: i18n.value.menuLabMap,
    route: "/laboratory",
  },
  {
    key: "goreTexForms",
    icon: "📝",
    label: i18n.value.menuGoreTexForms,
    route: "/laboratory/forms",
  },
]);

const ktpMenu = computed(() => [
  {
    key: "main",
    icon: "🚚",
    label: i18n.value.menuKtpMap,
    route: "/warehouse-ktp",
  },
]);

const navigate = (item) => {
  if (item.action) {
    item.action();
    return;
  }
  if (item.route) {
    router.push(item.route);
  }
};
</script>

<style lang="scss" scoped>
/* Biến màu sắc */
$blue: #3b82f6;
$amber: #f59e0b;
$purple: #8b5cf6;
$emerald: #10b981;
$cyan: #06b6d4;
$bg-gradient: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
$text-main: #0f172a;
$text-muted: #64748b;
$card-bg: rgba(255, 255, 255, 0.85);
$card-border: rgba(255, 255, 255, 0.6);
$card-shadow: 0 10px 40px rgba(15, 23, 42, 0.06);
$card-shadow-hover: 0 20px 50px rgba(15, 23, 42, 0.12);

.landing {
  min-height: 100vh;
  background: $bg-gradient;
  color: $text-main;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 24px 60px;
  position: relative;
  overflow: hidden;
  font-family: "Segoe UI", system-ui, sans-serif;
  z-index: 1;
}

.landing__bg-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  z-index: -1;
  opacity: 0.6;
}
.landing__bg-glow--1 {
  width: 500px;
  height: 500px;
  background: rgba(59, 130, 246, 0.2);
  top: -100px;
  left: -100px;
}
.landing__bg-glow--2 {
  width: 600px;
  height: 600px;
  background: rgba(139, 92, 246, 0.15);
  bottom: -150px;
  right: -100px;
}

.landing__header {
  width: 100%;
  max-width: 1600px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 0;
}

.landing__logo {
  height: 48px;
  object-fit: contain;
}

.landing__header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.landing__clock {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: $text-main;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.8);
  border-radius: 999px;
  padding: 8px 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);

  .clock-icon {
    width: 18px;
    height: 18px;
    color: $blue;
  }
}

.landing__lang {
  display: flex;
  gap: 6px;
  background: rgba(255, 255, 255, 0.6);
  padding: 4px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.8);
}

.lang-btn {
  padding: 6px 14px;
  border-radius: 999px;
  border: none;
  background: transparent;
  color: $text-muted;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    color: $text-main;
  }

  &.is-active {
    background: #fff;
    color: $blue;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  }
}

.landing__hero {
  text-align: center;
  margin: 60px 0 50px;
}

.landing__eyebrow {
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.2em;
  color: $blue;
  text-transform: uppercase;
  margin: 0 0 16px;
}

.landing__title {
  font-size: clamp(40px, 6vw, 64px);
  font-weight: 900;
  margin: 0;
  line-height: 1.2;
  color: $text-main;

  em {
    font-style: normal;
    background: linear-gradient(135deg, #3b82f6, #8b5cf6);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
}

/* --- CẤU TRÚC 5 CARD LUÔN NẰM NGANG HOẶC RỚT DÒNG KHI MÀN HÌNH NHỎ --- */
.landing__cards {
  width: 100%;
  max-width: 1800px;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap; /* Cho phép rớt dòng khi không đủ chỗ */
  justify-content: center;
  gap: 24px;
  padding: 10px 10px 30px 10px;

  &::-webkit-scrollbar {
    height: 8px;
  }
  &::-webkit-scrollbar-thumb {
    background: rgba(15, 23, 42, 0.2);
    border-radius: 10px;
  }
}

.wh-card {
  flex: 1 1 260px; /* Card tự động co giãn dựa trên min-width */
  min-width: 260px;
  max-width: 100%; /* Tránh việc card tràn ra khỏi viền màn hình trên mobile */
  border-radius: 24px;
  border: 1px solid $card-border;
  background: $card-bg;
  backdrop-filter: blur(12px);
  box-shadow: $card-shadow;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);

  display: flex;
  flex-direction: column;
  justify-content: flex-start;

  &:hover {
    transform: translateY(-6px);
    box-shadow: $card-shadow-hover;
  }
}

.wh-card--kvt {
  border-top: 4px solid $blue;
}
.wh-card--khode {
  border-top: 4px solid $amber;
}
.wh-card--khodevt {
  border-top: 4px solid $cyan;
}
.wh-card--ktp {
  border-top: 4px solid $purple;
}
.wh-card--lab {
  border-top: 4px solid $emerald;
}

.wh-card__head {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  cursor: pointer;
  user-select: none;
  min-height: 115px; 
}

.wh-card__icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: transform 0.3s;

  svg {
    width: 32px;
    height: 32px;
  }
}

.wh-card:hover .wh-card__icon {
  transform: scale(1.05);
}

.wh-card--kvt .wh-card__icon {
  color: $blue;
  background: rgba(59, 130, 246, 0.12);
}
.wh-card--khode .wh-card__icon {
  color: $amber;
  background: rgba(245, 158, 11, 0.12);
}
.wh-card--khodevt .wh-card__icon {
  color: $cyan;
  background: rgba(6, 182, 212, 0.12);
}
.wh-card--ktp .wh-card__icon {
  color: $purple;
  background: rgba(139, 92, 246, 0.12);
}
.wh-card--lab .wh-card__icon {
  color: $emerald;
  background: rgba(16, 185, 129, 0.12);
}

.wh-card__body {
  flex: 1;
}

.wh-card__label {
  font-size: 11px;
  font-weight: 800;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  margin: 0 0 4px;
}

.wh-card--kvt .wh-card__label {
  color: $blue;
}
.wh-card--khode .wh-card__label {
  color: $amber;
}
.wh-card--khodevt .wh-card__label {
  color: $cyan;
}
.wh-card--ktp .wh-card__label {
  color: $purple;
}
.wh-card--lab .wh-card__label {
  color: $emerald;
}

.wh-card__name {
  font-size: 22px;
  font-weight: 900;
  color: $text-main;
  margin: 0 0 6px;
  line-height: 1.1;
}
.wh-card__chevron {
  width: 28px;
  height: 28px;
  color: $text-muted;
  transition: transform 0.3s ease;
  background: rgba(15, 23, 42, 0.04);
  border-radius: 50%;
  padding: 4px;
  flex-shrink: 0;

  svg {
    width: 100%;
    height: 100%;
  }
}

.wh-card.is-open .wh-card__chevron {
  transform: rotate(180deg);
  background: rgba(15, 23, 42, 0.08);
}

/* --- SUBMENU VÀ NÚT BẤM (XẾP DỌC) --- */
.wh-submenu-wrapper {
  overflow: hidden;
}

.wh-submenu {
  display: flex;
  flex-direction: column; 
  gap: 10px;
  padding: 0 24px 24px;
}

.wh-submenu__item {
  width: 100%; 
  min-height: 44px;

  display: flex;
  align-items: center;
  justify-content: flex-start; 
  gap: 12px;
  padding: 10px 16px;

  background: #ffffff;
  border: 1px solid rgba(15, 23, 42, 0.06);
  border-radius: 12px;
  color: $text-main;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s ease;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.02);

  &:hover {
    background: #f8fafc;
    border-color: rgba(59, 130, 246, 0.3);
    color: $blue;
    transform: translateY(-2px);
    box-shadow: 0 6px 12px rgba(59, 130, 246, 0.08);

    .submenu-arrow {
      opacity: 1;
      transform: translateX(4px);
    }
  }
}

.submenu-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.submenu-label {
  flex: 1; 
  text-align: left;
  font-weight: 600;
}

.submenu-arrow {
  display: block; 
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.25s ease, transform 0.25s ease;
}

/* --- HIỆU ỨNG ĐÓNG/MỞ CARD --- */
.expand-enter-active,
.expand-leave-active {
  transition: max-height 0.35s ease-in-out, opacity 0.3s ease, padding 0.3s ease;
  max-height: 500px;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
  padding-bottom: 0;
}

/* --- RESPONSIVE CHO MÀN HÌNH NHỎ (TABLET) --- */
@media (max-width: 900px) {
  .landing__header {
    flex-direction: column;
    gap: 20px;
  }
  .landing__header-right {
    flex-direction: row; 
    align-items: center;
    justify-content: center;
    flex-wrap: wrap;
    gap: 16px;
  }
  .landing__hero {
    margin: 30px 0 30px;
  }
}

/* --- RESPONSIVE CHO MÀN HÌNH ĐIỆN THOẠI DI ĐỘNG --- */
@media (max-width: 600px) {
  .landing {
    padding: 0 16px 40px; 
  }
  .landing__header-right {
    flex-direction: column;
    gap: 12px;
  }
  .landing__title {
    font-size: 32px !important; 
  }
  .wh-card__head {
    padding: 16px; 
    min-height: 90px;
  }
  .wh-card__name {
    font-size: 18px; 
  }
}
</style>