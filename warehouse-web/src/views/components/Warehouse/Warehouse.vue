<template>
  <WarehouseLayout>
    <template #default>
      <el-card class="box-card">
        <!-- KhoDe -->
        <template v-if="state.warehouseSelect === 'KhoDe'">
          <div class="map-container">
            <div class="map-inner">
              <img
                class="responsive-map-image"
                src="@/assets/images/KhoDe.jpg"
                usemap="#image-map"
              />
              <div
                v-for="rackId in mapAreas"
                :key="rackId"
                v-show="getTonKhoByRack(rackId).total > 0"
                class="overlay-text"
                :class="{
                  'has-bg': [
                    'A',
                    'B',
                    'C',
                    'D',
                    'E',
                    'F',
                    'G',
                    'H',
                    'I',
                  ].includes(rackId),
                  zero: getTonKhoByRack(rackId).total === 0,
                }"
                :style="getOverlayStyle(rackId)"
              >
                <!-- MAIN RACK -->
                <template
                  v-if="
                    [
                      'A',
                      'B',
                      'C',
                      'D',
                      'E',
                      'F',
                      'G',
                      'H',
                      'I',
                      'K',
                      'L',
                    ].includes(rackId)
                  "
                >
                  <div class="fraction-box">
                    <div class="fraction-top">
                      {{ getTonKhoByRack(rackId).total }}
                    </div>

                    <div
                      v-if="getTonKhoByRack(rackId).limit > 0"
                      class="fraction-divider"
                    ></div>

                    <div
                      v-if="getTonKhoByRack(rackId).limit > 0"
                      class="fraction-bottom"
                    >
                      {{ getTonKhoByRack(rackId).limit }}
                    </div>
                  </div>
                </template>

                <!-- SUB RACK -->
                <template v-else>
                  {{ getTonKhoByRack(rackId).total }}
                </template>
              </div>
              <template v-for="rackId in mapAreas" :key="'sub-' + rackId">
                <div
                  v-if="rackId.startsWith('Sub')"
                  class="sub-progress-overlay"
                  :style="getSubOverlayStyle(rackId)"
                >
                  <div
                    class="progress-bar"
                    :class="{
                      full: getRackPercent(rackId.replace(/^Sub/i, '')) === 100,
                    }"
                  >
                    <div
                      class="progress-used"
                      :style="{
                        height:
                          getRackPercent(rackId.replace(/^Sub/i, '')) + '%',
                      }"
                    ></div>
                  </div>
                </div>
              </template>

              <map name="image-map">
                <area
                  id="A"
                  title="A"
                  href="#"
                  coords="3344,28,3479,90"
                  shape="rect"
                />
                <area
                  id="B"
                  title="B"
                  href="#"
                  coords="3024,28,3159,90"
                  shape="rect"
                />
                <area
                  id="C"
                  title="C"
                  href="#"
                  coords="2704,28,2839,90"
                  shape="rect"
                />
                <area
                  id="D"
                  title="D"
                  href="#"
                  coords="2384,28,2519,90"
                  shape="rect"
                />
                <area
                  id="E"
                  title="E"
                  href="#"
                  coords="2064,28,2199,90"
                  shape="rect"
                />
                <area
                  id="F"
                  title="F"
                  href="#"
                  coords="1744,28,1879,90"
                  shape="rect"
                />
                <area
                  id="G"
                  title="G"
                  href="#"
                  coords="1424,28,1559,90"
                  shape="rect"
                />
                <area
                  id="H"
                  title="H"
                  href="#"
                  coords="1104,28,1239,90"
                  shape="rect"
                />
                <area
                  id="I"
                  title="I"
                  href="#"
                  coords="784,28,919,90"
                  shape="rect"
                />
                <area id="SubA" coords="3309,106,3330,1406" shape="rect" />
                <area id="SubB" coords="2989,106,3010,1406" shape="rect" />
                <area id="SubC" coords="2664,106,2685,1406" shape="rect" />
                <area id="SubD" coords="2344,106,2365,1406" shape="rect" />
                <area id="SubE" coords="2024,106,2045,1406" shape="rect" />
                <area id="SubF" coords="1704,106,1725,1406" shape="rect" />
                <area id="SubG" coords="1384,106,1405,1406" shape="rect" />
                <area id="SubH" coords="1064,106,1085,1406" shape="rect" />
                <area id="SubI" coords="744,106,765,1406" shape="rect" />
                <!-- A -->
                <!-- Khu vực A (X: 3345 -> 3480) -->
                <area
                  id="A01"
                  title="A01"
                  href="#"
                  coords="3345,1145,3480,1395"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A02"
                  title="A02"
                  href="#"
                  coords="3345,885,3480,1135"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A03"
                  title="A03"
                  href="#"
                  coords="3345,625,3480,875"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A04"
                  title="A04"
                  href="#"
                  coords="3345,365,3480,615"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A05"
                  title="A05"
                  href="#"
                  coords="3345,105,3480,355"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- Khu vực B (X: 3025 -> 3160) -->
                <area
                  id="B01"
                  title="B01"
                  href="#"
                  coords="3025,1145,3160,1395"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B02"
                  title="B02"
                  href="#"
                  coords="3025,885,3160,1135"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B03"
                  title="B03"
                  href="#"
                  coords="3025,625,3160,875"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B04"
                  title="B04"
                  href="#"
                  coords="3025,365,3160,615"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B05"
                  title="B05"
                  href="#"
                  coords="3025,105,3160,355"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <!-- Khu vực C (X: 2705 -> 2840) -->
                <area
                  id="C01"
                  title="C01"
                  href="#"
                  coords="2705,1145,2840,1395"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="C02"
                  title="C02"
                  href="#"
                  coords="2705,885,2840,1135"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="C03"
                  title="C03"
                  href="#"
                  coords="2705,625,2840,875"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="C04"
                  title="C04"
                  href="#"
                  coords="2705,365,2840,615"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="C05"
                  title="C05"
                  href="#"
                  coords="2705,105,2840,355"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- Khu vực D (X: 2385 -> 2520) -->
                <area
                  id="D01"
                  title="D01"
                  href="#"
                  coords="2385,1145,2520,1395"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="D02"
                  title="D02"
                  href="#"
                  coords="2385,885,2520,1135"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="D03"
                  title="D03"
                  href="#"
                  coords="2385,625,2520,875"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="D04"
                  title="D04"
                  href="#"
                  coords="2385,365,2520,615"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="D05"
                  title="D05"
                  href="#"
                  coords="2385,105,2520,355"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- Khu vực E (X: 2065 -> 2200) -->
                <area
                  id="E01"
                  title="E01"
                  href="#"
                  coords="2065,1145,2200,1395"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="E02"
                  title="E02"
                  href="#"
                  coords="2065,885,2200,1135"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="E03"
                  title="E03"
                  href="#"
                  coords="2065,625,2200,875"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="E04"
                  title="E04"
                  href="#"
                  coords="2065,365,2200,615"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="E05"
                  title="E05"
                  href="#"
                  coords="2065,105,2200,355"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- Khu vực F (X: 1745 -> 1880) -->
                <area
                  id="F01"
                  title="F01"
                  href="#"
                  coords="1745,1145,1880,1395"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="F02"
                  title="F02"
                  href="#"
                  coords="1745,885,1880,1135"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="F03"
                  title="F03"
                  href="#"
                  coords="1745,625,1880,875"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="F04"
                  title="F04"
                  href="#"
                  coords="1745,365,1880,615"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="F05"
                  title="F05"
                  href="#"
                  coords="1745,105,1880,355"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- Khu vực G (X: 1425 -> 1560) -->
                <area
                  id="G01"
                  title="G01"
                  href="#"
                  coords="1425,1145,1560,1395"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="G02"
                  title="G02"
                  href="#"
                  coords="1425,885,1560,1135"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="G03"
                  title="G03"
                  href="#"
                  coords="1425,625,1560,875"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="G04"
                  title="G04"
                  href="#"
                  coords="1425,365,1560,615"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="G05"
                  title="G05"
                  href="#"
                  coords="1425,105,1560,355"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- Khu vực H (X: 1105 -> 1240) -->
                <area
                  id="H01"
                  title="H01"
                  href="#"
                  coords="1105,1145,1240,1395"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H02"
                  title="H02"
                  href="#"
                  coords="1105,885,1240,1135"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H03"
                  title="H03"
                  href="#"
                  coords="1105,625,1240,875"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H04"
                  title="H04"
                  href="#"
                  coords="1105,365,1240,615"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H05"
                  title="H05"
                  href="#"
                  coords="1105,105,1240,355"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- Khu vực I (X: 785 -> 920) -->
                <area
                  id="I01"
                  title="I01"
                  href="#"
                  coords="785,1145,920,1395"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="I02"
                  title="I02"
                  href="#"
                  coords="785,885,920,1135"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="I03"
                  title="I03"
                  href="#"
                  coords="785,625,920,875"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="I04"
                  title="I04"
                  href="#"
                  coords="785,365,920,615"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="I05"
                  title="I05"
                  href="#"
                  coords="785,105,920,355"
                  shape="rect"
                  @click="handleAreaClick"
                />
              </map>
            </div>
          </div>
        </template>

        <!-- A01-2F -->
        <template v-else> </template>
      </el-card>

      <RackInformationDialog
        v-if="state.rackInformationDialogVisible"
        :rackNo="state.rackNo"
        :dialogVisible="state.rackInformationDialogVisible"
        @close="state.rackInformationDialogVisible = false"
      />
    </template>
  </WarehouseLayout>
