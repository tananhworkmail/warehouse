import { locale } from '@/hooks/i18n'

const vi = {
  common: { back: 'Trở về', refresh: 'Làm mới', loading: 'Đang tải', submit: 'Hoàn tất', submitting: 'Đang lưu...', cancel: 'Hủy', reset: 'Xóa dữ liệu', draft: 'Lưu nháp', edit: 'Chỉnh sửa', close: 'Đóng', view: 'Xem', history: 'Lịch sử' },
  entry: { eyebrow: 'PHÒNG THÍ NGHIỆM · GORE-TEX', hint: 'Chọn khu vực bạn muốn truy cập.', portalTag: 'NHẬP VÀ TRA CỨU DỮ LIỆU', portalTitle: 'Portal nhập 3 biểu mẫu', portalDescription: 'Nhập mới, xem danh sách và chỉnh sửa các biểu mẫu GORE-TEX.', openPortal: 'Mở Portal', dashboardTag: 'TỔNG HỢP VÀ THEO DÕI', dashboardDescription: 'Khu vực dashboard tổng hợp dữ liệu GORE-TEX.', openDashboard: 'Mở Dashboard' },
  portal: { eyebrow: 'PHÒNG THÍ NGHIỆM · GORE-TEX', title: 'Chọn biểu mẫu', description: 'Nhập mới hoặc xem lại các biểu mẫu GORE-TEX đã lưu.', submitted: 'Đã lưu “{title}” thành công.', submittedDefault: 'Đã lưu biểu mẫu thành công.', list: 'Xem danh sách đã nhập', open: 'Mở biểu mẫu' },
  forms: {
    waterproof: { title: 'BÁO BIỂU KIỂM TRA CHẤT LƯỢNG GIÀY CHỐNG THẤM NƯỚC', short: 'Kiểm tra chống thấm' },
    centrifugal: { title: 'BÁO BIỂU GIÀY THÀNH PHẨM THỬ NGHIỆM LI TÂM', short: 'Thử nghiệm li tâm' },
    analysis: { title: 'BIỂU PHÂN TÍCH NGUYÊN NHÂN VÀ CẢI THIỆN THỬ NGHIỆM LI TÂM / THỬ NƯỚC', short: 'Phân tích và cải thiện' },
  },
  defects: { toe: 'Mũi thấm nước', heel: 'Gót thấm nước', medial: 'Hong trong thấm nước', lateral: 'Hong ngoài thấm nước', material: 'Vật tư không đạt', attaching: 'Dán đế lệch', wrinkled: 'Ép đế nhăn', zigzag: 'Zíc zắc hở', bonding: 'Hở keo đế' },
  dashboard: { eyebrow: 'GORE-TEX · BÁO CÁO THEO TUẦN', title: 'Dashboard kiểm nghiệm', description: 'Theo dõi tỷ lệ đạt và lỗi chống thấm theo từng tuần.', week: 'Tuần {week}/{year}', selectWeek: 'Chọn tuần', loadError: 'Không tải được dashboard.', aggregating: 'Đang tổng hợp dữ liệu tuần…', chart1: 'Tỷ lệ đạt theo dạng giày', chart2: 'Xu hướng tỷ lệ đạt SUTER', chart3: 'So sánh hai tuần', chart4: 'Xu hướng tỷ lệ đạt R.RDY', chart5: 'Kết quả lỗi chống thấm', noPass: 'Chưa có dữ liệu PASS/FAIL trong tuần này.', noTrend: 'Chưa có dữ liệu xu hướng trong tuần này.', noCompare: 'Chưa có dữ liệu để so sánh hai tuần.', noDefect: 'Chưa ghi nhận lỗi chống thấm trong tuần này.', target: 'Mục tiêu 99%', passed: 'đạt', defects: 'lỗi', cumulative: 'lũy kế' },
  history: { title: 'Danh sách biểu mẫu đã nhập', chooseForm: 'Chọn biểu mẫu', filter: 'Lọc loại biểu mẫu', records: 'bản ghi', loading: 'Đang tải danh sách...', emptyTitle: 'Chưa có dữ liệu', emptyText: 'Biểu mẫu này chưa có bản ghi nào được lưu.', line: 'Chuyền / Line', style: 'Dạng giày', testDate: 'Ngày thử nghiệm', improvementDate: 'Ngày cải thiện', inspectionDate: 'Ngày kiểm tra', updated: 'Cập nhật gần nhất', details: 'Xem chi tiết' },
  actions: { noDraft: 'Chưa có bản nháp', savedAt: 'Lưu lúc {time}', required: 'Cần nhập đầy đủ và đúng định dạng trước khi hoàn tất.', savedDraft: 'Đã lưu nháp.', missing: 'Còn {count} ô chưa nhập hoặc không hợp lệ.', confirmReset: 'Xóa toàn bộ dữ liệu trong biểu mẫu?', resetDone: 'Đã xóa dữ liệu.', addRow: 'Thêm dòng', choose: 'Chọn' },
  review: { revisionMode: 'Đang xem phiên bản lịch sử', editMode: 'Chế độ chỉnh sửa', reviewMode: 'Chế độ xem lại', revisionHint: 'Ô có dấu góc là nội dung mới sau chỉnh sửa. Bấm vào ô để xem nội dung cũ.', editHint: 'Bạn có thể thay đổi dữ liệu và gửi lại biểu mẫu.', requiredHint: 'Cần nhập đầy đủ tất cả các ô trước khi submit.', lockedHint: 'Dữ liệu đang được khóa để tránh chỉnh sửa ngoài ý muốn.', current: 'Về phiên bản hiện tại', cancelEdit: 'Hủy chỉnh sửa', sending: 'Đang gửi...' },
  revisions: { title: 'Lịch sử chỉnh sửa', versions: '{count} phiên bản trước đã được lưu', empty: 'Biểu mẫu chưa có lần chỉnh sửa nào.', version: 'Phiên bản {number}', savedBefore: 'Lưu trước lần chỉnh sửa lúc {time}', changes: '{count} nội dung thay đổi', viewVersion: 'Xem phiên bản', noTime: 'Không có thời gian' },
}

