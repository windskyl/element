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
      <label>验证码</label>
      <div class="captcha-container">
        <input type="text" v-model="form.captcha" placeholder="输入验证码" required>
        <img :src="captchaImage" class="captcha-img" @click="refreshCaptcha">
      </div>
    </div>
    
    <button type="submit" class="btn btn-block">
      <i class="fas fa-sign-in-alt"></i> 登录
    </button>
  </form>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/services/api'

const form = ref({
  username: '',
  password: '',
  captcha: ''
})

const captchaImage = ref('')
const captchaId = ref('')

const emit = defineEmits(['submit'])

const getCaptcha = async () => {
  const res = await api.getCaptcha()
  captchaImage.value = res.data.image_base64
  captchaId.value = res.data.captcha_id
}

onMounted(getCaptcha)

const refreshCaptcha = () => {
  getCaptcha()
}

const handleSubmit = () => {
  emit('submit', { ...form.value, captcha_id: captchaId.value })
}
</script>