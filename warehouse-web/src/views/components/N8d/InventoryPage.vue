<template>
  <div class="inventory-page">
    <N8dHeader back-to="/n8d">
      <template #context>
        <div class="header-context" :class="{ offline: !hasSession }">
          <span
            >{{ mode === "month" ? copy.modeMonth : copy.modeDay
            }}<template v-if="mode === 'month'">
              ·
              {{
                inventoryRole === "accounting"
                  ? copy.roleAccounting
                  : copy.roleUnit
              }}</template
            ></span
          >
          <strong>{{ gsbh || copy.sessionMissing }}</strong>
        </div>

        <span
          class="user-chip"
          :class="{ offline: !hasSession }"
          :title="
            hasSession ? `${copy.operator}: ${userId}` : copy.sessionMissing
          "
        >
          <span class="status-dot" />
          {{ userId || copy.guest }}
        </span>
      </template>

      <main class="inventory-workspace">
        <!-- YÊU CẦU ĐĂNG NHẬP -->
        <section v-if="!hasSession" class="auth-required" role="status">
          <span class="auth-required__icon" aria-hidden="true">
            <svg viewBox="0 0 24 24">
              <rect x="5" y="10" width="14" height="11" rx="3" />
              <path d="M8.5 10V7.5a3.5 3.5 0 0 1 7 0V10M12 14v3" />
            </svg>
          </span>
          <div class="auth-required__copy">
            <p>{{ copy.sessionEyebrow }}</p>
            <h2>{{ copy.signInTitle }}</h2>
            <span>{{ copy.signInDescription }}</span>
          </div>
          <div class="auth-required__actions">
            <button class="primary" type="button" @click="goToLogin">
              {{ copy.signInAction }}
              <svg viewBox="0 0 20 20" aria-hidden="true">
                <path d="M4 10h11m-4-4 4 4-4 4" />
              </svg>
            </button>
            <button type="button" @click="router.push('/n8d')">
              {{ copy.backToHome }}
            </button>
          </div>
        </section>

        <template v-else>
          <!-- BẢNG ĐIỀU KHIỂN -->
          <section class="control-panel">
            <div class="control-grid">
              <label class="field-group">
                <span>{{
                  mode === "month"
                    ? t("inventoryPage.monthLabel")
                    : t("inventoryPage.dateLabel")
                }}</span>
                <span class="input-shell input-shell--picker">
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <rect x="3" y="5" width="18" height="16" rx="3" />
                    <path d="M3 10h18M8 3v4m8-4v4" />
                  </svg>
                  <el-date-picker
                    v-model="searchValue"
                    class="period-picker"
                    popper-class="inv-picker-popper"
                    :type="mode === 'month' ? 'month' : 'date'"
                    :value-format="mode === 'month' ? 'YYYY-MM' : 'YYYY-MM-DD'"
                    :clearable="false"
                    :editable="false"
                    @change="onPeriodChange"
                  />
                </span>
              </label>

              <label class="field-group">
                <span>{{
                  mode === "month"
                    ? t("inventoryPage.selectRackFirst")
                    : t("inventoryPage.filterByRack")
                }}</span>
                <span
                  class="input-shell"
                  :class="{
                    'input-shell--required': mode === 'month' && !rackFilter,
                  }"
                >
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M4 4h16v16H4zM4 10h16M4 15h16M9 4v16" />
                  </svg>
                  <el-select
                    v-if="mode === 'month'"
                    v-model="rackFilter"
                    class="rack-select"
                    popper-class="inv-select-popper"
                    :disabled="loadingRacks"
                    :placeholder="t('inventoryPage.selectRackFirst')"
                    @change="onRackSelected"
                  >
                    <el-option
                      v-for="rack in rackOptions"
                      :key="rack.value"
                      :label="rack.value"
                      :value="rack.value"
                    />
                  </el-select>
                  <input
                    v-else
                    v-model="rackFilter"
                    list="inventory-rack-options"
                    :placeholder="t('inventoryPage.filterByRack')"
                    @input="currentPage = 1"
                  />
                  <datalist v-if="mode === 'day'" id="inventory-rack-options">
                    <option
                      v-for="rack in rackOptions"
                      :key="rack"
                      :value="rack"
                    />
                  </datalist>
                </span>
              </label>

              <label class="field-group scan-field">
                <span class="scan-label">
                  <span>{{ t("inventoryPage.scanBarcodeHere") }}</span>
                  <small :class="{ ready: canScanInput }"
                    ><i />
                    {{
                      canScanInput ? copy.scanReady : copy.scanWaiting
                    }}</small
                  >
                </span>
                <span
                  class="scan-shell"
                  :class="{
                    disabled: !canScanInput,
                    'is-ready': canScanInput,
                    'is-busy': scanBusy,
                  }"
                >
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M4 5v14M7 5v14M11 5v14M14 5v14M18 5v14M20 5v14" />
                  </svg>
                  <input
                    ref="barcodeInput"
                    v-model="barcode"
                    :placeholder="
                      canScanInput
                        ? t('inventoryPage.scanBarcodeHere')
                        : t('inventoryPage.selectRackBeforeScan')
                    "
                    :disabled="!canScanInput || scanBusy"
                    autocomplete="off"
                    autofocus
                    @keydown.enter.prevent="handleBarcodeScan"
                  />
                  <button
                    type="button"
                    :disabled="!canSubmitScan"
                    @click="handleBarcodeScan"
                  >
                    <span v-if="scanBusy" class="spinner" />
                    <span v-else>{{ copy.record }}</span>
                    <svg
                      v-if="!scanBusy"
                      viewBox="0 0 20 20"
                      aria-hidden="true"
                    >
                      <path d="M4 10h11m-4-4 4 4-4 4" />
                    </svg>
                  </button>
                </span>
              </label>
            </div>

            <div v-if="mode === 'month'" class="control-status">
              <div v-if="loadingRacks" class="inline-state info" role="status">
                <span class="spinner" /> {{ t("inventoryPage.loadingData") }}
              </div>
              <div
                v-else-if="rackError"
                class="inline-state error"
                role="alert"
              >
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <circle cx="12" cy="12" r="9" />
                  <path d="M12 7v6m0 4h.01" />
                </svg>
                <span>{{ rackError }}</span
                ><button type="button" @click="loadRackOptions">
                  {{ t("inventoryPage.retry") }}
                </button>
              </div>
              <div v-else-if="needsSync" class="inline-state warning">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <path d="M12 3 2.8 20h18.4L12 3Z" />
                  <path d="M12 9v5m0 3h.01" />
                </svg>
                <span v-if="inventoryRole === 'accounting'">{{
                  copy.syncNeededAccounting
                }}</span>
                <span v-else
                  ><strong>{{ copy.syncNeededUnitTitle }}</strong>
                  {{ copy.syncNeededUnit }}</span
                >
              </div>
              <div v-else-if="!rackFilter" class="inline-state warning">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <path d="M12 3 2.8 20h18.4L12 3Z" />
                  <path d="M12 9v5m0 3h.01" />
                </svg>
                <span
                  ><strong>{{ t("inventoryPage.rackWarningTitle") }}</strong>
                  {{ t("inventoryPage.rackWarningSubtitle") }}</span
                >
              </div>
            </div>
          </section>

          <!-- KHU VỰC KẾT QUẢ -->
          <section v-if="showContent" class="results-area">
            <div class="stats-grid">
              <article class="stat-card stat-card--stock">
                <span class="stat-icon"
                  ><svg viewBox="0 0 24 24" aria-hidden="true">
                    <path
                      d="m4 8 8-4 8 4-8 4-8-4Zm0 0v8l8 4 8-4V8M12 12v8"
                    /></svg
                ></span>
                <span class="stat-copy"
                  ><small>{{ t("inventoryPage.statRemaining") }}</small
                  ><strong>{{ stats.rem }}</strong></span
                >
              </article>
              <article class="stat-card stat-card--scanned">
                <span class="stat-icon"
                  ><svg viewBox="0 0 24 24" aria-hidden="true">
                    <circle cx="12" cy="12" r="9" />
                    <path d="m8 12 2.5 2.5L16 9" /></svg
                ></span>
                <span class="stat-copy"
                  ><small>{{ t("inventoryPage.statScanned") }}</small
                  ><strong>{{ stats.scan }}</strong></span
                >
              </article>
              <article class="stat-card stat-card--diff">
                <span class="stat-icon"
                  ><svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M5 7h14M5 17h14M8 4 5 7l3 3m8 4 3 3-3 3" /></svg
                ></span>
                <span class="stat-copy"
                  ><small>{{ t("inventoryPage.statDiff") }}</small
                  ><strong>{{ stats.diff }}</strong></span
                >
              </article>
            </div>

            <section class="list-panel">
              <div class="list-toolbar">
                <div class="list-heading">
                  <span v-if="mode === 'month'" class="rack-badge">{{
                    rackFilter
                  }}</span>
                  <div>
                    <h2>{{ listTitle }}</h2>
                    <span
                      >{{ tabFilteredCount }}
                      {{ t("inventoryPage.rows") }}</span
                    >
                  </div>
                </div>
                <div class="list-actions">
                  <!-- Ô TÌM KIẾM TẤT CẢ -->
                  <div class="search-box">
                    <svg viewBox="0 0 24 24" aria-hidden="true">
                      <circle cx="11" cy="11" r="8" />
                      <path d="m21 21-4.3-4.3" />
                    </svg>
                    <input
                      v-model="searchQuery"
                      type="text"
                      :placeholder="t('inventoryPage.searchPlaceholder')"
                    />
                  </div>

                  <label v-if="showList" class="page-size">
                    <span>{{ copy.rowsPerPage }}</span>
                    <select
                      v-model.number="pageSize"
                      @change="onPageSizeChange"
                    >
                      <option
                        v-for="size in pageSizeOptions"
                        :key="size"
                        :value="size"
                      >
                        {{ size }}
                      </option>
                    </select>
                  </label>
                  <button
                    class="visibility-action"
                    type="button"
                    @click="toggleListVisibility"
                  >
                    <svg v-if="showList" viewBox="0 0 24 24" aria-hidden="true">
                      <path d="M3 12s3.3-6 9-6 9 6 9 6-3.3 6-9 6-9-6-9-6Z" />
                      <circle cx="12" cy="12" r="2.5" />
                      <path d="m4 4 16 16" />
                    </svg>
                    <svg v-else viewBox="0 0 24 24" aria-hidden="true">
                      <path d="M3 12s3.3-6 9-6 9 6 9 6-3.3 6-9 6-9-6-9-6Z" />
                      <circle cx="12" cy="12" r="2.5" />
                    </svg>
                    {{
                      showList
                        ? t("inventoryPage.hideList")
                        : t("inventoryPage.showList")
                    }}
                  </button>
                </div>
              </div>

              <div class="status-tabs" role="tablist">
                <button
                  v-for="tab in statusTabs"
                  :key="tab.value"
                  type="button"
                  role="tab"
                  :aria-selected="statusTab === tab.value"
                  :class="[tab.value, { active: statusTab === tab.value }]"
                  @click="setStatusTab(tab.value)"
                >
                  <span>{{ tab.label }}</span
                  ><strong>{{ tab.count }}</strong>
                </button>
              </div>

              <div
                v-if="loadingInventory && inventoryLoaded"
                class="refresh-line"
              >
                <span />
              </div>

              <div
                v-if="inventoryError && displayedInventoryList.length > 0"
                class="data-alert"
                role="alert"
              >
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <circle cx="12" cy="12" r="9" />
                  <path d="M12 7v6m0 4h.01" />
                </svg>
                <span>{{ inventoryError }}</span
                ><button type="button" @click="retryInventory">
                  {{ t("inventoryPage.retry") }}
                </button>
              </div>

              <div v-if="!showList" class="data-state">
                <span class="state-icon"
                  ><svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M3 12s3.3-6 9-6 9 6 9 6-3.3 6-9 6-9-6-9-6Z" />
                    <path d="m4 4 16 16" /></svg
                ></span>
                <strong>{{ t("inventoryPage.listHidden") }}</strong>
                <p>{{ t("inventoryPage.listHiddenHint") }}</p>
              </div>

              <div
                v-else-if="loadingInventory && !inventoryLoaded"
                class="data-state"
                role="status"
              >
                <span class="large-spinner" /><strong>{{
                  t("inventoryPage.loadingData")
                }}</strong>
                <p>{{ copy.loadingHint }}</p>
              </div>

              <div
                v-else-if="
                  inventoryError && displayedInventoryList.length === 0
                "
                class="data-state data-state--error"
                role="alert"
              >
                <span class="state-icon"
                  ><svg viewBox="0 0 24 24" aria-hidden="true">
                    <circle cx="12" cy="12" r="9" />
                    <path d="M12 7v6m0 4h.01" /></svg
                ></span>
                <strong>{{ copy.loadFailed }}</strong>
                <p>{{ inventoryError }}</p>
                <button type="button" @click="retryInventory">
                  {{ t("inventoryPage.retry") }}
                </button>
              </div>

              <div v-else-if="tabFilteredCount === 0" class="data-state">
                <span class="state-icon"
                  ><svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M4 6h16v14H4zM8 3h8v6H8zM8 14h8" /></svg
                ></span>
                <strong>{{ t("inventoryPage.noData") }}</strong>
                <p>{{ copy.emptyHint }}</p>
              </div>

              <div v-else class="table-scroll">
                <table>
                  <thead>
                    <tr>
                      <th class="col-sticky--first">#</th>
                      <th class="col-sticky--second">
                        {{ t("inventoryPage.colItemInfo") }}
                      </th>
                      <th>{{ t("inventoryPage.colRack") }}</th>
                      <th class="col-num">{{ t("inventoryPage.colRem") }}</th>
                      <th>{{ t("inventoryPage.colImportDate") }}</th>
                      <th>{{ t("inventoryPage.colBuyNo") }}</th>
                      <th>{{ t("inventoryPage.colMemo") }}</th>
                      <th v-if="mode === 'month' && inventoryRole === 'unit'">
                        {{ copy.colUnitDate }}
                      </th>
                      <th v-if="mode === 'month' && inventoryRole === 'unit'">
                        {{ copy.colUnitUser }}
                      </th>
                      <th
                        v-if="
                          mode === 'month' && inventoryRole === 'accounting'
                        "
                      >
                        {{ copy.colAccDate }}
                      </th>
                      <th
                        v-if="
                          mode === 'month' && inventoryRole === 'accounting'
                        "
                      >
                        {{ copy.colAccUser }}
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    <template
                      v-for="(item, index) in pagedList"
                      :key="item._key"
                    >
                      <tr :class="{ scanned: isItemChecked(item) }">
                        <td class="row-index col-sticky--first" rowspan="2">
                          {{ (currentPage - 1) * pageSize + index + 1 }}
                        </td>
                        <td class="col-sticky--second" rowspan="2">
                          <div class="item-cell">
                            <span class="item-code">{{
                              item.CLBH || "—"
                            }}</span>
                            <span
                              v-if="isItemChecked(item)"
                              class="checked-badge"
                            >
                              <svg viewBox="0 0 16 16" aria-hidden="true">
                                <path d="m3 8 3 3 7-7" />
                              </svg>
                              {{ copy.checked }}
                            </span>
                          </div>
                          <small
                            >{{ t("inventoryPage.pack") }}:
                            {{ item.Pack || "—" }}</small
                          >
                        </td>
                        <td rowspan="2">
                          <span class="rack-cell">{{
                            item.Rack || "—"
                          }}</span>
                        </td>
                        <td class="quantity-cell col-num" rowspan="2">
                          {{ item._remQtyFmt }}
                        </td>
                        <td rowspan="2">{{ item.CFMDate || "—" }}</td>
                        <td rowspan="2">{{ item.BuyNos || "—" }}</td>
                        <td class="memo-row-cell">
                          <b>{{ t("inventoryPage.art") }}</b>
                          {{ item.Memo_Article || "—" }}
                        </td>
                        <td
                          v-if="mode === 'month' && inventoryRole === 'unit'"
                          rowspan="2"
                        >
                          {{ item.Round1Date || "—" }}
                        </td>
                        <td
                          v-if="mode === 'month' && inventoryRole === 'unit'"
                          rowspan="2"
                        >
                          {{ item.Round1UserID || "—" }}
                        </td>
                        <td
                          v-if="
                            mode === 'month' && inventoryRole === 'accounting'
                          "
                          rowspan="2"
                        >
                          {{ item.Round2Date || "—" }}
                        </td>
                        <td
                          v-if="
                            mode === 'month' && inventoryRole === 'accounting'
                          "
                          rowspan="2"
                        >
                          {{ item.Round2UserID || "—" }}
                        </td>
                      </tr>
                      <tr :class="{ scanned: isItemChecked(item) }">
                        <td class="memo-row-cell">
                          <b>{{ t("inventoryPage.ry") }}</b>
                          {{ item.Memo_RY || "—" }}
                        </td>
                      </tr>
                    </template>
                  </tbody>
                </table>
              </div>

              <nav
                v-if="showList && totalPages > 1"
                class="pagination"
                :aria-label="t('inventoryPage.page')"
              >
                <button
                  type="button"
                  :disabled="currentPage === 1"
                  @click="goToPage(1)"
                >
                  «
                </button>
                <button
                  type="button"
                  :disabled="currentPage === 1"
                  @click="goToPage(currentPage - 1)"
                >
                  ‹
                </button>
                <button
                  v-for="page in visiblePages"
                  :key="page"
                  type="button"
                  :class="{ active: page === currentPage }"
                  @click="goToPage(page)"
                >
                  {{ page }}
                </button>
                <button
                  type="button"
                  :disabled="currentPage === totalPages"
                  @click="goToPage(currentPage + 1)"
                >
                  ›
                </button>
                <button
                  type="button"
                  :disabled="currentPage === totalPages"
                  @click="goToPage(totalPages)"
                >
                  »
                </button>
                <span>{{ currentPage }} / {{ totalPages }}</span>
              </nav>
            </section>
          </section>

          <section v-else-if="mode === 'month'" class="rack-placeholder">
            <span
              ><svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M4 4h16v16H4zM4 10h16M4 15h16M9 4v16" /></svg
            ></span>
            <div>
              <strong>{{ t("inventoryPage.rackWarningTitle") }}</strong>
              <p>{{ t("inventoryPage.rackWarningSubtitle") }}</p>
            </div>
          </section>
        </template>
      </main>

      <!-- Thông báo nằm dưới header cố định -->
      <Transition name="toast">
        <div
          v-if="snackbar.show"
          class="scan-toast"
          :class="snackbar.color"
          role="status"
        >
          <svg
            v-if="snackbar.color === 'success'"
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
          <span>{{ snackbar.message }}</span>
          <button
            type="button"
            :aria-label="copy.close"
            @click="snackbar.show = false"
          >
            ×
          </button>
        </div>
      </Transition>
    </N8dHeader>
  </div>
