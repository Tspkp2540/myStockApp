<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">ประวัติการเคลื่อนไหว</h1>
      <p class="page-subtitle">ดูรายการรับเข้า-จ่ายออกทั้งหมด</p>
    </div>

    <div class="card">
      <div class="toolbar">
        <div class="toolbar-filters">
          <select v-model="filterType" class="form-control" @change="load">
            <option value="">ทุกประเภท</option>
            <option value="in">▲ รับเข้า</option>
            <option value="out">▼ จ่ายออก</option>
          </select>
          <select v-model="filterIngredient" class="form-control" @change="load">
            <option :value="''">ทุกวัตถุดิบ</option>
            <option v-for="ing in ingredients" :key="ing.id" :value="ing.id">{{ ing.name }}</option>
          </select>
        </div>
      </div>

      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>วันที่</th>
              <th>วัตถุดิบ</th>
              <th>ประเภท</th>
              <th>จำนวน</th>
              <th>หน่วย</th>
              <th>หมายเหตุ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(txn, i) in transactions" :key="txn.id">
              <td>{{ i + 1 }}</td>
              <td>{{ formatDate(txn.created_at) }}</td>
              <td>{{ txn.ingredient?.name }}</td>
              <td>
                <span :class="['badge', txn.type === 'in' ? 'badge-in' : 'badge-out']">
                  {{ txn.type === 'in' ? '▲ รับเข้า' : '▼ จ่ายออก' }}
                </span>
              </td>
              <td>{{ txn.quantity }}</td>
              <td>{{ txn.ingredient?.unit?.name }}</td>
              <td>{{ txn.note || '-' }}</td>
            </tr>
            <tr v-if="transactions.length === 0">
              <td colspan="7" class="table-empty">
                <span class="table-empty-icon">📋</span>
                ยังไม่มีรายการเคลื่อนไหว
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { getTransactions, getIngredients } from '../api/index.js'

export default {
  data() {
    return {
      transactions: [],
      ingredients: [],
      filterType: '',
      filterIngredient: ''
    }
  },
  async mounted() {
    const { data } = await getIngredients()
    this.ingredients = data
    await this.load()
  },
  methods: {
    async load() {
      const params = {}
      if (this.filterType) params.type = this.filterType
      if (this.filterIngredient) params.ingredient_id = this.filterIngredient
      const { data } = await getTransactions(params)
      this.transactions = data
    },
    formatDate(d) {
      if (!d) return '-'
      return new Date(d).toLocaleString('th-TH')
    }
  }
}
</script>
