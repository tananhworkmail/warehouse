<template>
  <div ref="viewerRef" class="viewer_container"></div>
</template>

<script>
export default {
  name: "SpreadsheetViewer",
};
</script>

<script setup>
import { onMounted, ref, watch } from "vue";
import { SpreadsheetViewer } from "spreadsheet-viewer";

const props = defineProps({
  fileURL: {
    type: String,
    required: true,
  },
});
const emit = defineEmits(["change"]);

const viewerRef = ref(null);

const init = () => {
  SpreadsheetViewer({
    container: viewerRef.value,
    assetsUrl: "/spreadsheet-viewer/sv/index.html",
  })
    .then((instance) => {
      instance.configure({
        license: "demo",
      });
      instance.loadWorkbook(props.fileURL, 0);
      emit("change", true);
    })
    .catch(() => {
      emit("change", false);
    });
};
watch(
  () => props.fileURL,
  (newUrl) => {
    if (instance && newUrl) {
      instance.loadWorkbook(newUrl, 0);
    }
  }
);

onMounted(() => {
  init();
});
</script>

<style lang="scss" scoped>
.viewer_container {
  height: 80vh;
}
</style>
