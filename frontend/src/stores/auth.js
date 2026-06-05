import { defineStore } from 'pinia';
import { ref } from 'vue';
import { login, setAuthHeader } from '../api/http';
export const useAuthStore = defineStore('auth', () => {
    const token = ref(localStorage.getItem('auth_token'));
    const user = ref(null);
    if (token.value) {
        setAuthHeader(token.value);
        user.value = { role: 'admin' };
    }
    async function signIn(email, password) {
        const result = await login(email, password);
        token.value = result.token;
        localStorage.setItem('auth_token', result.token);
        setAuthHeader(result.token);
        user.value = { role: 'admin' };
    }
    function signOut() {
        token.value = null;
        user.value = null;
        localStorage.removeItem('auth_token');
        setAuthHeader('');
    }
    return { token, user, signIn, signOut };
});
//# sourceMappingURL=auth.js.map