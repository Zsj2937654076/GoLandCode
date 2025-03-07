<template>
  <div class="class-detail-container">
    <div class="page-header">
      <h1 class="page-title">班级详情</h1>
      <div class="action-buttons">
        <el-button @click="$router.push('/classes')">
          返回班级列表
        </el-button>
        <el-button 
          type="primary" 
          @click="$router.push(`/classes/${id}/edit`)"
          v-if="classData"
        >
          编辑班级
        </el-button>
      </div>
    </div>
    
    <el-card v-loading="loading">
      <template v-if="classData">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="班级名称">
            {{ classData.name }}
          </el-descriptions-item>
          
          <el-descriptions-item label="描述">
            {{ classData.description || '无描述' }}
          </el-descriptions-item>
          
          <el-descriptions-item label="学生人数">
            {{ classData.student_count }}
          </el-descriptions-item>
          
          <el-descriptions-item label="创建时间">
            {{ formatDate(classData.created_at) }}
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="action-buttons mt-20">
          <el-button 
            type="danger" 
            @click="confirmDelete"
            :disabled="classData.student_count > 0"
          >
            删除班级
          </el-button>
        </div>
        
        <!-- Students in this class -->
        <div class="students-section mt-20">
          <h2>班级学生</h2>
          
          <el-table 
            :data="students" 
            v-loading="loadingStudents" 
            style="width: 100%"
            border
          >
            <el-table-column prop="student_id" label="学号" width="120" />
            <el-table-column prop="name" label="姓名" />
            <el-table-column prop="email" label="邮箱" />
            <el-table-column prop="phone" label="电话" width="150" />
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="scope">
                <el-button 
                  size="small" 
                  type="primary" 
                  @click="$router.push(`/students/${scope.row.id}`)"
                >
                  查看
                </el-button>
                <el-button 
                  size="small" 
                  type="warning" 
                  @click="$router.push(`/students/${scope.row.id}/edit`)"
                >
                  编辑
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <!-- Empty state for students -->
          <el-empty 
            v-if="students.length === 0 && !loadingStudents" 
            description="该班级暂无学生"
          >
            <el-button type="primary" @click="$router.push('/students/new')">
              添加学生
            </el-button>
          </el-empty>
        </div>
      </template>
      
      <el-empty 
        v-else-if="!loading" 
        description="未找到该班级"
      />
    </el-card>
    
    <!-- Delete Confirmation Dialog -->
    <el-dialog
      v-model="deleteDialog.visible"
      title="确认删除"
      width="30%"
    >
      <span>
        确定要删除这个班级吗？
        <br><br>
        <strong v-if="classData?.student_count > 0" class="text-danger">
          警告：该班级有 {{ classData.student_count }} 名学生。 
          不能删除包含学生的班级。
        </strong>
      </span>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="deleteDialog.visible = false">取消</el-button>
          <el-button 
            type="danger" 
            @click="deleteClass" 
            :loading="deleteDialog.loading"
            :disabled="classData?.student_count > 0"
          >
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
  name: 'ClassDetail',
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
    const loadingStudents = ref(false)
    
    // Delete dialog state
    const deleteDialog = reactive({
      visible: false,
      loading: false
    })
    
    // Get class and students from store
    const classData = computed(() => store.getters['classes/currentClass'])
    const students = computed(() => store.getters['classes/classStudents'])
    
    // Format date
    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      const date = new Date(dateString)
      return date.toLocaleString()
    }
    
    // Fetch class data
    const fetchClassData = async () => {
      loading.value = true
      try {
        await store.dispatch('classes/fetchClass', Number(props.id))
      } catch (error) {
        console.error('Error fetching class:', error)
        ElMessage.error('Failed to load class data')
      } finally {
        loading.value = false
      }
    }
    
    // Fetch students in this class
    const fetchClassStudents = async () => {
      loadingStudents.value = true
      try {
        await store.dispatch('classes/fetchClassStudents', Number(props.id))
      } catch (error) {
        console.error('Error fetching class students:', error)
        ElMessage.error('Failed to load students')
      } finally {
        loadingStudents.value = false
      }
    }
    
    // Delete class
    const confirmDelete = () => {
      if (classData.value.student_count > 0) {
        ElMessage.warning('Cannot delete a class that has students')
        return
      }
      deleteDialog.visible = true
    }
    
    const deleteClass = async () => {
      deleteDialog.loading = true
      try {
        await store.dispatch('classes/deleteClass', Number(props.id))
        ElMessage.success('Class deleted successfully')
        router.push('/classes')
      } catch (error) {
        console.error('Error deleting class:', error)
        if (error.response && error.response.status === 409) {
          ElMessage.error('Cannot delete a class that has associated students')
        } else {
          ElMessage.error('Failed to delete class')
        }
      } finally {
        deleteDialog.loading = false
        deleteDialog.visible = false
      }
    }
    
    // Fetch data on component mount
    onMounted(async () => {
      await fetchClassData()
      if (classData.value) {
        await fetchClassStudents()
      }
    })
    
    return {
      classData,
      students,
      loading,
      loadingStudents,
      deleteDialog,
      formatDate,
      confirmDelete,
      deleteClass
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

.students-section {
  margin-top: 30px;
}

.students-section h2 {
  margin-bottom: 15px;
  font-size: 18px;
}

.text-danger {
  color: #F56C6C;
}
</style> 