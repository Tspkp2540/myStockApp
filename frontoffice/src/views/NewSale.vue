<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">บันทึกการขาย</h1>
    </div>

    <!-- Menu Grid -->
    <div v-if="menuItems.length === 0" class="card">
      <div class="empty-state">
        <span class="material-symbols-outlined">restaurant_menu</span>
        <div>ยังไม่มีเมนูในระบบ</div>
        <div class="text-muted" style="font-size:.8rem;margin-top:4px;">กรุณาเพิ่มเมนูจาก Back-Office ก่อน</div>
      </div>
    </div>
    <div v-else>
      <!-- Menu Grid grouped by category -->
      <div v-for="group in groupedMenuItems" :key="group.category" style="margin-bottom:20px;">
        <h3 style="font-size:.95rem;font-weight:700;margin-bottom:8px;color:var(--color-text-secondary);">
          <span class="material-symbols-outlined" style="vertical-align:middle;margin-right:4px;font-size:18px;">category</span>
          {{ group.category }}
        </h3>
        <div class="menu-grid">
          <div
            v-for="item in group.items"
            :key="item.id"
            class="menu-card"
            :class="{ selected: getQty(item.id) > 0 }"
            @click="addItem(item.id)"
          >
            <div class="menu-card-name">{{ item.name }}</div>
            <div class="menu-card-price">{{ formatMoney(item.price) }}</div>
            <div v-if="getQty(item.id) > 0" class="menu-card-qty">
              <button class="qty-btn" @click.stop="decreaseItem(item.id)">−</button>
              <span class="qty-value">{{ getQty(item.id) }}</span>
              <button class="qty-btn" @click.stop="addItem(item.id)">+</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Order Summary -->
      <div v-if="orderItems.length > 0" class="card" style="margin-top:16px;">
        <h3 style="font-size:1rem;font-weight:700;margin-bottom:12px;">
          <span class="material-symbols-outlined" style="vertical-align:middle;margin-right:4px;">shopping_cart</span>
          สรุปรายการ
        </h3>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>เมนู</th>
                <th class="text-center">จำนวน</th>
                <th class="text-right">ราคา/จาน</th>
                <th class="text-right">รวม</th>
                <th class="text-center"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="oi in orderItems" :key="oi.menu_item_id">
                <td style="font-weight:600;">{{ getMenuName(oi.menu_item_id) }}</td>
                <td class="text-center">
                  <div style="display:flex;align-items:center;justify-content:center;gap:6px;">
                    <button class="qty-btn" @click="decreaseItem(oi.menu_item_id)">−</button>
                    <span class="qty-value" style="font-size:.95rem;">{{ oi.quantity }}</span>
                    <button class="qty-btn" @click="addItem(oi.menu_item_id)">+</button>
                  </div>
                </td>
                <td class="text-right font-mono">{{ formatMoney(getMenuPrice(oi.menu_item_id)) }}</td>
                <td class="text-right font-mono" style="font-weight:600;">{{ formatMoney(getMenuPrice(oi.menu_item_id) * oi.quantity) }}</td>
                <td class="text-center">
                  <button class="btn btn-icon btn-outline btn-sm" style="color:var(--color-danger)" @click="removeItem(oi.menu_item_id)">
                    <span class="material-symbols-outlined">delete</span>
                  </button>
                </td>
              </tr>
            </tbody>
            <tfoot>
              <tr>
                <td colspan="3" style="font-weight:700;font-size:1.1rem;text-align:right;">ยอดรวมทั้งหมด</td>
                <td class="text-right font-mono" style="font-weight:700;font-size:1.1rem;color:var(--color-primary);">{{ formatMoney(totalAmount) }}</td>
                <td></td>
              </tr>
            </tfoot>
          </table>
        </div>

        <div class="form-group" style="margin-top:16px;">
          <label>ประเภทการขาย *</label>
          <div style="display:flex;gap:12px;margin-top:4px;">
            <label style="display:flex;align-items:center;gap:4px;cursor:pointer;">
              <input type="radio" v-model="saleType" value="dine_in" /> หน้าร้าน
            </label>
            <label style="display:flex;align-items:center;gap:4px;cursor:pointer;">
              <input type="radio" v-model="saleType" value="delivery" /> ส่ง Delivery
            </label>
          </div>
          <span v-if="fieldErrors.saleType" class="field-error" style="color:var(--color-danger);font-size:.8rem;">{{ fieldErrors.saleType }}</span>
        </div>

        <div class="form-group" style="margin-top:12px;">
          <label>ชำระเงินด้วย *</label>
          <div style="display:flex;gap:12px;margin-top:4px;">
            <label style="display:flex;align-items:center;gap:4px;cursor:pointer;">
              <input type="radio" v-model="paymentMethod" value="cash" /> เงินสด
            </label>
            <label style="display:flex;align-items:center;gap:4px;cursor:pointer;">
              <input type="radio" v-model="paymentMethod" value="transfer" /> เงินโอน
            </label>
          </div>
          <span v-if="fieldErrors.paymentMethod" class="field-error" style="color:var(--color-danger);font-size:.8rem;">{{ fieldErrors.paymentMethod }}</span>
        </div>

        <div class="form-group" style="margin-top:12px;">
          <label>หมายเหตุ (ถ้ามี)</label>
          <input v-model="note" class="form-control" placeholder="เช่น โต๊ะ 5, ไข่ดาวเพิ่ม" />
        </div>

        <div v-if="error" class="alert alert-error">
          <span class="material-symbols-outlined">error</span>
          {{ error }}
        </div>
        <div v-if="success" class="alert alert-success">
          <span class="material-symbols-outlined">check_circle</span>
          {{ success }}
        </div>

        <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:12px;">
          <button class="btn btn-outline" @click="confirmClear">ล้างรายการ</button>
          <button class="btn btn-primary btn-lg" @click="confirmSubmit" :disabled="saving">
            <span class="material-symbols-outlined">check</span>
            {{ saving ? 'กำลังบันทึก...' : 'บันทึกการขาย' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Confirm Submit Sale -->
    <ConfirmModal
      v-if="showConfirmSubmit"
      title="ยืนยันการขาย"
      :message="`ยืนยันบันทึกการขาย ${orderItems.length} รายการ ยอดรวม ${formatMoney(totalAmount)} หรือไม่?`"
      confirmText="บันทึกการขาย"
      variant="success"
      @confirm="submitSale"
      @cancel="showConfirmSubmit = false"
    />

    <!-- Confirm Clear -->
    <ConfirmModal
      v-if="showConfirmClear"
      title="ล้างรายการทั้งหมด"
      message="ต้องการล้างรายการทั้งหมดหรือไม่?"
      confirmText="ล้างทั้งหมด"
      variant="warning"
      @confirm="clearOrder"
      @cancel="showConfirmClear = false"
    />
  </div>
</template>

<script>
import { getMenuItems, createSale } from '../api/index.js'
import ConfirmModal from '../components/ConfirmModal.vue'

export default {
  components: { ConfirmModal },
  data() {
    return {
      menuItems: [],
      orderItems: [],
      note: '',
      saleType: '',
      paymentMethod: '',
      saving: false,
      error: '',
      success: '',
      showConfirmSubmit: false,
      showConfirmClear: false,
      fieldErrors: { saleType: '', paymentMethod: '' }
    }
  },
  computed: {
    totalAmount() {
      return this.orderItems.reduce((sum, oi) => {
        return sum + this.getMenuPrice(oi.menu_item_id) * oi.quantity
      }, 0)
    },
    groupedMenuItems() {
      const groups = {}
      for (const item of this.menuItems) {
        const catName = item.menu_category?.name || 'อื่นๆ'
        if (!groups[catName]) groups[catName] = []
        groups[catName].push(item)
      }
      return Object.keys(groups).sort().map(cat => ({ category: cat, items: groups[cat] }))
    }
  },
  async created() {
    await this.loadMenu()
  },
  methods: {
    async loadMenu() {
      try {
        const { data } = await getMenuItems()
        this.menuItems = data || []
      } catch (e) {
        console.error(e)
      }
    },
    getQty(menuItemId) {
      const item = this.orderItems.find(i => i.menu_item_id === menuItemId)
      return item ? item.quantity : 0
    },
    getMenuName(menuItemId) {
      const m = this.menuItems.find(i => i.id === menuItemId)
      return m ? m.name : ''
    },
    getMenuPrice(menuItemId) {
      const m = this.menuItems.find(i => i.id === menuItemId)
      return m ? m.price : 0
    },
    addItem(menuItemId) {
      const existing = this.orderItems.find(i => i.menu_item_id === menuItemId)
      if (existing) {
        existing.quantity++
      } else {
        this.orderItems.push({ menu_item_id: menuItemId, quantity: 1 })
      }
      this.success = ''
    },
    decreaseItem(menuItemId) {
      const existing = this.orderItems.find(i => i.menu_item_id === menuItemId)
      if (existing) {
        existing.quantity--
        if (existing.quantity <= 0) {
          this.orderItems = this.orderItems.filter(i => i.menu_item_id !== menuItemId)
        }
      }
    },
    removeItem(menuItemId) {
      this.orderItems = this.orderItems.filter(i => i.menu_item_id !== menuItemId)
    },
    clearOrder() {
      this.showConfirmClear = false
      this.orderItems = []
      this.note = ''
      this.saleType = ''
      this.paymentMethod = ''
      this.error = ''
      this.success = ''
      this.fieldErrors = { saleType: '', paymentMethod: '' }
    },
    confirmClear() {
      this.showConfirmClear = true
    },
    confirmSubmit() {
      this.error = ''
      this.success = ''
      this.fieldErrors = { saleType: '', paymentMethod: '' }
      let valid = true
      if (this.orderItems.length === 0) {
        this.error = 'กรุณาเลือกเมนูอย่างน้อย 1 รายการ'
        return
      }
      if (!this.saleType) {
        this.fieldErrors.saleType = 'กรุณาเลือกประเภทการขาย'
        valid = false
      }
      if (!this.paymentMethod) {
        this.fieldErrors.paymentMethod = 'กรุณาเลือกวิธีชำระเงิน'
        valid = false
      }
      if (!valid) return
      this.showConfirmSubmit = true
    },
    async submitSale() {
      this.showConfirmSubmit = false
      this.error = ''
      this.saving = true
      try {
        const payload = {
          items: this.orderItems.map(i => ({
            menu_item_id: i.menu_item_id,
            quantity: i.quantity
          })),
          sale_type: this.saleType,
          payment_method: this.paymentMethod,
          note: this.note
        }
        await createSale(payload)
        this.success = `บันทึกการขายสำเร็จ — ยอดรวม ${this.formatMoney(this.totalAmount)}`
        this.orderItems = []
        this.note = ''
        this.saleType = ''
        this.paymentMethod = ''
      } catch (e) {
        this.error = e.response?.data?.error || 'เกิดข้อผิดพลาด'
      } finally {
        this.saving = false
      }
    },
    formatMoney(val) {
      if (!val && val !== 0) return '฿0'
      return '฿' + Number(val).toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    }
  }
}
</script>
