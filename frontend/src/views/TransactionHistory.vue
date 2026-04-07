<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">ประวัติการเคลื่อนไหว</h1>
      <p class="page-subtitle">ดูรายการรับเข้า-จ่ายออกทั้งหมด</p>
    </div>

    <!-- Cost Summary Cards -->
    <div class="summary-cards" v-if="costSummary.length > 0">
      <div class="card summary-overview">
        <div class="card-header">
          <h2 class="section-title">💰 สรุปยอดต้นทุนรายวัตถุดิบ</h2>
        </div>
        <div class="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>วัตถุดิบ</th>
                <th class="text-right">รับเข้า (จำนวน)</th>
                <th class="text-right">ต้นทุนรับเข้า (฿)</th>
                <th class="text-right">จ่ายออก (จำนวน)</th>
                <th class="text-right">ต้นทุนจ่ายออก (฿)</th>
                <th class="text-right">ต้นทุนสุทธิ (฿)</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="s in costSummaryWithNames" :key="s.ingredient_id">
                <td>{{ s.name }}</td>
                <td class="text-right">{{ s.total_in }}</td>
                <td class="text-right text-success">{{ formatMoney(s.cost_in) }}</td>
                <td class="text-right">{{ s.total_out }}</td>
                <td class="text-right text-danger">{{ formatMoney(s.cost_out) }}</td>
                <td class="text-right"><strong>{{ formatMoney(s.cost_in - s.cost_out) }}</strong></td>
              </tr>
              <tr class="total-row">
                <td><strong>รวมทั้งหมด</strong></td>
                <td class="text-right"><strong>{{ totalSummary.totalIn }}</strong></td>
                <td class="text-right text-success"><strong>{{ formatMoney(totalSummary.costIn) }}</strong></td>
                <td class="text-right"><strong>{{ totalSummary.totalOut }}</strong></td>
                <td class="text-right text-danger"><strong>{{ formatMoney(totalSummary.costOut) }}</strong></td>
                <td class="text-right"><strong>{{ formatMoney(totalSummary.costIn - totalSummary.costOut) }}</strong></td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="toolbar">
        <div class="toolbar-filters">
          <select v-model="filterType" class="form-control" @change="load">
            <option value="">ทุกประเภท</option>
            <option value="in">▲ รับเข้า</option>
            <option value="out">▼ จ่ายออก</option>
          </select>
          <select v-model="filterIngredient" class="form-control" @change="load">
            <option :value="''">ทุกวัตถุดิบ</option>
            <option v-for="ing in ingredients" :key="ing.id" :value="ing.id">{{ ing.name }}</option>
          </select>
        </div>
      </div>

      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>วันที่</th>
              <th>วัตถุดิบ</th>
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
              <td>{{ formatDate(txn.created_at) }}</td>
              <td>{{ txn.ingredient?.name }}</td>
              <td>
                <span :class="['badge', txn.type === 'in' ? 'badge-in' : 'badge-out']">
                  {{ txn.type === 'in' ? '▲ รับเข้า' : '▼ จ่ายออก' }}
                </span>
              </td>
              <td class="text-right">{{ txn.quantity }}</td>
              <td>{{ txn.ingredient?.unit?.name }}</td>
              <td class="text-right">{{ formatMoney(txn.price) }}</td>
              <td class="text-right"><strong>{{ formatMoney(txn.total_cost) }}</strong></td>
              <td>{{ txn.user?.full_name || txn.user?.username || '-' }}</td>
              <td>{{ txn.note || '-' }}</td>
            </tr>
            <tr v-if="transactions.length === 0">
              <td colspan="10" class="table-empty">
                <span class="table-empty-icon">📋</span>
                ยังไม่มีรายการเคลื่อนไหว
              </td>
            </tr>
            <tr v-if="transactions.length > 0" class="total-row">
              <td colspan="7" class="text-right"><strong>รวมทั้งหมด</strong></td>
              <td class="text-right"><strong>{{ formatMoney(transactionTotal) }}</strong></td>
              <td colspan="2"></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { getTransactions, getIngredients, getCostSummary } from '../api/index.js'

export default {
  data() {
    return {
      transactions: [],
      ingredients: [],
      costSummary: [],
      filterType: '',
      filterIngredient: ''
    }
  },
  computed: {
    costSummaryWithNames() {
      return this.costSummary.map(s => {
        const ing = this.ingredients.find(i => i.id === s.ingredient_id)
        return { ...s, name: ing?.name || 'ไม่ทราบ' }
      }).sort((a, b) => (b.cost_in - b.cost_out) - (a.cost_in - a.cost_out))
    },
    totalSummary() {
      return this.costSummary.reduce((acc, s) => ({
        totalIn: acc.totalIn + s.total_in,
        totalOut: acc.totalOut + s.total_out,
        costIn: acc.costIn + s.cost_in,
        costOut: acc.costOut + s.cost_out
      }), { totalIn: 0, totalOut: 0, costIn: 0, costOut: 0 })
    },
    transactionTotal() {
      return this.transactions.reduce((sum, t) => sum + (t.total_cost || 0), 0)
    }
  },
  async mounted() {
    const [ingRes, summaryRes] = await Promise.all([getIngredients(), getCostSummary()])
    this.ingredients = ingRes.data
    this.costSummary = summaryRes.data
    await this.load()
  },
  methods: {
    async load() {
      const params = {}
      if (this.filterType) params.type = this.filterType
      if (this.filterIngredient) params.ingredient_id = this.filterIngredient
      const { data } = await getTransactions(params)
      this.transactions = data
    },
    formatDate(d) {
      if (!d) return '-'
      return new Date(d).toLocaleString('th-TH')
    },
    formatMoney(val) {
      return Number(val || 0).toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    }
  }
}
</script>

<style scoped>
.summary-cards {
  margin-bottom: var(--space-lg);
}

.summary-overview .table-wrapper {
  margin-top: var(--space-sm);
}

.text-right {
  text-align: right;
}

.text-success {
  color: var(--success);
}

.text-danger {
  color: var(--danger);
}

.total-row {
  background: var(--gray-50);
  border-top: 2px solid var(--gray-300);
}

.total-row td {
  padding-top: var(--space-sm);
  padding-bottom: var(--space-sm);
}
</style>
