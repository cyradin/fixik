// stores/authStore.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '@/api/client'

export interface CurrentUser {
  id: number
  username: string
  roles: string[]
}

export const useAuthStore = defineStore('auth', () => {
  const isAuth = ref(false)
  const loading = ref(false)
  const error = ref('')
  const initialized = ref(false)
  const user = ref<CurrentUser | null>(null) // <-- сюда сохраняем данные юзера

  const login = async (username: string, password: string) => {
    loading.value = true
    try {
      const data = await authApi.authLoginPost({ request: { username, password } })
      isAuth.value = true
      user.value = data.user // предполагаем, что сервер вернул { user: { id, username, roles } }
      error.value = ''
    } catch (e: any) {
      error.value = e.message || 'Login failed'
      isAuth.value = false
      user.value = null
      throw e
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    await authApi.authLogoutPost()
    isAuth.value = false
    user.value = null
  }

  const refresh = async () => {
    try {
      const data = await authApi.authRefreshPost()
      isAuth.value = true
      user.value = data.user // <-- тоже из ответа сервера
    } catch (e) {
      isAuth.value = false
      user.value = null
    } finally {
      initialized.value = true
    }
  }

  return { isAuth, loading, error, initialized, user, login, logout, refresh }
})
