<template>
  <div class="class-form-container">
    <div class="page-header">
      <h1 class="page-title">{{ isEdit ? '编辑班级' : '添加班级' }}</h1>
      <el-button @click="$router.push('/classes')">
        返回班级列表
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
        <!-- Class Name -->
        <el-form-item label="班级名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入班级名称" />
        </el-form-item>
        
        <!-- Description -->
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="form.description" 
            placeholder="请输入班级描述" 
            type="textarea" 
            :rows="4" 
          />
        </el-form-item>
        
        <!-- Form Actions -->
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            {{ isEdit ? '更新班级' : '创建班级' }}
          </el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button @click="$router.push('/classes')" type="info">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

export default {
  name: 'ClassForm',
  props: {
    id: {
      type: Number,
      required: false
    },
    isEdit: {
      type: Boolean,
      default: false
    }
  },
  setup(props) {
    const store = useStore()
    const router = useRouter()
    const formRef = ref(null)
    const loading = ref(false)
    const submitting = ref(false)
    
    // Form data
    const form = reactive({
      id: null,
      name: '',
      description: ''
    })
    
    // Form validation rules
    const rules = {
      name: [
        { required: true, message: 'Please enter class name', trigger: 'blur' },
        { min: 2, max: 100, message: 'Length should be 2 to 100 characters', trigger: 'blur' }
      ]
    }
    
    // Fetch class data if in edit mode
    const fetchClassData = async () => {
      if (!props.isEdit || !props.id) return
      
      loading.value = true
      try {
        await store.dispatch('classes/fetchClass', props.id)
        const classData = store.getters['classes/currentClass']
        
        if (classData) {
          // Populate form with class data
          form.id = classData.id
          form.name = classData.name
          form.description = classData.description || ''
        } else {
          ElMessage.error('Class not found')
          router.push('/classes')
        }
      } catch (error) {
        console.error('Error fetching class:', error)
        ElMessage.error('Failed to load class data')
        router.push('/classes')
      } finally {
        loading.value = false
      }
    }
    
    // Submit form
    const submitForm = () => {
      formRef.value.validate(async valid => {
        if (!valid) return
        
        submitting.value = true
        try {
          if (props.isEdit) {
            // Update existing class
            await store.dispatch('classes/updateClass', { ...form })
            ElMessage.success('Class updated successfully')
          } else {
            // Create new class
            await store.dispatch('classes/createClass', { ...form })
            ElMessage.success('Class created successfully')
          }
          router.push('/classes')
        } catch (error) {
          console.error('Error saving class:', error)
          ElMessage.error('Failed to save class')
        } finally {
          submitting.value = false
        }
      })
    }
    
    // Reset form
    const resetForm = () => {
      if (props.isEdit) {
        // If editing, reset to original values
        fetchClassData()
      } else {
        // If creating new, clear all fields
        formRef.value.resetFields()
      }
    }
    
    // Fetch data on component mount
    onMounted(() => {
      if (props.isEdit) {
        fetchClassData()
      }
    })
    
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