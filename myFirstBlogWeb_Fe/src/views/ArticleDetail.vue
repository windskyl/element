<template>
  <div class="content-area">
    <div class="detail-controls">
      <div class="controls-row">
        <div class="controls-left">
          <button class="btn" @click="router.push({ name: 'home' })">
            <i class="fas fa-arrow-left"></i> 返回
          </button>
        </div>
        <div class="controls-right">
          <button v-if="userStore.isLoggedIn && isAuthor" class="btn btn-subtle" @click="editArticle">
            <i class="fas fa-edit"></i> 编辑
          </button>
          <button v-if="userStore.isLoggedIn && isAuthor" class="btn btn-subtle-danger" @click="deleteArticle">
            <i class="fas fa-trash"></i> 删除
          </button>
        </div>
      </div>
      <div class="title-row">
        <h2 class="detail-title">{{ article ? article.title : '' }}</h2>
      </div>
    </div>
    
    <div v-if="article">
      <div class="card-content" v-html="article.content"></div>
      <div class="card-meta">
        <span>{{ authorName || article.author_username || '' }}</span>
        <span>发布于: {{ formatDate(article.create_time) }}</span>
      </div>
      <!-- 编辑/删除 按钮已移动到页面顶部 actions 区，不再在正文重复显示 -->
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
    <div v-else>
      <div v-if="errorMessage" class="alert alert-error">加载文章失败：{{ errorMessage }}</div>
      <div v-else class="empty-state">文章正在加载中...</div>
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
const authorName = ref('')
const comments = ref([])
const newComment = ref('')
const errorMessage = ref('')

const isAuthor = computed(() => {
  try {
    return article.value && String(userStore.userId) === String(article.value.author_id)
  } catch (e) {
    return false
  }
})

const articleId = computed(() => route.params.id)

const loadArticle = async () => {
  errorMessage.value = ''
  try {
    const res = await api.getArticle(articleId.value)
    // 后端现在返回 { article: {...}, author_username: 'name' }
    article.value = res.data.article || res.data
    authorName.value = res.data.author_username || ''
    console.debug('loadArticle debug:', {
      articleAuthorId: article.value?.author_id,
      articleAuthorIdType: typeof article.value?.author_id,
      currentUserId: userStore.userId,
      currentUserIdType: typeof userStore.userId
    })
  } catch (err) {
    console.error('加载文章失败', err)
    // 尝试从后端响应中获取错误信息
    errorMessage.value = err?.response?.data?.error || err?.response?.data?.message || err.message || '未知错误'
    article.value = null
  }
}

const loadComments = async () => {
  try {
    const res = await api.getComments(articleId.value)
    comments.value = res.data.comments
  } catch (err) {
    console.error('加载评论失败', err)
    comments.value = []
  }
}

const editArticle = () => {
  router.push({ name: 'edit-article', params: { id: articleId.value } })
}

const deleteArticle = async () => {
  if (confirm(`确定要删除文章吗？`)) {
    await api.deleteArticle(articleId.value)
    // 刷新侧边栏计数并返回首页
    await userStore.refreshArticlesCount()
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

<style scoped>
.detail-controls{ display:flex; flex-direction:column; gap:8px; margin-bottom:12px }
.controls-row{ display:flex; justify-content:space-between; align-items:center }
.controls-left{ display:flex; align-items:center }
.controls-right{ display:flex; gap:8px; align-items:center }
.title-row{ display:flex; justify-content:center }
.detail-title{ margin:0; font-size:1.6rem; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; text-align:center }
.btn-subtle{ background:transparent; border:1px solid #ddd; color:#444; padding:6px 10px; border-radius:4px }
.btn-subtle:hover{ background:#f5f7fb }
.btn-subtle-danger{ background:transparent; border:1px solid #f2c6c6; color:#b12b2b; padding:6px 10px; border-radius:4px }
.btn-subtle-danger:hover{ background:#fff5f5 }

/* 让文章内容中的图片响应式，最大宽度不超过容器宽度 */
.card-content img{ max-width:100%; height:auto; display:block; margin:12px auto }
.card-content iframe{ max-width:100%; }
</style>