<template>
  <div class="lts-page">
    <!-- HEADER -->
    <header class="lts-header">
      <div class="lts-header__brand">
        <img src="@/assets/Logo.png" alt="LAIYIH" class="lts-logo" />
        <div>
          <p class="lts-eyebrow">{{ t('lossTemSize.page.eyebrow') }}</p>
          <h1 class="lts-title">{{ t('lossTemSize.page.title') }}</h1>
        </div>
      </div>

      <div class="lts-header__right">
        <el-button class="lts-btn-lang" type="info" plain @click="changeLanguage">
          🌐 {{ currentLangLabel }}
        </el-button>

        <div class="lts-login-info-box">
          <div class="lts-user-info">
            <span class="lts-user-name">{{ loggedUser.username || loggedUser.userid }}</span>
            <span class="lts-user-id">{{ loggedUser.userid }}</span>
          </div>
        </div>

        <el-button class="lts-btn-logout" type="primary" @click="logout">
          {{ t('lossTemSize.auth.logout') }}
        </el-button>

        <el-button class="lts-btn-back" type="warning" @click="goBack" plain>
          <img src="@/assets/undo.png" alt="" />
        </el-button>
      </div>
    </header>

    <!-- SEARCH CARD -->
    <div class="lts-card lts-search-card">
      <div class="lts-card__head">
        <div class="lts-section-label" style="margin-bottom: 0">
          <el-icon><Document /></el-icon>
          {{ t('lossTemSize.page.searchInstruction') }}
        </div>
      </div>

      <div class="lts-search-grid">
        <el-input
          v-model="lenhInput"
          :placeholder="t('lossTemSize.search.placeholder')"
          clearable
          :class="['lts-input', { 'lts-input--error': lenhError }]"
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>

        <div class="lts-search-actions">
          <el-button
            class="lts-search-btn lts-excel-btn"
            :loading="excelParsing"
            :disabled="searching || excelParsing"
            @click="triggerExcelInput"
          >
            📊 Import Excel
          </el-button>
          <input
            ref="excelFileInput"
            type="file"
            accept=".xlsx,.xls"
            style="display:none"
            @change="handleExcelFileChange"
          />

          <el-button
            type="primary"
            class="lts-search-btn"
            :loading="searching && mode === 'IN'"
            :disabled="searching"
            @click="handleInClick"
          >
            {{ t('lossTemSize.search.importButton') }}
          </el-button>

          <el-button
            type="warning"
            class="lts-search-btn"
            :loading="searching && mode === 'OUT'"
            :disabled="searching"
            @click="handleOutClick"
          >
            {{ t('lossTemSize.search.takeButton') }}
          </el-button>
        </div>
      </div>

      <div class="lts-error-row">
        <p v-if="lenhError" class="lts-error">{{ lenhError }}</p>
      </div>
    </div>

    <!-- OVERVIEW STATS -->
    <section class="lts-overview">
      <div class="lts-table-header">
        <div class="lts-section-label" style="margin-bottom: 0">
          <el-icon><DataAnalysis /></el-icon>
          {{ t('lossTemSize.page.overviewTitle') }}
        </div>

        <div class="lts-table-header__tools">
          <span class="lts-filter-count-inline">
            {{ t('lossTemSize.summary.page') }}
            <strong>{{ summaryPage }}</strong>/<strong>{{ summaryTotalPages }}</strong>
          </span>

          <el-select
            v-model="summaryPageSize"
            size="small"
            style="width: 110px"
            @change="handleSummaryPageSizeChange"
          >
            <el-option :value="5"  label="5 / page" />
            <el-option :value="10" label="10 / page" />
            <el-option :value="20" label="20 / page" />
            <el-option :value="50" label="50 / page" />
          </el-select>

          <el-button
            v-if="filteredSummaryData.length > 0"
            @click="refreshSummary"
            :icon="Refresh"
            type="info"
            plain
            size="small"
            :loading="summaryLoading"
          >
            {{ t('lossTemSize.summary.refresh') }}
          </el-button>

          <el-button
            v-if="filteredSummaryData.length > 0"
            @click="exportToExcel"
            type="success"
            plain
            size="small"
            :loading="exportLoading"
          >
            📊 Export Excel
          </el-button>
        </div>
      </div>

      <div class="lts-overview-grid">
        <div class="lts-stat-box is-in">
          <span>{{ t('lossTemSize.summary.totalIn') }}</span>
          <strong>{{ grandInTotal }}</strong>
        </div>
        <div class="lts-stat-box is-out">
          <span>{{ t('lossTemSize.summary.totalOut') }}</span>
          <strong>{{ grandOutTotal }}</strong>
        </div>
        <div class="lts-stat-box is-net">
          <span>{{ t('lossTemSize.summary.remain') }}</span>
          <strong>{{ grandNetTotal }}</strong>
        </div>
      </div>
    </section>

    <!-- MAIN LAYOUT -->
    <div class="lts-body-layout" :class="{ 'is-collapsed': filterCollapsed }">
      <!-- SIDEBAR -->
      <aside class="lts-sidebar" :class="{ 'is-collapsed': filterCollapsed }">
        <div class="lts-sidebar-header">
          <div class="lts-section-label" style="margin-bottom: 0; font-size: 11px;">
            <el-icon><Filter /></el-icon>
            {{ t('lossTemSize.page.filterTitle') }}
          </div>
          <div class="lts-sidebar-actions">
            <button class="lts-toggle-btn" type="button" @click="filterCollapsed = !filterCollapsed">
              {{ filterCollapsed ? '›' : '‹' }}
            </button>
            <div class="lts-filter-badge">
              {{ filteredSummaryData.length }}<span>/{{ summaryData.length }}</span>
            </div>
          </div>
        </div>

        <div class="lts-sidebar-filters">
          <div class="lts-filter-group">
            <label class="lts-filter-label">{{ t('lossTemSize.filters.month') }}</label>
            <el-date-picker
              v-model="filterMonth"
              type="month"
              :placeholder="t('lossTemSize.filters.monthPlaceholder')"
              format="MM/YYYY"
              value-format="YYYY-MM"
              clearable
              class="lts-filter-item"
              style="width: 100%"
            />
          </div>

          <div class="lts-filter-group">
            <label class="lts-filter-label">{{ t('lossTemSize.filters.keyword') }}</label>
            <el-input
              v-model="filterKeyword"
              clearable
              :placeholder="t('lossTemSize.filters.keywordPlaceholder')"
              class="lts-filter-item"
              size="small"
            >
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
          </div>

          <div class="lts-filter-group">
            <label class="lts-filter-label">{{ t('lossTemSize.filters.reason') }}</label>
            <el-select
              v-model="filterReason"
              :placeholder="t('lossTemSize.filters.allReasons')"
              clearable
              class="lts-filter-item"
              size="small"
            >
              <el-option :label="t('lossTemSize.filters.allReasons')" value="ALL" />
              <el-option v-for="r in reasons" :key="r.key" :label="r.label" :value="r.key" />
            </el-select>
          </div>

          <div class="lts-filter-group">
            <label class="lts-filter-label">{{ t('lossTemSize.filters.status') }}</label>
            <el-select
              v-model="filterStatus"
              :placeholder="t('lossTemSize.filters.allStatus')"
              clearable
              class="lts-filter-item"
              size="small"
            >
              <el-option :label="t('lossTemSize.filters.allStatus')" value="ALL" />
              <el-option :label="t('lossTemSize.filters.hasOut')" value="HAS_OUT" />
              <el-option :label="t('lossTemSize.filters.hasRemain')" value="HAS_REMAIN" />
              <el-option :label="t('lossTemSize.filters.zeroNet')" value="ZERO_NET" />
            </el-select>
          </div>

          <div class="lts-filter-group">
            <label class="lts-filter-label">{{ t('lossTemSize.filters.sort') }}</label>
            <el-select
              v-model="filterSort"
              :placeholder="t('lossTemSize.filters.sort')"
              clearable
              class="lts-filter-item"
              size="small"
            >
              <el-option :label="t('lossTemSize.filters.dateDesc')" value="date_desc" />
              <el-option :label="t('lossTemSize.filters.dateAsc')" value="date_asc" />
              <el-option :label="t('lossTemSize.filters.ddbhAsc')" value="ddbh_asc" />
              <el-option :label="t('lossTemSize.filters.ddbhDesc')" value="ddbh_desc" />
              <el-option :label="t('lossTemSize.filters.netDesc')" value="net_desc" />
              <el-option :label="t('lossTemSize.filters.netAsc')" value="net_asc" />
            </el-select>
          </div>
        </div>

        <el-button @click="resetFilters" type="warning" plain size="small" class="lts-reset-btn">
          <el-icon><Refresh /></el-icon>
          {{ t('lossTemSize.filters.clear') }}
        </el-button>
      </aside>

      <!-- MAIN TABLE -->
      <main class="lts-main">
        <div v-if="summaryLoading" class="lts-state">
          <el-icon class="is-loading"><Loading /></el-icon>
          {{ t('lossTemSize.page.loading') }}
        </div>

        <div v-else-if="summaryError" class="lts-state lts-state--error">
          <el-icon><Warning /></el-icon>{{ summaryError }}
        </div>

        <div v-else-if="summaryData.length === 0" class="lts-state">
          {{ t('lossTemSize.page.noData') }}
        </div>

        <div v-else-if="filteredSummaryData.length === 0" class="lts-state lts-state--error">
          <el-icon><Warning /></el-icon>
          {{ t('lossTemSize.page.noResult') }}
        </div>

        <template v-else>
  <div class="lts-scroll-hint">{{ t('lossTemSize.page.scrollHint') }}</div>

  <div class="lts-table-wrap">
    <table class="lts-table">
      <thead>
        <tr>
          <th class="col-lenh col-resizable" :style="lenhColStyle">
            {{ t('lossTemSize.table.command') }}
            <span
              class="resize-handle"
              @mousedown.prevent.stop="startResize"
              @dblclick.stop="autoResizeLenh"
            />
          </th>
          <th class="col-pc">Product Code</th>
          <th class="col-qty">{{ t('lossTemSize.table.orderQty') }}</th>
          <th class="col-loai">{{ t('lossTemSize.table.typeReason') }}</th>
          <th v-for="sz in summaryAllSizes" :key="sz" class="col-size">{{ sz }}</th>
          <th class="col-total">{{ t('lossTemSize.table.total') }}</th>
        </tr>
      </thead>

      <tbody v-for="s in paginatedSummaryData" :key="s.ddbh" class="lts-tbody-group">
        <tr class="tr-in tr-group-first">
          <td class="td-lenh" :style="lenhColStyle" :rowspan="2 + reasons.length">
            <div class="lenh-badge">{{ s.ddbh }}</div>
            <div class="lenh-sub">{{ s.article }}</div>
            <div class="lenh-meta" v-if="s.xie_ming">{{ s.xie_ming }}</div>
            <div class="lenh-meta" v-if="s.ywpm">{{ s.ywpm }}</div>
            <div class="lenh-meta" v-if="s.date">{{ s.date }}</div>
          </td>

          <td class="td-pc" :rowspan="2 + reasons.length">
            <strong>{{ s.product_code || '—' }}</strong>
          </td>

          <td class="td-qty" :rowspan="2 + reasons.length">
            <strong>{{ s.pairs }}</strong>
          </td>

          <td class="td-loai">
            <el-tag type="success" size="small" effect="dark">LOSS TEM SIZE 標</el-tag>
          </td>

          <td v-for="sz in summaryAllSizes" :key="sz" class="td-center">
            <strong>{{ s.sizes_in?.[sz] ?? 0 }}</strong>
          </td>

          <td class="td-total">
            <strong>{{ s.total_in }}</strong>
          </td>
        </tr>

        <tr v-for="r in reasons" :key="`${s.ddbh}-${r.key}`" :class="getReasonRowClass(r.key)">
          <td class="td-loai">
            <el-tag :type="reasonTagType(r.key)" size="small" effect="plain">
              {{ reasonShort(r.key) }}
            </el-tag>
          </td>

          <td v-for="sz in summaryAllSizes" :key="sz" class="td-center">
            <strong>{{ s.out_by_reason?.[r.key]?.[sz] ?? '' }}</strong>
          </td>

          <td class="td-total">
            <strong>
              {{ s.out_by_reason?.[r.key] ? sumObjectValues(s.out_by_reason[r.key]) : '' }}
            </strong>
          </td>
        </tr>

        <tr class="tr-subtotal">
          <td colspan="1" class="td-subtotal-label">
            <strong>{{ t('lossTemSize.table.remain') }}</strong>
          </td>

          <td v-for="sz in summaryAllSizes" :key="sz" class="td-center">
            <strong>{{ (s.sizes_in?.[sz] ?? 0) - (s.sizes_out?.[sz] ?? 0) }}</strong>
          </td>

          <td class="td-total">
            <strong>{{ s.total_net }}</strong>
          </td>
        </tr>
      </tbody>

      <tbody class="lts-tbody-grand">
        <tr v-for="r in reasons" :key="'g-' + r.key" class="tr-grand-reason">
          <td colspan="4" class="td-subtotal-label">
            <strong>TỔNG {{ reasonShort(r.key) }} / 总计 {{ reasonShort(r.key) }}</strong>
          </td>

          <td v-for="sz in summaryAllSizes" :key="sz" class="td-center">
            <strong>{{ grandReasonTotals[r.key]?.[sz] || '' }}</strong>
          </td>

          <td class="td-total">
            <strong>{{ grandReasonTotal(r.key) }}</strong>
          </td>
        </tr>

        <tr class="tr-grand-net">
          <td colspan="4" class="td-subtotal-label">
            <strong>TỔNG CÒN LẠI / 剩余总计</strong>
          </td>

          <td v-for="sz in summaryAllSizes" :key="'gn-' + sz" class="td-center">
            <strong>{{ summaryGrandSizeData[sz]?.net ?? '' }}</strong>
          </td>

          <td class="td-total">
            <strong>{{ summaryGrandNetTotal }}</strong>
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <div v-if="filteredSummaryData.length > summaryPageSize" class="lts-pagination-row">
    <el-pagination
      background
      layout="prev, pager, next, sizes, jumper, total"
      :total="filteredSummaryData.length"
      :current-page="summaryPage"
      :page-size="summaryPageSize"
      :page-sizes="[5, 10, 20, 50]"
      @current-change="handleSummaryPageChange"
      @size-change="handleSummaryPageSizeChange"
    />
  </div>
