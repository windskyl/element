<template>
  <form @submit.prevent="handleSubmit">
    <div class="form-group">
      <label>用户名</label>
      <input type="text" v-model="form.username" placeholder="输入用户名" required>
    </div>
    
    <div class="form-group">
      <label>密码</label>
      <input type="password" v-model="form.password" placeholder="输入密码" required>
    </div>
    
    <div class="form-group">
      <label>确认密码</label>
      <input type="password" v-model="form.confirmPassword" placeholder="再次输入密码" required>
    </div>
    
    <div class="form-group">
      <label>验证码</label>
      <div class="captcha-container">
        <input type="text" v-model="form.captcha" placeholder="输入验证码" required>
        <img :src="captchaImage" class="captcha-img" @click="refreshCaptcha">
        <span class="captcha-indicator" :class="{ ok: isValid === true, bad: isValid === false }"></span>
      </div>
    </div>
    
    <button type="submit" class="btn btn-block">
      <i class="fas fa-user-plus"></i> 注册
    </button>
  </form>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import api from '@/services/api'

const form = ref({
  username: '',
  password: '',
  confirmPassword: '',
  captcha: ''
})

const captchaImage = ref('')
const captchaId = ref('')
const isValid = ref(null)
let verifyTimeout = null

const emit = defineEmits(['submit'])

const getCaptcha = async () => {
  const res = await api.getCaptcha()
  console.log('getCaptcha response:', res.data)
  captchaImage.value = res.data.image_base64
  captchaId.value = res.data.id || res.data.captcha_id
  isValid.value = null
}

onMounted(getCaptcha)

const refreshCaptcha = () => {
  getCaptcha()
}

// watch captcha input for realtime validation
watch(() => form.value.captcha, (val) => {
  isValid.value = null
  if (verifyTimeout) clearTimeout(verifyTimeout)
  verifyTimeout = setTimeout(async () => {
    const input = (val || '').trim()
    if (!input) { isValid.value = null; return }
    try {
      const res = await api.verifyCaptcha({ captcha_id: captchaId.value, captcha: input })
      isValid.value = !!(res.data && res.data.ok)
    } catch (err) {
      console.error('captcha check error', err)
      isValid.value = false
    }
  }, 400)
})

const handleSubmit = () => {
  emit('submit', { ...form.value, captcha_id: captchaId.value })
}
</script>

<style scoped>
.captcha-img{ height:40px; border:1px solid #ddd; border-radius:4px; cursor:pointer }
.captcha-indicator{ width:12px; height:12px; border-radius:50%; display:inline-block; margin-left:8px; border:1px solid #ccc }
.captcha-indicator.ok{ background: #2ecc71; border-color: #2ecc71 }
.captcha-indicator.bad{ background: #e74c3c; border-color: #e74c3c }
</style>