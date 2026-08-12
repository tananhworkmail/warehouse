import "./assets/main.css";
import { createApp } from "vue";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css"; // Import CSS của Element Plus
import App from "./App.vue";
import router from "./router";
import { createPinia } from "pinia";
import Print from "vue3-print-nb";


const pinia = createPinia();
const app = createApp(App);
app.use(router);

app.use(pinia);
app.use(Print);
app.use(ElementPlus);
app.mount("#app");
