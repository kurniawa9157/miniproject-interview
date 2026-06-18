import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  withCredentials: true,
})

api.interceptors.response.use(
  (res) => res,
  (err) => {
    const url: string = err.config?.url ?? ''
    // Don't redirect on /auth/me — it's the auth probe and returns 401 when
    // logged out, which is expected. Let callers handle it.
    const isAuthProbe = url.includes('/auth/me')
    if (err.response?.status === 401 && !isAuthProbe) {
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

export default api
