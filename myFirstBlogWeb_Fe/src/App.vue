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
        <!-- 圆形突出写文章按钮 -->
        <router-link 
          v-if="isLoggedIn" 
          :to="{ name: 'create-article' }" 
          class="fab" 
          aria-label="写文章"
        >
          <span class="fab-plus">+</span>
        </router-link>
        <!-- 用户信息页面入口 -->
        <router-link v-if="isLoggedIn" :to="{ name: 'user-info' }" class="nav-btn" active-class="active">
          <i class="fas fa-user-circle"></i> 用户信息
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
    
    <main class="main-with-sidebar">
      <aside class="sidebar left">
        <!-- 左侧侧边栏占位，可用于后续扩展 -->
      </aside>
      <div class="content-area">
        <router-view></router-view>
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
    
  // const popularArticles = ref([]) // 已移除热门文章功能
  const userArticlesCount = computed(() => userStore.articlesCount)
    
    const isLoggedIn = computed(() => userStore.isLoggedIn)
    const username = computed(() => userStore.username)
    const userId = computed(() => userStore.userId)
    
    const logout = () => {
      userStore.logout()
    }
    
    // hot articles removed
    
    onMounted(async () => {
      if (userStore.isLoggedIn && userStore.isTokenValid && userStore.isTokenValid()) {
        await userStore.refreshArticlesCount()
      }
    })
    
    return {
      isLoggedIn,
      username,
      userId,
      userArticlesCount,
      logout,
      formatDate,
      // kept minimal
    }
  }
}
</script>

<style scoped>
/* Top FAB and left sidebar styles */
.main-with-sidebar{
  display:flex;
  gap:20px;
}
.sidebar.left{
  width:220px;
  flex: 0 0 220px;
}
.content-area{
  flex:1 1 auto;
}
.fab{
  display:inline-flex;
  align-items:center;
  justify-content:center;
  width:46px;
  height:46px;
  border-radius:50%;
  background: #2f8fef;
  color: #fff;
  box-shadow: 0 6px 14px rgba(47,143,239,0.25);
  text-decoration: none;
  margin: 0 8px;
}
.fab:hover{ transform: translateY(-2px); }
.fab-plus{ font-size: 22px; line-height:1; font-weight:700; }

/* tweak nav button spacing */
.nav-btn{ margin-right:10px }
</style>