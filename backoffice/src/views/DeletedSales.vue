<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">รายการขายที่ถูกลบ</h1>
    </div>

    <!-- Filters -->
    <div class="filters">
      <div class="form-group">
        <label>ตั้งแต่</label>
        <input v-model="dateFrom" class="form-control" type="date" @change="load" />
      </div>
      <div class="form-group">
        <label>ถึง</label>
        <input v-model="dateTo" class="form-control" type="date" @change="load" />
      </div>
    </div>

    <div class="card">
      <div v-if="items.length === 0" class="empty-state">
        <span class="material-symbols-outlined">delete_sweep</span>
        <div>ยังไม่มีรายการขายที่ถูกลบ</div>
      </div>
      <div class="table-wrap" v-else>
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>รายการเดิม</th>
              <th class="text-right">ยอดรวม</th>
              <th class="text-right">ต้นทุน</th>
              <th class="text-right">กำไร</th>
              <th>ผู้ขายเดิม</th>
              <th>เหตุผลที่ลบ</th>
              <th>ลบโดย</th>
              <th>วันที่ลบ</th>
              <th>วันที่ขายเดิม</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="d in items" :key="d.id">
              <td>{{ d.original_id }}</td>
              <td>{{ d.items_summary || '-' }}</td>
              <td class="text-right font-mono">{{ formatMoney(d.total) }}</td>
              <td class="text-right font-mono text-muted">{{ formatMoney(d.total_cost) }}</td>
              <td class="text-right font-mono" :class="d.profit >= 0 ? 'text-success' : 'text-danger'">
                {{ formatMoney(d.profit) }}
              </td>
              <td>{{ d.user?.full_name || d.user?.username || '-' }}</td>
              <td>
                <span style="color:var(--color-danger);font-weight:600;">{{ d.delete_reason }}</span>
              </td>
              <td>{{ d.deleted_by?.full_name || d.deleted_by?.username || '-' }}</td>
              <td>{{ formatDate(d.deleted_at) }}</td>
              <td>{{ formatDate(d.original_created_at) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { getDeletedSales } from '../api/index.js'

export default {
  data() {
    return { items: [], dateFrom: '', dateTo: '' }
  },
  async created() {
    await this.load()
  },
  methods: {
    async load() {
      try {
        const params = {}
        if (this.dateFrom) params.date_from = this.dateFrom
        if (this.dateTo) params.date_to = this.dateTo
        const { data } = await getDeletedSales(params)
        this.items = data || []
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
      return dt.toLocaleDateString('th-TH', { year: 'numeric', month: 'short', day: '2-digit' }) +
        ' ' + dt.toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' })
    }
  }
}
</script>
