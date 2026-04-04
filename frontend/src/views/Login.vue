<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <span class="login-icon">🍳</span>
        <h1 class="login-title">ระบบสต็อคร้านอาหาร</h1>
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
          />
        </div>
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <button class="btn btn-primary btn-login" type="submit" :disabled="loading">
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

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  padding: var(--space-md);
}

.login-card {
  background: var(--color-white);
  border-radius: var(--radius-xl);
  padding: var(--space-2xl);
  width: 100%;
  max-width: 400px;
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--color-border);
}

.login-header {
  text-align: center;
  margin-bottom: var(--space-xl);
}

.login-icon {
  font-size: 3rem;
  display: block;
  margin-bottom: var(--space-sm);
}

.login-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text);
  margin-bottom: var(--space-xs);
}

.login-subtitle {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.btn-login {
  width: 100%;
  padding: 0.75rem;
  font-size: var(--font-size-base);
  margin-top: var(--space-sm);
}

.btn-login:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.alert {
  margin-bottom: var(--space-sm);
}
</style>
