<template>
  <div>
    <div class="comment" v-for="comment in comments" :key="comment.comment_id">
      <div class="comment-content">{{ comment.content }}</div>
      <div class="comment-meta">
        <span class="comment-author">{{ comment.author_username || '匿名' }}</span>
        <span>{{ formatDate(comment.create_time) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue'
import { formatDate } from '@/utils'

defineProps({
  comments: {
    type: Array,
    required: true
  }
})
</script>