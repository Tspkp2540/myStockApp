<template>
  <div id="app" :class="{ 'app-layout': showLayout }">
    <!-- Sidebar -->
    <aside v-if="showLayout" class="sidebar" :class="{ 'sidebar--open': sidebarOpen }">
      <div class="sidebar-header">
        <span class="material-symbols-outlined sidebar-logo-icon">point_of_sale</span>
        <div>
          <span class="sidebar-brand-name">Front-Office</span>
          <span class="sidebar-brand-sub">ระบบขายหน้าร้าน</span>
        </div>
      </div>

      <nav class="sidebar-nav">
        <div class="sidebar-nav-section">
          <span class="sidebar-nav-label">เมนูหลัก</span>
          <router-link to="/" class="sidebar-link" @click="closeSidebar">
            <span class="material-symbols-outlined">dashboard</span>
            <span>แดชบอร์ด</span>
          </router-link>
          <router-link to="/new-sale" class="sidebar-link" @click="closeSidebar">
            <span class="material-symbols-outlined">add_shopping_cart</span>
            <span>บันทึกการขาย</span>
          </router-link>
          <router-link to="/history" class="sidebar-link" @click="closeSidebar">
            <span class="material-symbols-outlined">receipt_long</span>
            <span>ประวัติการขาย</span>
          </router-link>
        </div>
      </nav>

      <div class="sidebar-footer">
        <div class="sidebar-user-info">
          <div class="sidebar-avatar">
            <span class="material-symbols-outlined">person</span>
          </div>
          <div>
            <span class="sidebar-user-name">{{ user?.full_name || user?.username }}</span>
            <span class="sidebar-user-role">{{ user?.role === 'admin' ? 'ผู้ดูแลระบบ' : 'พนักงาน' }}</span>
          </div>
        </div>
        <button class="sidebar-logout-btn" @click="handleLogout" title="ออกจากระบบ">
          <span class="material-symbols-outlined">logout</span>
        </button>
      </div>
    </aside>

    <div v-if="sidebarOpen" class="sidebar-overlay" @click="closeSidebar"></div>

    <div :class="{ 'main-wrapper': showLayout }">
      <header v-if="showLayout" class="mobile-topbar">
        <button class="mobile-topbar-btn" @click="toggleSidebar">
          <span class="material-symbols-outlined">{{ sidebarOpen ? 'close' : 'menu' }}</span>
        </button>
        <span class="mobile-topbar-title">Front-Office</span>
      </header>
      <main class="main-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script>
import { useAuth } from './composables/useAuth.js'

export default {
  data() {
    return { sidebarOpen: false }
  },
  computed: {
    showLayout() {
      return this.$route.name !== 'Login' && this.user
    },
    user() {
      const { state } = useAuth()
      return state.user
    }
  },
  methods: {
    toggleSidebar() { this.sidebarOpen = !this.sidebarOpen },
    closeSidebar() { this.sidebarOpen = false },
    async handleLogout() {
      const { logout } = useAuth()
      await logout()
      this.$router.push('/login')
    }
  }
}
</script>
