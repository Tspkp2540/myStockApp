<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">ประวัติการขาย</h1>
      <div style="display:flex;gap:8px;">
        <button class="btn btn-outline" @click="doExport" :disabled="exporting">
          <span class="material-symbols-outlined">download</span>
          {{ exporting ? 'กำลัง Export...' : 'Export Excel' }}
        </button>
        <button class="btn btn-success" @click="openSale">
          <span class="material-symbols-outlined">point_of_sale</span>
          ขายเมนู
        </button>
      </div>
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
        <div class="stat-icon blue">
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
        <div class="stat-icon red">
          <span class="material-symbols-outlined">account_balance_wallet</span>
        </div>
        <div>
          <div class="stat-value">{{ formatMoney(totalCost) }}</div>
          <div class="stat-label">ต้นทุนรวม</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon purple">
          <span class="material-symbols-outlined">trending_up</span>
        </div>
        <div>
          <div class="stat-value" :class="totalProfit >= 0 ? 'text-success' : 'text-danger'">
            {{ formatMoney(totalProfit) }}
          </div>
          <div class="stat-label">กำไรรวม</div>
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
              <th class="text-right">ต้นทุน</th>
              <th class="text-right">กำไร</th>
              <th>ประเภท</th>
              <th>ชำระเงิน</th>
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
                <div v-for="item in sale.items" :key="item.id" style="font-size: .8rem;">
                  {{ item.menu_item?.name }} × {{ item.quantity }}
                  <span class="text-muted">({{ formatMoney(item.price) }})</span>
                </div>
              </td>
              <td class="text-right font-mono" style="font-weight: 600;">{{ formatMoney(sale.total) }}</td>
              <td class="text-right font-mono text-muted">{{ formatMoney(sale.total_cost) }}</td>
              <td class="text-right font-mono" :class="sale.profit >= 0 ? 'text-success' : 'text-danger'">
                {{ formatMoney(sale.profit) }}
              </td>
              <td>
                <span class="badge" :class="sale.sale_type === 'delivery' ? 'badge-blue' : 'badge-green'">
                  {{ sale.sale_type === 'delivery' ? 'Delivery' : 'หน้าร้าน' }}
                </span>
              </td>
              <td>
                <span class="badge" :class="sale.payment_method === 'transfer' ? 'badge-purple' : 'badge-orange'">
                  {{ sale.payment_method === 'transfer' ? 'เงินโอน' : 'เงินสด' }}
                </span>
              </td>
              <td>{{ sale.user?.full_name || sale.user?.username || '-' }}</td>
              <td class="text-muted">{{ sale.note || '-' }}</td>
              <td class="text-center">
                <button class="btn btn-outline btn-sm btn-icon" style="color:var(--color-danger)" @click="confirmDelete(sale)" title="ลบ">
                  <span class="material-symbols-outlined">delete</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- New Sale Modal -->
    <div v-if="showSaleModal" class="modal-overlay" @click.self="showSaleModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>ขายเมนู</h3>
          <button class="btn btn-icon" @click="showSaleModal = false">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <div class="modal-body">
          <div v-for="(row, idx) in saleForm.items" :key="idx" class="ingredient-row">
            <select v-model="row.menu_item_id" class="form-control">
              <option value="">-- เลือกเมนู --</option>
              <option v-for="m in menuItems" :key="m.id" :value="m.id">
                {{ m.name }} ({{ formatMoney(m.price) }})
              </option>
            </select>
            <input v-model.number="row.quantity" class="form-control" type="number" min="1" placeholder="จำนวน" style="max-width: 80px;" />
            <button class="btn btn-icon btn-outline" style="color:var(--color-danger)" @click="saleForm.items.splice(idx, 1)">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>
          <button class="btn btn-outline btn-sm" @click="saleForm.items.push({ menu_item_id: '', quantity: 1 })" style="margin-top: 4px;">
            <span class="material-symbols-outlined">add</span>
            เพิ่มรายการ
          </button>

          <div class="form-group" style="margin-top: 16px;">
            <label>หมายเหตุ</label>
            <input v-model="saleForm.note" class="form-control" placeholder="หมายเหตุ (ถ้ามี)" />
          </div>

          <div class="form-group" style="margin-top: 12px;">
            <label>ประเภทการขาย *</label>
            <div style="display:flex;gap:12px;margin-top:4px;">
              <label style="display:flex;align-items:center;gap:4px;cursor:pointer;">
                <input type="radio" v-model="saleForm.sale_type" value="dine_in" /> หน้าร้าน
              </label>
              <label style="display:flex;align-items:center;gap:4px;cursor:pointer;">
                <input type="radio" v-model="saleForm.sale_type" value="delivery" /> ส่ง Delivery
              </label>
            </div>
          </div>

          <div class="form-group" style="margin-top: 12px;">
            <label>ชำระเงินด้วย *</label>
            <div style="display:flex;gap:12px;margin-top:4px;">
              <label style="display:flex;align-items:center;gap:4px;cursor:pointer;">
                <input type="radio" v-model="saleForm.payment_method" value="cash" /> เงินสด
              </label>
              <label style="display:flex;align-items:center;gap:4px;cursor:pointer;">
                <input type="radio" v-model="saleForm.payment_method" value="transfer" /> เงินโอน
              </label>
            </div>
          </div>

          <!-- Total preview -->
          <div v-if="salePreviewTotal > 0" style="text-align: right; font-size: 1.1rem; font-weight: 700; margin-top: 12px;">
            ยอดรวม: {{ formatMoney(salePreviewTotal) }}
          </div>

          <div v-if="saleError" class="alert alert-error" style="margin-top: 12px;">
            <span class="material-symbols-outlined">error</span>
            {{ saleError }}
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showSaleModal = false">ยกเลิก</button>
          <button class="btn btn-success" @click="submitSale" :disabled="saleSaving">
            {{ saleSaving ? 'กำลังบันทึก...' : 'บันทึกการขาย' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Confirm Sale -->
    <ConfirmModal
      v-if="showConfirmSale"
      title="ยืนยันการขาย"
      :message="`ยืนยันบันทึกการขาย ยอดรวม ${formatMoney(salePreviewTotal)} หรือไม่?`"
      confirmText="บันทึกการขาย"
      variant="success"
      @confirm="doSubmitSale"
      @cancel="showConfirmSale = false"
    />

    <!-- Delete Confirm -->
    <ConfirmModal
      v-if="showDeleteModal"
      title="ยืนยันการลบ"
      :message="`ต้องการลบรายการขาย #${deleteTarget?.id} หรือไม่?`"
      confirmText="ลบ"
      variant="danger"
      @confirm="doDelete"
      @cancel="showDeleteModal = false"
    />
  </div>
</template>

<script>
import { getSales, deleteSale, createSale, getMenuItems, exportSalesExcel } from '../api/index.js'
import ConfirmModal from '../components/ConfirmModal.vue'

export default {
  components: { ConfirmModal },
  data() {
    return {
      sales: [],
      menuItems: [],
      dateFrom: '',
      dateTo: '',
      showSaleModal: false,
      showConfirmSale: false,
      showDeleteModal: false,
      deleteTarget: null,
      saleSaving: false,
      saleError: '',
      exporting: false,
      saleForm: {
        items: [{ menu_item_id: '', quantity: 1 }],
        note: '',
        sale_type: '',
        payment_method: ''
      }
    }
  },
  computed: {
    totalRevenue() { return this.sales.reduce((s, x) => s + (x.total || 0), 0) },
    totalCost() { return this.sales.reduce((s, x) => s + (x.total_cost || 0), 0) },
    totalProfit() { return this.sales.reduce((s, x) => s + (x.profit || 0), 0) },
    salePreviewTotal() {
      return this.saleForm.items.reduce((sum, row) => {
        const menu = this.menuItems.find(m => m.id === Number(row.menu_item_id))
        return sum + (menu ? menu.price * (row.quantity || 0) : 0)
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
    async openSale() {
      this.saleError = ''
      this.saleForm = { items: [{ menu_item_id: '', quantity: 1 }], note: '', sale_type: '', payment_method: '' }
      try {
        const { data } = await getMenuItems({ active: 'true' })
        this.menuItems = data || []
      } catch (e) {
        console.error(e)
      }
      this.showSaleModal = true
    },
    async submitSale() {
      this.saleError = ''
      const items = this.saleForm.items
        .filter(i => i.menu_item_id && i.quantity > 0)
        .map(i => ({ menu_item_id: Number(i.menu_item_id), quantity: Number(i.quantity) }))
      if (items.length === 0) {
        this.saleError = 'กรุณาเลือกเมนูอย่างน้อย 1 รายการ'
        return
      }
      if (!this.saleForm.sale_type) {
        this.saleError = 'กรุณาเลือกประเภทการขาย'
        return
      }
      if (!this.saleForm.payment_method) {
        this.saleError = 'กรุณาเลือกวิธีชำระเงิน'
        return
      }
      this.showConfirmSale = true
    },
    async doSubmitSale() {
      this.showConfirmSale = false
      const items = this.saleForm.items
        .filter(i => i.menu_item_id && i.quantity > 0)
        .map(i => ({ menu_item_id: Number(i.menu_item_id), quantity: Number(i.quantity) }))
      this.saleSaving = true
      try {
        await createSale({ items, note: this.saleForm.note, sale_type: this.saleForm.sale_type, payment_method: this.saleForm.payment_method })
        this.showSaleModal = false
        await this.load()
      } catch (e) {
        this.saleError = e.response?.data?.error || 'เกิดข้อผิดพลาด'
      } finally {
        this.saleSaving = false
      }
    },
    confirmDelete(sale) {
      this.deleteTarget = sale
      this.showDeleteModal = true
    },
    async doDelete() {
      if (!this.deleteTarget) return
      try {
        await deleteSale(this.deleteTarget.id)
        this.showDeleteModal = false
        this.deleteTarget = null
        await this.load()
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
    },
    async doExport() {
      this.exporting = true
      try {
        const params = {}
        if (this.dateFrom) params.date_from = this.dateFrom
        if (this.dateTo) params.date_to = this.dateTo
        const { data } = await exportSalesExcel(params)
        const url = window.URL.createObjectURL(new Blob([data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', 'sales_export.xlsx')
        document.body.appendChild(link)
        link.click()
        link.remove()
        window.URL.revokeObjectURL(url)
      } catch (e) {
        console.error(e)
        alert('ไม่สามารถ Export ได้')
      } finally {
        this.exporting = false
      }
    }
  }
}
</script>
