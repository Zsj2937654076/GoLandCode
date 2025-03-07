<template>
  <div class="dashboard-container">
    <h1 class="page-title">仪表盘</h1>
    
    <el-row :gutter="20">
      <!-- Students Card -->
      <el-col :xs="24" :sm="12" :md="8">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <h2>学生</h2>
            </div>
          </template>
          <div class="card-content">
            <div class="stat-value">{{ studentCount }}</div>
            <div class="stat-label">学生总数</div>
            <el-button type="primary" @click="$router.push('/students')">
              查看学生
            </el-button>
          </div>
        </el-card>
      </el-col>
      
      <!-- Classes Card -->
      <el-col :xs="24" :sm="12" :md="8">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <h2>班级</h2>
            </div>
          </template>
          <div class="card-content">
            <div class="stat-value">{{ classCount }}</div>
            <div class="stat-label">班级总数</div>
            <el-button type="primary" @click="$router.push('/classes')">
              查看班级
            </el-button>
          </div>
        </el-card>
      </el-col>
      
      <!-- User Card -->
      <el-col :xs="24" :sm="12" :md="8">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <h2>用户信息</h2>
            </div>
          </template>
          <div class="card-content">
            <div class="user-info">
              <p><strong>用户名：</strong> {{ user?.username }}</p>
              <p><strong>邮箱：</strong> {{ user?.email }}</p>
              <p><strong>角色：</strong> {{ user?.role === 'admin' ? '管理员' : '普通用户' }}</p>
            </div>
            <el-button type="primary" @click="$router.push('/profile')">
              查看资料
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- Recent Students -->
    <el-card class="mt-20">
      <template #header>
        <div class="card-header">
          <h2>最近添加的学生</h2>
        </div>
      </template>
      <el-table :data="recentStudents" v-loading="loading" style="width: 100%">
        <el-table-column prop="student_id" label="学号" width="120" />
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="class_name" label="班级" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button 
              size="small" 
              type="primary" 
              @click="$router.push(`/students/${scope.row.id}`)"
            >
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'

export default {
  name: 'Dashboard',
  setup() {
    const store = useStore()
    const loading = ref(false)
    const recentStudents = ref([])
    
    // Get user from store
    const user = computed(() => store.getters['auth/user'])
    
    // Computed properties for statistics
    const studentCount = ref(0)
    const classCount = ref(0)
    
    // Fetch dashboard data
    const fetchDashboardData = async () => {
      loading.value = true
      try {
        // Fetch classes
        await store.dispatch('classes/fetchClasses')
        classCount.value = store.getters['classes/allClasses'].length
        
        // Fetch recent students (first page with small page size)
        const params = {
          page: 1,
          page_size: 5
        }
        const response = await store.dispatch('students/fetchStudents')
        recentStudents.value = store.getters['students/allStudents']
        studentCount.value = store.getters['students/pagination'].total
      } catch (error) {
        console.error('Error fetching dashboard data:', error)
      } finally {
        loading.value = false
      }
    }
    
    // Fetch data when component is mounted
    onMounted(fetchDashboardData)
    
    return {
      user,
      loading,
      studentCount,
      classCount,
      recentStudents
    }
  }
}
</script>

<style scoped>
.dashboard-card {
  height: 100%;
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
  font-size: 18px;
}

.card-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
}

.stat-value {
  font-size: 36px;
  font-weight: bold;
  color: #409EFF;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 16px;
  color: #606266;
  margin-bottom: 15px;
}

.user-info {
  width: 100%;
  margin-bottom: 15px;
}

.user-info p {
  margin: 5px 0;
}
</style> 