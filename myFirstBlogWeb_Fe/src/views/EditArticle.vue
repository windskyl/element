<template>
  <div class="content-area">
    <div class="edit-header">
      <div class="edit-controls-row">
        <button class="btn" @click="router.back()"><i class="fas fa-arrow-left"></i> 返回</button>
      </div>
      <div class="edit-title-row">
        <h2 class="edit-title"><i class="fas fa-edit"></i> 编辑文章</h2>
      </div>
    </div>
    <ArticleForm :article="article" :isEditMode="true" :error="errorMessage" :isSubmitting="isSubmitting" @submitted="onSubmitted" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/services/api'
import { useUserStore } from '@/store/user'
import ArticleForm from '@/components/ArticleForm.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const article = ref(null)
const errorMessage = ref('')
const isSubmitting = ref(false)

onMounted(async () => {
  try {
    const res = await api.getArticle(route.params.id)
    // 后端返回 { article: {...}, author_username: '...' } 或直接 article
    article.value = res.data.article || res.data
    console.debug('EditArticle loaded article', article.value)
  } catch (err) {
    console.error('加载文章失败', err)
    errorMessage.value = err.response?.data?.error || '加载文章失败'
  }
})

const onSubmitted = async (formData) => {
  errorMessage.value = ''
  // 弹出确认对话框
  if (!confirm('确定要保存修改吗？')) return
  isSubmitting.value = true
  try {
    await api.updateArticle(route.params.id, formData)
    // 更新成功后刷新用户计数（严格来说计数不变，但保持同步）
    await userStore.refreshArticlesCount()
    router.push({ name: 'article-detail', params: { id: route.params.id } })
  } catch (err) {
    console.error('更新文章失败', err)
    errorMessage.value = err.response?.data?.error || err.response?.data?.reason || '更新文章失败'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.edit-header{ display:flex; flex-direction:column; gap:8px; margin-bottom:12px }
.edit-controls-row{ display:flex; justify-content:flex-start }
.edit-title-row{ display:flex; justify-content:center }
.edit-title{ margin:0 }
</style>