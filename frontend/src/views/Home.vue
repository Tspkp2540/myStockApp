<template>
  <div class="home-page">
    <div class="home-header">
      <span class="material-symbols-outlined home-logo">restaurant</span>
      <h1 class="home-title">ระบบจัดการร้านอาหาร</h1>
      <p class="home-subtitle">เลือกระบบที่ต้องการเข้าใช้งาน</p>
    </div>

    <!-- Role Selector -->
    <div class="home-role-selector">
      <button
        :class="['home-role-btn', { active: selectedRole === 'employee' }]"
        @click="selectedRole = 'employee'"
      >
        <span class="material-symbols-outlined">person</span>
        พนักงาน
      </button>
      <button
        :class="['home-role-btn', { active: selectedRole === 'admin' }]"
        @click="selectedRole = 'admin'"
      >
        <span class="material-symbols-outlined">shield_person</span>
        ผู้ดูแลระบบ
      </button>
    </div>

    <!-- System Cards -->
    <div class="home-cards">
      <div class="home-card home-card--stock" @click="goStock">
        <div class="home-card-icon">
          <span class="material-symbols-outlined">inventory_2</span>
        </div>
        <h2 class="home-card-title">Stock Login</h2>
        <p class="home-card-desc">ระบบจัดการสต็อควัตถุดิบ</p>
        <span class="home-card-arrow material-symbols-outlined">arrow_forward</span>
      </div>

      <div v-if="selectedRole === 'admin'" class="home-card home-card--backoffice" @click="goBackOffice">
        <div class="home-card-icon">
          <span class="material-symbols-outlined">admin_panel_settings</span>
        </div>
        <h2 class="home-card-title">Sell Back-Office</h2>
        <p class="home-card-desc">ระบบจัดการหลังร้าน</p>
        <span class="home-card-arrow material-symbols-outlined">arrow_forward</span>
      </div>

      <div class="home-card home-card--frontoffice" @click="goFrontOffice">
        <div class="home-card-icon">
          <span class="material-symbols-outlined">point_of_sale</span>
        </div>
        <h2 class="home-card-title">Sell Front-Office</h2>
        <p class="home-card-desc">ระบบขายหน้าร้าน</p>
        <span class="home-card-arrow material-symbols-outlined">arrow_forward</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      selectedRole: 'employee'
    }
  },
  methods: {
    goStock() {
      this.$router.push('/login')
    },
    goBackOffice() {
      // TODO: link to sell back-office
    },
    goFrontOffice() {
      // TODO: link to sell front-office
    }
  }
}
</script>
