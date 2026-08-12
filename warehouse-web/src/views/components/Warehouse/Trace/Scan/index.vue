<template>
  <el-dialog v-model="dialogVisible" :title="title" :fullscreen="true">
    <template #header>
      <div class="dialog-header">
        <img src="@/assets/Logo.png" alt="Logo" class="logo" />
        <div class="custom-title">{{ title }}</div>
      </div>
    </template>

    <div class="box">
     
      <!-- Input -->
      <el-input
        ref="inputRef"
        v-model="state.input"
        size="large"
        placeholder="Please Input Order No"
        :prefix-icon="Search"
        class="input"
        @keyup.enter="handleSearchClick"
      />

      <!-- Button -->
      <el-button
        :disabled="!state.input"
        size="large"
        color="#E87325"
        class="button"
        @click="handleSearchClick"
      >
        Search
      </el-button>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch, nextTick } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { Search } from "@element-plus/icons-vue";
const API_URL = import.meta.env.VITE_API_URL;
// Props và emit từ cha
const props = defineProps({
  traceDialogVisible: { type: Boolean, required: true },
});
const emit = defineEmits(["close"]);

// computed để bind v-model
const dialogVisible = computed({
  get: () => props.traceDialogVisible,
  set: (val) => emit("close", val),
});

const title = computed(() => "Trace");

const router = useRouter();
const inputRef = ref(null);
const state = reactive({ input: "" });

// Focus input mỗi khi dialog mở
watch(dialogVisible, async (val) => {
  if (val) {
    await nextTick();
    inputRef.value?.focus();
  } else {
    state.input = "";
  }
});

const handleSearchClick = async () => {
  if (!state.input) return;

  try {
    const res = await fetch(
      `${API_URL}/warehouse/check_order?orderNo=${encodeURIComponent(
        state.input
      )}`
    );

    if (!res.ok) {
      throw new Error(`HTTP error! status: ${res.status}`);
    }

    const data = await res.json();

    if (data.data && data.data.DDBH) {
      router.push(`/trace/${state.input}`);
      dialogVisible.value = false; // đóng dialog
    } else {
      ElMessage.error("Order does not exist");
    }
  } catch (err) {
    ElMessage.error("Error fetching order");
    console.error(err);
  }
};
</script>

<style scoped lang="scss">
.custom-title {
  text-align: center;
  width: 100%;
  font-size: 18px;
  font-weight: bold;
  color: black;
}

.box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.logo {
  width: 80px;
  height: auto;
}

.input {
  width: 30vw;
}

.button {
  width: 30vw; /* chỉnh chiều rộng nút */
  border-radius: 8px;
  font-weight: 600;
}
.dialog-header {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}
::v-deep(.el-dialog) {
  user-select: text !important;
}
</style>
