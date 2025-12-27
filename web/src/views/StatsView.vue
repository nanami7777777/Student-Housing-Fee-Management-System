<template>
  <div class="page">
    <div class="card">
      <div class="card-header">
        <div>
          <h2 class="card-title">公寓入住率统计</h2>
          <p class="card-subtitle">
            按公寓楼维度统计总床位数、已入住人数及入住率
          </p>
        </div>
        <button class="secondary" @click="load">刷新</button>
      </div>
      <p v-if="errorOccupancy" class="error">{{ errorOccupancy }}</p>
      <table class="table">
        <thead>
          <tr>
            <th>公寓楼号</th>
            <th>总床位数</th>
            <th>已入住人数</th>
            <th>入住率</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in occupancyList" :key="item.buildingID">
            <td>{{ item.buildingNo }}</td>
            <td>{{ item.totalCapacity }}</td>
            <td>{{ item.occupiedBeds }}</td>
            <td>{{ item.occupancyRate }}%</td>
          </tr>
          <tr v-if="!occupancyList.length">
            <td colspan="4" class="empty">暂无数据</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="card">
      <div class="card-header">
        <div>
          <h2 class="card-title">公寓收费统计</h2>
          <p class="card-subtitle">按公寓楼统计累计交费总金额</p>
        </div>
        <button class="secondary" @click="load">刷新</button>
      </div>
      <p v-if="errorPayment" class="error">{{ errorPayment }}</p>
      <table class="table">
        <thead>
          <tr>
            <th>公寓楼号</th>
            <th>累计交费金额</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in paymentList" :key="item.buildingID">
            <td>{{ item.buildingNo }}</td>
            <td>{{ formatAmount(item.totalAmount) }}</td>
          </tr>
          <tr v-if="!paymentList.length">
            <td colspan="2" class="empty">暂无数据</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import {
  getBuildingOccupancy,
  getBuildingPaymentSummary
} from "../api/stats";

const occupancyList = ref([]);
const paymentList = ref([]);
const errorOccupancy = ref("");
const errorPayment = ref("");

const load = async () => {
  try {
    errorOccupancy.value = "";
    errorPayment.value = "";
    const [occRes, payRes] = await Promise.all([
      getBuildingOccupancy(),
      getBuildingPaymentSummary()
    ]);
    occupancyList.value = occRes.data;
    paymentList.value = payRes.data;
  } catch (e) {
    errorOccupancy.value = "加载统计数据失败";
    errorPayment.value = "加载统计数据失败";
  }
};

const formatAmount = (value) => {
  if (value == null) return "0.00";
  return Number(value).toFixed(2);
};

onMounted(load);
</script>

<style scoped>
.page {
  max-width: 1000px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.card {
  background: #ffffff;
  border-radius: 8px;
  padding: 16px 20px 20px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.08);
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.card-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}
.card-subtitle {
  margin: 4px 0 0;
  font-size: 13px;
  color: #6b7280;
}
.secondary {
  padding: 6px 12px;
  border-radius: 6px;
  border: 1px solid #d1d5db;
  background: #ffffff;
  color: #374151;
  font-size: 14px;
  cursor: pointer;
}
.table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 12px;
}
.table th,
.table td {
  border: 1px solid #ccc;
  padding: 6px 10px;
  font-size: 13px;
}
.table thead {
  background: #f9fafb;
}
.table tbody tr:nth-child(even) {
  background: #f9fafb;
}
.error {
  color: red;
  margin-top: 8px;
}
.empty {
  text-align: center;
  color: #9ca3af;
  font-size: 13px;
}
</style>

