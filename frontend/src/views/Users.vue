<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">จัดการผู้ใช้งาน</h1>
      <p class="page-subtitle">เพิ่ม แก้ไข สิทธิ์ผู้ใช้ระบบ (เฉพาะผู้ดูแล)</p>
    </div>

    <div class="card">
      <div class="card-header">
        <h2 class="section-title">👤 รายชื่อผู้ใช้</h2>
        <button class="btn btn-primary" @click="openModal()">+ เพิ่มผู้ใช้</button>
      </div>
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>ชื่อผู้ใช้</th>
              <th>ชื่อ-นามสกุล</th>
              <th>สิทธิ์</th>
              <th>สถานะ</th>
              <th>จัดการ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(u, i) in users" :key="u.id">
              <td>{{ i + 1 }}</td>
              <td>{{ u.username }}</td>
              <td>{{ u.full_name || '-' }}</td>
              <td>
                <span :class="['badge', u.role === 'admin' ? 'badge-warning' : 'badge-info']">
                  {{ u.role === 'admin' ? '🔑 ผู้ดูแล' : '👤 พนักงาน' }}
                </span>
              </td>
              <td>
                <span :class="['badge', u.active ? 'badge-ok' : 'badge-low']">
                  {{ u.active ? '✓ ใช้งาน' : '✗ ระงับ' }}
                </span>
              </td>
              <td>
                <div class="btn-group">
                  <button class="btn btn-sm btn-warning" @click="openModal(u)">แก้ไข</button>
                  <button class="btn btn-sm btn-outline-danger" @click="remove(u.id)" v-if="u.id !== currentUser?.id">ลบ</button>
                </div>
              </td>
            </tr>
            <tr v-if="users.length === 0">
              <td colspan="6" class="table-empty">
                <span class="table-empty-icon">👤</span>
                ยังไม่มีผู้ใช้
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal -->
    <div class="modal-overlay" v-if="showModal" @click.self="cancelForm">
      <div class="modal">
        <h3 class="modal-title">{{ editing ? 'แก้ไขผู้ใช้' : 'เพิ่มผู้ใช้ใหม่' }}</h3>
        <div class="form-group" v-if="!editing">
          <label>ชื่อผู้ใช้ <span class="required">*</span></label>
          <input v-model="form.username" class="form-control" placeholder="username" />
        </div>
        <div class="form-group">
          <label>ชื่อ-นามสกุล</label>
          <input v-model="form.full_name" class="form-control" placeholder="เช่น สมชาย ใจดี" />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>{{ editing ? 'รหัสผ่านใหม่ (เว้นว่างถ้าไม่เปลี่ยน)' : 'รหัสผ่าน *' }}</label>
            <input v-model="form.password" class="form-control" type="password" placeholder="••••••••" />
          </div>
          <div class="form-group">
            <label>สิทธิ์</label>
            <select v-model="form.role" class="form-control">
              <option value="employee">👤 พนักงาน</option>
              <option value="admin">🔑 ผู้ดูแลระบบ</option>
            </select>
          </div>
        </div>
        <div class="form-group" v-if="editing">
          <label class="checkbox-label">
            <input type="checkbox" v-model="form.active" />
            เปิดใช้งานบัญชี
          </label>
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="cancelForm">ยกเลิก</button>
          <button class="btn btn-primary" @click="confirmSave">บันทึก</button>
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
import { getUsers, createUser, updateUser, deleteUser } from '../api/index.js'
import { useAuth } from '../composables/useAuth.js'
import ConfirmModal from '../components/ConfirmModal.vue'

export default {
  components: { ConfirmModal },
  data() {
    return {
      users: [],
      showModal: false,
      editing: null,
      form: { username: '', password: '', full_name: '', role: 'employee', active: true },
      currentUser: null,
      confirmAction: null
    }
  },
  async mounted() {
    const { state } = useAuth()
    this.currentUser = state.user
    await this.load()
  },
  methods: {
    async load() {
      const { data } = await getUsers()
      this.users = data
    },
    openModal(user = null) {
      this.editing = user
      if (user) {
        this.form = { full_name: user.full_name, role: user.role, active: user.active, password: '' }
      } else {
        this.form = { username: '', password: '', full_name: '', role: 'employee', active: true }
      }
      this.showModal = true
    },
    async save() {
      this.confirmAction = null
      if (!this.editing && (!this.form.username.trim() || !this.form.password)) return
      try {
        if (this.editing) {
          const payload = { full_name: this.form.full_name, role: this.form.role, active: this.form.active }
          if (this.form.password) payload.password = this.form.password
          await updateUser(this.editing.id, payload)
        } else {
          await createUser(this.form)
        }
        this.showModal = false
        await this.load()
      } catch (err) {
        alert(err.response?.data?.error || 'เกิดข้อผิดพลาด')
      }
    },
    confirmSave() {
      if (!this.editing && (!this.form.username.trim() || !this.form.password)) return
      this.confirmAction = {
        title: this.editing ? 'ยืนยันการแก้ไข' : 'ยืนยันการเพิ่ม',
        message: this.editing ? 'คุณต้องการบันทึกการแก้ไขผู้ใช้นี้หรือไม่?' : `ต้องการเพิ่มผู้ใช้ "${this.form.username}" หรือไม่?`,
        confirmText: 'บันทึก',
        variant: 'info',
        onConfirm: () => this.save()
      }
    },
    cancelForm() {
      const isDirty = this.editing
        ? (this.form.full_name || this.form.password)
        : (this.form.username.trim() || this.form.password || this.form.full_name)
      if (isDirty) {
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
        title: 'ลบผู้ใช้',
        message: 'คุณต้องการลบผู้ใช้นี้หรือไม่? การดำเนินการนี้ไม่สามารถย้อนกลับได้',
        confirmText: 'ลบ',
        variant: 'danger',
        onConfirm: async () => {
          this.confirmAction = null
          await deleteUser(id)
          await this.load()
        }
      }
    }
  }
}
</script>

<style scoped>
.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  font-size: var(--font-size-sm);
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: var(--color-primary);
}
</style>
