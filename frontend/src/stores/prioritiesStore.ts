import { defineStore } from 'pinia'
import { prioritiesApi } from '@/api/client'

export interface Priority {
  id: number
  name: string
  sort: number
}

interface PrioritiesState {
  items: Priority[]
}

export const usePrioritiesStore = defineStore('priorities', {
  state: (): PrioritiesState => ({
    items: [],
  }),

  actions: {
    async fetchAll(): Promise<void> {
      try {
        const resp = await prioritiesApi.prioritiesGet({})
        this.items = resp.items
      } catch (e) {
        console.error('priorities fetch error:', e)
      }
    },
  },
})