</template>
      </main>
    </div>

    <!-- DIALOG THỦ CÔNG -->
    <el-dialog
      v-model="showEntryDialog"
      :title="mode === 'OUT' ? t('lossTemSize.dialog.takeTitle') : t('lossTemSize.dialog.importTitle')"
      :width="dialogWidth"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      class="lts-dialog"
      @closed="onDialogClosed"
    >
      <div v-if="currentLenh" class="lts-dialog-body">
        <div class="lts-info-bar">
          <div class="lts-info-grid">
            <div class="lts-info-item">
              <span>{{ t('lossTemSize.dialog.commandNo') }}</span>
              <strong>{{ currentLenh.commandNo }}</strong>
            </div>
            <div v-if="mode !== 'OUT'" class="lts-info-item">
              <span>{{ t('lossTemSize.dialog.msnv') }}</span>
              <strong>{{ loggedUser.userid }}</strong>
            </div>
            <div class="lts-info-item">
              <span>{{ t('lossTemSize.dialog.sku') }}</span>
              <strong>{{ currentLenh.sku }}</strong>
            </div>
            <div class="lts-info-item">
              <span>{{ t('lossTemSize.dialog.shoeType') }}</span>
              <strong>{{ currentLenh.shoeType || '—' }}</strong>
            </div>
            <div class="lts-info-item">
              <span>{{ t('lossTemSize.dialog.pairs') }}</span>
              <strong>{{ currentLenh.orderQty }}</strong>
            </div>
            <div class="lts-info-item">
              <span>{{ t('lossTemSize.dialog.totalRows') }}</span>
              <strong :class="mode === 'OUT' ? 'txt-danger' : 'txt-success'">{{ subTotal }}</strong>
            </div>
          </div>
          <div class="lts-mode-badge" :class="mode === 'OUT' ? 'is-out' : 'is-in'">
            {{ mode === 'OUT' ? t('lossTemSize.dialog.takeMode') : t('lossTemSize.dialog.importMode') }}
          </div>
        </div>

        <div v-if="mode === 'OUT'" class="lts-step-block lts-step-msnv">
          <div class="lts-step-label">
            <span class="lts-step-num">1</span>
            {{ t('lossTemSize.dialog.stepReason') }}
          </div>
          <div class="lts-msnv-box">
            <el-input
              v-model="msnvInput"
              :placeholder="t('lossTemSize.dialog.msnvPlaceholder')"
              clearable
              size="large"
              class="lts-msnv-input"
            >
              <template #prefix><el-icon><User /></el-icon></template>
            </el-input>
            <div class="msnv-help">{{ t('lossTemSize.dialog.msnvHelp') }}</div>
          </div>
        </div>

        <transition name="fade-slide">
          <div v-if="mode === 'OUT'" class="lts-step-block">
            <div class="lts-step-label">
              <span class="lts-step-num">2</span>
              {{ t('lossTemSize.dialog.stepReason') }}
            </div>
            <div class="lts-reason-list">
              <div
                v-for="r in reasons"
                :key="r.key"
                class="lts-reason-item"
                :class="[`reason--${r.color}`, { 'is-active': currentReason === r.key }]"
                @click="selectReason(r.key)"
              >
                <span class="reason-icon">{{ r.icon }}</span>
                <div class="reason-text">
                  <div class="reason-name">{{ r.label }}</div>
                  <div class="reason-sub">{{ r.sub }}</div>
                </div>
                <div v-if="currentReason === r.key" class="reason-check">✓</div>
              </div>
            </div>
            <div v-if="currentReason" class="lts-reason-note">
              <el-icon><Warning /></el-icon>
              {{ t('lossTemSize.dialog.reasonSelected') }}
              <strong>{{ reasonShort(currentReason) }}</strong>
              {{ t('lossTemSize.dialog.willBeSubtracted') }}
            </div>
          </div>
        </transition>
        <!-- PRODUCT CODE — bắt buộc khi IN -->
        <div v-if="mode === 'IN'" class="lts-step-block lts-step-msnv">
          <div class="lts-step-label">
            <span class="lts-step-num">1</span>
            产品编码 PRODUCT CODE <span style="color:#dc2626">*</span>
          </div>
          <div class="lts-msnv-box">
            <el-input
              v-model="productCodeInput"
              placeholder="Nhập mã sản phẩm — 请输入产品编码 (vd: PD240700960)"
              clearable
              size="large"
              class="lts-msnv-input"
            >
              <template #prefix><el-icon><Document /></el-icon></template>
            </el-input>
            <div class="msnv-help">Bắt buộc nhập </div>
          </div>
        </div>
        <div class="lts-step-block">
          <div class="lts-step-label">
            <span class="lts-step-num">{{ mode === 'OUT' ? 3 : 2 }}</span>
            {{ t('lossTemSize.dialog.stepSize') }}
          </div>
          <div class="lts-size-grid">
            <div
              v-for="sz in currentLenh.xxccSizes"
              :key="sz"
              class="lts-size-cell"
              :class="{
                'is-filled': (sizeInputs[sz] || 0) > 0,
                'is-out': mode === 'OUT',
              }"
              @click.capture="handleSizeCellClick"
            >
              <span class="size-lbl">{{ sz }}</span>
              <!-- IN & OUT: đều hiện số còn lại -->
              <div class="size-remain" :class="{ 'size-remain--in': mode === 'IN' }">
                {{ t('lossTemSize.dialog.remainLabel') }}:
                {{ lenhSizeRemainingFromSummary(currentLenh.commandNo, sz) }}
              </div>

              <el-input-number
                v-model="sizeInputs[sz]"
                :min="0"
                :max="9999"
                :step="0.5"       
                :precision="1"    
                :controls="false"
                :disabled="mode === 'OUT' && !currentReason"
                size="small"
                class="size-input"
                @change="(val) => {
                  nextTick(() => {
                    sizeInputs[sz] = val ?? null
                    touchedSizes.add(sz)
                    updateSubTotal()
                  })
                }"
              />
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="lts-dialog-footer">
          <span v-if="mode === 'OUT' && currentReason" class="dlg-note">
            {{ t('lossTemSize.dialog.noteReason') }}:
            <strong>{{ reasonShort(currentReason) }}</strong>
          </span>
          <div class="dlg-actions">
            <el-button @click="closeEntryDialog">{{ t('lossTemSize.dialog.cancel') }}</el-button>
            <el-button @click="clearSizes">{{ t('lossTemSize.dialog.clear') }}</el-button>
            <el-button
              type="success"
              @click="saveEntry"
              :disabled="
              !hasAnyInput ||
              (mode === 'IN'  && !productCodeInput.trim()) ||
              (mode === 'OUT' && !msnvInput.trim()) ||
              (mode === 'OUT' && !currentReason)
              "
            >
              <el-icon><Check /></el-icon>
              {{ t('lossTemSize.dialog.save') }}
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>

    <!-- EXCEL IMPORT DIALOG -->
      <el-dialog
        v-model="showExcelDialog"
        title="📊 IMPORT EXCEL — LOSS TEM SIZE"
        width="96%"
        :close-on-click-modal="false"
        :close-on-press-escape="!excelSaving"
        class="lts-dialog lts-excel-dialog"
        @closed="resetExcelState"
      >
      <div v-if="excelParsing" class="lts-state">
        <el-icon class="is-loading"><Loading /></el-icon> 
        Loading Excel...
      </div>

      <div v-else-if="excelRows.length === 0" class="lts-state lts-state--error">
        <el-icon><Warning /></el-icon>
        {{ t('lossTemSize.excel.noRows') }}
      </div>

      <div v-else class="lts-excel-body">
        <div class="lts-excel-toolbar">
          <div class="lts-excel-info">
            <el-tag type="info" size="small" effect="plain">
              📄 {{ excelRows.length }} rows found
            </el-tag>
            <el-tag :type="excelSelectedRows.length > 0 ? 'success' : 'warning'" size="small" effect="plain">
              ✅ {{ excelSelectedRows.length }} rows selected
            </el-tag>
            <el-tag v-if="excelErrorRows > 0" type="danger" size="small" effect="plain">
              ⚠ {{ excelErrorRows }} rows with errors
            </el-tag>
          </div>
        </div>

        <div v-if="excelSaving" class="lts-excel-progress">
          <div class="lts-excel-progress-label">
            <span>🔄 {{ t('lossTemSize.excel.saving') }}...</span>
            <strong>{{ excelProgressDone }}/{{ excelSelectedRows.length }}</strong>
          </div>
          <el-progress
            :percentage="excelProgressPct"
            :stroke-width="10"
            striped
            striped-flow
            :duration="4"
            status="success"
          />
        </div>

        <div class="lts-excel-table-wrap">
          <table class="lts-excel-table">
            <thead>
              <tr>
                <th class="th-chk">
                  <el-checkbox
                    :model-value="excelAllSelected"
                    :indeterminate="excelSomeSelected"
                    :disabled="excelSaving"
                    @change="toggleSelectAll"
                  />
                </th>
                <th class="th-no">#</th>
                <th class="th-ddbh">Lệnh (J)</th>
                <th class="th-mode">Mode</th>
                <th class="th-pc">Product Code</th>
                <th v-for="sz in excelAllSizes" :key="sz" class="th-sz">{{ sz }}</th>
                <th class="th-loss">Tổng LOSS</th>
                <th class="th-status">Trạng thái</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(row, idx) in excelRows"
                :key="idx"
                class="excel-row"
                :class="{
                  'is-selected': excelSelectedRows.includes(idx),
                  'is-error': !!row._error,
                  'is-saved': row._saved,
                  'is-fail': row._fail,
                }"
                @click="!excelSaving && toggleExcelRow(idx)"
              >
                <td class="td-chk" @click.stop>
                  <el-checkbox
                    :model-value="excelSelectedRows.includes(idx)"
                    :disabled="!!row._error || excelSaving"
                    @change="toggleExcelRow(idx)"
                  />
                </td>
                <td class="td-no">{{ idx + 1 }}</td>
                <td class="td-ddbh">{{ row.ddbh }}</td>
                <td class="td-mode">
                  <el-tag :type="row.mode === 'IN' ? 'success' : 'danger'" size="small" effect="dark">
                    {{ row.mode === 'IN' ? '⬇ NHẬP' : '⬆ LẤY' }}
                  </el-tag>
                </td>
                <td class="td-pc">{{ row.product_code || '—' }}</td>
                <td class="td-msnv">{{ row.msnv_out || '—' }}</td>
                <td
                  v-for="sz in excelAllSizes"
                  :key="sz"
                  class="td-sz"
                  :class="{ 'has-val': (row.sizes[sz] || 0) > 0 }"
                >
                  {{ row.sizes[sz] || '' }}
                </td>
                <td class="td-loss"><strong>{{ row.totalLoss }}</strong></td>
                <td class="td-status">
                  <el-tag v-if="row._saved" type="success" size="small" effect="dark">✓ Đã lưu</el-tag>
                  <el-tag v-else-if="row._fail" type="danger" size="small" effect="dark">✗ Lỗi</el-tag>
                  <el-tag v-else-if="row._error" type="danger" size="small" effect="plain">{{ row._error }}</el-tag>
                  <el-tag v-else type="success" size="small" effect="plain">OK</el-tag>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <template #footer>
        <div class="lts-dialog-footer">
          <div class="dlg-note">
            <span v-if="excelModeNote">{{ excelModeNote }}</span>
          </div>
          <div class="dlg-actions">
            <el-button @click="showExcelDialog = false" :disabled="excelSaving">
              {{ t('lossTemSize.dialog.cancel') }}
            </el-button>
            <el-button
              type="success"
              :loading="excelSaving"
              :disabled="excelRows.length === 0 || excelSelectedRows.length === 0"
              @click="saveExcelRows"
            >
              <el-icon><Check /></el-icon>
              SAVE IMPORT ({{ excelSelectedRows.length }})
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from '@/hooks/i18n.js'
import {
  Search, Document, Check, DataAnalysis, Warning,
  Refresh, Loading, Filter, User,
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as XLSX from 'xlsx-js-style'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()

const API_BASE_URL = (import.meta.env.VITE_API_URL || 'http://192.168.71.87:3084/api/v1').replace(/\/$/, '')
const SSE_URL = (import.meta.env.VITE_SSE_URL || `${API_BASE_URL}/sse`).replace(/\/$/, '')

const templateUrl = new URL(
  'templates/TEM_SIZE_LOSS_5_Nguyen_Nhan.xlsx',
  `${window.location.origin}${import.meta.env.BASE_URL || '/'}`,
).href

// ─── Language ──────────────────────────────────────────────────
const langOrder = ['vi', 'en', 'zh']
const currentLangLabel = computed(() => {
  if (locale.value === 'vi') return 'VI'
  if (locale.value === 'en') return 'EN'
  if (locale.value === 'zh') return '中文'
  return 'VI'
})
const changeLanguage = () => {
  const next = langOrder[(langOrder.indexOf(locale.value) + 1) % langOrder.length]
  locale.value = next
  localStorage.setItem('lang', next)
}

const REASON_KEYS = [
  'Thí nghiêm giày B',
  'Thí nghiêm hàng ngày',
  'Báo phế hoán đổi',
]

const norm = (v) => String(v ?? '').replace(/\s+/g, ' ').trim()
const cloneCell = (cell) => JSON.parse(JSON.stringify(cell))

// ─── Auth ──────────────────────────────────────────────────────
const loggedUser = reactive({
  userid: localStorage.getItem('userid') || '',
  username: localStorage.getItem('username') || '',
})

const logout = async () => {
  try {
    await ElMessageBox.confirm(
      t('lossTemSize.auth.confirmLogoutMessage'),
      t('lossTemSize.auth.confirmLogoutTitle'),
      {
        confirmButtonText: t('lossTemSize.auth.confirmButton'),
        cancelButtonText: t('lossTemSize.auth.cancelButton'),
        type: 'warning',
      },
    )
    localStorage.removeItem('token')
    localStorage.removeItem('userid')
    localStorage.removeItem('username')
    localStorage.removeItem('role')
    localStorage.removeItem('lang')
    stopSSE()
    router.push({ name: 'login' })
  } catch {
    /* cancelled */
  }
}

// ─── Constants ─────────────────────────────────────────────────
const reasons = computed(() => [
  {
    key: 'Thí nghiêm giày B',
    label: t('lossTemSize.reasons.trialShoeB.label'),
    sub: t('lossTemSize.reasons.trialShoeB.sub'),
    icon: '◈',
    color: 'danger',
  },
  {
    key: 'Thí nghiêm hàng ngày',
    label: t('lossTemSize.reasons.dailyTrial.label'),
    sub: t('lossTemSize.reasons.dailyTrial.sub'),
    icon: '◉',
    color: 'warning',
  },
  {
    key: 'Báo phế hoán đổi',
    label: t('lossTemSize.reasons.scrapExchange.label'),
    sub: t('lossTemSize.reasons.scrapExchange.sub'),
    icon: '◎',
    color: 'success',
  },
])

// ─── UI State ──────────────────────────────────────────────────
const filterCollapsed = ref(false)
const isMobile = ref(window.innerWidth < 768)
const dialogWidth = computed(() => (isMobile.value ? '96%' : '82%'))
const onResize = () => { isMobile.value = window.innerWidth < 768 }

const colWidths = reactive({ lenh: 140 })
const lenhColStyle = computed(() => ({
  width: `${colWidths.lenh}px`,
  minWidth: `${colWidths.lenh}px`,
  maxWidth: `${colWidths.lenh}px`,
}))

const startResize = (e) => {
  e.preventDefault()
  e.stopPropagation()
  const startX = e.clientX
  const startW = colWidths.lenh
  const onMove = (ev) => { colWidths.lenh = Math.max(110, startW + (ev.clientX - startX)) }
  const onUp = () => {
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
  }
  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
}

const autoResizeLenh = () => {
  let maxWidth = 110
  document.querySelectorAll('.col-lenh, .td-lenh').forEach(el => {
    const clone = el.cloneNode(true)
    Object.assign(clone.style, {
      width: 'auto',
      maxWidth: 'none',
      minWidth: '0',
      position: 'absolute',
      visibility: 'hidden',
      whiteSpace: 'nowrap',
    })
    document.body.appendChild(clone)
    maxWidth = Math.max(maxWidth, clone.scrollWidth + 40)
    document.body.removeChild(clone)
  })
  colWidths.lenh = maxWidth
}

// ─── Search State ──────────────────────────────────────────────
const lenhInput = ref('')
const lenhError = ref('')
const msnvInput = ref('')
const searching = ref(false)

// ─── Dialog State ──────────────────────────────────────────────
const currentLenh = ref(null)
const currentReason = ref(null)
const mode = ref('IN')
const showEntryDialog = ref(false)
const sizeInputs = reactive({})
const productCodeInput = ref('')
const subTotal = ref(0)

// ─── hasAnyInput: đặt SAU currentLenh và sizeInputs ───────────
// true = có ít nhất 1 ô đã được nhập (kể cả nhập 0)
// false = tất cả ô vẫn là null (chưa chạm vào)
const hasAnyInput = computed(() => {
  if (mode.value === 'OUT') {
    return (currentLenh.value?.xxccSizes || []).some(
      sz => sizeInputs[sz] !== null && sizeInputs[sz] !== undefined,
    )
  }
  // mode IN: pre-filled > 0 cũng tính là có input
  return (currentLenh.value?.xxccSizes || []).some(
    sz => sizeInputs[sz] !== null && sizeInputs[sz] !== undefined,
  )
})

// ─── Filter State ──────────────────────────────────────────────
const filterMonth = ref(null)
const filterKeyword = ref('')
const filterMsnv = ref('')
const filterReason = ref('ALL')
const filterStatus = ref('ALL')
const filterSort = ref('date_desc')

// ─── Summary State ─────────────────────────────────────────────
const summaryData = ref([])
const summaryLoading = ref(false)
const summaryError = ref('')
const summaryPage = ref(1)
const summaryPageSize = ref(20)

// ─── SSE State ─────────────────────────────────────────────────
const sseSource = ref(null)
const sseConnected = ref(false)
let sseRetryTimer = null
let sseRetryCount = 0
const MAX_RETRY_DELAY = 30000

// ─── Excel Import State ────────────────────────────────────────
const excelFileInput = ref(null)
const excelParsing = ref(false)
const excelSaving = ref(false)
const showExcelDialog = ref(false)
const excelRows = ref([])
const excelAllSizes = ref([])
const excelSelectedRows = ref([])
const excelDefaultMsnv = ref('')
const excelProgressDone = ref(0)
const excelProgressPct = ref(0)
const exportLoading = ref(false)

const excelErrorRows = computed(() => excelRows.value.filter(r => r._error).length)

const excelAllSelected = computed(() => {
  const valid = excelRows.value.filter(r => !r._error)
  return valid.length > 0 && excelSelectedRows.value.length === valid.length
})

const excelSomeSelected = computed(() =>
  excelSelectedRows.value.length > 0 && !excelAllSelected.value,
)

const excelModeNote = computed(() => {
  const inCount = excelRows.value.filter(r => r.mode === 'IN').length
  const outCount = excelRows.value.filter(r => r.mode === 'OUT').length
  return `📥 Nhập: ${inCount} lệnh | 📤 Xuất: ${outCount} lệnh`
})

// ─── Utilities ─────────────────────────────────────────────────
const toNum = (v) => Number(v ?? 0) || 0

// ---- NORMALIZE SIZE KEY (CHUẨN HÓA SIZE: "5" -> "05", "4.5" -> "04.5") ----
const normalizeSizeKey = (raw) => {
  if (raw === undefined || raw === null) return ''
  const num = parseFloat(String(raw).replace(/[^\d.-]/g, ''))
  if (isNaN(num)) return String(raw).trim()
  if (Number.isInteger(num)) {
    return num.toString().padStart(2, '0')
  }
  const intPart = Math.floor(num)
  const frac = num - intPart
  if (Math.abs(frac - 0.5) < 0.001) {
    return `${intPart.toString().padStart(2, '0')}.5`
  }
  return String(num)
}

const parseXXCC = (xxcc) => {
  if (!xxcc || !String(xxcc).trim()) return []
  const rawSizes = String(xxcc).replace(/\n/g, ' ').split(/[\s,;|]+/).map(s => s.trim()).filter(Boolean)
  const normalized = new Set()
  rawSizes.forEach(sz => {
    const norm = normalizeSizeKey(sz)
    if (norm) normalized.add(norm)
  })
  return [...normalized].sort((a, b) => parseFloat(a) - parseFloat(b))
}

const normalizeSizeObject = (obj) => {
  if (!obj) return {}
  const res = {}
  Object.entries(obj).forEach(([k, v]) => {
    const normKey = normalizeSizeKey(k)
    if (normKey) res[normKey] = toNum(v)
  })
  return res
}

const normalizeNestedReasonMap = (map) => {
  const res = {}
  Object.entries(map || {}).forEach(([r, m]) => {
    res[r] = normalizeSizeObject(m)
  })
  return res
}

const monthKeyFromDate = (d) => {
  if (!d) return ''
  const s = String(d).trim()
  return s.length >= 7 ? s.slice(0, 7) : ''
}

const reasonShort = (r) => {
  if (r === 'Thí nghiêm giày B') return t('lossTemSize.reasons.trialShoeB.short')
  if (r === 'Thí nghiêm hàng ngày') return t('lossTemSize.reasons.dailyTrial.short')
  if (r === 'Báo phế hoán đổi') return t('lossTemSize.reasons.scrapExchange.short')
  return r
}

const reasonTagType = (r) => {
  if (r === 'Thí nghiêm giày B') return 'danger'
  if (r === 'Thí nghiêm hàng ngày') return 'warning'
  if (r === 'Báo phế hoán đổi') return 'success'
  return 'info'
}

const getReasonRowClass = (r) => {
  if (r === 'Thí nghiêm giày B') return 'tr-loss'
  if (r === 'Thí nghiêm hàng ngày') return 'tr-bao'
  if (r === 'Báo phế hoán đổi') return 'tr-rem'
  return ''
}

const sumObjectValues = (obj) =>
  Object.values(obj || {}).reduce((a, v) => a + toNum(v), 0)

const rowText = (ws, rowIndex, maxCol = 33) => {
  const parts = []
  for (let c = 0; c <= maxCol; c++) {
    const addr = XLSX.utils.encode_cell({ r: rowIndex, c })
    const v = ws[addr]?.v
    if (v !== undefined && v !== null && String(v).trim() !== '') parts.push(String(v))
  }
  return norm(parts.join(' '))
}

// ─── Computed: filtered / paginated / totals ───────────────────
const filteredSummaryData = computed(() => {
  const kw = filterKeyword.value.trim().toLowerCase()
  const ms = filterMsnv.value.trim().toLowerCase()
  const reason = filterReason.value
  const status = filterStatus.value
  const sort = filterSort.value

  const base = summaryData.value.filter(item => {
    if (filterMonth.value && monthKeyFromDate(item.date) !== filterMonth.value) return false
    if (ms && !String(item.msnv || '').toLowerCase().includes(ms)) return false
    if (kw) {
      const hay = [item.ddbh, item.article, item.xie_ming, item.ywpm, item.msnv, item.date].join(' ').toLowerCase()
      if (!hay.includes(kw)) return false
    }
    if (reason !== 'ALL' && !Object.prototype.hasOwnProperty.call(item.out_by_reason || {}, reason)) return false
    if (status === 'HAS_OUT' && toNum(item.total_out) <= 0) return false
    if (status === 'HAS_REMAIN' && toNum(item.total_net) <= 0) return false
    if (status === 'ZERO_NET' && toNum(item.total_net) !== 0) return false
    return true
  })

  return [...base].sort((a, b) => {
    if (sort === 'date_asc') return String(a.date || '').localeCompare(String(b.date || ''))
    if (sort === 'date_desc') return String(b.date || '').localeCompare(String(a.date || ''))
    if (sort === 'ddbh_asc') return String(a.ddbh || '').localeCompare(String(b.ddbh || ''))
    if (sort === 'ddbh_desc') return String(b.ddbh || '').localeCompare(String(a.ddbh || ''))
    if (sort === 'net_asc') return toNum(a.total_net) - toNum(b.total_net)
    if (sort === 'net_desc') return toNum(b.total_net) - toNum(a.total_net)
    return 0
  })
})

const summaryTotalPages = computed(() =>
  Math.max(1, Math.ceil(filteredSummaryData.value.length / summaryPageSize.value)),
)

const paginatedSummaryData = computed(() => {
  const start = (summaryPage.value - 1) * summaryPageSize.value
  return filteredSummaryData.value.slice(start, start + summaryPageSize.value)
})

const summaryAllSizes = computed(() => {
  const set = new Set()
  summaryData.value.forEach(s => {
    ;(s.xxcc_sizes || []).forEach(sz => set.add(sz))
    Object.keys(s.sizes_in || {}).forEach(sz => set.add(sz))
    Object.keys(s.sizes_out || {}).forEach(sz => set.add(sz))
  })
  return [...set].sort((a, b) => parseFloat(a) - parseFloat(b))
})

const summaryGrandSizeData = computed(() => {
  const res = {}
  summaryAllSizes.value.forEach(sz => { res[sz] = { in: 0, out: 0, net: 0 } })
  filteredSummaryData.value.forEach(s => {
    summaryAllSizes.value.forEach(sz => {
      res[sz].in += toNum(s.sizes_in?.[sz])
      res[sz].out += toNum(s.sizes_out?.[sz])
    })
  })
  summaryAllSizes.value.forEach(sz => { res[sz].net = res[sz].in - res[sz].out })
  return res
})

const summaryGrandNetTotal = computed(() =>
  Object.values(summaryGrandSizeData.value).reduce((a, v) => a + v.net, 0),
)

const grandInTotal = computed(() => filteredSummaryData.value.reduce((a, s) => a + toNum(s.total_in), 0))
const grandOutTotal = computed(() => filteredSummaryData.value.reduce((a, s) => a + toNum(s.total_out), 0))
const grandNetTotal = computed(() => grandInTotal.value - grandOutTotal.value)

const grandReasonTotals = computed(() => {
  const totals = {}
  reasons.value.forEach(r => {
    totals[r.key] = {}
    summaryAllSizes.value.forEach(sz => { totals[r.key][sz] = 0 })
  })
  filteredSummaryData.value.forEach(item => {
    reasons.value.forEach(r => {
      Object.entries(item.out_by_reason?.[r.key] || {}).forEach(([sz, val]) => {
        totals[r.key][sz] = (totals[r.key][sz] || 0) + toNum(val)
      })
    })
  })
  return totals
})

const grandReasonTotal = (reasonKey) =>
  Object.values(grandReasonTotals.value[reasonKey] || {}).reduce((a, v) => a + v, 0)

// ── Lấy sizes đã có sẵn từ summaryData cho lệnh hiện tại ──
const existingSizesIn = computed(() => {
  if (!currentLenh.value) return {}
  const found = summaryData.value.find(s => s.ddbh === currentLenh.value.customerOrderNo)
  return found?.sizes_in || {}
})

// ── initSizeInputs: giữ null = chưa nhập, không pre-fill value ──
const initSizeInputs = (forceReset = false) => {
  if (forceReset) {
    Object.keys(sizeInputs).forEach(k => delete sizeInputs[k])
    touchedSizes.clear()
  }

  ;(currentLenh.value?.xxccSizes || []).forEach(sz => {
    if (!(sz in sizeInputs)) {
      sizeInputs[sz] = null  // ← luôn null, không pre-fill
    }
  })
  updateSubTotal()
}
const updateSubTotal = () => {
  subTotal.value = (currentLenh.value?.xxccSizes || []).reduce(
    (acc, sz) => acc + toNum(sizeInputs[sz]), 0,
  )
}

const clearSizes = () => initSizeInputs()

const lenhSizeRemainingFromSummary = (_zlbh, sz) => {
  const found = summaryData.value.find(s => s.ddbh === currentLenh.value?.customerOrderNo)
  if (!found) return 0
  return toNum(found.sizes_in?.[sz]) - toNum(found.sizes_out?.[sz])
}

// ─── API ───────────────────────────────────────────────────────
const fetchLenhFromApi = async (zlbh) => {
  const url = `${API_BASE_URL}/loss/list?zlbh=${encodeURIComponent(zlbh)}&page=1&page_size=1`
  const res = await fetch(url, { headers: { 'Content-Type': 'application/json' } })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  const json = await res.json()
  if (json.code !== 200 || !json.data?.items?.length) return null
  const item = json.data.items[0]
  return {
    commandNo: item.zlbh,
    customerOrderNo: item.ddbh,
    sku: item.article,
    shoeType: item.xie_ming,
    orderQty: item.pairs ?? 0,
    xxcc: item.xxcc || '',
    ywpm: item.ywpm || '',
    xxccSizes: [...new Set(parseXXCC(item.xxcc))],
    productCode: item.product_code || '',
  }
}

const saveEntryToAPI = async (payload) => {
  const res = await fetch(`${API_BASE_URL}/loss/loss-tem-size/save`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })

  const json = await res.json().catch(() => null)

  if (!res.ok) {
    throw new Error(json?.message || `Lỗi HTTP ${res.status}`)
  }

  if (json.code !== 200) throw new Error(json.message || t('lossTemSize.messages.saveError'))
  return json
}
const loadSummaryFromAPI = async () => {
  summaryLoading.value = true
  summaryError.value = ''
  try {
    const res = await fetch(`${API_BASE_URL}/loss/loss-tem-size/summary`, {
      headers: { 'Content-Type': 'application/json' },
    })

    if (!res.ok) {
      summaryError.value = `Lỗi HTTP ${res.status}`
      return
    }

    const json = await res.json()
    if (json.code !== 200) {
      summaryError.value = json.message || 'API lỗi'
      return
    }

    summaryData.value = (json.data || []).map(item => ({
      ddbh: item.ddbh ?? item.DDBH ?? '',
      msnv: item.msnv ?? item.MSNV ?? '',
      article: item.article ?? item.Article ?? '',
      xie_ming: item.xie_ming ?? item.XieMing ?? '',
      ywpm: item.ywpm ?? item.YWPM ?? '',
      product_code: item.product_code ?? item.ProductCode ?? item.PRODUCT_CODE ?? '',
      pairs: toNum(item.pairs ?? item.Pairs ?? 0),
      sizes_in: normalizeSizeObject(item.sizes_in ?? item.SizesIn ?? {}),
      sizes_out: normalizeSizeObject(item.sizes_out ?? item.SizesOut ?? {}),
      sizes_net: normalizeSizeObject(item.sizes_net ?? item.SizesNet ?? {}),
      out_by_reason: normalizeNestedReasonMap(item.out_by_reason ?? item.OutByReason ?? {}),
      total_in: toNum(item.total_in ?? item.TotalIn ?? 0),
      total_out: toNum(item.total_out ?? item.TotalOut ?? 0),
      total_net: toNum(item.total_net ?? item.TotalNet ?? 0),
      date: item.date ?? item.DATE ?? item.Date ?? '',
      xxcc_sizes: parseXXCC(item.xxcc ?? item.XXCC ?? ''),
      qty_out: toNum(item.qty_out ?? 0),
      userid: item.userid ?? '',
    }))

    if (summaryPage.value > summaryTotalPages.value) {
      summaryPage.value = summaryTotalPages.value
    }
  } catch (err) {
    summaryError.value = `Lỗi kết nối: ${err.message}`
  } finally {
    summaryLoading.value = false
  }
}

