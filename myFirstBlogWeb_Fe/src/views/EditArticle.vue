<template>
  <div class="content-area">
    <h2><i class="fas fa-edit"></i> 编辑文章</h2>
    <ArticleForm :article="article" :isEditMode="true" @submitted="onSubmitted" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/services/api'

const route = useRoute()
const router = useRouter()
const article = ref(null)

onMounted(async () => {
  const res = await api.getArticle(route.params.id)
  article.value = res.data
})

const onSubmitted = async (formData) => {
  await api.updateArticle(route.params.id, formData)
  router.push({ name: 'article-detail', params: { id: route.params.id } })
}
</script>