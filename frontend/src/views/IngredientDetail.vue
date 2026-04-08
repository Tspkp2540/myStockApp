<template>
  <div v-if="ingredient">
    <div class="page-header">
      <div class="page-header-top">
        <router-link to="/ingredients" class="btn btn-outline btn-sm">← กลับ</router-link>
      </div>
      <h1 class="page-title">{{ ingredient.name }}</h1>
      <p class="page-subtitle">รายละเอียดการเคลื่อนไหวของวัตถุดิบ</p>
    </div>

    <!-- Ingredient Info Card -->
    <div class="card">
      <div class="card-header">
        <h2 class="section-title">📦 ข้อมูลวัตถุดิบ</h2>
      </div>
      <div class="info-grid">
        <div class="info-item"><strong>ชื่อ:</strong> {{ ingredient.name }}</div>
        <div class="info-item"><strong>หมวดหมู่:</strong> {{ ingredient.category?.name || '-' }}</div>
        <div class="info-item">
          <strong>คงเหลือ:</strong>
          <span :class="{ 'text-danger': ingredient.quantity <= ingredient.min_stock }">
            {{ ingredient.quantity }} {{ ingredient.unit?.name }}
          </span>
        </div>
        <div class="info-item"><strong>สต็อคขั้นต่ำ:</strong> {{ ingredient.min_stock }} {{ ingredient.unit?.name }}</div>
        <div class="info-item">
          <strong>สถานะ:</strong>
          <span v-if="ingredient.quantity <= ingredient.min_stock" class="badge badge-low">⚠ ต่ำ</span>
          <span v-else class="badge badge-ok">✓ ปกติ</span>
        </div>
      </div>
    </div>

    <!-- Stock Transaction Form -->
    <div class="card">
      <div class="card-header">
        <h2 class="section-title">📝 รับ / จ่าย สต็อค — {{ ingredient.name }}</h2>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label>ประเภท <span class="required">*</span></label>
          <select v-model="form.type" class="form-control">
            <option value="in">▲ รับเข้า (Stock In)</option>
            <option value="out">▼ จ่ายออก (Stock Out)</option>
          </select>
        </div>
        <div class="form-group">
          <label>จำนวน <span class="required">*</span></label>
          <input v-model.number="form.quantity" class="form-control" type="number" min="0.1" step="0.1" />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label>ราคาต่อหน่วย (฿) <span v-if="form.type === 'in'" class="required">*</span></label>
          <input v-model.number="form.price" class="form-control" type="number" min="0" step="0.01" placeholder="0.00" />
        </div>
        <div class="form-group">
          <label>รวมเป็นเงิน</label>
          <div class="total-cost-display">
            {{ formatMoney(computedTotal) }} ฿
          </div>
        </div>
      </div>
      <div class="form-group">
        <label>หมายเหตุ / เหตุผล <span class="required">*</span></label>
        <input v-model="form.note" class="form-control" :placeholder="form.type === 'in' ? 'เช่น ซื้อจากตลาด, รับจากซัพพลายเออร์' : 'เช่น ใช้ทำอาหาร, หมดอายุ, เสียหาย'" />
      </div>
      <button class="btn btn-lg btn-success" @click="confirmSubmit">
        บันทึกรายการ
      </button>
      <div v-if="message" :class="['alert', messageType]">
        {{ message }}
      </div>
    </div>

    <!-- Transaction History for this ingredient -->
    <div class="card">
      <div class="card-header" style="display:flex;justify-content:space-between;align-items:center;">
        <h2 class="section-title">📋 ประวัติการเคลื่อนไหว — {{ ingredient.name }}</h2>
        <button class="btn btn-sm btn-success" @click="exportHistory" v-if="transactions.length > 0">📥 Export Excel</button>
      </div>
      <div class="table-scroll">
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>วันที่</th>
              <th>ประเภท</th>
              <th class="text-right">จำนวน</th>
              <th>หน่วย</th>
              <th class="text-right">ราคา/หน่วย</th>
              <th class="text-right">รวมเงิน</th>
              <th>ผู้ทำรายการ</th>
              <th>หมายเหตุ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(txn, i) in transactions" :key="txn.id">
              <td>{{ i + 1 }}</td>
              <td class="nowrap">{{ formatDate(txn.created_at) }}</td>
              <td>
                <span :class="['badge', txn.type === 'in' ? 'badge-in' : 'badge-out']">
                  {{ txn.type === 'in' ? '▲ รับเข้า' : '▼ จ่ายออก' }}
                </span>
              </td>
              <td class="text-right">{{ txn.quantity }}</td>
              <td>{{ ingredient.unit?.name }}</td>
              <td class="text-right">{{ formatMoney(txn.price) }}</td>
              <td class="text-right"><strong>{{ formatMoney(txn.total_cost) }}</strong></td>
              <td class="nowrap">{{ txn.user?.full_name || txn.user?.username || '-' }}</td>
              <td>{{ txn.note || '-' }}</td>
            </tr>
            <tr v-if="transactions.length === 0">
              <td colspan="9" class="table-empty">
                <span class="table-empty-icon">📋</span>
                ยังไม่มีรายการเคลื่อนไหวสำหรับวัตถุดิบนี้
              </td>
            </tr>
            <tr v-if="transactions.length > 0" class="total-row">
              <td colspan="6" class="text-right"><strong>รวมทั้งหมด</strong></td>
              <td class="text-right"><strong>{{ formatMoney(transactionTotal) }}</strong></td>
              <td colspan="2"></td>
            </tr>
          </tbody>
        </table>
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
  <div v-else class="card" style="text-align:center;padding:2rem;">
    กำลังโหลด...
  </div>
