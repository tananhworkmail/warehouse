<template>
  <WarehouseLayoutKVT>
    <template #default>
      <el-card class="box-card">
        <!-- KhoDe -->
        <template v-if="state.warehouseSelect === 'KhoDe'">
          <div class="map-container">
            <div class="map-inner">
              <div v-show="isWarningLoading" class="warning-loading-badge">
                <span class="spinner"></span>
                <span>Đang tải dữ liệu kế hoạch 3 ngày và tồn kho...</span>
              </div>
              <img
                class="responsive-map-image"
                src="@/assets/images/KhoVatTu.jpg"
                usemap="#image-map"
              />

              <!-- Overlay số lượng cho các kệ chính/phụ -->
              <template v-for="rackId in mapAreas" :key="rackId">
                <div
                  v-if="rackWarnings[rackId]"
                  class="rack-color-overlay"
                  :style="rackWarnings[rackId].style"
                >
                  <div
                    v-if="rackWarnings[rackId].day180"
                    class="warning-dot red-dot"
                  ></div>
                  <div
                    v-if="rackWarnings[rackId].day3"
                    class="warning-dot orange-dot"
                  ></div>
                </div>

                <div
                  v-if="
                    isMainRack(rackId) || GetRackToltalColumn(rackId).total > 0
                  "
                  class="overlay-text"
                  :class="{ zero: GetRackToltalColumn(rackId).total === 0 }"
                  :style="getOverlayStyle(rackId)"
                >
                  <template v-if="isMainRack(rackId)">
                    <div class="fraction-box">
                      <div class="fraction-top">
                        {{ GetRackToltalColumn(rackId).total }}
                      </div>
                      <div
                        v-if="GetRackToltalColumn(rackId).limit > 0"
                        class="fraction-divider"
                      ></div>
                      <div
                        v-if="GetRackToltalColumn(rackId).limit > 0"
                        class="fraction-bottom"
                      >
                        {{ GetRackToltalColumn(rackId).limit }}
                      </div>
                    </div>
                  </template>
                  <template v-else>
                    <div class="sub-rack-text">
                      {{ GetRackToltalColumn(rackId).total }}
                    </div>
                  </template>
                </div>
              </template>

              <!-- Overlay thanh tiến trình cho Sub rack -->
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

              <!-- NÚT QR / KHUNG NHẬP LIỆU -->
              <div
                v-show="!showBarcodeInput"
                class="barcode-floating-button"
                :style="barcodeButton"
                @click.stop="openBarcodeInput"
              >
                <img
                  src="@/assets/qr-code.png"
                  alt="QR Icon"
                  style="width: 32px; height: 32px"
                />
              </div>

              <div
                v-show="showBarcodeInput"
                class="barcode-floating-box"
                :style="barcodeBoxStyle"
              >
                <div class="barcode-label">{{ t("kvt.scanCode") }}</div>
                <el-input
                  ref="barcodeInputRef"
                  v-model="barcodeValue"
                  :placeholder="t('kvt.scanPlaceholder')"
                  clearable
                  @keyup.enter="handleBarcodeSearch"
                  size="small"
                />
                <el-button
                  size="small"
                  @click="closeBarcodeInput"
                  style="margin-top: 6px"
                >
                  {{ t("kvt.close") }}
                </el-button>
              </div>

              <!-- Image Map (TOÀN BỘ AREA) -->
              <map name="image-map">
                <!-- B01 -->
                <area
                  id="B01"
                  title="B01"
                  href="#"
                  coords="2846,530,2737,481"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B02"
                  title="B02"
                  href="#"
                  coords="2846,425,2737,376"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B03"
                  title="B03"
                  href="#"
                  coords="2843,275,2734,226"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B04"
                  title="B04"
                  href="#"
                  coords="2843,122,2734,73"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- A01 → A17 (main) -->
                <area
                  id="A01"
                  title="A01"
                  href="#"
                  coords="2807,1092,2929,1155"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A02"
                  title="A02"
                  href="#"
                  coords="2698,1109,2768,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A03"
                  title="A03"
                  href="#"
                  coords="2566,1109,2636,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A04"
                  title="A04"
                  href="#"
                  coords="2428,1109,2498,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A05"
                  title="A05"
                  href="#"
                  coords="2299,1109,2371,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A06"
                  title="A06"
                  href="#"
                  coords="2167,1109,2239,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A07"
                  title="A07"
                  href="#"
                  coords="2035,1109,2107,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A08"
                  title="A08"
                  href="#"
                  coords="1903,1109,1975,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A09"
                  title="A09"
                  href="#"
                  coords="1771,1109,1843,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A10"
                  title="A10"
                  href="#"
                  coords="1633,1109,1708,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A11"
                  title="A11"
                  href="#"
                  coords="1501,1109,1576,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A12"
                  title="A12"
                  href="#"
                  coords="1367,1109,1446,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A13"
                  title="A13"
                  href="#"
                  coords="1232,1109,1311,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A14"
                  title="A14"
                  href="#"
                  coords="1103,1109,1182,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A15"
                  title="A15"
                  href="#"
                  coords="968,1109,1047,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A16"
                  title="A16"
                  href="#"
                  coords="836,1109,915,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A17"
                  title="A17"
                  href="#"
                  coords="701,1109,780,1145"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- A18 → A38 (main) -->
                <area
                  id="A18"
                  title="A18"
                  href="#"
                  coords="2827,605,2904,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A19"
                  title="A19"
                  href="#"
                  coords="2695,605,2772,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A20"
                  title="A20"
                  href="#"
                  coords="2563,605,2640,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A21"
                  title="A21"
                  href="#"
                  coords="2428,605,2505,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A22"
                  title="A22"
                  href="#"
                  coords="2293,605,2370,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A23"
                  title="A23"
                  href="#"
                  coords="2164,605,2241,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A24"
                  title="A24"
                  href="#"
                  coords="2030,605,2107,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A25"
                  title="A25"
                  href="#"
                  coords="1897,605,1974,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A26"
                  title="A26"
                  href="#"
                  coords="1768,605,1845,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A27"
                  title="A27"
                  href="#"
                  coords="1633,605,1710,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A28"
                  title="A28"
                  href="#"
                  coords="1498,605,1575,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A29"
                  title="A29"
                  href="#"
                  coords="1369,605,1446,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A30"
                  title="A30"
                  href="#"
                  coords="1237,605,1314,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A31"
                  title="A31"
                  href="#"
                  coords="1102,605,1179,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A32"
                  title="A32"
                  href="#"
                  coords="967,605,1044,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A33"
                  title="A33"
                  href="#"
                  coords="835,605,912,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A34"
                  title="A34"
                  href="#"
                  coords="704,605,781,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A35"
                  title="A35"
                  href="#"
                  coords="476,605,553,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A36"
                  title="A36"
                  href="#"
                  coords="341,605,418,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A37"
                  title="A37"
                  href="#"
                  coords="212,605,289,637"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A38"
                  title="A38"
                  href="#"
                  coords="77,605,154,637"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- Sub racks (dọc) -->
                <area
                  id="SubA01"
                  title="SubA01"
                  href="#"
                  coords="2811,1160,2824,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA02"
                  title="SubA02"
                  href="#"
                  coords="2680,1160,2693,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA03"
                  title="SubA03"
                  href="#"
                  coords="2544,1160,2557,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA04"
                  title="SubA04"
                  href="#"
                  coords="2413,1160,2426,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA05"
                  title="SubA05"
                  href="#"
                  coords="2284,1160,2297,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA06"
                  title="SubA06"
                  href="#"
                  coords="2152,1160,2165,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA07"
                  title="SubA07"
                  href="#"
                  coords="2020,1160,2033,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA08"
                  title="SubA08"
                  href="#"
                  coords="1888,1160,1901,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA09"
                  title="SubA09"
                  href="#"
                  coords="1754,1160,1767,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA10"
                  title="SubA10"
                  href="#"
                  coords="1621,1160,1634,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA11"
                  title="SubA11"
                  href="#"
                  coords="1487,1160,1500,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA12"
                  title="SubA12"
                  href="#"
                  coords="1353,1160,1366,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA13"
                  title="SubA13"
                  href="#"
                  coords="1218,1160,1231,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA14"
                  title="SubA14"
                  href="#"
                  coords="1089,1160,1102,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA15"
                  title="SubA15"
                  href="#"
                  coords="954,1160,967,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA16"
                  title="SubA16"
                  href="#"
                  coords="822,1160,835,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA17"
                  title="SubA17"
                  href="#"
                  coords="686,1160,699,1550"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA18"
                  title="SubA18"
                  href="#"
                  coords="2812,650,2825,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA19"
                  title="SubA19"
                  href="#"
                  coords="2677,650,2690,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA20"
                  title="SubA20"
                  href="#"
                  coords="2547,650,2560,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA21"
                  title="SubA21"
                  href="#"
                  coords="2414,650,2427,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA22"
                  title="SubA22"
                  href="#"
                  coords="2283,650,2296,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA23"
                  title="SubA23"
                  href="#"
                  coords="2150,650,2163,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA24"
                  title="SubA24"
                  href="#"
                  coords="2016,650,2029,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA25"
                  title="SubA25"
                  href="#"
                  coords="1884,650,1897,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA26"
                  title="SubA26"
                  href="#"
                  coords="1752,650,1765,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA27"
                  title="SubA27"
                  href="#"
                  coords="1620,650,1633,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA28"
                  title="SubA28"
                  href="#"
                  coords="1488,650,1501,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA29"
                  title="SubA29"
                  href="#"
                  coords="1356,650,1369,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA30"
                  title="SubA30"
                  href="#"
                  coords="1221,650,1234,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA31"
                  title="SubA31"
                  href="#"
                  coords="1090,650,1103,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA32"
                  title="SubA32"
                  href="#"
                  coords="958,650,971,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA33"
                  title="SubA33"
                  href="#"
                  coords="824,650,837,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA34"
                  title="SubA34"
                  href="#"
                  coords="694,650,707,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA35"
                  title="SubA35"
                  href="#"
                  coords="456,650,469,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA36"
                  title="SubA36"
                  href="#"
                  coords="324,650,337,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA37"
                  title="SubA37"
                  href="#"
                  coords="192,650,205,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubA38"
                  title="SubA38"
                  href="#"
                  coords="60,650,73,1045"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubB04"
                  title="SubB04"
                  href="#"
                  coords="1162,20,1175,155"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubB03"
                  title="SubB03"
                  href="#"
                  coords="1162,174,1175,309"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubB02"
                  title="SubB02"
                  href="#"
                  coords="1162,326,1175,459"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="SubB01"
                  title="SubB01"
                  href="#"
                  coords="1162,478,1175,525"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- B01 → B04 chi tiết -->
                <area
                  id="B0101"
                  title="B0101"
                  href="#"
                  coords="2725,531,1183,486"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B0201"
                  title="B0201"
                  href="#"
                  coords="2725,467,1183,422"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B0202"
                  title="B0202"
                  href="#"
                  coords="2725,422,1183,379"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B0203"
                  title="B0203"
                  href="#"
                  coords="2725,381,1183,336"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B0301"
                  title="B0301"
                  href="#"
                  coords="2725,314,1183,269"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B0302"
                  title="B0302"
                  href="#"
                  coords="2725,269,1183,226"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B0303"
                  title="B0303"
                  href="#"
                  coords="2725,222,1183,177"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B0401"
                  title="B0401"
                  href="#"
                  coords="2725,161,1183,116"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B0402"
                  title="B0402"
                  href="#"
                  coords="2725,116,1183,73"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="B0403"
                  title="B0403"
                  href="#"
                  coords="2725,72,1183,27"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- Các tầng con A0101... -->
                <area
                  id="A0101"
                  title="A0101"
                  href="#"
                  coords="2901,1549,2831,1416"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0102"
                  title="A0102"
                  href="#"
                  coords="2901,1416,2831,1285"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0103"
                  title="A0103"
                  href="#"
                  coords="2901,1281,2831,1154"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0201"
                  title="A0201"
                  href="#"
                  coords="2769,1549,2699,1416"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0202"
                  title="A0202"
                  href="#"
                  coords="2769,1416,2699,1285"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0203"
                  title="A0203"
                  href="#"
                  coords="2769,1281,2699,1154"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0301"
                  title="A0301"
                  href="#"
                  coords="2634,1546,2564,1413"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0302"
                  title="A0302"
                  href="#"
                  coords="2634,1416,2564,1285"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0303"
                  title="A0303"
                  href="#"
                  coords="2634,1281,2564,1154"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0401"
                  title="A0401"
                  href="#"
                  coords="2502,1552,2432,1419"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0402"
                  title="A0402"
                  href="#"
                  coords="2502,1416,2432,1285"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0403"
                  title="A0403"
                  href="#"
                  coords="2502,1281,2432,1154"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0501"
                  title="A0501"
                  href="#"
                  coords="2370,1549,2300,1416"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0502"
                  title="A0502"
                  href="#"
                  coords="2370,1416,2300,1285"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0503"
                  title="A0503"
                  href="#"
                  coords="2370,1284,2300,1157"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0601"
                  title="A0601"
                  href="#"
                  coords="2238,1543,2168,1410"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0602"
                  title="A0602"
                  href="#"
                  coords="2238,1416,2168,1285"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0603"
                  title="A0603"
                  href="#"
                  coords="2238,1281,2168,1154"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0701"
                  title="A0701"
                  href="#"
                  coords="2106,1549,2036,1416"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0702"
                  title="A0702"
                  href="#"
                  coords="2106,1416,2036,1285"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0703"
                  title="A0703"
                  href="#"
                  coords="2106,1281,2036,1154"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0801"
                  title="A0801"
                  href="#"
                  coords="1974,1549,1904,1416"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0802"
                  title="A0802"
                  href="#"
                  coords="1974,1416,1904,1285"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0803"
                  title="A0803"
                  href="#"
                  coords="1974,1281,1904,1154"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0901"
                  title="A0901"
                  href="#"
                  coords="1842,1549,1772,1416"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0902"
                  title="A0902"
                  href="#"
                  coords="1842,1416,1772,1285"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A0903"
                  title="A0903"
                  href="#"
                  coords="1842,1281,1772,1154"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1001"
                  title="A1001"
                  href="#"
                  coords="1711,1548,1635,1418"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1002"
                  title="A1002"
                  href="#"
                  coords="1706,1417,1633,1285"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1003"
                  title="A1003"
                  href="#"
                  coords="1706,1282,1633,1150"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1101"
                  title="A1101"
                  href="#"
                  coords="1576,1548,1500,1418"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1102"
                  title="A1102"
                  href="#"
                  coords="1576,1415,1501,1287"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1103"
                  title="A1103"
                  href="#"
                  coords="1576,1284,1500,1151"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1201"
                  title="A1201"
                  href="#"
                  coords="1446,1548,1367,1418"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1202"
                  title="A1202"
                  href="#"
                  coords="1444,1418,1365,1287"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1203"
                  title="A1203"
                  href="#"
                  coords="1445,1284,1366,1151"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1301"
                  title="A1301"
                  href="#"
                  coords="1311,1548,1232,1418"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1302"
                  title="A1302"
                  href="#"
                  coords="1312,1418,1233,1287"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1303"
                  title="A1303"
                  href="#"
                  coords="1313,1287,1234,1154"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1401"
                  title="A1401"
                  href="#"
                  coords="1177,1548,1098,1418"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1402"
                  title="A1402"
                  href="#"
                  coords="1177,1415,1098,1284"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1403"
                  title="A1403"
                  href="#"
                  coords="1177,1284,1098,1151"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1501"
                  title="A1501"
                  href="#"
                  coords="1047,1548,968,1418"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1502"
                  title="A1502"
                  href="#"
                  coords="1047,1415,968,1284"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1503"
                  title="A1503"
                  href="#"
                  coords="1047,1284,968,1151"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1601"
                  title="A1601"
                  href="#"
                  coords="915,1548,836,1418"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1602"
                  title="A1602"
                  href="#"
                  coords="915,1415,836,1284"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1603"
                  title="A1603"
                  href="#"
                  coords="915,1287,836,1154"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1701"
                  title="A1701"
                  href="#"
                  coords="780,1551,701,1421"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1702"
                  title="A1702"
                  href="#"
                  coords="780,1418,701,1287"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1703"
                  title="A1703"
                  href="#"
                  coords="780,1287,701,1154"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1801"
                  title="A1801"
                  href="#"
                  coords="2904,1039,2829,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1802"
                  title="A1802"
                  href="#"
                  coords="2904,905,2828,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1803"
                  title="A1803"
                  href="#"
                  coords="2904,772,2827,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1901"
                  title="A1901"
                  href="#"
                  coords="2771,1041,2694,905"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1902"
                  title="A1902"
                  href="#"
                  coords="2772,902,2696,777"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A1903"
                  title="A1903"
                  href="#"
                  coords="2772,772,2695,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2001"
                  title="A2001"
                  href="#"
                  coords="2638,1041,2561,905"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2002"
                  title="A2002"
                  href="#"
                  coords="2639,902,2562,777"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2003"
                  title="A2003"
                  href="#"
                  coords="2639,772,2562,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2101"
                  title="A2101"
                  href="#"
                  coords="2507,1041,2430,905"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2102"
                  title="A2102"
                  href="#"
                  coords="2505,902,2429,777"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2103"
                  title="A2103"
                  href="#"
                  coords="2505,772,2428,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2201"
                  title="A2201"
                  href="#"
                  coords="2374,1039,2295,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2202"
                  title="A2202"
                  href="#"
                  coords="2373,905,2297,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2203"
                  title="A2203"
                  href="#"
                  coords="2373,772,2296,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2301"
                  title="A2301"
                  href="#"
                  coords="2242,1039,2163,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2302"
                  title="A2302"
                  href="#"
                  coords="2238,905,2162,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2303"
                  title="A2303"
                  href="#"
                  coords="2241,775,2164,650"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2401"
                  title="A2401"
                  href="#"
                  coords="2109,1039,2030,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2402"
                  title="A2402"
                  href="#"
                  coords="2108,905,2032,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2403"
                  title="A2403"
                  href="#"
                  coords="2109,772,2032,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2501"
                  title="A2501"
                  href="#"
                  coords="1977,1039,1898,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2502"
                  title="A2502"
                  href="#"
                  coords="1976,905,1900,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2503"
                  title="A2503"
                  href="#"
                  coords="1977,772,1900,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2601"
                  title="A2601"
                  href="#"
                  coords="1844,1039,1765,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2602"
                  title="A2602"
                  href="#"
                  coords="1843,905,1767,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2603"
                  title="A2603"
                  href="#"
                  coords="1844,772,1767,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2701"
                  title="A2701"
                  href="#"
                  coords="1711,1039,1632,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2702"
                  title="A2702"
                  href="#"
                  coords="1710,905,1634,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2703"
                  title="A2703"
                  href="#"
                  coords="1711,772,1634,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2801"
                  title="A2801"
                  href="#"
                  coords="1578,1039,1499,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2802"
                  title="A2802"
                  href="#"
                  coords="1577,905,1501,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2803"
                  title="A2803"
                  href="#"
                  coords="1578,772,1501,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2901"
                  title="A2901"
                  href="#"
                  coords="1445,1039,1366,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2902"
                  title="A2902"
                  href="#"
                  coords="1444,905,1368,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A2903"
                  title="A2903"
                  href="#"
                  coords="1445,772,1368,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3001"
                  title="A3001"
                  href="#"
                  coords="1312,1039,1233,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3002"
                  title="A3002"
                  href="#"
                  coords="1311,905,1235,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3003"
                  title="A3003"
                  href="#"
                  coords="1312,772,1235,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3101"
                  title="A3101"
                  href="#"
                  coords="1179,1039,1100,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3102"
                  title="A3102"
                  href="#"
                  coords="1178,905,1102,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3103"
                  title="A3103"
                  href="#"
                  coords="1179,772,1102,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3201"
                  title="A3201"
                  href="#"
                  coords="1046,1039,967,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3202"
                  title="A3202"
                  href="#"
                  coords="1045,905,969,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3203"
                  title="A3203"
                  href="#"
                  coords="1046,772,969,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3301"
                  title="A3301"
                  href="#"
                  coords="913,1039,834,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3302"
                  title="A3302"
                  href="#"
                  coords="912,905,836,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3303"
                  title="A3303"
                  href="#"
                  coords="913,772,836,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3401"
                  title="A3401"
                  href="#"
                  coords="780,1039,701,908"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3402"
                  title="A3402"
                  href="#"
                  coords="779,905,703,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3403"
                  title="A3403"
                  href="#"
                  coords="780,772,703,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3501"
                  title="A3501"
                  href="#"
                  coords="555,1036,476,905"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3502"
                  title="A3502"
                  href="#"
                  coords="552,902,476,777"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3503"
                  title="A3503"
                  href="#"
                  coords="551,772,474,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3601"
                  title="A3601"
                  href="#"
                  coords="422,1036,343,905"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3602"
                  title="A3602"
                  href="#"
                  coords="421,902,345,777"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3603"
                  title="A3603"
                  href="#"
                  coords="422,772,345,647"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3701"
                  title="A3701"
                  href="#"
                  coords="285,1042,206,911"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3702"
                  title="A3702"
                  href="#"
                  coords="282,905,206,780"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3703"
                  title="A3703"
                  href="#"
                  coords="284,769,207,644"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3801"
                  title="A3801"
                  href="#"
                  coords="156,1036,77,905"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3802"
                  title="A3802"
                  href="#"
                  coords="155,902,79,777"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="A3803"
                  title="A3803"
                  href="#"
                  coords="156,772,79,647"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- Hàng H -->
                <area
                  id="H01"
                  title="H01"
                  href="#"
                  coords="3378,1149,3437,1549"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H02"
                  title="H02"
                  href="#"
                  coords="3316,1149,3375,1549"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H03"
                  title="H03"
                  href="#"
                  coords="3250,1149,3309,1549"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H04"
                  title="H04"
                  href="#"
                  coords="3186,1149,3245,1549"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H05"
                  title="H05"
                  href="#"
                  coords="3124,1149,3175,1549"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H06"
                  title="H06"
                  href="#"
                  coords="551,1208,61,1153"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H07"
                  title="H07"
                  href="#"
                  coords="551,1285,61,1228"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H08"
                  title="H08"
                  href="#"
                  coords="551,1361,64,1304"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H09"
                  title="H09"
                  href="#"
                  coords="551,1440,61,1376"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H10"
                  title="H10"
                  href="#"
                  coords="551,1509,61,1458"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="H11"
                  title="H11"
                  href="#"
                  coords="551,1588,61,1530"
                  shape="rect"
                  @click="handleAreaClick"
                />

                <!-- C -->
                <area
                  id="C01"
                  title="C01"
                  href="#"
                  coords="3384,642,3435,1041"
                  shape="rect"
                  @click="handleAreaClick"
                />
                <area
                  id="C02"
                  title="C02"
                  href="#"
                  coords="2987,1152,3036,1548"
                  shape="rect"
                  @click="handleAreaClick"
                />
              </map>
            </div>
          </div>
        </template>

        <!-- Dialog thông tin rack -->
        <RackInformationDialogKVT
          v-if="state.rackInformationDialogVisible"
          :rackNo="state.rackNo"
          :rack-data-list="rackDataList"
          :dialogVisible="state.rackInformationDialogVisible"
          :auto-select-rack="state.autoSelectRack"
          @close="
            state.rackInformationDialogVisible = false;
            state.autoSelectRack = '';
          "
        />
        <AlertKVT v-model="state.alertDialogVisible" />
      </el-card>
    </template>
  </WarehouseLayoutKVT>
