<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">จัดการเมนู</h1>
      <button class="btn btn-primary" @click="openAdd">
        <span class="material-symbols-outlined">add</span>
        เพิ่มเมนู
      </button>
    </div>

    <!-- Search -->
    <div class="filters">
      <div class="form-group">
        <label>ค้นหา</label>
        <input v-model="search" class="form-control" placeholder="ชื่อเมนู..." @input="load" />
      </div>
      <div class="form-group">
        <label>หมวดหมู่</label>
        <select v-model="filterCategory" class="form-control" @change="load">
          <option value="">-- ทั้งหมด --</option>
          <option v-for="cat in menuCategories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
        </select>
      </div>
    </div>

    <!-- Menu Table -->
    <div class="card">
      <div v-if="items.length === 0" class="empty-state">
        <span class="material-symbols-outlined">restaurant_menu</span>
        <div>ยังไม่มีเมนู</div>
      </div>
      <div class="table-wrap" v-else>
        <table>
          <thead>
            <tr>
              <th>ลำดับ</th>
              <th>ชื่อเมนู</th>
              <th>หมวดหมู่</th>
              <th class="text-right">ราคาขาย</th>
              <th class="text-right">ราคาต้นทุน</th>
              <th class="text-right">ต้นทุน(วัตถุดิบ)</th>
              <th class="text-right">กำไร</th>
              <th>สถานะ</th>
              <th class="text-center">จัดการ</th>
            </tr>
          </thead>
          <template v-for="group in groupedItems" :key="group.category">
            <tbody>
              <tr style="background:var(--color-bg-secondary);">
                <td colspan="9" style="font-weight:700;font-size:.95rem;padding:8px 12px;">
                  <span class="material-symbols-outlined" style="vertical-align:middle;margin-right:4px;font-size:18px;">category</span>
                  {{ group.category }}
                </td>
              </tr>
              <tr v-for="item in group.items" :key="item.id">
                <td>{{ item.sort_order }}</td>
                <td style="font-weight: 600;">{{ item.name }}</td>
                <td class="text-muted">{{ item.menu_category?.name || '-' }}</td>
                <td class="text-right font-mono">{{ formatMoney(item.price) }}</td>
                <td class="text-right font-mono">{{ formatMoney(item.cost_price) }}</td>
                <td class="text-right font-mono">{{ formatMoney(item.cost_per_unit) }}</td>
                <td class="text-right font-mono" :class="item.price - item.cost_per_unit >= 0 ? 'text-success' : 'text-danger'">
                  {{ formatMoney(item.price - item.cost_per_unit) }}
                </td>
                <td>
                  <span class="badge" :class="item.active ? 'badge-green' : 'badge-gray'">
                    {{ item.active ? 'เปิดขาย' : 'ปิด' }}
                  </span>
                </td>
                <td class="text-center">
                  <button class="btn btn-outline btn-sm btn-icon" @click="openEdit(item)" title="แก้ไข">
                    <span class="material-symbols-outlined">edit</span>
                  </button>
                  <button class="btn btn-outline btn-sm btn-icon" style="margin-left:4px;color:var(--color-danger)" @click="confirmDelete(item)" title="ลบ">
                    <span class="material-symbols-outlined">delete</span>
                  </button>
                </td>
              </tr>
            </tbody>
          </template>
        </table>
      </div>
    </div>

    <!-- Add / Edit Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editing ? 'แก้ไขเมนู' : 'เพิ่มเมนูใหม่' }}</h3>
          <button class="btn btn-icon" @click="showModal = false">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>ชื่อเมนู *</label>
            <input v-model="form.name" class="form-control" :class="{ 'is-invalid': fieldErrors.name }" placeholder="เช่น ข้าวผัดกระเพรา" @input="fieldErrors.name = ''" />
            <span v-if="fieldErrors.name" class="field-error">{{ fieldErrors.name }}</span>
          </div>
          <div class="form-group">
            <label>หมวดหมู่ *</label>
            <select v-model="form.menu_category_id" class="form-control" :class="{ 'is-invalid': fieldErrors.menu_category_id }" @change="fieldErrors.menu_category_id = ''">
              <option value="">-- เลือกหมวดหมู่ --</option>
              <option v-for="cat in menuCategories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
            <span v-if="fieldErrors.menu_category_id" class="field-error">{{ fieldErrors.menu_category_id }}</span>
          </div>
          <div class="form-group">
            <label>รายละเอียด</label>
            <input v-model="form.description" class="form-control" placeholder="คำอธิบายเมนู" />
          </div>
          <div style="display:flex;gap:12px;">
            <div class="form-group" style="flex:1;">
              <label>ราคาขาย (บาท) *</label>
              <input v-model.number="form.price" class="form-control" :class="{ 'is-invalid': fieldErrors.price }" type="number" step="0.01" min="0" @input="fieldErrors.price = ''" />
              <span v-if="fieldErrors.price" class="field-error">{{ fieldErrors.price }}</span>
            </div>
            <div class="form-group" style="flex:1;">
              <label>ราคาต้นทุน (บาท) *</label>
              <input v-model.number="form.cost_price" class="form-control" :class="{ 'is-invalid': fieldErrors.cost_price }" type="number" step="0.01" min="0" @input="fieldErrors.cost_price = ''" />
              <span v-if="fieldErrors.cost_price" class="field-error">{{ fieldErrors.cost_price }}</span>
            </div>
          </div>
          <div class="form-group">
            <label>ลำดับการแสดง</label>
            <input v-model.number="form.sort_order" class="form-control" type="number" min="0" placeholder="เช่น 1, 2, 3" />
          </div>
          <div class="form-group" v-if="editing">
            <label>สถานะ</label>
            <select v-model="form.active" class="form-control">
              <option :value="true">เปิดขาย</option>
              <option :value="false">ปิด</option>
            </select>
          </div>

          <!-- Ingredients -->
          <div class="form-group">
            <label>วัตถุดิบที่ใช้</label>
            <div v-for="(row, idx) in form.ingredients" :key="idx" class="ingredient-row">
              <select v-model="row.ingredient_id" class="form-control">
                <option value="">-- เลือกวัตถุดิบ --</option>
                <option v-for="ing in allIngredients" :key="ing.id" :value="ing.id">
                  {{ ing.name }} ({{ ing.unit?.name }})
                </option>
              </select>
              <input v-model.number="row.quantity" class="form-control" type="number" step="0.01" min="0" placeholder="จำนวน" />
              <button class="btn btn-icon btn-outline" style="color:var(--color-danger)" @click="removeIngredient(idx)">
                <span class="material-symbols-outlined">close</span>
              </button>
            </div>
            <button class="btn btn-outline btn-sm" @click="addIngredient" style="margin-top: 4px;">
              <span class="material-symbols-outlined">add</span>
              เพิ่มวัตถุดิบ
            </button>
          </div>

          <div v-if="formError" class="alert alert-error">
            <span class="material-symbols-outlined">error</span>
            {{ formError }}
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showModal = false">ยกเลิก</button>
          <button class="btn btn-primary" @click="save" :disabled="saving">
            {{ saving ? 'กำลังบันทึก...' : 'บันทึก' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirm Modal -->
    <ConfirmModal
      v-if="showDeleteModal"
      title="ยืนยันการลบ"
      :message="`ต้องการลบเมนู ${deleteTarget?.name} หรือไม่?`"
      confirmText="ลบ"
      variant="danger"
      @confirm="doDelete"
      @cancel="showDeleteModal = false"
    />
  </div>
</template>

<script>
import { getMenuItems, createMenuItem, updateMenuItem, deleteMenuItem, getIngredients, getMenuCategories } from '../api/index.js'
import ConfirmModal from '../components/ConfirmModal.vue'

export default {
  components: { ConfirmModal },
  data() {
    return {
      items: [],
      allIngredients: [],
      menuCategories: [],
      search: '',
      filterCategory: '',
      showModal: false,
      showDeleteModal: false,
      editing: null,
      deleteTarget: null,
      saving: false,
      formError: '',
      fieldErrors: { name: '', price: '', cost_price: '', menu_category_id: '' },
      form: {
        name: '',
        description: '',
        price: 0,
        cost_price: 0,
        menu_category_id: '',
        sort_order: 0,
        active: true,
        ingredients: []
      }
    }
  },
  computed: {
    groupedItems() {
      const groups = {}
      for (const item of this.items) {
        const catName = item.menu_category?.name || 'ไม่มีหมวดหมู่'
        if (!groups[catName]) groups[catName] = []
        groups[catName].push(item)
      }
      return Object.keys(groups).sort().map(cat => ({ category: cat, items: groups[cat] }))
    }
  },
  async created() {
    await this.loadCategories()
    await this.load()
    await this.loadIngredients()
  },
  methods: {
    async load() {
      try {
        const params = {}
        if (this.search) params.search = this.search
        if (this.filterCategory) params.menu_category_id = this.filterCategory
        const { data } = await getMenuItems(params)
        this.items = data || []
      } catch (e) {
        console.error(e)
      }
    },
    async loadIngredients() {
      try {
        const { data } = await getIngredients()
        this.allIngredients = data || []
      } catch (e) {
        console.error(e)
      }
    },
    async loadCategories() {
      try {
        const { data } = await getMenuCategories()
        this.menuCategories = data || []
      } catch (e) {
        console.error(e)
      }
    },
    openAdd() {
      this.editing = null
      this.formError = ''
      this.fieldErrors = { name: '', price: '', cost_price: '', menu_category_id: '' }
      this.form = { name: '', description: '', price: 0, cost_price: 0, menu_category_id: '', sort_order: 0, active: true, ingredients: [] }
      this.showModal = true
    },
    openEdit(item) {
      this.editing = item.id
      this.formError = ''
      this.fieldErrors = { name: '', price: '', cost_price: '', menu_category_id: '' }
      this.form = {
        name: item.name,
        description: item.description || '',
        price: item.price,
        cost_price: item.cost_price || 0,
        menu_category_id: item.menu_category_id || '',
        sort_order: item.sort_order || 0,
        active: item.active,
        ingredients: (item.ingredients || []).map(i => ({
          ingredient_id: i.ingredient_id,
          quantity: i.quantity
        }))
      }
      this.showModal = true
    },
    addIngredient() {
      this.form.ingredients.push({ ingredient_id: '', quantity: 0 })
    },
    removeIngredient(idx) {
      this.form.ingredients.splice(idx, 1)
    },
    async save() {
      this.formError = ''
      this.fieldErrors = { name: '', price: '', cost_price: '', menu_category_id: '' }
      let valid = true
      if (!this.form.name.trim()) { this.fieldErrors.name = 'กรุณาระบุชื่อเมนู'; valid = false }
      if (!this.form.menu_category_id) { this.fieldErrors.menu_category_id = 'กรุณาเลือกหมวดหมู่'; valid = false }
      if (!this.form.price || this.form.price <= 0) { this.fieldErrors.price = 'กรุณาระบุราคาขายที่มากกว่า 0'; valid = false }
      if (this.form.cost_price === '' || this.form.cost_price === null || this.form.cost_price < 0) { this.fieldErrors.cost_price = 'กรุณาระบุราคาต้นทุน'; valid = false }
      if (!valid) return

      this.saving = true
      try {
        const payload = {
          name: this.form.name,
          description: this.form.description,
          price: this.form.price,
          cost_price: Number(this.form.cost_price),
          menu_category_id: Number(this.form.menu_category_id),
          sort_order: Number(this.form.sort_order) || 0,
          active: this.form.active,
          ingredients: this.form.ingredients.filter(i => i.ingredient_id && i.quantity > 0).map(i => ({
            ingredient_id: Number(i.ingredient_id),
            quantity: Number(i.quantity)
          }))
        }
        if (this.editing) {
          await updateMenuItem(this.editing, payload)
        } else {
          await createMenuItem(payload)
        }
        this.showModal = false
        await this.load()
      } catch (e) {
        this.formError = e.response?.data?.error || 'เกิดข้อผิดพลาด'
      } finally {
        this.saving = false
      }
    },
    confirmDelete(item) {
      this.deleteTarget = item
      this.showDeleteModal = true
    },
    async doDelete() {
      if (!this.deleteTarget) return
      try {
        await deleteMenuItem(this.deleteTarget.id)
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
    }
  }
}
</script>
