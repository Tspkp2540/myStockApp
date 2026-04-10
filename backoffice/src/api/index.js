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
      window.location.href = '/backoffice/login'
    }
    return Promise.reject(err)
  }
)

// Auth
export const login = (data) => api.post('/auth/login', data)
export const logout = () => api.post('/auth/logout')
export const getMe = () => api.get('/auth/me')

// Menu Items
export const getMenuItems = (params) => api.get('/menu-items', { params })
export const getMenuItem = (id) => api.get(`/menu-items/${id}`)
export const createMenuItem = (data) => api.post('/menu-items', data)
export const updateMenuItem = (id, data) => api.put(`/menu-items/${id}`, data)
export const deleteMenuItem = (id) => api.delete(`/menu-items/${id}`)

// Sales
export const createSale = (data) => api.post('/sales', data)
export const getSales = (params) => api.get('/sales', { params })
export const getSale = (id) => api.get(`/sales/${id}`)
export const deleteSale = (id) => api.delete(`/sales/${id}`)

// Backoffice Dashboard
export const getBackofficeDashboard = () => api.get('/backoffice/dashboard')

// Deleted Sales (from front-office)
export const getDeletedSales = (params) => api.get('/sales/deleted', { params })

// Menu Categories
export const getMenuCategories = () => api.get('/menu-categories')
export const createMenuCategory = (data) => api.post('/menu-categories', data)
export const updateMenuCategory = (id, data) => api.put(`/menu-categories/${id}`, data)
export const deleteMenuCategory = (id) => api.delete(`/menu-categories/${id}`)

// Ingredients (for menu composition)
export const getIngredients = (params) => api.get('/ingredients', { params })

// Export Sales
export const exportSalesExcel = (params) => api.get('/sales/export', { params, responseType: 'blob' })
