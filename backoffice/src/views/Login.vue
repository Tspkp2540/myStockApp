<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <div class="login-icon">
          <span class="material-symbols-outlined">storefront</span>
        </div>
        <h1 class="login-title">Backoffice</h1>
        <p class="login-subtitle">เข้าสู่ระบบขายอาหาร</p>
      </div>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>ชื่อผู้ใช้</label>
          <input v-model="username" class="form-control" :class="{ 'is-invalid': errors.username }" placeholder="username" autocomplete="username" maxlength="15" autofocus @input="errors.username = ''" />
          <span v-if="errors.username" class="field-error">{{ errors.username }}</span>
        </div>
        <div class="form-group">
          <label>รหัสผ่าน</label>
          <input v-model="password" class="form-control" :class="{ 'is-invalid': errors.password }" type="password" placeholder="••••••••" autocomplete="current-password" maxlength="15" @input="errors.password = ''" />
          <span v-if="errors.password" class="field-error">{{ errors.password }}</span>
        </div>
        <div v-if="error" class="alert alert-error">
          <span class="material-symbols-outlined">error</span>
          {{ error }}
        </div>
        <button class="btn btn-primary btn-login" type="submit" :disabled="loading">
          {{ loading ? 'กำลังเข้าสู่ระบบ...' : 'เข้าสู่ระบบ' }}
        </button>
        <button type="button" class="btn btn-outline btn-login" style="margin-top:8px;" @click="goHome">
          <span class="material-symbols-outlined">home</span>
          กลับหน้าหลัก
        </button>
      </form>
    </div>
  </div>
</template>

<script>
import { useAuth } from '../composables/useAuth.js'

export default {
  data() {
    return { username: '', password: '', error: '', loading: false, errors: { username: '', password: '' } }
  },
  methods: {
    goHome() { window.location.href = '/home' },
    async handleLogin() {
      this.error = ''
      this.errors = { username: '', password: '' }
      let valid = true
      if (!this.username.trim()) { this.errors.username = 'กรุณาระบุชื่อผู้ใช้'; valid = false }
      if (!this.password) { this.errors.password = 'กรุณาระบุรหัสผ่าน'; valid = false }
      if (!valid) return
      this.loading = true
      try {
        const { login } = useAuth()
        const user = await login(this.username, this.password)
        if (user.role !== 'admin') {
          const { logout } = useAuth()
          await logout()
          this.error = 'เฉพาะผู้ดูแลระบบเท่านั้นที่สามารถเข้าใช้งานได้'
          return
        }
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
