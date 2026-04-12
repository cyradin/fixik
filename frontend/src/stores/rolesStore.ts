import { defineStore } from 'pinia'
import { rolesApi } from '@/api/client'
import { notifyError } from '@/utils/notify'

export interface Role {
  code: string
  description: string
  name: string
}

interface RolesState {
  items: Role[]
  pollingId: ReturnType<typeof setInterval> | null
}

export const useRolesStore = defineStore('roles', {
  state: (): RolesState => ({
    items: [],
    pollingId: null,
  }),

  actions: {
    async fetchAll(): Promise<void> {
      try {
        const resp = await rolesApi.rolesGet({})
        this.items = resp.items
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
