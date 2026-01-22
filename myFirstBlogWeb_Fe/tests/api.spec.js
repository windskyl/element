/// <reference types="vitest" />
globalThis.localStorage = {
  store: {},
  getItem(key) { return this.store[key] || ''; },
  setItem(key, value) { this.store[key] = value; },
  removeItem(key) { delete this.store[key]; },
  clear() { this.store = {}; }
}
import { describe, it, expect, beforeAll } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
setActivePinia(createPinia())
import api from '../src/services/api'
import { useUserStore } from '../src/store/user'

describe('API服务层', () => {
  let captchaId = ''
  let captchaAnswer = ''
  let token = ''
  let userID = ''
  let articleId = ''
  let uniqueUser = 'testuser_' + Math.floor(Math.random() * 10000)

  beforeAll(async () => {
    // 获取验证码
    const captchaRes = await api.getCaptcha()
    captchaId = captchaRes.data.id
    captchaAnswer = '1234'
    // 注册
    const regRes = await api.register({
      username: uniqueUser,
      password_hash: 'd41d8cd98f00b204e9800998ecf8427e',
      captcha: captchaAnswer,
      captcha_id: captchaId
    })
    userID = regRes.data.userID
    // 登录
    const loginRes = await api.login({
      username: uniqueUser,
      password_hash: 'd41d8cd98f00b204e9800998ecf8427e',
      captcha: captchaAnswer,
      captcha_id: captchaId
    })
    token = loginRes.data.token
    localStorage.setItem('token', token)
    localStorage.setItem('userId', userID)
    localStorage.setItem('username', uniqueUser)
  // 修复：同步 token 到 pinia userStore，确保 axios 请求拦截器能正确带上 token
  const userStore = useUserStore()
  userStore.token = token
  userStore.userId = userID
  userStore.username = uniqueUser
  })

  it('创建文章', async () => {
    try {
      const res = await api.createArticle({
        title: '测试文章' + Math.floor(Math.random() * 10000),
        content: '这是自动化测试创建的内容'
      })
  expect(res.data).toHaveProperty('ArticleID')
  articleId = res.data.ArticleID
    } catch (error) {
      console.error('创建文章失败:', error?.response?.data || error)
      throw error
    }
  })

  it('获取文章列表', async () => {
    try {
      const res = await api.getArticles(1, 10)
      expect(res.data).toHaveProperty('articles')
      expect(Array.isArray(res.data.articles)).toBe(true)
  // 检查刚创建的文章是否在列表中
  const found = res.data.articles.some(a => a.ArticleID === articleId)
  expect(found).toBe(true)
    } catch (error) {
      console.error('获取文章列表失败:', error?.response?.data || error)
      throw error
    }
  })

  it('发表评论', async () => {
    try {
      const res = await api.createComment({
        article_id: articleId,
        content: '单元测试评论'
      })
      expect(res.data).toHaveProperty('CommentID')
    } catch (error) {
      console.error('发表评论失败:', error?.response?.data || error)
      throw error
    }
  })

  it('删除文章', async () => {
    try {
      const res = await api.deleteArticle(articleId)
      // 删除成功后无内容返回
      expect(res.status).toBe(204)
    } catch (error) {
      console.error('删除文章失败:', error?.response?.data || error)
      throw error
    }
  })
})