const en = {
  common: { back: 'Back', refresh: 'Refresh', loading: 'Loading', submit: 'Submit', submitting: 'Saving...', cancel: 'Cancel', reset: 'Clear data', draft: 'Save draft', edit: 'Edit', close: 'Close', view: 'View', history: 'History' },
  entry: { eyebrow: 'LABORATORY · GORE-TEX', hint: 'Choose the area you want to access.', portalTag: 'DATA ENTRY AND LOOKUP', portalTitle: 'Three-form portal', portalDescription: 'Create, review and edit GORE-TEX forms.', openPortal: 'Open Portal', dashboardTag: 'SUMMARY AND MONITORING', dashboardDescription: 'GORE-TEX data summary dashboard.', openDashboard: 'Open Dashboard' },
  portal: { eyebrow: 'LABORATORY · GORE-TEX', title: 'Choose a form', description: 'Create a form or review saved GORE-TEX forms.', submitted: '“{title}” was saved successfully.', submittedDefault: 'The form was saved successfully.', list: 'View submitted forms', open: 'Open form' },
  forms: {
    waterproof: { title: 'DAILY QUALITY REPORT FOR WATERPROOF SHOES', short: 'Waterproof inspection' },
    centrifugal: { title: 'FINISHED SHOE CENTRIFUGAL TEST REPORT', short: 'Centrifugal test' },
    analysis: { title: 'CAUSE ANALYSIS AND IMPROVEMENT FORM FOR CENTRIFUGAL / WATER TEST', short: 'Analysis and improvement' },
  },
  defects: { toe: 'Toe leakage', heel: 'Heel leakage', medial: 'Medial quarter leakage', lateral: 'Lateral quarter leakage', material: 'Poor material', attaching: 'Sole attaching misalignment', wrinkled: 'Wrinkled sole', zigzag: 'Open zigzag stitch', bonding: 'Poor sole bonding' },
  dashboard: { eyebrow: 'GORE-TEX · WEEKLY REPORT', title: 'Test dashboard', description: 'Monitor weekly pass rates and waterproof defects.', week: 'Week {week}/{year}', selectWeek: 'Select week', loadError: 'Unable to load the dashboard.', aggregating: 'Aggregating weekly data…', chart1: 'Pass rate by style', chart2: 'SUTER pass-rate trend', chart3: 'Two-week comparison', chart4: 'R.RDY pass-rate trend', chart5: 'Waterproof defect results', noPass: 'No PASS/FAIL data for this week.', noTrend: 'No trend data for this week.', noCompare: 'No data available for two-week comparison.', noDefect: 'No waterproof defects recorded this week.', target: 'Target 99%', passed: 'passed', defects: 'defects', cumulative: 'cumulative' },
  history: { title: 'Submitted forms', chooseForm: 'Choose form', filter: 'Filter form type', records: 'records', loading: 'Loading forms...', emptyTitle: 'No data', emptyText: 'No records have been saved for this form.', line: 'Line', style: 'Style name', testDate: 'Test date', improvementDate: 'Improvement date', inspectionDate: 'Inspection date', updated: 'Last updated', details: 'View details' },
  actions: { noDraft: 'No draft saved', savedAt: 'Saved at {time}', required: 'Complete every field with valid data before submitting.', savedDraft: 'Draft saved.', missing: '{count} fields are empty or invalid.', confirmReset: 'Clear all data in this form?', resetDone: 'Data cleared.', addRow: 'Add row', choose: 'Select' },
  review: { revisionMode: 'Viewing a historical version', editMode: 'Edit mode', reviewMode: 'Review mode', revisionHint: 'Corner-marked cells contain newer values. Click a cell to view its old value.', editHint: 'You can update the data and submit the form again.', requiredHint: 'Complete every required field before submitting.', lockedHint: 'Data is locked to prevent unintended changes.', current: 'Current version', cancelEdit: 'Cancel editing', sending: 'Sending...' },
  revisions: { title: 'Edit history', versions: '{count} previous versions saved', empty: 'This form has no edit history.', version: 'Version {number}', savedBefore: 'Saved before the edit at {time}', changes: '{count} changes', viewVersion: 'View version', noTime: 'Time unavailable' },
}

