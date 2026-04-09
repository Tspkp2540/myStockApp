<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">แดชบอร์ด</h1>
      <router-link to="/new-sale" class="btn btn-primary btn-lg">
        <span class="material-symbols-outlined">add_shopping_cart</span>
        บันทึกการขาย
      </router-link>
    </div>

    <!-- Stats -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon orange">
          <span class="material-symbols-outlined">receipt_long</span>
        </div>
        <div>
          <div class="stat-value">{{ dashboard.my_sales_today || 0 }}</div>
          <div class="stat-label">ยอดขายของฉันวันนี้</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green">
          <span class="material-symbols-outlined">payments</span>
        </div>
        <div>
          <div class="stat-value">{{ formatMoney(dashboard.my_revenue_today) }}</div>
          <div class="stat-label">รายได้ของฉันวันนี้</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon blue">
          <span class="material-symbols-outlined">storefront</span>
        </div>
        <div>
          <div class="stat-value">{{ dashboard.all_sales_today || 0 }}</div>
          <div class="stat-label">ยอดขายทั้งร้านวันนี้</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green">
          <span class="material-symbols-outlined">monitoring</span>
        </div>
        <div>
          <div class="stat-value">{{ formatMoney(dashboard.all_revenue_today) }}</div>
          <div class="stat-label">รายได้ทั้งร้านวันนี้</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon red">
          <span class="material-symbols-outlined">restaurant</span>
        </div>
        <div>
          <div class="stat-value">{{ dashboard.total_items_today || 0 }}</div>
          <div class="stat-label">จานที่ขายได้วันนี้</div>
        </div>
      </div>
    </div>

    <div style="display:grid;grid-template-columns:1fr 1fr;gap:24px;">
      <!-- Top items today -->
      <div class="card">
        <h3 style="font-size:1rem;font-weight:700;margin-bottom:16px;">
          <span class="material-symbols-outlined" style="vertical-align:middle;margin-right:4px;">trophy</span>
          เมนูขายดีวันนี้
        </h3>
        <div v-if="!dashboard.top_items || dashboard.top_items.length === 0" class="empty-state" style="padding:24px 0;">
          <span class="material-symbols-outlined">restaurant</span>
          <div>ยังไม่มีข้อมูล</div>
        </div>
        <div class="table-wrap" v-else>
          <table>
            <thead>
              <tr>
                <th>#</th>
                <th>เมนู</th>
                <th class="text-right">จำนวน</th>
                <th class="text-right">ยอดขาย</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, idx) in dashboard.top_items" :key="item.menu_item_id">
                <td>{{ idx + 1 }}</td>
                <td>{{ item.name }}</td>
                <td class="text-right font-mono">{{ item.total_qty }}</td>
                <td class="text-right font-mono">{{ formatMoney(item.total_sales) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- My recent sales -->
      <div class="card">
        <h3 style="font-size:1rem;font-weight:700;margin-bottom:16px;">
          <span class="material-symbols-outlined" style="vertical-align:middle;margin-right:4px;">history</span>
          การขายล่าสุดของฉัน
        </h3>
        <div v-if="!dashboard.recent_sales || dashboard.recent_sales.length === 0" class="empty-state" style="padding:24px 0;">
          <span class="material-symbols-outlined">receipt_long</span>
          <div>ยังไม่มีการขาย</div>
        </div>
        <div class="table-wrap" v-else>
          <table>
            <thead>
              <tr>
                <th>เวลา</th>
                <th>รายการ</th>
                <th class="text-right">ยอดรวม</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="sale in dashboard.recent_sales" :key="sale.id">
                <td>{{ formatTime(sale.created_at) }}</td>
                <td>
                  <span v-for="(item, i) in sale.items" :key="item.id">
                    {{ item.menu_item?.name }}×{{ item.quantity }}<span v-if="i < sale.items.length - 1">, </span>
                  </span>
                </td>
                <td class="text-right font-mono" style="font-weight:600;">{{ formatMoney(sale.total) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getFrontofficeDashboard } from '../api/index.js'

export default {
  data() {
    return { dashboard: {} }
  },
  async created() {
    await this.load()
  },
  methods: {
    async load() {
      try {
        const { data } = await getFrontofficeDashboard()
        this.dashboard = data
      } catch (e) {
        console.error(e)
      }
    },
    formatMoney(val) {
      if (!val && val !== 0) return '฿0'
      return '฿' + Number(val).toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    },
    formatTime(d) {
      if (!d) return '-'
      const dt = new Date(d)
      return dt.toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' })
    }
  }
}
</script>

<style scoped>
@media (max-width: 768px) {
  div[style*="grid-template-columns:1fr 1fr"] {
    grid-template-columns: 1fr !important;
  }
}
</style>