// ─── SSE ──────────────────────────────────────────────────────
const stopSSE = () => {
  if (sseRetryTimer) {
    clearTimeout(sseRetryTimer)
    sseRetryTimer = null
  }
  if (sseSource.value) {
    sseSource.value.close()
    sseSource.value = null
  }
  sseConnected.value = false
}

const scheduleReconnect = () => {
  if (document.hidden) return
  const delay = Math.min(1000 * Math.pow(2, sseRetryCount), MAX_RETRY_DELAY)
  sseRetryCount = Math.min(sseRetryCount + 1, 5)
  if (sseRetryTimer) clearTimeout(sseRetryTimer)
  sseRetryTimer = setTimeout(() => connectSSE(), delay)
}

const connectSSE = () => {
  stopSSE()
  const es = new EventSource(SSE_URL)
  sseSource.value = es
  es.onopen = () => { sseConnected.value = true; sseRetryCount = 0 }
  es.onmessage = async (event) => {
    if (!event.data?.trim()) return
    try {
      const data = JSON.parse(event.data)
      if (data.event === 'loss_tem_size_updated') {
        await loadSummaryFromAPI()
        ElMessage({
          message: `📡 ${data.ddbh} — ${data.mode === 'OUT' ? '⬆ Lấy SIZE' : '⬇ Nhập SIZE'} (MSNV: ${data.msnv})`,
          type: data.mode === 'OUT' ? 'warning' : 'success',
          duration: 3000,
        })
      }
    } catch (err) {
      console.warn('SSE parse error:', err)
    }
  }
  es.onerror = () => {
    sseConnected.value = false
    es.close()
    sseSource.value = null
    scheduleReconnect()
  }
}

