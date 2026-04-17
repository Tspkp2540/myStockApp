<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">บันทึกการขาย</h1>
      <div class="new-sale-header-actions">
        <button class="btn btn-outline" type="button" @click="toggleSummary" :disabled="!menuItems.length">
          <span class="material-symbols-outlined">shopping_cart</span>
          {{ isSummaryOpen ? 'ซ่อนสรุป' : 'ดูสรุป' }}
          <span v-if="orderItems.length > 0" class="summary-count-badge">{{ orderItems.length }}</span>
        </button>
      </div>
    </div>

    <div v-if="menuItems.length === 0" class="card">
      <div class="empty-state">
        <span class="material-symbols-outlined">restaurant_menu</span>
        <div>ยังไม่มีเมนูในระบบ</div>
        <div class="text-muted" style="font-size:.8rem;margin-top:4px;">กรุณาเพิ่มเมนูจาก Back-Office ก่อน</div>
      </div>
    </div>
    <div v-else>
      <div class="new-sale-layout">
        <section class="new-sale-menu-section">
          <div class="card new-sale-tools">
            <div class="new-sale-search-wrap">
              <span class="material-symbols-outlined new-sale-search-icon">search</span>
              <input
                v-model.trim="searchQuery"
                class="form-control new-sale-search-input"
                placeholder="ค้นหาเมนูหรือหมวดหมู่"
              />
            </div>
            <div class="new-sale-tool-actions">
              <button class="btn btn-outline btn-sm" type="button" @click="setAllCategoriesCollapsed(false)">
                ขยายทั้งหมด
              </button>
              <button class="btn btn-outline btn-sm" type="button" @click="setAllCategoriesCollapsed(true)">
                ย่อทั้งหมด
              </button>
            </div>
          </div>

          <div v-if="filteredGroupedMenuItems.length > 0" class="category-chip-bar">
            <button
              v-for="group in filteredGroupedMenuItems"
              :key="`chip-${group.category}`"
              type="button"
              class="category-chip"
              :class="{ active: activeCategory === group.category }"
              @click="jumpToCategory(group.category)"
            >
              <span>{{ group.category }}</span>
              <span class="category-chip-count">{{ group.items.length }}</span>
            </button>
          </div>

          <div v-if="filteredGroupedMenuItems.length === 0" class="card">
            <div class="empty-state new-sale-empty-results">
              <span class="material-symbols-outlined">search_off</span>
              <div>ไม่พบเมนูที่ค้นหา</div>
              <div class="text-muted" style="font-size:.8rem;margin-top:4px;">ลองค้นหาด้วยชื่อเมนูหรือชื่อหมวดหมู่</div>
            </div>
          </div>

          <div
            v-for="group in filteredGroupedMenuItems"
            :key="group.category"
            class="menu-category-group"
            :ref="el => setCategoryRef(group.category, el)"
          >
            <button class="menu-category-toggle" type="button" @click="toggleCategory(group.category)">
              <div class="menu-category-title-wrap">
                <span class="material-symbols-outlined menu-category-icon">category</span>
                <div>
                  <div class="menu-category-title">{{ group.category }}</div>
                  <div class="menu-category-meta">{{ group.items.length }} เมนู <span v-if="getCategorySelectedCount(group.category) > 0">• เลือกแล้ว {{ getCategorySelectedCount(group.category) }}</span></div>
                </div>
              </div>
              <span class="material-symbols-outlined menu-category-arrow">{{ isCategoryCollapsed(group.category) ? 'expand_more' : 'expand_less' }}</span>
            </button>

            <div v-show="!isCategoryCollapsed(group.category)" class="menu-grid menu-category-grid">
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
        </section>

        <div v-if="isMobileSummaryVisible" class="summary-overlay" @click="closeSummary"></div>

        <aside class="order-summary-panel" :class="{ 'is-open': isSummaryOpen }">
          <div class="card order-summary-card">
            <div class="order-summary-header">
              <h3 class="order-summary-title">
                <span class="material-symbols-outlined">shopping_cart</span>
                สรุปรายการ
              </h3>
              <button class="btn btn-icon btn-outline order-summary-close" type="button" @click="closeSummary">
                <span class="material-symbols-outlined">close</span>
              </button>
            </div>

            <div v-if="orderItems.length === 0" class="empty-state order-summary-empty">
              <span class="material-symbols-outlined">shopping_cart_checkout</span>
              <div>ยังไม่มีรายการที่เลือก</div>
              <div class="text-muted" style="font-size:.8rem;margin-top:4px;">แตะเมนูเพื่อเพิ่มลงออเดอร์</div>
            </div>

            <template v-else>
              <div class="table-wrap order-summary-table-wrap">
                <table>
                  <thead>
                    <tr>
                      <th>เมนู</th>
                      <th class="text-center">จำนวน</th>
                      <th class="text-right">รวม</th>
                      <th class="text-center"></th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="oi in orderItems" :key="oi.menu_item_id">
                      <td style="font-weight:600;">
                        <div>{{ getMenuName(oi.menu_item_id) }}</div>
                        <div class="text-muted order-summary-item-price">{{ formatMoney(getMenuPrice(oi.menu_item_id)) }} / จาน</div>
                      </td>
                      <td class="text-center">
                        <div style="display:flex;align-items:center;justify-content:center;gap:6px;">
                          <button class="qty-btn" @click="decreaseItem(oi.menu_item_id)">−</button>
                          <span class="qty-value" style="font-size:.95rem;">{{ oi.quantity }}</span>
                          <button class="qty-btn" @click="addItem(oi.menu_item_id)">+</button>
                        </div>
                      </td>
                      <td class="text-right font-mono" style="font-weight:600;">{{ formatMoney(getMenuPrice(oi.menu_item_id) * oi.quantity) }}</td>
                      <td class="text-center">
                        <button class="btn btn-icon btn-outline btn-sm" style="color:var(--color-danger)" @click="removeItem(oi.menu_item_id)">
                          <span class="material-symbols-outlined">delete</span>
                        </button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </template>

            <div class="order-summary-total">
              <span>ยอดรวมทั้งหมด</span>
              <strong class="font-mono">{{ formatMoney(totalAmount) }}</strong>
            </div>

            <div class="form-group" style="margin-top:16px;">
              <div class="form-section-title">ประเภทการขาย *</div>
              <div class="order-summary-toggle-row">
                <button type="button" class="toggle-btn" :class="{ active: saleType === 'dine_in' }" @click="saleType = 'dine_in'">
                  <span class="material-symbols-outlined" style="font-size:18px;">storefront</span>
                  หน้าร้าน
                </button>
                <button type="button" class="toggle-btn" :class="{ active: saleType === 'delivery' }" @click="saleType = 'delivery'" style="--toggle-color:var(--color-info);--toggle-bg:var(--color-info-light);">
                  <span class="material-symbols-outlined" style="font-size:18px;">delivery_dining</span>
                  Delivery
                </button>
              </div>
              <span v-if="fieldErrors.saleType" class="field-error" style="color:var(--color-danger);font-size:.8rem;">{{ fieldErrors.saleType }}</span>
            </div>

            <div class="form-group" style="margin-top:12px;">
              <div class="form-section-title">ชำระเงินด้วย *</div>
              <div class="order-summary-toggle-row">
                <button type="button" class="toggle-btn" :class="{ active: paymentMethod === 'cash' }" @click="paymentMethod = 'cash'" style="--toggle-color:var(--color-success);--toggle-bg:var(--color-success-light);">
                  <span class="material-symbols-outlined" style="font-size:18px;">payments</span>
                  เงินสด
                </button>
                <button type="button" class="toggle-btn" :class="{ active: paymentMethod === 'transfer' }" @click="paymentMethod = 'transfer'" style="--toggle-color:#7c3aed;--toggle-bg:#f5f3ff;">
                  <span class="material-symbols-outlined" style="font-size:18px;">account_balance</span>
                  เงินโอน
                </button>
              </div>
              <span v-if="fieldErrors.paymentMethod" class="field-error" style="color:var(--color-danger);font-size:.8rem;">{{ fieldErrors.paymentMethod }}</span>
            </div>

            <div class="form-group" style="margin-top:12px;">
              <label for="sale-note">หมายเหตุ (ถ้ามี)</label>
              <input id="sale-note" v-model="note" class="form-control" placeholder="เช่น โต๊ะ 5, ไข่ดาวเพิ่ม" />
            </div>

            <div v-if="error" class="alert alert-error">
              <span class="material-symbols-outlined">error</span>
              {{ error }}
            </div>
            <div v-if="success" class="alert alert-success">
              <span class="material-symbols-outlined">check_circle</span>
              {{ success }}
            </div>

            <div class="order-summary-actions">
              <button class="btn btn-outline" @click="confirmClear" :disabled="orderItems.length === 0">ล้างรายการ</button>
              <button class="btn btn-primary btn-lg" @click="confirmSubmit" :disabled="saving || orderItems.length === 0">
                <span class="material-symbols-outlined">check</span>
                {{ saving ? 'กำลังบันทึก...' : 'บันทึกการขาย' }}
              </button>
            </div>
          </div>
        </aside>
      </div>

      <button v-if="menuItems.length > 0" class="mobile-summary-fab" type="button" @click="toggleSummary">
        <span class="material-symbols-outlined">shopping_cart</span>
        <span v-if="orderItems.length > 0" class="summary-count-badge">{{ orderItems.length }}</span>
      </button>
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
import { nextTick } from 'vue'
import { getMenuItems, createSale } from '../api/index.js'
import ConfirmModal from '../components/ConfirmModal.vue'

