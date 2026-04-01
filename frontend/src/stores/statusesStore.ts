import { defineStore } from 'pinia'
import { statusesApi } from '@/api/client'

export interface Status {
  id: number
  code: string
  name: string
  sort: number
  isFinal: boolean
}

interface StatusesState {
  items: Status[]
  loading: boolean
}

export const useStatusesStore = defineStore('statuses', {
  state: (): StatusesState => ({
    items: [],
    loading: false,
  }),

  actions: {
    async fetchAll(): Promise<void> {
      this.loading = true
      try {
        const resp = await statusesApi.statusesGet({})
        this.items = resp.items.sort((a, b) => a.sort - b.sort)
      } catch (e) {
        console.error('statuses fetch error:', e)
      } finally {
        this.loading = false
      }
    },
  },
})
