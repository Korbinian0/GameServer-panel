import { defineStore } from 'pinia';
import { ref } from 'vue';
export const useServerStore = defineStore('server', () => {
    const servers = ref([]);
    const selectedServer = ref(null);
    function setServers(items) {
        servers.value = items;
    }
    function selectServer(server) {
        selectedServer.value = server;
    }
    return { servers, selectedServer, setServers, selectServer };
});
//# sourceMappingURL=server.js.map