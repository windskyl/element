<template>
  <div>
    <h2><i class="fas fa-newspaper"></i> 最新文章</h2>
    
    <div v-if="loading" class="loading">
      <i class="fas fa-spinner fa-spin"></i> 正在加载...
    </div>
    
    <div v-else-if="articles.length === 0" class="empty-state">
      <i class="fas fa-folder-open"></i> 暂无文章
    </div>
    
    <div v-else>
      <ArticleCard
        v-for="article in articles"
        :key="article.article_id"
        :article="article"
        @view="viewArticle"
        @delete="deleteArticle"
      />
      <div class="pagination">
        <button
          v-for="page in totalPages"
          :key="page"
          class="page-btn"
          :class="{ active: page === currentPage }"
          @click="changePage(page)"
        >{{ page }}</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import ArticleCard from '@/components/ArticlaCard.vue'
import api from '@/services/api'

const router = useRouter()
const articles = ref([])
const loading = ref(false)
const currentPage = ref(1)
const limit = ref(10)
const total = ref(0)

const totalPages = computed(() => Math.ceil(total.value / limit.value))

const loadArticles = async (page = 1) => {
  loading.value = true
  try {
    const res = await api.getArticles(page, limit.value)
    articles.value = res.data.articles
    currentPage.value = res.data.page
    total.value = res.data.total
  } catch (error) {
    console.error('加载文章失败:', error)
  } finally {
    loading.value = false
  }
}

const viewArticle = (articleId) => {
  router.push({ name: 'article-detail', params: { id: articleId } })
}

const deleteArticle = async (articleId) => {
  if (confirm(`确定要删除文章吗？`)) {
    await api.deleteArticle(articleId)
    loadArticles(currentPage.value)
  }
}

const changePage = (page) => {
  loadArticles(page)
}

onMounted(() => {
  loadArticles()
})
</script>