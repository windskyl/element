import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/services/api'
import md5 from 'md5'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userId = ref(localStorage.getItem('userId') || '')
  const username = ref(localStorage.getItem('username') || '')
  // 顶层不再声明 router
  
  const isLoggedIn = ref(!!token.value)
  
  const setUser = (userData) => {
    token.value = userData.token || ''
    userId.value = userData.userID || ''
    username.value = userData.username || ''
    
    localStorage.setItem('token', token.value)
    localStorage.setItem('userId', userId.value)
    localStorage.setItem('username', username.value)
    
    isLoggedIn.value = !!token.value
  }
  
  const clearUser = () => {
    token.value = ''
    userId.value = ''
    username.value = ''
    
    localStorage.removeItem('token')
    localStorage.removeItem('userId')
    localStorage.removeItem('username')
    
    isLoggedIn.value = false
  }
  
    const login = async (credentials) => {
      try {
        const payload = {
          Username: credentials.username,
          PasswordHash: md5(credentials.password),
          Captcha: credentials.captcha,
          CaptchaID: credentials.captcha_id
        }
        console.log('login payload:', payload)
        const res = await api.login(payload)
        console.log('login response:', res)
        if (res.data.token) {
          setUser(res.data)
          return true
        } else {
          return false
        }
      } catch (error) {
        console.log('login error:', error?.response?.data || error)
        return false
      }
    }
  
  const register = async (userData) => {
    try {
      const payload = {
        Username: userData.username,
        PasswordHash: md5(userData.password),
        Captcha: userData.captcha,
        CaptchaID: userData.captcha_id
      }
      console.log('register payload:', payload)
      const res = await api.register(payload)
      console.log('register response:', res)
      if (res.data.userID) {
        setUser({ ...userData, userID: res.data.userID, token: '' })
        return true
      } else {
        return false
      }
    } catch (error) {
      console.log('register error:', error?.response?.data || error)
      return false
    }
    }
  
  const logout = () => {
    clearUser()
    // 不做路由跳转，路由跳转交由页面或组件处理
  }
  
  return {
    token,
    userId,
    username,
    isLoggedIn,
    setUser,
    clearUser,
    login,
    register,
    logout
  }
})