</template>

<script setup>
import {
  reactive,
  onMounted,
  nextTick,
  watch,
  ref,
  onUnmounted,
  onBeforeUnmount,
  computed,
} from "vue";
import imageMapResize from "image-map-resizer";
import axios from "axios";
import WarehouseLayoutKVT from "@/views/components/WarehouseKVT/WarehouseLayoutKVT.vue";
import RackInformationDialogKVT from "../WarehouseKVT/RackInformationDialogKVT/index.vue";
import { useWarehouseMapI18n } from "@/hooks/useWarehouseMapI18n";
import AlertKVT from "@/views/components/WarehouseKVT/AlertDialogKVT/AlertKVT.vue";
import { useRoute } from "vue-router";
const API_URL = import.meta.env.VITE_API_URL;
const { t } = useWarehouseMapI18n();
const route = useRoute();
const state = reactive({
  alertDialogVisible: false,
  warehouseSelect: "KhoDe",
  rackInformationDialogVisible: false,
  rackNo: "",
  autoSelectRack: "",
});
const isWarningLoading = ref(false);
const rackDataList = ref([]);
const tonKhoKeList = ref([]);
const tonKhoTangList = ref([]);
const barcodeValue = ref("");
const barcodeInputRef = ref(null);
const showBarcodeInput = ref(false);
const rackStatusList = ref([]);
const imgClientWidth = ref(0);
const imgClientHeight = ref(0);
const imgNaturalWidth = ref(0);
const imgNaturalHeight = ref(0);