</template>

<script setup>
import {
  computed,
  markRaw,
  nextTick,
  onBeforeUnmount,
  onMounted,
  ref,
  shallowRef,
  watch,
} from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import axios from "axios";
import { useInventory } from "@/hooks/useInventory";
import N8dHeader from "@/views/components/N8d/N8dHeader.vue";

const route = useRoute();
const router = useRouter();
const { locale, t } = useI18n();
const {
  inventoryList,
  rackList,
  fetchInventory,
  fetchRackList,
  saveInventory,
} = useInventory();

const locales = [
  { code: "vi", label: "Tiếng Việt", short: "VN" },
  { code: "en", label: "English", short: "EN" },
  { code: "zh", label: "中文", short: "CN" },
];

const interfaceMessages = {
  vi: {
    brandSubtitle: "Kho vật tư",
    modeLabel: "Chế độ kiểm kê",
    modeDay: "Theo ngày",
    modeMonth: "Theo tháng",
    operator: "Người thao tác",
    fullscreen: "Toàn màn hình",
    exitFullscreen: "Thoát toàn màn hình",
    home: "Trang chủ",
    scanReady: "Sẵn sàng quét",
    scanWaiting: "Đang chờ chọn kệ",
    record: "Ghi nhận",
    rowsPerPage: "Số dòng",
    loadingHint: "Dữ liệu đang được đồng bộ từ máy chủ.",
    loadFailed: "Không tải được dữ liệu",
    emptyHint: "Không có vật tư phù hợp với bộ lọc hiện tại.",
    checked: "Đã kiểm",
    close: "Đóng",
    guest: "Chưa đăng nhập",
    sessionMissing: "Chưa có phiên",
    sessionEyebrow: "YÊU CẦU ĐĂNG NHẬP",
    signInTitle: "Đăng nhập để bắt đầu kiểm kê",
    signInDescription:
      "Trang kiểm kê vẫn sẵn sàng. Hãy đăng nhập để tải tồn kho, chọn kệ và ghi nhận barcode.",
    signInAction: "Đăng nhập ngay",
    backToHome: "Về trang chức năng",
    requestFailed: "Không thể tải dữ liệu. Vui lòng thử lại.",
    rackLoadFailed: "Không thể tải danh sách kệ.",
    scanFailed: "Ghi nhận barcode thất bại.",
    sessionExpired: "Phiên đăng nhập đã hết. Vui lòng đăng nhập lại.",
    selectRack: "Vui lòng chọn kệ trước khi quét.",
    scanSuccess: "Quét thành công: {code}",
    periodMonthError: "Chỉ được quét tháng hiện tại: {month}",
    periodDayError: "Chỉ được quét ngày hôm nay: {date}",
    warning: "Cảnh báo",
    round2Confirm: "KT xác nhận",
    round2Confirmed: "Đã xác nhận",
    round2Success: "Đã xác nhận kiểm kê lần 2",
    round2Fail: "Xác nhận thất bại",
    roleUnit: "Đơn vị kiểm kê",
    roleAccounting: "Kế toán",
    colUnitDate: "Ngày ĐV kiểm",
    colUnitUser: "Người ĐV kiểm",
    colAccDate: "Ngày KT xác nhận",
    colAccUser: "Người KT xác nhận",
    syncNeededUnitTitle: "Chưa đồng bộ dữ liệu tháng này.",
    syncNeededUnit:
      "Vui lòng quay lại trang chủ và bấm \"Đồng bộ\" để tải tồn kho trước khi kiểm kê.",
    syncNeededAccounting:
      "Đơn vị chưa đồng bộ dữ liệu tháng này. Vui lòng chờ đơn vị đồng bộ trước khi xác nhận.",
    searchPlaceholder: "Tìm tất cả (mã, kệ, ngày, số lượng...)",
  },
  en: {
    brandSubtitle: "Material warehouse",
    modeLabel: "Inventory mode",
    modeDay: "Daily",
    modeMonth: "Monthly",
    operator: "Operator",
    fullscreen: "Enter fullscreen",
    exitFullscreen: "Exit fullscreen",
    home: "Home",
    scanReady: "Ready to scan",
    scanWaiting: "Waiting for rack",
    record: "Record",
    rowsPerPage: "Rows",
    loadingHint: "Data is being synchronized with the server.",
    loadFailed: "Unable to load data",
    emptyHint: "No materials match the current filter.",
    checked: "Checked",
    close: "Close",
    guest: "Not signed in",
    sessionMissing: "No session",
    sessionEyebrow: "SIGN-IN REQUIRED",
    signInTitle: "Sign in to start inventory",
    signInDescription:
      "The inventory workspace is ready. Sign in to load stock, choose a rack and record barcodes.",
    signInAction: "Sign in now",
    backToHome: "Back to workspace",
    requestFailed: "Unable to load data. Please try again.",
    rackLoadFailed: "Unable to load rack options.",
    scanFailed: "Barcode recording failed.",
    sessionExpired: "Your session has expired. Please sign in again.",
    selectRack: "Please select a rack before scanning.",
    scanSuccess: "Scan successful: {code}",
    periodMonthError: "You can only scan the current month: {month}",
    periodDayError: "You can only scan today: {date}",
    warning: "Warning",
    round2Confirm: "Accounting confirm",
    round2Confirmed: "Confirmed",
    round2Success: "Round 2 confirmed",
    round2Fail: "Confirmation failed",
    roleUnit: "Unit check",
    roleAccounting: "Accounting",
    colUnitDate: "Unit check date",
    colUnitUser: "Unit checker",
    colAccDate: "Accounting confirm date",
    colAccUser: "Accounting confirmer",
    syncNeededUnitTitle: "This month hasn't been synced yet.",
    syncNeededUnit:
      'Please go back to the home page and tap "Sync" to load stock before checking.',
    syncNeededAccounting:
      "The unit hasn't synced this month's data yet. Please wait before confirming.",
    searchPlaceholder: "Search all (code, rack, date, qty...)",
  },
  zh: {
    brandSubtitle: "材料仓库",
    modeLabel: "盘点模式",
    modeDay: "每日",
    modeMonth: "每月",
    operator: "操作员",
    fullscreen: "全屏",
    exitFullscreen: "退出全屏",
    home: "主页",
    scanReady: "可以扫描",
    scanWaiting: "等待选择货架",
    record: "记录",
    rowsPerPage: "行数",
    loadingHint: "正在与服务器同步资料。",
    loadFailed: "无法加载资料",
    emptyHint: "当前筛选条件没有符合的物料。",
    checked: "已盘",
    close: "关闭",
    guest: "尚未登录",
    sessionMissing: "没有会话",
    sessionEyebrow: "需要登录",
    signInTitle: "登录后开始盘点",
    signInDescription:
      "盘点页面已准备好。请登录以加载库存、选择货架并记录条码。",
    signInAction: "立即登录",
    backToHome: "返回功能页",
    requestFailed: "无法加载资料，请重试。",
    rackLoadFailed: "无法加载货架清单。",
    scanFailed: "条码记录失败。",
    sessionExpired: "登录会话已过期，请重新登录。",
    selectRack: "请先选择货架再扫描。",
    scanSuccess: "扫描成功：{code}",
    periodMonthError: "只能扫描当前月份：{month}",
    periodDayError: "只能扫描今天：{date}",
    warning: "警告",
    round2Confirm: "会计确认",
    round2Confirmed: "已确认",
    round2Success: "已确认第二次盘点",
    round2Fail: "确认失败",
    roleUnit: "单位盘点",
    roleAccounting: "会计",
    colUnitDate: "单位盘点日期",
    colUnitUser: "单位盘点人",
    colAccDate: "单位盘点日期",
    colAccUser: "会计确认人",
    syncNeededUnitTitle: "本月尚未同步资料。",
    syncNeededUnit: "请返回首页点击「同步」按钮加载库存后再盘点。",
    syncNeededAccounting: "单位尚未同步本月资料，请等待单位同步后再确认。",
    searchPlaceholder: "搜索全部 (代码、货架、日期、数量...)",
  },
};

