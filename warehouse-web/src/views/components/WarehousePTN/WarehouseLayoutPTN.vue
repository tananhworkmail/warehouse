<template>
  <div class="mes-container">
    <el-container>
      <el-container class="container">
        <el-header class="header fixed" height="70px">
          <img
            src="@/assets/Logo.png"
            alt="LAIYIH"
            style="height: 60px; margin-left: 20px; margin-top: 5px"
          />
          <h1 class="currentTime">{{ currentTime }}</h1>
          <div class="action">
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

            <el-dropdown @command="handleCommand">
              <span class="dropdown-trigger">
                {{ t("common.menu") }}
                <el-icon :size="24"><Menu /></el-icon>
              </span>

              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="goreTexForms">
                    <div class="dropdown-back">
                      <el-icon :size="22"><Document /></el-icon>
                      <span>{{ t("laboratory.forms") }}</span>
                    </div>
                  </el-dropdown-item>
                  <el-dropdown-item command="back" divided>
                    <div class="dropdown-back">
                      <img
                        src="@/assets/undo.png"
                        :alt="t('common.back')"
                        class="undo-icon"
                      />
                      <span>{{ t("common.back") }}</span>
                    </div>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <!-- Nút phóng to / thu nhỏ toàn màn hình -->
            <div
              class="fullscreen-btn"
              @click="toggleFullscreen"
              :title="
                isFullscreen
                  ? t('common.exitFullscreen')
                  : t('common.fullscreen')
              "
            >
              <el-icon :size="24" v-if="!isFullscreen"><FullScreen /></el-icon>
              <el-icon :size="24" v-else><CloseBold /></el-icon>
            </div>
          </div>
        </el-header>

        <el-main class="main">
          <div class="app-main-container">
            <slot></slot>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
export default {
  name: "WarehouseLayoutKVT",
};
</script>

<script setup>
import {
  ref,
  reactive,
  computed,
  onBeforeMount,
  onUnmounted,
  onMounted,
} from "vue";
import { useRouter, useRoute } from "vue-router";

import {
  Search,
  Menu,
  Bell,
  Document,
  FullScreen,
  CloseBold,
} from "@element-plus/icons-vue";
import axios from "axios";
import dayjs from "dayjs";
import { useWarehouseMapI18n } from "@/hooks/useWarehouseMapI18n";

const router = useRouter();
const route = useRoute();
const API_URL = import.meta.env.VITE_API_URL;
const { lang, langOptions, setLang, t } = useWarehouseMapI18n();

// Fullscreen
const isFullscreen = ref(false);
const handleCommand = (command) => {
  if (command === "goreTexForms") {
    router.push("/laboratory/forms");
  } else if (command === "back") {
    goBack();
  }
};
const goBack = () => {
  router.push({ name: "home" });
};

let evt = null;
let timer = null;

const updateFullscreenStatus = () => {
  isFullscreen.value = !!(
    document.fullscreenElement ||
    document.webkitFullscreenElement ||
    document.msFullscreenElement
  );
};

const toggleFullscreen = () => {
  const elem = document.documentElement;
  if (!isFullscreen.value) {
    if (elem.requestFullscreen) elem.requestFullscreen();
    else if (elem.webkitRequestFullscreen) elem.webkitRequestFullscreen();
    else if (elem.msRequestFullscreen) elem.msRequestFullscreen();
  } else {
    if (document.exitFullscreen) document.exitFullscreen();
    else if (document.webkitExitFullscreen) document.webkitExitFullscreen();
    else if (document.msExitFullscreen) document.msExitFullscreen();
  }
};

onMounted(() => {
  // Lắng nghe fullscreen change (kể cả F11, Esc)
  document.addEventListener("fullscreenchange", updateFullscreenStatus);
  document.addEventListener("webkitfullscreenchange", updateFullscreenStatus);
  document.addEventListener("msfullscreenchange", updateFullscreenStatus);
  updateFullscreenStatus();
});

