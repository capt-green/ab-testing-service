import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token'))
  const user = ref(null)

  const isAuthenticated = computed(() => !!token.value)

  async function login(email, password) {
    try {
      const response = await axios.post('/api/auth/login', { email, password })
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('token', token.value)
      
      // Set token for all future requests
      axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
      
      router.push('/')
    } catch (error) {
      throw error.response?.data?.error || 'Login failed'
    }
  }

  async function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    delete axios.defaults.headers.common['Authorization']
    router.push('/login')
  }

  async function register(email, password) {
    try {
      const response = await axios.post('/api/auth/register', { email, password })
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('token', token.value)
      
      // Set token for all future requests
      axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
      
      router.push('/')
    } catch (error) {
      throw error.response?.data?.error || 'Registration failed'
    }
  }

  // Initialize axios header if token exists
  if (token.value) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  return {
    token,
    user,
    isAuthenticated,
    login,
    logout,
    register
  }
})
