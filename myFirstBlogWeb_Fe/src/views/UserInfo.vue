<template>
  <div class="content-area">
    <h2><i class="fas fa-user-circle"></i> 用户信息</h2>

    <div v-if="isLoggedIn">
      <div class="card">
        <div class="card-title">{{ username }}</div>
        <div class="card-content">
          <p>已发表文章: {{ articlesCount }} 篇</p>
        </div>
      </div>

      <h3 style="margin-top: 16px;">最近发布的文章（最多 10 条）</h3>
      <div v-if="articles.length > 0">
        <div class="card" v-for="a in articles" :key="a.article_id">
          <div class="card-title"><router-link :to="{ name: 'article-detail', params: { id: a.article_id } }">{{ a.title }}</router-link></div>
          <div class="card-meta">
            <span>{{ formatDate(a.create_time) }}</span>
          </div>
        </div>
      </div>
      <div v-else class="empty-state">暂无文章</div>
    </div>

    <div v-else class="empty-state">
      <p>请先登录以查看用户信息</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '@/store/user'
import api from '@/services/api'
import { formatDate } from '@/utils'

const userStore = useUserStore()
const isLoggedIn = computed(() => userStore.isLoggedIn)
const username = computed(() => userStore.username)
const articlesCount = computed(() => userStore.articlesCount)

const articles = ref([])

onMounted(async () => {
  if (!userStore.isLoggedIn) return
  // refresh count
  await userStore.refreshArticlesCount()
  try {
    const res = await api.getArticles(1, 100)
    const list = (res.data.articles || []).filter(a => String(a.author_id) === String(userStore.userId))
    articles.value = list.slice(0, 10)
  } catch (err) {
    console.error('load user articles failed', err)
    articles.value = []
  }
})

</script>

<style scoped>
/* reuse global styles */
</style>
