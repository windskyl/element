<template>
  <div class="content-area">
    <h2><i class="fas fa-user-plus"></i> 用户注册</h2>
    <div class="alert alert-error" v-if="registerError">
      {{ registerError }}
    </div>
    <RegisterForm @submit="register" />
    <p style="text-align: center; margin-top: 20px;">
      已有账号？<router-link :to="{ name: 'login' }">去登录</router-link>
    </p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import RegisterForm from '@/components/RegisterForm.vue'

const userStore = useUserStore()
const router = useRouter()
const registerError = ref('')

const register = async (userData) => {
  registerError.value = ''
  const success = await userStore.register(userData)
  if (success) {
    router.push({ name: 'home' })
  } else {
    registerError.value = '注册失败，请检查信息或验证码'
  }
}
</script>