export default {
  components: { ConfirmModal },
  data() {
    return {
      menuItems: [],
      orderItems: [],
      searchQuery: '',
      collapsedCategories: {},
      activeCategory: '',
      categoryRefs: {},
      isSummaryOpen: false,
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
    },
    filteredGroupedMenuItems() {
      const query = this.searchQuery.trim().toLowerCase()
      if (!query) return this.groupedMenuItems

      return this.groupedMenuItems
        .map(group => {
          const categoryMatch = group.category.toLowerCase().includes(query)
          const items = categoryMatch
            ? group.items
            : group.items.filter(item => item.name.toLowerCase().includes(query))

          return { category: group.category, items }
        })
        .filter(group => group.items.length > 0)
    },
    isMobileSummaryVisible() {
      return this.isSummaryOpen
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
        this.initializeCollapsedCategories()
      } catch (e) {
        console.error(e)
      }
    },
    initializeCollapsedCategories() {
      const nextState = {}

      this.groupedMenuItems.forEach((group, index) => {
        nextState[group.category] = this.collapsedCategories[group.category] ?? index !== 0
      })

      this.collapsedCategories = nextState
      if (!this.activeCategory && this.groupedMenuItems.length > 0) {
        this.activeCategory = this.groupedMenuItems[0].category
      }
    },
    toggleCategory(category) {
      this.collapsedCategories = {
        ...this.collapsedCategories,
        [category]: !this.isCategoryCollapsed(category)
      }
      this.activeCategory = category
    },
    isCategoryCollapsed(category) {
      if (this.searchQuery) return false
      return !!this.collapsedCategories[category]
    },
    setAllCategoriesCollapsed(collapsed) {
      const nextState = {}
      this.groupedMenuItems.forEach(group => {
        nextState[group.category] = collapsed
      })
      this.collapsedCategories = nextState
    },
    setCategoryRef(category, el) {
      if (el) {
        this.categoryRefs = {
          ...this.categoryRefs,
          [category]: el
        }
        return
      }

      if (this.categoryRefs[category]) {
        const nextRefs = { ...this.categoryRefs }
        delete nextRefs[category]
        this.categoryRefs = nextRefs
      }
    },
    async jumpToCategory(category) {
      this.activeCategory = category

      if (this.isCategoryCollapsed(category)) {
        this.collapsedCategories = {
          ...this.collapsedCategories,
          [category]: false
        }
      }

      await nextTick()
      this.categoryRefs[category]?.scrollIntoView({ behavior: 'smooth', block: 'start' })
    },
    getCategorySelectedCount(category) {
      const group = this.groupedMenuItems.find(item => item.category === category)
      if (!group) return 0
      return group.items.reduce((count, item) => count + this.getQty(item.id), 0)
    },
    toggleSummary() {
      this.isSummaryOpen = !this.isSummaryOpen
    },
    closeSummary() {
      this.isSummaryOpen = false
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
      this.isSummaryOpen = false
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
        this.isSummaryOpen = false
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

<style scoped>
.new-sale-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.new-sale-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 16px;
  align-items: start;
}

.new-sale-tools {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.new-sale-search-wrap {
  position: relative;
  flex: 1;
}

.new-sale-search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-muted);
  font-size: 20px;
}

