<template>
  <div class="login-container">
    <div class="login-form">
      <h1 class="login-title">学生管理系统</h1>
      
      <el-alert
        v-if="error"
        :title="error"
        type="error"
        show-icon
        @close="error = ''"
        class="mb-20"
      />
      
      <el-form
        ref="formRef"
        :model="loginForm"
        :rules="rules"
        label-position="top"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            prefix-icon="el-icon-user"
            clearable
          />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="el-icon-lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            class="full-width"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
      
      <p class="text-center mt-20">
        默认管理员账号：admin / admin123
      </p>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

export default {
  name: 'Login',
  setup() {
    const store = useStore()
    const router = useRouter()
    const formRef = ref(null)
    const loginForm = reactive({
      username: '',
      password: ''
    })
    const loading = ref(false)
    const error = ref('')
    
    // Form validation rules
    const rules = {
      username: [
        { required: true, message: 'Please enter your username', trigger: 'blur' }
      ],
      password: [
        { required: true, message: 'Please enter your password', trigger: 'blur' }
      ]
    }
    
    // Handle login form submission
    const handleLogin = () => {
      formRef.value.validate(async valid => {
        if (!valid) return
        
        loading.value = true
        error.value = ''
        
        try {
          await store.dispatch('auth/login', loginForm)
          router.push('/dashboard')
        } catch (err) {
          error.value = err.response?.data || 'Login failed. Please check your credentials.'
        } finally {
          loading.value = false
        }
      })
    }
    
    return {
      formRef,
      loginForm,
      loading,
      error,
      rules,
      handleLogin
    }
  }
}
</script>

<style scoped>
/* Component specific styles would go here */
/* Most styling comes from global CSS in assets/styles/main.css */
</style> 