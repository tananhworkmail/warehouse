import { createRouter, createWebHistory } from "vue-router";
import Landingpage from "@/views/components/Landingpage.vue";
import Warehouse from "@/views/components/Warehouse/Warehouse.vue";
import WarehouseKVT from "@/views/components/WarehouseKVT/WarehouseKVT.vue";
import WarehouseKTP from "@/views/components/WarehouseKTP/WarehouseKTP.vue";
import Trace from "@/views/components/Warehouse/Trace/Scan/index.vue";
import TraceReport from "@/views/components/Warehouse/Trace/report.vue";
import LossTemSize from "@/views/components/WarehouseKVT/LossTemSize/index.vue";
import LoginPage from "@/views/components/N8d/LoginPage.vue";
import Laboratory from "@/views/components/WarehousePTN/Laboratory.vue";
import { ElMessage } from "element-plus";
import { clearSession } from "@/utils/appStoreAndSecurity";

const routes = [
  {
    path: "/",
    name: "home",
    component: Landingpage,
  },
  {
    path: "/login",
    name: "login",
    component: LoginPage,
    meta: { requiresGuest: true },
  },
  {
    path: "/laboratory",
    name: "laboratory",
    component: Laboratory,
  },
  {
    path: "/laboratory/forms",
    component: () =>
      import("@/views/components/WarehousePTN/GoreTex/GoreTexPortal.vue"),
    children: [
      {
        path: "",
        name: "laboratory-forms",
        component: () =>
          import("@/views/components/WarehousePTN/GoreTex/views/EntryView.vue"),
        meta: { title: "GORE-TEX" },
      },
      {
        path: "portal",
        name: "laboratory-forms-portal",
        component: () =>
          import("@/views/components/WarehousePTN/GoreTex/views/HomeView.vue"),
        meta: { title: "Chọn biểu mẫu GORE-TEX" },
      },
      {
        path: "dashboard",
        name: "laboratory-forms-dashboard",
        component: () =>
          import(
            "@/views/components/WarehousePTN/GoreTex/views/DashboardView.vue"
          ),
        meta: { title: "Dashboard GORE-TEX" },
      },
      {
        path: "history",
        name: "laboratory-forms-history",
        component: () =>
          import(
            "@/views/components/WarehousePTN/GoreTex/views/HistoryView.vue"
          ),
        meta: { title: "Biểu mẫu GORE-TEX đã nhập" },
      },
      {
        path: "kiem-tra-chong-tham",
        name: "laboratory-waterproof-form",
        component: () =>
          import(
            "@/views/components/WarehousePTN/GoreTex/views/WaterproofFormView.vue"
          ),
        meta: { title: "Kiểm tra chất lượng giày chống thấm" },
      },
      {
        path: "thu-nghiem-li-tam",
        name: "laboratory-centrifugal-form",
        component: () =>
          import(
            "@/views/components/WarehousePTN/GoreTex/views/CentrifugalFormView.vue"
          ),
        meta: { title: "Giày thành phẩm thử nghiệm li tâm" },
      },
      {
        path: "phan-tich-cai-thien",
        name: "laboratory-analysis-form",
        component: () =>
          import(
            "@/views/components/WarehousePTN/GoreTex/views/AnalysisFormView.vue"
          ),
        meta: { title: "Phân tích nguyên nhân và cải thiện" },
      },
    ],
  },
  {
    path: "/warehouse-khode",
    name: "warehouse-khode",
    component: Warehouse,
  },
  {
    path: "/trace",
    name: "trace",
    component: Trace,
  },
  {
    path: "/trace/:orderNo",
    name: "tracereport",
    component: TraceReport,
  },
  {
    path: "/warehouse-kvt",
    name: "warehouse-kvt",
    component: WarehouseKVT,
  },
  {
    path: "/warehouse-kvt/HumTemp",
    name: "warehouse-kvt-humtemp",
    component: WarehouseKVT,
  },
  {
    path: "/warehouse-ktp",
    name: "warehouse-ktp",
    component: WarehouseKTP,
  },
  {
    path: "/warehouse-ktp/move-order",
    name: "warehouse-ktp-move-order",
    component: () => import("@/views/components/WarehouseKTP/MoveOrderKTP.vue"),
  },
  {
    path: "/warehouse-ktp/scan-in/login",
    name: "warehouse-ktp-scan-in-login",
    component: () => import("@/views/components/WarehouseKTP/ScanInLogin.vue"),
    meta: {
      requiresGuest: true,
      authenticatedRedirect: "/warehouse-ktp/scan-in",
    },
  },
  {
    path: "/warehouse-ktp/scan-out/login",
    name: "warehouse-ktp-scan-out-login",
    component: () => import("@/views/components/WarehouseKTP/ScanOutLogin.vue"),
    meta: {
      requiresGuest: true,
      authenticatedRedirect: "/warehouse-ktp/scan-out",
    },
  },
  {
    path: "/warehouse-ktp/scan-ship/login",
    name: "warehouse-ktp-scan-ship-login",
    component: () => import("@/views/components/WarehouseKTP/ScanShipLogin.vue"),
    meta: {
      requiresGuest: true,
      authenticatedRedirect: "/warehouse-ktp/scan-ship",
    },
  },
  {
    path: "/warehouse-ktp/scan-in",
    name: "warehouse-ktp-scan-in",
    component: () => import("@/views/components/WarehouseKTP/ScanInKTP.vue"),
    meta: { requiresAuth: true, requiresKTPScanInPermission: true },
  },
  {
    path: "/warehouse-ktp/scan-out",
    name: "warehouse-ktp-scan-out",
    component: () => import("@/views/components/WarehouseKTP/ScanOutKTP.vue"),
    meta: { requiresAuth: true, requiresKTPScanOutPermission: true },
  },
  {
    path: "/warehouse-ktp/scan-ship",
    name: "warehouse-ktp-scan-ship",
    component: () => import("@/views/components/WarehouseKTP/ScanShipKTP.vue"),
    meta: { requiresAuth: true, requiresKTPScanShipPermission: true },
  },
  {
    path: "/loss-tem-size",
    name: "loss-tem-size",
    component: LossTemSize,
    meta: { requiresAuth: true },
  },

  // ===== MỚI: Module N8d (Kiểm kê Kho Vật Tư) =====
  {
    path: "/n8d",
    name: "n8d-home",
    component: () => import("@/views/components/N8d/MainPage.vue"),
  },
  {
    path: "/n8d/login",
    name: "n8d-login",
    component: LoginPage,
    meta: { requiresGuest: true },
  },
  {
    path: "/n8d/tracking",
    name: "n8d-tracking",
    component: () => import("@/views/components/N8d/TrackingPage.vue"),
  },
  {
    path: "/n8d/inventory",
    name: "n8d-inventory",
    component: () => import("@/views/components/N8d/InventoryPage.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

const AUTH_TOKEN_SCHEMA = "ktp-jwt-v2";

function ktpLoginRoute(path) {
  if (path.startsWith("/warehouse-ktp/scan-in")) {
    return "warehouse-ktp-scan-in-login";
  }
  if (path.startsWith("/warehouse-ktp/scan-out")) {
    return "warehouse-ktp-scan-out-login";
  }
  if (path.startsWith("/warehouse-ktp/scan-ship")) {
    return "warehouse-ktp-scan-ship-login";
  }
  return "login";
}

router.beforeEach(async (to) => {
  if (localStorage.getItem("AUTH_TOKEN_SCHEMA") !== AUTH_TOKEN_SCHEMA) {
    clearSession();
    localStorage.setItem("AUTH_TOKEN_SCHEMA", AUTH_TOKEN_SCHEMA);
  }
  const token = localStorage.getItem("token");
  const hasSession = Boolean(token && localStorage.getItem("USERID"));

  // Yêu cầu đăng nhập cho khu vực N8d
  if (to.meta.requiresN8dAuth && !hasSession) {
    return { name: "n8d-login", query: { redirect: to.fullPath } };
  }

  // Yêu cầu đăng nhập chung
  if (to.meta.requiresAuth && !hasSession) {
    return {
      name: ktpLoginRoute(to.path),
      query: { redirect: to.fullPath },
    };
  }

  // Trang chỉ dành cho khách (chưa đăng nhập) — nếu đã login thì redirect đi nơi khác
  if (to.meta.requiresGuest && hasSession) {
    const redirect =
      typeof to.query.redirect === "string" && to.query.redirect.startsWith("/")
        ? to.query.redirect
        : typeof to.meta.authenticatedRedirect === "string"
        ? to.meta.authenticatedRedirect
        : to.name === "n8d-login"
        ? "/n8d"
        : "/loss-tem-size";
    return redirect;
  }

  // Kiểm tra quyền Scan In KTP (gọi API)
  if (to.meta.requiresKTPScanInPermission) {
    try {
      const apiBase = (import.meta.env.VITE_API_URL || "/api/v1").replace(/\/$/, "");
      const result = await fetch(`${apiBase}/warehouse-ktp/scan-in/access`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (result.status === 401) {
        clearSession();
        ElMessage.error("Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.");
        return {
          name: "warehouse-ktp-scan-in-login",
          query: { redirect: to.fullPath },
        };
      }

      if (result.status === 403) {
        clearSession();
        ElMessage.error("Bạn không có quyền sử dụng chức năng Scan In KTP.");
        return {
          name: "warehouse-ktp-scan-in-login",
          query: {
            redirect: to.fullPath,
            reason: "no-permission",
          },
        };
      }

      if (!result.ok) {
        ElMessage.error("Không thể kiểm tra quyền Scan In KTP. Vui lòng thử lại.");
        return { name: "home" };
      }
    } catch {
      ElMessage.error("Không kết nối được máy chủ để kiểm tra quyền Scan In KTP.");
      return { name: "home" };
    }
  }
  if (to.meta.requiresKTPScanOutPermission) {
    try {
      const apiBase = (import.meta.env.VITE_API_URL || "/api/v1").replace(/\/$/, "");
      const result = await fetch(`${apiBase}/warehouse-ktp/scan-out/access`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (result.status === 401) {
        clearSession();
        ElMessage.error("Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.");
        return {
          name: "warehouse-ktp-scan-out-login",
          query: { redirect: to.fullPath },
        };
      }
      if (result.status === 403) {
        ElMessage.error("Bạn không có quyền sử dụng chức năng Scan Out KTP.");
        return { name: "warehouse-ktp" };
      }
      if (!result.ok) {
        ElMessage.error("Không thể kiểm tra quyền Scan Out KTP. Vui lòng thử lại.");
        return { name: "warehouse-ktp" };
      }
    } catch {
      ElMessage.error("Không kết nối được máy chủ để kiểm tra quyền Scan Out KTP.");
      return { name: "warehouse-ktp" };
    }
  }
  if (to.meta.requiresKTPScanShipPermission) {
    try {
      const apiBase = (import.meta.env.VITE_API_URL || "/api/v1").replace(/\/$/, "");
      const result = await fetch(`${apiBase}/warehouse-ktp/scan-ship/access`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (result.status === 401) {
        clearSession();
        ElMessage.error("Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.");
        return {
          name: "warehouse-ktp-scan-ship-login",
          query: { redirect: to.fullPath },
        };
      }
      if (result.status === 403) {
        ElMessage.error("Bạn không có quyền sử dụng chức năng Scan Ship KTP.");
        return { name: "warehouse-ktp" };
      }
      if (!result.ok) {
        ElMessage.error("Không thể kiểm tra quyền Scan Ship KTP. Vui lòng thử lại.");
        return { name: "warehouse-ktp" };
      }
    } catch {
      ElMessage.error("Không kết nối được máy chủ để kiểm tra quyền Scan Ship KTP.");
      return { name: "warehouse-ktp" };
    }
  }
  return true;
});

export default router;
