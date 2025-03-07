import axios from 'axios'

// API base URL
const API_URL = process.env.VUE_APP_API_URL || 'http://localhost:8080/api'

const state = {
  students: [],
  student: null,
  pagination: {
    total: 0,
    page: 1,
    pageSize: 10
  },
  filters: {
    classId: '',
    studentId: '',
    name: ''
  }
}

const getters = {
  allStudents: state => state.students,
  studentById: state => id => state.students.find(student => student.id === id),
  currentStudent: state => state.student,
  pagination: state => state.pagination,
  filters: state => state.filters
}

const actions = {
  // Fetch students with pagination and filters
  async fetchStudents({ commit, state }) {
    try {
      const { page, pageSize } = state.pagination
      const { classId, studentId, name } = state.filters
      
      // Build query params
      let params = `page=${page}&page_size=${pageSize}`
      if (classId) params += `&class_id=${classId}`
      if (studentId) params += `&student_id=${studentId}`
      if (name) params += `&name=${name}`
      
      const response = await axios.get(`${API_URL}/students?${params}`)
      const { data, pagination } = response.data
      
      commit('SET_STUDENTS', data)
      commit('SET_PAGINATION', pagination)
      
      return data
    } catch (error) {
      console.error('Error fetching students:', error)
      throw error
    }
  },
  
  // Fetch a single student by ID
  async fetchStudent({ commit }, id) {
    try {
      const response = await axios.get(`${API_URL}/students/${id}`)
      commit('SET_STUDENT', response.data)
      return response.data
    } catch (error) {
      console.error(`Error fetching student ${id}:`, error)
      throw error
    }
  },
  
  // Create a new student
  async createStudent({ commit }, student) {
    try {
      const response = await axios.post(`${API_URL}/students`, student)
      commit('ADD_STUDENT', response.data)
      return response.data
    } catch (error) {
      console.error('Error creating student:', error)
      throw error
    }
  },
  
  // Update an existing student
  async updateStudent({ commit }, student) {
    try {
      const response = await axios.put(`${API_URL}/students/${student.id}`, student)
      commit('UPDATE_STUDENT', response.data)
      return response.data
    } catch (error) {
      console.error(`Error updating student ${student.id}:`, error)
      throw error
    }
  },
  
  // Delete a student
  async deleteStudent({ commit }, id) {
    try {
      await axios.delete(`${API_URL}/students/${id}`)
      commit('DELETE_STUDENT', id)
      return id
    } catch (error) {
      console.error(`Error deleting student ${id}:`, error)
      throw error
    }
  },
  
  // Update filters
  setFilters({ commit, dispatch }, filters) {
    commit('SET_FILTERS', filters)
    // Reset to first page when filters change
    commit('SET_PAGE', 1)
    return dispatch('fetchStudents')
  },
  
  // Update pagination
  setPagination({ commit, dispatch }, pagination) {
    commit('SET_PAGINATION', pagination)
    return dispatch('fetchStudents')
  },
  
  // Set current page
  setPage({ commit, dispatch }, page) {
    commit('SET_PAGE', page)
    return dispatch('fetchStudents')
  }
}

const mutations = {
  SET_STUDENTS(state, students) {
    state.students = students
  },
  SET_STUDENT(state, student) {
    state.student = student
  },
  ADD_STUDENT(state, student) {
    state.students.unshift(student)
  },
  UPDATE_STUDENT(state, updatedStudent) {
    const index = state.students.findIndex(s => s.id === updatedStudent.id)
    if (index !== -1) {
      state.students.splice(index, 1, updatedStudent)
    }
    if (state.student && state.student.id === updatedStudent.id) {
      state.student = updatedStudent
    }
  },
  DELETE_STUDENT(state, id) {
    state.students = state.students.filter(student => student.id !== id)
    if (state.student && state.student.id === id) {
      state.student = null
    }
  },
  SET_PAGINATION(state, pagination) {
    state.pagination = { ...state.pagination, ...pagination }
  },
  SET_PAGE(state, page) {
    state.pagination.page = page
  },
  SET_FILTERS(state, filters) {
    state.filters = { ...state.filters, ...filters }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
} 
 