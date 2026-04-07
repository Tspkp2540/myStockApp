<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">จัดการวัตถุดิบ</h1>
      <p class="page-subtitle">เพิ่ม แก้ไข และดูสถานะสต็อควัตถุดิบ</p>
    </div>

    <div class="card">
      <div class="toolbar">
        <div class="toolbar-filters">
          <input
            v-model="search"
            class="form-control"
            placeholder="🔍 ค้นหาวัตถุดิบ..."
            @input="load"
          />
          <select v-model="filterCategory" class="form-control" @change="load">
            <option :value="''">ทุกหมวดหมู่</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
        <button v-if="isAdmin" class="btn btn-primary" @click="openModal()">+ เพิ่มวัตถุดิบ</button>
      </div>

      <div v-if="isAdmin && selectedIds.length > 0" class="bulk-bar">
        <span class="bulk-bar-text">เลือกแล้ว {{ selectedIds.length }} รายการ</span>
        <button class="btn btn-sm btn-danger" @click="confirmBulkDelete">ลบที่เลือก</button>
      </div>

      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th v-if="isAdmin" class="table-checkbox">
                <input type="checkbox" :checked="allSelected" @change="toggleAll" />
              </th>
              <th>#</th>
              <th>ชื่อ</th>
              <th>หมวดหมู่</th>
              <th>คงเหลือ</th>
              <th>หน่วย</th>
              <th>ขั้นต่ำ</th>
              <th>สถานะ</th>
              <th v-if="isAdmin">จัดการ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(ing, i) in ingredients" :key="ing.id">
              <td v-if="isAdmin" class="table-checkbox">
                <input type="checkbox" :value="ing.id" v-model="selectedIds" />
              </td>
              <td>{{ i + 1 }}</td>
              <td>{{ ing.name }}</td>
              <td>{{ ing.category?.name || '-' }}</td>
              <td>{{ ing.quantity }}</td>
              <td>{{ ing.unit?.name || '-' }}</td>
              <td>{{ ing.min_stock }}</td>
              <td>
                <span v-if="ing.quantity <= ing.min_stock" class="badge badge-low">⚠ ต่ำ</span>
                <span v-else class="badge badge-ok">✓ ปกติ</span>
              </td>
              <td v-if="isAdmin">
                <div class="btn-group">
                  <button class="btn btn-sm btn-warning" @click="openModal(ing)">แก้ไข</button>
                  <button class="btn btn-sm btn-outline-danger" @click="remove(ing.id)">ลบ</button>
                </div>
              </td>
            </tr>
            <tr v-if="ingredients.length === 0">
              <td :colspan="isAdmin ? 9 : 7" class="table-empty">
                <span class="table-empty-icon">📦</span>
                ยังไม่มีวัตถุดิบ — กดปุ่ม "+ เพิ่มวัตถุดิบ" เพื่อเริ่มต้น
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal -->
    <div class="modal-overlay" v-if="showModal" @click.self="cancelForm">
      <div class="modal">
        <h3 class="modal-title">{{ editing ? 'แก้ไขวัตถุดิบ' : 'เพิ่มวัตถุดิบใหม่' }}</h3>
        <div class="form-group">
          <label>ชื่อวัตถุดิบ <span class="required">*</span></label>
          <input v-model="form.name" class="form-control" placeholder="เช่น หมูสามชั้น, ผักกาด" />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>หมวดหมู่</label>
            <div class="input-with-addon">
              <select v-model="form.category_id" class="form-control">
                <option :value="null">-- เลือก --</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
              <button class="btn btn-icon btn-success" type="button" @click="showAddCategory = true" title="เพิ่มหมวดหมู่ใหม่">+</button>
            </div>
          </div>
          <div class="form-group">
            <label>หน่วย <span class="required">*</span></label>
            <div class="input-with-addon">
              <select v-model="form.unit_id" class="form-control">
                <option :value="null">-- เลือก --</option>
                <option v-for="u in units" :key="u.id" :value="u.id">{{ u.name }}</option>
              </select>
              <button class="btn btn-icon btn-success" type="button" @click="showAddUnit = true" title="เพิ่มหน่วยใหม่">+</button>
            </div>
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>จำนวนเริ่มต้น</label>
            <input v-model.number="form.quantity" class="form-control" type="number" min="0" step="0.1" />
          </div>
          <div class="form-group">
            <label>สต็อคขั้นต่ำ</label>
            <input v-model.number="form.min_stock" class="form-control" type="number" min="0" step="0.1" />
          </div>
        </div>
        <div class="form-group" v-if="!editing && form.quantity > 0">
          <label>ราคาต่อหน่วย (ฺ)</label>
          <input v-model.number="form.price" class="form-control" type="number" min="0" step="0.01" placeholder="0.00" />
        </div>
        <div class="form-group" v-if="!editing && form.quantity > 0">
          <label>เหตุผลในการเพิ่มสต็อค <span class="required">*</span></label>
          <input v-model="form.note" class="form-control" placeholder="เช่น สต็อคเปิดร้าน, รับสินค้าจากซัพพลายเออร์" />
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="cancelForm">ยกเลิก</button>
          <button class="btn btn-primary" @click="confirmSave">บันทึก</button>
        </div>
      </div>
    </div>

    <!-- Quick Add Category Modal -->
    <div class="modal-overlay" v-if="showAddCategory" @click.self="cancelAddCategory">
      <div class="modal">
        <h3 class="modal-title">เพิ่มหมวดหมู่ใหม่</h3>
        <div class="form-group">
          <label>ชื่อหมวดหมู่ <span class="required">*</span></label>
          <input v-model="newCategoryName" class="form-control" placeholder="เช่น เนื้อสัตว์, ผัก" @keyup.enter="addCategory" />
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="cancelAddCategory">ยกเลิก</button>
          <button class="btn btn-primary" @click="confirmAddCategory">เพิ่ม</button>
        </div>
      </div>
    </div>

    <!-- Quick Add Unit Modal -->
    <div class="modal-overlay" v-if="showAddUnit" @click.self="cancelAddUnit">
      <div class="modal">
        <h3 class="modal-title">เพิ่มหน่วยใหม่</h3>
        <div class="form-group">
          <label>ชื่อหน่วย <span class="required">*</span></label>
          <input v-model="newUnitName" class="form-control" placeholder="เช่น กก., ลิตร, ชิ้น" @keyup.enter="addUnit" />
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="cancelAddUnit">ยกเลิก</button>
          <button class="btn btn-primary" @click="confirmAddUnit">เพิ่ม</button>
        </div>
      </div>
    </div>

    <ConfirmModal
      v-if="confirmAction"
      :title="confirmAction.title"
      :message="confirmAction.message"
      :confirmText="confirmAction.confirmText"
      :variant="confirmAction.variant"
      @confirm="confirmAction.onConfirm"
      @cancel="confirmAction = null"
    />
  </div>
