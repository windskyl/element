import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user.js'

// 使用动态导入组件
const Home = () => import('@/views/Home.vue')
// About page removed
const Login = () => import('@/views/Login.vue')
const Register = () => import('@/views/Register.vue')
const CreateArticle = () => import('@/views/CreateArticle.vue')
const ArticleDetail = () => import('@/views/ArticleDetail.vue')
const EditArticle = () => import('@/views/EditArticle.vue')
const UserInfo = () => import('@/views/UserInfo.vue')

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home
  },
  {
    path: '/login',
    name: 'login',
    component: Login,
    meta: { guestOnly: true }
  },
  {
    path: '/register',
    name: 'register',
    component: Register,
    meta: { guestOnly: true }
  },
  {
    path: '/articles/create',
    name: 'create-article',
    component: CreateArticle,
    meta: { requiresAuth: true }
  },
  {
    path: '/articles/:id',
    name: 'article-detail',
    component: ArticleDetail,
    props: true
  },
  {
    path: '/articles/:id/edit',
    name: 'edit-article',
    component: EditArticle,
    props: true,
    meta: { requiresAuth: true }
  },
  {
    path: '/user',
    name: 'user-info',
    component: UserInfo,
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  // } else if (to.meta.guestOnly && userStore.isLoggedIn) {
  //   next({ name: 'home' })
  } else {
    next()
  }
})

export default router