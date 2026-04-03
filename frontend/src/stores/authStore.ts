// stores/authStore.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const isAuth = ref(false)
  const loading = ref(false)
  const error = ref('')
  const initialized = ref(false) // <- добавили сюда

  const login = async (username: string, password: string) => {
    loading.value = true
    try {
      await authApi.authLoginPost({ request: { username, password } })
      isAuth.value = true
      error.value = ''
    } catch (e: any) {
      error.value = e.message || 'Login failed'
      isAuth.value = false
      throw e
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    await authApi.authLogoutPost()
    isAuth.value = false
  }

  const refresh = async () => {
    try {
      await authApi.authRefreshPost()
      isAuth.value = true
    } catch (e) {
      isAuth.value = false
    } finally {
      initialized.value = true
    }
  }

  return { isAuth, loading, error, initialized, login, logout, refresh }
})
