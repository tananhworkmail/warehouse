import "./assets/main.css";
import { createApp } from "vue";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css"; // Import CSS của Element Plus
import App from "./App.vue";
import router from "./router";
import { createPinia } from "pinia";
import Print from "vue3-print-nb";

import vuetify from "./plugins/vuetify";
// ⚠️ Không dùng plugins/index.js (registerPlugins) nữa — i18n giờ lấy trực tiếp
// từ locales/{vi,en,zh}.js đã gộp n8d vào warehouse.n8d, xem src/hooks/i18n.js
import i18n from "@/hooks/i18n";

const pinia = createPinia();
const app = createApp(App);
app.use(router);

app.use(pinia);
app.use(Print);
app.use(ElementPlus);
app.use(vuetify);
app.use(i18n);
app.mount("#app");