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
          <label>ราคาต่อหน่วย (ฺ) <span v-if="form.type === 'in'" class="required">*</span></label>
          <input v-model.number="form.price" class="form-control" type="number" min="0" step="0.01" placeholder="0.00" />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label>หมายเหตุ / เหตุผล <span class="required">*</span></label>
          <input v-model="form.note" class="form-control" :placeholder="form.type === 'in' ? 'เช่น ซื้อจากตลาด, รับจากซัพพลายเออร์' : 'เช่น ใช้ทำอาหาร, หมดอายุ, เสียหาย'" /> />
        </div>
        <div class="form-group">
          <label>รวมเป็นเงิน</label>
          <div class="total-cost-display">
            {{ formatMoney(computedTotal) }} ฺ
          </div>
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
      form: { ingredient_id: null, type: 'in', quantity: 0, price: 0, note: '' },
      message: '',
      messageType: '',
      confirmAction: null
    }
  },
  computed: {
    selectedIngredient() {
      if (!this.form.ingredient_id) return null
      return this.ingredients.find(i => i.id === this.form.ingredient_id)
    },
    computedTotal() {
      return (this.form.price || 0) * (this.form.quantity || 0)
    }
  },
  async mounted() {
    await this.loadIngredients()
  },
  methods: {
    formatMoney(val) {
      return Number(val || 0).toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    },
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
      if (!this.form.note.trim()) {
        this.message = 'กรุณาระบุหมายเหตุ / เหตุผล'
        this.messageType = 'alert-error'
        return
      }
      if (this.form.type === 'in' && (!this.form.price || this.form.price <= 0)) {
        this.message = 'กรุณาระบุราคาต่อหน่วยสำหรับการรับเข้า'
        this.messageType = 'alert-error'
        return
      }
      const ing = this.selectedIngredient
      const typeName = this.form.type === 'in' ? 'รับเข้า' : 'จ่ายออก'
      const totalText = this.form.price > 0 ? ` (รวม ${this.formatMoney(this.computedTotal)} ฺ)` : ''
      this.confirmAction = {
        title: `ยืนยันการ${typeName}สต็อค`,
        message: `${typeName} "${ing?.name}" จำนวน ${this.form.quantity} ${ing?.unit?.name || ''}${totalText} ต้องการดำเนินการหรือไม่?`,
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
          price: this.form.price || 0,
          note: this.form.note
        })
        this.message = '✓ บันทึกสำเร็จ!'
        this.messageType = 'alert-success'
        this.form.quantity = 0
        this.form.price = 0
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

.total-cost-display {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--primary);
  padding: 0.6rem 1rem;
  background: var(--gray-50);
  border-radius: var(--radius);
  border: 1px solid var(--gray-200);
  min-height: 42px;
  display: flex;
  align-items: center;
}
</style>
