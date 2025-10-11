<template>
  <div class="content-area">
    <button class="btn" @click="router.push({ name: 'home' })">
      <i class="fas fa-arrow-left"></i> 返回
    </button>
    
    <div v-if="article">
      <h2>{{ article.title }}</h2>
      <div class="card-content">{{ article.content }}</div>
      <div class="card-meta">
        <span>作者: {{ article.author_id }}</span>
        <span>发布于: {{ formatDate(article.create_time) }}</span>
      </div>
      <div class="card-actions" v-if="userStore.isLoggedIn && userStore.userId === article.author_id">
        <button class="btn" @click="editArticle">
          <i class="fas fa-edit"></i> 编辑
        </button>
        <button class="btn btn-danger" @click="deleteArticle">
          <i class="fas fa-trash"></i> 删除
        </button>
      </div>
      <hr>
      <h3>评论</h3>
      <CommentList :comments="comments" />
      <div v-if="userStore.isLoggedIn" class="form-group">
        <textarea v-model="newComment" placeholder="写下你的评论..." rows="3"></textarea>
        <button class="btn btn-block" @click="postComment">发表评论</button>
      </div>
      <div v-else class="alert alert-error">
        请登录后发表评论
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import CommentList from '@/components/CommentList.vue'
import { formatDate } from '@/utils'
import api from '@/services/api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const article = ref(null)
const comments = ref([])
const newComment = ref('')

const articleId = computed(() => route.params.id)

const loadArticle = async () => {
  const res = await api.getArticle(articleId.value)
  article.value = res.data
}

const loadComments = async () => {
  const res = await api.getComments(articleId.value)
  comments.value = res.data.comments
}

const editArticle = () => {
  router.push({ name: 'edit-article', params: { id: articleId.value } })
}

const deleteArticle = async () => {
  if (confirm(`确定要删除文章吗？`)) {
    await api.deleteArticle(articleId.value)
    router.push({ name: 'home' })
  }
}

const postComment = async () => {
  if (!newComment.value.trim()) return
  await api.createComment({
    article_id: articleId.value,
    content: newComment.value
  })
  newComment.value = ''
  loadComments()
}

onMounted(() => {
  loadArticle()
  loadComments()
})
</script>