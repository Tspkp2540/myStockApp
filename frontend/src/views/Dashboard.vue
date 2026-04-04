<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">แดชบอร์ด</h1>
      <p class="page-subtitle">ภาพรวมสถานะสต็อควัตถุดิบของร้าน</p>
    </div>

    <div class="grid-3">
      <div class="card stat-card stat-info">
        <div class="stat-icon icon-info">📦</div>
        <div class="stat-number">{{ dashboard.total_ingredients || 0 }}</div>
        <div class="stat-label">วัตถุดิบทั้งหมด</div>
      </div>
      <div class="card stat-card stat-success">
        <div class="stat-icon icon-success">📂</div>
        <div class="stat-number">{{ dashboard.total_categories || 0 }}</div>
        <div class="stat-label">หมวดหมู่</div>
      </div>
      <div class="card stat-card stat-danger">
        <div class="stat-icon icon-danger">⚠️</div>
        <div class="stat-number">{{ (dashboard.low_stock_items || []).length }}</div>
        <div class="stat-label">สต็อคต่ำ</div>
      </div>
    </div>

    <div class="card" v-if="(dashboard.low_stock_items || []).length > 0">
      <div class="card-header">
        <h2 class="section-title text-danger">⚠️ วัตถุดิบที่สต็อคต่ำ</h2>
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
              <td><span class="badge badge-low">{{ item.quantity }} {{ item.unit?.name }}</span></td>
              <td>{{ item.min_stock }} {{ item.unit?.name }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <h2 class="section-title">📋 รายการเคลื่อนไหวล่าสุด</h2>
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
                  {{ txn.type === 'in' ? '▲ รับเข้า' : '▼ จ่ายออก' }}
                </span>
              </td>
              <td>{{ txn.quantity }}</td>
              <td>{{ txn.note || '-' }}</td>
            </tr>
            <tr v-if="(dashboard.recent_transactions || []).length === 0">
              <td colspan="5" class="table-empty">
                <span class="table-empty-icon">📋</span>
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
import { getDashboard } from '../api/index.js'

export default {
  data() {
    return { dashboard: {} }
  },
  async mounted() {
    await this.load()
  },
  methods: {
    async load() {
      const { data } = await getDashboard()
      this.dashboard = data
    },
    formatDate(d) {
      if (!d) return '-'
      return new Date(d).toLocaleString('th-TH')
    }
  }
}
</script>