const mainRackList = [
  "A01",
  "A02",
  "A03",
  "A04",
  "A05",
  "A06",
  "A07",
  "A08",
  "A09",
  "A10",
  "A11",
  "A12",
  "A13",
  "A14",
  "A15",
  "A16",
  "A17",
  "A18",
  "A19",
  "A20",
  "A21",
  "A22",
  "A23",
  "A24",
  "A25",
  "A26",
  "A27",
  "A28",
  "A29",
  "A30",
  "A31",
  "A32",
  "A33",
  "A34",
  "A35",
  "A36",
  "A37",
  "A38",
  "B01",
  "B02",
  "B03",
  "B04",
  "C01",
  "C02",
  "H01",
  "H02",
  "H03",
  "H04",
  "H05",
  "H06",
  "H07",
  "H08",
  "H09",
  "H10",
  "H11",
];
const closeBarcodeInput = () => {
  showBarcodeInput.value = false;
  barcodeValue.value = ""; // xoá mã đã nhập
};
const isMainRack = (rackId) => mainRackList.includes(rackId);

const GetRackToltalColumn = (rackId) => {
  const list = isMainRack(rackId) ? tonKhoKeList.value : tonKhoTangList.value;
  if (!Array.isArray(list)) return { total: 0, limit: 0 };
  if (isMainRack(rackId)) {
    const found = list.find((item) => item.make3 === rackId);
    return { total: found?.totalcolumn || 0, limit: found?.limit || 0 };
  } else {
    const total = list
      .filter((item) => item.make5?.startsWith(rackId))
      .reduce((sum, item) => sum + (item.totaltang || 0), 0);
    return { total, limit: 0 };
  }
};

