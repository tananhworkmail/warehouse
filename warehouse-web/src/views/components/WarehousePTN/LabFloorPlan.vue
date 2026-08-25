<template>
  <svg
    :viewBox="`0 0 ${W} ${H}`"
    xmlns="http://www.w3.org/2000/svg"
    class="lab-floor-plan"
    preserveAspectRatio="xMidYMid meet"
    role="img"
    :aria-label="ariaLabel"
  >
    <!-- ── Tiêu đề ─────────────────────────────────────────────── -->
    <g
      v-if="showTitle"
      text-anchor="middle"
      font-family="'Times New Roman', Georgia, serif"
      font-weight="500"
    >
      <text :x="W / 2" y="62" font-size="50" fill="#151bd6">
        {{ labels.titleLine1 }}
      </text>
      <text :x="W / 2" y="132" font-size="50" fill="#e2101c">
        {{ labels.titleLine2 }}
      </text>
    </g>

    <!-- ── Lưới các phòng ──────────────────────────────────────── -->
    <g
      stroke="#000"
      stroke-width="3"
      font-family="'Times New Roman', Georgia, serif"
    >
      <!-- các ô "trống" không có nhãn (đúng theo sơ đồ gốc) -->
      <rect
        v-for="b in blankRects"
        :key="b.key"
        :x="b.x"
        :y="b.y"
        :width="b.w"
        :height="b.h"
        :fill="COLORS.white"
      />

      <!-- các phòng có nhãn -->
      <template v-for="room in rooms" :key="room.key">
        <rect
          :x="room.x"
          :y="room.y"
          :width="room.w"
          :height="room.h"
          :fill="room.fill"
        />

        <text
          text-anchor="middle"
          fill="#111827"
          stroke="none"
          :font-size="room.fontSize"
          font-weight="400"
        >
          <tspan
            v-for="(line, i) in room.lines"
            :key="i"
            :x="room.textX"
            :y="room.textY + i * room.dy"
          >{{ line }}</tspan>
        </text>

        <!-- biểu tượng cấm (no-entry) -->
        <g v-if="room.icon">
          <circle
            :cx="room.icon.cx"
            :cy="room.icon.cy"
            :r="room.icon.r"
            fill="#fff"
            stroke="#e2101c"
            stroke-width="7"
          />
          <line
            :x1="room.icon.cx - room.icon.r * 0.62"
            :y1="room.icon.cy + room.icon.r * 0.62"
            :x2="room.icon.cx + room.icon.r * 0.62"
            :y2="room.icon.cy - room.icon.r * 0.62"
            stroke="#e2101c"
            stroke-width="7"
            stroke-linecap="round"
          />
        </g>

        <!-- vùng cắm dữ liệu động (nhiệt độ / độ ẩm) -->
        <foreignObject
          v-if="room.slotName"
          :x="room.x"
          :y="room.y"
          :width="room.w"
          :height="room.h"
          style="overflow: visible"
        >
          <div style="width: 100%; height: 100%">
            <slot :name="room.slotName" :room="room" />
          </div>
        </foreignObject>
      </template>
    </g>
  </svg>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  // 'vi' | 'en' | 'zh' — chỉ ảnh hưởng tới CHỮ trong sơ đồ (hình học giữ nguyên
  // cho mọi ngôn ngữ, nên không cần đo toạ độ overlay riêng theo từng ảnh nữa)
  locale: { type: String, default: "vi" },
  showTitle: { type: Boolean, default: true },
  ariaLabel: { type: String, default: "Sơ đồ phòng lab" },
});

// ── kích thước gốc (đo trực tiếp từ ảnh vi.png 1161x1354) ─────────────
const W = 1161;
const H = 1354;

const COLORS = {
  cyan: "#b3fdfb",
  khaki: "#e6e5b4",
  pink: "#fcb4e4",
  lavender: "#ece0fd",
  white: "#ffffff",
};

