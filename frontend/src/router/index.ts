import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue')
  },
  {
    path: '/article/:id',
    name: 'ArticleDetail',
    component: () => import('@/views/Article/Detail.vue')
  },
  {
    path: '/article/new',
    name: 'ArticleCreate',
    component: () => import('@/views/Article/Create.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Auth/Login.vue'),
    meta: { guestOnly: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Auth/Register.vue'),
    meta: { guestOnly: true }
  },
  {
    path: '/categories',
    name: 'CategoryList',
    component: () => import('@/views/Category/List.vue')
  },
  {
    path: '/tags',
    name: 'TagList',
    component: () => import('@/views/Tag/List.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  // 需要登录
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  // 仅游客访问
  if (to.meta.guestOnly && userStore.isLoggedIn) {
    next({ name: 'Home' })
    return
  }

  next()
})

export default router