const copy = computed(
  () => interfaceMessages[locale.value] || interfaceMessages.vi,
);
const userId = ref("");
const gsbh = ref("");
const hasSession = ref(false);
const mode = ref("day");
const searchValue = ref("");
const rackFilter = ref("");
const statusTab = ref("all");
const barcode = ref("");
const barcodeInput = ref(null);
const scanBusy = ref(false);
const loadingRacks = ref(false);
const racksLoaded = ref(false);
const rackError = ref("");
const loadingInventory = ref(false);
const inventoryLoaded = ref(false);
const inventoryError = ref("");
const showList = ref(true);
const snackbar = ref({ show: false, message: "", color: "success" });
const displayedInventoryList = shallowRef([]);
const currentPage = ref(1);
const pageSize = ref(20);
const pageSizeOptions = [20, 50, 100];
const inventoryRole = ref("unit");

// BIẾN TÌM KIẾM
const searchQuery = ref("");

let snackbarTimer = null;
let pollTimer = null;
let pollingBusy = false;
let fetchAbort = null;
let initialized = false;
let lastSignature = "";

const pad2 = (v) => String(v).padStart(2, "0");
const getLocalDate = () => {
  const d = new Date();
  return `${d.getFullYear()}-${pad2(d.getMonth() + 1)}-${pad2(d.getDate())}`;
};
const getLocalMonth = () => {
  const d = new Date();
  return `${d.getFullYear()}-${pad2(d.getMonth() + 1)}`;
};
const getRackCode = (r = "") =>
  String(r || "")
    .trim()
    .slice(0, 3)
    .toUpperCase();
