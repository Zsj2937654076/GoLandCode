<template>
  <div class="class-list-container">
    <div class="page-header">
      <h1 class="page-title">班级管理</h1>
      <el-button type="primary" @click="$router.push('/classes/new')">
        添加班级
      </el-button>
    </div>
    
    <!-- Classes Table -->
    <el-card>
      <el-table 
        :data="classes" 
        v-loading="loading" 
        style="width: 100%"
        border
      >
        <el-table-column prop="name" label="班级名称" sortable />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="student_count" label="学生人数" width="100" sortable />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="scope">
            <el-button 
              size="small" 
              type="primary" 
              @click="$router.push(`/classes/${scope.row.id}`)"
            >
              查看
            </el-button>
            <el-button 
              size="small" 
              type="warning" 
              @click="$router.push(`/classes/${scope.row.id}/edit`)"
            >
              编辑
            </el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="confirmDelete(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- Empty state -->
      <el-empty 
        v-if="classes.length === 0 && !loading" 
        description="暂无班级数据"
      >
        <el-button type="primary" @click="$router.push('/classes/new')">
          创建第一个班级
        </el-button>
      </el-empty>
    </el-card>
    
    <!-- Delete Confirmation Dialog -->
    <el-dialog
      v-model="deleteDialog.visible"
      title="确认删除"
      width="30%"
    >
      <span>
        确定要删除班级 "{{ deleteDialog.class?.name }}" 吗？
        <br><br>
        <strong v-if="deleteDialog.class?.student_count > 0" class="text-danger">
          警告：该班级有 {{ deleteDialog.class.student_count }} 名学生。 
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
            :disabled="deleteDialog.class?.student_count > 0"
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
import { ElMessage } from 'element-plus'

export default {
  name: 'ClassList',
  setup() {
    const store = useStore()
    const loading = ref(false)
    
    // Get classes from store
    const classes = computed(() => store.getters['classes/allClasses'])
    
    // Delete dialog state
    const deleteDialog = reactive({
      visible: false,
      class: null,
      loading: false
    })
    
    // Fetch classes
    const fetchClasses = async () => {
      loading.value = true
      try {
        await store.dispatch('classes/fetchClasses')
      } catch (error) {
        console.error('Error fetching classes:', error)
        ElMessage.error('Failed to load classes')
      } finally {
        loading.value = false
      }
    }
    
    // Delete class
    const confirmDelete = (classData) => {
      deleteDialog.class = classData
      deleteDialog.visible = true
    }
    
    const deleteClass = async () => {
      if (!deleteDialog.class) return
      
      // Check if class has students
      if (deleteDialog.class.student_count > 0) {
        ElMessage.error('Cannot delete a class that has students')
        return
      }
      
      deleteDialog.loading = true
      try {
        await store.dispatch('classes/deleteClass', deleteDialog.class.id)
        ElMessage.success('Class deleted successfully')
        deleteDialog.visible = false
      } catch (error) {
        console.error('Error deleting class:', error)
        if (error.response && error.response.status === 409) {
          ElMessage.error('Cannot delete a class that has associated students')
        } else {
          ElMessage.error('Failed to delete class')
        }
      } finally {
        deleteDialog.loading = false
      }
    }
    
    // Fetch data on component mount
    onMounted(fetchClasses)
    
    return {
      classes,
      loading,
      deleteDialog,
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

.text-danger {
  color: #F56C6C;
}
</style> 