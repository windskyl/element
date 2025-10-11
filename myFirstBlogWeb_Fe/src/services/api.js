import axios from 'axios'
import { useUserStore } from '@/store/user.js'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 10000
})

// 请求拦截器 - 添加认证token（仅 /ctx/ 路径自动加 token）
api.interceptors.request.use(config => {
  const userStore = useUserStore()
  if (userStore.token && config.url.startsWith('/ctx/')) {
    config.headers.Authorization = `Bearer ${userStore.token}`
  }
  return config
}, error => {
  return Promise.reject(error)
})

// 响应拦截器 - 处理错误
api.interceptors.response.use(response => {
  return response
}, error => {
  if (error.response && error.response.status === 401) {
    const userStore = useUserStore()
    userStore.logout()
  }
  return Promise.reject(error)
})

export default {
  // 获取验证码（POST），返回 captcha_id
  getCaptcha() {
    return api.post('/captcha')
  },

  // 用户注册
  register(data) {
    return api.post('/register', data)
  },

  // 用户登录
  login(data) {
    return api.post('/login', data)
  },

  // 创建文章
  createArticle(data) {
    return api.post('/ctx/articles', data)
  },

  // 获取文章列表
  getArticles(page = 1, limit = 10) {
    return api.get(`/ctx/articles?page=${page}&limit=${limit}`)
  },

  // 获取文章详情
  getArticle(id) {
    return api.get(`/ctx/articles/${id}`)
  },

  // 更新文章
  updateArticle(id, data) {
    return api.put(`/ctx/articles/${id}`, data)
  },

  // 删除文章
  deleteArticle(id) {
    return api.delete(`/ctx/articles/${id}`)
  },

  // 获取评论列表
  getComments(articleId, page = 1, limit = 10) {
    return api.get(`/ctx/comments/${articleId}?page=${page}&limit=${limit}`)
  },

  // 发表评论（POST /ctx/comments/{articleId}）
  createComment(data) {
    // data: { article_id, content }
    return api.post(`/ctx/comments/${data.article_id}`, { content: data.content, article_id: data.article_id })
  }
}