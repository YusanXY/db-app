import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, getCurrentUser } from '@/api/auth'
import { getToken, setToken, removeToken } from '@/utils/auth'
import type { User } from '@/api/types'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(getToken() || '')
  const userInfo = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 'admin')

  async function loginAction(username: string, password: string) {
    const result = await login({ username, password })
    token.value = result.token
    userInfo.value = result.user
    setToken(result.token)
  }

  async function fetchUserInfo() {
    if (!token.value) return
    userInfo.value = await getCurrentUser()
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    removeToken()
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    isAdmin,
    loginAction,
    fetchUserInfo,
    logout
  }
})

