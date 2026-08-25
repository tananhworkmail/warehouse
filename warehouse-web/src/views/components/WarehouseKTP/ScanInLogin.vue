<template>
  <KTPScanLogin
    mode="in"
    title="Đăng nhập Scan In"
    badge="NHẬP KHO · KTP"
    submit-label="Vào Scan In"
    default-redirect="/warehouse-ktp/scan-in"
    permission-message="Tài khoản không có quyền Scan In KTP."
  />
</template>

<script setup>
import KTPScanLogin from "./KTPScanLogin.vue";
</script>
