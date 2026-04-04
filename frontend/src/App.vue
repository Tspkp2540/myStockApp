<template>
  <div id="app">
    <nav class="navbar" v-if="isLoggedIn">
      <div class="nav-brand">🍳 ระบบสต็อคร้านอาหาร</div>
      <div class="nav-links">
        <router-link to="/">แดชบอร์ด</router-link>
        <router-link to="/ingredients">วัตถุดิบ</router-link>
        <router-link to="/categories">หมวดหมู่</router-link>
        <router-link to="/units">หน่วย</router-link>
        <router-link to="/stock">รับ/จ่ายสต็อค</router-link>
        <router-link to="/history">ประวัติ</router-link>
        <router-link to="/users" v-if="isAdmin">จัดการผู้ใช้</router-link>
      </div>
      <div class="nav-user">
        <span class="nav-user-name">
          <span :class="['badge', isAdmin ? 'badge-warning' : 'badge-info']">
            {{ isAdmin ? '🔑' : '👤' }}
          </span>
          {{ user?.full_name || user?.username }}
        </span>
        <button class="btn btn-sm btn-outline nav-logout" @click="showLogoutConfirm = true">ออกจากระบบ</button>
      </div>
    </nav>
    <main :class="isLoggedIn ? 'content' : ''">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

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
      showLogoutConfirm: false
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
    isAdmin() {
      return this.authState?.user?.role === 'admin'
    }
  },
  methods: {
    async handleLogout() {
      this.showLogoutConfirm = false
      const { logout } = useAuth()
      await logout()
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>
.nav-user {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.nav-user-name {
  color: var(--color-nav-link);
  font-size: var(--font-size-sm);
  display: flex;
  align-items: center;
  gap: var(--space-xs);
}

.nav-logout {
  color: var(--color-nav-link);
  border-color: rgba(148, 163, 184, 0.3);
}

.nav-logout:hover {
  color: var(--color-text-inverse);
  border-color: rgba(148, 163, 184, 0.6);
  background: var(--color-nav-hover);
}
</style>