const onVisibilityChange = () => {
  if (!document.hidden && !sseConnected.value) connectSSE()
}

// ─── Dialog helpers ────────────────────────────────────────────
const prepareEntryDialog = (entryMode, data) => {
  mode.value = entryMode
  currentLenh.value = data
  currentReason.value = null
  msnvInput.value = entryMode === 'OUT' ? '' : (loggedUser.userid || '')
  productCodeInput.value = data.productCode || ''
  initSizeInputs()
  showEntryDialog.value = true
}

const resetEntryState = () => {
  currentLenh.value = null
  currentReason.value = null
  msnvInput.value = ''
  productCodeInput.value = '' 
  subTotal.value = 0
  Object.keys(sizeInputs).forEach(k => delete sizeInputs[k])
}

const closeEntryDialog = () => { showEntryDialog.value = false }
const onDialogClosed = () => { resetEntryState() }

// ─── Actions (thủ công) ────────────────────────────────────────
const handleInClick = async () => {
  lenhError.value = ''
  const v = lenhInput.value.trim()
  if (!v) {
    lenhError.value = t('lossTemSize.search.emptyInput')
    return
  }
  if (!loggedUser.userid) {
    ElMessage.error(t('lossTemSize.auth.loginMissing'))
    router.push({ name: 'login' })
    return
  }
  searching.value = true
  try {
    const data = await fetchLenhFromApi(v)
    if (!data) {
      lenhError.value = `${t('lossTemSize.search.notFound')} "${v}"`
      return
    }
    prepareEntryDialog('IN', data)
    ElMessage.success(`${t('lossTemSize.search.foundSuccess')}: ${v} — ${data.xxccSizes.length} ${t('lossTemSize.search.sizes')}`)
  } catch (err) {
    lenhError.value = `${t('lossTemSize.search.apiError')}: ${err.message}`
  } finally {
    searching.value = false
  }
}

const handleOutClick = async () => {
  lenhError.value = ''
  const v = lenhInput.value.trim()
  if (!v) {
    lenhError.value = t('lossTemSize.search.emptyInput')
    return
  }
  searching.value = true
  try {
    const data = await fetchLenhFromApi(v)
    if (!data) {
      lenhError.value = `${t('lossTemSize.search.notFound')} "${v}"`
      return
    }
    prepareEntryDialog('OUT', data)
    ElMessage.success(`${t('lossTemSize.search.foundSuccess')}: ${v} — ${data.xxccSizes.length} ${t('lossTemSize.search.sizes')}`)
  } catch (err) {
    lenhError.value = `${t('lossTemSize.search.apiError')}: ${err.message}`
  } finally {
    searching.value = false
  }
}

const selectReason = (key) => {
  if (mode.value !== 'OUT') return
  currentReason.value = key
  initSizeInputs()
}

const saveEntry = async () => {
  if (!currentLenh.value) return
  const loginMsnv = loggedUser.userid.trim()
  const outMsnv = msnvInput.value.trim()

  if (!loginMsnv) {
    ElMessage.warning(t('lossTemSize.messages.loginMsnvMissing'))
    return
  }
  if (mode.value === 'IN' && !productCodeInput.value.trim()) {
  ElMessage.warning('Vui lòng nhập Mã sản phẩm (PRODUCT CODE) 产品编码')
  return
  }
  if (mode.value === 'OUT' && !outMsnv) {
    ElMessage.warning(t('lossTemSize.messages.takerMsnvMissing'))
    return
  }
  if (mode.value === 'OUT' && !currentReason.value) {
    ElMessage.warning(t('lossTemSize.messages.reasonMissing'))
    return
  }

  if (!hasAnyInput.value) {
    ElMessage.warning(t('lossTemSize.messages.quantityMissing'))
    return
  }

  if (mode.value === 'OUT') {
    for (const sz of currentLenh.value.xxccSizes) {
      const inputQty = toNum(sizeInputs[sz])
      const remaining = lenhSizeRemainingFromSummary(currentLenh.value.commandNo, sz)
      // ✅ Dùng epsilon 0.001 để tránh lỗi floating point
      if (inputQty - remaining > 0.001) {
        ElMessage.error(t('lossTemSize.messages.sizeExceeded', { size: sz, remain: remaining }))
        return
      }
    }
  }

  const sizesCopy = {}
  currentLenh.value.xxccSizes.forEach((sz) => {
    const rawVal = sizeInputs[sz]
    if (rawVal === null || rawVal === undefined) return
    // ✅ toFixed(4) làm sạch float trước khi gửi API
    sizesCopy[sz] = parseFloat(toNum(rawVal).toFixed(4))
  })

  if (Object.keys(sizesCopy).length === 0) {
    ElMessage.warning(t('lossTemSize.messages.quantityMissing'))
    return
  }

  const payload = {
    ddbh: currentLenh.value.customerOrderNo,
    msnv: loginMsnv,
    msnv_out: mode.value === 'OUT' ? outMsnv : '',
    ywpm: currentLenh.value.ywpm || '',
    article: currentLenh.value.sku,
    xie_ming: currentLenh.value.shoeType || '',
    pairs: currentLenh.value.orderQty,
    mode: mode.value,
    reason: mode.value === 'OUT' ? currentReason.value : '',
    note: '',
    sizes: sizesCopy,
    date: new Date().toISOString().split('T')[0],
    product_code: productCodeInput.value.trim(),
  }

  try {
    await saveEntryToAPI(payload)
    ElMessage.success(
      mode.value === 'OUT'
        ? t('lossTemSize.messages.saveSuccessTake', { commandNo: currentLenh.value.commandNo })
        : t('lossTemSize.messages.saveSuccessImport', { commandNo: currentLenh.value.commandNo }),
    )
    showEntryDialog.value = false
    resetEntryState()
    lenhInput.value = ''
    lenhError.value = ''
    await loadSummaryFromAPI()
  } catch (err) {
    ElMessageBox.alert(err.message, '⚠️ Không thể lưu', {
      confirmButtonText: 'Đã hiểu',
      type: 'error',
    })
  }
}
const refreshSummary = () => loadSummaryFromAPI()

const goBack = () => {
  const from = route.query.from
  if (from === 'warehouse-kvt') {
    router.push('/warehouse-kvt')
    return
  }
  router.push('/')
}

const resetFilters = () => {
  const now = new Date()
  const mm = String(now.getMonth() + 1).padStart(2, '0')
  filterMonth.value = `${now.getFullYear()}-${mm}`  // ← tháng hiện tại thay vì null
  filterKeyword.value = ''
  filterMsnv.value = ''
  filterReason.value = 'ALL'
  filterStatus.value = 'ALL'
  filterSort.value = 'date_desc'
  summaryPage.value = 1
}

// ─── EXPORT EXCEL ──────────────────────────────────────────────
const TEMPLATE_SIZE_FLOATS = [
  3.5, 4, 4.5, 5, 5.5, 6, 6.5, 7, 7.5, 8,
  8.5, 9, 9.5, 10, 10.5, 11, 11.5, 12, 12.5, 13, 14, 15, 16,
]
// Col index tương ứng: size float -> col index (10..32)
const SIZE_COL_START = 10   // cột K
const TOTAL_COL     = 33   // cột AH
 
// Chuyển normKey ("05", "05.5", "10"...) sang float để tra cứu col index
const normKeyToFloat = (key) => parseFloat(key)
 
// Tra col index cho 1 normKey (vd "05" -> col 13)
const sizeKeyToColIndex = (normKey) => {
  const f = normKeyToFloat(normKey)
  const idx = TEMPLATE_SIZE_FLOATS.indexOf(f)
  if (idx === -1) return -1
  return SIZE_COL_START + idx
}
 