const decodeTokenPayload = (token) => {
  try {
    const ep = token.split(".")[1] || "";
    if (!ep) return null;
    const n = ep.replace(/-/g, "+").replace(/_/g, "/");
    return JSON.parse(atob(n.padEnd(Math.ceil(n.length / 4) * 4, "=")));
  } catch {
    return null;
  }
};
const hydrateIdentity = () => {
  userId.value = "";
  gsbh.value = "";
  hasSession.value = false;
  const token = localStorage.getItem("token") || "";
  const payload = decodeTokenPayload(token);
  const uid = String(
    payload?.USERID ?? payload?.userid ?? localStorage.getItem("USERID") ?? "",
  ).trim();
  const gsb = String(
    payload?.GSBH ?? payload?.gsbh ?? localStorage.getItem("GSBH") ?? "",
  ).trim();
  if (!token || !uid || !gsb) return false;
  userId.value = uid;
  gsbh.value = gsb;
  localStorage.setItem("USERID", uid);
  localStorage.setItem("GSBH", gsb);
  hasSession.value = true;
  return true;
};
const isItemChecked = (item) => {
  if (mode.value === "month") {
    if (inventoryRole.value === "accounting")
      return Boolean(item?.Round2Date) && Boolean(item?.Round2UserID);
    return Boolean(item?.Round1Date) && Boolean(item?.Round1UserID);
  }
  return String(item?.Status) === "1";
};
const listTitle = computed(() =>
  mode.value === "month"
    ? t("inventoryPage.listTitleMonth")
    : t("inventoryPage.listTitleDay"),
);
const canScanInput = computed(
  () =>
    hasSession.value &&
    (mode.value === "day" || !needsSync.value),
);
const canSubmitScan = computed(
  () =>
    canScanInput.value &&
    Boolean(String(barcode.value || "").trim()) &&
    !scanBusy.value,
);
const showContent = computed(() =>
  mode.value === "day" ? Boolean(searchValue.value) : Boolean(rackFilter.value),
);
const rackOptions = computed(() => {
  if (mode.value === "month") {
    return (rackList.value || [])
      .map((i) => {
        const r = String(i?.Rack ?? i?.rack ?? i ?? "").trim();
        return {
          label: r,
          value: getRackCode(r),
          scanned: Number(i?.ScannedCount) || 0,
          total: Number(i?.TotalCount) || 0,
        };
      })
      .filter((i) => i.value);
  }
  const racks = new Set();
  for (const item of displayedInventoryList.value || []) {
    const r = getRackCode(item.Rack);
    if (r) racks.add(r);
  }
  return Array.from(racks).sort();
});
const needsSync = computed(() => {
  if (
    mode.value !== "month" ||
    !racksLoaded.value ||
    loadingRacks.value ||
    rackError.value
  )
    return false;
  const opts = rackOptions.value;
  if (opts.length === 0) return true;
  return opts.every((r) => (Number(r.total) || 0) === 0);
});
const processedData = computed(() => {
  const q = getRackCode(rackFilter.value);
  const f = [];
  let s = 0,
    rq = 0,
    sq = 0;
  const isRound2Accounting =
    mode.value === "month" && inventoryRole.value === "accounting";
  for (const item of displayedInventoryList.value || []) {
    if (q && q !== "TẤT CẢ" && getRackCode(item.Rack) !== q) continue;
    f.push(item);
    const checked = isItemChecked(item);
    if (checked) s += 1;
    rq += Number(item.RemQty) || 0;
    if (isRound2Accounting) {
      if (Boolean(item?.Round2Date) && Boolean(item?.Round2UserID))
        sq += Number(item.QtyScan) || 0;
    } else {
      sq += Number(item.QtyScan) || 0;
    }
  }
  return {
    filtered: f,
    scanned: s,
    unscanned: f.length - s,
    remainingQuantity: rq,
    scannedQuantity: sq,
  };
});
const stats = computed(() => ({
  rem: processedData.value.remainingQuantity.toLocaleString(),
  scan: processedData.value.scannedQuantity.toLocaleString(),
  diff: (
    processedData.value.remainingQuantity - processedData.value.scannedQuantity
  ).toLocaleString(),
}));

const tabFiltered = computed(() => {
  const l = processedData.value.filtered;
  if (statusTab.value === "scanned") return l.filter((i) => isItemChecked(i));
  if (statusTab.value === "unscanned")
    return l.filter((i) => !isItemChecked(i));
  return l;
});

const searchFiltered = computed(() => {
  const q = searchQuery.value.trim().toLowerCase();
  if (!q) return tabFiltered.value;

  return tabFiltered.value.filter((item) => {
    return [
      item.CLBH,
      item.Pack,
      item.Rack,
      String(item._remQtyFmt || ""),
      item.CFMDate,
      item.BuyNos,
      item.Memo_Article,
      item.Memo_RY,
      item.Round1Date,
      item.Round1UserID,
      item.Round2Date,
      item.Round2UserID,
    ].some((field) =>
      String(field || "")
        .toLowerCase()
        .includes(q),
    );
  });
});

const statusTabs = computed(() => [
  {
    value: "all",
    label: t("inventoryPage.tabAll"),
    count: processedData.value.filtered.length,
  },
  {
    value: "scanned",
    label: t("inventoryPage.tabScanned"),
    count: processedData.value.scanned,
  },
  {
    value: "unscanned",
    label: t("inventoryPage.tabUnscanned"),
    count: processedData.value.unscanned,
  },
]);
const tabFilteredCount = computed(() => searchFiltered.value.length);
const totalPages = computed(() =>
  Math.max(1, Math.ceil(tabFilteredCount.value / pageSize.value)),
);
const pagedList = computed(() => {
  const s = (currentPage.value - 1) * pageSize.value;
  return searchFiltered.value.slice(s, s + pageSize.value);
});
const visiblePages = computed(() => {
  const p = [];
  const s = Math.max(1, currentPage.value - 2);
  const e = Math.min(totalPages.value, currentPage.value + 2);
  for (let i = s; i <= e; i += 1) p.push(i);
  return p;
});

watch(searchQuery, () => {
  currentPage.value = 1;
});

const enrichItem = (item, i) =>
  markRaw({
    ...item,
    _key: `${item.CLBH || i}-${item.Rack || ""}-${item.BuyNos || ""}-${i}`,
    _remQtyFmt: (Number(item.RemQty) || 0).toLocaleString(),
  });
const buildSignature = (l = []) => {
  if (!l.length) return "EMPTY";
  return l
    .map((i, idx) =>
      [
        idx,
        i?.CLBH,
        i?.RemQty,
        i?.QtyScan,
        i?.Status,
        i?.Rack,
        i?.BuyNos,
        i?.CFMDate,
        i?.Round1Date,
        i?.Round1UserID,
        i?.Round2Date,
        i?.Round2UserID,
      ].join("\u001f"),
    )
    .join("\u001e");
};
const syncDisplayedInventory = () => {
  const n = Array.isArray(inventoryList.value) ? inventoryList.value : [];
  const sig = buildSignature(n);
  if (sig === lastSignature) return;
  lastSignature = sig;
  displayedInventoryList.value = n.map(enrichItem);
  if (currentPage.value > totalPages.value) currentPage.value = 1;
};
const resetInventoryView = () => {
  displayedInventoryList.value = [];
  lastSignature = "";
  inventoryLoaded.value = false;
  inventoryError.value = "";
};
const getRequestErrorMessage = (e) =>
  e?.response?.status === 401
    ? copy.value.sessionExpired
    : copy.value.requestFailed;
const showSnackbar = (msg, color = "success") => {
  window.clearTimeout(snackbarTimer);
  snackbar.value = { show: true, message: msg, color };
  snackbarTimer = window.setTimeout(() => {
    snackbar.value.show = false;
  }, 3500);
};
const abortCurrentFetch = () => {
  if (fetchAbort) {
    fetchAbort.abort();
    fetchAbort = null;
  }
  loadingInventory.value = false;
};
const fetchAndSync = async ({ silent = false } = {}) => {
  if (!hasSession.value || !searchValue.value || !gsbh.value) return;
  if (mode.value === "month" && !rackFilter.value) return;
  abortCurrentFetch();
  const c = new AbortController();
  fetchAbort = c;
  const sig = c.signal;
  loadingInventory.value = !silent || !inventoryLoaded.value;
  inventoryError.value = "";
  try {
    await fetchInventory(
      searchValue.value,
      gsbh.value,
      rackFilter.value || "",
      mode.value,
      sig,
    );
    if (!sig.aborted) {
      syncDisplayedInventory();
      inventoryLoaded.value = true;
    }
  } catch (e) {
    if (e?.code !== "ERR_CANCELED" && e?.name !== "AbortError")
      inventoryError.value = getRequestErrorMessage(e);
  } finally {
    if (fetchAbort === c) {
      fetchAbort = null;
      loadingInventory.value = false;
    }
  }
};
const retryInventory = () => fetchAndSync();
const stopPolling = () => {
  if (!pollTimer) return;
  window.clearTimeout(pollTimer);
  pollTimer = null;
};
const schedulePolling = () => {
  stopPolling();
  const tick = async () => {
    const skip =
      !searchValue.value ||
      !hasSession.value ||
      document.hidden ||
      pollingBusy ||
      (mode.value === "month" && !rackFilter.value);
    if (!skip) {
      pollingBusy = true;
      try {
        await fetchAndSync({ silent: true });
      } finally {
        pollingBusy = false;
      }
    }
    pollTimer = window.setTimeout(tick, 5000);
  };
  pollTimer = window.setTimeout(tick, 5000);
};
const loadRackOptions = async () => {
  if (!hasSession.value || mode.value !== "month" || !gsbh.value) return;
  loadingRacks.value = true;
  racksLoaded.value = false;
  rackError.value = "";
  try {
    await fetchRackList(gsbh.value);
    racksLoaded.value = true;
  } catch (e) {
    rackError.value =
      e?.response?.status === 401
        ? copy.value.sessionExpired
        : copy.value.rackLoadFailed;
  } finally {
    loadingRacks.value = false;
  }
};
const goToLogin = async () => {
  stopPolling();
  abortCurrentFetch();
  await router.push({ name: "n8d-login", query: { redirect: route.fullPath } });
};
const setDefaultValueByMode = () => {
  searchValue.value = mode.value === "month" ? getLocalMonth() : getLocalDate();
};
const applyModeAndLoad = async (nm) => {
  mode.value = nm === "month" ? "month" : "day";
  inventoryRole.value =
    mode.value === "month" && route.query.role === "accounting"
      ? "accounting"
      : "unit";
  setDefaultValueByMode();
  rackFilter.value = "";
  currentPage.value = 1;
  racksLoaded.value = false;
  rackError.value = "";
  resetInventoryView();
  stopPolling();
  if (!hasSession.value) return;
  if (mode.value === "month") {
    await loadRackOptions();
  } else {
    await fetchAndSync();
    schedulePolling();
  }
};
const selectMode = (nm) => {
  if (nm === mode.value) return;
  router.replace({ query: { ...route.query, type: nm } });
};
const onPeriodChange = async () => {
  currentPage.value = 1;
  rackFilter.value = "";
  resetInventoryView();
  stopPolling();
  if (!hasSession.value) return;
  if (mode.value === "month") {
    await loadRackOptions();
  } else {
    await fetchAndSync();
    schedulePolling();
  }
};
const onRackSelected = async (v) => {
  currentPage.value = 1;
  resetInventoryView();
  stopPolling();
  abortCurrentFetch();
  if (!v) return;
  await fetchAndSync();
  schedulePolling();
  await nextTick();
  barcodeInput.value?.focus();
};
const canScanSelectedPeriod = () =>
  mode.value === "month"
    ? searchValue.value === getLocalMonth()
    : searchValue.value === getLocalDate();
const getScanPeriodError = () =>
  mode.value === "month"
    ? copy.value.periodMonthError.replace("{month}", getLocalMonth())
    : copy.value.periodDayError.replace("{date}", getLocalDate());
