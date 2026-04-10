import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from './composables/useAuth.js'
import Login from './views/Login.vue'
import Dashboard from './views/Dashboard.vue'
import MenuItems from './views/MenuItems.vue'
import MenuCategories from './views/MenuCategories.vue'
import SalesHistory from './views/SalesHistory.vue'
import DeletedSales from './views/DeletedSales.vue'

const routes = [
  { path: '/login', name: 'Login', component: Login, meta: { public: true } },
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/menu', name: 'MenuItems', component: MenuItems },
  { path: '/menu-categories', name: 'MenuCategories', component: MenuCategories },
  { path: '/sales', name: 'SalesHistory', component: SalesHistory },
  { path: '/deleted-sales', name: 'DeletedSales', component: DeletedSales },
]

const router = createRouter({
  history: createWebHistory('/backoffice/'),
  routes
})

router.beforeEach(async (to, from, next) => {
  const { state, fetchUser, isLoggedIn } = useAuth()

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

  // Only admin can access backoffice
  if (state.user && state.user.role !== 'admin') {
    const { logout } = useAuth()
    await logout()
    return next('/login')
  }

  next()
})

export default router
