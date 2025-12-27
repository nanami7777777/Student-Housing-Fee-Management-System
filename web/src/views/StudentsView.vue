<template>
  <div class="page">
    <div class="card">
      <div class="card-header">
        <div>
          <h2 class="card-title">学生基本信息管理</h2>
          <p class="card-subtitle">管理学生入住、公寓和寝室分配等信息</p>
        </div>
        <button class="primary" @click="openCreate">新增学生</button>
      </div>
      <div class="toolbar">
        <input
          v-model="keyword"
          class="search-input"
          placeholder="按学号、姓名、班级搜索"
        />
        <button class="secondary" @click="load">刷新</button>
      </div>
      <p v-if="error" class="error">{{ error }}</p>
      <table class="table">
        <thead>
          <tr>
            <th>学号</th>
            <th>姓名</th>
            <th>性别</th>
            <th>民族</th>
            <th>专业</th>
            <th>班级</th>
            <th>联系方式</th>
            <th>所属公寓楼</th>
            <th>寝室</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in filteredList" :key="s.id">
            <td>{{ s.studentNo }}</td>
            <td>{{ s.name }}</td>
            <td>{{ s.gender }}</td>
            <td>{{ s.ethnicity }}</td>
            <td>{{ s.major }}</td>
            <td>{{ s.className }}</td>
            <td>{{ s.phone }}</td>
            <td>{{ buildingName(s.buildingID) }}</td>
            <td>{{ roomName(s.roomID) }}</td>
            <td>
              <button class="link" @click="openEdit(s)">编辑</button>
              <button class="link danger" @click="remove(s.id)">删除</button>
            </td>
          </tr>
          <tr v-if="!filteredList.length">
            <td colspan="10" class="empty">暂无数据</td>
          </tr>
        </tbody>
      </table>
      <div class="pagination" v-if="total > pageSize">
        <button class="secondary" :disabled="page === 1" @click="changePage(page - 1)">
          上一页
        </button>
        <span class="page-info">
          第 {{ page }} / {{ Math.ceil(total / pageSize) }} 页，共 {{ total }} 条
        </span>
        <button
          class="secondary"
          :disabled="page >= Math.ceil(total / pageSize)"
          @click="changePage(page + 1)"
        >
          下一页
        </button>
      </div>
    </div>

    <div v-if="showDialog" class="modal-mask">
      <div class="modal">
        <h3 class="modal-title">{{ form.id ? "编辑学生" : "新增学生" }}</h3>
        <form class="modal-form" @submit.prevent="save">
          <div class="modal-row">
            <label>学号</label>
            <input v-model="form.studentNo" placeholder="学号" />
          </div>
          <div class="modal-row">
            <label>姓名</label>
            <input v-model="form.name" placeholder="姓名" />
          </div>
          <div class="modal-row">
            <label>性别</label>
            <input v-model="form.gender" placeholder="性别" />
          </div>
          <div class="modal-row">
            <label>民族</label>
            <input v-model="form.ethnicity" placeholder="民族" />
          </div>
          <div class="modal-row">
            <label>专业</label>
            <input v-model="form.major" placeholder="专业" />
          </div>
          <div class="modal-row">
            <label>班级</label>
            <input v-model="form.className" placeholder="班级" />
          </div>
          <div class="modal-row">
            <label>联系方式</label>
            <input v-model="form.phone" placeholder="联系方式" />
          </div>
          <div class="modal-row">
            <label>所属公寓楼</label>
            <select v-model.number="form.buildingID">
              <option value="">选择公寓楼</option>
              <option
                v-for="b in buildingList"
                :key="b.id"
                :value="b.id"
              >
                {{ b.buildingNo }} (ID: {{ b.id }})
              </option>
            </select>
          </div>
          <div class="modal-row">
            <label>寝室</label>
            <select v-model.number="form.roomID">
              <option value="">选择寝室</option>
              <option
                v-for="r in filteredRooms"
                :key="r.id"
                :value="r.id"
              >
                {{ r.roomNo }} (ID: {{ r.id }})
              </option>
            </select>
          </div>
          <p v-if="error" class="error">{{ error }}</p>
          <div class="modal-footer">
            <button type="button" class="secondary" @click="closeDialog">
              取消
            </button>
            <button type="submit" class="primary">
              {{ form.id ? "保存修改" : "确认添加" }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from "vue";
import {
  listStudents,
  createStudent,
  updateStudent,
  deleteStudent
} from "../api/students";
import { listBuildings } from "../api/buildings";
import { listRooms } from "../api/rooms";

const list = ref([]);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const error = ref("");
const keyword = ref("");
const showDialog = ref(false);
const buildingList = ref([]);
const roomList = ref([]);
const form = ref({
  id: null,
  studentNo: "",
  name: "",
  gender: "",
  ethnicity: "",
  major: "",
  className: "",
  phone: "",
  buildingID: null,
  roomID: null
});

const filteredList = computed(() => {
  if (!keyword.value) return list.value;
  const k = keyword.value.toLowerCase();
  return list.value.filter(
    (s) =>
      String(s.studentNo || "").toLowerCase().includes(k) ||
      String(s.name || "").toLowerCase().includes(k) ||
      String(s.className || "").toLowerCase().includes(k)
  );
});

const filteredRooms = computed(() => {
  if (!form.value.buildingID) return roomList.value;
  return roomList.value.filter((r) => r.buildingID === form.value.buildingID);
});

const load = async () => {
  try {
    error.value = "";
    const [studentsRes, buildingsRes, roomsRes] = await Promise.all([
      listStudents({ page: page.value, pageSize: pageSize.value }),
      listBuildings(),
      listRooms()
    ]);
    const data = studentsRes.data;
    list.value = Array.isArray(data) ? data : data.items || [];
    total.value = Array.isArray(data) ? list.value.length : data.total || 0;
    buildingList.value = buildingsRes.data;
    roomList.value = roomsRes.data;
  } catch (e) {
    error.value = "加载学生列表失败";
  }
};

const changePage = async (p) => {
  if (p < 1) return;
  page.value = p;
  await load();
};

const reset = () => {
  form.value = {
    id: null,
    studentNo: "",
    name: "",
    gender: "",
    ethnicity: "",
    major: "",
    className: "",
    phone: "",
    buildingID: null,
    roomID: null
  };
};

const buildingName = (id) => {
  const b = buildingList.value.find((item) => item.id === id);
  return b ? b.buildingNo : id;
};

const roomName = (id) => {
  const r = roomList.value.find((item) => item.id === id);
  return r ? r.roomNo : id;
};

const openCreate = () => {
  error.value = "";
  reset();
  showDialog.value = true;
};

const openEdit = (s) => {
  error.value = "";
  form.value = {
    id: s.id,
    studentNo: s.studentNo,
    name: s.name,
    gender: s.gender,
    ethnicity: s.ethnicity,
    major: s.major,
    className: s.className,
    phone: s.phone,
    buildingID: s.buildingID,
    roomID: s.roomID
  };
  showDialog.value = true;
};

const closeDialog = () => {
  showDialog.value = false;
};

const save = async () => {
  try {
    error.value = "";
    if (!form.value.studentNo || !form.value.name) {
      error.value = "学号和姓名不能为空";
      return;
    }
    if (form.value.id) {
      await updateStudent(form.value.id, form.value);
    } else {
      const data = { ...form.value };
      delete data.id;
      await createStudent(data);
    }
    await load();
    reset();
    showDialog.value = false;
  } catch (e) {
    error.value = "保存学生信息失败";
  }
};

const remove = async (id) => {
  try {
    error.value = "";
    if (window.confirm("确认删除该学生吗？")) {
      await deleteStudent(id);
      await load();
    }
  } catch (e) {
    error.value = "删除学生失败";
  }
};

watch(
  () => form.value.buildingID,
  () => {
    if (!filteredRooms.value.some((r) => r.id === form.value.roomID)) {
      form.value.roomID = null;
    }
  }
);

onMounted(load);
</script>

<style scoped>
.page {
  max-width: 1200px;
  margin: 0 auto;
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
.toolbar {
  margin-top: 16px;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.search-input {
  width: 260px;
  padding: 6px 10px;
  border-radius: 6px;
  border: 1px solid #d1d5db;
  font-size: 14px;
}
.buttons {
  display: flex;
  gap: 8px;
  align-items: center;
}
.primary {
  padding: 6px 14px;
  border-radius: 6px;
  border: none;
  background: #2563eb;
  color: #fff;
  font-size: 14px;
  cursor: pointer;
}
.primary:hover {
  background: #1d4ed8;
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
.link {
  border: none;
  background: none;
  color: #2563eb;
  cursor: pointer;
  padding: 0 4px;
  font-size: 13px;
}
.link.danger {
  color: #dc2626;
}
.error {
  color: red;
  margin-top: 8px;
}
.pagination {
  margin-top: 10px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}
.page-info {
  color: #4b5563;
}
.empty {
  text-align: center;
  color: #9ca3af;
  font-size: 13px;
}
.modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  width: 520px;
  max-width: 90%;
  background: #ffffff;
  border-radius: 10px;
  padding: 20px 22px 18px;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.3);
}
.modal-title {
  margin: 0 0 12px;
  font-size: 16px;
  font-weight: 600;
}
.modal-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px 16px;
}
.modal-row {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.modal-row label {
  font-size: 13px;
  color: #4b5563;
}
.modal-row input,
.modal-row select {
  padding: 6px 8px;
  border-radius: 6px;
  border: 1px solid #d1d5db;
  font-size: 14px;
}
.modal-footer {
  grid-column: 1 / -1;
  margin-top: 4px;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
