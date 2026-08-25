<!-- src/components/Inventory/ScanToolBar2.vue -->
<template>
    <div class="navbar">
        <div class="navbar-left">
            <span v-if="title" class="page-title">{{ title }}</span>
        </div>

        <div class="navbar-right">
            <!-- Language switcher -->
            <div class="lang-switcher-wrap">
                <button
                    type="button"
                    class="lang-switcher hoverable-icon"
                    @click="langMenuOpen = !langMenuOpen"
                >
                    <svg
                        class="lang-icon-svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="1.6" />
                        <ellipse cx="12" cy="12" rx="4" ry="9" stroke="currentColor" stroke-width="1.6" />
                        <line x1="3" y1="12" x2="21" y2="12" stroke="currentColor" stroke-width="1.6" />
                        <path d="M5 7.5C6.8 8.7 9.3 9.4 12 9.4C14.7 9.4 17.2 8.7 19 7.5" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
                        <path d="M5 16.5C6.8 15.3 9.3 14.6 12 14.6C14.7 14.6 17.2 15.3 19 16.5" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
                    </svg>
                    <span class="lang-code">{{ currentLangLabel }}</span>
                </button>

                <ul v-if="langMenuOpen" v-click-outside="closeLangMenu" class="lang-list">
                    <li
                        v-for="lang in languages"
                        :key="lang.code"
                        class="lang-item"
                        :class="{ 'lang-item--active': locale === lang.code }"
                        @click="changeLanguage(lang.code)"
                    >
                        <span class="lang-flag">{{ lang.flag }}</span>
                        <span class="lang-item-label">{{ lang.label }}</span>
                        <span v-if="locale === lang.code" class="lang-check">✓</span>
                    </li>
                </ul>
            </div>

            <!-- Fullscreen toggle -->
            <button
                type="button"
                class="hoverable-icon"
                :title="isFullscreen ? 'Exit fullscreen' : 'Enter fullscreen'"
                @click="isFullscreen ? exitFullscreen() : enterFullscreen()"
            >
                <svg v-if="!isFullscreen" viewBox="0 0 24 24" width="26" height="26" fill="none" stroke="currentColor" stroke-width="1.8">
                    <path d="M8 3H3v5M16 3h5v5M8 21H3v-5M16 21h5v-5" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
                <svg v-else viewBox="0 0 24 24" width="26" height="26" fill="none" stroke="currentColor" stroke-width="1.8">
                    <path d="M3 8h5V3M21 8h-5V3M3 16h5v5M21 16h-5v5" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
            </button>

            <!-- Exit -->
            <button type="button" class="hoverable-icon exit-btn" title="Exit" @click="router.push('/n8d')">
                <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="1.8">
                    <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import useScreen from '@/hooks/useScreen'

// ⚠️ Trước đây thiếu defineProps nên :title="pageTitle" truyền từ Inventory.vue
// bị rớt (chỉ thành thuộc tính title HTML mặc định trên <div>, không hiển thị).
const props = defineProps({
    title: {
        type: String,
        default: '',
    },
})

const router = useRouter()
const { isFullscreen, enterFullscreen, exitFullscreen } = useScreen()

// ─────────────────────────────────────────────
// Language switcher (dropdown thuần, không phụ thuộc thư viện UI)
// ─────────────────────────────────────────────
const { locale } = useI18n()

const languages = [
    { code: 'vi', label: 'Tiếng Việt', flag: '🇻🇳' },
    { code: 'en', label: 'English',    flag: '🇬🇧' },
    { code: 'zh', label: '中文',        flag: '🇨🇳' },
]

const currentLangLabel = computed(() => locale.value.toUpperCase())
const langMenuOpen = ref(false)

const closeLangMenu = () => {
    langMenuOpen.value = false
}

const changeLanguage = (code) => {
    langMenuOpen.value = false
    if (locale.value === code) return
    locale.value = code
    // Lưu lại lựa chọn để giữ ngôn ngữ khi F5 / chuyển trang / vào lại Main
    localStorage.setItem('APP_LOCALE', code)
}

// Directive nhỏ gọn để đóng dropdown khi click ra ngoài (không cần thư viện ngoài).
// Vue 3 tự động đăng ký biến đặt tên vXxx trong <script setup> thành directive v-xxx trong template.
const vClickOutside = {
    mounted(el, binding) {
        el.__clickOutsideHandler__ = (event) => {
            if (!(el === event.target || el.contains(event.target))) {
                binding.value(event)
            }
        }
        document.addEventListener('click', el.__clickOutsideHandler__, true)
    },
    unmounted(el) {
        document.removeEventListener('click', el.__clickOutsideHandler__, true)
    },
}
</script>

<style scoped>
.navbar {
    position: relative;
    background: linear-gradient(to bottom, #000080, #4169e1);
    min-width: 100%;
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    border-top: 1px solid white;
    border-bottom: 1px solid white;
}

.navbar-left {
    display: flex;
    align-items: center;
    flex: 1;
}

.navbar-right {
    display: flex;
    align-items: center;
    gap: 16px;
}

.hoverable-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: transparent;
    padding: 0;
    color: white;
    cursor: pointer;
    transition: transform 0.2s, color 0.2s;
}

.hoverable-icon:hover {
    color: lightgreen;
    transform: scale(1.1);
}

.exit-btn:hover {
    color: #ff5252;
}

/* Language switcher */
.lang-switcher-wrap {
    position: relative;
}

.lang-switcher {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 16px;
    font-weight: bold;
    color: white;
}

.lang-icon-svg {
    display: block;
    width: 24px;
    height: 24px;
    color: currentColor;
    flex-shrink: 0;
}

.lang-code {
    letter-spacing: 0.5px;
}

.lang-list {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    z-index: 30;
    min-width: 170px;
    margin: 0;
    padding: 6px;
    list-style: none;
    background: rgba(24, 24, 24, 0.97);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 10px;
    box-shadow: 0 12px 30px rgba(0, 0, 0, 0.4);
}

.lang-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 10px;
    border-radius: 8px;
    color: white;
    font-size: 14px;
    cursor: pointer;
    transition: background 0.15s ease;
}

.lang-item:hover {
    background: rgba(255, 255, 255, 0.1);
}

.lang-item--active {
    background: rgba(76, 175, 80, 0.25);
}

.lang-flag {
    font-size: 16px;
}

.lang-item-label {
    flex: 1;
}

.lang-check {
    color: lightgreen;
    font-weight: 700;
}

.page-title {
    color: white;
    font-size: 18px;
    font-weight: 700;
    letter-spacing: 0.02em;
}

@media (max-width: 600px) {
    .navbar {
        padding: 0 12px;
    }

    .navbar-right {
        gap: 10px;
    }

    .page-title {
        font-size: 15px;
    }
}
</style>