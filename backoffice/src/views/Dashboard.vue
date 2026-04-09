<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">แดชบอร์ด</h1>
    </div>

    <!-- Stats -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon blue">
          <span class="material-symbols-outlined">receipt_long</span>
        </div>
        <div>
          <div class="stat-value">{{ dashboard.today_sales || 0 }}</div>
          <div class="stat-label">ยอดขายวันนี้ (รายการ)</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green">
          <span class="material-symbols-outlined">payments</span>
        </div>
        <div>
          <div class="stat-value">{{ formatMoney(dashboard.today_revenue) }}</div>
          <div class="stat-label">รายได้วันนี้</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon red">
          <span class="material-symbols-outlined">account_balance_wallet</span>
        </div>
        <div>
          <div class="stat-value">{{ formatMoney(dashboard.today_cost) }}</div>
          <div class="stat-label">ต้นทุนวันนี้</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon purple">
          <span class="material-symbols-outlined">trending_up</span>
        </div>
        <div>
          <div class="stat-value" :class="dashboard.today_profit >= 0 ? 'text-success' : 'text-danger'">
            {{ formatMoney(dashboard.today_profit) }}
          </div>
          <div class="stat-label">กำไรวันนี้</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon orange">
          <span class="material-symbols-outlined">restaurant_menu</span>
        </div>
        <div>
          <div class="stat-value">{{ dashboard.total_menu_items || 0 }}</div>
          <div class="stat-label">เมนูทั้งหมด</div>
        </div>
      </div>
    </div>

    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 24px;">
      <!-- Top Selling -->
      <div class="card">
        <h3 style="font-size: 1rem; font-weight: 700; margin-bottom: 16px;">
          <span class="material-symbols-outlined" style="vertical-align: middle; margin-right: 4px;">trophy</span>
          เมนูขายดี (30 วัน)
        </h3>
        <div v-if="!dashboard.top_items || dashboard.top_items.length === 0" class="empty-state" style="padding: 24px 0;">
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

      <!-- Recent Sales -->
      <div class="card">
        <h3 style="font-size: 1rem; font-weight: 700; margin-bottom: 16px;">
          <span class="material-symbols-outlined" style="vertical-align: middle; margin-right: 4px;">history</span>
          การขายล่าสุด
        </h3>
        <div v-if="!dashboard.recent_sales || dashboard.recent_sales.length === 0" class="empty-state" style="padding: 24px 0;">
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
                <th class="text-right">กำไร</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="sale in dashboard.recent_sales" :key="sale.id">
                <td>{{ formatDate(sale.created_at) }}</td>
                <td>
                  <span v-for="(item, i) in sale.items" :key="item.id">
                    {{ item.menu_item?.name }}×{{ item.quantity }}<span v-if="i < sale.items.length - 1">, </span>
                  </span>
                </td>
                <td class="text-right font-mono">{{ formatMoney(sale.total) }}</td>
                <td class="text-right font-mono" :class="sale.profit >= 0 ? 'text-success' : 'text-danger'">
                  {{ formatMoney(sale.profit) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getBackofficeDashboard } from '../api/index.js'

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
        const { data } = await getBackofficeDashboard()
        this.dashboard = data
      } catch (e) {
        console.error(e)
      }
    },
    formatMoney(val) {
      if (!val && val !== 0) return '฿0'
      return '฿' + Number(val).toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    },
    formatDate(d) {
      if (!d) return '-'
      const dt = new Date(d)
      return dt.toLocaleDateString('th-TH', { day: '2-digit', month: 'short' }) + ' ' + dt.toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' })
    }
  }
}
</script>

<style scoped>
@media (max-width: 768px) {
  div[style*="grid-template-columns: 1fr 1fr"] {
    grid-template-columns: 1fr !important;
  }
}
</style>
