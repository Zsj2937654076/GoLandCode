<template>
  <div id="app">
    <template v-if="$route.meta.requiresAuth === false">
      <!-- Unauthenticated routes (login page) -->
      <router-view />
    </template>
    <template v-else>
      <!-- Authenticated routes with main layout -->
      <el-container v-if="isAuthenticated">
        <el-header class="app-header">
          <el-row type="flex" justify="space-between" align="middle" style="height: 100%">
            <el-col :span="8">
              <h1>学生管理系统</h1>
            </el-col>
            <el-col :span="16">
              <el-menu mode="horizontal" :router="true" background-color="#409EFF" text-color="#fff" active-text-color="#ffd04b">
                <el-menu-item index="/dashboard">仪表盘</el-menu-item>
                <el-menu-item index="/students">学生管理</el-menu-item>
                <el-menu-item index="/classes">班级管理</el-menu-item>
                
                <el-sub-menu index="user" style="float: right;">
                  <template #title>
                    <span>{{ user?.username || '用户' }}</span>
                  </template>
                  <el-menu-item index="/profile">个人资料</el-menu-item>
                  <el-menu-item index="/change-password">修改密码</el-menu-item>
                  <el-menu-item @click="logout">退出登录</el-menu-item>
                </el-sub-menu>
              </el-menu>
            </el-col>
          </el-row>
        </el-header>
        
        <el-main class="app-main">
          <router-view />
        </el-main>
        
        <el-footer class="app-footer">
          <p>学生管理系统 &copy; {{ new Date().getFullYear() }}</p>
        </el-footer>
      </el-container>
      <div v-else>
        <!-- This will redirect to login if not authenticated but route requires auth -->
        <router-view />
      </div>
    </template>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

export default {
  name: 'App',
  setup() {
    const store = useStore()
    const router = useRouter()
    
    const isAuthenticated = computed(() => store.getters['auth/isAuthenticated'])
    const user = computed(() => store.getters['auth/user'])
    
    // Check if user is authenticated when component is mounted
    onMounted(() => {
      // Try to restore session from localStorage
      const token = localStorage.getItem('token')
      if (token) {
        store.dispatch('auth/restoreSession', token)
          .catch(() => {
            // If the token is invalid or expired, redirect to login
            router.push('/login')
          })
      } else if (router.currentRoute.value.meta.requiresAuth !== false) {
        // Redirect to login if current route requires auth and no token
        router.push('/login')
      }
    })
    
    // Logout function
    const logout = () => {
      store.dispatch('auth/logout')
      router.push('/login')
    }
    
    return {
      isAuthenticated,
      user,
      logout
    }
  }
}
</script>

<style>
/* App-specific styles would be here, but we use main.css for global styles */
</style> 