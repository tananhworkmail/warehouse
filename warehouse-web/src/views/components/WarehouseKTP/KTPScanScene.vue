<template>
  <div class="warehouse-scene" aria-hidden="true">
    <svg viewBox="0 0 1600 900" preserveAspectRatio="xMidYMid slice">
      <defs>
        <linearGradient id="ktp-sky" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0" stop-color="#071c19" />
          <stop offset="1" stop-color="#0b3027" />
        </linearGradient>
        <linearGradient id="ktp-floor" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0" stop-color="#163c34" />
          <stop offset="1" stop-color="#061512" />
        </linearGradient>
        <linearGradient id="ktp-carton" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0" stop-color="#ffd36f" />
          <stop offset="1" stop-color="#d98a2a" />
        </linearGradient>
        <filter id="ktp-glow" x="-80%" y="-80%" width="260%" height="260%">
          <feGaussianBlur stdDeviation="8" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
        <pattern id="ktp-grid" width="90" height="60" patternUnits="userSpaceOnUse">
          <path d="M90 0H0V60" fill="none" stroke="#58b99a" stroke-opacity=".12" stroke-width="2" />
        </pattern>
      </defs>

      <rect width="1600" height="620" fill="url(#ktp-sky)" />
      <path d="M0 610H1600V900H0Z" fill="url(#ktp-floor)" />
      <path d="M0 610H1600V900H0Z" fill="url(#ktp-grid)" />

      <g fill="#a7ffe0">
        <rect x="170" y="72" width="260" height="12" rx="6" opacity=".72" />
        <rect x="670" y="72" width="260" height="12" rx="6" opacity=".48" />
        <rect x="1170" y="72" width="260" height="12" rx="6" opacity=".72" />
      </g>

      <g>
        <path d="M45 145H410V650H45ZM45 315H410M45 480H410" fill="none" stroke="#39a58d" stroke-width="18" />
        <g fill="url(#ktp-carton)" stroke="#8f591f" stroke-width="3">
          <rect x="72" y="188" width="138" height="105" rx="5" />
          <rect x="225" y="177" width="155" height="116" rx="5" />
          <rect x="72" y="348" width="170" height="110" rx="5" />
          <rect x="257" y="360" width="123" height="98" rx="5" />
          <rect x="72" y="518" width="128" height="105" rx="5" />
          <rect x="215" y="511" width="165" height="112" rx="5" />
        </g>
        <g stroke="#ffe29a" stroke-width="5" opacity=".72">
          <path d="M141 188V293M302 177V293M157 348V458M318 360V458M136 518V623M298 511V623" />
        </g>
      </g>

      <g opacity=".64">
        <path d="M1328 145H1570V585H1328ZM1328 292H1570M1328 438H1570" fill="none" stroke="#39a58d" stroke-width="16" />
        <g fill="#bd762a" stroke="#7c4a1d" stroke-width="3">
          <rect x="1350" y="174" width="86" height="96" rx="4" />
          <rect x="1448" y="174" width="98" height="96" rx="4" />
          <rect x="1350" y="318" width="112" height="98" rx="4" />
          <rect x="1474" y="318" width="72" height="98" rx="4" />
          <rect x="1350" y="463" width="92" height="96" rx="4" />
          <rect x="1454" y="463" width="92" height="96" rx="4" />
        </g>
      </g>

      <g>
        <path d="M385 690H1490L1540 744H435Z" fill="#244d48" stroke="#52b7a0" stroke-width="8" />
        <path d="M435 744H1540V790H435Z" fill="#0b2421" stroke="#2b6f62" stroke-width="7" />
        <g fill="#9bcfc3" stroke="#163d37" stroke-width="5">
          <circle v-for="x in rollers" :key="x" :cx="x" cy="718" r="22" />
        </g>
        <path d="M520 790L485 884M1435 790L1470 884" stroke="#173d37" stroke-width="20" />
      </g>

      <g class="moving-carton">
        <rect x="1080" y="565" width="205" height="125" rx="9" fill="url(#ktp-carton)" stroke="#8f591f" stroke-width="5" />
        <path d="M1182 565V690" stroke="#ffe29a" stroke-width="8" opacity=".78" />
        <rect x="1210" y="600" width="58" height="55" rx="4" fill="#fff8dc" />
        <g stroke="#25382f" stroke-width="4">
          <path d="M1220 610V646M1228 610V646M1238 610V646M1248 610V646M1258 610V646" />
        </g>
      </g>

      <g transform="translate(1320 470) rotate(-12)">
        <path d="M115 55C170 42 210 62 238 100L196 131C169 105 143 95 104 103Z" fill="#2e6975" />
        <rect width="142" height="82" rx="23" fill="#182c32" stroke="#68d4c7" stroke-width="6" />
        <rect x="32" y="63" width="52" height="112" rx="19" fill="#132126" stroke="#477c7c" stroke-width="5" />
        <rect x="13" y="17" width="30" height="46" rx="8" fill="#ff6359" filter="url(#ktp-glow)" />
      </g>

      <g filter="url(#ktp-glow)">
        <path d="M1327 509L1090 604V636L1327 533Z" fill="#ff594f" opacity=".18">
          <animate attributeName="opacity" values=".08;.42;.08" dur="1.25s" repeatCount="indefinite" />
        </path>
        <path d="M1327 521L1090 620" stroke="#ff7065" stroke-width="6" stroke-linecap="round">
          <animate attributeName="stroke-width" values="3;8;3" dur="1.25s" repeatCount="indefinite" />
        </path>
      </g>

      <g transform="translate(1118 520)" filter="url(#ktp-glow)">
        <circle r="43" fill="#20d889" opacity=".18">
          <animate attributeName="r" values="34;53;34" dur="1.5s" repeatCount="indefinite" />
          <animate attributeName="opacity" values=".3;.05;.3" dur="1.5s" repeatCount="indefinite" />
        </circle>
        <circle r="28" fill="#0d3c30" stroke="#52f0af" stroke-width="5" />
        <path d="M-13 0L-3 11 15-12" fill="none" stroke="#8bffd0" stroke-width="7" stroke-linecap="round" stroke-linejoin="round" />
      </g>
    </svg>
  </div>
</template>

<script setup>
const rollers = [475, 570, 665, 760, 855, 950, 1045, 1140, 1235, 1330, 1425];
</script>

<style scoped>
.warehouse-scene {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.warehouse-scene svg {
  display: block;
  width: 100%;
  height: 100%;
}

.moving-carton {
  animation: carton-travel 8s linear infinite;
  will-change: transform;
}

@keyframes carton-travel {
  from { transform: translateX(150px); }
  to { transform: translateX(-760px); }
}
</style>