</template>

<script>
export default {
  name: "WarehouseIndex",
};
</script>

<script setup>
import { reactive, onMounted, nextTick, watch, ref, onUnmounted } from "vue";
import imageMapResize from "image-map-resizer";
import axios from "axios";
import WarehouseLayout from "@/views/components/Warehouse/WarehouseLayout.vue";
import RackInformationDialog from "../Warehouse/RackInformationDialog/index.vue";
const tonKhoList = ref([]);
const tonKhoTangList = ref([]);
const API_URL = import.meta.env.VITE_API_URL;
const state = reactive({
  warehouseSelect: "KhoDe",
  warehouseOptions: [
    { value: "KhoDe", label: "KhoDe" },
    { value: "A01-2F", label: "A01-2F" },
  ],
  rackInformationDialogVisible: false,
  rackNo: "",
  rackList: [],
  dialogRackListVisible: false,
});
const fetchTonKho = async () => {
  try {
    const res = await axios.get(`${API_URL}/warehouse/tonkhoke`);
    tonKhoList.value = Array.isArray(res.data.data) ? res.data.data : [];

    console.log("tonKhoList loaded:", tonKhoList.value);
  } catch (err) {
    console.error("Lỗi khi gọi API tồn kho:", err);
    tonKhoList.value = [];
  }
};
const getRackPercent = (rackId) => {
  const rack = getTonKhoByRack(rackId);

  if (!rack.limit || rack.limit === 0) return 0;

  return Math.min(Math.round((rack.total / rack.limit) * 100), 100);
};

