<template>
  <div class="content-area">
    <h2><i class="fas fa-plus"></i> 创建新文章</h2>
    <!-- 将错误和提交中状态传给子组件，让父组件负责真实的网络请求与导航 -->
    <ArticleForm @submitted="onSubmitted" :error="errorMessage" :isSubmitting="isSubmitting" />
    <div v-if="errorMessage" class="mt-2 text-red-600">{{ errorMessage }}</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import ArticleForm from '@/components/ArticleForm.vue'
import api from '@/services/api'
import { useUserStore } from '@/store/user'

const router = useRouter()
const errorMessage = ref('')
const isSubmitting = ref(false)
const userStore = useUserStore()

const onSubmitted = async (formData) => {
  errorMessage.value = ''
  isSubmitting.value = true
  try {
    // 使用封装好的 createArticle，确保会走 /ctx/ 前缀并触发 token 拦截
    await api.createArticle(formData)
    // 成功后刷新用户侧栏计数并跳转到首页
    await userStore.refreshArticlesCount()
    router.push({ name: 'home' })
  } catch (err) {
    errorMessage.value = err.response?.data?.error || err.response?.data?.reason || '创建文章失败，请重试'
  } finally {
    isSubmitting.value = false
  }
}
</script>