// Nhãn cột I cho từng loại dòng
const ROW_LABELS = {
  IN:          'LOSS_TEM SIZE 標',
  TRIAL_B:     'Thí nghiêm giày B',
  DAILY:       'Thí nghiêm hàng ngày',
  SCRAP:       'Báo phế hoán đổi',
  REMAIN:      'Số lượng còn lại',
  // Grand total rows
  GRAND_IN:    'TOTAL LOSS_TEM SIZE 標',
  GRAND_B:     'TOTAL Thí nghiêm giày B',
  GRAND_DAILY: 'TOTAL Thí nghiêm hàng ngày',
  GRAND_SCRAP: 'TOTAL Báo phế hoán đổi',
  GRAND_NET:   'TOTAL Số lượng còn lại',
}
 
const REASON_KEY_MAP = {
  TRIAL_B: 'Thí nghiêm giày B',
  DAILY:   'Thí nghiêm hàng ngày',
  SCRAP:   'Báo phế hoán đổi',
}
 
// ── Helper: set cell value vào wsNew ──────────────────────────
const setCell = (wsNew, rowIndex, colIndex, value, type = null) => {
  if (value === null || value === undefined || value === '') return
  const addr = XLSX.utils.encode_cell({ r: rowIndex, c: colIndex })
  const t = type || (typeof value === 'number' ? 'n' : 's')
  wsNew[addr] = { v: value, t }
}
 
// ── Helper: copy toàn bộ style từ ô mẫu (nếu có) ─────────────
const copyStyleCell = (sourceWs, sourceRow, sourceCol, wsNew, destRow, destCol, overrideValue) => {
  const srcAddr = XLSX.utils.encode_cell({ r: sourceRow, c: sourceCol })
  const destAddr = XLSX.utils.encode_cell({ r: destRow, c: destCol })
  const src = sourceWs[srcAddr]
  if (src) {
    const newCell = JSON.parse(JSON.stringify(src))
    if (overrideValue !== undefined && overrideValue !== null && overrideValue !== '') {
      newCell.v = overrideValue
      newCell.t = typeof overrideValue === 'number' ? 'n' : 's'
    } else {
      // Giữ nguyên cell nhưng xóa value (giữ style)
      delete newCell.v
      delete newCell.t
      delete newCell.f
    }
    wsNew[destAddr] = newCell
  } else if (overrideValue !== undefined && overrideValue !== null && overrideValue !== '') {
    const t = typeof overrideValue === 'number' ? 'n' : 's'
    wsNew[destAddr] = { v: overrideValue, t }
  }
}
 
// ── Main export function ───────────────────────────────────────

const exportToExcel = async () => {
  if (!filteredSummaryData.value.length) {
    ElMessage.warning('Không có dữ liệu để xuất')
    return
  }
  exportLoading.value = true

  try {
    const wb      = XLSX.utils.book_new()
    const wsNew   = {}
    const merges  = []

    const monthLabel = filterMonth.value
      ? filterMonth.value.replace('-', '/')
      : new Date().toISOString().slice(0, 7).replace('-', '/')

    // ── Lấy sizes tự động từ data thực tế ──────────────────────────
    const sortedSizes = [...summaryAllSizes.value].sort(
      (a, b) => parseFloat(a) - parseFloat(b)
    )

    const SIZE_COL_START = 8
    const TOTAL_COL      = SIZE_COL_START + sortedSizes.length  // động theo data

    // ── Bảng màu ────────────────────────────────────────────────────
    const C = {
      NAVY:      '1F3864',
      BLUE:      '2E75B6',
      WHITE:     'FFFFFF',
      BLACK:     '000000',
      GRAY_BORD: '595959',
      LOSS:      'FFFFFF',
      TRIAL_B:   'DEEAF1',
      DAILY:     'FFFFFF',
      SCRAP:     'FCE4D6',
      REMAIN:    'E2EFDA',
      REMAIN_FT: '1F5C1F',
      GRAND:     ['F4B942', 'BDD7EE', 'DDEBF7', 'F8CBAD', 'A9D18E'],
    }

    const side     = (style, color) => ({ style, color: { rgb: color } })
    const THIN     = (c = C.GRAY_BORD) => side('thin',   c)
    const MEDIUM   = (c = C.NAVY)      => side('medium', c)
    const WHITE_MED = side('medium', C.WHITE)

    const S = ({
      bg, fontColor = C.BLACK, bold = false, fontSize = 11,
      hAlign = 'center', vAlign = 'center', wrap = true,
      bTop, bBottom, bLeft, bRight,
    }) => ({
      fill:      { patternType: 'solid', fgColor: { rgb: bg } },
      font:      { name: 'Times New Roman', sz: fontSize, bold, color: { rgb: fontColor } },
      alignment: { horizontal: hAlign, vertical: vAlign, wrapText: wrap },
      border: {
        top:    bTop    || THIN(),
        bottom: bBottom || THIN(),
        left:   bLeft   || THIN(),
        right:  bRight  || THIN(),
      },
    })

    const S_TITLE = S({
      bg: C.NAVY, fontColor: C.WHITE, bold: true, fontSize: 14,
      bTop: WHITE_MED, bBottom: WHITE_MED, bLeft: WHITE_MED, bRight: WHITE_MED,
    })
    const S_H1 = S({
      bg: C.NAVY, fontColor: C.WHITE, bold: true, fontSize: 12,
      bTop: WHITE_MED, bBottom: WHITE_MED, bLeft: WHITE_MED, bRight: WHITE_MED,
    })
    const S_H2 = S({
      bg: C.BLUE, fontColor: C.WHITE, bold: true, fontSize: 11,
      bTop: WHITE_MED, bBottom: WHITE_MED, bLeft: WHITE_MED, bRight: WHITE_MED,
    })

    const rowStyle = (bg, fontColor, bold, isFirst, isLast, isLeftEdge, isRightEdge, hAlign = 'center') =>
      S({
        bg, fontColor, bold, fontSize: 11, hAlign, wrap: true,
        bTop:    isFirst     ? MEDIUM() : THIN(),
        bBottom: isLast      ? MEDIUM() : THIN(),
        bLeft:   isLeftEdge  ? MEDIUM() : THIN(),
        bRight:  isRightEdge ? MEDIUM() : THIN(),
      })

    const OFFSET_CFG = [
      { bg: C.LOSS,    fc: C.BLACK,      bold: true  },
      { bg: C.TRIAL_B, fc: C.BLACK,      bold: false },
      { bg: C.DAILY,   fc: C.BLACK,      bold: false },
      { bg: C.SCRAP,   fc: C.BLACK,      bold: false },
      { bg: C.REMAIN,  fc: C.REMAIN_FT,  bold: true  },
    ]

    // ── Helper ghi ô ────────────────────────────────────────────────
    const wc = (r, c, value, style) => {
      const addr = XLSX.utils.encode_cell({ r, c })
      const v    = value !== null && value !== undefined ? value : 0
      wsNew[addr] = { v, t: typeof v === 'number' ? 'n' : 's', s: style }
    }

    // ── ROW 0: Tiêu đề tháng ────────────────────────────────────────
    wc(0, 0, `THÁNG ${monthLabel}`, S_TITLE)
    for (let c = 1; c <= TOTAL_COL; c++) wc(0, c, '', S_TITLE)
    merges.push({ s: { r: 0, c: 0 }, e: { r: 0, c: TOTAL_COL } })

    // ── ROW 1: Header chính ─────────────────────────────────────────
    const H1_LABELS = {
      0: '产品名称(客户) CHI TIẾT',
      1: '产品编码 PRODUCT CODE',         // ← THÊM col 1
      2: '指令号 LỆNH',                   // ← dịch từ 1
      3: '型体配色 SKU',                  // ← dịch từ 2
      4: '鞋型 Dạng giày',               // ← dịch từ 3
      5: '订单量(客户) SL ĐƠN HÀNG',     // ← dịch từ 4
      6: 'GHI CHÚ 備註',                  // ← dịch từ 5
      7: 'SIZE',                          // ← dịch từ 6
    }
    for (let c = 0; c <= TOTAL_COL; c++) {
      let label = H1_LABELS[c] ?? ''
      if (c >= SIZE_COL_START && c < TOTAL_COL) {
        // Hiện size thực tế từ data
        label = parseFloat(sortedSizes[c - SIZE_COL_START])
      }
      if (c === TOTAL_COL) label = 'Loss数(客户) SL LOSS'
      wc(1, c, label, S_H1)
    }

    // ── ROW 2: Sub-header ───────────────────────────────────────────
    const H2_LABELS = {
      0: 'YWPM',
      1: 'PRODUCT_CODE',    // ← THÊM
      2: 'ZLBH/DDBH',      // ← dịch từ 1
      3: 'article',         // ← dịch từ 2
      5: 'pair',            // ← dịch từ 4
    }
    for (let c = 0; c <= TOTAL_COL; c++) {
      wc(2, c, H2_LABELS[c] ?? '', S_H2)
    }

    // ── ROW data ────────────────────────────────────────────────────
    let currentRow = 3

    for (const item of filteredSummaryData.value) {
      const rowStart = currentRow

      const rowTypes = [
        {
          label:    ROW_LABELS.IN,
          getSizes: (sz) => toNum(item.sizes_in?.[sz]),
          getTotal: ()   => toNum(item.total_in),
        },
        {
          label:    ROW_LABELS.TRIAL_B,
          getSizes: (sz) => toNum(item.out_by_reason?.['Thí nghiêm giày B']?.[sz]),
          getTotal: ()   => sumObjectValues(item.out_by_reason?.['Thí nghiêm giày B']),
        },
        {
          label:    ROW_LABELS.DAILY,
          getSizes: (sz) => toNum(item.out_by_reason?.['Thí nghiêm hàng ngày']?.[sz]),
          getTotal: ()   => sumObjectValues(item.out_by_reason?.['Thí nghiêm hàng ngày']),
        },
        {
          label:    ROW_LABELS.SCRAP,
          getSizes: (sz) => toNum(item.out_by_reason?.['Báo phế hoán đổi']?.[sz]),
          getTotal: ()   => sumObjectValues(item.out_by_reason?.['Báo phế hoán đổi']),
        },
        {
          label:    ROW_LABELS.REMAIN,
          getSizes: (sz) => toNum(item.sizes_in?.[sz]) - toNum(item.sizes_out?.[sz]),
          getTotal: ()   => toNum(item.total_net),
        },
      ]

      for (let offset = 0; offset < 5; offset++) {
        const r       = rowStart + offset
        const rt      = rowTypes[offset]
        const cfg     = OFFSET_CFG[offset]
        const isFirst = offset === 0
        const isLast  = offset === 4

      if (offset === 0) {
        wc(r, 0, item.ywpm,              rowStyle(cfg.bg, cfg.fc, cfg.bold, true, false, true,  false, 'left'))
        wc(r, 1, item.product_code || '',rowStyle(cfg.bg, cfg.fc, cfg.bold, true, false, false, false, 'center')) // ← THÊM
        wc(r, 2, item.ddbh,              rowStyle(cfg.bg, cfg.fc, cfg.bold, true, false, false, false, 'center'))
        wc(r, 3, item.article,           rowStyle(cfg.bg, cfg.fc, cfg.bold, true, false, false, false, 'center'))
        wc(r, 4, item.xie_ming || '',    rowStyle(cfg.bg, cfg.fc, cfg.bold, true, false, false, false, 'center'))
        wc(r, 5, item.pairs,             rowStyle(cfg.bg, cfg.fc, true,     true, false, false, false, 'center'))
        wc(r, 7, 'SIZE',                 rowStyle(cfg.bg, cfg.fc, cfg.bold, true, false, false, false, 'center'))
      }
      wc(r, 6, rt.label, rowStyle(cfg.bg, cfg.fc, cfg.bold, isFirst, isLast, false, false, 'left'))

        // ── Size columns: tự động theo sortedSizes ──
        sortedSizes.forEach((sz, sizeIdx) => {
          const colIndex = SIZE_COL_START + sizeIdx
          wc(r, colIndex, rt.getSizes(sz),
            rowStyle(cfg.bg, cfg.fc, cfg.bold, isFirst, isLast, false, false, 'center'))
        })

        wc(r, TOTAL_COL, rt.getTotal(),
          rowStyle(cfg.bg, cfg.fc, true, isFirst, isLast, false, true, 'center'))
      }

      for (const c of [0, 1, 2, 3, 4, 5, 7]) {    // ← đổi từ [0, 1, 2, 3, 4, 6]
        merges.push({ s: { r: rowStart, c }, e: { r: rowStart + 4, c } })
      }


      currentRow += 5
    }

    // ── Grand total ─────────────────────────────────────────────────
    const grandRows = [
      {
        label:    ROW_LABELS.GRAND_IN,
        getSizes: (sz) => filteredSummaryData.value.reduce((a, s) => a + toNum(s.sizes_in?.[sz]), 0),
        getTotal: ()   => filteredSummaryData.value.reduce((a, s) => a + toNum(s.total_in), 0),
      },
      {
        label:    ROW_LABELS.GRAND_B,
        getSizes: (sz) => filteredSummaryData.value.reduce(
          (a, s) => a + toNum(s.out_by_reason?.['Thí nghiêm giày B']?.[sz]), 0),
        getTotal: ()   => filteredSummaryData.value.reduce(
          (a, s) => a + sumObjectValues(s.out_by_reason?.['Thí nghiêm giày B']), 0),
      },
      {
        label:    ROW_LABELS.GRAND_DAILY,
        getSizes: (sz) => filteredSummaryData.value.reduce(
          (a, s) => a + toNum(s.out_by_reason?.['Thí nghiêm hàng ngày']?.[sz]), 0),
        getTotal: ()   => filteredSummaryData.value.reduce(
          (a, s) => a + sumObjectValues(s.out_by_reason?.['Thí nghiêm hàng ngày']), 0),
      },
      {
        label:    ROW_LABELS.GRAND_SCRAP,
        getSizes: (sz) => filteredSummaryData.value.reduce(
          (a, s) => a + toNum(s.out_by_reason?.['Báo phế hoán đổi']?.[sz]), 0),
        getTotal: ()   => filteredSummaryData.value.reduce(
          (a, s) => a + sumObjectValues(s.out_by_reason?.['Báo phế hoán đổi']), 0),
      },
      {
        label:    ROW_LABELS.GRAND_NET,
        getSizes: (sz) => filteredSummaryData.value.reduce(
          (a, s) => a + toNum(s.sizes_in?.[sz]) - toNum(s.sizes_out?.[sz]), 0),
        getTotal: ()   => filteredSummaryData.value.reduce((a, s) => a + toNum(s.total_net), 0),
      },
    ]

    for (let gi = 0; gi < grandRows.length; gi++) {
      const r       = currentRow + gi
      const gr      = grandRows[gi]
      const bg      = C.GRAND[gi]
      const isFirst = gi === 0
      const isLast  = gi === grandRows.length - 1

      for (const c of [0, 1, 2, 3, 4, 5, 7]) {    // ← đổi từ [0, 1, 2, 3, 4, 6]
        wc(r, c, '', rowStyle(bg, C.BLACK, true, isFirst, isLast, c === 0, false, 'center'))
      }

      wc(r, 6, gr.label, rowStyle(bg, C.BLACK, true, isFirst, isLast, false, false, 'left'))

      // ── Size columns grand total: tự động theo sortedSizes ──
      sortedSizes.forEach((sz, sizeIdx) => {
        const colIndex = SIZE_COL_START + sizeIdx
        wc(r, colIndex, gr.getSizes(sz),
          rowStyle(bg, C.BLACK, true, isFirst, isLast, false, false, 'center'))
      })

      wc(r, TOTAL_COL, gr.getTotal(),
        rowStyle(bg, C.BLACK, true, isFirst, isLast, false, true, 'center'))
    }

    const lastRow = currentRow + grandRows.length - 1

    // ── Column widths: tự động theo số size ─────────────────────────
    const colWidths = []
    colWidths[0] = { wch: 42 }   // ywpm
    colWidths[1] = { wch: 16 }   // product_code ← THÊM
    colWidths[2] = { wch: 14 }   // ddbh
    colWidths[3] = { wch: 16 }   // article
    colWidths[4] = { wch: 18 }   // xie_ming
    colWidths[5] = { wch: 8  }   // pairs
    colWidths[6] = { wch: 28 }   // label
    colWidths[7] = { wch: 6  }   // SIZE header
    sortedSizes.forEach((_, i) => { colWidths[SIZE_COL_START + i] = { wch: 5 } })
    colWidths[TOTAL_COL] = { wch: 10 }

    // ── Row heights ─────────────────────────────────────────────────
    const rowHeights = []
    rowHeights[0] = { hpt: 28 }
    rowHeights[1] = { hpt: 36 }
    rowHeights[2] = { hpt: 22 }
    for (let r = 3; r <= lastRow; r++) rowHeights[r] = { hpt: 20 }

    // ── Metadata ────────────────────────────────────────────────────
    wsNew['!ref']    = XLSX.utils.encode_range({ s: { r: 0, c: 0 }, e: { r: lastRow, c: TOTAL_COL } })
    wsNew['!merges'] = merges
    wsNew['!cols']   = colWidths
    wsNew['!rows']   = rowHeights

    // ── Download ────────────────────────────────────────────────────
    const newSheetName = `THANG_${monthLabel.replace('/', '_')}`
    XLSX.utils.book_append_sheet(wb, wsNew, newSheetName)

    const fileName = `TEM_SIZE_LOSS_${monthLabel.replace('/', '_')}_${Date.now()}.xlsx`
    XLSX.writeFile(wb, fileName, { cellStyles: true, bookSST: true })

    ElMessage.success(`✅ Đã xuất ${filteredSummaryData.value.length} lệnh → ${fileName}`)
  } catch (err) {
    ElMessage.error(`❌ Lỗi xuất Excel: ${err.message}`)
    console.error(err)
  } finally {
    exportLoading.value = false
  }
}
const handleSizeCellClick = () => {
  if (mode.value === 'OUT' && !currentReason.value) {
    ElMessage.warning(t('lossTemSize.messages.chooseReasonFirst'))
  }
}