.new-sale-search-input {
  padding-left: 40px;
}

.new-sale-tool-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.category-chip-bar {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 2px 0 12px;
  margin-bottom: 8px;
  scrollbar-width: thin;
}

.category-chip {
  flex: 0 0 auto;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  border: 1px solid var(--color-border);
  background: var(--color-white);
  color: var(--color-text-secondary);
  border-radius: 999px;
  padding: 8px 12px;
  cursor: pointer;
  box-shadow: var(--shadow-sm);
  font-weight: 600;
}

.category-chip.active {
  border-color: var(--color-primary);
  color: var(--color-primary);
  background: var(--color-primary-light);
}

.category-chip-count {
  min-width: 22px;
  height: 22px;
  padding: 0 6px;
  border-radius: 999px;
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: .75rem;
}

.category-chip.active .category-chip-count {
  background: var(--color-primary);
  color: var(--color-white);
}

.menu-category-group {
  margin-bottom: 16px;
}

.menu-category-toggle {
  width: 100%;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-white);
  padding: 14px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  cursor: pointer;
  text-align: left;
  box-shadow: var(--shadow-sm);
}

.menu-category-title-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.menu-category-icon {
  width: 36px;
  height: 36px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--color-primary-light);
  color: var(--color-primary);
}

.menu-category-title {
  font-size: .98rem;
  font-weight: 700;
}

