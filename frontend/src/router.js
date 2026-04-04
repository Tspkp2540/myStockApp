import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from './composables/useAuth.js'
import Login from './views/Login.vue'
import Dashboard from './views/Dashboard.vue'
import Ingredients from './views/Ingredients.vue'
import Categories from './views/Categories.vue'
import Units from './views/Units.vue'
import StockTransaction from './views/StockTransaction.vue'
import TransactionHistory from './views/TransactionHistory.vue'
import Users from './views/Users.vue'

const routes = [
  { path: '/login', name: 'Login', component: Login, meta: { public: true } },
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/ingredients', name: 'Ingredients', component: Ingredients },
  { path: '/categories', name: 'Categories', component: Categories },
  { path: '/units', name: 'Units', component: Units },
  { path: '/stock', name: 'StockTransaction', component: StockTransaction },
  { path: '/history', name: 'TransactionHistory', component: TransactionHistory },
  { path: '/users', name: 'Users', component: Users, meta: { admin: true } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  const { state, fetchUser, isLoggedIn, isAdmin } = useAuth()

  if (!state.loaded) {
    await fetchUser()
  }

  if (to.meta.public) {
    if (isLoggedIn()) return next('/')
    return next()
  }

  if (!isLoggedIn()) {
    return next('/login')
  }

  if (to.meta.admin && !isAdmin()) {
    return next('/')
  }

  next()
})

export default router
