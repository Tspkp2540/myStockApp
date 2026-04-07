import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true
})

// Redirect to login on 401
api.interceptors.response.use(
  res => res,
  err => {
    if (err.response && err.response.status === 401 && !err.config.url.includes('/auth/')) {
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

// Auth
export const login = (data) => api.post('/auth/login', data)
export const logout = () => api.post('/auth/logout')
export const getMe = () => api.get('/auth/me')

// Users (admin)
export const getUsers = () => api.get('/users')
export const createUser = (data) => api.post('/users', data)
export const updateUser = (id, data) => api.put(`/users/${id}`, data)
export const deleteUser = (id) => api.delete(`/users/${id}`)

// Categories
export const getCategories = () => api.get('/categories')
export const createCategory = (data) => api.post('/categories', data)
export const updateCategory = (id, data) => api.put(`/categories/${id}`, data)
export const deleteCategory = (id) => api.delete(`/categories/${id}`)
export const bulkDeleteCategories = (ids) => api.post('/categories/bulk-delete', { ids })

// Units
export const getUnits = () => api.get('/units')
export const createUnit = (data) => api.post('/units', data)
export const updateUnit = (id, data) => api.put(`/units/${id}`, data)
export const deleteUnit = (id) => api.delete(`/units/${id}`)
export const bulkDeleteUnits = (ids) => api.post('/units/bulk-delete', { ids })

// Ingredients
export const getIngredients = (params) => api.get('/ingredients', { params })
export const createIngredient = (data) => api.post('/ingredients', data)
export const updateIngredient = (id, data) => api.put(`/ingredients/${id}`, data)
export const deleteIngredient = (id) => api.delete(`/ingredients/${id}`)
export const bulkDeleteIngredients = (ids) => api.post('/ingredients/bulk-delete', { ids })

// Transactions
export const getTransactions = (params) => api.get('/transactions', { params })
export const createTransaction = (data) => api.post('/transactions', data)
export const getCostSummary = () => api.get('/transactions/cost-summary')

// Dashboard
export const getDashboard = () => api.get('/dashboard')
