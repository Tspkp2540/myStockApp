import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true
})

api.interceptors.response.use(
  res => res,
  err => {
    if (err.response && err.response.status === 401 && !err.config.url.includes('/auth/')) {
      window.location.href = '/front-office/login'
    }
    return Promise.reject(err)
  }
)

// Auth
export const login = (data) => api.post('/auth/login', data)
export const logout = () => api.post('/auth/logout')
export const getMe = () => api.get('/auth/me')

// Front-Office Dashboard
export const getFrontofficeDashboard = () => api.get('/frontoffice/dashboard')

// Menu Items (read-only — active only)
export const getMenuItems = () => api.get('/menu-items', { params: { active: 'true' } })

// Menu Categories
export const getMenuCategories = () => api.get('/menu-categories')

// Sales
export const createSale = (data) => api.post('/frontoffice/sales', data)
export const getSales = (params) => api.get('/frontoffice/sales', { params })
export const deleteSale = (id, reason) => api.delete(`/frontoffice/sales/${id}`, { data: { reason } })

// Export Sales (admin only)
export const exportSalesExcel = (params) => api.get('/frontoffice/sales/export', { params, responseType: 'blob' })
