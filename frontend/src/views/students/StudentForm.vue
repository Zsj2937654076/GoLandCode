<template>
  <div class="student-form-container">
    <div class="page-header">
      <h1 class="page-title">{{ isEdit ? '编辑学生' : '添加学生' }}</h1>
      <el-button @click="$router.push('/students')">
        返回学生列表
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
        <el-row :gutter="20">
          <!-- Student ID -->
          <el-col :span="12">
            <el-form-item label="学号" prop="student_id">
              <el-input v-model="form.student_id" placeholder="请输入学号" />
            </el-form-item>
          </el-col>
          
          <!-- Name -->
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="form.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          
          <!-- Class -->
          <el-col :span="12">
            <el-form-item label="班级" prop="class_id">
              <el-select 
                v-model="form.class_id" 
                placeholder="选择班级"
                style="width: 100%"
              >
                <el-option 
                  v-for="option in classOptions" 
                  :key="option.value" 
                  :label="option.label" 
                  :value="option.value" 
                />
              </el-select>
            </el-form-item>
          </el-col>
          
          <!-- Email -->
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱地址" type="email" />
            </el-form-item>
          </el-col>
          
          <!-- Phone -->
          <el-col :span="12">
            <el-form-item label="电话" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入电话号码" />
            </el-form-item>
          </el-col>
          
          <!-- Address -->
          <el-col :span="24">
            <el-form-item label="地址" prop="address">
              <el-input 
                v-model="form.address" 
                placeholder="请输入地址" 
                type="textarea" 
                :rows="3" 
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <!-- Form Actions -->
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            {{ isEdit ? '更新学生' : '创建学生' }}
          </el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button @click="$router.push('/students')" type="info">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

export default {
  name: 'StudentForm',
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
      student_id: '',
      name: '',
      class_id: '',
      email: '',
      phone: '',
      address: ''
    })
    
    // Form validation rules
    const rules = {
      student_id: [
        { required: true, message: 'Please enter student ID', trigger: 'blur' },
        { min: 3, max: 20, message: 'Length should be 3 to 20 characters', trigger: 'blur' }
      ],
      name: [
        { required: true, message: 'Please enter student name', trigger: 'blur' },
        { min: 2, max: 100, message: 'Length should be 2 to 100 characters', trigger: 'blur' }
      ],
      class_id: [
        { required: true, message: 'Please select a class', trigger: 'change' }
      ],
      email: [
        { type: 'email', message: 'Please enter a valid email address', trigger: 'blur' }
      ]
    }
    
    // Get class options from store
    const classOptions = computed(() => store.getters['classes/classOptions'])
    
    // Fetch student data if in edit mode
    const fetchStudentData = async () => {
      if (!props.isEdit || !props.id) return
      
      loading.value = true
      try {
        await store.dispatch('students/fetchStudent', props.id)
        const student = store.getters['students/currentStudent']
        
        if (student) {
          // Populate form with student data
          form.id = student.id
          form.student_id = student.student_id
          form.name = student.name
          form.class_id = student.class_id
          form.email = student.email || ''
          form.phone = student.phone || ''
          form.address = student.address || ''
        } else {
          ElMessage.error('Student not found')
          router.push('/students')
        }
      } catch (error) {
        console.error('Error fetching student:', error)
        ElMessage.error('Failed to load student data')
        router.push('/students')
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
            // Update existing student
            await store.dispatch('students/updateStudent', { ...form })
            ElMessage.success('Student updated successfully')
          } else {
            // Create new student
            await store.dispatch('students/createStudent', { ...form })
            ElMessage.success('Student created successfully')
          }
          router.push('/students')
        } catch (error) {
          console.error('Error saving student:', error)
          ElMessage.error('Failed to save student')
        } finally {
          submitting.value = false
        }
      })
    }
    
    // Reset form
    const resetForm = () => {
      if (props.isEdit) {
        // If editing, reset to original values
        fetchStudentData()
      } else {
        // If creating new, clear all fields
        formRef.value.resetFields()
      }
    }
    
    // Fetch data on component mount
    onMounted(async () => {
      loading.value = true
      try {
        // Fetch classes for dropdown
        await store.dispatch('classes/fetchClasses')
        
        // Fetch student data if in edit mode
        if (props.isEdit) {
          await fetchStudentData()
        }
      } catch (error) {
        console.error('Error initializing form:', error)
        ElMessage.error('Failed to initialize form')
      } finally {
        loading.value = false
      }
    })
    
    return {
      formRef,
      form,
      rules,
      loading,
      submitting,
      classOptions,
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