const scanRound2ByBarcode = async (code) => {
  const { data } = await axios.post(
    "/api/v1/scan/number-scan/inventory-month/confirm-round2-scan",
    {
      Barcode: code,
      USERDATE: searchValue.value,
      USERID: userId.value,
      GSBH: gsbh.value,
    },
  );
  return data;
};
const handleBarcodeScan = async () => {
  const code = String(barcode.value || "").trim();
  if (!code || scanBusy.value) {
    barcode.value = "";
    return;
  }
  if (!hydrateIdentity()) {
    barcode.value = "";
    stopPolling();
    abortCurrentFetch();
    showSnackbar(copy.value.sessionExpired, "error");
    return;
  }
  if (!canScanSelectedPeriod()) {
    barcode.value = "";
    showSnackbar(getScanPeriodError(), "error");
    nextTick(() => barcodeInput.value?.focus());
    return;
  }
  const isAccountingRound2 =
    mode.value === "month" && inventoryRole.value === "accounting";
  scanBusy.value = true;
  barcode.value = "";
  try {
    const res = isAccountingRound2
      ? await scanRound2ByBarcode(code)
      : await saveInventory({
          Barcode: code,
          USERDATE: searchValue.value,
          USERID: userId.value,
          GSBH: gsbh.value,
          Rack: rackFilter.value || "",
          YN: "1",
          MODE: mode.value,
        });
    if (res?.type === "warning") showSnackbar(copy.value.warning, "warning");
    else if (isAccountingRound2)
      showSnackbar(copy.value.round2Success, "success");
    else
      showSnackbar(
        copy.value.scanSuccess.replace("{code}", res?.data || code),
        "success",
      );
    await fetchAndSync();
  } catch (e) {
    showSnackbar(getRequestErrorMessage(e), "error");
  } finally {
    scanBusy.value = false;
    nextTick(() => barcodeInput.value?.focus());
  }
};
const setStatusTab = (v) => {
  statusTab.value = v;
  currentPage.value = 1;
};
const toggleListVisibility = () => {
  showList.value = !showList.value;
  if (showList.value && currentPage.value > totalPages.value)
    currentPage.value = 1;
};
const onPageSizeChange = () => {
  currentPage.value = 1;
};
const goToPage = (p) => {
  currentPage.value = Math.max(1, Math.min(p, totalPages.value));
};
const onVisibilityChange = () => {
  if (hasSession.value && !document.hidden && searchValue.value)
    fetchAndSync({ silent: true });
};
watch(
  () => `${route.query.type || ""}|${route.query.role || ""}`,
  async () => {
    if (initialized) await applyModeAndLoad(route.query.type);
  },
);
onMounted(async () => {
  const sl = localStorage.getItem("APP_LOCALE");
  if (locales.some((i) => i.code === sl)) locale.value = sl;
  hydrateIdentity();
  document.addEventListener("visibilitychange", onVisibilityChange);
  initialized = true;
  await applyModeAndLoad(route.query.type);
});
onBeforeUnmount(() => {
  window.clearTimeout(snackbarTimer);
  stopPolling();
  abortCurrentFetch();
  document.removeEventListener("visibilitychange", onVisibilityChange);
});
</script>

<style scoped>
/* ============================================================
   THIẾT KẾ MỚI — "OCEAN DECK" (Tông Xanh Biển & Đã Kiểm Xanh Lá Nhạt)
   ============================================================ */
*,
*::before,
*::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
button,
input,
select,
textarea {
  font: inherit;
  color: inherit;
}
button {
  background: none;
  border: none;
  cursor: pointer;
  -webkit-appearance: none;
  appearance: none;
  color: inherit;
}
img {
  display: block;
  max-width: 100%;
}
::selection {
  background: rgba(2, 132, 199, 0.18);
  color: #075985;
}

.inventory-page {
  --ink: #0b1220;
  --ink-soft: #3f4a63;
  --muted: #8a93a8;
  --paper: #f0f9ff;
  --surface: #ffffff;
  --surface-tint: #fbfcff;
  --line: rgba(15, 23, 42, 0.09);
  --line-strong: rgba(15, 23, 42, 0.16);

  /* MÀU XANH BIỂN CHỦ ĐẠO */
  --accent: #0284c7;
  --accent-2: #06b6d4;
  --accent-soft: #f0f9ff;
  --accent-dark: #0369a1;
  --accent-glow: rgba(2, 132, 199, 0.32);
  --accent2-glow: rgba(6, 182, 212, 0.28);

  /* MÀU XANH LÁ VÀNG (LIME) CHO ĐÃ KIỂM */
  --success: #65a30d;
  --success-soft: #f7fee7;
  --success-dark: #3f6212;
  --success-glow: rgba(101, 163, 13, 0.35);

  --warning: #d97706;
  --warning-soft: #fff7ed;
  --warning-dark: #b45309;
  --warning-glow: rgba(217, 119, 6, 0.28);
  --danger: #e11d48;
  --danger-soft: #fff1f2;
  --danger-dark: #be123c;
  --danger-glow: rgba(225, 29, 72, 0.28);
  --font-display: "Space Grotesk", "Inter", -apple-system, BlinkMacSystemFont,
    "Segoe UI", Roboto, sans-serif;
  --font-mono: "JetBrains Mono", ui-monospace, "SFMono-Regular", Menlo,
    monospace;
  --shadow-xs: 0 1px 2px rgba(15, 23, 42, 0.06);
  --shadow-sm: 0 3px 10px -2px rgba(15, 23, 42, 0.08),
    0 1px 3px rgba(15, 23, 42, 0.06);
  --shadow-md: 0 16px 36px -18px rgba(30, 41, 90, 0.35),
    0 3px 10px -2px rgba(15, 23, 42, 0.06);
  --shadow-lg: 0 26px 60px -20px rgba(30, 41, 90, 0.4);
  --radius-lg: 20px;
  --radius-md: 14px;
  --radius-sm: 10px;
  color-scheme: light;
  position: relative;
  isolation: isolate;
  height: 100dvh;
  height: 100vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background: radial-gradient(
      1000px 520px at 6% -12%,
      rgba(2, 132, 199, 0.11),
      transparent 60%
    ),
    radial-gradient(
      760px 460px at 104% -6%,
      rgba(6, 182, 212, 0.1),
      transparent 55%
    ),
    radial-gradient(
      900px 560px at 46% 118%,
      rgba(2, 132, 199, 0.07),
      transparent 62%
    ),
    var(--paper);
  color: var(--ink);
  font-family: var(--font-display);
  -webkit-font-smoothing: antialiased;
}

.header-context {
  display: flex;
  align-items: baseline;
  min-width: 0;
  gap: 10px;
}
.header-context span {
  color: rgba(186, 230, 253, 0.78);
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
}
.header-context strong {
  overflow: hidden;
  max-width: 220px;
  color: #f0f9ff;
  font-size: 13px;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 700;
  text-shadow: 0 1px 12px rgba(56, 189, 248, 0.35);
}
.header-context.offline strong {
  color: #fca5a5;
  text-shadow: 0 1px 12px rgba(248, 113, 113, 0.3);
}

.user-chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  max-width: 160px;
  overflow: hidden;
  padding: 7px 14px;
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  color: #e0f2fe;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
  text-overflow: ellipsis;
  flex-shrink: 0;
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}
.user-chip:hover {
  background: rgba(255, 255, 255, 0.14);
  color: #fff;
  box-shadow: 0 6px 18px -6px rgba(0, 0, 0, 0.35);
  transform: translateY(-1px);
}
.status-dot {
  width: 7px;
  height: 7px;
  flex: 0 0 auto;
  border-radius: 50%;
  background: #34d399;
  box-shadow: 0 0 0 3px rgba(52, 211, 153, 0.25),
    0 0 8px rgba(52, 211, 153, 0.7);
}
.user-chip.offline .status-dot {
  background: #f87171;
  box-shadow: 0 0 0 3px rgba(248, 113, 113, 0.25),
    0 0 8px rgba(248, 113, 113, 0.7);
}

.inventory-workspace {
  flex: 1 1 0%;
  display: flex;
  flex-direction: column;
  min-height: 0;
  width: 100%;
  padding: 26px 30px 30px;
  overflow-y: auto;
  overflow-x: hidden;
  scrollbar-gutter: stable;
}
.inventory-workspace::-webkit-scrollbar {
  width: 8px;
}
.inventory-workspace::-webkit-scrollbar-track {
  background: transparent;
}
.inventory-workspace::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, var(--accent), var(--accent-2));
  opacity: 0.5;
  border-radius: 99px;
  border: 2px solid transparent;
  background-clip: padding-box;
}
.inventory-workspace::-webkit-scrollbar-thumb:hover {
  opacity: 0.8;
}

.control-panel,
.list-panel {
  position: relative;
  border: 1px solid var(--line);
  border-radius: var(--radius-lg);
  background: linear-gradient(
    180deg,
    var(--surface-tint) 0%,
    var(--surface) 55%
  );
  box-shadow: var(--shadow-md), inset 0 1px 0 rgba(255, 255, 255, 0.7);
  transition: box-shadow 0.25s ease;
}
.control-panel {
  overflow: hidden;
}
.list-panel {
  overflow-x: clip;
  overflow-y: visible;
}
.control-panel::before,
.list-panel::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  z-index: 3;
  background: linear-gradient(
    90deg,
    var(--accent) 0%,
    #38bdf8 45%,
    var(--accent-2) 100%
  );
  box-shadow: 0 2px 16px -2px var(--accent-glow);
}

