<template>
  <div id="app" :class="{ 'app-layout': showLayout }">
    <!-- Sidebar -->
    <aside v-if="showLayout" class="sidebar" :class="{ 'sidebar--open': sidebarOpen }">
      <div class="sidebar-header">
        <span class="material-symbols-outlined sidebar-logo-icon">restaurant</span>
        <div class="sidebar-brand">
          <span class="sidebar-brand-name">StockPro</span>
          <span class="sidebar-brand-sub">Restaurant Manager</span>
        </div>
      </div>

      <nav class="sidebar-nav">
        <div class="sidebar-nav-section">
          <span class="sidebar-nav-label">เมนูหลัก</span>
          <router-link to="/" class="sidebar-link" @click="closeSidebar">
            <span class="material-symbols-outlined">dashboard</span>
            <span>แดชบอร์ด</span>
          </router-link>
          <router-link to="/ingredients" class="sidebar-link" @click="closeSidebar">
            <span class="material-symbols-outlined">inventory_2</span>
            <span>วัตถุดิบ</span>
          </router-link>
          <router-link to="/categories" class="sidebar-link" @click="closeSidebar">
            <span class="material-symbols-outlined">category</span>
            <span>หมวดหมู่</span>
          </router-link>
          <router-link to="/units" class="sidebar-link" @click="closeSidebar">
            <span class="material-symbols-outlined">straighten</span>
            <span>หน่วย</span>
          </router-link>
        </div>
        <div class="sidebar-nav-section">
          <span class="sidebar-nav-label">รายการ</span>
          <router-link to="/stock" class="sidebar-link" @click="closeSidebar">
            <span class="material-symbols-outlined">swap_vert</span>
            <span>รับ/จ่ายสต็อค</span>
          </router-link>
          <router-link to="/history" class="sidebar-link" @click="closeSidebar">
            <span class="material-symbols-outlined">history</span>
            <span>ประวัติ</span>
          </router-link>
        </div>
        <div class="sidebar-nav-section" v-if="isAdmin">
          <span class="sidebar-nav-label">ผู้ดูแลระบบ</span>
          <router-link to="/users" class="sidebar-link" @click="closeSidebar">
            <span class="material-symbols-outlined">group</span>
            <span>จัดการผู้ใช้</span>
          </router-link>
        </div>
      </nav>

      <div class="sidebar-footer">
        <div class="sidebar-user-info">
          <div class="sidebar-avatar">
            <span class="material-symbols-outlined">{{ isAdmin ? 'shield_person' : 'person' }}</span>
          </div>
          <div class="sidebar-user-detail">
            <span class="sidebar-user-name">{{ user?.full_name || user?.username }}</span>
            <span class="sidebar-user-role">{{ isAdmin ? 'ผู้ดูแลระบบ' : 'พนักงาน' }}</span>
          </div>
        </div>
        <button class="sidebar-logout-btn" @click="showLogoutConfirm = true" title="ออกจากระบบ">
          <span class="material-symbols-outlined">logout</span>
        </button>
      </div>
    </aside>

    <!-- Mobile overlay -->
    <div v-if="sidebarOpen" class="sidebar-overlay" @click="closeSidebar"></div>

    <!-- Main content area -->
    <div :class="{ 'main-wrapper': showLayout }">
      <header v-if="showLayout" class="mobile-topbar">
        <button class="mobile-topbar-btn" @click="toggleSidebar">
          <span class="material-symbols-outlined">{{ sidebarOpen ? 'close' : 'menu' }}</span>
        </button>
        <span class="mobile-topbar-title">StockPro</span>
      </header>

      <main :class="showLayout ? 'content' : ''">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <ConfirmModal
      v-if="showLogoutConfirm"
      title="ออกจากระบบ"
      message="คุณต้องการออกจากระบบหรือไม่?"
      confirmText="ออกจากระบบ"
      variant="warning"
      @confirm="handleLogout"
      @cancel="showLogoutConfirm = false"
    />
  </div>
</template>

<script>
import { useAuth } from './composables/useAuth.js'
import ConfirmModal from './components/ConfirmModal.vue'

export default {
  components: { ConfirmModal },
  data() {
    return {
      authState: null,
      showLogoutConfirm: false,
      sidebarOpen: false
    }
  },
  created() {
    const { state } = useAuth()
    this.authState = state
  },
  computed: {
    user() {
      return this.authState?.user
    },
    isLoggedIn() {
      return !!this.authState?.user
    },
    isHomePage() {
      return this.$route?.name === 'Home'
    },
    showLayout() {
      return this.isLoggedIn && !this.isHomePage
    },
    isAdmin() {
      return this.authState?.user?.role === 'admin'
    }
  },
  methods: {
    toggleSidebar() {
      this.sidebarOpen = !this.sidebarOpen
    },
    closeSidebar() {
      this.sidebarOpen = false
    },
    async handleLogout() {
      this.showLogoutConfirm = false
      const { logout } = useAuth()
      await logout()
      this.$router.push('/home')
    }
  }
}
</script>