const zh = {
  common: { back: '返回', refresh: '刷新', loading: '加载中', submit: '提交', submitting: '保存中...', cancel: '取消', reset: '清除数据', draft: '保存草稿', edit: '编辑', close: '关闭', view: '查看', history: '历史记录' },
  entry: { eyebrow: '实验室 · GORE-TEX', hint: '请选择要访问的区域。', portalTag: '数据录入与查询', portalTitle: '三份表单入口', portalDescription: '新建、查看及编辑 GORE-TEX 表单。', openPortal: '打开表单入口', dashboardTag: '汇总与跟踪', dashboardDescription: 'GORE-TEX 数据汇总仪表板。', openDashboard: '打开仪表板' },
  portal: { eyebrow: '实验室 · GORE-TEX', title: '选择表单', description: '新建或查看已保存的 GORE-TEX 表单。', submitted: '“{title}”保存成功。', submittedDefault: '表单保存成功。', list: '查看已录入表单', open: '打开表单' },
  forms: {
    waterproof: { title: '防水鞋品质检验日报表', short: '防水检验' },
    centrifugal: { title: '成品鞋离心试验报告', short: '离心试验' },
    analysis: { title: '离心／水测试原因分析与改善表', short: '分析与改善' },
  },
  defects: { toe: '鞋头漏水', heel: '后套漏水', medial: '内腰漏水', lateral: '外腰漏水', material: '材料不良', attaching: '贴底移位', wrinkled: '底皱', zigzag: '万能车线松开', bonding: '底开胶' },
  dashboard: { eyebrow: 'GORE-TEX · 周报', title: '测试仪表板', description: '按周跟踪合格率和防水缺陷。', week: '第 {week} 周 / {year}', selectWeek: '选择周', loadError: '无法加载仪表板。', aggregating: '正在汇总每周数据…', chart1: '各鞋型合格率', chart2: 'SUTER 合格率趋势', chart3: '两周比较', chart4: 'R.RDY 合格率趋势', chart5: '防水缺陷结果', noPass: '本周暂无 PASS/FAIL 数据。', noTrend: '本周暂无趋势数据。', noCompare: '暂无两周比较数据。', noDefect: '本周未记录防水缺陷。', target: '目标 99%', passed: '合格', defects: '缺陷', cumulative: '累计' },
  history: { title: '已录入表单列表', chooseForm: '选择表单', filter: '筛选表单类型', records: '条记录', loading: '正在加载表单...', emptyTitle: '暂无数据', emptyText: '此表单尚未保存任何记录。', line: '线别', style: '型体', testDate: '测试日期', improvementDate: '改善日期', inspectionDate: '检查日期', updated: '最近更新', details: '查看详情' },
  actions: { noDraft: '尚无草稿', savedAt: '保存于 {time}', required: '提交前请填写完整且有效的数据。', savedDraft: '草稿已保存。', missing: '还有 {count} 个栏位为空或无效。', confirmReset: '确定清除此表单的全部数据吗？', resetDone: '数据已清除。', addRow: '新增一行', choose: '请选择' },
  review: { revisionMode: '正在查看历史版本', editMode: '编辑模式', reviewMode: '查看模式', revisionHint: '带角标的单元格为修改后的新内容，点击可查看旧内容。', editHint: '您可以修改数据并重新提交表单。', requiredHint: '提交前请填写所有必填栏位。', lockedHint: '数据已锁定，以避免意外修改。', current: '返回当前版本', cancelEdit: '取消编辑', sending: '提交中...' },
  revisions: { title: '编辑历史', versions: '已保存 {count} 个旧版本', empty: '此表单尚无编辑历史。', version: '版本 {number}', savedBefore: '保存于 {time} 编辑之前', changes: '{count} 项变更', viewVersion: '查看版本', noTime: '无时间信息' },
}

const messages = { vi, en, zh }

function resolve(object, path) {
  return path.split('.').reduce((value, key) => value?.[key], object)
}

export function useGoreTexI18n() {
  const t = (path, params = {}) => {
    const value = resolve(messages[locale.value] || vi, path) ?? resolve(vi, path) ?? path
    if (typeof value !== 'string') return value
    return value.replace(/\{(\w+)\}/g, (_, key) => params[key] ?? `{${key}}`)
  }
  const dateLocale = () => ({ vi: 'vi-VN', en: 'en-GB', zh: 'zh-CN' })[locale.value] || 'vi-VN'
  return { t, locale, dateLocale }
}
