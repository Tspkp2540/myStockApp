import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from './composables/useAuth.js'
import Login from './views/Login.vue'
import Dashboard from './views/Dashboard.vue'
import NewSale from './views/NewSale.vue'
import SalesHistory from './views/SalesHistory.vue'

const routes = [
  { path: '/login', name: 'Login', component: Login, meta: { public: true } },
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/new-sale', name: 'NewSale', component: NewSale },
  { path: '/history', name: 'SalesHistory', component: SalesHistory },
]

const router = createRouter({
  history: createWebHistory('/front-office/'),
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

  next()
})

export default router
