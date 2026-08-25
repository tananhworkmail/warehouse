<template>
  <KTPScanLogin
    mode="ship"
    title="Đăng nhập Scan Ship"
    badge="XUẤT HÀNG · KTP"
    submit-label="Vào Scan Ship"
    default-redirect="/warehouse-ktp/scan-ship"
    permission-message="Tài khoản không có quyền Scan Ship KTP."
  />
</template>

<script setup>
import KTPScanLogin from "./KTPScanLogin.vue";
</script>