.menu-category-meta {
  color: var(--color-text-secondary);
  font-size: .82rem;
}

.menu-category-arrow {
  color: var(--color-text-secondary);
}

.menu-category-grid {
  margin-top: 12px;
  margin-bottom: 0;
}

.order-summary-panel {
  position: sticky;
  top: 24px;
}

.order-summary-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.order-summary-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.order-summary-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 1rem;
  font-weight: 700;
}

.order-summary-close {
  display: none;
}

.order-summary-table-wrap {
  max-height: 340px;
  overflow-y: auto;
}

.order-summary-item-price {
  font-size: .78rem;
}

.order-summary-total {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 16px;
  border-radius: var(--radius-md);
  background: var(--color-primary-light);
  color: var(--color-text);
}

.order-summary-total strong {
  color: var(--color-primary);
  font-size: 1.1rem;
}

.order-summary-toggle-row {
  display: flex;
  gap: 10px;
  margin-top: 4px;
  flex-wrap: wrap;
}

.form-section-title {
  display: block;
  font-weight: 600;
  margin-bottom: 6px;
  color: var(--color-text-secondary);
  font-size: .8rem;
}

.order-summary-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 4px;
  flex-wrap: wrap;
}

.summary-count-badge {
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  border-radius: 999px;
  background: var(--color-primary);
  color: var(--color-white);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: .72rem;
  font-weight: 700;
}

.summary-overlay {
  display: none;
}

.mobile-summary-fab {
  display: none;
}

@media (max-width: 1024px) {
  .new-sale-layout {
    grid-template-columns: minmax(0, 1fr);
  }

  .order-summary-panel {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    width: min(420px, 92vw);
    padding: 16px;
    background: transparent;
    z-index: 120;
    transform: translateX(100%);
    transition: transform .2s ease;
  }

  .order-summary-panel.is-open {
    transform: translateX(0);
  }

  .order-summary-card {
    height: 100%;
    border-radius: 20px 0 0 20px;
    box-shadow: var(--shadow-lg);
    overflow: hidden;
    overflow-y: auto;
  }

  .order-summary-close {
    display: inline-flex;
  }

  .summary-overlay {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(28, 25, 23, .45);
    z-index: 110;
  }

  .mobile-summary-fab {
    position: fixed;
    right: 16px;
    bottom: 16px;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    border: none;
    border-radius: 999px;
    background: var(--color-primary);
    color: var(--color-white);
    padding: 12px 16px;
    box-shadow: var(--shadow-lg);
    z-index: 100;
    cursor: pointer;
  }
}

@media (max-width: 768px) {
  .category-chip-bar {
    position: sticky;
    top: 64px;
    z-index: 20;
    background: var(--color-bg);
    margin: 0 -4px 8px;
    padding: 2px 4px 12px;
  }

  .new-sale-tools {
    flex-direction: column;
    align-items: stretch;
  }

  .new-sale-tool-actions {
    width: 100%;
  }

  .new-sale-tool-actions .btn {
    flex: 1;
    justify-content: center;
  }

  .menu-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .menu-card {
    padding: 14px 12px;
  }

  .order-summary-actions {
    flex-direction: column-reverse;
  }

  .order-summary-actions .btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
