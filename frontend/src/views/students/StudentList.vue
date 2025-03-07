<template>
  <div class="student-list-container">
    <div class="page-header">
      <h1 class="page-title">学生管理</h1>
      <el-button type="primary" @click="$router.push('/students/new')">
        添加学生
      </el-button>
    </div>
    
    <!-- Filter Form -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-container">
        <el-form-item label="班级">
          <el-select 
            v-model="filters.classId" 
            placeholder="选择班级" 
            clearable
            @change="handleFilterChange"
          >
            <el-option 
              v-for="option in classOptions" 
              :key="option.value" 
              :label="option.label" 
              :value="option.value" 
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="学号">
          <el-input 
            v-model="filters.studentId" 
            placeholder="搜索学号" 
            clearable
            @input="handleFilterChange"
          />
        </el-form-item>
        
        <el-form-item label="姓名">
          <el-input 
            v-model="filters.name" 
            placeholder="搜索姓名" 
            clearable
            @input="handleFilterChange"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleFilterChange">搜索</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- Students Table -->
    <el-card class="mt-20">
      <el-table 
        :data="students" 
        v-loading="loading" 
        style="width: 100%"
        border
      >
        <el-table-column prop="student_id" label="学号" width="120" sortable />
        <el-table-column prop="name" label="姓名" sortable />
        <el-table-column prop="class_name" label="班级" sortable />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="phone" label="电话" width="150" />
        <el-table-column label="操作" width="200" fixed="right">
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
        v-if="students.length === 0 && !loading" 
        description="暂无学生数据"
      />
      
      <!-- Pagination -->
      <div class="pagination-container">
        <el-pagination
          v-if="pagination.total > 0"
          background
          layout="total, sizes, prev, pager, next"
          :total="pagination.total"
          :page-size="pagination.pageSize"
          :current-page="pagination.page"
          :page-sizes="[10, 20, 50, 100]"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- Delete Confirmation Dialog -->
    <el-dialog
      v-model="deleteDialog.visible"
      title="确认删除"
      width="30%"
    >
      <span>确定要删除学生 "{{ deleteDialog.student?.name }}" 吗？</span>
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
import { ElMessage } from 'element-plus'

export default {
  name: 'StudentList',
  setup() {
    const store = useStore()
    const loading = ref(false)
    
    // Get students from store
    const students = computed(() => store.getters['students/allStudents'])
    const pagination = computed(() => store.getters['students/pagination'])
    const classOptions = computed(() => store.getters['classes/classOptions'])
    
    // Filter state
    const filters = reactive({
      classId: '',
      studentId: '',
      name: ''
    })
    
    // Delete dialog state
    const deleteDialog = reactive({
      visible: false,
      student: null,
      loading: false
    })
    
    // Fetch students and classes
    const fetchData = async () => {
      loading.value = true
      try {
        await Promise.all([
          store.dispatch('students/fetchStudents'),
          store.dispatch('classes/fetchClasses')
        ])
      } catch (error) {
        console.error('Error fetching data:', error)
        ElMessage.error('Failed to load data')
      } finally {
        loading.value = false
      }
    }
    
    // Handle filter changes
    const handleFilterChange = () => {
      store.dispatch('students/setFilters', {
        classId: filters.classId,
        studentId: filters.studentId,
        name: filters.name
      })
    }
    
    // Reset filters
    const resetFilters = () => {
      filters.classId = ''
      filters.studentId = ''
      filters.name = ''
      handleFilterChange()
    }
    
    // Pagination handlers
    const handleSizeChange = (size) => {
      store.dispatch('students/setPagination', { pageSize: size })
    }
    
    const handleCurrentChange = (page) => {
      store.dispatch('students/setPage', page)
    }
    
    // Delete student
    const confirmDelete = (student) => {
      deleteDialog.student = student
      deleteDialog.visible = true
    }
    
    const deleteStudent = async () => {
      if (!deleteDialog.student) return
      
      deleteDialog.loading = true
      try {
        await store.dispatch('students/deleteStudent', deleteDialog.student.id)
        ElMessage.success('Student deleted successfully')
        deleteDialog.visible = false
      } catch (error) {
        console.error('Error deleting student:', error)
        ElMessage.error('Failed to delete student')
      } finally {
        deleteDialog.loading = false
      }
    }
    
    // Fetch data on component mount
    onMounted(fetchData)
    
    return {
      students,
      loading,
      pagination,
      filters,
      classOptions,
      deleteDialog,
      handleFilterChange,
      resetFilters,
      handleSizeChange,
      handleCurrentChange,
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

.filter-card {
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style> 