import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import { useAuthStore } from './stores/auth';

const app = createApp(App);
const pinia = createPinia();
app.use(pinia);
app.use(router);

router.beforeEach((to, from, next) => {
  const auth = useAuthStore(pinia);
  if (to.meta.requiresAuth && !auth.token) {
    return next('/login');
  }
  if (to.path === '/login' && auth.token) {
    return next('/');
  }
  return next();
});

app.mount('#app');