</template>

<script>
import { getIngredients, createIngredient, updateIngredient, deleteIngredient, bulkDeleteIngredients, getCategories, createCategory, getUnits, createUnit } from '../api/index.js'
import { useAuth } from '../composables/useAuth.js'
import ConfirmModal from '../components/ConfirmModal.vue'

export default {
  components: { ConfirmModal },
  computed: {
    isAdmin() {
      const { isAdmin } = useAuth()
      return isAdmin()
    },
    allSelected() {
      return this.ingredients.length > 0 && this.selectedIds.length === this.ingredients.length
    }
  },
  data() {
    return {
      ingredients: [],
      categories: [],
      units: [],
      search: '',
      filterCategory: '',
      showModal: false,
      editing: null,
      form: { name: '', category_id: null, unit_id: null, quantity: 0, min_stock: 0, note: '', price: 0 },
      showAddCategory: false,
      showAddUnit: false,
      newCategoryName: '',
      newUnitName: '',
      selectedIds: [],
      confirmAction: null
    }
  },
  async mounted() {
    const [catRes, unitRes] = await Promise.all([getCategories(), getUnits()])
    this.categories = catRes.data
    this.units = unitRes.data
    await this.load()
  },
  methods: {
    async load() {
      const params = {}
      if (this.search) params.search = this.search
      if (this.filterCategory) params.category_id = this.filterCategory
      const { data } = await getIngredients(params)
      this.ingredients = data
    },
    openModal(ing = null) {
      this.editing = ing
      this.form = ing
        ? { name: ing.name, category_id: ing.category_id || null, unit_id: ing.unit_id || null, quantity: ing.quantity, min_stock: ing.min_stock, note: '', price: 0 }
        : { name: '', category_id: null, unit_id: null, quantity: 0, min_stock: 0, note: '', price: 0 }
      this.showModal = true
    },
    isFormDirty() {
      return this.form.name.trim() !== '' || this.form.category_id !== null || this.form.unit_id !== null || this.form.quantity > 0 || this.form.min_stock > 0
    },
    confirmSave() {
      if (!this.form.name.trim() || !this.form.unit_id) return
      if (!this.editing && this.form.quantity > 0 && !this.form.note.trim()) return
      this.confirmAction = {
        title: this.editing ? 'ยืนยันการแก้ไข' : 'ยืนยันการเพิ่ม',
        message: this.editing ? 'คุณต้องการบันทึกการแก้ไขวัตถุดิบนี้หรือไม่?' : 'คุณต้องการเพิ่มวัตถุดิบนี้หรือไม่?',
        confirmText: 'บันทึก',
        variant: 'info',
        onConfirm: () => this.save()
      }
    },
    async save() {
      this.confirmAction = null
      const payload = {
        name: this.form.name,
        category_id: Number(this.form.category_id) || 0,
        unit_id: Number(this.form.unit_id),
        quantity: Number(this.form.quantity) || 0,
        min_stock: Number(this.form.min_stock) || 0
      }
      if (this.editing) {
        await updateIngredient(this.editing.id, payload)
      } else {
        if (this.form.note.trim()) payload.note = this.form.note
        if (this.form.price > 0) payload.price = this.form.price
        await createIngredient(payload)
      }
      this.showModal = false
      await this.load()
    },
    cancelForm() {
      if (this.isFormDirty()) {
        this.confirmAction = {
          title: 'ยกเลิกการกรอกข้อมูล',
          message: 'คุณมีข้อมูลที่ยังไม่ได้บันทึก ต้องการยกเลิกหรือไม่?',
          confirmText: 'ยืนยัน',
          variant: 'warning',
          onConfirm: () => { this.confirmAction = null; this.showModal = false }
        }
      } else {
        this.showModal = false
      }
    },
    async remove(id) {
      this.confirmAction = {
        title: 'ลบวัตถุดิบ',
        message: 'คุณต้องการลบวัตถุดิบนี้หรือไม่? การดำเนินการนี้ไม่สามารถย้อนกลับได้',
        confirmText: 'ลบ',
        variant: 'danger',
        onConfirm: async () => {
          this.confirmAction = null
          await deleteIngredient(id)
          this.selectedIds = this.selectedIds.filter(sid => sid !== id)
          await this.load()
        }
      }
    },
    confirmBulkDelete() {
      this.confirmAction = {
        title: 'ลบวัตถุดิบที่เลือก',
        message: `คุณต้องการลบวัตถุดิบ ${this.selectedIds.length} รายการหรือไม่? การดำเนินการนี้ไม่สามารถย้อนกลับได้`,
        confirmText: `ลบ ${this.selectedIds.length} รายการ`,
        variant: 'danger',
        onConfirm: async () => {
          this.confirmAction = null
          await bulkDeleteIngredients(this.selectedIds)
          this.selectedIds = []
          await this.load()
        }
      }
    },
    toggleAll(e) {
      this.selectedIds = e.target.checked ? this.ingredients.map(i => i.id) : []
    },
    confirmAddCategory() {
      if (!this.newCategoryName.trim()) return
      this.confirmAction = {
        title: 'ยืนยันเพิ่มหมวดหมู่',
        message: `ต้องการเพิ่มหมวดหมู่ "${this.newCategoryName}" หรือไม่?`,
        confirmText: 'เพิ่ม',
        variant: 'info',
        onConfirm: () => this.addCategory()
      }
    },
    async addCategory() {
      this.confirmAction = null
      const { data } = await createCategory({ name: this.newCategoryName })
      this.categories.push(data)
      this.form.category_id = data.id
      this.newCategoryName = ''
      this.showAddCategory = false
    },
    cancelAddCategory() {
      if (this.newCategoryName.trim()) {
        this.confirmAction = {
          title: 'ยกเลิกการเพิ่มหมวดหมู่',
          message: 'คุณมีข้อมูลที่ยังไม่ได้บันทึก ต้องการยกเลิกหรือไม่?',
          confirmText: 'ยืนยีน',
          variant: 'warning',
          onConfirm: () => { this.confirmAction = null; this.newCategoryName = ''; this.showAddCategory = false }
        }
      } else {
        this.showAddCategory = false
      }
    },
    confirmAddUnit() {
      if (!this.newUnitName.trim()) return
      this.confirmAction = {
        title: 'ยืนยันเพิ่มหน่วย',
        message: `ต้องการเพิ่มหน่วย "${this.newUnitName}" หรือไม่?`,
        confirmText: 'เพิ่ม',
        variant: 'info',
        onConfirm: () => this.addUnit()
      }
    },
    async addUnit() {
      this.confirmAction = null
      const { data } = await createUnit({ name: this.newUnitName })
      this.units.push(data)
      this.form.unit_id = data.id
      this.newUnitName = ''
      this.showAddUnit = false
    },
    cancelAddUnit() {
      if (this.newUnitName.trim()) {
        this.confirmAction = {
          title: 'ยกเลิกการเพิ่มหน่วย',
          message: 'คุณมีข้อมูลที่ยังไม่ได้บันทึก ต้องการยกเลิกหรือไม่?',
          confirmText: 'ยกเลิก',
          variant: 'warning',
          onConfirm: () => { this.confirmAction = null; this.newUnitName = ''; this.showAddUnit = false }
        }
      } else {
        this.showAddUnit = false
      }
    }
  }
}
</script>
