import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { authApi, usersApi } from '@/api/client'
import { extractUserMessage } from '@/utils/errors'
import { useRolesStore } from './rolesStore'

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
  const initialized = ref(false)
  const user = ref<CurrentUser | null>(null)
  const rolesStore = useRolesStore()

  const login = async (username: string, password: string) => {
    loading.value = true
    try {
      const data = await authApi.authLoginPost({ request: { username, password } })
      isAuth.value = true
      user.value = data
    } catch (e: any) {
      isAuth.value = false
      user.value = null

      console.log('login error: ', e)

      throw new Error(await extractUserMessage(e))
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

      await authApi.authPasswordPost({
        request: {
          currentPassword,
          newPassword,
        },
      })
    } catch (e: any) {
      console.log('login error: ', e)

      throw new Error(await extractUserMessage(e))
    }
  }

  function can(permission: string): boolean
  function can(permission: readonly string[]): boolean
  function can(permission: string | readonly string[]): boolean {
    const role = rolesStore.byCode(user.value?.role)
    const userPermissions = new Set(role?.permissions?.map((p) => p.code) ?? [])
    const perms = Array.isArray(permission) ? permission : [permission]

    return perms.every((p) => userPermissions.has(p))
  }

  return {
    isAuth,
    loading,
    initialized,
    user,
    login,
    logout,
    refresh,
    changePassword,
    can,
  }
})