// ── nhãn theo ngôn ngữ ──────────────────────────────────────────────
// Nhãn "conditionRoom"/"bondingRoom" (2 phòng có dữ liệu nhiệt độ/độ ẩm)
// vẫn nên lấy qua slot props từ component cha (dùng hệ i18n hiện có, ví dụ t()),
// còn các phòng "tĩnh" còn lại lấy từ bảng dưới đây — chỉnh lại câu chữ
// cho khớp với en.png / zh.png gốc của bạn nếu cần.
const LABELS = {
  vi: {
    titleLine1: "SƠ ĐỒ PHÒNG LAB",
    titleLine2: "PHÒNG THÍ NGHIỆM",
    chem: ["Phòng Hóa Chất", "& Nước", "(Khu vực thử nghiệm)", "5*7M"],
    storage: ["Phòng Lưu Trữ", "(Kho mẫu thử)", "17.5*5M"],
    aging: ["Phòng", "Lão Hóa", "18*3M"],
    env: ["Phòng Xử Lý", "Môi Trường", "3.9*5m"],
    stable: ["Khu Vực Ổn", "9.75*5m"],
    mcs: ["Phòng MCS"],
    pull: ["Phòng Kéo", "(Khu vực kéo)", "16*8.15M"],
    office: ["Phòng Văn Phòng", "QC+LAB+QA", "(Khu làm việc)"],
    meeting: ["Phòng Họp", "QC+LAB+QA", "(Khu văn phòng)"],
  },
  en: {
    titleLine1: "LAB ROOM LAYOUT",
    titleLine2: "LABORATORY",
    chem: ["Chemical &", "Water Room", "(Test Area)", "5*7M"],
    storage: ["Storage Room", "(Sample Warehouse)", "17.5*5M"],
    aging: ["Aging", "Room", "18*3M"],
    env: ["Environmental", "Treatment Room", "3.9*5m"],
    stable: ["Stable Zone", "9.75*5m"],
    mcs: ["MCS Room"],
    pull: ["Pull Room", "(Pull Test Area)", "16*8.15M"],
    office: ["Office Room", "QC+LAB+QA", "(Work Area)"],
    meeting: ["Meeting Room", "QC+LAB+QA", "(Office Area)"],
  },
  zh: {
    titleLine1: "實驗室平面圖",
    titleLine2: "化驗室",
    chem: ["化學藥品", "與用水室", "(測試區)", "5*7M"],
    storage: ["儲藏室", "(樣品庫)", "17.5*5M"],
    aging: ["老化室", "", "18*3M"],
    env: ["環境處理室", "", "3.9*5m"],
    stable: ["穩定區", "9.75*5m"],
    mcs: ["MCS室"],
    pull: ["拉力室", "(拉力測試區)", "16*8.15M"],
    office: ["辦公室", "QC+LAB+QA", "(工作區)"],
    meeting: ["會議室", "QC+LAB+QA", "(辦公區)"],
  },
};

const labels = computed(() => LABELS[props.locale] || LABELS.vi);

// ── các ô trống không nhãn (giữ nguyên như sơ đồ gốc) ───────────────
const blankRects = computed(() => [
  { key: "aging-blank", x: 8, y: 1194, w: 151, h: 152 },
  { key: "pull-blank", x: 159, y: 1306, w: 511, h: 40 },
  // vùng "Văn Phòng" gộp với ô trống bên dưới "Phòng Họp" (không có vách ngăn
  // ở giữa trong ảnh gốc)
  { key: "office-base", x: 670, y: 649, w: 483, h: 697 },
]);

// ── các phòng có nhãn (toạ độ đo trực tiếp từ ảnh gốc 1161x1354) ────
const rooms = computed(() => {
  const l = labels.value;
  return [
    {
      key: "chem",
      x: 8,
      y: 162,
      w: 301,
      h: 215,
      fill: COLORS.cyan,
      lines: l.chem,
      textX: 158.5,
      textY: 228,
      dy: 25,
      fontSize: 20,
    },
    {
      key: "storage",
      x: 309,
      y: 162,
      w: 844,
      h: 215,
      fill: COLORS.white,
      lines: l.storage,
      textX: 731,
      textY: 236,
      dy: 34,
      fontSize: 27,
    },
    {
      key: "aging",
      x: 8,
      y: 377,
      w: 151,
      h: 817,
      fill: COLORS.pink,
      lines: l.aging,
      textX: 83.5,
      textY: 760,
      dy: 26,
      fontSize: 20,
    },
    {
      key: "env",
      x: 159,
      y: 377,
      w: 237,
      h: 272,
      fill: COLORS.lavender,
      lines: l.env,
      textX: 277.5,
      textY: 430,
      dy: 25,
      fontSize: 20,
      icon: { cx: 280, cy: 540, r: 26 },
      slotName: "condition-room",
    },
    {
      key: "stable",
      x: 396,
      y: 377,
      w: 560,
      h: 272,
      fill: COLORS.khaki,
      lines: l.stable,
      textX: 676,
      textY: 497,
      dy: 32,
      fontSize: 25,
    },
    {
      key: "mcs",
      x: 956,
      y: 377,
      w: 197,
      h: 272,
      fill: COLORS.white,
      lines: l.mcs,
      textX: 1054.5,
      textY: 521,
      dy: 28,
      fontSize: 24,
    },
    {
      key: "pull",
      x: 159,
      y: 649,
      w: 511,
      h: 657,
      fill: COLORS.cyan,
      lines: l.pull,
      textX: 414.5,
      textY: 836,
      dy: 32,
      fontSize: 25,
      icon: { cx: 414, cy: 999, r: 28 },
      slotName: "bonding-room",
    },
    {
      key: "meeting",
      x: 956,
      y: 649,
      w: 197,
      h: 409,
      fill: COLORS.white,
      lines: l.meeting,
      textX: 1054.5,
      textY: 826,
      dy: 28,
      fontSize: 22,
    },
    {
      key: "office",
      // không vẽ rect riêng (đã có trong blankRects "office-base"),
      // chỉ vẽ chữ đè lên
      x: 670,
      y: 649,
      w: 0,
      h: 0,
      fill: "none",
      lines: l.office,
      textX: 810,
      textY: 700,
      dy: 29,
      fontSize: 22,
    },
  ];
});
</script>

<style scoped>
.lab-floor-plan {
  display: block;
  width: 100%;
  height: 100%;
  user-select: none;
}
</style>