<template>
  <el-dialog  v-model="dialogVisible" :fullscreen="true" class="image-dialog">
    <!-- Header với logo -->
    <template #header>
      <div class="dialog-header">
        <img src="@/assets/Logo.png" alt="Logo" class="logo" />
        <el-tabs v-model="activeTab" class="tabs">
          <el-tab-pane label="KVT" name="KVT"></el-tab-pane>
          <el-tab-pane label="3rd Floor" name="3rd Floor"></el-tab-pane>
          <el-tab-pane label="4th Floor" name="4th Floor"></el-tab-pane>
        </el-tabs>
      </div>
    </template>

    <!-- Nội dung ảnh theo tab -->
    <div class="image-wrapper">
      <img :src="currentImage" class="responsive-image" />
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from "vue";

// Props & Emit
const props = defineProps({
  organizationalDialogVisible: Boolean,
});
const emit = defineEmits(["close"]);

const dialogVisible = computed({
  get: () => props.organizationalDialogVisible,
  set: (v) => emit("close", v),
});

// Tab active
const activeTab = ref("KVT");

// Import ảnh từ assets
import tab1Img from "@/assets/So-do-to-chuc-kvt.png";
import tab2Img from "@/assets/sodol3.png";
import tab3Img from "@/assets/sodol4.png";

// Lấy ảnh theo tab
const currentImage = computed(() => {
  if (activeTab.value === "KVT") return tab1Img;
  if (activeTab.value === "3rd Floor") return tab2Img;
  if (activeTab.value === "4th Floor") return tab3Img;
  return tab1Img;
});
</script>

<style scoped>
.image-wrapper {
  width: 100%;
  height: calc(100vh - 100px); /* điều chỉnh header + tab */
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

.responsive-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  padding-left: 20px;
}

.logo {
  width: 50px;
  height: auto;
  margin-right: 20px;
}

.tabs {
  flex: 1;
}
</style>
