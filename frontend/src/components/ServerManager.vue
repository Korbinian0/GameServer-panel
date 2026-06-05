<template>
  <div class="server-manager">
    <h2>Server-Verwaltung</h2>
    <div class="toolbar">
      <button @click="reload">Aktualisieren</button>
    </div>
    <div v-if="connectionError" class="notice error">{{ connectionError }}</div>
    <div class="list">
      <div v-for="server in servers" :key="server.id" class="server-card">
        <h3>{{ server.hostname }} ({{ server.platform }})</h3>
        <p>ID: {{ server.id }}</p>
        <p>Version: {{ server.version }}</p>
        <p>Letzte Verbindung: {{ server.lastSeen || 'unbekannt' }}</p>
        <div class="actions">
          <button @click="select(server)">Details</button>
        </div>
      </div>
      <div v-if="servers.length === 0" class="empty-state">
        Keine Knoten verfügbar. Bitte aktualisiere die Liste.
      </div>
    </div>
    <section class="events">
      <h3>Live-Ereignisse</h3>
      <ul>
        <li v-for="(event, index) in events" :key="index">{{ event }}</li>
      </ul>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useServerStore } from '../stores/server';
import { getNodes } from '../api/http';
import { useAuthStore } from '../stores/auth';

const store = useServerStore();
const auth = useAuthStore();
const events = ref<string[]>([]);
const connectionError = ref('');
let socket: WebSocket | null = null;

async function reload() {
  try {
    const nodes = await getNodes();
    store.setServers(nodes);
  } catch (err) {
    connectionError.value = 'Knoten konnten nicht geladen werden. Bitte prüfe die Authentifizierung.';
  }
}

function select(server: any) {
  store.selectServer(server);
}

function connectWebsocket() {
  if (!auth.token) {
    return;
  }
  const scheme = window.location.protocol === 'https:' ? 'wss' : 'ws';
  const wsUrl = `${scheme}://${window.location.host}/ws/events?token=${encodeURIComponent(auth.token)}`;
  socket = new WebSocket(wsUrl);

  socket.onmessage = (event) => {
    events.value.unshift(event.data);
    if (events.value.length > 20) {
      events.value.pop();
    }
  };
  socket.onerror = () => {
    connectionError.value = 'WebSocket-Verbindung fehlgeschlagen.';
  };
  socket.onclose = () => {
    connectionError.value = 'WebSocket geschlossen.';
  };
}

onMounted(() => {
  reload();
  connectWebsocket();
});

const servers = store.servers;
</script>

<style scoped>
.server-manager { background: #0f172a; padding: 1.5rem; border-radius: 1rem; }
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; }
button { padding: 0.75rem 1rem; border-radius: 0.75rem; background: #2563eb; color: white; border: none; cursor: pointer; }
.list { display: grid; gap: 1rem; margin-bottom: 1rem; }
.server-card { padding: 1rem; background: #1e293b; border-radius: 1rem; }
.actions { margin-top: 0.75rem; }
.events { padding: 1rem; background: #111827; border-radius: 1rem; }
.empty-state { color: #94a3b8; padding: 1rem; }
.notice.error { margin-bottom: 1rem; color: #f87171; }
</style>

<style scoped>
.server-manager { background: #0f172a; padding: 1.5rem; border-radius: 1rem; }
.toolbar { display: flex; justify-content: flex-end; margin-bottom: 1rem; }
button { padding: 0.75rem 1rem; border-radius: 0.75rem; background: #2563eb; color: white; border: none; }
.list { display: grid; gap: 1rem; }
.server-card { padding: 1rem; background: #1e293b; border-radius: 1rem; }
.actions { margin-top: 0.75rem; }
</style>
