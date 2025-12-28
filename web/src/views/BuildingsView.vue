<template>
  <div class="page">
    <div class="card">
      <div class="card-header">
        <div>
          <h2 class="card-title">公寓楼房基本信息管理</h2>
          <p class="card-subtitle">维护公寓楼的基本信息，便于寝室分配</p>
        </div>
        <button class="primary" @click="openCreate">新增公寓楼</button>
      </div>
      <div class="toolbar">
        <input
          v-model="keyword"
          class="search-input"
          placeholder="按公寓楼号搜索"
        />
        <button class="secondary" @click="load">刷新</button>
      </div>
      <p v-if="error" class="error">{{ error }}</p>
      <table class="table">
        <thead>
          <tr>
            <th>公寓楼号</th>
            <th>楼层数</th>
            <th>房间数</th>
            <th>启用时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="b in filteredList" :key="b.id">
            <td>{{ b.buildingNo }}</td>
            <td>{{ b.floorCount }}</td>
            <td>{{ b.roomCount }}</td>
            <td>{{ formatDate(b.startedAt) }}</td>
            <td>
              <button class="link" @click="openEdit(b)">编辑</button>
              <button class="link danger" @click="remove(b.id)">删除</button>
            </td>
          </tr>
          <tr v-if="!filteredList.length">
            <td colspan="5" class="empty">暂无数据</td>
          </tr>
        </tbody>
      </table>
      <div class="pagination" v-if="pageCount > 1 || filteredList.length">
        <button
          class="secondary"
          :disabled="displayPage === 1"
          @click="changePage(displayPage - 1)"
        >
          上一页
        </button>
        <span class="page-info">
          第 {{ displayPage }} / {{ displayPageCount }} 页，当前页
          {{ filteredList.length }} 条，共 {{ displayTotal }} 条
        </span>
        <button
          class="secondary"
          :disabled="displayPage >= displayPageCount"
          @click="changePage(displayPage + 1)"
        >
          下一页
        </button>
        <span class="page-info">
          跳转到
          <input v-model.number="pageInput" class="search-input" style="width: 60px" />
          页
          <button class="secondary" @click="jumpToPage">跳转</button>
        </span>
      </div>
    </div>

    <div v-if="showDialog" class="modal-mask">
      <div class="modal">
        <h3 class="modal-title">{{ form.id ? "编辑公寓楼" : "新增公寓楼" }}</h3>
        <form class="modal-form" @submit.prevent="save">
          <div class="modal-row">
            <label>公寓楼号</label>
            <input v-model="form.buildingNo" placeholder="公寓楼号" />
          </div>
          <div class="modal-row">
            <label>楼层数</label>
            <input v-model.number="form.floorCount" placeholder="楼层数" />
          </div>
          <div class="modal-row">
            <label>房间数</label>
            <input v-model.number="form.roomCount" placeholder="房间数" />
          </div>
          <div class="modal-row">
            <label>启用时间</label>
            <input v-model="form.startedAt" placeholder="如 2024-09-01" />
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
import { ref, computed, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listBuildings,
  createBuilding,
  updateBuilding,
  deleteBuilding
} from "../api/buildings";

const route = useRoute();
const router = useRouter();

const list = ref([]);
const page = computed(() => {
  const p = Number(route.query.page);
  return Number.isFinite(p) && p > 0 ? p : 1;
});
const pageSize = computed(() => {
  const s = Number(route.query.pageSize);
  return Number.isFinite(s) && s > 0 ? s : 10;
});
const total = ref(0);
const error = ref("");
const keyword = ref(route.query.keyword || "");
const showDialog = ref(false);
const form = ref({
  id: null,
  buildingNo: "",
  floorCount: null,
  roomCount: null,
  startedAt: ""
});

const filteredList = computed(() => list.value);

const pageCount = computed(() => {
  if (!pageSize.value) return 1;
  const n = Math.ceil(total.value / pageSize.value);
  return n > 0 ? n : 1;
});

const displayPage = computed(() => page.value);

const displayPageCount = computed(() => pageCount.value);

const displayTotal = computed(() => total.value);

const pageInput = ref(page.value);

const jumpToPage = () => {
  const p = Number(pageInput.value);
  if (!Number.isFinite(p) || p < 1) return;
  const target = p > pageCount.value ? pageCount.value : p;
  router.replace({
    path: route.path,
    query: {
      ...route.query,
      page: String(target),
      pageSize: String(pageSize.value)
    }
  });
};

const load = async () => {
  try {
    error.value = "";
    const res = await listBuildings({
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value || undefined
    });
    const data = res.data;
    list.value = Array.isArray(data) ? data : data.items || [];
    total.value = Array.isArray(data) ? list.value.length : data.total || 0;
  } catch (e) {
    error.value = "加载公寓楼信息失败";
  }
};

const changePage = async (p) => {
  if (p < 1) return;
  router.replace({
    path: route.path,
    query: {
      ...route.query,
      page: String(p),
      pageSize: String(pageSize.value)
    }
  });
};

const reset = () => {
  form.value = {
    id: null,
    buildingNo: "",
    floorCount: null,
    roomCount: null,
    startedAt: ""
  };
};

const openCreate = () => {
  error.value = "";
  reset();
  showDialog.value = true;
};

const openEdit = (b) => {
  error.value = "";
  form.value = {
    id: b.id,
    buildingNo: b.buildingNo,
    floorCount: b.floorCount,
    roomCount: b.roomCount,
    startedAt: formatDate(b.startedAt)
  };
  showDialog.value = true;
};

const closeDialog = () => {
  showDialog.value = false;
};

const save = async () => {
  try {
    error.value = "";
    if (!form.value.buildingNo) {
      error.value = "公寓楼号不能为空";
      return;
    }
    if (form.value.id) {
      await updateBuilding(form.value.id, form.value);
    } else {
      const data = { ...form.value };
      delete data.id;
      await createBuilding(data);
    }
    await load();
    reset();
    showDialog.value = false;
  } catch (e) {
    if (e.response && e.response.data && e.response.data.error) {
      error.value = e.response.data.error;
    } else {
      error.value = "保存公寓楼信息失败";
    }
  }
};

const formatDate = (value) => {
  if (!value) return "";
  const s = String(value);
  if (s.length >= 10) {
    return s.slice(0, 10);
  }
  return s;
};

const remove = async (id) => {
  try {
    error.value = "";
    if (window.confirm("确认删除该公寓楼吗？")) {
      await deleteBuilding(id);
      await load();
    }
  } catch (e) {
    error.value = "删除公寓楼失败";
  }
};

watch(
  () => [route.query.page, route.query.pageSize, route.query.keyword],
  () => {
    keyword.value = route.query.keyword || "";
    load();
    pageInput.value = page.value;
  },
  { immediate: true }
);

watch(
  () => keyword.value,
  () => {
    router.replace({
      path: route.path,
      query: {
        ...route.query,
        page: "1",
        pageSize: String(pageSize.value),
        keyword: keyword.value || undefined
      }
    });
  }
);
</script>

<style scoped>
.page {
  max-width: 900px;
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
  width: 220px;
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
  width: 420px;
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
.modal-row input {
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
