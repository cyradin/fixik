import { defineStore } from 'pinia'
import { prioritiesApi } from '@/api/client'
import { notifyError } from '@/utils/notify'

export interface Priority {
  id: number
  name: string
  sort: number
}

interface PrioritiesState {
  items: Priority[]
  pollingId: ReturnType<typeof setInterval> | null
}

export const usePrioritiesStore = defineStore('priorities', {
  state: (): PrioritiesState => ({
    items: [],
    pollingId: null,
  }),

  actions: {
    async fetchAll(): Promise<void> {
      try {
        const resp = await prioritiesApi.prioritiesGet({})
        this.items = resp.items
      } catch (e) {
        notifyError('Не удалось обновить список приоритетов инцидентов')
        console.error('priorities fetch error:', e)
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
