import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useServerStore = defineStore('server', () => {
  const servers = ref<Array<any>>([]);
  const selectedServer = ref<any>(null);

  function setServers(items: Array<any>) {
    servers.value = items;
  }

  function selectServer(server: any) {
    selectedServer.value = server;
  }

  return { servers, selectedServer, setServers, selectServer };
});
