import { createRouter, createWebHistory } from "vue-router";
import LandingPage from "@/views/components/Landingpage.vue";         // ← MỚI
import Warehouse from "@/views/components/Warehouse/Warehouse.vue";
import WarehouseKVT from "@/views/components/WarehouseKVT/WarehouseKVT.vue";
import WarehouseKTP from "@/views/components/WarehouseKTP/WarehouseKTP.vue";
import Trace from "@/views/components/Warehouse/Trace/Scan/index.vue";
import TraceReport from "@/views/components/Warehouse/Trace/report.vue";
import LossTemSize from "@/views/components/WarehouseKVT/LossTemSize/index.vue";
import LoginPage from "@/views/components/WarehouseKVT/LossTemSize/LoginPage.vue";
import Landingpage from "@/views/components/Landingpage.vue"; 
import Laboratory from "@/views/components/WarehousePTN/Laboratory.vue";
const routes = [
  {
    path: "/",
    name: "home",
    component: Landingpage, // ← đổi từ redirect sang landing page
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
      import(
        "@/views/components/WarehousePTN/GoreTex/GoreTexPortal.vue"
      ),
    children: [
      {
        path: "",
        name: "laboratory-forms",
        component: () =>
          import(
            "@/views/components/WarehousePTN/GoreTex/views/EntryView.vue"
          ),
        meta: { title: "GORE-TEX" },
      },
      {
        path: "portal",
        name: "laboratory-forms-portal",
        component: () =>
          import(
            "@/views/components/WarehousePTN/GoreTex/views/HomeView.vue"
          ),
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
    path: "/loss-tem-size",
    name: "loss-tem-size",
    component: LossTemSize,
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem("token");
  const isLoggedIn = !!token;

  if (to.meta.requiresAuth && !isLoggedIn) {
    next({ name: "login", query: { redirect: to.fullPath } });
  } else if (to.meta.requiresGuest && isLoggedIn) {
    next({ name: "warehouse-kvt" });
  } else {
    next();
  }
});

export default router;
