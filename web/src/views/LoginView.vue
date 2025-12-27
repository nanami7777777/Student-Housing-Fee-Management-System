<template>
  <div class="login-wrapper">
    <div class="login-card">
      <h1 class="title">学生公寓交费管理系统</h1>
      <p class="subtitle">请使用管理员账号登录系统</p>
      <form class="form" @submit.prevent="submit">
        <div class="form-item">
          <label>账号</label>
          <input v-model="username" placeholder="请输入账号" />
        </div>
        <div class="form-item">
          <label>密码</label>
          <input v-model="password" type="password" placeholder="请输入密码" />
        </div>
        <button type="submit" class="primary">登录</button>
        <p v-if="error" class="error">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import axios from "axios";
import { useRouter } from "vue-router";

const router = useRouter();
const username = ref("");
const password = ref("");
const error = ref("");

const submit = async () => {
  error.value = "";
  if (!username.value || !password.value) {
    error.value = "账号和密码不能为空";
    return;
  }
  try {
    const res = await axios.post("http://localhost:8080/api/login", {
      username: username.value,
      password: password.value
    });
    localStorage.setItem("token", res.data.token);
    router.push("/students");
  } catch (e) {
    error.value = "登录失败，请检查账号或密码";
  }
};
</script>

<style scoped>
.login-wrapper {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1f2933, #4b5563);
}
.login-card {
  width: 360px;
  padding: 32px 28px;
  border-radius: 12px;
  background: #ffffff;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.4);
}
.title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  text-align: center;
  color: #111827;
}
.subtitle {
  margin-top: 8px;
  margin-bottom: 24px;
  text-align: center;
  color: #6b7280;
  font-size: 13px;
}
.form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.form-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
label {
  font-size: 13px;
  color: #4b5563;
}
input {
  padding: 8px 10px;
  border-radius: 6px;
  border: 1px solid #d1d5db;
  outline: none;
  font-size: 14px;
}
input:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 1px rgba(37, 99, 235, 0.2);
}
.primary {
  margin-top: 4px;
  padding: 8px 0;
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
.error {
  margin-top: 6px;
  font-size: 13px;
  color: #dc2626;
  text-align: center;
}
</style>
