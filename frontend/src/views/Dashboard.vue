<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">แดชบอร์ด</h1>
      <p class="page-subtitle">ภาพรวมสถานะสต็อควัตถุดิบของร้าน</p>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-card-icon stat-card-icon--primary">
          <span class="material-symbols-outlined">inventory_2</span>
        </div>
        <div class="stat-card-content">
          <span class="stat-card-value">{{ dashboard.total_ingredients || 0 }}</span>
          <span class="stat-card-label">วัตถุดิบทั้งหมด</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-card-icon stat-card-icon--success">
          <span class="material-symbols-outlined">category</span>
        </div>
        <div class="stat-card-content">
          <span class="stat-card-value">{{ dashboard.total_categories || 0 }}</span>
          <span class="stat-card-label">หมวดหมู่</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-card-icon stat-card-icon--danger">
          <span class="material-symbols-outlined">warning</span>
        </div>
        <div class="stat-card-content">
          <span class="stat-card-value">{{ (dashboard.low_stock_items || []).length }}</span>
          <span class="stat-card-label">สต็อคต่ำ</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-card-icon stat-card-icon--info">
          <span class="material-symbols-outlined">account_balance_wallet</span>
        </div>
        <div class="stat-card-content">
          <span class="stat-card-value">{{ formatMoney(totalNetCost) }}</span>
          <span class="stat-card-label">มูลค่าสต็อครวม (฿)</span>
        </div>
      </div>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-card-icon stat-card-icon--success">
          <span class="material-symbols-outlined">trending_up</span>
        </div>
        <div class="stat-card-content">
          <span class="stat-card-value text-success">{{ formatMoney(totalCostIn) }}</span>
          <span class="stat-card-label">ยอดรับเข้ารวม (฿)</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-card-icon stat-card-icon--danger">
          <span class="material-symbols-outlined">trending_down</span>
        </div>
        <div class="stat-card-content">
          <span class="stat-card-value text-danger">{{ formatMoney(totalCostOut) }}</span>
          <span class="stat-card-label">ยอดจ่ายออกรวม (฿)</span>
        </div>
      </div>
    </div>

    <div class="card" v-if="(dashboard.low_stock_items || []).length > 0">
      <div class="card-header">
        <h2 class="section-title text-danger">
          <span class="material-symbols-outlined">warning</span>
          วัตถุดิบที่สต็อคต่ำ
        </h2>
      </div>
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>ชื่อ</th>
              <th>หมวดหมู่</th>
              <th>คงเหลือ</th>
              <th>ขั้นต่ำ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in dashboard.low_stock_items" :key="item.id">
              <td>{{ item.name }}</td>
              <td>{{ item.category?.name || '-' }}</td>
              <td><span class="badge badge-low"><span class="material-symbols-outlined">error</span> {{ item.quantity }} {{ item.unit?.name }}</span></td>
              <td>{{ item.min_stock }} {{ item.unit?.name }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <h2 class="section-title">
          <span class="material-symbols-outlined">receipt_long</span>
          รายการเคลื่อนไหวล่าสุด
        </h2>
      </div>
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>วันที่</th>
              <th>วัตถุดิบ</th>
              <th>ประเภท</th>
              <th>จำนวน</th>
              <th>หมายเหตุ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="txn in dashboard.recent_transactions || []" :key="txn.id">
              <td>{{ formatDate(txn.created_at) }}</td>
              <td>{{ txn.ingredient?.name }}</td>
              <td>
                <span :class="['badge', txn.type === 'in' ? 'badge-in' : 'badge-out']">
                  <span class="material-symbols-outlined">{{ txn.type === 'in' ? 'arrow_upward' : 'arrow_downward' }}</span>
                  {{ txn.type === 'in' ? 'รับเข้า' : 'จ่ายออก' }}
                </span>
              </td>
              <td>{{ txn.quantity }}</td>
              <td>{{ txn.note || '-' }}</td>
            </tr>
            <tr v-if="(dashboard.recent_transactions || []).length === 0">
              <td colspan="5" class="table-empty">
                <span class="table-empty-icon"><span class="material-symbols-outlined">inbox</span></span>
                ยังไม่มีรายการเคลื่อนไหว
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { getDashboard, getCostSummary } from '../api/index.js'

export default {
  data() {
    return {
      dashboard: {},
      costSummary: []
    }
  },
  computed: {
    totalCostIn() {
      return this.costSummary.reduce((sum, s) => sum + (s.cost_in || 0), 0)
    },
    totalCostOut() {
      return this.costSummary.reduce((sum, s) => sum + (s.cost_out || 0), 0)
    },
    totalNetCost() {
      return this.totalCostIn - this.totalCostOut
    }
  },
  async mounted() {
    await this.load()
  },
  methods: {
    async load() {
      const [dashRes, costRes] = await Promise.all([getDashboard(), getCostSummary()])
      this.dashboard = dashRes.data
      this.costSummary = costRes.data || []
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
