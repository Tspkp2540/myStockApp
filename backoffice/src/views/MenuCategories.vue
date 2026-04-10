<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">หมวดหมู่เมนู</h1>
        <p class="page-subtitle">จัดการหมวดหมู่สำหรับจัดกลุ่มเมนูอาหาร</p>
      </div>
      <button class="btn btn-primary" @click="openAdd">
        <span class="material-symbols-outlined">add</span> เพิ่มหมวดหมู่
      </button>
    </div>

    <div class="card">
      <div class="table-responsive">
        <table class="table">
          <thead>
            <tr>
              <th style="width:60px;">#</th>
              <th>ชื่อหมวดหมู่</th>
              <th>วันที่สร้าง</th>
              <th class="text-center" style="width:140px;">จัดการ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="4" class="text-center text-muted" style="padding:32px;">กำลังโหลด...</td>
            </tr>
            <tr v-else-if="categories.length === 0">
              <td colspan="4" class="text-center text-muted" style="padding:32px;">ยังไม่มีหมวดหมู่</td>
            </tr>
            <tr v-for="(cat, idx) in categories" :key="cat.id">
              <td class="text-muted">{{ idx + 1 }}</td>
              <td style="font-weight:600;">{{ cat.name }}</td>
              <td class="text-muted">{{ formatDate(cat.created_at) }}</td>
              <td class="text-center">
                <div class="action-buttons">
                  <button class="btn-icon" title="แก้ไข" @click="openEdit(cat)">
                    <span class="material-symbols-outlined">edit</span>
                  </button>
                  <button class="btn-icon btn-icon-danger" title="ลบ" @click="confirmDelete(cat)">
                    <span class="material-symbols-outlined">delete</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal" style="max-width:440px;">
        <div class="modal-header">
          <h3>{{ editId ? 'แก้ไขหมวดหมู่' : 'เพิ่มหมวดหมู่' }}</h3>
          <button class="modal-close" @click="showModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label for="catName">ชื่อหมวดหมู่ *</label>
            <input id="catName" v-model="form.name" type="text" class="form-control"
              :class="{ 'is-invalid': fieldErrors.name }" placeholder="เช่น อาหารจานเดียว"
              @keyup.enter="save" />
            <span v-if="fieldErrors.name" class="field-error">{{ fieldErrors.name }}</span>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showModal = false">ยกเลิก</button>
          <button class="btn btn-primary" :disabled="saving" @click="save">
            <span v-if="saving" class="material-symbols-outlined spin">progress_activity</span>
            {{ editId ? 'บันทึก' : 'เพิ่ม' }}
          </button>
        </div>
      </div>
    </div>

    <ConfirmModal v-if="showDeleteConfirm" title="ยืนยันการลบ"
      :message="`ต้องการลบหมวดหมู่ &quot;${deleteTarget?.name}&quot; ใช่หรือไม่?`"
      confirm-text="ลบ" type="danger" @confirm="doDelete" @cancel="showDeleteConfirm = false" />
  </div>
</template>

<script>
import { getMenuCategories, createMenuCategory, updateMenuCategory, deleteMenuCategory } from '../api/index.js'
import ConfirmModal from '../components/ConfirmModal.vue'

export default {
  name: 'MenuCategories',
  components: { ConfirmModal },
  data() {
    return {
      categories: [],
      loading: false,
      showModal: false,
      saving: false,
      editId: null,
      form: { name: '' },
      fieldErrors: {},
      showDeleteConfirm: false,
      deleteTarget: null,
    }
  },
  mounted() {
    this.fetchData()
  },
  methods: {
    async fetchData() {
      this.loading = true
      try {
        const res = await getMenuCategories()
        this.categories = res.data
      } catch (e) {
        console.error(e)
      } finally {
        this.loading = false
      }
    },
    openAdd() {
      this.editId = null
      this.form = { name: '' }
      this.fieldErrors = {}
      this.showModal = true
      this.$nextTick(() => document.getElementById('catName')?.focus())
    },
    openEdit(cat) {
      this.editId = cat.id
      this.form = { name: cat.name }
      this.fieldErrors = {}
      this.showModal = true
      this.$nextTick(() => document.getElementById('catName')?.focus())
    },
    async save() {
      this.fieldErrors = {}
      if (!this.form.name.trim()) {
        this.fieldErrors.name = 'กรุณากรอกชื่อหมวดหมู่'
        return
      }
      this.saving = true
      try {
        if (this.editId) {
          await updateMenuCategory(this.editId, this.form)
        } else {
          await createMenuCategory(this.form)
        }
        this.showModal = false
        await this.fetchData()
      } catch (e) {
        if (e.response?.data?.error) {
          this.fieldErrors.name = e.response.data.error
        }
      } finally {
        this.saving = false
      }
    },
    confirmDelete(cat) {
      this.deleteTarget = cat
      this.showDeleteConfirm = true
    },
    async doDelete() {
      try {
        await deleteMenuCategory(this.deleteTarget.id)
        this.showDeleteConfirm = false
        await this.fetchData()
      } catch (e) {
        alert(e.response?.data?.error || 'เกิดข้อผิดพลาด')
        this.showDeleteConfirm = false
      }
    },
    formatDate(d) {
      if (!d) return '-'
      return new Date(d).toLocaleDateString('th-TH', { year: 'numeric', month: 'short', day: 'numeric' })
    },
  },
}
</script>
