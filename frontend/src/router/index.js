import { createRouter, createWebHistory } from 'vue-router'
import store from '../store'

// Lazy load components
const Login = () => import('../views/auth/Login.vue')
const Dashboard = () => import('../views/Dashboard.vue')
const StudentList = () => import('../views/students/StudentList.vue')
const StudentForm = () => import('../views/students/StudentForm.vue')
const StudentDetail = () => import('../views/students/StudentDetail.vue')
const ClassList = () => import('../views/classes/ClassList.vue')
const ClassForm = () => import('../views/classes/ClassForm.vue')
const ClassDetail = () => import('../views/classes/ClassDetail.vue')
const Profile = () => import('../views/auth/Profile.vue')
const ChangePassword = () => import('../views/auth/ChangePassword.vue')
const NotFound = () => import('../views/NotFound.vue')

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard
  },
  // Student routes
  {
    path: '/students',
    name: 'StudentList',
    component: StudentList
  },
  {
    path: '/students/new',
    name: 'CreateStudent',
    component: StudentForm
  },
  {
    path: '/students/:id',
    name: 'StudentDetail',
    component: StudentDetail,
    props: true
  },
  {
    path: '/students/:id/edit',
    name: 'EditStudent',
    component: StudentForm,
    props: route => ({ id: parseInt(route.params.id), isEdit: true })
  },
  // Class routes
  {
    path: '/classes',
    name: 'ClassList',
    component: ClassList
  },
  {
    path: '/classes/new',
    name: 'CreateClass',
    component: ClassForm
  },
  {
    path: '/classes/:id',
    name: 'ClassDetail',
    component: ClassDetail,
    props: true
  },
  {
    path: '/classes/:id/edit',
    name: 'EditClass',
    component: ClassForm,
    props: route => ({ id: parseInt(route.params.id), isEdit: true })
  },
  // User routes
  {
    path: '/profile',
    name: 'Profile',
    component: Profile
  },
  {
    path: '/change-password',
    name: 'ChangePassword',
    component: ChangePassword
  },
  // 404 route
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

// Navigation guard to check authentication
router.beforeEach((to, from, next) => {
  // Check if the route requires authentication
  if (to.matched.some(record => record.meta.requiresAuth !== false)) {
    // This route requires auth, check if logged in
    if (!store.getters['auth/isAuthenticated']) {
      // Not logged in, redirect to login page
      next({ name: 'Login' })
    } else {
      next() // Logged in, proceed to route
    }
  } else {
    next() // No auth required, proceed to route
  }
})

export default router 