const getSubOverlayStyle = (rackId) => {
  const areaElement = document.querySelector(`area[id="${rackId}"]`);
  if (!areaElement) return {};

  const coords = areaElement.coords.split(",").map(Number);
  const [x1, y1, x2, y2] = coords;

  return {
    position: "absolute",
    left: `${x1}px`,
    top: `${y1}px`,
    width: `${x2 - x1}px`,
    height: `${y2 - y1}px`,
    zIndex: 1,
  };
};
// const getTonKhoByRack = (rackId) => {
//   const isMainRack = [
//     "A",
//     "B",
//     "C",
//     "D",
//     "E",
//     "F",
//     "G",
//     "H",
//     "I",

//   ].includes(rackId);

//   if (isMainRack) {
//     const list = tonKhoList.value;
//     if (!Array.isArray(list)) return "";
//     const found = list.find((item) => item.MA_KE === rackId);
//     return found ? found.TonKhoKe : 0;
//   } else {
//     const list = tonKhoTangList.value;
//     if (!Array.isArray(list)) return "";
//     const found = list.find((item) => item.MA_KE === rackId);
//     return found ? found.TonKhoKe : 0;
//   }
// };
const getTonKhoByRack = (rackId) => {
  const mainRacks = ["A", "B", "C", "D", "E", "F", "G", "H", "I"];

  // MAIN RACK
  if (mainRacks.includes(rackId)) {
    const found = tonKhoList.value.find((item) => item.MA_KE === rackId);

    return {
      total: found?.TonKhoKe || 0,
      limit: found?.Limit || 0,
    };
  }

  // SUB RACK
  const foundSub = tonKhoTangList.value.find((item) => item.MA_KE === rackId);

  // lấy rack cha (SubA → A)
  const parentRack = rackId.replace(/^Sub/i, "");
  const parent = tonKhoList.value.find((item) => item.MA_KE === parentRack);

  return {
    total: foundSub?.TonKhoKe || 0,
    limit: parent?.Limit || 0, // ✅ dùng limit của rack cha
  };
};
const fetchTonKhoTang = async () => {
  try {
    const res = await axios.get(`${API_URL}/warehouse/tonkhotang`);
    tonKhoTangList.value = Array.isArray(res.data.data) ? res.data.data : [];
    console.log("Tồn kho tầng:", tonKhoTangList.value);
  } catch (err) {
    console.error("Lỗi khi gọi API tồn kho tầng:", err);
    tonKhoTangList.value = [];
  }
};

