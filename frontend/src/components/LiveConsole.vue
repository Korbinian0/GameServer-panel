<template>
  <div class="console-shell">
    <h2>Live-Konsole</h2>
    <div class="output" ref="outputRef">{{ output }}</div>
    <form @submit.prevent="sendCommand">
      <input v-model="command" placeholder="Befehl eingeben" />
      <button type="submit">Senden</button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const command = ref('');
const output = ref('Konsolen-Stream wird vorbereitet...\n');
const outputRef = ref<HTMLElement | null>(null);

function sendCommand() {
  output.value += `> ${command.value}\n`;
  command.value = '';
}
</script>

<style scoped>
.console-shell { padding: 1.5rem; background: #0f172a; border-radius: 1rem; }
.output { min-height: 200px; padding: 1rem; background: #020617; border-radius: 0.75rem; margin-bottom: 1rem; color: #e2e8f0; font-family: monospace; white-space: pre-wrap; }
input { width: calc(100% - 110px); margin-right: 0.75rem; padding: 0.75rem; border-radius: 0.75rem; background: #1e293b; color: #f8fafc; border: 1px solid #475569; }
button { padding: 0.75rem 1rem; border-radius: 0.75rem; border: none; background: #2563eb; color: white; }
</style>
