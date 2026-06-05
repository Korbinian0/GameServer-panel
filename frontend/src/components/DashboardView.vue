<template>
  <div class="dashboard">
    <header>
      <div class="header-row">
        <h1>Server-Management Dashboard</h1>
        <button @click="logout">Abmelden</button>
      </div>
    </header>
    <section class="cards">
      <article class="card">Backend-Verwaltung</article>
      <article class="card">Serverstatus</article>
      <article class="card">Live-Konsole</article>
      <article class="card">Dateimanager</article>
    </section>
    <section class="panels">
      <ServerManager />
    </section>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import ServerManager from './ServerManager.vue';

const router = useRouter();
const auth = useAuthStore();

function logout() {
  auth.signOut();
  router.push('/login');
}
</script>

<style scoped>
.dashboard { padding: 2rem; }
header h1 { margin-bottom: 1.5rem; }
.cards { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 1rem; }
.card { padding: 1.25rem; background: #1e293b; border-radius: 1rem; color: #e2e8f0; }
.panels { margin-top: 2rem; }
</style>