const getRackPercent = (rackId) => {
  const rack = GetRackToltalColumn(rackId);
  if (!rack.limit || rack.limit === 0) return 0;
  return Math.min(Math.round((rack.total / rack.limit) * 100), 100);
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

const rackWarnings = computed(() => {
  const results = {};

  mapAreas.value.forEach((rackId) => {
    const areaElement = document.querySelector(`area[id="${rackId}"]`);
    if (!areaElement || !areaElement.coords) return;

    const coords = areaElement.coords.split(",").map(Number);
    if (coords.length < 4 || isNaN(coords[0])) return;

    const [x1, y1, x2, y2] = coords;

    // Luôn lấy điểm nhỏ nhất làm điểm xuất phát (Top - Left)
    const trueLeft = Math.min(x1, x2);
    const trueTop = Math.min(y1, y2);
    const trueWidth = Math.abs(x2 - x1);
    const trueHeight = Math.abs(y2 - y1);

    // Xử lý H và C
    let searchMake = rackId;
    if (rackId.startsWith("H") || rackId.startsWith("C")) {
      searchMake = rackId + "01";
    }

    const status = rackStatusList.value.find(
      (item) => item.make === searchMake,
    );

    // Nếu có cảnh báo day3 hoặc day180 thì mới lưu vào object
    if (status && (status.day3 === 1 || status.day180 === 1)) {
      results[rackId] = {
        style: {
          position: "absolute",
          left: `${trueLeft}px`,
          top: `${trueTop}px`,
          width: `${trueWidth}px`,
          height: `${trueHeight}px`,
          display: "flex",
          flexDirection: "row",
          justifyContent: "flex-end",
          alignItems: "flex-start",
          paddingTop: "2px",
          paddingRight: "2px",
          gap: "0px",
          pointerEvents: "none",
          zIndex: 2,
        },
        day3: status.day3 === 1,
        day180: status.day180 === 1,
      };
    }
  });

  return results;
});

const barcodeButton = computed(() => {
  if (!imgNaturalWidth.value) return { display: "none" };
  const scaleX = imgClientWidth.value / imgNaturalWidth.value;
  const scaleY = imgClientHeight.value / imgNaturalHeight.value;
  const x1 = 3200,
    y1 = 70,
    x2 = 3500,
    y2 = 100;
  const centerX = (x1 + x2) / 2;
  const centerY = (y1 + y2) / 2;
  return {
    position: "absolute",
    left: `${centerX * scaleX}px`,
    top: `${centerY * scaleY}px`,
    transform: "translate(-50%, -50%)",
    zIndex: 100,
  };
});

const barcodeBoxStyle = computed(() => {
  if (!imgNaturalWidth.value) return { display: "none" };
  const scaleX = imgClientWidth.value / imgNaturalWidth.value;
  const scaleY = imgClientHeight.value / imgNaturalHeight.value;
  const x1 = 3100,
    y1 = 15,
    x2 = 3400,
    y2 = 200;
  return {
    position: "absolute",
    left: `${x1 * scaleX}px`,
    top: `${y1 * scaleY}px`,
    width: `${(x2 - x1) * scaleX}px`,
    height: `${(y2 - y1) * scaleY}px`,
    zIndex: 100,
  };
});

const updateImageSize = () => {
  const img = document.querySelector(".responsive-map-image");
  if (!img || !img.naturalWidth) return;
  imgClientWidth.value = img.clientWidth;
  imgClientHeight.value = img.clientHeight;
  imgNaturalWidth.value = img.naturalWidth;
  imgNaturalHeight.value = img.naturalHeight;
};

const handleAreaClick = async (event) => {
  const rackCode = event.target.id;
  state.rackNo = rackCode;
  try {
    const response = await axios.get(`${API_URL}/warehouse/rackkvt`, {
      params: { rackCode },
    });
    rackDataList.value = (response.data.data || []).map((item) => item.make);
    state.rackInformationDialogVisible = true;
  } catch (error) {
    console.error("Lỗi khi gọi API /rackkvt:", error);
  }
};

const openBarcodeInput = () => {
  showBarcodeInput.value = true;
  nextTick(() => {
    barcodeInputRef.value?.focus();
  });
};

const handleBarcodeSearch = async () => {
  const raw = barcodeValue.value.trim().toUpperCase();
  if (!raw) return;

  let rackNo = raw;
  let autoSelect = "";

  const fullMatch = raw.match(/^([A-Z]\d{2})(\d{2})(\d{2})$/); // A31 01 01
  if (fullMatch) {
    rackNo = fullMatch[1]; // A31
    autoSelect = raw; // A310101  → dialog sẽ tự parse
  } else {
    const shortMatch = raw.match(/^([A-Z]\d{2})(\d{2})$/); // A3101
    if (shortMatch) {
      rackNo = shortMatch[1]; // A31
      autoSelect = raw; // A3101
    }
    // else: mã ngắn như "A31" → chỉ mở dialog bình thường
  }

  state.rackNo = rackNo;
  state.autoSelectRack = autoSelect;

  try {
    const response = await axios.get(`${API_URL}/warehouse/rackkvt`, {
      params: { rackCode: rackNo },
    });
    rackDataList.value = (response.data.data || []).map((item) => item.make);
    state.rackInformationDialogVisible = true;
  } catch (error) {
    console.error("Lỗi barcode search:", error);
  }

  barcodeValue.value = "";
  nextTick(() => {
    barcodeInputRef.value?.focus();
  });
};

const fetchTonKhoTang = async () => {
  try {
    const res = await axios.get(`${API_URL}/warehouse/totaltangkvt`);
    tonKhoTangList.value = Array.isArray(res.data.data) ? res.data.data : [];
  } catch (err) {
    console.error(err);
    tonKhoTangList.value = [];
  }
};

const fetchTonKhoKe = async () => {
  try {
    const res = await axios.get(`${API_URL}/warehouse/totalcolumnkvt`);
    tonKhoKeList.value = Array.isArray(res.data.data) ? res.data.data : [];
  } catch (err) {
    console.error(err);
    tonKhoKeList.value = [];
  }
};
const fetch3And180Day = async () => {
  isWarningLoading.value = true; // Bật loading khi bắt đầu gọi API
  try {
    const { data } = await axios.get(`${API_URL}/warehouse/3and180daykvt`);
    rackStatusList.value = data?.data ?? [];
  } catch (error) {
    console.error("Fetch 3 & 180 day KVT failed:", error);
    rackStatusList.value = [];
  } finally {
    isWarningLoading.value = false; // Tắt loading khi API chạy xong (dù lỗi hay thành công)
  }
};
const mapAreas = ref([]);
let resizeTimer;

const updateOverlay = () => {
  clearTimeout(resizeTimer);
  resizeTimer = setTimeout(() => {
    nextTick(() => {
      imageMapResize();
      updateImageSize();
      const areaIds = [...document.querySelectorAll("area")].map((a) => a.id);
      mapAreas.value = areaIds;
      rackStatusList.value = [...rackStatusList.value];
    });
  }, 200);
};

watch(
  () => route.path,
  async (newPath) => {
    if (newPath.includes("HumTemp")) {
      await nextTick();
      state.alertDialogVisible = true;
    }
  },
);
onMounted(async () => {
  console.log("Đường dẫn hiện tại:", route.path);
  console.log("Tên route hiện tại:", route.name);

  // Dùng route.name để check, hoặc dùng .includes() để bỏ qua lỗi sai khác viết hoa/thường
  if (
    route.name === "warehouse-kvt-humtemp" ||
    route.path.includes("HumTemp")
  ) {
    await nextTick(); // Đợi DOM của component con sẵn sàng
    state.alertDialogVisible = true;
    console.log("Trạng thái Dialog đã được bật thành TRUE");
  }
  // Chạy song song 2 API
  await Promise.all([fetchTonKhoTang(), fetchTonKhoKe()]);

  // Chạy nền, không chờ
  fetch3And180Day();

  await nextTick();

  imageMapResize();
  updateImageSize();

  const areaIds = [...document.querySelectorAll("area")].map((a) => a.id);
  mapAreas.value = areaIds;

  updateOverlay();
  window.addEventListener("resize", updateOverlay);
});
onUnmounted(() => {
  window.removeEventListener("resize", updateOverlay);
});

const intervalId = setInterval(() => {
  fetchTonKhoTang();
  fetchTonKhoKe();
}, 5 * 60 * 1000);

onBeforeUnmount(() => {
  clearInterval(intervalId);
});
</script>

<style lang="scss" scoped>
.map-container {
  position: relative;
  display: flex;
  width: 100%;
  height: 92vh;
  overflow: auto;
}
.map-inner {
  position: relative;
  display: inline-block;
  margin: 0 auto;
}
.responsive-map-image {
  display: block;
  width: auto;
  height: auto;
  max-width: 100%;
  max-height: 100%;
}

.overlay-text {
  position: absolute;
  font-weight: bold;
  padding: 2px 4px;
  border-radius: 4px;
  font-size: 8px;
  pointer-events: none;
  color: #000;
  transform: translate(-50%, -50%);
  text-align: center;
  white-space: nowrap;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 5;
}
.overlay-text.zero {
  color: gray;
}

.box-card {
  width: 100%;
  height: 100%;
}
:deep(.el-card__header) {
  padding: 20px 0px 0px 0px;
}
:deep(.el-card__body) {
  padding: 0;
}

.fraction-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  font-size: 9px;
  font-weight: bold;
  line-height: 1.1;
}
.fraction-top {
  background: #f5e663;
  padding: 1px 4px;
  border-radius: 3px 3px 0 0;
}
.fraction-divider {
  width: 100%;
  height: 1px;
  background: #000;
  margin: 1px 0;
}
.fraction-bottom {
  background: #ff6b6b;
  padding: 1px 4px;
  border-radius: 0 0 3px 3px;
  color: #fff;
}

