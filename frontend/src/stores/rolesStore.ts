import { defineStore } from 'pinia'
import { rolesApi } from '@/api/client'
import { notifyError } from '@/utils/notify'
import type { PermissionCode } from '@/constants/permissions'

export interface Role {
  code: string
  description: string
  name: string
  permissions: Permission[]
}

export interface Permission {
  code: PermissionCode
}
interface RolesState {
  initialized: boolean
  items: Role[]
  pollingId: ReturnType<typeof setInterval> | null
}

export const useRolesStore = defineStore('roles', {
  state: (): RolesState => ({
    items: [],
    pollingId: null,
    initialized: false,
  }),

  getters: {
    byCode: (state) => (code?: string) =>
      code ? state.items.find((r) => r.code === code) || null : null,

    nameByCode: (state) => (code?: string) => {
      const role = state.items.find((r) => r.code === code)
      return role?.name || code
    },
  },

  actions: {
    async fetchAll(): Promise<void> {
      try {
        const resp = await rolesApi.rolesGet({})
        this.items = resp.items.map((role) => ({
          ...role,
          permissions: role.permissions as { code: PermissionCode }[],
        }))
        this.initialized = true
      } catch (e) {
        notifyError('Не удалось обновить список ролей пользователей')
        console.error('roles fetch error:', e)
      }
    },

    startPolling(interval = 30000): void {
      this.stopPolling()
      this.fetchAll()
      this.pollingId = setInterval(() => this.fetchAll(), interval)
    },

    stopPolling(): void {
      if (this.pollingId) {
        clearInterval(this.pollingId)
        this.pollingId = null
      }
    },
  },
})
