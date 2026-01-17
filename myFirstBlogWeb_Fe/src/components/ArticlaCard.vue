<template>
  <div class="card">
    <div class="card-title large">{{ article.title }}</div>
    <div class="card-excerpt">{{ excerpt }}</div>
    <div class="card-meta">
      <span class="meta-author">{{ article.author_username || '' }}</span>
      <span class="meta-time">{{ formatDate(article.modify_time || article.create_time) }}</span>
      <span class="meta-comments">{{ article.comment_count || article.commentCount || 0 }} 评论</span>
    </div>
    <div class="card-actions">
      <button class="btn" @click="$emit('view', article.article_id)">
        <i class="fas fa-eye"></i> 查看详情
      </button>
    </div>
  </div>
</template>

<script setup>
import { defineProps, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { formatDate } from '@/utils'

const props = defineProps({
  article: {
    type: Object,
    required: true
  }
})

const userStore = useUserStore()
const isLoggedIn = computed(() => userStore.isLoggedIn)
const excerpt = computed(() => {
  const text = props.article?.content || ''
  if (!text) return ''
  if (text.length >= 100) return text.slice(0, 100) + '...'
  const n = Math.ceil(text.length / 2)
  return text.slice(0, n) + '...'
})
</script>

<style scoped>
.card-title.large{ font-size: 1.4rem; font-weight:700; margin-bottom:8px; line-height:1.2 }
.card-excerpt{ color: #444; margin-bottom:12px; line-height:1.6; max-height:6.4rem; overflow:hidden }
.card-meta{ display:flex; gap:12px; color:#888; font-size:0.9rem; align-items:center }
.meta-author{ font-weight:600; color:#333 }
.meta-comments{ margin-left:auto }

@media (max-width: 768px) {
  .card-title.large{ font-size: 1.15rem }
  .card-excerpt{ font-size: 0.95rem }
  .card{ padding:12px }
}

@media (max-width: 480px) {
  .card-meta{ flex-direction:column; align-items:flex-start; gap:6px }
  .meta-comments{ margin-left:0 }
}
</style>