<template>
  <div class="content-area">
    <h2><i class="fas fa-plus"></i> 创建新文章</h2>
    <ArticleForm @submitted="onSubmitted" />
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import ArticleForm from '@/components/ArticleForm.vue'
import api from '@/services/api'

const router = useRouter()

const onSubmitted = async (formData) => {
  await api.createArticle(formData)
  router.push({ name: 'home' })
}
</script>