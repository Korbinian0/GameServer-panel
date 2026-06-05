<template>
  <div class="login-shell">
    <h1>Gateway Login</h1>
    <form @submit.prevent="submit">
      <label>Email</label>
      <input v-model="email" type="email" required />
      <label>Passwort</label>
      <input v-model="password" type="password" required />
      <button type="submit">Anmelden</button>
    </form>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const email = ref('');
const password = ref('');
const error = ref('');
const router = useRouter();
const auth = useAuthStore();

async function submit() {
  try {
    await auth.signIn(email.value, password.value);
    router.push('/');
  } catch (err) {
    error.value = 'Login fehlgeschlagen';
  }
}
</script>

<style scoped>
.login-shell {
  max-width: 420px;
  margin: 10vh auto;
  background: rgba(15, 23, 42, 0.9);
  padding: 2rem;
  border-radius: 16px;
  box-shadow: 0 24px 80px rgba(15, 23, 42, 0.55);
}
label { display: block; margin-top: 1rem; color: #cbd5e1; }
input { width: 100%; padding: 0.75rem; margin-top: 0.5rem; border-radius: 0.75rem; border: 1px solid #334155; background: #0f172a; color: #f8fafc; }
button { margin-top: 1.5rem; width: 100%; padding: 0.9rem; border: none; border-radius: 0.75rem; background: #3b82f6; color: #fff; cursor: pointer; }
.error { margin-top: 1rem; color: #f87171; }
</style>
