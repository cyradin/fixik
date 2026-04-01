import { defineStore } from 'pinia'
import { usersApi } from '@/api/client'

export interface User {
  id: number
  name: string
  email: string
  username: string
}

interface UsersState {
  items: User[]
  loading: boolean
  pollingId: ReturnType<typeof setInterval> | null
}

export const useUsersStore = defineStore('users', {
  state: (): UsersState => ({
    items: [],
    loading: false,
    pollingId: null,
  }),

  actions: {
    async fetchAll(): Promise<void> {
      this.loading = true
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
        console.error('users fetch error:', e)
      } finally {
        this.loading = false
      }
    },

    startPolling(interval = 10000): void {
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
