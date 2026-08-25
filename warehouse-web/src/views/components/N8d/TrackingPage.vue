<template>
  <N8dHeader :back-to="{ name: 'n8d-home' }">
    <template #context>
      <span class="tracking-page-chip">{{ pageTitle }}</span>
    </template>

    <div class="tracking-page">
      <main class="tracking-shell">
        <SearchToolBar />

        <section class="tracking-results-region" aria-label="Tracking results">
          <TrackingLeftTable />
        </section>
      </main>

      <div v-if="print_file_layout_show" class="print-overlay">
        <PrintFileLayoutImport v-if="file_type === 'import'" />
        <PrintFileLayoutExport v-else />
      </div>
    </div>
  </N8dHeader>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import N8dHeader from '@/views/components/N8d/N8dHeader.vue'
import SearchToolBar from '@/components/TrackingPage/SearchToolBar.vue'
import TrackingLeftTable from '@/components/TrackingPage/TrackingLeftTable.vue'
import PrintFileLayoutImport from '@/components/TrackingPage/PrintFileLayoutImport.vue'
import PrintFileLayoutExport from '@/components/TrackingPage/PrintFileLayoutExport.vue'
import useTracking from '@/hooks/useTracking'

const { t } = useI18n()
const { print_file_layout_show, file_type } = useTracking()

const pageTitle = computed(() => t('warehouse.n8d.home.tracking'))
</script>
<style scoped>
/* Chip tiêu đề trang, hiển thị trong slot "context" của N8dHeader */
.tracking-page-chip {
  padding: 6px 14px;
  border: 1px solid rgba(255, 255, 255, 0.28);
  border-radius: 999px;
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0.02em;
  white-space: nowrap;
}

/* ---------------------------------------------------------
   ĐỒNG BỘ DESIGN SYSTEM VỚI TRANG HOME
   --------------------------------------------------------- */
.tracking-page {
  /* Sử dụng lại biến màu chung */
  --ink: #1E293B;
  --ink-soft: #475569;
  --muted: #94A3B8;
  --paper: #F1F5F9;
  --surface: #FFFFFF;
  --line: #E2E8F0;
  --line-strong: #CBD5E1;
  
  /* Biến局部 của trang này map sang biến chung */
  --tracking-ink: var(--ink);
  --tracking-muted: var(--muted);
  --tracking-line: var(--line);

  position: relative;
  overflow-y: auto; /* Cho phép cuộn dọc toàn trang khi bảng cao hơn viewport (thay vì hidden) */
  overflow-x: hidden;
  width: 100%;
  height: 100%; /* Lấp đầy khu vực nội dung của N8dHeader thay vì toàn viewport */
  color: var(--tracking-ink);
  background: var(--paper);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Inter", Roboto, "Helvetica Neue", Arial, sans-serif;
}

/* Hiệu ứng Ambient Glow nền (Giống hệt trang Home) */
.tracking-page::before {
  content: "";
  position: absolute;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse 80% 50% at 10% 0%, rgba(37, 99, 235, 0.12), transparent 50%),
    radial-gradient(ellipse 60% 40% at 90% 10%, rgba(5, 150, 105, 0.08), transparent 40%);
  mask-image: radial-gradient(ellipse 100% 80% at 50% 0%, black 20%, transparent 70%);
  -webkit-mask-image: radial-gradient(ellipse 100% 80% at 50% 0%, black 20%, transparent 70%);
}

/* ---------------------------------------------------------
   FIX FULL MÀN HÌNH & GRID LAYOUT
   --------------------------------------------------------- */
.tracking-shell {
  position: relative;
  z-index: 1; /* Nằm trên lớp nền glow */
  display: grid;
  /* Toolbar và khu vực bảng đều auto: chiều cao co giãn theo nội dung (số dòng dữ liệu) thay vì ép 1fr full viewport */
  grid-template-rows: auto auto;
  gap: 16px;
  width: 100%;
  min-height: 100%; /* Vẫn lấp đầy tối thiểu khu vực nội dung khi dữ liệu ít */
  margin: 0 auto;
  padding: 16px 24px 24px 24px;
}

/* ---------------------------------------------------------
   KHU VỰC BẢNG BIỂU (NỔI BẬT GIỐNG CARD HOME)
   --------------------------------------------------------- */
.tracking-results-region {
  position: relative;
  min-width: 0;
  min-height: 0; /* QUAN TRỌNG: Cho phép nội dung bên trong (table) cuộn linh hoạt */
  overflow: hidden; /* Không cho tràn ra ngoài bo góc */
  border-radius: 20px; /* Bo góc đồng bộ với Home */
  background: var(--surface);
  border: 1px solid transparent; /* Ẩn viền mặc định */
  
  /* Bóng đổ Multi-layer Pro (Giống Card Home) */
  box-shadow: 
    0 0 0 1px rgba(15, 23, 42, 0.03),
    0 2px 4px -1px rgba(15, 23, 42, 0.04),
    0 8px 16px -4px rgba(15, 23, 42, 0.06),
    0 24px 48px -12px rgba(15, 23, 42, 0.12);
}

/* Bảng dữ liệu bên trong co giãn theo nội dung, không ép full height của thẻ nữa */
.tracking-results-region > :deep(.tracking-results) {
  height: auto;
  min-height: 0; /* Ngăn table bị thrust */
}

/* ---------------------------------------------------------
   POPUP IN ẤN
   --------------------------------------------------------- */
.print-overlay {
  position: fixed;
  z-index: 9999;
  inset: 0;
  overflow: auto;
  color: #111827;
  background: #fff;
}

/* ---------------------------------------------------------
   RESPONSIVE
   --------------------------------------------------------- */
@media (max-width: 768px) {
  .tracking-shell {
    display: flex;
    flex-direction: column;
    width: 100%;
    height: auto;
    min-height: 100dvh;
    padding: 12px;
    gap: 12px;
  }

  .tracking-results-region {
    flex: 1;
    min-height: 400px; /* Chiều cao tối thiểu cho mobile */
    border-radius: 16px;
  }
}
</style>