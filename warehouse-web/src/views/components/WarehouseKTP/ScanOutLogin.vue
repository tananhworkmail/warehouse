<template>
  <KTPScanLogin
    mode="out"
    title="Đăng nhập Scan Out"
    badge="TÁI CHẾ SB 2 · KIỂM HÀNG SB 4"
    submit-label="Vào Scan Out"
    default-redirect="/warehouse-ktp/scan-out"
    permission-message="Tài khoản không có quyền Scan Out KTP."
  />
</template>

<script setup>
import KTPScanLogin from "./KTPScanLogin.vue";
</script>
