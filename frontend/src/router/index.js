import { createRouter, createWebHistory } from 'vue-router';
import LoginView from '../components/LoginView.vue';
import DashboardView from '../components/DashboardView.vue';
const routes = [
    { path: '/', component: DashboardView, meta: { requiresAuth: true } },
    { path: '/login', component: LoginView },
];
const router = createRouter({
    history: createWebHistory(),
    routes,
});
export default router;
//# sourceMappingURL=index.js.map