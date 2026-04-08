<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">จัดการหมวดหมู่</h1>
      <p class="page-subtitle">เพิ่ม แก้ไข และลบหมวดหมู่วัตถุดิบ</p>
    </div>

    <div class="card">
      <div class="card-header">
        <h2 class="section-title"><span class="material-symbols-outlined">folder</span> รายการหมวดหมู่</h2>
        <button v-if="isAdmin" class="btn btn-primary" @click="openModal()"><span class="material-symbols-outlined">add</span> เพิ่มหมวดหมู่</button>
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
              <th v-if="isAdmin">จัดการ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(cat, i) in categories" :key="cat.id">
              <td v-if="isAdmin" class="table-checkbox">
                <input type="checkbox" :value="cat.id" v-model="selectedIds" />
              </td>
              <td>{{ i + 1 }}</td>
              <td>{{ cat.name }}</td>
              <td v-if="isAdmin">
                <div class="btn-group">
                  <button class="btn btn-sm btn-warning" @click="openModal(cat)"><span class="material-symbols-outlined">edit</span> แก้ไข</button>
                  <button class="btn btn-sm btn-outline-danger" @click="remove(cat.id)"><span class="material-symbols-outlined">delete</span> ลบ</button>
                </div>
              </td>
            </tr>
            <tr v-if="categories.length === 0">
              <td :colspan="isAdmin ? 4 : 2" class="table-empty">
                <span class="table-empty-icon"><span class="material-symbols-outlined">folder_off</span></span>
                ยังไม่มีหมวดหมู่
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal -->
    <div class="modal-overlay" v-if="showModal" @click.self="cancelForm">
      <div class="modal">
        <h3 class="modal-title">{{ editing ? 'แก้ไขหมวดหมู่' : 'เพิ่มหมวดหมู่ใหม่' }}</h3>
        <div class="form-group">
          <label>ชื่อหมวดหมู่ <span class="required">*</span></label>
          <input v-model="form.name" class="form-control" :class="{ 'is-invalid': nameDuplicate }" placeholder="เช่น เนื้อสัตว์, ผัก, เครื่องปรุง" maxlength="100" @keyup.enter="confirmSave" @input="checkDuplicateName" />
          <div v-if="nameDuplicate" class="inline-error"><span class="material-symbols-outlined">error</span> ชื่อหมวดหมู่นี้มีอยู่แล้ว</div>
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="cancelForm">ยกเลิก</button>
          <button class="btn btn-primary" :disabled="nameDuplicate" @click="confirmSave">บันทึก</button>
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
import { getCategories, createCategory, updateCategory, deleteCategory, bulkDeleteCategories } from '../api/index.js'
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
      return this.categories.length > 0 && this.selectedIds.length === this.categories.length
    }
  },
  data() {
    return {
      categories: [],
      showModal: false,
      editing: null,
      form: { name: '' },
      selectedIds: [],
      confirmAction: null,
      nameDuplicate: false
    }
  },
  async mounted() {
    await this.load()
  },
  methods: {
    async load() {
      const { data } = await getCategories()
      this.categories = data
    },
    checkDuplicateName() {
      const name = this.form.name.trim().toLowerCase()
      if (!name) { this.nameDuplicate = false; return }
      const isDup = this.categories.some(c =>
        c.name.trim().toLowerCase() === name && (!this.editing || c.id !== this.editing.id)
      )
      this.nameDuplicate = isDup
      if (isDup) {
        this.confirmAction = {
          title: 'ชื่อซ้ำ',
          message: `หมวดหมู่ "${this.form.name.trim()}" มีอยู่ในระบบแล้ว กรุณาใช้ชื่ออื่น`,
          confirmText: 'ตกลง',
          variant: 'warning',
          onConfirm: () => { this.confirmAction = null }
        }
      }
    },
    openModal(cat = null) {
      this.editing = cat
      this.form = cat ? { name: cat.name } : { name: '' }
      this.nameDuplicate = false
      this.showModal = true
    },
    async save() {
      this.confirmAction = null
      if (this.editing) {
        await updateCategory(this.editing.id, this.form)
      } else {
        await createCategory(this.form)
      }
      this.showModal = false
      await this.load()
    },
    confirmSave() {
      if (!this.form.name.trim() || this.nameDuplicate) return
      this.confirmAction = {
        title: this.editing ? 'ยืนยันการแก้ไข' : 'ยืนยันการเพิ่ม',
        message: this.editing ? 'คุณต้องการบันทึกการแก้ไขหมวดหมู่นี้หรือไม่?' : `ต้องการเพิ่มหมวดหมู่ "${this.form.name}" หรือไม่?`,
        confirmText: 'บันทึก',
        variant: 'info',
        onConfirm: () => this.save()
      }
    },
    cancelForm() {
      if (this.form.name.trim()) {
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
        title: 'ลบหมวดหมู่',
        message: 'คุณต้องการลบหมวดหมู่นี้หรือไม่? การดำเนินการนี้ไม่สามารถย้อนกลับได้',
        confirmText: 'ลบ',
        variant: 'danger',
        onConfirm: async () => {
          this.confirmAction = null
          await deleteCategory(id)
          this.selectedIds = this.selectedIds.filter(sid => sid !== id)
          await this.load()
        }
      }
    },
    confirmBulkDelete() {
      this.confirmAction = {
        title: 'ลบหมวดหมู่ที่เลือก',
        message: `คุณต้องการลบหมวดหมู่ ${this.selectedIds.length} รายการหรือไม่? การดำเนินการนี้ไม่สามารถย้อนกลับได้`,
        confirmText: `ลบ ${this.selectedIds.length} รายการ`,
        variant: 'danger',
        onConfirm: async () => {
          this.confirmAction = null
          await bulkDeleteCategories(this.selectedIds)
          this.selectedIds = []
          await this.load()
        }
      }
    },
    toggleAll(e) {
      this.selectedIds = e.target.checked ? this.categories.map(c => c.id) : []
    }
  }
}
</script>
