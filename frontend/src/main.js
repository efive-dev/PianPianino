import { createApp } from "vue";
import App from "./App.vue";
import "./style.css";
import naive from "naive-ui";
import router from "../router";

const app = createApp(App);
app.use(naive);
app.use(router)
app.mount("#app");
