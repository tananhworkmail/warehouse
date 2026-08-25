// src/plugins/axios.js
import axios from "axios";
import { clearSession, getToken } from "@/utils/appStoreAndSecurity";

// 🧩 Sử dụng đúng biến .env
// API modules use VITE_API_URL as their path prefix. Keep axios rooted at the
// current origin unless a deployment explicitly provides a host base URL.
axios.defaults.baseURL = import.meta.env.VITE_BASE_URL_API || "";
axios.defaults.withCredentials = false; // không dùng cookie

// 🧠 Request Interceptor - tự động thêm token vào header
axios.interceptors.request.use(
  (config) => {
    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// 🧩 Response Interceptor - xử lý lỗi chung
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    // Kiểm tra lỗi 401 (token hết hạn)
    if (error.response && error.response.status === 401) {
      console.warn("⚠️ Unauthorized - JWT expired or invalid");
      clearSession();

      const currentPath = `${window.location.pathname}${window.location.search}${window.location.hash}`;
      const pathname = window.location.pathname;
      const loginPath = pathname.startsWith("/warehouse-ktp/scan-in")
        ? "/warehouse-ktp/scan-in/login"
        : pathname.startsWith("/warehouse-ktp/scan-out")
        ? "/warehouse-ktp/scan-out/login"
        : pathname.startsWith("/warehouse-ktp/scan-ship")
        ? "/warehouse-ktp/scan-ship/login"
        : pathname.startsWith("/n8d")
        ? "/n8d/login"
        : "/login";
      const isLoginPage = window.location.pathname === loginPath;

      if (!isLoginPage) {
        window.location.assign(
          `${loginPath}?redirect=${encodeURIComponent(currentPath)}`,
        );
      }
    }
    return Promise.reject(error);
  },
);

export default axios;
