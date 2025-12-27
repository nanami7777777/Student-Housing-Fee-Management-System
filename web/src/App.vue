<template>
  <div class="layout">
    <header class="header" v-if="showNav">
      <div class="logo">学生公寓交费管理系统</div>
      <nav class="menu">
        <router-link to="/students" class="menu-item">学生信息</router-link>
        <router-link to="/buildings" class="menu-item">公寓楼房</router-link>
        <router-link to="/rooms" class="menu-item">公寓寝室</router-link>
        <router-link to="/payments" class="menu-item">交费信息</router-link>
        <router-link to="/users" class="menu-item">系统用户</router-link>
        <router-link to="/stats" class="menu-item">统计分析</router-link>
      </nav>
      <button class="logout" @click="logout">退出登录</button>
    </header>
    <main class="content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();

const showNav = computed(() => route.path !== "/login");

const logout = () => {
  localStorage.removeItem("token");
  router.push("/login");
};
</script>

<style scoped>
.layout {
  min-height: 100vh;
  background: #f3f4f6;
}
.header {
  display: flex;
  align-items: center;
  padding: 0 24px;
  height: 56px;
  background: #1f2933;
  color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
}
.logo {
  font-size: 18px;
  font-weight: 600;
}
.menu {
  display: flex;
  gap: 16px;
  margin-left: 40px;
}
.menu-item {
  color: #e5e7eb;
  text-decoration: none;
  padding: 6px 10px;
  border-radius: 4px;
  font-size: 14px;
}
.menu-item.router-link-active {
  background: #111827;
  color: #fff;
}
.logout {
  margin-left: auto;
  background: transparent;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  padding: 4px 12px;
  color: #e5e7eb;
  cursor: pointer;
}
.logout:hover {
  background: rgba(249, 250, 251, 0.1);
}
.content {
  padding: 24px;
}
</style>
