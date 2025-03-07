import axios from 'axios'

// API base URL
const API_URL = process.env.VUE_APP_API_URL || 'http://localhost:8080/api'

const state = {
  classes: [],
  class: null,
  classStudents: []
}

const getters = {
  allClasses: state => state.classes,
  classById: state => id => state.classes.find(cls => cls.id === id),
  currentClass: state => state.class,
  classStudents: state => state.classStudents,
  // Create a map of classId to className for easy access
  classOptions: state => state.classes.map(cls => ({
    value: cls.id,
    label: cls.name
  }))
}

const actions = {
  // Fetch all classes
  async fetchClasses({ commit }) {
    try {
      const response = await axios.get(`${API_URL}/classes`)
      commit('SET_CLASSES', response.data)
      return response.data
    } catch (error) {
      console.error('Error fetching classes:', error)
      throw error
    }
  },
  
  // Fetch a single class by ID
  async fetchClass({ commit }, id) {
    try {
      const response = await axios.get(`${API_URL}/classes/${id}`)
      commit('SET_CLASS', response.data)
      return response.data
    } catch (error) {
      console.error(`Error fetching class ${id}:`, error)
      throw error
    }
  },
  
  // Fetch students in a class
  async fetchClassStudents({ commit }, classId) {
    try {
      const response = await axios.get(`${API_URL}/classes/${classId}/students`)
      commit('SET_CLASS_STUDENTS', response.data)
      return response.data
    } catch (error) {
      console.error(`Error fetching students for class ${classId}:`, error)
      throw error
    }
  },
  
  // Create a new class
  async createClass({ commit }, classData) {
    try {
      const response = await axios.post(`${API_URL}/classes`, classData)
      commit('ADD_CLASS', response.data)
      return response.data
    } catch (error) {
      console.error('Error creating class:', error)
      throw error
    }
  },
  
  // Update an existing class
  async updateClass({ commit }, classData) {
    try {
      const response = await axios.put(`${API_URL}/classes/${classData.id}`, classData)
      commit('UPDATE_CLASS', response.data)
      return response.data
    } catch (error) {
      console.error(`Error updating class ${classData.id}:`, error)
      throw error
    }
  },
  
  // Delete a class
  async deleteClass({ commit }, id) {
    try {
      await axios.delete(`${API_URL}/classes/${id}`)
      commit('DELETE_CLASS', id)
      return id
    } catch (error) {
      console.error(`Error deleting class ${id}:`, error)
      throw error
    }
  }
}

const mutations = {
  SET_CLASSES(state, classes) {
    state.classes = classes
  },
  SET_CLASS(state, classData) {
    state.class = classData
  },
  SET_CLASS_STUDENTS(state, students) {
    state.classStudents = students
  },
  ADD_CLASS(state, classData) {
    state.classes.push(classData)
  },
  UPDATE_CLASS(state, updatedClass) {
    const index = state.classes.findIndex(c => c.id === updatedClass.id)
    if (index !== -1) {
      state.classes.splice(index, 1, updatedClass)
    }
    if (state.class && state.class.id === updatedClass.id) {
      state.class = updatedClass
    }
  },
  DELETE_CLASS(state, id) {
    state.classes = state.classes.filter(classData => classData.id !== id)
    if (state.class && state.class.id === id) {
      state.class = null
    }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
} 