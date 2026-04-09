<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <div class="login-icon">
          <span class="material-symbols-outlined">restaurant</span>
        </div>
        <h1 class="login-title">StockPro</h1>
        <p class="login-subtitle">เข้าสู่ระบบเพื่อจัดการสต็อค</p>
      </div>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>ชื่อผู้ใช้</label>
          <input
            v-model="username"
            class="form-control"
            placeholder="username"
            autocomplete="username"
            maxlength="15"
            autofocus
          />
        </div>
        <div class="form-group">
          <label>รหัสผ่าน</label>
          <input
            v-model="password"
            class="form-control"
            type="password"
            placeholder="••••••••"
            autocomplete="current-password"
            maxlength="15"
          />
        </div>
        <div v-if="error" class="alert alert-error">
          <span class="material-symbols-outlined">error</span>
          {{ error }}
        </div>
        <button class="btn btn-primary btn-login" type="submit" :disabled="loading">
          <span v-if="loading" class="material-symbols-outlined">progress_activity</span>
          {{ loading ? 'กำลังเข้าสู่ระบบ...' : 'เข้าสู่ระบบ' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script>
import { useAuth } from '../composables/useAuth.js'

export default {
  data() {
    return {
      username: '',
      password: '',
      error: '',
      loading: false
    }
  },
  methods: {
    async handleLogin() {
      this.error = ''
      if (!this.username || !this.password) {
        this.error = 'กรุณาระบุชื่อผู้ใช้และรหัสผ่าน'
        return
      }
      this.loading = true
      try {
        const { login } = useAuth()
        await login(this.username, this.password)
        this.$router.push('/')
      } catch (err) {
        this.error = err.response?.data?.error || 'เข้าสู่ระบบไม่สำเร็จ'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
