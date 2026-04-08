<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">จัดการผู้ใช้งาน</h1>
      <p class="page-subtitle">เพิ่ม แก้ไข สิทธิ์ผู้ใช้ระบบ (เฉพาะผู้ดูแล)</p>
    </div>

    <div class="card">
      <div class="card-header">
        <h2 class="section-title"><span class="material-symbols-outlined">group</span> รายชื่อผู้ใช้</h2>
        <button class="btn btn-primary" @click="openModal()"><span class="material-symbols-outlined">person_add</span> เพิ่มผู้ใช้</button>
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
                  <span class="material-symbols-outlined">{{ u.role === 'admin' ? 'shield_person' : 'person' }}</span>
                  {{ u.role === 'admin' ? 'ผู้ดูแล' : 'พนักงาน' }}
                </span>
              </td>
              <td>
                <span :class="['badge', u.active ? 'badge-ok' : 'badge-low']">
                  <span class="material-symbols-outlined">{{ u.active ? 'check_circle' : 'cancel' }}</span>
                  {{ u.active ? 'ใช้งาน' : 'ระงับ' }}
                </span>
              </td>
              <td>
                <div class="btn-group">
                  <button class="btn btn-sm btn-warning" @click="openModal(u)"><span class="material-symbols-outlined">edit</span> แก้ไข</button>
                  <button class="btn btn-sm btn-outline-danger" @click="remove(u.id)" v-if="u.id !== currentUser?.id"><span class="material-symbols-outlined">delete</span> ลบ</button>
                </div>
              </td>
            </tr>
            <tr v-if="users.length === 0">
              <td colspan="6" class="table-empty">
                <span class="table-empty-icon"><span class="material-symbols-outlined">person_off</span></span>
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
          <input v-model="form.username" class="form-control" :class="{ 'is-invalid': usernameError }" placeholder="username (a-z, 0-9)" maxlength="15" @input="validateUsername" />
          <div v-if="usernameError" class="inline-error"><span class="material-symbols-outlined">error</span> {{ usernameError }}</div>
        </div>
        <div class="form-group">
          <label>ชื่อ-นามสกุล</label>
          <input v-model="form.full_name" class="form-control" placeholder="เช่น สมชาย ใจดี" maxlength="100" />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>{{ editing ? 'รหัสผ่านใหม่ (เว้นว่างถ้าไม่เปลี่ยน)' : 'รหัสผ่าน *' }}</label>
            <input v-model="form.password" class="form-control" :class="{ 'is-invalid': passwordError }" type="password" placeholder="••••••••" maxlength="15" @input="validatePassword" />
            <div v-if="passwordError" class="inline-error"><span class="material-symbols-outlined">error</span> {{ passwordError }}</div>
          </div>
          <div class="form-group">
            <label>สิทธิ์</label>
            <select v-model="form.role" class="form-control">
              <option value="employee">พนักงาน</option>
              <option value="admin">ผู้ดูแลระบบ</option>
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
      confirmAction: null,
      usernameError: '',
      passwordError: ''
    }
  },
  async mounted() {
    const { state } = useAuth()
    this.currentUser = state.user
    await this.load()
  },
  methods: {
    validateUsername() {
      const v = this.form.username
      // Strip non-alphanumeric characters
      this.form.username = v.replace(/[^a-zA-Z0-9]/g, '')
      const val = this.form.username
      if (!val) {
        this.usernameError = ''
        return
      }
      if (val.length < 4) {
        this.usernameError = 'ต้องมีอย่างน้อย 4 ตัวอักษร'
      } else {
        this.usernameError = ''
      }
    },
    validatePassword() {
      const v = this.form.password
      if (/[^a-zA-Z0-9]/.test(v)) {
        this.form.password = ''
        this.passwordError = ''
        this.confirmAction = {
          title: 'รหัสผ่านไม่ถูกต้อง',
          message: 'รหัสผ่านต้องเป็นภาษาอังกฤษหรือตัวเลขเท่านั้น กรุณากรอกใหม่',
          confirmText: 'ตกลง',
          variant: 'warning',
          onConfirm: () => { this.confirmAction = null }
        }
        return
      }
      const val = this.form.password
      if (!val) {
        this.passwordError = ''
        return
      }
      if (val.length < 6) {
        this.passwordError = 'ต้องมีอย่างน้อย 6 ตัวอักษร'
      } else {
        this.passwordError = ''
      }
    },
    async load() {
      const { data } = await getUsers()
      this.users = data
    },
    openModal(user = null) {
      this.editing = user
      this.usernameError = ''
      this.passwordError = ''
      if (user) {
        this.form = { full_name: user.full_name, role: user.role, active: user.active, password: '' }
      } else {
        this.form = { username: '', password: '', full_name: '', role: 'employee', active: true }
      }
      this.showModal = true
    },
    async save() {
      this.confirmAction = null
      if (this.hasValidationErrors()) return
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
    hasValidationErrors() {
      if (!this.editing) {
        const u = this.form.username.trim()
        if (!u || u.length < 4 || !/^[a-zA-Z0-9]+$/.test(u)) return true
        if (!this.form.password || this.form.password.length < 6 || !/^[a-zA-Z0-9]+$/.test(this.form.password)) return true
      } else if (this.form.password) {
        if (this.form.password.length < 6 || !/^[a-zA-Z0-9]+$/.test(this.form.password)) return true
      }
      return false
    },
    confirmSave() {
      // Trigger validation display
      if (!this.editing) {
        this.validateUsername()
        this.validatePassword()
        if (!this.form.username.trim()) this.usernameError = 'กรุณากรอกชื่อผู้ใช้'
        if (!this.form.password) this.passwordError = 'กรุณากรอกรหัสผ่าน'
      } else if (this.form.password) {
        this.validatePassword()
      }
      if (this.hasValidationErrors()) return
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
