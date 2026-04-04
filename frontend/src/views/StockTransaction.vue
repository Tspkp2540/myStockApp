<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">รับ / จ่าย สต็อค</h1>
      <p class="page-subtitle">บันทึกรายการรับเข้าหรือจ่ายออกวัตถุดิบ</p>
    </div>

    <div class="card">
      <div class="form-row">
        <div class="form-group">
          <label>วัตถุดิบ <span class="required">*</span></label>
          <select v-model="form.ingredient_id" class="form-control">
            <option :value="null">-- เลือกวัตถุดิบ --</option>
            <option v-for="ing in ingredients" :key="ing.id" :value="ing.id">
              {{ ing.name }} (คงเหลือ: {{ ing.quantity }} {{ ing.unit?.name }})
            </option>
          </select>
        </div>
        <div class="form-group">
          <label>ประเภท <span class="required">*</span></label>
          <select v-model="form.type" class="form-control">
            <option value="in">▲ รับเข้า (Stock In)</option>
            <option value="out">▼ จ่ายออก (Stock Out)</option>
          </select>
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label>จำนวน <span class="required">*</span></label>
          <input v-model.number="form.quantity" class="form-control" type="number" min="0.1" step="0.1" />
        </div>
        <div class="form-group">
          <label>หมายเหตุ</label>
          <input v-model="form.note" class="form-control" placeholder="เช่น ซื้อจากตลาด, ใช้ทำอาหาร" />
        </div>
      </div>
      <button class="btn btn-lg btn-success" @click="confirmSubmit">
        บันทึกรายการ
      </button>
      <div v-if="message" :class="['alert', messageType]">
        {{ message }}
      </div>
    </div>

    <div class="card" v-if="selectedIngredient">
      <div class="card-header">
        <h2 class="section-title">📦 ข้อมูลวัตถุดิบที่เลือก</h2>
      </div>
      <div class="info-grid">
        <div class="info-item"><strong>ชื่อ:</strong> {{ selectedIngredient.name }}</div>
        <div class="info-item"><strong>หมวดหมู่:</strong> {{ selectedIngredient.category?.name || '-' }}</div>
        <div class="info-item"><strong>คงเหลือ:</strong> {{ selectedIngredient.quantity }} {{ selectedIngredient.unit?.name }}</div>
        <div class="info-item"><strong>สต็อคขั้นต่ำ:</strong> {{ selectedIngredient.min_stock }} {{ selectedIngredient.unit?.name }}</div>
      </div>
    </div>

    <ConfirmModal
      v-if="confirmAction"
      :title="confirmAction.title"
      :message="confirmAction.message"
      :confirmText="confirmAction.confirmText"
      :variant="confirmAction.variant"
      @confirm="confirmAction.onConfirm"
      @cancel="confirmAction = null"
    />
  </div>
</template>

<script>
import { getIngredients, createTransaction } from '../api/index.js'
import ConfirmModal from '../components/ConfirmModal.vue'

export default {
  components: { ConfirmModal },
  data() {
    return {
      ingredients: [],
      form: { ingredient_id: null, type: 'in', quantity: 0, note: '' },
      message: '',
      messageType: '',
      confirmAction: null
    }
  },
  computed: {
    selectedIngredient() {
      if (!this.form.ingredient_id) return null
      return this.ingredients.find(i => i.id === this.form.ingredient_id)
    }
  },
  async mounted() {
    await this.loadIngredients()
  },
  methods: {
    async loadIngredients() {
      const { data } = await getIngredients()
      this.ingredients = data
    },
    confirmSubmit() {
      if (!this.form.ingredient_id || this.form.quantity <= 0) {
        this.message = 'กรุณาเลือกวัตถุดิบและระบุจำนวน'
        this.messageType = 'alert-error'
        return
      }
      const ing = this.selectedIngredient
      const typeName = this.form.type === 'in' ? 'รับเข้า' : 'จ่ายออก'
      this.confirmAction = {
        title: `ยืนยันการ${typeName}สต็อค`,
        message: `${typeName} "${ing?.name}" จำนวน ${this.form.quantity} ${ing?.unit?.name || ''} ต้องการดำเนินการหรือไม่?`,
        confirmText: 'บันทึก',
        variant: this.form.type === 'in' ? 'info' : 'warning',
        onConfirm: () => this.submit()
      }
    },
    async submit() {
      this.confirmAction = null
      try {
        await createTransaction({
          ingredient_id: Number(this.form.ingredient_id),
          type: this.form.type,
          quantity: this.form.quantity,
          note: this.form.note
        })
        this.message = '✓ บันทึกสำเร็จ!'
        this.messageType = 'alert-success'
        this.form.quantity = 0
        this.form.note = ''
        await this.loadIngredients()
      } catch (err) {
        this.message = err.response?.data?.error || 'เกิดข้อผิดพลาด'
        this.messageType = 'alert-error'
      }
    }
  }
}
</script>

<style scoped>
.alert {
  margin-top: var(--space-md);
}

.btn-lg {
  margin-top: var(--space-sm);
}
</style>
