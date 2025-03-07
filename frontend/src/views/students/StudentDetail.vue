<template>
  <div class="student-detail-container">
    <div class="page-header">
      <h1 class="page-title">学生详情</h1>
      <div class="action-buttons">
        <el-button @click="$router.push('/students')">
          返回学生列表
        </el-button>
        <el-button 
          type="primary" 
          @click="$router.push(`/students/${id}/edit`)"
          v-if="student"
        >
          编辑学生
        </el-button>
      </div>
    </div>
    
    <el-card v-loading="loading">
      <template v-if="student">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="学号">
            {{ student.student_id }}
          </el-descriptions-item>
          
          <el-descriptions-item label="姓名">
            {{ student.name }}
          </el-descriptions-item>
          
          <el-descriptions-item label="班级">
            {{ student.class_name }}
          </el-descriptions-item>
          
          <el-descriptions-item label="邮箱">
            {{ student.email || '无' }}
          </el-descriptions-item>
          
          <el-descriptions-item label="电话">
            {{ student.phone || '无' }}
          </el-descriptions-item>
          
          <el-descriptions-item label="创建时间">
            {{ formatDate(student.created_at) }}
          </el-descriptions-item>
          
          <el-descriptions-item label="地址" :span="2">
            {{ student.address || '无' }}
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="action-buttons mt-20">
          <el-button 
            type="danger" 
            @click="confirmDelete"
          >
            删除学生
          </el-button>
        </div>
      </template>
      
      <el-empty 
        v-else-if="!loading" 
        description="未找到该学生"
      />
    </el-card>
    
    <!-- Delete Confirmation Dialog -->
    <el-dialog
      v-model="deleteDialog.visible"
      title="确认删除"
      width="30%"
    >
      <span>确定要删除这个学生吗？</span>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="deleteDialog.visible = false">取消</el-button>
          <el-button type="danger" @click="deleteStudent" :loading="deleteDialog.loading">
            删除
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

export default {
  name: 'StudentDetail',
  props: {
    id: {
      type: [Number, String],
      required: true
    }
  },
  setup(props) {
    const store = useStore()
    const router = useRouter()
    const loading = ref(false)
    
    // Delete dialog state
    const deleteDialog = reactive({
      visible: false,
      loading: false
    })
    
    // Get student from store
    const student = computed(() => store.getters['students/currentStudent'])
    
    // Format date
    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      const date = new Date(dateString)
      return date.toLocaleString()
    }
    
    // Fetch student data
    const fetchStudentData = async () => {
      loading.value = true
      try {
        await store.dispatch('students/fetchStudent', Number(props.id))
      } catch (error) {
        console.error('Error fetching student:', error)
        ElMessage.error('Failed to load student data')
      } finally {
        loading.value = false
      }
    }
    
    // Delete student
    const confirmDelete = () => {
      deleteDialog.visible = true
    }
    
    const deleteStudent = async () => {
      deleteDialog.loading = true
      try {
        await store.dispatch('students/deleteStudent', Number(props.id))
        ElMessage.success('Student deleted successfully')
        router.push('/students')
      } catch (error) {
        console.error('Error deleting student:', error)
        ElMessage.error('Failed to delete student')
      } finally {
        deleteDialog.loading = false
        deleteDialog.visible = false
      }
    }
    
    // Fetch data on component mount
    onMounted(fetchStudentData)
    
    return {
      student,
      loading,
      deleteDialog,
      formatDate,
      confirmDelete,
      deleteStudent
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

.action-buttons {
  display: flex;
  gap: 10px;
}
</style> 