const handleAreaClick = (event) => {
  const rackCode = event.target.id; // A01, B03, ...

  state.rackNo = rackCode;
  state.rackInformationDialogVisible = true;
};
const getOverlayStyle = (rackId) => {
  const areaElement = document.querySelector(`area[id="${rackId}"]`);
  if (!areaElement) return {};
  const coords = areaElement.coords.split(",").map(Number);
  const [x1, y1, x2, y2] = coords;
  const left = (x1 + x2) / 2;
  const top = (y1 + y2) / 2;
  return {
    position: "absolute",
    left: `${left}px`,
    top: `${top}px`,
    transform: "translate(-50%, -50%)",
  };
};

const mapAreas = ref([]);

let resizeTimer;
const updateOverlay = () => {
  clearTimeout(resizeTimer);
  resizeTimer = setTimeout(() => {
    nextTick(() => {
      imageMapResize();
      const areaIds = [...document.querySelectorAll("area")].map((a) => a.id);
      mapAreas.value = areaIds;
      tonKhoList.value = [...tonKhoList.value];
    });
  }, 200);
};

onMounted(async () => {
  await fetchTonKho();
  await fetchTonKhoTang();

  nextTick(() => {
    imageMapResize();
    const areaIds = [...document.querySelectorAll("area")].map((a) => a.id);
    mapAreas.value = areaIds;
  });

  updateOverlay();
  window.addEventListener("resize", updateOverlay);
});
const intervalId = setInterval(() => {
  fetchTonKho();
  fetchTonKhoTang();
}, 5 * 60 * 1000);

onUnmounted(() => {
  window.removeEventListener("resize", updateOverlay);
  clearInterval(intervalId);
});

watch(
  () => state.warehouseSelect,
  () => {
    nextTick(() => {
      imageMapResize();
    });
  },
);
</script>

<style lang="scss" scoped>
.map-container {
  position: relative;
  display: flex;
  width: 100%;
  height: 90vh;
  overflow: auto;
}

.map-inner {
  position: relative;
  margin-left: 25px;
  display: inline-block;
}

.responsive-map-image {
  display: block;
  max-width: 100%;
  max-height: 100%;
}

/* ================= OVERLAY ================= */

.overlay-text {
  position: absolute;
  font-weight: bold;
  font-size: 12px;
  pointer-events: none;
  color: #000;
  transform: translate(-50%, -50%);
  text-align: center;
  white-space: nowrap;

  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 5;

  &.has-bg {
    background: transparent;
    padding: 0;
  }

  &.zero {
    color: gray;
  }
}

/* ================= FRACTION ================= */

.fraction-box {
  display: inline-flex;
  flex-direction: column;
  align-items: stretch;
  font-size: 11px;
  font-weight: bold;
  line-height: 1;
  border-radius: 4px;
  overflow: hidden; // QUAN TRỌNG
}
.fraction-top {
  background: #f5e663;
  padding: 2px 6px;
  text-align: center;
}

.fraction-divider {
  height: 1px;
  background: #000;
}

.fraction-bottom {
  background: #ff6b6b;
  padding: 2px 6px;
  text-align: center;
  color: #fff;
}

/* ================= SUB PROGRESS ================= */

.sub-progress-overlay {
  position: absolute;
  display: flex;
  align-items: flex-end;
  z-index: 1;
}

.progress-bar {
  width: 100%;
  height: 100%;
  background: transparent;
  border: 1px solid #014074; /* viền xanh */
  border-radius: 3px;
  overflow: hidden;
  display: flex;
  flex-direction: column-reverse;
}

.progress-used {
  width: 100%;
  background: #014074;
  transition: height 0.3s ease;
}

.progress-bar.full .progress-used {
  background: rgb(228, 0, 0);
}

/* ================= CARD ================= */

.box-card {
  width: 100%;
  height: 100%;
}

:deep(.el-card__header) {
  padding: 20px 0 0 0;
}

:deep(.el-card__body) {
  padding: 0;
}
</style>
