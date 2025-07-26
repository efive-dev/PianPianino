import { createWebHistory, createRouter } from "vue-router";
import HomeView from "../src/view/HomeView.vue";
import RegisterView from "../src/view/RegisterView.vue";
import LoginView from "../src/view/LoginView.vue";
import DashboardView from "../src/view/DashboardView.vue";

const routes = [
  { path: "/", component: HomeView },
  { path: "/register", component: RegisterView },
  { path: "/login", component: LoginView },
  { path: "/dashboard", component: DashboardView },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
