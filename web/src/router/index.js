import { createRouter, createWebHistory } from "vue-router";
import LoginView from "../views/LoginView.vue";
import StudentsView from "../views/StudentsView.vue";
import BuildingsView from "../views/BuildingsView.vue";
import RoomsView from "../views/RoomsView.vue";
import PaymentsView from "../views/PaymentsView.vue";
import UsersView from "../views/UsersView.vue";
import StatsView from "../views/StatsView.vue";

const routes = [
  { path: "/", redirect: "/login" },
  { path: "/login", component: LoginView },
  { path: "/students", component: StudentsView },
  { path: "/buildings", component: BuildingsView },
  { path: "/rooms", component: RoomsView },
  { path: "/payments", component: PaymentsView },
  { path: "/users", component: UsersView },
  { path: "/stats", component: StatsView }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem("token");
  if (to.path !== "/login" && !token) {
    next("/login");
  } else {
    next();
  }
});

export default router;
