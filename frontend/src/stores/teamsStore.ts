import { defineStore } from 'pinia'
import { teamsApi } from '@/api/client'

export interface Team {
  id: number
  code: string
  name: string
}

interface TeamsState {
  items: Team[]
  pollingId: ReturnType<typeof setInterval> | null
}

export const useTeamsStore = defineStore('teams', {
  state: (): TeamsState => ({
    items: [],
    pollingId: null,
  }),

  actions: {
    async fetchAll(): Promise<void> {
      try {
        const resp = await teamsApi.teamsGet({})
        this.items = resp.items
      } catch (e) {
        console.error('teams fetch error:', e)
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