const handleSummaryPageChange = (page) => { summaryPage.value = page }
const handleSummaryPageSizeChange = (size) => { summaryPageSize.value = size; summaryPage.value = 1 }

watch(
  [() => filteredSummaryData.value.length, summaryPageSize],
  () => {
    const total = summaryTotalPages.value
    if (summaryPage.value > total) summaryPage.value = total
    if (summaryPage.value < 1) summaryPage.value = 1
  },
)

// ══════════════════════════════════════════════════════════════
// EXCEL IMPORT — ĐỌC FILE CHI_TIET (SHEET: TEM)
function extractSizeNumbers(raw) {
  if (raw == null) return []
  const s = String(raw).trim()
  if (!s) return []
  // Lấy phần đầu tiên nếu có dấu / hoặc - hoặc ~
  const first = s.split(/[\/\-~]/)[0].trim()
  // Đổi dấu phẩy thành dấu chấm (kiểu EU: "11,5" => "11.5")
  const normalized = first.replace(',', '.')
  // Bỏ tất cả ký tự không phải số và dấu chấm
  const num = normalized.replace(/[^0-9.]/g, '').trim()
  // Nếu có nhiều hơn 1 dấu chấm thì chỉ giữ phần trước dấu chấm thứ 2
  const clean = num.replace(/^(\d+\.?\d*).*$/, '$1')
  // Kiểm tra có phải số hợp lệ không
  return clean && !isNaN(parseFloat(clean)) ? [clean] : []
}
const sortSizes = (sizes) =>
  [...sizes].sort((a, b) => parseFloat(a) - parseFloat(b))

const triggerExcelInput = () => {
  if (excelFileInput.value) excelFileInput.value.value = ''
  excelFileInput.value?.click()
}

const resetExcelState = () => {
  excelRows.value = []
  excelAllSizes.value = []
  excelSelectedRows.value = []
  excelDefaultMsnv.value = ''
  excelProgressDone.value = 0
  excelProgressPct.value = 0
}

const handleExcelFileChange = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return

  excelParsing.value = true
  showExcelDialog.value = true
  excelRows.value = []
  excelAllSizes.value = []

  try {
    const buffer = await file.arrayBuffer()
    const wb = XLSX.read(buffer, { type: 'array' })

    // ──────────────────────────────────────────────
    // Tìm sheet hợp lệ (ưu tiên tên có "TEM")
    // ──────────────────────────────────────────────
    const findSheet = (mustHaveTem) => {
      for (const name of wb.SheetNames) {
        if (mustHaveTem && !name.toUpperCase().includes('TEM')) continue

        const sheet = wb.Sheets[name]
        const rows = XLSX.utils.sheet_to_json(sheet, { header: 1, defval: null })
        if (rows.length < 2) continue

        const header = rows[0] || []
        // Col I (index 8) = 指令号, Col L (index 11) = 尺码
        const hasLenhCol = String(header[8] ?? '').includes('指令')
        const hasSizeCol = String(header[11] ?? '').includes('尺码')
        if (!hasLenhCol || !hasSizeCol) continue

        return { sheet, name }
      }
      return null
    }

    const found = findSheet(true) || findSheet(false)
    if (!found) throw new Error('Không tìm thấy sheet hợp lệ trong file Excel')

    const { sheet: ws, name: wsName } = found
    const raw = XLSX.utils.sheet_to_json(ws, { header: 1, defval: null })

    ElMessage.info(`📋 Đọc sheet: "${wsName}"`)

    // Cột (0-indexed): G=6, I=8, K=10, L=11
    const COL_PRODUCT_CODE = 6
    const COL_LENH        = 8
    const COL_LOSS        = 10
    const COL_SIZE        = 11

    const groupMap = new Map()
    let skippedRows = 0

    for (let i = 1; i < raw.length; i++) {
      const row = raw[i]
      if (!row) continue

      const lenh        = String(row[COL_LENH] ?? '').trim()
      const lossQty     = toNum(row[COL_LOSS])
      const rawSize     = row[COL_SIZE]
      const productCode = String(row[COL_PRODUCT_CODE] ?? '').trim()

      if (!lenh) { skippedRows++; continue }

      if (!groupMap.has(lenh)) {
        groupMap.set(lenh, { sizes: {}, temCount: 0, productCode: '' })
      }
      const entry = groupMap.get(lenh)
      if (!entry.productCode && productCode) entry.productCode = productCode

      const parsedSizes = extractSizeNumbers(rawSize)
      for (const sz of parsedSizes) {
        entry.sizes[sz] = (entry.sizes[sz] || 0) + lossQty
      }
      entry.temCount += lossQty
    }

    if (groupMap.size === 0) {
      throw new Error(
        `Không tìm thấy lệnh hợp lệ nào trong sheet (${skippedRows} dòng bỏ qua vì thiếu lệnh)`
      )
    }

    const allSizeSet = new Set()
    for (const [, entry] of groupMap)
      Object.keys(entry.sizes).forEach(sz => allSizeSet.add(sz))
    excelAllSizes.value = sortSizes([...allSizeSet])

    const parsed = []
    for (const [lenh, entry] of groupMap) {
      const totalLoss = Object.values(entry.sizes).reduce((a, v) => a + v, 0)
      let _error = null
      if (!lenh) _error = 'Thiếu Lệnh'
      else if (totalLoss === 0 && Object.keys(entry.sizes).length === 0)
        _error = 'Loss = 0, không có size'

      parsed.push({
        ddbh: lenh,
        zlbh: lenh,
        sizes: entry.sizes,
        temCount: entry.temCount,
        totalLoss,
        mode: 'IN',
        reason: '',
        msnv_out: '',
        product_code: entry.productCode,
        _error,
        _saved: false,
        _fail: false,
      })
    }

    excelRows.value = parsed
    excelSelectedRows.value = parsed
      .map((r, idx) => (!r._error ? idx : null))
      .filter(v => v !== null)

    const validCount = excelSelectedRows.value.length
    const errorCount = parsed.length - validCount
    ElMessage.success(
      `📊 Đọc được ${parsed.length} lệnh — ${validCount} hợp lệ` +
      (errorCount > 0 ? ` — ${errorCount} lỗi` : '') +
      (skippedRows > 0 ? ` (${skippedRows} dòng trống bỏ qua)` : '')
    )
  } catch (err) {
    ElMessage.error(`❌ Lỗi đọc file Excel: ${err.message}`)
    showExcelDialog.value = false
    resetExcelState()
  } finally {
    excelParsing.value = false
  }
}

const toggleExcelRow = (idx) => {
  if (excelRows.value[idx]?._error || excelSaving.value) return
  const pos = excelSelectedRows.value.indexOf(idx)
  if (pos === -1) excelSelectedRows.value.push(idx)
  else excelSelectedRows.value.splice(pos, 1)
}

const toggleSelectAll = (val) => {
  if (val) {
    excelSelectedRows.value = excelRows.value
      .map((r, i) => (!r._error ? i : null))
      .filter(v => v !== null)
  } else {
    excelSelectedRows.value = []
  }
}

const saveExcelRows = async () => {
  if (excelSelectedRows.value.length === 0) return

  const loginMsnv = loggedUser.userid.trim()
  if (!loginMsnv) {
    ElMessage.warning(t('lossTemSize.messages.loginMsnvMissing'))
    return
  }

  const indicesToSave = [...excelSelectedRows.value].sort((a, b) => {
    const modeA = excelRows.value[a]?.mode === 'IN' ? 0 : 1
    const modeB = excelRows.value[b]?.mode === 'IN' ? 0 : 1
    return modeA - modeB
  })

  excelSaving.value = true
  excelProgressDone.value = 0
  excelProgressPct.value = 0

  const total = indicesToSave.length
  let done = 0
  let errors = 0

  for (const idx of indicesToSave) {
    const row = excelRows.value[idx]

    let lenhInfo = null
    try {
      lenhInfo = await fetchLenhFromApi(row.ddbh)
    } catch {
      // nếu không lấy được thì để trống
    }

    const payload = {
      ddbh:     row.ddbh,
      msnv:     loginMsnv,
      msnv_out: '',
      ywpm:     lenhInfo?.ywpm    || '',   // ← bổ sung
      article:  lenhInfo?.sku     || '',   // ← bổ sung
      xie_ming: lenhInfo?.shoeType || '',  // ← bổ sung
      pairs:    lenhInfo?.orderQty || 0,   // ← bổ sung
      mode:     row.mode,
      reason:   row.reason,
      note:     'import_excel_chi_tiet',
      sizes:    row.sizes,
      date:     new Date().toISOString().split('T')[0],
      product_code: row.product_code || lenhInfo?.productCode || '',
    }

    try {
      await saveEntryToAPI(payload)
      excelRows.value[idx]._saved = true
      done++
    } catch (err) {
      console.error(`Lỗi lưu lệnh ${row.ddbh}:`, err)
      excelRows.value[idx]._fail = true
      errors++
    }

    excelProgressDone.value = done + errors
    excelProgressPct.value = Math.round(((done + errors) / total) * 100)
  }

  excelSaving.value = false

  if (errors === 0) {
    ElMessage.success(`✅ Đã lưu thành công ${done}/${total} lệnh từ file chi_tiet`)
    await new Promise(r => setTimeout(r, 800))
    showExcelDialog.value = false
    resetExcelState()
    if (!sseConnected.value) await loadSummaryFromAPI()
  } else {
    ElMessage.warning(`⚠️ Thành công ${done}, thất bại ${errors}/${total} — Kiểm tra các dòng đỏ`)
    excelSelectedRows.value = indicesToSave.filter(i => excelRows.value[i]?._fail)
  }
}

