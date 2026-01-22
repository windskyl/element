/// <reference types="vitest" />
globalThis.localStorage = {
  store: {},
  getItem(key) { return this.store[key] || ''; },
  setItem(key, value) { this.store[key] = value; },
  removeItem(key) { delete this.store[key]; },
  clear() { this.store = {}; }
}
import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUserStore } from '../src/store/user'
import api from '../src/services/api'

describe('用户Store', () => {
  let captchaId = ''
  let captchaAnswer = ''

  beforeEach(async () => {
    setActivePinia(createPinia())
    // 获取验证码
    const res = await api.getCaptcha()
    captchaId = res.data.id
    captchaAnswer = '1234' // 如后端返回内容字段可自动赋值，否则需手动填写
  })

  it('注册流程', async () => {
    const store = useUserStore()
    const uniqueUser = 'testuser_' + Math.floor(Math.random() * 10000)
    const success = await store.register({
      username: uniqueUser,
      password: '123456',
      captcha: captchaAnswer,
      captcha_id: captchaId
    })
    expect(success).toBe(true)
    expect(store.userId).toBeTruthy()
  })

  it('登录流程', async () => {
    const store = useUserStore()
    const uniqueUser = 'testuser_' + Math.floor(Math.random() * 10000)
    // 先注册
    await store.register({
      username: uniqueUser,
      password: '123456',
      captcha: captchaAnswer,
      captcha_id: captchaId
    })
    // 再登录
    const success = await store.login({
      username: uniqueUser,
      password: '123456',
      captcha: captchaAnswer,
      captcha_id: captchaId
    })
    expect(success).toBe(true)
    expect(store.token).toBeTruthy()
  })

  it('登出流程', () => {
    const store = useUserStore()
    store.setUser({ token: 'token', userID: 'uid', username: 'testuser' })
    store.logout()
    expect(store.token).toBe('')
    expect(store.userId).toBe('')
    expect(store.username).toBe('')
  })
})
