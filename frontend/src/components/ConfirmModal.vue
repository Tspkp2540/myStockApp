<template>
  <div class="modal-overlay" @click.self="$emit('cancel')">
    <div class="modal confirm-modal">
      <div class="confirm-icon" :class="iconClass">{{ icon }}</div>
      <h3 class="confirm-title">{{ title }}</h3>
      <p class="confirm-message">{{ message }}</p>
      <div class="modal-actions">
        <button class="btn btn-outline" @click="$emit('cancel')">ยกเลิก</button>
        <button :class="['btn', confirmBtnClass]" @click="$emit('confirm')">{{ confirmText }}</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    title: { type: String, default: 'ยืนยันการดำเนินการ' },
    message: { type: String, default: 'คุณแน่ใจหรือไม่?' },
    confirmText: { type: String, default: 'ยืนยัน' },
    variant: { type: String, default: 'danger' } // danger, warning, info
  },
  emits: ['confirm', 'cancel'],
  computed: {
    icon() {
      const icons = { danger: '⚠️', warning: '⚡', info: 'ℹ️' }
      return icons[this.variant] || icons.danger
    },
    iconClass() {
      return `confirm-icon-${this.variant}`
    },
    confirmBtnClass() {
      const classes = { danger: 'btn-danger', warning: 'btn-warning', info: 'btn-primary' }
      return classes[this.variant] || 'btn-danger'
    }
  }
}
</script>
