import axios from 'axios'

// In production behind a reverse proxy, set VITE_API_URL="" so requests are
// relative (/api/...) and proxied to the backend on the same origin. The
// `?? 'http://localhost:8080'` only applies when the var is undefined (dev).
const apiBaseURL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

const api = axios.create({
  baseURL: apiBaseURL,
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
