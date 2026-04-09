import { reactive } from 'vue'
import { getMe, login as apiLogin, logout as apiLogout } from '../api/index.js'

const state = reactive({
  user: null,
  loaded: false
})

export function useAuth() {
  async function fetchUser() {
    try {
      const { data } = await getMe()
      state.user = data.user
    } catch {
      state.user = null
    }
    state.loaded = true
  }

  async function login(username, password) {
    const { data } = await apiLogin({ username, password })
    state.user = data.user
    return data.user
  }

  async function logout() {
    await apiLogout()
    state.user = null
  }

  function isLoggedIn() {
    return !!state.user
  }

  return { state, fetchUser, login, logout, isLoggedIn }
}
