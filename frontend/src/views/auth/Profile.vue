<template>
  <div class="profile-container">
    <h1 class="page-title">用户资料</h1>
    
    <el-card v-loading="loading">
      <template v-if="user">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="用户名">
            {{ user.username }}
          </el-descriptions-item>
          
          <el-descriptions-item label="邮箱">
            {{ user.email }}
          </el-descriptions-item>
          
          <el-descriptions-item label="角色">
            {{ user.role === 'admin' ? '管理员' : '普通用户' }}
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="action-buttons mt-20">
          <el-button 
            type="primary" 
            @click="$router.push('/change-password')"
          >
            修改密码
          </el-button>
        </div>
      </template>
      
      <el-empty 
        v-else-if="!loading" 
        description="未找到用户资料"
      />
    </el-card>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'

export default {
  name: 'Profile',
  setup() {
    const store = useStore()
    const loading = ref(false)
    
    // Get user from store
    const user = computed(() => store.getters['auth/user'])
    
    // Fetch user profile
    const fetchProfile = async () => {
      if (user.value) return // Already have user data
      
      loading.value = true
      try {
        // Get token from localStorage
        const token = localStorage.getItem('token')
        if (token) {
          await store.dispatch('auth/restoreSession', token)
        } else {
          ElMessage.error('Not authenticated')
        }
      } catch (error) {
        console.error('Error fetching profile:', error)
        ElMessage.error('Failed to load profile')
      } finally {
        loading.value = false
      }
    }
    
    // Fetch data on component mount
    onMounted(fetchProfile)
    
    return {
      user,
      loading
    }
  }
}
</script>

<style scoped>
.action-buttons {
  display: flex;
  gap: 10px;
}
</style>

<style>
/* Add any additional styles here */
</style>
 