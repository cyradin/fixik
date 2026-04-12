import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi, usersApi } from '@/api/client'

export interface CurrentUser {
  id: number
  name: string
  username: string
  email: string
  role: string
}

export const useAuthStore = defineStore('auth', () => {
  const isAuth = ref(false)
  const loading = ref(false)
  const error = ref('')
  const initialized = ref(false)
  const user = ref<CurrentUser | null>(null)

  const login = async (username: string, password: string) => {
    loading.value = true
    try {
      const data = await authApi.authLoginPost({ request: { username, password } })
      isAuth.value = true
      user.value = data
      error.value = ''
    } catch (e: any) {
      if (e.response?.status === 401) {
        error.value = 'Неверный логин или пароль'
      } else {
        error.value = 'Что-то пошло не так. Попробуйте позже'
      }

      isAuth.value = false
      user.value = null
      throw e
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    try {
      await authApi.authLogoutPost()
      isAuth.value = false
      user.value = null
    } catch (e: any) {
      throw e
    }
  }

  const refresh = async () => {
    try {
      const data = await authApi.authRefreshPost()
      isAuth.value = true
      user.value = data
    } catch (e) {
      isAuth.value = false
      user.value = null
    } finally {
      initialized.value = true
    }
  }

  const changePassword = async (currentPassword: string, newPassword: string) => {
    try {
      const id = user.value?.id
      if (!id) throw new Error('User not initialized')

      await usersApi.usersIdPasswordPost({
        id,
        request: {
          currentPassword,
          newPassword,
        },
      })
    } catch (e: any) {
      if (e?.response?.status === 401) {
        throw new Error('Неверный текущий пароль')
      }

      throw e
    }
  }

  const clearError = () => {
    error.value = ''
  }

  return {
    isAuth,
    loading,
    error,
    initialized,
    user,
    login,
    logout,
    refresh,
    clearError,
    changePassword,
  }
})
