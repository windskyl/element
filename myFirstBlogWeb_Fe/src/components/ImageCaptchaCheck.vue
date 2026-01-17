<template>
  <div class="captcha-check">
    <div class="form-group">
      <label>图片验证码校验</label>
      <div class="captcha-container">
        <input type="text" v-model="answer" placeholder="输入图片上的验证码" />
        <img :src="captchaImage" class="captcha-img" @click="refreshCaptcha" v-if="captchaImage" />
        <span class="captcha-indicator" :class="{ ok: isValid === true, bad: isValid === false }"></span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import api from '@/services/api'

const captchaImage = ref('')
const captchaId = ref('')
const answer = ref('')
const isVerified = ref(false)
const lastVerifiedAnswer = ref('')
const isValid = ref(null) // null=unknown, true=valid, false=invalid
const verifyDebounce = { id: null }
const emit = defineEmits(['verified'])

const getCaptcha = async () => {
  const res = await api.getCaptcha()
  const raw = (res && res.data) ? (res.data.image_base64 || res.data.image || res.data.b64s || '') : ''
  // ensure we have a proper data URI for the <img src>
  let final = raw || ''
  if (final && !final.startsWith('data:')) {
    final = 'data:image/png;base64,' + final
  }
  captchaImage.value = final
  captchaId.value = (res && res.data) ? (res.data.captcha_id || res.data.id || '') : ''
  // reset verified status when fetching a new captcha
  isVerified.value = false
  lastVerifiedAnswer.value = ''
  isValid.value = null
}

// fetch a captcha when component mounts so image shows immediately
onMounted(() => {
  getCaptcha()
})

// watch input and perform debounced CheckCaptcha calls
watch(() => answer.value, (val) => {
  isValid.value = null
  if (verifyDebounce.id) clearTimeout(verifyDebounce.id)
  verifyDebounce.id = setTimeout(async () => {
    if (!captchaId.value) return
    const input = (val || '').trim()
    if (!input) { isValid.value = null; return }
    try {
      const res = await api.verifyCaptcha({ captcha_id: captchaId.value, captcha: input })
      isValid.value = !!(res.data && res.data.ok)
      // emit verified payload once when a new valid input is detected
      if (isValid.value && lastVerifiedAnswer.value !== input) {
        lastVerifiedAnswer.value = input
        isVerified.value = true
        emit('verified', { captcha: input, captcha_id: captchaId.value })
      }
    } catch (err) {
      console.error('captcha check error', err)
      isValid.value = false
    }
  }, 400)
})
const refreshCaptcha = () => getCaptcha()


</script>

<style scoped>
.captcha-check .captcha-container{ display:flex; gap:8px; align-items:center }
.captcha-img{ height:40px; border:1px solid #ddd; border-radius:4px; cursor:pointer }
.captcha-indicator{ width:12px; height:12px; border-radius:50%; display:inline-block; margin-left:8px; border:1px solid #ccc }
.captcha-indicator.ok{ background: #2ecc71; border-color: #2ecc71 }
.captcha-indicator.bad{ background: #e74c3c; border-color: #e74c3c }
</style>