.sub-progress-overlay {
  position: absolute;
  display: flex;
  align-items: flex;
  z-index: 1;
}
.progress-bar {
  width: 100%;
  height: 100%;
  background: transparent;
  border: 1px solid #014074;
  border-radius: 3px;
  overflow: hidden;
  position: relative;
  display: flex;
  flex-direction: column-reverse;
}
.progress-used {
  width: 100%;
  background: #014074;
  transition: height 0.3s ease;
}
.progress-bar.full .progress-used {
  background: red;
}

/* Nút QR nổi */
.barcode-floating-button {
  position: absolute;
  cursor: pointer;
}

/* Khung nhập liệu nổi */
.barcode-floating-box {
  background: #ebecf3eb;
  border: 4px solid black;
  border-radius: 8px;
  padding: 8px 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 6px;
  position: absolute; /* được override bởi inline style */
}

.barcode-label {
  font-size: 11px;
  font-weight: bold;
  color: #014074;
}
/* Nút QR nổi bật */
.barcode-floating-button {
  position: absolute;
  cursor: pointer;
  width: 50px; /* Tăng kích thước vùng bấm */
  height: 50px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 50%; /* Hình tròn */
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 15px rgba(0, 123, 255, 0.6), 0 4px 15px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
  border: 2px solid #007bff;
  animation: pulse 2s infinite;
}
.rack-color-overlay {
  z-index: 2 !important;
  pointer-events: none;
}

