import axios from "@/plugins/axios";
import { clearSession, setSession } from "@/utils/appStoreAndSecurity";

const API_PREFIX = (import.meta.env.VITE_API_URL || "/api/v1").replace(
  /\/+$/,
  "",
);
const API_URL = `${API_PREFIX}/auth/login`;

const LOGIN_ERRORS = {
  400: {
    code: "invalid_request",
    message:
      "Yêu cầu đăng nhập không hợp lệ. Vui lòng kiểm tra thông tin đã nhập.",
  },
  401: {
    code: "invalid_credentials",
    message: "Mã nhân viên hoặc mật khẩu không đúng.",
  },
  404: {
    code: "login_endpoint_not_found",
    message: "Không tìm thấy dịch vụ đăng nhập trên máy chủ.",
  },
};

const resolveLoginError = (status, serverMessage = "") => {
  const knownError = LOGIN_ERRORS[status];
  if (knownError) return { status, ...knownError };

  const normalizedMessage = String(serverMessage)
    .trim()
    .toLocaleLowerCase("vi");
  if (
    normalizedMessage.includes("account does not exist") ||
    normalizedMessage.includes("sai mật khẩu") ||
    normalizedMessage.includes("incorrect username or password")
  ) {
    return { status: status || 401, ...LOGIN_ERRORS[401] };
  }

  if (normalizedMessage.includes("invalid form")) {
    return { status: status || 400, ...LOGIN_ERRORS[400] };
  }

  return {
    status,
    code: "server_error",
    message: serverMessage || "Máy chủ đang gặp sự cố. Vui lòng thử lại sau.",
  };
};

export const useLogin = async (credentials) => {
  const userId = String(credentials?.USERID || "").trim();
  const password = String(credentials?.PWD || "");
  const branch = String(credentials?.GSBH || "")
    .trim()
    .toUpperCase();

  try {
    const { data } = await axios.post(API_URL, {
      GSBH: branch,
      USERID: userId,
      PWD: password,
    });

    // N8d uses CommonResponse ({ code, data, message }). Accept the direct
    // token shape too so the client remains compatible with the legacy login.
    const tokenValue =
      (typeof data?.data === "string" && data.data) ||
      (typeof data?.data?.token === "string" && data.data.token) ||
      (typeof data?.token === "string" && data.token) ||
      "";
    const token = tokenValue.trim();

    if (!token) {
      clearSession();
      const responseStatus =
        Number(data?.code) >= 400 ? Number(data.code) : undefined;
      return {
        success: false,
        ...resolveLoginError(responseStatus, data?.message || data?.error),
      };
    }

    const authenticatedUserId = String(
      data?.data?.USERID || data?.USERID || userId,
    ).trim();
    const username = String(data?.data?.USERNAME || data?.USERNAME || "");
    const role = String(data?.data?.role || data?.role || "");

    setSession({
      token,
      userId: authenticatedUserId,
      username,
      role,
      branch,
      sn8: data?.data?.SN8 || data?.SN8 || "",
    });

    return {
      success: true,
      user: {
        USERID: authenticatedUserId,
        USERNAME: username,
        role,
        GSBH: branch,
      },
    };
  } catch (error) {
    console.error("Login error:", error);
    clearSession();

    const status = error.response?.status;
    const serverMessage =
      error.response?.data?.message || error.response?.data?.error;

    if (!error.response) {
      return {
        success: false,
        code: "network_error",
        message:
          "Không thể kết nối tới máy chủ. Vui lòng kiểm tra mạng hoặc thử lại.",
      };
    }

    return {
      success: false,
      ...resolveLoginError(status, serverMessage),
    };
  }
};
