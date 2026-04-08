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
            placeholder="ค้นหาวัตถุดิบ..."
            maxlength="100"
            @input="load"
          />
          <select v-model="filterCategory" class="form-control" @change="load">
            <option :value="''">ทุกหมวดหมู่</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
        <button v-if="isAdmin" class="btn btn-primary" @click="openModal()"><span class="material-symbols-outlined">add</span> เพิ่มวัตถุดิบ</button>
      </div>

      <div v-if="isAdmin && selectedIds.length > 0" class="bulk-bar">
        <span class="bulk-bar-text">เลือกแล้ว {{ selectedIds.length }} รายการ</span>
        <button class="btn btn-sm btn-danger" @click="confirmBulkDelete">ลบที่เลือก</button>
      </div>

      <div class="table-wrapper">
        <table v-for="group in groupedIngredients" :key="group.category" class="grouped-table">
          <thead>
            <tr class="category-header-row">
              <th :colspan="isAdmin ? 9 : 7" class="category-group-header">
                <span class="material-symbols-outlined">folder</span> {{ group.category }}
                <span class="category-count">({{ group.items.length }} รายการ)</span>
              </th>
            </tr>
            <tr>
              <th v-if="isAdmin" class="table-checkbox">
                <input type="checkbox" :checked="isGroupAllSelected(group.items)" @change="toggleGroup($event, group.items)" />
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
            <tr v-for="(ing, i) in group.items" :key="ing.id">
              <td v-if="isAdmin" class="table-checkbox">
                <input type="checkbox" :value="ing.id" v-model="selectedIds" />
              </td>
              <td>{{ i + 1 }}</td>
              <td>
                <router-link :to="`/ingredients/${ing.id}`" class="ingredient-link">{{ ing.name }}</router-link>
              </td>
              <td>{{ ing.category?.name || '-' }}</td>
              <td>{{ ing.quantity }}</td>
              <td>{{ ing.unit?.name || '-' }}</td>
              <td>{{ ing.min_stock }}</td>
              <td>
                <span v-if="ing.quantity <= ing.min_stock" class="badge badge-low"><span class="material-symbols-outlined">error</span> ต่ำ</span>
                <span v-else class="badge badge-ok"><span class="material-symbols-outlined">check_circle</span> ปกติ</span>
              </td>
              <td v-if="isAdmin">
                <div class="btn-group">
                  <router-link :to="`/ingredients/${ing.id}`" class="btn btn-sm btn-primary"><span class="material-symbols-outlined">bar_chart</span> ดูรายละเอียด</router-link>
                  <button class="btn btn-sm btn-warning" @click="openModal(ing)"><span class="material-symbols-outlined">edit</span> แก้ไข</button>
                  <button class="btn btn-sm btn-outline-danger" @click="remove(ing.id)"><span class="material-symbols-outlined">delete</span> ลบ</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-if="ingredients.length === 0" class="table-empty">
          <span class="table-empty-icon"><span class="material-symbols-outlined">inventory_2</span></span>
          ยังไม่มีวัตถุดิบ — กดปุ่ม "+ เพิ่มวัตถุดิบ" เพื่อเริ่มต้น
        </div>
      </div>
    </div>

    <!-- Modal -->
    <div class="modal-overlay" v-if="showModal" @click.self="cancelForm">
      <div class="modal">
        <h3 class="modal-title">{{ editing ? 'แก้ไขวัตถุดิบ' : 'เพิ่มวัตถุดิบใหม่' }}</h3>
        <div class="form-group">
          <label>ชื่อวัตถุดิบ <span class="required">*</span></label>
          <input v-model="form.name" class="form-control" :class="{ 'is-invalid': ingredientDuplicate }" placeholder="เช่น หมูสามชั้น, ผักกาด" maxlength="100" @input="checkIngredientDuplicate" />
          <div v-if="ingredientDuplicate" class="inline-error"><span class="material-symbols-outlined">error</span> ชื่อวัตถุดิบนี้มีอยู่แล้ว</div>
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
            <input v-model.number="form.quantity" class="form-control" type="number" min="0" step="0.1" @input="limitNumber($event, 'quantity')" />
          </div>
          <div class="form-group">
            <label>สต็อคขั้นต่ำ</label>
            <input v-model.number="form.min_stock" class="form-control" type="number" min="0" step="0.1" @input="limitNumber($event, 'min_stock')" />
          </div>
        </div>
        <div class="form-group" v-if="!editing && form.quantity > 0">
          <label>ราคาต่อหน่วย (ฺ)</label>
          <input v-model.number="form.price" class="form-control" type="number" min="0" step="0.01" placeholder="0.00" @input="limitNumber($event, 'price')" />
        </div>
        <div class="form-group" v-if="!editing && form.quantity > 0">
          <label>วันที่เพิ่ม <small>(ถ้าไม่กรอกจะใช้วันปัจจุบัน)</small></label>
          <input v-model="form.date" class="form-control" type="date" />
        </div>
        <div class="form-group" v-if="!editing && form.quantity > 0">
          <label>เหตุผลในการเพิ่มสต็อค <span class="required">*</span></label>
          <input v-model="form.note" class="form-control" placeholder="เช่น สต็อคเปิดร้าน, รับสินค้าจากซัพพลายเออร์" maxlength="100" />
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="cancelForm">ยกเลิก</button>
          <button class="btn btn-primary" :disabled="ingredientDuplicate" @click="confirmSave">บันทึก</button>
        </div>
      </div>
    </div>

    <!-- Quick Add Category Modal -->
    <div class="modal-overlay" v-if="showAddCategory" @click.self="cancelAddCategory">
      <div class="modal">
        <h3 class="modal-title">เพิ่มหมวดหมู่ใหม่</h3>
        <div class="form-group">
          <label>ชื่อหมวดหมู่ <span class="required">*</span></label>
          <input v-model="newCategoryName" class="form-control" :class="{ 'is-invalid': categoryDuplicate }" placeholder="เช่น เนื้อสัตว์, ผัก" maxlength="100" @keyup.enter="confirmAddCategory" @input="checkCategoryDuplicate" />
          <div v-if="categoryDuplicate" class="inline-error"><span class="material-symbols-outlined">error</span> ชื่อหมวดหมู่นี้มีอยู่แล้ว</div>
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="cancelAddCategory">ยกเลิก</button>
          <button class="btn btn-primary" :disabled="categoryDuplicate" @click="confirmAddCategory">เพิ่ม</button>
        </div>
      </div>
    </div>

    <!-- Quick Add Unit Modal -->
    <div class="modal-overlay" v-if="showAddUnit" @click.self="cancelAddUnit">
      <div class="modal">
        <h3 class="modal-title">เพิ่มหน่วยใหม่</h3>
        <div class="form-group">
          <label>ชื่อหน่วย <span class="required">*</span></label>
          <input v-model="newUnitName" class="form-control" :class="{ 'is-invalid': unitDuplicate }" placeholder="เช่น กก., ลิตร, ชิ้น" maxlength="100" @keyup.enter="confirmAddUnit" @input="checkUnitDuplicate" />
          <div v-if="unitDuplicate" class="inline-error"><span class="material-symbols-outlined">error</span> ชื่อหน่วยนี้มีอยู่แล้ว</div>
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="cancelAddUnit">ยกเลิก</button>
          <button class="btn btn-primary" :disabled="unitDuplicate" @click="confirmAddUnit">เพิ่ม</button>
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
    },
    groupedIngredients() {
      const groups = {}
      this.ingredients.forEach(ing => {
        const cat = ing.category?.name || 'ไม่มีหมวดหมู่'
        if (!groups[cat]) groups[cat] = []
        groups[cat].push(ing)
      })
      return Object.keys(groups).sort().map(cat => ({
        category: cat,
        items: groups[cat]
      }))
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
      form: { name: '', category_id: null, unit_id: null, quantity: 0, min_stock: 0, note: '', price: 0, date: '' },
      showAddCategory: false,
      showAddUnit: false,
      newCategoryName: '',
      newUnitName: '',
      selectedIds: [],
      confirmAction: null,
      ingredientDuplicate: false,
      categoryDuplicate: false,
      unitDuplicate: false
    }
  },
  async mounted() {
    const [catRes, unitRes] = await Promise.all([getCategories(), getUnits()])
    this.categories = catRes.data
    this.units = unitRes.data
    await this.load()
  },
  methods: {
    limitNumber(e, field) {
      const raw = e.target.value.replace(/[^0-9.]/g, '')
      const parts = raw.split('.')
      let clean = parts[0].slice(0, 7)
      if (parts.length > 1) clean += '.' + parts[1].slice(0, 2)
      e.target.value = clean
      this.form[field] = clean === '' ? 0 : Number(clean)
    },
    checkIngredientDuplicate() {
      const name = this.form.name.trim().toLowerCase()
      if (!name) { this.ingredientDuplicate = false; return }
      const isDup = this.ingredients.some(i =>
        i.name.trim().toLowerCase() === name && (!this.editing || i.id !== this.editing.id)
      )
      this.ingredientDuplicate = isDup
      if (isDup) {
        this.confirmAction = {
          title: 'ชื่อซ้ำ',
          message: `วัตถุดิบ "${this.form.name.trim()}" มีอยู่ในระบบแล้ว กรุณาใช้ชื่ออื่น`,
          confirmText: 'ตกลง',
          variant: 'warning',
          onConfirm: () => { this.confirmAction = null }
        }
      }
    },
    checkCategoryDuplicate() {
      const name = this.newCategoryName.trim().toLowerCase()
      if (!name) { this.categoryDuplicate = false; return }
      const isDup = this.categories.some(c => c.name.trim().toLowerCase() === name)
      this.categoryDuplicate = isDup
      if (isDup) {
        this.confirmAction = {
          title: 'ชื่อซ้ำ',
          message: `หมวดหมู่ "${this.newCategoryName.trim()}" มีอยู่ในระบบแล้ว กรุณาใช้ชื่ออื่น`,
          confirmText: 'ตกลง',
          variant: 'warning',
          onConfirm: () => { this.confirmAction = null }
        }
      }
    },
    checkUnitDuplicate() {
      const name = this.newUnitName.trim().toLowerCase()
      if (!name) { this.unitDuplicate = false; return }
      const isDup = this.units.some(u => u.name.trim().toLowerCase() === name)
      this.unitDuplicate = isDup
      if (isDup) {
        this.confirmAction = {
          title: 'ชื่อซ้ำ',
          message: `หน่วย "${this.newUnitName.trim()}" มีอยู่ในระบบแล้ว กรุณาใช้ชื่ออื่น`,
          confirmText: 'ตกลง',
          variant: 'warning',
          onConfirm: () => { this.confirmAction = null }
        }
      }
    },
    async load() {
      const params = {}
      if (this.search) params.search = this.search
      if (this.filterCategory) params.category_id = this.filterCategory
      const { data } = await getIngredients(params)
      this.ingredients = data
    },
    openModal(ing = null) {
      this.editing = ing
      this.ingredientDuplicate = false
      this.form = ing
        ? { name: ing.name, category_id: ing.category_id || null, unit_id: ing.unit_id || null, quantity: ing.quantity, min_stock: ing.min_stock, note: '', price: 0, date: '' }
        : { name: '', category_id: null, unit_id: null, quantity: 0, min_stock: 0, note: '', price: 0, date: '' }
      this.showModal = true
    },
    isFormDirty() {
      return this.form.name.trim() !== '' || this.form.category_id !== null || this.form.unit_id !== null || this.form.quantity > 0 || this.form.min_stock > 0
    },
    confirmSave() {
      if (!this.form.name.trim() || !this.form.unit_id || this.ingredientDuplicate) return
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
        if (this.form.date) payload.date = this.form.date
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
    isGroupAllSelected(items) {
      return items.length > 0 && items.every(i => this.selectedIds.includes(i.id))
    },
    toggleGroup(e, items) {
      const ids = items.map(i => i.id)
      if (e.target.checked) {
        const newIds = ids.filter(id => !this.selectedIds.includes(id))
        this.selectedIds = [...this.selectedIds, ...newIds]
      } else {
        this.selectedIds = this.selectedIds.filter(id => !ids.includes(id))
      }
    },
    confirmAddCategory() {
      if (!this.newCategoryName.trim() || this.categoryDuplicate) return
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
      if (!this.newUnitName.trim() || this.unitDuplicate) return
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
