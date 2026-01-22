<template>
  <div class="content-area">
    <h2><i class="fas fa-sign-in-alt"></i> 用户登录</h2>
    <div class="alert alert-error" v-if="loginError">
      {{ loginError }}
    </div>
    <LoginForm @submit="login" :initialCaptcha="verifiedCaptcha" :initialCaptchaId="verifiedCaptchaId" />
    <hr />
    <!-- 新增图片验证码校验工具：仅用于在登录页面测试/验证验证码 -->
    <ImageCaptchaCheck @verified="onCaptchaVerified" />
    <p style="text-align: center; margin-top: 20px;">
      没有账号？<router-link :to="{ name: 'register' }">立即注册</router-link>
    </p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import LoginForm from '@/components/LoginForm.vue'
import ImageCaptchaCheck from '@/components/ImageCaptchaCheck.vue'

const userStore = useUserStore()
const router = useRouter()
const loginError = ref('')

// 图片验证码验证结果（ImageCaptchaCheck 会在验证成功时通知）
const verifiedCaptcha = ref('')
const verifiedCaptchaId = ref('')

const onCaptchaVerified = (payload) => {
  verifiedCaptcha.value = payload.captcha || ''
  verifiedCaptchaId.value = payload.captcha_id || ''
} 



const login = async (credentials) => {


  loginError.value = ''
  const success = await userStore.login(credentials)
  if (success) {
    router.push({ name: 'home' })
  } else {
    loginError.value = '登录失败，请检查用户名、密码或验证码'
  }
}
</script>