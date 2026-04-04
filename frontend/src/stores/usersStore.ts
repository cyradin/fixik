import { defineStore } from 'pinia'
import { usersApi } from '@/api/client'
import { notifyError } from '@/utils/notify'

export interface User {
  id: number
  name: string
  email: string
  username: string
}

interface UsersState {
  items: User[]
  pollingId: ReturnType<typeof setInterval> | null
}

export const useUsersStore = defineStore('users', {
  state: (): UsersState => ({
    items: [],
    pollingId: null,
  }),

  getters: {
    byTeam: (state) => {
      return (teamId?: number) => {
        if (!teamId || teamId === 0) return state.items
        return state.items.filter((u: any) => u.teamId === teamId)
      }
    },
  },

  actions: {
    async fetchAll(): Promise<void> {
      try {
        let allItems: User[] = []
        let offset = 0
        const limit = 100

        while (true) {
          const resp = await usersApi.usersGet({ limit, offset })
          const items = resp.items ?? []
          allItems = allItems.concat(items)

          if (items.length < limit) break

          offset += limit
        }

        this.items = allItems
      } catch (e) {
        notifyError('Не удалось обновить список пользователей')
        console.error('users fetch error:', e)
      }
    },

    startPolling(interval = 5000): void {
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
