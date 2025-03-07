import axios from 'axios'

// API base URL
const API_URL = process.env.VUE_APP_API_URL || 'http://localhost:8080/api'

// Helper to set auth header
const setAuthHeader = (token) => {
  if (token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
  } else {
    delete axios.defaults.headers.common['Authorization']
  }
}

const state = {
  token: localStorage.getItem('token') || '',
  user: JSON.parse(localStorage.getItem('user')) || null,
  status: ''
}

const getters = {
  isAuthenticated: state => !!state.token,
  authStatus: state => state.status,
  user: state => state.user
}

const actions = {
  // Login action
  async login({ commit }, user) {
    commit('AUTH_REQUEST')
    try {
      const response = await axios.post(`${API_URL}/auth/login`, user)
      const { token, user: userData } = response.data
      
      // Store token in localStorage
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(userData))
      
      // Set auth header
      setAuthHeader(token)
      
      commit('AUTH_SUCCESS', { token, user: userData })
      return response
    } catch (error) {
      commit('AUTH_ERROR', error)
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      throw error
    }
  },
  
  // Logout action
  async logout({ commit }) {
    try {
      // Call logout endpoint if needed
      await axios.post(`${API_URL}/auth/logout`)
    } catch (error) {
      console.error('Logout error:', error)
    }
    
    // Clean up regardless of the API call result
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setAuthHeader(null)
    commit('AUTH_LOGOUT')
  },
  
  // Restore session from token
  async restoreSession({ commit }, token) {
    setAuthHeader(token)
    commit('AUTH_REQUEST')
    
    try {
      const response = await axios.get(`${API_URL}/auth/profile`)
      const userData = response.data
      
      // Update user data in localStorage
      localStorage.setItem('user', JSON.stringify(userData))
      
      commit('AUTH_SUCCESS', { token, user: userData })
      return response
    } catch (error) {
      commit('AUTH_ERROR', error)
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      setAuthHeader(null)
      throw error
    }
  },
  
  // Change password
  async changePassword({ commit }, { oldPassword, newPassword }) {
    try {
      const response = await axios.post(`${API_URL}/auth/change-password`, {
        old_password: oldPassword,
        new_password: newPassword
      })
      return response
    } catch (error) {
      throw error
    }
  }
}

const mutations = {
  AUTH_REQUEST(state) {
    state.status = 'loading'
  },
  AUTH_SUCCESS(state, { token, user }) {
    state.status = 'success'
    state.token = token
    state.user = user
  },
  AUTH_ERROR(state) {
    state.status = 'error'
    state.token = ''
    state.user = null
  },
  AUTH_LOGOUT(state) {
    state.status = ''
    state.token = ''
    state.user = null
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
} 