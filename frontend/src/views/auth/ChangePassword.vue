<template>
  <div class="change-password-container">
    <div class="page-header">
      <h1 class="page-title">修改密码</h1>
      <el-button @click="$router.push('/profile')">
        返回资料
      </el-button>
    </div>
    
    <el-card>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        v-loading="loading"
      >
        <!-- Current Password -->
        <el-form-item label="当前密码" prop="oldPassword">
          <el-input 
            v-model="form.oldPassword" 
            placeholder="请输入当前密码" 
            type="password"
            show-password
          />
        </el-form-item>
        
        <!-- New Password -->
        <el-form-item label="新密码" prop="newPassword">
          <el-input 
            v-model="form.newPassword" 
            placeholder="请输入新密码" 
            type="password"
            show-password
          />
        </el-form-item>
        
        <!-- Confirm New Password -->
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input 
            v-model="form.confirmPassword" 
            placeholder="请确认新密码" 
            type="password"
            show-password
          />
        </el-form-item>
        
        <!-- Form Actions -->
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            修改密码
          </el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button @click="$router.push('/profile')" type="info">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

export default {
  name: 'ChangePassword',
  setup() {
    const store = useStore()
    const router = useRouter()
    const formRef = ref(null)
    const loading = ref(false)
    const submitting = ref(false)
    
    // Form data
    const form = reactive({
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    })
    
    // Validate password match
    const validatePasswordMatch = (rule, value, callback) => {
      if (value !== form.newPassword) {
        callback(new Error('Passwords do not match'))
      } else {
        callback()
      }
    }
    
    // Form validation rules
    const rules = {
      oldPassword: [
        { required: true, message: 'Please enter your current password', trigger: 'blur' }
      ],
      newPassword: [
        { required: true, message: 'Please enter your new password', trigger: 'blur' },
        { min: 6, message: 'Password must be at least 6 characters', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: 'Please confirm your new password', trigger: 'blur' },
        { validator: validatePasswordMatch, trigger: 'blur' }
      ]
    }
    
    // Submit form
    const submitForm = () => {
      formRef.value.validate(async valid => {
        if (!valid) return
        
        submitting.value = true
        try {
          await store.dispatch('auth/changePassword', {
            oldPassword: form.oldPassword,
            newPassword: form.newPassword
          })
          
          ElMessage.success('Password changed successfully')
          router.push('/profile')
        } catch (error) {
          console.error('Error changing password:', error)
          if (error.response && error.response.status === 400) {
            ElMessage.error('Current password is incorrect')
          } else {
            ElMessage.error('Failed to change password')
          }
        } finally {
          submitting.value = false
        }
      })
    }
    
    // Reset form
    const resetForm = () => {
      formRef.value.resetFields()
    }
    
    return {
      formRef,
      form,
      rules,
      loading,
      submitting,
      submitForm,
      resetForm
    }
  }
}
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style> 