.control-panel {
  flex-shrink: 0;
  margin-bottom: 22px;
  padding: 22px 24px;
}
.control-grid {
  display: grid;
  grid-template-columns: minmax(180px, 0.8fr) minmax(190px, 0.9fr) minmax(
      340px,
      2fr
    );
  gap: 16px;
}
.field-group {
  display: grid;
  min-width: 0;
  gap: 8px;
  color: var(--ink-soft);
  font-size: 11px;
  font-weight: 800;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.input-shell,
.scan-shell {
  display: flex;
  min-width: 0;
  min-height: 48px;
  align-items: center;
  border: 1px solid var(--line-strong);
  border-radius: var(--radius-sm);
  background: var(--surface);
  box-shadow: var(--shadow-xs);
  transition: border-color 0.15s ease, box-shadow 0.15s ease,
    transform 0.15s ease;
}
.input-shell {
  gap: 10px;
  padding: 0 14px;
}
.input-shell:focus-within,
.scan-shell:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 0 4px rgba(2, 132, 199, 0.14), var(--shadow-sm);
  transform: translateY(-1px);
}
.input-shell--required {
  border-color: var(--warning);
  background: var(--warning-soft);
  box-shadow: 0 0 0 3px rgba(217, 119, 6, 0.1);
}
.input-shell svg,
.scan-shell > svg {
  width: 17px;
  flex: 0 0 auto;
  fill: none;
  stroke: var(--muted);
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
  transition: stroke 0.15s;
}
.input-shell:focus-within svg {
  stroke: var(--accent);
}
.input-shell input,
.input-shell select,
.scan-shell input {
  min-width: 0;
  height: 46px;
  flex: 1;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--ink);
  font-size: 14px;
  font-weight: 600;
}
.input-shell select option {
  background: #fff;
  color: var(--ink);
}

.input-shell--picker {
  padding-right: 6px;
}
.period-picker,
.rack-select {
  flex: 1;
  min-width: 0;
}
.period-picker :deep(.el-input__wrapper),
.rack-select :deep(.el-select__wrapper) {
  min-height: 46px;
  padding: 0;
  border: 0;
  box-shadow: none !important;
  background: transparent;
}
.period-picker :deep(.el-input__inner) {
  height: 46px;
  padding: 0;
  border: 0;
  background: transparent;
  color: var(--ink);
  font-size: 14px;
  font-weight: 600;
}
.period-picker :deep(.el-input__prefix),
.period-picker :deep(.el-input__suffix) {
  display: none;
}
.rack-select :deep(.el-select__selected-item),
.rack-select :deep(.el-select__placeholder) {
  color: var(--ink);
  font-size: 14px;
  font-weight: 600;
}
.rack-select :deep(.el-select__caret) {
  color: var(--muted);
}
.rack-select.is-disabled :deep(.el-select__wrapper) {
  background: transparent;
}

.scan-label {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.scan-label small {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
  font-weight: 800;
  color: var(--muted);
  text-transform: none;
  letter-spacing: 0;
}
.scan-label small i {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--muted);
}
.scan-label small.ready {
  color: var(--success);
}
.scan-label small.ready i {
  background: var(--success);
  box-shadow: 0 0 0 3px var(--success-soft), 0 0 8px var(--success-glow);
  animation: pulseDot 1.8s ease-in-out infinite;
}
@keyframes pulseDot {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.45;
  }
}

.scan-shell {
  gap: 0;
  padding: 0;
  border-color: var(--line-strong);
  background: #fff;
  border-width: 1.5px;
}
.scan-shell.disabled {
  opacity: 0.55;
  pointer-events: none;
}
.scan-shell.is-ready {
  border-color: var(--success);
  box-shadow: 0 0 0 3px var(--success-soft);
}
.scan-shell.is-ready:focus-within {
  border-color: var(--success-dark);
  box-shadow: 0 0 0 4px var(--success-glow), var(--shadow-sm);
}
.scan-shell.is-ready > svg {
  stroke: var(--success);
}
.scan-shell.is-busy {
  border-color: var(--accent);
  box-shadow: 0 0 0 4px rgba(2, 132, 199, 0.14);
}

