<template>
  <div class="card">
    <div class="card-title">{{ article.title }}</div>
    <div class="card-content">{{ truncate(article.content, 150) }}</div>
    <div class="card-meta">
      <span>作者: {{ article.author_id }}</span>
      <span>发布于: {{ formatDate(article.create_time) }}</span>
    </div>
    <div class="card-actions">
      <button class="btn" @click="$emit('view', article.article_id)">
        <i class="fas fa-eye"></i> 查看详情
      </button>
      <button v-if="isLoggedIn" class="btn btn-danger" @click="$emit('delete', article.article_id)">
        <i class="fas fa-trash"></i> 删除
      </button>
    </div>
  </div>
</template>

<script setup>
import { defineProps, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { formatDate, truncate } from '@/utils'

const props = defineProps({
  article: {
    type: Object,
    required: true
  }
})

const userStore = useUserStore()
const isLoggedIn = computed(() => userStore.isLoggedIn)
</script>