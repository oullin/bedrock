import { createApp, h, ref, onMounted, onUnmounted, computed, defineAsyncComponent } from "vue";
import App from "./App.vue";
import "./styles.css";

const app = createApp(App);
app.mount("#app");