// ─── Lifecycle ─────────────────────────────────────────────────
onMounted(() => {
  window.addEventListener('resize', onResize)
  document.addEventListener('visibilitychange', onVisibilityChange)
  const savedLang = localStorage.getItem('lang')
  if (savedLang && ['vi', 'en', 'zh'].includes(savedLang)) locale.value = savedLang

  // Set tháng hiện tại làm mặc định
  const now = new Date()
  const mm = String(now.getMonth() + 1).padStart(2, '0')
  filterMonth.value = `${now.getFullYear()}-${mm}`

  loadSummaryFromAPI()
  connectSSE()
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
  document.removeEventListener('visibilitychange', onVisibilityChange)
  stopSSE()
})
</script>
<style lang="scss" scoped>
$blue:         #2563eb;
$blue-dark:    #1e40af;
$blue-light:   #eff6ff;
$blue-mid:     #bfdbfe;
$green:        #16a34a;
$green-light:  #f0fdf4;
$orange:       #ea580c;
$orange-light: #fff7ed;
$red:          #dc2626;
$red-light:    #fee2e2;
$yellow:       #d97706;
$yellow-light: #fef3c7;
$emerald:      #059669;
$emerald-light:#d1fae5;
$border:       #e5e7eb;
$border-light: #f3f4f6;
$bg:           #f8fafc;
$surface:      #ffffff;
$text:         #111827;
$text-2:       #6b7280;
$text-3:       #9ca3af;
$radius:       10px;
$radius-lg:    14px;
$shadow:       0 1px 3px rgba(0,0,0,.08), 0 1px 2px rgba(0,0,0,.04);
$shadow-md:    0 4px 8px rgba(0,0,0,.08), 0 2px 4px rgba(0,0,0,.04);

.col-resizable { position: relative; user-select: none; }
.resize-handle {
  position: absolute; right: 0; top: 0; bottom: 0; width: 5px;
  cursor: col-resize; background: rgba(255,255,255,.3);
  &:hover { background: rgba(255,255,255,.7); }
}

.lts-page {
  min-height: 100vh; background: $bg; color: $text;
  font-size: 13px; padding: 12px;
  display: flex; flex-direction: column; gap: 10px;
}

