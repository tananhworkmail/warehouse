<template>
  <el-dialog v-model="dialogVisible" :title="title" :fullscreen="true">
    <template #header>
      <div class="dialog-header">
        <img src="@/assets/Logo.png" alt="Logo" class="logo" />
        <div class="custom-title">{{ title }}</div>
      </div>
    </template>

    <el-tabs v-model="tabsActiveName">
      <el-tab-pane label="NCU" name="NCU">
        <SpreadsheetViewer
          v-if="tabsActiveName === 'NCU' && fileBlobUrl"
          :fileURL="fileBlobUrl"
        />
      </el-tab-pane>

      <el-tab-pane label="Upload" name="upload">
        <el-card shadow="hover" class="mb-3">
          <span style="font-size: 16px; font-weight: bold">Chọn File:</span>

          <el-upload
            v-loading="loading"
            drag
            accept="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
            :http-request="handleFileUpload"
            :limit="1"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">Choose File <em>Here</em></div>
          </el-upload>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onUnmounted } from "vue";
import { ElMessage } from "element-plus";
import SpreadsheetViewer from "../../SpreadsheetViewer/index.vue";

const API_URL = import.meta.env.VITE_API_URL;

/* ===== PROPS / EMIT ===== */
const props = defineProps({
  ncuDialogVisible: {
    type: Boolean,
    default: false,
    required: true,
  },
});
const emit = defineEmits(["close"]);

const dialogVisible = computed({
  get() {
    return props.ncuDialogVisible;
  },
  set(value) {
    emit("close", value);
  },
});

const title = computed(() => "Nhà cung ứng");

/* ===== FIX CỨNG ===== */
const buildingNo = "NCU";
const FILE_NAME = "NCU.xlsx";

const tabsActiveName = ref("NCU");
const loading = ref(false);

/* ===== FILE VIEW ===== */
const fileBlobUrl = ref(null);

const getScheduleFileURL = (bd) => {
  return `${API_URL}/warehouse/schedule?buildingNo=${bd}`;
};

// load file + đổi tên
const loadScheduleFile = async () => {
  try {
    loading.value = true;

    // cleanup blob cũ
    if (fileBlobUrl.value) {
      URL.revokeObjectURL(fileBlobUrl.value);
      fileBlobUrl.value = null;
    }

    const res = await fetch(getScheduleFileURL(buildingNo));
    const blob = await res.blob();

    // ÉP TÊN FILE
    const file = new File([blob], FILE_NAME, {
      type: blob.type,
    });

    fileBlobUrl.value = URL.createObjectURL(file);
  } catch (err) {
    ElMessage.error("Không tải được file NCU");
  } finally {
    loading.value = false;
  }
};

// khi mở tab KVT thì load file
watch(
  () => dialogVisible.value,
  (open) => {
    if (open) {
      loadScheduleFile();
    }
  },
  { immediate: true },
);

// cleanup khi đóng component
onUnmounted(() => {
  if (fileBlobUrl.value) {
    URL.revokeObjectURL(fileBlobUrl.value);
  }
});

/* ===== UPLOAD ===== */
const handleFileUpload = async (uploadInfo) => {
  const { file, onSuccess, onError } = uploadInfo;
  const formData = new FormData();
  formData.append("buildingNo", buildingNo);
  formData.append("file", file);

  loading.value = true;
  try {
    const res = await fetch(`${API_URL}/warehouse/schedule`, {
      method: "POST",
      body: formData,
    });
    const data = await res.json();

    if (data.code === 200) {
      ElMessage.success("Upload thành công");
      await loadScheduleFile();
      onSuccess(data);
    } else {
      ElMessage.error("Upload thất bại");
      onError(new Error("Upload failed"));
    }
  } catch (error) {
    ElMessage.error("Upload lỗi");
    onError(error);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>

.custom-title {
  text-align: center;
  width: 100%;
  font-size: 18px;
  font-weight: bold;
  color: black;
}
.mb-3 {
  margin-bottom: 1rem;
}
.building-form {
  display: flex;
  align-items: center;
  gap: 16px;
}

.building-label {
  font-weight: bold;
  font-size: 16px;
}

.building-select {
  width: 220px;
  font-size: 14px;
}

.confirm-btn {
  width: 120px; /* cố định để button nhìn gọn */
  font-weight: bold;
}
.logo {
  width: 80px;
  height: auto;
}
.dialog-header {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}
</style>
