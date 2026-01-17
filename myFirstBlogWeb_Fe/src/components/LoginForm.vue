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
    
    <!--
    <div class="form-group">
      <label>验证码</label>
      <div class="captcha-container">
        <input type="text" v-model="form.captcha" placeholder="输入验证码" required>
      </div>
    </div>
    -->
    
    <button type="submit" class="btn btn-block">
      <i class="fas fa-sign-in-alt"></i> 登录
    </button>
  </form>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/services/api'

const props = defineProps({
  initialCaptcha: { type: String, default: '' },
  initialCaptchaId: { type: String, default: '' }
})

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
  // support multiple possible response field names for compatibility
  captchaImage.value = res.data.image_base64 || res.data.image || res.data.b64s || ''
  captchaId.value = res.data.captcha_id || res.data.id || ''
}

// only fetch a fresh captcha if not already provided by parent
onMounted(() => {
  if (!props.initialCaptcha && !props.initialCaptchaId) getCaptcha()
})

const refreshCaptcha = () => {
  getCaptcha()
}

const handleSubmit = () => {
  // 如果父组件提供了已验证的图片验证码，则优先使用它；否则使用本地输入（若未填写则为空）
  const payload = {
    ...form.value,
    captcha_id: props.initialCaptchaId || captchaId.value,
    captcha: props.initialCaptcha || form.value.captcha
  }
  emit('submit', payload)
}
</script>