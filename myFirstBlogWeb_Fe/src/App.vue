<template>
  <div class="container">
    <header>
      <div class="logo">
        <i class="fas fa-blog"></i>
        <span>博客系统</span>
      </div>
      <nav>
        <router-link :to="{ name: 'home' }" class="nav-btn" active-class="active">
          <i class="fas fa-home"></i> 首页
        </router-link>
        <router-link 
          v-if="isLoggedIn" 
          :to="{ name: 'create-article' }" 
          class="nav-btn" 
          active-class="active"
        >
          <i class="fas fa-plus"></i> 写文章
        </router-link>
        <router-link :to="{ name: 'about' }" class="nav-btn" active-class="active">
          <i class="fas fa-info-circle"></i> 关于
        </router-link>
      </nav>
      <div class="user-info" v-if="isLoggedIn">
        <span>欢迎, {{ username }}</span>
        <div class="auth-buttons">
          <button class="logout-btn" @click="logout">
            <i class="fas fa-sign-out-alt"></i> 退出
          </button>
        </div>
      </div>
      <div class="auth-buttons" v-else>
        <router-link :to="{ name: 'login' }" class="auth-btn">
          <i class="fas fa-sign-in-alt"></i> 登录
        </router-link>
        <router-link :to="{ name: 'register' }" class="auth-btn">
          <i class="fas fa-user-plus"></i> 注册
        </router-link>
      </div>
    </header>
    
    <main>
      <div class="content-area">
        <router-view></router-view>
      </div>
      
      <div class="sidebar">
        <h2><i class="fas fa-user"></i> 用户信息</h2>
        <div v-if="isLoggedIn" class="card">
          <div class="card-title">{{ username }}</div>
          <div class="card-content">
            <p>用户ID: {{ userId }}</p>
            <p>已发表文章: {{ userArticlesCount }} 篇</p>
          </div>
        </div>
        <div v-else class="empty-state">
          <i class="fas fa-user"></i>
          <p>请登录查看用户信息</p>
        </div>
        
        <h2 style="margin-top: 30px;"><i class="fas fa-fire"></i> 热门文章</h2>
        <div v-if="popularArticles.length > 0">
          <div class="card" v-for="article in popularArticles" :key="article.article_id">
            <div class="card-title">{{ article.title }}</div>
            <div class="card-meta">
              <span>{{ article.commentCount }} 评论</span>
              <span>{{ formatDate(article.create_time) }}</span>
            </div>
          </div>
        </div>
        <div v-else class="empty-state">
          <i class="fas fa-file-alt"></i>
          <p>暂无热门文章</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { formatDate } from '@/utils'

export default {
  name: 'App',
  setup() {
    const userStore = useUserStore()
    const router = useRouter()
    
    const popularArticles = ref([])
    const userArticlesCount = ref(5)
    
    const isLoggedIn = computed(() => userStore.isLoggedIn)
    const username = computed(() => userStore.username)
    const userId = computed(() => userStore.userId)
    
    const logout = () => {
      userStore.logout()
    }
    
    const loadPopularArticles = () => {
      // 模拟热门文章数据
      popularArticles.value = [
        {
          article_id: 'pop-1',
          title: 'Vue 3 新特性解析',
          create_time: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
          commentCount: 15
        },
        {
          article_id: 'pop-2',
          title: 'Pinia 状态管理实践',
          create_time: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000).toISOString(),
          commentCount: 8
        },
        {
          article_id: 'pop-3',
          title: '前端工程化最佳实践',
          create_time: new Date().toISOString(),
          commentCount: 12
        }
      ]
    }
    
    onMounted(() => {
      loadPopularArticles()
    })
    
    return {
      isLoggedIn,
      username,
      userId,
      popularArticles,
      userArticlesCount,
      logout,
      formatDate
    }
  }
}
</script>

<style scoped>
/* 这里不再重复，使用公共的style.css */
</style>