.lts-header {
  display: flex; align-items: center; justify-content: space-between;
  gap: 12px; background: $surface; border: 1px solid $border;
  border-radius: $radius-lg; padding: 12px 16px; box-shadow: $shadow;
  &__brand { display: flex; align-items: center; gap: 12px; min-width: 0; flex: 1; }
  &__right { display: flex; align-items: center; gap: 10px; flex-shrink: 0; }
}
.lts-login-info-box {
  display: flex; align-items: center; padding: 8px 12px;
  background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 12px;
  box-shadow: 0 2px 6px rgba(37,99,235,.08);
}
.lts-user-info { display: flex; flex-direction: column; line-height: 1.25; }
.lts-user-name { font-size: 12px; font-weight: 800; color: $blue-dark; word-break: break-word; }
.lts-user-id   { font-size: 10px; color: $text-2; }
.lts-btn-logout {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 10px 14px !important; font-weight: 800; border-radius: 12px;
  box-shadow: 0 3px 10px rgba(220,38,38,.18); height: 45px;
}
.lts-logo  { height: 52px; width: auto; object-fit: contain; flex-shrink: 0; }
.lts-eyebrow {
  font-size: 10px; font-weight: 700; letter-spacing: .12em; color: $blue;
  text-transform: uppercase; margin: 0 0 4px;
}
.lts-title {
  font-size: clamp(15px,2.5vw,20px); font-weight: 800; margin: 0;
  background: linear-gradient(135deg,$blue 0%,$blue-dark 100%);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
  background-clip: text; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.lts-btn-back {
  flex-shrink: 0; width: 38px; min-height: 38px; padding: 0 !important;
  border-radius: 10px; display: flex; align-items: center; justify-content: center; height: 45px;
  img { width: 18px; height: 18px; object-fit: contain; display: block; }
}

.lts-table-header {
  display: flex; justify-content: space-between; align-items: center;
  gap: 12px;
  padding: 14px 18px;                         // ← tăng padding
  border-bottom: 1px solid $border;
  background: #fafafa;
  margin-bottom: 0;

  &__tools {
    display: flex; align-items: center; gap: 10px; flex-shrink: 0; flex-wrap: wrap;

    // ← tăng size các nút và select bên trong
    :deep(.el-button) {
      height: 34px;
      font-size: 13px;
      font-weight: 600;
      padding: 0 14px;
    }

    :deep(.el-select) {
      .el-input__wrapper { height: 34px; }
      .el-input__inner   { font-size: 13px; }
    }
  }
}

.lts-filter-count-inline {
  font-size: 14px;     // ← to hơn
  font-weight: 600;
  color: $text;
  white-space: nowrap;
}
.lts-overview-grid {
  display: grid; grid-template-columns: repeat(3,1fr);
  gap: 0;        // ← bỏ gap, dùng border thay thế
  padding: 0;
}

.lts-stat-box {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding: 16px 8px;
  border-radius: 0;                            // ← bỏ radius, nằm trong grid liền nhau
  border: none;
  border-right: 1px solid $border;            // ← chỉ giữ border phải
  text-align: center; transition: transform .2s, box-shadow .2s;
  &:last-child { border-right: none; }
  &:hover { transform: translateY(-1px); box-shadow: $shadow-md; }
  span   { font-size: 10px; color: $text-2; margin-bottom: 4px; font-weight: 500; line-height: 1.3; text-transform: uppercase; letter-spacing: .04em; }
  strong { font-size: 22px; font-weight: 800; }
  &.is-in  { background: $green-light;  border-color: #bbf7d0; strong { color: $green; } }
  &.is-out { background: $orange-light; border-color: #fed7aa; strong { color: $orange; } }
  &.is-net { background: $blue-light;   border-color: $blue-mid; strong { color: $blue; } }
}

.lts-search-card {
  background: $surface; border: 1px solid $border; border-left: 4px solid $blue;
  border-radius: $radius-lg; padding: 14px; box-shadow: $shadow;
}
.lts-card__head {
  display: flex; align-items: center; justify-content: space-between;
  gap: 8px; margin-bottom: 12px;
}
.lts-search-grid {
  display: grid; grid-template-columns: 1fr auto; gap: 8px; align-items: center;
}
.lts-search-actions {
  display: flex; gap: 8px;
  .lts-search-btn { height: 38px; font-weight: 600; border-radius: $radius; white-space: nowrap; }
}
.lts-error-row { margin-top: 4px; }
.lts-error { color: $red; font-size: 12px; font-weight: 500; margin: 2px 0 0; }
.lts-input {
  :deep(.el-input__wrapper) {
    border-radius: $radius; border: 1.5px solid $border; box-shadow: none !important;
    transition: border-color .2s, box-shadow .2s;
  }
  :deep(.el-input__wrapper:hover) { border-color: $blue; }
  :deep(.el-input__wrapper.is-focus) {
    border-color: $blue; box-shadow: 0 0 0 3px rgba(37,99,235,.12) !important;
  }
  &--error :deep(.el-input__wrapper) {
    border-color: $red !important; background: #fef2f2 !important;
    box-shadow: 0 0 0 3px rgba(220,38,38,.1) !important;
  }
}

.lts-excel-btn {
  background: linear-gradient(135deg, #16a34a 0%, #15803d 100%) !important;
  border-color: #15803d !important;
  color: #fff !important;
  font-weight: 700 !important;
  &:hover {
    background: linear-gradient(135deg, #15803d 0%, #166534 100%) !important;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(22,163,74,.35) !important;
  }
}

.lts-body-layout {
  display: grid; grid-template-columns: 220px 1fr; gap: 12px;
  align-items: stretch; min-width: 0;
  transition: grid-template-columns .25s ease;
  &.is-collapsed { grid-template-columns: 52px 1fr; }
}

.lts-sidebar {
  background: $surface; border: 1px solid $border; border-radius: $radius-lg;
  padding: 12px; box-shadow: $shadow; display: flex; flex-direction: column;
  justify-content: flex-start; gap: 10px;
  position: sticky; top: 12px; height: fit-content; align-self: start;
  overflow: hidden; transition: width .25s ease, padding .25s ease;
}

.lts-sidebar-header {
  display: flex; align-items: center; justify-content: space-between;
  gap: 6px; padding-bottom: 8px; border-bottom: 1px solid $border;
}

.lts-sidebar-actions { display: flex; align-items: center; gap: 6px; flex-shrink: 0; }

.lts-toggle-btn {
  width: 30px; height: 30px; padding: 0; border: 1px solid $blue-mid;
  border-radius: 8px; background: #fff; color: $blue-dark; font-size: 18px;
  font-weight: 800; cursor: pointer; display: inline-flex; align-items: center;
  justify-content: center; transition: all .2s;
  &:hover { transform: translateY(-1px); background: $blue-light; border-color: $blue; }
}

.lts-filter-badge {
  display: inline-flex; align-items: center; background: $blue; color: #fff;
  border-radius: 999px; font-size: 11px; font-weight: 700; padding: 2px 8px;
  span { opacity: .65; font-weight: 500; margin-left: 1px; }
}

.lts-sidebar-filters { display: flex; flex-direction: column; gap: 6px; }

.lts-filter-group { display: flex; flex-direction: column; gap: 3px; }

.lts-filter-label {
  font-size: 10px; font-weight: 700; color: $text-2;
  text-transform: uppercase; letter-spacing: .06em;
}

// ← chỉnh toàn bộ input/select bên trong sidebar
.lts-filter-item {
  width: 100%;

  // el-input
  :deep(.el-input__wrapper) {
    border-radius: 8px;
    border: 1.5px solid $border;
    box-shadow: none !important;
    background: #fafafa;
    padding: 0 10px;
    height: 32px;
    transition: border-color .2s, background .2s;
    &:hover  { border-color: $blue; background: #fff; }
  }
  :deep(.el-input__wrapper.is-focus) {
    border-color: $blue;
    background: #fff;
    box-shadow: 0 0 0 3px rgba(37,99,235,.1) !important;
  }
  :deep(.el-input__inner) {
    font-size: 12px;
    color: $text;
    height: 32px;
  }

  // el-select
  :deep(.el-select__wrapper) {
    border-radius: 8px;
    border: 1.5px solid $border;
    box-shadow: none !important;
    background: #fafafa;
    height: 32px;
    font-size: 12px;
    transition: border-color .2s, background .2s;
    &:hover { border-color: $blue; background: #fff; }
  }
  :deep(.el-select__wrapper.is-focused) {
    border-color: $blue;
    background: #fff;
    box-shadow: 0 0 0 3px rgba(37,99,235,.1) !important;
  }

  // el-date-picker
  :deep(.el-date-editor.el-input__wrapper) {
    border-radius: 8px;
    border: 1.5px solid $border;
    box-shadow: none !important;
    background: #fafafa;
    height: 32px;
    padding: 0 10px;
    &:hover { border-color: $blue; background: #fff; }
  }
}

.lts-reset-btn {
  width: 100%; margin-top: 4px;
  border-radius: 8px !important;
  font-size: 12px !important;
  height: 32px !important;
}

.lts-sidebar.is-collapsed {
  padding: 10px 6px; align-items: center;
  .lts-sidebar-header { width: 100%; justify-content: center; padding-bottom: 0; border-bottom: none; }
  .lts-section-label, .lts-sidebar-filters, .lts-reset-btn, .lts-filter-badge { display: none; }
  .lts-sidebar-actions { width: 100%; justify-content: center; }
}

.lts-main { display: flex; flex-direction: column; gap: 8px; min-width: 0; }
.lts-page-size { min-width: 110px; }
.lts-section-label {
  display: flex; align-items: center; gap: 7px; font-size: 13px; font-weight: 700;
  letter-spacing: .07em; text-transform: uppercase; color: $text; margin-bottom: 10px;
}
.lts-state {
  display: flex; align-items: center; justify-content: center; gap: 8px;
  padding: 28px; text-align: center; color: $text-2; font-size: 13px;
  background: $surface; border: 1px solid $border; border-radius: $radius-lg;
  &--error { color: $red; background: $red-light; border-color: #fecaca; }
}
.lts-scroll-hint {
  display: none; font-size: 11px; color: $text-2; padding: 6px 10px;
  background: $blue-light; border: 1px dashed $blue-mid; border-radius: $radius;
}
.lts-table-wrap {
  overflow-x: auto;
  overflow-y: auto;        // ← đổi từ visible → auto
  max-height: 65vh;        // ← bắt buộc phải có, không có thì sticky không kích hoạt
  -webkit-overflow-scrolling: touch;
  overscroll-behavior-x: contain;
  border: 1px solid $border;
  border-radius: $radius-lg;
  box-shadow: $shadow;
  background: $surface;
}
.lts-table {
  width: 100%; border-collapse: collapse; font-size: 12px;
  min-width: 1050px; table-layout: auto;
  th, td {
    padding: 9px 8px; border-bottom: 1px solid $border;
    border-right: 1px solid $border-light; white-space: nowrap;
    &:last-child { border-right: none; }
  }
  th {
    position: sticky; top: 0; z-index: 3;
    background: linear-gradient(135deg,$blue 0%,$blue-dark 100%);
    color: #fff; font-weight: 700; text-align: center;
    font-size: 11px; text-transform: uppercase; letter-spacing: .04em;
  }
}
.lts-tbody-group > tr:first-child td { border-top: 2px solid $blue-mid; }
.lts-tbody-grand  > tr:first-child td { border-top: 3px solid $blue; }
.col-lenh  { text-align: left !important; padding-left: 12px !important; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.col-qty   { min-width: 80px; }
.col-loai  { min-width: 160px; }
.col-size  { min-width: 48px; font-size: 10px !important; }
.col-total { min-width: 56px; }
.td-lenh {
  vertical-align: top; padding: 10px 12px !important;
  background: linear-gradient(135deg,$blue-light 0%,#f0f9ff 100%) !important;
  border-right: 2px solid $blue-mid !important; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;}
.col-pc {
  min-width: 80px;
  max-width: 100px;
  width: 90px;
  text-align: center;
}
.td-pc {
  text-align: center;
  vertical-align: middle;
  background: #f8fafc;
  border-right: 2px solid $blue-mid !important;
  font-size: 11px;
  min-width: 80px;
  max-width: 100px;
  width: 90px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.td-pc strong {
  color: $blue-dark;
  font-weight: 700;
  font-size: 12px;
  word-break: break-word;
  white-space: normal;
}
.td-qty    { vertical-align: middle; text-align: center; background: #f8fafc; border-right: 2px solid $blue-mid !important; }
.td-loai   { text-align: center; vertical-align: middle; }
.td-center { text-align: center; font-weight: 600; }
.td-total  { text-align: center; font-weight: 700; background: $border-light; border-right: 2px solid $blue-mid !important; }
.td-subtotal-label { text-align: left; font-size: 11px; padding-left: 14px !important; font-weight: 700; }
.lenh-badge { font-size: 12px; font-weight: 800; color: $blue; margin-bottom: 3px; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.lenh-sub   { font-size: 11px; color: $text-2; margin-bottom: 2px; font-weight: 500; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.lenh-meta  { font-size: 10px; color: $text-3; font-style: italic; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
// override el-tag màu cho khớp với nền row
.tr-loss .el-tag--primary { background: #bdd7ee; border-color: #9ec6e0; color: #1e4d78; }
.tr-bao  .el-tag--warning  { background: #ffe699; border-color: #ffd966; color: #7d5a00; }
.tr-rem  .el-tag--success  { background: #c6e0b4; border-color: #a9d18e; color: #375623; }
.tr-subtotal td  { background: linear-gradient(180deg,$blue-light 0%,#dbeafe 100%); color: $blue-dark; border-top: 1.5px solid $blue-mid; }
.tr-grand-reason td { background: #f0fdf4; border-top: 1.5px solid #86efac; font-weight: 600; }
.tr-grand-net td { background: linear-gradient(135deg,#bfdbfe 0%,#93c5fd 100%); color: #1e40af; border-top: 2px solid #2563eb; font-weight: 700; }
.lts-pagination-row { display: flex; justify-content: flex-end; padding: 6px 0; }

.lts-dialog {
  :deep(.el-dialog)         { margin: 4vh auto !important; border-radius: 16px !important; overflow: hidden; }
  :deep(.el-dialog__header) { padding: 14px 18px; border-bottom: 1px solid $border; }
  :deep(.el-dialog__title)  { font-size: 16px; font-weight: 800; color: $blue-dark; }
  :deep(.el-dialog__body)   { padding: 14px 16px; max-height: 72vh; overflow-y: auto; }
  :deep(.el-dialog__footer) { padding: 12px 16px; border-top: 1px solid $border; }
}
.lts-dialog-body { display: flex; flex-direction: column; gap: 12px; }
.lts-info-bar {
  display: flex; gap: 12px; align-items: flex-start; justify-content: space-between;
  flex-wrap: wrap; padding: 12px 14px; border: 1px solid $blue-mid;
  border-radius: $radius; background: linear-gradient(135deg,$blue-light 0%,#f0f9ff 100%);
}
.lts-info-grid { display: grid; grid-template-columns: repeat(3,auto); gap: 4px 20px; font-size: 12px; }
.lts-info-item {
  display: flex; gap: 6px; align-items: center;
  span   { color: $text-2; font-size: 11px; }
  strong { color: $text; font-weight: 700; }
}
.lts-mode-badge {
  padding: 6px 12px; border-radius: 999px; font-size: 11px; font-weight: 800;
  letter-spacing: .04em; flex-shrink: 0;
  &.is-in  { background: $green-light;  color: $green; border: 1px solid #a7f3d0; }
  &.is-out { background: $red-light;    color: $red;   border: 1px solid #fecaca; }
}
.lts-step-block  { border: 1px solid $border; border-radius: $radius-lg; padding: 12px; background: $surface; }
.lts-step-msnv   { background: #fffbeb; border-color: #f59e0b; }
.lts-step-label  {
  display: flex; align-items: center; gap: 8px; font-size: 12px; font-weight: 700;
  color: $blue-dark; text-transform: uppercase; letter-spacing: .06em; margin-bottom: 10px;
}
.lts-step-num {
  display: inline-flex; align-items: center; justify-content: center;
  width: 22px; height: 22px; background: $blue; color: #fff;
  border-radius: 50%; font-size: 11px; font-weight: 700; flex-shrink: 0;
}
.lts-msnv-box { display: flex; flex-direction: column; gap: 6px; }
.lts-msnv-input {
  width: 100%;
  :deep(.el-input__wrapper) { border-radius: $radius; border: 1.5px solid #f59e0b; box-shadow: none !important; background: #fff; }
  :deep(.el-input__inner)   { font-size: 14px; font-weight: 600; }
  :deep(.el-input__wrapper:hover) { border-color: #d97706; }
  :deep(.el-input__wrapper.is-focus) { border-color: #d97706; box-shadow: 0 0 0 3px rgba(245,158,11,.15) !important; }
}
.msnv-help { margin-top: 2px; font-size: 11px; font-weight: 600; color: #92400e; }
.lts-reason-list { display: flex; flex-direction: column; gap: 8px; }
.lts-reason-item {
  display: flex; align-items: center; gap: 10px; padding: 10px 12px;
  border: 1.5px solid $border; border-radius: $radius; cursor: pointer;
  transition: all .2s; background: #fff; user-select: none;
  &:hover { border-color: $blue; background: $blue-light; transform: translateX(3px); }
  .reason-icon  { font-size: 16px; color: $text-2; flex-shrink: 0; }
  .reason-text  { flex: 1; min-width: 0; }
  .reason-name  { font-weight: 600; font-size: 12px; }
  .reason-sub   { font-size: 10px; color: $text-3; margin-top: 1px; }
  .reason-check { font-size: 14px; font-weight: 800; flex-shrink: 0; }
  &.reason--danger.is-active  { border-color: $red;     background: $red-light;     box-shadow: 0 0 0 3px rgba(220,38,38,.08);  .reason-icon,.reason-name,.reason-check { color: $red; } }
  &.reason--warning.is-active { border-color: $yellow;  background: $yellow-light;  box-shadow: 0 0 0 3px rgba(217,119,6,.08);  .reason-icon,.reason-name,.reason-check { color: $yellow; } }
  &.reason--success.is-active { border-color: $emerald; background: $emerald-light; box-shadow: 0 0 0 3px rgba(5,150,105,.08);  .reason-icon,.reason-name,.reason-check { color: $emerald; } }
}
.lts-reason-note {
  display: flex; align-items: center; gap: 6px; margin-top: 10px; padding: 8px 10px;
  background: $yellow-light; border: 1px solid #fcd34d; border-radius: $radius;
  font-size: 11px; color: #78350f; font-weight: 500;
}
.lts-size-grid { display: grid; grid-template-columns: repeat(auto-fill,minmax(78px,1fr)); gap: 8px; }
.lts-size-cell {
  display: flex; flex-direction: column; align-items: center; gap: 4px;
  padding: 8px 4px; border: 1.5px solid $border; border-radius: $radius;
  background: #fff; transition: all .2s;
  .size-lbl    { font-size: 11px; font-weight: 700; }
  .size-remain { font-size: 9px; color: $text-2; font-weight: 500; }
  :deep(.el-input-number)  { width: 100%; }
  :deep(.el-input__inner)  { text-align: center; font-size: 13px; font-weight: 600; padding: 4px; }
  &.is-filled { border-color: $emerald; background: $emerald-light; box-shadow: 0 0 0 2px rgba(5,150,105,.1); .size-lbl { color: $emerald; font-weight: 800; } :deep(.el-input__inner) { color: $emerald; font-weight: 700; } }
  &.is-out    { background: linear-gradient(135deg,$yellow-light 0%,#fffbeb 100%); border-color: #fcd34d; }
}
.lts-dialog-footer { display: flex; align-items: center; justify-content: space-between; gap: 10px; flex-wrap: wrap; width: 100%; }
.dlg-note    { font-size: 12px; color: $text-2; }
.dlg-actions { display: flex; gap: 8px; flex-wrap: wrap; }
.txt-danger  { color: $red; }
.txt-success { color: $green; }

.lts-excel-dialog {
  :deep(.el-dialog) { max-height: 95vh; margin: 2vh auto !important; }
  :deep(.el-dialog__body) { padding: 12px 14px; max-height: 78vh; overflow-y: auto; }
}
.lts-excel-body { display: flex; flex-direction: column; gap: 10px; }
.lts-excel-toolbar {
  display: flex; align-items: center; flex-wrap: wrap; gap: 10px 16px;
  padding: 10px 14px; background: $blue-light; border: 1px solid $blue-mid;
  border-radius: $radius; position: sticky; top: 0; z-index: 5;
}
.lts-excel-info { display: flex; gap: 6px; align-items: center; flex-wrap: wrap; }
.lts-excel-group {
  display: flex; align-items: center; gap: 8px;
}
.lts-excel-group-label {
  font-size: 11px; font-weight: 700; color: $text-2;
  text-transform: uppercase; letter-spacing: .05em; white-space: nowrap;
}
.lts-excel-progress {
  display: flex; flex-direction: column; gap: 6px;
  padding: 10px 14px; background: $yellow-light;
  border: 1px solid #fcd34d; border-radius: $radius;
}
.lts-excel-progress-label {
  display: flex; justify-content: space-between; align-items: center;
  font-size: 12px; color: #92400e; font-weight: 600;
}
.lts-excel-table-wrap {
  overflow-x: auto; overflow-y: auto;
  max-height: 48vh; border: 1px solid $border;
  border-radius: $radius; -webkit-overflow-scrolling: touch;
}
.lts-excel-table {
  width: 100%; border-collapse: collapse; font-size: 11px; min-width: 900px;
  th, td {
    padding: 6px 8px; border-bottom: 1px solid $border;
    border-right: 1px solid $border-light; white-space: nowrap;
    &:last-child { border-right: none; }
  }
  th {
    position: sticky; top: 0; z-index: 3;
    background: linear-gradient(135deg,$blue 0%,$blue-dark 100%);
    color: #fff; font-weight: 700; text-align: center;
    font-size: 10px; text-transform: uppercase; letter-spacing: .04em;
  }
  .th-chk    { width: 36px; }
  .th-no     { width: 32px; }
  .th-ddbh   { min-width: 130px; }
  .th-mode   { min-width: 70px; }
  .th-msnv   { min-width: 90px; }
  .th-sz     { min-width: 38px; }
  .th-loss   { min-width: 70px; }
  .th-status { min-width: 80px; }
  .excel-row {
    cursor: pointer; transition: background .15s;
    &:hover:not(.is-error) { background: $blue-light; }
    &.is-selected { background: #f0fdf4 !important; }
    &.is-saved    { background: #dcfce7 !important; opacity: .8; }
    &.is-fail     { background: $red-light !important; }
    &.is-error    { opacity: .55; cursor: not-allowed; background: #fef2f2; }
  }
  .td-chk    { text-align: center; }
  .td-no     { text-align: center; color: $text-3; font-size: 10px; }
  .td-ddbh   { font-weight: 700; color: $blue; font-size: 11px; }
  .td-mode   { text-align: center; }
  .td-msnv   { font-family: monospace; font-size: 10px; text-align: center; }
  .td-sz     { text-align: center; color: $text-2; }
  .td-sz.has-val { color: $green; font-weight: 700; background: $green-light; }
  .td-loss   { text-align: center; font-weight: 600; }
  .td-status { text-align: center; }
}
.fade-slide-enter-active, .fade-slide-leave-active {
  transition: opacity .25s ease, transform .25s cubic-bezier(.34,1.56,.64,1);
}
.fade-slide-enter-from, .fade-slide-leave-to { opacity: 0; transform: translateY(-10px); }

@media (max-width: 1024px) { .lts-body-layout { grid-template-columns: 190px 1fr; } }

@media (max-width: 768px) {
  .lts-page   { padding: 8px; }
  .lts-header { padding: 10px 12px; }
  .lts-logo   { height: 40px; }
  .lts-title  { font-size: 14px; }
  .lts-login-info-box { display: none; }
  .lts-overview-grid  { grid-template-columns: repeat(3,1fr); }
  .lts-stat-box { padding: 8px 4px; span { font-size: 9px; } strong { font-size: 15px; } }
  .lts-search-grid    { grid-template-columns: 1fr; .lts-search-actions { flex-direction: row; flex-wrap: wrap; } }
  .lts-body-layout    { grid-template-columns: 1fr; }
  .lts-sidebar        { position: static; }
  .lts-sidebar-filters { display: grid; grid-template-columns: repeat(2,1fr); gap: 8px; }
  .col-lenh, .td-lenh { position: sticky; left: 0; z-index: 2; background: #f0f9ff !important; }
  .col-lenh           { z-index: 4; }
  .lts-scroll-hint    { display: flex; }
  .lts-info-grid      { grid-template-columns: repeat(2,auto); }
  .lts-dialog {
    :deep(.el-dialog) { width: 97vw !important; margin: 2vh auto !important; }
    :deep(.el-dialog__body) { padding: 10px 12px; max-height: 75vh; }
  }
  .lts-dialog-footer  { flex-direction: column; align-items: stretch; .dlg-actions { justify-content: stretch; .el-button { flex: 1; } } }
  .lts-size-grid      { grid-template-columns: repeat(auto-fill,minmax(68px,1fr)); }
  .lts-table-header   { flex-direction: column; align-items: stretch; gap: 8px; &__tools { justify-content: space-between; } }
  .lts-pagination-row { justify-content: center; overflow-x: auto; }
  .lts-excel-toolbar  { flex-direction: column; align-items: flex-start; }
}

@media (max-width: 480px) {
  .lts-overview-grid  { gap: 6px; }
  .lts-stat-box       { span { font-size: 8px; } strong { font-size: 14px; } }
  .lts-info-grid      { grid-template-columns: 1fr 1fr; }
  .lts-size-grid      { grid-template-columns: repeat(auto-fill,minmax(60px,1fr)); }
  .lts-sidebar-filters { grid-template-columns: 1fr; }
}
</style>