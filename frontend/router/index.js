import { createWebHistory, createRouter } from "vue-router";
import HomeView from "../src/view/HomeView.vue";
import RegisterView from "../src/view/RegisterView.vue";

const routes = [
  { path: "/", component: HomeView },
  { path: "/register", component: RegisterView },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
