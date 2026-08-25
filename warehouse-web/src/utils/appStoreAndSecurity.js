// appStoreAndSecurity.js
import { defineStore } from "pinia";
import CryptoJS from "crypto-js";

const SECRET_KEY = import.meta.env.VITE_SECRET_KEY;
const TOKEN_KEY = "token";
const STORAGE_KEY = "app_token";

// === 🟩 Pinia Store ===
export const useDashboardStore = defineStore("dashboard", {
  state: () => ({
    activeMenu: localStorage.getItem("activeMenu") || "home",
    isCollapse: localStorage.getItem("navbarCollapse")
      ? localStorage.getItem("navbarCollapse") === "true"
      : true,
    selectedQC: localStorage.getItem("selectedQC") || "",
    arrDepName: JSON.parse(localStorage.getItem("arrDepName") || "[]"),
  }),
  actions: {
    updateActiveMenu(newMenu) {
      this.activeMenu = newMenu;
      localStorage.setItem("activeMenu", newMenu);
    },
    setCollapse(state) {
      this.isCollapse = state;
      localStorage.setItem("navbarCollapse", state);
    },
    handleQCSelect(value) {
      this.selectedQC = value;
      localStorage.setItem("selectedQC", value);
    },
    updateDepartments(departments) {
      this.arrDepName = departments;
      localStorage.setItem("arrDepName", JSON.stringify(departments));
    },
  },
});

// === 🔐 Encryption ===
export function encryptAES(text) {
  const salt = CryptoJS.lib.WordArray.random(8); // 8 bytes
  const keyIv = CryptoJS.EvpKDF(SECRET_KEY, salt, {
    keySize: 256 / 32 + 128 / 32,
    iterations: 1,
  });
  const key = CryptoJS.lib.WordArray.create(keyIv.words.slice(0, 8)); // 32 bytes
  const iv = CryptoJS.lib.WordArray.create(keyIv.words.slice(8, 12)); // 16 bytes

  const encrypted = CryptoJS.AES.encrypt(text, key, { iv });
  const result = CryptoJS.enc.Utf8.parse("Salted__")
    .concat(salt)
    .concat(encrypted.ciphertext);
  return CryptoJS.enc.Base64.stringify(result);
}

// === 🔐 Token Handling ===
export function setToken(token) {
  if (!token) return;
  localStorage.setItem(TOKEN_KEY, token);
}
export function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}

export function removeToken() {
  localStorage.removeItem(STORAGE_KEY);
  localStorage.removeItem(TOKEN_KEY);
}

export function setSession({
  token,
  userId,
  username = "",
  role = "",
  branch = "",
  sn8 = "",
}) {
  setToken(token);
  const normalizedUserId = String(userId || "").trim();
  const normalizedUsername = String(username || "");
  const normalizedRole = String(role || "");
  localStorage.setItem("USERID", normalizedUserId);
  localStorage.setItem("userid", normalizedUserId);
  localStorage.setItem("USERNAME", normalizedUsername);
  localStorage.setItem("username", normalizedUsername);
  localStorage.setItem("ROLE", normalizedRole);
  localStorage.setItem("role", normalizedRole);
  localStorage.setItem("GSBH", branch || "");
  localStorage.setItem("SN8", sn8 || "");
  localStorage.setItem("AUTH_VERSION", `${Date.now()}`);
  window.dispatchEvent(new Event("auth-changed"));
}

export function clearSession() {
  removeToken();
  ["USERID", "USERNAME", "ROLE", "GSBH", "SN8", "userid", "userId", "username", "role"].forEach(
    (key) => {
      localStorage.removeItem(key);
    },
  );
  localStorage.setItem("AUTH_VERSION", `${Date.now()}`);
  window.dispatchEvent(new Event("auth-changed"));
}
