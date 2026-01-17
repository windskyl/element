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
  const articlesCount = ref(0)
  
  const setUser = (userData) => {
    token.value = userData.token || ''
    userId.value = userData.userID || ''
    username.value = userData.username || ''
    // keep articlesCount in sync if provided
    if (typeof userData.articlesCount !== 'undefined') {
      articlesCount.value = userData.articlesCount
    }
    
    localStorage.setItem('token', token.value)
    localStorage.setItem('userId', userId.value)
    localStorage.setItem('username', username.value)
    
    isLoggedIn.value = !!token.value
  }

  const refreshArticlesCount = async () => {
    if (!token.value) {
      articlesCount.value = 0
      return
    }
    try {
      const res = await api.getMyArticlesCount()
      articlesCount.value = res.data.count || 0
    } catch (err) {
      console.error('refreshArticlesCount error', err)
      articlesCount.value = 0
    }
  }

  // 检查 token 的有效期（基于 JWT exp 字段），避免使用已过期的 token 发起受保护请求
  const isTokenValid = () => {
    if (!token.value) return false
    try {
      const parts = token.value.split('.')
      if (parts.length < 2) return false
      const payload = JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')))
      if (!payload.exp) return true
      return payload.exp > Math.floor(Date.now() / 1000)
    } catch (err) {
      console.error('isTokenValid parse error', err)
      return false
    }
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
          username: credentials.username,
          password_hash: md5(credentials.password),
          captcha: credentials.captcha,
          captcha_id: credentials.captcha_id
        }
        console.log('login payload:', payload)
        const res = await api.login(payload)
        console.log('login response:', res)
        if (res.data.token) {
          setUser(res.data)
          // refresh the user's article count after login
          await refreshArticlesCount()
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
        username: userData.username,
        password_hash: md5(userData.password),
        captcha: userData.captcha,
        captcha_id: userData.captcha_id
      }
      console.log('register payload:', payload)
      const res = await api.register(payload)
      console.log('register response:', res)
      if (res.data.userID) {
        setUser({ ...userData, userID: res.data.userID, token: '' })
        // after registering, count is zero
        articlesCount.value = 0
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
    articlesCount,
    setUser,
    clearUser,
    login,
    register,
    logout,
    refreshArticlesCount,
    isTokenValid
  }
})