onUnmounted(() => {
  if (evt) evt.close();
  if (timer) clearInterval(timer);
  document.removeEventListener("fullscreenchange", updateFullscreenStatus);
  document.removeEventListener(
    "webkitfullscreenchange",
    updateFullscreenStatus,
  );
  document.removeEventListener("msfullscreenchange", updateFullscreenStatus);
});

const currentTime = ref(dayjs().format("YYYY/MM/DD HH:mm:ss"));

const syncCurrentTime = () => {
  timer = setInterval(() => {
    currentTime.value = dayjs().format("YYYY/MM/DD HH:mm:ss");
  }, 1000);
};

onBeforeMount(() => {
  syncCurrentTime();
});
</script>

<style lang="scss" scoped>
$base-z-index-999: 999;

.fullscreen-btn {
  cursor: pointer;
  color: white;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  transition: color 0.2s;
}

.fullscreen-btn:hover {
  color: #409eff;
}

.search-toolbar {
  position: sticky;
  top: 0;
  background: white;
  z-index: 10;
  padding-bottom: 8px;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;
  border-bottom: 1px solid #e4e7ed;
}

.search-result-table {
  max-height: 70vh;
  overflow-y: auto;
}

.el-table {
  user-select: text;
}

.dropdown-back {
  display: flex;
  align-items: center;
  gap: 8px;
}

.undo-icon {
  width: 22px;
  height: 22px;
  object-fit: contain;
}

:deep(.el-dropdown-menu__item:hover .undo-icon) {
  transform: scale(1.08);
  transition: transform 0.2s ease;
}

.action {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  gap: 20px;
}

.search {
  cursor: pointer;
  color: white;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.alert-btn {
  cursor: pointer;
  color: white;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.lang-switch {
  display: flex;
  gap: 4px;
  padding: 3px;
  border: 1px solid rgba(255, 255, 255, 0.28);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.1);
}

.lang-btn {
  min-width: 38px;
  height: 28px;
  padding: 0 10px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: rgba(255, 255, 255, 0.78);
  font-size: 12px;
  font-weight: 800;
  cursor: pointer;
}

.lang-btn.is-active {
  background: #fff;
  color: #241f20;
}

.dropdown-trigger {
  cursor: pointer;
  color: white;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.dropdown-trigger:hover {
  color: #409eff;
}

.summary-box {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
  padding: 4px 8px;
  background-color: #f9f9f9;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  min-height: 30px;
  align-items: center;
}

.summary-item {
  display: flex;
  align-items: center;
  min-width: 120px;
}

.label {
  font-weight: bold;
  color: #606266;
  margin-right: 4px;
}

.value {
  color: #409eff;
  font-weight: bold;
}

.mes-container {
  user-select: none;

  .header {
    padding: 0;
    background-color: #241f20;

    .currentTime {
      position: absolute;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      color: white;
      font-weight: 900;
      text-align: center;
    }

    .menu {
      position: absolute;
      top: 5px;
      right: 20px;

      h2 {
        color: #fff;
        font-size: 18px;
        font-weight: 800;
      }
    }

    &.fixed {
      position: relative;
      width: 100%;
      z-index: $base-z-index-999;
    }
  }

  .main {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 70px);
    padding: 0;

    .app-main-container {
      flex: 1;
      overflow: auto;
      box-sizing: border-box;
      width: 100%;
      text-align: left;

      .app-main-height {
        height: 100%;
      }

      .footer-copyright {
        height: 40px;
        line-height: 40px;
        flex-shrink: 0;
        color: #999;
        text-align: center;
        border-top: 1px dashed #ccc;
      }
    }
  }
}

:deep(.el-table thead th) {
  background-color: #3b82f6 !important;
  color: white !important;
  font-weight: bold;
  text-align: center;
}

@keyframes shake {
  0% {
    transform: rotate(0deg);
  }
  25% {
    transform: rotate(15deg);
  }
  50% {
    transform: rotate(-15deg);
  }
  75% {
    transform: rotate(10deg);
  }
  100% {
    transform: rotate(0deg);
  }
}

:deep(.shaking) {
  animation: shake 0.5s infinite;
  color: #f56c6c;
}
</style>