.scan-shell input {
  padding: 0 14px;
  font-family: var(--font-mono);
  letter-spacing: 0.02em;
}
.scan-shell button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 100%;
  padding: 0 22px;
  border: none;
  border-left: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 0 9px 9px 0;
  background: linear-gradient(135deg, var(--accent), var(--accent-dark));
  color: #fff;
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0.01em;
  white-space: nowrap;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.15);
  transition: background 0.15s ease, filter 0.15s ease;
}
.scan-shell button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.scan-shell button:not(:disabled):hover {
  filter: brightness(1.08);
}
.scan-shell button:not(:disabled):active {
  filter: brightness(0.96);
}
.scan-shell button svg {
  width: 15px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2.5;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.control-status {
  margin-top: 15px;
}
.inline-state {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 11px 16px;
  border: 1px solid transparent;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 700;
  box-shadow: var(--shadow-xs);
}
.inline-state.info {
  background: var(--accent-soft);
  border-color: rgba(2, 132, 199, 0.22);
  color: var(--accent-dark);
}
.inline-state.error {
  background: var(--danger-soft);
  border-color: rgba(225, 29, 72, 0.22);
  color: var(--danger-dark);
}
.inline-state.warning {
  background: var(--warning-soft);
  border-color: rgba(217, 119, 6, 0.25);
  color: var(--warning-dark);
}
.inline-state svg {
  width: 18px;
  height: 18px;
  flex: 0 0 auto;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
}
.inline-state button {
  padding: 5px 12px;
  border: 1px solid currentColor;
  border-radius: 8px;
  font-size: 11px;
  font-weight: 800;
  transition: 0.15s;
  background: rgba(255, 255, 255, 0.7);
}
.inline-state button:hover {
  background: rgba(255, 255, 255, 0.98);
}

.results-area {
  display: flex;
  flex: 1;
  min-height: 0;
  flex-direction: column;
  gap: 22px;
}
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 18px;
  flex-shrink: 0;
}
.stat-card {
  position: relative;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  background: var(--surface);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.stat-card::after {
  content: "";
  position: absolute;
  inset: 0 auto 0 0;
  width: 4px;
}
.stat-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
}
.stat-icon {
  display: grid;
  width: 52px;
  height: 52px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-xs);
}
.stat-icon svg {
  width: 24px;
  fill: none;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.stat-card--stock {
  background: linear-gradient(160deg, var(--accent-soft) 0%, #fff 68%);
  border-color: rgba(2, 132, 199, 0.18);
}
.stat-card--stock::after {
  background: linear-gradient(180deg, var(--accent), var(--accent-2));
}
.stat-card--stock .stat-icon {
  background: linear-gradient(135deg, var(--accent), var(--accent-2));
  color: #fff;
  box-shadow: 0 0 0 6px rgba(2, 132, 199, 0.09), var(--shadow-xs);
}
.stat-card--stock .stat-icon svg {
  stroke: currentColor;
}
.stat-card--stock .stat-copy small {
  color: var(--accent-dark);
}
.stat-card--stock .stat-copy strong {
  color: var(--accent-dark);
}

/* Thẻ Đã Kiểm màu Xanh Lá Vàng */
.stat-card--scanned {
  background: linear-gradient(160deg, var(--success-soft) 0%, #fff 68%);
  border-color: rgba(101, 163, 13, 0.18);
}
.stat-card--scanned::after {
  background: linear-gradient(180deg, var(--success), #a3e635);
}
.stat-card--scanned .stat-icon {
  background: linear-gradient(135deg, var(--success), #a3e635);
  color: #fff;
  box-shadow: 0 0 0 6px rgba(101, 163, 13, 0.09), var(--shadow-xs);
}
.stat-card--scanned .stat-icon svg {
  stroke: currentColor;
}
.stat-card--scanned .stat-copy small {
  color: var(--success-dark);
}
.stat-card--scanned .stat-copy strong {
  color: var(--success-dark);
}

.stat-card--diff {
  background: linear-gradient(160deg, var(--warning-soft) 0%, #fff 68%);
  border-color: rgba(217, 119, 6, 0.2);
}
.stat-card--diff::after {
  background: linear-gradient(180deg, var(--warning), #fbbf24);
}
.stat-card--diff .stat-icon {
  background: linear-gradient(135deg, var(--warning), #f59e0b);
  color: #fff;
  box-shadow: 0 0 0 6px rgba(217, 119, 6, 0.09), var(--shadow-xs);
}
.stat-card--diff .stat-icon svg {
  stroke: currentColor;
}
.stat-card--diff .stat-copy small {
  color: var(--warning-dark);
}
.stat-card--diff .stat-copy strong {
  color: var(--warning-dark);
}
.stat-copy {
  display: grid;
  gap: 4px;
  min-width: 0;
}
.stat-copy small {
  color: var(--muted);
  font-size: 11px;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}
.stat-copy strong {
  color: var(--ink);
  font-size: 28px;
  font-weight: 800;
  letter-spacing: -0.02em;
  line-height: 1;
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
}

.list-panel {
  display: flex;
  flex: 1;
  min-height: 0;
  flex-direction: column;
  overflow-x: clip;
}
.list-toolbar {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 18px 24px;
  border-bottom: 1px solid var(--line);
  flex-shrink: 0;
  background: linear-gradient(180deg, #f0f9ff 0%, var(--paper) 100%);
}
.list-toolbar::after {
  content: "";
  position: absolute;
  left: 24px;
  right: 24px;
  bottom: -1px;
  height: 1px;
  background: linear-gradient(
    90deg,
    rgba(2, 132, 199, 0.35),
    rgba(6, 182, 212, 0.35),
    transparent
  );
}
.list-heading {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}
.rack-badge {
  display: inline-flex;
  align-items: center;
  padding: 6px 13px;
  border: 1px solid rgba(2, 132, 199, 0.22);
  border-radius: 8px;
  background: linear-gradient(135deg, var(--accent-soft), #fff);
  color: var(--accent-dark);
  font-size: 13px;
  font-weight: 800;
  font-family: var(--font-mono);
  flex-shrink: 0;
  box-shadow: var(--shadow-xs), inset 0 1px 0 rgba(255, 255, 255, 0.8);
}
.list-heading div {
  display: grid;
  gap: 2px;
}
.list-heading h2 {
  font-size: 16px;
  font-weight: 800;
  letter-spacing: -0.01em;
  background: linear-gradient(135deg, var(--ink) 30%, var(--accent-dark) 130%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}
.list-heading span {
  color: var(--muted);
  font-size: 12px;
  font-weight: 600;
}

.list-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}
.page-size {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--muted);
  font-size: 11px;
  font-weight: 800;
}
.page-size select {
  padding: 7px 11px;
  border: 1px solid var(--line-strong);
  border-radius: 8px;
  background: #fff;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: var(--shadow-xs);
  transition: border-color 0.15s ease;
}
.page-size select:focus {
  border-color: var(--accent);
  outline: none;
}
.visibility-action {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 8px 15px;
  border: 1px solid var(--line-strong);
  border-radius: 8px;
  background: #fff;
  color: var(--ink-soft);
  font-size: 12px;
  font-weight: 800;
  box-shadow: var(--shadow-xs);
  transition: border-color 0.15s ease, color 0.15s ease, transform 0.15s ease;
}
.visibility-action:hover {
  border-color: var(--accent);
  color: var(--accent);
  transform: translateY(-1px);
}
.visibility-action svg {
  width: 16px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

/* CSS CHO Ô TÌM KIẾM */
.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 260px;
  height: 38px;
  padding: 0 12px;
  border-radius: 8px;
  border: 1px solid var(--line-strong);
  background: #fff;
  box-shadow: var(--shadow-xs);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}
.search-box:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 0 4px rgba(2, 132, 199, 0.14);
}
.search-box svg {
  width: 16px;
  flex: 0 0 auto;
  fill: none;
  stroke: var(--muted);
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}
.search-box:focus-within svg {
  stroke: var(--accent);
}
.search-box input {
  flex: 1;
  min-width: 0;
  height: 100%;
  border: 0;
  outline: 0;
  background: transparent;
  font-size: 13px;
  font-weight: 500;
  color: var(--ink);
}
.search-box input::placeholder {
  color: var(--muted);
}


.status-tabs {
  display: flex;
  gap: 8px;
  padding: 14px 22px;
  flex-shrink: 0;
  background: linear-gradient(180deg, #f0f9ff, var(--paper));
  border-bottom: 1px solid var(--line);
}
.status-tabs button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 9px 18px;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.7);
  color: var(--muted);
  font-size: 13px;
  font-weight: 700;
  transition: color 0.15s ease, background 0.15s ease, box-shadow 0.15s ease,
    transform 0.15s ease, border-color 0.15s ease;
}
.status-tabs button strong {
  font-family: var(--font-mono);
  font-size: 11.5px;
  font-weight: 800;
  padding: 1px 8px;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.06);
}
.status-tabs button:not(.active):hover {
  color: var(--ink);
  background: #fff;
  border-color: var(--line-strong);
  transform: translateY(-1px);
}
.status-tabs button.active {
  background: linear-gradient(135deg, var(--ink), #1e293b);
  color: #fff;
  font-weight: 800;
  border-color: transparent;
  box-shadow: 0 10px 20px -9px rgba(15, 23, 42, 0.45);
  transform: translateY(-1px);
}
.status-tabs button.active strong {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}
.status-tabs button.scanned.active {
  background: linear-gradient(135deg, var(--success), var(--success-dark));
  box-shadow: 0 10px 20px -9px var(--success-glow);
}
.status-tabs button.unscanned.active {
  background: linear-gradient(135deg, var(--warning), var(--warning-dark));
  box-shadow: 0 10px 20px -9px var(--warning-glow);
}

/* ============================================================
   TABLE — Viền đen đậm rõ nét, đã kiểm màu xanh nhạt (ko xen kẽ)
   ============================================================ */
.table-scroll {
  flex: 1 1 auto;
  min-height: 0;
  max-height: min(640px, 58vh);
  overflow-x: auto;
  overflow-y: auto;
  padding: 0;
  scrollbar-gutter: stable;
  background: #fff;
}
.table-scroll::-webkit-scrollbar {
  height: 8px;
}
.table-scroll::-webkit-scrollbar-thumb {
  background: var(--line-strong);
  border-radius: 99px;
}

table {
  border-collapse: collapse;
  width: 100%;
  min-width: 1000px;
  font-variant-numeric: tabular-nums;
}

thead {
  position: sticky;
  top: 0;
  z-index: 4;
}

/* Nâng cấp Header Bảng */
th {
  background: linear-gradient(180deg, #0284c7 0%, #0369a1 100%);
  color: #ffffff;
  padding: 16px;
  border: 1.5px solid #0f172a;
  border-bottom: 3px solid #0f172a;
  font-size: 11px;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  text-align: left;
  white-space: nowrap;
  user-select: none;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

th.col-num {
  text-align: right;
}

/* Đã bỏ màu xen kẽ (nth-child even) -> Chưa kiểm thì trắng tinh khiết */
td {
  background: #ffffff;
  padding: 14px 16px;
  border: 1px solid #0f172a;
  font-size: 13.5px;
  vertical-align: top;
  line-height: 1.5;
  color: var(--ink-soft);
  transition: background 0.15s ease;
}

tbody tr {
  transition: box-shadow 0.15s ease, background 0.15s ease;
}
tbody tr:hover td {
  background: #f8fafc;
}

/* Dòng đã quét: Chỉ một màu xanh lá nhạt, KHÔNG XEN KẼ */
tbody tr.scanned td {
  background: #f7fee7;
  color: #3f6212;
  font-weight: 600;
  border-color: #0f172a;
}
tbody tr.scanned:hover td {
  background: #ecfccb;
}
tbody tr.scanned td.col-sticky--first {
  box-shadow: inset 4px 0 0 var(--success);
}

.col-sticky--first {
  position: sticky;
  left: 0;
  z-index: 2;
  width: 56px;
  min-width: 56px;
  max-width: 56px;
  text-align: center;
}

/* Cố định header của 2 cột đầu khi cuộn ngang */
thead .col-sticky--first {
  z-index: 8;
  background: linear-gradient(180deg, #0284c7 0%, #0369a1 100%);
}
thead .col-sticky--second {
  z-index: 8;
  background: linear-gradient(180deg, #0284c7 0%, #0369a1 100%);
}

td.col-sticky--first {
  font-family: var(--font-mono);
  font-size: 12.5px;
  font-weight: 800;
  color: var(--muted);
}

.col-sticky--second {
  position: sticky;
  left: 56px;
  z-index: 2;
  min-width: 230px;
  border-right: 2px solid #0f172a;
}

tbody tr:hover .col-sticky--first,
tbody tr:hover .col-sticky--second {
  background: #f8fafc;
}
tbody tr.scanned .col-sticky--first,
tbody tr.scanned .col-sticky--second {
  background: #f7fee7;
}
tbody tr.scanned:hover .col-sticky--first,
tbody tr.scanned:hover .col-sticky--second {
  background: #ecfccb;
}
td.col-sticky--first,
td.col-sticky--second {
  background: #ffffff;
}

.item-cell {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.item-code {
  font-family: var(--font-mono);
  font-size: 13.5px;
  font-weight: 800;
  color: var(--ink);
}

.checked-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 3px 10px;
  border: 1px solid #84cc16;
  border-radius: 999px;
  background: linear-gradient(135deg, #f7fee7, #fff);
  color: #4d7c0f;
  font-size: 10px;
  font-weight: 800;
}
.checked-badge svg {
  width: 11px;
  fill: none;
  stroke: currentColor;
  stroke-width: 3;
  stroke-linecap: round;
  stroke-linejoin: round;
}

td small {
  color: var(--muted);
  font-size: 12px;
}
.rack-cell {
  display: inline-block;
  padding: 5px 12px;
  border: 1px solid #0ea5e9;
  border-radius: 8px;
  background: #ecfeff;
  font-family: var(--font-mono);
  font-size: 12.5px;
  font-weight: 800;
  color: #0369a1;
}
.quantity-cell {
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 800;
  color: var(--ink);
}
.col-num {
  text-align: right;
  white-space: nowrap;
}
.memo-row-cell {
  padding-top: 8px;
  padding-bottom: 8px;
  font-size: 12.5px;
  line-height: 1.4;
}
.memo-row-cell b {
  margin-right: 4px;
  color: var(--ink);
  font-weight: 800;
}

.pagination {
  flex-shrink: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  padding: 15px 20px;
  background: var(--paper);
  border-top: 2px solid #0f172a;
  flex-shrink: 0;
}
.pagination button {
  min-width: 38px;
  height: 38px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--line-strong);
  border-radius: 10px;
  background: #fff;
  color: var(--ink-soft);
  font-size: 13px;
  font-weight: 800;
  font-family: var(--font-mono);
  box-shadow: var(--shadow-xs);
  transition: border-color 0.15s ease, color 0.15s ease, transform 0.15s ease,
    box-shadow 0.15s ease;
}
.pagination button:hover:not(:disabled) {
  border-color: var(--accent);
  color: var(--accent);
  transform: translateY(-1px);
}
.pagination button.active {
  background: linear-gradient(135deg, var(--accent), var(--accent-dark));
  color: #fff;
  border-color: transparent;
  box-shadow: 0 8px 18px -8px var(--accent-glow);
}
.pagination button:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}
.pagination span {
  margin-left: 10px;
  color: var(--muted);
  font-size: 13px;
  font-weight: 700;
  font-family: var(--font-mono);
}

.data-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  padding: 64px 20px;
  text-align: center;
  background: #fff;
}
.data-alert {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 14px 20px;
  background: var(--danger-soft);
  border-bottom: 1px solid rgba(225, 29, 72, 0.2);
  flex-shrink: 0;
  flex-wrap: wrap;
}
.data-alert svg {
  width: 20px;
  fill: none;
  stroke: var(--danger);
  stroke-width: 2;
}
.data-alert span {
  font-size: 13px;
  color: var(--danger-dark);
  font-weight: 700;
}
.data-alert button {
  padding: 5px 14px;
  border: 1px solid var(--danger);
  border-radius: 8px;
  color: var(--danger);
  font-size: 12px;
  font-weight: 800;
  background: #fff;
}
.data-alert button:hover {
  background: var(--danger-soft);
}

.state-icon {
  display: grid;
  width: 62px;
  height: 62px;
  place-items: center;
  border: 1px solid var(--line);
  border-radius: 17px;
  background: var(--paper);
  color: var(--muted);
  box-shadow: var(--shadow-xs);
}
.state-icon svg {
  width: 26px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.5;
  stroke-linecap: round;
  stroke-linejoin: round;
}
.data-state strong {
  font-size: 15px;
  color: var(--ink);
  font-weight: 800;
}
.data-state p {
  font-size: 13px;
  color: var(--muted);
  max-width: 400px;
  line-height: 1.6;
}
.data-state--error .state-icon {
  background: var(--danger-soft);
  color: var(--danger);
  border-color: rgba(225, 29, 72, 0.2);
}
.data-state button {
  margin-top: 6px;
  padding: 10px 24px;
  border: 1px solid var(--line-strong);
  border-radius: 10px;
  background: #fff;
  color: var(--ink-soft);
  font-size: 13px;
  font-weight: 800;
  box-shadow: var(--shadow-xs);
  transition: border-color 0.15s ease, color 0.15s ease;
}
.data-state button:hover {
  border-color: var(--accent);
  color: var(--accent);
}

.refresh-line {
  height: 2.5px;
  flex-shrink: 0;
  background: var(--line);
  overflow: hidden;
}
.refresh-line span {
  display: block;
  width: 30%;
  height: 100%;
  background: linear-gradient(90deg, var(--accent), var(--accent-2));
  animation: refreshSlide 1.5s ease-in-out infinite;
}
@keyframes refreshSlide {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(400%);
  }
}

.spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2.5px solid rgba(255, 255, 255, 0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}
.inline-state .spinner {
  border-color: rgba(2, 132, 199, 0.2);
  border-top-color: var(--accent);
}
.large-spinner {
  display: block;
  width: 38px;
  height: 38px;
  border: 3px solid var(--line);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.rack-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 18px;
  padding: 58px 40px;
  text-align: center;
  border: 1px solid var(--line);
  border-radius: var(--radius-lg);
  background: linear-gradient(180deg, var(--surface-tint), var(--surface));
  box-shadow: var(--shadow-md);
}
.rack-placeholder > span {
  display: grid;
  width: 74px;
  height: 74px;
  place-items: center;
  border-radius: 19px;
  background: linear-gradient(135deg, var(--warning), #f59e0b);
  color: #fff;
  box-shadow: 0 12px 26px -10px var(--warning-glow);
}
.rack-placeholder svg {
  width: 32px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.5;
  stroke-linecap: round;
  stroke-linejoin: round;
}
.rack-placeholder strong {
  font-size: 17px;
  color: var(--ink);
  letter-spacing: -0.01em;
  font-weight: 800;
}
.rack-placeholder p {
  font-size: 14px;
  color: var(--muted);
  max-width: 420px;
  line-height: 1.6;
}

/* THÔNG BÁO NẰM DƯỚI HEADER CỐ ĐỊNH CAO 70PX */
.scan-toast {
  position: fixed;
  top: 86px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 100;
  display: flex;
  align-items: center;
  gap: 12px;
  max-width: 560px;
  padding: 15px 22px;
  border-radius: 14px;
  background: rgba(11, 18, 32, 0.88);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  color: #f8fafc;
  font-size: 14px;
  font-weight: 700;
  box-shadow: 0 20px 44px -14px rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.12);
}
.scan-toast.success {
  box-shadow: 0 20px 44px -14px rgba(0, 0, 0, 0.5),
    0 0 0 1px rgba(132, 204, 22, 0.18) inset;
}
.scan-toast.warning {
  box-shadow: 0 20px 44px -14px rgba(0, 0, 0, 0.5),
    0 0 0 1px rgba(251, 191, 36, 0.18) inset;
}
.scan-toast.error {
  box-shadow: 0 20px 44px -14px rgba(0, 0, 0, 0.5),
    0 0 0 1px rgba(248, 113, 113, 0.18) inset;
}
.scan-toast.success svg {
  width: 20px;
  fill: none;
  stroke: #a3e635;
  stroke-width: 2;
  filter: drop-shadow(0 0 6px rgba(132, 204, 22, 0.6));
}
.scan-toast.warning svg {
  width: 20px;
  fill: none;
  stroke: #fbbf24;
  stroke-width: 2;
  filter: drop-shadow(0 0 6px rgba(251, 191, 36, 0.6));
}
.scan-toast.error svg {
  width: 20px;
  fill: none;
  stroke: #f87171;
  stroke-width: 2;
  filter: drop-shadow(0 0 6px rgba(248, 113, 113, 0.6));
}
.scan-toast span {
  flex: 1;
}
.scan-toast button {
  width: 28px;
  height: 28px;
  display: grid;
  place-items: center;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  color: #94a3b8;
  font-size: 18px;
  transition: 0.15s;
}
.scan-toast button:hover {
  background: rgba(255, 255, 255, 0.18);
  color: #fff;
}
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.25, 1, 0.5, 1);
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(-50%) translateY(-20px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-10px);
}

@media (max-width: 768px) {
  .scan-toast {
    top: 82px;
    right: 16px;
    left: 16px;
    max-width: none;
    transform: none;
  }

  .toast-enter-from,
  .toast-leave-to {
    transform: translateY(-10px);
  }
}

.auth-required {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: center;
  gap: 28px;
  padding: 42px;
  border: 1px solid var(--line);
  border-radius: var(--radius-lg);
  background: linear-gradient(180deg, var(--surface-tint), var(--surface));
  box-shadow: var(--shadow-md);
}
.auth-required__icon {
  display: grid;
  width: 66px;
  height: 66px;
  place-items: center;
  border-radius: 18px;
  background: linear-gradient(135deg, var(--accent), var(--accent-2));
  color: #fff;
  box-shadow: 0 14px 28px -10px var(--accent-glow);
}
.auth-required__icon svg {
  width: 30px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
}
.auth-required__copy {
  display: grid;
  gap: 8px;
}
.auth-required__copy p {
  color: var(--accent);
  font-size: 11px;
  font-weight: 800;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}
.auth-required__copy h2 {
  color: var(--ink);
  font-size: 25px;
  letter-spacing: -0.02em;
  line-height: 1.25;
  font-weight: 800;
}
.auth-required__copy span {
  color: var(--muted);
  font-size: 14px;
  line-height: 1.6;
  max-width: 600px;
}
.auth-required__actions {
  display: flex;
  gap: 10px;
}
.auth-required__actions button {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 0 24px;
  height: 46px;
  border-radius: 11px;
  background: #fff;
  color: var(--ink-soft);
  font-size: 13px;
  font-weight: 800;
  border: 1px solid var(--line-strong);
  box-shadow: var(--shadow-xs);
  transition: border-color 0.15s ease, background 0.15s ease,
    transform 0.15s ease;
}
.auth-required__actions button.primary {
  background: linear-gradient(135deg, var(--accent), var(--accent-dark));
  color: #fff;
  border-color: transparent;
  box-shadow: 0 12px 26px -10px var(--accent-glow);
}
.auth-required__actions button.primary:hover {
  filter: brightness(1.08);
  transform: translateY(-1px);
}
.auth-required__actions button:not(.primary):hover {
  border-color: var(--accent);
  background: var(--accent-soft);
  color: var(--accent-dark);
}
.auth-required__actions svg {
  width: 17px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
}

@media (max-width: 1024px) {
  .control-grid {
    grid-template-columns: 1fr 1fr;
  }
  .control-grid .scan-field {
    grid-column: 1 / -1;
  }
  .header-context {
    display: none;
  }
}
@media (max-width: 768px) {
  .user-chip {
    max-width: 90px;
    padding: 6px 10px;
    font-size: 11px;
  }
  .inventory-workspace {
    padding: 16px 16px 20px;
  }
  .control-panel {
    padding: 16px;
    margin-bottom: 16px;
  }
  .control-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  .stats-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  .stat-copy strong {
    font-size: 22px;
  }
  .list-toolbar {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
    padding: 14px;
  }
  .list-actions {
    width: 100%;
    flex-wrap: wrap;
  }
  .search-box {
    width: 100%;
  }
  .status-tabs {
    padding: 10px 14px 0;
    overflow-x: auto;
  }
  .status-tabs button {
    padding: 8px 14px;
    font-size: 12px;
    white-space: nowrap;
  }
  .table-scroll {
    padding: 0;
  }
  .col-sticky--second {
    position: static;
    left: auto;
  }
  .auth-required {
    grid-template-columns: 1fr;
    text-align: center;
    justify-items: center;
    gap: 18px;
    padding: 26px;
  }
  .auth-required__copy span {
    max-width: 100%;
  }
  .auth-required__actions {
    flex-direction: column;
    width: 100%;
  }
  .auth-required__actions button {
    width: 100%;
    justify-content: center;
  }
}
@media (max-width: 480px) {
  .user-chip {
    display: none;
  }
}
</style>

<style>
.inv-select-popper.el-popper,
.inv-picker-popper.el-popper {
  border: 1px solid #bae6fd !important;
  border-radius: 14px !important;
  box-shadow: 0 20px 48px -16px rgba(30, 41, 90, 0.28) !important;
  overflow: hidden;
}
.inv-select-popper .el-popper__arrow::before,
.inv-picker-popper .el-popper__arrow::before {
  border-color: #bae6fd !important;
}

.inv-select-popper .el-select-dropdown__wrap {
  max-height: 264px;
}
.inv-select-popper .el-select-dropdown__item {
  height: 38px;
  padding: 0 14px;
  color: #0b1220;
  font-size: 13.5px;
  font-weight: 600;
  line-height: 38px;
}
.inv-select-popper .el-select-dropdown__item.hover,
.inv-select-popper .el-select-dropdown__item:hover {
  background: #f0f9ff;
  color: #0284c7;
}
.inv-select-popper .el-select-dropdown__item.is-selected {
  background: linear-gradient(135deg, #f0f9ff, #bae6fd);
  color: #0369a1;
  font-weight: 800;
}

.inv-picker-popper .el-picker-panel {
  border: 0;
  border-radius: 14px;
}
.inv-picker-popper
  .el-date-table
  td.current:not(.disabled)
  .el-date-table-cell__text,
.inv-picker-popper .el-month-table td.current:not(.disabled) .cell {
  background: linear-gradient(135deg, #0284c7, #0369a1);
}
.inv-picker-popper .el-date-table td.today .el-date-table-cell__text,
.inv-picker-popper .el-month-table td.today .cell {
  color: #0284c7;
  font-weight: 800;
}
</style>