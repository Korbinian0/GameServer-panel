import axios from 'axios';
const api = axios.create({
    baseURL: '/api',
    headers: {
        'Content-Type': 'application/json',
    },
});
export function setAuthHeader(token) {
    api.defaults.headers.common.Authorization = `Bearer ${token}`;
}
export async function login(email, password) {
    const response = await api.post('/login', { email, password });
    return response.data;
}
export async function getNodes() {
    const response = await api.get('/nodes');
    return response.data;
}
export async function getRoles() {
    const response = await api.get('/roles');
    return response.data;
}
export default api;
//# sourceMappingURL=http.js.map