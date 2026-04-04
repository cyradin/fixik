import { defineStore } from 'pinia'
import { statusesApi } from '@/api/client'
import { notifyError } from '@/utils/notify'

export interface Status {
  id: number
  code: string
  name: string
  sort: number
  isFinal: boolean
}

interface StatusesState {
  items: Status[]
  pollingId: ReturnType<typeof setInterval> | null
}

export const useStatusesStore = defineStore('statuses', {
  state: (): StatusesState => ({
    items: [],
    pollingId: null,
  }),

  actions: {
    async fetchAll(): Promise<void> {
      try {
        const resp = await statusesApi.statusesGet({})
        this.items = resp.items.sort((a, b) => a.sort - b.sort)
      } catch (e) {
        notifyError('Не удалось обновить список статусов инцидентов')
        console.error('statuses fetch error:', e)
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