</template>

<script>
import { getIngredients, getTransactions, createTransaction } from '../api/index.js'
import ConfirmModal from '../components/ConfirmModal.vue'
import * as XLSX from 'xlsx'

export default {
  components: { ConfirmModal },
  data() {
    return {
      ingredient: null,
      transactions: [],
      form: { type: 'in', quantity: 0, price: 0, note: '' },
      message: '',
      messageType: '',
      confirmAction: null
    }
  },
  computed: {
    computedTotal() {
      return (this.form.price || 0) * (this.form.quantity || 0)
    },
    transactionTotal() {
      return this.transactions.reduce((sum, t) => sum + (t.total_cost || 0), 0)
    }
  },
  async mounted() {
    const id = Number(this.$route.params.id)
    await this.loadIngredient(id)
    await this.loadTransactions(id)
  },
  methods: {
    async loadIngredient(id) {
      const { data } = await getIngredients()
      this.ingredient = data.find(i => i.id === id)
    },
    async loadTransactions(id) {
      const { data } = await getTransactions({ ingredient_id: id || this.ingredient?.id })
      this.transactions = data
    },
    formatDate(d) {
      if (!d) return '-'
      return new Date(d).toLocaleString('th-TH')
    },
    formatMoney(val) {
      return Number(val || 0).toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    },
    confirmSubmit() {
      if (this.form.quantity <= 0) {
        this.message = 'กรุณาระบุจำนวน'
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
      const typeName = this.form.type === 'in' ? 'รับเข้า' : 'จ่ายออก'
      const totalText = this.form.price > 0 ? ` (รวม ${this.formatMoney(this.computedTotal)} ฿)` : ''
      this.confirmAction = {
        title: `ยืนยันการ${typeName}สต็อค`,
        message: `${typeName} "${this.ingredient.name}" จำนวน ${this.form.quantity} ${this.ingredient.unit?.name || ''}${totalText} ต้องการดำเนินการหรือไม่?`,
        confirmText: 'บันทึก',
        variant: this.form.type === 'in' ? 'info' : 'warning',
        onConfirm: () => this.submit()
      }
    },
    async submit() {
      this.confirmAction = null
      try {
        await createTransaction({
          ingredient_id: this.ingredient.id,
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
        await this.loadIngredient(this.ingredient.id)
        await this.loadTransactions(this.ingredient.id)
      } catch (err) {
        this.message = err.response?.data?.error || 'เกิดข้อผิดพลาด'
        this.messageType = 'alert-error'
      }
    },
    exportHistory() {
      const rows = this.transactions.map((txn, i) => ({
        '#': i + 1,
        'วันที่': this.formatDate(txn.created_at),
        'ประเภท': txn.type === 'in' ? 'รับเข้า' : 'จ่ายออก',
        'จำนวน': txn.quantity,
        'หน่วย': this.ingredient.unit?.name || '',
        'ราคา/หน่วย': txn.price || 0,
        'รวมเงิน': txn.total_cost || 0,
        'ผู้ทำรายการ': txn.user?.full_name || txn.user?.username || '-',
        'หมายเหตุ': txn.note || '-'
      }))
      rows.push({
        '#': '', 'วันที่': '', 'ประเภท': '', 'จำนวน': '', 'หน่วย': '',
        'ราคา/หน่วย': 'รวมทั้งหมด', 'รวมเงิน': this.transactionTotal,
        'ผู้ทำรายการ': '', 'หมายเหตุ': ''
      })
      const ws = XLSX.utils.json_to_sheet(rows)
      const wb = XLSX.utils.book_new()
      XLSX.utils.book_append_sheet(wb, ws, this.ingredient.name)
      const dateStr = new Date().toISOString().slice(0, 10)
      XLSX.writeFile(wb, `${this.ingredient.name}_${dateStr}.xlsx`)
    }
  }
}
</script>

<style scoped>
.page-header-top {
  margin-bottom: var(--space-sm);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--space-sm);
}

.info-item {
  padding: var(--space-xs);
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

.table-scroll {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.table-scroll table {
  min-width: 700px;
}

.text-right {
  text-align: right;
}

.text-danger {
  color: var(--danger);
}

.total-row {
  background: var(--gray-50);
  border-top: 2px solid var(--gray-300);
}

.nowrap {
  white-space: nowrap;
}

.alert {
  margin-top: var(--space-md);
}

.btn-lg {
  margin-top: var(--space-sm);
}
</style>