.warning-dot {
  width: 10px; /* Đường kính chấm tròn (có thể tăng giảm) */
  height: 10px;
  border-radius: 50%; /* Bo tròn hoàn toàn */
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.4); /* Thêm bóng đổ nhẹ để chấm tròn nổi bật trên nền bản đồ */
}

/* Màu đỏ cho day180 */
.red-dot {
  background-color: rgba(255, 0, 0, 0.85);
}

/* Màu cam cho day3 */
.orange-dot {
  background-color: rgba(255, 165, 0, 0.85);
}

/* Chữ số hiển thị tổng số lượng */
.overlay-text {
  z-index: 5 !important; /* Đảm bảo text luôn hiển thị đè lên trên lớp cảnh báo */
}

/* Thanh tiến trình của Sub rack (nếu có) */
.sub-progress-overlay {
  z-index: 5 !important; /* Đảm bảo hiển thị cao nhất cùng với text */
}
/* Hiệu ứng phát sáng nhẹ */
@keyframes pulse {
  0% {
    box-shadow: 0 0 10px rgba(0, 123, 255, 0.5), 0 4px 12px rgba(0, 0, 0, 0.2);
  }
  50% {
    box-shadow: 0 0 25px rgba(0, 123, 255, 0.8), 0 6px 20px rgba(0, 0, 0, 0.4);
  }
  100% {
    box-shadow: 0 0 10px rgba(0, 123, 255, 0.5), 0 4px 12px rgba(0, 0, 0, 0.2);
  }
}

.barcode-floating-button:hover {
  transform: scale(1.15);
  box-shadow: 0 0 30px rgba(0, 123, 255, 0.9), 0 8px 25px rgba(0, 0, 0, 0.5);
}

/* Ảnh QR bên trong */
.qr-icon-img {
  width: 70%; /* Ảnh chiếm 70% vùng tròn */
  height: 70%;
  object-fit: contain;
  border-radius: 4px; /* Bo nhẹ nếu ảnh vuông */
}
.warning-loading-badge {
  position: absolute;
  top: 10px;
  left: 10px;
  background-color: rgba(0, 0, 0, 0.6);
  color: #fff;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 13px;
  z-index: 100;
  display: flex;
  align-items: center;
  gap: 8px;
  pointer-events: none; /* Không cản trở click chuột */
}

/* Vòng xoay (spinner) tự tạo bằng CSS */
.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
