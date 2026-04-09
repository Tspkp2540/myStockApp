<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">ประวัติการขาย</h1>
      <router-link to="/new-sale" class="btn btn-primary">
        <span class="material-symbols-outlined">add_shopping_cart</span>
        บันทึกการขาย
      </router-link>
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

    <!-- Summary -->
    <div class="stats-grid" v-if="sales.length > 0">
      <div class="stat-card">
        <div class="stat-icon orange">
          <span class="material-symbols-outlined">receipt_long</span>
        </div>
        <div>
          <div class="stat-value">{{ sales.length }}</div>
          <div class="stat-label">จำนวนบิล</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green">
          <span class="material-symbols-outlined">payments</span>
        </div>
        <div>
          <div class="stat-value">{{ formatMoney(totalRevenue) }}</div>
          <div class="stat-label">รายได้รวม</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon blue">
          <span class="material-symbols-outlined">restaurant</span>
        </div>
        <div>
          <div class="stat-value">{{ totalItems }}</div>
          <div class="stat-label">จำนวนจานรวม</div>
        </div>
      </div>
    </div>

    <!-- Sales Table -->
    <div class="card">
      <div v-if="sales.length === 0" class="empty-state">
        <span class="material-symbols-outlined">receipt_long</span>
        <div>ยังไม่มีประวัติการขาย</div>
      </div>
      <div class="table-wrap" v-else>
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>วันที่/เวลา</th>
              <th>รายการ</th>
              <th class="text-right">ยอดรวม</th>
              <th>ผู้ขาย</th>
              <th>หมายเหตุ</th>
              <th class="text-center">จัดการ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="sale in sales" :key="sale.id">
              <td>{{ sale.id }}</td>
              <td>{{ formatDate(sale.created_at) }}</td>
              <td>
                <div v-for="item in sale.items" :key="item.id" style="font-size:.8rem;">
                  {{ item.menu_item?.name }} × {{ item.quantity }}
                  <span class="text-muted">({{ formatMoney(item.price) }})</span>
                </div>
              </td>
              <td class="text-right font-mono" style="font-weight:600;">{{ formatMoney(sale.total) }}</td>
              <td>{{ sale.user?.full_name || sale.user?.username || '-' }}</td>
              <td class="text-muted">{{ sale.note || '-' }}</td>
              <td class="text-center">
                <button class="btn btn-outline btn-sm btn-icon" style="color:var(--color-danger)" @click="openDeleteModal(sale)" title="ลบ">
                  <span class="material-symbols-outlined">delete</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Delete Modal with Reason -->
    <div v-if="showDeleteModal" class="modal-overlay" @click.self="showDeleteModal = false">
      <div class="modal" style="max-width:460px;">
        <div class="modal-header">
          <h3>ลบรายการขาย #{{ deleteTarget?.id }}</h3>
          <button class="btn btn-icon" @click="showDeleteModal = false">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <div class="modal-body">
          <div style="margin-bottom:16px;padding:12px;background:var(--color-danger-light);border-radius:var(--radius-md);font-size:.85rem;color:var(--color-danger);">
            <span class="material-symbols-outlined" style="vertical-align:middle;margin-right:4px;font-size:18px;">warning</span>
            รายการที่ลบจะถูกเก็บไว้ให้ผู้ดูแลระบบตรวจสอบ
          </div>

          <div style="margin-bottom:12px;">
            <strong>รายการ:</strong>
            <div v-for="item in deleteTarget?.items" :key="item.id" style="font-size:.85rem;margin-left:8px;">
              {{ item.menu_item?.name }} × {{ item.quantity }} = {{ formatMoney(item.price * item.quantity) }}
            </div>
            <div style="margin-top:4px;font-weight:700;">ยอดรวม: {{ formatMoney(deleteTarget?.total) }}</div>
          </div>

          <div class="form-group">
            <label>เหตุผลในการลบ *</label>
            <textarea v-model="deleteReason" class="form-control" placeholder="กรุณาระบุเหตุผล เช่น ลูกค้ายกเลิก, กรอกผิด" rows="3"></textarea>
          </div>

          <div v-if="deleteError" class="alert alert-error" style="margin-bottom:0;">
            <span class="material-symbols-outlined">error</span>
            {{ deleteError }}
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showDeleteModal = false">ยกเลิก</button>
          <button class="btn btn-danger" @click="doDelete" :disabled="deleting">
            <span class="material-symbols-outlined">delete</span>
            {{ deleting ? 'กำลังลบ...' : 'ยืนยันลบ' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getSales, deleteSale } from '../api/index.js'

export default {
  data() {
    return {
      sales: [],
      dateFrom: '',
      dateTo: '',
      showDeleteModal: false,
      deleteTarget: null,
      deleteReason: '',
      deleteError: '',
      deleting: false
    }
  },
  computed: {
    totalRevenue() { return this.sales.reduce((s, x) => s + (x.total || 0), 0) },
    totalItems() {
      return this.sales.reduce((s, sale) => {
        return s + (sale.items || []).reduce((si, item) => si + item.quantity, 0)
      }, 0)
    }
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
        const { data } = await getSales(params)
        this.sales = data || []
      } catch (e) {
        console.error(e)
      }
    },
    openDeleteModal(sale) {
      this.deleteTarget = sale
      this.deleteReason = ''
      this.deleteError = ''
      this.showDeleteModal = true
    },
    async doDelete() {
      this.deleteError = ''
      if (!this.deleteReason.trim()) {
        this.deleteError = 'กรุณาระบุเหตุผลในการลบ'
        return
      }
      this.deleting = true
      try {
        await deleteSale(this.deleteTarget.id, this.deleteReason.trim())
        this.showDeleteModal = false
        this.deleteTarget = null
        await this.load()
      } catch (e) {
        this.deleteError = e.response?.data?.error || 'เกิดข้อผิดพลาด'
      } finally {
        this.deleting = false
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
